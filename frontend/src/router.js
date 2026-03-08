import { createRouter, createWebHistory } from 'vue-router'
import DashboardPage from './pages/DashboardPage.vue'
import TicketsPage from './pages/TicketsPage.vue'
import AssetsPage from './pages/AssetsPage.vue'
import AccessPage from './pages/AccessPage.vue'
import TenantSetupPage from './pages/TenantSetupPage.vue'
import LoginPage from './pages/LoginPage.vue'
import SignupPage from './pages/SignupPage.vue'
import { getTenantSetupCompleted, getUser, hasPermission, isAuthed, setTenantSetupCompleted } from './lib/auth'
import { getTenantById } from './lib/api'

const routes = [
  { path: '/', redirect: '/dashboard' },
  { path: '/login', component: LoginPage, meta: { public: true, title: 'Login' } },
  { path: '/signup', component: SignupPage, meta: { public: true, title: 'Sign Up' } },
  { path: '/dashboard', component: DashboardPage, meta: { title: 'Dashboard', permission: 'dashboard:read' } },
  { path: '/tickets', component: TicketsPage, meta: { title: 'Service Desk', permission: 'tickets:read' } },
  { path: '/assets', component: AssetsPage, meta: { title: 'Asset Management', permission: 'assets:read' } },
  { path: '/access', component: AccessPage, meta: { title: 'Users & Access', permission: 'users:read' } },
  { path: '/setup', component: TenantSetupPage, meta: { title: 'Tenant Setup' } },
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
  if (me?.tenant_id && hasPermission('users:write')) {
    let completed = getTenantSetupCompleted()
    if (!completed) {
      try {
        const tenant = await getTenantById(me.tenant_id)
        completed = Boolean(tenant?.setup_completed_at)
        setTenantSetupCompleted(completed)
      } catch {
        // if tenant lookup fails, do not break all routes
      }
    }

    if (!completed && to.path !== '/setup') return '/setup'
    if (completed && to.path === '/setup') return '/dashboard'
  }

  return true
})

export default router
