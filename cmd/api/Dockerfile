# syntax=docker/dockerfile:1

# --- Builder Stage ---
FROM golang:1.20-alpine AS builder

# Set working directory
WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o compiai-engine main.go


# --- Runtime Stage ---
FROM alpine:latest

# Install certs for HTTPS
RUN apk add --no-cache ca-certificates

# Create app directory
WORKDIR /root/

# Copy the compiled binary and config file
COPY --from=builder /app/compiai-engine ./compiai-engine
COPY config.yaml ./config.yaml

# Expose the application port
EXPOSE 8080

# Default entrypoint
ENTRYPOINT ["./compiai-engine", "-config", "config.yaml"]
