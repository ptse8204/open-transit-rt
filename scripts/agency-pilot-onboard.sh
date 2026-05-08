#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

COMPOSE_FILE="${COMPOSE_FILE:-deploy/docker-compose.yml}"
PROFILE="app"
APP_SERVICES="agency-config telemetry-ingest feed-vehicle-positions feed-trip-updates feed-alerts gtfs-studio local-proxy"

MODE="${MODE:-local-compose}"
AGENCY_ID="${AGENCY_ID:-}"
GTFS_URL="${GTFS_URL:-}"
PUBLIC_BASE_URL="${PUBLIC_BASE_URL:-http://localhost:8080}"
ADMIN_BASE_URL="${ADMIN_BASE_URL:-}"
ADMIN_TOKEN="${ADMIN_TOKEN:-}"
IMPORT_TIMEOUT="${GTFS_IMPORT_TIMEOUT:-15m}"
ADMIN_SUBJECT="${ADMIN_SUBJECT:-admin@example.com}"
TECHNICAL_CONTACT_EMAIL="${TECHNICAL_CONTACT_EMAIL:-operator-placeholder@example.invalid}"
FEED_LICENSE_NAME="${FEED_LICENSE_NAME:-replace-with-agency-approved-license}"
FEED_LICENSE_URL="${FEED_LICENSE_URL:-https://example.invalid/replace-with-agency-license}"
PUBLICATION_ENVIRONMENT="${PUBLICATION_ENVIRONMENT:-dev}"
RESET_LOCAL_STATE="false"
FORCE_RESET="false"
STRICT_VALIDATORS="${STRICT_VALIDATORS:-false}"
SKIP_VALIDATORS="${SKIP_VALIDATORS:-false}"
DRY_RUN="${DRY_RUN:-false}"

usage() {
  cat <<'EOF'
Usage:
  scripts/agency-pilot-onboard.sh --agency-id <id> --gtfs-url <url> [options]

Required:
  --agency-id ID              GTFS agency_id to import and publish
  --gtfs-url URL              http(s) URL for a GTFS ZIP

Options:
  --mode MODE                 local-compose or running (default: local-compose)
  --public-base-url URL       public root, default http://localhost:8080
  --admin-base-url URL        admin API root; defaults to public base URL in local-compose mode only
  --admin-token TOKEN         required in running mode unless already available
  --admin-subject SUBJECT     local admin JWT subject, default admin@example.com
  --import-timeout DURATION   0, 90s, 15m, or 1h style duration, default 15m
  --technical-contact-email EMAIL
  --feed-license-name NAME
  --feed-license-url URL
  --reset-local-state         destructive local Compose reset after confirmation
  --force                     skip reset confirmation when used with --reset-local-state
  --strict-validators         fail if validators are missing or a validator run fails
  --skip-validators           skip validator API calls
  --dry-run                   validate inputs, print planned paths/URLs, then exit
  -h, --help                  show this help

Environment fallbacks:
  AGENCY_ID, GTFS_URL, PUBLIC_BASE_URL, ADMIN_BASE_URL, ADMIN_TOKEN,
  GTFS_IMPORT_TIMEOUT, TECHNICAL_CONTACT_EMAIL, FEED_LICENSE_NAME,
  FEED_LICENSE_URL, PUBLICATION_ENVIRONMENT, ADMIN_SUBJECT, MODE,
  STRICT_VALIDATORS, SKIP_VALIDATORS

Notes:
  local-compose mode starts/builds/migrates/services directly and imports only
  the requested GTFS. It does not call make agency-app-up or import the demo
  sample feed.

  running mode requires --admin-base-url or ADMIN_BASE_URL. The running-mode
  admin URL should normally be loopback, VPN, SSH tunnel, or another
  private/admin-protected URL.

  Publication metadata is local/reference placeholder metadata unless the operator supplied agency-approved values.
EOF
}

log() {
  printf '\n==> %s\n' "$1"
}

warn() {
  printf 'WARNING: %s\n' "$1" >&2
}

