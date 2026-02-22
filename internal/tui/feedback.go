package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gataky/greekmaster/internal/models"
)

// RenderFeedback renders the feedback screen after an answer
func RenderFeedback(isCorrect bool, userAnswer string, sentence *models.Sentence, explanation *models.Explanation, terminalWidth int) string {
	var s strings.Builder

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1, 2)

	correctStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("120"))

	incorrectStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("196"))

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Bold(true)

	// Calculate content width based on terminal width
	// Account for border (2 chars), padding (4 chars), and some margin
	contentWidth := terminalWidth - 10
	if contentWidth < 40 {
		contentWidth = 40 // Minimum width
	}
	if contentWidth > 100 {
		contentWidth = 100 // Maximum width for readability
	}

	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Width(contentWidth)

	answerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("120")).
		Bold(true)

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196"))

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		MarginTop(1)

	// Header
	if isCorrect {
		s.WriteString(correctStyle.Render("✓ Correct!"))
	} else {
		s.WriteString(incorrectStyle.Render("✗ Incorrect"))
		s.WriteString("\n\n")
		s.WriteString(labelStyle.Render("You entered: "))
		s.WriteString(errorStyle.Render(userAnswer))
		s.WriteString("\n")
		s.WriteString(labelStyle.Render("Correct answer: "))
		s.WriteString(answerStyle.Render(sentence.CorrectAnswer))
	}

	s.WriteString("\n\n")
	s.WriteString(strings.Repeat("─", 50))
	s.WriteString("\n\n")

	// Show explanation if available
	if explanation != nil {
		s.WriteString(labelStyle.Render("Translation: "))
		s.WriteString(textStyle.Render(explanation.Translation))
		s.WriteString("\n\n")

		s.WriteString(labelStyle.Render("Syntactic Role: "))
		s.WriteString(textStyle.Render(explanation.SyntacticRole))
		s.WriteString("\n\n")

		s.WriteString(labelStyle.Render("Morphology: "))
		s.WriteString(textStyle.Render(explanation.Morphology))
	} else {
		s.WriteString(textStyle.Render("Full Greek sentence: " + sentence.GreekSentence))
	}

	s.WriteString("\n")
	s.WriteString(hintStyle.Render("\n[Press any key to continue]"))

	return borderStyle.Render(s.String())
}
