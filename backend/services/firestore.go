package services

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"read-robin/models"

	"cloud.google.com/go/firestore"
)

// FirestoreClient is a wrapper around the Firestore client
type FirestoreClient struct {
	Client *firestore.Client
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

	return &FirestoreClient{Client: client}, nil
}

// SaveQuiz saves a quiz to Firestore
func (fc *FirestoreClient) SaveQuiz(ctx context.Context, url string, quiz models.Quiz) (string, error) {
	collection := "quizzes"
	if os.Getenv("ENV") == "development" {
		collection = "dev_quizzes"
	}

	contentID := GenerateID(url)
	docRef := fc.Client.Collection(collection).Doc(contentID)

	// Get the existing document or create a new one
	doc, err := docRef.Get(ctx)
	var content models.Content
	if err == nil {
		if err := doc.DataTo(&content); err != nil {
			return "", fmt.Errorf("failed to parse existing content: %v", err)
		}
	} else {
		content = models.Content{
			URL:       url,
			Timestamp: time.Now(),
			ContentID: contentID,
			Quizzes:   []models.Quiz{},
		}
	}

	nextQuizID := GetNextQuizID(content.Quizzes)
	quiz.QuizID = nextQuizID

	content.Quizzes = append(content.Quizzes, quiz)

	_, err = docRef.Set(ctx, content)
	if err != nil {
		return "", fmt.Errorf("failed adding quiz: %v", err)
	}
	return contentID, nil
}

// GetQuiz retrieves a quiz from Firestore by contentID and quizID
func (fc *FirestoreClient) GetQuiz(ctx context.Context, contentID, quizID string) (*models.Quiz, error) {
	collection := "quizzes"
	if os.Getenv("ENV") == "development" {
		collection = "dev_quizzes"
	}

	doc, err := fc.Client.Collection(collection).Doc(contentID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving quiz: %v", err)
	}

	var content models.Content
	if err := doc.DataTo(&content); err != nil {
		return nil, fmt.Errorf("dataTo: %v", err)
	}

	for _, quiz := range content.Quizzes {
		if quiz.QuizID == quizID {
			return &quiz, nil
		}
	}

	return nil, fmt.Errorf("no quiz found for quizID: %s", quizID)
}

// GetNextQuizID generates the next sequential quiz ID for the given quizzes
func GetNextQuizID(quizzes []models.Quiz) string {
	maxID := 0
	for _, quiz := range quizzes {
		id, err := strconv.Atoi(quiz.QuizID)
		if err == nil && id > maxID {
			maxID = id
		}
	}
	return fmt.Sprintf("%04d", maxID+1)
}

// GenerateID creates a unique ID based on the URL
func GenerateID(url string) string {
	return fmt.Sprintf("%x", hash(url))
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
