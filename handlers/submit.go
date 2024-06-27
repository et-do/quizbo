package handlers

import (
    "encoding/json"
    "html/template"
    "log"
    "net/http"

    "read-robin/utils"
)

type URLRequest struct {
    URL string `json:"url"`
}

type Response struct {
    Status string `json:"status"`
    URL    string `json:"url"`
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
    var urlRequest URLRequest
    if err := json.NewDecoder(r.Body).Decode(&urlRequest); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    log.Printf("Received URL: %s\n", urlRequest.URL)

    response := Response{Status: "success", URL: urlRequest.URL}

    tmplPath := utils.GetTemplatePath("index.html")
    tmpl, err := template.ParseFiles(tmplPath)
    if err != nil {
        log.Printf("Error parsing template: %v", err)
        http.Error(w, "Error loading template", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "text/html")
    err = tmpl.Execute(w, response)
    if err != nil {
        log.Printf("Error executing template: %v", err)
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
    }
}
