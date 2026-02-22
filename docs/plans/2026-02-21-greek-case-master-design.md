# Greek Case Master - Design Document

**Date:** 2026-02-21
**Status:** Approved

## Architecture Overview

**Greek Case Master** is a single-binary Go CLI application using:

- **Bubble Tea TUI framework** for the interactive practice interface
- **SQLite database** for storing noun declensions, generated sentences, and explanations
- **Claude API** for one-time generation during import (declensions + contextual sentences + grammar explanations)
- **Cobra** for command routing (`practice`, `import`, `add`, `list`)

### Key Architectural Principles

1. **Import-time generation, practice-time retrieval** - All AI generation happens during `import`, stored in SQLite. Practice sessions are entirely offline with instant feedback.

2. **Single binary distribution** - Compile to a single executable with embedded migrations. SQLite database created at `~/.greekmaster/greekmaster.db` on first run.

3. **Robust import process** - Checkpoint-based resume, retry with exponential backoff, progress indicators for long imports.

4. **Clean separation** - Domain models (Noun, Sentence), storage layer (SQLite repo), AI client (Claude API wrapper), TUI components (Bubble Tea models).

## Data Model & Database Schema

### SQLite Schema

**1. `nouns` table:**
```sql
CREATE TABLE nouns (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    english TEXT NOT NULL,
    gender TEXT NOT NULL, -- "masculine", "feminine", "neuter", "invariable"
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
```

**2. `sentences` table:**
```sql
CREATE TABLE sentences (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    noun_id INTEGER NOT NULL,
    english_prompt TEXT NOT NULL, -- "I see ___ (the teacher)"
    greek_sentence TEXT NOT NULL, -- "Βλέπω τον δάσκαλο"
    correct_answer TEXT NOT NULL, -- "τον δάσκαλο"
    case_type TEXT NOT NULL, -- "nominative", "genitive", "accusative"
    number TEXT NOT NULL, -- "singular", "plural"
    difficulty_phase INTEGER NOT NULL, -- 1, 2, or 3
    context_type TEXT NOT NULL, -- "direct_object", "possession", "preposition"
    preposition TEXT, -- "σε", "από", "για" if applicable
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (noun_id) REFERENCES nouns(id)
);
```

**3. `explanations` table:**
```sql
CREATE TABLE explanations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sentence_id INTEGER NOT NULL,
    translation TEXT NOT NULL, -- full Greek sentence translation
    syntactic_role TEXT NOT NULL, -- "direct object of βλέπω, requires accusative"
    morphology TEXT NOT NULL, -- "ο δάσκαλος (nom.) → τον δάσκαλο (acc.)"
    FOREIGN KEY (sentence_id) REFERENCES sentences(id)
);
```

