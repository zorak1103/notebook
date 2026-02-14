package web

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/zorak1103/notebook/internal/db/models"
	"github.com/zorak1103/notebook/internal/db/repositories"
)

// Config key constants
const (
	configKeyLLMProviderURL = "llm_provider_url"
	configKeyLLMAPIKey      = "llm_api_key" // #nosec G101 - config key name, not credential
	configKeyLLMModel       = "llm_model"
)

// ConfigData represents the configuration response
type ConfigData struct {
	LLMProviderURL string `json:"llm_provider_url"`
	LLMAPIKey      string `json:"llm_api_key"`
	LLMModel       string `json:"llm_model"`
}

// ConfigUpdateRequest represents the configuration update request
type ConfigUpdateRequest struct {
	LLMProviderURL string `json:"llm_provider_url"`
	LLMAPIKey      string `json:"llm_api_key"`
	LLMModel       string `json:"llm_model"`
}

// handleGetConfig returns the current configuration with masked API key
func (s *Server) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	repo := repositories.NewConfigRepository(s.database.DB)

	configs, err := repo.GetAll()
	if err != nil {
		s.logError(r, "failed to get configuration", err)
		writeError(w, http.StatusInternalServerError, "failed to get configuration")
		return
	}

	data := buildConfigData(configs)
	writeJSON(w, http.StatusOK, data)
}

// handleUpdateConfig updates the configuration
func (s *Server) handleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	var req ConfigUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate provider URL if provided
	if req.LLMProviderURL != "" {
		if _, err := url.ParseRequestURI(req.LLMProviderURL); err != nil {
			writeError(w, http.StatusBadRequest, "invalid provider URL")
			return
		}
	}

	repo := repositories.NewConfigRepository(s.database.DB)

	// Save non-empty, non-masked fields
	if err := saveConfigFields(repo, &req); err != nil {
		s.logError(r, "failed to save configuration", err)
		writeError(w, http.StatusInternalServerError, "failed to save configuration")
		return
	}

	// Return updated configuration via handleGetConfig
	s.handleGetConfig(w, r)
}

// buildConfigData constructs ConfigData from config models
func buildConfigData(configs []*models.Config) ConfigData {
	data := ConfigData{}

	for _, cfg := range configs {
		switch cfg.Key {
		case configKeyLLMProviderURL:
			data.LLMProviderURL = cfg.Value
		case configKeyLLMAPIKey:
			data.LLMAPIKey = maskAPIKey(cfg.Value)
		case configKeyLLMModel:
			data.LLMModel = cfg.Value
		}
	}

	return data
}

// saveConfigFields saves non-empty, non-masked fields to the repository
func saveConfigFields(repo *repositories.ConfigRepository, req *ConfigUpdateRequest) error {
	if req.LLMProviderURL != "" {
		if err := repo.Set(configKeyLLMProviderURL, req.LLMProviderURL); err != nil {
			return err
		}
	}

	if req.LLMAPIKey != "" && !isMasked(req.LLMAPIKey) {
		if err := repo.Set(configKeyLLMAPIKey, req.LLMAPIKey); err != nil {
			return err
		}
	}

	if req.LLMModel != "" {
		if err := repo.Set(configKeyLLMModel, req.LLMModel); err != nil {
			return err
		}
	}

	return nil
}

// maskAPIKey masks an API key, showing first 4 and last 4 characters
func maskAPIKey(key string) string {
	if key == "" {
		return ""
	}

	if len(key) <= 8 {
		return strings.Repeat("*", len(key))
	}

	first4 := key[:4]
	last4 := key[len(key)-4:]
	masked := strings.Repeat("*", len(key)-8)

	return first4 + masked + last4
}

// isMasked checks if a string contains asterisks (indicating it's masked)
func isMasked(s string) bool {
	return strings.Contains(s, "*")
}
