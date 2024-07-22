package gemini

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractContentFromPDF(t *testing.T) {
	ctx := context.Background()
	client, err := NewGeminiClient(ctx)
	assert.NoError(t, err)

	pdfPath := "gs://read-robin-2e150.appspot.com/pdfs/test_document.pdf"
	prompt := pdfPrompt{pdfPath: pdfPath}
	var output strings.Builder

	content, err := client.extractContentFromPDF(ctx, &output, prompt, modelName)
	assert.NoError(t, err)
	assert.NotEmpty(t, content)
}

func TestGenerateQuizFromPDF(t *testing.T) {
	ctx := context.Background()
	client, err := NewGeminiClient(ctx)
	assert.NoError(t, err)

	pdfPath := "gs://read-robin-2e150.appspot.com/pdfs/test_document.pdf"
	personaName := "Test Persona"
	personaRole := "Test Role"
	personaLanguage := "English"
	personaDifficulty := "Easy"

	quizContent, err := client.GenerateQuizFromPDF(ctx, pdfPath, personaName, personaRole, personaLanguage, personaDifficulty)
	assert.NoError(t, err)
	assert.NotEmpty(t, quizContent)
}
