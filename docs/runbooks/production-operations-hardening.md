# Runbook: Production Operations Hardening

This runbook defines the docs-first Phase 28 operations checklist for small-agency pilots and production-directed deployments. It makes day-to-day operations more repeatable without claiming hosted SaaS availability, paid support, SLA coverage, universal production readiness, production multi-tenant hosting, consumer acceptance, CAL-ITP/Caltrans compliance, agency endorsement, marketplace/vendor equivalence, or production-grade ETA quality.

Use this with:

- `docs/runbooks/small-agency-pilot-operations.md`
- `docs/runbooks/backup-and-restore.md`
- `docs/runbooks/monitoring-and-alerting.md`
- `docs/runbooks/validator-evidence.md`
- `docs/upgrade-and-rollback.md`
- `docs/evidence/redaction-policy.md`
- `SECURITY.md`

## Operating Cadence

### Daily Checks

- Public feed availability: fetch `feeds.json`, schedule ZIP, Vehicle Positions, Trip Updates, and Alerts from the public root and record HTTP status, response size, and timestamp.
- Service readiness: check `/readyz` for agency-config, telemetry-ingest, Vehicle Positions, Trip Updates, Alerts, and GTFS Studio on the private/origin side.
- Telemetry freshness: review latest accepted telemetry time, stale vehicle count, unmatched vehicle count, and low-confidence assignment count.
- Trip Updates diagnostics: review withheld count, coverage changes, stale inputs, degraded cases, and cancellation/alert linkage diagnostics.
- Alerts health: confirm Alerts feed availability and review active/published alert lifecycle state when alerts are expected.
- Validator status: review latest schedule and realtime validator results; treat missing validator tooling as an operations finding, not a pass.
- Backup status: confirm the last scheduled backup completed and produced private evidence in the expected operator-owned location.
- Incident queue: check open incidents, prediction review items, backup failures, validator failures, and unresolved alert delivery failures.

### Weekly Checks

- Review monitoring history for public feed availability, feed freshness, readiness failures, validator failures, backup failures, telemetry staleness, Trip Updates withheld/coverage changes, Alerts availability, disk growth, database connectivity, and log growth.
- Confirm alert delivery proof for at least one representative check, or record `notification not configured` if alert destinations are intentionally absent.
- Review redacted evidence summaries that are candidates for commit; keep raw logs, private dumps, credentials, private paths, and private operator artifacts outside the repo.
- Confirm backup retention cleanup is running and storage growth remains within the deployment threshold.
- Review operator access, admin roles, and any pending device-token rotations.

### Monthly Checks

- Run a restore drill into an isolated restore database or isolated restore environment. Do not restore over the live database unless executing an approved incident recovery.
- Verify post-restore public feed fetches and post-restore validator runs before marking the drill complete.
- Review secret inventory and rotation schedule for admin JWT, CSRF, device token pepper, device tokens, DB password, TLS/ACME material, and webhook/notification credentials if used.
- Review capacity thresholds for disk, database size, backup storage, logs, and evidence artifacts.
- Refresh production-directed evidence only when current operator-reviewed artifacts exist, then run `EVIDENCE_PACKET_DIR=<packet> make audit-hosted-evidence`.
- Review operator handover material and update release/version, feed URLs, backup location, restore process, validator cadence, monitoring cadence, known blockers, agency-owned-domain status, and multi-agency limitations.

### After-Upgrade Checks

- Confirm the pre-upgrade backup exists, is private, and has a checksum.
- Record source tag, commit SHA, dirty/clean state, release notes link, artifact checksum if any, and migration status before and after upgrade.
- Run public feed fetch checks for all five public feed paths.
- Run schedule and realtime validator checks after upgrade.
- Review Trip Updates coverage/withheld diagnostics and Alerts availability after upgrade.
- If rollback or restore was needed, complete the restore-event template and retain redacted evidence.

## Capacity Guidance

Operators should set local thresholds that match the host size and agency data volume. The defaults below are starting points for small deployments, not guarantees.

| Area | Warning threshold | Critical threshold | Next action |
| --- | ---: | ---: | --- |
| Disk space | 75% used | 90% used | Remove expired backups/logs, move private evidence off host, expand disk, or reduce retention after policy review. |
| Database size growth | 25% growth in 30 days | 50% growth in 30 days or storage pressure | Inspect telemetry retention, feed history, validation reports, incidents, and indexes; plan archiving or storage expansion. |
| Backup storage | 75% of backup volume | 90% of backup volume | Verify retention cleanup, move older backups to approved private storage, or increase backup storage. |
| Log growth | logs exceed expected retention size | logs threaten service disk | Rotate/compress logs, lower noisy debug logging, and preserve only redacted summaries in evidence packets. |
| Evidence artifacts | packet directory grows unexpectedly | raw artifacts contain private data or threaten disk | Move raw artifacts to private operator storage, commit only redacted summaries, and refresh checksums for committed artifacts. |

