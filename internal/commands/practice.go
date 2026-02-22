package commands

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gataky/greekmaster/internal/storage"
	"github.com/gataky/greekmaster/internal/tui"
	"github.com/spf13/cobra"
)

// NewPracticeCmd creates the practice command
func NewPracticeCmd() *cobra.Command {
	var dbPath string

	cmd := &cobra.Command{
		Use:   "practice",
		Short: "Start a practice session",
		Long: `Start an interactive practice session to learn Greek noun declension.

You'll be presented with English prompts and must provide the correctly
declined Greek article + noun combination.

The application will guide you through difficulty selection, plural inclusion,
and session type (quick, standard, long, or endless).

After each answer, you'll receive detailed grammar explanations including
translation, syntactic role, and morphology.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize repository
			repo, err := storage.NewSQLiteRepository(dbPath)
			if err != nil {
				return fmt.Errorf("failed to initialize database: %w", err)
			}
			defer repo.Close()

			// Check if database has content
			nouns, err := repo.ListNouns()
			if err != nil {
				return fmt.Errorf("failed to check database: %w", err)
			}
			if len(nouns) == 0 {
				return fmt.Errorf("no nouns found in database. Please run 'greekmaster import <csv-file>' first")
			}

			// Start with setup screen
			setupModel := tui.NewSetupModel()
			p := tea.NewProgram(setupModel)

			finalModel, err := p.Run()
			if err != nil {
				return fmt.Errorf("error running setup: %w", err)
			}

			// Check if user quit during setup
			if setup, ok := finalModel.(tui.SetupModel); ok && setup.View() == "" {
				return nil
			}

			// Get session config
			var config tui.SessionConfigMsg
			setupFinal := finalModel.(tui.SetupModel)
			// Wait for next update to get the config
			_, teaCmd := setupFinal.Update(nil)
			if teaCmd != nil {
				msg := teaCmd()
				if cfgMsg, ok := msg.(tui.SessionConfigMsg); ok {
					config = cfgMsg
				}
			}

			// If we got a config, start practice
			if config.Config.DifficultyLevel != "" {
				practiceModel, err := tui.NewPracticeModel(repo, config.Config)
				if err != nil {
					return fmt.Errorf("failed to initialize practice session: %w", err)
				}

				p = tea.NewProgram(practiceModel)
				if _, err := p.Run(); err != nil {
					return fmt.Errorf("error running practice: %w", err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&dbPath, "db-path", "", "Path to database file (default: ~/.greekmaster/greekmaster.db)")

	return cmd
}
