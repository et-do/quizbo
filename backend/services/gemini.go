package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

const (
	location                    = "northamerica-northeast1"
	modelName                   = "gemini-1.5-pro"
	quizModelSystemInstructions = `You are a highly skilled model that generates quiz questions and answers from summarized content. Your task is to generate questions and answers based on the summarized content provided. You should also generate a small piece of reference text that was used to create your question/answer pair. Omit any backticks or format reference. Return everything in a JSON dictionary with 'quiz' being an array of objects containing 'question', 'answer', and 'reference' strings. The structure should look like this:
{
	"quiz": [
		{
			"question": "question",
			"answer": "answer",
			"reference": "reference"
		},
		{
			"question": "question",
			"answer": "answer",
			"reference": "reference"
		}
	]
}`
	webscrapeModelSystemInstructions = "You are a highly skilled model that extracts readable text from HTML content. Your task is to extract the given HTML content and output into a clear and concise article, ignoring any unnecessary HTML tags or irrelevant content."
)

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

// ExtractContent extracts the given HTML text using the Gemini model
func (gc *GeminiClient) ExtractContent(ctx context.Context, htmlText string) (string, string, error) {
	return gc.generateContent(ctx, webscrapeModelSystemInstructions, htmlText)
}

// GenerateQuiz generates quiz questions and answers from the summarized content
func (gc *GeminiClient) GenerateQuiz(ctx context.Context, summarizedContent string) (string, string, error) {
	return gc.generateContent(ctx, quizModelSystemInstructions, summarizedContent)
}

// ExtractAndGenerateQuiz extracts content and generates a quiz using the Gemini client
func (gc *GeminiClient) ExtractAndGenerateQuiz(ctx context.Context, htmlContent string) (map[string]interface{}, error) {
	extractedContent, _, err := gc.ExtractContent(ctx, htmlContent)
	if err != nil {
		return nil, err
	}
	quizContent, _, err := gc.GenerateQuiz(ctx, extractedContent)
	if err != nil {
		return nil, err
	}
	var quizContentMap map[string]interface{}
	if err := json.Unmarshal([]byte(quizContent), &quizContentMap); err != nil {
		return nil, err
	}
	return quizContentMap, nil
}
