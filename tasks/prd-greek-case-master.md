# PRD: Greek Case Master CLI Application

## 1. Introduction/Overview

Greek Case Master is a command-line educational tool designed to help English speakers learning Modern Greek transition from rote memorization of noun declension tables to instinctive application of cases (Nominative, Genitive, Accusative) in real sentence contexts.

The application presents English sentence prompts with missing Greek nouns. Users must provide the correctly declined Greek article + noun combination. The system validates answers with exact matching and provides detailed grammar explanations (translation, syntactic role, morphology) for every answer.

**Problem it solves:** Learners struggle to apply noun cases "on the fly" during sentence construction because cases feel abstract when studied in isolation.

**What needs to be built:** A Go CLI application with an interactive TUI that generates practice sentences using Claude API during import, stores them in SQLite, and provides offline practice with instant feedback.

## 2. Goals

### Functional Goals

1. Enable learners to practice Greek noun declension in realistic sentence contexts
2. Provide immediate, analytical grammar feedback on every answer (correct or incorrect)
3. Support progressive difficulty from basic accusative cases to complex prepositional usage
4. Generate varied, natural practice sentences from a user-provided word bank

### Technical Goals

1. Import and process 50+ nouns with AI-generated declensions and sentences in under 10 minutes
2. Deliver instant practice feedback (< 100ms response time from SQLite queries)
3. Produce a single, portable binary under 20MB for cross-platform distribution
4. Handle API failures gracefully with checkpoint-based resume capability
5. Support proper Greek Unicode rendering in terminal environments

## 3. Tech Stack & Architecture

### Languages
- **Go 1.21+** - Primary language for CLI application

### Frameworks & Libraries

**Core Dependencies:**
- `github.com/spf13/cobra v1.8.0` - CLI command routing and argument parsing
- `github.com/charmbracelet/bubbletea v0.25.0` - Terminal User Interface framework for interactive practice
- `github.com/charmbracelet/lipgloss v0.9.1` - TUI styling and layout
- `github.com/anthropic-ai/anthropic-sdk-go v0.1.0` - Claude API client
- `github.com/mattn/go-sqlite3 v1.14.18` - SQLite database driver
- `github.com/jmoiron/sqlx v1.3.5` - SQL convenience wrapper with struct scanning

### Architectural Pattern

**Layered architecture with clear separation:**

1. **Command Layer** (`cmd/greekmaster/`) - Cobra command definitions and entry point
2. **Domain Layer** (`internal/models/`) - Core business entities (Noun, Sentence, Explanation, Session)
3. **Storage Layer** (`internal/storage/`) - SQLite repository pattern with embedded migrations
4. **AI Layer** (`internal/ai/`) - Claude API client with retry logic and prompt templates
5. **Import Layer** (`internal/importer/`) - CSV processing and AI generation orchestration
6. **TUI Layer** (`internal/tui/`) - Bubble Tea models for interactive practice interface

**Data Flow:**
```
CSV Import → Claude API (declensions + sentences + explanations) → SQLite
Practice Session → SQLite (sentence retrieval) → TUI (user interaction) → SQLite validation
```

### Existing Patterns to Follow

This is a new codebase. Establish patterns:
- Repository pattern for all database operations
- Dependency injection via constructor functions
- Embedded SQL migrations using `//go:embed`
- Error wrapping with context using `fmt.Errorf`
- Structured logging using standard library `log/slog`

### Data Storage

- **SQLite 3** - Single-file relational database
- **Location:** `~/.greekmaster/greekmaster.db`
- **Schema:** 4 tables (nouns, sentences, explanations, import_checkpoints)
- **Migrations:** Embedded in binary using `//go:embed`, run on first launch

### External Services/APIs

- **Claude API (Anthropic)** - Used during import only
  - Model: `claude-3-5-sonnet-20241022`
  - Temperature: 0.3
  - Authentication: `ANTHROPIC_API_KEY` environment variable
  - Rate limiting: Exponential backoff on 429 errors
  - Three prompt types: declension generation, sentence generation, explanation generation

## 4. Functional Requirements

### FR1: CSV Import with AI Generation

The system must support importing a CSV file containing Greek nouns and generating complete practice data.

