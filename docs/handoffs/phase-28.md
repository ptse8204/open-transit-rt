# Phase 28 Handoff

## Phase

Phase 28 — Production Operations Hardening

## Status

- Complete for the docs-first operations hardening scope.
- Active phase after this handoff: Phase 29 — Real-World Realtime Quality Expansion.

## What Was Implemented

- Added `docs/runbooks/production-operations-hardening.md` with daily, weekly, monthly, and after-upgrade operations checks.
- Added template-only incident, restore, secret rotation, and operator handover records under `docs/runbooks/templates/`.
- Hardened backup/restore, monitoring/alerting, validator evidence, deployment evidence, release checklist, and upgrade/rollback guidance.
- Updated docs navigation, phase status, current status, and changelog.

## What Was Designed But Intentionally Not Implemented Yet

- No runtime API, database schema, public feed URL, GTFS-RT contract, consumer status, external integration, systemd, Docker, or evidence-claim changes.
- No hosted monitoring SaaS, Prometheus/Grafana deployment, Kubernetes, tenant-safe backup/restore/export tooling, consumer submission automation, or external portal contact.
- No fake incidents, fake outage evidence, fake alert delivery proof, fake rotation records, fake restore events, or placeholder operational artifacts.

## Runbooks And Templates Added

- `docs/runbooks/production-operations-hardening.md`
- `docs/runbooks/templates/README.md`
- `docs/runbooks/templates/feed-outage-incident-template.md`
- `docs/runbooks/templates/validator-failure-incident-template.md`
- `docs/runbooks/templates/telemetry-staleness-incident-template.md`
- `docs/runbooks/templates/trip-updates-quality-incident-template.md`
- `docs/runbooks/templates/secret-exposure-incident-template.md`
- `docs/runbooks/templates/consumer-complaint-or-rejection-template.md`
- `docs/runbooks/templates/restore-event-template.md`
- `docs/runbooks/templates/secret-rotation-record-template.md`
- `docs/runbooks/templates/operator-handover-template.md`

## Backup/Restore Hardening

- Documented daily backup cadence, minimum retention guidance, monthly restore-drill cadence, pre-upgrade backup rule, restore verification, post-restore feed checks, post-restore validation checks, and redaction of backup paths and DB URLs.
- Preserved the Phase 27 boundary: current backup/restore/export/evidence workflows are deployment/DB scoped and are not tenant-safe multi-agency workflows.

## Upgrade/Rollback Hardening

- Added pre-upgrade backup evidence requirements, migration status before/after upgrade, post-upgrade public feed fetches, validator reruns, Trip Updates/Alerts checks, rollback limits, restore-event record guidance, and release-note record expectations.
- Irreversible or untested migrations remain backup/restore recovery events, not `migrate-down` assumptions.

## Monitoring/Alerting Hardening

- Covered public feed HTTP status, public feed freshness, `/readyz`, validator failure, backup failure, disk space, database connectivity, telemetry staleness, Trip Updates withheld/coverage changes, and Alerts availability.
- Added alert delivery proof pattern: monitor/check name, expected destination, test timestamp, delivery result, redacted proof location, `notification not configured` state, and follow-up on failed delivery.
- Added capacity guidance for disk, DB growth, backup storage, logs, and evidence artifacts.

## Validator Failure Response

- Added schedule validator failure, realtime validator failure, validator tooling unavailable, warning-vs-failed handling, rerun criteria, evidence packet implications, public claim implications, and continue/stop publishing decision guidance.

## Secret Rotation Guidance

- Covered admin JWT secret, CSRF secret, device token pepper, device tokens, DB password, TLS/ACME material, optional webhook/notification credentials, and Phase 15 `.cache` secret findings.
- Documented that deleting a file is not enough when a real secret was exposed; rotate or revoke the credential and verify the old value no longer works.

## Operator Handover Guidance

- Added required handover fields: current release version, deployment environment, public feed URLs, admin access process without secrets, secret storage location without secrets, backup location, restore process, validator cadence, monitoring cadence, evidence packet location, known blockers, consumer status boundaries, agency-owned-domain status, and multi-agency limitations.

## Schema And Interface Changes

- None.

## Dependency Changes

- None.

## Migrations Added

- None.

## Script Changes

- None.

## Tests Added And Results

- No code tests were added because Phase 28 was docs-first.
- Required checks run:
  - Pre-edit `make validate` — passed.
  - Pre-edit `make test` — passed.
  - Pre-edit `make test-integration` — passed.
  - Pre-edit `make realtime-quality` — passed.
  - Pre-edit `make smoke` — passed.
  - Pre-edit `docker compose -f deploy/docker-compose.yml config` — passed.
  - Pre-edit `git diff --check` — passed.
  - Post-edit `make validate` — passed.
  - Post-edit `make test` — passed.
  - Post-edit `make test-integration` — passed.
  - Post-edit `make realtime-quality` — passed.
  - Post-edit `make smoke` — passed.
  - Post-edit `docker compose -f deploy/docker-compose.yml config` — passed.
  - Post-edit `git diff --check` — passed.
  - Post-edit targeted redaction and forbidden-claim scan — passed.

## Checks Run And Blocked Checks

- Blocked commands: none.

## Known Remaining Production-Operations Gaps

- Alert delivery proof requires real deployment/operator evidence.
- Restore drills require real private backups and isolated restore targets.
- Tenant-safe backup/restore/export/evidence workflows remain unimplemented.
- Agency-owned final feed root proof remains blocked until an agency-owned or agency-approved root exists.
- Consumer targets remain `prepared` only without target-originated submission, review, rejection, blocker, acceptance, or ingestion evidence.

## Exact Next-Step Recommendation

- First files to read:
  - `docs/phase-29-realtime-quality-expansion.md`
  - `docs/handoffs/latest.md`
  - `docs/current-status.md`
  - `docs/runbooks/production-operations-hardening.md`
- First files likely to edit:
  - `internal/realtimequality/`
  - `testdata/replay/`
  - `internal/prediction/`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/handoffs/phase-29.md`
- Commands to run before coding:
  - `make validate`
  - `make realtime-quality`
  - `make test`
  - `make smoke`
  - `make test-integration`
  - `git diff --check`
- Known blockers:
  - No real-world agency telemetry or observed stop-time evidence is currently committed.
  - Production-grade ETA quality must not be claimed from small fixtures or replay-only evidence.
- Recommended first implementation slice:
  - Add richer replay fixtures for after-midnight service, frequency-window trips, block continuity, sparse telemetry, noisy GPS, stale/ambiguous inputs, cancellation/alert linkage, and manual overrides over time while preserving visible unknown/withheld outcomes.
