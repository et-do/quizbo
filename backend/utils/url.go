package utils

import (
	"net/url"
	"strings"
)

// NormalizeURL normalizes the URL by making it lowercase and removing the scheme and www.
func NormalizeURL(rawURL string) (string, error) {
	// Ensure the URL has a scheme for proper parsing
	if !strings.HasPrefix(strings.ToLower(rawURL), "http://") && !strings.HasPrefix(strings.ToLower(rawURL), "https://") {
		rawURL = "http://" + rawURL
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
	normalizedURL := host + path

	return normalizedURL, nil
}
