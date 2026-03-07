<template>
  <div class="app-shell">
    <aside class="sidebar">
      <div class="brand">
        <div class="brand-icon">IT</div>
        <div>
          <h1>ITMS</h1>
          <p>Gateway System</p>
        </div>
      </div>

      <div class="user-mini">
        <div class="avatar">R</div>
        <div>
          <strong>Riski</strong>
          <p>Admin</p>
        </div>
      </div>

      <nav>
        <RouterLink to="/dashboard" class="nav-item">Dashboard <span>▾</span></RouterLink>
        <RouterLink to="/tickets" class="nav-item">Service Desk <span>▾</span></RouterLink>
        <RouterLink to="/assets" class="nav-item">Manajemen Aset <span>▾</span></RouterLink>
        <a class="nav-item ghost">Manajemen Pengguna & Akses <span>▾</span></a>
        <a class="nav-item ghost">Manajemen Cabang <span>▾</span></a>
        <a class="nav-item ghost">Manajemen Permission <span>▾</span></a>
        <a class="nav-item ghost">Pengumuman <span>▾</span></a>
        <RouterLink to="/onboarding" class="nav-item">Tenant Onboarding <span>▾</span></RouterLink>
      </nav>
    </aside>

    <main class="content">
      <header class="topbar">
        <div class="top-left">
          <button class="icon-btn">☰</button>
          <div class="searchbox">
            <span>🔎</span>
            <input placeholder="Search..." />
          </div>
        </div>

        <div class="top-right">
          <button class="icon-btn">☾</button>
          <button class="icon-btn">⚙️</button>
          <button class="icon-btn">🔔</button>
          <div class="profile-chip">
            <div class="avatar small">R</div>
            <span>Riski</span>
          </div>
        </div>
      </header>

      <section class="page-wrap">
        <RouterView />
      </section>
    </main>
  </div>
</template>
