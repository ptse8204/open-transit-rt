#!/usr/bin/env sh
# deploy/oci/setup-instance.sh — First-time Oracle Cloud Ubuntu instance setup
#
# Run via: scripts/oci-pilot.sh setup
# Or manually: sudo ./setup-instance.sh
#
# This script is IDEMPOTENT — safe to run more than once.
# It does NOT write the application env file or secrets.
# After this script, the operator must run:
#   scripts/oci-pilot.sh env-init   (generates env file with secrets)

set -eu

echo "==> open-transit-rt OCI instance setup"
echo "    $(date -u)"

# ---------------------------------------------------------------------------
# 1. Swap (4 GB — critical for 1 GB RAM)
# ---------------------------------------------------------------------------

if swapon --show | grep -q /swapfile; then
  echo "==> Swap already enabled — skipping swap creation"
else
  echo "==> Creating 4 GB swap file..."
  fallocate -l 4G /swapfile
  chmod 0600 /swapfile
  mkswap /swapfile
  swapon /swapfile
  if ! grep -q '/swapfile' /etc/fstab; then
    echo '/swapfile none swap sw 0 0' >> /etc/fstab
  fi
  echo "vm.swappiness=10" >> /etc/sysctl.conf
  sysctl -p
  echo "==> Swap enabled: $(free -h | grep Swap)"
fi

# ---------------------------------------------------------------------------
# 2. iptables — open ports 80 and 443 persistently
# ---------------------------------------------------------------------------

echo "==> Configuring iptables for ports 80 and 443..."
# Add rules only if not already present
iptables -C INPUT -p tcp -m state --state NEW -m tcp --dport 80  -j ACCEPT 2>/dev/null \
  || iptables -I INPUT 6 -p tcp -m state --state NEW -m tcp --dport 80  -j ACCEPT
iptables -C INPUT -p tcp -m state --state NEW -m tcp --dport 443 -j ACCEPT 2>/dev/null \
  || iptables -I INPUT 6 -p tcp -m state --state NEW -m tcp --dport 443 -j ACCEPT

apt-get install -y -q iptables-persistent 2>/dev/null || true
netfilter-persistent save 2>/dev/null || iptables-save > /etc/iptables/rules.v4 || true
echo "==> iptables rules saved."

# ---------------------------------------------------------------------------
# 3. System packages
# ---------------------------------------------------------------------------

echo "==> Installing system packages..."
apt-get update -q
apt-get install -y -q \
  curl wget unzip zip git ca-certificates gnupg \
  openjdk-17-jre-headless \
  python3 \
  postgresql postgresql-contrib

# PostGIS: add PostgreSQL apt repo which has the postgis packages
if ! dpkg -l | grep -q postgresql-14-postgis-3 2>/dev/null && \
   ! dpkg -l | grep -q postgresql-16-postgis-3 2>/dev/null; then
  PG_VER=$(pg_lsclusters -h 2>/dev/null | awk '{print $1}' | head -1 || echo "16")
  apt-get install -y -q "postgresql-${PG_VER}-postgis-3" || \
    apt-get install -y -q postgis || true
fi

echo "==> System packages installed."

# ---------------------------------------------------------------------------
# 4. Caddy (from official repo)
# ---------------------------------------------------------------------------

if command -v caddy > /dev/null 2>&1; then
  echo "==> Caddy already installed: $(caddy version)"
else
  echo "==> Installing Caddy from official repository..."
  curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' \
    | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
  curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' \
    | tee /etc/apt/sources.list.d/caddy-stable.list
  apt-get update -q
  apt-get install -y -q caddy
  echo "==> Caddy installed: $(caddy version)"
fi

# ---------------------------------------------------------------------------
# 5. Go (match the version in go.mod — 1.23.x)
# ---------------------------------------------------------------------------

if command -v go > /dev/null 2>&1; then
  echo "==> Go already installed: $(go version)"
else
  GO_VERSION="1.23.8"
  echo "==> Installing Go ${GO_VERSION}..."
  wget -q "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz" -O /tmp/go.tar.gz
  tar -C /usr/local -xzf /tmp/go.tar.gz
  rm /tmp/go.tar.gz
  if ! grep -q '/usr/local/go/bin' /etc/profile.d/go.sh 2>/dev/null; then
    echo 'export PATH=$PATH:/usr/local/go/bin' > /etc/profile.d/go.sh
  fi
  echo "==> Go installed: $(/usr/local/go/bin/go version)"
fi

# ---------------------------------------------------------------------------
# 6. Create application system user
# ---------------------------------------------------------------------------

if id open-transit > /dev/null 2>&1; then
  echo "==> System user 'open-transit' already exists."
else
  echo "==> Creating system user 'open-transit'..."
  useradd -r -m -d /opt/open-transit-rt -s /sbin/nologin open-transit
fi

