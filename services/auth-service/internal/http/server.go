package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	JWTSecret        string
	TenantServiceURL string
}

type loginRequest struct {
	Email string `json:"email"`
}

type registerTenantRequest struct {
	CompanyName string `json:"company_name"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Plan        string `json:"plan"`
}

type tenantResponse struct {
	Data struct {
		Tenant struct {
			TenantID    string `json:"tenant_id"`
			CompanyName string `json:"company_name"`
			AdminEmail  string `json:"admin_email"`
			Plan        string `json:"plan"`
		} `json:"tenant"`
	} `json:"data"`
}

func NewServer(cfg Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}

		var req loginRequest
		_ = json.NewDecoder(r.Body).Decode(&req)

		token, err := signToken(cfg.JWTSecret, "u-1", "tnt_demo", "owner", req.Email)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "token_sign_failed", "failed to sign token")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"token": token,
			"user":  map[string]any{"id": "u-1", "role": "owner", "email": req.Email},
		})
	})

	mux.HandleFunc("/api/v1/auth/register-tenant", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}

		var req registerTenantRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
			return
		}

		req.CompanyName = strings.TrimSpace(req.CompanyName)
		req.FullName = strings.TrimSpace(req.FullName)
		req.Email = strings.TrimSpace(req.Email)
		req.Plan = strings.TrimSpace(req.Plan)

		if req.CompanyName == "" || req.FullName == "" || req.Email == "" || req.Password == "" {
			writeError(w, http.StatusBadRequest, "validation_error", "company_name, full_name, email, password are required")
			return
		}
		if req.Plan == "" {
			req.Plan = "starter"
		}

		tenantURL := strings.TrimSuffix(cfg.TenantServiceURL, "/") + "/api/v1/tenants/onboard"
		bodyBytes, _ := json.Marshal(map[string]any{
			"company_name": req.CompanyName,
			"admin_email":  req.Email,
			"plan":         req.Plan,
		})

		onboardResp, err := http.Post(tenantURL, "application/json", bytes.NewReader(bodyBytes))
		if err != nil {
			writeError(w, http.StatusBadGateway, "tenant_service_unreachable", err.Error())
			return
		}
		defer onboardResp.Body.Close()

		if onboardResp.StatusCode >= 400 {
			writeError(w, http.StatusBadGateway, "tenant_onboard_failed", fmt.Sprintf("tenant service status %d", onboardResp.StatusCode))
			return
		}

		var tr tenantResponse
		if err := json.NewDecoder(onboardResp.Body).Decode(&tr); err != nil {
			writeError(w, http.StatusBadGateway, "tenant_decode_failed", "failed to decode tenant response")
			return
		}

		tenantID := tr.Data.Tenant.TenantID
		if tenantID == "" {
			tenantID = "tnt_demo"
		}

		token, err := signToken(cfg.JWTSecret, "u-owner", tenantID, "owner", req.Email)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "token_sign_failed", "failed to sign token")
			return
		}

		writeJSON(w, http.StatusCreated, map[string]any{
			"token": token,
			"user": map[string]any{
				"id":    "u-owner",
				"name":  req.FullName,
				"email": req.Email,
				"role":  "owner",
			},
			"tenant": tr.Data.Tenant,
		})
	})

	return mux
}

func signToken(secret, sub, tenantID, role, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       sub,
		"tenant_id": tenantID,
		"role":      role,
		"email":     email,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	})
	return token.SignedString([]byte(secret))
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
	})
}
