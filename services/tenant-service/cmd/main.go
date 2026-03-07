package main

import (
	"log"
	"net/http"
	"os"

	apphttp "github.com/riskidwi1711/itms-saas/services/tenant-service/internal/http"
	"github.com/riskidwi1711/itms-saas/services/tenant-service/internal/store"
)

func main() {
	port := getenv("PORT", "8082")
	st := store.NewMemoryStore()
	h := apphttp.NewServer(st)
	log.Printf("tenant-service listening on :%s", port)
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
