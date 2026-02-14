package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zorak1103/notebook/internal/db/repositories"
)

func TestHandleGetConfig_Empty(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	w := httptest.NewRecorder()

	srv.handleGetConfig(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp ConfigData
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.LLMProviderURL != "" {
		t.Errorf("expected empty provider URL, got %q", resp.LLMProviderURL)
	}
	if resp.LLMAPIKey != "" {
		t.Errorf("expected empty API key, got %q", resp.LLMAPIKey)
	}
	if resp.LLMModel != "" {
		t.Errorf("expected empty model, got %q", resp.LLMModel)
	}
	if resp.Language != "en" {
		t.Errorf("expected default language 'en', got %q", resp.Language)
	}
	// Prompt fields should have default values from migration
	if resp.LLMPromptSummary == "" {
		t.Error("expected non-empty summary prompt from migration")
	}
	if resp.LLMPromptEnhance == "" {
		t.Error("expected non-empty enhance prompt from migration")
	}
}

func TestHandleGetConfig_WithValues(t *testing.T) {
	srv := newTestServer(t)
	repo := repositories.NewConfigRepository(srv.database.DB)

	// Set config values
	if err := repo.Set("llm_provider_url", "https://api.openai.com/v1"); err != nil {
		t.Fatalf("failed to set provider URL: %v", err)
	}
	if err := repo.Set("llm_api_key", "sk-1234567890abcdefghijklmnop"); err != nil {
		t.Fatalf("failed to set API key: %v", err)
	}
	if err := repo.Set("llm_model", "gpt-4o"); err != nil {
		t.Fatalf("failed to set model: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	w := httptest.NewRecorder()

	srv.handleGetConfig(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp ConfigData
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.LLMProviderURL != "https://api.openai.com/v1" {
		t.Errorf("expected provider URL 'https://api.openai.com/v1', got %q", resp.LLMProviderURL)
	}
	if resp.LLMModel != "gpt-4o" {
		t.Errorf("expected model 'gpt-4o', got %q", resp.LLMModel)
	}
	// API key should be masked
	if resp.LLMAPIKey == "sk-1234567890abcdefghijklmnop" {
		t.Error("API key should be masked, got plain text")
	}
	if resp.LLMAPIKey != "sk-1*********************mnop" {
		t.Errorf("expected masked key 'sk-1*********************mnop', got %q", resp.LLMAPIKey)
	}
}

func TestHandleGetConfig_MasksShortKey(t *testing.T) {
	srv := newTestServer(t)
	repo := repositories.NewConfigRepository(srv.database.DB)

	// Set short API key (8 chars)
	if err := repo.Set("llm_api_key", "short123"); err != nil {
		t.Fatalf("failed to set API key: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	w := httptest.NewRecorder()

	srv.handleGetConfig(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp ConfigData
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Short keys should be all asterisks
	if resp.LLMAPIKey != "********" {
		t.Errorf("expected all asterisks for short key, got %q", resp.LLMAPIKey)
	}
}

func TestHandleUpdateConfig_Success(t *testing.T) {
	srv := newTestServer(t)

	reqBody := ConfigUpdateRequest{
		LLMProviderURL: "https://api.anthropic.com/v1",
		LLMAPIKey:      "sk-ant-1234567890abcdefghijklmnop",
		LLMModel:       "claude-opus-4-6",
		Language:       "de",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/config", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.handleUpdateConfig(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp ConfigData
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.LLMProviderURL != "https://api.anthropic.com/v1" {
		t.Errorf("expected provider URL 'https://api.anthropic.com/v1', got %q", resp.LLMProviderURL)
	}
	if resp.LLMModel != "claude-opus-4-6" {
		t.Errorf("expected model 'claude-opus-4-6', got %q", resp.LLMModel)
	}
	if resp.Language != "de" {
		t.Errorf("expected language 'de', got %q", resp.Language)
	}
	// API key should be masked
	if resp.LLMAPIKey == "sk-ant-1234567890abcdefghijklmnop" {
		t.Error("API key should be masked, got plain text")
	}
}

func TestHandleUpdateConfig_InvalidJSON(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest(http.MethodPost, "/api/config", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.handleUpdateConfig(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleUpdateConfig_InvalidURL(t *testing.T) {
	srv := newTestServer(t)

	reqBody := ConfigUpdateRequest{
		LLMProviderURL: "not-a-valid-url",
		LLMAPIKey:      "sk-test",
		LLMModel:       "model",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/config", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.handleUpdateConfig(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleUpdateConfig_SkipsMaskedKey(t *testing.T) {
	srv := newTestServer(t)
	repo := repositories.NewConfigRepository(srv.database.DB)

	// Set initial API key
	originalKey := "sk-original-key-12345678"
	if err := repo.Set("llm_api_key", originalKey); err != nil {
		t.Fatalf("failed to set initial API key: %v", err)
	}

	// Try to update with masked key
	reqBody := ConfigUpdateRequest{
		LLMProviderURL: "https://api.example.com",
		LLMAPIKey:      "sk-o***************5678", // masked
		LLMModel:       "test-model",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/config", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.handleUpdateConfig(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify original key is preserved
	storedKey, err := repo.Get("llm_api_key")
	if err != nil {
		t.Fatalf("failed to get API key: %v", err)
	}
	if storedKey == nil || storedKey.Value != originalKey {
		t.Errorf("expected original key to be preserved, got %v", storedKey)
	}
}

func TestHandleUpdateConfig_EmptyFieldsSkipped(t *testing.T) {
	srv := newTestServer(t)
	repo := repositories.NewConfigRepository(srv.database.DB)

	// Set initial values
	if err := repo.Set("llm_provider_url", "https://api.initial.com"); err != nil {
		t.Fatalf("failed to set initial URL: %v", err)
	}
	if err := repo.Set("llm_api_key", "initial-key"); err != nil {
		t.Fatalf("failed to set initial key: %v", err)
	}
	if err := repo.Set("llm_model", "initial-model"); err != nil {
		t.Fatalf("failed to set initial model: %v", err)
	}

	// Update with only model field
	reqBody := ConfigUpdateRequest{
		LLMProviderURL: "", // empty - should be skipped
		LLMAPIKey:      "", // empty - should be skipped
		LLMModel:       "new-model",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/config", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.handleUpdateConfig(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify URL and key are preserved
	url, err := repo.Get("llm_provider_url")
	if err != nil {
		t.Fatalf("failed to get URL: %v", err)
	}
	if url == nil || url.Value != "https://api.initial.com" {
		t.Errorf("expected URL to be preserved, got %v", url)
	}

	key, err := repo.Get("llm_api_key")
	if err != nil {
		t.Fatalf("failed to get key: %v", err)
	}
	if key == nil || key.Value != "initial-key" {
		t.Errorf("expected key to be preserved, got %v", key)
	}

	// Verify model was updated
	model, err := repo.Get("llm_model")
	if err != nil {
		t.Fatalf("failed to get model: %v", err)
	}
	if model == nil || model.Value != "new-model" {
		t.Errorf("expected model 'new-model', got %v", model)
	}
}

func TestMaskAPIKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "short key (8 chars)",
			input:    "short123",
			expected: "********",
		},
		{
			name:     "short key (less than 8 chars)",
			input:    "short",
			expected: "*****",
		},
		{
			name:     "long key",
			input:    "sk-1234567890abcdefghijklmnop",
			expected: "sk-1*********************mnop",
		},
		{
			name:     "exactly 9 chars",
			input:    "123456789",
			expected: "1234*6789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskAPIKey(tt.input)
			if result != tt.expected {
				t.Errorf("maskAPIKey(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHandleUpdateConfig_LanguagePersistence(t *testing.T) {
	srv := newTestServer(t)

	// Update language to French
	reqBody := ConfigUpdateRequest{
		Language: "fr",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/config", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.handleUpdateConfig(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify language is persisted via GET
	req2 := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	w2 := httptest.NewRecorder()

	srv.handleGetConfig(w2, req2)

	var resp ConfigData
	if err := json.NewDecoder(w2.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Language != "fr" {
		t.Errorf("expected language 'fr' to be persisted, got %q", resp.Language)
	}
}

func TestHandleUpdateConfig_InvalidLanguage(t *testing.T) {
	srv := newTestServer(t)

	reqBody := ConfigUpdateRequest{
		Language: "invalid",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/config", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.handleUpdateConfig(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleUpdateConfig_PromptPersistence(t *testing.T) {
	srv := newTestServer(t)

	customSummaryPrompt := "Custom summary prompt with {{subject}} and {{notes}}"
	customEnhancePrompt := "Custom enhance prompt with {{content}}"

	// Update prompts
	reqBody := ConfigUpdateRequest{
		LLMPromptSummary: customSummaryPrompt,
		LLMPromptEnhance: customEnhancePrompt,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/config", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.handleUpdateConfig(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify prompts are persisted via GET
	req2 := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	w2 := httptest.NewRecorder()

	srv.handleGetConfig(w2, req2)

	var resp ConfigData
	if err := json.NewDecoder(w2.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.LLMPromptSummary != customSummaryPrompt {
		t.Errorf("expected summary prompt %q, got %q", customSummaryPrompt, resp.LLMPromptSummary)
	}
	if resp.LLMPromptEnhance != customEnhancePrompt {
		t.Errorf("expected enhance prompt %q, got %q", customEnhancePrompt, resp.LLMPromptEnhance)
	}
}

func TestIsMasked(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "plain text",
			input:    "sk-1234567890abcdefghijklmnop",
			expected: false,
		},
		{
			name:     "masked key",
			input:    "sk-1*********************mnop",
			expected: true,
		},
		{
			name:     "all asterisks",
			input:    "********",
			expected: true,
		},
		{
			name:     "single asterisk",
			input:    "*",
			expected: true,
		},
		{
			name:     "mixed with asterisks",
			input:    "abc*def",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isMasked(tt.input)
			if result != tt.expected {
				t.Errorf("isMasked(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
