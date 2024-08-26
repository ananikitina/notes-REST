package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ananikitina/notes-rest/internal/domain"
)

// Define a custom type for context keys
type contextKey string

const (
	UserIDKey contextKey = "userID"
)

// AuthMiddleware ensures that the request is authenticated.
func AuthMiddleware(jwtService domain.JWTServiceInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			userID := claims.UserID
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
