package storage

import (
	_ "embed"
	"fmt"

	"github.com/jmoiron/sqlx"
)

//go:embed migrations/001_initial_schema.sql
var initialSchema string

//go:embed migrations/002_drop_explanations.sql
var dropExplanations string

//go:embed migrations/003_create_templates.sql
var createTemplates string

// RunMigrations executes all database migrations
func RunMigrations(db *sqlx.DB) error {
	// Execute the initial schema
	_, err := db.Exec(initialSchema)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Drop explanations table (no longer needed)
	_, err = db.Exec(dropExplanations)
	if err != nil {
		return fmt.Errorf("failed to run migration 002: %w", err)
	}

	// Create sentence_templates table
	_, err = db.Exec(createTemplates)
	if err != nil {
		return fmt.Errorf("failed to run migration 003: %w", err)
	}

	return nil
}
