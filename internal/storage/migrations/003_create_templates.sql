-- Create sentence_templates table for template-based sentence generation
CREATE TABLE IF NOT EXISTS sentence_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    english_template TEXT NOT NULL,
    greek_template TEXT NOT NULL,
    article_field TEXT NOT NULL,
    noun_form_field TEXT NOT NULL,
    case_type TEXT NOT NULL CHECK(case_type IN ('nominative', 'genitive', 'accusative')),
    number TEXT NOT NULL CHECK(number IN ('singular', 'plural', 'both')),
    difficulty_phase INTEGER NOT NULL CHECK(difficulty_phase IN (1, 2, 3)),
    context_type TEXT NOT NULL CHECK(context_type IN ('direct_object', 'possession', 'preposition')),
    preposition TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_templates_difficulty ON sentence_templates(difficulty_phase, number);
