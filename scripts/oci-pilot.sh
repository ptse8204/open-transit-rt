#!/usr/bin/env sh
# scripts/oci-pilot.sh — Oracle Cloud VM.Standard.E2.1.Micro pilot orchestration (Oracle Linux 9)
#
# Runs from the local Mac. Manages the full lifecycle of the OCI pilot deployment:
# build, push, setup, migrate, systemd services, DNS, and evidence collection.
#
# Usage:  scripts/oci-pilot.sh <subcommand> [args]
# Run:    scripts/oci-pilot.sh help   for a full list of subcommands.
#
# Required env for most commands:
#   OCI_HOST      — public IP of the OCI instance        (default: 192.9.142.92)
#   OCI_USER      — SSH username                         (default: opc)
#   OCI_KEY       — path to SSH private key              (default: use ssh-agent)
#
# Required env for DNS update:
#   DUCKDNS_TOKEN — DuckDNS API token
#
# Required env for evidence collection:
#   ADMIN_TOKEN   — admin JWT (or run: scripts/oci-pilot.sh token)

set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

# ---------------------------------------------------------------------------
# Configuration
# ---------------------------------------------------------------------------

OCI_HOST="${OCI_HOST:-192.9.142.92}"
OCI_USER="${OCI_USER:-opc}"          # Oracle Linux default; Ubuntu instances use 'ubuntu'
OCI_KEY="${OCI_KEY:-}"                          # empty = use ssh-agent / ~/.ssh/id_rsa
OCI_REMOTE_DIR="${OCI_REMOTE_DIR:-/opt/open-transit-rt}"
OCI_APP_USER="${OCI_APP_USER:-open-transit}"
DOMAIN="${DOMAIN:-open-transit-pilot.duckdns.org}"
DUCKDNS_DOMAIN="${DUCKDNS_DOMAIN:-open-transit-pilot}"
ENVIRONMENT_NAME="${ENVIRONMENT_NAME:-oci-pilot}"
PUBLIC_BASE_URL="https://${DOMAIN}"
BIN_DIR="$ROOT_DIR/deploy/bin"

SERVICES="agency-config telemetry-ingest feed-vehicle-positions feed-trip-updates feed-alerts"
SERVICE_PORTS="8081 8082 8083 8084 8085"

# ---------------------------------------------------------------------------
# SSH / SCP helpers
# ---------------------------------------------------------------------------

ssh_opts() {
  if [ -n "${OCI_KEY:-}" ]; then
    printf -- '-i %s' "$OCI_KEY"
  fi
}

# ssh_run <remote-command>
ssh_run() {
  if [ -n "${OCI_KEY:-}" ]; then
    ssh -i "$OCI_KEY" "${OCI_USER}@${OCI_HOST}" "$@"
  else
    ssh "${OCI_USER}@${OCI_HOST}" "$@"
  fi
}

# scp_to <local-src> <remote-dest>
scp_to() {
  local src="$1" dest="$2"
  if [ -n "${OCI_KEY:-}" ]; then
    scp -r -i "$OCI_KEY" "$src" "${OCI_USER}@${OCI_HOST}:${dest}"
  else
    scp -r "$src" "${OCI_USER}@${OCI_HOST}:${dest}"
  fi
}

# copy_dir_to_remote <local-dir> <remote-dir>
# Streams a directory through tar because OCI_REMOTE_DIR is owned by the service account.
copy_dir_to_remote() {
  local src_dir="${1%/}" dest_dir="$2"
  ssh_run "sudo mkdir -p '${dest_dir}'"
  if [ -n "${OCI_KEY:-}" ]; then
    COPYFILE_DISABLE=1 tar --disable-copyfile -C "$src_dir" -cf - . | ssh -i "$OCI_KEY" "${OCI_USER}@${OCI_HOST}" "sudo tar -C '${dest_dir}' -xf -"
  else
    COPYFILE_DISABLE=1 tar --disable-copyfile -C "$src_dir" -cf - . | ssh "${OCI_USER}@${OCI_HOST}" "sudo tar -C '${dest_dir}' -xf -"
  fi
}

