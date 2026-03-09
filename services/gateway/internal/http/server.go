package http

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Config struct {
	AuthServiceURL   string
	TenantServiceURL string
	TicketServiceURL string
	AssetServiceURL  string
	JWTSecret        string
}

func NewServer(cfg Config) http.Handler {
	mux := http.NewServeMux()

	authProxy := mustProxy(cfg.AuthServiceURL)
	tenantProxy := mustProxy(cfg.TenantServiceURL)
	ticketProxy := mustProxy(cfg.TicketServiceURL)
	assetProxy := mustProxy(cfg.AssetServiceURL)
	jwtGuard := requireJWT(cfg.JWTSecret)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/api/v1/auth/refresh", withRequestID(jwtGuard(func(w http.ResponseWriter, r *http.Request) {
		authProxy.ServeHTTP(w, r)
	})))

	mux.HandleFunc("/api/v1/auth/change-password", withRequestID(jwtGuard(func(w http.ResponseWriter, r *http.Request) {
		authProxy.ServeHTTP(w, r)
	})))

	mux.HandleFunc("/api/v1/auth/", withRequestID(func(w http.ResponseWriter, r *http.Request) {
		authProxy.ServeHTTP(w, r)
	}))

	mux.HandleFunc("/api/v1/roles", withRequestID(jwtGuard(requirePermission("users:read")(func(w http.ResponseWriter, r *http.Request) {
		authProxy.ServeHTTP(w, r)
	}))))

	mux.HandleFunc("/api/v1/users", withRequestID(jwtGuard(requirePermission("users:read", "users:write")(func(w http.ResponseWriter, r *http.Request) {
		authProxy.ServeHTTP(w, r)
	}))))
	mux.HandleFunc("/api/v1/users/", withRequestID(jwtGuard(requirePermission("users:write")(func(w http.ResponseWriter, r *http.Request) {
		authProxy.ServeHTTP(w, r)
	}))))

	mux.HandleFunc("/api/v1/tenants/onboard", withRequestID(func(w http.ResponseWriter, r *http.Request) {
		tenantProxy.ServeHTTP(w, r)
	}))

	mux.HandleFunc("/api/v1/tenants", withRequestID(jwtGuard(requirePermission("tenants:read")(func(w http.ResponseWriter, r *http.Request) {
		tenantProxy.ServeHTTP(w, r)
	}))))

	mux.HandleFunc("/api/v1/tenants/", withRequestID(jwtGuard(permissionByMethod(
		func(w http.ResponseWriter, r *http.Request) { tenantProxy.ServeHTTP(w, r) },
		map[string][]string{http.MethodGet: []string{"tenants:read"}, http.MethodPut: []string{"users:write"}, http.MethodPatch: []string{"users:write"}},
	))))

	mux.HandleFunc("/api/v1/tickets", withRequestID(jwtGuard(permissionByMethod(
		func(w http.ResponseWriter, r *http.Request) { ticketProxy.ServeHTTP(w, r) },
		map[string][]string{http.MethodGet: []string{"tickets:read"}, http.MethodPost: []string{"tickets:write"}},
	))))
	mux.HandleFunc("/api/v1/tickets/", withRequestID(jwtGuard(permissionByMethod(
		func(w http.ResponseWriter, r *http.Request) { ticketProxy.ServeHTTP(w, r) },
		map[string][]string{http.MethodGet: []string{"tickets:read"}, http.MethodPut: []string{"tickets:write"}, http.MethodPatch: []string{"tickets:write"}},
	))))

	mux.HandleFunc("/api/v1/assets", withRequestID(jwtGuard(permissionByMethod(
		func(w http.ResponseWriter, r *http.Request) { assetProxy.ServeHTTP(w, r) },
		map[string][]string{http.MethodGet: []string{"assets:read"}, http.MethodPost: []string{"assets:write"}},
	))))
	mux.HandleFunc("/api/v1/assets/", withRequestID(jwtGuard(permissionByMethod(
		func(w http.ResponseWriter, r *http.Request) { assetProxy.ServeHTTP(w, r) },
		map[string][]string{http.MethodGet: []string{"assets:read"}, http.MethodPut: []string{"assets:write"}, http.MethodPatch: []string{"assets:write"}, http.MethodDelete: []string{"assets:write"}},
	))))

	return mux
}

func permissionByMethod(next http.HandlerFunc, rules map[string][]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		required, ok := rules[r.Method]
		if !ok {
			next(w, r)
			return
		}
		requirePermission(required...)(next)(w, r)
	}
}

func mustProxy(raw string) *httputil.ReverseProxy {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	origDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		origDirector(req)
		req.Header.Set("X-Forwarded-By", "itms-gateway")
		req.URL.Path = strings.Replace(req.URL.Path, "//", "/", -1)
	}
	return proxy
}
