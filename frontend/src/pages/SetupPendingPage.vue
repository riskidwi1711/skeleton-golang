<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getTenantById } from '../lib/api'
import { clearToken, getUser, setTenantSetupCompleted } from '../lib/auth'

const router = useRouter()
const me = getUser() || {}
const loading = ref(true)
const error = ref('')
const tenant = ref(null)

async function checkStatus() {
  loading.value = true
  error.value = ''
  try {
    const t = await getTenantById(me.tenant_id)
    tenant.value = t
    if (t?.setup_completed_at) {
      setTenantSetupCompleted(true)
      router.push('/dashboard')
    }
  } catch (e) {
    error.value = e.message || 'Failed to check setup status'
  } finally {
    loading.value = false
  }
}

function logout() {
  clearToken()
  router.push('/login')
}

onMounted(() => {
  checkStatus()
  const interval = setInterval(checkStatus, 5000)
  return () => clearInterval(interval)
})
</script>

<template>
  <div class="login-shell">
    <div class="login-card">
      <div class="login-head">
        <div class="brand-icon">IT</div>
        <div>
          <h2>Workspace Setup Pending</h2>
          <p>Your workspace owner needs to complete initial setup</p>
        </div>
      </div>

      <div class="pending-content">
        <div v-if="loading" class="spinner-box">
          <div class="spinner"></div>
          <p>Checking setup status...</p>
        </div>

        <div v-else-if="error" class="error-box">
          <p class="danger-text">{{ error }}</p>
          <button class="btn secondary" @click="checkStatus">Retry</button>
        </div>

        <div v-else class="info-box">
          <div class="status-icon">⏳</div>
          <h3>Setup Not Completed</h3>
          <p class="muted">
            Your company workspace <strong>{{ tenant?.company_name || 'this tenant' }}</strong> 
            requires initial configuration by an owner or admin.
          </p>
          <p class="muted">
            Current status: <span class="badge">{{ tenant?.status || 'trialing' }}</span>
          </p>
          <p class="muted small">
            This page will automatically redirect once setup is complete.
            Checking every 5 seconds...
          </p>
          <div class="button-group">
            <button class="btn secondary" @click="checkStatus">Check Now</button>
            <button class="btn danger-outline" @click="logout">Switch Account</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.pending-content {
  padding: 2rem 0;
  text-align: center;
}

.spinner-box {
  padding: 2rem;
}

.spinner {
  width: 48px;
  height: 48px;
  border: 3px solid rgba(0, 0, 0, 0.1);
  border-top-color: #0f6cbd;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.info-box {
  padding: 1rem;
}

.status-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.info-box h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.info-box .muted {
  margin: 0.5rem 0;
  color: #6c757d;
}

.info-box .small {
  font-size: 0.875rem;
  font-style: italic;
}

.badge {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  background: #f8f9fa;
  border: 1px solid #dee2e6;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 600;
  color: #495057;
}

.error-box {
  padding: 2rem;
}

.btn.secondary {
  margin-top: 1rem;
}

.button-group {
  display: flex;
  gap: 0.5rem;
  justify-content: center;
  margin-top: 1.5rem;
}

.btn.danger-outline {
  background: transparent;
  border: 1px solid #dc3545;
  color: #dc3545;
}

.btn.danger-outline:hover {
  background: rgba(220, 53, 69, 0.1);
}
</style>