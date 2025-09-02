package runware

// Task types
const (
	ImageInference            = "imageInference"
	TextToImage               = "textToImage"
	ImageToImage              = "imageCaption"
	Inpainting                = "inpainting"
	ImageToText               = "imageToText"
	PromptEnhancer            = "promptEnhancer"
	ImageUpscale              = "imageUpscale"
	ImageUpload               = "imageUpload"
	RemoveBackground          = "imageBackgroundRemoval"
	ControlNetTextToImage     = "controlNetTextToImage"
	ControlNetImageToImage    = "controlNetImageToImage"
	ControlNetPreprocessImage = "controlNetPreprocessImage"
)

// Output types
const (
	OutputTypeURL        = "URL"
	OutputTypeBase64Data = "base64Data"
	OutputTypeDataURI    = "dataURI"
)

// Output formats
const (
	OutputFormatJPG  = "JPG"
	OutputFormatPNG  = "PNG"
	OutputFormatWEBP = "WEBP"
)

// Delivery methods
const (
	DeliveryMethodSync  = "sync"
	DeliveryMethodAsync = "async"
)

// Prompt weighting syntax
const (
	PromptWeightingCompel   = "compel"
	PromptWeightingSDEmbeds = "sdEmbeds"
)

// Available models (legacy constants for backward compatibility)
const (
	ModelSDXL               = 4
	ModelRevAnimated        = 13
	ModelAbsolutereality    = 18
	ModelCyberrealistic     = 19
	ModelDreamshaper        = 20
	ModelGhostmixBakedvae   = 22
	ModelSamaritan3DCartoon = 25
)

// Available processors
const (
	ProcessorCanny        = "canny"
	ProcessorDepth        = "depth"
	ProcessorMlsd         = "mlsd"
	ProcessorNormalbae    = "normalbae"
	ProcessorOpenpose     = "openpose"
	ProcessorTile         = "tile"
	ProcessorSeg          = "seg"
	ProcessorLineart      = "lineart"
	ProcessorLineartAnime = "lineart_anime"
	ProcessorShuffle      = "shuffle"
	ProcessorScribble     = "scribble"
	ProcessorSoftedge     = "softedge"
)

// Available sizes (legacy constants for backward compatibility)
const (
	SizeSquare512          = 1
	SizePortrait2to3       = 2
	SizePortrait1to2       = 3
	SizeLandscape2to3      = 4
	SizeLandscape2to1      = 5
	SizeLandscape4to3      = 6
	SizeLandscape16to9     = 7
	SizePortrait9to16      = 8
	SizePortrait3to4       = 9
	SizeSquare1024SDXL     = 11
	SizeLandscape16to9SDXL = 16
	SizePortrait9to16SDXL  = 17
	SizePortrait2to3SDXL   = 20
	SizeLandscape3to2SDXL  = 21
)

// ACE++ types
const (
	ACETypePlus         = "portrait"
	ACETypeSubject      = "subject"
	ACETypeLocalEditing = "local_editing"
)

// ControlNet represents a ControlNet configuration for guided image generation
type ControlNet struct {
	Model               string  `json:"model"`
	GuideImage          string  `json:"guideImage"`
	Weight              float64 `json:"weight,omitempty"`
	StartStep           *int    `json:"startStep,omitempty"`
	StartStepPercentage *int    `json:"startStepPercentage,omitempty"`
	EndStep             *int    `json:"endStep,omitempty"`
	EndStepPercentage   *int    `json:"endStepPercentage,omitempty"`
	ControlMode         string  `json:"controlMode,omitempty"`
}

// Lora represents a LoRA (Low-Rank Adaptation) configuration
type Lora struct {
	Model  string  `json:"model"`
	Weight float64 `json:"weight,omitempty"`
}

// Refiner represents SDXL refiner configuration for two-stage generation
type Refiner struct {
	Model               string `json:"model"`
	StartStep           *int   `json:"startStep,omitempty"`
	StartStepPercentage *int   `json:"startStepPercentage,omitempty"`
}

// Embedding represents an embedding (Textual Inversion) configuration
type Embedding struct {
	Model  string  `json:"model"`
	Weight float64 `json:"weight,omitempty"`
}

// IPAdapter represents an IP-Adapter configuration for image-prompted generation
type IPAdapter struct {
	Model      string  `json:"model"`
	GuideImage string  `json:"guideImage"`
	Weight     float64 `json:"weight,omitempty"`
}

// Outpaint represents outpainting configuration for image extension
type Outpaint struct {
	Top    int `json:"top,omitempty"`
	Right  int `json:"right,omitempty"`
	Bottom int `json:"bottom,omitempty"`
	Left   int `json:"left,omitempty"`
	Blur   int `json:"blur,omitempty"`
}

// AdvancedFeatures contains specialized features for image generation
type AdvancedFeatures struct {
	LayerDiffuse bool `json:"layerDiffuse,omitempty"`
}