# ---------------------------------------------------------------------------
# Subcommand: build
#   Cross-compile all service binaries for linux/amd64 into deploy/bin/
# ---------------------------------------------------------------------------

cmd_build() {
  echo "==> Building linux/amd64 binaries..."
  mkdir -p "$BIN_DIR"
  export GOOS=linux
  export GOARCH=amd64
  export CGO_ENABLED=0
  for svc in \
    agency-config \
    telemetry-ingest \
    feed-vehicle-positions \
    feed-trip-updates \
    feed-alerts \
    migrate \
    admin-token \
    gtfs-import
  do
    echo "    compiling cmd/${svc}..."
    go build -trimpath -ldflags "-s -w" \
      -o "${BIN_DIR}/${svc}" \
      "./cmd/${svc}"
  done
  echo "==> Binaries written to deploy/bin/:"
  ls -lh "$BIN_DIR"
}

# ---------------------------------------------------------------------------
# Subcommand: push
#   Upload binaries and deployment assets to the OCI instance.
#   Creates the remote directory structure and sets correct ownership.
# ---------------------------------------------------------------------------

cmd_push() {
  echo "==> Checking that deploy/bin/ is populated..."
  if [ ! -f "${BIN_DIR}/agency-config" ]; then
    echo "ERROR: deploy/bin/agency-config not found. Run: scripts/oci-pilot.sh build" >&2
    exit 1
  fi

  echo "==> Creating remote directory structure..."
  ssh_run "sudo mkdir -p ${OCI_REMOTE_DIR}/{bin,app/db/migrations,app/testdata,app/deploy/systemd,app/deploy/oci,.cache/validators,data}"
  ssh_run "sudo chown -R ${OCI_APP_USER}:${OCI_APP_USER} ${OCI_REMOTE_DIR}"
  ssh_run "sudo chmod 750 ${OCI_REMOTE_DIR}"

  echo "==> Uploading binaries..."
  copy_dir_to_remote "${BIN_DIR}" "${OCI_REMOTE_DIR}/bin"
  ssh_run "sudo sh -c 'chmod +x ${OCI_REMOTE_DIR}/bin/*'"

  echo "==> Uploading migrations..."
  copy_dir_to_remote "db/migrations" "${OCI_REMOTE_DIR}/app/db/migrations"

  echo "==> Uploading testdata (GTFS fixtures)..."
  copy_dir_to_remote "testdata" "${OCI_REMOTE_DIR}/app/testdata"

  echo "==> Uploading deploy assets (systemd units, OCI config)..."
  copy_dir_to_remote "deploy/systemd" "${OCI_REMOTE_DIR}/app/deploy/systemd"
  copy_dir_to_remote "deploy/oci"     "${OCI_REMOTE_DIR}/app/deploy/oci"

  echo "==> Re-chowning after upload..."
  ssh_run "sudo chown -R ${OCI_APP_USER}:${OCI_APP_USER} ${OCI_REMOTE_DIR}"

  echo "==> Push complete."
}

# ---------------------------------------------------------------------------
# Subcommand: setup
#   First-time instance setup: Oracle Linux packages, swap, PostgreSQL, Caddy, user.
#   Runs deploy/oci/setup-instance.sh on the remote host via sudo.
#   Safe to run more than once (idempotent steps).
# ---------------------------------------------------------------------------

cmd_setup() {
  echo "==> Uploading setup-instance.sh..."
  scp_to "deploy/oci/setup-instance.sh" "/tmp/oci-setup-instance.sh"
  ssh_run "chmod +x /tmp/oci-setup-instance.sh"

  echo "==> Running setup-instance.sh on remote host (requires sudo)..."
  ssh_run "sudo /tmp/oci-setup-instance.sh"
  echo "==> Instance setup complete."
}

