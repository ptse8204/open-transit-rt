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

## What To Alert On

Document active alert conditions, such as:

- Endpoint unavailable.
- Feed freshness breach.
- Elevated validation failures.
- Repeated feed generation errors.

Phase 17 feed monitor dry-run:

```sh
ENVIRONMENT_NAME=<environment> \
EVIDENCE_OUTPUT_DIR=/opt/open-transit-rt/evidence/<UTC-date> \
PUBLIC_BASE_URL=https://feeds.example.org \
scripts/pilot-ops.sh feed-monitor --dry-run
```

The live monitor writes `feed-monitor-YYYY-MM-DD.txt` to `EVIDENCE_OUTPUT_DIR`.

Notification destinations are optional placeholders. If `NOTIFY_WEBHOOK_URL` or `NOTIFY_EMAIL_TO` is missing, report `notification not configured`; do not treat that as a feed failure. Never commit real webhook URLs, email credentials, or notification tokens.

## Evidence To Capture

- Dashboard screenshots or exports with timestamps.
- Alert rule definitions and thresholds.
- Notification destination (pager/email/chat).
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
