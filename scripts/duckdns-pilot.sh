#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

DOMAIN="${DUCKDNS_DOMAIN:-open-transit-pilot.duckdns.org}"
ENVIRONMENT_NAME="${ENVIRONMENT_NAME:-duckdns-pilot}"
PUBLIC_BASE_URL="${PUBLIC_BASE_URL:-https://$DOMAIN}"
ADMIN_BASE_URL="${ADMIN_BASE_URL:-http://127.0.0.1:8081}"
WORK_DIR="${DUCKDNS_PILOT_WORK_DIR:-$ROOT_DIR/.cache/duckdns-pilot}"
LOG_DIR="$WORK_DIR/logs"
RUN_DIR="$WORK_DIR/run"
ENV_FILE="$WORK_DIR/env"
CADDYFILE="$WORK_DIR/Caddyfile"
PUBLICATION_ENVIRONMENT="${PUBLICATION_ENVIRONMENT:-production}"

mkdir -p "$WORK_DIR" "$LOG_DIR" "$RUN_DIR"

random_secret() {
  if command -v openssl >/dev/null 2>&1; then
    openssl rand -hex 32
  else
    date +%s | shasum -a 256 | awk '{print $1}'
  fi
}

java_binary() {
  if [ -n "${JAVA_BINARY:-}" ] && [ -x "$JAVA_BINARY" ]; then
    printf '%s\n' "$JAVA_BINARY"
    return
  fi
  for candidate in \
    /usr/local/opt/openjdk@17/bin/java \
    /opt/homebrew/opt/openjdk@17/bin/java \
    /usr/local/opt/openjdk/bin/java \
    /opt/homebrew/opt/openjdk/bin/java
  do
    if [ -x "$candidate" ]; then
      printf '%s\n' "$candidate"
      return
    fi
  done
  if command -v java >/dev/null 2>&1; then
    command -v java
  fi
}

write_env() {
  if [ -f "$ENV_FILE" ]; then
    return
  fi
  cat >"$ENV_FILE" <<EOF
DATABASE_URL=postgres://postgres:postgres@localhost:55432/open_transit_rt?sslmode=disable
MIGRATIONS_DIR=db/migrations
APP_ENV=production
BIND_ADDR=127.0.0.1
AGENCY_ID=demo-agency
ADMIN_JWT_SECRET=$(random_secret)
ADMIN_JWT_ISSUER=open-transit-rt-duckdns-pilot
ADMIN_JWT_AUDIENCE=open-transit-rt-admin
ADMIN_JWT_TTL=8h
CSRF_SECRET=$(random_secret)
DEVICE_TOKEN_PEPPER=$(random_secret)
PUBLIC_BASE_URL=$PUBLIC_BASE_URL
FEED_BASE_URL=$PUBLIC_BASE_URL/public
REALTIME_VALIDATION_BASE_URL=$PUBLIC_BASE_URL/public
TECHNICAL_CONTACT_EMAIL='${TECHNICAL_CONTACT_EMAIL:-ops@example.org}'
FEED_LICENSE_NAME='${FEED_LICENSE_NAME:-CC BY 4.0}'
FEED_LICENSE_URL='${FEED_LICENSE_URL:-https://creativecommons.org/licenses/by/4.0/}'
PUBLICATION_ENVIRONMENT=$PUBLICATION_ENVIRONMENT
VALIDATOR_TOOLING_MODE=pinned
GTFS_VALIDATOR_PATH=$ROOT_DIR/.cache/validators/gtfs-validator-7.1.0-cli.jar
GTFS_VALIDATOR_VERSION=v7.1.0
GTFS_RT_VALIDATOR_PATH=$ROOT_DIR/.cache/validators/gtfs-rt-validator-wrapper.sh
GTFS_RT_VALIDATOR_VERSION=ghcr.io/mobilitydata/gtfs-realtime-validator@sha256:5d2a3c14fba49983e1968c4a715e8ca624d4062bf4afede74aeca26322436c89
JAVA_BINARY=$(java_binary)
METRICS_ENABLED=false
EOF
  chmod 600 "$ENV_FILE"
}

load_env() {
  set -a
  # shellcheck disable=SC1090
  . "$ENV_FILE"
  set +a
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
    if [ "$attempt" -ge 60 ]; then
      echo "timed out waiting for $url to return $want; last status was $status" >&2
      return 1
    fi
    sleep 1
  done
}

