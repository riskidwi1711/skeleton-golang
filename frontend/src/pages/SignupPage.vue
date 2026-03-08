<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { registerTenant } from '../lib/api'
import { setSession } from '../lib/auth'

const router = useRouter()
const form = ref({
  company_name: '',
  full_name: '',
  email: '',
  password: '',
  plan: 'starter',
})

const loading = ref(false)
const error = ref('')

async function submit() {
  loading.value = true
  error.value = ''
  try {
    const result = await registerTenant(form.value)
    setSession(result)
    router.push('/setup')
  } catch (e) {
    error.value = e.message || 'Registration failed'
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
          <h2>Create Tenant Account</h2>
          <p>Daftar perusahaan kamu dan mulai trial</p>
        </div>
      </div>

      <form @submit.prevent="submit" class="login-form">
        <label>
          Company Name
          <input v-model="form.company_name" required placeholder="PT Contoh Nusantara" />
        </label>

        <label>
          Full Name
          <input v-model="form.full_name" required placeholder="Riski Dwi Patrio" />
        </label>

        <label>
          Email
          <input v-model="form.email" type="email" required placeholder="owner@contoh.com" />
        </label>

        <label>
          Password
          <input v-model="form.password" type="password" required placeholder="Minimum 8 karakter" />
        </label>

        <label>
          Plan
          <select v-model="form.plan">
            <option value="starter">starter</option>
            <option value="pro">pro</option>
            <option value="enterprise">enterprise</option>
          </select>
        </label>

        <button class="btn" :disabled="loading">{{ loading ? 'Creating account...' : 'Create Account' }}</button>
        <p v-if="error" class="danger-text">{{ error }}</p>

        <p class="muted" style="margin:0; font-size:12px;">
          Sudah punya akun?
          <RouterLink to="/login" style="color:#0f6cbd; font-weight:600;">Masuk di sini</RouterLink>
        </p>
      </form>
    </div>
  </div>
</template>
