package runware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	
	"github.com/google/uuid"
)

type NewImageInferenceReq struct {
	// Core task parameters
	TaskType        string `json:"taskType"`
	TaskUUID        string `json:"taskUUID"`
	DeliveryMethod  string `json:"deliveryMethod,omitempty"`
	WebhookURL      string `json:"webhookURL,omitempty"`
	UploadEndpoint  string `json:"uploadEndpoint,omitempty"`

	// Output configuration
	OutputType    string `json:"outputType,omitempty"`
	OutputFormat  string `json:"outputFormat,omitempty"`
	OutputQuality int    `json:"outputQuality,omitempty"`

	// Content and safety
	CheckNSFW   bool `json:"checkNSFW,omitempty"`
	IncludeCost bool `json:"includeCost,omitempty"`

	// Core generation parameters
	PositivePrompt string `json:"positivePrompt"`
	NegativePrompt string `json:"negativePrompt,omitempty"`
	Model          string `json:"model"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`

	// Image inputs for workflows
	SeedImage        string   `json:"seedImage,omitempty"`
	MaskImage        string   `json:"maskImage,omitempty"`
	MaskMargin       int      `json:"maskMargin,omitempty"`
	ReferenceImages  []string `json:"referenceImages,omitempty"`
	Strength         float64  `json:"strength,omitempty"`

	// Generation control parameters
	Steps              int     `json:"steps,omitempty"`
	Scheduler          string  `json:"scheduler,omitempty"`
	Seed               *int64  `json:"seed,omitempty"`
	CFGScale           float64 `json:"CFGScale,omitempty"`
	ClipSkip           *int    `json:"clipSkip,omitempty"`
	PromptWeighting    string  `json:"promptWeighting,omitempty"`
	NumberResults      int     `json:"numberResults,omitempty"`
	VAE                string  `json:"vae,omitempty"`

	// Outpainting
	Outpaint *Outpaint `json:"outpaint,omitempty"`

	// Advanced features and acceleration
	AdvancedFeatures     *AdvancedFeatures     `json:"advancedFeatures,omitempty"`
	AcceleratorOptions   *AcceleratorOptions   `json:"acceleratorOptions,omitempty"`

	// Identity and character consistency
	PuLID       *PuLID       `json:"puLID,omitempty"`
	ACEPlusPlus *ACEPlusPlus `json:"acePlusPlus,omitempty"`

	// Quality enhancement
	Refiner *Refiner `json:"refiner,omitempty"`

	// Style and control arrays
	Embeddings  []Embedding  `json:"embeddings,omitempty"`
	ControlNet  []ControlNet `json:"controlNet,omitempty"`
	Lora        []Lora       `json:"lora,omitempty"`
	IPAdapters  []IPAdapter  `json:"ipAdapters,omitempty"`

	// Provider-specific settings
	ProviderSettings *ProviderSettings `json:"providerSettings,omitempty"`
}

type NewImageInferenceResp struct {
	TaskType        string  `json:"taskType"`
	TaskUUID        string  `json:"taskUUID"`
	ImageUUID       string  `json:"imageUUID"`
	ImageURL        string  `json:"imageURL,omitempty"`
	ImageBase64Data string  `json:"imageBase64Data,omitempty"`
	ImageDataURI    string  `json:"imageDataURI,omitempty"`
	Seed            int64   `json:"seed,omitempty"`
	NSFWContent     bool    `json:"NSFWContent,omitempty"`
	Cost            float64 `json:"cost,omitempty"`
	TimedOut        bool    `json:"timedOut"`
}

