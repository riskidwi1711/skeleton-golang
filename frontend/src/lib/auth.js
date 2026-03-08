const TOKEN_KEY = 'itms_token'
const USER_KEY = 'itms_user'
const TENANT_SETUP_KEY = 'itms_tenant_setup_completed'

function parseJwt(token) {
  if (!token || token.split('.').length < 2) return null
  try {
    const payload = token.split('.')[1]
    const base64 = payload.replace(/-/g, '+').replace(/_/g, '/')
    return JSON.parse(atob(base64))
  } catch {
    return null
  }
}

export function getToken() {
  return localStorage.getItem(TOKEN_KEY) || ''
}

export function setToken(token) {
  localStorage.setItem(TOKEN_KEY, token)
}

export function setSession({ token, user }) {
  if (!token) return
  setToken(token)
  const claims = parseJwt(token) || {}
  const sessionUser = {
    id: user?.id || claims.sub || 'u-1',
    name: user?.name || user?.email?.split('@')?.[0] || claims.email?.split('@')?.[0] || 'User',
    email: user?.email || claims.email || '',
    role: user?.role || claims.role || 'viewer',
    permissions: user?.permissions || claims.permissions || [],
    tenant_id: user?.tenant_id || claims.tenant_id || 'tnt_demo',
  }
  localStorage.setItem(USER_KEY, JSON.stringify(sessionUser))
  const completed = Boolean(user?.setup_completed_at)
  localStorage.setItem(TENANT_SETUP_KEY, completed ? '1' : '0')
}

export function getUser() {
  try {
    const raw = localStorage.getItem(USER_KEY)
    if (raw) return JSON.parse(raw)
  } catch {
    // ignore parse error
  }
  const token = getToken()
  const claims = parseJwt(token)
  if (!claims) return null
  return {
    id: claims.sub,
    name: claims.email?.split('@')?.[0] || 'User',
    email: claims.email,
    role: claims.role || 'viewer',
    permissions: claims.permissions || [],
    tenant_id: claims.tenant_id || 'tnt_demo',
  }
}

export function hasPermission(permission) {
  const user = getUser()
  if (!user) return false
  return (user.permissions || []).includes(permission)
}

export function getTenantSetupCompleted() {
  return localStorage.getItem(TENANT_SETUP_KEY) === '1'
}

export function setTenantSetupCompleted(v) {
  localStorage.setItem(TENANT_SETUP_KEY, v ? '1' : '0')
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
  localStorage.removeItem(TENANT_SETUP_KEY)
}

export function isAuthed() {
  return !!getToken()
}
