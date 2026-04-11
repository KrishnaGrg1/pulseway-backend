package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/KrishnaGrg1/pulseway/internal/response"
	"github.com/KrishnaGrg1/pulseway/internal/store"
	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func JwtAuth(s *store.Store, jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractTokenFromHeader(r)
			if tokenStr == "" {
				response.Error(w, http.StatusUnauthorized, "Unauthorized", "AUTH_001", "No authentication token found")
				return
			}
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				response.Error(w, http.StatusUnauthorized, "Unauthorized", "AUTH_002", "Invalid token")
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				response.Error(w, http.StatusUnauthorized, "Unauthorized", "AUTH_002", "Invalid token")
				return
			}

			userIDFloat, ok := claims["userId"].(float64)
			if !ok {
				response.Error(w, http.StatusUnauthorized, "Unauthorized", "AUTH_002", "Invalid token")
				return
			}

			userID := int64(userIDFloat)

			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					response.Error(w, http.StatusUnauthorized, "Unauthorized", "AUTH_003", "Token expired")
					return
				}
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractTokenFromHeader(r *http.Request) string {
	if cookie, err := r.Cookie("token"); err == nil {
		if cookie.Value != "" {
			return cookie.Value
		}
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}
	return ""

}
