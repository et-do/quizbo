package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"read-robin/models"
)

const (
	quizModelSystemInstructions = `You are a highly skilled model that generates quiz questions and answers from summarized content tailored for a specific user persona. The persona details include Name, Role (profession, age, etc.), Language, and Difficulty (beginner, intermediate, expert). Your task is to generate questions and answers based on the summarized content provided, considering the persona details. You should also generate a small piece of reference text that was used to create your question/answer pair. Omit any backticks or format reference. Return everything in a JSON dictionary with 'quiz' being an array of objects containing 'question', 'answer', and 'reference' strings. The structure should look like this:
{
	"quiz": [
		{
			"question": "question",
			"answer": "answer",
			"reference": "reference"
		},
		{
			"question": "question",
			"answer": "answer",
			"reference": "reference"
		}
	]
}`
)

// GenerateQuiz generates quiz questions and answers from the summarized content
func (gc *GeminiClient) GenerateQuiz(ctx context.Context, summarizedContent string, persona models.Persona) (string, string, error) {
	promptText := fmt.Sprintf("Generate a quiz for a %s (%s) at %s difficulty level based on the following content: %s", persona.Role, persona.Language, persona.Difficulty, summarizedContent)
	return gc.generateContent(ctx, quizModelSystemInstructions, promptText)
}

// ExtractAndGenerateQuiz extracts content and generates a quiz using the Gemini client
func (gc *GeminiClient) ExtractAndGenerateQuizFromHtml(ctx context.Context, htmlContent string, persona models.Persona) (map[string]interface{}, string, error) {
	contentMap, _, err := gc.ExtractContentFromHtml(ctx, htmlContent)
	if err != nil {
		return nil, "", err
	}

	quizContent, _, err := gc.GenerateQuiz(ctx, contentMap["content"], persona)
	if err != nil {
		return nil, "", err
	}
	var quizContentMap map[string]interface{}
	if err := json.Unmarshal([]byte(quizContent), &quizContentMap); err != nil {
		return nil, "", err
	}

	return quizContentMap, contentMap["title"], nil
}

// ExtractAndGenerateQuiz extracts content and generates a quiz using the Gemini client
func (gc *GeminiClient) ExtractAndGenerateQuizFromPdf(ctx context.Context, pdfPath string, persona models.Persona) (map[string]interface{}, string, error) {
	contentMap, _, err := gc.ExtractContentFromPdf(ctx, pdfPath)
	if err != nil {
		return nil, "", err
	}

	quizContent, _, err := gc.GenerateQuiz(ctx, contentMap["content"], persona)
	if err != nil {
		return nil, "", err
	}
	var quizContentMap map[string]interface{}
	if err := json.Unmarshal([]byte(quizContent), &quizContentMap); err != nil {
		return nil, "", err
	}

	return quizContentMap, contentMap["title"], nil
}

// ExtractAndGenerateQuiz extracts content and generates a quiz using the Gemini client
func (gc *GeminiClient) ExtractAndGenerateQuizFromAudio(ctx context.Context, audioPath string, persona models.Persona) (map[string]interface{}, string, error) {
	contentMap, _, err := gc.ExtractContentFromPdf(ctx, audioPath)
	if err != nil {
		return nil, "", err
	}

	quizContent, _, err := gc.GenerateQuiz(ctx, contentMap["content"], persona)
	if err != nil {
		return nil, "", err
	}
	var quizContentMap map[string]interface{}
	if err := json.Unmarshal([]byte(quizContent), &quizContentMap); err != nil {
		return nil, "", err
	}

	return quizContentMap, contentMap["title"], nil
}
