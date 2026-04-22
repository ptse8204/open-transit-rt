#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

if [ -f ".env" ]; then
  set -a
  # shellcheck disable=SC1091
  . ".env"
  set +a
fi

DATABASE_URL="${DATABASE_URL:-postgres://postgres:postgres@localhost:55432/open_transit_rt?sslmode=disable}"
MIGRATIONS_DIR="${MIGRATIONS_DIR:-db/migrations}"
ADMIN_JWT_SECRET="${ADMIN_JWT_SECRET:-dev-admin-jwt-secret-change-me}"
ADMIN_JWT_ISSUER="${ADMIN_JWT_ISSUER:-open-transit-rt-local}"
ADMIN_JWT_AUDIENCE="${ADMIN_JWT_AUDIENCE:-open-transit-rt-admin}"
ADMIN_JWT_TTL="${ADMIN_JWT_TTL:-8h}"
CSRF_SECRET="${CSRF_SECRET:-dev-csrf-secret-change-me}"
DEVICE_TOKEN_PEPPER="${DEVICE_TOKEN_PEPPER:-dev-device-token-pepper-change-me}"
export ADMIN_JWT_SECRET ADMIN_JWT_ISSUER ADMIN_JWT_AUDIENCE ADMIN_JWT_TTL CSRF_SECRET DEVICE_TOKEN_PEPPER

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing required tool: $1" >&2
    exit 1
  fi
}

need docker
need go

echo "Starting Postgres/PostGIS..."
docker compose -f deploy/docker-compose.yml up -d postgres

echo "Waiting for database readiness..."
attempt=0
until docker compose -f deploy/docker-compose.yml exec -T postgres pg_isready -U postgres -d open_transit_rt >/dev/null 2>&1; do
  attempt=$((attempt + 1))
  if [ "$attempt" -ge 30 ]; then
    echo "database did not become ready after 30 attempts" >&2
    exit 1
  fi
  sleep 2
done

echo "Applying migrations..."
DATABASE_URL="$DATABASE_URL" MIGRATIONS_DIR="$MIGRATIONS_DIR" go run ./cmd/migrate up

echo "Seeding development agencies..."
docker compose -f deploy/docker-compose.yml exec -T postgres psql -U postgres -d open_transit_rt < scripts/seed-dev.sql

ADMIN_TOKEN="$(go run ./cmd/admin-token -sub admin@example.com -agency-id demo-agency | sed -n 's/^token=//p')"

cat <<URLS

Open Transit RT local bootstrap complete.

Local admin API token:
  Authorization: Bearer $ADMIN_TOKEN

Local telemetry device credential:
  device_id=device-1 vehicle_id=bus-1 token=dev-device-token

Cookie auth is only for browser-admin flows. Bearer auth is the default for machine/API admin calls.

Core service commands:
  make run-agency-config          http://localhost:8081/healthz
  make run-telemetry-ingest       http://localhost:8082/healthz
  make run-feed-vehicle-positions http://localhost:8083/healthz

Public feed URLs planned by contract:
  http://localhost:8083/public/gtfsrt/vehicle_positions.pb
  http://localhost:8083/public/gtfsrt/vehicle_positions.json

Fixtures:
  testdata/

Hardening smoke:
  make smoke
URLS
