package services

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"read-robin/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func TestSaveQuiz(t *testing.T) {
	ctx := context.Background()

	// Set the environment variable for development
	os.Setenv("ENV", "development")

	// Ensure the environment variable is set for the test
	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		t.Fatal("GCP_PROJECT environment variable not set")
	}

	firestoreClient, err := NewFirestoreClient(ctx)
	if err != nil {
		t.Fatalf("NewFirestoreClient: expected no error, got %v", err)
	}
	defer firestoreClient.client.Close()

	quiz := models.Quiz{
		Question:  "What is the purpose of the 'Example Domain'?",
		Answer:    "The 'Example Domain' is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission.",
		Reference: "This domain is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission.",
	}

	quizzes := models.Quizzes{
		QuizID: generateQuizID(),
		Quiz:   []models.Quiz{quiz},
	}

	content := models.Content{
		URL:       "http://example.com",
		Timestamp: time.Now(),
		Quizzes:   []models.Quizzes{quizzes},
	}

	docID, err := firestoreClient.SaveQuiz(ctx, content.URL, quizzes.Quiz)
	if err != nil {
		t.Fatalf("SaveQuiz: expected no error, got %v", err)
	}

	// Verify the quiz was saved correctly
	collection := "dev_quizzes"
	if os.Getenv("ENV") != "development" {
		collection = "quizzes"
	}

	doc, err := firestoreClient.client.Collection(collection).Doc(docID).Get(ctx)
	if err != nil {
		t.Fatalf("Failed to retrieve quiz: %v", err)
	}

	var savedContent models.Content
	if err := doc.DataTo(&savedContent); err != nil {
		t.Fatalf("DataTo: %v", err)
	}

	// Log the saved quiz fields
	t.Logf("Saved Content URL: %v", savedContent.URL)
	t.Logf("Saved Content Timestamp: %v", savedContent.Timestamp)
	for _, savedQuiz := range savedContent.Quizzes {
		for _, qa := range savedQuiz.Quiz {
			t.Logf("QuizID: %v", qa.QuizID)
			t.Logf("Question: %v", qa.Question)
			t.Logf("Answer: %v", qa.Answer)
			t.Logf("Reference: %v", qa.Reference)
		}
	}

	if savedContent.URL != content.URL {
		t.Errorf("expected URL %v, got %v", content.URL, savedContent.URL)
	}

	if savedContent.Quizzes[0].Quiz[0].Question != quizzes.Quiz[0].Question {
		t.Errorf("expected question %v, got %v", quizzes.Quiz[0].Question, savedContent.Quizzes[0].Quiz[0].Question)
	}

	// Clean up: Delete the document
	_, err = firestoreClient.client.Collection(collection).Doc(docID).Delete(ctx)
	if err != nil {
		t.Fatalf("Failed to delete test document: %v", err)
	}
}

// Additional cleanup for Firestore emulator
func cleanupFirestoreEmulator(t *testing.T, client *firestore.Client, collection string) {
	ctx := context.Background()
	iter := client.Collection(collection).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Fatalf("Failed to iterate documents: %v", err)
		}
		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			t.Fatalf("Failed to delete document: %v", err)
		}
	}
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		fmt.Println("GCP_PROJECT environment variable not set")
		os.Exit(1)
	}

	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Printf("firestore.NewClient: %v\n", err)
		os.Exit(1)
	}
	defer firestoreClient.Close()

	// Clean up collections before and after tests
	cleanupFirestoreEmulator(nil, firestoreClient, "dev_quizzes")
	cleanupFirestoreEmulator(nil, firestoreClient, "quizzes")

	exitCode := m.Run()

	cleanupFirestoreEmulator(nil, firestoreClient, "dev_quizzes")
	cleanupFirestoreEmulator(nil, firestoreClient, "quizzes")

	os.Exit(exitCode)
}
