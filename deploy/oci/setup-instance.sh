#!/usr/bin/env sh
# deploy/oci/setup-instance.sh — First-time Oracle Linux 9 instance setup
#
# Target OS: Oracle Linux Server 9.x (RHEL 9 compatible, uses dnf/firewalld)
# Run via:   scripts/oci-pilot.sh setup
# Or manually on the instance: sudo ./setup-instance.sh
#
# Idempotent — safe to run more than once.
# Does NOT write the application env file or secrets.

set -eu

echo "==> open-transit-rt OCI instance setup (Oracle Linux 9)"
echo "    $(date -u)"

# ---------------------------------------------------------------------------
# 1. Firewall — open ports 80 and 443 via firewalld
# ---------------------------------------------------------------------------

echo "==> Opening ports 80 and 443 in firewalld..."
firewall-cmd --permanent --add-port=80/tcp  2>/dev/null || true
firewall-cmd --permanent --add-port=443/tcp 2>/dev/null || true
firewall-cmd --reload 2>/dev/null || true
echo "==> Firewalld ports: $(firewall-cmd --list-ports 2>/dev/null || echo '(check manually)')"

# ---------------------------------------------------------------------------
# 2. Swap — check/create (OCI may have pre-configured swap)
# ---------------------------------------------------------------------------

if swapon --show | grep -q swap 2>/dev/null; then
  echo "==> Swap already enabled: $(free -h | awk '/Swap/{print $2}')"
else
  echo "==> Creating 4G swap file..."
  fallocate -l 4G /swapfile
  chmod 0600 /swapfile
  mkswap /swapfile
  swapon /swapfile
  grep -q '/swapfile' /etc/fstab || echo '/swapfile none swap sw 0 0' >> /etc/fstab
  echo "vm.swappiness=10" >> /etc/sysctl.conf && sysctl -p
  echo "==> Swap created: $(free -h | awk '/Swap/{print $2}')"
fi

# ---------------------------------------------------------------------------
# 3. System packages (dnf)
# ---------------------------------------------------------------------------

echo "==> Installing base system packages via dnf..."
dnf install -y --setopt=install_weak_deps=False \
  curl wget unzip zip git ca-certificates python3
echo "==> Base system packages installed."

if ! swapon --show | grep -q '/swapfile-extra' 2>/dev/null; then
  echo "==> Creating temporary 4G extra swap for Oracle Linux package transactions..."
  fallocate -l 4G /swapfile-extra
  chmod 0600 /swapfile-extra
  mkswap /swapfile-extra
  swapon /swapfile-extra
  echo "==> Extra swap enabled for setup: $(free -h | awk '/Swap/{print $2}')"
fi

# ---------------------------------------------------------------------------
# 4. PostgreSQL 15 + PostGIS (PGDG repo)
# ---------------------------------------------------------------------------

if rpm -q pgdg-redhat-repo > /dev/null 2>&1; then
  echo "==> PGDG repo already installed."
else
  echo "==> Installing PGDG repo..."
  rpm -Uvh --nodeps \
    https://download.postgresql.org/pub/repos/yum/reporpms/EL-9-x86_64/pgdg-redhat-repo-latest.noarch.rpm
fi

if rpm -q postgresql15-server > /dev/null 2>&1; then
  echo "==> PostgreSQL 15 already installed."
else
  echo "==> Installing PostgreSQL 15 from PGDG repo..."
  dnf install -y --setopt=install_weak_deps=False \
    --disablerepo='*' \
    --enablerepo=pgdg15 \
    --enablerepo=ol9_baseos_latest \
    postgresql15-server
  echo "==> PostgreSQL 15 installed."
fi

if rpm -q postgis34_15 > /dev/null 2>&1; then
  echo "==> PostGIS already installed."
else
  echo "==> Installing PostGIS 3.4 for PostgreSQL 15..."
  dnf install -y --nobest --setopt=install_weak_deps=False \
    --disablerepo='*' \
    --enablerepo=pgdg15 \
    --enablerepo=pgdg-common \
    --enablerepo=ol9_appstream \
    --enablerepo=ol9_baseos_latest \
    --enablerepo=ol9_codeready_builder \
    --enablerepo=ol9_developer_EPEL \
    postgis34_15
fi

# Initialize PostgreSQL cluster if not done yet
if [ ! -f /var/lib/pgsql/15/data/PG_VERSION ]; then
  echo "==> Initializing PostgreSQL 15 data directory..."
  PGSETUP_INITDB_OPTIONS="--encoding=UTF8 --locale=C" \
    /usr/pgsql-15/bin/postgresql-15-setup initdb
fi

# Apply memory tuning before first start
PG_CONF_DIR="/var/lib/pgsql/15/data"
TUNING_CONF="${PG_CONF_DIR}/conf.d/open-transit-rt.conf"
mkdir -p "${PG_CONF_DIR}/conf.d"

# Enable conf.d inclusion in postgresql.conf
if ! grep -q "include_dir" "${PG_CONF_DIR}/postgresql.conf" 2>/dev/null; then
  echo "include_dir = 'conf.d'" >> "${PG_CONF_DIR}/postgresql.conf"
fi

if [ ! -f "$TUNING_CONF" ]; then
  cat > "$TUNING_CONF" <<'PGEOF'
# Open Transit RT — VM.Standard.E2.1.Micro tuning (~503 MB usable RAM)
shared_buffers              = 96MB
effective_cache_size        = 384MB
work_mem                    = 3MB
maintenance_work_mem        = 24MB
max_connections             = 25
wal_buffers                 = 4MB
checkpoint_completion_target = 0.9
random_page_cost            = 1.1
log_min_duration_statement  = 500
PGEOF
  echo "==> PostgreSQL tuning applied."
