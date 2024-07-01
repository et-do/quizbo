package services

import (
	"encoding/json"
	"fmt"
	"read-robin/models"
	"strings"
)

// ParseQuizResponse parses the response from the Gemini model into a Quizzes struct
func ParseQuizResponse(response string) (models.Quizzes, error) {
	// Remove the backticks and leading/trailing whitespace
	cleanedResponse := strings.ReplaceAll(response, "```json", "")
	cleanedResponse = strings.ReplaceAll(cleanedResponse, "```", "")
	cleanedResponse = strings.TrimSpace(cleanedResponse)

	// Parse the cleaned response into a map
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(cleanedResponse), &result); err != nil {
		return models.Quizzes{}, fmt.Errorf("error unmarshalling quiz response: %w", err)
	}

	// Extract the quiz array
	quizInterface, ok := result["quiz"].([]interface{})
	if !ok {
		return models.Quizzes{}, fmt.Errorf("quiz field missing or not an array")
	}

	var quizzes []models.Quiz
	for _, qa := range quizInterface {
		qaMap, ok := qa.(map[string]interface{})
		if !ok {
			return models.Quizzes{}, fmt.Errorf("error parsing question and answer pair")
		}

		question, ok := qaMap["question"].(string)
		if !ok {
			return models.Quizzes{}, fmt.Errorf("question field missing or not a string")
		}

		answer, ok := qaMap["answer"].(string)
		if !ok {
			return models.Quizzes{}, fmt.Errorf("answer field missing or not a string")
		}

		reference, ok := qaMap["reference"].(string)
		if !ok {
			return models.Quizzes{}, fmt.Errorf("reference field missing or not a string")
		}

		quizzes = append(quizzes, models.Quiz{
			QuizID:    generateQuizID(),
			Question:  question,
			Answer:    answer,
			Reference: reference,
		})
	}

	return models.Quizzes{
		QuizID: generateQuizID(),
		Quiz:   quizzes,
	}, nil
}
