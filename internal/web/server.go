package web

import (
	"fmt"
	"net/http"

	"github.com/zorak1103/notebook/internal/db"
	"github.com/zorak1103/notebook/internal/tsapp"
)

// Server manages the HTTP server and routes for the notebook application
type Server struct {
	tsapp    *tsapp.App
	database *db.DB
	devMode  bool
	verbose  bool
}

// NewServer creates a new web server instance
func NewServer(app *tsapp.App, database *db.DB, devMode, verbose bool) *Server {
	return &Server{
		tsapp:    app,
		database: database,
		devMode:  devMode,
		verbose:  verbose,
	}
}

// Handler returns the configured HTTP handler with all routes and middleware
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("GET /api/whoami", s.handleWhoAmI)

	// Meeting CRUD
	mux.HandleFunc("GET /api/meetings", s.handleListMeetings)
	mux.HandleFunc("POST /api/meetings", s.handleCreateMeeting)
	mux.HandleFunc("GET /api/meetings/{id}", s.handleGetMeeting)
	mux.HandleFunc("PUT /api/meetings/{id}", s.handleUpdateMeeting)
	mux.HandleFunc("DELETE /api/meetings/{id}", s.handleDeleteMeeting)

	// Static files and SPA fallback
	mux.HandleFunc("/", s.handleStatic)

	// Apply middleware
	var handler http.Handler = mux
	handler = s.loggingMiddleware(handler)

	if s.devMode {
		handler = s.corsMiddleware(handler)
	}

	return handler
}

// logRequest logs an HTTP request if verbose mode is enabled
func (s *Server) logRequest(method, path string, status int) {
	if s.verbose {
		fmt.Printf("[%s] %s - %d\n", method, path, status)
	}
}

// logError logs an error with request context
func (s *Server) logError(r *http.Request, msg string, err error) {
	fmt.Printf("[ERROR] %s %s: %s: %v\n", r.Method, r.URL.Path, msg, err)
}
