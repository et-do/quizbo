package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"read-robin/handlers"
	"read-robin/middleware" // Import the middleware package

	gorillahandlers "github.com/gorilla/handlers" // Alias the gorilla/handlers package
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Define your routes
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/submit", handlers.SubmitHandler).Methods("POST")
	r.HandleFunc("/quiz/{contentID}/{quizID}", handlers.GetQuizHandler).Methods("GET")
	r.HandleFunc("/submit-response", handlers.ResponseHandler).Methods("POST")

	// Apply logging middleware
	r.Use(middleware.LoggingMiddleware)

	// Set up CORS
	corsAllowedOrigins := gorillahandlers.AllowedOrigins([]string{
		"http://localhost:3000",
		"http://127.0.0.1:5000",
		"https://read-robin-2e150.web.app",
		"https://read-robin-dev-6yudia4zva-nn.a.run.app",
		"https://read-robin-6yudia4zva-nn.a.run.app",
	})
	corsAllowedMethods := gorillahandlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})
	corsAllowedHeaders := gorillahandlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	// Apply CORS middleware to the router
	corsHandler := gorillahandlers.CORS(corsAllowedOrigins, corsAllowedMethods, corsAllowedHeaders)(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), corsHandler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
