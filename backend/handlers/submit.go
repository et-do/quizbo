package handlers

import (
	"encoding/json"
	"io/ioutil"
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
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("SubmitHandler: Error reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		log.Printf("SubmitHandler: Raw request body: %s", string(body))

		var urlRequest models.URLRequest
		var pdfRequest models.PDFRequest

		err = json.Unmarshal(body, &urlRequest)
		if err == nil && urlRequest.URL != "" {
			response, err = processURLSubmission(ctx, urlRequest, geminiClient, firestoreClient)
		} else {
			err = json.Unmarshal(body, &pdfRequest)
			if err == nil && pdfRequest.GCSURI != "" {
				// Use the GCS URI directly
				response, err = processPDFSubmission(ctx, pdfRequest.GCSURI, pdfRequest.Persona, geminiClient, firestoreClient)
			} else {
				log.Printf("SubmitHandler: Error unmarshalling request: %v", err)
				http.Error(w, "Unable to parse request", http.StatusBadRequest)
				return
			}
		}

		if err != nil {
			log.Printf("SubmitHandler: Error processing submission: %v", err)
			http.Error(w, "Error processing submission", http.StatusInternalServerError)
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
