package storage

import (
	"strings"
	"testing"

	"github.com/gataky/greekmaster/internal/models"
)

func TestGetFieldValue(t *testing.T) {
	noun := &models.Noun{
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
	}

	tests := []struct {
		name      string
		fieldName string
		want      string
		wantErr   bool
	}{
		{"AccSgArticle", "AccSgArticle", "τον", false},
		{"AccusativeSg", "AccusativeSg", "δάσκαλο", false},
		{"English", "English", "teacher", false},
		{"GenPlArticle", "GenPlArticle", "των", false},
		{"Invalid field", "InvalidField", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFieldValue(noun, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFieldValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getFieldValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubstituteTemplate(t *testing.T) {
	noun := &models.Noun{
		ID:           1,
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
	}

	tests := []struct {
		name     string
		template *models.SentenceTemplate
		wantEng  string
		wantGr   string
		wantAns  string
	}{
		{
			name: "Accusative singular",
			template: &models.SentenceTemplate{
				EnglishTemplate: "I see {noun}",
				GreekTemplate:   "Βλέπω {article} {noun_form}",
				ArticleField:    "AccSgArticle",
				NounFormField:   "AccusativeSg",
				CaseType:        "accusative",
				Number:          "singular",
				DifficultyPhase: 1,
				ContextType:     "direct_object",
			},
			wantEng: "I see teacher",
			wantGr:  "Βλέπω τον δάσκαλο",
			wantAns: "τον δάσκαλο",
		},
		{
			name: "Genitive plural",
			template: &models.SentenceTemplate{
				EnglishTemplate: "The books of {noun}",
				GreekTemplate:   "Τα βιβλία {article} {noun_form}",
				ArticleField:    "GenPlArticle",
				NounFormField:   "GenitivePl",
				CaseType:        "genitive",
				Number:          "plural",
				DifficultyPhase: 2,
				ContextType:     "possession",
			},
			wantEng: "The books of teacher",
			wantGr:  "Τα βιβλία των δασκάλων",
			wantAns: "των δασκάλων",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := substituteTemplate(tt.template, noun)
			if err != nil {
				t.Fatalf("substituteTemplate() error = %v", err)
			}

			if got.EnglishPrompt != tt.wantEng {
				t.Errorf("EnglishPrompt = %v, want %v", got.EnglishPrompt, tt.wantEng)
			}

			if got.GreekSentence != tt.wantGr {
				t.Errorf("GreekSentence = %v, want %v", got.GreekSentence, tt.wantGr)
			}

			if got.CorrectAnswer != tt.wantAns {
				t.Errorf("CorrectAnswer = %v, want %v", got.CorrectAnswer, tt.wantAns)
			}

			if got.CaseType != tt.template.CaseType {
				t.Errorf("CaseType = %v, want %v", got.CaseType, tt.template.CaseType)
			}

			if got.NounID != noun.ID {
				t.Errorf("NounID = %v, want %v", got.NounID, noun.ID)
			}
		})
	}
}

func TestSubstituteTemplateWithPreposition(t *testing.T) {
	noun := &models.Noun{
		ID:           1,
		English:      "teacher",
		AccSgArticle: "τον",
		AccusativeSg: "δάσκαλο",
	}

	prep := "σε"
	template := &models.SentenceTemplate{
		EnglishTemplate: "He is talking to {noun}",
		GreekTemplate:   "Μιλάει σε {article} {noun_form}",
		ArticleField:    "AccSgArticle",
		NounFormField:   "AccusativeSg",
		CaseType:        "accusative",
		Number:          "singular",
		DifficultyPhase: 3,
		ContextType:     "preposition",
		Preposition:     &prep,
	}

	got, err := substituteTemplate(template, noun)
	if err != nil {
		t.Fatalf("substituteTemplate() error = %v", err)
	}

	if got.Preposition == nil || *got.Preposition != prep {
		t.Errorf("Preposition = %v, want %v", got.Preposition, prep)
	}

	if got.ContextType != "preposition" {
		t.Errorf("ContextType = %v, want preposition", got.ContextType)
	}
}

func TestGeneratePracticeSentences(t *testing.T) {
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
		{
			English:      "book",
			Gender:       "neuter",
			NominativeSg: "βιβλίο",
			GenitiveSg:   "βιβλίου",
			AccusativeSg: "βιβλίο",
			NominativePl: "βιβλία",
			GenitivePl:   "βιβλίων",
			AccusativePl: "βιβλία",
			NomSgArticle: "το",
			GenSgArticle: "του",
			AccSgArticle: "το",
			NomPlArticle: "τα",
			GenPlArticle: "των",
			AccPlArticle: "τα",
		},
	}

	for _, noun := range nouns {
		if err := repo.CreateNoun(noun); err != nil {
			t.Fatalf("CreateNoun() error = %v", err)
		}
	}

	// Create test templates
	templates := []*models.SentenceTemplate{
		{
			EnglishTemplate: "I see {noun}",
			GreekTemplate:   "Βλέπω {article} {noun_form}",
			ArticleField:    "AccSgArticle",
			NounFormField:   "AccusativeSg",
			CaseType:        "accusative",
			Number:          "singular",
			DifficultyPhase: 1,
			ContextType:     "direct_object",
		},
		{
			EnglishTemplate: "She wants {noun}",
			GreekTemplate:   "Θέλει {article} {noun_form}",
			ArticleField:    "AccSgArticle",
			NounFormField:   "AccusativeSg",
			CaseType:        "accusative",
			Number:          "singular",
			DifficultyPhase: 1,
			ContextType:     "direct_object",
		},
	}

	for _, tmpl := range templates {
		if err := repo.CreateTemplate(tmpl); err != nil {
			t.Fatalf("CreateTemplate() error = %v", err)
		}
	}

	// Test generation
	sentences, err := repo.GeneratePracticeSentences(1, "singular", 5)
	if err != nil {
		t.Fatalf("GeneratePracticeSentences() error = %v", err)
	}

	if len(sentences) == 0 {
		t.Error("GeneratePracticeSentences() returned no sentences")
	}

	// Verify sentences have correct structure
	for _, sent := range sentences {
		if sent.EnglishPrompt == "" {
			t.Error("EnglishPrompt is empty")
		}
		if sent.GreekSentence == "" {
			t.Error("GreekSentence is empty")
		}
		if sent.CorrectAnswer == "" {
			t.Error("CorrectAnswer is empty")
		}
		if sent.CaseType != "accusative" {
			t.Errorf("CaseType = %v, want accusative", sent.CaseType)
		}
		if sent.DifficultyPhase != 1 {
			t.Errorf("DifficultyPhase = %v, want 1", sent.DifficultyPhase)
		}
		// Verify no placeholders remain
		if strings.Contains(sent.EnglishPrompt, "{") || strings.Contains(sent.GreekSentence, "{") {
			t.Error("Generated sentence contains placeholders")
		}
	}
}

func TestGeneratePracticeSentencesFiltering(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Create test noun
	noun := &models.Noun{
		English:      "teacher",
		Gender:       "masculine",
		AccSgArticle: "τον",
		AccusativeSg: "δάσκαλο",
		AccPlArticle: "τους",
		AccusativePl: "δασκάλους",
	}

	if err := repo.CreateNoun(noun); err != nil {
		t.Fatalf("CreateNoun() error = %v", err)
	}

	// Create templates for different phases
	templates := []*models.SentenceTemplate{
		{
			EnglishTemplate: "I see {noun}",
			GreekTemplate:   "Βλέπω {article} {noun_form}",
			ArticleField:    "AccSgArticle",
			NounFormField:   "AccusativeSg",
			CaseType:        "accusative",
			Number:          "singular",
			DifficultyPhase: 1,
			ContextType:     "direct_object",
		},
		{
			EnglishTemplate: "I see {noun}",
			GreekTemplate:   "Βλέπω {article} {noun_form}",
			ArticleField:    "AccPlArticle",
			NounFormField:   "AccusativePl",
			CaseType:        "accusative",
			Number:          "plural",
			DifficultyPhase: 1,
			ContextType:     "direct_object",
		},
		{
			EnglishTemplate: "The bag of {noun}",
			GreekTemplate:   "Η τσάντα {article} {noun_form}",
			ArticleField:    "GenSgArticle",
			NounFormField:   "GenitiveSg",
			CaseType:        "genitive",
			Number:          "singular",
			DifficultyPhase: 2,
			ContextType:     "possession",
		},
	}

	for _, tmpl := range templates {
		if err := repo.CreateTemplate(tmpl); err != nil {
			t.Fatalf("CreateTemplate() error = %v", err)
		}
	}

	// Test phase filtering
	phase1, err := repo.GeneratePracticeSentences(1, "", 10)
	if err != nil {
		t.Fatalf("GeneratePracticeSentences(phase=1) error = %v", err)
	}

	for _, sent := range phase1 {
		if sent.DifficultyPhase != 1 {
			t.Errorf("Phase 1 filtering failed: got phase %d", sent.DifficultyPhase)
		}
	}

	// Test number filtering
	singular, err := repo.GeneratePracticeSentences(1, "singular", 10)
	if err != nil {
		t.Fatalf("GeneratePracticeSentences(number=singular) error = %v", err)
	}

	for _, sent := range singular {
		if sent.Number != "singular" {
			t.Errorf("Singular filtering failed: got number %s", sent.Number)
		}
	}
}
