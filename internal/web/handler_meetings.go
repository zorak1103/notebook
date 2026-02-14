package web

import (
	"encoding/json"
	"net/http"

	"github.com/zorak1103/notebook/internal/db/models"
	"github.com/zorak1103/notebook/internal/db/repositories"
)

const devModeCreatedBy = "dev@example.com"

// handleListMeetings handles GET /api/meetings with optional sorting
func (s *Server) handleListMeetings(w http.ResponseWriter, r *http.Request) {
	sortColumn := r.URL.Query().Get("sort")
	if sortColumn == "" {
		sortColumn = "meeting_date"
	}

	order := r.URL.Query().Get("order")
	ascending := order == "asc"

	repo := repositories.NewMeetingRepository(s.database.DB)
	meetings, err := repo.List(sortColumn, ascending)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list meetings")
		return
	}

	// Coerce nil to empty slice for JSON response
	if meetings == nil {
		meetings = []*models.Meeting{}
	}

	writeJSON(w, http.StatusOK, meetings)
}

// handleGetMeeting handles GET /api/meetings/{id}
func (s *Server) handleGetMeeting(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid meeting ID")
		return
	}

	repo := repositories.NewMeetingRepository(s.database.DB)
	meeting, err := repo.GetByID(int(id))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get meeting")
		return
	}

	if meeting == nil {
		writeError(w, http.StatusNotFound, "meeting not found")
		return
	}

	writeJSON(w, http.StatusOK, meeting)
}

// handleCreateMeeting handles POST /api/meetings
func (s *Server) handleCreateMeeting(w http.ResponseWriter, r *http.Request) {
	var meeting models.Meeting
	if err := json.NewDecoder(r.Body).Decode(&meeting); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if meeting.Subject == "" || meeting.MeetingDate == "" || meeting.StartTime == "" {
		writeError(w, http.StatusBadRequest, "missing required fields: subject, meeting_date, start_time")
		return
	}

	// Set created_by (dev mode placeholder for now)
	meeting.CreatedBy = devModeCreatedBy

	repo := repositories.NewMeetingRepository(s.database.DB)
	if err := repo.Create(&meeting); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create meeting")
		return
	}

	writeJSON(w, http.StatusCreated, meeting)
}

// handleUpdateMeeting handles PUT /api/meetings/{id}
func (s *Server) handleUpdateMeeting(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid meeting ID")
		return
	}

	var meeting models.Meeting
	err = json.NewDecoder(r.Body).Decode(&meeting)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if meeting.Subject == "" || meeting.MeetingDate == "" || meeting.StartTime == "" {
		writeError(w, http.StatusBadRequest, "missing required fields: subject, meeting_date, start_time")
		return
	}

	// Check if meeting exists
	repo := repositories.NewMeetingRepository(s.database.DB)
	existing, err := repo.GetByID(int(id))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to check meeting existence")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "meeting not found")
		return
	}

	// Set ID from path parameter
	meeting.ID = int(id)

	err = repo.Update(&meeting)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update meeting")
		return
	}

	// Fetch updated meeting to return with all fields
	updated, err := repo.GetByID(int(id))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch updated meeting")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// handleDeleteMeeting handles DELETE /api/meetings/{id}
func (s *Server) handleDeleteMeeting(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid meeting ID")
		return
	}

	// Check if meeting exists
	repo := repositories.NewMeetingRepository(s.database.DB)
	existing, err := repo.GetByID(int(id))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to check meeting existence")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "meeting not found")
		return
	}

	if err := repo.Delete(int(id)); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete meeting")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
