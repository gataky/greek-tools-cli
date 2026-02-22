package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gataky/greekmaster/internal/ai"
	"github.com/gataky/greekmaster/internal/models"
	"github.com/gataky/greekmaster/internal/storage"
	"github.com/spf13/cobra"
)

// NewAddCmd creates the add command
func NewAddCmd() *cobra.Command {
	var dbPath string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a single noun interactively",
		Long: `Add a single Greek noun to the database with AI-generated declensions and sentences.

You'll be prompted to enter:
- English translation
- Greek nominative singular form
- Gender (masculine, feminine, neuter, or invariable)

The application will then use Claude API to generate all declined forms
and practice sentences.

This command requires the ANTHROPIC_API_KEY environment variable to be set.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			reader := bufio.NewReader(os.Stdin)

			// Prompt for English translation
			fmt.Print("English translation: ")
			english, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			english = strings.TrimSpace(english)
			if english == "" {
				return fmt.Errorf("english translation cannot be empty")
			}

			// Prompt for Greek nominative singular
			fmt.Print("Greek (nominative singular): ")
			greek, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			greek = strings.TrimSpace(greek)
			if greek == "" {
				return fmt.Errorf("greek form cannot be empty")
			}

			// Prompt for gender
			fmt.Println("\nGender:")
			fmt.Println("  1. Masculine")
			fmt.Println("  2. Feminine")
			fmt.Println("  3. Neuter")
			fmt.Println("  4. Invariable")
			fmt.Print("\nSelect (1-4): ")

			genderInput, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			genderInput = strings.TrimSpace(genderInput)

			var gender string
			switch genderInput {
			case "1":
				gender = "masculine"
			case "2":
				gender = "feminine"
			case "3":
				gender = "neuter"
			case "4":
				gender = "invariable"
			default:
				return fmt.Errorf("invalid gender selection")
			}

			fmt.Printf("\nAdding: %s (%s, %s)\n\n", english, greek, gender)

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

			// Generate declensions
			fmt.Print("Generating declensions... ")
			declensions, err := client.GenerateDeclensions(greek, english, gender)
			if err != nil {
				fmt.Println("FAILED")
				return fmt.Errorf("failed to generate declensions: %w", err)
			}
			fmt.Println("✓")

			// Create noun record
			noun := &models.Noun{
				English:      english,
				Gender:       gender,
				NominativeSg: declensions.NominativeSg,
				GenitiveSg:   declensions.GenitiveSg,
				AccusativeSg: declensions.AccusativeSg,
				NominativePl: declensions.NominativePl,
				GenitivePl:   declensions.GenitivePl,
				AccusativePl: declensions.AccusativePl,
				NomSgArticle: declensions.NomSgArticle,
				GenSgArticle: declensions.GenSgArticle,
				AccSgArticle: declensions.AccSgArticle,
				NomPlArticle: declensions.NomPlArticle,
				GenPlArticle: declensions.GenPlArticle,
				AccPlArticle: declensions.AccPlArticle,
			}

			if err := repo.CreateNoun(noun); err != nil {
				return fmt.Errorf("failed to store noun: %w", err)
			}

			// Generate sentences
			fmt.Print("Generating practice sentences... ")
			sentences, err := client.GenerateSentences(greek, english, gender)
			if err != nil {
				fmt.Println("FAILED")
				return fmt.Errorf("failed to generate sentences: %w", err)
			}
			fmt.Printf("✓ (%d sentences)\n", len(sentences))

			// Store sentences
			fmt.Print("Storing in database... ")
			for j, sentenceResp := range sentences {
				sentence := &models.Sentence{
					NounID:          noun.ID,
					EnglishPrompt:   sentenceResp.EnglishPrompt,
					GreekSentence:   sentenceResp.GreekSentence,
					CorrectAnswer:   sentenceResp.CorrectAnswer,
					CaseType:        sentenceResp.CaseType,
					Number:          sentenceResp.Number,
					DifficultyPhase: sentenceResp.DifficultyPhase,
					ContextType:     sentenceResp.ContextType,
					Preposition:     sentenceResp.Preposition,
				}

				if err := repo.CreateSentence(sentence); err != nil {
					return fmt.Errorf("failed to store sentence %d: %w", j+1, err)
				}
			}
			fmt.Println("✓")

			fmt.Printf("\n✓ Successfully added '%s' with %d practice sentences\n", english, len(sentences))

			return nil
		},
	}

	cmd.Flags().StringVar(&dbPath, "db-path", "", "Path to database file (default: ~/.greekmaster/greekmaster.db)")

	return cmd
}
