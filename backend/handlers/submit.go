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

// Response is a struct to hold the response to be sent back to the user
type SubmitResponse struct {
	Status string `json:"status"`
	URL    string `json:"url"`
	QuizID string `json:"quiz_id"`
}

// SubmitHandler handles the form submission and responds with JSON
func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	var urlRequest URLRequest

	// Check the Content-Type header to determine how to parse the request body
	if r.Header.Get("Content-Type") == "application/json" {
		// Parse JSON data
		if err := json.NewDecoder(r.Body).Decode(&urlRequest); err != nil {
			log.Printf("SubmitHandler: Unable to parse JSON request: %v", err)
			http.Error(w, "Unable to parse JSON request", http.StatusBadRequest)
			return
		}
	} else {
		// Parse form data
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

	// Log the received URL to the server logs
	log.Printf("SubmitHandler: Received URL: %s", urlRequest.URL)

	// Fetch the HTML content from the URL
	htmlContent, err := utils.FetchHTML(urlRequest.URL)
	if err != nil {
		log.Printf("SubmitHandler: Error fetching HTML content: %v", err)
		http.Error(w, "Error fetching HTML content", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	// Initialize the Gemini client
	geminiClient, err := services.NewGeminiClient(ctx)
	if err != nil {
		log.Printf("SubmitHandler: Error creating Gemini client: %v", err)
		http.Error(w, "Error creating Gemini client", http.StatusInternalServerError)
		return
	}

	// Extract content using the Gemini client
	extractedContent, _, err := geminiClient.ExtractContent(ctx, htmlContent)
	if err != nil {
		log.Printf("SubmitHandler: Error extracting content: %v", err)
		http.Error(w, "Error extracting content", http.StatusInternalServerError)
		return
	}

	// Generate quiz using the Gemini client
	quizContent, _, err := geminiClient.GenerateQuiz(ctx, extractedContent)
	if err != nil {
		log.Printf("SubmitHandler: Error generating quiz: %v", err)
		http.Error(w, "Error generating quiz", http.StatusInternalServerError)
		return
	}

	// Convert the quiz content to a map for parsing
	var quizContentMap map[string]interface{}
	if err := json.Unmarshal([]byte(quizContent), &quizContentMap); err != nil {
		log.Printf("SubmitHandler: Error unmarshalling quiz content: %v", err)
		http.Error(w, "Error unmarshalling quiz content", http.StatusInternalServerError)
		return
	}

	// Parse quiz content
	quiz, err := services.ParseQuizResponse(quizContentMap)
	if err != nil {
		log.Printf("SubmitHandler: Error parsing quiz response: %v", err)
		http.Error(w, "Error parsing quiz response", http.StatusInternalServerError)
		return
	}

	// Initialize the Firestore client
	firestoreClient, err := services.NewFirestoreClient(ctx)
	if err != nil {
		log.Printf("SubmitHandler: Error creating Firestore client: %v", err)
		http.Error(w, "Error creating Firestore client", http.StatusInternalServerError)
		return
	}

	// Save quiz to Firestore
	quizID, err := firestoreClient.SaveQuiz(ctx, urlRequest.URL, []models.Quiz{quiz})
	if err != nil {
		log.Printf("SubmitHandler: Error saving quiz to Firestore: %v", err)
		http.Error(w, "Error saving quiz to Firestore", http.StatusInternalServerError)
		return
	}

	// Include the quiz ID in the response
	response := SubmitResponse{
		Status: "success",
		URL:    urlRequest.URL,
		QuizID: quizID,
	}

	// Set the content type of the response to JSON
	w.Header().Set("Content-Type", "application/json")
	// Write the JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("SubmitHandler: Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
	log.Println("SubmitHandler: Response sent successfully")
}