**Input CSV format:**
- Columns: `english`, `greek`, `attribute`
- `english`: English translation (e.g., "teacher")
- `greek`: Greek nominative singular form (e.g., "δάσκαλος")
- `attribute`: Gender - "masculine", "feminine", "neuter", or "invariable"

**Process flow:**
1. Validate CSV structure and `ANTHROPIC_API_KEY` presence
2. Check for existing import checkpoint (resume capability)
3. For each row:
   - Call Claude API to generate all 6 declined forms (nom/gen/acc in sg/pl) with articles
   - Call Claude API to generate 12 practice sentences covering all three difficulty phases
   - Call Claude API to generate grammar explanations for each sentence
   - Store all data in SQLite (nouns, sentences, explanations tables)
   - Update checkpoint after each noun (atomic progress tracking)
4. Display progress: current noun, completion percentage, API call count
5. Handle failures: retry 3 times with exponential backoff (1s, 2s, 4s), skip noun if all fail, log to `~/.greekmaster/import.log`

**Validation rules:**
- CSV must have required columns
- Each row must have non-empty values for all three columns
- Claude API must return valid JSON matching expected schema
- Database writes must succeed (transaction rollback on failure)

**Edge cases:**
- Network interruption: Resume from last checkpoint
- Invalid JSON from Claude: Log error, skip noun, continue
- Duplicate nouns in CSV: Import as separate entries (allow practice variety)
- Empty CSV: Show error "No nouns found in CSV file"

### FR2: Interactive Practice Sessions

The system must provide an interactive terminal UI for practicing noun declension.

**Session setup:**
1. Prompt for difficulty level:
   - Beginner: Accusative case focus (direct objects after transitive verbs)
   - Intermediate: Genitive case focus (possession contexts)
   - Advanced: Mixed cases with prepositional usage
2. Prompt for plural inclusion: "Include plural forms? (Y/n)"
3. Prompt for session type:
   - Quick (10 questions)
   - Standard (25 questions)
   - Long (50 questions)
   - Endless (until user quits)

**Question presentation:**
- Display English prompt with blank: "I see ___ (the teacher)"
- Show input field for user to type Greek answer
- Accept answer on Enter key press
- Allow quit with Ctrl+C or 'q' key

**Answer validation:**
- Exact string matching (case-sensitive Unicode comparison)
- Must match article + noun exactly including all accents and spacing
- No fuzzy matching or typo tolerance

**Feedback display (immediate after each answer):**
- Show "✓ Correct!" or "✗ Incorrect"
- If incorrect, show: "You entered: [user_input]" and "Correct answer: [correct_answer]"
- Always display full grammar explanation:
  - Translation: Full Greek sentence with English translation
  - Syntactic Role: Why this case is required (e.g., "direct object of βλέπω requires accusative")
  - Morphology: Form transformation (e.g., "ο δάσκαλος (nom.) → τον δάσκαλο (acc.)")
- Wait for any key press to continue to next question

**Session completion (fixed-length sessions only):**
- Display: "Session Complete!"
- Show: Total answered (X/Y), Accuracy percentage
- Offer: 'q' to quit or 'r' to start new session

**Sentence selection logic:**
1. Query sentences matching: selected difficulty_phase + number (sg/pl based on user toggle)
2. Shuffle result set using random seed
3. Present in shuffled order without repetition within session
4. If endless mode exceeds available sentences, reshuffle and continue

### FR3: Manual Noun Addition

The system must allow adding individual nouns interactively.

**Process flow:**
1. Prompt for English translation
2. Prompt for Greek nominative singular form
3. Prompt for gender (provide options: masculine/feminine/neuter/invariable)
4. Call Claude API to generate declensions, sentences, and explanations (same as import)
5. Store in database
6. Display confirmation with generated sentence count

### FR4: Noun Listing

The system must display all nouns currently in the database.

**Display format:**
- Table with columns: ID, English, Greek (nominative), Gender
- Sort by ID (insertion order)
- Paginate if more than 50 nouns (show 50 per page with navigation)

### FR5: Database Initialization

The system must automatically initialize the database on first run.

**Process:**
1. Check if `~/.greekmaster/` directory exists, create if not
2. Check if `greekmaster.db` exists
3. If not, create database and run embedded SQL migrations
4. Verify schema with test query

## 5. Technical Specifications

### Data Models

**Go struct definitions:**

