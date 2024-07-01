package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// HomeHandler serves a simple JSON response
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Welcome to the API"}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
	log.Println("HomeHandler: Response sent successfully")
}
