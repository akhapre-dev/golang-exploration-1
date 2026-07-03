// Command agent is the entry point for the ADK exploration agent.
//
// Usage:
//
//	go run ./cmd/agent/           → interactive console chat
//	go run ./cmd/agent/ web       → start local Web UI server
//
// Required environment variables:
//
//	GOOGLE_API_KEY    Google AI Studio API key
//
// Optional environment variables:
//
//	AGENT_MODEL       Gemini model name (default: gemini-2.0-flash)
//	AGENT_NAME        Agent identifier  (default: exploration_agent)
//	AGENT_INSTRUCTION System prompt     (default: helpful assistant)
package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/cmd/launcher"
	"google.golang.org/adk/v2/cmd/launcher/full"

	internalagent "github.com/akhapre-dev/golang-exploration-1/internal/agent"
	"github.com/akhapre-dev/golang-exploration-1/pkg/config"
)

func main() {
	ctx := context.Background()

	// Load configuration from environment variables.
	cfg := config.Load()

	// Build the ADK agent.
	a, err := internalagent.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// Wrap the agent in a single-agent loader (required by the launcher).
	agentLoader := agent.NewSingleLoader(a)

	// Build the launcher config.
	launcherCfg := &launcher.Config{
		AgentLoader: agentLoader,
	}

	// Launch — supports both console (default) and web sub-launchers.
	l := full.NewLauncher()
	if err := l.Execute(ctx, launcherCfg, os.Args[1:]); err != nil {
		log.Fatalf("Agent run failed: %v", err)
	}
}