# ---------------------------------------------------------------------------
# 7. PostgreSQL — configure for 1 GB RAM + PostGIS
# ---------------------------------------------------------------------------

PG_VER=$(pg_lsclusters -h 2>/dev/null | awk '{print $1}' | head -1 || echo "16")
PG_CONF_DIR="/etc/postgresql/${PG_VER}/main"
TUNING_CONF="${PG_CONF_DIR}/conf.d/open-transit-rt.conf"

if [ ! -f "$TUNING_CONF" ]; then
  echo "==> Applying PostgreSQL memory tuning..."
  mkdir -p "${PG_CONF_DIR}/conf.d"
  cat > "$TUNING_CONF" <<'PGEOF'
# Open Transit RT — VM.Standard.E2.1.Micro tuning (1 GB RAM)
shared_buffers              = 128MB
effective_cache_size        = 512MB
work_mem                    = 4MB
maintenance_work_mem        = 32MB
max_connections             = 25
wal_buffers                 = 4MB
checkpoint_completion_target = 0.9
random_page_cost            = 1.1
log_min_duration_statement  = 500
PGEOF
  echo "==> PostgreSQL tuning written to ${TUNING_CONF}"
fi

# Bind PostgreSQL to loopback only
if grep -q "listen_addresses = '\*'" "${PG_CONF_DIR}/postgresql.conf" 2>/dev/null; then
  sed -i "s/listen_addresses = '\*'/listen_addresses = '127.0.0.1'/" \
    "${PG_CONF_DIR}/postgresql.conf"
fi

systemctl restart postgresql
systemctl enable postgresql
echo "==> PostgreSQL running: $(systemctl is-active postgresql)"

# ---------------------------------------------------------------------------
# 8. Create Postgres DB and user (idempotent)
# ---------------------------------------------------------------------------

# We cannot generate the password here (it must match the env file).
# Print instructions and create the DB structure without a known password first.
# The operator will set the password when they run env-init and then re-run this step.
echo "==> Ensuring database 'open_transit_rt' and role 'open_transit' exist..."
sudo -u postgres psql -tc "SELECT 1 FROM pg_roles WHERE rolname='open_transit'" \
  | grep -q 1 || \
  sudo -u postgres psql -c "CREATE ROLE open_transit WITH LOGIN PASSWORD 'changeme-run-env-init';"

sudo -u postgres psql -tc "SELECT 1 FROM pg_database WHERE datname='open_transit_rt'" \
  | grep -q 1 || \
  sudo -u postgres psql -c "CREATE DATABASE open_transit_rt OWNER open_transit;"

sudo -u postgres psql -d open_transit_rt -c "CREATE EXTENSION IF NOT EXISTS postgis;" 2>/dev/null || true
sudo -u postgres psql -d open_transit_rt -c "CREATE EXTENSION IF NOT EXISTS postgis_topology;" 2>/dev/null || true
sudo -u postgres psql -c "GRANT ALL ON DATABASE open_transit_rt TO open_transit;" 2>/dev/null || true

echo ""
echo "==> IMPORTANT: after running 'scripts/oci-pilot.sh env-init', update the DB password:"
echo "    sudo -u postgres psql -c \"ALTER ROLE open_transit WITH PASSWORD '<password-from-env-file>';\""

# ---------------------------------------------------------------------------
# 9. Caddy systemd enable
# ---------------------------------------------------------------------------

systemctl enable caddy
echo "==> Caddy enabled."

# ---------------------------------------------------------------------------
# 10. Docker (optional — for GTFS-RT validator only, not running as a daemon)
# ---------------------------------------------------------------------------

if command -v docker > /dev/null 2>&1; then
  echo "==> Docker already installed."
else
  echo "==> Installing Docker Engine (for GTFS-RT validator only)..."
  curl -fsSL https://get.docker.com | sh
  usermod -aG docker open-transit || true
  # Do NOT auto-start Docker on boot — only start it when running validation
  systemctl disable docker    2>/dev/null || true
  systemctl disable containerd 2>/dev/null || true
  echo "==> Docker installed but disabled at boot."
fi

# ---------------------------------------------------------------------------
# Done
# ---------------------------------------------------------------------------

echo ""
echo "==> setup-instance.sh complete."
echo ""
echo "NEXT STEPS:"
echo "  1. Run: scripts/oci-pilot.sh env-init"
echo "  2. Update Postgres password to match the generated DB password:"
echo "       sudo -u postgres psql -c \"ALTER ROLE open_transit WITH PASSWORD '<from-env-file>';\""
echo "  3. Run: scripts/oci-pilot.sh push"
echo "  4. Run: scripts/oci-pilot.sh units"
echo "  5. Run: scripts/oci-pilot.sh migrate"
echo "  6. Run: scripts/oci-pilot.sh start"
echo "  7. Configure Caddy: see deploy/oci/Caddyfile"
echo "  8. Run: scripts/oci-pilot.sh status"
