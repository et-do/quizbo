package gemini

import (
	"context"
	"fmt"
	"read-robin/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

const pdfPath = "gs://read-robin-examples/pdfs/chemistry_chapter_page.pdf"

func TestExtractContentFromPDF(t *testing.T) {
	ctx := context.Background()
	client, err := NewGeminiClient(ctx)
	assert.NoError(t, err)

	pdf_path := pdfPrompt{pdfPath: pdfPath}

	contentMap, fullHTML, err := client.ExtractContentFromPdf(ctx, pdf_path.pdfPath)
	fmt.Print(contentMap)
	fmt.Print(fullHTML)
	assert.NoError(t, err)
	assert.NotEmpty(t, contentMap)
}

func TestGenerateQuizFromPDF(t *testing.T) {
	ctx := context.Background()
	client, err := NewGeminiClient(ctx)
	assert.NoError(t, err)

	persona := models.Persona{
		ID:         "Test_ID",
		Name:       "Test Persona",
		Role:       "Test Role",
		Language:   "English",
		Difficulty: "Easy"}

	quizContent, err := client.GenerateQuizFromPDF(ctx, pdfPath, persona)
	fmt.Print(quizContent)
	assert.NoError(t, err)
	assert.NotEmpty(t, quizContent)
}
