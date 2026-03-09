<script setup>
import { ref } from 'vue'
import { changePassword } from '../lib/api'

const form = ref({ current: '', next: '', confirm: '' })
const loading = ref(false)
const error = ref('')
const success = ref('')

async function submit() {
  error.value = ''
  success.value = ''
  if (!form.value.current || !form.value.next) {
    error.value = 'Current password and new password are required'
    return
  }
  if (form.value.next.length < 6) {
    error.value = 'New password must be at least 6 characters'
    return
  }
  if (form.value.next !== form.value.confirm) {
    error.value = 'Password confirmation does not match'
    return
  }
  loading.value = true
  try {
    await changePassword(form.value.current, form.value.next)
    success.value = 'Password updated successfully'
    form.value = { current: '', next: '', confirm: '' }
  } catch (e) {
    error.value = e.message || 'Failed to change password'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="page">
    <div class="card" style="max-width: 640px;">
      <h3 class="section-title">Security</h3>
      <p class="muted">Change your account password.</p>

      <form class="form-grid mt-16" @submit.prevent="submit">
        <label>
          Current Password
          <input v-model="form.current" type="password" required />
        </label>
        <label>
          New Password
          <input v-model="form.next" type="password" required minlength="6" />
        </label>
        <label>
          Confirm New Password
          <input v-model="form.confirm" type="password" required minlength="6" />
        </label>

        <div class="actions full">
          <button class="btn" :disabled="loading">{{ loading ? 'Updating...' : 'Update Password' }}</button>
        </div>

        <p v-if="success" class="ok full">{{ success }}</p>
        <p v-if="error" class="danger-text full">{{ error }}</p>
      </form>
    </div>
  </div>
</template>