# ---------------------------------------------------------------------------
# Subcommand: units
#   Install systemd unit files and enable all services.
#   Runs deploy/oci/install-units.sh on the remote host.
# ---------------------------------------------------------------------------

cmd_units() {
  echo "==> Uploading install-units.sh..."
  scp_to "deploy/oci/install-units.sh" "/tmp/oci-install-units.sh"
  ssh_run "chmod +x /tmp/oci-install-units.sh"

  echo "==> Installing systemd units on remote host..."
  ssh_run "sudo OCI_REMOTE_DIR=${OCI_REMOTE_DIR} OCI_APP_USER=${OCI_APP_USER} DOMAIN=${DOMAIN} /tmp/oci-install-units.sh"
  echo "==> Systemd units installed and enabled."
}

# ---------------------------------------------------------------------------
# Subcommand: env-init
#   Write an initial environment file on the OCI instance if one does not
#   exist yet. Secrets are generated randomly; operators must fill in
#   AGENCY_ID, TECHNICAL_CONTACT_EMAIL, and feed license fields.
# ---------------------------------------------------------------------------

cmd_env_init() {
  echo "==> Checking for existing env file on remote..."
  if ssh_run "test -f ${OCI_REMOTE_DIR}/env" 2>/dev/null; then
    echo "  env file already exists at ${OCI_REMOTE_DIR}/env — skipping."
    echo "  To regenerate, delete it first: ssh ${OCI_USER}@${OCI_HOST} sudo rm ${OCI_REMOTE_DIR}/env"
    return 0
  fi

  echo "==> Generating env file with random secrets..."
  _secret() { openssl rand -hex 32; }

  ADMIN_JWT_SECRET=$(_secret)
  CSRF_SECRET=$(_secret)
  DEVICE_TOKEN_PEPPER=$(_secret)
  DB_PASSWORD=$(_secret)

  ENV_CONTENT="# Open Transit RT — OCI pilot environment
# Generated by scripts/oci-pilot.sh env-init. Keep mode 600.
# TODO: fill in AGENCY_ID, TECHNICAL_CONTACT_EMAIL, FEED_LICENSE_* below.

DATABASE_URL=postgres://open_transit:${DB_PASSWORD}@127.0.0.1:5432/open_transit_rt?sslmode=disable
MIGRATIONS_DIR=${OCI_REMOTE_DIR}/app/db/migrations

APP_ENV=production
BIND_ADDR=127.0.0.1

# REQUIRED: set your agency identifier
AGENCY_ID=your-agency-id

PUBLIC_BASE_URL=${PUBLIC_BASE_URL}
FEED_BASE_URL=${PUBLIC_BASE_URL}/public
VEHICLE_POSITIONS_FEED_URL=${PUBLIC_BASE_URL}/public/gtfsrt/vehicle_positions.pb
REALTIME_VALIDATION_BASE_URL=${PUBLIC_BASE_URL}/public

# REQUIRED: set contact and license
TECHNICAL_CONTACT_EMAIL=ops@example.org
FEED_LICENSE_NAME=CC-BY-4.0
FEED_LICENSE_URL=https://creativecommons.org/licenses/by/4.0/
PUBLICATION_ENVIRONMENT=production

ADMIN_JWT_SECRET=${ADMIN_JWT_SECRET}
ADMIN_JWT_ISSUER=open-transit-rt-oci-pilot
ADMIN_JWT_AUDIENCE=open-transit-rt-admin
ADMIN_JWT_TTL=8h
CSRF_SECRET=${CSRF_SECRET}
DEVICE_TOKEN_PEPPER=${DEVICE_TOKEN_PEPPER}

METRICS_ENABLED=false
LOG_LEVEL=info

MATCH_CONFIDENCE_THRESHOLD=0.75
STALE_TELEMETRY_TTL_SECONDS=90
SUPPRESS_STALE_VEHICLE_AFTER_SECONDS=300
VEHICLE_POSITIONS_MAX_VEHICLES=2000
VEHICLE_POSITIONS_TRIP_CONFIDENCE_THRESHOLD=0.65
TRIP_UPDATES_ADAPTER=deterministic
TRIP_UPDATES_MAX_VEHICLES=2000
TRIP_UPDATES_STALE_TELEMETRY_TTL_SECONDS=90
TRIP_UPDATES_ASSIGNMENT_CONFIDENCE_THRESHOLD=0.65
TRIP_UPDATES_MAX_SCHEDULE_DEVIATION_SECONDS=2700

# Validator tooling (Java is installed by setup-instance.sh)
VALIDATOR_TOOLING_MODE=pinned
JAVA_BINARY=/usr/bin/java
GTFS_VALIDATOR_PATH=${OCI_REMOTE_DIR}/.cache/validators/gtfs-validator-7.1.0-cli.jar
GTFS_VALIDATOR_VERSION=v7.1.0
GTFS_RT_VALIDATOR_PATH=${OCI_REMOTE_DIR}/.cache/validators/gtfs-rt-validator-wrapper.sh
GTFS_RT_VALIDATOR_VERSION=ghcr.io/mobilitydata/gtfs-realtime-validator@sha256:5d2a3c14fba49983e1968c4a715e8ca624d4062bf4afede74aeca26322436c89
"

  # Write via stdin to avoid exposing secrets in process list
  printf '%s\n' "$ENV_CONTENT" \
    | ssh_run "sudo tee ${OCI_REMOTE_DIR}/env > /dev/null && \
               sudo chown ${OCI_APP_USER}:${OCI_APP_USER} ${OCI_REMOTE_DIR}/env && \
               sudo chmod 600 ${OCI_REMOTE_DIR}/env"

  echo "==> Synchronizing Postgres role password with generated env secret..."
  ssh_run "sudo -u postgres /usr/pgsql-15/bin/psql -v ON_ERROR_STOP=1 -c \
    \"DO \\\$\\\$ BEGIN
       IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'open_transit') THEN
         CREATE ROLE open_transit WITH LOGIN PASSWORD '${DB_PASSWORD}';
       ELSE
         ALTER ROLE open_transit WITH LOGIN PASSWORD '${DB_PASSWORD}';
       END IF;
     END \\\$\\\$;\""

  echo "==> Ensuring database exists..."
  ssh_run "sudo -u postgres /usr/pgsql-15/bin/psql -tc \
    \"SELECT 1 FROM pg_database WHERE datname='open_transit_rt'\" | grep -q 1 || \
    sudo -u postgres /usr/pgsql-15/bin/createdb -O open_transit open_transit_rt"

  echo "==> env file written at ${OCI_REMOTE_DIR}/env."
  echo "    Edit ${OCI_REMOTE_DIR}/env to fill in AGENCY_ID, TECHNICAL_CONTACT_EMAIL, etc."
}

