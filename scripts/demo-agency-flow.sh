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
AGENCY_ID="${AGENCY_ID:-demo-agency}"
ADMIN_JWT_SECRET="${ADMIN_JWT_SECRET:-dev-admin-jwt-secret-change-me}"
ADMIN_JWT_ISSUER="${ADMIN_JWT_ISSUER:-open-transit-rt-local}"
ADMIN_JWT_AUDIENCE="${ADMIN_JWT_AUDIENCE:-open-transit-rt-admin}"
ADMIN_JWT_TTL="${ADMIN_JWT_TTL:-8h}"
CSRF_SECRET="${CSRF_SECRET:-dev-csrf-secret-change-me}"
DEVICE_TOKEN_PEPPER="${DEVICE_TOKEN_PEPPER:-dev-device-token-pepper-change-me}"
PUBLIC_PROXY_PORT="${PUBLIC_PROXY_PORT:-8090}"
PUBLIC_BASE_URL="${PUBLIC_BASE_URL:-http://localhost:${PUBLIC_PROXY_PORT}}"
FEED_BASE_URL="${FEED_BASE_URL:-http://localhost:${PUBLIC_PROXY_PORT}/public}"
REALTIME_VALIDATION_BASE_URL="${REALTIME_VALIDATION_BASE_URL:-$FEED_BASE_URL}"
TECHNICAL_CONTACT_EMAIL="${TECHNICAL_CONTACT_EMAIL:-dev@example.com}"
FEED_LICENSE_NAME="${FEED_LICENSE_NAME:-CC BY 4.0}"
FEED_LICENSE_URL="${FEED_LICENSE_URL:-https://creativecommons.org/licenses/by/4.0/}"
PUBLICATION_ENVIRONMENT="${PUBLICATION_ENVIRONMENT:-dev}"
VALIDATOR_TOOLING_MODE="${VALIDATOR_TOOLING_MODE:-pinned}"
GTFS_VALIDATOR_PATH="${GTFS_VALIDATOR_PATH:-$ROOT_DIR/.cache/validators/gtfs-validator-7.1.0-cli.jar}"
GTFS_VALIDATOR_VERSION="${GTFS_VALIDATOR_VERSION:-v7.1.0}"
GTFS_RT_VALIDATOR_PATH="${GTFS_RT_VALIDATOR_PATH:-$ROOT_DIR/.cache/validators/gtfs-rt-validator-wrapper.sh}"
GTFS_RT_VALIDATOR_VERSION="${GTFS_RT_VALIDATOR_VERSION:-ghcr.io/mobilitydata/gtfs-realtime-validator@sha256:5d2a3c14fba49983e1968c4a715e8ca624d4062bf4afede74aeca26322436c89}"
if [ -z "${JAVA_BINARY:-}" ]; then
  for candidate in \
    /usr/local/opt/openjdk@17/bin/java \
    /opt/homebrew/opt/openjdk@17/bin/java \
    /usr/local/opt/openjdk/bin/java \
    /opt/homebrew/opt/openjdk/bin/java
  do
    if [ -x "$candidate" ]; then
      JAVA_BINARY="$candidate"
      break
    fi
  done
fi

export DATABASE_URL MIGRATIONS_DIR AGENCY_ID
export ADMIN_JWT_SECRET ADMIN_JWT_ISSUER ADMIN_JWT_AUDIENCE ADMIN_JWT_TTL CSRF_SECRET DEVICE_TOKEN_PEPPER
export PUBLIC_BASE_URL FEED_BASE_URL REALTIME_VALIDATION_BASE_URL TECHNICAL_CONTACT_EMAIL FEED_LICENSE_NAME FEED_LICENSE_URL PUBLICATION_ENVIRONMENT
export VALIDATOR_TOOLING_MODE GTFS_VALIDATOR_PATH GTFS_VALIDATOR_VERSION GTFS_RT_VALIDATOR_PATH GTFS_RT_VALIDATOR_VERSION JAVA_BINARY
export APP_ENV="${APP_ENV:-dev}"
export BIND_ADDR="${BIND_ADDR:-127.0.0.1}"
export METRICS_ENABLED="${METRICS_ENABLED:-false}"
export PUBLIC_PROXY_PORT

TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/open-transit-rt-demo.XXXXXX")"
PIDS=""

cleanup() {
  for pid in $PIDS; do
    if kill "$pid" >/dev/null 2>&1; then
      wait "$pid" >/dev/null 2>&1 || true
    fi
  done
}
trap cleanup EXIT INT TERM

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing required tool: $1" >&2
    exit 1
  fi
}

