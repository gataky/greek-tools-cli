package importer

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// CSVRow represents a single row from the CSV file
type CSVRow struct {
	English string
	Greek   string
	Gender  string
	RowNum  int // For checkpoint tracking
}

// ValidateGender checks if the gender value is valid
func ValidateGender(gender string) error {
	gender = strings.ToLower(strings.TrimSpace(gender))
	validGenders := map[string]bool{
		"masculine":  true,
		"feminine":   true,
		"neuter":     true,
		"invariable": true,
	}

	if !validGenders[gender] {
		return fmt.Errorf("invalid gender '%s', must be one of: masculine, feminine, neuter, invariable", gender)
	}
	return nil
}

// ParseCSV reads and validates a CSV file
// Returns a slice of CSVRow structs
func ParseCSV(filepath string) ([]CSVRow, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read header
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV headers: %w", err)
	}

	// Validate headers
	if len(headers) < 3 {
		return nil, fmt.Errorf("CSV must have at least 3 columns (english, greek, attribute)")
	}

	// Find column indices
	englishIdx := -1
	greekIdx := -1
	genderIdx := -1

	for i, header := range headers {
		header = strings.ToLower(strings.TrimSpace(header))
		switch header {
		case "english":
			englishIdx = i
		case "greek":
			greekIdx = i
		case "attribute":
			genderIdx = i
		}
	}

	if englishIdx == -1 || greekIdx == -1 || genderIdx == -1 {
		return nil, fmt.Errorf("CSV must have 'english', 'greek', and 'attribute' columns")
	}

	// Read all rows
	rows := []CSVRow{}
	rowNum := 1 // Start at 1 (header is row 0)

	for {
		record, err := reader.Read()
		if err != nil {
			// EOF is expected
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("error reading CSV at row %d: %w", rowNum+1, err)
		}

		rowNum++

		// Validate row has enough columns
		if len(record) <= englishIdx || len(record) <= greekIdx || len(record) <= genderIdx {
			return nil, fmt.Errorf("row %d has missing columns", rowNum)
		}

		english := strings.TrimSpace(record[englishIdx])
		greek := strings.TrimSpace(record[greekIdx])
		gender := strings.TrimSpace(record[genderIdx])

		// Validate required fields
		if english == "" {
			return nil, fmt.Errorf("row %d: 'english' field is empty", rowNum)
		}
		if greek == "" {
			return nil, fmt.Errorf("row %d: 'greek' field is empty", rowNum)
		}
		if gender == "" {
			return nil, fmt.Errorf("row %d: 'attribute' (gender) field is empty", rowNum)
		}

		// Normalize gender to lowercase
		gender = strings.ToLower(gender)

		// Validate gender
		if err := ValidateGender(gender); err != nil {
			return nil, fmt.Errorf("row %d: %w", rowNum, err)
		}

		rows = append(rows, CSVRow{
			English: english,
			Greek:   greek,
			Gender:  gender,
			RowNum:  rowNum,
		})
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("CSV file contains no data rows")
	}

	return rows, nil
}