fail() {
  printf '\nERROR: %s\n' "$1" >&2
  exit 1
}

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    fail "missing required tool: $1"
  fi
}

dc() {
  docker compose -f "$COMPOSE_FILE" "$@"
}

sql_quote() {
  printf "%s" "$1" | sed "s/'/''/g"
}

json_body() {
  python3 - "$@" <<'PY'
import json
import sys

keys = [
    "public_base_url",
    "feed_base_url",
    "technical_contact_email",
    "license_name",
    "license_url",
    "publication_environment",
]
print(json.dumps({k: v for k, v in zip(keys, sys.argv[1:])}, separators=(",", ":")))
PY
}

sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

validate_url() {
  name="$1"
  value="$2"
  case "$value" in
    http://*|https://*) return 0 ;;
    *) fail "$name must start with http:// or https://" ;;
  esac
}

validate_inputs() {
  case "$MODE" in
    local-compose|running) ;;
    *) fail "--mode must be exactly local-compose or running" ;;
  esac
  if [ -z "$AGENCY_ID" ]; then
    fail "--agency-id or AGENCY_ID is required"
  fi
  case "$AGENCY_ID" in
    .|..|.*)
      fail "AGENCY_ID must not be '.', '..', or start with a dot"
      ;;
  esac
  case "$AGENCY_ID" in
    *[!A-Za-z0-9._-]*)
      fail "AGENCY_ID may contain only letters, numbers, dot, underscore, and hyphen"
      ;;
  esac
  if [ -z "$GTFS_URL" ]; then
    fail "--gtfs-url or GTFS_URL is required"
  fi
  validate_url "GTFS_URL" "$GTFS_URL"
  validate_url "PUBLIC_BASE_URL" "$PUBLIC_BASE_URL"
  validate_url "FEED_LICENSE_URL" "$FEED_LICENSE_URL"
  if [ -z "$ADMIN_BASE_URL" ]; then
    if [ "$MODE" = "running" ]; then
      fail "running mode requires --admin-base-url or ADMIN_BASE_URL; use a loopback, VPN, SSH tunnel, or otherwise private/admin-protected URL"
    fi
    ADMIN_BASE_URL="$PUBLIC_BASE_URL"
  fi
  validate_url "ADMIN_BASE_URL" "$ADMIN_BASE_URL"
  if [ "$MODE" = "running" ]; then
    case "$ADMIN_BASE_URL" in
      http://127.0.0.1:*|http://localhost:*|https://127.0.0.1:*|https://localhost:*) ;;
      *) warn "running-mode admin URL should normally be loopback, VPN, SSH tunnel, or another private/admin-protected URL" ;;
    esac
  fi
  if ! printf "%s" "$IMPORT_TIMEOUT" | grep -Eq '^(0|[0-9]+(ns|us|ms|s|m|h))$'; then
    fail "--import-timeout must be 0 or a simple Go duration such as 90s, 15m, or 1h"
  fi
  if [ "$MODE" = "running" ] && [ -z "$ADMIN_TOKEN" ]; then
    fail "running mode requires --admin-token or ADMIN_TOKEN"
  fi
  case "$STRICT_VALIDATORS" in true|false) ;;
    *) fail "STRICT_VALIDATORS must be true or false" ;;
  esac
  case "$SKIP_VALIDATORS" in true|false) ;;
    *) fail "SKIP_VALIDATORS must be true or false" ;;
  esac
  if [ "$STRICT_VALIDATORS" = "true" ] && [ "$SKIP_VALIDATORS" = "true" ]; then
    fail "--strict-validators and --skip-validators cannot be used together"
  fi
  case "$PUBLICATION_ENVIRONMENT" in
    dev|production) ;;
    reference)
      warn "PUBLICATION_ENVIRONMENT=reference is a deployment label; storing publication_environment=dev because the current schema accepts dev or production"
      PUBLICATION_ENVIRONMENT="dev"
      ;;
    *) fail "PUBLICATION_ENVIRONMENT must be dev or production" ;;
  esac
}

