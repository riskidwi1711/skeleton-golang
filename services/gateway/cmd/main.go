package main

import (
	"log"
	"net/http"
	"os"

	apphttp "github.com/riskidwi1711/itms-saas/services/gateway/internal/http"
)

func main() {
	port := getenv("PORT", "8080")
	cfg := apphttp.Config{
		AuthServiceURL:   getenv("AUTH_SERVICE_URL", "http://auth-service:8081"),
		TenantServiceURL: getenv("TENANT_SERVICE_URL", "http://tenant-service:8082"),
		TicketServiceURL: getenv("TICKET_SERVICE_URL", "http://ticket-service:8083"),
		AssetServiceURL:  getenv("ASSET_SERVICE_URL", "http://asset-service:8084"),
		JWTSecret:        getenv("JWT_SECRET", "dev-secret"),
	}

	h := apphttp.NewServer(cfg)
	log.Printf("gateway listening on :%s", port)
	if err := http.ListenAndServe(":"+port, h); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
