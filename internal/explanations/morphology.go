package explanations

import (
	"fmt"

	"github.com/gataky/greekmaster/internal/models"
)

// FormatMorphology shows the declension transformation
func FormatMorphology(noun *models.Noun, caseType string, number string) string {
	// Get nominative form (starting point)
	nomArticle := noun.NomSgArticle
	nomNoun := noun.NominativeSg

	// Get target form based on case and number
	var targetArticle, targetNoun string

	switch caseType {
	case "nominative":
		if number == "singular" {
			targetArticle = noun.NomSgArticle
			targetNoun = noun.NominativeSg
		} else {
			targetArticle = noun.NomPlArticle
			targetNoun = noun.NominativePl
		}

	case "genitive":
		if number == "singular" {
			targetArticle = noun.GenSgArticle
			targetNoun = noun.GenitiveSg
		} else {
			targetArticle = noun.GenPlArticle
			targetNoun = noun.GenitivePl
		}

	case "accusative":
		if number == "singular" {
			targetArticle = noun.AccSgArticle
			targetNoun = noun.AccusativeSg
		} else {
			targetArticle = noun.AccPlArticle
			targetNoun = noun.AccusativePl
		}

	default:
		return fmt.Sprintf("%s %s", nomArticle, nomNoun)
	}

	return fmt.Sprintf("%s %s → %s %s", nomArticle, nomNoun, targetArticle, targetNoun)
}
