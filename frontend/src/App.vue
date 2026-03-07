<script setup>
import { computed, onMounted, onBeforeUnmount, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { clearToken } from './lib/auth'

const route = useRoute()
const router = useRouter()
const isAuthPage = computed(() => ['/login', '/signup'].includes(route.path))

const userMenuOpen = ref(false)
const userMenuRef = ref(null)

function toggleUserMenu() {
  userMenuOpen.value = !userMenuOpen.value
}

function closeUserMenu() {
  userMenuOpen.value = false
}

function onDocClick(e) {
  if (!userMenuRef.value) return
  if (!userMenuRef.value.contains(e.target)) {
    closeUserMenu()
  }
}

function logout() {
  clearToken()
  closeUserMenu()
  router.push('/login')
}

onMounted(() => {
  document.addEventListener('click', onDocClick)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onDocClick)
})
</script>

<template>
  <RouterView v-if="isAuthPage" />

  <div v-else class="app-shell">
    <aside class="sidebar">
      <div class="brand">
        <div class="brand-icon">IT</div>
        <div>
          <h1>ITMS</h1>
          <p>Gateway System</p>
        </div>
      </div>

      <div class="user-mini">
        <div class="avatar">R</div>
        <div>
          <strong>Riski</strong>
          <p>Admin</p>
        </div>
      </div>

      <nav>
        <RouterLink to="/dashboard" class="nav-item">Dashboard <span>▾</span></RouterLink>
        <RouterLink to="/tickets" class="nav-item">Service Desk <span>▾</span></RouterLink>
        <RouterLink to="/assets" class="nav-item">Manajemen Aset <span>▾</span></RouterLink>
        <a class="nav-item ghost">Manajemen Pengguna & Akses <span>▾</span></a>
        <a class="nav-item ghost">Manajemen Cabang <span>▾</span></a>
        <a class="nav-item ghost">Manajemen Permission <span>▾</span></a>
        <a class="nav-item ghost">Pengumuman <span>▾</span></a>
        <RouterLink to="/onboarding" class="nav-item">Tenant Onboarding <span>▾</span></RouterLink>
      </nav>
    </aside>

    <main class="content">
      <header class="topbar">
        <div class="top-left">
          <button class="icon-btn" aria-label="Menu">
            <svg viewBox="0 0 24 24" class="icon-svg"><path d="M4 7h16M4 12h16M4 17h16"/></svg>
          </button>
          <div class="searchbox">
            <svg viewBox="0 0 24 24" class="icon-svg muted-icon"><circle cx="11" cy="11" r="7"/><path d="M20 20l-3.5-3.5"/></svg>
            <input placeholder="Search..." />
          </div>
        </div>

        <div class="top-right">
          <button class="icon-btn" aria-label="Theme">
            <svg viewBox="0 0 24 24" class="icon-svg"><path d="M21 12.8A9 9 0 1111.2 3a7 7 0 009.8 9.8z"/></svg>
          </button>
          <button class="icon-btn" aria-label="Settings">
            <svg viewBox="0 0 24 24" class="icon-svg"><path d="M12 8a4 4 0 100 8 4 4 0 000-8zm8 4a8.9 8.9 0 00-.1-1l2-1.5-2-3.4-2.4 1a8 8 0 00-1.7-1L15.5 2h-4l-.3 2.1a8 8 0 00-1.7 1l-2.4-1-2 3.4 2 1.5a8.9 8.9 0 000 2l-2 1.5 2 3.4 2.4-1a8 8 0 001.7 1l.3 2.1h4l.3-2.1a8 8 0 001.7-1l2.4 1 2-3.4-2-1.5c.1-.3.1-.7.1-1z"/></svg>
          </button>
          <button class="icon-btn" aria-label="Notifications">
            <svg viewBox="0 0 24 24" class="icon-svg"><path d="M18 16V11a6 6 0 10-12 0v5l-2 2h16l-2-2zm-8 4a2 2 0 004 0"/></svg>
          </button>

          <div class="user-menu-wrap" ref="userMenuRef">
            <button class="profile-chip" @click.stop="toggleUserMenu">
              <div class="avatar small">R</div>
              <span>Riski</span>
              <svg viewBox="0 0 24 24" class="icon-svg tiny-chevron"><path d="M6 9l6 6 6-6"/></svg>
            </button>

            <div class="user-dropdown" v-if="userMenuOpen">
              <button class="dropdown-item" @click="logout">
                <svg viewBox="0 0 24 24" class="icon-svg"><path d="M10 17l-5-5 5-5M5 12h14M14 7v-2h5v14h-5v-2"/></svg>
                Logout
              </button>
            </div>
          </div>
        </div>
      </header>

      <section class="page-wrap">
        <RouterView />
      </section>
    </main>
  </div>
</template>
