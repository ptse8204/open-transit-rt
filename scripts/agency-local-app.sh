#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

COMPOSE_FILE="deploy/docker-compose.yml"
PROFILE="app"
PUBLIC_ROOT="${PUBLIC_ROOT:-http://localhost:8080}"
FEED_ROOT="$PUBLIC_ROOT/public"
AGENCY_ID="${AGENCY_ID:-demo-agency}"
ADMIN_SUBJECT="${ADMIN_SUBJECT:-admin@example.com}"
SERVICES="agency-config telemetry-ingest feed-vehicle-positions feed-trip-updates feed-alerts gtfs-studio"
APP_SERVICES="$SERVICES local-proxy"

dc() {
  docker compose -f "$COMPOSE_FILE" "$@"
}

usage() {
  cat <<'EOF'
Usage:
  scripts/agency-local-app.sh up
  scripts/agency-local-app.sh down
  scripts/agency-local-app.sh logs
  scripts/agency-local-app.sh reset [--force]

Local app packaging starts the full Open Transit RT demo stack behind
http://localhost:8080. It is for local evaluation only, not production TLS or
admin network-boundary configuration.
EOF
}

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    fail "Missing required tool: $1

Next action:
  Install $1 and rerun: make agency-app-up"
  fi
}

log() {
  printf '\n==> %s\n' "$1"
}

fail() {
  printf '\nERROR: %s\n' "$1" >&2
  printf '\nRecovery:\n' >&2
  printf '  make agency-app-logs\n' >&2
  printf '  make agency-app-down\n' >&2
  printf '  make agency-app-reset   # destructive local reset with confirmation\n' >&2
  exit 1
}

check_docker() {
  need docker
  if ! docker info >/dev/null 2>&1; then
    fail "Docker is not available or the Docker daemon is not running.

Next action:
  Start Docker Desktop or your Docker daemon, then rerun:
    make agency-app-up"
  fi
}

wait_for_postgres() {
  attempt=0
  until dc exec -T postgres pg_isready -U postgres -d open_transit_rt >/dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [ "$attempt" -ge 45 ]; then
      fail "Postgres did not become ready in time."
    fi
    sleep 2
  done
}

wait_for_container_url() {
  service="$1"
  url="$2"
  label="$3"
  attempt=0
  until dc exec -T "$service" wget -qO- "$url" >/dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [ "$attempt" -ge 45 ]; then
      fail "$label did not become ready at $url."
    fi
    sleep 2
  done
}

wait_for_url() {
  url="$1"
  label="$2"
  attempt=0
  until curl -fsS "$url" >/dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [ "$attempt" -ge 45 ]; then
      fail "$label did not become ready at $url."
    fi
    sleep 2
  done
}

fetch_nonempty() {
  url="$1"
  label="$2"
  out="$(mktemp "${TMPDIR:-/tmp}/open-transit-local-fetch.XXXXXX")"
  if ! curl -fsS -o "$out" "$url"; then
    rm -f "$out"
    fail "Could not fetch $label from $url."
  fi
  if [ ! -s "$out" ]; then
    rm -f "$out"
    fail "$label fetched from $url was empty."
  fi
  rm -f "$out"
}

validator_status() {
  if ./scripts/check-validators.sh >/dev/null 2>&1; then
    printf 'ready: pinned validator tooling is installed'
  else
    printf 'not run: validator tooling is optional for startup; run make validators-install validators-check'
  fi
}

print_existing_state_note() {
  if docker volume inspect deploy_postgres-data >/dev/null 2>&1; then
    cat <<'EOF'

Local state note:
  A previous local demo database volume exists. This command will reuse it,
  apply migrations safely, seed demo records, and publish a fresh sample GTFS
  import as the active local feed.

  For a clean local database, run:
    make agency-app-reset
EOF
  fi
}

build_and_migrate() {
  log "Start Postgres/PostGIS"
  if ! dc up -d postgres; then
    fail "Could not start Postgres. Check Docker and port 55432 availability."
  fi
  wait_for_postgres

  log "Build local app image"
  if ! dc --profile "$PROFILE" build agency-config; then
    fail "Local app image build failed."
  fi

  log "Apply migrations"
  if ! dc --profile "$PROFILE" run --rm --no-deps agency-config /app/bin/migrate up; then
    fail "Migrations failed. Check database logs with make agency-app-logs."
  fi

  log "Seed local demo agency, roles, and device binding"
  if ! dc exec -T postgres psql -U postgres -d open_transit_rt < scripts/seed-dev.sql; then
    fail "Demo seed failed."
  fi
}

