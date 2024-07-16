package services

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"read-robin/models"

	"cloud.google.com/go/firestore"
)

// TestSaveQuizWithTitle tests the SaveQuizWithTitle function to ensure it correctly saves a quiz and its title to Firestore
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
		Timestamp: time.Now(),
	}

	contentURL := "example.com"
	contentTitle := "Example Domain"
	contentID := GenerateID(contentURL)

	err = firestoreClient.SaveQuiz(ctx, contentURL, contentTitle, quiz)
	if err != nil {
		t.Fatalf("SaveQuizWithTitle: expected no error, got %v", err)
	}

	// Retrieve the content to verify the title and quiz
	doc, err := firestoreClient.Client.Collection("quizzes").Doc(contentID).Get(ctx)
	if err != nil {
		t.Fatalf("Failed to retrieve document: %v", err)
	}

	var content models.Content
	if err := doc.DataTo(&content); err != nil {
		t.Fatalf("Failed to parse content: %v", err)
	}

	if content.Title != contentTitle {
		t.Errorf("SaveQuizWithTitle: expected title %v, got %v", contentTitle, content.Title)
	}

	if len(content.Quizzes) == 0 {
		t.Fatalf("SaveQuizWithTitle: expected quizzes to be saved, got none")
	}

	if content.Quizzes[0].QuizID != quiz.QuizID {
		t.Errorf("SaveQuizWithTitle: expected quizID %v, got %v", quiz.QuizID, content.Quizzes[0].QuizID)
	}
}

// TestGetQuiz tests the GetQuiz function to ensure it correctly retrieves a quiz from Firestore
func TestGetQuiz(t *testing.T) {
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
		Timestamp: time.Now(),
	}

	contentURL := "http://example.com"
	contentTitle := "Example Domain"
	contentID := GenerateID(contentURL)

	// Save the quiz to Firestore first
	err = firestoreClient.SaveQuiz(ctx, contentURL, contentTitle, quiz)
	if err != nil {
		t.Fatalf("SaveQuizWithTitle: expected no error, got %v", err)
	}

	// Retrieve the quiz from Firestore
	retrievedQuiz, err := firestoreClient.GetQuiz(ctx, contentID, quiz.QuizID)
	if err != nil {
		t.Fatalf("GetQuiz: expected no error, got %v", err)
	}

	if retrievedQuiz.QuizID != quiz.QuizID {
		t.Errorf("GetQuiz: expected quizID %v, got %v", quiz.QuizID, retrievedQuiz.QuizID)
	}

	if len(retrievedQuiz.Questions) != len(quiz.Questions) {
		t.Errorf("GetQuiz: expected %d questions, got %d", len(quiz.Questions), len(retrievedQuiz.Questions))
	}

	for i, q := range retrievedQuiz.Questions {
		if q.Question != quiz.Questions[i].Question {
			t.Errorf("GetQuiz: expected question %v, got %v", quiz.Questions[i].Question, q.Question)
		}
		if q.Answer != quiz.Questions[i].Answer {
			t.Errorf("GetQuiz: expected answer %v, got %v", quiz.Questions[i].Answer, q.Answer)
		}
		if q.Reference != quiz.Questions[i].Reference {
			t.Errorf("GetQuiz: expected reference %v, got %v", quiz.Questions[i].Reference, q.Reference)
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

	exitCode := m.Run()

	os.Exit(exitCode)
}
