package services

import (
	"encoding/json"
	"read-robin/models"
	"testing"
)

const sampleResponse = `
{
	"quiz": [
		{
			"question": "What is the purpose of the 'example' domain?",
			"answer": "The 'example' domain is intended for use in illustrative examples within documents, allowing users to employ it without needing prior permission or coordination.",
			"reference": "This domain is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission."
		},
		{
			"question": "Where can you find more information about the 'example' domain?",
			"answer": "You can find more information about the 'example' domain by following the link provided: [More information...](https://www.iana.org/domains/example)",
			"reference": "[More information...](https://www.iana.org/domains/example)"
		}
	]
}
`

func TestParseQuizResponse(t *testing.T) {
	var sampleResponseMap map[string]interface{}
	err := json.Unmarshal([]byte(sampleResponse), &sampleResponseMap)
	if err != nil {
		t.Fatalf("Failed to unmarshal sample response: %v", err)
	}

	expectedQuiz := models.Quiz{
		QuizID: "0001",
		Questions: []models.Question{
			{
				QuestionID: generateQuestionID(),
				Question:   "What is the purpose of the 'example' domain?",
				Answer:     "The 'example' domain is intended for use in illustrative examples within documents, allowing users to employ it without needing prior permission or coordination.",
				Reference:  "This domain is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission.",
			},
			{
				QuestionID: generateQuestionID(),
				Question:   "Where can you find more information about the 'example' domain?",
				Answer:     "You can find more information about the 'example' domain by following the link provided: [More information...](https://www.iana.org/domains/example)",
				Reference:  "[More information...](https://www.iana.org/domains/example)",
			},
		},
	}

	quiz, err := ParseQuizResponse(sampleResponseMap, "0001")
	if err != nil {
		t.Fatalf("ParseQuizResponse: expected no error, got %v", err)
	}

	if len(quiz.Questions) != len(expectedQuiz.Questions) {
		t.Errorf("ParseQuizResponse: expected %d questions and answers, got %d", len(expectedQuiz.Questions), len(quiz.Questions))
	}

	for i, expectedQA := range expectedQuiz.Questions {
		actualQA := quiz.Questions[i]
		if expectedQA.Question != actualQA.Question {
			t.Errorf("ParseQuizResponse: expected question %q, got %q", expectedQA.Question, actualQA.Question)
		}
		if expectedQA.Answer != actualQA.Answer {
			t.Errorf("ParseQuizResponse: expected answer %q, got %q", expectedQA.Answer, actualQA.Answer)
		}
		if expectedQA.Reference != actualQA.Reference {
			t.Errorf("ParseQuizResponse: expected reference %q, got %q", expectedQA.Reference, actualQA.Reference)
		}
	}
}
