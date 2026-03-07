<template>
  <div>
    <div class="page-row">
      <div>
        <h3 class="page-title">Inventaris Hardware</h3>
        <p class="muted">Kelola perangkat hardware perusahaan Anda.</p>
      </div>
      <div class="actions">
        <button class="btn success">Export Excel</button>
        <button class="btn warning">Import</button>
        <button class="btn">+ Tambah Hardware</button>
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
          <tr>
            <td>HRW-01-01</td>
            <td>Laptop Lenovo</td>
            <td>Laptop</td>
            <td><span class="pill green">Active</span></td>
            <td>482149</td>
            <td>Infrastruktur IT</td>
            <td>👁️ ✏️ 🗑️</td>
          </tr>
          <tr>
            <td>HRW-01-02</td>
            <td>PC Asus ROG</td>
            <td>PC</td>
            <td><span class="pill green">Active</span></td>
            <td>92323</td>
            <td>Manajemen Aset</td>
            <td>👁️ ✏️ 🗑️</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
