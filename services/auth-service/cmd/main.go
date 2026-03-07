package main

import (
	"log"
	"net/http"
	"os"

	apphttp "github.com/riskidwi1711/itms-saas/services/auth-service/internal/http"
)

func main() {
	port := getenv("PORT", "8081")
	h := apphttp.NewServer()
	log.Printf("auth-service listening on :%s", port)
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
