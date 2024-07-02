package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs the incoming HTTP requests and responses.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request
		log.Printf("Incoming request: %s %s from %s", r.Method, r.RequestURI, r.RemoteAddr)

		// Create a response writer to capture the response
		lrw := NewLoggingResponseWriter(w)

		// Call the next handler
		startTime := time.Now()
		next.ServeHTTP(lrw, r)
		duration := time.Since(startTime)

		// Log the response details
		log.Printf("Completed %s %s in %v with status %d", r.Method, r.RequestURI, duration, lrw.statusCode)
	})
}

// LoggingResponseWriter is a custom response writer to capture the status code.
type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewLoggingResponseWriter creates a new logging response writer.
func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

// WriteHeader captures the status code for logging.
func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
