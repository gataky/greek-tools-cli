# PRD: Template-Based Explanation System

## 1. Introduction/Overview

Replace the current AI-generated explanation system with a template-based approach that generates grammatical explanations on-demand during practice sessions. This eliminates expensive API calls for explanations (~12 calls per noun) while maintaining educational value through rule-based explanations derived from sentence metadata (case, number, context type, preposition).

## 2. Goals

### Functional Goals
- Generate clear, consistent grammatical explanations without AI API calls
- Maintain the three-field explanation structure: Translation, Syntactic Role, Morphology
- Provide simple, rule-based explanations suitable for learners
- Support all existing context types: direct_object, possession, preposition

### Technical Goals
- Reduce API costs by 92% (from 13 calls/noun to 2 calls/noun)
- Eliminate database storage for explanations
- Maintain sub-50ms explanation generation time
- Zero external dependencies for explanation generation

## 3. Tech Stack & Architecture

### Languages
- Go 1.23+

### Frameworks & Libraries
- Standard library only (no external dependencies for templating)
- Existing: `internal/models` for data structures
- Existing: `internal/storage` repository pattern

### Architectural Pattern
- **Template-based generation**: Pure functions that map sentence metadata → explanation text
- **On-demand generation**: No storage, generate during practice session feedback display
- **Hard-coded rules**: Greek grammar rules embedded in code (no configuration files)

### Existing Patterns to Follow
- Package structure: Follow `internal/ai` pattern (client + functions)
- Error handling: Return `error` from all fallible operations
- Testing: Follow `internal/storage/repository_test.go` pattern with table-driven tests

### Data Storage
- **Remove**: Explanations table and all related storage logic
- **Keep**: Sentences and Nouns tables (source data for templates)

### External Services/APIs
- None (eliminating Claude API dependency for explanations)

## 4. Functional Requirements

### Core Functionality

1. **Template Engine**: Generate explanations from Sentence + Noun data with three components:
   - **Translation**: Template-based English translation from Greek sentence structure
   - **Syntactic Role**: Simple rule explaining why the case is used
   - **Morphology**: Show nominative → target case transformation

2. **Context-Specific Rules**:
   - **direct_object**: "Direct objects use accusative case"
   - **possession**: "Possession requires genitive case"
   - **preposition**: Hard-coded preposition → case rules (σε→accusative, από→genitive, για→genitive, etc.)

3. **Translation Generation**: Parse English prompt and Greek sentence to create natural translation
   - Example: Prompt "I see ___ (the teacher)" + Greek "Βλέπω τον δάσκαλο" → "I see the teacher"

4. **Morphology Display**: Show article + noun transformation
   - Format: `{nominative_article} {nominative_noun} → {target_article} {target_noun}`
   - Example: `ο δάσκαλος → τον δάσκαλο`

### Data Validation
5. All case types must be supported: nominative, genitive, accusative
6. All number types must be supported: singular, plural
7. All context types must be supported: direct_object, possession, preposition
8. Handle nil prepositions gracefully (non-preposition contexts)

### Error Handling
9. Return descriptive errors for unsupported case/context combinations
10. Fall back to minimal explanation if translation parsing fails
11. Never panic - all errors must be handled gracefully

### Integration Requirements
12. Provide interface compatible with existing `*models.Explanation` structure
13. Integrate with practice TUI at feedback display time
14. Remove all database explanation storage logic
15. Remove explanation-related API client code

## 5. Technical Specifications

### Data Models/Schema

**Input**: `*models.Sentence` + `*models.Noun`

**Output**: `*models.Explanation` (generated, not from DB)
```go
type Explanation struct {
    Translation   string // Generated from prompt + Greek sentence
    SyntacticRole string // Rule-based explanation
    Morphology    string // Transformation display
}
```

### Module Structure

**New Package**: `internal/explanations/`

```
internal/explanations/
├── generator.go      // Main generation logic
├── templates.go      // Rule templates for each context type
├── translation.go    // Translation parsing logic
├── morphology.go     // Morphology transformation formatting
└── generator_test.go // Comprehensive tests
```

### Core Functions

**generator.go**:
```go
// Generate creates an explanation from sentence and noun data
func Generate(sentence *models.Sentence, noun *models.Noun) (*models.Explanation, error)
```

**templates.go**:
```go
// SyntacticRoleTemplate returns the rule explanation for a given context
func SyntacticRoleTemplate(contextType string, caseType string, prep *string) string

// Examples:
// ("direct_object", "accusative", nil) → "Direct objects use accusative case"
// ("preposition", "accusative", "σε") → "The preposition 'σε' requires accusative case"
// ("possession", "genitive", nil) → "Possession requires genitive case"
```

**translation.go**:
```go
// GenerateTranslation creates English translation from prompt and Greek sentence
func GenerateTranslation(englishPrompt string, greekSentence string) string

// Logic:
// 1. Find the blank placeholder "___" in englishPrompt
// 2. Extract the article + noun from greekSentence
// 3. Translate to English using noun.English
// 4. Replace "___" with "the {english}"
// 5. Clean up parenthetical hints like "(the teacher)"
```

