.PHONY: up down logs test-onboard

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
