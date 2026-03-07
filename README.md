# ITMS SaaS (Microservice Foundation)

Foundation setup for a multi-tenant IT Management SaaS.

## Structure
- `services/gateway`
- `services/auth-service`
- `services/tenant-service`
- `frontend/`
- `deploy/docker-compose.yml`
- `docs/`

## Quick Start
```bash
cd deploy
docker compose up -d
```

## Gateway Endpoints
- `GET /health`
- `POST /api/v1/auth/login`
- `POST /api/v1/tenants/onboard` (public onboarding)
- `GET /api/v1/tenants` (JWT required)

## Auth + Tenant Flow Example
```bash
# 1) login to get token
TOKEN=$(curl -sS -X POST http://localhost/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"owner@acme.com"}' | jq -r '.token')

# 2) onboard tenant (public)
curl -sS -X POST http://localhost/api/v1/tenants/onboard \
  -H 'Content-Type: application/json' \
  -d '{"company_name":"Acme Corp","admin_email":"owner@acme.com","plan":"starter"}' | jq

# 3) list tenants (protected)
curl -sS http://localhost/api/v1/tenants \
  -H "Authorization: Bearer $TOKEN" | jq
```