**4. `import_checkpoints` table:**
```sql
CREATE TABLE import_checkpoints (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    csv_filename TEXT NOT NULL,
    last_processed_row INTEGER NOT NULL,
    status TEXT NOT NULL, -- "in_progress", "completed", "failed"
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Import Process Flow

**Command:** `greekmaster import words.csv`

### Step-by-Step Execution

**1. Validation phase:**
- Check CSV exists and is readable
- Validate headers (english, greek, attribute)
- Check `ANTHROPIC_API_KEY` is set
- Check for existing import checkpoint

**2. Resume or start:**
- If checkpoint exists: ask user "Resume from row X or start fresh?"
- Create/update checkpoint record with status "in_progress"

**3. Process each row with checkpoint:**

For each noun in CSV (starting from checkpoint):

a. **Generate declensions** via Claude:
   - Prompt: "Given the Greek noun '{greek}' with gender '{gender}', provide all declined forms with articles in JSON format"
   - Parse JSON response and validate

b. **Generate 10-15 sentences** via Claude:
   - 3 accusative sentences (direct object after different transitive verbs like βλέπω, ψάχνω, θέλω)
   - 3 genitive sentences (possession contexts)
   - 3 accusative with prepositions (σε, από)
   - 3 genitive with preposition (για)
   - Mix singular and plural forms
   - Return as JSON array with all metadata

c. **Generate explanations** via Claude:
   - For each sentence, provide translation, syntactic role, and morphology explanation
   - Return as JSON matching sentence array

d. **Parse and insert** into SQLite (nouns, sentences, explanations tables)

e. **Update checkpoint** after each noun (enables resume)

f. **Retry logic**: 3 attempts with exponential backoff (1s, 2s, 4s) on API failures

**4. Progress display:**
- Show progress bar: "Processing: 45/200 nouns (22%)"
- Show current noun: "Generating sentences for 'teacher' (δάσκαλος)..."
- Show API call counts and estimated time remaining

**5. Completion:**
- Mark checkpoint as "completed"
- Display summary: "Imported 200 nouns, 2,847 sentences generated"

## Practice Session Flow

**Command:** `greekmaster practice`

### Interactive TUI Flow (Bubble Tea)

**1. Session setup screen:**

Prompts:
- "Select difficulty level:"
  - Beginner: Focus on accusative (direct objects)
  - Intermediate: Focus on genitive (possession)
  - Advanced: Mixed cases with prepositions
- "Include plural forms? (Y/n)"
- "Select session type:"
  - Quick (10 questions)
  - Standard (25 questions)
  - Long (50 questions)
  - Endless (practice until quit)

**2. Question screen layout:**
```
┌─────────────────────────────────────────────────────┐
│  Greek Case Master - Question 5/25                  │
├─────────────────────────────────────────────────────┤
│                                                      │
│  I see ___ (the teacher)                            │
│                                                      │
│  Your answer:                                        │
│  > τον δάσκαλο_                                     │
│                                                      │
│  [Enter to submit] [Ctrl+C to quit]                │
└─────────────────────────────────────────────────────┘
```

**3. Feedback screen (correct answer):**
```
┌─────────────────────────────────────────────────────┐
│  ✓ Correct!                                         │
├─────────────────────────────────────────────────────┤
│  Translation: Βλέπω τον δάσκαλο                     │
│                                                      │
│  Syntactic Role: 'τον δάσκαλο' is the direct       │
│  object of βλέπω (I see). Direct objects require   │
│  the accusative case.                               │
│                                                      │
│  Morphology: ο δάσκαλος (nominative) becomes       │
│  τον δάσκαλο (accusative). The article changes     │
│  from ο → τον and the ending changes -ος → -ο.     │
│                                                      │
│  [Press any key to continue]                        │
└─────────────────────────────────────────────────────┘
```

**4. Feedback screen (incorrect answer):**
```
┌─────────────────────────────────────────────────────┐
│  ✗ Incorrect                                        │
├─────────────────────────────────────────────────────┤
│  You entered:  ο δάσκαλος                           │
│  Correct answer: τον δάσκαλο                        │
│                                                      │
│  Translation: Βλέπω τον δάσκαλο                     │
│                                                      │
│  [Full grammar explanation as above]                │
│                                                      │
│  [Press any key to continue]                        │
└─────────────────────────────────────────────────────┘
```

**5. Session completion (fixed mode only):**
```
┌─────────────────────────────────────────────────────┐
│  Session Complete!                                  │
├─────────────────────────────────────────────────────┤
│  Answered: 25/25                                    │
│  Accuracy: 84%                                      │
│                                                      │
│  [Press 'q' to quit or 'r' to start new session]   │
└─────────────────────────────────────────────────────┘
```

### Sentence Selection Logic

- Query sentences matching: difficulty_phase + number (singular/plural based on user choice)
- Shuffle the result set
- Present in shuffled order (no repetition until all exhausted)
- Reshuffle if user continues past available sentences

### Answer Validation

- **Exact match only** - must match exactly including accents, spacing (e.g., "τον δάσκαλο")
- No fuzzy matching or typo tolerance
- Case-sensitive Unicode comparison

## Claude API Integration

### Configuration

- **SDK:** `github.com/anthropic-ai/anthropic-sdk-go`
- **API Key:** Read from `ANTHROPIC_API_KEY` environment variable
- **Model:** `claude-3-5-sonnet-20241022`
- **Temperature:** 0.3 (consistent, accurate linguistic output)

### Three API Call Types

**1. Declension generation:**

```
Prompt: "You are a Modern Greek grammar expert. Given the noun '{greek}'
({english}) with gender '{gender}', provide all declined forms with their
articles in this JSON format:
{
  "nominative_sg": "...", "nom_sg_article": "...",
  "genitive_sg": "...", "gen_sg_article": "...",
  "accusative_sg": "...", "acc_sg_article": "...",
  "nominative_pl": "...", "nom_pl_article": "...",
  "genitive_pl": "...", "gen_pl_article": "...",
  "accusative_pl": "...", "acc_pl_article": "..."
}
Return only valid JSON, no explanation."
```

**2. Sentence generation:**

```
Prompt: "Generate 12 practice sentences for learning Greek cases using
the noun '{greek}' ({english}, {gender}). Create:
- 3 sentences with accusative case (direct object after different
  transitive verbs like βλέπω, ψάχνω, θέλω)
