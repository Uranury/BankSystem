package middleware

import (
	"MockBankGo/auth"
	"context"
	"net/http"
	"strings"
)

// Define a custom type for context keys to avoid collisions
type contextKey string

// Define a non-exported variable for the user context key
var userIDKey = contextKey("userID")

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := auth.VerifyJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) (int64, bool) {
	UserID, ok := ctx.Value(userIDKey).(int64)
	return UserID, ok
}
