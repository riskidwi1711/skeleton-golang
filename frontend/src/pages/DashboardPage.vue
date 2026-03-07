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
    <h3>Dashboard Overview</h3>
    <p class="muted">Realtime ringkasan tenant onboarding dari API gateway.</p>

    <div v-if="loading" class="card">Loading dashboard...</div>
    <div v-else-if="error" class="card danger">{{ error }}</div>
    <div v-else class="grid">
      <div class="card">
        <p class="label">Total Tenants</p>
        <p class="value">{{ tenants.length }}</p>
      </div>
      <div class="card">
        <p class="label">Open Tickets</p>
        <p class="value">14</p>
      </div>
      <div class="card">
        <p class="label">Assets Tracked</p>
        <p class="value">289</p>
      </div>
    </div>

    <div class="card mt-16">
      <h4>Latest Tenants</h4>
      <ul class="tenant-list" v-if="tenants.length">
        <li v-for="tenant in tenants.slice(0, 6)" :key="tenant.tenant_id">
          <strong>{{ tenant.company_name }}</strong>
          <span>{{ tenant.admin_email }}</span>
        </li>
      </ul>
      <p v-else class="muted">Belum ada tenant.</p>
    </div>
  </div>
</template>
