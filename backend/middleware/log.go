package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// continue with next handler
		next.ServeHTTP(w, r)

		// log after response is served
		duration := time.Since(start)
		log.Printf("%s %s %s", r.Method, r.URL.Path, duration)
	})
}
