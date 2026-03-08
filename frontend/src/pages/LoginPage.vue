<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { login } from '../lib/api'
import { setSession } from '../lib/auth'

const router = useRouter()
const email = ref('owner@acme.com')
const password = ref('demo12345')
const loading = ref(false)
const error = ref('')

async function submit() {
  loading.value = true
  error.value = ''
  try {
    const session = await login(email.value, password.value)
    setSession(session)
    router.push('/dashboard')
  } catch (e) {
    error.value = e.message || 'Login failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-shell">
    <div class="login-card">
      <div class="login-head">
        <div class="brand-icon">IT</div>
        <div>
          <h2>ITMS Login</h2>
          <p>Sign in to access dashboard</p>
        </div>
      </div>

      <form @submit.prevent="submit" class="login-form">
        <label>
          Email
          <input v-model="email" type="email" placeholder="owner@acme.com" required />
        </label>

        <label>
          Password
          <input v-model="password" type="password" placeholder="Enter password" required />
        </label>

        <button class="btn" :disabled="loading">{{ loading ? 'Signing in...' : 'Sign In' }}</button>
        <p v-if="error" class="danger-text">{{ error }}</p>

        <p class="muted" style="margin:0; font-size:12px;">
          Belum punya tenant?
          <RouterLink to="/signup" style="color:#0f6cbd; font-weight:600;">Daftar sekarang</RouterLink>
        </p>
      </form>
    </div>
  </div>
</template>
