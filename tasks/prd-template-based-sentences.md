# PRD: Template-Based Sentence System

## Introduction/Overview

Replace the current duplicative sentence storage model with a template-based system to eliminate redundancy and reduce database size. Currently, the application stores ~11 practice sentences per noun (2,291 total sentences for 204 nouns), resulting in significant duplication where similar sentence patterns are repeated for every noun. This refactor will consolidate these patterns into ~100 reusable templates that dynamically substitute noun forms at runtime, reducing the sentences table by ~80% while maintaining functional parity.

The system will migrate from storing pre-generated sentences (each with `english_prompt`, `greek_sentence`, `correct_answer`) to storing templates with placeholder-based patterns that reference noun declension data from the `nouns` table.

## Goals

### Functional Goals
1. Maintain exact feature parity with current practice session functionality
2. Preserve all existing practice sentences through pattern analysis and template consolidation
3. Support the same difficulty progression (beginner/intermediate/advanced mapping to phases 1/2/3)
4. Maintain the same variety of sentence structures (~11 patterns per noun across case/number/context combinations)
5. Continue supporting plural inclusion filtering and case-specific practice

### Technical Goals
1. Reduce database size by at least 80% (from 2,291 sentences to ~100 templates)
2. Eliminate AI API calls during noun import for sentence generation
3. Implement auto-migration script that converts existing 204 nouns to template system without data loss
4. Maintain sub-100ms query performance for practice session initialization
5. Remove all AI sentence generation code from the import pipeline

## Tech Stack & Architecture

### Languages & Frameworks
- **Go 1.25**: Core application language (existing)
- **SQLite 3**: Database storage via `mattn/go-sqlite3` driver (existing)
- **sqlx**: SQL extension library for named parameters and struct scanning (existing)

### Existing Patterns to Follow
- **Repository Pattern**: Follow `internal/storage/repository.go` interface pattern
  - Add new methods to `Repository` interface
  - Implement in `SQLiteRepository` struct
- **Model Structs**: Create new model in `internal/models/` following existing `Noun`, `Sentence` patterns
  - Use struct tags for `db` field mapping
  - Include `time.Time` fields for `created_at`
- **Migrations**: Add migration in `internal/storage/migrations.go` following existing pattern
  - Use sequential numbering (next: `003_create_templates.sql`)
  - Include `CREATE TABLE IF NOT EXISTS` and indexes
- **Testing**: Follow `internal/storage/repository_test.go` test structure
  - Use in-memory SQLite (`:memory:`) for tests
  - Test CRUD operations and edge cases

### Architectural Pattern
- **Template Substitution at Runtime**:
  - Templates stored in DB as strings with placeholders
  - Go code performs string substitution when generating practice questions
  - No complex template engine needed—simple `strings.Replace()` or equivalent
- **Data Storage**: SQLite with new `sentence_templates` table
- **Backward Compatibility**: During migration, analyze existing sentences to extract patterns, then delete redundant data

### External Services/APIs
- **Remove**: AI sentence generation (`internal/ai/prompts.go` sentence generation functions)
- **Keep**: AI declension generation (still needed for noun import)

## Functional Requirements

### FR1: Template Storage
The system must store sentence templates with the following properties:
- Unique template ID
- English prompt template (e.g., "I see {noun}")
- Greek sentence template (e.g., "Βλέπω {article} {noun_form}")
- Placeholder mappings specifying which noun fields to use:
  - `{article}` → which article field (e.g., `acc_sg_article`)
  - `{noun_form}` → which noun form field (e.g., `accusative_sg`)
- Case type (nominative/genitive/accusative)
- Number (singular/plural)
- Difficulty phase (1/2/3)
- Context type (direct_object/possession/preposition)
- Optional preposition field (for preposition context types)

### FR2: Template Substitution
The system must generate practice sentences by:
1. Selecting a random template matching session filters (phase, number)
2. Selecting a random noun from the `nouns` table
3. Substituting placeholders in both English and Greek templates:
   - Replace `{noun}` in English template with `noun.English`
   - Replace `{article}` in Greek template with appropriate article from noun (e.g., `noun.AccSgArticle`)
   - Replace `{noun_form}` in Greek template with appropriate declension (e.g., `noun.AccusativeSg`)
