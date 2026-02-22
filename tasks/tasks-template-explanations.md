# Tasks: Template-Based Explanation System

## Relevant Files

### New Files
- `internal/explanations/generator.go` - Main explanation generation logic
- `internal/explanations/templates.go` - Rule templates for syntactic role explanations
- `internal/explanations/translation.go` - Translation generation from prompts
- `internal/explanations/morphology.go` - Morphology transformation formatting
- `internal/explanations/generator_test.go` - Comprehensive unit tests for explanation generation

### Modified Files
- `internal/tui/practice.go` - Update to use template generator instead of database lookup
- `internal/importer/processor.go` - Remove explanation generation and storage logic
- `internal/commands/import.go` - Remove --skip-explanations flag
- `internal/commands/add.go` - Remove explanation generation/storage and --skip-explanations flag
- `internal/storage/repository.go` - Remove explanation-related methods from interface and implementation
- `internal/storage/repository_test.go` - Remove explanation-related tests
- `internal/ai/client.go` - Remove GenerateExplanations method and ExplanationResponse type
- `internal/ai/prompts.go` - Remove GenerateExplanationPrompt function

### Database Migration
- `internal/storage/migrations/002_drop_explanations.sql` - SQL migration to drop explanations table

### Notes

- This is a refactoring task that replaces AI-generated explanations with template-based generation
- Run tests with: `go test ./...` (all packages) or `go test ./internal/explanations` (specific package)
- Build with: `go build -o greekmaster cmd/greekmaster/main.go`
- The goal is to eliminate ~12 API calls per noun (92% cost reduction)

## Instructions for Completing Tasks

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [ ]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

Example:
- `- [ ] 1.1 Read file` → `- [x] 1.1 Read file` (after completing)

Update the file after completing each sub-task, not just after completing an entire parent task.

## Tasks

- [x] 0.0 Create feature branch
  - [x] 0.1 Create and checkout a new branch with `git checkout -b feature/template-explanations`

- [x] 1.0 Create explanation generation package
  - [x] 1.1 Create directory `internal/explanations/`
  - [x] 1.2 Create `internal/explanations/generator.go` with package declaration, imports (models package), and empty Generate function signature
  - [x] 1.3 Create `internal/explanations/templates.go` with package declaration and empty SyntacticRoleTemplate function
  - [x] 1.4 Create `internal/explanations/translation.go` with package declaration and empty GenerateTranslation function
  - [x] 1.5 Create `internal/explanations/morphology.go` with package declaration and empty FormatMorphology function
  - [x] 1.6 Create `internal/explanations/generator_test.go` with package declaration and imports

- [x] 2.0 Implement core template logic
  - [x] 2.1 Implement `SyntacticRoleTemplate(contextType, caseType string, prep *string) string` in templates.go with switch statements for all context types (direct_object, possession, preposition)
  - [x] 2.2 Add preposition case map to templates.go with entries for σε, από, για, με, χωρίς, μετά, πριν and their expected cases
  - [x] 2.3 Implement `FormatMorphology(noun *models.Noun, caseType, number string) string` in morphology.go to build nominative→target transformation string using noun fields
  - [x] 2.4 Implement `GenerateTranslation(englishPrompt, greekSentence string) string` in translation.go to find "___" placeholder, extract English from prompt parentheses, and construct translation
  - [x] 2.5 Implement `Generate(sentence *models.Sentence, noun *models.Noun) (*models.Explanation, error)` in generator.go to call all helper functions and return populated Explanation struct

- [x] 3.0 Remove database explanation storage
  - [x] 3.1 Create `internal/storage/migrations/002_drop_explanations.sql` with `DROP TABLE IF EXISTS explanations;`
  - [x] 3.2 Update `internal/storage/migrations.go` to embed and execute the new migration file
  - [x] 3.3 Remove `CreateExplanation(explanation *models.Explanation) error` method signature from Repository interface in repository.go
  - [x] 3.4 Remove `GetExplanationBySentenceID(sentenceID int64) (*models.Explanation, error)` method signature from Repository interface in repository.go
  - [x] 3.5 Remove `CreateExplanation` implementation method from SQLiteRepository struct in repository.go
  - [x] 3.6 Remove `GetExplanationBySentenceID` implementation method from SQLiteRepository struct in repository.go
  - [x] 3.7 Remove `TestCreateAndGetExplanation` and all explanation-related test functions from repository_test.go

