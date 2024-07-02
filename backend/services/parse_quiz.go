package services

import (
	"fmt"
	"read-robin/models"
)

// ParseQuizResponse parses the response from the Gemini model into a Quiz struct
func ParseQuizResponse(response map[string]interface{}, quizID string) (models.Quiz, error) {
	quizInterface, ok := response["quiz"].([]interface{})
	if !ok {
		return models.Quiz{}, fmt.Errorf("quiz field missing or not an array")
	}

	var questions []models.Question
	for _, qa := range quizInterface {
		qaMap, ok := qa.(map[string]interface{})
		if !ok {
			return models.Quiz{}, fmt.Errorf("error parsing question and answer pair")
		}

		questionText, ok := qaMap["question"].(string)
		if !ok {
			return models.Quiz{}, fmt.Errorf("question field missing or not a string")
		}

		answer, ok := qaMap["answer"].(string)
		if !ok {
			return models.Quiz{}, fmt.Errorf("answer field missing or not a string")
		}

		reference, ok := qaMap["reference"].(string)
		if !ok {
			return models.Quiz{}, fmt.Errorf("reference field missing or not a string")
		}

		questions = append(questions, models.Question{
			QuestionID: generateQuestionID(),
			Question:   questionText,
			Answer:     answer,
			Reference:  reference,
		})
	}

	return models.Quiz{
		QuizID:    quizID,
		Questions: questions,
	}, nil
}
