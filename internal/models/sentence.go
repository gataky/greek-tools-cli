package models

import "time"

// Sentence represents a practice sentence
type Sentence struct {
	ID              int64     `db:"id"`
	NounID          int64     `db:"noun_id"`
	EnglishPrompt   string    `db:"english_prompt"`
	GreekSentence   string    `db:"greek_sentence"`
	CorrectAnswer   string    `db:"correct_answer"`
	CaseType        string    `db:"case_type"`
	Number          string    `db:"number"`
	DifficultyPhase int       `db:"difficulty_phase"`
	ContextType     string    `db:"context_type"`
	Preposition     *string   `db:"preposition"`
	CreatedAt       time.Time `db:"created_at"`
}
