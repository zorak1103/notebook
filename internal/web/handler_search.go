package web

import (
	"net/http"

	"github.com/zorak1103/notebook/internal/db/models"
	"github.com/zorak1103/notebook/internal/db/repositories"
)

// handleSearch handles GET /api/search?q=<query>
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	// Empty or missing query returns empty array
	if query == "" {
		writeJSON(w, http.StatusOK, []*models.Meeting{})
		return
	}

	repo := repositories.NewMeetingRepository(s.database.DB)
	meetings, err := repo.Search(query)
	if err != nil {
		s.logError(r, "failed to search meetings", err)
		writeError(w, http.StatusInternalServerError, "failed to search meetings")
		return
	}

	// Coerce nil slice to empty slice for JSON response
	if meetings == nil {
		meetings = []*models.Meeting{}
	}

	writeJSON(w, http.StatusOK, meetings)
}
