package storage

import (
	"testing"

	"github.com/gataky/greekmaster/internal/models"
)

func TestDetectFields(t *testing.T) {
	noun := &models.Noun{
		NominativeSg: "δάσκαλος",
		GenitiveSg:   "δασκάλου",
		AccusativeSg: "δάσκαλο",
		NominativePl: "δάσκαλοι",
		GenitivePl:   "δασκάλων",
		AccusativePl: "δασκάλους",
		NomSgArticle: "ο",
		GenSgArticle: "του",
		AccSgArticle: "τον",
		NomPlArticle: "οι",
		GenPlArticle: "των",
		AccPlArticle: "τους",
	}

	tests := []struct {
		name            string
		sentence        *models.Sentence
		wantArticle     string
		wantNounForm    string
		wantErr         bool
	}{
		{
			name: "Accusative singular",
			sentence: &models.Sentence{
				CorrectAnswer: "τον δάσκαλο",
			},
			wantArticle:  "AccSgArticle",
			wantNounForm: "AccusativeSg",
			wantErr:      false,
		},
		{
			name: "Genitive plural",
			sentence: &models.Sentence{
				CorrectAnswer: "των δασκάλων",
			},
			wantArticle:  "GenPlArticle",
			wantNounForm: "GenitivePl",
			wantErr:      false,
		},
		{
			name: "Invalid format",
			sentence: &models.Sentence{
				CorrectAnswer: "invalid",
			},
			wantErr: true,
		},
		{
			name: "Unknown article",
			sentence: &models.Sentence{
				CorrectAnswer: "xyz δάσκαλο",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArticle, gotNounForm, err := detectFields(tt.sentence, noun)
			if (err != nil) != tt.wantErr {
				t.Errorf("detectFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if gotArticle != tt.wantArticle {
					t.Errorf("detectFields() article = %v, want %v", gotArticle, tt.wantArticle)
				}
				if gotNounForm != tt.wantNounForm {
					t.Errorf("detectFields() nounForm = %v, want %v", gotNounForm, tt.wantNounForm)
				}
			}
		})
	}
}

