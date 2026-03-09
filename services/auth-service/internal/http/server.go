package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	JWTSecret        string
	TenantServiceURL string
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerTenantRequest struct {
	CompanyName string `json:"company_name"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Plan        string `json:"plan"`
}

type createUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

type updateUserRoleRequest struct {
	Role string `json:"role"`
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
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

type appUser struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	TenantID     string    `json:"tenant_id"`
	PasswordHash string    `json:"-"`
	Created      time.Time `json:"created_at"`
}

type roleDefinition struct {
	Name        string   `json:"name"`
	Label       string   `json:"label"`
	Permissions []string `json:"permissions"`
}

type customClaims struct {
	TenantID    string   `json:"tenant_id"`
	Role        string   `json:"role"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

var rolePermissions = map[string][]string{
	"owner":  {"dashboard:read", "tickets:read", "tickets:write", "assets:read", "assets:write", "users:read", "users:write", "tenants:read"},
	"admin":  {"dashboard:read", "tickets:read", "tickets:write", "assets:read", "assets:write", "users:read", "tenants:read"},
	"agent":  {"dashboard:read", "tickets:read", "tickets:write", "assets:read"},
	"viewer": {"dashboard:read", "tickets:read", "assets:read"},
}

type userStore struct {
	mu    sync.RWMutex
	items map[string]map[string]appUser
}

func newUserStore() *userStore {
	now := time.Now()
	// Hash the default demo password
	demoHash, _ := bcrypt.GenerateFromPassword([]byte("demo12345"), bcrypt.DefaultCost)
	return &userStore{items: map[string]map[string]appUser{
		"tnt_demo": {
			"owner@acme.com":  {ID: "u-1", Name: "Owner Demo", Email: "owner@acme.com", Role: "owner", TenantID: "tnt_demo", PasswordHash: string(demoHash), Created: now},
			"admin@acme.com":  {ID: "u-2", Name: "Admin Demo", Email: "admin@acme.com", Role: "admin", TenantID: "tnt_demo", PasswordHash: string(demoHash), Created: now},
			"agent@acme.com":  {ID: "u-3", Name: "Agent Demo", Email: "agent@acme.com", Role: "agent", TenantID: "tnt_demo", PasswordHash: string(demoHash), Created: now},
			"viewer@acme.com": {ID: "u-4", Name: "Viewer Demo", Email: "viewer@acme.com", Role: "viewer", TenantID: "tnt_demo", PasswordHash: string(demoHash), Created: now},
		},
	}}
}

func (s *userStore) findByEmail(tenantID, email string) (appUser, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tenant := s.items[tenantID]
	u, ok := tenant[email]
	return u, ok
}

func (s *userStore) listUsers(tenantID string) []appUser {
	s.mu.RLock()
	defer s.mu.RUnlock()
	users := make([]appUser, 0, len(s.items[tenantID]))
	for _, u := range s.items[tenantID] {
		u.PasswordHash = ""
		users = append(users, u)
	}
	sort.Slice(users, func(i, j int) bool { return users[i].Created.After(users[j].Created) })
	return users
}

func (s *userStore) createUser(tenantID, actorRole string, req createUserRequest) (appUser, error) {
	if actorRole != "owner" && actorRole != "admin" {
		return appUser{}, fmt.Errorf("forbidden")
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Name = strings.TrimSpace(req.Name)
	req.Role = strings.TrimSpace(strings.ToLower(req.Role))
	req.Password = strings.TrimSpace(req.Password)
	if req.Email == "" || req.Name == "" || req.Role == "" {
		return appUser{}, fmt.Errorf("name, email, role are required")
	}
	if req.Password == "" {
		req.Password = "welcome123"
	}
	if len(req.Password) < 6 {
		return appUser{}, fmt.Errorf("password must be at least 6 characters")
	}
	if _, ok := rolePermissions[req.Role]; !ok {
		return appUser{}, fmt.Errorf("invalid role")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return appUser{}, fmt.Errorf("failed to hash password")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.items[tenantID] == nil {
		s.items[tenantID] = map[string]appUser{}
	}
	if _, exists := s.items[tenantID][req.Email]; exists {
		return appUser{}, fmt.Errorf("email already exists")
	}
	u := appUser{ID: fmt.Sprintf("u-%d", len(s.items[tenantID])+1), Name: req.Name, Email: req.Email, Role: req.Role, TenantID: tenantID, PasswordHash: string(hashedPassword), Created: time.Now()}
	s.items[tenantID][req.Email] = u
	u.PasswordHash = ""
	return u, nil
}

func (s *userStore) upsertUser(u appUser) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.items[u.TenantID] == nil {
		s.items[u.TenantID] = map[string]appUser{}
	}
	s.items[u.TenantID][u.Email] = u
}

func (s *userStore) updateRole(tenantID, actorRole, userID, role string) (appUser, error) {
	if actorRole != "owner" {
		return appUser{}, fmt.Errorf("forbidden")
	}
	role = strings.TrimSpace(strings.ToLower(role))
	if _, ok := rolePermissions[role]; !ok {
		return appUser{}, fmt.Errorf("invalid role")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for email, u := range s.items[tenantID] {
		if u.ID == userID {
			u.Role = role
			s.items[tenantID][email] = u
			u.PasswordHash = ""
			return u, nil
		}
	}
	return appUser{}, fmt.Errorf("user not found")
}

func (s *userStore) changePassword(tenantID, email, currentPassword, newPassword string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, exists := s.items[tenantID][email]
	if !exists {
		return fmt.Errorf("user not found")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(currentPassword)); err != nil {
		return fmt.Errorf("invalid current password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}

	u.PasswordHash = string(hashedPassword)
	s.items[tenantID][email] = u
	return nil
}

func NewServer(cfg Config) http.Handler {
	mux := http.NewServeMux()
	users := newUserStore()

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
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
			return
		}
		req.Email = strings.TrimSpace(strings.ToLower(req.Email))
		req.Password = strings.TrimSpace(req.Password)
		if req.Email == "" || req.Password == "" {
			writeError(w, http.StatusBadRequest, "validation_error", "email and password are required")
			return
		}
		u, ok := users.findByEmail("tnt_demo", req.Email)
		if !ok {
			writeError(w, http.StatusUnauthorized, "invalid_credentials", "invalid email or password")
			return
		}
		// Compare password with bcrypt
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
			writeError(w, http.StatusUnauthorized, "invalid_credentials", "invalid email or password")
			return
		}
		token, err := signToken(cfg.JWTSecret, u)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "token_sign_failed", "failed to sign token")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"token": token,
			"user":  map[string]any{"id": u.ID, "role": u.Role, "email": u.Email, "name": u.Name, "permissions": rolePermissions[u.Role]},
		})
	})

	mux.HandleFunc("/api/v1/auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}
		tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
		email := strings.TrimSpace(strings.ToLower(r.Header.Get("X-User-Email")))
		if tenantID == "" || email == "" {
			writeError(w, http.StatusUnauthorized, "missing_identity", "missing user identity")
			return
		}
		u, ok := users.findByEmail(tenantID, email)
		if !ok {
			writeError(w, http.StatusUnauthorized, "user_not_found", "user not found")
			return
		}
		token, err := signToken(cfg.JWTSecret, u)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "token_sign_failed", "failed to sign token")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"token": token,
			"user":  map[string]any{"id": u.ID, "role": u.Role, "email": u.Email, "name": u.Name, "permissions": rolePermissions[u.Role]},
		})
	})

	mux.HandleFunc("/api/v1/auth/change-password", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}
		tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
		email := strings.TrimSpace(strings.ToLower(r.Header.Get("X-User-Email")))
		if tenantID == "" || email == "" {
			writeError(w, http.StatusUnauthorized, "missing_identity", "missing user identity")
			return
		}

		var req changePasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
			return
		}

		if req.CurrentPassword == "" || req.NewPassword == "" {
			writeError(w, http.StatusBadRequest, "validation_error", "current_password and new_password are required")
			return
		}

		if len(req.NewPassword) < 6 {
			writeError(w, http.StatusBadRequest, "validation_error", "new password must be at least 6 characters")
			return
		}

		if err := users.changePassword(tenantID, email, req.CurrentPassword, req.NewPassword); err != nil {
			if err.Error() == "invalid current password" {
				writeError(w, http.StatusUnauthorized, "invalid_password", err.Error())
				return
			}
			writeError(w, http.StatusInternalServerError, "change_failed", err.Error())
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"message": "password changed successfully"})
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
		req.Email = strings.TrimSpace(strings.ToLower(req.Email))
		req.Password = strings.TrimSpace(req.Password)
		req.Plan = strings.TrimSpace(req.Plan)
		if req.CompanyName == "" || req.FullName == "" || req.Email == "" || req.Password == "" {
			writeError(w, http.StatusBadRequest, "validation_error", "company_name, full_name, email, password are required")
			return
		}
		if len(req.Password) < 6 {
			writeError(w, http.StatusBadRequest, "validation_error", "password must be at least 6 characters")
			return
		}
		if req.Plan == "" {
			req.Plan = "starter"
		}

		tenantURL := strings.TrimSuffix(cfg.TenantServiceURL, "/") + "/api/v1/tenants/onboard"
		bodyBytes, _ := json.Marshal(map[string]any{"company_name": req.CompanyName, "admin_email": req.Email, "plan": req.Plan})
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

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "password_hash_failed", "failed to hash password")
			return
		}

		u := appUser{ID: "u-owner", Name: req.FullName, Email: req.Email, Role: "owner", TenantID: tenantID, PasswordHash: string(hashedPassword), Created: time.Now()}
		users.upsertUser(u)
		token, err := signToken(cfg.JWTSecret, u)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "token_sign_failed", "failed to sign token")
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{
			"token":  token,
			"user":   map[string]any{"id": u.ID, "name": u.Name, "email": u.Email, "role": u.Role},
			"tenant": tr.Data.Tenant,
		})
	})

	mux.HandleFunc("/api/v1/roles", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}
		roles := make([]roleDefinition, 0, len(rolePermissions))
		for name, perms := range rolePermissions {
			roles = append(roles, roleDefinition{Name: name, Label: strings.Title(name), Permissions: perms})
		}
		sort.Slice(roles, func(i, j int) bool { return roles[i].Name < roles[j].Name })
		writeJSON(w, http.StatusOK, map[string]any{"data": map[string]any{"items": roles}})
	})

	mux.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
		if tenantID == "" {
			tenantID = "tnt_demo"
		}
		actorRole := strings.TrimSpace(r.Header.Get("X-User-Role"))
		if actorRole == "" {
			actorRole = "owner"
		}
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, http.StatusOK, map[string]any{"data": map[string]any{"items": users.listUsers(tenantID)}})
		case http.MethodPost:
			var req createUserRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
				return
			}
			u, err := users.createUser(tenantID, actorRole, req)
			if err != nil {
				if err.Error() == "forbidden" {
					writeError(w, http.StatusForbidden, "forbidden", "insufficient role")
					return
				}
				writeError(w, http.StatusBadRequest, "validation_error", err.Error())
				return
			}
			writeJSON(w, http.StatusCreated, map[string]any{"data": map[string]any{"user": u}})
		default:
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		}
	})

	mux.HandleFunc("/api/v1/users/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}
		if !strings.HasSuffix(r.URL.Path, "/role") {
			writeError(w, http.StatusNotFound, "not_found", "endpoint not found")
			return
		}
		userID := strings.Trim(strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/users/"), "/role"), "/")
		if userID == "" {
			writeError(w, http.StatusBadRequest, "validation_error", "user id is required")
			return
		}
		tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
		if tenantID == "" {
			tenantID = "tnt_demo"
		}
		actorRole := strings.TrimSpace(r.Header.Get("X-User-Role"))
		var req updateUserRoleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
			return
		}
		u, err := users.updateRole(tenantID, actorRole, userID, req.Role)
		if err != nil {
			switch err.Error() {
			case "forbidden":
				writeError(w, http.StatusForbidden, "forbidden", "only owner can change role")
			case "invalid role":
				writeError(w, http.StatusBadRequest, "validation_error", err.Error())
			case "user not found":
				writeError(w, http.StatusNotFound, "not_found", err.Error())
			default:
				writeError(w, http.StatusBadRequest, "validation_error", err.Error())
			}
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"data": map[string]any{"user": u}})
	})

	return mux
}

func signToken(secret string, u appUser) (string, error) {
	perms := rolePermissions[u.Role]
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims{
		TenantID:    u.TenantID,
		Role:        u.Role,
		Email:       u.Email,
		Permissions: perms,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   u.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	return token.SignedString([]byte(secret))
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]any{"code": code, "message": message}})
}
