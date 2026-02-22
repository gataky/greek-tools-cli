package explanations

import (
	"strings"
)

// GenerateTranslation creates English translation from prompt and Greek sentence
func GenerateTranslation(englishPrompt string, greekSentence string) string {
	// Find the blank placeholder "___"
	blankIndex := strings.Index(englishPrompt, "___")
	if blankIndex == -1 {
		// No blank found, return the Greek sentence as-is
		return greekSentence
	}

	// Extract the English translation from parentheses (if present)
	// Example: "I see ___ (the teacher)" -> extract "the teacher"
	openParen := strings.Index(englishPrompt, "(")
	closeParen := strings.Index(englishPrompt, ")")

	var replacement string
	if openParen != -1 && closeParen != -1 && closeParen > openParen {
		// Extract text between parentheses
		replacement = strings.TrimSpace(englishPrompt[openParen+1 : closeParen])
	} else {
		// No parentheses, use a generic placeholder
		replacement = "[answer]"
	}

	// Replace the blank with the extracted text
	before := englishPrompt[:blankIndex]
	after := englishPrompt[blankIndex+3:] // Skip "___"

	// Remove the parenthetical hint from the after part
	if openParen != -1 {
		after = strings.TrimSpace(after[:openParen-blankIndex-3])
	}

	translation := strings.TrimSpace(before + replacement + after)

	return translation
}
