package store

import (
	"context"
	"time"
)

type Tenant struct {
	TenantID          string     `json:"tenant_id"`
	CompanyName       string     `json:"company_name"`
	AdminEmail        string     `json:"admin_email"`
	Plan              string     `json:"plan"`
	Status            string     `json:"status"`
	Timezone          string     `json:"timezone,omitempty"`
	PrimaryBranch     string     `json:"primary_branch,omitempty"`
	PrimaryDepartment string     `json:"primary_department,omitempty"`
	TicketPrefix      string     `json:"ticket_prefix,omitempty"`
	SetupCompletedAt  *time.Time `json:"setup_completed_at,omitempty"`
	CreatedAtUTC      time.Time  `json:"created_at_utc"`
}

type TenantSetupInput struct {
	CompanyName       string `json:"company_name"`
	Timezone          string `json:"timezone"`
	PrimaryBranch     string `json:"primary_branch"`
	PrimaryDepartment string `json:"primary_department"`
	TicketPrefix      string `json:"ticket_prefix"`
}

type Store interface {
	SaveTenant(ctx context.Context, tenant Tenant) (Tenant, error)
	ListTenant(ctx context.Context) ([]Tenant, error)
	UpdateTenantSetup(ctx context.Context, tenantID string, in TenantSetupInput) (Tenant, error)
}
