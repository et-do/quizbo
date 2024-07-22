package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

func (gc *GeminiClient) generateContent(ctx context.Context, systemInstructions, promptText string) (string, string, error) {
	geminiModel := gc.client.GenerativeModel(modelName)
	geminiModel.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemInstructions)},
	}

	prompt := genai.Text(promptText)

	fmt.Printf("Generating content with prompt: %s\n", promptText)
	resp, err := geminiModel.GenerateContent(ctx, prompt)
	if err != nil {
		return "", "", fmt.Errorf("error generating content: %w", err)
	}

	fullResponse, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return "", "", fmt.Errorf("json.MarshalIndent: %w", err)
	}

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
