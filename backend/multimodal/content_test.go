package multimodalpdf

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TODO: Try to convert PDFS to image and upload those or extract the text
func Test_generateContentFromPDF(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if projectID == "" {
		t.Skip("GOLANG_SAMPLES_PROJECT_ID not set")
	}

	buf := new(bytes.Buffer)
	prompt := pdfPrompt{
		pdfPath: "gs://read-robin-testing/test_document.pdf",
		question: `
            You are a very professional document summarization specialist.
            Please summarize the given document.
        `,
	}
	location := "us-central1"
	modelName := "gemini-1.5-flash-001"

	err := generateContentFromPDF(buf, prompt, projectID, location, modelName)
	if err != nil {
		t.Errorf("Test_generateContentFromPDF: %v", err.Error())
	}

	generatedSummary := buf.String()
	generatedSummaryLowercase := strings.ToLower(generatedSummary)
	// We expect these important topics in the video to be correctly covered
	// in the generated summary
	for _, word := range []string{
		"gemini",
		"tokens",
	} {
		if !strings.Contains(generatedSummaryLowercase, word) {
			t.Errorf("expected the word %q in the description of %s", word, prompt.pdfPath)
		}
	}
}