```go
// Noun represents a Greek noun with all declined forms
type Noun struct {
    ID            int64     `db:"id"`
    English       string    `db:"english"`
    Gender        string    `db:"gender"`
    NominativeSg  string    `db:"nominative_sg"`
    GenitiveSg    string    `db:"genitive_sg"`
    AccusativeSg  string    `db:"accusative_sg"`
    NominativePl  string    `db:"nominative_pl"`
    GenitivePl    string    `db:"genitive_pl"`
    AccusativePl  string    `db:"accusative_pl"`
    NomSgArticle  string    `db:"nom_sg_article"`
    GenSgArticle  string    `db:"gen_sg_article"`
    AccSgArticle  string    `db:"acc_sg_article"`
    NomPlArticle  string    `db:"nom_pl_article"`
    GenPlArticle  string    `db:"gen_pl_article"`
    AccPlArticle  string    `db:"acc_pl_article"`
    CreatedAt     time.Time `db:"created_at"`
}

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

// Explanation represents grammar analysis for a sentence
type Explanation struct {
    ID            int64  `db:"id"`
    SentenceID    int64  `db:"sentence_id"`
    Translation   string `db:"translation"`
    SyntacticRole string `db:"syntactic_role"`
    Morphology    string `db:"morphology"`
}

// ImportCheckpoint tracks import progress for resume capability
type ImportCheckpoint struct {
    ID                int64     `db:"id"`
    CSVFilename       string    `db:"csv_filename"`
    LastProcessedRow  int       `db:"last_processed_row"`
    Status            string    `db:"status"`
    UpdatedAt         time.Time `db:"updated_at"`
}

// SessionConfig holds practice session configuration
type SessionConfig struct {
    DifficultyLevel string // "beginner", "intermediate", "advanced"
    IncludePlural   bool
    QuestionCount   int // 0 for endless mode
}
```

### API Contracts (Claude API)

**1. Declension Generation Request:**

```json
{
  "model": "claude-3-5-sonnet-20241022",
  "temperature": 0.3,
  "max_tokens": 1000,
  "messages": [
    {
      "role": "user",
      "content": "You are a Modern Greek grammar expert. Given the noun 'δάσκαλος' (teacher) with gender 'masculine', provide all declined forms with their articles in this JSON format:\n{\n  \"nominative_sg\": \"...\", \"nom_sg_article\": \"...\",\n  \"genitive_sg\": \"...\", \"gen_sg_article\": \"...\",\n  \"accusative_sg\": \"...\", \"acc_sg_article\": \"...\",\n  \"nominative_pl\": \"...\", \"nom_pl_article\": \"...\",\n  \"genitive_pl\": \"...\", \"gen_pl_article\": \"...\",\n  \"accusative_pl\": \"...\", \"acc_pl_article\": \"...\"\n}\nReturn only valid JSON, no explanation."
    }
  ]
}
```

**Expected Response:**
```json
{
  "nominative_sg": "δάσκαλος", "nom_sg_article": "ο",
  "genitive_sg": "δασκάλου", "gen_sg_article": "του",
  "accusative_sg": "δάσκαλο", "acc_sg_article": "τον",
  "nominative_pl": "δάσκαλοι", "nom_pl_article": "οι",
  "genitive_pl": "δασκάλων", "gen_pl_article": "των",
  "accusative_pl": "δασκάλους", "acc_pl_article": "τους"
}
```

**2. Sentence Generation Request:**

```
Prompt: "Generate 12 practice sentences for learning Greek cases using the noun 'δάσκαλος' (teacher, masculine). Create:
- 3 sentences with accusative case (direct object after different transitive verbs like βλέπω, ψάχνω, θέλω)
- 3 sentences with genitive case (possession contexts)
- 3 sentences with accusative after prepositions (σε, από)
- 3 sentences with genitive after preposition (για)

Mix singular and plural forms. Return as JSON array with this structure:
[{
  \"english_prompt\": \"I see ___ (the teacher)\",
  \"greek_sentence\": \"Βλέπω τον δάσκαλο\",
  \"correct_answer\": \"τον δάσκαλο\",
  \"case_type\": \"accusative\",
  \"number\": \"singular\",
  \"difficulty_phase\": 1,
  \"context_type\": \"direct_object\",
  \"preposition\": null
}, ...]"
```

**Expected Response:** JSON array of 12 sentence objects