- 3 sentences with genitive case (possession contexts)
- 3 sentences with accusative after prepositions (σε, από)
- 3 sentences with genitive after preposition (για)

Mix singular and plural forms. Return as JSON array:
[{
  "english_prompt": "I see ___ (the teacher)",
  "greek_sentence": "Βλέπω τον δάσκαλο",
  "correct_answer": "τον δάσκαλο",
  "case_type": "accusative",
  "number": "singular",
  "difficulty_phase": 1,
  "context_type": "direct_object",
  "preposition": null
}, ...]"
```

**3. Explanation generation:**

```
Prompt: "For the sentence '{greek_sentence}' where the answer is
'{correct_answer}', provide a grammar explanation as the Modern Greek
Grammar Analyst. Return JSON:
{
  "translation": "Full English translation",
  "syntactic_role": "Explain why this case is required",
  "morphology": "Explain the form transformation"
}
Be concise and analytical, no conversational filler."
```

### Error Handling

- **Rate limit errors (429):** Exponential backoff (1s, 2s, 4s)
- **API errors (500, 503):** Retry 3 times
- **Invalid JSON response:** Log error, skip noun, continue
- **Network timeouts:** Retry with backoff
- **All failures logged** to `~/.greekmaster/import.log`

## Project Structure

```
greekmaster/
├── cmd/
│   └── greekmaster/
│       └── main.go              # Entry point, Cobra root command
├── internal/
│   ├── models/
│   │   ├── noun.go              # Noun, Sentence, Explanation structs
│   │   └── session.go           # Session configuration
│   ├── storage/
│   │   ├── migrations.go        # Embedded SQL migrations
│   │   ├── repository.go        # SQLite CRUD operations
│   │   └── checkpoint.go        # Import checkpoint logic
│   ├── ai/
│   │   ├── client.go            # Claude API wrapper
│   │   ├── prompts.go           # Prompt templates
│   │   └── retry.go             # Retry/backoff logic
│   ├── importer/
│   │   ├── csv.go               # CSV parsing
│   │   └── processor.go         # Import orchestration
│   └── tui/
│       ├── practice.go          # Practice session Bubble Tea model
│       ├── setup.go             # Session setup screens
│       └── feedback.go          # Answer feedback screens
├── go.mod
└── README.md
```

## Dependencies

```go
require (
    github.com/spf13/cobra v1.8.0           // CLI framework
    github.com/charmbracelet/bubbletea v0.25.0  // TUI framework
    github.com/charmbracelet/lipgloss v0.9.1    // TUI styling
    github.com/anthropic-ai/anthropic-sdk-go v0.1.0  // Claude API
    github.com/mattn/go-sqlite3 v1.14.18    // SQLite driver
    github.com/jmoiron/sqlx v1.3.5          // SQL extensions
)
```

## Build & Distribution

- **Single binary:** `go build -o greekmaster cmd/greekmaster/main.go`
- **Cross-platform:** Use `GOOS` and `GOARCH` for Linux/macOS/Windows
- **Embed migrations:** Use `//go:embed` to include SQL schema in binary
- **Database location:** `~/.greekmaster/greekmaster.db` (created on first run)

## Configuration

- **API key:** `ANTHROPIC_API_KEY` environment variable only
- **Database path:** Default `~/.greekmaster/`, can override with `--db-path` flag
- No other configuration needed for MVP

