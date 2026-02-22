package commands

import (
	"fmt"
	"strings"

	"github.com/gataky/greekmaster/internal/storage"
	"github.com/spf13/cobra"
)

// NewListCmd creates the list command
func NewListCmd() *cobra.Command {
	var dbPath string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all nouns in the database",
		Long: `Display a table of all Greek nouns currently in the database.

Shows:
- ID
- English translation
- Greek nominative singular form
- Gender

The list is ordered by ID (insertion order).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize repository
			repo, err := storage.NewSQLiteRepository(dbPath)
			if err != nil {
				return fmt.Errorf("failed to initialize database: %w", err)
			}
			defer repo.Close()

			// Get all nouns
			nouns, err := repo.ListNouns()
			if err != nil {
				return fmt.Errorf("failed to list nouns: %w", err)
			}

			if len(nouns) == 0 {
				fmt.Println("No nouns found in database.")
				fmt.Println("Run 'greekmaster import <csv-file>' or 'greekmaster add' to add nouns.")
				return nil
			}

			// Display table
			fmt.Printf("\nTotal nouns: %d\n\n", len(nouns))

			// Calculate column widths
			maxEnglish := len("English")
			maxGreek := len("Greek")
			for _, noun := range nouns {
				if len(noun.English) > maxEnglish {
					maxEnglish = len(noun.English)
				}
				if len(noun.NominativeSg) > maxGreek {
					maxGreek = len(noun.NominativeSg)
				}
			}

			// Add padding
			maxEnglish += 2
			maxGreek += 2

			// Print header
			header := fmt.Sprintf("%-4s  %-*s  %-*s  %-12s",
				"ID", maxEnglish, "English", maxGreek, "Greek", "Gender")
			fmt.Println(header)
			fmt.Println(strings.Repeat("-", len(header)))

			// Print rows
			const maxRows = 50
			displayCount := len(nouns)
			if displayCount > maxRows {
				displayCount = maxRows
			}

			for i := 0; i < displayCount; i++ {
				noun := nouns[i]
				fmt.Printf("%-4d  %-*s  %-*s  %-12s\n",
					noun.ID,
					maxEnglish, noun.English,
					maxGreek, noun.NominativeSg,
					noun.Gender)
			}

			// Show pagination info if needed
			if len(nouns) > maxRows {
				fmt.Printf("\n... and %d more (showing first %d)\n", len(nouns)-maxRows, maxRows)
			}

			fmt.Println()

			return nil
		},
	}

	cmd.Flags().StringVar(&dbPath, "db-path", "", "Path to database file (default: ~/.greekmaster/greekmaster.db)")

	return cmd
}
