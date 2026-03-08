<script setup>
import { ref, computed, onMounted } from 'vue'
import { listTickets, createTicket as apiCreateTicket, ticketStats } from '../lib/api'

const activeTab = ref('list')
const searchQuery = ref('')
const filterStatus = ref('all')
const filterPriority = ref('all')
const loading = ref(false)

const tickets = ref([])
const statsData = ref({ total: 0, open: 0, inProgress: 0, resolved: 0 })

const filteredTickets = computed(() => {
  return tickets.value.filter((ticket) => {
    const q = searchQuery.value.toLowerCase()
    const matchesSearch = ticket.title?.toLowerCase().includes(q) || ticket.id?.toLowerCase().includes(q)
    const matchesStatus = filterStatus.value === 'all' || ticket.status === filterStatus.value
    const matchesPriority = filterPriority.value === 'all' || ticket.priority === filterPriority.value
    return matchesSearch && matchesStatus && matchesPriority
  })
})

const statusColor = (status) => {
  const colors = { open: 'red', 'in-progress': 'yellow', pending: 'blue', resolved: 'green', closed: 'gray' }
  return colors[status] || 'gray'
}

const priorityColor = (priority) => {
  const colors = { high: 'red', medium: 'yellow', low: 'green' }
  return colors[priority] || 'gray'
}

const stats = computed(() => statsData.value)

const showNewTicket = ref(false)
const newTicket = ref({ title: '', description: '', category: '', priority: 'medium' })

function normalizeTicket(t) {
  return {
    ...t,
    assignee: t.assignee || 'Unassigned',
    requester: t.requester || 'Current User',
    created: t.created_at || t.created,
    updated: t.updated_at || t.updated,
    sla: {
      response: t.sla_response || '2h',
      resolution: t.sla_resolve || '8h',
      remaining: t.status === 'resolved' || t.status === 'closed' ? 'Completed' : t.sla_resolve || '8h',
    },
  }
}

async function loadTickets() {
  loading.value = true
  try {
    const data = await listTickets()
    tickets.value = (data || []).map(normalizeTicket)
    const st = await ticketStats()
    statsData.value = {
      total: st.total || 0,
      open: st.open || 0,
      inProgress: st.in_progress || 0,
      resolved: st.resolved || 0,
    }
  } catch {
    // keep UI usable even when backend is down
  } finally {
    loading.value = false
  }
}

async function createTicket() {
  if (!newTicket.value.title || !newTicket.value.category) return
  try {
    await apiCreateTicket(newTicket.value)
    showNewTicket.value = false
    newTicket.value = { title: '', description: '', category: '', priority: 'medium' }
    await loadTickets()
  } catch {
    // fallback local optimistic create
    tickets.value.unshift(
      normalizeTicket({
        id: `INC-2024-${String(tickets.value.length + 1).padStart(3, '0')}`,
        ...newTicket.value,
        status: 'open',
      })
    )
    showNewTicket.value = false
    newTicket.value = { title: '', description: '', category: '', priority: 'medium' }
  }
}

onMounted(loadTickets)
</script>

