package storage

import (
	"database/sql"
	"fmt"

	"github.com/gataky/greekmaster/internal/models"
)

// CreateTemplate inserts a new template into the database
func (r *SQLiteRepository) CreateTemplate(template *models.SentenceTemplate) error {
	query := `
		INSERT INTO sentence_templates (
			english_template, greek_template, article_field, noun_form_field,
			case_type, number, difficulty_phase, context_type, preposition
		) VALUES (
			:english_template, :greek_template, :article_field, :noun_form_field,
			:case_type, :number, :difficulty_phase, :context_type, :preposition
		)
	`
	result, err := r.db.NamedExec(query, template)
	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	template.ID = id
	return nil
}

// GetTemplate retrieves a template by ID
func (r *SQLiteRepository) GetTemplate(id int64) (*models.SentenceTemplate, error) {
	var template models.SentenceTemplate
	query := "SELECT * FROM sentence_templates WHERE id = ?"
	err := r.db.Get(&template, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("template not found with id %d", id)
		}
		return nil, fmt.Errorf("failed to get template: %w", err)
	}
	return &template, nil
}

// ListTemplates retrieves all templates ordered by ID
func (r *SQLiteRepository) ListTemplates() ([]*models.SentenceTemplate, error) {
	var templates []*models.SentenceTemplate
	query := "SELECT * FROM sentence_templates ORDER BY id"
	err := r.db.Select(&templates, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}
	return templates, nil
}

// GetRandomTemplates retrieves random templates filtered by phase and number
func (r *SQLiteRepository) GetRandomTemplates(phase int, number string, limit int) ([]*models.SentenceTemplate, error) {
	var templates []*models.SentenceTemplate
	var query string
	var args []interface{}

	// Build query based on number filter
	if number == "" || number == "both" {
		// Include templates for singular, plural, and both
		query = "SELECT * FROM sentence_templates WHERE difficulty_phase = ? AND (number = 'singular' OR number = 'plural' OR number = 'both') ORDER BY RANDOM() LIMIT ?"
		args = []interface{}{phase, limit}
	} else {
		// Filter by specific number (singular or plural), but also include templates marked as 'both'
		query = "SELECT * FROM sentence_templates WHERE difficulty_phase = ? AND (number = ? OR number = 'both') ORDER BY RANDOM() LIMIT ?"
		args = []interface{}{phase, number, limit}
	}

	err := r.db.Select(&templates, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get random templates: %w", err)
	}
	return templates, nil
}
