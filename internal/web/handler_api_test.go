package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleVersion(t *testing.T) {
	server := &Server{
		version: "1.2.3",
		commit:  "abc1234567890",
		date:    "2026-02-27",
	}

	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	w := httptest.NewRecorder()

	server.handleVersion(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp versionResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Version != "1.2.3" {
		t.Errorf("expected version 1.2.3, got %q", resp.Version)
	}
	if resp.Commit != "abc1234567890" {
		t.Errorf("expected commit abc1234567890, got %q", resp.Commit)
	}
	if resp.Date != "2026-02-27" {
		t.Errorf("expected date 2026-02-27, got %q", resp.Date)
	}
}

func TestHandleVersion_Defaults(t *testing.T) {
	server := &Server{
		version: "dev",
		commit:  "none",
		date:    "unknown",
	}

	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	w := httptest.NewRecorder()

	server.handleVersion(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp versionResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Version != "dev" {
		t.Errorf("expected version dev, got %q", resp.Version)
	}
}
