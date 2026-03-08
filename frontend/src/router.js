import { createRouter, createWebHistory } from 'vue-router'
import DashboardPage from './pages/DashboardPage.vue'
import TicketsPage from './pages/TicketsPage.vue'
import AssetsPage from './pages/AssetsPage.vue'
import AccessPage from './pages/AccessPage.vue'
import OnboardingPage from './pages/OnboardingPage.vue'
import LoginPage from './pages/LoginPage.vue'
import SignupPage from './pages/SignupPage.vue'
import { hasPermission, isAuthed } from './lib/auth'

const routes = [
  { path: '/', redirect: '/dashboard' },
  { path: '/login', component: LoginPage, meta: { public: true, title: 'Login' } },
  { path: '/signup', component: SignupPage, meta: { public: true, title: 'Sign Up' } },
  { path: '/dashboard', component: DashboardPage, meta: { title: 'Dashboard', permission: 'dashboard:read' } },
  { path: '/tickets', component: TicketsPage, meta: { title: 'Service Desk', permission: 'tickets:read' } },
  { path: '/assets', component: AssetsPage, meta: { title: 'Asset Management', permission: 'assets:read' } },
  { path: '/access', component: AccessPage, meta: { title: 'Users & Access', permission: 'users:read' } },
  { path: '/onboarding', component: OnboardingPage, meta: { title: 'Tenant Onboarding', permission: 'tenants:read' } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  if (to.meta.public) return true
  if (!isAuthed()) return '/login'
  if (to.meta.permission && !hasPermission(to.meta.permission)) return '/dashboard'
  return true
})

export default router
