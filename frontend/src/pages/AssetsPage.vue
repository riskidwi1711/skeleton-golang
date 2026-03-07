<script setup>
const rows = [
  { code: 'HRW-01-01', name: 'Laptop Lenovo', category: 'Laptop', status: 'Active', serial: '482149', division: 'Infrastruktur IT' },
  { code: 'HRW-01-02', name: 'PC Asus ROG', category: 'Desktop', status: 'Active', serial: '92323', division: 'Manajemen Aset' },
]
</script>

<template>
  <div>
    <div class="page-row">
      <div>
        <h3 class="page-title">Inventaris Hardware</h3>
        <p class="muted">Kelola perangkat hardware perusahaan Anda secara terstruktur.</p>
      </div>
      <div class="actions">
        <button class="btn secondary">Import CSV</button>
        <button class="btn">Tambah Hardware</button>
      </div>
    </div>

    <div class="card mt-16">
      <div class="filters-grid">
        <input placeholder="Kode, nama, serial, IP..." />
        <select><option>Semua Status</option></select>
        <select><option>Semua Kategori</option></select>
        <select><option>Semua Divisi</option></select>
      </div>
    </div>

    <div class="card mt-16">
      <table class="table">
        <thead>
          <tr>
            <th>Kode</th>
            <th>Nama Hardware</th>
            <th>Kategori</th>
            <th>Status</th>
            <th>Serial Number</th>
            <th>Divisi</th>
            <th>Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rows" :key="row.code">
            <td>{{ row.code }}</td>
            <td>{{ row.name }}</td>
            <td>{{ row.category }}</td>
            <td><span class="pill green">{{ row.status }}</span></td>
            <td>{{ row.serial }}</td>
            <td>{{ row.division }}</td>
            <td>
              <div class="table-actions">
                <button class="btn tiny secondary">Detail</button>
                <button class="btn tiny secondary">Edit</button>
                <button class="btn tiny danger-solid">Hapus</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