wait_for_public_https() {
  attempt=0
  while :; do
    if curl --connect-timeout 5 --max-time 10 -fsS "$PUBLIC_BASE_URL/public/feeds.json" >/dev/null 2>&1; then
      return 0
    fi
    attempt=$((attempt + 1))
    if [ "$attempt" -ge 60 ]; then
      echo "public HTTPS feed root is not reachable yet after waiting: $PUBLIC_BASE_URL" >&2
      return 1
    fi
    sleep 2
  done
}

stop_pid() {
  pid_file="$1"
  if [ ! -f "$pid_file" ]; then
    return
  fi
  pid="$(cat "$pid_file")"
  if [ -n "$pid" ] && kill "$pid" >/dev/null 2>&1; then
    kill "$pid" >/dev/null 2>&1 || true
    wait "$pid" >/dev/null 2>&1 || true
  fi
  rm -f "$pid_file"
}

stop_all() {
  for name in caddy gtfs-studio feed-alerts feed-trip-updates feed-vehicle-positions telemetry-ingest agency-config
  do
    stop_pid "$RUN_DIR/$name.pid"
  done
}

start_service() {
  name="$1"
  port="$2"
  shift 2
  stop_pid "$RUN_DIR/$name.pid"
  nohup sh -c '. "$1"; export PORT="$2"; shift 2; exec "$@"' \
    duckdns-pilot "$ENV_FILE" "$port" "$@" \
    >"$LOG_DIR/$name.log" 2>&1 </dev/null &
  echo "$!" >"$RUN_DIR/$name.pid"
  wait_for_status "http://127.0.0.1:$port/healthz" "200"
  echo "$name listening on 127.0.0.1:$port"
}

