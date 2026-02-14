package web

import (
	"io/fs"
	"net/http"
	"strings"
)

// handleStatic serves the embedded frontend files and implements SPA fallback.
// For unknown paths (not /api/ and not a file), it serves index.html to support
// client-side routing in React.
func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	// Strip the "frontend/dist" prefix from the embedded FS
	distFS, err := fs.Sub(FrontendFS, "frontend/dist")
	if err != nil {
		http.Error(w, "failed to access frontend files", http.StatusInternalServerError)
		return
	}

	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	// Try to serve the requested file
	if file, err := distFS.Open(strings.TrimPrefix(path, "/")); err == nil {
		_ = file.Close()
		// File exists, serve it
		http.FileServer(http.FS(distFS)).ServeHTTP(w, r)
		s.logRequest(r.Method, r.URL.Path, http.StatusOK)
		return
	}

	// File doesn't exist - SPA fallback to index.html
	// This allows React Router to handle the route
	r.URL.Path = "/index.html"
	http.FileServer(http.FS(distFS)).ServeHTTP(w, r)
	s.logRequest(r.Method, path, http.StatusOK)
}
