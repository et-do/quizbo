package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"read-robin/models"
	"read-robin/services"
	"read-robin/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func decodeURLRequest(r *http.Request) (models.URLRequest, error) {
	var urlRequest models.URLRequest
	err := utils.DecodeJSONBody(r, &urlRequest)
	return urlRequest, err
}

func decodePDFRequest(r *http.Request) (models.PDFRequest, error) {
	var pdfRequest models.PDFRequest
	err := utils.DecodeJSONBody(r, &pdfRequest)
	return pdfRequest, err
}

func processURLSubmission(ctx context.Context, urlRequest models.URLRequest, geminiClient *services.GeminiClient, firestoreClient *services.FirestoreClient) (models.SubmitResponse, error) {
	htmlContent, err := utils.FetchHTML(urlRequest.URL)
	if err != nil {
		return models.SubmitResponse{}, err
	}

	normalizedURL, contentID, err := normalizeAndGenerateID(urlRequest.URL)
	if err != nil {
		return models.SubmitResponse{}, err
	}

	existingQuizzes, err := firestoreClient.GetExistingQuizzes(ctx, contentID)
	isFirstQuiz := false
	if err != nil {
		if status.Code(err) == codes.NotFound {
			existingQuizzes = []models.Quiz{}
			isFirstQuiz = true
		} else {
			return models.SubmitResponse{}, err
		}
	}

	quizContentMap, title, err := geminiClient.ExtractAndGenerateQuiz(ctx, htmlContent, urlRequest.Persona)
	if err != nil {
		return models.SubmitResponse{}, err
	}

	latestQuizID := services.GetLatestQuizID(existingQuizzes)

	quiz, err := utils.ParseQuizResponse(quizContentMap, latestQuizID)
	if err != nil {
		return models.SubmitResponse{}, err
	}

	err = firestoreClient.SaveQuiz(ctx, normalizedURL, title, quiz)
	if err != nil {
		return models.SubmitResponse{}, err
	}

	return models.SubmitResponse{
		Status:      "success",
		URL:         urlRequest.URL,
		ContentID:   contentID,
		QuizID:      latestQuizID,
		Title:       title,
		IsFirstQuiz: isFirstQuiz,
	}, nil
}

func processPDFSubmission(ctx context.Context, pdfURL string, persona models.Persona, geminiClient *services.GeminiClient, firestoreClient *services.FirestoreClient) (models.SubmitResponse, error) {
	contentID := utils.GenerateID(pdfURL)

	existingQuizzes, err := firestoreClient.GetExistingQuizzes(ctx, contentID)
	isFirstQuiz := false
	if err != nil {
		if status.Code(err) == codes.NotFound {
			existingQuizzes = []models.Quiz{}
			isFirstQuiz = true
		} else {
			return models.SubmitResponse{}, err
		}
	}

	quizContent, err := geminiClient.GenerateQuizFromPDF(ctx, pdfURL, "Generate a quiz", persona.Name, persona.Role, persona.Language, persona.Difficulty)
	if err != nil {
		return models.SubmitResponse{}, err
	}
	var quizContentMap map[string]interface{}
	if err := json.Unmarshal([]byte(quizContent), &quizContentMap); err != nil {
		return models.SubmitResponse{}, err
	}

	latestQuizID := services.GetLatestQuizID(existingQuizzes)

	quiz, err := utils.ParseQuizResponse(quizContentMap, latestQuizID)
	if err != nil {
		return models.SubmitResponse{}, err
	}

	err = firestoreClient.SaveQuiz(ctx, "", "", quiz) // No normalized URL for PDF
	if err != nil {
		return models.SubmitResponse{}, err
	}

	return models.SubmitResponse{
		Status:      "success",
		URL:         "", // No URL for PDF
		ContentID:   contentID,
		QuizID:      latestQuizID,
		Title:       "", // Title is not extracted for PDF
		IsFirstQuiz: isFirstQuiz,
	}, nil
}

func createFirestoreClient(ctx context.Context) (*services.FirestoreClient, error) {
	return services.NewFirestoreClient(ctx)
}

func createGeminiClient(ctx context.Context) (*services.GeminiClient, error) {
	return services.NewGeminiClient(ctx)
}

func normalizeAndGenerateID(url string) (string, string, error) {
	normalizedURL, err := utils.NormalizeURL(url)
	if err != nil {
		return "", "", err
	}
	contentID := utils.GenerateID(normalizedURL)
	return normalizedURL, contentID, nil
}
