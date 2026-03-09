import { createRouter, createWebHistory } from 'vue-router'
import DashboardPage from './pages/DashboardPage.vue'
import TicketsPage from './pages/TicketsPage.vue'
import AssetsPage from './pages/AssetsPage.vue'
import AccessPage from './pages/AccessPage.vue'
import TenantSetupPage from './pages/TenantSetupPage.vue'
import SetupPendingPage from './pages/SetupPendingPage.vue'
import SecurityPage from './pages/SecurityPage.vue'
import LoginPage from './pages/LoginPage.vue'
import SignupPage from './pages/SignupPage.vue'
import { getTenantSetupCompleted, getUser, hasPermission, isAuthed, setTenantSetupCompleted } from './lib/auth'
import { getTenantById } from './lib/api'

const routes = [
  { path: '/', redirect: '/dashboard' },
  { path: '/login', component: LoginPage, meta: { public: true, title: 'Login' } },
  { path: '/signup', component: SignupPage, meta: { public: true, title: 'Sign Up' } },
  { path: '/dashboard', component: DashboardPage, meta: { title: 'Dashboard', permission: 'dashboard:read', requireSetup: true } },
  { path: '/tickets', component: TicketsPage, meta: { title: 'Service Desk', permission: 'tickets:read', requireSetup: true } },
  { path: '/assets', component: AssetsPage, meta: { title: 'Asset Management', permission: 'assets:read', requireSetup: true } },
  { path: '/access', component: AccessPage, meta: { title: 'Users & Access', permission: 'users:read', requireSetup: true } },
  { path: '/setup', component: TenantSetupPage, meta: { title: 'Tenant Setup' } },
  { path: '/setup-pending', component: SetupPendingPage, meta: { title: 'Setup Pending' } },
  { path: '/security', component: SecurityPage, meta: { title: 'Security', requireSetup: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to) => {
  if (to.meta.public) return true
  if (!isAuthed()) return '/login'

  if (to.meta.permission && !hasPermission(to.meta.permission)) return '/dashboard'

  const me = getUser()
  const isOwnerOrAdmin = hasPermission('users:write')

  if (me?.tenant_id && to.meta.requireSetup) {
    let completed = getTenantSetupCompleted()
    if (!completed) {
      try {
        const tenant = await getTenantById(me.tenant_id)
        completed = Boolean(tenant?.setup_completed_at)
        setTenantSetupCompleted(completed)
      } catch {
        completed = false
      }
    }

    if (!completed) {
      return isOwnerOrAdmin ? '/setup' : '/setup-pending'
    }
  }

  if (me?.tenant_id) {
    const completed = getTenantSetupCompleted()
    if (completed && (to.path === '/setup' || to.path === '/setup-pending')) {
      return '/dashboard'
    }
    if (!completed && to.path === '/setup' && !isOwnerOrAdmin) {
      return '/setup-pending'
    }
  }

  return true
})

export default router
