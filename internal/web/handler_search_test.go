package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zorak1103/notebook/internal/db/models"
	"github.com/zorak1103/notebook/internal/db/repositories"
)

func TestHandleSearch_EmptyQuery(t *testing.T) {
	s := newTestServer(t)
	defer s.database.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/search?q=", nil)
	w := httptest.NewRecorder()

	s.handleSearch(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var meetings []*models.Meeting
	err := json.NewDecoder(w.Body).Decode(&meetings)
	require.NoError(t, err)
	assert.Empty(t, meetings)
}

func TestHandleSearch_NoResults(t *testing.T) {
	s := newTestServer(t)
	defer s.database.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/search?q=nonexistent", nil)
	w := httptest.NewRecorder()

	s.handleSearch(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var meetings []*models.Meeting
	err := json.NewDecoder(w.Body).Decode(&meetings)
	require.NoError(t, err)
	assert.Empty(t, meetings)
}

func TestHandleSearch_MatchesSubject(t *testing.T) {
	s := newTestServer(t)
	defer s.database.Close()

	// Create test meeting
	repo := repositories.NewMeetingRepository(s.database.DB)
	meeting := &models.Meeting{
		CreatedBy:    "test@example.com",
		Subject:      "Sprint Planning Meeting",
		MeetingDate:  "2026-02-14",
		StartTime:    "10:00",
		EndTime:      stringPtr("11:00"),
		Participants: stringPtr("team@example.com"),
		Summary:      stringPtr("Planning for Q1"),
	}
	err := repo.Create(meeting)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/search?q=Sprint", nil)
	w := httptest.NewRecorder()

	s.handleSearch(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var meetings []*models.Meeting
	err = json.NewDecoder(w.Body).Decode(&meetings)
	require.NoError(t, err)
	require.Len(t, meetings, 1)
	assert.Equal(t, "Sprint Planning Meeting", meetings[0].Subject)
}

func TestHandleSearch_MatchesSummary(t *testing.T) {
	s := newTestServer(t)
	defer s.database.Close()

	// Create test meeting
	repo := repositories.NewMeetingRepository(s.database.DB)
	meeting := &models.Meeting{
		CreatedBy:    "test@example.com",
		Subject:      "Weekly Standup",
		MeetingDate:  "2026-02-14",
		StartTime:    "09:00",
		EndTime:      stringPtr("09:30"),
		Participants: stringPtr("team@example.com"),
		Summary:      stringPtr("Discussed upcoming product launch"),
	}
	err := repo.Create(meeting)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/search?q=launch", nil)
	w := httptest.NewRecorder()

	s.handleSearch(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var meetings []*models.Meeting
	err = json.NewDecoder(w.Body).Decode(&meetings)
	require.NoError(t, err)
	require.Len(t, meetings, 1)
	assert.Equal(t, "Weekly Standup", meetings[0].Subject)
}

func TestHandleSearch_NoQueryParam(t *testing.T) {
	s := newTestServer(t)
	defer s.database.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/search", nil)
	w := httptest.NewRecorder()

	s.handleSearch(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var meetings []*models.Meeting
	err := json.NewDecoder(w.Body).Decode(&meetings)
	require.NoError(t, err)
	assert.Empty(t, meetings)
}

// stringPtr is a helper to create string pointers
func stringPtr(s string) *string {
	return &s
}
