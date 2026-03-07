<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { login } from '../lib/api'
import { setToken } from '../lib/auth'

const router = useRouter()
const email = ref('owner@acme.com')
const loading = ref(false)
const error = ref('')

async function submit() {
  loading.value = true
  error.value = ''
  try {
    const token = await login(email.value)
    setToken(token)
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

        <button class="btn" :disabled="loading">{{ loading ? 'Signing in...' : 'Sign In' }}</button>
        <p v-if="error" class="danger-text">{{ error }}</p>
      </form>
    </div>
  </div>
</template>
