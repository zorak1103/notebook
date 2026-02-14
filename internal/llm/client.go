package llm

import (
	"context"
	"fmt"
	"strings"
)

// Provider defines the interface for LLM completion providers
type Provider interface {
	Complete(ctx context.Context, prompt string) (string, error)
}

// Client wraps an LLM provider for completions
type Client struct {
	provider Provider
}

// New creates a new LLM client based on the provider URL
// URL-based detection: anthropic.com -> AnthropicProvider, everything else -> OpenAIProvider
func New(providerURL, apiKey, model string) (*Client, error) {
	var provider Provider
	var err error

	if strings.Contains(providerURL, "anthropic.com") {
		provider, err = NewAnthropicProvider(apiKey, model)
	} else {
		// Everything else uses OpenAI-compatible format (OpenAI, Azure, Ollama, LM Studio, vLLM, etc.)
		provider, err = NewOpenAIProvider(providerURL, apiKey, model)
	}

	if err != nil {
		return nil, fmt.Errorf("create provider: %w", err)
	}

	return &Client{provider: provider}, nil
}

// Complete sends a prompt to the LLM and returns the completion
func (c *Client) Complete(ctx context.Context, prompt string) (string, error) {
	return c.provider.Complete(ctx, prompt)
}
