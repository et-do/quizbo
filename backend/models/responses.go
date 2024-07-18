package models

type URLRequest struct {
	URL     string `json:"url"`
	Persona `json:"persona"`
}

type PDFRequest struct {
	PDFURL  string  `json:"pdf_url"`
	Persona Persona `json:"persona"`
}

type SubmitResponse struct {
	Status      string `json:"status"`
	URL         string `json:"url"`
	ContentID   string `json:"content_id"`
	QuizID      string `json:"quiz_id"`
	Title       string `json:"title"`
	IsFirstQuiz bool   `json:"is_first_quiz"`
}

type QuizResponse struct {
	QuizID    string     `json:"quiz_id"`
	Questions []Question `json:"questions"`
}
