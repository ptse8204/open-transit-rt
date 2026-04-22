# Monitoring and Alerting Evidence

- Environment: `local-demo`
- Capture date (UTC): 2026-04-22
- Operator: Codex local run

## Metrics and Dashboards

- Availability evidence: service request logs in `artifacts/logs/` show 200 responses for public feed fetches and 401 responses for anonymous protected routes during the local demo.
- Freshness evidence: local JSON debug artifacts from the demo recorded `generated_at` timestamps for Vehicle Positions, Trip Updates, and Alerts; the final scorecard recorded health timestamps for Trip Updates and Alerts.
- Feed latency evidence: request logs include `duration_ms` for public feed requests in local services.
- Validator trend evidence: validation JSON records exist for schedule and all three realtime feeds, but all validator runs failed in this environment.

## Alert Rules

- Rule set reference: Missing.
- Notification destination(s): Missing.
- Threshold summary: Missing.

## Real Alert Lifecycle Example

- Alert ID/title: Missing.
- Detected at (UTC): Missing.
- Acknowledged at (UTC): Missing.
- Mitigated at (UTC): Missing.
- Resolved at (UTC): Missing.
- Root cause summary: Missing.
- Follow-up action: Configure a real monitoring stack with alert rules, delivery destination, and lifecycle retention.

## Evidence Links

- Dashboard export or screenshot: Missing.
- Incident/ticket reference: Missing.
- Local service logs: `artifacts/logs/agency-config.log`, `artifacts/logs/feed-vehicle-positions.log`, `artifacts/logs/feed-trip-updates.log`, and `artifacts/logs/feed-alerts.log`.

## Blocker

This repository currently exposes request logs, readiness checks, and optional Prometheus-format metrics, but no Prometheus server, Grafana dashboard, alert rules, notification destination, or incident system is configured in this workspace. A real alert lifecycle cannot be fabricated from local request logs.
