package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://Example.com", "//example.com"},
		{"http://Example.com/Path", "//example.com/path"},
		{"HTTP://EXAMPLE.COM/PATH", "//example.com/path"},
		{"example.com", "example.com"},
	}

	for _, test := range tests {
		result, err := NormalizeURL(test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, result)
	}
}
