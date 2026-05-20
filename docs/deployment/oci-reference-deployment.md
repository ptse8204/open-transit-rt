# OCI/OCL Reference Deployment

This guide is documentation for reproducing a self-hosted deployment pattern. It is not evidence that a deployment was run, accepted, compliant, agency-approved, or production-ready.

This is the Phase 36 self-hosted OCI/OCL-style reference path for operators who
want to reuse the existing pilot-server pattern. It documents one repeatable
single-server shape using compiled Go binaries, PostgreSQL/PostGIS, systemd,
and Caddy or an equivalent reverse proxy.

It is not hosted SaaS, agency adoption, agency endorsement, final-root proof,
consumer submission/review/acceptance, CAL-ITP/Caltrans compliance, production
readiness, vendor compatibility, or production-grade ETA quality evidence.

## Operator Copy/Paste Path

1. Prepare server.
2. Create user and directories.
3. Install dependencies.
4. Build or upload binaries.
5. Configure env.
6. Migrate database.
7. Install systemd units.
8. Configure Caddy.
9. Import GTFS.
10. Run smoke checklist.

Use the detailed sections below for the exact boundaries and checks. The
`scripts/oci-pilot.sh` script and Makefile `oci-*` targets are repo-supported
helpers for the current reference path. They are not hosted product automation
and do not prove a deployment was run.

For one guided local/reference evaluation that combines this deployment path,
`make agency-pilot-up`, five public feed checks, readiness review, validator
handling, and the synthetic AVL dry-run adapter, use
[Self-Hosted Operator Trial](../tutorials/self-hosted-operator-trial.md).

