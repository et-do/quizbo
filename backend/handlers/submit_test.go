package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"read-robin/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubmitHandler(t *testing.T) {

	tests := []struct {
		name             string
		requestBody      interface{}
		expectedStatus   int
		expectedResponse models.SubmitResponse
	}{
		{
			name: "Valid URL Submission",
			requestBody: models.URLRequest{
				URL: "http://example.com",
				Persona: models.Persona{
					Name:       "Test Persona",
					Role:       "Test Role",
					Language:   "English",
					Difficulty: "Easy",
				},
			},
			expectedStatus: http.StatusOK,
			expectedResponse: models.SubmitResponse{
				Status: "success",
				URL:    "http://example.com",
			},
		},
		{
			name: "Valid PDF Submission",
			requestBody: models.PDFRequest{
				GCSURI: "gs://test-bucket/test.pdf",
				Persona: models.Persona{
					Name:       "Test Persona",
					Role:       "Test Role",
					Language:   "English",
					Difficulty: "Easy",
				},
			},
			expectedStatus: http.StatusOK,
			expectedResponse: models.SubmitResponse{
				Status: "success",
				URL:    "",
			},
		},
		{
			name:           "Invalid Submission",
			requestBody:    map[string]string{"invalid": "data"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, err := http.NewRequest("POST", "/submit", bytes.NewBuffer(body))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(SubmitHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var response models.SubmitResponse
				err = json.NewDecoder(rr.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse.Status, response.Status)
				assert.Equal(t, tt.expectedResponse.URL, response.URL)
			}
		})
	}
}
