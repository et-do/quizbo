package gemini

import (
	"context"
	"fmt"
	"read-robin/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	audioPath = "gs://read-robin-2e150.appspot.com/audio/Porsche+Macan+July+5+2018+(1).mp3"
)

func TestExtractContentFromAudio(t *testing.T) {
	ctx := context.Background()
	client, err := NewGeminiClient(ctx)
	assert.NoError(t, err)

	audio_path := audioPrompt{audioPath: audioPath}

	contentMap, fullText, err := client.ExtractContentFromAudio(ctx, audio_path.audioPath)
	fmt.Print(contentMap)
	fmt.Print(fullText)
	assert.NoError(t, err)
	assert.NotEmpty(t, contentMap)
}

func TestGenerateQuizFromAudio(t *testing.T) {
	ctx := context.Background()
	client, err := NewGeminiClient(ctx)
	assert.NoError(t, err)

	persona := models.Persona{
		ID:         "Test_ID",
		Name:       "Test Persona",
		Role:       "Test Role",
		Language:   "English",
		Difficulty: "Easy"}

	quizContent, err := client.GenerateQuizFromAudio(ctx, audioPath, persona)
	fmt.Print(quizContent)
	assert.NoError(t, err)
	assert.NotEmpty(t, quizContent)
}
