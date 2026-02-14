package llm

import (
	"testing"
)

func TestNew_ProviderDetection(t *testing.T) {
	tests := []struct {
		name         string
		providerURL  string
		apiKey       string
		model        string
		expectError  bool
		providerType string
	}{
		{
			name:         "anthropic URL",
			providerURL:  "https://api.anthropic.com/v1",
			apiKey:       "sk-ant-test",
			model:        "claude-opus-4-6",
			expectError:  false,
			providerType: "*llm.AnthropicProvider",
		},
		{
			name:         "openai URL",
			providerURL:  "https://api.openai.com/v1",
			apiKey:       "sk-test",
			model:        "gpt-4o",
			expectError:  false,
			providerType: "*llm.OpenAIProvider",
		},
		{
			name:         "unknown URL defaults to OpenAI",
			providerURL:  "https://api.example.com",
			apiKey:       "test-key",
			model:        "test-model",
			expectError:  false,
			providerType: "*llm.OpenAIProvider",
		},
		{
			name:        "empty API key for Anthropic",
			providerURL: "https://api.anthropic.com/v1",
			apiKey:      "",
			model:       "claude-opus-4-6",
			expectError: true,
		},
		{
			name:        "empty API key for OpenAI",
			providerURL: "https://api.openai.com/v1",
			apiKey:      "",
			model:       "gpt-4o",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.providerURL, tt.apiKey, tt.model)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if client == nil {
				t.Fatal("expected non-nil client")
			}

			// Check provider type
			switch tt.providerType {
			case "*llm.AnthropicProvider":
				if _, ok := client.provider.(*AnthropicProvider); !ok {
					t.Errorf("expected AnthropicProvider, got %T", client.provider)
				}
			case "*llm.OpenAIProvider":
				if _, ok := client.provider.(*OpenAIProvider); !ok {
					t.Errorf("expected OpenAIProvider, got %T", client.provider)
				}
			}
		})
	}
}

func TestRenderPrompt(t *testing.T) {
	tests := []struct {
		name     string
		template string
		vars     map[string]string
		expected string
	}{
		{
			name:     "single variable",
			template: "Hello {{name}}",
			vars:     map[string]string{"name": "World"},
			expected: "Hello World",
		},
		{
			name:     "multiple variables",
			template: "Meeting: {{subject}} on {{date}}",
			vars:     map[string]string{"subject": "Standup", "date": "2024-01-15"},
			expected: "Meeting: Standup on 2024-01-15",
		},
		{
			name:     "no variables",
			template: "Static text",
			vars:     map[string]string{},
			expected: "Static text",
		},
		{
			name:     "unused variables",
			template: "Hello {{name}}",
			vars:     map[string]string{"name": "World", "unused": "value"},
			expected: "Hello World",
		},
		{
			name:     "missing variables",
			template: "Hello {{name}} and {{other}}",
			vars:     map[string]string{"name": "World"},
			expected: "Hello World and {{other}}",
		},
		{
			name:     "multiline template",
			template: "Subject: {{subject}}\nNotes:\n{{notes}}",
			vars:     map[string]string{"subject": "Meeting", "notes": "Point 1\nPoint 2"},
			expected: "Subject: Meeting\nNotes:\nPoint 1\nPoint 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderPrompt(tt.template, tt.vars)
			if result != tt.expected {
				t.Errorf("RenderPrompt() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestNewOpenAIProvider_DefaultModel(t *testing.T) {
	provider, err := NewOpenAIProvider("https://api.openai.com/v1", "test-key", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if provider.model != "gpt-4o-mini" {
		t.Errorf("expected default model 'gpt-4o-mini', got %q", provider.model)
	}
}

func TestNewAnthropicProvider_DefaultModel(t *testing.T) {
	provider, err := NewAnthropicProvider("test-key", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if provider.model != "claude-sonnet-4-20250514" {
		t.Errorf("expected default model 'claude-sonnet-4-20250514', got %q", provider.model)
	}
}
