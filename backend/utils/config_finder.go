package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// findConfigFile attempts to locate the configuration file in various directories
func FindConfigFile(filename string) (string, error) {
	paths := []string{
		filename,
		filepath.Join("services", filename),
		filepath.Join("..", "services", filename),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("config file %s not found in expected paths", filename)
}
