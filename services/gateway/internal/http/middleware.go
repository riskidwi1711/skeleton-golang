package http

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type customClaims struct {
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func withRequestID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get("X-Request-ID")
		if rid == "" {
			buf := make([]byte, 8)
			_, _ = rand.Read(buf)
			rid = hex.EncodeToString(buf)
		}
		w.Header().Set("X-Request-ID", rid)
		r.Header.Set("X-Request-ID", rid)
		next(w, r)
	}
}

func requireJWT(secret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeJSON(w, http.StatusUnauthorized, map[string]any{
					"error": map[string]any{"code": "missing_bearer_token"},
				})
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims := &customClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				writeJSON(w, http.StatusUnauthorized, map[string]any{
					"error": map[string]any{"code": "invalid_token"},
				})
				return
			}

			if claims.TenantID != "" {
				r.Header.Set("X-Tenant-ID", claims.TenantID)
			}
			if claims.Role != "" {
				r.Header.Set("X-User-Role", claims.Role)
			}

			next(w, r)
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
