package gemini

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	webscrapeModelSystemInstructions = `You are a highly skilled model that extracts readable text from HTML content and generates a title for the content. Your task is to extract the given HTML content and output it into a clear and concise article, ignoring any unnecessary HTML tags or irrelevant content. Additionally, generate a title from the URL to objectively define the site's host and page names (e.g., www.example.com would be Example, and https://en.wikipedia.org/wiki/The_World%27s_Largest_Lobster would be Wikipedia - The World's Largest Lobster). Return everything in a JSON dictionary with 'content' and 'title' keys. Exclude any markdown code fences in your response. The structure should look like this:
{
	"content": "extracted content",
	"title": "generated title"
}`
)

// ExtractContentFromHtml extracts the given HTML text using the Gemini model and returns both the content and title
func (gc *GeminiClient) ExtractContentFromHtml(ctx context.Context, htmlText string) (map[string]string, string, error) {
	extractedContent, fullResponse, err := gc.generateContent(ctx, webscrapeModelSystemInstructions, htmlText)
	if err != nil {
		return nil, "", fmt.Errorf("error extracting content: %w", err)
	}

	var contentMap map[string]string
	if err := json.Unmarshal([]byte(extractedContent), &contentMap); err != nil {
		return nil, "", fmt.Errorf("json.Unmarshal: %w", err)
	}

	return contentMap, fullResponse, nil
}