log() {
  printf '\n==> %s\n' "$1"
}

wait_for_database() {
  attempt=0
  until docker compose -f deploy/docker-compose.yml exec -T postgres pg_isready -U postgres -d open_transit_rt >/dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [ "$attempt" -ge 30 ]; then
      echo "database did not become ready after 30 attempts" >&2
      exit 1
    fi
    sleep 2
  done
  attempt=0
  until DATABASE_URL="$DATABASE_URL" MIGRATIONS_DIR="$MIGRATIONS_DIR" go run ./cmd/migrate status >/dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [ "$attempt" -ge 30 ]; then
      echo "database host connection did not become ready after 30 attempts" >&2
      exit 1
    fi
    sleep 2
  done
}

wait_for_status() {
  url="$1"
  want="$2"
  attempt=0
  while :; do
    status="$(curl -s -o /dev/null -w '%{http_code}' "$url" || true)"
    if [ "$status" = "$want" ]; then
      return 0
    fi
    attempt=$((attempt + 1))
    if [ "$attempt" -ge 40 ]; then
      echo "timed out waiting for $url to return $want; last status was $status" >&2
      exit 1
    fi
    sleep 1
  done
}

expect_status() {
  want="$1"
  url="$2"
  status="$(curl -s -o /dev/null -w '%{http_code}' "$url" || true)"
  if [ "$status" != "$want" ]; then
    echo "expected $url to return $want, got $status" >&2
    exit 1
  fi
  echo "verified $url -> HTTP $status"
}

start_service() {
  name="$1"
  port="$2"
  shift 2
  log_file="$TMP_DIR/$name.log"
  (
    PORT="$port" "$@"
  ) >"$log_file" 2>&1 &
  pid="$!"
  PIDS="$PIDS $pid"
  wait_for_status "http://localhost:$port/healthz" "200"
  echo "$name listening on http://localhost:$port (log: $log_file)"
}

fetch_nonempty() {
  url="$1"
  out="$2"
  curl -fsS -o "$out" "$url"
  if [ ! -s "$out" ]; then
    echo "fetched empty artifact from $url" >&2
    exit 1
  fi
  echo "verified non-empty fetch: $url"
}

need curl
need docker
need go
need zip
need unzip

log "Bootstrap database and pinned validators"
make db-up
wait_for_database
make migrate-up
make seed
make validators-install
make validators-check

log "Import sample GTFS"
GTFS_ZIP="$TMP_DIR/valid-small.zip"
(cd testdata/gtfs/valid-small && zip -qr "$GTFS_ZIP" .)
go run ./cmd/gtfs-import -agency-id "$AGENCY_ID" -zip "$GTFS_ZIP" -actor-id demo-script -notes "Phase 10 agency demo flow" >"$TMP_DIR/import-result.json"
grep -q '"status":"published"' "$TMP_DIR/import-result.json"
echo "imported testdata/gtfs/valid-small"

log "Start services"
start_service agency-config 8081 env go run ./cmd/agency-config
start_service telemetry-ingest 8082 env go run ./cmd/telemetry-ingest
start_service feed-vehicle-positions 8083 env go run ./cmd/feed-vehicle-positions
start_service feed-trip-updates 8084 env VEHICLE_POSITIONS_FEED_URL="$FEED_BASE_URL/gtfsrt/vehicle_positions.pb" go run ./cmd/feed-trip-updates
start_service feed-alerts 8085 env go run ./cmd/feed-alerts
start_service gtfs-studio 8086 env go run ./cmd/gtfs-studio

log "Start demo public feed proxy"
cat >"$TMP_DIR/public-proxy.go" <<'GO'
package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func proxy(target string) http.Handler {
	u, err := url.Parse(target)
	if err != nil {
		log.Fatal(err)
	}
	return httputil.NewSingleHostReverseProxy(u)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/public/gtfs/", proxy("http://localhost:8081"))
	mux.Handle("/public/feeds.json", proxy("http://localhost:8081"))
	mux.Handle("/public/gtfsrt/vehicle_positions.pb", proxy("http://localhost:8083"))
	mux.Handle("/public/gtfsrt/trip_updates.pb", proxy("http://localhost:8084"))
	mux.Handle("/public/gtfsrt/alerts.pb", proxy("http://localhost:8085"))
	port := os.Getenv("PUBLIC_PROXY_PORT")
	if port == "" {
		port = "8090"
	}
	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, mux))
}
GO
go run "$TMP_DIR/public-proxy.go" >"$TMP_DIR/public-proxy.log" 2>&1 &
PIDS="$PIDS $!"
wait_for_status "$PUBLIC_BASE_URL/public/gtfs/schedule.zip" "200"
echo "public proxy listening on $PUBLIC_BASE_URL (log: $TMP_DIR/public-proxy.log)"

