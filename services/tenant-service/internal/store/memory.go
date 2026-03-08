package store

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type MemoryStore struct {
	mu      sync.RWMutex
	tenants []Tenant
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) SaveTenant(_ context.Context, tenant Tenant) (Tenant, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tenants = append(m.tenants, tenant)
	return tenant, nil
}

func (m *MemoryStore) ListTenant(_ context.Context) ([]Tenant, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Tenant, len(m.tenants))
	copy(out, m.tenants)
	return out, nil
}

func (m *MemoryStore) GetTenantByID(_ context.Context, tenantID string) (Tenant, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, t := range m.tenants {
		if t.TenantID == tenantID {
			return t, nil
		}
	}
	return Tenant{}, fmt.Errorf("tenant not found")
}

func (m *MemoryStore) UpdateTenantSetup(_ context.Context, tenantID string, in TenantSetupInput) (Tenant, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, t := range m.tenants {
		if t.TenantID != tenantID {
			continue
		}
		now := time.Now().UTC()
		t.CompanyName = in.CompanyName
		t.Timezone = in.Timezone
		t.PrimaryBranch = in.PrimaryBranch
		t.PrimaryDepartment = in.PrimaryDepartment
		t.TicketPrefix = in.TicketPrefix
		t.Status = "active"
		t.SetupCompletedAt = &now
		m.tenants[i] = t
		return t, nil
	}
	return Tenant{}, fmt.Errorf("tenant not found")
}