print_plan() {
  cache_dir="${AGENCY_PILOT_CACHE_DIR:-.cache/agency-pilot}/$AGENCY_ID"
  cat <<EOF
Agency pilot onboarding plan:
  mode: $MODE
  agency_id: $AGENCY_ID
  gtfs_url: $GTFS_URL
  cache_dir: $cache_dir
  source_zip: $cache_dir/source.zip
  public_base_url: $PUBLIC_BASE_URL
  admin_base_url: $ADMIN_BASE_URL
  import_timeout: $IMPORT_TIMEOUT
  technical_contact_email: $TECHNICAL_CONTACT_EMAIL
  feed_license_name: $FEED_LICENSE_NAME
  feed_license_url: $FEED_LICENSE_URL
  publication_environment: $PUBLICATION_ENVIRONMENT
  reset_local_state: $RESET_LOCAL_STATE
  validators: $(if [ "$SKIP_VALIDATORS" = "true" ]; then printf 'skipped'; elif [ "$STRICT_VALIDATORS" = "true" ]; then printf 'strict'; else printf 'best-effort'; fi)

Publication metadata is local/reference placeholder metadata unless the operator supplied agency-approved values.

No download, import, database write, service start, validator run, or evidence creation was performed.
EOF
}

parse_args() {
  while [ "$#" -gt 0 ]; do
    case "$1" in
      --agency-id) AGENCY_ID="${2:-}"; shift 2 ;;
      --gtfs-url) GTFS_URL="${2:-}"; shift 2 ;;
      --public-base-url) PUBLIC_BASE_URL="${2:-}"; shift 2 ;;
      --admin-base-url) ADMIN_BASE_URL="${2:-}"; shift 2 ;;
      --admin-token) ADMIN_TOKEN="${2:-}"; shift 2 ;;
      --admin-subject) ADMIN_SUBJECT="${2:-}"; shift 2 ;;
      --import-timeout) IMPORT_TIMEOUT="${2:-}"; shift 2 ;;
      --technical-contact-email) TECHNICAL_CONTACT_EMAIL="${2:-}"; shift 2 ;;
      --feed-license-name) FEED_LICENSE_NAME="${2:-}"; shift 2 ;;
      --feed-license-url) FEED_LICENSE_URL="${2:-}"; shift 2 ;;
      --mode) MODE="${2:-}"; shift 2 ;;
      --reset-local-state) RESET_LOCAL_STATE="true"; shift ;;
      --force) FORCE_RESET="true"; shift ;;
      --strict-validators) STRICT_VALIDATORS="true"; shift ;;
      --skip-validators) SKIP_VALIDATORS="true"; shift ;;
      --dry-run) DRY_RUN="true"; shift ;;
      -h|--help|help) usage; exit 0 ;;
      *) usage; fail "unknown argument: $1" ;;
    esac
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
      fail "$label did not become ready at $url"
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
      fail "$label did not become ready at $url"
    fi
    sleep 2
  done
}

fetch_nonempty() {
  url="$1"
  out="$2"
  label="$3"
  if ! curl -fsS -o "$out" "$url"; then
    fail "could not fetch $label from $url"
  fi
  if [ ! -s "$out" ]; then
    fail "$label fetched from $url was empty"
  fi
}

maybe_reset_local_state() {
  if [ "$RESET_LOCAL_STATE" != "true" ]; then
    if docker volume inspect deploy_postgres-data >/dev/null 2>&1; then
      warn "existing local Compose database volume deploy_postgres-data will be reused non-destructively"
    fi
    return 0
  fi

  if [ "$FORCE_RESET" != "true" ]; then
    printf 'This will delete the local Compose Postgres volume. Type reset-agency-pilot to continue: '
    read answer
    if [ "$answer" != "reset-agency-pilot" ]; then
      echo "Reset canceled."
      exit 0
    fi
  fi
  log "Reset local Compose containers and volumes"
  dc --profile "$PROFILE" down -v --remove-orphans
}

