package tui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gataky/greekmaster/internal/explanations"
	"github.com/gataky/greekmaster/internal/models"
	"github.com/gataky/greekmaster/internal/storage"
)

// PracticeModel represents the practice session
type PracticeModel struct {
	repo            storage.Repository
	config          models.SessionConfig
	sentences       []*models.Sentence
	currentIndex    int
	userInput       string
	state           string // "question", "feedback", "complete"
	correctCount    int
	incorrectCount  int
	currentSentence *models.Sentence
	isCorrect       bool
	explanation     *models.Explanation
	err             error
	rng             *rand.Rand // Random number generator
	width           int        // Terminal width
}

// NewPracticeModel creates a new practice model
func NewPracticeModel(repo storage.Repository, config models.SessionConfig) (*PracticeModel, error) {
	// Map difficulty to phase
	phaseMap := map[string]int{
		"beginner":     1,
		"intermediate": 2,
		"advanced":     3,
	}
	phase := phaseMap[config.DifficultyLevel]

	// Determine number filter
	numberFilter := "singular"
	if config.IncludePlural {
		numberFilter = "" // Empty means both
	}

	// Load sentences from database
	limit := 1000
	if config.QuestionCount > 0 {
		limit = config.QuestionCount * 2 // Get extra for variety
	}

	sentences, err := repo.GetRandomSentences(phase, numberFilter, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to load sentences: %w", err)
	}

	if len(sentences) == 0 {
		return nil, fmt.Errorf("no sentences found for this difficulty level. Please run 'greekmaster import' first")
	}

	// Create random number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Shuffle sentences
	rng.Shuffle(len(sentences), func(i, j int) {
		sentences[i], sentences[j] = sentences[j], sentences[i]
	})

	// Limit to question count if set
	if config.QuestionCount > 0 && len(sentences) > config.QuestionCount {
		sentences = sentences[:config.QuestionCount]
	}

	model := &PracticeModel{
		repo:         repo,
		config:       config,
		sentences:    sentences,
		currentIndex: 0,
		state:        "question",
		rng:          rng,
		width:        80, // Default width
	}

	// Load first sentence
	model.loadCurrentSentence()

	return model, nil
}

func (m *PracticeModel) loadCurrentSentence() {
	if m.currentIndex < len(m.sentences) {
		m.currentSentence = m.sentences[m.currentIndex]
	}
}

func (m PracticeModel) Init() tea.Cmd {
	return nil
}

func (m PracticeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case "question":
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit

			case "enter":
				// Submit answer
				m.isCorrect = m.validateAnswer()
				if m.isCorrect {
					m.correctCount++
				} else {
					m.incorrectCount++
				}

				// Generate explanation using template
				noun, err := m.repo.GetNoun(m.currentSentence.NounID)
				if err != nil {
					m.err = err
				} else {
					m.explanation, err = explanations.Generate(m.currentSentence, noun)
					if err != nil {
						m.err = err
					}
				}

				m.state = "feedback"
				return m, nil

			case "backspace":
				if len(m.userInput) > 0 {
					// Convert to runes to handle multi-byte Unicode characters properly
					runes := []rune(m.userInput)
					if len(runes) > 0 {
						m.userInput = string(runes[:len(runes)-1])
					}
				}

			default:
				// Add character to input
				if len(msg.Runes) > 0 {
					m.userInput += string(msg.Runes)
				}
			}

		case "feedback":
			// Any key continues to next question
			m.currentIndex++
			m.userInput = ""

			// Check if we're done
			if m.config.QuestionCount > 0 && m.currentIndex >= len(m.sentences) {
				m.state = "complete"
			} else if m.currentIndex >= len(m.sentences) {
				// Endless mode - reshuffle and continue
				m.rng.Shuffle(len(m.sentences), func(i, j int) {
					m.sentences[i], m.sentences[j] = m.sentences[j], m.sentences[i]
				})
				m.currentIndex = 0
				m.loadCurrentSentence()
				m.state = "question"
			} else {
				m.loadCurrentSentence()
				m.state = "question"
			}
			return m, nil

		case "complete":
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "r":
				// Restart session
				m.currentIndex = 0
				m.correctCount = 0
				m.incorrectCount = 0
				m.userInput = ""
				m.rng.Shuffle(len(m.sentences), func(i, j int) {
					m.sentences[i], m.sentences[j] = m.sentences[j], m.sentences[i]
				})
				m.loadCurrentSentence()
				m.state = "question"
			}
		}
	}

	return m, nil
}

func (m *PracticeModel) validateAnswer() bool {
	// Exact Unicode string comparison
	return strings.TrimSpace(m.userInput) == strings.TrimSpace(m.currentSentence.CorrectAnswer)
}

func (m PracticeModel) View() string {
	switch m.state {
	case "question":
		return m.renderQuestion()
	case "feedback":
		return m.renderFeedback()
	case "complete":
		return m.renderComplete()
	default:
		return ""
	}
}

func (m PracticeModel) renderQuestion() string {
	var s strings.Builder

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1, 2)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		MarginTop(1).
		MarginBottom(1)

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("120")).
		Bold(true)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		MarginTop(1)

	// Header
	var header string
	if m.config.QuestionCount > 0 {
		header = fmt.Sprintf("Greek Case Master - Question %d/%d", m.currentIndex+1, len(m.sentences))
	} else {
		header = fmt.Sprintf("Greek Case Master - Question %d (Endless)", m.currentIndex+1)
	}
	s.WriteString(titleStyle.Render(header))
	s.WriteString("\n\n")

	// Prompt
	s.WriteString(promptStyle.Render(m.currentSentence.EnglishPrompt))
	s.WriteString("\n\n")

	// Input
	s.WriteString("Your answer:\n")
	s.WriteString(inputStyle.Render("> " + m.userInput + "_"))
	s.WriteString("\n")

	// Hints
	s.WriteString(hintStyle.Render("[Enter to submit] [Ctrl+C or q to quit]"))

	return borderStyle.Render(s.String())
}

func (m PracticeModel) renderFeedback() string {
	// Delegate to feedback.go
	return RenderFeedback(m.isCorrect, m.userInput, m.currentSentence, m.explanation, m.width)
}

func (m PracticeModel) renderComplete() string {
	var s strings.Builder

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1, 2)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("120")).
		MarginBottom(1)

	statsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1)

	s.WriteString(titleStyle.Render("Session Complete!"))
	s.WriteString("\n\n")

	total := m.correctCount + m.incorrectCount
	accuracy := 0
	if total > 0 {
		accuracy = (m.correctCount * 100) / total
	}

	s.WriteString(statsStyle.Render(fmt.Sprintf("Answered: %d/%d", total, len(m.sentences))))
	s.WriteString("\n")
	s.WriteString(statsStyle.Render(fmt.Sprintf("Accuracy: %d%%", accuracy)))
	s.WriteString("\n")

	s.WriteString(hintStyle.Render("\n[q] Quit  [r] Restart session"))

	return borderStyle.Render(s.String())
}
