# Tasks: Greek Case Master Implementation

## Relevant Files

### Main Application
- `cmd/greekmaster/main.go` - Application entry point, Cobra CLI setup
- `go.mod` - Go module definition with dependencies
- `go.sum` - Dependency checksums
- `README.md` - User documentation and setup instructions

### Models Package
- `internal/models/noun.go` - Noun struct and related types
- `internal/models/sentence.go` - Sentence struct and related types
- `internal/models/explanation.go` - Explanation struct
- `internal/models/session.go` - SessionConfig struct for practice configuration

### Storage Package
- `internal/storage/repository.go` - Database interface and SQLite implementation
- `internal/storage/migrations.go` - Embedded SQL migrations using go:embed
- `internal/storage/checkpoint.go` - Import checkpoint management
- `migrations/001_initial_schema.sql` - Database schema SQL file

### AI Package
- `internal/ai/client.go` - Claude API client wrapper
- `internal/ai/prompts.go` - Prompt template functions
- `internal/ai/retry.go` - Exponential backoff and retry logic

### Importer Package
- `internal/importer/csv.go` - CSV file parsing
- `internal/importer/processor.go` - Import orchestration and AI generation coordination

### TUI Package
- `internal/tui/practice.go` - Practice session Bubble Tea model
- `internal/tui/setup.go` - Session configuration screens
- `internal/tui/feedback.go` - Answer feedback display components

### Commands Package
- `internal/commands/import.go` - Import command implementation
- `internal/commands/practice.go` - Practice command implementation
- `internal/commands/add.go` - Add single noun command
- `internal/commands/list.go` - List nouns command

### Test Files
- `internal/storage/repository_test.go` - Repository unit tests
- `internal/ai/retry_test.go` - Retry logic unit tests
- `internal/importer/csv_test.go` - CSV parsing tests
- `internal/models/noun_test.go` - Model validation tests

### Test Data
- `testdata/sample_words.csv` - Sample CSV for testing (5-10 nouns)
- `testdata/mock_responses.json` - Mock Claude API responses

### Notes

- This is a new Go project starting from scratch
- Tests will use Go's standard `testing` package
- Run tests with: `go test ./...` (all packages) or `go test ./internal/storage` (specific package)
- Build with: `go build -o greekmaster cmd/greekmaster/main.go`
- SQLite database will be created at `~/.greekmaster/greekmaster.db`

## Instructions for Completing Tasks

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [ ]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

Example:
- `- [ ] 1.1 Read file` → `- [x] 1.1 Read file` (after completing)

Update the file after completing each sub-task, not just after completing an entire parent task.

## Tasks

- [x] 0.0 Create feature branch
  - [x] 0.1 Initialize git repository if not already done with `git init`
  - [x] 0.2 Create and checkout a new branch with `git checkout -b feature/greek-case-master-mvp`

- [x] 1.0 Initialize Go project and dependencies
  - [x] 1.1 Create project directory structure: `cmd/greekmaster/`, `internal/{models,storage,ai,importer,tui,commands}/`, `migrations/`, `testdata/`
  - [x] 1.2 Initialize Go module with `go mod init github.com/gataky/greekmaster`
  - [x] 1.3 Add core dependencies to go.mod:
    - `github.com/spf13/cobra@v1.10.2` for CLI framework
    - `github.com/charmbracelet/bubbletea@v0.25.0` for TUI
    - `github.com/charmbracelet/lipgloss@v0.9.1` for TUI styling
    - `github.com/anthropics/anthropic-sdk-go@v1.26.0` for Claude API (note: correct org is anthropics not anthropic-ai)
    - `github.com/mattn/go-sqlite3@v1.14.18` for SQLite driver
    - `github.com/jmoiron/sqlx@v1.3.5` for SQL extensions
  - [x] 1.4 Run `go mod tidy` to download dependencies and create go.sum
  - [x] 1.5 Create basic `cmd/greekmaster/main.go` with Cobra root command setup and version flag
  - [x] 1.6 Verify project builds with `go build -o greekmaster cmd/greekmaster/main.go`

