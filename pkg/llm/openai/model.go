package openai

// Model represents an OpenAI model
// https://platform.openai.com/docs/api-reference/models
// https://platform.openai.com/docs/api-reference/models/retrieve

type Model struct {
	ID         string            `json:"id"`
	Object     string            `json:"object"`
	OwnedBy    string            `json:"owned_by"`
	Permission []ModelPermission `json:"permission"`
	Root       string            `json:"root,omitempty"`
	Parent     string            `json:"parent,omitempty"`
}

type ModelPermission struct {
	ID                 string `json:"id"`
	Object             string `json:"object"`
	Created            int64  `json:"created"`
	AllowCreateEngine  bool   `json:"allow_create_engine"`
	AllowSampling      bool   `json:"allow_sampling"`
	AllowLogprobs      bool   `json:"allow_logprobs"`
	AllowSearchIndices bool   `json:"allow_search_indices"`
	AllowView          bool   `json:"allow_view"`
	AllowFineTuning    bool   `json:"allow_fine_tuning"`
	Organization       string `json:"organization"`
	Group              string `json:"group,omitempty"`
	IsBlocking         bool   `json:"is_blocking"`
}

type ListModelsResponse struct {
	Data   []Model `json:"data"`
	Object string  `json:"object"`
}

type RetrieveModelResponse = Model

// Usage represents token usage

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Completion endpoint structs
// https://platform.openai.com/docs/api-reference/completions

// CompletionRequest is used to create a text completion

type CompletionRequest struct {
	Model            string         `json:"model"`
	Prompt           interface{}    `json:"prompt,omitempty"` // string or []string
	Suffix           string         `json:"suffix,omitempty"`
	MaxTokens        int            `json:"max_tokens,omitempty"`
	Temperature      float64        `json:"temperature,omitempty"`
	TopP             float64        `json:"top_p,omitempty"`
	N                int            `json:"n,omitempty"`
	Stream           bool           `json:"stream,omitempty"`
	Logprobs         int            `json:"logprobs,omitempty"`
	Echo             bool           `json:"echo,omitempty"`
	Stop             interface{}    `json:"stop,omitempty"` // string or []string
	PresencePenalty  float64        `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64        `json:"frequency_penalty,omitempty"`
	BestOf           int            `json:"best_of,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	User             string         `json:"user,omitempty"`
}

type CompletionResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
	Usage   Usage              `json:"usage,omitempty"`
}

type CompletionChoice struct {
	Text         string    `json:"text"`
	Index        int       `json:"index"`
	Logprobs     *Logprobs `json:"logprobs,omitempty"`
	FinishReason string    `json:"finish_reason,omitempty"`
}

type Logprobs struct {
	Tokens        []string             `json:"tokens"`
	TokenLogprobs []float64            `json:"token_logprobs"`
	TopLogprobs   []map[string]float64 `json:"top_logprobs"`
	TextOffset    []int                `json:"text_offset"`
}

// Chat endpoint structs
// https://platform.openai.com/docs/api-reference/chat

type ChatCompletionRequest struct {
	Model            string         `json:"model"`
	Messages         []ChatMessage  `json:"messages"`
	Temperature      float64        `json:"temperature,omitempty"`
	TopP             float64        `json:"top_p,omitempty"`
	N                int            `json:"n,omitempty"`
	Stream           bool           `json:"stream,omitempty"`
	Stop             interface{}    `json:"stop,omitempty"`
	MaxTokens        int            `json:"max_tokens,omitempty"`
	PresencePenalty  float64        `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64        `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	User             string         `json:"user,omitempty"`
}

type ChatMessage struct {
	Role         string            `json:"role"`
	Content      string            `json:"content"`
	Name         string            `json:"name,omitempty"`
	FunctionCall *ChatFunctionCall `json:"function_call,omitempty"`
}

type ChatFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ChatCompletionResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
	Usage   Usage        `json:"usage,omitempty"`
}

type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason,omitempty"`
}

// Edit endpoint structs
// https://platform.openai.com/docs/api-reference/edits

