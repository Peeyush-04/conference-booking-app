package middleware

import (
	"backend/auth"
	"context"
	"net/http"
	"strings"
)

// reason: prevents other modules over-writing
type contextKey string

const (
	userIDKey contextKey = "user_id"
	roleKey   contextKey = "role"
)

// checks and validates from auth header
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get token from auth
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// validate token
		token, err := auth.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// extract claims
		userID, role, err := auth.ExtractClaims(token)
		if err != nil {
			http.Error(w, "Invalid claims", http.StatusUnauthorized)
			return
		}

		// attach user id and role to request context
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx = context.WithValue(ctx, roleKey, role)

		// call next handler with new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
