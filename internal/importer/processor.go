package importer

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gataky/greekmaster/internal/ai"
	"github.com/gataky/greekmaster/internal/models"
	"github.com/gataky/greekmaster/internal/storage"
)

// ImportProcessor orchestrates the CSV import process
type ImportProcessor struct {
	repo   storage.Repository
	client *ai.ClaudeClient
}

// NewImportProcessor creates a new import processor
func NewImportProcessor(repo storage.Repository, client *ai.ClaudeClient) *ImportProcessor {
	return &ImportProcessor{
		repo:   repo,
		client: client,
	}
}

// ProcessImport imports nouns from a CSV file with AI generation
func (p *ImportProcessor) ProcessImport(csvPath string) error {
	// Parse CSV file
	fmt.Println("Parsing CSV file...")
	rows, err := ParseCSV(csvPath)
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %w", err)
	}

	fmt.Printf("Found %d nouns to import\n", len(rows))

	// Check for existing checkpoint
	filename := filepath.Base(csvPath)
	checkpoint, err := p.repo.(*storage.SQLiteRepository).GetCheckpointByFilename(filename)
	if err != nil {
		return fmt.Errorf("failed to check for checkpoint: %w", err)
	}

	startRow := 0
	if checkpoint != nil && checkpoint.Status == "in_progress" {
		fmt.Printf("\nFound existing import in progress (last processed row: %d)\n", checkpoint.LastProcessedRow)
		fmt.Print("Resume from checkpoint? (y/n): ")
		var response string
		fmt.Scanln(&response)
		if response == "y" || response == "Y" {
			startRow = checkpoint.LastProcessedRow
			fmt.Printf("Resuming from row %d\n", startRow+1)
		} else {
			fmt.Println("Starting fresh import")
			startRow = 0
		}
	}

	// Create or update checkpoint
	if checkpoint == nil {
		checkpoint = &storage.ImportCheckpoint{
			CSVFilename:      filename,
			LastProcessedRow: 0,
			Status:           "in_progress",
		}
		if err := p.repo.(*storage.SQLiteRepository).CreateCheckpoint(checkpoint); err != nil {
			return fmt.Errorf("failed to create checkpoint: %w", err)
		}
	} else {
		checkpoint.Status = "in_progress"
		checkpoint.LastProcessedRow = startRow
		if err := p.repo.(*storage.SQLiteRepository).UpdateCheckpoint(checkpoint); err != nil {
			return fmt.Errorf("failed to update checkpoint: %w", err)
		}
	}

	// Track statistics
	totalSentences := 0
	apiCalls := 0
	startTime := time.Now()

	// Process each row
	for i := startRow; i < len(rows); i++ {
		row := rows[i]

		fmt.Printf("\n[%d/%d] Processing '%s' (%s)...\n", i+1, len(rows), row.English, row.Greek)

		// Generate declensions
		fmt.Print("  → Generating declensions... ")
		declensions, err := p.client.GenerateDeclensions(row.Greek, row.English, row.Gender)
		if err != nil {
			fmt.Printf("FAILED\n")
			fmt.Printf("     Error: %v\n", err)
			fmt.Printf("     Skipping this noun and continuing...\n")
			continue
		}
		apiCalls++
		fmt.Println("✓")

		// Create noun record
		noun := &models.Noun{
			English:      row.English,
			Gender:       row.Gender,
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

		if err := p.repo.CreateNoun(noun); err != nil {
			fmt.Printf("     Error storing noun: %v\n", err)
			fmt.Printf("     Skipping this noun and continuing...\n")
			continue
		}

		// Generate sentences
		fmt.Print("  → Generating practice sentences... ")
		sentences, err := p.client.GenerateSentences(row.Greek, row.English, row.Gender)
		if err != nil {
			fmt.Printf("FAILED\n")
			fmt.Printf("     Error: %v\n", err)
			fmt.Printf("     Skipping sentences for this noun...\n")
			// Continue with next noun
			continue
		}
		apiCalls++
		fmt.Printf("✓ (%d sentences)\n", len(sentences))

		// Store sentences
		fmt.Print("  → Storing in database... ")
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

			if err := p.repo.CreateSentence(sentence); err != nil {
				fmt.Printf("\n     Warning: Failed to store sentence %d: %v\n", j+1, err)
				continue
			}

			totalSentences++
		}
		fmt.Println("✓")

		// Update checkpoint after each noun
		checkpoint.LastProcessedRow = i + 1
		if err := p.repo.(*storage.SQLiteRepository).UpdateCheckpoint(checkpoint); err != nil {
			fmt.Printf("     Warning: Failed to update checkpoint: %v\n", err)
		}
	}

	// Mark checkpoint as completed
	checkpoint.Status = "completed"
	if err := p.repo.(*storage.SQLiteRepository).UpdateCheckpoint(checkpoint); err != nil {
		fmt.Printf("Warning: Failed to mark checkpoint as completed: %v\n", err)
	}

	// Print summary
	duration := time.Since(startTime)
	fmt.Print("\n" + strings.Repeat("=", 50) + "\n")
	fmt.Println("Import Complete!")
	fmt.Printf("  Nouns imported: %d\n", len(rows)-startRow)
	fmt.Printf("  Sentences generated: %d\n", totalSentences)
	fmt.Printf("  API calls made: %d\n", apiCalls)
	fmt.Printf("  Time elapsed: %s\n", duration.Round(time.Second))
	fmt.Println(strings.Repeat("=", 50))

	return nil
}
