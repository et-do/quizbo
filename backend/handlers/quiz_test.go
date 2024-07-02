package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"read-robin/services"

	"github.com/gorilla/mux"
)

func setupFirestoreClient(t *testing.T) *services.FirestoreClient {
	ctx := context.Background()
	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		t.Fatal("GCP_PROJECT environment variable not set")
	}
	firestoreClient, err := services.NewFirestoreClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	return firestoreClient
}

func TestGetQuizHandler(t *testing.T) {
	// Known document and quiz IDs
	contentID := "-5a93f065626cbf2c"
	quizID := "0001"

	// Ensure the working directory is the project root
	err := os.Chdir("..")
	if err != nil {
		t.Fatal(err)
	}

	// Set up the Firestore client
	firestoreClient := setupFirestoreClient(t)
	defer firestoreClient.Client.Close()

	// Create a new GET request to the /quiz/{contentID}/{quizID} endpoint with the known contentID and quizID
	getRequest, err := http.NewRequest("GET", "/quiz/"+contentID+"/"+quizID, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new ResponseRecorder to record the response
	responseRecorder := httptest.NewRecorder()
	// Wrap the GetQuizHandler function with http.HandlerFunc and set up the router
	router := mux.NewRouter()
	router.HandleFunc("/quiz/{contentID}/{quizID}", GetQuizHandler)
	router.ServeHTTP(responseRecorder, getRequest)

	// Check if the status code returned by the handler is 200 OK
	if statusCode := responseRecorder.Code; statusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", statusCode, http.StatusOK)
		return
	}

	// Parse the response body into QuizResponse struct
	var quizResponse QuizResponse
	if err := json.NewDecoder(responseRecorder.Body).Decode(&quizResponse); err != nil {
		t.Fatalf("failed to parse quiz response body: %v", err)
	}

	// Log the full response for debugging
	t.Logf("Quiz response body: %+v", quizResponse)

	// Check if the response contains the expected quiz_id and questions
	if quizResponse.QuizID != quizID {
		t.Errorf("handler returned unexpected quiz_id: got %v want %v", quizResponse.QuizID, quizID)
	}
}
