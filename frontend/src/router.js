import { createRouter, createWebHistory } from 'vue-router'
import DashboardPage from './pages/DashboardPage.vue'
import TicketsPage from './pages/TicketsPage.vue'
import AssetsPage from './pages/AssetsPage.vue'
import OnboardingPage from './pages/OnboardingPage.vue'
import LoginPage from './pages/LoginPage.vue'
import SignupPage from './pages/SignupPage.vue'
import { isAuthed } from './lib/auth'

const routes = [
  { path: '/', redirect: '/dashboard' },
  { path: '/login', component: LoginPage, meta: { public: true, title: 'Login' } },
  { path: '/signup', component: SignupPage, meta: { public: true, title: 'Sign Up' } },
  { path: '/dashboard', component: DashboardPage, meta: { title: 'Dashboard' } },
  { path: '/tickets', component: TicketsPage, meta: { title: 'Service Desk' } },
  { path: '/assets', component: AssetsPage, meta: { title: 'Manajemen Aset' } },
  { path: '/onboarding', component: OnboardingPage, meta: { title: 'Tenant Onboarding' } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  if (to.meta.public) return true
  if (!isAuthed()) return '/login'
  return true
})

export default router
