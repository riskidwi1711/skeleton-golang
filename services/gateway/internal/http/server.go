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
	JWTSecret        string
}

func NewServer(cfg Config) http.Handler {
	mux := http.NewServeMux()

	authProxy := mustProxy(cfg.AuthServiceURL)
	tenantProxy := mustProxy(cfg.TenantServiceURL)
	jwtGuard := requireJWT(cfg.JWTSecret)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/api/v1/auth/", withRequestID(func(w http.ResponseWriter, r *http.Request) {
		authProxy.ServeHTTP(w, r)
	}))

	mux.HandleFunc("/api/v1/tenants/onboard", withRequestID(func(w http.ResponseWriter, r *http.Request) {
		tenantProxy.ServeHTTP(w, r)
	}))

	mux.HandleFunc("/api/v1/tenants", withRequestID(jwtGuard(func(w http.ResponseWriter, r *http.Request) {
		tenantProxy.ServeHTTP(w, r)
	})))

	mux.HandleFunc("/api/v1/tenants/", withRequestID(jwtGuard(func(w http.ResponseWriter, r *http.Request) {
		tenantProxy.ServeHTTP(w, r)
	})))

	return mux
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
