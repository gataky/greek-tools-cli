package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gataky/greekmaster/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Repository defines the interface for data storage operations
type Repository interface {
	// Noun operations
	CreateNoun(noun *models.Noun) error
	GetNoun(id int64) (*models.Noun, error)
	ListNouns() ([]*models.Noun, error)

	// Sentence operations
	CreateSentence(sentence *models.Sentence) error
	GetRandomSentences(phase int, number string, limit int) ([]*models.Sentence, error)

	// Close database connection
	Close() error
}

// SQLiteRepository implements Repository using SQLite
type SQLiteRepository struct {
	db *sqlx.DB
}

// NewSQLiteRepository creates a new SQLite repository
// It creates the ~/.greekmaster directory and database if they don't exist
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	// If no path provided, use default
	if dbPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		dbDir := filepath.Join(homeDir, ".greekmaster")

		// Create directory if it doesn't exist
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create data directory: %w", err)
		}

		dbPath = filepath.Join(dbDir, "greekmaster.db")
	}

	// Open database connection
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Run migrations
	if err := RunMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &SQLiteRepository{db: db}, nil
}

// CreateNoun inserts a new noun into the database
func (r *SQLiteRepository) CreateNoun(noun *models.Noun) error {
	query := `
		INSERT INTO nouns (
			english, gender,
			nominative_sg, genitive_sg, accusative_sg,
			nominative_pl, genitive_pl, accusative_pl,
			nom_sg_article, gen_sg_article, acc_sg_article,
			nom_pl_article, gen_pl_article, acc_pl_article
		) VALUES (
			:english, :gender,
			:nominative_sg, :genitive_sg, :accusative_sg,
			:nominative_pl, :genitive_pl, :accusative_pl,
			:nom_sg_article, :gen_sg_article, :acc_sg_article,
			:nom_pl_article, :gen_pl_article, :acc_pl_article
		)
	`
	result, err := r.db.NamedExec(query, noun)
	if err != nil {
		return fmt.Errorf("failed to create noun: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	noun.ID = id
	return nil
}

// GetNoun retrieves a noun by ID
func (r *SQLiteRepository) GetNoun(id int64) (*models.Noun, error) {
	var noun models.Noun
	query := "SELECT * FROM nouns WHERE id = ?"
	err := r.db.Get(&noun, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("noun not found with id %d", id)
		}
		return nil, fmt.Errorf("failed to get noun: %w", err)
	}
	return &noun, nil
}

// ListNouns retrieves all nouns ordered by ID
func (r *SQLiteRepository) ListNouns() ([]*models.Noun, error) {
	var nouns []*models.Noun
	query := "SELECT * FROM nouns ORDER BY id"
	err := r.db.Select(&nouns, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list nouns: %w", err)
	}
	return nouns, nil
}

// CreateSentence inserts a new sentence into the database
func (r *SQLiteRepository) CreateSentence(sentence *models.Sentence) error {
	query := `
		INSERT INTO sentences (
			noun_id, english_prompt, greek_sentence, correct_answer,
			case_type, number, difficulty_phase, context_type, preposition
		) VALUES (
			:noun_id, :english_prompt, :greek_sentence, :correct_answer,
			:case_type, :number, :difficulty_phase, :context_type, :preposition
		)
	`
	result, err := r.db.NamedExec(query, sentence)
	if err != nil {
		return fmt.Errorf("failed to create sentence: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	sentence.ID = id
	return nil
}

// GetRandomSentences retrieves random sentences filtered by phase and number
func (r *SQLiteRepository) GetRandomSentences(phase int, number string, limit int) ([]*models.Sentence, error) {
	var sentences []*models.Sentence
	var query string
	var args []interface{}

	if number == "" {
		// Include both singular and plural
		query = "SELECT * FROM sentences WHERE difficulty_phase = ? ORDER BY RANDOM() LIMIT ?"
		args = []interface{}{phase, limit}
	} else {
		// Filter by specific number
		query = "SELECT * FROM sentences WHERE difficulty_phase = ? AND number = ? ORDER BY RANDOM() LIMIT ?"
		args = []interface{}{phase, number, limit}
	}

	err := r.db.Select(&sentences, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get random sentences: %w", err)
	}
	return sentences, nil
}

// Close closes the database connection
func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}
