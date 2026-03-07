package store

import (
	"sync"
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
	SaveTenant(tenant Tenant)
	ListTenant() []Tenant
}

type MemoryStore struct {
	mu      sync.RWMutex
	tenants []Tenant
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) SaveTenant(tenant Tenant) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tenants = append(m.tenants, tenant)
}

func (m *MemoryStore) ListTenant() []Tenant {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Tenant, len(m.tenants))
	copy(out, m.tenants)
	return out
}
