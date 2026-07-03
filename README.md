# golang-exploration-1

A Go project exploring Google's [Agent Development Kit (ADK) for Go v2](https://github.com/google/adk-go).

## Prerequisites

- Go 1.25+
- A [Google AI Studio API key](https://aistudio.google.com/app/api-keys)

## Setup

1. **Clone the repo**
   ```bash
   git clone https://github.com/ashishkhapre/golang-exploration-1.git
   cd golang-exploration-1
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set your API key**
   ```bash
   export GOOGLE_API_KEY="your-api-key-here"
   ```

## Running the Agent

```bash
go run ./cmd/agent/
```

The `full` launcher provides two interaction modes:

- **CLI mode** (default): Interactive chat in the terminal.
- **Web UI mode**: Run `go run ./cmd/agent/ web` to start a local web server.

## Configuration

The agent can be tuned via environment variables:

| Variable           | Default                                         | Description                  |
|--------------------|-------------------------------------------------|------------------------------|
| `GOOGLE_API_KEY`   | *(required)*                                    | Google AI Studio API key     |
| `AGENT_MODEL`      | `gemini-2.5-flash`                              | Gemini model name            |
| `AGENT_NAME`       | `exploration_agent`                             | Unique agent identifier      |
| `AGENT_INSTRUCTION`| `You are a helpful assistant. Answer user questions clearly and concisely.` | System prompt |

## Project Structure

```
.
├── cmd/
│   └── agent/
│       └── main.go          # Entry point — wires config + agent + launcher
├── internal/
│   └── agent/
│       └── agent.go         # ADK agent construction (module-private)
├── pkg/
│   └── config/
│       └── config.go        # Config loading from environment variables
├── go.mod                   # Go module definition (google.golang.org/adk/v2)
├── go.sum                   # Dependency checksums
└── README.md                # This file
```

## References

- [ADK Go Documentation](https://google.github.io/adk-docs/)
- [ADK Go GitHub Repository](https://github.com/google/adk-go)
- [ADK Go Examples](https://github.com/google/adk-go/tree/main/examples)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
