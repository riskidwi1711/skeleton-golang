.PHONY: up down logs test-onboard test-flow

up:
	cd deploy && docker compose up -d

down:
	cd deploy && docker compose down

logs:
	cd deploy && docker compose logs -f --tail=200

test-onboard:
	curl -sS -X POST http://localhost/api/v1/tenants/onboard \
	-H 'Content-Type: application/json' \
	-d '{"company_name":"Acme Corp","admin_email":"owner@acme.com","plan":"starter"}' | jq

test-flow:
	TOKEN=$$(curl -sS -X POST http://localhost/api/v1/auth/login -H 'Content-Type: application/json' -d '{"email":"owner@acme.com"}' | jq -r '.token'); \
	echo "token: $$TOKEN"; \
	curl -sS http://localhost/api/v1/tenants -H "Authorization: Bearer $$TOKEN" | jq
