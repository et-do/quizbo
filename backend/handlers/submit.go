package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"read-robin/services"
	"read-robin/utils"

	"golang.org/x/net/context"
)

// URLRequest is a struct to hold the URL submitted by the user
type URLRequest struct {
	URL string `json:"url"`
}

// SubmitResponse is a struct to hold the response to be sent back to the user
type SubmitResponse struct {
	Status    string `json:"status"`
	URL       string `json:"url"`
	ContentID string `json:"content_id"`
	QuizID    string `json:"quiz_id"`
}

// decodeURLRequest decodes the URL request from the HTTP request
func decodeURLRequest(r *http.Request) (URLRequest, error) {
	var urlRequest URLRequest
	if r.Header.Get("Content-Type") == "application/json" {
		err := utils.DecodeJSONBody(r, &urlRequest)
		return urlRequest, err
	} else {
		err := utils.DecodeFormBody(r, "url", &urlRequest.URL)
		return urlRequest, err
	}
}

// normalizeAndGenerateID normalizes the URL and generates a content ID
func normalizeAndGenerateID(url string) (string, string, error) {
	normalizedURL, err := utils.NormalizeURL(url)
	if err != nil {
		return "", "", err
	}
	contentID := services.GenerateID(normalizedURL)
	return normalizedURL, contentID, nil
}

// createFirestoreClient creates a new Firestore client
func createFirestoreClient(ctx context.Context) (*services.FirestoreClient, error) {
	return services.NewFirestoreClient(ctx)
}

// createGeminiClient creates a new Gemini client
func createGeminiClient(ctx context.Context) (*services.GeminiClient, error) {
	return services.NewGeminiClient(ctx)
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	urlRequest, err := decodeURLRequest(r)
	if err != nil {
		log.Printf("SubmitHandler: Unable to parse request: %v", err)
		http.Error(w, "Unable to parse request", http.StatusBadRequest)
		return
	}

	log.Printf("SubmitHandler: Received URL: %s", urlRequest.URL)

	htmlContent, err := utils.FetchHTML(urlRequest.URL)
	if err != nil {
		log.Printf("SubmitHandler: Error fetching HTML content: %v", err)
		http.Error(w, "Error fetching HTML content", http.StatusInternalServerError)
		return
	}

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

	normalizedURL, contentID, err := normalizeAndGenerateID(urlRequest.URL)
	if err != nil {
		log.Printf("SubmitHandler: Error normalizing URL: %v", err)
		http.Error(w, "Error normalizing URL", http.StatusInternalServerError)
		return
	}

	existingQuizzes, err := firestoreClient.GetExistingQuizzes(ctx, contentID) // TODO: Figure out why content IDs with no quizzes aren't generating quizzes
	if err != nil && err.Error() != "firestore: document not found" {
		log.Printf("SubmitHandler: Error fetching existing quizzes: %v", err)
		http.Error(w, "Error fetching existing quizzes", http.StatusInternalServerError)
		return
	}

	quizContentMap, err := geminiClient.ExtractAndGenerateQuiz(ctx, htmlContent)
	if err != nil {
		log.Printf("SubmitHandler: Error generating quiz content: %v", err)
		http.Error(w, "Error generating quiz content", http.StatusInternalServerError)
		return
	}

	latestQuizID := services.GetLatestQuizID(existingQuizzes)

	quiz, err := services.ParseQuizResponse(quizContentMap, latestQuizID)
	if err != nil {
		log.Printf("SubmitHandler: Error parsing quiz response: %v", err)
		http.Error(w, "Error parsing quiz response", http.StatusInternalServerError)
		return
	}

	_, err = firestoreClient.SaveQuiz(ctx, normalizedURL, quiz)
	if err != nil {
		log.Printf("SubmitHandler: Error saving quiz to Firestore: %v", err)
		http.Error(w, "Error saving quiz to Firestore", http.StatusInternalServerError)
		return
	}

	response := SubmitResponse{
		Status:    "success",
		URL:       urlRequest.URL,
		ContentID: contentID,
		QuizID:    latestQuizID,
	}

	log.Printf("SubmitHandler: Response - %v\n", response)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("SubmitHandler: Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
	log.Println("SubmitHandler: Response sent successfully")
}
