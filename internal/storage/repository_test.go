package storage

import (
	"testing"

	"github.com/gataky/greekmaster/internal/models"
)

func setupTestDB(t *testing.T) *SQLiteRepository {
	// Use in-memory database for testing
	repo, err := NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	return repo
}

func TestCreateAndGetNoun(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

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

	// Create noun
	err := repo.CreateNoun(noun)
	if err != nil {
		t.Fatalf("CreateNoun() error = %v", err)
	}

	if noun.ID == 0 {
		t.Error("Expected noun ID to be set after creation")
	}

	// Get noun
	retrieved, err := repo.GetNoun(noun.ID)
	if err != nil {
		t.Fatalf("GetNoun() error = %v", err)
	}

	if retrieved.English != noun.English {
		t.Errorf("Expected English %q, got %q", noun.English, retrieved.English)
	}

	if retrieved.NominativeSg != noun.NominativeSg {
		t.Errorf("Expected NominativeSg %q, got %q", noun.NominativeSg, retrieved.NominativeSg)
	}
}

func TestListNouns(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Create multiple nouns
	nouns := []*models.Noun{
		{
			English: "teacher", Gender: "masculine",
			NominativeSg: "δάσκαλος", GenitiveSg: "δασκάλου", AccusativeSg: "δάσκαλο",
			NominativePl: "δάσκαλοι", GenitivePl: "δασκάλων", AccusativePl: "δασκάλους",
			NomSgArticle: "ο", GenSgArticle: "του", AccSgArticle: "τον",
			NomPlArticle: "οι", GenPlArticle: "των", AccPlArticle: "τους",
		},
		{
			English: "book", Gender: "neuter",
			NominativeSg: "βιβλίο", GenitiveSg: "βιβλίου", AccusativeSg: "βιβλίο",
			NominativePl: "βιβλία", GenitivePl: "βιβλίων", AccusativePl: "βιβλία",
			NomSgArticle: "το", GenSgArticle: "του", AccSgArticle: "το",
			NomPlArticle: "τα", GenPlArticle: "των", AccPlArticle: "τα",
		},
	}

	for _, noun := range nouns {
		if err := repo.CreateNoun(noun); err != nil {
			t.Fatalf("CreateNoun() error = %v", err)
		}
	}

	// List all nouns
	list, err := repo.ListNouns()
	if err != nil {
		t.Fatalf("ListNouns() error = %v", err)
	}

	if len(list) != 2 {
		t.Errorf("Expected 2 nouns, got %d", len(list))
	}

	// Verify order by ID
	if list[0].ID > list[1].ID {
		t.Error("Expected nouns to be ordered by ID")
	}
}

func TestCreateAndGetSentence(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Create a noun first
	noun := &models.Noun{
		English: "teacher", Gender: "masculine",
		NominativeSg: "δάσκαλος", GenitiveSg: "δασκάλου", AccusativeSg: "δάσκαλο",
		NominativePl: "δάσκαλοι", GenitivePl: "δασκάλων", AccusativePl: "δασκάλους",
		NomSgArticle: "ο", GenSgArticle: "του", AccSgArticle: "τον",
		NomPlArticle: "οι", GenPlArticle: "των", AccPlArticle: "τους",
	}
	if err := repo.CreateNoun(noun); err != nil {
		t.Fatal(err)
	}

	// Create sentence
	sentence := &models.Sentence{
		NounID:          noun.ID,
		EnglishPrompt:   "I see ___ (the teacher)",
		GreekSentence:   "Βλέπω τον δάσκαλο",
		CorrectAnswer:   "τον δάσκαλο",
		CaseType:        "accusative",
		Number:          "singular",
		DifficultyPhase: 1,
		ContextType:     "direct_object",
		Preposition:     nil,
	}

	err := repo.CreateSentence(sentence)
	if err != nil {
		t.Fatalf("CreateSentence() error = %v", err)
	}

	if sentence.ID == 0 {
		t.Error("Expected sentence ID to be set after creation")
	}
}

