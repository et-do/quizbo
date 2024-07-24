package gemini

import (
	"context"
	"os"
	"read-robin/models"
	"testing"
)

func TestGenerateQuiz(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Ensure the environment variable is set for the test
	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		t.Fatal("GCP_PROJECT environment variable not set")
	}

	geminiClient, err := NewGeminiClient(ctx)
	if err != nil {
		t.Fatalf("NewGeminiClient: expected no error, got %v", err)
	}

	// Use the previously tested method to get the extracted content
	contentMap, _, err := geminiClient.ExtractContentFromHtml(ctx, testHTML)
	if err != nil {
		t.Fatalf("ExtractContent: expected no error, got %v", err)
	}

	// Define a test persona
	testPersona := models.Persona{
		ID:         "test_persona_id",
		Name:       "Test User",
		Role:       "Student",
		Language:   "English",
		Difficulty: "Intermediate",
	}

	quiz, fullQuiz, err := geminiClient.GenerateQuiz(ctx, contentMap["content"], testPersona)
	if err != nil {
		t.Fatalf("GenerateQuiz: expected no error, got %v", err)
	}

	if quiz == "" {
		t.Errorf("GenerateQuiz: expected a quiz, got an empty string")
	} else {
		t.Logf("Quiz: %s", quiz)
		t.Logf("Full Quiz Response: %s", fullQuiz)
	}
}
