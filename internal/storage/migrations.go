package storage

import (
	_ "embed"
	"fmt"

	"github.com/jmoiron/sqlx"
)

//go:embed ../../migrations/001_initial_schema.sql
var initialSchema string

// RunMigrations executes all database migrations
func RunMigrations(db *sqlx.DB) error {
	// Execute the initial schema
	_, err := db.Exec(initialSchema)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}