// AcceleratorOptions contains caching mechanisms for faster generation
type AcceleratorOptions struct {
	TeaCache           bool    `json:"teaCache,omitempty"`
	TeaCacheDistance   float64 `json:"teaCacheDistance,omitempty"`
	DeepCache          bool    `json:"deepCache,omitempty"`
	DeepCacheInterval  int     `json:"deepCacheInterval,omitempty"`
	DeepCacheBranchId  int     `json:"deepCacheBranchId,omitempty"`
}

// PuLID represents PuLID identity customization configuration
type PuLID struct {
	InputImages             []string `json:"inputImages"`
	IdWeight                int      `json:"idWeight,omitempty"`
	TrueCFGScale            float64  `json:"trueCFGScale,omitempty"`
	CFGStartStep            int      `json:"CFGStartStep,omitempty"`
	CFGStartStepPercentage  int      `json:"CFGStartStepPercentage,omitempty"`
}

// ACEPlusPlus represents ACE++ character-consistent generation configuration
type ACEPlusPlus struct {
	Type            string   `json:"type,omitempty"`
	InputImages     []string `json:"inputImages,omitempty"`
	InputMasks      []string `json:"inputMasks,omitempty"`
	RepaintingScale float64  `json:"repaintingScale,omitempty"`
}

// BFLSettings represents Black Forest Labs provider-specific settings
type BFLSettings struct {
	PromptUpsampling bool `json:"promptUpsampling,omitempty"`
	SafetyTolerance  int  `json:"safetyTolerance,omitempty"`
	Raw              bool `json:"raw,omitempty"`
}

// ProviderSettings contains provider-specific configurations
type ProviderSettings struct {
	BFL *BFLSettings `json:"bfl,omitempty"`
}

// ImageInferenceRequest represents a comprehensive image generation request
type ImageInferenceRequest struct {
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

// ImageInferenceResponse represents the response from image generation
type ImageInferenceResponse struct {
	TaskType        string  `json:"taskType"`
	TaskUUID        string  `json:"taskUUID"`
	ImageUUID       string  `json:"imageUUID"`
	ImageURL        string  `json:"imageURL,omitempty"`
	ImageBase64Data string  `json:"imageBase64Data,omitempty"`
	ImageDataURI    string  `json:"imageDataURI,omitempty"`
	Seed            int64   `json:"seed,omitempty"`
	NSFWContent     bool    `json:"NSFWContent,omitempty"`
	Cost            float64 `json:"cost,omitempty"`
}

// Legacy types for backward compatibility

// Image represents a legacy image response
type Image struct {
	ImageSrc     string `json:"imageSrc"`
	ImageUUID    string `json:"imageUUID"`
	BNSFWContent bool   `json:"bNSFWContent"`
	ImageAltText string `json:"imageAltText"`
	TaskUUID     string `json:"taskUUID"`
}

// Text represents a text response
type Text struct {
	TaskUUID string `json:"taskUUID"`
	Text     string `json:"text"`
}

// PreProcessControlNet represents ControlNet preprocessing parameters
type PreProcessControlNet struct {
	TaskUUID           string `json:"taskUUID"`
	PreProcessorType   string `json:"preProcessorType"`
	GuideImageUUID     string `json:"guideImageUUID"`
	TaskType           string `json:"taskType"`
	Width              int    `json:"width"`
	Height             int    `json:"height"`
	LowThresholdCanny  int    `json:"lowThresholdCanny"`
	HighThresholdCanny int    `json:"highThresholdCanny"`
}

// Task represents a legacy task structure (deprecated, use ImageInferenceRequest instead)
type Task struct {
	TaskUUID           string       `json:"taskUUID"`
	ImageInitiatorUUID string       `json:"imageInitiatorUUID,omitempty"`
	ImageMaskUUID      string       `json:"imageMaskUUID,omitempty"`
	PositivePrompt     string       `json:"positivePrompt"`
	NegativePrompt     string       `json:"negativePrompt,omitempty"`
	Width              int          `json:"width"`
	Height             int          `json:"height"`
	Steps              int          `json:"steps,omitempty"`
	CFGScale           float64      `json:"CFGScale,omitempty"`
	Scheduler          string       `json:"scheduler,omitempty"`
	Seed               *int64       `json:"seed,omitempty"`
	Model              string       `json:"model"`
	Vae                string       `json:"vae,omitempty"`
	ClipSkip           *int         `json:"clipSkip,omitempty"`
	NumberResults      int          `json:"numberResults,omitempty"`
	TaskType           string       `json:"taskType"`
	PromptLanguageId   *string      `json:"promptLanguageId,omitempty"`
	Offset             int          `json:"offset,omitempty"`
	Refiner            *Refiner     `json:"refiner,omitempty"`
	Lora               []Lora       `json:"lora,omitempty"`
	ControlNet         []ControlNet `json:"controlNet,omitempty"`
	Embeddings         []Embedding  `json:"embeddings,omitempty"`
}
