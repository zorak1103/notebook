package web

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const contentTypeJSON = "application/json"

// errorResponse represents a JSON error response.
type errorResponse struct {
	Error string `json:"error"`
}

// writeJSON writes a JSON response with the given status code and data.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// writeError writes a JSON error response with the given status code and message.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

// parseIDParam extracts and parses the "id" path parameter.
func parseIDParam(r *http.Request) (int64, error) {
	idStr := r.PathValue("id")
	return strconv.ParseInt(idStr, 10, 64)
}
