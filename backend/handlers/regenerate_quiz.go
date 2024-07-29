package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"read-robin/models"
	"read-robin/services"
	"read-robin/utils"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RegenerateQuizRequest is a struct to hold the content text and persona details submitted by the user
type RegenerateQuizRequest struct {
	ContentID   string         `json:"content_id"`
	ContentText string         `json:"content_text"`
	Title       string         `json:"title"`
	URL         string         `json:"url"`
	Persona     models.Persona `json:"persona"`
}

// RegenerateQuizHandler handles the regeneration of quizzes from text content
func RegenerateQuizHandler(w http.ResponseWriter, r *http.Request) {
	var request RegenerateQuizRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("RegenerateQuizHandler: Unable to parse request: %v", err)
		http.Error(w, "Unable to parse request", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	geminiClient, err := createGeminiClient(ctx)
	if err != nil {
		log.Printf("RegenerateQuizHandler: Error creating Gemini client: %v", err)
		http.Error(w, "Error creating Gemini client", http.StatusInternalServerError)
		return
	}

	firestoreClient, err := createFirestoreClient(ctx)
	if err != nil {
		log.Printf("RegenerateQuizHandler: Error creating Firestore client: %v", err)
		http.Error(w, "Error creating Firestore client", http.StatusInternalServerError)
		return
	}

	contentID := request.ContentID
	existingQuizzes, err := firestoreClient.GetExistingQuizzes(ctx, contentID)
	if err != nil && status.Code(err) != codes.NotFound {
		log.Printf("RegenerateQuizHandler: Error fetching existing quizzes: %v", err)
		http.Error(w, "Error fetching existing quizzes", http.StatusInternalServerError)
		return
	}

	var quizContentMap map[string]interface{}
	var contentMap map[string]string

	quizContentMap, contentMap, err = geminiClient.GenerateQuizFromText(ctx, contentID, request.ContentText, request.Persona)
	if err != nil {
		log.Printf("RegenerateQuizHandler: Error generating quiz content from text: %v", err)
		http.Error(w, "Error generating quiz content from text", http.StatusInternalServerError)
		return
	}

	title := request.Title
	contentText := contentMap["content"]
	url := request.URL
	latestQuizID := services.GetLatestQuizID(existingQuizzes)

	quiz, err := utils.ParseQuizResponse(quizContentMap, latestQuizID)
	if err != nil {
		log.Printf("RegenerateQuizHandler: Error parsing quiz response: %v", err)
		http.Error(w, "Error parsing quiz response", http.StatusInternalServerError)
		return
	}

	err = firestoreClient.SaveQuiz(ctx, url, title, contentText, quiz)
	if err != nil {
		log.Printf("RegenerateQuizHandler: Error saving quiz to Firestore: %v", err)
		http.Error(w, "Error saving quiz to Firestore", http.StatusInternalServerError)
		return
	}

	response := SubmitResponse{
		Status:      "success",
		ContentID:   contentID,
		QuizID:      latestQuizID,
		Title:       title,
		ContentText: contentText,
		IsFirstQuiz: len(existingQuizzes) == 0,
	}

	log.Printf("RegenerateQuizHandler: Response - %v\n", response)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("RegenerateQuizHandler: Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
	log.Println("RegenerateQuizHandler: Response sent successfully")
}
