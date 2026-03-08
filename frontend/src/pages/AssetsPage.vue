<script setup>
import { ref, computed } from 'vue'

const activeTab = ref('hardware')
const searchQuery = ref('')
const filterStatus = ref('all')
const filterCategory = ref('all')
const filterDepartment = ref('all')

// Mock asset data
const assets = ref([
  {
    id: 'AST-2024-001',
    name: 'Dell Latitude 5520',
    category: 'Laptop',
    type: 'hardware',
    status: 'active',
    assignedTo: 'John Smith',
    department: 'Engineering',
    location: 'Building A, Floor 3',
    serialNumber: 'DL5520X4821',
    purchaseDate: '2023-06-15',
    warrantyExpiry: '2026-06-14',
    value: '$1,200'
  },
  {
    id: 'AST-2024-002',
    name: 'HP EliteDesk 800 G6',
    category: 'Desktop',
    type: 'hardware',
    status: 'active',
    assignedTo: 'Sarah Johnson',
    department: 'Finance',
    location: 'Building B, Floor 2',
    serialNumber: 'HP800G692323',
    purchaseDate: '2023-08-20',
    warrantyExpiry: '2026-08-19',
    value: '$900'
  },
  {
    id: 'AST-2024-003',
    name: 'Microsoft Office 365',
    category: 'Productivity',
    type: 'software',
    status: 'active',
    assignedTo: 'Multiple Users',
    department: 'Company-wide',
    location: 'Cloud',
    licenseKey: 'XXXX-XXXX-XXXX',
    purchaseDate: '2024-01-01',
    expiryDate: '2024-12-31',
    value: '$150/year'
  },
  {
    id: 'AST-2024-004',
    name: 'iPhone 13 Pro',
    category: 'Mobile',
    type: 'hardware',
    status: 'active',
    assignedTo: 'Mike Chen',
    department: 'Sales',
    location: 'Remote',
    serialNumber: 'IP13PRO7654',
    purchaseDate: '2023-11-10',
    warrantyExpiry: '2024-11-09',
    value: '$999'
  },
  {
    id: 'AST-2024-005',
    name: 'Adobe Creative Cloud',
    category: 'Design',
    type: 'software',
    status: 'active',
    assignedTo: 'Design Team',
    department: 'Marketing',
    location: 'Cloud',
    licenseKey: 'XXXX-XXXX-XXXX',
    purchaseDate: '2024-02-01',
    expiryDate: '2025-01-31',
    value: '$600/year'
  },
  {
    id: 'AST-2024-006',
    name: 'ThinkPad X1 Carbon',
    category: 'Laptop',
    type: 'hardware',
    status: 'maintenance',
    assignedTo: 'Unassigned',
    department: 'IT',
    location: 'IT Storage',
    serialNumber: 'TPX1C3892',
    purchaseDate: '2022-04-15',
    warrantyExpiry: '2025-04-14',
    value: '$1,500'
  }
])

const filteredAssets = computed(() => {
  return assets.value.filter(asset => {
    const matchesTab = activeTab.value === 'all' || asset.type === activeTab.value
    const matchesSearch = asset.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                         asset.id.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                         asset.serialNumber?.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesStatus = filterStatus.value === 'all' || asset.status === filterStatus.value
    const matchesCategory = filterCategory.value === 'all' || asset.category === filterCategory.value
    const matchesDepartment = filterDepartment.value === 'all' || asset.department === filterDepartment.value
    return matchesTab && matchesSearch && matchesStatus && matchesCategory && matchesDepartment
  })
})

const statusColor = (status) => {
  const colors = {
    'active': 'green',
    'inactive': 'gray',
    'maintenance': 'yellow',
    'retired': 'red'
  }
  return colors[status] || 'gray'
}

const stats = computed(() => {
  const hardwareCount = assets.value.filter(a => a.type === 'hardware').length
  const softwareCount = assets.value.filter(a => a.type === 'software').length
  const activeCount = assets.value.filter(a => a.status === 'active').length
  const totalValue = assets.value.reduce((sum, asset) => {
    const value = parseInt(asset.value.replace(/[^0-9]/g, '')) || 0
    return sum + value
  }, 0)

  return {
    total: assets.value.length,
    hardware: hardwareCount,
    software: softwareCount,
    active: activeCount,
    totalValue: `$${totalValue.toLocaleString()}`
  }
})

// New asset form
const showNewAsset = ref(false)
const newAsset = ref({
  name: '',
  category: '',
  type: 'hardware',
  serialNumber: '',
  department: '',
  location: '',
  value: ''
})

const createAsset = () => {
  const asset = {
    id: `AST-2024-${String(assets.value.length + 1).padStart(3, '0')}`,
    name: newAsset.value.name,
    category: newAsset.value.category,
    type: newAsset.value.type,
    status: 'active',
    assignedTo: 'Unassigned',
    department: newAsset.value.department,
    location: newAsset.value.location,
    serialNumber: newAsset.value.serialNumber,
    purchaseDate: new Date().toISOString().split('T')[0],
    warrantyExpiry: '',
    value: newAsset.value.value
  }
  assets.value.unshift(asset)
  showNewAsset.value = false
  newAsset.value = {
    name: '',
    category: '',
    type: 'hardware',
    serialNumber: '',
    department: '',
    location: '',
    value: ''
  }
}

const categories = computed(() => {
  const cats = new Set(assets.value.map(a => a.category))
  return Array.from(cats)
})

const departments = computed(() => {
  const depts = new Set(assets.value.map(a => a.department))
  return Array.from(depts)
})
</script>

