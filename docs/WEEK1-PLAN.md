# Week 1 Execution Plan

## Done in this scaffold
- [x] Monorepo baseline with service folders
- [x] Gateway service
- [x] Auth service (health + login mock)
- [x] Tenant service with `POST /api/v1/tenants/onboard`
- [x] Docker Compose with Traefik + Redis + NATS + Postgres + Mongo

## Next coding tasks
1. Implement real auth + user/tenant repository in auth-service (replace mock user)
2. Add RBAC authorization rules in gateway (role checks per endpoint)
3. Add NATS consumers for onboarding downstream actions (welcome email, provisioning)
4. Add integration tests for login -> onboard -> list flow
