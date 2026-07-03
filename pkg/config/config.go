// Package config provides application-level configuration loaded from
// environment variables, with sensible defaults.
package config

import "os"

// Config holds the configuration for the ADK agent application.
type Config struct {
	// APIKey is the Google AI Studio API key (GOOGLE_API_KEY env var).
	APIKey string

	// ModelName is the Gemini model to use (AGENT_MODEL env var).
	// Defaults to "gemini-2.5-flash".
	ModelName string

	// AgentName is the unique name for the agent (AGENT_NAME env var).
	// Defaults to "exploration_agent".
	AgentName string

	// Instruction is the system instruction for the agent (AGENT_INSTRUCTION env var).
	// Defaults to a helpful assistant prompt.
	Instruction string
}

// Load reads configuration from environment variables and returns a Config.
// Missing optional values fall back to their defaults.
func Load() Config {
	modelName := os.Getenv("AGENT_MODEL")
	if modelName == "" {
		modelName = "gemini-2.5-flash"
	}

	agentName := os.Getenv("AGENT_NAME")
	if agentName == "" {
		agentName = "exploration_agent"
	}

	instruction := os.Getenv("AGENT_INSTRUCTION")
	if instruction == "" {
		instruction = "You are a helpful assistant. Answer user questions clearly and concisely."
	}

	return Config{
		APIKey:      os.Getenv("GOOGLE_API_KEY"),
		ModelName:   modelName,
		AgentName:   agentName,
		Instruction: instruction,
	}
}