type EditRequest struct {
	Model       string  `json:"model"`
	Input       string  `json:"input,omitempty"`
	Instruction string  `json:"instruction"`
	N           int     `json:"n,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
}

type EditResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Choices []EditChoice `json:"choices"`
	Usage   Usage        `json:"usage,omitempty"`
}

type EditChoice struct {
	Text  string `json:"text"`
	Index int    `json:"index"`
}

// Embedding endpoint structs
// https://platform.openai.com/docs/api-reference/embeddings

type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
	User  string   `json:"user,omitempty"`
}

type EmbeddingResponse struct {
	Object string          `json:"object"`
	Data   []EmbeddingData `json:"data"`
	Model  string          `json:"model"`
	Usage  Usage           `json:"usage,omitempty"`
}

type EmbeddingData struct {
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

// Moderation endpoint structs
// https://platform.openai.com/docs/api-reference/moderations

type ModerationRequest struct {
	Input interface{} `json:"input"` // string or []string
	Model string      `json:"model,omitempty"`
}

type ModerationResponse struct {
	ID      string             `json:"id"`
	Model   string             `json:"model"`
	Results []ModerationResult `json:"results"`
}

type ModerationResult struct {
	Categories     map[string]bool    `json:"categories"`
	CategoryScores map[string]float64 `json:"category_scores"`
	Flagged        bool               `json:"flagged"`
}

// Image endpoint structs
// https://platform.openai.com/docs/api-reference/images

type ImageCreateRequest struct {
	Prompt         string `json:"prompt"`
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}

type ImageData struct {
	URL     string `json:"url,omitempty"`
	B64JSON string `json:"b64_json,omitempty"`
}

type ImageCreateResponse struct {
	Created int64       `json:"created"`
	Data    []ImageData `json:"data"`
}

type ImageEditRequest struct {
	Image          string `json:"image"`
	Mask           string `json:"mask,omitempty"`
	Prompt         string `json:"prompt"`
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}

type ImageEditResponse struct {
	Created int64       `json:"created"`
	Data    []ImageData `json:"data"`
}

type ImageVariationRequest struct {
	Image          string `json:"image"`
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}

type ImageVariationResponse struct {
	Created int64       `json:"created"`
	Data    []ImageData `json:"data"`
}

// File endpoint structs
// https://platform.openai.com/docs/api-reference/files

type File struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int    `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
}

type FileListResponse struct {
	Object string `json:"object"`
	Data   []File `json:"data"`
}

// Fine-tune endpoint structs
// https://platform.openai.com/docs/api-reference/fine-tunes

type FineTuneCreateRequest struct {
	TrainingFile                 string    `json:"training_file"`
	ValidationFile               string    `json:"validation_file,omitempty"`
	Model                        string    `json:"model,omitempty"`
	NEpochs                      int       `json:"n_epochs,omitempty"`
	BatchSize                    int       `json:"batch_size,omitempty"`
	LearningRateMultiplier       float64   `json:"learning_rate_multiplier,omitempty"`
	PromptLossWeight             float64   `json:"prompt_loss_weight,omitempty"`
	ComputeClassificationMetrics bool      `json:"compute_classification_metrics,omitempty"`
	ClassificationNClasses       int       `json:"classification_n_classes,omitempty"`
	ClassificationPositiveClass  string    `json:"classification_positive_class,omitempty"`
	ClassificationBetas          []float64 `json:"classification_betas,omitempty"`
	Suffix                       string    `json:"suffix,omitempty"`
}

type FineTune struct {
	ID              string              `json:"id"`
	Object          string              `json:"object"`
	Model           string              `json:"model"`
	CreatedAt       int64               `json:"created_at"`
	FineTunedModel  string              `json:"fine_tuned_model"`
	OrganizationID  string              `json:"organization_id"`
	Status          string              `json:"status"`
	Hyperparams     FineTuneHyperparams `json:"hyperparams"`
	TrainingFiles   []File              `json:"training_files"`
	ValidationFiles []File              `json:"validation_files"`
	ResultFiles     []File              `json:"result_files"`
	Events          []FineTuneEvent     `json:"events"`
}

type FineTuneHyperparams struct {
	NEpochs                int     `json:"n_epochs"`
	BatchSize              int     `json:"batch_size"`
	LearningRateMultiplier float64 `json:"learning_rate_multiplier"`
}

type FineTuneEvent struct {
	Object    string `json:"object"`
	CreatedAt int64  `json:"created_at"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

type FineTuneListResponse struct {
	Object string     `json:"object"`
	Data   []FineTune `json:"data"`
}
