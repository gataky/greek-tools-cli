package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// ImportCheckpoint represents import progress tracking
type ImportCheckpoint struct {
	ID               int64     `db:"id"`
	CSVFilename      string    `db:"csv_filename"`
	LastProcessedRow int       `db:"last_processed_row"`
	Status           string    `db:"status"`
	UpdatedAt        time.Time `db:"updated_at"`
}

// CreateCheckpoint inserts a new import checkpoint
func (r *SQLiteRepository) CreateCheckpoint(checkpoint *ImportCheckpoint) error {
	query := `
		INSERT INTO import_checkpoints (
			csv_filename, last_processed_row, status
		) VALUES (
			:csv_filename, :last_processed_row, :status
		)
	`
	result, err := r.db.NamedExec(query, checkpoint)
	if err != nil {
		return fmt.Errorf("failed to create checkpoint: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	checkpoint.ID = id
	return nil
}

// UpdateCheckpoint updates an existing checkpoint
func (r *SQLiteRepository) UpdateCheckpoint(checkpoint *ImportCheckpoint) error {
	query := `
		UPDATE import_checkpoints
		SET last_processed_row = :last_processed_row,
		    status = :status,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = :id
	`
	_, err := r.db.NamedExec(query, checkpoint)
	if err != nil {
		return fmt.Errorf("failed to update checkpoint: %w", err)
	}
	return nil
}

// GetCheckpointByFilename retrieves the most recent checkpoint for a CSV file
func (r *SQLiteRepository) GetCheckpointByFilename(filename string) (*ImportCheckpoint, error) {
	var checkpoint ImportCheckpoint
	query := `
		SELECT * FROM import_checkpoints
		WHERE csv_filename = ?
		ORDER BY id DESC
		LIMIT 1
	`
	err := r.db.Get(&checkpoint, query, filename)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No checkpoint found, not an error
		}
		return nil, fmt.Errorf("failed to get checkpoint: %w", err)
	}
	return &checkpoint, nil
}
