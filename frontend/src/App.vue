<template>
  <div class="app-shell">
    <aside class="sidebar">
      <div class="brand">
        <div class="dot" />
        <div>
          <h1>ITMS</h1>
          <p>Riski Labs</p>
        </div>
      </div>

      <nav>
        <RouterLink to="/dashboard" class="nav-item">Dashboard</RouterLink>
        <RouterLink to="/tickets" class="nav-item">Tickets</RouterLink>
        <RouterLink to="/assets" class="nav-item">Assets</RouterLink>
        <RouterLink to="/onboarding" class="nav-item">Tenant Onboarding</RouterLink>
      </nav>
    </aside>

    <main class="content">
      <header class="topbar">
        <h2>IT Management System</h2>
        <span class="badge">MVP UI</span>
      </header>

      <section class="page-wrap">
        <RouterView />
      </section>
    </main>
  </div>
</template>
