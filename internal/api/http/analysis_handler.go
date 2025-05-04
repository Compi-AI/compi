package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log/slog"

	"github.com/compiai/engine/internal/core/domain/agent/stat_analyzer"
)

// AnalysisRequest defines the payload for analysis.
type AnalysisRequest struct {
	UserID uuid.UUID `json:"userId"`
}

// Validate checks required fields in AnalysisRequest.
func (r *AnalysisRequest) Validate() error {
	if r.UserID == uuid.Nil {
		return errors.New("userId is required")
	}
	return nil
}

// AnalysisResponse is sent for each analysis segment.
type AnalysisResponse struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}

// RegisterRoutes mounts stat analyzer endpoints onto the router.
func RegisterRoutes(r chi.Router, agent stat_analyzer.Agent, logger *slog.Logger) {
	r.Route("/analysis", func(r chi.Router) {
		r.Post("/", makeAnalysisHandler(agent, logger))
	})
}

// makeAnalysisHandler streams analysis via Server-Sent Events.
func makeAnalysisHandler(agent stat_analyzer.Agent, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Decode and validate JSON request body.
		var reqModel AnalysisRequest
		if err := json.NewDecoder(req.Body).Decode(&reqModel); err != nil {
			logger.Error("invalid request body", "err", err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := reqModel.Validate(); err != nil {
			logger.Error("validation error", "err", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Set SSE headers.
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache, no-transform")
		w.Header().Set("Connection", "keep-alive")
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}

		// Initiate analysis stream.
		ctx := req.Context()
		stream, err := agent.BuildAnalysis(ctx, stat_analyzer.BuildAnalysisRequest{UserID: reqModel.UserID})
		if err != nil {
			logger.Error("analysis build error", "err", err)
			http.Error(w, "analysis error", http.StatusInternalServerError)
			return
		}

		// Stream SSE events.
		for resp := range stream {
			// Build and marshal response.
			res := AnalysisResponse{ID: resp.ID, Content: resp.Content}
			if resp.Error != nil {
				res.Error = resp.Error.Error()
			}
			payload, err := json.Marshal(res)
			if err != nil {
				logger.Error("marshal response failed", "err", err)
				continue
			}

			// Send SSE event.
			w.Write([]byte("event: analysis\n"))
			w.Write([]byte("data: "))
			w.Write(payload)
			w.Write([]byte("\n\n"))
			flusher.Flush()

			// Small delay to regulate stream.
			time.Sleep(10 * time.Millisecond)
		}

		// Final end event.
		w.Write([]byte("event: end\n"))
		w.Write([]byte("data: {}\n\n"))
		flusher.Flush()
	}
}