4. Generating the correct answer by combining article + noun form

### FR3: Migration from Existing Sentences
The system must provide an auto-migration script that:
1. Analyzes existing 2,291 sentences to identify unique patterns
2. Groups sentences by:
   - Sentence structure (ignoring noun-specific content)
   - Case type, number, difficulty phase, context type, preposition
3. Creates templates for each unique pattern (target: ~100 templates)
4. Validates that each existing sentence can be regenerated from templates
5. Deletes migrated sentences after validation succeeds
6. Preserves migration metadata (timestamp, sentence count migrated)

### FR4: Practice Session Integration
The system must modify practice session logic to:
1. Replace `GetRandomSentences()` repository method with new `GeneratePracticeSentences(phase, number, limit)` method
2. Generate sentences at runtime using template substitution
3. Return same `Sentence` struct format for backward compatibility with TUI
4. Support same filtering: difficulty phase, singular/plural inclusion
5. Maintain randomization: both template and noun selection must be random

### FR5: Template Variety
The system must create ~100 templates covering:
- **Difficulty Phase 1 (Beginner - Accusative focus)**: ~35 templates
  - Direct object patterns: "I see {noun}", "She wants {noun}", "They have {noun}"
  - Preposition patterns with accusative: "He is speaking to {noun}", "The letter came from {noun}"
- **Difficulty Phase 2 (Intermediate - Genitive focus)**: ~35 templates
  - Possession patterns: "The bag belongs to {noun}", "The opinion of {noun} matters"
  - Genitive plural patterns: "The names of {noun} are on the list"
- **Difficulty Phase 3 (Advanced - Mixed cases)**: ~30 templates
  - Complex preposition contexts with varied cases
  - Mixed nominative/genitive/accusative patterns
  - Rare prepositions (χωρίς, μετά, πριν)

### FR6: Plural Handling
The system must handle plural-specific patterns:
- Most templates work for both singular and plural (just swap noun form)
- Some templates are plural-specific (e.g., "The names of {noun} are on the list" only makes sense for plural)
- Templates have a `number` field: "singular", "plural", or "both"
- Practice session respects user's plural inclusion setting

### FR7: Preposition Support
The system must maintain preposition functionality:
- Templates with `context_type = "preposition"` must include `preposition` field
- Support existing prepositions: σε, από, για, με, χωρίς, μετά, πριν
- Preposition maps to correct case via existing `prepositionCaseMap` in `internal/explanations/templates.go`
- Greek templates include preposition: "Μιλάει σε {article} {noun_form}"

### FR8: Remove AI Sentence Generation
The system must remove:
- `GenerateSentences()` function from `internal/ai/client.go`
- Sentence generation API calls from `internal/importer/processor.go`
- Sentence generation prompts from `internal/ai/prompts.go`
- Related error handling and retry logic for sentence generation

### FR9: Validation and Error Handling
The system must validate:
- Template placeholders match available noun fields
- Each template can generate valid Greek sentences for all noun genders
- Migration script validates 100% of sentences can be regenerated
- Runtime substitution handles nil/empty values gracefully
- Appropriate errors when no templates exist for requested filters

## Technical Specifications

### Data Models/Schema

#### New Table: sentence_templates
```sql
CREATE TABLE sentence_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    english_template TEXT NOT NULL,
    greek_template TEXT NOT NULL,
    article_field TEXT NOT NULL,        -- e.g., "acc_sg_article"
    noun_form_field TEXT NOT NULL,      -- e.g., "accusative_sg"
    case_type TEXT NOT NULL CHECK(case_type IN ('nominative', 'genitive', 'accusative')),
    number TEXT NOT NULL CHECK(number IN ('singular', 'plural', 'both')),
    difficulty_phase INTEGER NOT NULL CHECK(difficulty_phase IN (1, 2, 3)),
    context_type TEXT NOT NULL CHECK(context_type IN ('direct_object', 'possession', 'preposition')),
    preposition TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_templates_difficulty ON sentence_templates(difficulty_phase, number);
```

