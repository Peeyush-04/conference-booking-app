package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

// attaches unique request id to each request
// helps in debugging logs and distributed tracing
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()

		w.Header().Set("X-Request-ID", requestID)
		r.Header.Set("X-Request-ID", requestID) // for downstream handler

		next.ServeHTTP(w, r)
	})
}
