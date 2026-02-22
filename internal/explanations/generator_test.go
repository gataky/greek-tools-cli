package explanations

import (
	"testing"

	"github.com/gataky/greekmaster/internal/models"
)

func TestGenerate(t *testing.T) {
	// Create a test noun
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
		name     string
		sentence *models.Sentence
		noun     *models.Noun
		wantErr  bool
	}{
		{
			name: "accusative singular direct object",
			sentence: &models.Sentence{
				NounID:          1,
				EnglishPrompt:   "I see ___ (the teacher)",
				GreekSentence:   "Βλέπω τον δάσκαλο",
				CorrectAnswer:   "τον δάσκαλο",
				CaseType:        "accusative",
				Number:          "singular",
				DifficultyPhase: 1,
				ContextType:     "direct_object",
				Preposition:     nil,
			},
			noun:    noun,
			wantErr: false,
		},
		{
			name: "genitive singular possession",
			sentence: &models.Sentence{
				NounID:          1,
				EnglishPrompt:   "The book of ___ (the teacher)",
				GreekSentence:   "Το βιβλίο του δασκάλου",
				CorrectAnswer:   "του δασκάλου",
				CaseType:        "genitive",
				Number:          "singular",
				DifficultyPhase: 2,
				ContextType:     "possession",
				Preposition:     nil,
			},
			noun:    noun,
			wantErr: false,
		},
		{
			name: "accusative plural direct object",
			sentence: &models.Sentence{
				NounID:          1,
				EnglishPrompt:   "I see ___ (the teachers)",
				GreekSentence:   "Βλέπω τους δασκάλους",
				CorrectAnswer:   "τους δασκάλους",
				CaseType:        "accusative",
				Number:          "plural",
				DifficultyPhase: 1,
				ContextType:     "direct_object",
				Preposition:     nil,
			},
			noun:    noun,
			wantErr: false,
		},
		{
			name: "accusative with preposition σε",
			sentence: &models.Sentence{
				NounID:          1,
				EnglishPrompt:   "I go to ___ (the teacher)",
				GreekSentence:   "Πηγαίνω στον δάσκαλο",
				CorrectAnswer:   "τον δάσκαλο",
				CaseType:        "accusative",
				Number:          "singular",
				DifficultyPhase: 3,
				ContextType:     "preposition",
				Preposition:     stringPtr("σε"),
			},
			noun:    noun,
			wantErr: false,
		},
		{
			name:     "nil sentence",
			sentence: nil,
			noun:     noun,
			wantErr:  true,
		},
		{
			name: "nil noun",
			sentence: &models.Sentence{
				EnglishPrompt: "test",
				CaseType:      "accusative",
				Number:        "singular",
				ContextType:   "direct_object",
			},
			noun:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			explanation, err := Generate(tt.sentence, tt.noun)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if explanation == nil {
					t.Error("Generate() returned nil explanation")
					return
				}
				if explanation.Translation == "" {
					t.Error("Generate() returned empty Translation")
				}
				if explanation.SyntacticRole == "" {
					t.Error("Generate() returned empty SyntacticRole")
				}
				if explanation.Morphology == "" {
					t.Error("Generate() returned empty Morphology")
				}
			}
		})
	}
}

func TestSyntacticRoleTemplate(t *testing.T) {
	tests := []struct {
		name        string
		contextType string
		caseType    string
		prep        *string
		want        string
	}{
		{
			name:        "direct object",
			contextType: "direct_object",
			caseType:    "accusative",
			prep:        nil,
			want:        "Direct objects use accusative case",
		},
		{
			name:        "possession",
			contextType: "possession",
			caseType:    "genitive",
			prep:        nil,
			want:        "Possession requires genitive case",
		},
		{
			name:        "preposition σε",
			contextType: "preposition",
			caseType:    "accusative",
			prep:        stringPtr("σε"),
			want:        "The preposition 'σε' requires accusative case",
		},
		{
			name:        "preposition από",
			contextType: "preposition",
			caseType:    "genitive",
			prep:        stringPtr("από"),
			want:        "The preposition 'από' requires genitive case",
		},
		{
			name:        "preposition για",
			contextType: "preposition",
			caseType:    "genitive",
			prep:        stringPtr("για"),
			want:        "The preposition 'για' requires genitive case",
		},
		{
			name:        "preposition με",
			contextType: "preposition",
			caseType:    "accusative",
			prep:        stringPtr("με"),
			want:        "The preposition 'με' requires accusative case",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SyntacticRoleTemplate(tt.contextType, tt.caseType, tt.prep)
			if got != tt.want {
				t.Errorf("SyntacticRoleTemplate() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGenerateTranslation(t *testing.T) {
	tests := []struct {
		name          string
		englishPrompt string
		greekSentence string
		want          string
	}{
		{
			name:          "simple with parentheses",
			englishPrompt: "I see ___ (the teacher)",
			greekSentence: "Βλέπω τον δάσκαλο",
			want:          "I see the teacher",
		},
		{
			name:          "with parentheses at end",
			englishPrompt: "The book of ___ (the teacher)",
			greekSentence: "Το βιβλίο του δασκάλου",
			want:          "The book of the teacher",
		},
		{
			name:          "no blank",
			englishPrompt: "I see the teacher",
			greekSentence: "Βλέπω τον δάσκαλο",
			want:          "Βλέπω τον δάσκαλο",
		},
		{
			name:          "no parentheses",
			englishPrompt: "I see ___",
			greekSentence: "Βλέπω τον δάσκαλο",
			want:          "I see [answer]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateTranslation(tt.englishPrompt, tt.greekSentence)
			if got != tt.want {
				t.Errorf("GenerateTranslation() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatMorphology(t *testing.T) {
	// Masculine noun
	masculine := &models.Noun{
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

	// Neuter noun
	neuter := &models.Noun{
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
	}

	tests := []struct {
		name     string
		noun     *models.Noun
		caseType string
		number   string
		want     string
	}{
		{
			name:     "masculine accusative singular",
			noun:     masculine,
			caseType: "accusative",
			number:   "singular",
			want:     "ο δάσκαλος → τον δάσκαλο",
		},
		{
			name:     "masculine genitive singular",
			noun:     masculine,
			caseType: "genitive",
			number:   "singular",
			want:     "ο δάσκαλος → του δασκάλου",
		},
		{
			name:     "masculine accusative plural",
			noun:     masculine,
			caseType: "accusative",
			number:   "plural",
			want:     "ο δάσκαλος → τους δασκάλους",
		},
		{
			name:     "masculine genitive plural",
			noun:     masculine,
			caseType: "genitive",
			number:   "plural",
			want:     "ο δάσκαλος → των δασκάλων",
		},
		{
			name:     "neuter accusative singular",
			noun:     neuter,
			caseType: "accusative",
			number:   "singular",
			want:     "το βιβλίο → το βιβλίο",
		},
		{
			name:     "neuter genitive plural",
			noun:     neuter,
			caseType: "genitive",
			number:   "plural",
			want:     "το βιβλίο → των βιβλίων",
		},
		{
			name:     "nominative singular (no change)",
			noun:     masculine,
			caseType: "nominative",
			number:   "singular",
			want:     "ο δάσκαλος → ο δάσκαλος",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatMorphology(tt.noun, tt.caseType, tt.number)
			if got != tt.want {
				t.Errorf("FormatMorphology() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