log "Generate admin token"
ADMIN_TOKEN="$(go run ./cmd/admin-token -sub admin@example.com -agency-id "$AGENCY_ID" | sed -n 's/^token=//p')"
if [ -z "$ADMIN_TOKEN" ]; then
  echo "admin token generation failed" >&2
  exit 1
fi
AUTH_HEADER="Authorization: Bearer $ADMIN_TOKEN"

log "Bootstrap publication metadata"
curl -fsS -X POST "http://localhost:8081/admin/publication/bootstrap" \
  -H "$AUTH_HEADER" \
  -H "Content-Type: application/json" \
  --data '{}' >"$TMP_DIR/publication-bootstrap.json"
grep -q '"stored":true' "$TMP_DIR/publication-bootstrap.json"

log "Seed consumer-ingestion workflow records through admin API"
for consumer in "Google Maps" "Apple Maps" "Transit App" "Bing Maps" "Moovit"; do
  curl -fsS -X POST "http://localhost:8081/admin/consumer-ingestion" \
    -H "$AUTH_HEADER" \
    -H "Content-Type: application/json" \
    --data "{
      \"consumer_name\": \"$consumer\",
      \"status\": \"not_started\",
      \"notes\": \"local demo workflow record\",
      \"packet\": {\"source\": \"phase_10_demo\"}
    }" >"$TMP_DIR/consumer-${consumer}.json"
done

log "Verify protected admin/debug surfaces reject anonymous requests"
expect_status 401 "http://localhost:8081/admin/compliance/scorecard"
expect_status 401 "http://localhost:8082/v1/events?limit=10"
expect_status 401 "http://localhost:8083/public/gtfsrt/vehicle_positions.json"
expect_status 401 "http://localhost:8084/public/gtfsrt/trip_updates.json"
expect_status 401 "http://localhost:8085/public/gtfsrt/alerts.json"
expect_status 401 "http://localhost:8086/admin/gtfs-studio"
expect_status 401 "http://localhost:8086/admin/gtfs-studio/drafts/demo-draft"

log "Verify protected GTFS Studio access with admin token"
curl -fsS -H "$AUTH_HEADER" "http://localhost:8086/admin/gtfs-studio" >"$TMP_DIR/gtfs-studio.html"
grep -q "GTFS Studio" "$TMP_DIR/gtfs-studio.html"
echo "verified /admin/gtfs-studio requires admin auth and succeeds with Bearer JWT"

log "Ingest authenticated telemetry"
OBSERVED_AT="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
DEVICE_TOKEN="dev-device-token"
cat >"$TMP_DIR/telemetry.json" <<EOF
{
  "agency_id": "$AGENCY_ID",
  "device_id": "device-1",
  "vehicle_id": "bus-1",
  "timestamp": "$OBSERVED_AT",
  "lat": 49.2827,
  "lon": -123.1207,
  "bearing": 120.0,
  "speed_mps": 8.4,
  "accuracy_m": 7.5,
  "trip_hint": "trip-10-0800"
}
EOF
curl -fsS -X POST "http://localhost:8082/v1/telemetry" \
  -H "Authorization: Bearer ${DEVICE_TOKEN}" \
  -H "Content-Type: application/json" \
  --data @"$TMP_DIR/telemetry.json" >"$TMP_DIR/telemetry-response.json"
grep -q "\"agency_id\":\"$AGENCY_ID\"" "$TMP_DIR/telemetry-response.json"

log "Fetch and verify public feeds"
fetch_nonempty "$PUBLIC_BASE_URL/public/gtfs/schedule.zip" "$TMP_DIR/schedule.zip"
unzip -t "$TMP_DIR/schedule.zip" >/dev/null
echo "verified schedule.zip is a readable ZIP"
fetch_nonempty "$PUBLIC_BASE_URL/public/feeds.json" "$TMP_DIR/feeds.json"
grep -q '"feeds"' "$TMP_DIR/feeds.json"
fetch_nonempty "$PUBLIC_BASE_URL/public/gtfsrt/vehicle_positions.pb" "$TMP_DIR/vehicle_positions.pb"
fetch_nonempty "$PUBLIC_BASE_URL/public/gtfsrt/trip_updates.pb" "$TMP_DIR/trip_updates.pb"
fetch_nonempty "$PUBLIC_BASE_URL/public/gtfsrt/alerts.pb" "$TMP_DIR/alerts.pb"