download_gtfs() {
  RUN_DIR="${AGENCY_PILOT_CACHE_DIR:-.cache/agency-pilot}/$AGENCY_ID"
  FETCH_DIR="$RUN_DIR/fetched"
  SOURCE_ZIP="$RUN_DIR/source.zip"
  SOURCE_SHA="$RUN_DIR/source.sha256"
  mkdir -p "$RUN_DIR" "$FETCH_DIR"
  log "Download GTFS ZIP into ignored storage"
  curl -L --fail --silent --show-error -o "$SOURCE_ZIP" "$GTFS_URL"
  source_hash="$(sha256_file "$SOURCE_ZIP")"
  printf '%s  %s\n' "$source_hash" "$(basename "$SOURCE_ZIP")" >"$SOURCE_SHA"
}

start_local_compose() {
  need docker
  if ! docker info >/dev/null 2>&1; then
    fail "Docker is not available or the Docker daemon is not running"
  fi

  maybe_reset_local_state

  export AGENCY_ID PUBLIC_BASE_URL
  export FEED_BASE_URL="${FEED_BASE_URL:-$PUBLIC_BASE_URL/public}"
  export SCHEDULE_FEED_URL="${SCHEDULE_FEED_URL:-$PUBLIC_BASE_URL/public/gtfs/schedule.zip}"
  export VEHICLE_POSITIONS_FEED_URL="${VEHICLE_POSITIONS_FEED_URL:-$PUBLIC_BASE_URL/public/gtfsrt/vehicle_positions.pb}"
  export TRIP_UPDATES_FEED_URL="${TRIP_UPDATES_FEED_URL:-$PUBLIC_BASE_URL/public/gtfsrt/trip_updates.pb}"
  export ALERTS_FEED_URL="${ALERTS_FEED_URL:-$PUBLIC_BASE_URL/public/gtfsrt/alerts.pb}"
  export REALTIME_VALIDATION_BASE_URL="${REALTIME_VALIDATION_BASE_URL:-$PUBLIC_BASE_URL/public}"
  export TECHNICAL_CONTACT_EMAIL FEED_LICENSE_NAME FEED_LICENSE_URL PUBLICATION_ENVIRONMENT

  log "Start Postgres/PostGIS"
  dc up -d postgres
  attempt=0
  until dc exec -T postgres pg_isready -U postgres -d open_transit_rt >/dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [ "$attempt" -ge 45 ]; then
      fail "Postgres did not become ready in time"
    fi
    sleep 2
  done

  log "Build local app image"
  dc --profile "$PROFILE" build agency-config

  log "Apply migrations"
  dc --profile "$PROFILE" run --rm --no-deps agency-config /app/bin/migrate up

  log "Seed requested agency and admin roles"
  seed_agency_admin_local

  log "Start local app services"
  dc --profile "$PROFILE" up -d $APP_SERVICES

  wait_for_container_url agency-config "http://127.0.0.1:8081/healthz" "agency-config health"
  wait_for_container_url telemetry-ingest "http://127.0.0.1:8082/healthz" "telemetry-ingest health"
  wait_for_container_url feed-vehicle-positions "http://127.0.0.1:8083/healthz" "Vehicle Positions health"
  wait_for_container_url feed-trip-updates "http://127.0.0.1:8084/healthz" "Trip Updates health"
  wait_for_container_url feed-alerts "http://127.0.0.1:8085/healthz" "Alerts health"
  wait_for_container_url gtfs-studio "http://127.0.0.1:8086/healthz" "GTFS Studio health"
  wait_for_url "$PUBLIC_BASE_URL/healthz" "local reverse proxy"
}