<template>
  <div>
    <!-- Header -->
    <div class="page-row">
      <div>
        <h3 class="page-title">Asset Management</h3>
        <p class="muted">Track and manage IT hardware and software assets</p>
      </div>
      <div class="actions">
        <button class="btn secondary">Import CSV</button>
        <button class="btn" @click="showNewAsset = true">Add Asset</button>
      </div>
    </div>

    <!-- Stats cards -->
    <div class="stats-grid mt-16">
      <div class="stat-card">
        <div class="stat-value">{{ stats.total }}</div>
        <div class="stat-label">Total Assets</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ stats.hardware }}</div>
        <div class="stat-label">Hardware</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ stats.software }}</div>
        <div class="stat-label">Software</div>
      </div>
      <div class="stat-card">
        <div class="stat-value text-green">{{ stats.active }}</div>
        <div class="stat-label">Active</div>
      </div>
    </div>

    <!-- Tabs -->
    <div class="tabs mt-16">
      <button 
        :class="['tab', activeTab === 'all' ? 'active' : '']"
        @click="activeTab = 'all'"
      >
        All Assets
      </button>
      <button 
        :class="['tab', activeTab === 'hardware' ? 'active' : '']"
        @click="activeTab = 'hardware'"
      >
        Hardware
      </button>
      <button 
        :class="['tab', activeTab === 'software' ? 'active' : '']"
        @click="activeTab = 'software'"
      >
        Software
      </button>
    </div>

    <!-- Filters -->
    <div class="card mt-16">
      <div class="filters-grid">
        <input 
          v-model="searchQuery" 
          placeholder="Search by name, ID, or serial number..." 
        />
        <select v-model="filterStatus">
          <option value="all">All Status</option>
          <option value="active">Active</option>
          <option value="inactive">Inactive</option>
          <option value="maintenance">Maintenance</option>
          <option value="retired">Retired</option>
        </select>
        <select v-model="filterCategory">
          <option value="all">All Categories</option>
          <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
        </select>
        <select v-model="filterDepartment">
          <option value="all">All Departments</option>
          <option v-for="dept in departments" :key="dept" :value="dept">{{ dept }}</option>
        </select>
      </div>
    </div>

    <!-- Assets table -->
    <div class="card mt-16">
      <table class="table">
        <thead>
          <tr>
            <th>Asset ID</th>
            <th>Name</th>
            <th>Category</th>
            <th>Status</th>
            <th>Assigned To</th>
            <th>Department</th>
            <th>Location</th>
            <th>Value</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="asset in filteredAssets" :key="asset.id">
            <td class="font-mono">{{ asset.id }}</td>
            <td>
              <div class="asset-name">{{ asset.name }}</div>
              <div class="asset-serial" v-if="asset.serialNumber">{{ asset.serialNumber }}</div>
            </td>
            <td>{{ asset.category }}</td>
            <td>
              <span :class="`pill ${statusColor(asset.status)}`">
                {{ asset.status }}
              </span>
            </td>
            <td>{{ asset.assignedTo }}</td>
            <td>{{ asset.department }}</td>
            <td>{{ asset.location }}</td>
            <td>{{ asset.value }}</td>
            <td>
              <div class="table-actions">
                <button class="btn tiny secondary">View</button>
                <button class="btn tiny secondary">Edit</button>
                <button class="btn tiny danger-outline">Retire</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- New asset modal -->
    <div v-if="showNewAsset" class="modal-overlay" @click.self="showNewAsset = false">
      <div class="modal">
        <div class="modal-header">
          <h3>Add New Asset</h3>
          <button class="close-btn" @click="showNewAsset = false">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>Asset Name</label>
            <input v-model="newAsset.name" placeholder="e.g., Dell Latitude 5520">
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>Type</label>
              <select v-model="newAsset.type">
                <option value="hardware">Hardware</option>
                <option value="software">Software</option>
              </select>
            </div>
            <div class="form-group">
              <label>Category</label>
              <input v-model="newAsset.category" placeholder="e.g., Laptop, Desktop, Mobile">
            </div>
          </div>
          <div class="form-group" v-if="newAsset.type === 'hardware'">
            <label>Serial Number</label>
            <input v-model="newAsset.serialNumber" placeholder="e.g., SN123456789">
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>Department</label>
              <input v-model="newAsset.department" placeholder="e.g., Engineering, Finance">
            </div>
            <div class="form-group">
              <label>Location</label>
              <input v-model="newAsset.location" placeholder="e.g., Building A, Floor 2">
            </div>
          </div>
          <div class="form-group">
            <label>Value</label>
            <input v-model="newAsset.value" placeholder="e.g., $1,200">
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn secondary" @click="showNewAsset = false">Cancel</button>
          <button class="btn" @click="createAsset">Add Asset</button>
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

.text-green { color: #22c55e; }

.tabs {
  display: flex;
  gap: 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.tab {
  background: none;
  border: none;
  color: #888;
  padding: 0.75rem 1rem;
  cursor: pointer;
  font-size: 0.875rem;
  position: relative;
  transition: color 0.2s;
}

.tab:hover {
  color: #fff;
}

.tab.active {
  color: #fff;
}

.tab.active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 2px;
  background: #3b82f6;
}

.asset-name {
  font-weight: 500;
  color: #fff;
}

.asset-serial {
  font-size: 0.75rem;
  color: #888;
  margin-top: 0.25rem;
  font-family: monospace;
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
.form-group select {
  width: 100%;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  padding: 0.75rem;
  color: #fff;
  font-size: 0.875rem;
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

.danger-outline {
  background: transparent;
  border-color: #ef4444;
  color: #ef4444;
}

.danger-outline:hover {
  background: rgba(239, 68, 68, 0.1);
}
</style>