#### Updated Model: internal/models/template.go
```go
package models

import "time"

// SentenceTemplate represents a reusable sentence pattern
type SentenceTemplate struct {
    ID              int64     `db:"id"`
    EnglishTemplate string    `db:"english_template"`
    GreekTemplate   string    `db:"greek_template"`
    ArticleField    string    `db:"article_field"`      // Which article to use
    NounFormField   string    `db:"noun_form_field"`    // Which declension to use
    CaseType        string    `db:"case_type"`
    Number          string    `db:"number"`
    DifficultyPhase int       `db:"difficulty_phase"`
    ContextType     string    `db:"context_type"`
    Preposition     *string   `db:"preposition"`
    CreatedAt       time.Time `db:"created_at"`
}
```

### API Contracts (Repository Interface)

#### New Repository Methods
```go
// Template operations
CreateTemplate(template *models.SentenceTemplate) error
GetTemplate(id int64) (*models.SentenceTemplate, error)
ListTemplates() ([]*models.SentenceTemplate, error)
GetRandomTemplates(phase int, number string, limit int) ([]*models.SentenceTemplate, error)

// Template-based sentence generation
GeneratePracticeSentences(phase int, number string, limit int) ([]*models.Sentence, error)

// Migration support
AnalyzeExistingSentences() (map[string][]*models.Sentence, error) // Groups by pattern
DeleteSentencesByIDs(ids []int64) error
```

### Component/Module Structure

#### New File: internal/storage/templates.go
```go
// Template CRUD operations
// Implement: CreateTemplate, GetTemplate, ListTemplates, GetRandomTemplates
```

#### New File: internal/storage/template_generator.go
```go
// Runtime sentence generation from templates
// Core function: GeneratePracticeSentences(phase, number, limit)
// Helper function: substituteTemplate(template, noun) -> Sentence
```

#### New File: internal/storage/migration_templates.go
```go
// Migration logic
// Function: MigrateToTemplates() error
// Function: analyzeSentencePatterns(sentences) -> []SentenceTemplate
// Function: validateMigration() error
```

#### Modified File: internal/importer/processor.go
```go
// Remove: AI sentence generation calls
// Keep: AI declension generation
// Simplify: ProcessImport now only creates nouns, no sentences
```

#### New File: cmd/migrate.go (or add to main.go)
```go
// CLI command: greekmaster migrate-to-templates
// Runs the one-time migration script
// Shows progress and validation results
```

### Integration Points

#### TUI Practice Session (internal/tui/practice.go)
- **Current**: Calls `repo.GetRandomSentences(phase, numberFilter, limit)`
- **New**: Calls `repo.GeneratePracticeSentences(phase, numberFilter, limit)`
- **Change**: Single line replacement in `NewPracticeModel` function (~line 56)
- **Impact**: Returns same `[]*models.Sentence` type, no downstream changes needed

#### Explanation Generator (internal/explanations/)
- **No changes needed**: Templates preserve `case_type`, `context_type`, `preposition` fields
- **Validation**: Ensure generated sentences work with existing `generator.go` functions

### Complex Logic: Template Substitution Algorithm

```go
// Pseudo-code for substituteTemplate function
func substituteTemplate(template *SentenceTemplate, noun *Noun) (*Sentence, error) {
    // 1. Get the article value from noun using reflection or map lookup
    article := getFieldValue(noun, template.ArticleField) // e.g., noun.AccSgArticle

    // 2. Get the noun form value
    nounForm := getFieldValue(noun, template.NounFormField) // e.g., noun.AccusativeSg

    // 3. Substitute English template
    englishPrompt := strings.Replace(template.EnglishTemplate, "{noun}", noun.English, -1)

    // 4. Substitute Greek template
    greekSentence := template.GreekTemplate
    greekSentence = strings.Replace(greekSentence, "{article}", article, -1)
    greekSentence = strings.Replace(greekSentence, "{noun_form}", nounForm, -1)

    // 5. Generate correct answer (article + noun form)
    correctAnswer := article + " " + nounForm

    // 6. Create Sentence struct
    return &Sentence{
        NounID:          noun.ID,
        EnglishPrompt:   englishPrompt,
        GreekSentence:   greekSentence,
        CorrectAnswer:   correctAnswer,
        CaseType:        template.CaseType,
        Number:          template.Number,
        DifficultyPhase: template.DifficultyPhase,
        ContextType:     template.ContextType,
        Preposition:     template.Preposition,
    }, nil
}
```

