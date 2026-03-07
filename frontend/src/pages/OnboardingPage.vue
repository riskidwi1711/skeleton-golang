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
    <h3>Tenant Onboarding</h3>
    <p class="muted">Form onboarding ke endpoint <code>/api/v1/tenants/onboard</code>.</p>

    <form class="card mt-16" @submit.prevent="submit">
      <label>
        Company Name
        <input v-model="form.company_name" required placeholder="Acme Corp" />
      </label>

      <label>
        Admin Email
        <input v-model="form.admin_email" type="email" required placeholder="owner@acme.com" />
      </label>

      <label>
        Plan
        <select v-model="form.plan">
          <option value="starter">starter</option>
          <option value="pro">pro</option>
          <option value="enterprise">enterprise</option>
        </select>
      </label>

      <button :disabled="loading" class="btn">{{ loading ? 'Creating...' : 'Create Tenant' }}</button>

      <p v-if="success" class="ok">{{ success }}</p>
      <p v-if="error" class="danger-text">{{ error }}</p>
    </form>
  </div>
</template>
