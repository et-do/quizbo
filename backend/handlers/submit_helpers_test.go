package handlers

import (
	"context"
	"read-robin/models"
	"read-robin/services"
	"read-robin/services/gemini"
	"read-robin/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestProcessURLSubmission(t *testing.T) {
	ctx := context.Background()

	// Mock HTML fetcher
	utils.FetchHTML = func(url string) (string, error) {
		return "<html><body>Test HTML content</body></html>", nil
	}

	// Mock Firestore client
	firestoreClient := &services.FirestoreClient{}
	firestoreClient.GetExistingQuizzes = func(ctx context.Context, contentID string) ([]models.Quiz, error) {
		return nil, status.Error(codes.NotFound, "not found")
	}
	firestoreClient.SaveQuiz = func(ctx context.Context, normalizedURL, title string, quiz models.Quiz) error {
		return nil
	}

	// Mock Gemini client
	geminiClient := &gemini.GeminiClient{}
	geminiClient.ExtractAndGenerateQuiz = func(ctx context.Context, htmlContent string, persona models.Persona) (map[string]interface{}, string, error) {
		return map[string]interface{}{"content": "quiz content"}, "Test Title", nil
	}

	urlRequest := models.URLRequest{
		URL: "http://example.com",
		Persona: models.Persona{
			Name:       "Test Persona",
			Role:       "Test Role",
			Language:   "English",
			Difficulty: "Easy",
		},
	}

	response, err := processURLSubmission(ctx, urlRequest, geminiClient, firestoreClient)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "http://example.com", response.URL)
	assert.NotEmpty(t, response.ContentID)
	assert.NotEmpty(t, response.QuizID)
	assert.Equal(t, "Test Title", response.Title)
	assert.True(t, response.IsFirstQuiz)
}

func TestProcessPDFSubmission(t *testing.T) {
	ctx := context.Background()

	// Mock Firestore client
	firestoreClient := &services.FirestoreClient{}
	firestoreClient.GetExistingQuizzes = func(ctx context.Context, contentID string) ([]models.Quiz, error) {
		return nil, status.Error(codes.NotFound, "not found")
	}
	firestoreClient.SaveQuiz = func(ctx context.Context, normalizedURL, title string, quiz models.Quiz) error {
		return nil
	}

	// Mock Gemini client
	geminiClient := &gemini.GeminiClient{}
	geminiClient.GenerateQuizFromPDF = func(ctx context.Context, pdfURL, personaName, personaRole, personaLanguage, personaDifficulty string) (string, error) {
		return `{"content": "quiz content"}`, nil
	}

	pdfURL := "gs://test-bucket/test.pdf"
	persona := models.Persona{
		Name:       "Test Persona",
		Role:       "Test Role",
		Language:   "English",
		Difficulty: "Easy",
	}

	response, err := processPDFSubmission(ctx, pdfURL, persona, geminiClient, firestoreClient)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Empty(t, response.URL)
	assert.NotEmpty(t, response.ContentID)
	assert.NotEmpty(t, response.QuizID)
	assert.Empty(t, response.Title)
	assert.True(t, response.IsFirstQuiz)
}
