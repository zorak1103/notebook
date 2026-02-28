package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/zorak1103/notebook/internal/db/models"
	"github.com/zorak1103/notebook/internal/db/repositories"
	"github.com/zorak1103/notebook/internal/llm"
)

type enhanceNoteRequest struct {
	Content string `json:"content"`
}

type enhanceNoteResponse struct {
	Content string `json:"content"`
}

// handleSummarizeMeeting generates an LLM summary for a meeting based on its notes
func (s *Server) handleSummarizeMeeting(w http.ResponseWriter, r *http.Request) {
	meetingID, err := parseIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid meeting ID")
		return
	}

	// Load LLM config
	configRepo := repositories.NewConfigRepository(s.database.DB)
	llmURL, llmAPIKey, llmModel, summaryPrompt, err := loadLLMConfig(configRepo)
	if err != nil {
		s.logError(r, "failed to load LLM config", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Load meeting and notes
	meeting, notes, meetingRepo, err := s.loadMeetingWithNotes(int(meetingID))
	if err != nil {
		s.logError(r, "failed to load meeting data", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if meeting == nil {
		writeError(w, http.StatusNotFound, "meeting not found")
		return
	}
	if len(notes) == 0 {
		writeError(w, http.StatusBadRequest, "no notes to summarize")
		return
	}

	// Generate summary using LLM
	summary, err := s.generateSummary(r, llmURL, llmAPIKey, llmModel, summaryPrompt, meeting, notes)
	if err != nil {
		s.logError(r, "failed to generate summary", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Update meeting with summary
	meeting.Summary = &summary
	if err := meetingRepo.Update(meeting); err != nil {
		s.logError(r, "failed to update meeting", err)
		writeError(w, http.StatusInternalServerError, "failed to update meeting")
		return
	}

	writeJSON(w, http.StatusOK, meeting)
}

// handleEnhanceNote transforms note content via LLM and returns the result.
// It does not persist to DB — the caller decides whether to save.
func (s *Server) handleEnhanceNote(w http.ResponseWriter, r *http.Request) {
	noteID, err := parseIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid note ID")
		return
	}

	var req enhanceNoteRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(req.Content) == "" {
		writeError(w, http.StatusBadRequest, "content is required")
		return
	}

	configRepo := repositories.NewConfigRepository(s.database.DB)
	llmURL, llmAPIKey, llmModel, enhancePrompt, err := loadLLMConfigForEnhance(configRepo)
	if err != nil {
		s.logError(r, "failed to load LLM config", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	noteRepo := repositories.NewNoteRepository(s.database.DB)
	note, err := noteRepo.GetByID(int(noteID))
	if err != nil {
		s.logError(r, "failed to get note", err)
		writeError(w, http.StatusInternalServerError, "failed to get note")
		return
	}
	if note == nil {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}

	enhanced, err := s.enhanceContent(r, llmURL, llmAPIKey, llmModel, enhancePrompt, req.Content)
	if err != nil {
		s.logError(r, "LLM completion failed", err)
		writeError(w, http.StatusInternalServerError, "LLM completion failed")
		return
	}

	writeJSON(w, http.StatusOK, enhanceNoteResponse{Content: enhanced})
}

// enhanceContent uses LLM to transform a piece of text
func (s *Server) enhanceContent(r *http.Request, llmURL, llmAPIKey, llmModel, enhancePrompt, content string) (string, error) {
	client, err := llm.New(llmURL, llmAPIKey, llmModel)
	if err != nil {
		return "", fmt.Errorf("failed to create LLM client: %w", err)
	}

	prompt := llm.RenderPrompt(enhancePrompt, map[string]string{
		"content": content,
	})

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	enhanced, err := client.Complete(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("LLM completion failed: %w", err)
	}

	return enhanced, nil
}

// loadMeetingWithNotes loads a meeting and its notes from the database
func (s *Server) loadMeetingWithNotes(meetingID int) (*models.Meeting, []*models.Note, *repositories.MeetingRepository, error) {
	meetingRepo := repositories.NewMeetingRepository(s.database.DB)
	meeting, err := meetingRepo.GetByID(meetingID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get meeting: %w", err)
	}

	noteRepo := repositories.NewNoteRepository(s.database.DB)
	notes, err := noteRepo.ListByMeeting(meetingID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get notes: %w", err)
	}

	return meeting, notes, meetingRepo, nil
}

// generateSummary creates an LLM summary from meeting and notes
func (s *Server) generateSummary(r *http.Request, llmURL, llmAPIKey, llmModel, summaryPrompt string, meeting *models.Meeting, notes []*models.Note) (string, error) {
	client, err := llm.New(llmURL, llmAPIKey, llmModel)
	if err != nil {
		return "", fmt.Errorf("failed to create LLM client: %w", err)
	}

	notesText := formatNotes(notes)
	participants := ""
	if meeting.Participants != nil {
		participants = *meeting.Participants
	}
	prompt := llm.RenderPrompt(summaryPrompt, map[string]string{
		"subject":      meeting.Subject,
		"date":         meeting.MeetingDate,
		"participants": participants,
		"notes":        notesText,
	})

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	summary, err := client.Complete(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("LLM completion failed: %w", err)
	}

	return summary, nil
}

// loadLLMConfig loads and validates LLM configuration for summarization
func loadLLMConfig(repo *repositories.ConfigRepository) (url, apiKey, model, summaryPrompt string, err error) {
	configs, err := repo.GetAll()
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to get config: %w", err)
	}

	for _, cfg := range configs {
		switch cfg.Key {
		case configKeyLLMProviderURL:
			url = cfg.Value
		case configKeyLLMAPIKey:
			apiKey = cfg.Value
		case configKeyLLMModel:
			model = cfg.Value
		case configKeyLLMPromptSummary:
			summaryPrompt = cfg.Value
		}
	}

	if url == "" || apiKey == "" {
		return "", "", "", "", fmt.Errorf("LLM provider not configured")
	}

	return url, apiKey, model, summaryPrompt, nil
}

// loadLLMConfigForEnhance loads and validates LLM configuration for note enhancement
func loadLLMConfigForEnhance(repo *repositories.ConfigRepository) (url, apiKey, model, enhancePrompt string, err error) {
	configs, err := repo.GetAll()
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to get config: %w", err)
	}

	for _, cfg := range configs {
		switch cfg.Key {
		case configKeyLLMProviderURL:
			url = cfg.Value
		case configKeyLLMAPIKey:
			apiKey = cfg.Value
		case configKeyLLMModel:
			model = cfg.Value
		case configKeyLLMPromptEnhance:
			enhancePrompt = cfg.Value
		}
	}

	if url == "" || apiKey == "" {
		return "", "", "", "", fmt.Errorf("LLM provider not configured")
	}

	return url, apiKey, model, enhancePrompt, nil
}

// formatNotes converts a list of notes to a formatted string
func formatNotes(notes []*models.Note) string {
	var parts []string
	for _, note := range notes {
		parts = append(parts, fmt.Sprintf("%d. %s", note.NoteNumber, note.Content))
	}
	return strings.Join(parts, "\n")
}
