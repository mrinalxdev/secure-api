package middleware

import (
	"context"
	"net/http"
	"secure-api/utils"
	"strings"

	"github.com/gorilla/mux"
)

func AuthMiddleware(role string) mux.MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if !strings.HasPrefix(authHeader, "Bearer ") {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }

            tokenString := strings.TrimPrefix(authHeader, "Bearer ")
            claims, err := utils.ValidateToken(tokenString)
            if err != nil {
                http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
                return
            }

            // Role check
            if role != "" && claims.Role != role {
                http.Error(w, "Insufficient permissions", http.StatusForbidden)
                return
            }

            // Pass user info to next handler
            ctx := context.WithValue(r.Context(), "userID", claims.UserID)
            ctx = context.WithValue(ctx, "role", claims.Role)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}