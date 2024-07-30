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

func TestRegenerateQuizHandler(t *testing.T) {
	// Ensure the working directory is the project root
	err := os.Chdir("..")
	if err != nil {
		t.Fatal(err)
	}

	// Create test cases for Text content type
	testCases := []struct {
		name               string
		contentID          string
		contentText        string
		title              string
		url                string
		expectedStatusCode int
	}{
		{
			name:        "Text content type",
			contentID:   "-53fd7eb86d4e84fd",
			contentText: "The World's Largest Lobster (French: Le plus grand homard du monde) is a concrete and reinforced steel sculpture in Shediac, New Brunswick, Canada sculpted by Canadian artist Winston Bronnum. Despite being known by its name The World's Largest Lobster, it is not actually the largest lobster sculpture. Description The sculpture is 11 metres long and 5 metres tall, weighing 90 tonnes.[1] The sculpture was commissioned by the Shediac Rotary Club as a tribute to the town's lobster fishing industry.[2] The sculpture took three years to complete,[2] at a cost of $170,000.[3] It attracts 500,000 visitors per year.[2] Contrary to popular belief, this is not actually the \"World's Largest Lobster\" as that title went to the Big Lobster sculpture in Kingston, South Australia, until 2015 when Qianjiang, Hubei, China built a 100-tonne lobster/crayfish.[4] See also * List of world's largest roadside attractions * Betsy the Lobster, another large lobster sculpture",
			title:       "Wikipedia - The World's Largest Lobster",
			url:         "en.wikipedia.org/wiki/the_world's_largest_lobster",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a RegenerateQuizRequest payload with persona details to be sent in the POST request
			regenerateQuizRequestPayload := RegenerateQuizRequest{
				ContentID:   tc.contentID,
				ContentText: tc.contentText,
				Title:       tc.title,
				URL:         tc.url,
				Persona: models.Persona{
					ID:         "test_persona_id",
					Name:       "Test User",
					Role:       "Student",
					Language:   "Japanese",
					Difficulty: "Intermediate",
				},
			}
			// Marshal the payload into JSON format
			regenerateQuizRequestPayloadBytes, err := json.Marshal(regenerateQuizRequestPayload)
			if err != nil {
				t.Fatal(err)
			}

			// Create a new POST request to the /regenerate-quiz endpoint with the JSON payload
			postRequest, err := http.NewRequest("POST", "/regenerate-quiz", bytes.NewBuffer(regenerateQuizRequestPayloadBytes))
			if err != nil {
				t.Fatal(err)
			}
			// Set the Content-Type header to application/json
			postRequest.Header.Set("Content-Type", "application/json")

			// Create a ResponseRecorder to record the response
			responseRecorder := httptest.NewRecorder()
			// Wrap the RegenerateQuizHandler function with http.HandlerFunc
			regenerateQuizHandler := http.HandlerFunc(RegenerateQuizHandler)

			// Serve the HTTP request using the handler
			regenerateQuizHandler.ServeHTTP(responseRecorder, postRequest)

			// Check if the status code returned by the handler is 200 OK
			if statusCode := responseRecorder.Code; statusCode != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", statusCode, http.StatusOK)
				return
			}

			// Parse the response body into SubmitResponse struct
			var regenerateQuizResponse SubmitResponse
			if err := json.NewDecoder(responseRecorder.Body).Decode(&regenerateQuizResponse); err != nil {
				t.Fatalf("failed to parse response body: %v", err)
			}

			// Check if the response body contains the expected status, contentID, and quizID
			if regenerateQuizResponse.Status != "success" || regenerateQuizResponse.ContentID != tc.contentID {
				t.Errorf("handler returned unexpected body: got %v", regenerateQuizResponse)
			}

			// Log the full response for debugging
			t.Logf("Regenerate quiz response body: %v", regenerateQuizResponse)

			// Now make a GET request to the /quiz/{contentID}/{quizID} endpoint using the content_id and quiz_id from the response
			getRequest, err := http.NewRequest("GET", "/quiz/"+regenerateQuizResponse.ContentID+"/"+regenerateQuizResponse.QuizID, nil)
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
