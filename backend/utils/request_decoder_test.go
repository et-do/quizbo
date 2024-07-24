package utils

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestRequest struct {
	URL string `json:"url"`
}

func TestDecodeJSONBody(t *testing.T) {
	t.Parallel()
	jsonStr := `{"url": "http://example.com"}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	var testReq TestRequest
	err := DecodeJSONBody(req, &testReq)

	assert.NoError(t, err)
	assert.Equal(t, "http://example.com", testReq.URL)
}

func TestDecodeJSONBody_InvalidJSON(t *testing.T) {
	t.Parallel()
	jsonStr := `{"url": "http://example.com"`
	req := httptest.NewRequest("POST", "/", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	var testReq TestRequest
	err := DecodeJSONBody(req, &testReq)

	assert.Error(t, err)
}

func TestDecodeFormBody(t *testing.T) {
	t.Parallel()
	formStr := "url=http://example.com"
	req := httptest.NewRequest("POST", "/", strings.NewReader(formStr))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var url string
	err := DecodeFormBody(req, "url", &url)

	assert.NoError(t, err)
	assert.Equal(t, "http://example.com", url)
}

func TestDecodeFormBody_InvalidForm(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("POST", "/", strings.NewReader("invalid=formdata"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var url string
	err := DecodeFormBody(req, "url", &url)

	assert.NoError(t, err)
	assert.Equal(t, "", url)
}
