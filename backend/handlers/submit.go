package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"read-robin/models"
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

// SubmitHandler handles the form submission and responds with JSON
func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	var urlRequest URLRequest

	if r.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&urlRequest); err != nil {
			log.Printf("SubmitHandler: Unable to parse JSON request: %v", err)
			http.Error(w, "Unable to parse JSON request", http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			log.Printf("SubmitHandler: Unable to parse form: %v", err)
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		urlRequest.URL = r.FormValue("url")
		if urlRequest.URL == "" {
			log.Println("SubmitHandler: Missing URL")
			http.Error(w, "Missing URL", http.StatusBadRequest)
			return
		}
	}

	log.Printf("SubmitHandler: Received URL: %s", urlRequest.URL)

	htmlContent, err := utils.FetchHTML(urlRequest.URL)
	if err != nil {
		log.Printf("SubmitHandler: Error fetching HTML content: %v", err)
		http.Error(w, "Error fetching HTML content", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	geminiClient, err := services.NewGeminiClient(ctx)
	if err != nil {
		log.Printf("SubmitHandler: Error creating Gemini client: %v", err)
		http.Error(w, "Error creating Gemini client", http.StatusInternalServerError)
		return
	}

	extractedContent, _, err := geminiClient.ExtractContent(ctx, htmlContent)
	if err != nil {
		log.Printf("SubmitHandler: Error extracting content: %v", err)
		http.Error(w, "Error extracting content", http.StatusInternalServerError)
		return
	}

	quizContent, _, err := geminiClient.GenerateQuiz(ctx, extractedContent)
	if err != nil {
		log.Printf("SubmitHandler: Error generating quiz: %v", err)
		http.Error(w, "Error generating quiz", http.StatusInternalServerError)
		return
	}

	var quizContentMap map[string]interface{}
	if err := json.Unmarshal([]byte(quizContent), &quizContentMap); err != nil {
		log.Printf("SubmitHandler: Error unmarshalling quiz content: %v", err)
		http.Error(w, "Error unmarshalling quiz content", http.StatusInternalServerError)
		return
	}

	firestoreClient, err := services.NewFirestoreClient(ctx)
	if err != nil {
		log.Printf("SubmitHandler: Error creating Firestore client: %v", err)
		http.Error(w, "Error creating Firestore client", http.StatusInternalServerError)
		return
	}

	normalizedURL, err := utils.NormalizeURL(urlRequest.URL)
	if err != nil {
		log.Printf("SubmitHandler: Error normalizing URL: %v", err)
		http.Error(w, "Error normalizing URL", http.StatusInternalServerError)
		return
	}

	contentID := services.GenerateID(normalizedURL)
	doc, err := firestoreClient.Client.Collection("dev_quizzes").Doc(contentID).Get(ctx)
	var existingQuizzes []models.Quiz
	if err == nil {
		var existingContent models.Content
		if err := doc.DataTo(&existingContent); err != nil {
			log.Printf("SubmitHandler: Error parsing existing content: %v", err)
			http.Error(w, "Error parsing existing content", http.StatusInternalServerError)
			return
		}
		existingQuizzes = existingContent.Quizzes
	}

	nextQuizID := services.GetNextQuizID(existingQuizzes)

	quiz, err := services.ParseQuizResponse(quizContentMap, nextQuizID)
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
		QuizID:    nextQuizID,
	}

	log.Printf("SubmitHandler: Response - %v\n", response)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("SubmitHandler: Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
	log.Println("SubmitHandler: Response sent successfully")
}
