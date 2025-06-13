package middleware

import (
	"log"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// Override WriteHeader to capture status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the ResponseWriter
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		log.Printf("incoming request %s %s from %s - status %d",
			r.Method, r.URL.Path, r.RemoteAddr, rw.statusCode)
	})
}
