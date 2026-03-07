package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	JWTSecret string
}

type loginRequest struct {
	Email string `json:"email"`
}

func NewServer(cfg Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
			return
		}

		var req loginRequest
		_ = json.NewDecoder(r.Body).Decode(&req)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":       "u-1",
			"tenant_id": "tnt_demo",
			"role":      "owner",
			"email":     req.Email,
			"exp":       time.Now().Add(24 * time.Hour).Unix(),
			"iat":       time.Now().Unix(),
		})

		signed, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			writeError(w, http.StatusInternalServerError, "token_sign_failed")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"token": signed,
			"user":  map[string]any{"id": "u-1", "role": "owner"},
		})
	})
	return mux
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]any{
			"code": code,
		},
	})
}
