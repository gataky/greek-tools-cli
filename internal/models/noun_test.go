package models

import (
	"testing"
)

func TestNounStruct(t *testing.T) {
	// Test that Noun struct has all required fields
	noun := Noun{
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

	if noun.English != "teacher" {
		t.Errorf("Expected English 'teacher', got %q", noun.English)
	}

	if noun.Gender != "masculine" {
		t.Errorf("Expected Gender 'masculine', got %q", noun.Gender)
	}

	// Verify Greek Unicode is preserved
	if noun.NominativeSg != "δάσκαλος" {
		t.Errorf("Expected NominativeSg 'δάσκαλος', got %q", noun.NominativeSg)
	}
}

func TestSentenceStruct(t *testing.T) {
	// Test that Sentence struct handles nullable Preposition
	prep := "σε"
	sentence := Sentence{
		NounID:          1,
		EnglishPrompt:   "I go to ___ (the house)",
		GreekSentence:   "Πηγαίνω στο σπίτι",
		CorrectAnswer:   "το σπίτι",
		CaseType:        "accusative",
		Number:          "singular",
		DifficultyPhase: 3,
		ContextType:     "preposition",
		Preposition:     &prep,
	}

	if sentence.Preposition == nil {
		t.Error("Expected non-nil Preposition")
	}

	if *sentence.Preposition != "σε" {
		t.Errorf("Expected Preposition 'σε', got %q", *sentence.Preposition)
	}

	// Test nil preposition
	sentence2 := Sentence{
		NounID:          1,
		EnglishPrompt:   "I see ___",
		GreekSentence:   "Βλέπω το σπίτι",
		CorrectAnswer:   "το σπίτι",
		CaseType:        "accusative",
		Number:          "singular",
		DifficultyPhase: 1,
		ContextType:     "direct_object",
		Preposition:     nil,
	}

	if sentence2.Preposition != nil {
		t.Error("Expected nil Preposition for direct object")
	}
}

func TestSessionConfigStruct(t *testing.T) {
	config := SessionConfig{
		DifficultyLevel: "beginner",
		IncludePlural:   true,
		QuestionCount:   25,
	}

	if config.DifficultyLevel != "beginner" {
		t.Errorf("Expected DifficultyLevel 'beginner', got %q", config.DifficultyLevel)
	}

	if !config.IncludePlural {
		t.Error("Expected IncludePlural to be true")
	}

	if config.QuestionCount != 25 {
		t.Errorf("Expected QuestionCount 25, got %d", config.QuestionCount)
	}

	// Test endless mode (QuestionCount = 0)
	endlessConfig := SessionConfig{
		DifficultyLevel: "advanced",
		IncludePlural:   false,
		QuestionCount:   0,
	}

	if endlessConfig.QuestionCount != 0 {
		t.Errorf("Expected QuestionCount 0 for endless mode, got %d", endlessConfig.QuestionCount)
	}
}