seed_agency_admin_local() {
  agency_sql="$(sql_quote "$AGENCY_ID")"
  subject_sql="$(sql_quote "$ADMIN_SUBJECT")"
  contact_sql="$(sql_quote "$TECHNICAL_CONTACT_EMAIL")"
  public_base_sql="$(sql_quote "$PUBLIC_BASE_URL")"
  dc exec -T postgres psql -U postgres -d open_transit_rt <<SQL
INSERT INTO agency (id, name, timezone, contact_email, public_url)
VALUES ('$agency_sql', '$agency_sql local/reference placeholder', 'Etc/UTC', '$contact_sql', '$public_base_sql')
ON CONFLICT (id) DO UPDATE
SET contact_email = COALESCE(agency.contact_email, EXCLUDED.contact_email),
    public_url = EXCLUDED.public_url,
    updated_at = now();

WITH upserted AS (
  INSERT INTO agency_user (agency_id, email, display_name, auth_subject)
  VALUES ('$agency_sql', '$subject_sql', 'Agency Pilot Admin', '$subject_sql')
  ON CONFLICT (agency_id, email) DO UPDATE
  SET display_name = EXCLUDED.display_name,
      auth_subject = EXCLUDED.auth_subject
  RETURNING id, agency_id
)
INSERT INTO role_binding (agency_id, agency_user_id, role)
SELECT agency_id, id, role
FROM upserted
CROSS JOIN (VALUES ('admin'), ('editor'), ('operator'), ('read_only')) AS roles(role)
ON CONFLICT (agency_id, agency_user_id, role) DO NOTHING;
SQL
}

seed_agency_admin_running() {
  need psql
  : "${DATABASE_URL:?running mode requires DATABASE_URL for agency/admin upsert and gtfs-import}"
  log "Upsert requested agency and admin roles through DATABASE_URL"
  agency_sql="$(sql_quote "$AGENCY_ID")"
  subject_sql="$(sql_quote "$ADMIN_SUBJECT")"
  contact_sql="$(sql_quote "$TECHNICAL_CONTACT_EMAIL")"
  public_base_sql="$(sql_quote "$PUBLIC_BASE_URL")"
  psql "$DATABASE_URL" <<SQL
INSERT INTO agency (id, name, timezone, contact_email, public_url)
VALUES ('$agency_sql', '$agency_sql local/reference placeholder', 'Etc/UTC', '$contact_sql', '$public_base_sql')
ON CONFLICT (id) DO UPDATE
SET contact_email = COALESCE(agency.contact_email, EXCLUDED.contact_email),
    public_url = EXCLUDED.public_url,
    updated_at = now();

WITH upserted AS (
  INSERT INTO agency_user (agency_id, email, display_name, auth_subject)
  VALUES ('$agency_sql', '$subject_sql', 'Agency Pilot Admin', '$subject_sql')
  ON CONFLICT (agency_id, email) DO UPDATE
  SET display_name = EXCLUDED.display_name,
      auth_subject = EXCLUDED.auth_subject
  RETURNING id, agency_id
)
INSERT INTO role_binding (agency_id, agency_user_id, role)
SELECT agency_id, id, role
FROM upserted
CROSS JOIN (VALUES ('admin'), ('editor'), ('operator'), ('read_only')) AS roles(role)
ON CONFLICT (agency_id, agency_user_id, role) DO NOTHING;
SQL
}

import_local_compose() {
  log "Copy GTFS ZIP into agency-config container"
  container_id="$(dc ps -q agency-config)"
  if [ -z "$container_id" ]; then
    fail "agency-config container is not running"
  fi
  docker cp "$SOURCE_ZIP" "$container_id:/tmp/agency-pilot-source.zip"

  log "Import requested GTFS and publish it as active feed"
  import_output="$RUN_DIR/import-result.json"
  if ! dc exec -T agency-config /app/bin/gtfs-import \
      -agency-id "$AGENCY_ID" \
      -zip /tmp/agency-pilot-source.zip \
      -actor-id agency-pilot-onboard \
      -notes "Phase 37 reusable agency onboarding local/reference flow" \
      -timeout "$IMPORT_TIMEOUT" >"$import_output"; then
    cat "$import_output" >&2 || true
    fail "GTFS import failed"
  fi
  if ! grep -q '"status":"published"' "$import_output"; then
    cat "$import_output" >&2 || true
    fail "GTFS import did not publish successfully"
  fi
}

