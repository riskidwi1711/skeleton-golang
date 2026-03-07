<template>
  <div>
    <h3>Ticketing</h3>
    <p class="muted">UI draft untuk modul helpdesk + SLA.</p>

    <div class="card mt-16">
      <table class="table">
        <thead>
          <tr>
            <th>Ticket</th>
            <th>Priority</th>
            <th>Status</th>
            <th>SLA</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>#INC-1021 Printer HQ down</td>
            <td>High</td>
            <td><span class="pill warning">In Progress</span></td>
            <td>2h left</td>
          </tr>
          <tr>
            <td>#REQ-880 New laptop request</td>
            <td>Medium</td>
            <td><span class="pill">Waiting Approval</span></td>
            <td>18h left</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
