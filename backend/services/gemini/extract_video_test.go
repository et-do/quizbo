package gemini

import (
	"context"
	"fmt"
	"read-robin/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	videoPath = "gs://read-robin-examples/video/happiness_a_very_short_story.mp4"
)

func TestExtractContentFromVideo(t *testing.T) {
	ctx := context.Background()
	client, err := NewGeminiClient(ctx)
	assert.NoError(t, err)

	video_path := videoPrompt{videoPath: videoPath}

	contentMap, fullText, err := client.ExtractContentFromVideo(ctx, video_path.videoPath)
	fmt.Print(contentMap)
	fmt.Print(fullText)
	assert.NoError(t, err)
	assert.NotEmpty(t, contentMap)
}

func TestGenerateQuizFromVideo(t *testing.T) {
	ctx := context.Background()
	client, err := NewGeminiClient(ctx)
	assert.NoError(t, err)

	persona := models.Persona{
		ID:         "Test_ID",
		Name:       "Test Persona",
		Role:       "Test Role",
		Language:   "English",
		Difficulty: "Easy"}

	quizContent, err := client.GenerateQuizFromVideo(ctx, videoPath, persona)
	fmt.Print(quizContent)
	assert.NoError(t, err)
	assert.NotEmpty(t, quizContent)
}
