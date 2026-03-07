<template>
  <div>
    <div class="page-row">
      <div>
        <h3 class="page-title">Manajemen Tiket</h3>
        <p class="muted">Kelola semua tiket support layanan IT.</p>
      </div>
      <div class="actions">
        <button class="btn secondary">History Tiket</button>
        <button class="btn">+ Tiket Baru</button>
      </div>
    </div>

    <div class="card mt-16">
      <div class="filters-grid">
        <input placeholder="No. tiket / judul..." />
        <select>
          <option>Semua Status</option>
        </select>
        <select>
          <option>Semua Prioritas</option>
        </select>
        <select>
          <option>Semua Kategori</option>
        </select>
      </div>
    </div>

    <div class="card mt-16">
      <table class="table">
        <thead>
          <tr>
            <th>No. Tiket</th>
            <th>Judul</th>
            <th>Kategori</th>
            <th>Status</th>
            <th>Prioritas</th>
            <th>Dibuat</th>
            <th>Komentar</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>TKT-000007</td>
            <td>Jaringan Error</td>
            <td>Network</td>
            <td><span class="pill warning">Dalam Proses</span></td>
            <td><span class="pill yellow">Sedang</span></td>
            <td>08 Mar 2026</td>
            <td>3</td>
          </tr>
          <tr>
            <td>TKT-000008</td>
            <td>Email Gateway Delay</td>
            <td>Mail</td>
            <td><span class="pill">Open</span></td>
            <td><span class="pill red">Tinggi</span></td>
            <td>08 Mar 2026</td>
            <td>5</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
