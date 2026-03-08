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

type setupRequest struct {
	CompanyName       string `json:"company_name"`
	Timezone          string `json:"timezone"`
	PrimaryBranch     string `json:"primary_branch"`
	PrimaryDepartment string `json:"primary_department"`
	TicketPrefix      string `json:"ticket_prefix"`
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
	mux.HandleFunc("/api/v1/tenants/", s.tenantRoutes)
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
		Status:       "trialing",
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

func (s *Server) tenantRoutes(w http.ResponseWriter, r *http.Request) {
	trimmed := strings.TrimPrefix(r.URL.Path, "/api/v1/tenants/")
	trimmed = strings.Trim(trimmed, "/")
	if trimmed == "" {
		writeError(w, http.StatusNotFound, "not_found", "endpoint not found")
		return
	}

	if !strings.HasSuffix(trimmed, "/setup") {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}
		tenantID := trimmed
		headerTenant := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
		if headerTenant != "" && headerTenant != tenantID {
			writeError(w, http.StatusForbidden, "forbidden", "cross-tenant read not allowed")
			return
		}
		tenant, err := s.store.GetTenantByID(r.Context(), tenantID)
		if err != nil {
			if err.Error() == "tenant not found" {
				writeError(w, http.StatusNotFound, "not_found", err.Error())
				return
			}
			writeError(w, http.StatusInternalServerError, "query_failed", err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"data": map[string]any{"tenant": tenant}})
		return
	}
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	tenantID := strings.TrimSuffix(trimmed, "/setup")
	tenantID = strings.Trim(tenantID, "/")
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "tenant id required")
		return
	}

	headerTenant := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
	if headerTenant != "" && headerTenant != tenantID {
		writeError(w, http.StatusForbidden, "forbidden", "cross-tenant setup not allowed")
		return
	}
	role := strings.ToLower(strings.TrimSpace(r.Header.Get("X-User-Role")))
	if role != "owner" && role != "admin" {
		writeError(w, http.StatusForbidden, "forbidden", "only owner/admin can complete setup")
		return
	}

	var req setupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}

	if strings.TrimSpace(req.CompanyName) == "" || strings.TrimSpace(req.Timezone) == "" || strings.TrimSpace(req.PrimaryBranch) == "" || strings.TrimSpace(req.PrimaryDepartment) == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "company_name, timezone, primary_branch, primary_department are required")
		return
	}
	if strings.TrimSpace(req.TicketPrefix) == "" {
		req.TicketPrefix = "INC"
	}

	updated, err := s.store.UpdateTenantSetup(r.Context(), tenantID, store.TenantSetupInput{
		CompanyName:       strings.TrimSpace(req.CompanyName),
		Timezone:          strings.TrimSpace(req.Timezone),
		PrimaryBranch:     strings.TrimSpace(req.PrimaryBranch),
		PrimaryDepartment: strings.TrimSpace(req.PrimaryDepartment),
		TicketPrefix:      strings.ToUpper(strings.TrimSpace(req.TicketPrefix)),
	})
	if err != nil {
		if err.Error() == "tenant not found" {
			writeError(w, http.StatusNotFound, "not_found", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "update_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": map[string]any{"tenant": updated}})
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
