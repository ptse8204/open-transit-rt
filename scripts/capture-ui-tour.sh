#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

BASE_URL="${BASE_URL:-http://localhost:8080}"
OUT_DIR="${OUT_DIR:-site/assets/screenshots}"
STATE_FILE="${STATE_FILE:-.cache/ui-tour-storage.json}"
AGENCY_ID="${AGENCY_ID:-demo-agency}"
ADMIN_SUBJECT="${ADMIN_SUBJECT:-admin@example.com}"
COMPOSE_FILE="${COMPOSE_FILE:-deploy/docker-compose.yml}"
KEEP_CAPTURE_STATE="${KEEP_CAPTURE_STATE:-false}"

usage() {
  cat <<'EOF'
Usage:
  scripts/capture-ui-tour.sh [--check]

Captures public-safe tutorial screenshots from the running local evaluator app.

Environment:
  BASE_URL       Local app URL. Default: http://localhost:8080
  OUT_DIR        Screenshot output directory. Default: site/assets/screenshots
  ADMIN_TOKEN    Optional pre-generated admin token. If absent, the script asks
                 the running agency-config container for a short-lived token.
  STATE_FILE     Temporary browser storage state. Default: .cache/ui-tour-storage.json
  KEEP_CAPTURE_STATE true|false; keep STATE_FILE after capture. Default: false.

Prerequisite:
  Start the local app first:
    make agency-app-up
EOF
}

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required tool: $1" >&2
    exit 1
  fi
}

if [ "${1:-}" = "--help" ] || [ "${1:-}" = "-h" ]; then
  usage
  exit 0
fi

need curl
need npx
need python3

case "$KEEP_CAPTURE_STATE" in
  true|false) ;;
  *) echo "KEEP_CAPTURE_STATE must be true or false" >&2; exit 2 ;;
esac

cleanup_state() {
  if [ "$KEEP_CAPTURE_STATE" != "true" ]; then
    rm -f "$STATE_FILE"
  fi
}
trap cleanup_state EXIT INT TERM

if [ "${1:-}" = "--check" ]; then
  curl -fsS "$BASE_URL/admin/local-login" >/dev/null
  npx playwright --version >/dev/null
  echo "capture prerequisites are available"
  exit 0
elif [ "${1:-}" != "" ]; then
  usage
  exit 2
fi

if ! curl -fsS "$BASE_URL/admin/local-login" >/dev/null; then
  echo "Local app is not reachable at $BASE_URL/admin/local-login. Run make agency-app-up first." >&2
  exit 1
fi

mkdir -p "$OUT_DIR" "$(dirname "$STATE_FILE")"

if [ -n "${ADMIN_TOKEN:-}" ]; then
  token="$ADMIN_TOKEN"
else
  need docker
  token="$(docker compose -f "$COMPOSE_FILE" --profile app exec -T agency-config /app/bin/admin-token -sub "$ADMIN_SUBJECT" -agency-id "$AGENCY_ID" | sed -n 's/^token=//p')"
fi
if [ -z "$token" ]; then
  echo "Could not obtain an admin token for screenshot capture." >&2
  exit 1
fi

STATE_FILE="$STATE_FILE" ADMIN_TOKEN="$token" python3 - <<'PY'
import json
import os
import time

state = {
    "cookies": [
        {
            "name": "admin_session",
            "value": os.environ["ADMIN_TOKEN"],
            "domain": "localhost",
            "path": "/admin",
            "expires": int(time.time()) + 3600,
            "httpOnly": True,
            "secure": False,
            "sameSite": "Lax",
        }
    ],
    "origins": [],
}
with open(os.environ["STATE_FILE"], "w", encoding="utf-8") as f:
    json.dump(state, f)
PY

capture_public() {
  path="$1"
  file="$2"
  selector="$3"
  npx playwright screenshot \
    --viewport-size 1440,1000 \
    --wait-for-selector "$selector" \
    "$BASE_URL$path" "$OUT_DIR/$file"
}

capture_admin() {
  path="$1"
  file="$2"
  npx playwright screenshot \
    --load-storage "$STATE_FILE" \
    --viewport-size 1440,1000 \
    --wait-for-selector "#operations-main" \
    "$BASE_URL$path" "$OUT_DIR/$file"
}

capture_public "/admin/local-login" "local-login.png" "#local-login-title"
capture_admin "/admin/operations" "start-page.png"
capture_admin "/admin/operations/setup-wizard" "agency-setup.png"
capture_admin "/admin/operations/gtfs-workbench" "gtfs-workbench.png"
capture_admin "/admin/operations/feed-health" "feed-health.png"
capture_admin "/admin/operations/devices" "devices.png"
capture_admin "/admin/operations/realtime" "realtime.png"
capture_admin "/admin/operations/connectors" "connectors.png"
capture_admin "/admin/operations/readiness" "readiness.png"
capture_admin "/admin/operations/maintenance" "maintenance.png"
capture_admin "/admin/operations/help" "help.png"

echo "captured UI tour screenshots in $OUT_DIR"
