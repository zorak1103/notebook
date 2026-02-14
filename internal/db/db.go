package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"

	_ "modernc.org/sqlite" // SQLite driver for database/sql
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// DB wraps sql.DB for our database operations
type DB struct {
	*sql.DB
}

// Open opens or creates the SQLite database
func Open(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Pragmas for better performance and consistency
	pragmas := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
	}

	ctx := context.Background()
	for _, pragma := range pragmas {
		if _, err := db.ExecContext(ctx, pragma); err != nil {
			return nil, fmt.Errorf("set pragma: %w", err)
		}
	}

	return &DB{db}, nil
}

// Migrate runs all embedded migrations
func (db *DB) Migrate() error {
	ctx := context.Background()

	// Create schema version table
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_version (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("create schema_version table: %w", err)
	}

	// Get current version
	var currentVersion int
	err = db.QueryRowContext(ctx, "SELECT COALESCE(MAX(version), 0) FROM schema_version").Scan(&currentVersion)
	if err != nil {
		return fmt.Errorf("get current version: %w", err)
	}

	// Define migrations
	migrations := []struct {
		version int
		file    string
	}{
		{1, "migrations/001_initial_schema.sql"},
		{2, "migrations/002_add_language_config.sql"},
	}

	// Apply migrations
	for _, m := range migrations {
		if m.version <= currentVersion {
			continue
		}

		log.Printf("Applying migration %d: %s", m.version, m.file)

		migrationSQL, err := migrationsFS.ReadFile(m.file)
		if err != nil {
			return fmt.Errorf("read migration %d: %w", m.version, err)
		}

		if _, err := db.ExecContext(ctx, string(migrationSQL)); err != nil {
			return fmt.Errorf("apply migration %d: %w", m.version, err)
		}

		if _, err := db.ExecContext(ctx, "INSERT INTO schema_version (version) VALUES (?)", m.version); err != nil {
			return fmt.Errorf("record migration %d: %w", m.version, err)
		}
	}

	return nil
}
