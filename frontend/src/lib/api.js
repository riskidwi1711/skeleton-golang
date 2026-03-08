import { getToken } from './auth'

const API_BASE = import.meta.env.VITE_API_BASE || ''

function authHeaders() {
  const token = getToken()
  return token ? { Authorization: `Bearer ${token}` } : {}
}

export async function login(email = 'owner@acme.com', password = 'demo12345') {
  const res = await fetch(`${API_BASE}/api/v1/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({}))
    throw new Error(error?.error?.message || 'Login failed')
  }
  return res.json()
}

export async function refreshSession() {
  const res = await fetch(`${API_BASE}/api/v1/auth/refresh`, {
    method: 'POST',
    headers: { ...authHeaders() },
  })
  if (!res.ok) throw new Error('Session refresh failed')
  return res.json()
}

export async function listTenants(token) {
  const res = await fetch(`${API_BASE}/api/v1/tenants`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })

  if (!res.ok) throw new Error('Fetch tenants failed')
  const data = await res.json()
  return data?.data?.items || []
}

export async function onboardTenant(payload) {
  const res = await fetch(`${API_BASE}/api/v1/tenants/onboard`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({}))
    throw new Error(error?.error?.message || 'Onboarding failed')
  }

  return res.json()
}

export async function registerTenant(payload) {
  const res = await fetch(`${API_BASE}/api/v1/auth/register-tenant`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({}))
    throw new Error(error?.error?.message || 'Registration failed')
  }

  return res.json()
}

export async function listRoles() {
  const res = await fetch(`${API_BASE}/api/v1/roles`, { headers: { ...authHeaders() } })
  if (!res.ok) throw new Error('Fetch roles failed')
  const data = await res.json()
  return data?.data?.items || []
}

export async function listUsers() {
  const res = await fetch(`${API_BASE}/api/v1/users`, { headers: { ...authHeaders() } })
  if (!res.ok) throw new Error('Fetch users failed')
  const data = await res.json()
  return data?.data?.items || []
}

export async function createUser(payload) {
  const res = await fetch(`${API_BASE}/api/v1/users`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...authHeaders() },
    body: JSON.stringify(payload),
  })
  if (!res.ok) {
    const error = await res.json().catch(() => ({}))
    throw new Error(error?.error?.message || 'Create user failed')
  }
  const data = await res.json()
  return data?.data?.user
}

export async function updateUserRole(id, role) {
  const res = await fetch(`${API_BASE}/api/v1/users/${id}/role`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json', ...authHeaders() },
    body: JSON.stringify({ role }),
  })
  if (!res.ok) {
    const error = await res.json().catch(() => ({}))
    throw new Error(error?.error?.message || 'Update role failed')
  }
  const data = await res.json()
  return data?.data?.user
}

export async function listTickets(params = {}) {
  const query = new URLSearchParams(params).toString()
  const url = `${API_BASE}/api/v1/tickets${query ? `?${query}` : ''}`
  const res = await fetch(url, { headers: { ...authHeaders() } })
  if (!res.ok) throw new Error('Fetch tickets failed')
  return res.json()
}

export async function createTicket(payload) {
  const res = await fetch(`${API_BASE}/api/v1/tickets`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...authHeaders() },
    body: JSON.stringify(payload),
  })
  if (!res.ok) throw new Error('Create ticket failed')
  return res.json()
}

export async function updateTicket(id, payload) {
  const res = await fetch(`${API_BASE}/api/v1/tickets/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...authHeaders() },
    body: JSON.stringify(payload),
  })
  if (!res.ok) throw new Error('Update ticket failed')
  return res.json()
}

export async function ticketStats() {
  const res = await fetch(`${API_BASE}/api/v1/tickets/stats`, { headers: { ...authHeaders() } })
  if (!res.ok) throw new Error('Fetch ticket stats failed')
  return res.json()
}

export async function listAssets(params = {}) {
  const query = new URLSearchParams(params).toString()
  const url = `${API_BASE}/api/v1/assets${query ? `?${query}` : ''}`
  const res = await fetch(url, { headers: { ...authHeaders() } })
  if (!res.ok) throw new Error('Fetch assets failed')
  return res.json()
}

export async function createAsset(payload) {
  const res = await fetch(`${API_BASE}/api/v1/assets`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...authHeaders() },
    body: JSON.stringify(payload),
  })
  if (!res.ok) throw new Error('Create asset failed')
  return res.json()
}

export async function updateAsset(id, payload) {
  const res = await fetch(`${API_BASE}/api/v1/assets/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...authHeaders() },
    body: JSON.stringify(payload),
  })
  if (!res.ok) throw new Error('Update asset failed')
  return res.json()
}

export async function assetStats() {
  const res = await fetch(`${API_BASE}/api/v1/assets/stats`, { headers: { ...authHeaders() } })
  if (!res.ok) throw new Error('Fetch asset stats failed')
  return res.json()
}