start_services() {
  log "Start local app services and local reverse proxy"
  if ! dc --profile "$PROFILE" up -d $APP_SERVICES; then
    fail "Could not start local app services. Port 8080 may already be in use by another process."
  fi

  log "Wait for service health checks"
  wait_for_container_url agency-config "http://127.0.0.1:8081/healthz" "agency-config health"
  wait_for_container_url telemetry-ingest "http://127.0.0.1:8082/healthz" "telemetry-ingest health"
  wait_for_container_url feed-vehicle-positions "http://127.0.0.1:8083/healthz" "Vehicle Positions health"
  wait_for_container_url feed-trip-updates "http://127.0.0.1:8084/healthz" "Trip Updates health"
  wait_for_container_url feed-alerts "http://127.0.0.1:8085/healthz" "Alerts health"
  wait_for_container_url gtfs-studio "http://127.0.0.1:8086/healthz" "GTFS Studio health"
  wait_for_url "$PUBLIC_ROOT/healthz" "local reverse proxy"
}

import_and_bootstrap() {
  log "Import sample GTFS and publish it as the active local feed"
  import_output="$(dc exec -T agency-config sh -lc 'rm -f /tmp/open-transit-valid-small.zip && cd /app/testdata/gtfs/valid-small && zip -qr /tmp/open-transit-valid-small.zip . && /app/bin/gtfs-import -agency-id demo-agency -zip /tmp/open-transit-valid-small.zip -actor-id agency-local-app -notes "Phase 16 local app startup"')"
  printf '%s\n' "$import_output" | grep -q '"status":"published"' || fail "Sample GTFS import did not publish successfully."

  log "Generate admin token in memory for local bootstrap"
  admin_token="$(dc exec -T agency-config /app/bin/admin-token -sub "$ADMIN_SUBJECT" -agency-id "$AGENCY_ID" | sed -n 's/^token=//p')"
  if [ -z "$admin_token" ]; then
    fail "Could not generate a local admin token for bootstrap."
  fi

  log "Bootstrap publication metadata"
  bootstrap_out="$(mktemp "${TMPDIR:-/tmp}/open-transit-bootstrap.XXXXXX")"
  if ! curl -fsS -X POST "$PUBLIC_ROOT/admin/publication/bootstrap" \
      -H "Authorization: Bearer $admin_token" \
      -H "Content-Type: application/json" \
      --data '{}' >"$bootstrap_out"; then
    rm -f "$bootstrap_out"
    fail "Publication metadata bootstrap failed."
  fi
  if ! grep -Eq '"stored"[[:space:]]*:[[:space:]]*true' "$bootstrap_out"; then
    rm -f "$bootstrap_out"
    fail "Publication metadata bootstrap did not report stored=true."
  fi
  rm -f "$bootstrap_out"
}

wait_for_readiness_and_feeds() {
  log "Wait for service readiness"
  wait_for_container_url agency-config "http://127.0.0.1:8081/readyz" "agency-config readiness"
  wait_for_container_url telemetry-ingest "http://127.0.0.1:8082/readyz" "telemetry-ingest readiness"
  wait_for_container_url feed-vehicle-positions "http://127.0.0.1:8083/readyz" "Vehicle Positions readiness"
  wait_for_container_url feed-trip-updates "http://127.0.0.1:8084/readyz" "Trip Updates readiness"
  wait_for_container_url feed-alerts "http://127.0.0.1:8085/readyz" "Alerts readiness"
  wait_for_container_url gtfs-studio "http://127.0.0.1:8086/readyz" "GTFS Studio readiness"

  log "Verify public feed discovery and feed URLs"
  fetch_nonempty "$FEED_ROOT/feeds.json" "feeds.json"
  fetch_nonempty "$FEED_ROOT/gtfs/schedule.zip" "schedule.zip"
  fetch_nonempty "$FEED_ROOT/gtfsrt/vehicle_positions.pb" "Vehicle Positions protobuf"
  fetch_nonempty "$FEED_ROOT/gtfsrt/trip_updates.pb" "Trip Updates protobuf"
  fetch_nonempty "$FEED_ROOT/gtfsrt/alerts.pb" "Alerts protobuf"
}