func (sdk *SDK) ImageInference(ctx context.Context, req NewImageInferenceReq) (*NewImageInferenceResp, error) {
	req = *mergeImageInferenceReqWithDefaults(&req)
	if err := validateImageInferenceReq(req); err != nil {
		return nil, err
	}
	
	sendReq := Request{
		ID:            uuid.New().String(),
		Event:         NewTask,
		ResponseEvent: NewImage,
		Data:          []NewImageInferenceReq{req},
	}
	
	newImageInferenceResp := &NewImageInferenceResp{}
	
	responseChan := make(chan *NewImageInferenceResp)
	errChan := make(chan error)
	
	go func() {
		defer close(responseChan)
		defer close(errChan)
		
		for msg := range sdk.Client.Listen() {
			var msgData map[string]interface{}
			if err := json.Unmarshal(msg, &msgData); err != nil {
				errChan <- fmt.Errorf("%w:[%s]", ErrDecodeMessage, err.Error())
				return
			}
			
			// Check if is an error message first
			if errMsg, ok := sdk.OnError(msgData); ok {
				errChan <- errMsg
				return
			}
			
			for k, v := range msgData {
				if k != sendReq.ResponseEvent {
					log.Println("Skipping event", k, "Currently handling", sendReq.ResponseEvent)
					continue
				}
				
				// Handle array response
				vArray, ok := v.([]interface{})
				if !ok || len(vArray) == 0 {
					continue
				}
				
				bValue, err := interfaceToByte(vArray[0])
				if err != nil {
					errChan <- err
					return
				}
				
				err = json.Unmarshal(bValue, &newImageInferenceResp)
				if err != nil {
					errChan <- err
					return
				}
				
				responseChan <- newImageInferenceResp
				return
			}
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
		newImageInferenceResp.TimedOut = true
		return newImageInferenceResp, fmt.Errorf("%w:[%s]", ErrRequestTimeout, sendReq.Event)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func NewImageInferenceReqDefaults() *NewImageInferenceReq {
	return &NewImageInferenceReq{
		TaskType:       ImageInference,
		TaskUUID:       uuid.New().String(),
		DeliveryMethod: DeliveryMethodSync,
		OutputType:     OutputTypeURL,
		OutputFormat:   OutputFormatJPG,
		OutputQuality:  95,
		Steps:          20,
		CFGScale:       7,
		NumberResults:  1,
		Strength:       0.8,
	}
}

func mergeImageInferenceReqWithDefaults(req *NewImageInferenceReq) *NewImageInferenceReq {
	_ = MergeEventRequestsWithDefaults[*NewImageInferenceReq](req, NewImageInferenceReqDefaults())
	return req
}

func validateImageInferenceReq(req NewImageInferenceReq) error {
	if req.PositivePrompt == "" {
		return fmt.Errorf("%w:[%s]", ErrFieldRequired, "positivePrompt")
	}
	
	if req.Model == "" {
		return fmt.Errorf("%w:[%s]", ErrFieldRequired, "model")
	}
	
	if req.Width < 128 || req.Width > 2048 || req.Width%64 != 0 {
		return fmt.Errorf("%w:[%s][128-2048, divisible by 64]", ErrFieldIncorrectVal, "width")
	}
	
	if req.Height < 128 || req.Height > 2048 || req.Height%64 != 0 {
		return fmt.Errorf("%w:[%s][128-2048, divisible by 64]", ErrFieldIncorrectVal, "height")
	}
	
	if req.Steps < 1 || req.Steps > 100 {
		return fmt.Errorf("%w:[%s][1-100]", ErrFieldIncorrectVal, "steps")
	}
	
	if req.CFGScale < 0 || req.CFGScale > 50 {
		return fmt.Errorf("%w:[%s][0-50]", ErrFieldIncorrectVal, "CFGScale")
	}
	
	if req.ClipSkip != nil && (*req.ClipSkip < 0 || *req.ClipSkip > 2) {
		return fmt.Errorf("%w:[%s][0-2]", ErrFieldIncorrectVal, "clipSkip")
	}
	
	if req.OutputQuality < 20 || req.OutputQuality > 99 {
		return fmt.Errorf("%w:[%s][20-99]", ErrFieldIncorrectVal, "outputQuality")
	}
	
	if req.NumberResults < 1 || req.NumberResults > 20 {
		return fmt.Errorf("%w:[%s][1-20]", ErrFieldIncorrectVal, "numberResults")
	}
	
	if req.Strength < 0 || req.Strength > 1 {
		return fmt.Errorf("%w:[%s][0-1]", ErrFieldIncorrectVal, "strength")
	}
	
	// Validate workflow-specific requirements
	if req.SeedImage != "" && req.Strength == 0 {
		req.Strength = 0.8
	}
	
	if req.MaskImage != "" && req.SeedImage == "" {
		return fmt.Errorf("%w:[%s when maskImage is provided]", ErrFieldRequired, "seedImage")
	}
	
	if req.MaskMargin != 0 && (req.MaskMargin < 32 || req.MaskMargin > 128) {
		return fmt.Errorf("%w:[%s][32-128]", ErrFieldIncorrectVal, "maskMargin")
	}
	
	// Validate outpaint dimensions
	if req.Outpaint != nil {
		if req.SeedImage == "" {
			return fmt.Errorf("%w:[%s when outpaint is provided]", ErrFieldRequired, "seedImage")
		}
		
		if req.Outpaint.Top%64 != 0 || req.Outpaint.Right%64 != 0 || 
		   req.Outpaint.Bottom%64 != 0 || req.Outpaint.Left%64 != 0 {
			return fmt.Errorf("%w:[%s][all values must be divisible by 64]", ErrFieldIncorrectVal, "outpaint")
		}
		
		if req.Outpaint.Blur < 0 || req.Outpaint.Blur > 32 {
			return fmt.Errorf("%w:[%s][0-32]", ErrFieldIncorrectVal, "outpaint.blur")
		}
	}
	
	return nil
}
