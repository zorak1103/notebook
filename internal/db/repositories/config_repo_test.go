package repositories_test

import (
	"testing"

	"github.com/zorak1103/notebook/internal/db/repositories"
)

func TestConfigRepository_Get(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repositories.NewConfigRepository(database.DB)

	// Migration seeds llm_provider_url with empty string
	cfg, err := repo.Get("llm_provider_url")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("expected config to be found")
	}

	if cfg.Key != "llm_provider_url" {
		t.Errorf("expected key %q, got %q", "llm_provider_url", cfg.Key)
	}
}

func TestConfigRepository_Get_NotFound(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repositories.NewConfigRepository(database.DB)

	cfg, err := repo.Get("nonexistent_key")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if cfg != nil {
		t.Error("expected config to be nil for non-existent key")
	}
}

func TestConfigRepository_GetAll(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repositories.NewConfigRepository(database.DB)

	configs, err := repo.GetAll()
	if err != nil {
		t.Fatalf("getAll failed: %v", err)
	}

	// Migration seeds 6 config entries (llm_provider_url, llm_api_key, llm_model, language, llm_prompt_summary, llm_prompt_enhance)
	if len(configs) != 6 {
		t.Errorf("expected 6 configs, got %d", len(configs))
	}
}

func TestConfigRepository_Set_Update(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repositories.NewConfigRepository(database.DB)

	// Update existing key
	if err := repo.Set("llm_provider_url", "https://api.example.com"); err != nil {
		t.Fatalf("set failed: %v", err)
	}

	// Verify
	cfg, err := repo.Get("llm_provider_url")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if cfg.Value != "https://api.example.com" {
		t.Errorf("expected value %q, got %q", "https://api.example.com", cfg.Value)
	}
}

func TestConfigRepository_Set_Insert(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repositories.NewConfigRepository(database.DB)

	// Insert new key
	if err := repo.Set("custom_key", "custom_value"); err != nil {
		t.Fatalf("set failed: %v", err)
	}

	// Verify
	cfg, err := repo.Get("custom_key")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("expected config to be found")
	}

	if cfg.Value != "custom_value" {
		t.Errorf("expected value %q, got %q", "custom_value", cfg.Value)
	}
}
