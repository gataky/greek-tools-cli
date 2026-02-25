package commands

import (
	"fmt"

	"github.com/gataky/greekmaster/internal/storage"
	"github.com/spf13/cobra"
)

// NewMigrateCmd creates the migrate-to-templates command
func NewMigrateCmd() *cobra.Command {
	var dbPath string

	cmd := &cobra.Command{
		Use:   "migrate-to-templates",
		Short: "Migrate from sentence-based to template-based system",
		Long: `Performs a one-time migration that converts the existing sentence storage
model to a template-based system.

This migration will:
1. Analyze all existing sentences to extract unique patterns
2. Create ~100 reusable templates in the sentence_templates table
3. Validate that all sentences can be regenerated from templates
4. Delete the old sentence data (reducing database size by ~80%)

The migration is wrapped in a transaction and will automatically rollback
if any step fails. Your data is safe.

NOTE: This is a one-time operation. After migration, practice sessions will
generate sentences on-demand using templates + noun data.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize repository
			repo, err := storage.NewSQLiteRepository(dbPath)
			if err != nil {
				return fmt.Errorf("failed to initialize database: %w", err)
			}
			defer repo.Close()

			// Check if migration is needed
			fmt.Println("Checking if migration is needed...")
			templates, err := repo.ListTemplates()
			if err != nil {
				return fmt.Errorf("failed to check templates: %w", err)
			}

			if len(templates) > 0 {
				fmt.Printf("Migration already completed. Found %d templates in database.\n", len(templates))
				return nil
			}

			// Confirm with user
			fmt.Println("\n⚠️  IMPORTANT: This operation will modify your database.")
			fmt.Println("A transaction is used to ensure data safety, but it's recommended to backup first.")
			fmt.Print("\nProceed with migration? (yes/no): ")
			var response string
			fmt.Scanln(&response)

			if response != "yes" && response != "y" && response != "Y" {
				fmt.Println("Migration cancelled.")
				return nil
			}

			fmt.Println()

			// Run migration
			err = storage.MigrateToTemplates(repo)
			if err != nil {
				return fmt.Errorf("migration failed: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&dbPath, "db-path", "", "Path to database file (default: ~/.greekmaster/greekmaster.db)")

	return cmd
}
