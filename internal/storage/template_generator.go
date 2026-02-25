package storage

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"

	"github.com/gataky/greekmaster/internal/models"
)

// getFieldValue extracts a field value from a noun struct by field name
func getFieldValue(noun *models.Noun, fieldName string) (string, error) {
	// Use reflection to get the field value
	v := reflect.ValueOf(noun).Elem()
	field := v.FieldByName(fieldName)

	if !field.IsValid() {
		return "", fmt.Errorf("field %s not found in Noun struct", fieldName)
	}

	// Convert to string (all our fields are strings)
	return field.String(), nil
}

// substituteTemplate generates a Sentence from a template and noun
func substituteTemplate(template *models.SentenceTemplate, noun *models.Noun) (*models.Sentence, error) {
	// 1. Get the article value from noun
	article, err := getFieldValue(noun, template.ArticleField)
	if err != nil {
		return nil, fmt.Errorf("failed to get article field: %w", err)
	}

	// 2. Get the noun form value
	nounForm, err := getFieldValue(noun, template.NounFormField)
	if err != nil {
		return nil, fmt.Errorf("failed to get noun form field: %w", err)
	}

	// 3. Substitute English template
	englishPrompt := strings.ReplaceAll(template.EnglishTemplate, "{noun}", noun.English)

	// 4. Substitute Greek template
	greekSentence := template.GreekTemplate
	greekSentence = strings.ReplaceAll(greekSentence, "{article}", article)
	greekSentence = strings.ReplaceAll(greekSentence, "{noun_form}", nounForm)

	// 5. Generate correct answer (article + noun form)
	correctAnswer := article + " " + nounForm

	// 6. Determine the number for the sentence (map 'both' to actual number)
	number := template.Number
	if number == "both" {
		// Infer from the noun form field which number was used
		if strings.Contains(template.NounFormField, "Sg") {
			number = "singular"
		} else if strings.Contains(template.NounFormField, "Pl") {
			number = "plural"
		}
	}

	// 7. Create Sentence struct
	return &models.Sentence{
		NounID:          noun.ID,
		EnglishPrompt:   englishPrompt,
		GreekSentence:   greekSentence,
		CorrectAnswer:   correctAnswer,
		CaseType:        template.CaseType,
		Number:          number,
		DifficultyPhase: template.DifficultyPhase,
		ContextType:     template.ContextType,
		Preposition:     template.Preposition,
	}, nil
}

// GeneratePracticeSentences generates practice sentences from templates
func (r *SQLiteRepository) GeneratePracticeSentences(phase int, number string, limit int) ([]*models.Sentence, error) {
	// 1. Get all nouns
	nouns, err := r.ListNouns()
	if err != nil {
		return nil, fmt.Errorf("failed to get nouns: %w", err)
	}

	if len(nouns) == 0 {
		return nil, fmt.Errorf("no nouns found in database")
	}

	// 2. Get templates matching filters (get more than needed for variety)
	templateLimit := limit * 2
	if templateLimit < 100 {
		templateLimit = 100
	}
	templates, err := r.GetRandomTemplates(phase, number, templateLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to get templates: %w", err)
	}

	if len(templates) == 0 {
		return nil, fmt.Errorf("no templates found for phase %d and number %s", phase, number)
	}

	// 3. Generate sentences by combining templates with random nouns
	sentences := make([]*models.Sentence, 0, limit)
	used := make(map[string]bool) // Track used combinations to avoid duplicates

	// Try to generate the requested number of sentences
	attempts := 0
	maxAttempts := limit * 10 // Prevent infinite loop

	for len(sentences) < limit && attempts < maxAttempts {
		attempts++

		// Pick random template and noun
		template := templates[rand.Intn(len(templates))]
		noun := nouns[rand.Intn(len(nouns))]

		// Create unique key for this combination
		key := fmt.Sprintf("%d-%d", template.ID, noun.ID)
		if used[key] {
			continue // Skip if we've already used this combination
		}

		// Generate sentence
		sentence, err := substituteTemplate(template, noun)
		if err != nil {
			// Skip invalid combinations (e.g., field mismatch)
			continue
		}

		sentences = append(sentences, sentence)
		used[key] = true
	}

	// If we couldn't generate enough unique combinations, that's okay
	// Just return what we have
	return sentences, nil
}
