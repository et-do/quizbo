package services

import (
	"context"
	"encoding/json"
	"os"
	"read-robin/models"
	"testing"
)

const testHTML string = `<!doctype html>
<html>
<head>
    <title>Example Domain</title>

    <meta charset="utf-8" />
    <meta http-equiv="Content-type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style type="text/css">
    body {
        background-color: #f0f0f2;
        margin: 0;
        padding: 0;
        font-family: -apple-system, system-ui, BlinkMacSystemFont, "Segoe UI", "Open Sans", "Helvetica Neue", Helvetica, Arial, sans-serif;
    }
    div {
        width: 600px;
        margin: 5em auto;
        padding: 2em;
        background-color: #fdfdff;
        border-radius: 0.5em;
        box-shadow: 2px 3px 7px 2px rgba(0,0,0,0.02);
    }
    a:link, a:visited {
        color: #38488f;
        text-decoration: none;
    }
    @media (max-width: 700px) {
        div {
            margin: 0 auto;
            width: auto;
        }
    }
    </style>
</head>

<body>
<div>
    <h1>Example Domain</h1>
    <p>This domain is for use in illustrative examples in documents. You may use this
    domain in literature without prior coordination or asking for permission.</p>
    <p><a href="https://www.iana.org/domains/example">More information...</a></p>
</div>
</body>
</html>`

func TestExtractContent(t *testing.T) {
	ctx := context.Background()

	// Ensure the environment variable is set for the test
	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		t.Fatal("GCP_PROJECT environment variable not set")
	}

	geminiClient, err := NewGeminiClient(ctx)
	if err != nil {
		t.Fatalf("NewGeminiClient: expected no error, got %v", err)
	}

	contentMap, fullHTML, err := geminiClient.ExtractContent(ctx, testHTML)
	if err != nil {
		t.Fatalf("ExtractContent: expected no error, got %v", err)
	}

	if contentMap["content"] == "" {
		t.Errorf("ExtractContent: expected extracted content, got an empty string")
	} else {
		t.Logf("Extracted Content: %s", contentMap["content"])
	}

	if contentMap["title"] == "" {
		t.Errorf("ExtractContent: expected generated title, got an empty string")
	} else {
		t.Logf("Generated Title: %s", contentMap["title"])
	}

	t.Logf("Full HTML Response: %s", fullHTML)
}

func TestGenerateQuiz(t *testing.T) {
	ctx := context.Background()

	// Ensure the environment variable is set for the test
	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		t.Fatal("GCP_PROJECT environment variable not set")
	}

	geminiClient, err := NewGeminiClient(ctx)
	if err != nil {
		t.Fatalf("NewGeminiClient: expected no error, got %v", err)
	}

	// Use the previously tested method to get the extracted content
	contentMap, _, err := geminiClient.ExtractContent(ctx, testHTML)
	if err != nil {
		t.Fatalf("ExtractContent: expected no error, got %v", err)
	}

	// Define a test persona
	testPersona := models.Persona{
		ID:         "test_persona_id",
		Name:       "Test User",
		Role:       "Student",
		Language:   "English",
		Difficulty: "Intermediate",
	}

	quiz, fullQuiz, err := geminiClient.GenerateQuiz(ctx, contentMap["content"], testPersona.Name, testPersona.Role, testPersona.Language, testPersona.Difficulty)
	if err != nil {
		t.Fatalf("GenerateQuiz: expected no error, got %v", err)
	}

	if quiz == "" {
		t.Errorf("GenerateQuiz: expected a quiz, got an empty string")
	} else {
		t.Logf("Quiz: %s", quiz)
		t.Logf("Full Quiz Response: %s", fullQuiz)
	}
}

func TestReviewResponse(t *testing.T) {
	ctx := context.Background()

	// Ensure the environment variable is set for the test
	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		t.Fatal("GCP_PROJECT environment variable not set")
	}

	geminiClient, err := NewGeminiClient(ctx)
	if err != nil {
		t.Fatalf("NewGeminiClient: expected no error, got %v", err)
	}

	testCases := []struct {
		name           string
		reviewData     map[string]string
		expectedResult string
	}{
		{
			name: "Passing case",
			reviewData: map[string]string{
				"question":        "What is the purpose of the 'Example Domain'?",
				"user_response":   "The 'Example Domain' is used in illustrative examples.",
				"expected_answer": "The 'Example Domain' is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission.",
				"reference":       "This domain is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission.",
			},
			expectedResult: "PASS",
		},
		{
			name: "Failing case",
			reviewData: map[string]string{
				"question":        "What is the purpose of the 'Example Domain'?",
				"user_response":   "It is a domain for testing purposes.",
				"expected_answer": "The 'Example Domain' is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission.",
				"reference":       "This domain is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission.",
			},
			expectedResult: "FAIL",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reviewDataJSON, err := json.Marshal(tc.reviewData)
			if err != nil {
				t.Fatalf("Error marshaling review data: %v", err)
			}

			reviewResult, err := geminiClient.ReviewResponse(ctx, string(reviewDataJSON))
			if err != nil {
				t.Fatalf("ReviewResponse: expected no error, got %v", err)
			}

			if reviewResult != tc.expectedResult {
				t.Errorf("ReviewResponse: expected %s, got %s", tc.expectedResult, reviewResult)
			} else {
				t.Logf("Review Result for %s: %s", tc.name, reviewResult)
			}
		})
	}
}
