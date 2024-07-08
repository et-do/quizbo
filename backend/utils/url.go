package utils

import (
	"net/url"
	"strings"
)

// NormalizeURL standardizes the URL by converting it to lowercase and removing the scheme (http/https)
func NormalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	u.Scheme = ""
	normalizedURL := strings.ToLower(u.String())
	return normalizedURL, nil
}