import_running() {
  need go
  : "${DATABASE_URL:?running mode requires DATABASE_URL for gtfs-import}"
  log "Import requested GTFS through local Go command against configured DATABASE_URL"
  import_output="$RUN_DIR/import-result.json"
  if ! DATABASE_URL="$DATABASE_URL" go run ./cmd/gtfs-import \
      -agency-id "$AGENCY_ID" \
      -zip "$SOURCE_ZIP" \
      -actor-id agency-pilot-onboard \
      -notes "Phase 37 reusable agency onboarding running-mode flow" \
      -timeout "$IMPORT_TIMEOUT" >"$import_output"; then
    cat "$import_output" >&2 || true
    fail "GTFS import failed"
  fi
  if ! grep -q '"status":"published"' "$import_output"; then
    cat "$import_output" >&2 || true
    fail "GTFS import did not publish successfully"
  fi
}

preflight_running() {
  : "${DATABASE_URL:?running mode requires DATABASE_URL for agency/admin upsert and gtfs-import}"
  need go
  need psql
}

admin_token_local_compose() {
  ADMIN_TOKEN="$(dc exec -T agency-config /app/bin/admin-token -sub "$ADMIN_SUBJECT" -agency-id "$AGENCY_ID" | sed -n 's/^token=//p')"
  if [ -z "$ADMIN_TOKEN" ]; then
    fail "could not generate local admin token"
  fi
}

bootstrap_metadata() {
  log "Bootstrap publication metadata"
  body="$(json_body "$PUBLIC_BASE_URL" "$PUBLIC_BASE_URL/public" "$TECHNICAL_CONTACT_EMAIL" "$FEED_LICENSE_NAME" "$FEED_LICENSE_URL" "$PUBLICATION_ENVIRONMENT")"
  out="$RUN_DIR/publication-bootstrap.json"
  if ! curl -fsS -X POST "$ADMIN_BASE_URL/admin/publication/bootstrap" \
      -H "Authorization: Bearer $ADMIN_TOKEN" \
      -H "Content-Type: application/json" \
      --data "$body" >"$out"; then
    fail "publication metadata bootstrap failed"
  fi
  if ! grep -Eq '"stored"[[:space:]]*:[[:space:]]*true' "$out"; then
    cat "$out" >&2 || true
    fail "publication metadata bootstrap did not report stored=true"
  fi
}

summarize_zip() {
  zip_path="$1"
  out_path="$2"
  python3 - "$zip_path" "$out_path" <<'PY'
import csv
import json
import sys
import zipfile

zip_path, out_path = sys.argv[1:3]

def rows(zf, name):
    try:
        with zf.open(name) as f:
            return list(csv.DictReader(line.decode("utf-8-sig") for line in f))
    except KeyError:
        return []

with zipfile.ZipFile(zip_path) as zf:
    agency = rows(zf, "agency.txt")
    routes = rows(zf, "routes.txt")
    stops = rows(zf, "stops.txt")
    trips = rows(zf, "trips.txt")
    calendar = rows(zf, "calendar.txt")
    calendar_dates = rows(zf, "calendar_dates.txt")

dates = []
for row in calendar:
    for key in ("start_date", "end_date"):
        value = row.get(key, "")
        if value:
            dates.append(value)
for row in calendar_dates:
    value = row.get("date", "")
    if value:
        dates.append(value)

summary = {
    "agency_ids": sorted({row.get("agency_id", "") for row in agency if row.get("agency_id", "")}),
    "agency_names": sorted({row.get("agency_name", "") for row in agency if row.get("agency_name", "")}),
    "timezones": sorted({row.get("agency_timezone", "") for row in agency if row.get("agency_timezone", "")}),
    "route_count": len(routes),
    "stop_count": len(stops),
    "trip_count": len(trips),
    "route_type_counts": {},
    "service_date_min": min(dates) if dates else "",
    "service_date_max": max(dates) if dates else "",
}
for row in routes:
    route_type = row.get("route_type", "")
    summary["route_type_counts"][route_type] = summary["route_type_counts"].get(route_type, 0) + 1
with open(out_path, "w", encoding="utf-8") as f:
    json.dump(summary, f, indent=2, sort_keys=True)
    f.write("\n")
PY
}

