package models

import "time"

// Quiz represents a single question and answer pair with a reference
type Quiz struct {
	QuizID    string `json:"quiz_id" firestore:"quiz_id"`
	Question  string `json:"question" firestore:"question"`
	Answer    string `json:"answer" firestore:"answer"`
	Reference string `json:"reference" firestore:"reference"`
}

// Quizzes represents the structure of a quiz collection
type Quizzes struct {
	QuizID string `json:"quiz_id" firestore:"quiz_id"`
	Quiz   []Quiz `json:"quiz" firestore:"questions_and_answers"`
}

// Content represents the structure of content with multiple quizzes
type Content struct {
	Timestamp time.Time `firestore:"timestamp"`
	ContentID string    `firestore:"content_id"`
	URL       string    `firestore:"url"`
	Quizzes   []Quizzes `firestore:"quizzes"`
}
