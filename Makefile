# Makefile for compiai-engine

# Binary name
BINARY_NAME := compiai-engine
# Configuration file
CONFIG_FILE := config.yaml

# Go settings
GO      := go
GOCMD   := $(GO)
GOBUILD := $(GO) build
GOCLEAN := $(GO) clean
GOTEST  := $(GO) test
GOFMT   := $(GO) fmt
GOVET   := $(GO) vet
GOMOD   := $(GO) mod

# Directories
PKG     := ./...
BIN_DIR := bin

# Docker settings
DOCKER_IMAGE := compiai/engine:latest
DOCKERFILE   := Dockerfile

.PHONY: all build run test fmt vet lint tidy clean docker-build help

all: build

## build: compile the binary
build:
	@echo "==> Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) -ldflags="-s -w" main.go

## run: run the application using the config file
run: build
	@echo "==> Running $(BINARY_NAME)..."
	@$(BIN_DIR)/$(BINARY_NAME) -config $(CONFIG_FILE)

## test: run unit tests
test:
	@echo "==> Running tests..."
	$(GOTEST) -v $(PKG)

## fmt: format Go code
fmt:
	@echo "==> Formatting code..."
	$(GOFMT) $(PKG)

## vet: report potential issues
vet:
	@echo "==> Vetting code..."
	$(GOVET) $(PKG)

## lint: run linter (requires golangci-lint)
lint:
	@echo "==> Linting code..."
	golangci-lint run

## tidy: ensure go.mod matches imports
tidy:
	@echo "==> Tidying modules..."
	$(GOMOD) tidy

## clean: remove build artifacts
clean:
	@echo "==> Cleaning up..."
	@rm -rf $(BIN_DIR)
	$(GOCLEAN)

## docker-build: build Docker image
docker-build:
	@echo "==> Building Docker image $(DOCKER_IMAGE)..."
	docker build -t $(DOCKER_IMAGE) -f $(DOCKERFILE) .

## help: display this help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%]()*_]()