**3. Explanation Generation Request:**

```
Prompt: "For the sentence 'Βλέπω τον δάσκαλο' where the answer is 'τον δάσκαλο', provide a grammar explanation as the Modern Greek Grammar Analyst. Return JSON:
{
  \"translation\": \"I see the teacher\",
  \"syntactic_role\": \"'τον δάσκαλο' is the direct object of βλέπω (I see). Direct objects in Greek require the accusative case.\",
  \"morphology\": \"The masculine noun ο δάσκαλος (nominative) becomes τον δάσκαλο (accusative). The article changes from ο → τον and the noun ending changes from -ος → -ο.\"
}
Be concise and analytical, no conversational filler."
```

**Status Codes:**
- 200: Success
- 429: Rate limit (retry with exponential backoff)
- 500/503: Server error (retry 3 times)
- 401: Invalid API key (fail immediately with clear error)

### Component/Module Structure

**Project layout:**
```
greekmaster/
├── cmd/greekmaster/main.go          # Entry point, Cobra setup
├── internal/
│   ├── models/
│   │   ├── noun.go                  # Noun struct and methods
│   │   ├── sentence.go              # Sentence struct and methods
│   │   ├── explanation.go           # Explanation struct
│   │   └── session.go               # SessionConfig struct
│   ├── storage/
│   │   ├── migrations.go            # Embedded SQL with //go:embed
│   │   ├── repository.go            # Database interface + SQLite impl
│   │   └── checkpoint.go            # Checkpoint management
│   ├── ai/
│   │   ├── client.go                # Claude API client wrapper
│   │   ├── prompts.go               # Prompt template functions
│   │   └── retry.go                 # Exponential backoff logic
│   ├── importer/
│   │   ├── csv.go                   # CSV parsing
│   │   └── processor.go             # Import orchestration
│   └── tui/
│       ├── practice.go              # Practice session Bubble Tea model
│       ├── setup.go                 # Session configuration screens
│       └── feedback.go              # Answer feedback display
├── migrations/
│   └── 001_initial_schema.sql       # Embedded via //go:embed
├── go.mod
├── go.sum
└── README.md
```

**Repository interface:**
```go
type Repository interface {
    // Noun operations
    CreateNoun(noun *Noun) error
    GetNoun(id int64) (*Noun, error)
    ListNouns() ([]*Noun, error)

    // Sentence operations
    CreateSentence(sentence *Sentence) error
    GetRandomSentences(phase int, number string, limit int) ([]*Sentence, error)

    // Explanation operations
    CreateExplanation(explanation *Explanation) error
    GetExplanationBySentenceID(sentenceID int64) (*Explanation, error)

    // Checkpoint operations
    CreateCheckpoint(checkpoint *ImportCheckpoint) error
    UpdateCheckpoint(checkpoint *ImportCheckpoint) error
    GetCheckpointByFilename(filename string) (*ImportCheckpoint, error)
}
```

### Integration Points

**Database integration:**
- Use `sqlx.DB` for connection pooling
- Prepared statements for all queries (prevent SQL injection)
- Transactions for multi-table inserts (import process)
- Foreign key constraints enforced at schema level

**Claude API integration:**
- Single shared client with timeout (30 seconds per request)
- Rate limiting: 50 requests per minute (enforced client-side)
- Retry logic: Exponential backoff starting at 1 second, max 3 retries
- Request/response logging to `~/.greekmaster/import.log` for debugging

**Terminal integration:**
- Use Bubble Tea's `tea.Program` for TUI lifecycle
- Greek Unicode rendering (ensure UTF-8 locale)
- Handle terminal resize events gracefully
- Ctrl+C handling for clean shutdown

### Complex Logic: Sentence Selection Algorithm

**Pseudo code:**
```
function selectSentences(config SessionConfig) []Sentence:
    // Map difficulty to phase numbers
    phase := map[string]int{
        "beginner": 1,      // Accusative focus
        "intermediate": 2,  // Genitive focus
        "advanced": 3       // Mixed with prepositions
    }[config.DifficultyLevel]

    // Determine number filter
    numberFilter := "singular"
    if config.IncludePlural:
        numberFilter = "" // Empty means both sg and pl

    // Query database
    sentences := repository.GetRandomSentences(phase, numberFilter, 1000)

    // Shuffle using Fisher-Yates
    shuffle(sentences)

    // For fixed sessions, limit to question count
    if config.QuestionCount > 0:
        sentences = sentences[:min(len(sentences), config.QuestionCount)]

    return sentences
```

