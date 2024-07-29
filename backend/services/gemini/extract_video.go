package gemini

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"path/filepath"
	"read-robin/models"

	"cloud.google.com/go/vertexai/genai"
)

const (
	videoModelSystemInstructions = `You are a highly skilled model that generates a full text transcript from Video content and generates a title for the content. Your task is to extract the given Video content and output it into a clear and concise article, ignoring any unnecessary formatting or irrelevant content. Additionally, generate a title that objectively defines the main topic of the Video. Return everything in a JSON dictionary with 'content' and 'title' keys, omit any markdown backticks. The structure should look like this:
    {
        "content": "extracted content",
        "title": "generated title"
    }`
)

type videoPrompt struct {
	videoPath string
}

// ExtractContentFromVideo extracts readable text and title from Video content using the Gemini model
func (gc *GeminiClient) ExtractContentFromVideo(ctx context.Context, videoPath string) (map[string]string, string, error) {
	model := gc.client.GenerativeModel(modelName)

	part := genai.FileData{
		MIMEType: mime.TypeByExtension(filepath.Ext(videoPath)),
		FileURI:  videoPath,
	}

	fmt.Printf("Extracting content from Video: %s\n", videoPath)
	res, err := model.GenerateContent(ctx, genai.Text(videoModelSystemInstructions), part)
	if err != nil {
		return nil, "", fmt.Errorf("unable to generate contents: %w", err)
	}

	if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
		return nil, "", errors.New("empty response from model")
	}

	content := res.Candidates[0].Content.Parts[0]

	contentText := fmt.Sprintf("%s", content)

	// Parse the JSON response to extract content and title
	var contentMap map[string]string
	if err := json.Unmarshal([]byte(contentText), &contentMap); err != nil {
		return nil, "", fmt.Errorf("json.Unmarshal: %w", err)
	}

	// Convert the response to a readable format
	fullResponse, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return nil, "", fmt.Errorf("json.MarshalIndent: %w", err)
	}

	return contentMap, string(fullResponse), nil
}

func (gc *GeminiClient) GenerateQuizFromVideo(ctx context.Context, videoPath string, persona models.Persona) (string, error) {
	prompt := videoPrompt{
		videoPath: videoPath,
	}

	contentMap, fullHTML, err := gc.ExtractContentFromVideo(ctx, prompt.videoPath)
	if err != nil {
		return "", fmt.Errorf("error generating content from Video: %w", err)
	}

	fmt.Printf("Generating quiz from Video content: %s\n", contentMap)
	fmt.Print(fullHTML)
	promptText := fmt.Sprintf("Generate a quiz for a %s (%s) at %s difficulty level based on the following content: %s", persona.Role, persona.Language, persona.Difficulty, contentMap)
	quizContent, _, err := gc.generateContent(ctx, quizModelSystemInstructions, promptText)
	if err != nil {
		return "", fmt.Errorf("error generating quiz: %w", err)
	}

	return quizContent, nil
}
