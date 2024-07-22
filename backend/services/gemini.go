package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"read-robin/models"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

const (
	location                    = "northamerica-northeast1"
	modelName                   = "gemini-1.5-pro"
	quizModelSystemInstructions = `You are a highly skilled model that generates quiz questions and answers from summarized content tailored for a specific user persona. The persona details include Name, Role (profession, age, etc.), Language, and Difficulty (beginner, intermediate, expert). Your task is to generate questions and answers based on the summarized content provided, considering the persona details. You should also generate a small piece of reference text that was used to create your question/answer pair. Omit any backticks or format reference. Return everything in a JSON dictionary with 'quiz' being an array of objects containing 'question', 'answer', and 'reference' strings. The structure should look like this:
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
	webscrapeModelSystemInstructions = `You are a highly skilled model that extracts readable text from HTML content and generates a title for the content. Your task is to extract the given HTML content and output it into a clear and concise article, ignoring any unnecessary HTML tags or irrelevant content. Additionally, generate a title from the URL to objectively define the site's host and page names (e.g., www.example.com would be Example, and https://en.wikipedia.org/wiki/The_World%27s_Largest_Lobster would be Wikipedia - The World's Largest Lobster). Return everything in a JSON dictionary with 'content' and 'title' keys. Exclude any markdown code fences in your response. The structure should look like this:
{
	"content": "extracted content",
	"title": "generated title"
}`
	reviewModelSystemInstructions = `You are a highly skilled model that reviews quiz responses. Your task is to determine if the user's response captures the essence of the expected answer based on the reference provided. As long as the user's response includes the key points or main ideas of the expected answer, it should be considered a "PASS". Yes or No is an acceptable response for yes and no questions. Use a lenient approach, focusing on the main concepts rather than exact wording. Return only "PASS" or "FAIL" as the response.

	Examples:
	1. Expected Answer: "The 'Example Domain' is for use in illustrative examples in documents."
	   User Response: "Example Domain is used in documents for examples."
	   Response: "PASS"
	
	2. Expected Answer: "The 'Example Domain' is for use in illustrative examples in documents."
	   User Response: "It is a domain for examples."
	   Response: "PASS"
	
	3. Expected Answer: "The 'Example Domain' is for use in illustrative examples in documents."
	   User Response: "This domain is for examples in documents."
	   Response: "PASS"
	
	4. Expected Answer: "The 'Example Domain' is for use in illustrative examples in documents."
	   User Response: "It is a domain used in documents."
	   Response: "FAIL"`
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

// ExtractContent extracts the given HTML text using the Gemini model and returns both the content and title
func (gc *GeminiClient) ExtractContent(ctx context.Context, htmlText string) (map[string]string, string, error) {
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

// GenerateQuiz generates quiz questions and answers from the summarized content
func (gc *GeminiClient) GenerateQuiz(ctx context.Context, summarizedContent, personaName, personaRole, personaLanguage, personaDifficulty string) (string, string, error) {
	promptText := fmt.Sprintf("Generate a quiz for a %s (%s) at %s difficulty level based on the following content: %s", personaRole, personaLanguage, personaDifficulty, summarizedContent)
	return gc.generateContent(ctx, quizModelSystemInstructions, promptText)
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
	reviewResult, _, err := gc.generateContent(ctx, reviewModelSystemInstructions, reviewData)
	if err != nil {
		return "", fmt.Errorf("error reviewing response: %w", err)
	}
	return strings.TrimSpace(reviewResult), nil
}
