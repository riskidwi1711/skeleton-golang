# ITMS SaaS Architecture (Phase 0)

## Topology
- FE (to be implemented)
- Traefik (edge routing)
- Gateway (Go)
- Auth Service (Go + PostgreSQL)
- Tenant Service (Go + PostgreSQL)
- Ticket/Asset Services (planned)
- MongoDB (activity, notification, audit/event)
- Redis (cache, rate limit, queue helper)
- NATS (event bus)

## Request Flow
1. FE -> Gateway (`/api/v1/...`)
2. Gateway injects request-id + proxies to service
3. Service handles business logic and emits events (next phase)

## Data Split (Hybrid)
- PostgreSQL: tenant, auth, roles/permissions, asset master
- MongoDB: ticket activity stream, notifications, audit/event logs
