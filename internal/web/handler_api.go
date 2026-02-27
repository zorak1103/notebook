package web

import (
	"encoding/json"
	"net/http"

	"github.com/zorak1103/notebook/internal/tsapp"
)

// versionResponse holds build-time version information.
type versionResponse struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

// handleVersion returns the build-time version information.
func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	resp := versionResponse{
		Version: s.version,
		Commit:  s.commit,
		Date:    s.date,
	}
	writeJSON(w, http.StatusOK, resp)
	s.logRequest(r.Method, r.URL.Path, http.StatusOK)
}

// handleWhoAmI returns the authenticated user's Tailscale information.
// In dev mode, it returns mock data since Tailscale is not available.
func (s *Server) handleWhoAmI(w http.ResponseWriter, r *http.Request) {
	var userInfo *tsapp.UserInfo
	var err error

	if s.devMode {
		// Mock data for development without Tailscale
		userInfo = &tsapp.UserInfo{
			DisplayName:   "Dev User",
			LoginName:     "dev@example.com",
			ProfilePicURL: "https://ui-avatars.com/api/?name=Dev+User&size=128",
			NodeName:      "dev-machine",
			NodeID:        "dev-node-12345",
		}
	} else {
		// Real Tailscale WhoIs lookup
		if s.tsapp == nil {
			http.Error(w, "Tailscale not initialized", http.StatusInternalServerError)
			return
		}

		userInfo, err = s.tsapp.WhoIs(r)
		if err != nil {
			http.Error(w, "failed to authenticate user: "+err.Error(), http.StatusUnauthorized)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userInfo); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	s.logRequest(r.Method, r.URL.Path, http.StatusOK)
}
