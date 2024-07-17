package handlers

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
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

func handlePDFUpload(file multipart.File) ([]byte, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}

func handleMultipartForm(r *http.Request, ctx context.Context, geminiClient *services.GeminiClient, firestoreClient *services.FirestoreClient) (models.SubmitResponse, error) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		return models.SubmitResponse{}, err
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return models.SubmitResponse{}, err
	}
	defer file.Close()

	var persona models.Persona
	if err := json.Unmarshal([]byte(r.FormValue("persona")), &persona); err != nil {
		return models.SubmitResponse{}, err
	}

	return processPDFSubmission(ctx, file, persona, geminiClient, firestoreClient)
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

	quiz, err := services.ParseQuizResponse(quizContentMap, latestQuizID)
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

func processPDFSubmission(ctx context.Context, file multipart.File, persona models.Persona, geminiClient *services.GeminiClient, firestoreClient *services.FirestoreClient) (models.SubmitResponse, error) {
	fileBytes, err := handlePDFUpload(file)
	if err != nil {
		return models.SubmitResponse{}, err
	}

	contentID := services.GenerateID(string(fileBytes))

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

	quizContentMap, title, err := geminiClient.ExtractAndGenerateQuiz(ctx, string(fileBytes), persona)
	if err != nil {
		return models.SubmitResponse{}, err
	}

	latestQuizID := services.GetLatestQuizID(existingQuizzes)

	quiz, err := services.ParseQuizResponse(quizContentMap, latestQuizID)
	if err != nil {
		return models.SubmitResponse{}, err
	}

	err = firestoreClient.SaveQuiz(ctx, "", title, quiz) // No normalized URL for PDF
	if err != nil {
		return models.SubmitResponse{}, err
	}

	return models.SubmitResponse{
		Status:      "success",
		URL:         "", // No URL for PDF
		ContentID:   contentID,
		QuizID:      latestQuizID,
		Title:       title,
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
	contentID := services.GenerateID(normalizedURL)
	return normalizedURL, contentID, nil
}