verify_schedule_identity() {
  log "Verify fetched schedule matches imported GTFS summary"
  summarize_zip "$SOURCE_ZIP" "$RUN_DIR/source-summary.json"
  summarize_zip "$FETCH_DIR/schedule.zip" "$RUN_DIR/fetched-schedule-summary.json"
  python3 - "$AGENCY_ID" "$RUN_DIR/source-summary.json" "$RUN_DIR/fetched-schedule-summary.json" <<'PY'
import json
import sys

agency_id, source_path, fetched_path = sys.argv[1:4]
source = json.load(open(source_path, encoding="utf-8"))
fetched = json.load(open(fetched_path, encoding="utf-8"))
if agency_id not in fetched.get("agency_ids", []):
    raise SystemExit(f"fetched schedule does not contain requested agency_id {agency_id!r}")
keys = ["agency_ids", "route_count", "stop_count", "trip_count", "service_date_min", "service_date_max"]
mismatches = {key: (source.get(key), fetched.get(key)) for key in keys if source.get(key) != fetched.get(key)}
if mismatches:
    raise SystemExit(f"fetched schedule summary does not match source summary: {mismatches}")
if agency_id != "demo-agency" and "demo-agency" in fetched.get("agency_ids", []):
    raise SystemExit("fetched schedule unexpectedly contains demo-agency")
print("schedule identity check passed")
PY
}

fetch_public_paths() {
  log "Fetch five public paths"
  fetch_nonempty "$PUBLIC_BASE_URL/public/feeds.json" "$FETCH_DIR/feeds.json" "feeds.json"
  fetch_nonempty "$PUBLIC_BASE_URL/public/gtfs/schedule.zip" "$FETCH_DIR/schedule.zip" "schedule.zip"
  fetch_nonempty "$PUBLIC_BASE_URL/public/gtfsrt/vehicle_positions.pb" "$FETCH_DIR/vehicle_positions.pb" "Vehicle Positions protobuf"
  fetch_nonempty "$PUBLIC_BASE_URL/public/gtfsrt/trip_updates.pb" "$FETCH_DIR/trip_updates.pb" "Trip Updates protobuf"
  fetch_nonempty "$PUBLIC_BASE_URL/public/gtfsrt/alerts.pb" "$FETCH_DIR/alerts.pb" "Alerts protobuf"
  for file in feeds.json schedule.zip vehicle_positions.pb trip_updates.pb alerts.pb; do
    printf '%s  %s\n' "$(sha256_file "$FETCH_DIR/$file")" "$file"
  done >"$FETCH_DIR/SHA256SUMS.txt"
  verify_schedule_identity
}

run_validators() {
  if [ "$SKIP_VALIDATORS" = "true" ]; then
    VALIDATOR_STATUS="skipped: --skip-validators was supplied"
    return 0
  fi
  if ! ./scripts/check-validators.sh >/dev/null 2>&1; then
    VALIDATOR_STATUS="blocked: pinned validator tooling is missing or misconfigured; run make validators-install && make validators-check"
    if [ "$STRICT_VALIDATORS" = "true" ]; then
      fail "$VALIDATOR_STATUS"
    fi
    return 0
  fi

  log "Run validator API calls"
  validation_dir="$RUN_DIR/validation"
  mkdir -p "$validation_dir"
  failures=0
  warning_statuses=0
  for feed_type in schedule vehicle_positions trip_updates alerts; do
    if [ "$feed_type" = "schedule" ]; then
      validator_id="static-mobilitydata"
    else
      validator_id="realtime-mobilitydata"
    fi
    out="$validation_dir/$feed_type.json"
    if ! curl -fsS -X POST "$ADMIN_BASE_URL/admin/validation/run" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        --data "{\"validator_id\":\"$validator_id\",\"feed_type\":\"$feed_type\"}" >"$out"; then
      failures=$((failures + 1))
      continue
    fi
    if ! grep -q '"validator_name"' "$out"; then
      failures=$((failures + 1))
      continue
    fi
    status="$(python3 - "$out" <<'PY'
import json
import sys
try:
    print(json.load(open(sys.argv[1], encoding="utf-8")).get("status", "unknown"))
except Exception:
    print("unknown")
PY
)"
    case "$status" in
      passed) ;;
      warning) warning_statuses=$((warning_statuses + 1)) ;;
      *) failures=$((failures + 1)) ;;
    esac
  done
  if [ "$failures" -gt 0 ]; then
    VALIDATOR_STATUS="blocked: $failures validator API call(s) failed or returned failed/not_run/unknown; run make validators-install && make validators-check, then configure validator paths for this runtime"
    if [ "$STRICT_VALIDATORS" = "true" ]; then
      fail "$VALIDATOR_STATUS"
    fi
  elif [ "$warning_statuses" -gt 0 ]; then
    VALIDATOR_STATUS="passed: validator API calls completed with $warning_statuses warning status(es); review $validation_dir"
  else
    VALIDATOR_STATUS="passed: validator API calls completed"
  fi
}

