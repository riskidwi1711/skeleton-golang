package store

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(databaseURL string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	st := &PostgresStore{db: db}
	if err := st.ensureSchema(); err != nil {
		return nil, err
	}

	return st, nil
}

func (p *PostgresStore) ensureSchema() error {
	query := `
CREATE TABLE IF NOT EXISTS tenants (
  tenant_id TEXT PRIMARY KEY,
  company_name TEXT NOT NULL,
  admin_email TEXT NOT NULL,
  plan TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'trialing',
  timezone TEXT,
  primary_branch TEXT,
  primary_department TEXT,
  ticket_prefix TEXT,
  setup_completed_at TIMESTAMPTZ,
  created_at_utc TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'trialing';
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS timezone TEXT;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS primary_branch TEXT;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS primary_department TEXT;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS ticket_prefix TEXT;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS setup_completed_at TIMESTAMPTZ;
CREATE INDEX IF NOT EXISTS idx_tenants_created_at ON tenants (created_at_utc DESC);
`
	_, err := p.db.Exec(query)
	return err
}

func (p *PostgresStore) SaveTenant(ctx context.Context, tenant Tenant) (Tenant, error) {
	query := `
INSERT INTO tenants (tenant_id, company_name, admin_email, plan, status, created_at_utc)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING tenant_id, company_name, admin_email, plan, status, timezone, primary_branch, primary_department, ticket_prefix, setup_completed_at, created_at_utc
`

	var out Tenant
	var setup sql.NullTime
	err := p.db.QueryRowContext(
		ctx,
		query,
		tenant.TenantID,
		tenant.CompanyName,
		tenant.AdminEmail,
		tenant.Plan,
		tenant.Status,
		tenant.CreatedAtUTC,
	).Scan(&out.TenantID, &out.CompanyName, &out.AdminEmail, &out.Plan, &out.Status, &out.Timezone, &out.PrimaryBranch, &out.PrimaryDepartment, &out.TicketPrefix, &setup, &out.CreatedAtUTC)
	if err != nil {
		return Tenant{}, err
	}
	if setup.Valid {
		t := setup.Time
		out.SetupCompletedAt = &t
	}
	return out, nil
}

func (p *PostgresStore) ListTenant(ctx context.Context) ([]Tenant, error) {
	rows, err := p.db.QueryContext(ctx, `
SELECT tenant_id, company_name, admin_email, plan, status, timezone, primary_branch, primary_department, ticket_prefix, setup_completed_at, created_at_utc
FROM tenants
ORDER BY created_at_utc DESC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Tenant, 0)
	for rows.Next() {
		var t Tenant
		var setup sql.NullTime
		if err := rows.Scan(&t.TenantID, &t.CompanyName, &t.AdminEmail, &t.Plan, &t.Status, &t.Timezone, &t.PrimaryBranch, &t.PrimaryDepartment, &t.TicketPrefix, &setup, &t.CreatedAtUTC); err != nil {
			return nil, err
		}
		if setup.Valid {
			ts := setup.Time
			t.SetupCompletedAt = &ts
		}
		items = append(items, t)
	}
	return items, rows.Err()
}

func (p *PostgresStore) GetTenantByID(ctx context.Context, tenantID string) (Tenant, error) {
	query := `
SELECT tenant_id, company_name, admin_email, plan, status, timezone, primary_branch, primary_department, ticket_prefix, setup_completed_at, created_at_utc
FROM tenants
WHERE tenant_id = $1
`
	var out Tenant
	var setup sql.NullTime
	err := p.db.QueryRowContext(ctx, query, tenantID).
		Scan(&out.TenantID, &out.CompanyName, &out.AdminEmail, &out.Plan, &out.Status, &out.Timezone, &out.PrimaryBranch, &out.PrimaryDepartment, &out.TicketPrefix, &setup, &out.CreatedAtUTC)
	if err != nil {
		if err == sql.ErrNoRows {
			return Tenant{}, fmt.Errorf("tenant not found")
		}
		return Tenant{}, err
	}
	if setup.Valid {
		ts := setup.Time
		out.SetupCompletedAt = &ts
	}
	return out, nil
}

func (p *PostgresStore) UpdateTenantSetup(ctx context.Context, tenantID string, in TenantSetupInput) (Tenant, error) {
	query := `
UPDATE tenants
SET company_name = $1,
    timezone = $2,
    primary_branch = $3,
    primary_department = $4,
    ticket_prefix = $5,
    status = 'active',
    setup_completed_at = NOW()
WHERE tenant_id = $6
RETURNING tenant_id, company_name, admin_email, plan, status, timezone, primary_branch, primary_department, ticket_prefix, setup_completed_at, created_at_utc
`
	var out Tenant
	var setup sql.NullTime
	err := p.db.QueryRowContext(ctx, query, in.CompanyName, in.Timezone, in.PrimaryBranch, in.PrimaryDepartment, in.TicketPrefix, tenantID).
		Scan(&out.TenantID, &out.CompanyName, &out.AdminEmail, &out.Plan, &out.Status, &out.Timezone, &out.PrimaryBranch, &out.PrimaryDepartment, &out.TicketPrefix, &setup, &out.CreatedAtUTC)
	if err != nil {
		if err == sql.ErrNoRows {
			return Tenant{}, fmt.Errorf("tenant not found")
		}
		return Tenant{}, err
	}
	if setup.Valid {
		ts := setup.Time
		out.SetupCompletedAt = &ts
	}
	return out, nil
}
