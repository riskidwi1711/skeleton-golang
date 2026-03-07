# Week 1 Execution Plan

## Done in this scaffold
- [x] Monorepo baseline with service folders
- [x] Gateway service
- [x] Auth service (health + login mock)
- [x] Tenant service with `POST /api/v1/tenants/onboard`
- [x] Docker Compose with Traefik + Redis + NATS + Postgres + Mongo

## Next coding tasks
1. Replace tenant in-memory store with PostgreSQL repository
2. Add JWT validation in gateway middleware
3. Add RBAC model in auth-service
4. Add event publish `tenant.onboarded` to NATS
