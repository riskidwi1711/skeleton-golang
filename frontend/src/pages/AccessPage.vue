<script setup>
import { onMounted, ref } from 'vue'
import { createUser, listRoles, listUsers, updateUserRole } from '../lib/api'
import { getUser } from '../lib/auth'

const me = getUser()
const users = ref([])
const roles = ref([])
const loading = ref(false)
const error = ref('')

const form = ref({ name: '', email: '', role: 'agent' })

async function loadData() {
  loading.value = true
  error.value = ''
  try {
    const [u, r] = await Promise.all([listUsers(), listRoles()])
    users.value = u
    roles.value = r
    if (r.length && !form.value.role) form.value.role = r[0].name
  } catch (e) {
    error.value = e.message || 'Failed loading access data'
  } finally {
    loading.value = false
  }
}

async function submitCreate() {
  error.value = ''
  try {
    await createUser(form.value)
    form.value = { name: '', email: '', role: 'agent' }
    await loadData()
  } catch (e) {
    error.value = e.message || 'Create user failed'
  }
}

async function changeRole(userId, role) {
  error.value = ''
  try {
    await updateUserRole(userId, role)
    await loadData()
  } catch (e) {
    error.value = e.message || 'Update role failed'
  }
}

onMounted(loadData)
</script>

<template>
  <div class="page">
    <div class="card">
      <div class="card-head">
        <h3>Users & Access</h3>
        <p class="muted">Role-based access control per tenant.</p>
      </div>

      <p v-if="error" class="danger-text">{{ error }}</p>

      <div class="grid-2" style="margin-top:12px; gap:14px; align-items:flex-start;">
        <form class="card muted-box" @submit.prevent="submitCreate" v-if="me?.role === 'owner' || me?.role === 'admin'">
          <h4 style="margin:0 0 8px 0;">Invite User</h4>
          <label>
            Full Name
            <input v-model="form.name" required placeholder="Jane Doe" />
          </label>
          <label>
            Email
            <input v-model="form.email" required type="email" placeholder="jane@company.com" />
          </label>
          <label>
            Role
            <select v-model="form.role">
              <option v-for="r in roles" :key="r.name" :value="r.name">{{ r.label || r.name }}</option>
            </select>
          </label>
          <button class="btn" type="submit">Create User</button>
        </form>

        <div class="card muted-box">
          <h4 style="margin:0 0 8px 0;">Role Catalog</h4>
          <div v-for="r in roles" :key="r.name" class="list-row" style="padding:10px 0; border-bottom:1px solid #e9eef4;">
            <strong>{{ r.label || r.name }}</strong>
            <small class="muted">{{ (r.permissions || []).join(', ') }}</small>
          </div>
        </div>
      </div>

      <div class="table-wrap" style="margin-top:16px;">
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Email</th>
              <th>Role</th>
              <th style="width:180px;">Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="4">Loading...</td>
            </tr>
            <tr v-for="u in users" :key="u.id">
              <td>{{ u.name }}</td>
              <td>{{ u.email }}</td>
              <td>
                <select :value="u.role" @change="changeRole(u.id, $event.target.value)" :disabled="me?.role !== 'owner'">
                  <option v-for="r in roles" :key="r.name" :value="r.name">{{ r.label || r.name }}</option>
                </select>
              </td>
              <td>
                <span class="muted">{{ me?.role === 'owner' ? 'Owner can update roles' : 'View only' }}</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