- [x] 2.0 Implement database layer with migrations
  - [x] 2.1 Create `migrations/001_initial_schema.sql` with all four table definitions (nouns, sentences, explanations, import_checkpoints) as specified in PRD section 5
  - [x] 2.2 Create `internal/models/noun.go` with Noun struct including all fields (ID, English, Gender, 6 declined forms, 6 articles, CreatedAt) with db tags
  - [x] 2.3 Create `internal/models/sentence.go` with Sentence struct including all fields (ID, NounID, EnglishPrompt, GreekSentence, CorrectAnswer, CaseType, Number, DifficultyPhase, ContextType, Preposition, CreatedAt) with db tags
  - [x] 2.4 Create `internal/models/explanation.go` with Explanation struct (ID, SentenceID, Translation, SyntacticRole, Morphology) with db tags
  - [x] 2.5 Create `internal/models/session.go` with SessionConfig struct (DifficultyLevel, IncludePlural, QuestionCount)
  - [x] 2.6 Create `internal/storage/migrations.go` with embedded SQL migrations using `//go:embed` directive and migration execution function
  - [x] 2.7 Create `internal/storage/repository.go` with Repository interface defining all CRUD methods (CreateNoun, GetNoun, ListNouns, CreateSentence, GetRandomSentences, CreateExplanation, GetExplanationBySentenceID)
  - [x] 2.8 Implement SQLite repository struct in `repository.go` with sqlx.DB connection, constructor function that creates ~/.greekmaster/ directory and initializes database
  - [x] 2.9 Implement all Repository interface methods using prepared statements and transactions where appropriate
  - [x] 2.10 Create `internal/storage/checkpoint.go` with checkpoint CRUD functions (CreateCheckpoint, UpdateCheckpoint, GetCheckpointByFilename)

- [x] 3.0 Implement Claude API client with retry logic
  - [x] 3.1 Create `internal/ai/retry.go` with exponential backoff function that takes retry count and returns wait duration (1s, 2s, 4s progression)
  - [x] 3.2 Implement retry wrapper function that accepts an API call function and executes it with up to 3 retries on failures (429, 500, 503 status codes)
  - [x] 3.3 Create `internal/ai/prompts.go` with three template functions:
    - `GenerateDeclensionPrompt(greek, english, gender string) string` - returns formatted prompt for declension generation
    - `GenerateSentencesPrompt(greek, english, gender string) string` - returns formatted prompt for sentence generation
    - `GenerateExplanationPrompt(greekSentence, correctAnswer string) string` - returns formatted prompt for explanation generation
  - [x] 3.4 Create `internal/ai/client.go` with ClaudeClient struct containing API key and Anthropic SDK client
  - [x] 3.5 Implement `NewClaudeClient() (*ClaudeClient, error)` constructor that reads ANTHROPIC_API_KEY from environment and initializes SDK client
  - [x] 3.6 Implement `GenerateDeclensions(greek, english, gender string) (*DeclensionResponse, error)` method that calls Claude API with declension prompt, parses JSON response, and returns structured data
  - [x] 3.7 Implement `GenerateSentences(greek, english, gender string) ([]SentenceResponse, error)` method that calls Claude API with sentence prompt and parses JSON array response
  - [x] 3.8 Implement `GenerateExplanations(sentences []SentenceResponse) ([]ExplanationResponse, error)` method that generates explanations for multiple sentences in batch
  - [x] 3.9 Add error handling for invalid JSON responses, log errors to ~/.greekmaster/import.log

