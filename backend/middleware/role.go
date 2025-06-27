package middleware

import "net/http"

// checks role
func RequireRole(expectedRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleVal, ok := r.Context().Value(roleKey).(string)
			if !ok || roleVal != expectedRole {
				http.Error(w, "Forbidden: insufficient permission", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
