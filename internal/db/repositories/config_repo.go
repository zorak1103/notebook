package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zorak1103/notebook/internal/db/models"
)

// ConfigRepository handles config CRUD operations
type ConfigRepository struct {
	db *sql.DB
}

// NewConfigRepository creates a new config repository
func NewConfigRepository(db *sql.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

// Get retrieves a config value by key
func (r *ConfigRepository) Get(key string) (*models.Config, error) {
	ctx := context.Background()
	c := &models.Config{}
	err := r.db.QueryRowContext(ctx, `
		SELECT key, value, updated_at FROM config WHERE key = ?
	`, key).Scan(&c.Key, &c.Value, &c.UpdatedAt)

	if err == sql.ErrNoRows {
		//nolint:nilnil // Intentional: not found is not an error
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get config: %w", err)
	}

	return c, nil
}

// GetAll retrieves all config entries
func (r *ConfigRepository) GetAll() ([]*models.Config, error) {
	ctx := context.Background()
	rows, err := r.db.QueryContext(ctx, `SELECT key, value, updated_at FROM config`)
	if err != nil {
		return nil, fmt.Errorf("get all config: %w", err)
	}
	defer rows.Close()

	var configs []*models.Config
	for rows.Next() {
		c := &models.Config{}
		if err := rows.Scan(&c.Key, &c.Value, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan config: %w", err)
		}
		configs = append(configs, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return configs, nil
}

// Set sets a config value (upsert)
func (r *ConfigRepository) Set(key, value string) error {
	ctx := context.Background()
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO config (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`, key, value)

	if err != nil {
		return fmt.Errorf("set config: %w", err)
	}

	return nil
}