- [x] 4.0 Update practice TUI to use template generator
  - [x] 4.1 Add import for `github.com/gataky/greekmaster/internal/explanations` package in practice.go
  - [x] 4.2 In practice.go Update() method feedback state, add `noun, err := m.repo.GetNoun(m.currentSentence.NounID)` call before explanation generation
  - [x] 4.3 Replace `m.repo.GetExplanationBySentenceID(m.currentSentence.ID)` with `explanations.Generate(m.currentSentence, noun)` in practice.go
  - [x] 4.4 Update error handling to check both GetNoun and Generate errors, setting m.err appropriately

- [x] 5.0 Remove AI explanation generation
  - [x] 5.1 Remove `GenerateExplanations(sentences []SentenceResponse) ([]ExplanationResponse, error)` method from ClaudeClient in client.go
  - [x] 5.2 Remove `ExplanationResponse` type definition from client.go
  - [x] 5.3 Remove `GenerateExplanationPrompt(greekSentence, correctAnswer string) string` function from prompts.go
  - [x] 5.4 Remove all `GenerateExplanations()` calls from processor.go (around line 141-149)
  - [x] 5.5 Remove all explanation storage logic from processor.go (loop storing explanations around line 172-183)
  - [x] 5.6 Remove all explanation generation logic from add.go (around line 142-148)
  - [x] 5.7 Remove all explanation storage logic from add.go (loop storing explanations around line 170-181)

- [x] 6.0 Remove --skip-explanations flag
  - [x] 6.1 Remove `skipExplanations bool` variable declaration from import.go NewImportCmd function
  - [x] 6.2 Remove `cmd.Flags().BoolVar(&skipExplanations, "skip-explanations", ...)` line from import.go
  - [x] 6.3 Remove `processor.SetSkipExplanations(skipExplanations)` call from import.go
  - [x] 6.4 Remove `skipExplanations bool` variable declaration from add.go NewAddCmd function
  - [x] 6.5 Remove `cmd.Flags().BoolVar(&skipExplanations, "skip-explanations", ...)` line from add.go
  - [x] 6.6 Remove `SetSkipExplanations(skip bool)` method from processor.go
  - [x] 6.7 Remove `skipExplanations bool` field from ImportProcessor struct in processor.go
  - [x] 6.8 Remove all conditional logic checking `p.skipExplanations` from processor.go and add.go

- [x] 7.0 Implement comprehensive tests
  - [x] 7.1 Implement `TestGenerate` in generator_test.go with table-driven tests covering all combinations: 3 case types × 2 numbers × 3 context types, plus preposition variations
  - [x] 7.2 Implement `TestSyntacticRoleTemplate` in generator_test.go testing all context types (direct_object, possession) and all prepositions (σε, από, για, etc.)
  - [x] 7.3 Implement `TestGenerateTranslation` in generator_test.go with test cases for various prompt formats, with/without parentheses, edge cases like no blank or multiple blanks
  - [x] 7.4 Implement `TestFormatMorphology` in generator_test.go testing all case/number combinations (accusative singular, genitive plural, etc.) for all genders (masculine, feminine, neuter)
  - [x] 7.5 Run `go test ./internal/explanations -v` and verify all tests pass with clear output
  - [x] 7.6 Run `go test ./...` to verify no regressions in other packages

- [ ] 8.0 Verify and cleanup
  - [ ] 8.1 Run `go build -o greekmaster cmd/greekmaster/main.go` and verify successful compilation
  - [ ] 8.2 Manual test: Run `./greekmaster add` to add a noun interactively and verify no --skip-explanations flag exists and command works
  - [ ] 8.3 Manual test: Run `./greekmaster practice` and complete several questions, verify explanations display with Translation, Syntactic Role, and Morphology fields
  - [ ] 8.4 Manual test: During practice, verify all context types work (direct_object sentences, possession sentences, preposition sentences)
  - [ ] 8.5 Manual test: Verify both singular and plural forms show correct explanations
  - [ ] 8.6 Run `go fmt ./...` to format all code
  - [ ] 8.7 Run `go vet ./...` to check for common mistakes
  - [ ] 8.8 Commit all changes with message: "Replace AI-generated explanations with template-based system - eliminates 12 API calls per noun (92% cost reduction)"