func TestGetRandomSentences(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Create a noun
	noun := &models.Noun{
		English: "teacher", Gender: "masculine",
		NominativeSg: "δάσκαλος", GenitiveSg: "δασκάλου", AccusativeSg: "δάσκαλο",
		NominativePl: "δάσκαλοι", GenitivePl: "δασκάλων", AccusativePl: "δασκάλους",
		NomSgArticle: "ο", GenSgArticle: "του", AccSgArticle: "τον",
		NomPlArticle: "οι", GenPlArticle: "των", AccPlArticle: "τους",
	}
	if err := repo.CreateNoun(noun); err != nil {
		t.Fatal(err)
	}

	// Create sentences with different phases and numbers
	sentences := []*models.Sentence{
		{
			NounID: noun.ID, EnglishPrompt: "Test 1", GreekSentence: "Test",
			CorrectAnswer: "answer1", CaseType: "accusative", Number: "singular",
			DifficultyPhase: 1, ContextType: "direct_object",
		},
		{
			NounID: noun.ID, EnglishPrompt: "Test 2", GreekSentence: "Test",
			CorrectAnswer: "answer2", CaseType: "genitive", Number: "singular",
			DifficultyPhase: 2, ContextType: "possession",
		},
		{
			NounID: noun.ID, EnglishPrompt: "Test 3", GreekSentence: "Test",
			CorrectAnswer: "answer3", CaseType: "accusative", Number: "plural",
			DifficultyPhase: 1, ContextType: "direct_object",
		},
	}

	for _, s := range sentences {
		if err := repo.CreateSentence(s); err != nil {
			t.Fatal(err)
		}
	}

	// Test filtering by phase and number
	result, err := repo.GetRandomSentences(1, "singular", 10)
	if err != nil {
		t.Fatalf("GetRandomSentences() error = %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 sentence (phase 1, singular), got %d", len(result))
	}

	// Test including all numbers (empty string)
	result, err = repo.GetRandomSentences(1, "", 10)
	if err != nil {
		t.Fatalf("GetRandomSentences() error = %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 sentences (phase 1, all numbers), got %d", len(result))
	}
}

func TestCreateAndGetExplanation(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Create noun and sentence first
	noun := &models.Noun{
		English: "teacher", Gender: "masculine",
		NominativeSg: "δάσκαλος", GenitiveSg: "δασκάλου", AccusativeSg: "δάσκαλο",
		NominativePl: "δάσκαλοι", GenitivePl: "δασκάλων", AccusativePl: "δασκάλους",
		NomSgArticle: "ο", GenSgArticle: "του", AccSgArticle: "τον",
		NomPlArticle: "οι", GenPlArticle: "των", AccPlArticle: "τους",
	}
	if err := repo.CreateNoun(noun); err != nil {
		t.Fatal(err)
	}

	sentence := &models.Sentence{
		NounID: noun.ID, EnglishPrompt: "Test", GreekSentence: "Test",
		CorrectAnswer: "answer", CaseType: "accusative", Number: "singular",
		DifficultyPhase: 1, ContextType: "direct_object",
	}
	if err := repo.CreateSentence(sentence); err != nil {
		t.Fatal(err)
	}

	// Create explanation
	explanation := &models.Explanation{
		SentenceID:    sentence.ID,
		Translation:   "I see the teacher",
		SyntacticRole: "Direct object requires accusative",
		Morphology:    "ο δάσκαλος → τον δάσκαλο",
	}

	err := repo.CreateExplanation(explanation)
	if err != nil {
		t.Fatalf("CreateExplanation() error = %v", err)
	}

	// Get explanation
	retrieved, err := repo.GetExplanationBySentenceID(sentence.ID)
	if err != nil {
		t.Fatalf("GetExplanationBySentenceID() error = %v", err)
	}

	if retrieved.Translation != explanation.Translation {
		t.Errorf("Expected Translation %q, got %q", explanation.Translation, retrieved.Translation)
	}
}

func TestForeignKeyConstraint(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Try to create sentence with non-existent noun ID
	sentence := &models.Sentence{
		NounID: 999, // Non-existent
		EnglishPrompt: "Test", GreekSentence: "Test",
		CorrectAnswer: "answer", CaseType: "accusative", Number: "singular",
		DifficultyPhase: 1, ContextType: "direct_object",
	}

	err := repo.CreateSentence(sentence)
	if err == nil {
		t.Error("Expected error for foreign key constraint violation, got nil")
	}
}
