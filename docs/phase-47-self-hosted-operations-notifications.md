# Phase 47 — Self-Hosted Operations Notifications

Phase 47 adds a private local notification draft for self-hosted operators. It
summarizes already-created `validator-health` and `deployment-doctor`
summaries, then writes bounded local files under `.cache/operations-notify/` by
default.

This phase does not send notifications. It does not call webhooks, email,
admin routes, validators, databases, Docker, applications, consumers, or public
services.

## Operator Surface

- Script: `scripts/operations-notify.sh`
- Make target: `make operations-notify`
- Default output: `.cache/operations-notify/<timestamp>`

The script writes exactly:

- `summary.json`
- `summary.md`
- `manifest.json`
- `manifest.md`
- `notification.txt`

`--dry-run` writes the same five local draft files and requires no source
summaries, network, webhook, email, database, Docker, admin token, app, or
running service.

## Source Model

Default source discovery selects the latest lexicographic timestamp-named
summary from:

- `.cache/validator-health/*/summary.json`
- `.cache/deployment-doctor/*/summary.json`

Explicit source paths may be supplied with `VALIDATOR_HEALTH_SUMMARY` and
`DEPLOYMENT_DOCTOR_SUMMARY`. Sources must stay under `.cache` unless
`ALLOW_UNIGNORED_SOURCE_DIR=true`; evidence-like paths are always rejected.
The script rejects symlink source files or directories and reads only
documented `summary.json` files plus safe manifest metadata where needed. It
does not recursively copy raw reports, private logs, environment files, or
source JSON contents.

Each source file is capped by `MAX_SOURCE_BYTES`, default `5242880`.
Oversized, missing, or malformed sources produce safe next actions in default
mode. `STRICT_OPERATIONS_NOTIFY=true` exits nonzero for missing, malformed,
oversized, blocked, or unhealthy source states.

## Severity

Notification severity uses this precedence:

```text
blocked > needs_review > info > unknown
```

Missing or malformed sources are `needs_review` by default and `blocked` in
strict mode. The script does not invent health status for missing or malformed
inputs.

## Output Contract

`summary.json` includes:

- `generated_at`
- `source_summaries`
- `notification`
- `destinations`
- `counts`
- `next_actions`
- `overflow_count`
- false claim flags

The exact false flags are:

- `external_evidence_created=false`
- `consumer_statuses_changed=false`
- `compliance_claimed=false`
- `production_readiness_claimed=false`
- `hosted_saas_claimed=false`
- `agency_adoption_claimed=false`
- `consumer_acceptance_claimed=false`
- `vendor_compatibility_claimed=false`
- `production_grade_eta_claimed=false`
- `notification_sent=false`

`notification.not_sent` is always `true`. Destination values are never written;
only `destinations.webhook_present` and `destinations.email_present` booleans
are recorded.

Each `next_actions` item uses only:

- `source`
- `severity`
- `title`
- `action`
- `count`
- `overflow_count`

Lists are capped. `overflow_count` records omitted source items. `summary.json`,
`summary.md`, and `notification.txt` contain counts and bounded derived text
only, not raw source objects.

`notification.txt` begins with:

```text
DRAFT — NOT SENT
```

It states that the draft was not sent to webhook, email, consumers, agency, or
public service. It also states that the draft is not evidence, not compliance
proof, not production readiness proof, and not consumer acceptance.

## Privacy And Redaction

Generated output is scanned before completion. The scan fails on value-shaped
secrets or private values including bearer headers, cookie values, admin
sessions, CSRF values, database URLs, password-bearing Postgres URLs, private
keys, token/secret/password assignments, token-bearing webhook URLs, private
absolute paths, and raw report/stdout/stderr/argv fields copied as data.

The script may mention categories excluded from output, such as raw reports,
tokens, secrets, and evidence packets. Those category names alone are not a
redaction failure.

Terminal output is limited to a repo-relative `.cache` output path or a
redacted name. It does not print source contents, destination values, absolute
private paths, tokens, database URLs, or raw report data.

## Boundaries

Phase 47 output is private operator diagnostics only. It is not:

- evidence;
- final-root evidence;
- a consumer submission artifact;
- consumer acceptance or ingestion proof;
- CAL-ITP/Caltrans compliance proof;
- agency adoption or approval proof;
- hosted SaaS proof;
- production readiness proof;
- vendor compatibility proof;
- production-grade ETA proof;
- a public API;
- a publish blocker.

No consumer statuses are changed. No GTFS files are auto-edited.
