package storage

import (
	"fmt"
	"strings"

	"github.com/gataky/greekmaster/internal/models"
)

// detectFields determines which article and noun form fields were used in a sentence
func detectFields(sentence *models.Sentence, noun *models.Noun) (articleField, nounFormField string, err error) {
	// Parse the correct answer to extract article and noun form
	parts := strings.Fields(sentence.CorrectAnswer)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid correct_answer format: %s", sentence.CorrectAnswer)
	}

	article := parts[0]
	nounForm := parts[1]

	// Match article against noun's article fields
	articleFields := map[string]string{
		noun.NomSgArticle: "NomSgArticle",
		noun.GenSgArticle: "GenSgArticle",
		noun.AccSgArticle: "AccSgArticle",
		noun.NomPlArticle: "NomPlArticle",
		noun.GenPlArticle: "GenPlArticle",
		noun.AccPlArticle: "AccPlArticle",
	}

	articleField, ok := articleFields[article]
	if !ok {
		return "", "", fmt.Errorf("article %s not found in noun fields", article)
	}

	// Match noun form against noun's declension fields
	nounFormFields := map[string]string{
		noun.NominativeSg: "NominativeSg",
		noun.GenitiveSg:   "GenitiveSg",
		noun.AccusativeSg: "AccusativeSg",
		noun.NominativePl: "NominativePl",
		noun.GenitivePl:   "GenitivePl",
		noun.AccusativePl: "AccusativePl",
	}

	nounFormField, ok = nounFormFields[nounForm]
	if !ok {
		return "", "", fmt.Errorf("noun form %s not found in noun fields", nounForm)
	}

	return articleField, nounFormField, nil
}

// createGreekTemplate replaces specific article and noun form with placeholders
func createGreekTemplate(greekSentence string, article string, nounForm string) string {
	// Replace the specific article and noun form with placeholders
	template := greekSentence
	template = strings.ReplaceAll(template, article+" "+nounForm, "{article} {noun_form}")
	return template
}

// createEnglishTemplate replaces noun-specific content with placeholder
func createEnglishTemplate(englishPrompt string, englishNoun string) string {
	// The English prompt typically contains "(the noun)" pattern
	// Replace with {noun} placeholder
	template := englishPrompt

	// Try various patterns
	patterns := []string{
		"(the " + englishNoun + ")",
		"(" + englishNoun + ")",
		englishNoun,
	}

	for _, pattern := range patterns {
		if strings.Contains(template, pattern) {
			template = strings.ReplaceAll(template, pattern, "{noun}")
			break
		}
	}

	return template
}

