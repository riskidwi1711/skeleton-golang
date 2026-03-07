package store

import (
	"context"
	"database/sql"

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
  created_at_utc TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_tenants_created_at ON tenants (created_at_utc DESC);
`
	_, err := p.db.Exec(query)
	return err
}

func (p *PostgresStore) SaveTenant(ctx context.Context, tenant Tenant) (Tenant, error) {
	query := `
INSERT INTO tenants (tenant_id, company_name, admin_email, plan, created_at_utc)
VALUES ($1, $2, $3, $4, $5)
RETURNING tenant_id, company_name, admin_email, plan, created_at_utc
`

	var out Tenant
	err := p.db.QueryRowContext(
		ctx,
		query,
		tenant.TenantID,
		tenant.CompanyName,
		tenant.AdminEmail,
		tenant.Plan,
		tenant.CreatedAtUTC,
	).Scan(&out.TenantID, &out.CompanyName, &out.AdminEmail, &out.Plan, &out.CreatedAtUTC)
	if err != nil {
		return Tenant{}, err
	}

	return out, nil
}

func (p *PostgresStore) ListTenant(ctx context.Context) ([]Tenant, error) {
	rows, err := p.db.QueryContext(ctx, `
SELECT tenant_id, company_name, admin_email, plan, created_at_utc
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
		if err := rows.Scan(&t.TenantID, &t.CompanyName, &t.AdminEmail, &t.Plan, &t.CreatedAtUTC); err != nil {
			return nil, err
		}
		items = append(items, t)
	}
	return items, rows.Err()
}
