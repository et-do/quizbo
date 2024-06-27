package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"read-robin/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/submit", handlers.SubmitHandler).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
