package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		schema := GetActiveSchema(r.Context())
		duration := time.Since(start)
		log.Printf("[HTTP] %s %s | status=%d | duration=%s | schema=%s | ip=%s",
			r.Method, r.URL.Path, wrapped.statusCode, duration, schema, r.RemoteAddr)
	})
}
