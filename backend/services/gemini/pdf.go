package gemini

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

type pdfPrompt struct {
	pdfPath string
}

func (gc *GeminiClient) extractContentFromPDF(ctx context.Context, w io.Writer, prompt pdfPrompt, modelName string) (string, error) {
	model := gc.client.GenerativeModel(modelName)

	part := genai.FileData{
		MIMEType: "application/pdf",
		FileURI:  prompt.pdfPath,
	}

	fmt.Printf("Extracting content from PDF: %s", prompt.pdfPath)
	res, err := model.GenerateContent(ctx, genai.Text(instructions.PDFModelSystemInstructions), part)
	if err != nil {
		return "", fmt.Errorf("unable to generate contents: %w", err)
	}

	if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("empty response from model")
	}

	fmt.Fprintf(w, "generated response: %s\n", res.Candidates[0].Content.Parts[0])
	content := res.Candidates[0].Content.Parts[0]

	contentText := fmt.Sprintf("%s", content)

	return contentText, nil
}

func (gc *GeminiClient) GenerateQuizFromPDF(ctx context.Context, pdfPath, personaName, personaRole, personaLanguage, personaDifficulty string) (string, error) {
	prompt := pdfPrompt{
		pdfPath: pdfPath,
	}

	var output strings.Builder

	content, err := gc.extractContentFromPDF(ctx, &output, prompt, modelName)
	if err != nil {
		return "", fmt.Errorf("error generating content from PDF: %w", err)
	}

	fmt.Printf("Generating quiz from PDF content: %s\n", content)
	promptText := fmt.Sprintf("Generate a quiz for a %s (%s) at %s difficulty level based on the following content: %s", personaRole, personaLanguage, personaDifficulty, content)
	quizContent, _, err := gc.generateContent(ctx, instructions.QuizModelSystemInstructions, promptText)
	if err != nil {
		return "", fmt.Errorf("error generating quiz: %w", err)
	}

	return quizContent, nil
}
