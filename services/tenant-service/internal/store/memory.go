package store

import (
	"context"
	"sync"
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