print_success() {
  generated_at="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
  validation="$(validator_status)"
  cat <<EOF

Open Transit RT local app is running.

Public feed root:
  $FEED_ROOT

Public feed URLs:
  Schedule ZIP:               $FEED_ROOT/gtfs/schedule.zip
  Feed discovery:             $FEED_ROOT/feeds.json
  Vehicle Positions protobuf: $FEED_ROOT/gtfsrt/vehicle_positions.pb
  Trip Updates protobuf:      $FEED_ROOT/gtfsrt/trip_updates.pb
  Alerts protobuf:            $FEED_ROOT/gtfsrt/alerts.pb

GTFS Studio/admin URL:
  $PUBLIC_ROOT/admin/gtfs-studio

Admin token instructions:
  Generate a local admin token only when you need one:
    docker compose -f deploy/docker-compose.yml --profile app exec -T agency-config /app/bin/admin-token -sub $ADMIN_SUBJECT -agency-id $AGENCY_ID

  Use the printed token as:
    Authorization: Bearer <token>

Demo device token instructions:
  The local seed includes a demo device binding for device-1 on bus-1.
  To avoid printing long-lived tokens here, use the device helper:
    scripts/device-onboarding.sh sample
    scripts/device-onboarding.sh simulate --dry-run
    scripts/device-onboarding.sh rebind --device-id device-1 --vehicle-id bus-1

Logs:
  make agency-app-logs

Validation:
  $validation

Exact next action:
  Open $FEED_ROOT/feeds.json, then run scripts/device-onboarding.sh sample to send a demo telemetry event.

Copy/paste support summary:
  generated_at=$generated_at
  feed_root=$FEED_ROOT
  admin_url=$PUBLIC_ROOT/admin/gtfs-studio
  app_profile=$PROFILE
  status=running
  log_location="make agency-app-logs"
  next_action="Open $FEED_ROOT/feeds.json, then run scripts/device-onboarding.sh sample"

Local scope:
  http://localhost:8080 is local-demo packaging only. Admin/debug routes may be proxied locally, but they still require auth. Production deployments need HTTPS/TLS and deployment-owned admin network controls.
EOF
}

cmd_up() {
  need curl
  check_docker
  print_existing_state_note
  build_and_migrate
  start_services
  import_and_bootstrap
  wait_for_readiness_and_feeds
  print_success
}

cmd_down() {
  check_docker
  log "Stop local app containers"
  dc --profile "$PROFILE" down
}

cmd_logs() {
  check_docker
  dc --profile "$PROFILE" logs --tail=160 $APP_SERVICES
}

cmd_reset() {
  check_docker
  force="false"
  if [ "${1:-}" = "--force" ]; then
    force="true"
  elif [ "${1:-}" != "" ]; then
    usage
    exit 2
  fi

  cat <<'EOF'
This is a destructive local reset.

It will remove:
  - containers: postgres, agency-config, telemetry-ingest, feed-vehicle-positions, feed-trip-updates, feed-alerts, gtfs-studio, local-proxy
  - volumes: deploy_postgres-data
  - generated local env files, if present: .env.local-app
  - local demo database state: all data stored in the Compose Postgres volume
  - logs, if applicable: Docker container logs attached to removed containers

It will not remove:
  - tracked repository files
  - .env
  - .cache
  - validator downloads
EOF

  if [ "$force" != "true" ]; then
    printf '\nType reset-local-app to continue: '
    read answer
    if [ "$answer" != "reset-local-app" ]; then
      echo "Reset canceled."
      exit 0
    fi
  fi

  log "Remove local app containers and volumes"
  dc --profile "$PROFILE" down -v --remove-orphans
  if [ -f ".env.local-app" ]; then
    rm -f ".env.local-app"
  fi
  echo "Local app reset complete."
}

case "${1:-}" in
  up) cmd_up ;;
  down) cmd_down ;;
  logs) cmd_logs ;;
  reset) shift; cmd_reset "$@" ;;
  -h|--help|help|"") usage ;;
  *) usage; exit 2 ;;
esac
