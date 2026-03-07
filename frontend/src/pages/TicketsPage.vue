<script setup>
const services = [
  {
    title: 'Instalasi Windows',
    category: 'Sistem Operasi',
    description: 'Layanan instalasi sistem operasi Windows lengkap dengan driver & software standar.',
    price: 'Rp 500.000/jam',
    response: '2 jam',
    resolution: '4 jam',
    requests: '4 permintaan total (3 selesai)',
  },
  {
    title: 'Instalasi Linux',
    category: 'Sistem Operasi',
    description: 'Instalasi Linux (Ubuntu/CentOS) untuk server dan workstation dengan konfigurasi networking.',
    price: 'Rp 750.000/jam',
    response: '2 jam',
    resolution: '6 jam',
    requests: '0 permintaan total (0 selesai)',
  },
  {
    title: 'Troubleshooting Hardware',
    category: 'Hardware',
    description: 'Diagnosis dan perbaikan masalah hardware (RAM, storage, motherboard, power supply).',
    price: 'Rp 600.000/jam',
    response: '1 jam',
    resolution: '8 jam',
    requests: '0 permintaan total (0 selesai)',
  },
]
</script>

<template>
  <div>
    <div class="page-row">
      <div>
        <h3 class="page-title">Katalog Layanan</h3>
        <p class="muted">Kelola layanan yang tersedia dalam sistem.</p>
      </div>
      <div class="actions">
        <button class="btn success">Daftar Permintaan</button>
        <button class="btn purple">Buat Permintaan</button>
        <button class="btn">+ Tambah Layanan</button>
      </div>
    </div>

    <div class="card mt-16">
      <div class="filters-grid service-filter">
        <label>
          Cari Layanan
          <input placeholder="Nama atau deskripsi..." />
        </label>
        <label>
          Kategori
          <select>
            <option>Semua Kategori</option>
          </select>
        </label>
        <label>
          Status
          <select>
            <option>Semua Status</option>
          </select>
        </label>
        <div class="filter-buttons">
          <button class="btn">Cari</button>
          <button class="btn secondary">Reset</button>
        </div>
      </div>
    </div>

    <div class="service-grid mt-16">
      <article class="service-card" v-for="item in services" :key="item.title">
        <div class="service-head">
          <h4>{{ item.title }}</h4>
          <span class="service-cat">🏷 {{ item.category }}</span>
        </div>

        <div class="service-body">
          <p>{{ item.description }}</p>

          <div class="service-meta">
            <div><span>💵 Harga</span><strong>{{ item.price }}</strong></div>
            <div><span>⏱ SLA Response</span><strong>{{ item.response }}</strong></div>
            <div><span>🧪 SLA Resolusi</span><strong>{{ item.resolution }}</strong></div>
            <div><span>🟢 Status</span><strong class="ok">Aktif</strong></div>
          </div>

          <p class="service-req">📋 {{ item.requests }}</p>
        </div>

        <div class="service-actions">
          <button class="btn tiny">Lihat</button>
          <button class="btn tiny success">Request</button>
          <button class="btn tiny warning">Edit</button>
          <button class="btn tiny danger-solid">Hapus</button>
        </div>
      </article>
    </div>
  </div>
</template>
