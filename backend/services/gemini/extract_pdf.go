package gemini

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"cloud.google.com/go/vertexai/genai"
)

const (
	pdfModelSystemInstructions = `You are a highly skilled model that extracts readable text from PDF content and generates a title for the content. Your task is to extract the given PDF content and output it into a clear and concise article, ignoring any unnecessary formatting or irrelevant content. Additionally, generate a title that objectively defines the main topic of the PDF. Return everything in a JSON dictionary with 'content' and 'title' keys, omit any markdown backticks. The structure should look like this:
    {
        "content": "extracted content",
        "title": "generated title"
    }`
)

type pdfPrompt struct {
	pdfPath string
}

// extractContentFromPDF extracts readable text and title from PDF content using the Gemini model
func (gc *GeminiClient) ExtractContentFromPdf(ctx context.Context, pdfPath string) (map[string]string, string, error) {
	model := gc.client.GenerativeModel(modelName)

	part := genai.FileData{
		MIMEType: "application/pdf",
		FileURI:  pdfPath,
	}

	fmt.Printf("Extracting content from PDF: %s\n", pdfPath)
	res, err := model.GenerateContent(ctx, genai.Text(pdfModelSystemInstructions), part)
	if err != nil {
		return nil, "", fmt.Errorf("unable to generate contents: %w", err)
	}

	if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
		return nil, "", errors.New("empty response from model")
	}

	content := res.Candidates[0].Content.Parts[0]

	contentText := fmt.Sprintf("%s", content)

	// Parse the JSON response to extract content and title
	var contentMap map[string]string
	if err := json.Unmarshal([]byte(contentText), &contentMap); err != nil {
		return nil, "", fmt.Errorf("json.Unmarshal: %w", err)
	}

	// Convert the response to a readable format
	fullResponse, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return nil, "", fmt.Errorf("json.MarshalIndent: %w", err)
	}

	return contentMap, string(fullResponse), nil
}

func (gc *GeminiClient) GenerateQuizFromPDF(ctx context.Context, pdfPath, personaName, personaRole, personaLanguage, personaDifficulty string) (string, error) {
	prompt := pdfPrompt{
		pdfPath: pdfPath,
	}

	contentMap, fullHTML, err := gc.ExtractContentFromPdf(ctx, prompt.pdfPath)
	if err != nil {
		return "", fmt.Errorf("error generating content from PDF: %w", err)
	}

	fmt.Printf("Generating quiz from PDF content: %s\n", contentMap)
	fmt.Print(fullHTML)
	promptText := fmt.Sprintf("Generate a quiz for a %s (%s) at %s difficulty level based on the following content: %s", personaRole, personaLanguage, personaDifficulty, contentMap)
	quizContent, _, err := gc.generateContent(ctx, quizModelSystemInstructions, promptText)
	if err != nil {
		return "", fmt.Errorf("error generating quiz: %w", err)
	}

	return quizContent, nil
}
