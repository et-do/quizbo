package handlers

import (
	"html/template"
	"log"
	"net/http"

	"read-robin/utils"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := utils.GetTemplatePath("index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
