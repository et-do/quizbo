package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"read-robin/models"

	"golang.org/x/net/context"
)

// SubmitHandler handles both URL and PDF submissions.
func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	geminiClient, err := createGeminiClient(ctx)
	if err != nil {
		log.Printf("SubmitHandler: Error creating Gemini client: %v", err)
		http.Error(w, "Error creating Gemini client", http.StatusInternalServerError)
		return
	}

	firestoreClient, err := createFirestoreClient(ctx)
	if err != nil {
		log.Printf("SubmitHandler: Error creating Firestore client: %v", err)
		http.Error(w, "Error creating Firestore client", http.StatusInternalServerError)
		return
	}

	contentType := r.Header.Get("Content-Type")
	var response models.SubmitResponse

	if contentType == "application/json" {
		urlRequest, err := decodeURLRequest(r)
		if err != nil {
			log.Printf("SubmitHandler: Unable to parse request: %v", err)
			http.Error(w, "Unable to parse request", http.StatusBadRequest)
			return
		}

		response, err = processURLSubmission(ctx, urlRequest, geminiClient, firestoreClient)
		if err != nil {
			log.Printf("SubmitHandler: Error processing URL submission: %v", err)
			http.Error(w, "Error processing URL submission", http.StatusInternalServerError)
			return
		}

	} else if contentType == "multipart/form-data" {
		response, err = handleMultipartForm(r, ctx, geminiClient, firestoreClient)
		if err != nil {
			log.Printf("SubmitHandler: Error processing PDF submission: %v", err)
			http.Error(w, "Error processing PDF submission", http.StatusInternalServerError)
			return
		}

	} else {
		http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
		return
	}

	log.Printf("SubmitHandler: Response - %v\n", response)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("SubmitHandler: Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
	log.Println("SubmitHandler: Response sent successfully")
}
