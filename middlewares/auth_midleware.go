package middlewares

import (
	"akmmp241/go-jwt/configs"
	"akmmp241/go-jwt/helpers"
	"context"
	"net/http"
)

type AuthMiddleware struct {
	config *configs.Config
}

func NewAuthMiddleware(config *configs.Config) *AuthMiddleware {
	return &AuthMiddleware{config: config}
}

func (m AuthMiddleware) JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")

		if accessToken == "" {
			helpers.Response(w, http.StatusUnauthorized, "Invalid Token", nil)
			return
		}

		key := m.config.C.GetString("JWT_KEY")

		claims, err := helpers.ValidateToken(accessToken, []byte(key))
		if err != nil {
			helpers.Response(w, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
