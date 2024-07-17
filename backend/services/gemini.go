package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"read-robin/models"
	"strings"

	"cloud.google.com/go/vertexai/genai"
	"github.com/BurntSushi/toml"
)

const (
	location  = "northamerica-northeast1"
	modelName = "gemini-1.5-pro"
)

type SystemInstructions struct {
	QuizModelSystemInstructions      string `toml:"quizModelSystemInstructions"`
	WebscrapeModelSystemInstructions string `toml:"webscrapeModelSystemInstructions"`
	PDFModelSystemInstructions       string `toml:"pdfModelSystemInstructions"`
	ReviewModelSystemInstructions    string `toml:"reviewModelSystemInstructions"`
}

var instructions SystemInstructions

func init() {
	configFile := "gemini_system_instructions.toml"
	configData, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	err = toml.Unmarshal(configData, &instructions)
	if err != nil {
		fmt.Printf("Error unmarshaling config file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Loaded System Instructions: %+v\n", instructions) // Debug log
}

// GeminiClient is a wrapper around the Vertex AI GenAI client
type GeminiClient struct {
	client *genai.Client
}

// NewGeminiClient creates a new GeminiClient
func NewGeminiClient(ctx context.Context) (*GeminiClient, error) {
	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		return nil, fmt.Errorf("GCP_PROJECT environment variable not set")
	}

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	return &GeminiClient{client: client}, nil
}

// Helper function to generate content using Gemini model
func (gc *GeminiClient) generateContent(ctx context.Context, systemInstructions, promptText string) (string, string, error) {
	fmt.Printf("Generating content with system instructions: %s\n", systemInstructions) // Debug log
	geminiModel := gc.client.GenerativeModel(modelName)
	geminiModel.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemInstructions)},
	}

	prompt := genai.Text(promptText)

	resp, err := geminiModel.GenerateContent(ctx, prompt)
	if err != nil {
		return "", "", fmt.Errorf("error generating content: %w", err)
	}

	// Extract the full response as JSON
	fullResponse, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return "", "", fmt.Errorf("json.MarshalIndent: %w", err)
	}

	// Debug: Print the full response structure
	fmt.Printf("Full response: %s\n", fullResponse)

	// Extract the text from the parts
	var partContent strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				partContent.WriteString(fmt.Sprintf("%s", part))
			}
		}
	}

	return partContent.String(), string(fullResponse), nil
}

// ExtractContent extracts the given HTML text using the Gemini model and returns both the content and title
func (gc *GeminiClient) ExtractContent(ctx context.Context, htmlText string) (map[string]string, string, error) {
	extractedContent, fullResponse, err := gc.generateContent(ctx, instructions.WebscrapeModelSystemInstructions, htmlText)
	if err != nil {
		return nil, "", fmt.Errorf("error extracting content: %w", err)
	}

	var contentMap map[string]string
	if err := json.Unmarshal([]byte(extractedContent), &contentMap); err != nil {
		return nil, "", fmt.Errorf("json.Unmarshal: %w", err)
	}

	return contentMap, fullResponse, nil
}

// ExtractPDFContent extracts the given PDF text using the Gemini model and returns both the content and title
func (gc *GeminiClient) ExtractPDFContent(ctx context.Context, pdfText string) (map[string]string, string, error) {
	extractedContent, fullResponse, err := gc.generateContent(ctx, instructions.PDFModelSystemInstructions, pdfText)
	if err != nil {
		return nil, "", fmt.Errorf("error extracting content: %w", err)
	}

	var contentMap map[string]string
	if err := json.Unmarshal([]byte(extractedContent), &contentMap); err != nil {
		return nil, "", fmt.Errorf("json.Unmarshal: %w", err)
	}

	return contentMap, fullResponse, nil
}

// GenerateQuiz generates quiz questions and answers from the summarized content
func (gc *GeminiClient) GenerateQuiz(ctx context.Context, summarizedContent, personaName, personaRole, personaLanguage, personaDifficulty string) (string, string, error) {
	promptText := fmt.Sprintf("Generate a quiz for a %s (%s) at %s difficulty level based on the following content: %s", personaRole, personaLanguage, personaDifficulty, summarizedContent)
	return gc.generateContent(ctx, instructions.QuizModelSystemInstructions, promptText)
}

// ExtractAndGenerateQuiz extracts content and generates a quiz using the Gemini client
func (gc *GeminiClient) ExtractAndGenerateQuiz(ctx context.Context, htmlContent string, persona models.Persona) (map[string]interface{}, string, error) {
	contentMap, _, err := gc.ExtractContent(ctx, htmlContent)
	if err != nil {
		return nil, "", err
	}

	quizContent, _, err := gc.GenerateQuiz(ctx, contentMap["content"], persona.Name, persona.Role, persona.Language, persona.Difficulty)
	if err != nil {
		return nil, "", err
	}
	var quizContentMap map[string]interface{}
	if err := json.Unmarshal([]byte(quizContent), &quizContentMap); err != nil {
		return nil, "", err
	}

	return quizContentMap, contentMap["title"], nil
}

// ExtractAndGeneratePDFQuiz extracts content from a PDF and generates a quiz using the Gemini client
func (gc *GeminiClient) ExtractAndGeneratePDFQuiz(ctx context.Context, pdfContent string, persona models.Persona) (map[string]interface{}, string, error) {
	contentMap, _, err := gc.ExtractPDFContent(ctx, pdfContent)
	if err != nil {
		return nil, "", err
	}

	quizContent, _, err := gc.GenerateQuiz(ctx, contentMap["content"], persona.Name, persona.Role, persona.Language, persona.Difficulty)
	if err != nil {
		return nil, "", err
	}
	var quizContentMap map[string]interface{}
	if err := json.Unmarshal([]byte(quizContent), &quizContentMap); err != nil {
		return nil, "", err
	}

	return quizContentMap, contentMap["title"], nil
}

// ReviewResponse reviews the user's response using the Gemini model
func (gc *GeminiClient) ReviewResponse(ctx context.Context, reviewData string) (string, error) {
	reviewResult, _, err := gc.generateContent(ctx, instructions.ReviewModelSystemInstructions, reviewData)
	if err != nil {
		return "", fmt.Errorf("error reviewing response: %w", err)
	}
	return strings.TrimSpace(reviewResult), nil
}
