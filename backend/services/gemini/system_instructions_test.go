package gemini

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadSystemInstructions(t *testing.T) {
	LoadSystemInstructions()
	assert.NotEmpty(t, instructions.QuizModelSystemInstructions)
	assert.NotEmpty(t, instructions.WebscrapeModelSystemInstructions)
	assert.NotEmpty(t, instructions.PDFModelSystemInstructions)
	assert.NotEmpty(t, instructions.ReviewModelSystemInstructions)
}
