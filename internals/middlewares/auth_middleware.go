package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(authService AuthServiceInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			cookie, err := r.Cookie("sid")
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			userID, err := authService.GetUserIDFromSession(cookie.Value)
			if err != nil {
				http.Error(w, "invalid session", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) int64 {
	userID, ok := r.Context().Value(UserIDKey).(int64)
	if !ok {
		return 0
	}
	return userID
}
