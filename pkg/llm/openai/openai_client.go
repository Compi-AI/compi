package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/compiai/engine/pkg/llm"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type Config struct {
	ApiKey   string `yaml:"apiKey"`
	Endpoint string `yaml:"endpoint"`
	Model    string `yaml:"model"`
}

// Client implements the llm.Streamer interface using OpenAI's Chat Completions API.
type Client struct {
	logger     *slog.Logger
	config     Config
	httpClient *http.Client
}

// NewClient creates a new OpenAI client. Uses gpt-3.5-turbo by default.
func NewClient(logger *slog.Logger, config Config) *Client {
	return &Client{
		logger:     logger.WithGroup("openai-client"),
		config:     config,
		httpClient: http.DefaultClient,
	}
}

// Stream sends a streaming generation request and returns a channel of incremental responses.
func (c *Client) Stream(ctx context.Context, request llm.GenerateRequest) (<-chan llm.GenerateStreamResponse, error) {
	// Internal types for OpenAI API payload
	type chatMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	type streamRequest struct {
		Model    string        `json:"model"`
		Messages []chatMessage `json:"messages"`
		Stream   bool          `json:"stream"`
	}

	// Build the messages sequence: system, history, then new user prompt
	msgs := []chatMessage{{Role: "system", Content: request.Prompt.System}}
	for _, conv := range request.History {
		msgs = append(msgs, chatMessage{Role: "user", Content: conv.Request})
		msgs = append(msgs, chatMessage{Role: "assistant", Content: conv.Response})
	}
	msgs = append(msgs, chatMessage{Role: "user", Content: request.Prompt.User})

	reqBody := streamRequest{Model: c.config.Model, Messages: msgs, Stream: true}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	completionsEndpoint := fmt.Sprintf("%s/chat/completions", c.config.Endpoint)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", completionsEndpoint, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.config.ApiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("stream request error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		errBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("stream request failed: %s", string(errBody))
	}

	ch := make(chan llm.GenerateStreamResponse)
	// Start goroutine to read the SSE stream
	go func() {
		defer close(ch)
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					ch <- llm.GenerateStreamResponse{Error: err}
				}
				return
			}
			line = strings.TrimSpace(line)
			if line == "" || !strings.HasPrefix(line, "data:") {
				continue
			}
			payload := strings.TrimPrefix(line, "data: ")
			if payload == "[DONE]" {
				return
			}
			// Parse the streamed JSON chunk
			var event struct {
				ID      string `json:"id"`
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}
			if err := json.Unmarshal([]byte(payload), &event); err != nil {
				ch <- llm.GenerateStreamResponse{Error: fmt.Errorf("invalid stream response: %w", err)}
				return
			}

			// Emit each non-empty content delta
			for _, choice := range event.Choices {
				if choice.Delta.Content != "" {
					ch <- llm.GenerateStreamResponse{ID: event.ID, Response: choice.Delta.Content}
				}
			}
		}
	}()

	return ch, nil
}