### Complex Logic: Migration Pattern Analysis

```go
// Pseudo-code for analyzing existing sentences to extract patterns
func analyzeSentencePatterns(sentences []*Sentence) ([]SentenceTemplate, error) {
    patternMap := make(map[string]*SentenceTemplate)

    for _, sentence := range sentences {
        // 1. Get the noun used in this sentence
        noun := getNoun(sentence.NounID)

        // 2. Detect which article and noun form were used by comparing correct_answer
        // Example: if correct_answer is "τον ενήλικο", detect it's acc_sg_article + accusative_sg
        articleField, nounFormField := detectFields(sentence, noun)

        // 3. Create template by replacing noun-specific content with placeholders
        englishTemplate := strings.Replace(sentence.EnglishPrompt, "(the " + noun.English + ")", "{noun}", -1)
        greekTemplate := createGreekTemplate(sentence.GreekSentence, articleField, nounFormField, noun)

        // 4. Create a unique key for this pattern
        patternKey := fmt.Sprintf("%s|%s|%s|%s|%d|%s",
            englishTemplate, greekTemplate, sentence.CaseType,
            sentence.Number, sentence.DifficultyPhase, sentence.ContextType)

        // 5. Store pattern (deduplicate)
        if _, exists := patternMap[patternKey]; !exists {
            patternMap[patternKey] = &SentenceTemplate{
                EnglishTemplate: englishTemplate,
                GreekTemplate:   greekTemplate,
                ArticleField:    articleField,
                NounFormField:   nounFormField,
                CaseType:        sentence.CaseType,
                Number:          sentence.Number,
                DifficultyPhase: sentence.DifficultyPhase,
                ContextType:     sentence.ContextType,
                Preposition:     sentence.Preposition,
            }
        }
    }

    // 6. Convert map to slice
    templates := make([]SentenceTemplate, 0, len(patternMap))
    for _, template := range patternMap {
        templates = append(templates, *template)
    }

    return templates, nil
}
```

