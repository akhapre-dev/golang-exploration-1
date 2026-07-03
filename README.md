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
go run main.go
```

The `full` launcher provides two interaction modes:

- **CLI mode** (default): Interactive chat in the terminal.
- **Web UI mode**: Run `go run main.go web` to start a local web server.

## Project Structure

```
.
├── main.go       # Agent definition and entry point
├── go.mod        # Go module definition (google.golang.org/adk/v2)
├── go.sum        # Dependency checksums
└── README.md     # This file
```

## References

- [ADK Go Documentation](https://google.github.io/adk-docs/)
- [ADK Go GitHub Repository](https://github.com/google/adk-go)
- [ADK Go Examples](https://github.com/google/adk-go/tree/main/examples)