## CLI Commands

### Standard Command Set

1. **`greekmaster practice`** - Launch interactive practice session (Bubble Tea TUI)
2. **`greekmaster import <csv-file>`** - Import word bank from CSV with AI generation
3. **`greekmaster add`** - Interactively add a single noun with AI generation
4. **`greekmaster list`** - Show all nouns in database

## Testing Strategy

### Unit Tests

- CSV parsing logic (handle malformed rows, missing columns)
- Answer validation (exact string matching, Unicode handling)
- Sentence selection logic (shuffling, no repetition)
- Retry/backoff algorithms

### Integration Tests

- SQLite operations (CRUD, transactions, migrations)
- Mock Claude API responses (test JSON parsing without real API calls)
- Import checkpoint resume logic

### Manual Testing

- Full import flow with sample CSV (5-10 nouns)
- Practice session with real Bubble Tea interface
- Error scenarios (missing API key, network failures)

## Error Handling

### Import Errors

- **Missing API key:** "ANTHROPIC_API_KEY environment variable not set. Please set it and try again."
- **CSV format errors:** "Invalid CSV at row 5: missing 'attribute' column"
- **API failures:** "Failed to generate declensions for 'δάσκαλος' after 3 attempts. Skipping. See import.log for details."
- **Detailed logs:** Write to `~/.greekmaster/import.log`

### Practice Errors

- **Empty database:** "No sentences found. Please run 'greekmaster import <csv>' first."
- **Database corruption:** "Database error. Try deleting ~/.greekmaster/greekmaster.db and reimporting."

### Graceful Degradation

- If a sentence has missing explanation, show answer validation only
- If Greek input has encoding issues, show helpful error with Unicode info

## Non-Goals (Out of Scope for MVP)

1. **User statistics and progress tracking** - No accuracy history, streak tracking, or performance analytics
2. **Spaced repetition algorithms** - No SRS scheduling or difficulty adjustment
3. **Web or mobile UI** - CLI only
4. **Multi-user support** - Single user per database instance
5. **Verb conjugation practice** - Focus exclusively on noun declension
6. **Audio/pronunciation features** - Text-only interface
7. **Customizable lesson plans** - Fixed three-phase difficulty structure
8. **Offline sentence generation** - Claude API required for import
9. **Export/sharing of progress** - No export of session results

## Success Metrics

### Functional Completeness

- Can successfully import a 50+ noun CSV file
- Generates 10-15 varied sentences per noun
- Practice sessions run smoothly with immediate feedback
- Answer validation works correctly for all Greek Unicode characters

### Technical Success

- Import handles API failures gracefully with resume capability
- Single binary under 20MB
- Practice feedback appears instantly (< 100ms database queries)
- No crashes during normal operation

### Learning Effectiveness (Qualitative)

- Grammar explanations are accurate and clear
- Sentence variety feels natural (not repetitive)
- Difficulty levels appropriately challenge learners

## Design Decisions Summary

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Platform | CLI application | Validate core logic before UI investment |
| Language | Go | Fast, single binary, excellent Unicode support |
| TUI Framework | Bubble Tea | Rich interactive experience for practice |
| Storage | SQLite | Single file, structured queries, progress tracking |
| AI Generation | Claude API at import time | Generate once, practice offline with instant feedback |
| CSV Format | english, greek, attribute | Current format, expand during import |
| Declension Generation | Claude API | Handles irregular forms and articles |
| Sentence Count | 10-15 per noun | Good variety without repetition |
| Answer Input | Single field (article + noun) | Natural, matches real usage |
| Answer Validation | Exact match with accents | Ensures proper orthography learning |
| Feedback | Full grammar explanation | Core learning mechanism (FR2) |
| Difficulty | Session-based level selection | User control with appropriate mixing |
| Session Types | Fixed (10/25/50) + Endless | Flexibility for different practice styles |
| Plural Forms | User preference toggle | Accommodates different skill levels |
| Statistics | None for MVP | Focus on core learning experience |
| API Key | Environment variable | Simple, follows standard practices |
| Error Handling | Retry + checkpoint/resume | Robust for large imports |
