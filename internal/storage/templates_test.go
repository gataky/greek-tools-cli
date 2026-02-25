package storage

import (
	"testing"

	"github.com/gataky/greekmaster/internal/models"
)

func TestCreateTemplate(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	template := &models.SentenceTemplate{
		EnglishTemplate: "I see {noun}",
		GreekTemplate:   "Βλέπω {article} {noun_form}",
		ArticleField:    "AccSgArticle",
		NounFormField:   "AccusativeSg",
		CaseType:        "accusative",
		Number:          "singular",
		DifficultyPhase: 1,
		ContextType:     "direct_object",
	}

	err := repo.CreateTemplate(template)
	if err != nil {
		t.Fatalf("CreateTemplate() error = %v", err)
	}

	if template.ID == 0 {
		t.Error("Expected template ID to be set after creation")
	}
}

func TestGetTemplate(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Create template first
	template := &models.SentenceTemplate{
		EnglishTemplate: "I see {noun}",
		GreekTemplate:   "Βλέπω {article} {noun_form}",
		ArticleField:    "AccSgArticle",
		NounFormField:   "AccusativeSg",
		CaseType:        "accusative",
		Number:          "singular",
		DifficultyPhase: 1,
		ContextType:     "direct_object",
	}

	err := repo.CreateTemplate(template)
	if err != nil {
		t.Fatalf("CreateTemplate() error = %v", err)
	}

	// Get template
	retrieved, err := repo.GetTemplate(template.ID)
	if err != nil {
		t.Fatalf("GetTemplate() error = %v", err)
	}

	if retrieved.EnglishTemplate != template.EnglishTemplate {
		t.Errorf("EnglishTemplate = %v, want %v", retrieved.EnglishTemplate, template.EnglishTemplate)
	}

	if retrieved.GreekTemplate != template.GreekTemplate {
		t.Errorf("GreekTemplate = %v, want %v", retrieved.GreekTemplate, template.GreekTemplate)
	}

	// Test not found error
	_, err = repo.GetTemplate(99999)
	if err == nil {
		t.Error("Expected error for non-existent template")
	}
}

func TestListTemplates(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Create multiple templates
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
		err := repo.CreateTemplate(tmpl)
		if err != nil {
			t.Fatalf("CreateTemplate() error = %v", err)
		}
	}

	// List all templates
	list, err := repo.ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates() error = %v", err)
	}

	if len(list) != len(templates) {
		t.Errorf("ListTemplates() returned %d templates, want %d", len(list), len(templates))
	}
}

func TestGetRandomTemplates(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Create templates for different phases and numbers
	templates := []*models.SentenceTemplate{
		// Phase 1, singular
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
		// Phase 1, plural
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
		// Phase 2, singular
		{
			EnglishTemplate: "The bag of {noun} is red",
			GreekTemplate:   "Η τσάντα {article} {noun_form} είναι κόκκινη",
			ArticleField:    "GenSgArticle",
			NounFormField:   "GenitiveSg",
			CaseType:        "genitive",
			Number:          "singular",
			DifficultyPhase: 2,
			ContextType:     "possession",
		},
	}

	for _, tmpl := range templates {
		err := repo.CreateTemplate(tmpl)
		if err != nil {
			t.Fatalf("CreateTemplate() error = %v", err)
		}
	}

	// Test filtering by phase and number
	tests := []struct {
		name     string
		phase    int
		number   string
		wantMin  int
		wantMax  int
	}{
		{"Phase 1 singular", 1, "singular", 1, 1},
		{"Phase 1 plural", 1, "plural", 1, 1},
		{"Phase 1 both", 1, "", 2, 2},
		{"Phase 2 singular", 2, "singular", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := repo.GetRandomTemplates(tt.phase, tt.number, 10)
			if err != nil {
				t.Fatalf("GetRandomTemplates() error = %v", err)
			}

			if len(results) < tt.wantMin || len(results) > tt.wantMax {
				t.Errorf("GetRandomTemplates() returned %d templates, want between %d and %d",
					len(results), tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestTemplateConstraints(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	tests := []struct {
		name     string
		template *models.SentenceTemplate
		wantErr  bool
	}{
		{
			name: "Invalid case_type",
			template: &models.SentenceTemplate{
				EnglishTemplate: "Test",
				GreekTemplate:   "Test",
				ArticleField:    "AccSgArticle",
				NounFormField:   "AccusativeSg",
				CaseType:        "invalid",
				Number:          "singular",
				DifficultyPhase: 1,
				ContextType:     "direct_object",
			},
			wantErr: true,
		},
		{
			name: "Invalid number",
			template: &models.SentenceTemplate{
				EnglishTemplate: "Test",
				GreekTemplate:   "Test",
				ArticleField:    "AccSgArticle",
				NounFormField:   "AccusativeSg",
				CaseType:        "accusative",
				Number:          "invalid",
				DifficultyPhase: 1,
				ContextType:     "direct_object",
			},
			wantErr: true,
		},
		{
			name: "Invalid difficulty_phase",
			template: &models.SentenceTemplate{
				EnglishTemplate: "Test",
				GreekTemplate:   "Test",
				ArticleField:    "AccSgArticle",
				NounFormField:   "AccusativeSg",
				CaseType:        "accusative",
				Number:          "singular",
				DifficultyPhase: 5,
				ContextType:     "direct_object",
			},
			wantErr: true,
		},
		{
			name: "Invalid context_type",
			template: &models.SentenceTemplate{
				EnglishTemplate: "Test",
				GreekTemplate:   "Test",
				ArticleField:    "AccSgArticle",
				NounFormField:   "AccusativeSg",
				CaseType:        "accusative",
				Number:          "singular",
				DifficultyPhase: 1,
				ContextType:     "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateTemplate(tt.template)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
