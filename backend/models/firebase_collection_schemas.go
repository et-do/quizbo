package models

import "time"

// Question represents a single question and answer pair with a reference
type Question struct {
	QuestionID string `json:"question_id" firestore:"question_id"`
	Question   string `json:"question" firestore:"question"`
	Answer     string `json:"answer" firestore:"answer"`
	Reference  string `json:"reference" firestore:"reference"`
}

// Quiz represents the structure of a quiz with a list of questions and a timestamp
type Quiz struct {
	QuizID    string     `json:"quiz_id" firestore:"quiz_id"`
	Questions []Question `json:"questions" firestore:"questions"`
	Timestamp time.Time  `json:"timestamp" firestore:"timestamp"`
}

// Content represents the structure of content with multiple quizzes
type Content struct {
	Timestamp time.Time `json:"timestamp" firestore:"timestamp"`
	ContentID string    `json:"content_id" firestore:"content_id"`
	URL       string    `json:"url" firestore:"url"`
	Title     string    `json:"title" firestore:"title"` // Add this line
	Quizzes   []Quiz    `json:"quizzes" firestore:"quizzes"`
}

type Persona struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	Language   string `json:"language"`
	Difficulty string `json:"difficulty"`
}
