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
- `POST /api/v1/tenants/onboard`
- `GET /api/v1/tenants`

## Onboard Example
```bash
curl -X POST http://localhost/api/v1/tenants/onboard \
  -H 'Content-Type: application/json' \
  -d '{"company_name":"Acme Corp","admin_email":"owner@acme.com","plan":"starter"}'
```
