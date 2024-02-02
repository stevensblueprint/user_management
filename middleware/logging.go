package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs the HTTP requests.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start timer
		start := time.Now()

		// Wrap the response writer to capture the status code
		wrappedWriter := NewResponseWriter(w)

		// Process request
		next.ServeHTTP(wrappedWriter, r)

		// Log request details
		log.Printf("%s %s %s %d %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			wrappedWriter.status,
			time.Since(start),
		)
	})
}

// responseWriter is a custom http.ResponseWriter that captures the status code.
type responseWriter struct {
	http.ResponseWriter
	status int
}

// NewResponseWriter creates a new responseWriter instance.
func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK} // Default to 200 OK
}

// WriteHeader captures the status code and calls the original WriteHeader method.
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
