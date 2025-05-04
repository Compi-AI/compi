package stat_analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	prompts "github.com/compiai/engine/internal/core/domain/agent/stat_analyzer/prompt"
	"github.com/compiai/engine/internal/core/domain/user"
	"github.com/compiai/engine/pkg/llm"
	"github.com/google/uuid"
	"log/slog"
	"math"
	"time"
)

type BuildAnalysisRequest struct {
	UserID uuid.UUID
}

type BuildAnalysisStreamResponse struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Error   error  `json:"error"`
}

type Agent interface {
	BuildAnalysis(ctx context.Context, request BuildAnalysisRequest) (<-chan BuildAnalysisStreamResponse, error)
}

type agent struct {
	logger       *slog.Logger
	llmStreamer  llm.Streamer
	promptLoader prompts.PromptLoader
	userService  user.Service
}

// Agent defines the streaming analysis interface
func NewAgent(
	logger *slog.Logger,
	streamer llm.Streamer,
	loader prompts.PromptLoader,
	usrSvc user.Service,
) Agent {
	return &agent{
		logger:       logger.WithGroup("stat-analyzer-agent"),
		llmStreamer:  streamer,
		promptLoader: loader,
		userService:  usrSvc,
	}
}

// BuildAnalysis performs a multi-stage stats processing and streams an LLM-based analysis
func (a *agent) BuildAnalysis(ctx context.Context, req BuildAnalysisRequest) (<-chan BuildAnalysisStreamResponse, error) {
	// Stage 1: Load raw user profile
	usr, err := a.userService.FindOne(ctx, user.SingleFilter{ID: &req.UserID})
	if err != nil {
		a.logger.Error("user lookup failed", "err", err)
		return nil, fmt.Errorf("user lookup: %w", err)
	}

	// Stage 2: Derive advanced metrics
	advanced := deriveMetrics(usr)

	// Stage 3: Compose combined data for prompting
	data := struct {
		Profile  user.User              `json:"profile"`
		Advanced map[string]interface{} `json:"advancedMetrics"`
		Request  BuildAnalysisRequest   `json:"request"`
	}{usr, advanced, req}

	a.logger.Info("build started...", "data", data)

	// Stage 4: Render prompts
	sys := a.promptLoader.GetSystemPrompt()

	usrPr := a.promptLoader.GetUserPrompt()

	// Stage 5: Initiate LLM streaming
	genReq := llm.GenerateRequest{Prompt: llm.Prompt{System: sys, User: usrPr}}
	stream, err := a.llmStreamer.Stream(ctx, genReq)
	if err != nil {
		return nil, fmt.Errorf("LLM stream: %w", err)
	}

	// Stage 6: Process and enrich stream responses
	out := make(chan BuildAnalysisStreamResponse)
	go func() {
		defer close(out)
		for msg := range stream {
			// Inject timestamp metadata and segment tags
			segment := time.Now().Format(time.RFC3339Nano)
			meta := map[string]string{"segment": segment}
			if payload, err := json.Marshal(meta); err == nil {
				msg.Response = string(payload) + "\n" + msg.Response
			}
			out <- BuildAnalysisStreamResponse{ID: msg.ID, Content: msg.Response, Error: msg.Error}
		}
	}()

	return out, nil
}

// deriveMetrics computes advanced analytic metrics from user data
func deriveMetrics(u user.User) map[string]interface{} {
	metrics := make(map[string]interface{})

	// Example: Kill/Death Ratio
	kd := float64(countTotalKills(u.Games)) / math.Max(1, float64(countTotalDeaths(u.Games)))
	metrics["killDeathRatio"] = fmt.Sprintf("%.2f", kd)

	// Example: Engagement consistency score via variance
	scores := extractEngagementScores(u.Games)
	metrics["engagementConsistency"] = variance(scores)

	// Example: Popular role synergy cluster (mocked)
	metrics["synergyCluster"] = clusterRoleSynergy(u.Games)

	return metrics
}

// countTotalKills tallies kills (stub implementation)
func countTotalKills(games []string) int {
	// placeholder: parse and sum
	return len(games) * 5
}

// countTotalDeaths tallies deaths (stub implementation)
func countTotalDeaths(games []string) int {
	return len(games) * 3
}

// extractEngagementScores returns mock engagement values
func extractEngagementScores(games []string) []float64 {
	var s []float64
	for i := range games {
		s = append(s, float64((i%5)+1))
	}
	return s
}

// variance computes variance of a float64 slice
func variance(data []float64) float64 {
	mean := 0.0
	for _, v := range data {
		mean += v
	}
	mean /= float64(len(data))
	varSum := 0.0
	for _, v := range data {
		varSum += (v - mean) * (v - mean)
	}
	return varSum / float64(len(data))
}

// clusterRoleSynergy returns a stub cluster name
func clusterRoleSynergy(games []string) string {
	// pretend clustering logic
	return "Alpha-Synergy"
}