### Security Requirements

**Input Sanitization:**
- Validate all user input (CSV data, manual entry) before database insertion
- Escape special characters in SQL queries (use parameterized queries)
- Limit input lengths (noun: 100 chars, english: 100 chars, gender: 20 chars)

**API Key Security:**
- Never log or display API key
- Validate format before making requests
- Store in environment variable only (no config file storage in MVP)

**Database Security:**
- SQLite file permissions: 0600 (read/write owner only)
- No external network access to database
- Validate schema on startup (detect corruption)

### Performance Requirements

**Import performance:**
- Process 50 nouns (with 12 sentences each = 600 sentences total) in under 10 minutes
- API calls: ~150 total (50 declensions + 50 sentence batches + 50 explanation batches)
- Database inserts: Batch transactions (insert 50 nouns, 600 sentences, 600 explanations in 3 transactions)

**Practice performance:**
- Sentence retrieval from SQLite: < 50ms
- Answer validation: < 10ms (simple string comparison)
- TUI render time: < 16ms (60 FPS target)
- Total response time: < 100ms from answer submission to feedback display

**Binary size:**
- Compiled binary: < 20MB (including embedded migrations and dependencies)

**Memory usage:**
- Import process: < 100MB RAM
- Practice session: < 50MB RAM

## 6. Non-Goals (Out of Scope)

1. **User statistics and progress tracking** - No historical accuracy data, learning curves, or performance analytics. Each session is independent.

2. **Spaced repetition system (SRS)** - No algorithms for optimal review scheduling based on past performance.

3. **Web or mobile interfaces** - CLI only. GUI/web extensions are future work.

4. **Multi-user support** - Single database per user. No user accounts, authentication, or multi-tenancy.

5. **Verb conjugation** - Focus exclusively on noun declension. Verbs are separate feature.

6. **Audio features** - No text-to-speech, pronunciation guides, or audio recording.

7. **Lesson plan customization** - Fixed three-phase difficulty structure. No user-created exercise templates.

8. **Local LLM support** - Requires Claude API. No Ollama or other local model integration.

9. **Progress export/import** - No backup, export, or sharing of session data.

10. **Adjective agreement** - Nouns only. Adjectives are separate grammatical concept.

11. **Gamification** - No points, badges, leaderboards, or achievements.

## 7. Testing Requirements

### Unit Tests

**Must have:**
1. CSV parsing with malformed data (missing columns, empty rows, encoding issues)
2. Answer validation logic (exact match, Unicode comparison, whitespace handling)
3. Sentence selection algorithm (shuffling, no repetition, filtering by phase/number)
4. Exponential backoff calculation (verify timing: 1s, 2s, 4s, 8s)
5. JSON response parsing from Claude API (valid and invalid responses)

**Test data:**
- Sample CSV with 10 test nouns (mix of genders)
- Mock Claude API responses (JSON files)
- Edge cases: empty strings, very long strings, special Greek characters

### Integration Tests

**Must have:**
1. SQLite operations:
   - Schema creation and migrations
   - CRUD operations for all tables
   - Foreign key constraint enforcement
   - Transaction rollback on error
