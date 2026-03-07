package store

import (
	"context"
	"time"
)

type Tenant struct {
	TenantID     string    `json:"tenant_id"`
	CompanyName  string    `json:"company_name"`
	AdminEmail   string    `json:"admin_email"`
	Plan         string    `json:"plan"`
	CreatedAtUTC time.Time `json:"created_at_utc"`
}

type Store interface {
	SaveTenant(ctx context.Context, tenant Tenant) (Tenant, error)
	ListTenant(ctx context.Context) ([]Tenant, error)
}