func TestCreateGreekTemplate(t *testing.T) {
	tests := []struct {
		name          string
		greekSentence string
		article       string
		nounForm      string
		want          string
	}{
		{
			name:          "Simple sentence",
			greekSentence: "Βλέπω τον δάσκαλο",
			article:       "τον",
			nounForm:      "δάσκαλο",
			want:          "Βλέπω {article} {noun_form}",
		},
		{
			name:          "Sentence with multiple words",
			greekSentence: "Η τσάντα του δασκάλου είναι κόκκινη",
			article:       "του",
			nounForm:      "δασκάλου",
			want:          "Η τσάντα {article} {noun_form} είναι κόκκινη",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createGreekTemplate(tt.greekSentence, tt.article, tt.nounForm)
			if got != tt.want {
				t.Errorf("createGreekTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateEnglishTemplate(t *testing.T) {
	tests := []struct {
		name          string
		englishPrompt string
		englishNoun   string
		want          string
	}{
		{
			name:          "With parentheses",
			englishPrompt: "I see ___ (the teacher)",
			englishNoun:   "teacher",
			want:          "I see ___ {noun}",
		},
		{
			name:          "Simple replacement",
			englishPrompt: "The teacher is here",
			englishNoun:   "teacher",
			want:          "The {noun} is here",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createEnglishTemplate(tt.englishPrompt, tt.englishNoun)
			if got != tt.want {
				t.Errorf("createEnglishTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnalyzeSentencePatterns(t *testing.T) {
	nouns := map[int64]*models.Noun{
		1: {
			ID:           1,
			English:      "teacher",
			AccSgArticle: "τον",
			AccusativeSg: "δάσκαλο",
		},
		2: {
			ID:           2,
			English:      "book",
			AccSgArticle: "το",
			AccusativeSg: "βιβλίο",
		},
	}

	sentences := []*models.Sentence{
		// Same pattern, different nouns - should deduplicate
		{
			ID:            1,
			NounID:        1,
			EnglishPrompt: "I see ___ (the teacher)",
			GreekSentence: "Βλέπω τον δάσκαλο",
			CorrectAnswer: "τον δάσκαλο",
			CaseType:      "accusative",
			Number:        "singular",
		},
		{
			ID:            2,
			NounID:        2,
			EnglishPrompt: "I see ___ (the book)",
			GreekSentence: "Βλέπω το βιβλίο",
			CorrectAnswer: "το βιβλίο",
			CaseType:      "accusative",
			Number:        "singular",
		},
		// Different pattern
		{
			ID:            3,
			NounID:        1,
			EnglishPrompt: "She wants ___ (the teacher)",
			GreekSentence: "Θέλει τον δάσκαλο",
			CorrectAnswer: "τον δάσκαλο",
			CaseType:      "accusative",
			Number:        "singular",
		},
	}

	templates, err := analyzeSentencePatterns(sentences, nouns)
	if err != nil {
		t.Fatalf("analyzeSentencePatterns() error = %v", err)
	}

	// Should have 2 unique patterns (sentence 1 and 2 are the same pattern)
	if len(templates) != 2 {
		t.Errorf("analyzeSentencePatterns() returned %d templates, want 2", len(templates))
	}

	// Verify templates have placeholders
	for _, tmpl := range templates {
		if tmpl.EnglishTemplate == "" {
			t.Error("EnglishTemplate is empty")
		}
		if tmpl.GreekTemplate == "" {
			t.Error("GreekTemplate is empty")
		}
		if tmpl.ArticleField == "" {
			t.Error("ArticleField is empty")
		}
		if tmpl.NounFormField == "" {
			t.Error("NounFormField is empty")
		}
	}
}

func TestMigrateToTemplatesEndToEnd(t *testing.T) {
	// Skip this test since migration 003 already ran in setupTestDB
	// which means the sentences table doesn't exist anymore
	t.Skip("Skipping end-to-end migration test - migration already applied in test setup")

	repo := setupTestDB(t)
	defer repo.Close()

	// Create test nouns
	nouns := []*models.Noun{
		{
			English:      "teacher",
			Gender:       "masculine",
			NominativeSg: "δάσκαλος",
			GenitiveSg:   "δασκάλου",
			AccusativeSg: "δάσκαλο",
			NominativePl: "δάσκαλοι",
			GenitivePl:   "δασκάλων",
			AccusativePl: "δασκάλους",
			NomSgArticle: "ο",
			GenSgArticle: "του",
			AccSgArticle: "τον",
			NomPlArticle: "οι",
			GenPlArticle: "των",
			AccPlArticle: "τους",
		},
	}

	for _, noun := range nouns {
		if err := repo.CreateNoun(noun); err != nil {
			t.Fatalf("CreateNoun() error = %v", err)
		}
	}

	// Create test sentences
	sentences := []*models.Sentence{
		{
			NounID:          nouns[0].ID,
			EnglishPrompt:   "I see ___ (the teacher)",
			GreekSentence:   "Βλέπω τον δάσκαλο",
			CorrectAnswer:   "τον δάσκαλο",
			CaseType:        "accusative",
			Number:          "singular",
			DifficultyPhase: 1,
			ContextType:     "direct_object",
		},
		{
			NounID:          nouns[0].ID,
			EnglishPrompt:   "She wants ___ (the teacher)",
			GreekSentence:   "Θέλει τον δάσκαλο",
			CorrectAnswer:   "τον δάσκαλο",
			CaseType:        "accusative",
			Number:          "singular",
			DifficultyPhase: 1,
			ContextType:     "direct_object",
		},
	}

	for _, sent := range sentences {
		if err := repo.CreateSentence(sent); err != nil {
			t.Fatalf("CreateSentence() error = %v", err)
		}
	}

	// Run migration
	err := MigrateToTemplates(repo)
	if err != nil {
		t.Fatalf("MigrateToTemplates() error = %v", err)
	}

	// Verify templates were created
	templates, err := repo.ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates() error = %v", err)
	}

	if len(templates) == 0 {
		t.Error("No templates created")
	}

	// Verify sentences were deleted
	var count int
	err = repo.db.Get(&count, "SELECT COUNT(*) FROM sentences")
	if err != nil {
		t.Fatalf("Failed to count sentences: %v", err)
	}

	if count != 0 {
		t.Errorf("Sentences not deleted, found %d sentences", count)
	}

	// Verify we can generate sentences from templates
	generated, err := repo.GeneratePracticeSentences(1, "singular", 5)
	if err != nil {
		t.Fatalf("GeneratePracticeSentences() after migration error = %v", err)
	}

	if len(generated) == 0 {
		t.Error("Could not generate sentences after migration")
	}
}

func TestValidateMigration(t *testing.T) {
	templates := []*models.SentenceTemplate{
		{
			CaseType:        "accusative",
			DifficultyPhase: 1,
			ContextType:     "direct_object",
		},
		{
			CaseType:        "genitive",
			DifficultyPhase: 2,
			ContextType:     "possession",
		},
	}

	sentences := []*models.Sentence{
		{
			CaseType:        "accusative",
			DifficultyPhase: 1,
			ContextType:     "direct_object",
		},
		{
			CaseType:        "genitive",
			DifficultyPhase: 2,
			ContextType:     "possession",
		},
	}

	nouns := make(map[int64]*models.Noun)

	// Should pass - all combinations covered
	err := validateMigration(templates, sentences, nouns)
	if err != nil {
		t.Errorf("validateMigration() error = %v, want nil", err)
	}

	// Test with missing combination
	sentencesWithExtra := append(sentences, &models.Sentence{
		CaseType:        "nominative",
		DifficultyPhase: 3,
		ContextType:     "direct_object",
	})

	err = validateMigration(templates, sentencesWithExtra, nouns)
	if err == nil {
		t.Error("validateMigration() expected error for missing combination")
	}
}
