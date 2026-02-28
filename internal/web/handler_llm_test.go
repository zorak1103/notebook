package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zorak1103/notebook/internal/db/models"
	"github.com/zorak1103/notebook/internal/db/repositories"
)

func TestHandleSummarizeMeeting_MissingConfig(t *testing.T) {
	srv := newTestServer(t)
	meetingRepo := repositories.NewMeetingRepository(srv.database.DB)

	// Create a meeting without setting LLM config
	meeting := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Test Meeting",
		MeetingDate: "2024-01-15",
		StartTime:   "10:00",
	}
	if err := meetingRepo.Create(meeting); err != nil {
		t.Fatalf("failed to create test meeting: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/meetings/1/summarize", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	srv.handleSummarizeMeeting(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleSummarizeMeeting_MeetingNotFound(t *testing.T) {
	srv := newTestServer(t)
	setTestLLMConfig(t, srv)

	req := httptest.NewRequest(http.MethodPost, "/api/meetings/999/summarize", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	srv.handleSummarizeMeeting(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestHandleSummarizeMeeting_NoNotes(t *testing.T) {
	srv := newTestServer(t)
	setTestLLMConfig(t, srv)
	meetingRepo := repositories.NewMeetingRepository(srv.database.DB)

	// Create a meeting without notes
	meeting := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Test Meeting",
		MeetingDate: "2024-01-15",
		StartTime:   "10:00",
	}
	if err := meetingRepo.Create(meeting); err != nil {
		t.Fatalf("failed to create test meeting: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/meetings/1/summarize", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	srv.handleSummarizeMeeting(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleEnhanceNote_MissingConfig(t *testing.T) {
	srv := newTestServer(t)
	meetingRepo := repositories.NewMeetingRepository(srv.database.DB)
	noteRepo := repositories.NewNoteRepository(srv.database.DB)

	// Create a meeting and note without setting LLM config
	meeting := &models.Meeting{
		CreatedBy:   "test@example.com",
		Subject:     "Test Meeting",
		MeetingDate: "2024-01-15",
		StartTime:   "10:00",
	}
	if err := meetingRepo.Create(meeting); err != nil {
		t.Fatalf("failed to create test meeting: %v", err)
	}

	note := &models.Note{
		MeetingID: meeting.ID,
		Content:   "Test note content",
	}
	if err := noteRepo.Create(note); err != nil {
		t.Fatalf("failed to create test note: %v", err)
	}

	body, _ := json.Marshal(enhanceNoteRequest{Content: "Test note content"})
	req := httptest.NewRequest(http.MethodPost, "/api/notes/1/enhance", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	srv.handleEnhanceNote(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleEnhanceNote_NoteNotFound(t *testing.T) {
	srv := newTestServer(t)
	setTestLLMConfig(t, srv)

	body, _ := json.Marshal(enhanceNoteRequest{Content: "Some content"})
	req := httptest.NewRequest(http.MethodPost, "/api/notes/999/enhance", bytes.NewReader(body))
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	srv.handleEnhanceNote(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestHandleEnhanceNote_EmptyContent(t *testing.T) {
	srv := newTestServer(t)

	body, _ := json.Marshal(enhanceNoteRequest{Content: ""})
	req := httptest.NewRequest(http.MethodPost, "/api/notes/1/enhance", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	srv.handleEnhanceNote(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleEnhanceNote_WhitespaceContent(t *testing.T) {
	srv := newTestServer(t)

	body, _ := json.Marshal(enhanceNoteRequest{Content: "   "})
	req := httptest.NewRequest(http.MethodPost, "/api/notes/1/enhance", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	srv.handleEnhanceNote(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleEnhanceNote_InvalidBody(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest(http.MethodPost, "/api/notes/1/enhance", bytes.NewReader([]byte("not json {")))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	srv.handleEnhanceNote(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestFormatNotes(t *testing.T) {
	notes := []*models.Note{
		{NoteNumber: 1, Content: "First note"},
		{NoteNumber: 2, Content: "Second note"},
		{NoteNumber: 3, Content: "Third note"},
	}

	result := formatNotes(notes)
	expected := "1. First note\n2. Second note\n3. Third note"

	if result != expected {
		t.Errorf("formatNotes() = %q, expected %q", result, expected)
	}
}

func TestFormatNotes_Empty(t *testing.T) {
	notes := []*models.Note{}
	result := formatNotes(notes)

	if result != "" {
		t.Errorf("formatNotes([]) = %q, expected empty string", result)
	}
}

// Helper function to set test LLM config
func setTestLLMConfig(t *testing.T, srv *Server) {
	t.Helper()

	repo := repositories.NewConfigRepository(srv.database.DB)

	configs := map[string]string{
		"llm_provider_url": "https://api.openai.com/v1",
		"llm_api_key":      "sk-test-key",
		"llm_model":        "gpt-4o",
	}

	for key, value := range configs {
		if err := repo.Set(key, value); err != nil {
			t.Fatalf("failed to set config %s: %v", key, err)
		}
	}
}

// TestLoadLLMConfig tests the config loading helper
func TestLoadLLMConfig_Success(t *testing.T) {
	srv := newTestServer(t)
	repo := repositories.NewConfigRepository(srv.database.DB)

	// Set all required config
	configs := map[string]string{
		"llm_provider_url":   "https://api.openai.com/v1",
		"llm_api_key":        "sk-test-key",
		"llm_model":          "gpt-4o",
		"llm_prompt_summary": "Test summary prompt with {{notes}}",
	}

	for key, value := range configs {
		if err := repo.Set(key, value); err != nil {
			t.Fatalf("failed to set config: %v", err)
		}
	}

	url, apiKey, model, prompt, err := loadLLMConfig(repo)
	if err != nil {
		t.Fatalf("loadLLMConfig() error = %v", err)
	}

	if url != "https://api.openai.com/v1" {
		t.Errorf("expected URL 'https://api.openai.com/v1', got %q", url)
	}
	if apiKey != "sk-test-key" {
		t.Errorf("expected API key 'sk-test-key', got %q", apiKey)
	}
	if model != "gpt-4o" {
		t.Errorf("expected model 'gpt-4o', got %q", model)
	}
	if prompt != "Test summary prompt with {{notes}}" {
		t.Errorf("expected prompt, got %q", prompt)
	}
}

func TestLoadLLMConfig_MissingURL(t *testing.T) {
	srv := newTestServer(t)
	repo := repositories.NewConfigRepository(srv.database.DB)

	// Set only API key, missing URL
	if err := repo.Set("llm_api_key", "sk-test-key"); err != nil {
		t.Fatalf("failed to set API key: %v", err)
	}

	_, _, _, _, err := loadLLMConfig(repo)
	if err == nil {
		t.Error("expected error for missing URL, got nil")
	}
}

func TestLoadLLMConfig_MissingAPIKey(t *testing.T) {
	srv := newTestServer(t)
	repo := repositories.NewConfigRepository(srv.database.DB)

	// Set only URL, missing API key
	if err := repo.Set("llm_provider_url", "https://api.openai.com/v1"); err != nil {
		t.Fatalf("failed to set URL: %v", err)
	}

	_, _, _, _, err := loadLLMConfig(repo)
	if err == nil {
		t.Error("expected error for missing API key, got nil")
	}
}

func TestHandleSummarizeMeeting_InvalidID(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest(http.MethodPost, "/api/meetings/invalid/summarize", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	srv.handleSummarizeMeeting(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleEnhanceNote_InvalidID(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest(http.MethodPost, "/api/notes/invalid/enhance", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	srv.handleEnhanceNote(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestHandleUpdateConfig_AllFields tests that all config fields can be updated including prompts
func TestHandleUpdateConfig_AllFields(t *testing.T) {
	srv := newTestServer(t)

	reqBody := ConfigUpdateRequest{
		LLMProviderURL:   "https://api.anthropic.com/v1",
		LLMAPIKey:        "sk-ant-test",
		LLMModel:         "claude-opus-4-6",
		Language:         "de",
		LLMPromptSummary: "Custom summary prompt",
		LLMPromptEnhance: "Custom enhance prompt",
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

	if resp.LLMPromptSummary != "Custom summary prompt" {
		t.Errorf("expected summary prompt 'Custom summary prompt', got %q", resp.LLMPromptSummary)
	}
	if resp.LLMPromptEnhance != "Custom enhance prompt" {
		t.Errorf("expected enhance prompt 'Custom enhance prompt', got %q", resp.LLMPromptEnhance)
	}
}
