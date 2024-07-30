package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"read-robin/models"
	"read-robin/services"
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

// ExtractAndGenerateQuizFromHtml extracts content and generates a quiz using the Gemini client
func (gc *GeminiClient) ExtractAndGenerateQuizFromHtml(ctx context.Context, htmlContent string, persona models.Persona) (map[string]interface{}, map[string]string, error) {
	contentMap, _, err := gc.ExtractContentFromHtml(ctx, htmlContent)
	if err != nil {
		return nil, nil, err
	}

	quizContent, _, err := gc.GenerateQuiz(ctx, contentMap["content"], persona)
	if err != nil {
		return nil, nil, err
	}
	var quizContentMap map[string]interface{}
	if err := json.Unmarshal([]byte(quizContent), &quizContentMap); err != nil {
		return nil, nil, err
	}

	return quizContentMap, contentMap, nil
}

// ExtractAndGenerateQuizFromPdf extracts content and generates a quiz using the Gemini client
func (gc *GeminiClient) ExtractAndGenerateQuizFromPdf(ctx context.Context, pdfPath string, persona models.Persona) (map[string]interface{}, map[string]string, error) {
	contentMap, _, err := gc.ExtractContentFromPdf(ctx, pdfPath)
	if err != nil {
		return nil, nil, err
	}

	quizContent, _, err := gc.GenerateQuiz(ctx, contentMap["content"], persona)
	if err != nil {
		return nil, nil, err
	}
	var quizContentMap map[string]interface{}
	if err := json.Unmarshal([]byte(quizContent), &quizContentMap); err != nil {
		return nil, nil, err
	}

	return quizContentMap, contentMap, nil
}

// ExtractAndGenerateQuizFromAudio extracts content and generates a quiz using the Gemini client
func (gc *GeminiClient) ExtractAndGenerateQuizFromAudio(ctx context.Context, audioPath string, persona models.Persona) (map[string]interface{}, map[string]string, error) {
	contentMap, _, err := gc.ExtractContentFromAudio(ctx, audioPath)
	if err != nil {
		return nil, nil, err
	}

	quizContent, _, err := gc.GenerateQuiz(ctx, contentMap["content"], persona)
	if err != nil {
		return nil, nil, err
	}
	var quizContentMap map[string]interface{}
	if err := json.Unmarshal([]byte(quizContent), &quizContentMap); err != nil {
		return nil, nil, err
	}

	return quizContentMap, contentMap, nil
}

// ExtractAndGenerateQuizFromVideo extracts content and generates a quiz using the Gemini client
func (gc *GeminiClient) ExtractAndGenerateQuizFromVideo(ctx context.Context, videoPath string, persona models.Persona) (map[string]interface{}, map[string]string, error) {
	contentMap, _, err := gc.ExtractContentFromVideo(ctx, videoPath)
	if err != nil {
		return nil, nil, err
	}

	quizContent, _, err := gc.GenerateQuiz(ctx, contentMap["content"], persona)
	if err != nil {
		return nil, nil, err
	}
	var quizContentMap map[string]interface{}
	if err := json.Unmarshal([]byte(quizContent), &quizContentMap); err != nil {
		return nil, nil, err
	}

	return quizContentMap, contentMap, nil
}

// GenerateQuizFromText generates quiz content directly from text
func (gc *GeminiClient) GenerateQuizFromText(ctx context.Context, contentID string, textContent string, persona models.Persona) (map[string]interface{}, map[string]string, error) {
	// Fetch title from Firestore using contentID
	firestoreClient, err := services.NewFirestoreClient(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating Firestore client: %w", err)
	}

	quizRef := firestoreClient.Client.Collection("quizzes").Doc(contentID)
	quizDoc, err := quizRef.Get(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching quiz from Firestore: %w", err)
	}

	var quizData map[string]interface{}
	if err := quizDoc.DataTo(&quizData); err != nil {
		return nil, nil, fmt.Errorf("error parsing quiz data: %w", err)
	}

	title, ok := quizData["title"].(string)
	if !ok {
		title = "Generated Quiz from Text"
	}

	quizContent, _, err := gc.GenerateQuiz(ctx, textContent, persona)
	if err != nil {
		return nil, nil, err
	}
	var quizContentMap map[string]interface{}
	if err := json.Unmarshal([]byte(quizContent), &quizContentMap); err != nil {
		return nil, nil, err
	}

	contentMap := map[string]string{
		"title":   title,
		"content": textContent,
	}

	return quizContentMap, contentMap, nil
}
