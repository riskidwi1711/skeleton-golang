#!/bin/bash

set -e

BASE_DIR="$(cd "$(dirname "$0")" && pwd)"
GO_BIN="/usr/local/go/bin/go"

echo "🚀 Starting ITMS Development Environment..."

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

if ! command -v pg_isready >/dev/null 2>&1; then
  PG_PORT=5433
elif ! pg_isready -h localhost -p 5432 >/dev/null 2>&1; then
  echo -e "${YELLOW}⚠️  PostgreSQL not detected on port 5432${NC}"
  echo "Using existing PostgreSQL on port 5433..."
  PG_PORT=5433
else
  PG_PORT=5432
fi

echo "Cleaning up old processes..."
pkill -f "go run.*auth-service" 2>/dev/null || true
pkill -f "go run.*tenant-service" 2>/dev/null || true
pkill -f "go run.*ticket-service" 2>/dev/null || true
pkill -f "go run.*asset-service" 2>/dev/null || true
pkill -f "go run.*gateway" 2>/dev/null || true
sleep 1

mkdir -p "$BASE_DIR/logs"

echo -e "${GREEN}Starting Auth Service...${NC}"
(
  cd "$BASE_DIR/services/auth-service"
  PORT=8081 JWT_SECRET=dev-secret-key-2024 "$GO_BIN" run cmd/main.go
) > "$BASE_DIR/logs/auth.log" 2>&1 &
AUTH_PID=$!

echo -e "${GREEN}Starting Tenant Service...${NC}"
(
  cd "$BASE_DIR/services/tenant-service"
  PORT=8082 DATABASE_URL="postgres://mailhop_user:Mailhop2024!@localhost:${PG_PORT}/mailhop?sslmode=disable" NATS_URL="" "$GO_BIN" run cmd/main.go
) > "$BASE_DIR/logs/tenant.log" 2>&1 &
TENANT_PID=$!

echo -e "${GREEN}Starting Ticket Service...${NC}"
(
  cd "$BASE_DIR/services/ticket-service"
  PORT=8083 DATABASE_URL="postgres://mailhop_user:Mailhop2024!@localhost:${PG_PORT}/mailhop?sslmode=disable" JWT_SECRET=dev-secret-key-2024 "$GO_BIN" run cmd/main.go
) > "$BASE_DIR/logs/ticket.log" 2>&1 &
TICKET_PID=$!

echo -e "${GREEN}Starting Asset Service...${NC}"
(
  cd "$BASE_DIR/services/asset-service"
  PORT=8084 DATABASE_URL="postgres://mailhop_user:Mailhop2024!@localhost:${PG_PORT}/mailhop?sslmode=disable" JWT_SECRET=dev-secret-key-2024 "$GO_BIN" run cmd/main.go
) > "$BASE_DIR/logs/asset.log" 2>&1 &
ASSET_PID=$!

echo -e "${GREEN}Starting API Gateway...${NC}"
(
  cd "$BASE_DIR/services/gateway"
  PORT=8080 AUTH_SERVICE_URL=http://localhost:8081 TENANT_SERVICE_URL=http://localhost:8082 TICKET_SERVICE_URL=http://localhost:8083 ASSET_SERVICE_URL=http://localhost:8084 JWT_SECRET=dev-secret-key-2024 "$GO_BIN" run cmd/main.go
) > "$BASE_DIR/logs/gateway.log" 2>&1 &
GATEWAY_PID=$!

echo "Waiting for services to start..."
sleep 4

echo -e "\n📊 Service Status:"
for svc in "Auth:$AUTH_PID" "Tenant:$TENANT_PID" "Ticket:$TICKET_PID" "Asset:$ASSET_PID" "Gateway:$GATEWAY_PID"; do
  NAME="${svc%%:*}"
  PID="${svc##*:}"
  if ps -p "$PID" >/dev/null 2>&1; then
    echo -e "  $NAME Service: ${GREEN}✓ Running${NC} (PID: $PID)"
  else
    echo -e "  $NAME Service: ${RED}✗ Failed${NC}"
  fi
done

echo -e "\n🔍 Testing API Gateway health..."
if curl -s http://localhost:8080/health >/dev/null 2>&1; then
  echo -e "  Gateway Health: ${GREEN}✓ OK${NC}"
else
  echo -e "  Gateway Health: ${RED}✗ Not responding${NC}"
fi

echo "$AUTH_PID $TENANT_PID $TICKET_PID $ASSET_PID $GATEWAY_PID" > "$BASE_DIR/.pids"

echo -e "\n✅ Development environment started!"
echo "  API Gateway: http://localhost:8080"
echo "  Frontend: http://localhost:5173"
echo "To stop all services: ./stop-dev.sh"