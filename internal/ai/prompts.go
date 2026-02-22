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

// GenerateSentencesPrompt creates a prompt for generating practice sentences
func GenerateSentencesPrompt(greek, english, gender string) string {
	return fmt.Sprintf(`Generate 12 practice sentences for learning Greek cases using the noun '%s' (%s, %s). Create:
- 3 sentences with accusative case (direct object after different transitive verbs like βλέπω, ψάχνω, θέλω)
- 3 sentences with genitive case (possession contexts)
- 3 sentences with accusative after prepositions (σε, από)
- 3 sentences with genitive after preposition (για)

Mix singular and plural forms. Return as JSON array with this structure:
[{
  "english_prompt": "I see ___ (the teacher)",
  "greek_sentence": "Βλέπω τον δάσκαλο",
  "correct_answer": "τον δάσκαλο",
  "case_type": "accusative",
  "number": "singular",
  "difficulty_phase": 1,
  "context_type": "direct_object",
  "preposition": null
}, ...]

IMPORTANT: For context_type, use ONLY these exact values:
- "direct_object" - for direct objects (transitive verbs)
- "possession" - for genitive possession contexts
- "preposition" - for any sentence with a preposition

For case_type, use ONLY: "nominative", "genitive", or "accusative"
For number, use ONLY: "singular" or "plural"

Return only valid JSON array, no explanation.`, greek, english, gender)
}

// GenerateExplanationPrompt creates a prompt for generating grammar explanations
func GenerateExplanationPrompt(greekSentence, correctAnswer string) string {
	return fmt.Sprintf(`For the sentence '%s' where the answer is '%s', provide a grammar explanation as the Modern Greek Grammar Analyst. Return JSON:
{
  "translation": "Full English translation",
  "syntactic_role": "Explain why this case is required",
  "morphology": "Explain the form transformation"
}
Be concise and analytical, no conversational filler. Return only valid JSON, no explanation.`, greekSentence, correctAnswer)
}
