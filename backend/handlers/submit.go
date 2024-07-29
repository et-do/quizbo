package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"read-robin/models"
	"read-robin/services"
	"read-robin/services/gemini"
	"read-robin/utils"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SubmitRequest is a struct to hold the URL and persona details submitted by the user
type SubmitRequest struct {
	URL         string         `json:"url"`
	Persona     models.Persona `json:"persona"`
	ContentType string         `json:"content_type"`
}

// SubmitResponse is a struct to hold the response to be sent back to the user
type SubmitResponse struct {
	Status      string `json:"status"`
	URL         string `json:"url"`
	ContentID   string `json:"content_id"`
	QuizID      string `json:"quiz_id"`
	Title       string `json:"title"`
	ContentText string `json:"content_text"`
	IsFirstQuiz bool   `json:"is_first_quiz"`
}

// decodeSubmitRequest decodes the URL request from the HTTP request
func decodeSubmitRequest(r *http.Request) (SubmitRequest, error) {
	var submitRequest SubmitRequest
	if r.Header.Get("Content-Type") == "application/json" {
		err := utils.DecodeJSONBody(r, &submitRequest)
		return submitRequest, err
	} else {
		err := utils.DecodeFormBody(r, "url", &submitRequest.URL)
		return submitRequest, err
	}
}

// normalizeAndGenerateID normalizes the URL and generates a content ID
func normalizeAndGenerateID(url string) (string, string, error) {
	normalizedURL, err := utils.NormalizeURL(url)
	if err != nil {
		return "", "", err
	}
	contentID := utils.GenerateID(normalizedURL)
	return normalizedURL, contentID, nil
}

// createFirestoreClient creates a new Firestore client
func createFirestoreClient(ctx context.Context) (*services.FirestoreClient, error) {
	return services.NewFirestoreClient(ctx)
}

// createGeminiClient creates a new Gemini client
func createGeminiClient(ctx context.Context) (*gemini.GeminiClient, error) {
	return gemini.NewGeminiClient(ctx)
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	submitRequest, err := decodeSubmitRequest(r)
	if err != nil {
		log.Printf("SubmitHandler: Unable to parse request: %v", err)
		http.Error(w, "Unable to parse request", http.StatusBadRequest)
		return
	}

	log.Printf("SubmitHandler: Received Request: %s", submitRequest)

	ctx := context.Background()

	geminiClient, err := createGeminiClient(ctx)
	if err != nil {
		log.Printf("SubmitHandler: Error creating Gemini client: %v", err)
		http.Error(w, "Error creating Gemini client", http.StatusInternalServerError)
		return
	}

	firestoreClient, err := createFirestoreClient(ctx)
	if err != nil {
		log.Printf("SubmitHandler: Error creating Firestore client: %v", err)
		http.Error(w, "Error creating Firestore client", http.StatusInternalServerError)
		return
	}

	normalizedURL, contentID, err := normalizeAndGenerateID(submitRequest.URL)
	if err != nil {
		log.Printf("SubmitHandler: Error normalizing URL: %v", err)
		http.Error(w, "Error normalizing URL", http.StatusInternalServerError)
		return
	}

	existingQuizzes, err := firestoreClient.GetExistingQuizzes(ctx, contentID)
	isFirstQuiz := false
	if err != nil {
		if status.Code(err) == codes.NotFound {
			existingQuizzes = []models.Quiz{}
			isFirstQuiz = true
		} else {
			log.Printf("SubmitHandler: Error fetching existing quizzes: %v", err)
			http.Error(w, "Error fetching existing quizzes", http.StatusInternalServerError)
			return
		}
	}

	var quizContentMap map[string]interface{}
	var contentMap map[string]string

	switch submitRequest.ContentType {
	case "URL":
		htmlContent, err := utils.FetchHTML(submitRequest.URL)
		if err != nil {
			log.Printf("SubmitHandler: Error fetching HTML content: %v", err)
			http.Error(w, "Error fetching HTML content", http.StatusInternalServerError)
			return
		}

		quizContentMap, contentMap, err = geminiClient.ExtractAndGenerateQuizFromHtml(ctx, htmlContent, submitRequest.Persona)
		if err != nil {
			log.Printf("SubmitHandler: Error generating quiz content: %v", err)
			http.Error(w, "Error generating quiz content", http.StatusInternalServerError)
			return
		}
	case "PDF":
		quizContentMap, contentMap, err = geminiClient.ExtractAndGenerateQuizFromPdf(ctx, submitRequest.URL, submitRequest.Persona)
		if err != nil {
			log.Printf("SubmitHandler: Error generating quiz content from PDF: %v", err)
			http.Error(w, "Error generating quiz content from PDF", http.StatusInternalServerError)
			return
		}
	case "Audio":
		quizContentMap, contentMap, err = geminiClient.ExtractAndGenerateQuizFromAudio(ctx, submitRequest.URL, submitRequest.Persona)
		if err != nil {
			log.Printf("SubmitHandler: Error generating quiz content from Audio: %v", err)
			http.Error(w, "Error generating quiz content from Audio", http.StatusInternalServerError)
			return
		}
	case "Video":
		quizContentMap, contentMap, err = geminiClient.ExtractAndGenerateQuizFromVideo(ctx, submitRequest.URL, submitRequest.Persona)
		if err != nil {
			log.Printf("SubmitHandler: Error generating quiz content from Video: %v", err)
			http.Error(w, "Error generating quiz content from Video", http.StatusInternalServerError)
			return
		}
	default:
		log.Printf("SubmitHandler: Unsupported content type: %v", submitRequest.ContentType)
		http.Error(w, "Unsupported content type", http.StatusBadRequest)
		return
	}

	title := contentMap["title"]
	contentText := contentMap["content"]
	latestQuizID := services.GetLatestQuizID(existingQuizzes)

	quiz, err := utils.ParseQuizResponse(quizContentMap, latestQuizID)
	if err != nil {
		log.Printf("SubmitHandler: Error parsing quiz response: %v", err)
		http.Error(w, "Error parsing quiz response", http.StatusInternalServerError)
		return
	}

	err = firestoreClient.SaveQuiz(ctx, normalizedURL, title, contentText, quiz)
	if err != nil {
		log.Printf("SubmitHandler: Error saving quiz to Firestore: %v", err)
		http.Error(w, "Error saving quiz to Firestore", http.StatusInternalServerError)
		return
	}

	response := SubmitResponse{
		Status:      "success",
		URL:         submitRequest.URL,
		ContentID:   contentID,
		QuizID:      latestQuizID,
		Title:       title,
		ContentText: contentText,
		IsFirstQuiz: isFirstQuiz,
	}

	log.Printf("SubmitHandler: Response - %v\n", response)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("SubmitHandler: Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
	log.Println("SubmitHandler: Response sent successfully")
}
