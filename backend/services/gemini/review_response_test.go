package gemini

import (
	"context"
	"encoding/json"
	"os"
	"testing"
)

func TestReviewResponse(t *testing.T) {
	t.Parallel()
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
		expectedStatus string
	}{
		{
			name: "Passing case",
			reviewData: map[string]string{
				"question":        "What is the purpose of the 'Example Domain'?",
				"user_response":   "The 'Example Domain' is used in illustrative examples.",
				"expected_answer": "The 'Example Domain' is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission.",
				"reference":       "This domain is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission.",
			},
			expectedStatus: "PASS",
		},
		{
			name: "Failing case",
			reviewData: map[string]string{
				"question":        "What is the purpose of the 'Example Domain'?",
				"user_response":   "It is a domain for testing purposes.",
				"expected_answer": "The 'Example Domain' is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission.",
				"reference":       "This domain is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission.",
			},
			expectedStatus: "FAIL",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reviewDataJSON, err := json.Marshal(tc.reviewData)
			if err != nil {
				t.Fatalf("Error marshaling review data: %v", err)
			}

			status, explanation, err := geminiClient.ReviewResponse(ctx, string(reviewDataJSON))
			if err != nil {
				t.Fatalf("ReviewResponse: expected no error, got %v", err)
			}

			if status != tc.expectedStatus {
				t.Errorf("ReviewResponse: expected status %s, got %s", tc.expectedStatus, status)
			} else {
				t.Logf("Review Status for %s: %s", tc.name, status)
			}

			if explanation == "" {
				t.Errorf("ReviewResponse: expected non-empty explanation for %s", tc.name)
			} else {
				t.Logf("Review Explanation for %s: %s", tc.name, explanation)
			}
		})
	}
}
