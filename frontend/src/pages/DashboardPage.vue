<script setup>
import { computed, onMounted, ref } from 'vue'
import { login, listTenants } from '../lib/api'

const loading = ref(true)
const error = ref('')
const tenants = ref([])
const now = ref(new Date())

onMounted(async () => {
  try {
    const token = await login()
    tenants.value = await listTenants(token)
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
})

const dateText = computed(() =>
  now.value.toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
)
const timeText = computed(() =>
  now.value.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
)
</script>

<template>
  <div>
    <div class="hero-row">
      <div>
        <h3 class="page-title">Selamat Datang, Riski!</h3>
        <p class="muted">🗓 {{ dateText }} &nbsp;&nbsp;•&nbsp;&nbsp; 🕒 {{ timeText }}</p>
      </div>
      <div class="muted right-note">Admin<br />Head Office Jakarta</div>
    </div>

    <div v-if="loading" class="card mt-16">Loading dashboard...</div>
    <div v-else-if="error" class="card danger mt-16">{{ error }}</div>

    <template v-else>
      <div class="stats-grid mt-16">
        <div class="stat-card blue">
          <div>
            <p>Total Users</p>
            <h4>{{ Math.max(tenants.length * 2, 6) }}</h4>
            <small>↑ +86.7% this week</small>
          </div>
          <div class="stat-icon">👥</div>
        </div>
        <div class="stat-card green">
          <div>
            <p>Total Cabang</p>
            <h4>38</h4>
            <small>8 active</small>
          </div>
          <div class="stat-icon">🏢</div>
        </div>
        <div class="stat-card purple">
          <div>
            <p>Total Divisi</p>
            <h4>48</h4>
            <small>Across all branches</small>
          </div>
          <div class="stat-icon">👨‍💼</div>
        </div>
        <div class="stat-card orange">
          <div>
            <p>Total Roles</p>
            <h4>3</h4>
            <small>Permission groups</small>
          </div>
          <div class="stat-icon">🛡️</div>
        </div>
      </div>

      <div class="grid-2 mt-16">
        <div class="card chart-card">
          <div class="card-head">
            <h4>User Growth</h4>
            <span class="pill">Last 7 days ▾</span>
          </div>
          <div class="chart-grid">
            <div class="chart-line" />
          </div>
        </div>

        <div class="card">
          <h4>System Health</h4>
          <div class="progress-row"><span>Server</span><div><i style="width: 23%" /></div><em>23%</em></div>
          <div class="progress-row"><span>Database</span><div><i style="width: 45%" /></div><em>45%</em></div>
          <div class="progress-row"><span>Memory</span><div><i style="width: 67%" /></div><em>67%</em></div>
          <div class="progress-row"><span>Storage</span><div><i style="width: 58%" /></div><em>58%</em></div>
        </div>
      </div>
    </template>
  </div>
</template>
