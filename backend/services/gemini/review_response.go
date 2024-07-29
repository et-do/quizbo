package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

const (
	reviewModelSystemInstructions = `You are a friendly tutor reviewing quiz responses. Determine if the user's response captures the main idea of the expected answer and provide a conversational explanation on why it was right or wrong, including where it was found in the text. Offer additional advice or resources for further learning. Use a lenient approach, focusing on main concepts rather than exact wording. Return the response as a JSON object with "status" and "explanation" keys, without any backticks or markdown formatting.

Examples:
1. Expected Answer: "The 'Example Domain' is for use in illustrative examples in documents."
   User Response: "Example Domain is used for examples in documents."
   Response: {"status": "PASS", "explanation": "Great job! Your answer captures the main idea that Example Domain is used for examples in documents. For more, see example domains in technical writing."}
   
2. Expected Answer: "The 'Example Domain' is for use in illustrative examples in documents."
   User Response: "It is used for examples."
   Response: {"status": "PASS", "explanation": "Good effort! You captured the essence that it is used for examples, but remember it is specifically for use in documents. Check MDN Web Docs for more info."}
   
3. Expected Answer: "The 'Example Domain' is for use in illustrative examples in documents."
   User Response: "This domain is used in documents."
   Response: {"status": "FAIL", "explanation": "Almost there! You mentioned documents but missed that it's for illustrative examples. For details, explore RFC 2606."}
   
4. Expected Answer: "The 'Example Domain' is for use in illustrative examples in documents."
   User Response: "It is a domain used in documents."
   Response: {"status": "FAIL", "explanation": "Not quite. Your answer is too vague. It is used for illustrative examples in documents. Review example domains in technical documentation."}`
)

// ReviewResponse reviews the user's response using the Gemini model
func (gc *GeminiClient) ReviewResponse(ctx context.Context, reviewData string) (string, string, error) {
	reviewResult, _, err := gc.generateContent(ctx, reviewModelSystemInstructions, reviewData)
	if err != nil {
		return "", "", fmt.Errorf("error reviewing response: %w", err)
	}

	log.Printf("Raw LLM response: %s", reviewResult)

	// Sanitize the response to ensure it is valid JSON
	reviewResult = strings.TrimSpace(reviewResult)
	reviewResult = strings.TrimPrefix(reviewResult, "`")
	reviewResult = strings.TrimSuffix(reviewResult, "`")

	var result map[string]string
	if err := json.Unmarshal([]byte(reviewResult), &result); err != nil {
		return "", "", fmt.Errorf("error unmarshaling review result: %w", err)
	}

	return result["status"], result["explanation"], nil
}
