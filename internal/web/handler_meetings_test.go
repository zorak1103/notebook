package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/zorak1103/notebook/internal/db"
	"github.com/zorak1103/notebook/internal/db/models"
	"github.com/zorak1103/notebook/internal/db/repositories"
)

// newTestServer creates a test server with an in-memory database.
func newTestServer(t *testing.T) *Server {
	t.Helper()

	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	if err := database.Migrate(); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return &Server{database: database}
}

func TestHandleListMeetings_Empty(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/meetings", nil)
	w := httptest.NewRecorder()

	server.handleListMeetings(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var meetings []models.Meeting
	if err := json.NewDecoder(w.Body).Decode(&meetings); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if meetings == nil {
		t.Error("expected empty slice, got nil")
	}
	if len(meetings) != 0 {
		t.Errorf("expected 0 meetings, got %d", len(meetings))
	}
}

func TestHandleListMeetings_WithData(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	// Create test meetings
	repo := repositories.NewMeetingRepository(server.database.DB)

	meeting1 := &models.Meeting{
		CreatedBy:   "user1@example.com",
		Subject:     "Meeting 1",
		MeetingDate: time.Now().Format("2006-01-02"),
		StartTime:   "10:00",
	}
	meeting2 := &models.Meeting{
		CreatedBy:   "user2@example.com",
		Subject:     "Meeting 2",
		MeetingDate: time.Now().Format("2006-01-02"),
		StartTime:   "14:00",
	}

	if err := repo.Create(meeting1); err != nil {
		t.Fatalf("failed to create meeting1: %v", err)
	}
	if err := repo.Create(meeting2); err != nil {
		t.Fatalf("failed to create meeting2: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/meetings", nil)
	w := httptest.NewRecorder()

	server.handleListMeetings(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var meetings []models.Meeting
	if err := json.NewDecoder(w.Body).Decode(&meetings); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(meetings) != 2 {
		t.Errorf("expected 2 meetings, got %d", len(meetings))
	}
}

func TestHandleListMeetings_Sorting(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	// Create test meetings with different subjects
	repo := repositories.NewMeetingRepository(server.database.DB)

	meeting1 := &models.Meeting{
		CreatedBy:   "user1@example.com",
		Subject:     "Zebra Meeting",
		MeetingDate: time.Now().Format("2006-01-02"),
		StartTime:   "10:00",
	}
	meeting2 := &models.Meeting{
		CreatedBy:   "user2@example.com",
		Subject:     "Alpha Meeting",
		MeetingDate: time.Now().Format("2006-01-02"),
		StartTime:   "14:00",
	}

	if err := repo.Create(meeting1); err != nil {
		t.Fatalf("failed to create meeting1: %v", err)
	}
	if err := repo.Create(meeting2); err != nil {
		t.Fatalf("failed to create meeting2: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/meetings?sort=subject&order=asc", nil)
	w := httptest.NewRecorder()

	server.handleListMeetings(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var meetings []models.Meeting
	if err := json.NewDecoder(w.Body).Decode(&meetings); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(meetings) != 2 {
		t.Fatalf("expected 2 meetings, got %d", len(meetings))
	}

	if meetings[0].Subject != "Alpha Meeting" {
		t.Errorf("expected first meeting to be 'Alpha Meeting', got '%s'", meetings[0].Subject)
	}
	if meetings[1].Subject != "Zebra Meeting" {
		t.Errorf("expected second meeting to be 'Zebra Meeting', got '%s'", meetings[1].Subject)
	}
}

func TestHandleGetMeeting_Success(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	// Create test meeting
	repo := repositories.NewMeetingRepository(server.database.DB)

	meeting := &models.Meeting{
		CreatedBy:   "user@example.com",
		Subject:     "Test Meeting",
		MeetingDate: time.Now().Format("2006-01-02"),
		StartTime:   "10:00",
	}

	if err := repo.Create(meeting); err != nil {
		t.Fatalf("failed to create meeting: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/meetings/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	server.handleGetMeeting(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var result models.Meeting
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.Subject != "Test Meeting" {
		t.Errorf("expected subject 'Test Meeting', got '%s'", result.Subject)
	}
}

func TestHandleGetMeeting_NotFound(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/meetings/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	server.handleGetMeeting(w, req)

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

func TestHandleGetMeeting_InvalidID(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/meetings/abc", nil)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	server.handleGetMeeting(w, req)

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

func TestHandleCreateMeeting_Success(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	payload := map[string]interface{}{
		"subject":      "New Meeting",
		"meeting_date": time.Now().Format("2006-01-02"),
		"start_time":   "10:00",
		"participants": "Alice, Bob",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/meetings", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.handleCreateMeeting(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var result models.Meeting
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.ID == 0 {
		t.Error("expected meeting to have ID assigned")
	}
	if result.Subject != "New Meeting" {
		t.Errorf("expected subject 'New Meeting', got '%s'", result.Subject)
	}
	if result.CreatedBy == "" {
		t.Error("expected created_by to be set")
	}
}

func TestHandleCreateMeeting_MissingFields(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	payload := map[string]interface{}{
		"meeting_date": time.Now().Format("2006-01-02"),
		// Missing required subject and start_time
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/meetings", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.handleCreateMeeting(w, req)

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

func TestHandleCreateMeeting_InvalidJSON(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/meetings", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	server.handleCreateMeeting(w, req)

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

func TestHandleUpdateMeeting_Success(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	// Create test meeting
	repo := repositories.NewMeetingRepository(server.database.DB)

	meeting := &models.Meeting{
		CreatedBy:   "user@example.com",
		Subject:     "Original Subject",
		MeetingDate: time.Now().Format("2006-01-02"),
		StartTime:   "10:00",
	}

	if err := repo.Create(meeting); err != nil {
		t.Fatalf("failed to create meeting: %v", err)
	}

	// Update meeting
	payload := map[string]interface{}{
		"subject":      "Updated Subject",
		"meeting_date": time.Now().Format("2006-01-02"),
		"start_time":   "11:00",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/api/meetings/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	server.handleUpdateMeeting(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var result models.Meeting
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.Subject != "Updated Subject" {
		t.Errorf("expected subject 'Updated Subject', got '%s'", result.Subject)
	}
	if result.StartTime != "11:00" {
		t.Errorf("expected start_time '11:00', got '%s'", result.StartTime)
	}
}

func TestHandleUpdateMeeting_NotFound(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	payload := map[string]interface{}{
		"subject":      "Updated Subject",
		"meeting_date": time.Now().Format("2006-01-02"),
		"start_time":   "11:00",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/api/meetings/999", bytes.NewReader(body))
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	server.handleUpdateMeeting(w, req)

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

func TestHandleDeleteMeeting_Success(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	// Create test meeting
	repo := repositories.NewMeetingRepository(server.database.DB)

	meeting := &models.Meeting{
		CreatedBy:   "user@example.com",
		Subject:     "Test Meeting",
		MeetingDate: time.Now().Format("2006-01-02"),
		StartTime:   "10:00",
	}

	if err := repo.Create(meeting); err != nil {
		t.Fatalf("failed to create meeting: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/meetings/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	server.handleDeleteMeeting(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}

	// Verify meeting is deleted
	deleted, err := repo.GetByID(meeting.ID)
	if err != nil {
		t.Fatalf("unexpected error checking deleted meeting: %v", err)
	}
	if deleted != nil {
		t.Error("expected meeting to be deleted, but it still exists")
	}
}

func TestHandleDeleteMeeting_NotFound(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	req := httptest.NewRequest(http.MethodDelete, "/api/meetings/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	server.handleDeleteMeeting(w, req)

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

func TestHandleCreateMeeting_InvalidDateFormat(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	payload := map[string]interface{}{
		"subject":      "Test Meeting",
		"meeting_date": "2023-13-45",
		"start_time":   "10:00",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/meetings", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.handleCreateMeeting(w, req)

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

func TestHandleCreateMeeting_InvalidTimeFormat(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	payload := map[string]interface{}{
		"subject":      "Test Meeting",
		"meeting_date": time.Now().Format("2006-01-02"),
		"start_time":   "25:99",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/meetings", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.handleCreateMeeting(w, req)

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

func TestHandleCreateMeeting_InvalidEndTimeFormat(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	endTime := "99:99"
	payload := map[string]interface{}{
		"subject":      "Test Meeting",
		"meeting_date": time.Now().Format("2006-01-02"),
		"start_time":   "10:00",
		"end_time":     endTime,
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/meetings", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.handleCreateMeeting(w, req)

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

func TestHandleUpdateMeeting_InvalidDateFormat(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	// Create test meeting
	repo := repositories.NewMeetingRepository(server.database.DB)

	meeting := &models.Meeting{
		CreatedBy:   "user@example.com",
		Subject:     "Original Subject",
		MeetingDate: time.Now().Format("2006-01-02"),
		StartTime:   "10:00",
	}

	if err := repo.Create(meeting); err != nil {
		t.Fatalf("failed to create meeting: %v", err)
	}

	// Update with invalid date
	payload := map[string]interface{}{
		"subject":      "Updated Subject",
		"meeting_date": "2023-99-99",
		"start_time":   "11:00",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/api/meetings/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	server.handleUpdateMeeting(w, req)

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

func TestHandleCreateMeeting_SubjectTooLong(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	// Create subject with 256 characters (exceeds 255 limit)
	buf := make([]byte, 256)
	for i := range buf {
		buf[i] = byte('a' + (i % 26))
	}
	longSubject := string(buf)

	payload := map[string]interface{}{
		"subject":      longSubject,
		"meeting_date": time.Now().Format("2006-01-02"),
		"start_time":   "10:00",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/meetings", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.handleCreateMeeting(w, req)

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

func TestHandleCreateMeeting_ParticipantsTooLong(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	// Create participants with 1001 characters (exceeds 1000 limit)
	buf := make([]byte, 1001)
	for i := range buf {
		buf[i] = byte('a' + (i % 26))
	}
	longParticipants := string(buf)

	payload := map[string]interface{}{
		"subject":      "Test Meeting",
		"meeting_date": time.Now().Format("2006-01-02"),
		"start_time":   "10:00",
		"participants": longParticipants,
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/meetings", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.handleCreateMeeting(w, req)

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

func TestHandleUpdateMeeting_SubjectTooLong(t *testing.T) {
	server := newTestServer(t)
	defer server.database.Close()

	// Create test meeting
	repo := repositories.NewMeetingRepository(server.database.DB)

	meeting := &models.Meeting{
		CreatedBy:   "user@example.com",
		Subject:     "Original Subject",
		MeetingDate: time.Now().Format("2006-01-02"),
		StartTime:   "10:00",
	}

	if err := repo.Create(meeting); err != nil {
		t.Fatalf("failed to create meeting: %v", err)
	}

	// Create subject with 256 characters (exceeds 255 limit)
	buf := make([]byte, 256)
	for i := range buf {
		buf[i] = byte('a' + (i % 26))
	}
	longSubject := string(buf)

	payload := map[string]interface{}{
		"subject":      longSubject,
		"meeting_date": time.Now().Format("2006-01-02"),
		"start_time":   "11:00",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/api/meetings/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	server.handleUpdateMeeting(w, req)

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