- [ ] 4.0 Implement CSV import and AI generation orchestration
  - [ ] 4.1 Create `internal/importer/csv.go` with function to parse CSV file, validate headers (english, greek, attribute), and return slice of CSVRow structs
  - [ ] 4.2 Add CSV validation logic to check for empty required fields and valid gender values (masculine, feminine, neuter, invariable)
  - [ ] 4.3 Create `internal/importer/processor.go` with ImportProcessor struct that holds ClaudeClient and Repository dependencies
  - [ ] 4.4 Implement `ProcessImport(csvPath string) error` function that orchestrates the full import workflow
  - [ ] 4.5 Add checkpoint detection logic - check if import was previously started for this CSV file and prompt user to resume or start fresh
  - [ ] 4.6 Implement progress tracking - display current noun being processed, completion percentage, and API call count
  - [ ] 4.7 For each CSV row:
    - Call ClaudeClient.GenerateDeclensions() and store result
    - Call ClaudeClient.GenerateSentences() and store 12 sentences
    - Call ClaudeClient.GenerateExplanations() for all sentences
    - Insert Noun into database with all declined forms
    - Insert all 12 Sentences with foreign key to noun
    - Insert all 12 Explanations with foreign keys to sentences
    - Update checkpoint with current row number
  - [ ] 4.8 Add error handling with retry logic - if API call fails after 3 retries, log error to import.log, skip noun, and continue
  - [ ] 4.9 Display completion summary showing total nouns imported and total sentences generated
  - [ ] 4.10 Create `internal/commands/import.go` with Cobra command definition that calls ImportProcessor

- [ ] 5.0 Implement practice TUI with Bubble Tea
  - [ ] 5.1 Create `internal/tui/setup.go` with SetupModel Bubble Tea model for session configuration
  - [ ] 5.2 Implement setup screens that prompt for:
    - Difficulty level selection (Beginner/Intermediate/Advanced) using arrow keys
    - Plural inclusion toggle (Y/n)
    - Session type selection (Quick 10/Standard 25/Long 50/Endless)
  - [ ] 5.3 Parse user selections and create SessionConfig struct
  - [ ] 5.4 Create `internal/tui/practice.go` with PracticeModel Bubble Tea model that holds SessionConfig, current sentence, user input, question counter, correct answer count
  - [ ] 5.5 Implement Init() method that queries sentences from database based on SessionConfig (filter by difficulty_phase and number), shuffles results, and stores in model
  - [ ] 5.6 Implement Update() method to handle:
    - Key presses for text input (Greek characters)
    - Enter key to submit answer
    - Ctrl+C or 'q' to quit
    - After feedback, any key to continue to next question
  - [ ] 5.7 Implement View() method to render question screen with English prompt, input field, and instructions
  - [ ] 5.8 Implement answer validation logic - exact Unicode string comparison between user input and correct_answer from database
  - [ ] 5.9 Create `internal/tui/feedback.go` with FeedbackModel Bubble Tea model for displaying results
  - [ ] 5.10 Implement feedback screen that displays:
    - Correct/Incorrect indicator (✓ or ✗)
    - User's answer and correct answer (if incorrect)
    - Full explanation retrieved from database (translation, syntactic role, morphology)
    - Instruction to press any key to continue
  - [ ] 5.11 For fixed sessions, implement completion screen showing total answered and accuracy percentage, with options to quit or restart
  - [ ] 5.12 Add Greek Unicode rendering support with proper UTF-8 handling
  - [ ] 5.13 Style TUI with lipgloss for borders, colors, and layout
  - [ ] 5.14 Create `internal/commands/practice.go` with Cobra command that checks database is not empty, creates Bubble Tea program with SetupModel, and runs it

