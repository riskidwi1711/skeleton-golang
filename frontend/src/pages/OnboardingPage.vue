<script setup>
import { ref } from 'vue'
import { onboardTenant } from '../lib/api'

const form = ref({
  company_name: '',
  admin_email: '',
  plan: 'starter',
})

const loading = ref(false)
const success = ref('')
const error = ref('')

async function submit() {
  loading.value = true
  success.value = ''
  error.value = ''

  try {
    const result = await onboardTenant(form.value)
    const tenant = result?.data?.tenant
    success.value = `Tenant created: ${tenant?.company_name} (${tenant?.tenant_id})`
    form.value = { company_name: '', admin_email: '', plan: 'starter' }
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div>
    <div class="page-row">
      <div>
        <h3 class="page-title">Tenant Onboarding</h3>
        <p class="muted">Buat company baru untuk model SaaS multi-tenant.</p>
      </div>
      <div class="actions">
        <button class="btn secondary">View Logs</button>
      </div>
    </div>

    <div class="grid-2 mt-16">
      <div class="card">
        <h4 class="section-title">Create New Tenant</h4>
        <form class="form-grid mt-16" @submit.prevent="submit">
          <label>
            Company Name
            <input v-model="form.company_name" required placeholder="PT Contoh Nusantara" />
          </label>

          <label>
            Admin Email
            <input v-model="form.admin_email" type="email" required placeholder="owner@contoh.com" />
          </label>

          <label>
            Plan
            <select v-model="form.plan">
              <option value="starter">starter</option>
              <option value="pro">pro</option>
              <option value="enterprise">enterprise</option>
            </select>
          </label>

          <div class="actions full">
            <button :disabled="loading" class="btn">{{ loading ? 'Creating...' : 'Create Tenant' }}</button>
          </div>

          <p v-if="success" class="ok full">{{ success }}</p>
          <p v-if="error" class="danger-text full">{{ error }}</p>
        </form>
      </div>

      <div class="card">
        <h4 class="section-title">Onboarding Checklist</h4>
        <ul class="check-list mt-16">
          <li><span>01</span> Company profile diisi lengkap</li>
          <li><span>02</span> Admin owner terverifikasi</li>
          <li><span>03</span> Branch default dibuat</li>
          <li><span>04</span> Role & permission seed aktif</li>
        </ul>
      </div>
    </div>
  </div>
</template>
