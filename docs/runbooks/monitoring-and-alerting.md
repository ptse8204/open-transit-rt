# Runbook: Monitoring and Alerting Evidence

This runbook defines the minimum deployment evidence for feed monitoring and alert operations.

Latest captured packets:

- `docs/evidence/captured/local-demo/2026-04-22/monitoring-alert-2026-04-22.md`
- `docs/evidence/captured/oci-pilot/2026-04-24/monitoring-alert-2026-04-24.md`

The local packet captures request-log and scorecard evidence only. The OCI pilot packet records a deployment/operator alert lifecycle for that recorded pilot scope.

## What To Measure

Capture metrics/evidence for at least:

- Public feed endpoint availability.
- Feed freshness (entity/header timestamps).
- Feed generation latency.
- Validator run status over time.
- Stale/unmatched/low-confidence counts where available.
- Service readiness through private/origin `/readyz` checks.
- Backup success/failure.
- Disk space and backup storage growth.
- Database connectivity.
- Telemetry staleness.
- Trip Updates withheld counts and coverage changes.
- Alerts feed availability and active alert lifecycle health.

## What To Alert On

Document active alert conditions, such as:

- Endpoint unavailable.
- Feed freshness breach.
- Elevated validation failures.
- Repeated feed generation errors.
- Backup failure or missing backup evidence.
- Disk or backup storage threshold crossed.
- Database connectivity failure.
- Telemetry staleness threshold crossed.
- Trip Updates withheld/coverage regression.
- Alerts feed unavailable.

Phase 17 feed monitor dry-run:

```sh
ENVIRONMENT_NAME=<environment> \
EVIDENCE_OUTPUT_DIR=/opt/open-transit-rt/evidence/<UTC-date> \
PUBLIC_BASE_URL=https://feeds.example.org \
scripts/pilot-ops.sh feed-monitor --dry-run
```

The live monitor writes `feed-monitor-YYYY-MM-DD.txt` to `EVIDENCE_OUTPUT_DIR`.

Notification destinations are optional placeholders. If `NOTIFY_WEBHOOK_URL` or `NOTIFY_EMAIL_TO` is missing, report `notification not configured`; do not treat that as a feed failure. Never commit real webhook URLs, email credentials, or notification tokens.

## Alert Delivery Proof Pattern

Operators should periodically prove that a monitor can notify the expected destination. This can be a controlled test alert, a scheduled monitor test, or a real alert lifecycle. It does not require hosted monitoring SaaS, Prometheus, or Grafana.

Record:

- Monitor/check name.
- Expected notification destination, redacted if private.
- Test timestamp UTC.
- Delivery result: delivered / not delivered / notification not configured.
- Redacted proof location, such as a private ticket, redacted screenshot, or operator-reviewed summary.
- Follow-up if delivery failed.

If no destination is configured, record `notification not configured` and the operator decision. Missing notification configuration is an operations gap, but it is not itself proof that a public feed failed.

## Capacity Checks

Review capacity at least weekly for long-running pilots:

- Disk space: warning at 75% used, critical at 90% used.
- Database size growth: warning at 25% growth in 30 days, critical at 50% growth in 30 days or when storage pressure appears.
- Backup storage growth: warning at 75% of backup volume, critical at 90%.
- Log growth: warning when logs exceed expected retention size, critical when logs threaten service disk.
- Evidence artifact growth: warning when packet directories grow unexpectedly, critical when raw artifacts include private data or threaten disk.

When thresholds are crossed, remove expired backups/logs according to policy, move raw private evidence to approved operator storage, expand storage, reduce noisy logging, or plan data retention changes. Record the action and redaction review.

## Evidence To Capture

- Dashboard screenshots or exports with timestamps.
- Alert rule definitions and thresholds.
- Notification destination (pager/email/chat).
- Alert delivery proof using the pattern above.
- One real alert lifecycle example:
  - detected
  - acknowledged
  - mitigated
  - resolved
  - post-incident notes

## Output Artifact

Use:

- `docs/evidence/templates/monitoring-alert-template.md`

Store under `docs/evidence/captured/<environment>/monitoring-alert-YYYY-MM-DD.md`.

Evidence labels:

- `feed-monitor-YYYY-MM-DD.txt`: `safe-to-commit-after-review` if it contains only public URLs/statuses.
- Dashboard exports and alert lifecycle summaries: `safe-to-commit-after-review` after redaction.
- Real webhook URLs, notification credentials, private incident links, raw client IP logs, and private infrastructure names: `never-commit`.

## Truthfulness Guardrail

If there is no real alert lifecycle example yet, state that clearly as pending evidence.
