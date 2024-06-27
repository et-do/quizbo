package utils

import (
    "path/filepath"
    "runtime"
)

// GetTemplatePath returns the absolute path to a template file.
func GetTemplatePath(filename string) string {
    _, b, _, _ := runtime.Caller(0)
    basepath := filepath.Dir(b)
    return filepath.Join(basepath, "../templates", filename)
}
