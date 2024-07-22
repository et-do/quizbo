package gemini

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateContent(t *testing.T) {
	ctx := context.Background()
	client, err := NewGeminiClient(ctx)
	assert.NoError(t, err)

	systemInstructions := "Test system instructions"
	promptText := "Test prompt text"

	content, response, err := client.generateContent(ctx, systemInstructions, promptText)
	assert.NoError(t, err)
	assert.NotEmpty(t, content)
	assert.NotEmpty(t, response)
}

func TestExtractContent(t *testing.T) {
	ctx := context.Background()
	client, err := NewGeminiClient(ctx)
	assert.NoError(t, err)

	htmlText := "<html><body>Test HTML content</body></html>"

	contentMap, response, err := client.ExtractContent(ctx, htmlText)
	assert.NoError(t, err)
	assert.NotNil(t, contentMap)
	assert.NotEmpty(t, response)
}
