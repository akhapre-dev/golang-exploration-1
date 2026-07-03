// Package agent provides construction and configuration of the ADK LLM agent.
// This package is internal and not intended for use outside this module.
package agent

import (
	"context"
	"fmt"

	adkagent "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model/gemini"
	"google.golang.org/genai"

	"github.com/akhapre-dev/golang-exploration-1/pkg/config"
)

// New constructs and returns a new LLM agent using the provided configuration.
// It initialises the Gemini model and wires it into an ADK llmagent.
func New(ctx context.Context, cfg config.Config) (adkagent.Agent, error) {
	model, err := gemini.NewModel(ctx, cfg.ModelName, &genai.ClientConfig{
		APIKey: cfg.APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("creating gemini model: %w", err)
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        cfg.AgentName,
		Model:       model,
		Instruction: cfg.Instruction,
	})
	if err != nil {
		return nil, fmt.Errorf("creating llm agent: %w", err)
	}

	return a, nil
}
