package web

import (
	"fmt"
	"net/http"

	"github.com/zorak1103/notebook/internal/tsapp"
)

// Server manages the HTTP server and routes for the notebook application
type Server struct {
	tsapp   *tsapp.App
	devMode bool
	verbose bool
}

// NewServer creates a new web server instance
func NewServer(app *tsapp.App, devMode, verbose bool) *Server {
	return &Server{
		tsapp:   app,
		devMode: devMode,
		verbose: verbose,
	}
}

// Handler returns the configured HTTP handler with all routes and middleware
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("GET /api/whoami", s.handleWhoAmI)

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
