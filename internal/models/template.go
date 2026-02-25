package models

import "time"

// SentenceTemplate represents a reusable sentence pattern
type SentenceTemplate struct {
	ID              int64     `db:"id"`
	EnglishTemplate string    `db:"english_template"`
	GreekTemplate   string    `db:"greek_template"`
	ArticleField    string    `db:"article_field"`
	NounFormField   string    `db:"noun_form_field"`
	CaseType        string    `db:"case_type"`
	Number          string    `db:"number"`
	DifficultyPhase int       `db:"difficulty_phase"`
	ContextType     string    `db:"context_type"`
	Preposition     *string   `db:"preposition"`
	CreatedAt       time.Time `db:"created_at"`
}
