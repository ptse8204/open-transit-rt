# Phase Handoff

## Phase

Phase 17 — Deployment Automation And Pilot Operations

## Status

- Complete for the approved deployment automation, pilot operations runbook, helper-script, systemd-example, and evidence-refresh scope.
- Active phase after this handoff: Phase 18 — Admin UX And Agency Operations Console.

## What Was Implemented

- Added `docs/runbooks/small-agency-pilot-operations.md` as the primary Phase 17 operator guide.
- Added an explicit deployment environment variable matrix covering required/optional status, example values, secret classification, and service/script usage.
- Added evidence output locations, evidence commit-safety labels, and naming conventions:
  - `validator-cycle-YYYY-MM-DD.json`
  - `backup-run-YYYY-MM-DD.txt`
  - `restore-drill-YYYY-MM-DD.txt`
  - `feed-monitor-YYYY-MM-DD.txt`
  - `scorecard-export-YYYY-MM-DD.json`
- Added `scripts/pilot-ops.sh` with dry-run-capable `validator-cycle`, `backup`, `restore-drill`, `feed-monitor`, and `scorecard-export` subcommands.
- Added systemd service/timer examples for validator cycle, backup, feed monitor, and scorecard export.
- Updated `scripts/oci-pilot.sh` so remote/state-changing commands require explicit target env vars and print the target before doing remote work.
- Updated hosted evidence collection so admin-authenticated collection requires explicit `ADMIN_BASE_URL` and the workflow points operators to `EVIDENCE_PACKET_DIR=<packet> make audit-hosted-evidence`.

## What Was Designed But Intentionally Not Implemented Yet

- No backend API contracts were changed.
- No database schema or migration was added.
- No public feed URL or GTFS-RT contract was changed.
- No consumer-submission status was advanced.
- No hosted SaaS, Kubernetes, external predictor integration, consumer submission API, or product feature was added.
- No live webhook/email notification delivery integration was added; notification destinations remain deployment-owned placeholders.
- No new hosted evidence packet was collected in this pass.

## Deployment/Profile/Runbook Changes

- `docs/runbooks/small-agency-pilot-operations.md` documents the Phase 17 systemd/Caddy pilot profile and how it differs from the Phase 16 local app package.
- `docs/runbooks/deployment-evidence-overview.md` now defines Phase 17 helper output names and evidence labels.
- `docs/tutorials/deploy-with-docker-compose.md` and `docs/tutorials/production-checklist.md` now point operators to the Phase 17 operations profile and evidence refresh audit.
- Systemd examples use `EnvironmentFile={{OCI_REMOTE_DIR}}/ops/pilot-ops.env` and do not inline live secrets.

## Proxy/TLS Guidance Added

- `docs/runbooks/reverse-proxy-and-tls.md` now records the public-only default route set.
- `deploy/oci/Caddyfile` remains the redacted Caddy example exposing only:
  - `/public/gtfs/*`
  - `/public/feeds.json`
  - `/public/gtfsrt/vehicle_positions.pb`
  - `/public/gtfsrt/trip_updates.pb`
  - `/public/gtfsrt/alerts.pb`
- Admin/debug/JSON/metrics surfaces remain absent from the public edge, SSH-tunneled, or separately auth-protected.

## Validator Scheduling Changes

- Added `scripts/pilot-ops.sh validator-cycle`.
- Added `deploy/systemd/open-transit-validator-cycle.service`.
- Added `deploy/systemd/open-transit-validator-cycle.timer`.
- Validator evidence writes `validator-cycle-YYYY-MM-DD.json` plus per-feed response files to `EVIDENCE_OUTPUT_DIR`.
- Dry-run is required before enabling the timer.

## Backup/Restore Changes

- Added `scripts/pilot-ops.sh backup`.
- Added `scripts/pilot-ops.sh restore-drill`.
- Added `deploy/systemd/open-transit-backup.service`.
- Added `deploy/systemd/open-transit-backup.timer`.
- Backup evidence writes `backup-run-YYYY-MM-DD.txt`; raw dumps remain `never-commit`.
- Restore evidence writes `restore-drill-YYYY-MM-DD.txt`.
- Restore operations warn clearly and require typed confirmation unless `--force` is passed.

## Monitoring/Alerting Examples

- Added `scripts/pilot-ops.sh feed-monitor`.
- Added `deploy/systemd/open-transit-feed-monitor.service`.
- Added `deploy/systemd/open-transit-feed-monitor.timer`.
- Feed monitor evidence writes `feed-monitor-YYYY-MM-DD.txt`.
- Missing webhook/email destination is reported as `notification not configured`, not as feed failure.
- Real webhook URLs, notification credentials, and private incident links remain `never-commit`.