# ---------------------------------------------------------------------------
# Subcommand: migrate
#   Run database migrations on the OCI instance.
# ---------------------------------------------------------------------------

cmd_migrate() {
  echo "==> Running migrations on OCI instance..."
  ssh_run "sudo -u ${OCI_APP_USER} sh -c \
    'set -a; . ${OCI_REMOTE_DIR}/env; set +a; \
     ${OCI_REMOTE_DIR}/bin/migrate up'"
  echo "==> Migration status:"
  ssh_run "sudo -u ${OCI_APP_USER} sh -c \
    'set -a; . ${OCI_REMOTE_DIR}/env; set +a; \
     ${OCI_REMOTE_DIR}/bin/migrate status'"
}

# ---------------------------------------------------------------------------
# Subcommand: start / stop / restart
# ---------------------------------------------------------------------------

cmd_start() {
  echo "==> Starting all Open Transit RT services..."
  for svc in $SERVICES; do
    ssh_run "sudo systemctl start open-transit-${svc}"
    echo "    started: open-transit-${svc}"
  done
}

cmd_stop() {
  echo "==> Stopping all Open Transit RT services..."
  for svc in $SERVICES; do
    ssh_run "sudo systemctl stop open-transit-${svc}" || true
    echo "    stopped: open-transit-${svc}"
  done
}

