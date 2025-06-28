package middleware

import (
	"net/http"
	"sync"
	"time"
)

var (
	visitor = make(map[string]time.Time)
	mu      sync.Mutex
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// transaction check
		mu.Lock()
		lastSeen, found := visitor[r.RemoteAddr]
		if found && time.Since(lastSeen) < time.Second {
			mu.Unlock()
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		visitor[r.RemoteAddr] = time.Now()
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
