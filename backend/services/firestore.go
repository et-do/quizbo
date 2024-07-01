package services

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"read-robin/models"

	"cloud.google.com/go/firestore"
)

// FirestoreClient is a wrapper around the Firestore client
type FirestoreClient struct {
	client *firestore.Client
}

// NewFirestoreClient creates a new Firestore client
func NewFirestoreClient(ctx context.Context) (*FirestoreClient, error) {
	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		return nil, fmt.Errorf("GCP_PROJECT environment variable not set")
	}

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %v", err)
	}

	return &FirestoreClient{client: client}, nil
}

// SaveQuiz saves a quiz to Firestore
func (fc *FirestoreClient) SaveQuiz(ctx context.Context, url string, quizzes []models.Quiz) (string, error) {
	collection := "quizzes"
	if os.Getenv("ENV") == "development" {
		collection = "dev_quizzes"
	}

	contentID := generateID(url)

	// Assign a QuizID and QuestionID to each Quiz and Question
	for i := range quizzes {
		quizzes[i].QuizID = generateQuizID()
		for j := range quizzes[i].Questions {
			quizzes[i].Questions[j].QuestionID = generateQuestionID()
		}
	}

	content := models.Content{
		URL:       url,
		Timestamp: time.Now(),
		ContentID: contentID,
		Quizzes:   quizzes,
	}

	_, err := fc.client.Collection(collection).Doc(contentID).Set(ctx, content)
	if err != nil {
		return "", fmt.Errorf("failed adding quiz: %v", err)
	}
	return contentID, nil
}

// generateID creates a unique ID based on the URL
func generateID(url string) string {
	// Implement your own logic to generate a unique ID based on the URL
	// For simplicity, let's use a hash or similar technique
	return fmt.Sprintf("%x", hash(url))
}

// generateQuizID generates a 4-digit random quiz ID
func generateQuizID() string {
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

// generateQuestionID generates a 4-digit random question ID
func generateQuestionID() string {
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

// hash is a simple hash function for generating URL IDs
func hash(s string) int {
	h := 0
	for _, c := range s {
		h = int(c) + ((h << 5) - h)
	}
	return h
}