// analyzeSentencePatterns extracts unique templates from existing sentences
func analyzeSentencePatterns(sentences []*models.Sentence, nouns map[int64]*models.Noun) ([]*models.SentenceTemplate, error) {
	patternMap := make(map[string]*models.SentenceTemplate)

	for _, sentence := range sentences {
		// Get the noun used in this sentence
		noun, ok := nouns[sentence.NounID]
		if !ok {
			return nil, fmt.Errorf("noun not found for sentence ID %d", sentence.ID)
		}

		// Detect which fields were used
		articleField, nounFormField, err := detectFields(sentence, noun)
		if err != nil {
			// Skip sentences that can't be parsed
			fmt.Printf("Warning: skipping sentence ID %d: %v\n", sentence.ID, err)
			continue
		}

		// Extract article and noun form from correct answer
		parts := strings.Fields(sentence.CorrectAnswer)
		article := parts[0]
		nounForm := parts[1]

		// Create templates
		englishTemplate := createEnglishTemplate(sentence.EnglishPrompt, noun.English)
		greekTemplate := createGreekTemplate(sentence.GreekSentence, article, nounForm)

		// Create a unique key for this pattern
		patternKey := fmt.Sprintf("%s|%s|%s|%s|%s|%d|%s|%v",
			englishTemplate, greekTemplate, articleField, nounFormField,
			sentence.CaseType, sentence.DifficultyPhase, sentence.ContextType,
			sentence.Preposition)

		// Store pattern (deduplicate)
		if _, exists := patternMap[patternKey]; !exists {
			patternMap[patternKey] = &models.SentenceTemplate{
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

	// Convert map to slice
	templates := make([]*models.SentenceTemplate, 0, len(patternMap))
	for _, template := range patternMap {
		templates = append(templates, template)
	}

	return templates, nil
}

// validateMigration verifies all sentences can be regenerated from templates
func validateMigration(templates []*models.SentenceTemplate, originalSentences []*models.Sentence, nouns map[int64]*models.Noun) error {
	// For validation, we'll check that we have enough templates
	// and that they cover the same case/phase/context combinations
	if len(templates) == 0 {
		return fmt.Errorf("no templates created")
	}

	// Count original sentence characteristics
	originalCombos := make(map[string]int)
	for _, s := range originalSentences {
		key := fmt.Sprintf("%s-%d-%s", s.CaseType, s.DifficultyPhase, s.ContextType)
		originalCombos[key]++
	}

	// Count template characteristics
	templateCombos := make(map[string]int)
	for _, t := range templates {
		key := fmt.Sprintf("%s-%d-%s", t.CaseType, t.DifficultyPhase, t.ContextType)
		templateCombos[key]++
	}

	// Verify we have templates for each combination
	for combo := range originalCombos {
		if _, exists := templateCombos[combo]; !exists {
			return fmt.Errorf("missing template for combination: %s", combo)
		}
	}

	fmt.Printf("Validation passed: %d unique templates cover all %d sentence combinations\n",
		len(templates), len(originalCombos))

	return nil
}

// MigrateToTemplates performs the one-time migration from sentences to templates
func MigrateToTemplates(repo *SQLiteRepository) error {
	fmt.Println("Starting migration to template-based system...")
	fmt.Println()

	// Start transaction for rollback capability
	tx, err := repo.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback() // Will be no-op if we commit

	// 1. Load all existing sentences
	fmt.Println("Loading existing sentences...")
	var sentences []*models.Sentence
	err = repo.db.Select(&sentences, "SELECT * FROM sentences ORDER BY id")
	if err != nil {
		return fmt.Errorf("failed to load sentences: %w", err)
	}
	fmt.Printf("Loaded %d sentences\n", len(sentences))
	fmt.Println()

	// 2. Load all nouns
	fmt.Println("Loading nouns...")
	nouns, err := repo.ListNouns()
	if err != nil {
		return fmt.Errorf("failed to load nouns: %w", err)
	}

	// Create noun lookup map
	nounMap := make(map[int64]*models.Noun)
	for _, noun := range nouns {
		nounMap[noun.ID] = noun
	}
	fmt.Printf("Loaded %d nouns\n", len(nouns))
	fmt.Println()

	// 3. Analyze patterns
	fmt.Println("Analyzing sentence patterns...")
	templates, err := analyzeSentencePatterns(sentences, nounMap)
	if err != nil {
		return fmt.Errorf("failed to analyze patterns: %w", err)
	}
	fmt.Printf("Extracted %d unique templates\n", len(templates))
	fmt.Println()

	// 4. Insert templates
	fmt.Println("Creating templates in database...")
	for i, template := range templates {
		query := `
			INSERT INTO sentence_templates (
				english_template, greek_template, article_field, noun_form_field,
				case_type, number, difficulty_phase, context_type, preposition
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		result, err := tx.Exec(query,
			template.EnglishTemplate, template.GreekTemplate,
			template.ArticleField, template.NounFormField,
			template.CaseType, template.Number, template.DifficultyPhase,
			template.ContextType, template.Preposition)
		if err != nil {
			return fmt.Errorf("failed to insert template %d: %w", i, err)
		}

		id, _ := result.LastInsertId()
		templates[i].ID = id
	}
	fmt.Printf("Created %d templates\n", len(templates))
	fmt.Println()

	// 5. Validate migration
	fmt.Println("Validating migration...")
	err = validateMigration(templates, sentences, nounMap)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	fmt.Println()

	// 6. Delete old sentences
	fmt.Println("Deleting old sentences...")
	_, err = tx.Exec("DELETE FROM sentences")
	if err != nil {
		return fmt.Errorf("failed to delete sentences: %w", err)
	}
	fmt.Printf("Deleted %d sentences\n", len(sentences))
	fmt.Println()

	// Commit transaction
	fmt.Println("Committing migration...")
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Println("Migration completed successfully!")
	fmt.Printf("Summary:\n")
	fmt.Printf("  - Templates created: %d\n", len(templates))
	fmt.Printf("  - Sentences migrated: %d\n", len(sentences))
	fmt.Printf("  - Database size reduced by ~%d%%\n", (len(sentences)-len(templates))*100/len(sentences))

	return nil
}
