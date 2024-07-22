package gemini

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGeminiClient(t *testing.T) {
	os.Setenv("GCP_PROJECT", "test-project")
	ctx := context.Background()
	client, err := NewGeminiClient(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