2. Import checkpoint flow:
   - Create checkpoint, fail midway, resume from checkpoint
   - Verify idempotency (resuming doesn't duplicate data)
3. Claude API with mocks:
   - Parse all three response types (declension, sentences, explanations)
   - Handle API errors (rate limits, invalid JSON, network timeout)

**Test environment:**
- Use in-memory SQLite database (`:memory:`)
- Mock HTTP client for Claude API calls
- Temporary directories for test databases and logs

### End-to-End Tests (Manual)

**Test scenarios:**
1. **Happy path import:**
   - Import sample CSV with 5 nouns
   - Verify database contains 5 nouns and ~60 sentences
   - Spot-check declensions and explanations for accuracy

2. **Interrupted import:**
   - Start import, kill process after 2 nouns
   - Restart import, verify resume from checkpoint
   - Ensure no duplicate data

3. **Practice session - beginner:**
   - Start Quick (10 questions) session
   - Answer mix of correct and incorrect
   - Verify feedback displays properly with Greek Unicode
   - Verify session completion summary

4. **Practice session - advanced with plurals:**
   - Select advanced difficulty + include plurals
   - Verify sentences include prepositions and plural forms
   - Test endless mode (quit with Ctrl+C)

5. **Error handling:**
   - Run import without API key → verify clear error message
   - Run practice with empty database → verify error message
   - Provide CSV with invalid gender value → verify validation error

### Test Coverage Expectations

- Unit test coverage: > 80% for business logic (models, storage, ai, importer packages)
- Integration test coverage: All repository methods, import orchestration
- Manual E2E: All commands (`import`, `practice`, `add`, `list`) with various inputs

## 8. Success Metrics

### Functional Success Criteria

1. **Import capability:**
   - Successfully import CSV with 50+ nouns without errors
   - Generate 10-15 sentences per noun (500-750 total)
   - Complete import in under 10 minutes (real-world API timing)

2. **Practice functionality:**
   - Answer validation correctly identifies exact matches (including Greek accents)
   - Grammar explanations display properly with Greek Unicode
   - Session types work correctly (10/25/50/endless)
   - Difficulty filtering produces appropriate sentence distribution

3. **Data quality:**
   - Spot-check 20 random sentences: all have accurate declensions
   - Spot-check 20 random explanations: all are grammatically correct and clear
   - No duplicate sentences in a single practice session

### Technical Success Criteria

1. **Performance:**
   - Import: < 10 minutes for 50 nouns (including API latency)
   - Practice feedback: < 100ms from answer submission to display
   - Binary size: < 20MB
   - Memory usage: < 100MB during import, < 50MB during practice

2. **Reliability:**
   - Import survives network interruption and resumes correctly
   - No crashes during normal operation (1000 practice questions without crash)
   - Checkpoint data remains consistent (no corruption after interrupted imports)

3. **Usability:**
   - Clear error messages for all failure modes (missing API key, empty DB, invalid CSV)
   - Greek characters render correctly in standard terminal emulators (macOS Terminal, iTerm2, GNOME Terminal)
   - TUI is responsive (no input lag)

### Learning Effectiveness (Qualitative Assessment)

1. **Grammar explanations:**
   - Manually review 20 random explanations
   - Verify syntactic role analysis is accurate
   - Verify morphology explanations match Greek grammar rules
   - Confirm concise, analytical tone (no conversational filler)

2. **Sentence variety:**
   - Review 50 random sentences for a single noun
   - Confirm diverse verbs, prepositions, and contexts
   - Ensure no repetitive patterns
   - Verify natural Greek sentence structure

3. **Difficulty appropriateness:**
   - Beginner level focuses on accusative (direct objects)
   - Intermediate focuses on genitive (possession)
   - Advanced includes prepositional usage and mixed cases

## 9. Open Questions

**Resolved during brainstorming:**
- Platform: CLI ✓
- Language: Go ✓
- TUI framework: Bubble Tea ✓
- Storage: SQLite ✓
- AI approach: Claude API at import ✓
- CSV format: english, greek, attribute ✓
- Sentence count: 10-15 per noun ✓
- Answer input: Single field (article + noun) ✓
- Validation strictness: Exact match ✓
- Feedback detail: Full grammar explanation ✓
- Difficulty progression: User-selected levels ✓
- Session structure: Fixed + endless modes ✓
- Plural forms: User toggle ✓
- Statistics: None for MVP ✓
- API key: Environment variable ✓

**Remaining open questions:**
None at this time. All requirements are well-defined for MVP implementation.

---

## Implementation Notes

This PRD is designed to be actionable by an AI assistant or junior developer. All technical decisions have been made, dependencies specified, and data structures defined. The implementation should follow this document strictly to ensure the MVP delivers core learning value without scope creep.

For implementation, consider using the following task breakdown:
1. Project setup (Go module, dependencies, directory structure)
2. Database layer (migrations, repository implementation)
3. Claude API client (with retry logic and prompt templates)
4. Import command (CSV parsing, orchestration, checkpoint logic)
5. Practice TUI (Bubble Tea models, session flow)
6. Add and List commands (simpler CRUD operations)
7. Testing (unit tests, integration tests, manual E2E)
8. Documentation (README with usage examples)
