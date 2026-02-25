package ai

import "fmt"

// GenerateDeclensionPrompt creates a prompt for generating all declined forms of a Greek noun
func GenerateDeclensionPrompt(greek, english, gender string) string {
	return fmt.Sprintf(`You are a Modern Greek grammar expert. Given the noun '%s' (%s) with gender '%s', provide all declined forms with their articles in this JSON format:
{
  "nominative_sg": "...", "nom_sg_article": "...",
  "genitive_sg": "...", "gen_sg_article": "...",
  "accusative_sg": "...", "acc_sg_article": "...",
  "nominative_pl": "...", "nom_pl_article": "...",
  "genitive_pl": "...", "gen_pl_article": "...",
  "accusative_pl": "...", "acc_pl_article": "..."
}
Return only valid JSON, no explanation.`, greek, english, gender)
}