cmd_restart() {
  echo "==> Restarting all Open Transit RT services..."
  for svc in $SERVICES; do
    ssh_run "sudo systemctl restart open-transit-${svc}"
    echo "    restarted: open-transit-${svc}"
  done
}

# ---------------------------------------------------------------------------
# Subcommand: status
#   Check systemd status, local health endpoints, and public HTTPS feed URLs.
# ---------------------------------------------------------------------------

cmd_status() {
  echo "==> Systemd service status:"
  for svc in $SERVICES; do
    ssh_run "systemctl is-active open-transit-${svc} 2>/dev/null || echo inactive" \
      | awk -v s="open-transit-${svc}" '{printf "    %-45s %s\n", s, $0}'
  done

  echo ""
  echo "==> Local health endpoints (via SSH):"
  for port in $SERVICE_PORTS; do
    result=$(ssh_run "curl -s -o /dev/null -w '%{http_code}' http://127.0.0.1:${port}/healthz 2>/dev/null || echo 'ERR'")
    printf "    http://127.0.0.1:%s/healthz  =>  %s\n" "$port" "$result"
  done

  echo ""
  echo "==> Public HTTPS feed endpoints (from local machine):"
  for path in \
    "/public/feeds.json" \
    "/public/gtfs/schedule.zip" \
    "/public/gtfsrt/vehicle_positions.pb" \
    "/public/gtfsrt/trip_updates.pb" \
    "/public/gtfsrt/alerts.pb"
  do
    result=$(curl -s -o /dev/null -w '%{http_code}' \
      --connect-timeout 5 --max-time 10 \
      "${PUBLIC_BASE_URL}${path}" 2>/dev/null || echo "ERR")
    printf "    %s  %s%s\n" "$result" "$PUBLIC_BASE_URL" "$path"
  done

  echo ""
  echo "==> TLS certificate info:"
  echo | openssl s_client -connect "${DOMAIN}:443" 2>/dev/null \
    | openssl x509 -noout -dates -issuer 2>/dev/null \
    | sed 's/^/    /' || echo "    (TLS not yet available)"
}

# ---------------------------------------------------------------------------
# Subcommand: update-dns
#   Update the DuckDNS A record to point at the OCI instance's public IP.
#   Requires DUCKDNS_TOKEN env var.
# ---------------------------------------------------------------------------

cmd_update_dns() {
  if [ -z "${DUCKDNS_TOKEN:-}" ]; then
    echo "ERROR: DUCKDNS_TOKEN is not set." >&2
    echo "  Export it: export DUCKDNS_TOKEN=<your-token>" >&2
    exit 1
  fi
  echo "==> Updating DuckDNS record for ${DUCKDNS_DOMAIN} to ${OCI_HOST}..."
  RESPONSE=$(curl -s \
    "https://www.duckdns.org/update?domains=${DUCKDNS_DOMAIN}&token=${DUCKDNS_TOKEN}&ip=${OCI_HOST}")
  echo "    DuckDNS response: ${RESPONSE}"
  if [ "$RESPONSE" = "OK" ]; then
    echo "    DNS updated. Verify with: dig +short ${DOMAIN} A"
    dig +short "${DOMAIN}" A 2>/dev/null || true
  else
    echo "ERROR: DuckDNS update failed." >&2
    exit 1
  fi
}

# ---------------------------------------------------------------------------
# Subcommand: token
#   Generate an admin JWT on the OCI instance and print it locally.
#   Saves the token to .cache/oci-admin-token on this machine.
# ---------------------------------------------------------------------------

