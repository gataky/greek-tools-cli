package importer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateGender(t *testing.T) {
	tests := []struct {
		name    string
		gender  string
		wantErr bool
	}{
		{"masculine valid", "masculine", false},
		{"feminine valid", "feminine", false},
		{"neuter valid", "neuter", false},
		{"invariable valid", "invariable", false},
		{"uppercase masculine", "MASCULINE", false},
		{"mixed case", "Feminine", false},
		{"with spaces", "  neuter  ", false},
		{"invalid gender", "unknown", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGender(tt.gender)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGender(%q) error = %v, wantErr %v", tt.gender, err, tt.wantErr)
			}
		})
	}
}

func TestParseCSV_Valid(t *testing.T) {
	// Create a temporary CSV file
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")

	content := `english,greek,attribute
teacher,δάσκαλος,masculine
book,βιβλίο,neuter`

	if err := os.WriteFile(csvPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	rows, err := ParseCSV(csvPath)
	if err != nil {
		t.Fatalf("ParseCSV() error = %v", err)
	}

	if len(rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(rows))
	}

	// Check first row
	if rows[0].English != "teacher" {
		t.Errorf("Expected english 'teacher', got %q", rows[0].English)
	}
	if rows[0].Greek != "δάσκαλος" {
		t.Errorf("Expected greek 'δάσκαλος', got %q", rows[0].Greek)
	}
	if rows[0].Gender != "masculine" {
		t.Errorf("Expected gender 'masculine', got %q", rows[0].Gender)
	}
}

func TestParseCSV_MissingHeaders(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")

	content := `english,greek
teacher,δάσκαλος`

	if err := os.WriteFile(csvPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ParseCSV(csvPath)
	if err == nil {
		t.Error("Expected error for missing 'attribute' column, got nil")
	}
}

func TestParseCSV_EmptyField(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")

	content := `english,greek,attribute
teacher,,masculine`

	if err := os.WriteFile(csvPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ParseCSV(csvPath)
	if err == nil {
		t.Error("Expected error for empty 'greek' field, got nil")
	}
}

func TestParseCSV_InvalidGender(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")

	content := `english,greek,attribute
teacher,δάσκαλος,invalid_gender`

	if err := os.WriteFile(csvPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ParseCSV(csvPath)
	if err == nil {
		t.Error("Expected error for invalid gender, got nil")
	}
}

func TestParseCSV_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")

	content := `english,greek,attribute`

	if err := os.WriteFile(csvPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ParseCSV(csvPath)
	if err == nil {
		t.Error("Expected error for CSV with no data rows, got nil")
	}
}

func TestParseCSV_UnicodeHandling(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")

	content := `english,greek,attribute
woman,γυναίκα,feminine
student,μαθητής,masculine`

	if err := os.WriteFile(csvPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	rows, err := ParseCSV(csvPath)
	if err != nil {
		t.Fatalf("ParseCSV() error = %v", err)
	}

	if len(rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(rows))
	}

	// Verify Greek Unicode was preserved
	if rows[0].Greek != "γυναίκα" {
		t.Errorf("Expected greek 'γυναίκα', got %q", rows[0].Greek)
	}
	if rows[1].Greek != "μαθητής" {
		t.Errorf("Expected greek 'μαθητής', got %q", rows[1].Greek)
	}
}
