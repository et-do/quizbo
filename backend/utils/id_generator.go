package utils

import (
	"fmt"
	"math/rand"
)

// GenerateID creates a unique ID based on the URL
func GenerateID(url string) string {
	return fmt.Sprintf("%x", hash(url))
}

// generateQuestionID generates a 4-digit random question ID
func GenerateQuestionID() string {
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

// hash is a simple hash function for generating URL IDs
func hash(s string) int {
	h := 0
	for _, c := range s {
		h = int(c) + ((h << 5) - h)
	}
	return h
}
