package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/CRAZYKAYZY/chirpy/internal/config"
	"github.com/golang-jwt/jwt"
)

type MyConfig struct {
	config.ApiConfig
}

func (cfg *MyConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Return the secret key used to sign the token
			return []byte(cfg.JwtSecret), nil
		})
		if err != nil {
			http.Error(w, "Unauthorized Token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Unauthorized, Not Valid", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
