package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// DeclensionResponse represents the API response for declension generation
type DeclensionResponse struct {
	NominativeSg string `json:"nominative_sg"`
	NomSgArticle string `json:"nom_sg_article"`
	GenitiveSg   string `json:"genitive_sg"`
	GenSgArticle string `json:"gen_sg_article"`
	AccusativeSg string `json:"accusative_sg"`
	AccSgArticle string `json:"acc_sg_article"`
	NominativePl string `json:"nominative_pl"`
	NomPlArticle string `json:"nom_pl_article"`
	GenitivePl   string `json:"genitive_pl"`
	GenPlArticle string `json:"gen_pl_article"`
	AccusativePl string `json:"accusative_pl"`
	AccPlArticle string `json:"acc_pl_article"`
}

// SentenceResponse represents a single practice sentence
type SentenceResponse struct {
	EnglishPrompt   string  `json:"english_prompt"`
	GreekSentence   string  `json:"greek_sentence"`
	CorrectAnswer   string  `json:"correct_answer"`
	CaseType        string  `json:"case_type"`
	Number          string  `json:"number"`
	DifficultyPhase int     `json:"difficulty_phase"`
	ContextType     string  `json:"context_type"`
	Preposition     *string `json:"preposition"`
}

// ClaudeClient wraps the Anthropic SDK client
type ClaudeClient struct {
	client *anthropic.Client
	model  string
}

// NewClaudeClient creates a new Claude API client
// Reads ANTHROPIC_API_KEY from environment variable
func NewClaudeClient() (*ClaudeClient, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
	}

	client := anthropic.NewClient(option.WithAPIKey(apiKey))

	return &ClaudeClient{
		client: &client,
		model:  "claude-sonnet-4-6",
	}, nil
}

// callAPI makes an API call with the given prompt and parses the JSON response
func (c *ClaudeClient) callAPI(ctx context.Context, prompt string) (string, error) {
	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(c.model),
		MaxTokens: 2000,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})

	if err != nil {
		return "", fmt.Errorf("API call failed: %w", err)
	}

	// Extract text from response
	if len(message.Content) == 0 {
		return "", fmt.Errorf("empty response from API")
	}

	// Get the text content from the first content block
	contentBlock := message.Content[0]
	if contentBlock.Type == "text" && contentBlock.Text != "" {
		return cleanJSONResponse(contentBlock.Text), nil
	}

	return "", fmt.Errorf("unexpected response type from API: %s", contentBlock.Type)
}

// cleanJSONResponse removes markdown code fences from JSON responses
func cleanJSONResponse(text string) string {
	// Remove leading/trailing whitespace
	text = strings.TrimSpace(text)

	// Check for markdown code fences with json language identifier
	if strings.HasPrefix(text, "```json") {
		text = strings.TrimPrefix(text, "```json")
		text = strings.TrimSuffix(text, "```")
		text = strings.TrimSpace(text)
		return text
	}

	// Check for plain markdown code fences
	if strings.HasPrefix(text, "```") {
		text = strings.TrimPrefix(text, "```")
		text = strings.TrimSuffix(text, "```")
		text = strings.TrimSpace(text)
		return text
	}

	return text
}

// logError logs errors to import.log file
func (c *ClaudeClient) logError(format string, args ...any) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	logPath := filepath.Join(homeDir, ".greekmaster", "import.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[%s] "+format+"\n", append([]any{timestamp}, args...)...)
	f.WriteString(logMsg)
}

// GenerateDeclensions generates all declined forms for a Greek noun
func (c *ClaudeClient) GenerateDeclensions(greek, english, gender string) (*DeclensionResponse, error) {
	prompt := GenerateDeclensionPrompt(greek, english, gender)

	var response *DeclensionResponse
	err := RetryWithBackoff(func() error {
		ctx := context.Background()
		text, err := c.callAPI(ctx, prompt)
		if err != nil {
			c.logError("Declension API call failed for '%s': %v", greek, err)
			return err
		}

		// Parse JSON response
		var decl DeclensionResponse
		if err := json.Unmarshal([]byte(text), &decl); err != nil {
			c.logError("Failed to parse declension JSON for '%s': %v\nResponse: %s", greek, err, text)
			return fmt.Errorf("invalid JSON response: %w", err)
		}

		response = &decl
		return nil
	}, 3) // Max 3 retries

	if err != nil {
		return nil, err
	}

	return response, nil
}

// GenerateSentences generates practice sentences for a noun
func (c *ClaudeClient) GenerateSentences(greek, english, gender string) ([]SentenceResponse, error) {
	prompt := GenerateSentencesPrompt(greek, english, gender)

	var response []SentenceResponse
	err := RetryWithBackoff(func() error {
		ctx := context.Background()
		text, err := c.callAPI(ctx, prompt)
		if err != nil {
			c.logError("Sentences API call failed for '%s': %v", greek, err)
			return err
		}

		// Parse JSON response
		var sentences []SentenceResponse
		if err := json.Unmarshal([]byte(text), &sentences); err != nil {
			c.logError("Failed to parse sentences JSON for '%s': %v\nResponse: %s", greek, err, text)
			return fmt.Errorf("invalid JSON response: %w", err)
		}

		response = sentences
		return nil
	}, 3) // Max 3 retries

	if err != nil {
		return nil, err
	}

	return response, nil
}

// GenerateExplanations generates explanations for multiple sentences