### Security Requirements
- **Input Validation**: Template placeholders must be validated against allowed field names to prevent SQL injection or field access errors
- **SQL Injection Prevention**: Use parameterized queries (existing pattern with sqlx named parameters)
- **Data Integrity**: Foreign key constraints maintained (templates don't reference specific nouns, so no FK needed)

### Performance Requirements
- **Practice Session Startup**: Must remain under 100ms for loading 10-20 practice sentences
- **Template Substitution**: Single substitution operation should complete in <1ms
- **Migration Script**: Should complete full migration of 2,291 sentences in under 60 seconds
- **Database Size**: Target <500KB for templates table (vs current ~2MB for sentences)

## Non-Goals (Out of Scope)

1. **Advanced Template Engine**: Not implementing Jinja2-like template syntax with conditionals/loops
2. **User-Editable Templates**: Users cannot add/edit templates through CLI (hardcoded in migration)
3. **Multi-Language Support**: Templates remain Greek-English only
4. **AI-Generated Templates**: Templates are hardcoded/extracted from existing data, not AI-generated
5. **Dynamic Grammar Rules**: Not implementing full Greek morphology engine (no article fusion like σε+τον=στον)
6. **Template Versioning**: No version history or rollback for template changes
7. **Template Analytics**: Not tracking which templates are most effective for learning
8. **Backwards Compatibility Layer**: Old sentence-based system is fully removed after migration
9. **Gradual Migration**: All-or-nothing migration (cannot partially migrate)

## Testing Requirements

### Unit Tests

#### Template Repository Tests (internal/storage/templates_test.go)
- Test `CreateTemplate()`: Insert valid template, verify ID assigned
- Test `GetTemplate()`: Retrieve by ID, handle not found error
- Test `ListTemplates()`: Retrieve all templates, verify ordering
- Test `GetRandomTemplates()`: Filter by phase/number, verify randomization
- Test constraint violations: invalid case_type, invalid number, invalid phase

#### Template Generator Tests (internal/storage/template_generator_test.go)
- Test `substituteTemplate()`: Valid noun + template → correct Sentence struct
- Test placeholder substitution: All placeholders replaced, no remaining `{}`
- Test article/noun form field resolution: Correct fields fetched from noun
- Test nil handling: Graceful handling of nil preposition
- Test `GeneratePracticeSentences()`: Returns requested count, correct filtering

#### Migration Tests (internal/storage/migration_templates_test.go)
- Test `analyzeSentencePatterns()`: Consolidates duplicates correctly
- Test pattern detection: Correctly identifies article and noun form fields
- Test template extraction: Placeholders inserted at correct positions
- Test validation: All original sentences can be regenerated from templates
- Test with different noun genders: Templates work for masculine/feminine/neuter

### Integration Tests

#### End-to-End Migration Test
1. Set up test database with sample nouns and sentences
2. Run migration script
3. Verify templates created (count matches expected unique patterns)
4. Verify sentences deleted
5. Generate practice sentences using templates
6. Compare generated sentences to original sentences (should match or be equivalent)

#### Practice Session Integration Test
1. Initialize repository with templates
2. Create practice session with various configs (beginner/singular, advanced/plural)
3. Generate 100 practice sentences
4. Verify all sentences are valid Greek
5. Verify filtering works (correct phase, correct number)
6. Verify explanations still generate correctly

### Test Coverage Expectations
- **Repository methods**: 100% coverage (CRUD operations are critical)
- **Template substitution logic**: 100% coverage (core feature)
- **Migration logic**: 90%+ coverage (complex logic, edge cases)
- **Overall**: Target 85%+ test coverage for new code

## Success Metrics

### Functional Success Criteria
1. **✅ Zero functionality regression**: All existing practice session features work identically
2. **✅ Migration completes successfully**: 100% of 2,291 sentences converted to templates without data loss
3. **✅ Template variety maintained**: ~100 unique templates created, covering all case/number/difficulty/context combinations
4. **✅ Explanations still work**: Grammar explanations generate correctly for template-based sentences
5. **✅ Filtering works**: Difficulty and plural inclusion filters produce correct sentence sets

### Technical Success Criteria
1. **✅ Database size reduction**: Sentences table reduced by >80% (from 2,291 rows to ~100 template rows)
2. **✅ Performance maintained**: Practice session startup remains <100ms
3. **✅ All tests pass**: Unit and integration tests at 85%+ coverage
4. **✅ AI code removed**: No AI API calls during noun import, sentence generation code deleted
5. **✅ Import simplified**: Noun import process reduced by ~50% execution time (no sentence generation)

### Verification Steps
1. Run full test suite: `make test` passes with 85%+ coverage
2. Manual practice session test: Start sessions at all difficulty levels, verify sentences are natural and correct
3. Database size check: `du -h ~/.greekmaster/greekmaster.db` shows <1MB total size
4. Performance test: Measure practice session startup time with `time` command
5. Import test: Import a new noun, verify no AI sentence generation calls made

## Open Questions

1. **Template seeding**: Should we pre-seed the 100 templates in a migration file, or generate them dynamically during migration? *(Recommendation: Pre-seed for deterministic behavior)*

2. **Template IDs**: Should template IDs be stable/predictable (e.g., 1-100) or auto-incremented? *(Recommendation: Auto-increment for flexibility)*

3. **Sentence struct backwards compatibility**: Should we keep returning the `Sentence` struct with a `nil` ID (since sentences aren't stored anymore), or create a new `GeneratedSentence` struct? *(Recommendation: Keep `Sentence` struct, set ID to 0 or template ID)*

4. **Migration rollback**: Should we implement a rollback mechanism if migration fails halfway? *(Recommendation: Yes, use database transactions)*

5. **Template validation tooling**: Should we add a CLI command to validate all templates generate correct Greek? *(Recommendation: Yes, useful for debugging)*

6. **Future noun additions**: After migration, should new nouns added via `greekmaster add` still use AI, or skip sentence generation entirely? *(Recommendation: Skip, templates cover all nouns)*

7. **Gender-specific templates**: Do we need separate templates for masculine/feminine/neuter, or can one template work for all genders via field substitution? *(Recommendation: One template works for all, article fields handle gender)*

8. **Preposition template organization**: Should preposition-based templates be in a separate group for easier management? *(Recommendation: No, same table, filter by context_type)*
