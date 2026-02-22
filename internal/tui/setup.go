package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gataky/greekmaster/internal/models"
)

// SetupModel represents the session setup screen
type SetupModel struct {
	step          int // 0=difficulty, 1=plural, 2=session type
	difficulty    string
	includePlural bool
	questionCount int
	quitting      bool
	complete      bool
	err           error
}

// NewSetupModel creates a new setup model
func NewSetupModel() SetupModel {
	return SetupModel{
		step:          0,
		difficulty:    "",
		includePlural: false,
		questionCount: 0,
	}
}

func (m SetupModel) Init() tea.Cmd {
	return nil
}

func (m SetupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "1", "2", "3":
			if m.step == 0 {
				// Difficulty selection
				switch msg.String() {
				case "1":
					m.difficulty = "beginner"
				case "2":
					m.difficulty = "intermediate"
				case "3":
					m.difficulty = "advanced"
				}
				m.step = 1
			} else if m.step == 2 {
				// Session type selection
				switch msg.String() {
				case "1":
					m.questionCount = 10
				case "2":
					m.questionCount = 25
				case "3":
					m.questionCount = 50
				}
				// Setup complete
				m.complete = true
				return m, tea.Quit
			}

		case "4":
			if m.step == 2 {
				// Endless mode
				m.questionCount = 0
				m.complete = true
				return m, tea.Quit
			}

		case "y", "Y":
			if m.step == 1 {
				m.includePlural = true
				m.step = 2
			}

		case "n", "N":
			if m.step == 1 {
				m.includePlural = false
				m.step = 2
			}
		}
	}

	return m, nil
}

func (m SetupModel) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		MarginBottom(1)

	questionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		MarginBottom(1)

	optionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	s.WriteString(titleStyle.Render("Greek Case Master - Session Setup"))
	s.WriteString("\n\n")

	switch m.step {
	case 0:
		s.WriteString(questionStyle.Render("Select difficulty level:"))
		s.WriteString("\n\n")
		s.WriteString(optionStyle.Render("  1. Beginner     - Focus on accusative (direct objects)"))
		s.WriteString("\n")
		s.WriteString(optionStyle.Render("  2. Intermediate - Focus on genitive (possession)"))
		s.WriteString("\n")
		s.WriteString(optionStyle.Render("  3. Advanced     - Mixed cases with prepositions"))
		s.WriteString("\n\n")

	case 1:
		s.WriteString(questionStyle.Render(fmt.Sprintf("Difficulty: %s", m.difficulty)))
		s.WriteString("\n\n")
		s.WriteString(questionStyle.Render("Include plural forms? (Y/n)"))
		s.WriteString("\n\n")

	case 2:
		s.WriteString(questionStyle.Render(fmt.Sprintf("Difficulty: %s", m.difficulty)))
		s.WriteString("\n")
		s.WriteString(questionStyle.Render(fmt.Sprintf("Plural forms: %v", m.includePlural)))
		s.WriteString("\n\n")
		s.WriteString(questionStyle.Render("Select session type:"))
		s.WriteString("\n\n")
		s.WriteString(optionStyle.Render("  1. Quick (10 questions)"))
		s.WriteString("\n")
		s.WriteString(optionStyle.Render("  2. Standard (25 questions)"))
		s.WriteString("\n")
		s.WriteString(optionStyle.Render("  3. Long (50 questions)"))
		s.WriteString("\n")
		s.WriteString(optionStyle.Render("  4. Endless (practice until quit)"))
		s.WriteString("\n\n")
	}

	s.WriteString(optionStyle.Render("\n[q] Quit"))

	return s.String()
}

// GetConfig returns the session configuration if setup is complete
func (m SetupModel) GetConfig() (models.SessionConfig, bool) {
	if !m.complete {
		return models.SessionConfig{}, false
	}
	return models.SessionConfig{
		DifficultyLevel: m.difficulty,
		IncludePlural:   m.includePlural,
		QuestionCount:   m.questionCount,
	}, true
}

// SessionConfigMsg is sent when setup is complete
type SessionConfigMsg struct {
	Config models.SessionConfig
}
