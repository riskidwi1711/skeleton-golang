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
  { path: '/login', component: LoginPage, meta: { public: true } },
  { path: '/signup', component: SignupPage, meta: { public: true } },
  { path: '/dashboard', component: DashboardPage },
  { path: '/tickets', component: TicketsPage },
  { path: '/assets', component: AssetsPage },
  { path: '/onboarding', component: OnboardingPage },
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
