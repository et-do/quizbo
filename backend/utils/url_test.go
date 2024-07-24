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
		{"https://Example.com", "example.com"},
		{"http://Example.com/Path", "example.com/path"},
		{"HTTP://EXAMPLE.COM/PATH", "example.com/path"},
		{"example.com", "example.com"},
		{"http://www.example.com", "example.com"},
		{"https://www.example.com", "example.com"},
		{"www.example.com", "example.com"},
		{"http://example.com", "example.com"},
		{"https://example.com", "example.com"},
		{"HTTP://WWW.EXAMPLE.COM/TEST", "example.com/test"},
	}

	for _, test := range tests {
		result, err := NormalizeURL(test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, result)
	}
}
