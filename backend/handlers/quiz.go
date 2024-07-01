package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"read-robin/services"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

type QuizResponse struct {
	Questions []string `json:"questions"`
}

// GetQuizHandler retrieves a quiz from Firestore by QuizID
func GetQuizHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	quizID := vars["quizID"]

	if quizID == "" {
		http.Error(w, "quizID is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	firestoreClient, err := services.NewFirestoreClient(ctx)
	if err != nil {
		log.Printf("GetQuizHandler: Error creating Firestore client: %v", err)
		http.Error(w, "Error creating Firestore client", http.StatusInternalServerError)
		return
	}
	defer firestoreClient.Client.Close()

	// Retrieve the quiz from Firestore
	quiz, err := firestoreClient.GetQuiz(ctx, quizID)
	if err != nil {
		log.Printf("GetQuizHandler: Error retrieving quiz from Firestore: %v", err)
		http.Error(w, "Error retrieving quiz from Firestore", http.StatusInternalServerError)
		return
	}

	// Extract questions
	var questions []string
	for _, question := range quiz.Questions {
		questions = append(questions, question.Question)
	}

	// Send response
	response := QuizResponse{
		Questions: questions,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("GetQuizHandler: Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
