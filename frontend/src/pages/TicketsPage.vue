<script setup>
import { ref, computed } from 'vue'

const activeTab = ref('list')
const searchQuery = ref('')
const filterStatus = ref('all')
const filterPriority = ref('all')

// Mock ticket data
const tickets = ref([
  {
    id: 'INC-2024-001',
    title: 'Email service down for marketing team',
    description: 'Multiple users reporting inability to send/receive emails since 9 AM',
    status: 'open',
    priority: 'high',
    category: 'Email',
    assignee: 'John Smith',
    requester: 'Sarah Johnson',
    created: '2024-03-08 09:15',
    updated: '2024-03-08 10:30',
    sla: { response: '1h', resolution: '4h', remaining: '2h 30m' }
  },
  {
    id: 'INC-2024-002',
    title: 'Slow network performance in Building B',
    description: 'Users experiencing intermittent connectivity and slow speeds',
    status: 'in-progress',
    priority: 'medium',
    category: 'Network',
    assignee: 'Mike Chen',
    requester: 'David Lee',
    created: '2024-03-08 10:00',
    updated: '2024-03-08 11:00',
    sla: { response: '2h', resolution: '8h', remaining: '6h' }
  },
  {
    id: 'REQ-2024-003',
    title: 'New laptop request for new hire',
    description: 'Standard developer setup needed for new team member starting Monday',
    status: 'pending',
    priority: 'low',
    category: 'Hardware',
    assignee: 'Unassigned',
    requester: 'HR Department',
    created: '2024-03-08 11:30',
    updated: '2024-03-08 11:30',
    sla: { response: '4h', resolution: '48h', remaining: '47h' }
  },
  {
    id: 'INC-2024-004',
    title: 'Cannot access shared drive',
    description: 'Finance team unable to access Q1 reports on shared drive',
    status: 'resolved',
    priority: 'high',
    category: 'File Server',
    assignee: 'Lisa Wang',
    requester: 'Finance Team',
    created: '2024-03-07 14:00',
    updated: '2024-03-07 16:45',
    sla: { response: '1h', resolution: '4h', remaining: 'Completed' }
  }
])

const filteredTickets = computed(() => {
  return tickets.value.filter(ticket => {
    const matchesSearch = ticket.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                         ticket.id.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesStatus = filterStatus.value === 'all' || ticket.status === filterStatus.value
    const matchesPriority = filterPriority.value === 'all' || ticket.priority === filterPriority.value
    return matchesSearch && matchesStatus && matchesPriority
  })
})

const statusColor = (status) => {
  const colors = {
    'open': 'red',
    'in-progress': 'yellow',
    'pending': 'blue',
    'resolved': 'green',
    'closed': 'gray'
  }
  return colors[status] || 'gray'
}

const priorityColor = (priority) => {
  const colors = {
    'high': 'red',
    'medium': 'yellow',
    'low': 'green'
  }
  return colors[priority] || 'gray'
}

const stats = computed(() => {
  return {
    total: tickets.value.length,
    open: tickets.value.filter(t => t.status === 'open').length,
    inProgress: tickets.value.filter(t => t.status === 'in-progress').length,
    resolved: tickets.value.filter(t => t.status === 'resolved').length
  }
})

// New ticket form
const showNewTicket = ref(false)
const newTicket = ref({
  title: '',
  description: '',
  category: '',
  priority: 'medium'
})

const createTicket = () => {
  const ticket = {
    id: `INC-2024-${String(tickets.value.length + 1).padStart(3, '0')}`,
    title: newTicket.value.title,
    description: newTicket.value.description,
    status: 'open',
    priority: newTicket.value.priority,
    category: newTicket.value.category,
    assignee: 'Unassigned',
    requester: 'Current User',
    created: new Date().toLocaleString('en-GB').slice(0, -3),
    updated: new Date().toLocaleString('en-GB').slice(0, -3),
    sla: { response: '2h', resolution: '8h', remaining: '8h' }
  }
  tickets.value.unshift(ticket)
  showNewTicket.value = false
  newTicket.value = { title: '', description: '', category: '', priority: 'medium' }
}
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