package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"read-robin/models"
	"read-robin/services"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

// GetQuizHandler retrieves a quiz from Firestore by contentID and quizID
func GetQuizHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contentID := vars["contentID"]
	quizID := vars["quizID"]

	if contentID == "" || quizID == "" {
		http.Error(w, "contentID and quizID are required", http.StatusBadRequest)
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
	quiz, err := firestoreClient.GetQuiz(ctx, contentID, quizID)
	if err != nil {
		log.Printf("GetQuizHandler: Error retrieving quiz from Firestore: %v", err)
		http.Error(w, "Error retrieving quiz from Firestore", http.StatusInternalServerError)
		return
	}

	// Send response
	response := models.QuizResponse{
		QuizID:    quizID,
		Questions: quiz.Questions,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("GetQuizHandler: Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
