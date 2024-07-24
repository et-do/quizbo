package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"read-robin/models"
	"read-robin/services"
	"read-robin/services/gemini"

	"golang.org/x/net/context"
)

type ResponseSubmission struct {
	ContentID    string `json:"content_id"`
	QuizID       string `json:"quiz_id"`
	QuestionID   string `json:"question_id"`
	UserResponse string `json:"user_response"`
}

type ReviewResponse struct {
	Status string `json:"status"`
}

func SubmitResponseHandler(w http.ResponseWriter, r *http.Request) {
	var responseSubmission ResponseSubmission
	if err := json.NewDecoder(r.Body).Decode(&responseSubmission); err != nil {
		log.Printf("SubmitResponseHandler: Unable to parse request: %v", err)
		http.Error(w, "Unable to parse request", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	// Create Firestore client
	firestoreClient, err := services.NewFirestoreClient(ctx)
	if err != nil {
		log.Printf("SubmitResponseHandler: Error creating Firestore client: %v", err)
		http.Error(w, "Error creating Firestore client", http.StatusInternalServerError)
		return
	}

	// Fetch the quiz from Firestore
	quiz, err := firestoreClient.GetQuiz(ctx, responseSubmission.ContentID, responseSubmission.QuizID)
	if err != nil {
		log.Printf("SubmitResponseHandler: Error fetching quiz: %v", err)
		http.Error(w, "Error fetching quiz", http.StatusInternalServerError)
		return
	}

	// Find the specific question
	var question models.Question
	found := false
	for _, q := range quiz.Questions {
		if q.QuestionID == responseSubmission.QuestionID {
			question = q
			found = true
			break
		}
	}
	if !found {
		log.Printf("SubmitResponseHandler: Question not found")
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}

	// Create Gemini client
	geminiClient, err := gemini.NewGeminiClient(ctx)
	if err != nil {
		log.Printf("SubmitResponseHandler: Error creating Gemini client: %v", err)
		http.Error(w, "Error creating Gemini client", http.StatusInternalServerError)
		return
	}

	// Prepare data for Gemini
	reviewData := map[string]string{
		"question":        question.Question,
		"user_response":   responseSubmission.UserResponse,
		"expected_answer": question.Answer,
		"reference":       question.Reference,
	}

	reviewDataJSON, err := json.Marshal(reviewData)
	if err != nil {
		log.Printf("SubmitResponseHandler: Error marshaling review data: %v", err)
		http.Error(w, "Error preparing review data", http.StatusInternalServerError)
		return
	}

	// Call Gemini LLM for review
	reviewResult, err := geminiClient.ReviewResponse(ctx, string(reviewDataJSON))
	if err != nil {
		log.Printf("SubmitResponseHandler: Error reviewing response: %v", err)
		http.Error(w, "Error reviewing response", http.StatusInternalServerError)
		return
	}

	// Return the review result to the frontend
	reviewResponse := ReviewResponse{
		Status: reviewResult,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reviewResponse); err != nil {
		log.Printf("SubmitResponseHandler: Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
