package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchHTML(t *testing.T) {
	t.Parallel()
	// Create a test server that returns some HTML
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<html><body><h1>Test Page</h1></body></html>"))
	}))
	defer ts.Close()

	html, err := FetchHTML(ts.URL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedHTML := "<html><body><h1>Test Page</h1></body></html>"
	if html != expectedHTML {
		t.Errorf("expected %s, got %s", expectedHTML, html)
	}
}
