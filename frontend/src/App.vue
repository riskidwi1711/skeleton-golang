<script setup>
import { computed, onMounted, onBeforeUnmount, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { clearToken } from './lib/auth'

const route = useRoute()
const router = useRouter()
const isAuthPage = computed(() => ['/login', '/signup'].includes(route.path))

const userMenuOpen = ref(false)
const notifOpen = ref(false)
const menuWrapRef = ref(null)
const notifWrapRef = ref(null)

const notifications = ref([
  { title: 'Ticket baru #INC-1029', desc: 'Priority tinggi - cabang Jakarta', time: '2m ago' },
  { title: 'Asset warranty akan habis', desc: '5 perangkat bulan ini', time: '15m ago' },
  { title: 'Tenant signup baru', desc: 'PT Maju Digital', time: '1h ago' },
])

const unreadCount = computed(() => notifications.value.length)
const pageTitle = computed(() => route.meta?.title || 'Dashboard')
const breadcrumb = computed(() => `ITMS / ${pageTitle.value}`)

function toggleUserMenu() {
  userMenuOpen.value = !userMenuOpen.value
  if (userMenuOpen.value) notifOpen.value = false
}

function toggleNotif() {
  notifOpen.value = !notifOpen.value
  if (notifOpen.value) userMenuOpen.value = false
}

function closeMenus() {
  userMenuOpen.value = false
  notifOpen.value = false
}

function onDocClick(e) {
  const insideUser = menuWrapRef.value?.contains(e.target)
  const insideNotif = notifWrapRef.value?.contains(e.target)
  if (!insideUser && !insideNotif) closeMenus()
}

function logout() {
  clearToken()
  closeMenus()
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

          <div>
            <p class="breadcrumb">{{ breadcrumb }}</p>
            <h4 class="page-caption">{{ pageTitle }}</h4>
          </div>

          <div class="searchbox">
            <svg viewBox="0 0 24 24" class="icon-svg muted-icon"><circle cx="11" cy="11" r="7"/><path d="M20 20l-3.5-3.5"/></svg>
            <input placeholder="Search menu, ticket, asset..." />
          </div>
        </div>

        <div class="top-right">
          <div class="command-bar">
            <button class="btn tiny">New</button>
            <button class="btn secondary tiny">Export</button>
            <button class="btn secondary tiny">Refresh</button>
          </div>

          <div class="notif-wrap" ref="notifWrapRef">
            <button class="icon-btn" aria-label="Notifications" @click.stop="toggleNotif">
              <svg viewBox="0 0 24 24" class="icon-svg"><path d="M18 16V11a6 6 0 10-12 0v5l-2 2h16l-2-2zm-8 4a2 2 0 004 0"/></svg>
              <span class="badge-dot" v-if="unreadCount">{{ unreadCount }}</span>
            </button>

            <div class="notif-dropdown" v-if="notifOpen">
              <p class="notif-title">Notifications</p>
              <div class="notif-item" v-for="n in notifications" :key="n.title + n.time">
                <strong>{{ n.title }}</strong>
                <span>{{ n.desc }}</span>
                <em>{{ n.time }}</em>
              </div>
            </div>
          </div>

          <div class="context-chip">Tenant: Demo • Role: Owner</div>

          <div class="user-menu-wrap" ref="menuWrapRef">
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