cmd_token() {
  AGENCY_ID="${1:-}"
  if [ -z "$AGENCY_ID" ]; then
    # Try to read from the remote env file
    AGENCY_ID=$(ssh_run "sudo -u ${OCI_APP_USER} grep '^AGENCY_ID=' ${OCI_REMOTE_DIR}/env | cut -d= -f2" 2>/dev/null || echo "")
  fi
  if [ -z "$AGENCY_ID" ]; then
    echo "ERROR: could not determine AGENCY_ID. Pass it as an argument:" >&2
    echo "  scripts/oci-pilot.sh token <agency-id>" >&2
    exit 1
  fi

  echo "==> Generating admin token for agency: ${AGENCY_ID}..."
  TOKEN=$(ssh_run "sudo -u ${OCI_APP_USER} sh -c \
    'set -a; . ${OCI_REMOTE_DIR}/env; set +a; \
     ${OCI_REMOTE_DIR}/bin/admin-token -sub admin@example.com -agency-id ${AGENCY_ID}' \
    | sed -n 's/^token=//p'")
  if [ -z "$TOKEN" ]; then
    echo "ERROR: token generation returned empty output." >&2
    exit 1
  fi
  mkdir -p "$ROOT_DIR/.cache"
  printf '%s\n' "$TOKEN" > "$ROOT_DIR/.cache/oci-admin-token"
  chmod 600 "$ROOT_DIR/.cache/oci-admin-token"
  echo "    Token written to .cache/oci-admin-token"
  echo "    Export for use: export ADMIN_TOKEN=\$(cat .cache/oci-admin-token)"
}

# ---------------------------------------------------------------------------
# Subcommand: bootstrap
#   Run the publication metadata bootstrap API call.
#   Requires ADMIN_TOKEN to be set (run 'scripts/oci-pilot.sh token' first).
# ---------------------------------------------------------------------------

cmd_bootstrap() {
  ADMIN_TOKEN="${ADMIN_TOKEN:-$(cat "$ROOT_DIR/.cache/oci-admin-token" 2>/dev/null || echo '')}"
  if [ -z "$ADMIN_TOKEN" ]; then
    echo "ERROR: ADMIN_TOKEN is not set. Run: scripts/oci-pilot.sh token" >&2
    exit 1
  fi
  echo "==> Running publication metadata bootstrap over SSH loopback..."
  # The public Caddy edge intentionally exposes anonymous feed paths only.
  ssh_run "curl -fsS -X POST http://127.0.0.1:8081/admin/publication/bootstrap \
    -H 'Authorization: Bearer ${ADMIN_TOKEN}' \
    -H 'Content-Type: application/json' \
    --data '{}'" | python3 -m json.tool
  echo ""
  echo "==> Bootstrap complete. Verify feeds.json:"
  curl -s "${PUBLIC_BASE_URL}/public/feeds.json" | python3 -m json.tool
}

# ---------------------------------------------------------------------------
# Subcommand: collect
#   Run the repo's hosted evidence collector against the live OCI deployment.
# ---------------------------------------------------------------------------

cmd_collect() {
  ADMIN_TOKEN="${ADMIN_TOKEN:-$(cat "$ROOT_DIR/.cache/oci-admin-token" 2>/dev/null || echo '')}"
  if [ -z "$ADMIN_TOKEN" ]; then
    echo "WARN: ADMIN_TOKEN is not set — admin-authenticated evidence steps will be skipped." >&2
    echo "  Run: scripts/oci-pilot.sh token   to generate a token first." >&2
  fi

  echo "==> Collecting hosted evidence for environment: ${ENVIRONMENT_NAME}..."
  ENVIRONMENT_NAME="$ENVIRONMENT_NAME" \
  PUBLIC_BASE_URL="$PUBLIC_BASE_URL" \
  ADMIN_BASE_URL="$PUBLIC_BASE_URL" \
  ADMIN_TOKEN="${ADMIN_TOKEN:-}" \
    ./scripts/collect-hosted-evidence.sh

  PACKET_DIR="docs/evidence/captured/${ENVIRONMENT_NAME}/$(date -u +%Y-%m-%d)"
  echo ""
  echo "==> Auditing collected packet: ${PACKET_DIR}"
  if [ -d "$PACKET_DIR" ]; then
    EVIDENCE_PACKET_DIR="$PACKET_DIR" ./scripts/audit-hosted-evidence.sh || true
  else
    echo "    Packet directory not found: ${PACKET_DIR}"
  fi
}

