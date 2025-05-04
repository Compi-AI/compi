package claude

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/compiai/engine/pkg/llm"
	"io"
	"net/http"
	"strings"
)

type Config struct {
	ApiKey            string  `yaml:"apiKey"`
	Endpoint          string  `yaml:"endpoint"`
	Model             string  `yaml:"model"`
	Temperature       float64 `yaml:"temperature"`
	MaxTokensToSample int     `yaml:"maxTokensToSample"`
}

// ClaudeClient implements llm.Streamer using Anthropic's Claude API.
type ClaudeClient struct {
	config     Config
	httpClient *http.Client
}

// NewClaudeClient creates a new Claude client using provided Config.
func NewClaudeClient(cfg Config) *ClaudeClient {
	return &ClaudeClient{config: cfg, httpClient: http.DefaultClient}
}

func (c *ClaudeClient) Stream(ctx context.Context, request llm.GenerateRequest) (<-chan llm.GenerateStreamResponse, error) {

	// Build the Anthropic prompt
	var sb strings.Builder
	if request.Prompt.System != "" {
		sb.WriteString("\n\nHuman: ")
		sb.WriteString(request.Prompt.System)
	}
	for _, conv := range request.History {
		sb.WriteString("\n\nHuman: ")
		sb.WriteString(conv.Request)
		sb.WriteString("\n\nAssistant: ")
		sb.WriteString(conv.Response)
	}
	sb.WriteString("\n\nHuman: ")
	sb.WriteString(request.Prompt.User)
	sb.WriteString("\n\nAssistant:")

	reqBody := apiRequest{
		Model:             c.config.Model,
		Prompt:            sb.String(),
		MaxTokensToSample: c.config.MaxTokensToSample,
		Temperature:       c.config.Temperature,
		Stream:            true,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("Claude marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		strings.TrimRight(c.config.Endpoint, "/"), bytes.NewReader(data),
	)
	if err != nil {
		return nil, fmt.Errorf("Claude new request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", c.config.ApiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Claude do request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		errBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Claude error response: %s", string(errBody))
	}

	ch := make(chan llm.GenerateStreamResponse)
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
			if !strings.HasPrefix(line, "data:") {
				continue
			}
			payload := strings.TrimPrefix(line, "data: ")
			if payload == "[DONE]" {
				return
			}
			var event struct {
				Completion string `json:"completion"`
			}
			if err := json.Unmarshal([]byte(payload), &event); err != nil {
				ch <- llm.GenerateStreamResponse{Error: fmt.Errorf("Claude invalid stream: %w", err)}
				return
			}
			if event.Completion != "" {
				ch <- llm.GenerateStreamResponse{Response: event.Completion}
			}
		}
	}()

	return ch, nil
}
