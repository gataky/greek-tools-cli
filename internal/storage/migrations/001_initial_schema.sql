-- Greek Case Master Database Schema
-- Creates tables for nouns, sentences, explanations, and import checkpoints

CREATE TABLE IF NOT EXISTS nouns (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    english TEXT NOT NULL,
    gender TEXT NOT NULL CHECK(gender IN ('masculine', 'feminine', 'neuter', 'invariable')),
    nominative_sg TEXT NOT NULL,
    genitive_sg TEXT NOT NULL,
    accusative_sg TEXT NOT NULL,
    nominative_pl TEXT NOT NULL,
    genitive_pl TEXT NOT NULL,
    accusative_pl TEXT NOT NULL,
    nom_sg_article TEXT NOT NULL,
    gen_sg_article TEXT NOT NULL,
    acc_sg_article TEXT NOT NULL,
    nom_pl_article TEXT NOT NULL,
    gen_pl_article TEXT NOT NULL,
    acc_pl_article TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sentences (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    noun_id INTEGER NOT NULL,
    english_prompt TEXT NOT NULL,
    greek_sentence TEXT NOT NULL,
    correct_answer TEXT NOT NULL,
    case_type TEXT NOT NULL CHECK(case_type IN ('nominative', 'genitive', 'accusative')),
    number TEXT NOT NULL CHECK(number IN ('singular', 'plural')),
    difficulty_phase INTEGER NOT NULL CHECK(difficulty_phase IN (1, 2, 3)),
    context_type TEXT NOT NULL CHECK(context_type IN ('direct_object', 'possession', 'preposition')),
    preposition TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (noun_id) REFERENCES nouns(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_sentences_difficulty ON sentences(difficulty_phase, number);
CREATE INDEX IF NOT EXISTS idx_sentences_noun_id ON sentences(noun_id);

CREATE TABLE IF NOT EXISTS explanations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sentence_id INTEGER NOT NULL,
    translation TEXT NOT NULL,
    syntactic_role TEXT NOT NULL,
    morphology TEXT NOT NULL,
    FOREIGN KEY (sentence_id) REFERENCES sentences(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_explanations_sentence_id ON explanations(sentence_id);

CREATE TABLE IF NOT EXISTS import_checkpoints (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    csv_filename TEXT NOT NULL,
    last_processed_row INTEGER NOT NULL,
    status TEXT NOT NULL CHECK(status IN ('in_progress', 'completed', 'failed')),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_checkpoints_filename ON import_checkpoints(csv_filename);
