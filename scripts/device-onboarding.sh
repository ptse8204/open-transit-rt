#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

TARGET="${TARGET:-http://localhost:8080}"
AGENCY_ID="${AGENCY_ID:-demo-agency}"
DEVICE_ID="${DEVICE_ID:-device-1}"
VEHICLE_ID="${VEHICLE_ID:-bus-1}"
ADMIN_SUBJECT="${ADMIN_SUBJECT:-admin@example.com}"
DEVICE_TOKEN_VALUE="${DEVICE_TOKEN:-dev-device-token}"
DRY_RUN="false"

usage() {
  cat <<'EOF'
Usage:
  scripts/device-onboarding.sh help
  scripts/device-onboarding.sh rebind [--target URL] [--device-id ID] [--vehicle-id ID] [--admin-token TOKEN]
  scripts/device-onboarding.sh sample [--target URL] [--device-token TOKEN] [--dry-run]
  scripts/device-onboarding.sh simulate [--target URL] [--device-token TOKEN] [--dry-run]

Defaults target the local app package at http://localhost:8080 with:
  agency_id=demo-agency device_id=device-1 vehicle_id=bus-1

Security notes:
  - Device tokens are Bearer credentials.
  - Rebind prints the one-time token only because the existing API intentionally returns it.
  - This helper does not print JWT secrets, CSRF secrets, DB passwords, private keys, or .cache material.
EOF
}

log() {
  printf '\n==> %s\n' "$1"
}

fail() {
  printf '\nERROR: %s\n' "$1" >&2
  exit 1
}

parse_common() {
  ADMIN_TOKEN="${ADMIN_TOKEN:-}"
  while [ "$#" -gt 0 ]; do
    case "$1" in
      --target) TARGET="$2"; shift 2 ;;
      --agency-id) AGENCY_ID="$2"; shift 2 ;;
      --device-id) DEVICE_ID="$2"; shift 2 ;;
      --vehicle-id) VEHICLE_ID="$2"; shift 2 ;;
      --admin-token) ADMIN_TOKEN="$2"; shift 2 ;;
      --device-token) DEVICE_TOKEN_VALUE="$2"; shift 2 ;;
      --dry-run) DRY_RUN="true"; shift ;;
      -h|--help) usage; exit 0 ;;
      *) fail "Unknown option: $1" ;;
    esac
  done
}

local_admin_token() {
  if [ -n "${ADMIN_TOKEN:-}" ]; then
    printf '%s' "$ADMIN_TOKEN"
    return
  fi
  if ! command -v docker >/dev/null 2>&1; then
    fail "ADMIN_TOKEN is required when Docker is unavailable."
  fi
  token="$(docker compose -f deploy/docker-compose.yml --profile app exec -T agency-config /app/bin/admin-token -sub "$ADMIN_SUBJECT" -agency-id "$AGENCY_ID" 2>/dev/null | sed -n 's/^token=//p' || true)"
  if [ -z "$token" ]; then
    fail "Could not generate a local admin token. Start the local app with make agency-app-up or set ADMIN_TOKEN."
  fi
  printf '%s' "$token"
}

payload_for() {
  lat="$1"
  lon="$2"
  bearing="$3"
  speed="$4"
  observed_at="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
  cat <<EOF
{
  "agency_id": "$AGENCY_ID",
  "device_id": "$DEVICE_ID",
  "vehicle_id": "$VEHICLE_ID",
  "timestamp": "$observed_at",
  "lat": $lat,
  "lon": $lon,
  "bearing": $bearing,
  "speed_mps": $speed,
  "accuracy_m": 7.5,
  "trip_hint": "trip-10-0800"
}
EOF
}

send_payload() {
  payload="$1"
  endpoint="$TARGET/v1/telemetry"
  echo "Target: $endpoint"
  echo "Payload:"
  printf '%s\n' "$payload"
  if [ "$DRY_RUN" = "true" ]; then
    echo "Dry run only; no telemetry was sent."
    return
  fi
  printf '%s\n' "$payload" | curl -fsS -X POST "$endpoint" \
    -H "Authorization: Bearer ${DEVICE_TOKEN_VALUE}" \
    -H "Content-Type: application/json" \
    --data @-
  echo
}

cmd_rebind() {
  parse_common "$@"
  admin_token="$(local_admin_token)"
  endpoint="$TARGET/admin/devices/rebind"
  log "Rotate and bind device credential"
  echo "Target: $endpoint"
  echo "Device: $DEVICE_ID"
  echo "Vehicle: $VEHICLE_ID"
  response="$(curl -fsS -X POST "$endpoint" \
    -H "Authorization: Bearer $admin_token" \
    -H "Content-Type: application/json" \
    --data "{
      \"device_id\": \"$DEVICE_ID\",
      \"vehicle_id\": \"$VEHICLE_ID\",
      \"reason\": \"local device onboarding helper\"
    }")"
  token="$(printf '%s\n' "$response" | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')"
  if [ -z "$token" ]; then
    echo "Rebind succeeded, but the API response did not include a one-time token."
    echo "Response:"
    printf '%s\n' "$response"
    return
  fi
  cat <<EOF

Device rebound. Store this one-time device token securely now:
  $token

Use it as:
  Authorization: Bearer <device-token>

The previous token for this device binding is invalid after rotation.
EOF
}

cmd_sample() {
  parse_common "$@"
  log "Send one sample telemetry event"
  send_payload "$(payload_for 49.2827 -123.1207 120.0 8.4)"
}

cmd_simulate() {
  parse_common "$@"
  log "Send a short demo telemetry sequence"
  send_payload "$(payload_for 49.2827 -123.1207 118.0 8.0)"
  sleep 1
  send_payload "$(payload_for 49.2832 -123.1195 121.0 8.5)"
  sleep 1
  send_payload "$(payload_for 49.2838 -123.1184 124.0 8.1)"
}

case "${1:-help}" in
  help|-h|--help) usage ;;
  rebind) shift; cmd_rebind "$@" ;;
  sample) shift; cmd_sample "$@" ;;
  simulate) shift; cmd_simulate "$@" ;;
  *) usage; exit 2 ;;
esac
