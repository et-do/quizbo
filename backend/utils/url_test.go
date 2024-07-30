package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected string
	}{
		{"https://Example.com", "https://example.com"},
		{"http://Example.com/Path", "https://example.com/path"},
		{"HTTP://EXAMPLE.COM/PATH", "https://example.com/path"},
		{"example.com", "https://example.com"},
		{"http://www.example.com", "https://example.com"},
		{"https://www.example.com", "https://example.com"},
		{"www.example.com", "https://example.com"},
		{"http://example.com", "https://example.com"},
		{"https://example.com", "https://example.com"},
		{"HTTP://WWW.EXAMPLE.COM/TEST", "https://example.com/test"},
	}

	for _, test := range tests {
		result, err := NormalizeURL(test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, result)
	}
}
