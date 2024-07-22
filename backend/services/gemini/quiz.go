package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"read-robin/models"
	"strings"
)

func (gc *GeminiClient) GenerateQuiz(ctx context.Context, summarizedContent, personaName, personaRole, personaLanguage, personaDifficulty string) (string, string, error) {
	promptText := fmt.Sprintf("Generate a quiz for a %s (%s) at %s difficulty level based on the following content: %s", personaRole, personaLanguage, personaDifficulty, summarizedContent)
	return gc.generateContent(ctx, instructions.QuizModelSystemInstructions, promptText)
}

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

func (gc *GeminiClient) ReviewResponse(ctx context.Context, reviewData string) (string, error) {
	reviewResult, _, err := gc.generateContent(ctx, instructions.ReviewModelSystemInstructions, reviewData)
	if err != nil {
		return "", fmt.Errorf("error reviewing response: %w", err)
	}
	return strings.TrimSpace(reviewResult), nil
}
