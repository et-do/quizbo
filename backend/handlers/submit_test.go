package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"read-robin/models"

	"github.com/gorilla/mux"
)

func TestSubmitHandler(t *testing.T) {
	// Ensure the working directory is the project root
	err := os.Chdir("..")
	if err != nil {
		t.Fatal(err)
	}

	// Create test cases for URL and PDF content types
	testCases := []struct {
		name               string
		contentType        string
		url                string
		expectedStatusCode int
	}{
		{
			name:        "URL content type",
			contentType: "URL",
			url:         "http://www.example.com",
		},
		{
			name:        "PDF content type",
			contentType: "PDF",
			url:         "gs://read-robin-2e150.appspot.com/pdfs/test_document.pdf",
		},
		{
			name:        "Audio content type",
			contentType: "Audio",
			url:         "gs://read-robin-2e150.appspot.com/audio/porsche_macan_ad.mp3",
		},
		{
			name:        "Video content type",
			contentType: "Video",
			url:         "gs://read-robin-2e150.appspot.com/video/happiness_a_very_short_story.mp4",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a SubmitRequest payload with persona details to be sent in the POST request
			submitRequestPayload := SubmitRequest{
				URL: tc.url,
				Persona: models.Persona{
					ID:         "test_persona_id",
					Name:       "Test User",
					Role:       "Student",
					Language:   "Japanese",
					Difficulty: "Intermediate",
				},
				ContentType: tc.contentType,
			}
			// Marshal the payload into JSON format
			submitRequestPayloadBytes, err := json.Marshal(submitRequestPayload)
			if err != nil {
				t.Fatal(err)
			}

			// Create a new POST request to the /submit endpoint with the JSON payload
			postRequest, err := http.NewRequest("POST", "/submit", bytes.NewBuffer(submitRequestPayloadBytes))
			if err != nil {
				t.Fatal(err)
			}
			// Set the Content-Type header to application/json
			postRequest.Header.Set("Content-Type", "application/json")

			// Create a ResponseRecorder to record the response
			responseRecorder := httptest.NewRecorder()
			// Wrap the SubmitHandler function with http.HandlerFunc
			submitHandler := http.HandlerFunc(SubmitHandler)

			// Serve the HTTP request using the handler
			submitHandler.ServeHTTP(responseRecorder, postRequest)

			// Check if the status code returned by the handler is 200 OK
			if statusCode := responseRecorder.Code; statusCode != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", statusCode, http.StatusOK)
				return
			}

			// Parse the response body into SubmitResponse struct
			var submitResponse SubmitResponse
			if err := json.NewDecoder(responseRecorder.Body).Decode(&submitResponse); err != nil {
				t.Fatalf("failed to parse response body: %v", err)
			}

			// Check if the response body contains the expected status, URL, contentID, and quizID
			if submitResponse.Status != "success" || submitResponse.URL != tc.url {
				t.Errorf("handler returned unexpected body: got %v", submitResponse)
			}

			// Log the full response for debugging
			t.Logf("Submit response body: %v", submitResponse)

			// Now make a GET request to the /quiz/{contentID}/{quizID} endpoint using the content_id and quiz_id from the response
			getRequest, err := http.NewRequest("GET", "/quiz/"+submitResponse.ContentID+"/"+submitResponse.QuizID, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a new ResponseRecorder to record the response
			getResponseRecorder := httptest.NewRecorder()

			// Use mux to set up the router and route variables
			router := mux.NewRouter()
			router.HandleFunc("/quiz/{contentID}/{quizID}", GetQuizHandler)

			// Serve the HTTP request using the router
			router.ServeHTTP(getResponseRecorder, getRequest)

			// Check if the status code returned by the handler is 200 OK
			if statusCode := getResponseRecorder.Code; statusCode != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", statusCode, http.StatusOK)
				return
			}

			// Parse the response body into QuizResponse struct
			var quizResponse QuizResponse
			if err := json.NewDecoder(getResponseRecorder.Body).Decode(&quizResponse); err != nil {
				t.Fatalf("failed to parse response body: %v", err)
			}

			// Check if the response body contains questions
			if len(quizResponse.Questions) == 0 {
				t.Errorf("handler returned no questions: got %v", quizResponse)
			}

			// Log the full response for debugging
			t.Logf("Quiz response body: %v", quizResponse)
		})
	}
}
