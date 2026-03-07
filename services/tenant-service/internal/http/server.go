package http

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/riskidwi1711/itms-saas/services/tenant-service/internal/events"
	"github.com/riskidwi1711/itms-saas/services/tenant-service/internal/store"
)

type onboardRequest struct {
	CompanyName string `json:"company_name"`
	AdminEmail  string `json:"admin_email"`
	Plan        string `json:"plan"`
}

type Server struct {
	store     store.Store
	publisher events.Publisher
}

func NewServer(st store.Store, pub events.Publisher) http.Handler {
	s := &Server{store: st, publisher: pub}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/api/v1/tenants/onboard", s.onboard)
	mux.HandleFunc("/api/v1/tenants", s.list)
	return withRequestID(mux)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) onboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	var req onboardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	if strings.TrimSpace(req.CompanyName) == "" || strings.TrimSpace(req.AdminEmail) == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "company_name and admin_email are required")
		return
	}
	if strings.TrimSpace(req.Plan) == "" {
		req.Plan = "starter"
	}

	tenant := store.Tenant{
		TenantID:     "tnt_" + randomHex(6),
		CompanyName:  strings.TrimSpace(req.CompanyName),
		AdminEmail:   strings.TrimSpace(req.AdminEmail),
		Plan:         strings.TrimSpace(req.Plan),
		CreatedAtUTC: time.Now().UTC(),
	}

	tenant, err := s.store.SaveTenant(r.Context(), tenant)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "persist_failed", err.Error())
		return
	}

	_ = s.publisher.PublishTenantOnboarded(r.Context(), tenant)

	writeJSON(w, http.StatusCreated, map[string]any{
		"data": map[string]any{
			"tenant": tenant,
		},
	})
}

func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	items, err := s.store.ListTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"items": items,
		},
	})
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

func randomHex(length int) string {
	buf := make([]byte, length)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}
