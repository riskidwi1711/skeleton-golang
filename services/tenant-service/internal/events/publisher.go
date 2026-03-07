package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/riskidwi1711/itms-saas/services/tenant-service/internal/store"
)

type Publisher interface {
	PublishTenantOnboarded(ctx context.Context, tenant store.Tenant) error
}

type NatsPublisher struct {
	conn    *nats.Conn
	subject string
}

type NoopPublisher struct{}

func NewNatsPublisher(url, subject string) (*NatsPublisher, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsPublisher{conn: conn, subject: subject}, nil
}

func (n *NatsPublisher) PublishTenantOnboarded(_ context.Context, tenant store.Tenant) error {
	payload, err := json.Marshal(map[string]any{
		"event": "tenant.onboarded",
		"data":  tenant,
	})
	if err != nil {
		return err
	}
	return n.conn.Publish(n.subject, payload)
}

func (NoopPublisher) PublishTenantOnboarded(_ context.Context, _ store.Tenant) error {
	log.Printf("noop publisher: tenant.onboarded not sent")
	return nil
}