fi

# Bind PostgreSQL to loopback only
sed -i "s|#listen_addresses = 'localhost'|listen_addresses = '127.0.0.1'|" \
  "${PG_CONF_DIR}/postgresql.conf" 2>/dev/null || \
  sed -i "s|listen_addresses = '\*'|listen_addresses = '127.0.0.1'|" \
  "${PG_CONF_DIR}/postgresql.conf" 2>/dev/null || true

# Adjust pg_hba.conf to allow md5/scram auth for local connections
if ! grep -q "open_transit" "${PG_CONF_DIR}/pg_hba.conf" 2>/dev/null; then
  # Add a line before the default ident local auth lines
  sed -i '/^local.*all.*all.*peer/i host    open_transit_rt    open_transit    127.0.0.1/32    scram-sha-256' \
    "${PG_CONF_DIR}/pg_hba.conf" || true
fi

systemctl enable postgresql-15
systemctl restart postgresql-15
echo "==> PostgreSQL running: $(systemctl is-active postgresql-15)"

# ---------------------------------------------------------------------------
# 5. Create Postgres DB and user
# ---------------------------------------------------------------------------

echo "==> Ensuring database 'open_transit_rt' and role 'open_transit' exist..."
sudo -u postgres /usr/pgsql-15/bin/psql -tc \
  "SELECT 1 FROM pg_roles WHERE rolname='open_transit'" \
  | grep -q 1 || \
  sudo -u postgres /usr/pgsql-15/bin/psql -c \
  "CREATE ROLE open_transit WITH LOGIN PASSWORD 'changeme-run-env-init';"

sudo -u postgres /usr/pgsql-15/bin/psql -tc \
  "SELECT 1 FROM pg_database WHERE datname='open_transit_rt'" \
  | grep -q 1 || \
  sudo -u postgres /usr/pgsql-15/bin/psql -c \
  "CREATE DATABASE open_transit_rt OWNER open_transit;"

sudo -u postgres /usr/pgsql-15/bin/psql -d open_transit_rt \
  -c "CREATE EXTENSION IF NOT EXISTS postgis;" 2>/dev/null || true
sudo -u postgres /usr/pgsql-15/bin/psql -d open_transit_rt \
  -c "CREATE EXTENSION IF NOT EXISTS postgis_topology;" 2>/dev/null || true
sudo -u postgres /usr/pgsql-15/bin/psql \
  -c "GRANT ALL ON DATABASE open_transit_rt TO open_transit;" 2>/dev/null || true

echo "==> IMPORTANT: After running env-init, update Postgres password:"
echo "    sudo -u postgres /usr/pgsql-15/bin/psql -c \"ALTER ROLE open_transit WITH PASSWORD '<from-env-file>';\""

# ---------------------------------------------------------------------------
# 6. Caddy (official Cloudsmith RPM repo)
# ---------------------------------------------------------------------------

if command -v caddy > /dev/null 2>&1; then
  echo "==> Caddy already installed: $(caddy version)"
else
  echo "==> Installing Caddy..."
  dnf install -y --setopt=install_weak_deps=False \
    --disablerepo='*' \
    --enablerepo='copr:copr.fedorainfracloud.org:group_caddy:caddy' \
    --enablerepo=ol9_baseos_latest \
    caddy
  echo "==> Caddy installed: $(caddy version)"
fi

# Allow Caddy to bind ports 80/443 without root (SELinux + capability)
setcap cap_net_bind_service=+ep "$(command -v caddy)" 2>/dev/null || true

# Allow Caddy to connect to loopback backend (SELinux httpd_can_network_connect)
setsebool -P httpd_can_network_connect 1 2>/dev/null || true

systemctl enable caddy
echo "==> Caddy enabled."

# ---------------------------------------------------------------------------
# 7. Go toolchain
# ---------------------------------------------------------------------------

echo "==> Go is not installed by default; scripts/oci-pilot.sh deploys local cross-compiled binaries."

# ---------------------------------------------------------------------------
# 8. Create application system user
# ---------------------------------------------------------------------------

if id open-transit > /dev/null 2>&1; then
  echo "==> System user 'open-transit' already exists."
else
  echo "==> Creating system user 'open-transit'..."
  useradd -r -m -d /opt/open-transit-rt -s /sbin/nologin open-transit
fi
mkdir -p /opt/open-transit-rt
chown open-transit:open-transit /opt/open-transit-rt
chmod 750 /opt/open-transit-rt

# ---------------------------------------------------------------------------
# 9. Docker (optional — for GTFS-RT validator only)
# ---------------------------------------------------------------------------

echo "==> Docker is intentionally not installed on this 503MiB RAM instance."

# ---------------------------------------------------------------------------
# Done
# ---------------------------------------------------------------------------

echo ""
echo "==> setup-instance.sh complete (Oracle Linux 9)."
echo ""
echo "NEXT STEPS:"
echo "  1. scripts/oci-pilot.sh env-init"
echo "  2. Update Postgres password to match generated DB password"
echo "  3. scripts/oci-pilot.sh push"
echo "  4. scripts/oci-pilot.sh units"
echo "  5. scripts/oci-pilot.sh migrate"
echo "  6. scripts/oci-pilot.sh start"
echo "  7. scripts/oci-pilot.sh token <agency-id> && scripts/oci-pilot.sh bootstrap"
