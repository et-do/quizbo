package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"cloud.google.com/go/vertexai/genai"
)

const (
	location  = "northamerica-northeast1"
	modelName = "gemini-1.5-flash-001"
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

// ExtractContent extracts the given HTML text using the Gemini model
func (gc *GeminiClient) ExtractContent(ctx context.Context, htmlText string) (string, error) {
	systemInstructions := "You are a highly skilled model that extracts readable text from HTML content. Your task is to extract the given HTML content and output into a clear and concise article, ignoring any unnecessary HTML tags or irrelevant content."
	geminiModel := gc.client.GenerativeModel(modelName)
	geminiModel.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemInstructions)},
	}

	prompt := genai.Text(htmlText)

	resp, err := geminiModel.GenerateContent(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("error generating content: %w", err)
	}

	// Extract the summary from the response
	summary, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json.MarshalIndent: %w", err)
	}

	return string(summary), nil
}
