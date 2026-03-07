const API_BASE = import.meta.env.VITE_API_BASE || '/api'

export async function login(email = 'owner@acme.com') {
  const res = await fetch(`${API_BASE}/api/v1/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email }),
  })

  if (!res.ok) throw new Error('Login failed')
  const data = await res.json()
  return data.token
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
