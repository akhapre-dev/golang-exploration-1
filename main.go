package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/cmd/launcher"
	"google.golang.org/adk/v2/cmd/launcher/full"
	"google.golang.org/adk/v2/model/gemini"
	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()

	// Initialize the Gemini model.
	// Requires GOOGLE_API_KEY environment variable to be set.
	model, err := gemini.NewModel(ctx, "gemini-2.0-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	// Define the agent.
	myAgent, err := llmagent.New(llmagent.Config{
		Name:        "exploration_agent",
		Model:       model,
		Instruction: "You are a helpful assistant. Answer user questions clearly and concisely.",
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// Wrap the agent in a loader (required by the launcher config).
	agentLoader := agent.NewSingleLoader(myAgent)

	// Build the launcher config.
	cfg := &launcher.Config{
		AgentLoader: agentLoader,
	}

	// Launch the agent.
	// The full launcher supports both console (default) and web sub-launchers.
	// Usage:
	//   go run main.go           → interactive console chat
	//   go run main.go web       → start local Web UI server
	l := full.NewLauncher()
	if err := l.Execute(ctx, cfg, os.Args[1:]); err != nil {
		log.Fatalf("Agent run failed: %v", err)
	}
}
