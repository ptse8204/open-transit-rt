# Small Agency Maintenance Guide

This guide explains routine maintenance for a small agency using Open Transit
RT as a self-hosted GTFS and GTFS-Realtime operations platform.

Use the private browser UI first:

```text
/admin/operations/maintenance
```

This guide creates no evidence, contacts no external party, changes no
consumer status, and makes no compliance, adoption, consumer-acceptance,
hosted-service, production-readiness, vendor, SLA, uptime, or ETA-quality
claim.

## Weekly Tasks

| Task | Browser page | What to check | Technical helper needed when |
| --- | --- | --- | --- |
| Check feed health | `/admin/operations/feed-health` | All five configured public route paths, recorded HTTP status, byte count, checksum, validator state, health state, next action. | The page says public checks are missing, stale, blocked, or a proxy/server route needs repair. |
| Review validators | `/admin/operations/validation-health` | Static GTFS validator state, GTFS-Realtime validator state, stale reports, missing tooling, blocked artifacts. | Java, Docker, pinned validator assets, or off-host validation need setup. |
| Review GTFS quality | `/admin/operations/gtfs-quality` | Error vs warning, likely owner, affected files, what to fix first, verification step. | Source GTFS must be edited outside the browser or a large re-import is needed. |
| Review telemetry freshness | `/admin/operations/devices` and `/admin/operations/telemetry` | Latest accepted telemetry time, stale rows, unknown-device or rejected summaries when available, assignment state. | A device token must be rotated, a device is not sending, or an AVL adapter needs changes. |
| Check alerts | `/admin/alerts/console` and `/admin/operations/feed-health` | Active alerts count, empty/non-empty Alerts feed state, stale alerts, next action. | Alert workflow owners need to confirm real service notices or policy. |
| Review maintenance tasks | `/admin/operations/maintenance` | Backup/restore configuration presence, last GTFS import, last five-feed check, telemetry freshness, service health availability. | Any row is blocked, missing, or not available and the browser cannot resolve it. |

## Monthly Tasks

| Task | Browser page | What to check | Technical helper needed when |
| --- | --- | --- | --- |
| Update GTFS | `/admin/operations/gtfs-import` | Source, checksum, byte count, active feed version, import counts, warnings, schedule identity. | Staged comparison, rollback, or a large/scheduled import must be run outside the browser. |
| Review configured feed URLs | `/admin/operations/feed-health` | The five configured public route paths still match the expected root and return current artifacts. | DNS, TLS, reverse proxy, or off-host checks need repair. |
| Review off-host validation | `docs/deployment/off-host-validation.md` | Latest private `.cache/validate-public-feeds/.../summary.json` if a helper ran it. | Validators are missing, blocked, or cannot run on the server. |
| Check backup/restore status | `/admin/operations/maintenance` | Backup and restore-drill configuration presence only, not secret values. | Actual backup/restore commands or a restore drill need to run. |
| Generate support summary when needed | `/admin/operations/maintenance` | Support-bundle instructions and redaction boundaries. | A maintainer needs private diagnostics for a blocker. |

## How To Read Status Labels

- Ready for local review: the private source record is present and currently
  looks usable for local/operator review.
- Needs review: the record exists but the next action should be reviewed.
- Missing: the expected source record is not configured or not observed.
- Blocked: a known blocker prevents the workflow from continuing.
- Unknown: the UI cannot read the source state in this runtime.
- Diagnostic only: useful private context that does not prove production or
  external acceptance.

## What To Do When A Feed Is Empty

Vehicle Positions may be empty when no telemetry has been accepted, telemetry
is stale, devices are unbound, or matching withheld uncertain assignments.
Review Devices, Telemetry, and Telemetry Simulator before changing Trip
Updates or connector settings.

Trip Updates may be empty when there is no active schedule context, no usable
telemetry, withheld/unknown assignments, no deterministic prediction output,
or an external adapter is disabled or failing closed. Empty Trip Updates can
be valid when the safe fallback is to withhold uncertain predictions.

Alerts may be empty when there are no active alerts. That is different from an
Alerts service failure; use Feed Health and the Alerts Console to distinguish
the two.

## Safe Sharing

Public feed URLs can be shared when the agency intends to share them. Private
admin pages, `.cache` diagnostics, support bundles, logs, tokens, database
URLs, raw telemetry, and validator raw reports stay private unless a separate
retained-evidence approval exists.

## Maintenance Principle

Do not turn missing data into OK. If a row says not configured, blocked, or not
available, keep that status visible and either fix the source issue or ask a
technical helper for the next private diagnostic step.