- [ ] 6.0 Implement add and list commands
  - [ ] 6.1 Create `internal/commands/add.go` with interactive prompts for:
    - English translation input
    - Greek nominative singular input
    - Gender selection (provide options: 1. masculine, 2. feminine, 3. neuter, 4. invariable)
  - [ ] 6.2 After collecting input, call ClaudeClient to generate declensions, sentences, and explanations (same as import process)
  - [ ] 6.3 Insert generated data into database and display confirmation with sentence count
  - [ ] 6.4 Create `internal/commands/list.go` with table display of all nouns
  - [ ] 6.5 Implement query to retrieve all nouns ordered by ID
  - [ ] 6.6 Format output as table with columns: ID, English, Greek (nominative), Gender
  - [ ] 6.7 Add pagination if more than 50 nouns (display 50 per page with navigation prompt)
  - [ ] 6.8 Wire both commands into main.go Cobra command tree

- [ ] 7.0 Implement testing suite
  - [ ] 7.1 Create `testdata/sample_words.csv` with 5 test nouns (mix of genders: masculine, feminine, neuter)
  - [ ] 7.2 Create `testdata/mock_responses.json` with sample Claude API responses for declensions, sentences, and explanations
  - [ ] 7.3 Create `internal/storage/repository_test.go` with unit tests for:
    - Database initialization and migrations
    - CreateNoun, GetNoun, ListNouns operations
    - CreateSentence, GetRandomSentences with filtering
    - CreateExplanation, GetExplanationBySentenceID
    - Foreign key constraints
    - Transaction rollback on error
  - [ ] 7.4 Create `internal/ai/retry_test.go` with unit tests for:
    - Exponential backoff calculation (verify 1s, 2s, 4s progression)
    - Retry wrapper with mock functions that fail N times then succeed
    - Max retry limit (verify stops after 3 attempts)
  - [ ] 7.5 Create `internal/importer/csv_test.go` with unit tests for:
    - Valid CSV parsing
    - Invalid CSV (missing columns, empty rows, invalid gender values)
    - Unicode handling for Greek characters
  - [ ] 7.6 Create `internal/models/noun_test.go` with validation tests for struct field types and db tags
  - [ ] 7.7 Run full test suite with `go test ./...` and verify all tests pass
  - [ ] 7.8 Manually test import command with testdata/sample_words.csv (requires ANTHROPIC_API_KEY set)
  - [ ] 7.9 Manually test practice command with beginner difficulty and verify:
    - Greek Unicode renders correctly
    - Exact answer matching works with accents
    - Explanations display properly
    - Session completion shows correct accuracy
  - [ ] 7.10 Manually test add command by adding one noun interactively
  - [ ] 7.11 Manually test list command to verify table display
  - [ ] 7.12 Test error scenarios:
    - Run import without ANTHROPIC_API_KEY (verify clear error message)
    - Run practice with empty database (verify error message)
    - Provide invalid CSV format (verify validation error)

- [ ] 8.0 Create documentation and final verification
  - [ ] 8.1 Create `README.md` with sections:
    - Project description and goals
    - Prerequisites (Go 1.21+, ANTHROPIC_API_KEY)
    - Installation instructions (clone repo, go build)
    - Configuration (setting ANTHROPIC_API_KEY)
    - Usage examples for all four commands (import, practice, add, list)
    - Database location (~/.greekmaster/)
    - Troubleshooting section (common errors)
  - [ ] 8.2 Add code comments to all exported functions and types
  - [ ] 8.3 Run `go fmt ./...` to format all code
  - [ ] 8.4 Run `go vet ./...` to check for common mistakes
  - [ ] 8.5 Build final binary with `go build -o greekmaster cmd/greekmaster/main.go`
  - [ ] 8.6 Verify binary size is under 20MB with `ls -lh greekmaster`
  - [ ] 8.7 Test full workflow end-to-end:
    - Set ANTHROPIC_API_KEY
    - Import sample CSV with 5 nouns
    - Verify database created at ~/.greekmaster/greekmaster.db
    - Run practice session and complete 10 questions
    - Add one noun manually
    - List all nouns (should show 6 total)
  - [ ] 8.8 Commit all code with message: "Implement Greek Case Master MVP - CLI educational tool for Greek noun declension"
  - [ ] 8.9 Create git tag `v0.1.0` for MVP release
