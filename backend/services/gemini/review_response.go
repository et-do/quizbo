package gemini

import (
	"context"
	"fmt"
	"strings"
)

const (
	reviewModelSystemInstructions = `You are a highly skilled model that reviews quiz responses. Your task is to determine if the user's response captures the essence of the expected answer based on the reference provided. As long as the user's response includes the key points or main ideas of the expected answer, it should be considered a "PASS". Yes or No is an acceptable response for yes and no questions. Use a lenient approach, focusing on the main concepts rather than exact wording. Return only "PASS" or "FAIL" as the response.

	Examples:
	1. Expected Answer: "The 'Example Domain' is for use in illustrative examples in documents."
	   User Response: "Example Domain is used in documents for examples."
	   Response: "PASS"
	
	2. Expected Answer: "The 'Example Domain' is for use in illustrative examples in documents."
	   User Response: "It is a domain for examples."
	   Response: "PASS"
	
	3. Expected Answer: "The 'Example Domain' is for use in illustrative examples in documents."
	   User Response: "This domain is for examples in documents."
	   Response: "PASS"
	
	4. Expected Answer: "The 'Example Domain' is for use in illustrative examples in documents."
	   User Response: "It is a domain used in documents."
	   Response: "FAIL"`
)

// ReviewResponse reviews the user's response using the Gemini model
func (gc *GeminiClient) ReviewResponse(ctx context.Context, reviewData string) (string, error) {
	reviewResult, _, err := gc.generateContent(ctx, reviewModelSystemInstructions, reviewData)
	if err != nil {
		return "", fmt.Errorf("error reviewing response: %w", err)
	}
	return strings.TrimSpace(reviewResult), nil
}
