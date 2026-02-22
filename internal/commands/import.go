package commands

import (
	"fmt"
	"os"

	"github.com/gataky/greekmaster/internal/ai"
	"github.com/gataky/greekmaster/internal/importer"
	"github.com/gataky/greekmaster/internal/storage"
	"github.com/spf13/cobra"
)

// NewImportCmd creates the import command
func NewImportCmd() *cobra.Command {
	var dbPath string
	var skipExplanations bool

	cmd := &cobra.Command{
		Use:   "import <csv-file>",
		Short: "Import nouns from a CSV file",
		Long: `Import Greek nouns from a CSV file and generate practice sentences.

The CSV file must have three columns: english, greek, and attribute (gender).
The attribute column should contain: masculine, feminine, neuter, or invariable.

Example CSV format:
  english,greek,attribute
  teacher,δάσκαλος,masculine
  book,βιβλίο,neuter
  woman,γυναίκα,feminine

This command requires the ANTHROPIC_API_KEY environment variable to be set.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			csvPath := args[0]

			// Check if file exists
			if _, err := os.Stat(csvPath); os.IsNotExist(err) {
				return fmt.Errorf("CSV file not found: %s", csvPath)
			}

			// Initialize repository
			repo, err := storage.NewSQLiteRepository(dbPath)
			if err != nil {
				return fmt.Errorf("failed to initialize database: %w", err)
			}
			defer repo.Close()

			// Initialize Claude client
			client, err := ai.NewClaudeClient()
			if err != nil {
				return fmt.Errorf("failed to initialize Claude API client: %w\n\nMake sure ANTHROPIC_API_KEY environment variable is set", err)
			}

			// Create processor and run import
			processor := importer.NewImportProcessor(repo, client)
			processor.SetSkipExplanations(skipExplanations)
			if err := processor.ProcessImport(csvPath); err != nil {
				return fmt.Errorf("import failed: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&dbPath, "db-path", "", "Path to database file (default: ~/.greekmaster/greekmaster.db)")
	cmd.Flags().BoolVar(&skipExplanations, "skip-explanations", false, "Skip generating explanations (saves API calls and cost)")

	return cmd
}
