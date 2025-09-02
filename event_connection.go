package runware

import (
	"context"
	"fmt"
	"log"
	"time"
	
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

type NewConnectReq struct {
	APIKey                string `json:"apiKey"`
	ConnectionSessionUUID string `json:"connectionSessionUUID,omitempty"`
	TaskType              string `json:"taskType"`
}

type NewConnectResp struct {
	ConnectionSessionUUID string `json:"connectionSessionUUID"`
	TimedOut              bool   `json:"timedOut"`
}

func (sdk *SDK) Connect(ctx context.Context, req NewConnectReq) (*NewConnectResp, error) {
	req = *mergeNewConnectReqWithDefaults(&req)
	if err := validateNewConnectReq(req); err != nil {
		return nil, err
	}
	
	sendReq := Request{
		ID:            uuid.New().String(),
		Event:         NewConnection,
		ResponseEvent: NewConnectionSessionUUID,
		Data:          req,
	}
	
	newConnectResp := &NewConnectResp{}
	
	responseChan := make(chan *NewConnectResp)
	errChan := make(chan error)
	
	go func() {
		defer close(responseChan)
		defer close(errChan)
		
		for msg := range sdk.Client.Listen() {
			msgStr := string(msg)
			
			// Check if is an error message first
			if gjson.Get(msgStr, "error").Bool() {
				errorId := gjson.Get(msgStr, "errorId").Float()
				errorMessage := gjson.Get(msgStr, "errorMessage").String()
				
				var err error
				switch errorId {
				case 19:
					err = ErrInvalidApiKey
				default:
					err = ErrWsUnknownError
				}
				
				errChan <- fmt.Errorf("%w:[%v:%s]", err, errorId, errorMessage)
				return
			}
			
			// Check if message contains our response event
			if gjson.Get(msgStr, sendReq.ResponseEvent).Exists() {
				sessionUUID := gjson.Get(msgStr, sendReq.ResponseEvent).String()
				newConnectResp.ConnectionSessionUUID = sessionUUID
				responseChan <- newConnectResp
				return
			}
			
			// Skip if not our event
			log.Println("Skipping message, waiting for", sendReq.ResponseEvent)
		}
	}()
	
	bSendReq, err := sendReq.ToEvent()
	if err != nil {
		return nil, err
	}
	
	if err = sdk.Client.Send(bSendReq); err != nil {
		return nil, err
	}
	
	select {
	case resp := <-responseChan:
		return resp, nil
	case err = <-errChan:
		return nil, err
	case <-time.After(timeoutSendResponse * time.Second):
		newConnectResp.TimedOut = true
		return newConnectResp, fmt.Errorf("%w:[%s]", ErrRequestTimeout, sendReq.Event)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func NewConnectReqDefaults() *NewConnectReq {
	return &NewConnectReq{
		TaskType: "ping",
	}
}

func mergeNewConnectReqWithDefaults(req *NewConnectReq) *NewConnectReq {
	_ = MergeEventRequestsWithDefaults[*NewConnectReq](req, NewConnectReqDefaults())
	return req
}

func validateNewConnectReq(req NewConnectReq) error {
	if req.APIKey == "" {
		return fmt.Errorf("%w:[%s]", ErrFieldRequired, "apiKey")
	}
	return nil
}
