<script setup>
import { onMounted, ref } from 'vue'
import { login, listTenants } from '../lib/api'

const loading = ref(true)
const error = ref('')
const tenants = ref([])

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
</script>

<template>
  <div>
    <h3 class="page-title">Selamat Datang, Riski!</h3>
    <p class="muted">Dashboard IT Operations - ringkasan tenant, tiket, dan kesehatan sistem.</p>

    <div v-if="loading" class="card mt-16">Loading dashboard...</div>
    <div v-else-if="error" class="card danger mt-16">{{ error }}</div>

    <template v-else>
      <div class="stats-grid mt-16">
        <div class="stat-card blue">
          <p>Total Tenants</p>
          <h4>{{ tenants.length }}</h4>
        </div>
        <div class="stat-card green">
          <p>Total Branch</p>
          <h4>38</h4>
        </div>
        <div class="stat-card purple">
          <p>Total Divisi</p>
          <h4>48</h4>
        </div>
        <div class="stat-card orange">
          <p>Total Role</p>
          <h4>7</h4>
        </div>
      </div>

      <div class="grid-2 mt-16">
        <div class="card chart-card">
          <div class="card-head">
            <h4>User Growth</h4>
            <span class="pill">Last 7 days</span>
          </div>
          <div class="fake-chart">
            <div class="line" />
          </div>
        </div>

        <div class="card">
          <h4>System Health</h4>
          <div class="progress-row"><span>Server</span><div><i style="width: 42%" /></div></div>
          <div class="progress-row"><span>Database</span><div><i style="width: 63%" /></div></div>
          <div class="progress-row"><span>Memory</span><div><i style="width: 57%" /></div></div>
          <div class="progress-row"><span>Storage</span><div><i style="width: 34%" /></div></div>
        </div>
      </div>

      <div class="card mt-16">
        <div class="card-head">
          <h4>Latest Tenants</h4>
          <span class="muted">{{ tenants.length }} tenants</span>
        </div>
        <ul class="tenant-list" v-if="tenants.length">
          <li v-for="tenant in tenants.slice(0, 8)" :key="tenant.tenant_id">
            <strong>{{ tenant.company_name }}</strong>
            <span>{{ tenant.admin_email }}</span>
          </li>
        </ul>
        <p v-else class="muted">Belum ada tenant.</p>
      </div>
    </template>
  </div>
</template>
