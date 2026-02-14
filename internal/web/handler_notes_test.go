package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/zorak1103/notebook/internal/db/models"
	"github.com/zorak1103/notebook/internal/db/repositories"
	"github.com/zorak1103/notebook/internal/validation"
)

// createTestMeeting creates a test meeting and returns its ID
func createTestMeeting(t *testing.T, repo *repositories.MeetingRepository) int {
	t.Helper()

	meeting := &models.Meeting{
		CreatedBy:   "user@example.com",
		Subject:     "Test Meeting",
		MeetingDate: time.Now().Format("2006-01-02"),
		StartTime:   "10:00",
	}

	if err := repo.Create(meeting); err != nil {
		t.Fatalf("failed to create test meeting: %v", err)
	}

	return meeting.ID
}

func TestHandleListNotes_Empty(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	// Create meeting (required FK)
	meetingRepo := repositories.NewMeetingRepository(server.database.DB)
	meetingID := createTestMeeting(t, meetingRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/meetings/1/notes", nil)
	req.SetPathValue("meetingId", "1")
	w := httptest.NewRecorder()

	server.handleListNotes(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var notes []models.Note
	if err := json.NewDecoder(w.Body).Decode(&notes); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if notes == nil {
		t.Error("expected empty slice, got nil")
	}
	if len(notes) != 0 {
		t.Errorf("expected 0 notes, got %d (meetingID=%d)", len(notes), meetingID)
	}
}

func TestHandleListNotes_WithData(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	// Create meeting
	meetingRepo := repositories.NewMeetingRepository(server.database.DB)
	meetingID := createTestMeeting(t, meetingRepo)

	// Create notes
	noteRepo := repositories.NewNoteRepository(server.database.DB)

	note1 := &models.Note{
		MeetingID: meetingID,
		Content:   "First note",
	}
	note2 := &models.Note{
		MeetingID: meetingID,
		Content:   "Second note",
	}

	if err := noteRepo.Create(note1); err != nil {
		t.Fatalf("failed to create note1: %v", err)
	}
	if err := noteRepo.Create(note2); err != nil {
		t.Fatalf("failed to create note2: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/meetings/1/notes", nil)
	req.SetPathValue("meetingId", "1")
	w := httptest.NewRecorder()

	server.handleListNotes(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var notes []models.Note
	if err := json.NewDecoder(w.Body).Decode(&notes); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(notes) != 2 {
		t.Errorf("expected 2 notes, got %d", len(notes))
	}

	// Verify ordering by note_number
	if len(notes) == 2 {
		if notes[0].NoteNumber != 1 {
			t.Errorf("expected first note to have note_number 1, got %d", notes[0].NoteNumber)
		}
		if notes[1].NoteNumber != 2 {
			t.Errorf("expected second note to have note_number 2, got %d", notes[1].NoteNumber)
		}
	}
}

func TestHandleGetNote_Success(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	// Create meeting and note
	meetingRepo := repositories.NewMeetingRepository(server.database.DB)
	meetingID := createTestMeeting(t, meetingRepo)

	noteRepo := repositories.NewNoteRepository(server.database.DB)
	note := &models.Note{
		MeetingID: meetingID,
		Content:   "Test note content",
	}

	if err := noteRepo.Create(note); err != nil {
		t.Fatalf("failed to create note: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/notes/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	server.handleGetNote(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var result models.Note
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.Content != "Test note content" {
		t.Errorf("expected content 'Test note content', got '%s'", result.Content)
	}
}

func TestHandleGetNote_NotFound(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/notes/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	server.handleGetNote(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var errResp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Error == "" {
		t.Error("expected error message, got empty string")
	}
}

func TestHandleGetNote_InvalidID(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/notes/abc", nil)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	server.handleGetNote(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var errResp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Error == "" {
		t.Error("expected error message, got empty string")
	}
}

func TestHandleCreateNote_Success(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	// Create meeting
	meetingRepo := repositories.NewMeetingRepository(server.database.DB)
	meetingID := createTestMeeting(t, meetingRepo)

	payload := map[string]interface{}{
		"meeting_id": meetingID,
		"content":    "New note content",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/notes", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.handleCreateNote(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var result models.Note
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.ID == 0 {
		t.Error("expected note to have ID assigned")
	}
	if result.NoteNumber != 1 {
		t.Errorf("expected note_number 1, got %d", result.NoteNumber)
	}
	if result.Content != "New note content" {
		t.Errorf("expected content 'New note content', got '%s'", result.Content)
	}
}

func TestHandleCreateNote_AutoNumbering(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	// Create meeting
	meetingRepo := repositories.NewMeetingRepository(server.database.DB)
	meetingID := createTestMeeting(t, meetingRepo)

	// Create first note
	payload1 := map[string]interface{}{
		"meeting_id": meetingID,
		"content":    "First note",
	}

	body1, _ := json.Marshal(payload1)
	req1 := httptest.NewRequest(http.MethodPost, "/api/notes", bytes.NewReader(body1))
	w1 := httptest.NewRecorder()

	server.handleCreateNote(w1, req1)

	var note1 models.Note
	if err := json.NewDecoder(w1.Body).Decode(&note1); err != nil {
		t.Fatalf("failed to decode first note response: %v", err)
	}

	// Create second note
	payload2 := map[string]interface{}{
		"meeting_id": meetingID,
		"content":    "Second note",
	}

	body2, _ := json.Marshal(payload2)
	req2 := httptest.NewRequest(http.MethodPost, "/api/notes", bytes.NewReader(body2))
	w2 := httptest.NewRecorder()

	server.handleCreateNote(w2, req2)

	var note2 models.Note
	if err := json.NewDecoder(w2.Body).Decode(&note2); err != nil {
		t.Fatalf("failed to decode second note response: %v", err)
	}

	// Verify auto-numbering
	if note1.NoteNumber != 1 {
		t.Errorf("expected first note to have note_number 1, got %d", note1.NoteNumber)
	}
	if note2.NoteNumber != 2 {
		t.Errorf("expected second note to have note_number 2, got %d", note2.NoteNumber)
	}
}

func TestHandleCreateNote_MissingFields(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	payload := map[string]interface{}{
		// Missing required meeting_id and content
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/notes", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.handleCreateNote(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var errResp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Error == "" {
		t.Error("expected error message, got empty string")
	}
}

func TestHandleCreateNote_InvalidJSON(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/notes", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	server.handleCreateNote(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var errResp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Error == "" {
		t.Error("expected error message, got empty string")
	}
}

func TestHandleCreateNote_ContentTooLong(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	// Create meeting
	meetingRepo := repositories.NewMeetingRepository(server.database.DB)
	meetingID := createTestMeeting(t, meetingRepo)

	// Create content exceeding max length
	buf := make([]byte, validation.MaxNoteContentLength+1)
	for i := range buf {
		buf[i] = byte('a' + (i % 26))
	}
	longContent := string(buf)

	payload := map[string]interface{}{
		"meeting_id": meetingID,
		"content":    longContent,
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/notes", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.handleCreateNote(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var errResp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Error == "" {
		t.Error("expected error message, got empty string")
	}
}

func TestHandleUpdateNote_Success(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	// Create meeting and note
	meetingRepo := repositories.NewMeetingRepository(server.database.DB)
	meetingID := createTestMeeting(t, meetingRepo)

	noteRepo := repositories.NewNoteRepository(server.database.DB)
	note := &models.Note{
		MeetingID: meetingID,
		Content:   "Original content",
	}

	if err := noteRepo.Create(note); err != nil {
		t.Fatalf("failed to create note: %v", err)
	}

	// Update note
	payload := map[string]interface{}{
		"content": "Updated content",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/api/notes/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	server.handleUpdateNote(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var result models.Note
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.Content != "Updated content" {
		t.Errorf("expected content 'Updated content', got '%s'", result.Content)
	}
}

func TestHandleUpdateNote_NotFound(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	payload := map[string]interface{}{
		"content": "Updated content",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/api/notes/999", bytes.NewReader(body))
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	server.handleUpdateNote(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var errResp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Error == "" {
		t.Error("expected error message, got empty string")
	}
}

func TestHandleUpdateNote_MissingContent(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	// Create meeting and note
	meetingRepo := repositories.NewMeetingRepository(server.database.DB)
	meetingID := createTestMeeting(t, meetingRepo)

	noteRepo := repositories.NewNoteRepository(server.database.DB)
	note := &models.Note{
		MeetingID: meetingID,
		Content:   "Original content",
	}

	if err := noteRepo.Create(note); err != nil {
		t.Fatalf("failed to create note: %v", err)
	}

	// Update with empty content
	payload := map[string]interface{}{
		"content": "",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/api/notes/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	server.handleUpdateNote(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var errResp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Error == "" {
		t.Error("expected error message, got empty string")
	}
}

func TestHandleDeleteNote_Success(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	// Create meeting and note
	meetingRepo := repositories.NewMeetingRepository(server.database.DB)
	meetingID := createTestMeeting(t, meetingRepo)

	noteRepo := repositories.NewNoteRepository(server.database.DB)
	note := &models.Note{
		MeetingID: meetingID,
		Content:   "Test note",
	}

	if err := noteRepo.Create(note); err != nil {
		t.Fatalf("failed to create note: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/notes/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	server.handleDeleteNote(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}

	// Verify note is deleted
	deleted, err := noteRepo.GetByID(note.ID)
	if err != nil {
		t.Fatalf("unexpected error checking deleted note: %v", err)
	}
	if deleted != nil {
		t.Error("expected note to be deleted, but it still exists")
	}
}

func TestHandleDeleteNote_NotFound(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	req := httptest.NewRequest(http.MethodDelete, "/api/notes/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	server.handleDeleteNote(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var errResp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Error == "" {
		t.Error("expected error message, got empty string")
	}
}
