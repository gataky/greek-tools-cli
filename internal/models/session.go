package models

// SessionConfig holds practice session configuration
type SessionConfig struct {
	DifficultyLevel string // "beginner", "intermediate", "advanced"
	IncludePlural   bool
	QuestionCount   int // 0 for endless mode
}
