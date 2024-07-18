package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

// ReviewResponse reviews the user's response using the Gemini model
func (gc *GeminiClient) ReviewResponse(ctx context.Context, reviewData string) (string, error) {
	reviewResult, _, err := gc.generateContent(ctx, instructions.ReviewModelSystemInstructions, reviewData)
	if err != nil {
		return "", fmt.Errorf("error reviewing response: %w", err)
	}
	return strings.TrimSpace(reviewResult), nil
}

type pdfPrompt struct {
	// pdfPath is a Google Cloud Storage path starting with "gs://"
	pdfPath string
	// question asked to the model
	question string
}

// GenerateContentFromPDF generates a response based on the provided PDF asset and question
func (gc *GeminiClient) generateContentFromPDF(ctx context.Context, w io.Writer, prompt pdfPrompt, modelName string) (string, error) {
	model := gc.client.GenerativeModel(modelName)

	part := genai.FileData{
		MIMEType: "application/pdf",
		FileURI:  prompt.pdfPath,
	}

	res, err := model.GenerateContent(ctx, part, genai.Text(prompt.question))
	if err != nil {
		return "", fmt.Errorf("unable to generate contents: %w", err)
	}

	if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("empty response from model")
	}

	fmt.Fprintf(w, "generated response: %s\n", res.Candidates[0].Content.Parts[0])
	content := res.Candidates[0].Content.Parts[0]

	// Use fmt.Sprintf to convert genai.Part to string
	contentText := fmt.Sprintf("%s", content)

	return contentText, nil
}

// GenerateQuizFromPDF extracts content from a PDF and generates quiz questions
func (gc *GeminiClient) GenerateQuizFromPDF(ctx context.Context, pdfPath, question, personaName, personaRole, personaLanguage, personaDifficulty string) (string, error) {
	// Create a pdfPrompt object
	prompt := pdfPrompt{
		pdfPath:  pdfPath,
		question: question,
	}

	// Use an io.Writer to capture the output
	var output strings.Builder

	content, err := gc.generateContentFromPDF(ctx, &output, prompt, modelName)
	if err != nil {
		return "", fmt.Errorf("error generating content from PDF: %w", err)
	}

	promptText := fmt.Sprintf("Generate a quiz for a %s (%s) at %s difficulty level based on the following content: %s", personaRole, personaLanguage, personaDifficulty, content)
	quizContent, _, err := gc.generateContent(ctx, instructions.QuizModelSystemInstructions, promptText)
	if err != nil {
		return "", fmt.Errorf("error generating quiz: %w", err)
	}

	return quizContent, nil
}