When a threshold is crossed, record the detection time, operator, affected environment, action taken, evidence retained, and redaction review. Do not publish private capacity screenshots, private paths, raw logs, or credentials.

## Secret Rotation

Rotate or revoke exposed credentials. Deleting a file from the working tree is not enough when a real secret was exposed; operators must also rotate or revoke the credential, assess whether git history or backups contain it, and record the response.

Required rotation coverage:

- Admin JWT secret: set a new `ADMIN_JWT_SECRET`, temporarily place the prior value in `ADMIN_JWT_OLD_SECRETS` only for the approved transition window, restart admin services, verify old sessions expire, then remove the old secret.
- CSRF secret: set a new `CSRF_SECRET`, restart browser-admin services, verify unsafe cookie-authenticated admin forms require fresh CSRF tokens.
- Device token pepper: rotate only with a planned device-token reissue because existing token hashes depend on the pepper; rebind devices and verify old tokens fail.
- Device tokens: rotate on device loss, operator turnover, vendor change, suspected exposure, or scheduled maintenance; verify telemetry ingest rejects old tokens and accepts only the new binding.
- DB password: change the database role password, update private environment files, restart services, run `/readyz`, and verify backups still run.
- TLS/ACME material: rotate or renew through the deployment-owned TLS process; verify HTTPS certificate metadata and public anonymous feed fetches.
- Webhook/notification credentials: rotate in the notification provider, update private ops environment files, run an alert delivery proof, and verify old webhook URLs or tokens no longer work.
- Phase 15 `.cache` secret findings: do not reuse old local `.cache` credentials; rotate or revoke any real credential found there before further pilot use.

Use `docs/runbooks/templates/secret-rotation-record-template.md` for planned rotations and `docs/runbooks/templates/secret-exposure-incident-template.md` for suspected or confirmed exposures.

## Incident Response Workflow

1. Open the relevant template under `docs/runbooks/templates/`.
2. Record start time, affected environment, affected agency, affected public URLs or services, detection source, operator, and severity.
3. Stabilize public feed availability and operator safety first. Prefer valid but degraded feeds over invalid output, and prefer unknown/withheld realtime state over false certainty.
4. Record timeline entries as actions happen.
5. Retain raw evidence privately; commit only redacted summaries after review.
6. Run follow-up checks after mitigation: public feed fetches, `/readyz`, validators, telemetry freshness, Trip Updates diagnostics, Alerts availability, and backup/restore checks as applicable.
7. Record claim boundary. Validator success, pilot evidence, or restored service does not prove consumer acceptance, compliance, hosted SaaS availability, or production multi-tenant readiness.

## Evidence Refresh Workflow

- Store private raw outputs under deployment-owned `EVIDENCE_OUTPUT_DIR`.
- Commit only operator-reviewed, redacted summaries under `docs/evidence/captured/<environment>/<UTC-date>/`.
- Refresh `SHA256SUMS.txt` or per-file checksums whenever committed artifacts change.
- Run `EVIDENCE_PACKET_DIR=docs/evidence/captured/<environment>/<UTC-date> make audit-hosted-evidence` before calling a packet complete.
- Label evidence as local demo, pilot, hosted/operator, agency-owned-domain, or production-directed. Do not convert pilot evidence into agency-owned-domain proof.
- Preserve Phase 27 agency-boundary language: current backup, restore, export, and evidence workflows are deployment/DB scoped and are not tenant-safe multi-agency operations.

## Operator Handover

Use `docs/runbooks/templates/operator-handover-template.md` before an operator leaves, before a deployment changes owners, and before a pilot becomes a longer-running deployment. Do not include secrets in the handover record. Point to the private secret manager, backup location, evidence location, and access process instead.

## Phase 27 Operations Boundary

Current backup, restore, export, and evidence workflows are deployment/DB scoped. They are not tenant-safe multi-agency workflows. Phase 27 selected isolation tests prove repository-level isolation for selected paths only; they do not prove production multi-tenant operations, hosted SaaS availability, or tenant-safe backup/restore/export/evidence handling.
