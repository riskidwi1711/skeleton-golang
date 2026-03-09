<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { completeTenantSetup } from '../lib/api'
import { getUser, setTenantSetupCompleted } from '../lib/auth'

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
const toast = ref({ show: false, type: 'error', message: '' })
let toastTimer = null

function showToast(message, type = 'error') {
  toast.value = { show: true, type, message }
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    toast.value.show = false
  }, 3500)
}

async function submit() {
  loading.value = true
  try {
    await completeTenantSetup(me.tenant_id, form.value)
    setTenantSetupCompleted(true)
    showToast('Tenant setup completed. Redirecting...', 'success')
    setTimeout(() => router.push('/dashboard'), 500)
  } catch (e) {
    const msg = e.message || 'Setup failed'
    if (msg.toLowerCase().includes('tenant not found')) {
      showToast('Tenant tidak ditemukan. Silakan signup ulang atau login ulang.', 'error')
    } else {
      showToast(msg, 'error')
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-shell">
    <div v-if="toast.show" class="toast" :class="toast.type === 'success' ? 'ok' : 'err'">
      {{ toast.message }}
    </div>

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
      </form>
    </div>
  </div>
</template>

<style scoped>
.toast {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 2000;
  padding: 12px 16px;
  border-radius: 10px;
  color: #fff;
  font-weight: 600;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.18);
}

.toast.err {
  background: #d13438;
}

.toast.ok {
  background: #107c10;
}
</style>
