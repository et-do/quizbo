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

	os.Setenv("ENV", "development")

	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		t.Fatal("GCP_PROJECT environment variable not set")
	}

	firestoreClient, err := NewFirestoreClient(ctx)
	if err != nil {
		t.Fatalf("NewFirestoreClient: expected no error, got %v", err)
	}
	defer firestoreClient.Client.Close()

	questions := []models.Question{
		{
			QuestionID: generateQuestionID(),
			Question:   "What is the purpose of the 'Example Domain'?",
			Answer:     "The 'Example Domain' is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission.",
			Reference:  "This domain is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission.",
		},
	}

	quiz := models.Quiz{
		QuizID:    "0001",
		Questions: questions,
	}

	content := models.Content{
		URL:       "http://example.com",
		Timestamp: time.Now(),
		Quizzes:   []models.Quiz{quiz},
	}

	docID, err := firestoreClient.SaveQuiz(ctx, content.URL, quiz)
	if err != nil {
		t.Fatalf("SaveQuiz: expected no error, got %v", err)
	}

	collection := "dev_quizzes"
	if os.Getenv("ENV") != "development" {
		collection = "quizzes"
	}

	doc, err := firestoreClient.Client.Collection(collection).Doc(docID).Get(ctx)
	if err != nil {
		t.Fatalf("Failed to retrieve quiz: %v", err)
	}

	var savedContent models.Content
	if err := doc.DataTo(&savedContent); err != nil {
		t.Fatalf("DataTo: %v", err)
	}

	t.Logf("Saved Content URL: %v", savedContent.URL)
	t.Logf("Saved Content Timestamp: %v", savedContent.Timestamp)
	for _, savedQuiz := range savedContent.Quizzes {
		for _, qa := range savedQuiz.Questions {
			t.Logf("QuestionID: %v", qa.QuestionID)
			t.Logf("Question: %v", qa.Question)
			t.Logf("Answer: %v", qa.Answer)
			t.Logf("Reference: %v", qa.Reference)
		}
	}

	if savedContent.URL != content.URL {
		t.Errorf("expected URL %v, got %v", content.URL, savedContent.URL)
	}

	if savedContent.Quizzes[0].Questions[0].Question != quiz.Questions[0].Question {
		t.Errorf("expected question %v, got %v", quiz.Questions[0].Question, savedContent.Quizzes[0].Questions[0].Question)
	}

	_, err = firestoreClient.Client.Collection(collection).Doc(docID).Delete(ctx)
	if err != nil {
		t.Fatalf("Failed to delete test document: %v", err)
	}
}

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

	cleanupFirestoreEmulator(nil, firestoreClient, "dev_quizzes")
	cleanupFirestoreEmulator(nil, firestoreClient, "quizzes")

	exitCode := m.Run()

	cleanupFirestoreEmulator(nil, firestoreClient, "dev_quizzes")
	cleanupFirestoreEmulator(nil, firestoreClient, "quizzes")

	os.Exit(exitCode)
}
