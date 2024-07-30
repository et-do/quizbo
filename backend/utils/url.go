package utils

import (
	"net/url"
	"strings"
)

// NormalizeURL normalizes the URL by ensuring it has https://, making it lowercase, and removing the www.
func NormalizeURL(rawURL string) (string, error) {
	// Ensure the URL has a scheme for proper parsing
	if !strings.HasPrefix(strings.ToLower(rawURL), "http://") && !strings.HasPrefix(strings.ToLower(rawURL), "https://") {
		rawURL = "https://" + rawURL
	} else if strings.HasPrefix(strings.ToLower(rawURL), "http://") {
		rawURL = "https://" + rawURL[len("http://"):]
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	host := strings.ToLower(u.Hostname())
	path := strings.ToLower(u.Path)

	// Unconditionally remove "www." from the host
	host = strings.TrimPrefix(host, "www.")

	// Combine the host and path
	normalizedURL := "https://" + host + path

	return normalizedURL, nil
}
