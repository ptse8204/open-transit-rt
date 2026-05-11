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

COMPOSE_FILE="deploy/docker-compose.yml"

usage() {
  cat <<'EOF'
Usage:
  scripts/bootstrap-dev.sh
  scripts/bootstrap-dev.sh --check
  scripts/bootstrap-dev.sh --help

Development bootstrap starts the local Postgres/PostGIS dependency, applies
migrations, seeds demo records, and prints local development tokens and service
commands. It is local development setup only; it is not hosted service and
not production-readiness proof, not consumer-acceptance proof, not
agency-approval proof, and not compliance proof.

Preflight:
  --check verifies required local tools, required repository files, Docker
  daemon availability, and Docker Compose config without starting services or
  changing database state.

Common blockers:
  - Docker CLI missing: install Docker Desktop or a compatible Docker engine.
  - Docker daemon stopped: start Docker Desktop or the Docker daemon.
  - Go missing: install the Go version expected by go.mod.
  - Port 55432 occupied: stop the conflicting local database or change the
    compose port mapping intentionally in deploy/docker-compose.yml.
  - Database not ready: inspect make agency-app-logs or docker compose logs.
EOF
}

log() {
  printf '\n==> %s\n' "$1"
}

fail() {
  printf '\nERROR: %s\n' "$1" >&2
  printf '\nNext actions:\n' >&2
  printf '  scripts/bootstrap-dev.sh --check\n' >&2
  printf '  docker compose -f %s ps\n' "$COMPOSE_FILE" >&2
  printf '  docker compose -f %s logs postgres\n' "$COMPOSE_FILE" >&2
  printf '  docs/tutorials/local-quickstart.md\n' >&2
  exit 1
}

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    fail "Missing required tool: $1. Install it, then rerun scripts/bootstrap-dev.sh --check."
  fi
}

require_path() {
  path="$1"
  if [ ! -e "$path" ]; then
    fail "Required repository path is missing: $path"
  fi
}

port_hint() {
  port="$1"
  label="$2"
  if command -v lsof >/dev/null 2>&1 && lsof -nP -iTCP:"$port" -sTCP:LISTEN >/dev/null 2>&1; then
    printf 'NOTICE: %s port %s already has a listener. If bootstrap fails, stop the conflicting process or review %s.\n' "$label" "$port" "$COMPOSE_FILE"
  fi
}

check_docker() {
  need docker
  if ! docker info >/dev/null 2>&1; then
    fail "Docker is installed but the Docker daemon is not available. Start Docker Desktop or your Docker daemon, then rerun bootstrap."
  fi
  if ! docker compose -f "$COMPOSE_FILE" config >/dev/null; then
    fail "Docker Compose config did not render from $COMPOSE_FILE."
  fi
}

preflight() {
  log "Check required tools"
  need go
  check_docker
  for path in "$COMPOSE_FILE" scripts/seed-dev.sql "$MIGRATIONS_DIR" cmd/migrate .env.example; do
    require_path "$path"
  done
  port_hint 55432 "Postgres"
  printf 'Go: %s\n' "$(go version)"
  printf 'Docker: available\n'
  printf 'Compose file: %s renders\n' "$COMPOSE_FILE"
  printf 'Database URL: configured for local bootstrap target\n'
  printf 'Bootstrap preflight passed. Next action: run make dev or scripts/bootstrap-dev.sh.\n'
}

case "${1:-}" in
  "")
    ;;
  --check)
    preflight
    exit 0
    ;;
  --help|-h)
    usage
    exit 0
    ;;
  *)
    usage >&2
    fail "Unknown argument: $1"
    ;;
esac

preflight

log "Starting Postgres/PostGIS"
if ! docker compose -f "$COMPOSE_FILE" up -d postgres; then
  fail "Could not start Postgres/PostGIS. Docker may be unavailable, the PostGIS image may not be present, or host port 55432 may be in use."
fi

log "Waiting for database readiness"
attempt=0
until docker compose -f "$COMPOSE_FILE" exec -T postgres pg_isready -U postgres -d open_transit_rt >/dev/null 2>&1; do
  attempt=$((attempt + 1))
  if [ "$attempt" -ge 30 ]; then
    fail "Postgres did not become ready after 30 attempts."
  fi
  sleep 2
done
attempt=0
until DATABASE_URL="$DATABASE_URL" MIGRATIONS_DIR="$MIGRATIONS_DIR" go run ./cmd/migrate status >/dev/null 2>&1; do
  attempt=$((attempt + 1))
  if [ "$attempt" -ge 30 ]; then
    fail "Host database connection did not become ready after 30 attempts."
  fi
  sleep 2
done

log "Applying migrations"
if ! DATABASE_URL="$DATABASE_URL" MIGRATIONS_DIR="$MIGRATIONS_DIR" go run ./cmd/migrate up; then
  fail "Migration command failed. Check database readiness and migration output above."
fi

log "Seeding development agencies"
if ! docker compose -f "$COMPOSE_FILE" exec -T postgres psql -U postgres -d open_transit_rt < scripts/seed-dev.sql; then
  fail "Demo seed failed. Check Postgres logs and scripts/seed-dev.sql."
fi

ADMIN_TOKEN="$(go run ./cmd/admin-token -sub admin@example.com -agency-id demo-agency | sed -n 's/^token=//p')"
if [ -z "$ADMIN_TOKEN" ]; then
  fail "Could not generate a local admin token."
fi

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
  make run-feed-trip-updates      http://localhost:8084/healthz
  make run-feed-alerts            http://localhost:8085/healthz
  make run-gtfs-studio            http://localhost:8086/healthz

Public feed URLs:
  http://localhost:8081/public/gtfs/schedule.zip
  http://localhost:8081/public/feeds.json
  http://localhost:8083/public/gtfsrt/vehicle_positions.pb
  http://localhost:8084/public/gtfsrt/trip_updates.pb
  http://localhost:8085/public/gtfsrt/alerts.pb

Protected debug/admin routes require the admin token above, for example:
  curl -H "Authorization: Bearer \$ADMIN_TOKEN" http://localhost:8083/public/gtfsrt/vehicle_positions.json
  curl -H "Authorization: Bearer \$ADMIN_TOKEN" http://localhost:8086/admin/gtfs-studio

Pinned validators:
  make validators-install
  make validators-check

Fixtures:
  testdata/

Executable agency demo:
  make demo-agency-flow

Hardening smoke:
  make smoke
URLS
