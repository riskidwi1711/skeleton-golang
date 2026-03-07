package main

import (
	"log"
	"net/http"
	"os"

	"github.com/riskidwi1711/itms-saas/services/tenant-service/internal/events"
	apphttp "github.com/riskidwi1711/itms-saas/services/tenant-service/internal/http"
	"github.com/riskidwi1711/itms-saas/services/tenant-service/internal/store"
)

func main() {
	port := getenv("PORT", "8082")
	dbURL := getenv("DATABASE_URL", "postgres://itms:itms@postgres-core:5432/core_db?sslmode=disable")
	natsURL := getenv("NATS_URL", "nats://nats:4222")
	eventSubject := getenv("NATS_SUBJECT_TENANT", "tenant.events")

	var st store.Store
	pgStore, err := store.NewPostgresStore(dbURL)
	if err != nil {
		log.Printf("postgres unavailable (%v), fallback to memory store", err)
		st = store.NewMemoryStore()
	} else {
		st = pgStore
	}

	var pub events.Publisher
	natsPub, err := events.NewNatsPublisher(natsURL, eventSubject)
	if err != nil {
		log.Printf("nats unavailable (%v), fallback to noop publisher", err)
		pub = events.NoopPublisher{}
	} else {
		pub = natsPub
	}

	h := apphttp.NewServer(st, pub)
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
