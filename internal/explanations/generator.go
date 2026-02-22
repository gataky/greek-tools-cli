package explanations

import (
	"fmt"

	"github.com/gataky/greekmaster/internal/models"
)

// Generate creates an explanation from sentence and noun data
func Generate(sentence *models.Sentence, noun *models.Noun) (*models.Explanation, error) {
	if sentence == nil {
		return nil, fmt.Errorf("sentence cannot be nil")
	}
	if noun == nil {
		return nil, fmt.Errorf("noun cannot be nil")
	}

	// Generate translation
	translation := GenerateTranslation(sentence.EnglishPrompt, sentence.GreekSentence)

	// Generate syntactic role explanation
	syntacticRole := SyntacticRoleTemplate(sentence.ContextType, sentence.CaseType, sentence.Preposition)

	// Generate morphology transformation
	morphology := FormatMorphology(noun, sentence.CaseType, sentence.Number)

	return &models.Explanation{
		Translation:   translation,
		SyntacticRole: syntacticRole,
		Morphology:    morphology,
	}, nil
}
