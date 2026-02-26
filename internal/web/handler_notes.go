package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zorak1103/notebook/internal/db/models"
	"github.com/zorak1103/notebook/internal/db/repositories"
	"github.com/zorak1103/notebook/internal/validation"
)

const (
	errInvalidNoteID    = "invalid note ID"
	errNoteNotFound     = "note not found"
	errInvalidMeetingID = "invalid meeting ID"
	errInvalidDirection = "invalid direction: must be 'up' or 'down'"
)

// parseMeetingIDParam extracts and parses the "meetingId" path parameter.
func parseMeetingIDParam(r *http.Request) (int64, error) {
	idStr := r.PathValue("meetingId")
	return strconv.ParseInt(idStr, 10, 64)
}

// validateNoteContent validates note content is not empty and within length limit.
func validateNoteContent(content string) error {
	if content == "" {
		return fmt.Errorf("content is required")
	}

	if len(content) > validation.MaxNoteContentLength {
		return fmt.Errorf("content exceeds maximum length of %d characters", validation.MaxNoteContentLength)
	}

	return nil
}

// reorderNoteRequest is the request body for reordering a note
type reorderNoteRequest struct {
	Direction string `json:"direction"`
}

// findAdjacentNote returns the note adjacent to noteID in the given direction within the sorted list.
func findAdjacentNote(notes []*models.Note, noteID int, direction string) (*models.Note, error) {
	for i, n := range notes {
		if n.ID != noteID {
			continue
		}
		switch direction {
		case "up":
			if i == 0 {
				return nil, fmt.Errorf("note is already first")
			}
			return notes[i-1], nil
		case "down":
			if i == len(notes)-1 {
				return nil, fmt.Errorf("note is already last")
			}
			return notes[i+1], nil
		default:
			return nil, errors.New(errInvalidDirection)
		}
	}
	return nil, fmt.Errorf("note not found in list")
}

// handleReorderNote handles PUT /api/notes/{id}/reorder
func (s *Server) handleReorderNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, errInvalidNoteID)
		return
	}

	var req reorderNoteRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Direction != "up" && req.Direction != "down" {
		writeError(w, http.StatusBadRequest, errInvalidDirection)
		return
	}

	repo := repositories.NewNoteRepository(s.database.DB)
	note, err := repo.GetByID(int(id))
	if err != nil {
		s.logError(r, "failed to get note", err)
		writeError(w, http.StatusInternalServerError, "failed to get note")
		return
	}
	if note == nil {
		writeError(w, http.StatusNotFound, errNoteNotFound)
		return
	}

	noteList, err := repo.ListByMeeting(note.MeetingID)
	if err != nil {
		s.logError(r, "failed to list notes", err)
		writeError(w, http.StatusInternalServerError, "failed to list notes")
		return
	}

	adjacent, err := findAdjacentNote(noteList, int(id), req.Direction)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = repo.SwapNoteOrder(int(id), adjacent.ID); err != nil {
		s.logError(r, "failed to swap note order", err)
		writeError(w, http.StatusInternalServerError, "failed to reorder note")
		return
	}

	updated, err := repo.ListByMeeting(note.MeetingID)
	if err != nil {
		s.logError(r, "failed to list updated notes", err)
		writeError(w, http.StatusInternalServerError, "failed to list updated notes")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// handleListNotes handles GET /api/meetings/{meetingId}/notes
func (s *Server) handleListNotes(w http.ResponseWriter, r *http.Request) {
	meetingID, err := parseMeetingIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, errInvalidMeetingID)
		return
	}

	repo := repositories.NewNoteRepository(s.database.DB)
	notes, err := repo.ListByMeeting(int(meetingID))
	if err != nil {
		s.logError(r, "failed to list notes", err)
		writeError(w, http.StatusInternalServerError, "failed to list notes")
		return
	}

	// Coerce nil to empty slice for JSON response
	if notes == nil {
		notes = []*models.Note{}
	}

	writeJSON(w, http.StatusOK, notes)
}

// handleGetNote handles GET /api/notes/{id}
func (s *Server) handleGetNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, errInvalidNoteID)
		return
	}

	repo := repositories.NewNoteRepository(s.database.DB)
	note, err := repo.GetByID(int(id))
	if err != nil {
		s.logError(r, "failed to get note", err)
		writeError(w, http.StatusInternalServerError, "failed to get note")
		return
	}

	if note == nil {
		writeError(w, http.StatusNotFound, errNoteNotFound)
		return
	}

	writeJSON(w, http.StatusOK, note)
}

// handleCreateNote handles POST /api/notes
func (s *Server) handleCreateNote(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if note.MeetingID == 0 {
		writeError(w, http.StatusBadRequest, "missing required field: meeting_id")
		return
	}

	// Validate content
	if err := validateNoteContent(note.Content); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	repo := repositories.NewNoteRepository(s.database.DB)
	if err := repo.Create(&note); err != nil {
		s.logError(r, "failed to create note", err)
		writeError(w, http.StatusInternalServerError, "failed to create note")
		return
	}

	writeJSON(w, http.StatusCreated, note)
}

// handleUpdateNote handles PUT /api/notes/{id}
func (s *Server) handleUpdateNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, errInvalidNoteID)
		return
	}

	var note models.Note
	err = json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate content
	err = validateNoteContent(note.Content)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if note exists
	repo := repositories.NewNoteRepository(s.database.DB)
	existing, err := repo.GetByID(int(id))
	if err != nil {
		s.logError(r, "failed to check note existence", err)
		writeError(w, http.StatusInternalServerError, "failed to check note existence")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, errNoteNotFound)
		return
	}

	// Set ID from path parameter
	note.ID = int(id)

	err = repo.Update(&note)
	if err != nil {
		s.logError(r, "failed to update note", err)
		writeError(w, http.StatusInternalServerError, "failed to update note")
		return
	}

	// Fetch updated note to return with all fields
	updated, err := repo.GetByID(int(id))
	if err != nil {
		s.logError(r, "failed to fetch updated note", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch updated note")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// handleDeleteNote handles DELETE /api/notes/{id}
func (s *Server) handleDeleteNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, errInvalidNoteID)
		return
	}

	// Check if note exists
	repo := repositories.NewNoteRepository(s.database.DB)
	existing, err := repo.GetByID(int(id))
	if err != nil {
		s.logError(r, "failed to check note existence", err)
		writeError(w, http.StatusInternalServerError, "failed to check note existence")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, errNoteNotFound)
		return
	}

	if err := repo.Delete(int(id)); err != nil {
		s.logError(r, "failed to delete note", err)
		writeError(w, http.StatusInternalServerError, "failed to delete note")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