# ---------------------------------------------------------------------------
# Subcommand: logs
#   Tail recent logs from all services on the OCI instance.
# ---------------------------------------------------------------------------

cmd_logs() {
  echo "==> Tailing service logs (Ctrl+C to stop)..."
  UNITS=$(printf ' open-transit-%s' $SERVICES)
  # shellcheck disable=SC2086
  ssh_run "sudo journalctl -f -n 50 -u caddy $UNITS"
}

# ---------------------------------------------------------------------------
# Subcommand: deploy
#   Convenience: push + migrate + restart + status.
# ---------------------------------------------------------------------------

cmd_deploy() {
  cmd_push
  cmd_migrate
  cmd_restart
  cmd_status
}

# ---------------------------------------------------------------------------
# Subcommand: help
# ---------------------------------------------------------------------------

cmd_help() {
  cat <<'EOF'
scripts/oci-pilot.sh <subcommand>

FIRST-TIME SETUP (run in this order):
  update-dns   Update DuckDNS A record to OCI_HOST           (needs DUCKDNS_TOKEN)
  build        Cross-compile linux/amd64 binaries -> deploy/bin/
  setup        First-time instance setup (OL9 packages, swap, postgres, caddy) via SSH
  push         Upload binaries + migrations + config to OCI instance
  units        Install and enable systemd unit files on OCI instance
  env-init     Write initial env file on OCI instance (generates secrets)
  migrate      Run database migrations on OCI instance
  start        Start all Open Transit RT services
  token        Generate admin JWT and save to .cache/oci-admin-token
  bootstrap    POST /admin/publication/bootstrap (needs token)

DAY-TO-DAY:
  deploy       push + migrate + restart + status
  start        Start all services
  stop         Stop all services
  restart      Restart all services
  status       Health check: systemd + local + public HTTPS endpoints + TLS cert

EVIDENCE:
  collect      Run collect-hosted-evidence.sh against live OCI deployment

DIAGNOSTICS:
  logs         Tail journalctl for all services + Caddy

ENVIRONMENT:
  OCI_HOST          (default: 192.9.142.92)
  OCI_USER          (default: opc)
  OCI_KEY           SSH key path (default: ssh-agent)
  DUCKDNS_TOKEN     Required for update-dns
  ADMIN_TOKEN       Required for bootstrap / collect (or use 'token' subcommand)
  DOMAIN            (default: open-transit-pilot.duckdns.org)
EOF
}

# ---------------------------------------------------------------------------
# Dispatch
# ---------------------------------------------------------------------------

case "${1:-help}" in
  build)       cmd_build ;;
  push)        cmd_push ;;
  setup)       cmd_setup ;;
  units)       cmd_units ;;
  env-init)    cmd_env_init ;;
  migrate)     cmd_migrate ;;
  start)       cmd_start ;;
  stop)        cmd_stop ;;
  restart)     cmd_restart ;;
  status)      cmd_status ;;
  update-dns)  cmd_update_dns ;;
  token)       shift; cmd_token "${1:-}" ;;
  bootstrap)   cmd_bootstrap ;;
  collect)     cmd_collect ;;
  logs)        cmd_logs ;;
  deploy)      cmd_deploy ;;
  help|--help|-h) cmd_help ;;
  *)
    echo "Unknown subcommand: ${1}" >&2
    echo "Run: scripts/oci-pilot.sh help" >&2
    exit 2
    ;;
esac
