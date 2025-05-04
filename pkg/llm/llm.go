package llm

import "context"

type GenerateRequest struct {
	Prompt  Prompt         `json:"prompt"`
	History []Conversation `json:"conversation"`
}

type GenerateStreamResponse struct {
	ID       string `json:"id"`
	Response string `json:"response"`
	Error    error  // if no error this should be null/nil
}

// GenerateResponse
// aggregated response, if someone wants to get the result without streaming
type GenerateResponse struct {
	Response string `json:"response"`
}

type Streamer interface {
	Stream(ctx context.Context, request GenerateRequest) (<-chan GenerateStreamResponse, error)
}
