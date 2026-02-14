package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zorak1103/notebook/internal/db/models"
	"github.com/zorak1103/notebook/internal/db/repositories"
)

const (
	devModeCreatedBy      = "dev@example.com"
	dateFormat            = "2006-01-02"
	timeFormat            = "15:04"
	maxSubjectLength      = 255
	maxParticipantsLength = 1000
	maxSummaryLength      = 10000
	maxKeywordsLength     = 500
)

// validateMeetingDateTime validates date and time formats
func validateMeetingDateTime(meetingDate, startTime string, endTime *string) error {
	_, err := time.Parse(dateFormat, meetingDate)
	if err != nil {
		return fmt.Errorf("invalid meeting_date format, expected YYYY-MM-DD")
	}

	_, err = time.Parse(timeFormat, startTime)
	if err != nil {
		return fmt.Errorf("invalid start_time format, expected HH:MM")
	}

	if endTime != nil && *endTime != "" {
		_, err = time.Parse(timeFormat, *endTime)
		if err != nil {
			return fmt.Errorf("invalid end_time format, expected HH:MM")
		}
	}

	return nil
}

// validateMeetingFieldLengths validates field length limits
func validateMeetingFieldLengths(m *models.Meeting) error {
	if len(m.Subject) > maxSubjectLength {
		return fmt.Errorf("subject exceeds maximum length of %d characters", maxSubjectLength)
	}

	if m.Participants != nil && len(*m.Participants) > maxParticipantsLength {
		return fmt.Errorf("participants exceeds maximum length of %d characters", maxParticipantsLength)
	}

	if m.Summary != nil && len(*m.Summary) > maxSummaryLength {
		return fmt.Errorf("summary exceeds maximum length of %d characters", maxSummaryLength)
	}

	if m.Keywords != nil && len(*m.Keywords) > maxKeywordsLength {
		return fmt.Errorf("keywords exceeds maximum length of %d characters", maxKeywordsLength)
	}

	return nil
}

// validateMeeting validates both date/time formats and field lengths
func validateMeeting(m *models.Meeting) error {
	if err := validateMeetingDateTime(m.MeetingDate, m.StartTime, m.EndTime); err != nil {
		return err
	}
	return validateMeetingFieldLengths(m)
}

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
		s.logError(r, "failed to list meetings", err)
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
		s.logError(r, "failed to get meeting", err)
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

	// Validate formats and lengths
	if err := validateMeeting(&meeting); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Set created_by (dev mode placeholder for now)
	meeting.CreatedBy = devModeCreatedBy

	repo := repositories.NewMeetingRepository(s.database.DB)
	if err := repo.Create(&meeting); err != nil {
		s.logError(r, "failed to create meeting", err)
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

	// Validate formats and lengths
	err = validateMeeting(&meeting)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if meeting exists
	repo := repositories.NewMeetingRepository(s.database.DB)
	existing, err := repo.GetByID(int(id))
	if err != nil {
		s.logError(r, "failed to check meeting existence", err)
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
		s.logError(r, "failed to update meeting", err)
		writeError(w, http.StatusInternalServerError, "failed to update meeting")
		return
	}

	// Fetch updated meeting to return with all fields
	updated, err := repo.GetByID(int(id))
	if err != nil {
		s.logError(r, "failed to fetch updated meeting", err)
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
		s.logError(r, "failed to check meeting existence", err)
		writeError(w, http.StatusInternalServerError, "failed to check meeting existence")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "meeting not found")
		return
	}

	if err := repo.Delete(int(id)); err != nil {
		s.logError(r, "failed to delete meeting", err)
		writeError(w, http.StatusInternalServerError, "failed to delete meeting")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
