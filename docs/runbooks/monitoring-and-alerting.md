# Runbook: Monitoring and Alerting Evidence

This runbook defines the minimum deployment evidence for feed monitoring and alert operations.

Latest captured packet:

- `docs/evidence/captured/local-demo/2026-04-22/monitoring-alert-2026-04-22.md`

The local packet captures request-log and scorecard evidence only. A real monitoring dashboard, alert rules, notification destination, and alert lifecycle remain missing.

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

## Truthfulness Guardrail

If there is no real alert lifecycle example yet, state that clearly as pending evidence.