**morphology.go**:
```go
// FormatMorphology shows the declension transformation
func FormatMorphology(noun *models.Noun, caseType string, number string) string

// Logic:
// 1. Get nominative form: noun.NomSgArticle + " " + noun.NominativeSg
// 2. Get target form based on caseType + number:
//    - "accusative" + "singular" → noun.AccSgArticle + " " + noun.AccusativeSg
//    - "genitive" + "plural" → noun.GenPlArticle + " " + noun.GenitivePl
// 3. Return "{nom} → {target}"
```

### Preposition Rules (Hard-coded)

```go
var prepositionCaseMap = map[string]string{
    "σε":   "accusative",  // to, at, in
    "από":  "genitive",    // from
    "για":  "genitive",    // for
    "με":   "accusative",  // with
    "χωρίς": "accusative", // without
    "μετά": "accusative",  // after
    "πριν": "accusative",  // before
}

func ValidatePrepositionCase(prep string, caseType string) error {
    expected, ok := prepositionCaseMap[prep]
    if !ok {
        return fmt.Errorf("unknown preposition: %s", prep)
    }
    if expected != caseType {
        return fmt.Errorf("preposition '%s' expects %s case, got %s", prep, expected, caseType)
    }
    return nil
}
```

### Integration Changes

**Remove from `internal/storage/repository.go`**:
- `CreateExplanation()` method
- `GetExplanationBySentenceID()` method
- Remove Explanation from Repository interface

**Update `internal/tui/practice.go`**:
```go
// OLD:
m.explanation, err = m.repo.GetExplanationBySentenceID(m.currentSentence.ID)

// NEW:
noun, err := m.repo.GetNoun(m.currentSentence.NounID)
if err != nil {
    m.err = err
} else {
    m.explanation, err = explanations.Generate(m.currentSentence, noun)
    if err != nil {
        m.err = err
    }
}
```

**Remove from `internal/importer/processor.go`**:
- All `GenerateExplanations()` calls
- All `CreateExplanation()` calls
- Remove explanation tracking in progress display

**Remove from `internal/commands/add.go`**:
- All explanation generation logic
- All explanation storage logic

**Remove from `internal/ai/client.go`**:
- `GenerateExplanations()` method
- `GenerateExplanationPrompt()` from prompts.go
- `ExplanationResponse` type

**Remove flags**:
- `--skip-explanations` from import and add commands (no longer needed)

### Database Migration

**Drop table**:
```sql
-- migrations/002_drop_explanations.sql
DROP TABLE IF EXISTS explanations;
```

**Update repository initialization**:
- Run migration automatically on next database access
- No data migration needed (regenerate on-demand)

### Security Requirements
- Input sanitization: None needed (all inputs from trusted database)
- No user input processed in explanation generation

### Performance Requirements
- Explanation generation: < 50ms per explanation
- Memory: Single explanation struct ~200 bytes
- No caching needed (generation is cheap)

## 6. Non-Goals (Out of Scope)

- Internationalization (English explanations only)
- User-customizable explanation templates
- Advanced grammatical analysis (declension classes, stem analysis)
- Explanation history or versioning
- A/B testing different explanation styles
- Fallback to AI explanations for edge cases
- Configuration file for grammar rules
- Support for languages other than Modern Greek

## 7. Testing Requirements

### Unit Tests

**`generator_test.go`**:
```go
func TestGenerate(t *testing.T) {
    // Table-driven tests covering:
    // - All case types (nominative, genitive, accusative)
    // - All numbers (singular, plural)
    // - All context types (direct_object, possession, preposition)
    // - All common prepositions
    // - Error cases (nil inputs, invalid combinations)
}

func TestSyntacticRoleTemplate(t *testing.T) {
    // Test all context + case combinations
    // Test preposition-specific rules
}

func TestGenerateTranslation(t *testing.T) {
    // Test various prompt formats
    // Test with/without parenthetical hints
    // Test edge cases (no blank, multiple blanks)
}

func TestFormatMorphology(t *testing.T) {
    // Test all case + number combinations
    // Test all genders (masculine, feminine, neuter)
}
```

### Integration Tests

**`internal/tui/practice_test.go`** (new):
- Test feedback rendering with generated explanations
- Verify explanation display in TUI
- Test error handling when generation fails

### End-to-End Tests

**Manual testing checklist**:
1. Import nouns without `--skip-explanations` flag (flag should be removed)
2. Start practice session
3. Answer question correctly and incorrectly
4. Verify explanation displays three fields
5. Verify explanation accuracy for all context types
6. Test with singular and plural forms
7. Test with all prepositions
8. Verify no errors logged

### Test Coverage Expectations
- `internal/explanations/`: 95%+ coverage
- All exported functions must have tests
- All preposition rules must be tested
- All case/number/context combinations tested

## 8. Success Metrics

### Technical Success Criteria
- Zero API calls for explanation generation
- All existing practice tests still pass
- New explanation tests achieve 95%+ coverage
- Explanation generation < 50ms measured via benchmark tests
- Zero panics or errors during practice sessions
- Database size reduction (explanations table dropped)

### Functional Success Criteria
- Explanations are grammatically accurate for all test cases
- Users can complete practice sessions without errors
- Explanation display is readable and helpful
- All three explanation fields populated correctly

### Performance Success Criteria
- Import time reduced by ~70% (12 fewer API calls per noun)
- Import cost reduced by ~92% (from 13 calls to 2 calls per noun)
- Practice session feedback displays within 100ms

## 9. Open Questions

None - all requirements clarified through Q&A.