log "Verify protected debug/admin reads with admin token"
curl -fsS -H "$AUTH_HEADER" "http://localhost:8082/v1/events?limit=10" >"$TMP_DIR/events.json"
curl -fsS -H "$AUTH_HEADER" "http://localhost:8083/public/gtfsrt/vehicle_positions.json" >"$TMP_DIR/vehicle_positions.json"
curl -fsS -H "$AUTH_HEADER" "http://localhost:8084/public/gtfsrt/trip_updates.json" >"$TMP_DIR/trip_updates.json"
curl -fsS -H "$AUTH_HEADER" "http://localhost:8085/public/gtfsrt/alerts.json" >"$TMP_DIR/alerts.json"

log "Create and publish a simple Service Alert"
curl -fsS -X POST "http://localhost:8085/admin/alerts" \
  -H "$AUTH_HEADER" \
  -H "Content-Type: application/json" \
  --data '{
    "alert_key": "demo-alert",
    "cause": "OTHER_CAUSE",
    "effect": "OTHER_EFFECT",
    "header_text": "Demo service notice",
    "description_text": "This alert is created by the Phase 10 local demo flow.",
    "source_type": "operator",
    "entities": [{"route_id": "route-10"}],
    "publish": true
  }' >"$TMP_DIR/alert.json"
grep -q '"Status":"published"' "$TMP_DIR/alert.json"

log "Run validation flow"
curl -fsS -X POST "http://localhost:8081/admin/validation/run" \
  -H "$AUTH_HEADER" \
  -H "Content-Type: application/json" \
  --data '{"validator_id":"static-mobilitydata","feed_type":"schedule"}' >"$TMP_DIR/validate-schedule.json"
grep -q '"validator_name"' "$TMP_DIR/validate-schedule.json"
curl -fsS -X POST "http://localhost:8081/admin/validation/run" \
  -H "$AUTH_HEADER" \
  -H "Content-Type: application/json" \
  --data '{"validator_id":"realtime-mobilitydata","feed_type":"vehicle_positions"}' >"$TMP_DIR/validate-vehicle-positions.json"
grep -q '"validator_name"' "$TMP_DIR/validate-vehicle-positions.json"
curl -fsS -X POST "http://localhost:8081/admin/validation/run" \
  -H "$AUTH_HEADER" \
  -H "Content-Type: application/json" \
  --data '{"validator_id":"realtime-mobilitydata","feed_type":"trip_updates"}' >"$TMP_DIR/validate-trip-updates.json"
grep -q '"validator_name"' "$TMP_DIR/validate-trip-updates.json"
curl -fsS -X POST "http://localhost:8081/admin/validation/run" \
  -H "$AUTH_HEADER" \
  -H "Content-Type: application/json" \
  --data '{"validator_id":"realtime-mobilitydata","feed_type":"alerts"}' >"$TMP_DIR/validate-alerts.json"
grep -q '"validator_name"' "$TMP_DIR/validate-alerts.json"

log "Inspect scorecard and consumer-ingestion workflow records"
curl -fsS -X POST "http://localhost:8081/admin/compliance/scorecard" \
  -H "$AUTH_HEADER" \
  -H "Content-Type: application/json" \
  --data '{}' >"$TMP_DIR/scorecard.json"
grep -q '"overall_status"' "$TMP_DIR/scorecard.json"
curl -fsS -H "$AUTH_HEADER" "http://localhost:8081/admin/compliance/scorecard" >"$TMP_DIR/latest-scorecard.json"
curl -fsS -H "$AUTH_HEADER" "http://localhost:8081/admin/consumer-ingestion" >"$TMP_DIR/consumer-ingestion.json"
grep -q '"consumers"' "$TMP_DIR/consumer-ingestion.json"

cat <<EOF

Phase 10 agency demo flow completed.

Verified:
  - database bootstrap, migrations, and seed data
  - pinned validator install/check
  - sample GTFS import
  - publication metadata bootstrap
  - device-token telemetry ingest
  - public schedule.zip, feeds.json, and realtime protobuf fetches
  - protected debug/admin routes, including GTFS Studio
  - validation run flow
  - scorecard and consumer-ingestion visibility

Artifacts and service logs remain under:
  $TMP_DIR

EOF
