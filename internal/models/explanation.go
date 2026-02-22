package models

// Explanation represents grammar analysis for a sentence
type Explanation struct {
	ID            int64  `db:"id"`
	SentenceID    int64  `db:"sentence_id"`
	Translation   string `db:"translation"`
	SyntacticRole string `db:"syntactic_role"`
	Morphology    string `db:"morphology"`
}
