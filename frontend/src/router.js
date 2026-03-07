import { createRouter, createWebHistory } from 'vue-router'
import DashboardPage from './pages/DashboardPage.vue'
import TicketsPage from './pages/TicketsPage.vue'
import AssetsPage from './pages/AssetsPage.vue'
import OnboardingPage from './pages/OnboardingPage.vue'

const routes = [
  { path: '/', redirect: '/dashboard' },
  { path: '/dashboard', component: DashboardPage },
  { path: '/tickets', component: TicketsPage },
  { path: '/assets', component: AssetsPage },
  { path: '/onboarding', component: OnboardingPage },
]

export default createRouter({
  history: createWebHistory(),
  routes,
})
