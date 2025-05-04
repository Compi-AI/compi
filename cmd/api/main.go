package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	http2 "github.com/compiai/engine/internal/api/http"
	"github.com/compiai/engine/internal/core/domain/agent/stat_analyzer"
	"github.com/compiai/engine/pkg/llm"
	"github.com/compiai/engine/pkg/llm/claude"
	"github.com/compiai/engine/pkg/llm/openai"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"

	promptloader "github.com/compiai/engine/internal/core/domain/agent/stat_analyzer/prompt"
	domainuser "github.com/compiai/engine/internal/core/domain/user"
	"github.com/compiai/engine/internal/core/ext/storage"
	"log/slog"
)

// AppConfig holds application configuration loaded from YAML.
type AppConfig struct {
	Application struct {
		Version string `yaml:"version"`

		Database struct {
			Postgres struct {
				Addr string `yaml:"addr"`
				Auth struct {
					Username   string `yaml:"username"`
					Password   string `yaml:"password"`
					TLSEnabled bool   `yaml:"tlsEnabled"`
				} `yaml:"auth"`
			} `yaml:"postgres"`
		} `yaml:"database"`

		Clients struct {
			Claude claude.Config `yaml:"claude"`
			OpenAI openai.Config `yaml:"openai"`
		} `yaml:"clients"`

		Server struct {
			Public struct {
				Addr    string        `yaml:"addr"`
				Timeout time.Duration `yaml:"timeout"`
			} `yaml:"public"`
		} `yaml:"server"`

		Auth struct {
			PrivateKey string `yaml:"privateKey"`
			PublicKey  string `yaml:"publicKey"`
		} `yaml:"auth"`
	} `yaml:"application"`
}

// loadConfig reads and parses the YAML config file at the given path.
func loadConfig(path string) (*AppConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg AppConfig
	if err := yaml.UnmarshalStrict(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &cfg, nil
}

func main() {
	// Parse config file path
	configPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	// Load configuration
	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	logger.Info("starting application", "version", cfg.Application.Version)

	// Build Postgres DSN
	sslMode := "disable"
	if cfg.Application.Database.Postgres.Auth.TLSEnabled {
		sslMode = "require"
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s/?sslmode=%s",
		cfg.Application.Database.Postgres.Auth.Username,
		cfg.Application.Database.Postgres.Auth.Password,
		cfg.Application.Database.Postgres.Addr,
		sslMode,
	)

	// Connect to Postgres
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("database connection failed", "err", err)
		os.Exit(1)
	}
	defer db.Close()
	if err := db.PingContext(context.Background()); err != nil {
		logger.Error("database ping failed", "err", err)
		os.Exit(1)
	}

	// Initialize storage and domain services
	userStorage := storage.NewPostgresStorage(db)
	userService := domainuser.NewService(logger, userStorage)

	// Initialize LLM clients
	openaiClient := openai.NewClient(logger, cfg.Application.Clients.OpenAI)
	// claudeClient := claude.NewClaudeClient(logger, cfg.Application.Clients.Claude)

	// Choose streamer; here using OpenAI but could swap to Claude
	var llmStreamer llm.Streamer = openaiClient

	// Initialize prompt loader
	pl, err := promptloader.NewPromptLoader()
	if err != nil {
		logger.Error("prompt loader init failed", "err", err)
		os.Exit(1)
	}

	// Initialize agent
	statAgent := stat_analyzer.NewAgent(logger, llmStreamer, *pl, userService)

	// Setup HTTP router
	r := chi.NewRouter()
	http2.RegisterRoutes(r, statAgent, logger)

	// Start HTTP server
	srv := &http.Server{
		Addr:         cfg.Application.Server.Public.Addr,
		Handler:      r,
		ReadTimeout:  cfg.Application.Server.Public.Timeout,
		WriteTimeout: cfg.Application.Server.Public.Timeout,
		IdleTimeout:  cfg.Application.Server.Public.Timeout,
	}
	logger.Info("listening on", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", "err", err)
	}
}
