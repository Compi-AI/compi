**COMPI AI**
*Competitive Gaming Performance Analysis Platform*

---

## Socials

- [Website](https://compiai.xyz)
- [Twitter](https://x.com/CompiAI)

## Table of Contents

1. [Overview](#overview)
2. [Features](#features)
3. [Architecture](#architecture)
4. [Getting Started](#getting-started)

    * [Prerequisites](#prerequisites)
    * [Installation](#installation)
    * [Configuration](#configuration)
    * [Running the Server](#running-the-server)
5. [Usage](#usage)

    * [API Endpoints](#api-endpoints)
    * [Example Request](#example-request)
6. [Development](#development)

    * [Project Layout](#project-layout)
    * [Testing](#testing)
7. [Roadmap](#roadmap)
8. [Future Plans](#future-plans)
9. [Contributing](#contributing)
10. [License](#license)

---

## Overview

**compiai-engine** is a Go-powered backend service that ingests player performance data, computes advanced analytics, and generates in-depth coaching feedback via large language models (OpenAI or Anthropic Claude). It exposes a Server-Sent Events (SSE) API for real-time streaming of step-by-step analysis, enabling integration into dashboards, coaching tools, or direct in-game overlays.

---

## Features

* **Advanced Metrics Derivation**
  Calculates K/D ratio, variance-based consistency, and role-synergy clustering.
* **LLM-Powered Narrative**
  Wraps data in tailored system/user prompts and streams expert analysis via Chat Completions.
* **SSE Streaming API**
  Delivers incremental coaching advice in real time.
* **Pluggable LLM Clients**
  Swap easily between OpenAI and Claude via a common `llm.Streamer` interface.
* **PostgreSQL Storage**
  Secure, scalable user profile persistence with JSON/array support.
* **Go Chi Router**
  Lightweight, idiomatic HTTP routing and middleware support.
* **Config-Driven**
  All endpoints, timeouts, and LLM parameters defined via YAML for Kubernetes or local runs.

---

## Architecture

```
+-------------+       +---------------+       +-------------+
|  HTTP Client| <---> |   HTTP Server | <---> |   Chi Router|
+-------------+       +---------------+       +-------------+
                                |
                                v
                    +-----------------------+
                    |    stat_analyzer      |
                    |   (Agent + Prompts)   |
                    +-----------------------+
                     |       |        |
                     v       v        v
             +--------+  +--------+  +----------------+
             | User   |  | LLM    |  | AdvancedMetrics|
             | Service|  | Stream |  | Derivation     |
             +--------+  +--------+  +----------------+
                 |                        |
                 v                        |
           +-------------+                |
           | PostgreSQL  |<---------------+
           |  Storage    |
           +-------------+
```

---

## Getting Started

### Prerequisites

* Go 1.20+
* PostgreSQL 12+
* Kubernetes (optional for production)

### Installation

```bash
git clone https://github.com/compiai/engine.git
cd engine
go mod download
```

### Configuration

Copy the sample configuration and fill in real values:

```yaml
# config.yaml
application:
  version: 0.1.0-snapshot
  database:
    postgres:
      addr: db.example.com:5432
      auth:
        username: youruser
        password: yourpass
        tlsEnabled: false
  clients:
    openai:
      apiKey: YOUR_OPENAI_KEY
      endpoint: https://api.openai.com/v1
      model: gpt-4
      temperature: 0.7
      maxTokensToSample: 500
    claude:
      apiKey: YOUR_CLAUDE_KEY
      endpoint: https://api.anthropic.com/v1/complete
      model: claude-v1
      temperature: 1.0
      maxTokensToSample: 300
  server:
    public:
      addr: :8080
      timeout: 5s
  auth:
    privateKey: file://path/to/private.key
    publicKey: file://path/to/public.key
```

### Running the Server

```bash
go run main.go -config ./config.yaml
```

---

## Usage

### API Endpoints

* **POST** `/analysis/`

    * **Request Body**:

      ```json
      { "userId": "00000000-0000-0000-0000-000000000000" }
      ```
    * **Response**: Server-Sent Events (`event: analysis`) streaming chunks of analysis JSON.

### Example Request

```bash
curl -N -X POST http://localhost:8080/analysis/ \
  -H "Content-Type: application/json" \
  -d '{"userId":"123e4567-e89b-12d3-a456-426614174000"}'
```

---

## Development

### Project Layout

```
.
├── cmd/                # main.go entrypoint
├── internal/
│   ├── core/
│   │   ├── ext/user/           # Postgres storage + user service
│   │   └── domain/
│   │       └── agent/stat_analyzer/
│   │           ├── prompt/     # text/template files
│   │           ├── agent.go
│   │           └── http/       # Chi HTTP handlers
├── pkg/
│   └── llm/             # Streamer interface
│   └── openai/          # OpenAI/Claude client implementations
└── config.yaml
```

### Testing

```bash
go test ./...
```

---

## Roadmap

* **v0.2.0**

    * Add support for batch analysis of multiple users concurrently.
    * Implement real‐time WebSocket fallback.
    * Integrate optional GPU-accelerated metric pipelines.

* **v0.3.0**

    * Dashboard UI with live graphs and session recordings.
    * Role‐based access control (coach vs. player).
    * Plugins for popular game titles (Dota 2, Valorant, League of Legends).

---

## Future Plans

* **Adaptive Learning**
  Leverage user progress data to auto-tune prompt recommendations and drill difficulty over time.
* **Multi-Modal Inputs**
  Incorporate gameplay video analysis (computer vision) alongside statistics.
* **Cross-Platform SDKs**
  Provide JavaScript and Python SDKs for easy integration into third-party tools.
* **Enterprise Features**
  SSO integration, analytics dashboards, and exportable PDF coaching reports.

---

## Contributing

We welcome contributions! Please:

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/XYZ`).
3. Commit your changes (`git commit -m "Add feature XYZ"`).
4. Push to your fork (`git push origin feature/XYZ`).
5. Open a pull request and describe your changes.

Please follow the existing code style and include tests for new functionality.

---

## License

[MIT License](LICENSE)
Feel free to use, modify, and distribute under the terms of the MIT license.
