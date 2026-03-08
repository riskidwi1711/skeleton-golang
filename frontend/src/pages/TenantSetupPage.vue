<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { completeTenantSetup } from '../lib/api'
import { getUser } from '../lib/auth'

const router = useRouter()
const me = getUser() || {}

const form = ref({
  company_name: '',
  timezone: 'Asia/Jakarta',
  primary_branch: 'Head Office',
  primary_department: 'IT Operations',
  ticket_prefix: 'INC',
})

const loading = ref(false)
const error = ref('')

async function submit() {
  loading.value = true
  error.value = ''
  try {
    await completeTenantSetup(me.tenant_id, form.value)
    router.push('/dashboard')
  } catch (e) {
    error.value = e.message || 'Setup failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-shell">
    <div class="login-card" style="max-width: 760px;">
      <div class="login-head">
        <div class="brand-icon">IT</div>
        <div>
          <h2>Tenant Initial Setup</h2>
          <p>Complete this once to activate your workspace</p>
        </div>
      </div>

      <form @submit.prevent="submit" class="login-form">
        <label>
          Company Display Name
          <input v-model="form.company_name" required placeholder="PT Contoh Nusantara" />
        </label>

        <label>
          Timezone
          <select v-model="form.timezone">
            <option value="Asia/Jakarta">Asia/Jakarta</option>
            <option value="Asia/Shanghai">Asia/Shanghai</option>
            <option value="Asia/Singapore">Asia/Singapore</option>
            <option value="UTC">UTC</option>
          </select>
        </label>

        <label>
          Primary Branch
          <input v-model="form.primary_branch" required placeholder="Head Office" />
        </label>

        <label>
          Primary Department
          <input v-model="form.primary_department" required placeholder="IT Operations" />
        </label>

        <label>
          Ticket Prefix
          <input v-model="form.ticket_prefix" placeholder="INC" maxlength="6" />
        </label>

        <button class="btn" :disabled="loading">{{ loading ? 'Saving setup...' : 'Complete Setup' }}</button>
        <p v-if="error" class="danger-text">{{ error }}</p>
      </form>
    </div>
  </div>
</template>
