package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"read-robin/utils" // Import the utils package
)

// URLRequest is a struct to hold the URL submitted by the user
type URLRequest struct {
	URL string `json:"url"`
}

// Response is a struct to hold the response to be sent back to the user
type Response struct {
	Status string `json:"status"`
	URL    string `json:"url"`
	HTML   string `json:"html"`
}

// SubmitHandler handles the form submission and responds with JSON
func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	var urlRequest URLRequest

	// Check the Content-Type header to determine how to parse the request body
	if r.Header.Get("Content-Type") == "application/json" {
		// Parse JSON data
		if err := json.NewDecoder(r.Body).Decode(&urlRequest); err != nil {
			log.Printf("SubmitHandler: Unable to parse JSON request: %v", err)
			http.Error(w, "Unable to parse JSON request", http.StatusBadRequest)
			return
		}
	} else {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			log.Printf("SubmitHandler: Unable to parse form: %v", err)
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		urlRequest.URL = r.FormValue("url")
		if urlRequest.URL == "" {
			log.Println("SubmitHandler: Missing URL")
			http.Error(w, "Missing URL", http.StatusBadRequest)
			return
		}
	}

	// Log the received URL to the server logs
	log.Printf("SubmitHandler: Received URL: %s", urlRequest.URL)

	// Fetch the HTML content from the URL
	htmlContent, err := utils.FetchHTML(urlRequest.URL)
	if err != nil {
		log.Printf("SubmitHandler: Error fetching HTML content: %v", err)
		http.Error(w, "Error fetching HTML content", http.StatusInternalServerError)
		return
	}

	// Create a Response struct with the status, URL, and HTML content
	response := Response{Status: "success", URL: urlRequest.URL, HTML: htmlContent}
	// Log the received URL to the server logs
	log.Println("Constructed Response: ", response)
	// Set the content type of the response to JSON
	w.Header().Set("Content-Type", "application/json")
	// Write the JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("SubmitHandler: Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
	log.Println("SubmitHandler: Response sent successfully")
}
