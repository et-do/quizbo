package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

const (
	modelName = "gemini-1.5-pro"
)

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