## Evidence Refresh Process

- `docs/evidence/README.md` and `docs/runbooks/deployment-evidence-overview.md` now require hosted refresh to end with:

```sh
EVIDENCE_PACKET_DIR=<packet> make audit-hosted-evidence
```

- Refreshed evidence must not be called complete unless this audit passes.
- Passing audit remains deployment/operator proof only; it is not CAL-ITP compliance, consumer acceptance, agency endorsement, hosted SaaS availability, or universal production readiness.

## Schema And Interface Changes

- No schema changes.
- No backend API contract changes.
- No public feed URL changes.
- New operator CLI interface:
  - `scripts/pilot-ops.sh validator-cycle [--dry-run]`
  - `scripts/pilot-ops.sh backup [--dry-run]`
  - `scripts/pilot-ops.sh restore-drill [--dry-run] [--force]`
  - `scripts/pilot-ops.sh feed-monitor [--dry-run]`
  - `scripts/pilot-ops.sh scorecard-export [--dry-run]`

## Dependency Changes

- Added Phase 17 pilot operations helper documentation to `docs/dependencies.md`.
- Added ADR-0025 to `docs/decisions.md`.
- No new external runtime dependency was added beyond existing shell tooling, Postgres client tools, `curl`, and systemd examples for deployments that choose to use them.

## Migrations Added

- None.

## Tests Added And Results

- No Go tests were added because Phase 17 changed deployment docs, shell helpers, and systemd examples without backend behavior changes.
- Script syntax and dry-run checks cover the new helper interfaces.

## Checks Run And Blocked Checks

Pre-edit baseline:

- `make validate`: passed.
- `make test`: passed.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `git diff --check`: passed.

Post-edit checks:

- `sh -n scripts/pilot-ops.sh scripts/oci-pilot.sh scripts/collect-hosted-evidence.sh scripts/audit-hosted-evidence.sh`: passed.
- `scripts/pilot-ops.sh help`: passed.
- `scripts/pilot-ops.sh validator-cycle --dry-run` with explicit placeholder env: passed.
- `scripts/pilot-ops.sh backup --dry-run` with explicit placeholder env: passed.
- `scripts/pilot-ops.sh restore-drill --dry-run` with explicit placeholder env: passed and printed destructive warning.
- `scripts/pilot-ops.sh feed-monitor --dry-run` with explicit placeholder env: passed and reported notification not configured.
- `scripts/pilot-ops.sh scorecard-export --dry-run` with explicit placeholder env: passed.
- `scripts/pilot-ops.sh feed-monitor --dry-run` without required env: failed as expected with a clear missing `ENVIRONMENT_NAME` message.
- `make validate`: passed.
- `make test`: passed.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `git diff --check`: passed.
- Targeted changed-file secret scan: no private key, GitHub token, Slack token/webhook, Discord webhook, or AWS access-key patterns found; broader token-name scan matched only documented placeholders and dev example values.

Blocked commands:

- None at handoff-writing time.

## Known Issues

- The Phase 17 helpers are deployment examples and do not replace a full observability stack.
- `scripts/pilot-ops.sh feed-monitor` checks public feed availability but does not parse protobuf freshness internally yet.
- Notification delivery is deployment-owned; the repo only documents placeholders and does not include a real webhook/email sender.
- The OCI pilot credentials from earlier local `.cache` files must not be reused without rotation/revocation.
- Optional systemd timers require a private, deployment-owned `pilot-ops.env` before they are enabled.

## Exact Next-Step Recommendation

- First files to read:
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/runbooks/small-agency-pilot-operations.md`
  - `docs/phase-18-admin-ux-agency-operations-console.md`
- First files likely to edit:
  - `cmd/agency-config/`
  - `cmd/gtfs-studio/`
  - `internal/compliance/`
  - `internal/devices/`
  - `docs/handoffs/phase-18.md`
- Commands to run before coding:
  - `make validate`
  - `make test`
  - `docker compose -f deploy/docker-compose.yml config`
  - `git diff --check`
- Known blockers:
  - No consumer acceptance evidence exists; all consumer records remain `not_started`.
  - Old ignored `.cache` secrets from Phase 15 still require rotation/revocation before real pilot reuse.
- Recommended first implementation slice:
  - Start Phase 18 by turning the Phase 17 operator evidence and health concepts into a minimal authenticated admin operations console, without changing public feed URLs or evidence claims.