<template>
  <div>
    <!-- Header with stats -->
    <div class="page-row">
      <div>
        <h3 class="page-title">Service Desk</h3>
        <p class="muted">Manage support tickets and service requests</p>
      </div>
      <div class="actions">
        <button class="btn secondary">Export</button>
        <button class="btn" @click="showNewTicket = true">New Ticket</button>
      </div>
    </div>

    <!-- Stats cards -->
    <div class="stats-grid mt-16">
      <div class="stat-card">
        <div class="stat-value">{{ stats.total }}</div>
        <div class="stat-label">Total Tickets</div>
      </div>
      <div class="stat-card">
        <div class="stat-value text-red">{{ stats.open }}</div>
        <div class="stat-label">Open</div>
      </div>
      <div class="stat-card">
        <div class="stat-value text-yellow">{{ stats.inProgress }}</div>
        <div class="stat-label">In Progress</div>
      </div>
      <div class="stat-card">
        <div class="stat-value text-green">{{ stats.resolved }}</div>
        <div class="stat-label">Resolved</div>
      </div>
    </div>

    <!-- Filters -->
    <div class="card mt-16">
      <div class="filters-grid">
        <input 
          v-model="searchQuery" 
          placeholder="Search by ticket ID or title..." 
          class="search-input"
        />
        <select v-model="filterStatus">
          <option value="all">All Status</option>
          <option value="open">Open</option>
          <option value="in-progress">In Progress</option>
          <option value="pending">Pending</option>
          <option value="resolved">Resolved</option>
          <option value="closed">Closed</option>
        </select>
        <select v-model="filterPriority">
          <option value="all">All Priority</option>
          <option value="high">High</option>
          <option value="medium">Medium</option>
          <option value="low">Low</option>
        </select>
        <button class="btn secondary">Reset Filters</button>
      </div>
    </div>

    <!-- Tickets table -->
    <div class="card mt-16">
      <table class="table">
        <thead>
          <tr>
            <th>Ticket ID</th>
            <th>Title</th>
            <th>Status</th>
            <th>Priority</th>
            <th>Category</th>
            <th>Assignee</th>
            <th>SLA</th>
            <th>Updated</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="ticket in filteredTickets" :key="ticket.id">
            <td class="font-mono">{{ ticket.id }}</td>
            <td>
              <div class="ticket-title">{{ ticket.title }}</div>
              <div class="ticket-requester">{{ ticket.requester }}</div>
            </td>
            <td>
              <span :class="`pill ${statusColor(ticket.status)}`">
                {{ ticket.status.replace('-', ' ') }}
              </span>
            </td>
            <td>
              <span :class="`pill ${priorityColor(ticket.priority)}`">
                {{ ticket.priority }}
              </span>
            </td>
            <td>{{ ticket.category }}</td>
            <td>{{ ticket.assignee }}</td>
            <td>
              <div class="sla-info">{{ ticket.sla.remaining }}</div>
            </td>
            <td>{{ ticket.updated }}</td>
            <td>
              <div class="table-actions">
                <button class="btn tiny secondary">View</button>
                <button class="btn tiny secondary">Edit</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- New ticket modal -->
    <div v-if="showNewTicket" class="modal-overlay" @click.self="showNewTicket = false">
      <div class="modal">
        <div class="modal-header">
          <h3>Create New Ticket</h3>
          <button class="close-btn" @click="showNewTicket = false">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>Title</label>
            <input v-model="newTicket.title" placeholder="Brief description of the issue">
          </div>
          <div class="form-group">
            <label>Description</label>
            <textarea 
              v-model="newTicket.description" 
              rows="4" 
              placeholder="Detailed description of the issue or request"
            ></textarea>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>Category</label>
              <select v-model="newTicket.category">
                <option value="">Select category</option>
                <option>Hardware</option>
                <option>Software</option>
                <option>Network</option>
                <option>Email</option>
                <option>File Server</option>
                <option>Access Request</option>
                <option>Other</option>
              </select>
            </div>
            <div class="form-group">
              <label>Priority</label>
              <select v-model="newTicket.priority">
                <option value="high">High</option>
                <option value="medium">Medium</option>
                <option value="low">Low</option>
              </select>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn secondary" @click="showNewTicket = false">Cancel</button>
          <button class="btn" @click="createTicket">Create Ticket</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
}

.stat-card {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  padding: 1.5rem;
  text-align: center;
}

.stat-value {
  font-size: 2rem;
  font-weight: bold;
  color: #fff;
  margin-bottom: 0.5rem;
}

.stat-label {
  color: #888;
  font-size: 0.875rem;
}

.text-red { color: #ef4444; }
.text-yellow { color: #eab308; }
.text-green { color: #22c55e; }

.ticket-title {
  font-weight: 500;
  color: #fff;
}

.ticket-requester {
  font-size: 0.75rem;
  color: #888;
  margin-top: 0.25rem;
}

.sla-info {
  font-size: 0.875rem;
  color: #888;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: #1a1a1a;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-header h3 {
  margin: 0;
  color: #fff;
}

.close-btn {
  background: none;
  border: none;
  color: #888;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  color: #fff;
}

.modal-body {
  padding: 1.5rem;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #ccc;
  font-size: 0.875rem;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  padding: 0.75rem;
  color: #fff;
  font-size: 0.875rem;
}

.form-group textarea {
  resize: vertical;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.modal-footer {
  padding: 1.5rem;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}
</style>