For a capture-backed browser walkthrough of the same private workflow, use the
[Operator Workflow Tutorial](https://ptse8204.github.io/open-transit-rt/ui-tour.html).

For browser-first agency operations after deployment, use the same private
product routes as the local evaluator:

```text
/admin/operations
/admin/operations/setup-wizard
/admin/operations/feed-health
/admin/operations/gtfs-import
/admin/operations/gtfs-workbench
/admin/operations/connectors
/admin/operations/validation-health
/admin/operations/readiness
/admin/operations/maintenance
```

These private pages are the normal operator surfaces. Shell diagnostics remain
administrator or deployment-owner paths for startup, deployment, off-host
validation, support summaries, and secure token handling.

## Server Prerequisites

Target one operator-owned Linux server with:

- SSH and sudo access for the operator.
- A deployment-owned DNS name pointed at the server.
- Public ports `80` and `443` open at the firewall and cloud network edge.
- Loopback-only service binding for Open Transit RT services.
- PostgreSQL with PostGIS available locally or on an operator-controlled
  database endpoint.
- Caddy, nginx, or another TLS reverse proxy.
- `curl`, `git`, `unzip`, `zip`, `ca-certificates`, `python3`, and `openssl`.
- Java and Docker only on a machine that intentionally runs validators. Tiny
  reference servers should usually keep heavy validator workloads off host and
  use `docs/deployment/off-host-validation.md` from an operator workstation.
- Enough disk for database growth, backups, validator cache, and private
  operator evidence under `/opt/open-transit-rt`.

The current reference implementation has been shaped around Oracle Linux 9 on
an OCI-style VM, but the path is a self-hosted deployment pattern rather than
an OCI-only product.

## User And Group Layout

Create a dedicated service account and keep operator access separate:

```sh
sudo useradd --system --home-dir /opt/open-transit-rt --shell /sbin/nologin open-transit
sudo mkdir -p /opt/open-transit-rt
sudo chown open-transit:open-transit /opt/open-transit-rt
sudo chmod 750 /opt/open-transit-rt
```

Operators should use sudo for lifecycle actions. Do not run public services as
`root`, and do not put live secrets in shell history, unit files, or committed
docs.

## Directory Layout

Use this layout under `/opt/open-transit-rt`:

```text
/opt/open-transit-rt/
  bin/                         # compiled service binaries
  env                          # private runtime env, mode 0600
  app/
    db/migrations/             # repo migrations
    deploy/systemd/            # unit templates
    deploy/oci/                # setup/proxy helper assets
    scripts/                   # repo-supported helper scripts
    testdata/                  # optional local GTFS fixtures
  .cache/validators/           # pinned validator tooling
  data/                        # operator-owned import workspace
  ops/pilot-ops.env            # private helper env, mode 0600
  backups/                     # private DB dumps
  evidence/                    # private operator run output before review
```

Private paths stay operator-owned until reviewed against
`docs/evidence/redaction-policy.md`.

## Install Dependencies

Install OS packages, PostgreSQL/PostGIS, and Caddy using the server's package
manager. For the OCI-style reference, `deploy/oci/setup-instance.sh` captures
the current helper path for Oracle Linux 9. Treat it as repo-supported helper
tooling for this reference path, not as a managed-hosting promise.

If using the helper:

```sh
ENVIRONMENT_NAME=reference \
OCI_HOST=replace-with-host \
DOMAIN=feeds.example.org \
OCI_REMOTE_DIR=/opt/open-transit-rt \
scripts/oci-pilot.sh setup
```

Review the script output before proceeding, especially database, firewall,
service-user, and Caddy install steps.

## Build Or Upload Binaries

Build Linux binaries from the repo on a trusted workstation or CI runner:

```sh
make oci-build
```

Upload binaries, migrations, systemd templates, OCI helper assets, scripts, and
optional test fixtures to the server:

```sh
ENVIRONMENT_NAME=reference \
OCI_HOST=replace-with-host \
DOMAIN=feeds.example.org \
OCI_REMOTE_DIR=/opt/open-transit-rt \
make oci-push
```

Equivalent manual uploads are acceptable if they preserve ownership, modes, and
the directory layout above.

## Configure Env With Placeholders Only

Start from [oci-reference-env.example](oci-reference-env.example). Copy it to
private server paths and replace every placeholder with deployment-owned values:

```sh
sudo install -o open-transit -g open-transit -m 0600 \
  docs/deployment/oci-reference-env.example /opt/open-transit-rt/env
sudo install -o open-transit -g open-transit -m 0600 \
  docs/deployment/oci-reference-env.example /opt/open-transit-rt/ops/pilot-ops.env
```

Do not commit the populated files. Do not reuse historical pilot hostnames,
real database URLs, admin tokens, webhook URLs, or evidence values.

## Generate Secrets

Generate a unique value for each secret:

```sh
openssl rand -hex 32
```

Use separate generated values for at least:

- database role password;
- `ADMIN_JWT_SECRET`;
- `CSRF_SECRET`;
- `DEVICE_TOKEN_PEPPER`;
- admin helper token;
- notification webhook credential, if configured.

Rotate secrets after accidental disclosure. Follow `SECURITY.md` and
`docs/evidence/redaction-policy.md` if a secret reaches tracked files or public
artifacts.

## Database Setup And Migrations

Create a dedicated database and role. Example shape:

```sql
CREATE ROLE open_transit LOGIN PASSWORD 'replace-with-generated-secret';
CREATE DATABASE open_transit_rt OWNER open_transit;
\c open_transit_rt
CREATE EXTENSION IF NOT EXISTS postgis;
```

Set `DATABASE_URL`, conservative `DB_MAX_CONNS`, and `MIGRATIONS_DIR` in
`/opt/open-transit-rt/env`, then run:

```sh
sudo -u open-transit sh -lc '
  set -a
  . /opt/open-transit-rt/env
  set +a
  /opt/open-transit-rt/bin/migrate status
  /opt/open-transit-rt/bin/migrate up
  /opt/open-transit-rt/bin/migrate status
'
```

The helper equivalent is:

```sh
ENVIRONMENT_NAME=reference \
OCI_HOST=replace-with-host \
DOMAIN=feeds.example.org \
OCI_REMOTE_DIR=/opt/open-transit-rt \
scripts/oci-pilot.sh migrate
```

## First Admin Bootstrap Link

After migrations are applied, create the first admin user and one-time setup
link from the server/operator console:

```sh
sudo -u open-transit sh -lc '
  set -a
  . /opt/open-transit-rt/env
  set +a
  /opt/open-transit-rt/bin/agency-config bootstrap-admin-link \
    --agency-id "$AGENCY_ID" \
    --email admin@example.org \
    --base-url http://127.0.0.1:8081 \
    --ttl 30m
'
```

The command stores only a token hash, binds the user to the existing `admin`
role, and prints the setup URL once. Do not paste the generated URL into
tracked docs, evidence folders, issue comments, screenshots, or public logs.
Open the link through the deployment-owned admin access path, set the admin
password, then use `/admin/login` for normal browser sign-in and `/admin/logout`
to end the session. Successful password login issues the same signed
`admin_session` cookie used by the admin middleware. `/admin/local-login`
remains local/demo-only and production-disabled. SSO/OIDC is not implemented in
this roadmap; a future identity provider would issue the same internal
`admin_session` cookie after identity verification and role mapping.

## Service Supervision

The reference service set is:

- `agency-config` on loopback port `8081`;
- `telemetry-ingest` on loopback port `8082`;
- `feed-vehicle-positions` on loopback port `8083`;
- `feed-trip-updates` on loopback port `8084`;
- `feed-alerts` on loopback port `8085`.

Install and enable systemd units from `deploy/systemd/` using the helper or an
equivalent operator-controlled installation:

```sh
ENVIRONMENT_NAME=reference \
OCI_HOST=replace-with-host \
DOMAIN=feeds.example.org \
OCI_REMOTE_DIR=/opt/open-transit-rt \
make oci-units
```

Verify services:

```sh
sudo systemctl status open-transit-agency-config.service
sudo systemctl status open-transit-telemetry-ingest.service
sudo systemctl status open-transit-feed-vehicle-positions.service
sudo systemctl status open-transit-feed-trip-updates.service
sudo systemctl status open-transit-feed-alerts.service
```

## Caddy Or Equivalent Reverse Proxy

The public edge should expose only anonymous feed and discovery paths:

- `/public/feeds.json`
- `/public/gtfs/schedule.zip`
- `/public/gtfsrt/vehicle_positions.pb`
- `/public/gtfsrt/trip_updates.pb`
- `/public/gtfsrt/alerts.pb`

The reference Caddyfile is `deploy/oci/Caddyfile`. It proxies those paths to
loopback services and returns `404` for everything else. If using nginx,
Traefik, or another reverse proxy, preserve the same public/private boundary.

## Public And Private Route Boundaries

Public unauthenticated routes:

- schedule GTFS ZIP;
- feeds metadata JSON;
- GTFS-Realtime Vehicle Positions protobuf;
- GTFS-Realtime Trip Updates protobuf;
- GTFS-Realtime Alerts protobuf.

Private or internal routes:

- `/admin/*`;
- `/admin/debug/*`;
- `/admin/operations/*`;
- GTFS Studio routes;
- `/v1/events`;
- JSON debug realtime routes;
- `/metrics`, when enabled.

Expose private routes only through SSH tunnels, VPN, or an auth-aware proxy
under operator control. Do not route admin/debug/studio surfaces through the
anonymous public feed edge.

## Validator Install And Check

Install pinned validator tooling from the repo before publication checks:

```sh
make validators-install
make validators-check
```

On the server, keep validator artifacts under
`/opt/open-transit-rt/.cache/validators/` and set the validator env variables
from [oci-reference-env.example](oci-reference-env.example). If a deployment
uses equivalent tooling instead of the repo wrapper, document the pinned
version or digest, command boundary, and failure behavior in operator notes
before relying on it.

## GTFS Import

Place the GTFS ZIP in private operator storage, such as:

```text
/opt/open-transit-rt/data/gtfs/agency-schedule.zip
```

Import through the deployed binary:

```sh
sudo -u open-transit sh -lc '
  set -a
  . /opt/open-transit-rt/env
  set +a
  /opt/open-transit-rt/bin/gtfs-import \
    -agency-id "$AGENCY_ID" \
    -zip /opt/open-transit-rt/data/gtfs/agency-schedule.zip \
    -actor-id reference-operator \
    -notes "Reference deployment import"
'
```

Raw GTFS ZIPs should remain private unless licensing and redaction have been
reviewed.

For local/reference deployments where the operator has an `AGENCY_ID` and
public GTFS URL, the reusable onboarding helper can perform the download,
checksum, import, metadata bootstrap, five-path fetch, and validator/blocker
summary:

```sh
scripts/agency-pilot-onboard.sh \
  --mode running \
  --agency-id "$AGENCY_ID" \
  --gtfs-url https://example.org/gtfs.zip \
  --public-base-url https://feeds.example.org \
  --admin-base-url http://127.0.0.1:8081 \
  --admin-token "$ADMIN_TOKEN" \
  --technical-contact-email ops@example.org \
  --feed-license-name "replace-with-approved-license" \
  --feed-license-url https://example.org/license
```

Publication metadata is local/reference placeholder metadata unless the
operator supplied agency-approved values. Do not treat placeholder metadata as
agency-approved, final-root evidence, consumer evidence, CAL-ITP/Caltrans
compliance, or production readiness. The helper writes raw GTFS and fetch
summaries under ignored `.cache/` storage by default; that output is not
final-root evidence unless a future evidence phase explicitly reviews, redacts,
and retains it.

## Five Public Feed URL Verification

After migrations, service start, metadata bootstrap, and GTFS import, verify
all five public URLs from outside the server. Prefer the off-host helper when
an operator machine has the repo checkout:

```sh
PUBLIC_BASE_URL=https://feeds.example.org make validate-public-feeds
```

This writes private diagnostics under `.cache/validate-public-feeds/<timestamp>`
and distinguishes failed fetches, missing validator tooling, validator
failures, and skipped validation. It does not require Java or Docker on the
tiny reference server.

Manual equivalent:

```sh
BASE=https://feeds.example.org
curl -fsS -D feeds.headers "$BASE/public/feeds.json" -o feeds.json
curl -fsS -D schedule.headers "$BASE/public/gtfs/schedule.zip" -o schedule.zip
curl -fsS -D vp.headers "$BASE/public/gtfsrt/vehicle_positions.pb" -o vehicle_positions.pb
curl -fsS -D tu.headers "$BASE/public/gtfsrt/trip_updates.pb" -o trip_updates.pb
curl -fsS -D alerts.headers "$BASE/public/gtfsrt/alerts.pb" -o alerts.pb
shasum -a 256 feeds.json schedule.zip vehicle_positions.pb trip_updates.pb alerts.pb
```

HTTP success, byte counts, and checksums are operator verification data only
until reviewed and redacted. They do not prove consumer acceptance,
compliance, final-root ownership, or production readiness.

## OCI Reference Diagnostic

Use the reference check as product-support diagnostics, not evidence:

```sh
PUBLIC_BASE_URL=https://feeds.example.org make oci-reference-check
```

When `OCI_HOST` access is configured, the helper also checks loopback health
for the five services through SSH. It records deployment helper presence,
backup/restore configuration presence, public fetch/validator state, telemetry
simulator guidance, and safety booleans without printing populated env values.

See [OCI Reference Check](oci-reference-check.md).

## Backup And Restore

Run dry-runs first:

```sh
ENVIRONMENT_NAME=reference \
EVIDENCE_OUTPUT_DIR=/opt/open-transit-rt/evidence/$(date -u +%Y-%m-%d) \
DATABASE_URL=replace-with-database-url-containing-generated-secret \
BACKUP_DIR=/opt/open-transit-rt/backups \
/opt/open-transit-rt/app/scripts/pilot-ops.sh backup --dry-run
```

Then configure `/opt/open-transit-rt/ops/pilot-ops.env` and enable
`open-transit-backup.timer` only after paths and permissions are verified.

Restore drills must target a restore database, not the live database, unless
the operator is intentionally performing an incident restore:

```sh
/opt/open-transit-rt/app/scripts/pilot-ops.sh restore-drill --dry-run
```

Use `docs/runbooks/templates/restore-event-template.md` for private operator
records, and redact before committing any summary.

Before using backup, restore, upgrade, or rollback procedures, run the
deployment doctor and review the `install_recovery` category plus the private
stop-point summary:

```sh
PUBLIC_BASE_URL=https://feeds.example.org make deployment-doctor
```

The summary is written under
`.cache/deployment-doctor/<timestamp>/operations/install-upgrade-recovery-plan.summary.json`
and groups the current state into environment, migration, backup/restore,
service/proxy, and post-change verification stop points. Treat any blocker as
a pause before migration, service restart, rollback, restore, or public feed
sharing. The summary is private diagnostic guidance only; it does not execute
or prove backup, restore, upgrade, rollback, production readiness, SLA/uptime,
compliance, or consumer acceptance.

## Feed Monitor

Run the feed monitor dry-run before enabling its timer:

```sh
ENVIRONMENT_NAME=reference \
EVIDENCE_OUTPUT_DIR=/opt/open-transit-rt/evidence/$(date -u +%Y-%m-%d) \
PUBLIC_BASE_URL=https://feeds.example.org \
/opt/open-transit-rt/app/scripts/pilot-ops.sh feed-monitor --dry-run
```

If notification placeholders are configured, keep webhook URLs and notification
credentials private. A successful monitor run is an operator health check, not
consumer ingestion or compliance evidence.

## Scorecard Export

Run scorecard export through the private admin boundary:

```sh
ENVIRONMENT_NAME=reference \
EVIDENCE_OUTPUT_DIR=/opt/open-transit-rt/evidence/$(date -u +%Y-%m-%d) \
ADMIN_BASE_URL=http://127.0.0.1:8081 \
ADMIN_TOKEN=replace-with-generated-secret \
/opt/open-transit-rt/app/scripts/pilot-ops.sh scorecard-export --dry-run
```

Do not expose `ADMIN_TOKEN` in public logs, issue comments, terminal captures,
or committed artifacts.

## Readiness Workflow

After the five public feed URLs, publication metadata, validation records,
telemetry, and operations helpers are in place, review the authenticated
Operations Console readiness page:

```text
/admin/operations/readiness
```

The page shows CAL-ITP-style readiness rows with status sources, current
signals, next actions, and claim boundaries. It is an operator workflow view
only. It does not create external evidence, submit to consumers, prove an
agency-owned final root, or claim CAL-ITP/Caltrans compliance.

Also review:

```text
/admin/operations/maintenance
```

The Maintenance Center shows active feed version, last GTFS import, latest
five-feed check when recorded, validator state, backup/restore configuration
presence, telemetry freshness, service-health availability, support-summary
instructions, and weekly/monthly tasks.

For the full guided trial around this readiness review, see
[Self-Hosted Operator Trial](../tutorials/self-hosted-operator-trial.md).
For repeatable smoke checks and redaction-safe diagnostics, see
[Operator Smoke And Support Bundle](../tutorials/operator-smoke-and-support-bundle.md).
For read-only reference deployment preflight diagnostics, see
[Reference Deployment Doctor](reference-deployment-doctor.md).
For public-safe summaries of completed reference diagnostics, see
[OCI Reference Diagnostic Runs](oci-reference-diagnostic-runs.md). Those
records are not retained evidence packets and must keep the same claim
boundaries as this guide.
For Phase 56 route/proxy boundary diagnostics, run `make multi-agency-hosting`.
The diagnostic output is private `.cache` review material only and does not
prove hosted SaaS availability, production multi-tenant hosting, compliance,
SLA/uptime, agency adoption, consumer acceptance, or production readiness.

## Update And Rollback

Before an update:

1. Record the deployed git tag or commit.
2. Take a database backup.
3. Retain previous binaries in a private operator directory.
4. Run `migrate status`.
5. Confirm the five public feed URLs are healthy.
6. Run `make deployment-doctor` and review the `install_recovery` stop-point
   summary before migration, service restart, rollback, or public sharing work.

Update by uploading new binaries and migrations, running migrations, restarting
services, and running the smoke checklist.

Rollback depends on what changed:

- Binary-only issue: restore previous binaries and restart services.
- Migration issue before data changes: follow `docs/upgrade-and-rollback.md`
  and the migration rollback guidance for the specific release.
- Data issue after migration: restore from the pre-update backup into a known
  database state before restarting services.

Do not claim rollback proof from documentation alone. Retained operator output
must be reviewed and redacted before it becomes evidence.

## Redacted Support Bundle Guidance

The Phase 41 helper can collect a redaction-safe default bundle:

```sh
PUBLIC_BASE_URL=https://feeds.example.org \
ADMIN_BASE_URL=https://admin.example.org \
make support-bundle
```

Only set `ADMIN_TOKEN` and `INCLUDE_ADMIN_READINESS=true` when the admin URL is
private/safe and the summary is needed. The script uses bearer auth only and
does not support cookies.

A support bundle may include:

- deployed git tag or commit;
- service unit names and active/inactive status;
- Caddy or reverse-proxy route map with private hosts redacted;
- env variable names, never secret values;
- five public feed URL status codes, byte counts, and checksums;
- validator summary status and tool versions;
- recent service logs with tokens, client IPs, and private hosts redacted;
- backup and restore-drill summaries;
- feed-monitor and scorecard-export summaries.

Exclude:

- database dumps;
- raw credentials, JWTs, API keys, device tokens, webhook URLs, and DB URLs;
- private keys, certificates, SSH material, and ACME account keys;
- raw access logs with client IPs unless explicitly reviewed;
- raw GTFS ZIPs unless license and privacy review says they are safe;
- private correspondence, ticket links, or consumer portal details.

Apply `docs/evidence/redaction-policy.md` before sharing or committing any
bundle.