print_success() {
  generated_at="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
  cat <<EOF

Agency pilot onboarding completed.

Public feed URLs:
  Feed discovery:             $PUBLIC_BASE_URL/public/feeds.json
  Schedule ZIP:               $PUBLIC_BASE_URL/public/gtfs/schedule.zip
  Vehicle Positions protobuf: $PUBLIC_BASE_URL/public/gtfsrt/vehicle_positions.pb
  Trip Updates protobuf:      $PUBLIC_BASE_URL/public/gtfsrt/trip_updates.pb
  Alerts protobuf:            $PUBLIC_BASE_URL/public/gtfsrt/alerts.pb

Admin URL:
  $ADMIN_BASE_URL/admin/operations

Local output:
  Source ZIP:       $SOURCE_ZIP
  Source SHA-256:   $SOURCE_SHA
  Fetch summaries:  $FETCH_DIR
  Import summary:   $RUN_DIR/import-result.json

Validator status:
  $VALIDATOR_STATUS

Metadata warning:
  Publication metadata is local/reference placeholder metadata unless the operator supplied agency-approved values.
  Do not treat placeholder metadata as agency-approved, final-root evidence, consumer evidence, CAL-ITP/Caltrans compliance, or production readiness.

Next steps:
  1. Review $PUBLIC_BASE_URL/public/feeds.json for agency-approved contact and license metadata.
  2. Run validators again after metadata and GTFS fixes if validator status is blocked or warning-bearing.
  3. Connect real telemetry only after device credentials and public-safe identifiers are reviewed.
  4. Do not treat this run as agency approval, consumer acceptance, CAL-ITP/Caltrans compliance, production readiness, vendor compatibility, or production-grade ETA proof.

Copy/paste support summary:
  generated_at=$generated_at
  agency_id=$AGENCY_ID
  mode=$MODE
  public_base_url=$PUBLIC_BASE_URL
  admin_base_url=$ADMIN_BASE_URL
  source_sha256=$(sha256_file "$SOURCE_ZIP")
  validator_status="$VALIDATOR_STATUS"
  external_evidence_created=false
  consumer_statuses_changed=false
EOF
}

main() {
  parse_args "$@"
  validate_inputs
  if [ "$DRY_RUN" = "true" ]; then
    print_plan
    exit 0
  fi

  need curl
  need grep
  need sed
  need awk
  need python3
  need unzip

  if [ "$MODE" = "running" ]; then
    preflight_running
  fi

  download_gtfs

  if [ "$MODE" = "local-compose" ]; then
    start_local_compose
    import_local_compose
    admin_token_local_compose
  else
    seed_agency_admin_running
    import_running
    wait_for_url "$PUBLIC_BASE_URL/public/feeds.json" "public feed discovery"
  fi

  bootstrap_metadata
  fetch_public_paths
  run_validators
  print_success
}

main "$@"
