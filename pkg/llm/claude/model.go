package claude

type apiRequest struct {
	Model             string  `json:"model"`
	Prompt            string  `json:"prompt"`
	MaxTokensToSample int     `json:"max_tokens_to_sample"`
	Temperature       float64 `json:"temperature"`
	Stream            bool    `json:"stream"`
}