write_caddyfile() {
  cat >"$CADDYFILE" <<EOF
$DOMAIN {
  encode zstd gzip

  handle /public/gtfs/* {
    reverse_proxy 127.0.0.1:8081
  }

  handle /public/feeds.json {
    reverse_proxy 127.0.0.1:8081
  }

  handle /public/gtfsrt/vehicle_positions.pb {
    reverse_proxy 127.0.0.1:8083
  }

  handle /public/gtfsrt/trip_updates.pb {
    reverse_proxy 127.0.0.1:8084
  }

  handle /public/gtfsrt/alerts.pb {
    reverse_proxy 127.0.0.1:8085
  }

  handle {
    respond 404
  }
}
EOF
}

start_caddy() {
  if ! command -v caddy >/dev/null 2>&1; then
    echo "missing caddy; install it with: brew install caddy" >&2
    exit 1
  fi
  stop_pid "$RUN_DIR/caddy.pid"
  env XDG_DATA_HOME="$WORK_DIR/caddy-data" HOME="$WORK_DIR/caddy-home" \
    caddy start --config "$CADDYFILE" --adapter caddyfile --pidfile "$RUN_DIR/caddy.pid" \
    >"$LOG_DIR/caddy.log" 2>&1
  sleep 3
  if ! kill "$(cat "$RUN_DIR/caddy.pid")" >/dev/null 2>&1; then
    echo "caddy failed to start; see $LOG_DIR/caddy.log" >&2
    exit 1
  fi
  echo "caddy started for $DOMAIN"
}

bootstrap_data() {
  make db-up
  until docker compose -f deploy/docker-compose.yml exec -T postgres pg_isready -U postgres -d open_transit_rt >/dev/null 2>&1; do
    sleep 1
  done
  make migrate-up
  make seed
  make validators-install
  make validators-check

  gtfs_zip="$WORK_DIR/valid-small.zip"
  (cd testdata/gtfs/valid-small && zip -qr "$gtfs_zip" .)
  go run ./cmd/gtfs-import -agency-id "$AGENCY_ID" -zip "$gtfs_zip" -actor-id duckdns-pilot -notes "DuckDNS pilot bootstrap" >"$WORK_DIR/import-result.json"
}

post_start_data() {
  admin_token="$(go run ./cmd/admin-token -sub admin@example.com -agency-id "$AGENCY_ID" | sed -n 's/^token=//p')"
  printf '%s\n' "$admin_token" >"$WORK_DIR/admin-token"
  chmod 600 "$WORK_DIR/admin-token"

  curl -fsS -X POST "$ADMIN_BASE_URL/admin/publication/bootstrap" \
    -H "Authorization: Bearer $admin_token" \
    -H "Content-Type: application/json" \
    --data '{}' >"$WORK_DIR/publication-bootstrap.json"

  curl -fsS -X POST "$ADMIN_BASE_URL/admin/devices/rebind" \
    -H "Authorization: Bearer $admin_token" \
    -H "Content-Type: application/json" \
    --data '{
      "agency_id": "demo-agency",
      "device_id": "device-1",
      "vehicle_id": "bus-1",
      "reason": "DuckDNS pilot bootstrap token rotation"
    }' >"$WORK_DIR/device-rebind.json"
  device_token="$(sed -n 's/.*"token"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' "$WORK_DIR/device-rebind.json")"
  if [ -z "$device_token" ]; then
    echo "device rebind did not return a token" >&2
    exit 1
  fi
  printf '%s\n' "$device_token" >"$WORK_DIR/device-token"
  chmod 600 "$WORK_DIR/device-token"

  observed_at="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
  cat >"$WORK_DIR/telemetry.json" <<EOF
{
  "agency_id": "$AGENCY_ID",
  "device_id": "device-1",
  "vehicle_id": "bus-1",
  "timestamp": "$observed_at",
  "lat": 49.2827,
  "lon": -123.1207,
  "bearing": 120.0,
  "speed_mps": 8.4,
  "accuracy_m": 7.5,
  "trip_hint": "trip-10-0800"
}
EOF
  curl -fsS -X POST "http://127.0.0.1:8082/v1/telemetry" \
    -H "Authorization: Bearer $device_token" \
    -H "Content-Type: application/json" \
    --data @"$WORK_DIR/telemetry.json" >"$WORK_DIR/telemetry-response.json"
}

collect_if_reachable() {
  if wait_for_public_https; then
    ADMIN_TOKEN="$(cat "$WORK_DIR/admin-token")" \
    ENVIRONMENT_NAME="$ENVIRONMENT_NAME" \
    PUBLIC_BASE_URL="$PUBLIC_BASE_URL" \
    ADMIN_BASE_URL="$ADMIN_BASE_URL" \
      make collect-hosted-evidence
  else
    echo "public HTTPS feed root is not reachable yet: $PUBLIC_BASE_URL" >&2
    echo "Forward TCP 80 and 443 on the router/firewall to this machine, then rerun:" >&2
    echo "  scripts/duckdns-pilot.sh collect" >&2
    return 2
  fi
}

status() {
  echo "Domain: $DOMAIN"
  echo "Public base URL: $PUBLIC_BASE_URL"
  echo "Work dir: $WORK_DIR"
  if [ -f "$RUN_DIR/caddy.pid" ]; then
    caddy_pid="$(cat "$RUN_DIR/caddy.pid")"
    if kill "$caddy_pid" >/dev/null 2>&1; then
      echo "caddy pid: $caddy_pid running"
    else
      echo "caddy pid: $caddy_pid stopped"
    fi
  else
    echo "caddy pid: none"
  fi
  dig +short "$DOMAIN" A || true
  for port in 8081 8082 8083 8084 8085 8086
  do
    curl -s -o /dev/null -w "localhost:$port %{http_code}\n" "http://127.0.0.1:$port/healthz" || true
  done
  curl --connect-timeout 5 --max-time 10 -sS -D - -o /dev/null "$PUBLIC_BASE_URL/public/feeds.json" | sed -n '1,12p' || true
}

case "${1:-start}" in
  start)
    write_env
    load_env
    bootstrap_data
    start_service agency-config 8081 env go run ./cmd/agency-config
    start_service telemetry-ingest 8082 env go run ./cmd/telemetry-ingest
    start_service feed-vehicle-positions 8083 env go run ./cmd/feed-vehicle-positions
    start_service feed-trip-updates 8084 env VEHICLE_POSITIONS_FEED_URL="$FEED_BASE_URL/gtfsrt/vehicle_positions.pb" go run ./cmd/feed-trip-updates
    start_service feed-alerts 8085 env go run ./cmd/feed-alerts
    start_service gtfs-studio 8086 env go run ./cmd/gtfs-studio
    post_start_data
    write_caddyfile
    start_caddy
    status
    collect_if_reachable || true
    ;;
  collect)
    write_env
    load_env
    collect_if_reachable
    ;;
  proxy)
    write_env
    load_env
    write_caddyfile
    exec env XDG_DATA_HOME="$WORK_DIR/caddy-data" HOME="$WORK_DIR/caddy-home" \
      caddy run --config "$CADDYFILE" --adapter caddyfile
    ;;
  status)
    write_env
    load_env
    status
    ;;
  stop)
    stop_all
    ;;
  *)
    echo "usage: $0 [start|collect|proxy|status|stop]" >&2
    exit 2
    ;;
esac
