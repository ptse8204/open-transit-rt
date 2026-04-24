#!/usr/bin/env sh
# deploy/oci/install-units.sh — Install and enable Open Transit RT systemd units
#
# Run via: scripts/oci-pilot.sh units
# Or manually on the OCI instance:
#   sudo OCI_REMOTE_DIR=/opt/open-transit-rt OCI_APP_USER=open-transit DOMAIN=open-transit-pilot.duckdns.org \
#     ./install-units.sh
#
# Idempotent: can be run after code updates to refresh unit definitions.

set -eu

OCI_REMOTE_DIR="${OCI_REMOTE_DIR:-/opt/open-transit-rt}"
OCI_APP_USER="${OCI_APP_USER:-open-transit}"
DOMAIN="${DOMAIN:-open-transit-pilot.duckdns.org}"
UNIT_SRC="${OCI_REMOTE_DIR}/app/deploy/systemd"

echo "==> Installing Open Transit RT systemd units..."
echo "    Source: ${UNIT_SRC}"
echo "    Remote dir: ${OCI_REMOTE_DIR}"
echo "    App user: ${OCI_APP_USER}"

# ---------------------------------------------------------------------------
# 1. Install all service unit files
# ---------------------------------------------------------------------------

for unit_file in \
  open-transit-agency-config.service \
  open-transit-telemetry-ingest.service \
  open-transit-feed-vehicle-positions.service \
  open-transit-feed-trip-updates.service \
  open-transit-feed-alerts.service
do
  src="${UNIT_SRC}/${unit_file}"
  dest="/etc/systemd/system/${unit_file}"

  if [ ! -f "$src" ]; then
    echo "  ERROR: unit source not found: ${src}" >&2
    echo "  Run 'scripts/oci-pilot.sh push' first." >&2
    exit 1
  fi

  # Substitute placeholders before installing
  sed \
    -e "s|{{OCI_REMOTE_DIR}}|${OCI_REMOTE_DIR}|g" \
    -e "s|{{OCI_APP_USER}}|${OCI_APP_USER}|g" \
    -e "s|{{DOMAIN}}|${DOMAIN}|g" \
    "$src" > "$dest"

  echo "    installed: ${dest}"
done

# ---------------------------------------------------------------------------
# 2. Install and configure Caddyfile
# ---------------------------------------------------------------------------

CADDYFILE_SRC="${OCI_REMOTE_DIR}/app/deploy/oci/Caddyfile"
if [ -f "$CADDYFILE_SRC" ]; then
  sed \
    -e "s|{{DOMAIN}}|${DOMAIN}|g" \
    "$CADDYFILE_SRC" > /etc/caddy/Caddyfile
  echo "    installed: /etc/caddy/Caddyfile (domain: ${DOMAIN})"
else
  echo "  WARN: Caddyfile source not found at ${CADDYFILE_SRC} — skipping Caddy config." >&2
fi

# ---------------------------------------------------------------------------
# 3. Reload systemd and enable all services
# ---------------------------------------------------------------------------

systemctl daemon-reload
echo "==> Systemd daemon reloaded."

for svc in \
  open-transit-agency-config \
  open-transit-telemetry-ingest \
  open-transit-feed-vehicle-positions \
  open-transit-feed-trip-updates \
  open-transit-feed-alerts
do
  systemctl enable "${svc}"
  echo "    enabled: ${svc}"
done

# ---------------------------------------------------------------------------
# 4. Reload Caddy (apply Caddyfile)
# ---------------------------------------------------------------------------

if systemctl is-active caddy > /dev/null 2>&1; then
  systemctl reload caddy
  echo "    caddy reloaded with new Caddyfile"
else
  echo "    caddy is not running — start it with: systemctl start caddy"
fi

echo ""
echo "==> Units installed. Services are enabled but NOT started."
echo "    To start: scripts/oci-pilot.sh start"
echo "    Or:       sudo systemctl start open-transit-agency-config ..."
