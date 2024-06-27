package handlers

import (
	"html/template"
	"log"
	"net/http"

	"read-robin/utils"
)

// URLRequest is a struct to hold the URL submitted by the user
type URLRequest struct {
	URL string `json:"url"`
}

// Response is a struct to hold the response to be sent back to the user
type Response struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

// SubmitHandler handles the form submission and responds with a template
func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form data from the request
	// r.ParseForm() parses the form data and populates r.Form and r.PostForm
	if err := r.ParseForm(); err != nil {
		// If parsing the form fails, return a 400 Bad Request error
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve the URL value from the form data
	// r.FormValue("url") returns the first value for the named component of the query
	url := r.FormValue("url")
	if url == "" {
		// If the URL is empty, return a 400 Bad Request error
		http.Error(w, "Missing URL", http.StatusBadRequest)
		return
	}

	// Log the received URL to the server logs
	log.Printf("Received URL: %s\n", url)

	// Create a Response struct with the status and URL
	response := Response{Status: "success", URL: url}

	// Get the path to the template file
	tmplPath := utils.GetTemplatePath("index.html")
	// Parse the template file
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		// If parsing the template fails, log the error and return a 500 Internal Server Error
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Set the content type of the response to HTML
	w.Header().Set("Content-Type", "text/html")
	// Execute the template with the response data
	err = tmpl.Execute(w, response)
	if err != nil {
		// If executing the template fails, log the error and return a 500 Internal Server Error
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
