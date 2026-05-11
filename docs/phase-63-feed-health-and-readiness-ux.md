# Phase 63 -- Feed Health And Readiness UX

## Status

Closed. Checkpoint 000001 added this plan. Checkpoint 000002 added the private
feed-health dashboard and JSON export. Checkpoint 000003 added the private
readiness checklist v2 page and JSON export. Checkpoint 000004 closed the
phase after validation and claim-boundary review.

Phase 63 makes feed health, validator state, freshness, and CAL-ITP-style
readiness easier to understand inside the private Operations Console. It must
reuse existing feed discovery, validation-health, reliability, GTFS quality,
and Trip Updates diagnostics models; preserve public feed URLs and runtime
contracts; avoid schema changes; and avoid compliance, SLA, uptime,
production-readiness, consumer-acceptance, or evidence claims.

## Checkpoints

- Completed: `Phase 63 -- Checkpoint 000001: add feed health and readiness UX plan`
- Completed: `Phase 63 -- Checkpoint 000002: implement feed health dashboard`
- Completed: `Phase 63 -- Checkpoint 000003: implement readiness checklist v2`
- Completed: `Phase 63 -- Checkpoint 000004: close feed health and readiness UX`

## Checkpoint Scope

### Feed Health Dashboard

- Added a private Operations Console feed-health dashboard at
  `/admin/operations/feed-health`.
- Added private JSON output at `/admin/operations/feed-health.json`.
- Shows exactly five operator-facing rows:
  - `feeds.json`
  - `schedule`
  - `vehicle_positions`
  - `trip_updates`
  - `alerts`
- For each row, shows:
  - plain-language status;
  - what the status means;
  - freshness or generated/checked time where available;
  - validator context where applicable;
  - feed health or diagnostic context where available;
  - next action;
  - what the row does not prove.
- Link feed-health dashboard from the Operations Console dashboard, Launchpad,
  setup wizard, feeds page, and readiness page.
- Keep the dashboard private and read-only. It creates no retained evidence,
  changes no consumer status, and makes no compliance, SLA, uptime,
  production-readiness, public-launch, consumer-acceptance, final-root, agency
  approval/adoption, vendor-compatibility, hardware-certification,
  production-AVL-reliability, or ETA-quality claim.

### Readiness Checklist V2

- Added private readiness checklist v2 at `/admin/operations/readiness`.
- Added private JSON output at `/admin/operations/readiness.json` that returns
  only the v2 readiness model and false claim flags.
- The v2 model includes rows for:
  - feed discovery and metadata;
  - plain-language feed health;
  - static GTFS quality;
  - Vehicle Positions readiness;
  - Trip Updates adapter boundary;
  - Service Alerts readiness;
  - validator health;
  - operations reliability diagnostics;
  - telemetry and device setup;
  - operations scorecard;
  - consumer prepared tracker.
- Each row includes:
  - readiness item;
  - status;
  - current signal;
  - what this means;
  - why it matters;
  - what to do next;
  - what it does not prove;
  - source route/doc links where useful.
- Keeps CAL-ITP-style readiness wording bounded to readiness workflows and
  supporting signals.
- Keeps consumer tracker states prepared-only unless retained target-originated
  evidence exists.

## Non-Goals

- No database migrations.
- No public feed URL changes.
- No change to GTFS-RT protobuf semantics.
- No change to telemetry ingest.
- No change to validator execution semantics.
- No retained evidence.
- No `docs/evidence` writes.
- No consumer status changes.
- No external contact.
- No SLA/uptime proof.
- No CAL-ITP/Caltrans compliance claim.
- No agency approval, agency adoption, consumer acceptance, hosted SaaS,
  vendor compatibility, production-readiness, public-launch, production AVL
  reliability, or production-grade ETA claim.

## Closeout

Phase 63 closed with private Operations Console routes only:

- `/admin/operations/feed-health`
- `/admin/operations/feed-health.json`
- `/admin/operations/readiness`
- `/admin/operations/readiness.json`

The implementation reuses existing private discovery, validation-health,
reliability, GTFS quality, Trip Updates diagnostics, telemetry, device,
scorecard, and consumer prepared-tracker signals. It adds no public route,
schema change, public feed URL change, telemetry ingest change, validator
execution change, evidence write, consumer status change, or external contact.

The closeout boundary is unchanged: Phase 63 supports CAL-ITP-style readiness
workflows and private readiness signals only. It does not claim CAL-ITP or
Caltrans compliance, final-root proof, agency approval/adoption, consumer
submission/review/acceptance/listing/display/ingestion, hosted SaaS,
service-level or uptime proof, vendor compatibility, hardware certification,
production readiness, production AVL reliability, public launch, or
production-grade ETA quality.

## Closeout Validation

Closeout validation for Checkpoint 000004:

- `git diff --check`
- `go test ./cmd/agency-config`
- `make check`
- `make test`
- `make external-connection-check`
- `make adapter-conformance`
- `make test-connector-examples`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json`
- `git diff --exit-code -- docs/evidence/captured`
- `git diff --exit-code -- db/migrations go.mod go.sum`
- `docker compose -f deploy/docker-compose.yml config`

Protected evidence/status paths remained unchanged. All seven consumer and
aggregator targets remained `prepared`.

## Files Expected To Change

- `cmd/agency-config/main.go`
- `cmd/agency-config/operations.go`
- new or existing `cmd/agency-config/operations_feed_health.go`
- new or existing `cmd/agency-config/operations_readiness_v2.go`
- `cmd/agency-config/operations_launchpad.go`
- `cmd/agency-config/operations_setup_wizard.go`
- `cmd/agency-config/main_test.go`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- this phase file

Protected paths remain untouched:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/**`
- `db/migrations/**`
- `go.mod`
- `go.sum`

## Validation Plan

- `git diff --check`
- `go test ./cmd/agency-config`
- `make check`
- `make test`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json`
- `git diff --exit-code -- docs/evidence/captured`

Run broader checks when relevant and available:

- `make external-connection-check`
- `make adapter-conformance`
- `make test-connector-examples`
- `docker compose -f deploy/docker-compose.yml config`

## Rollback Path

Phase 63 should remain code-and-docs only. If rollback is needed, revert the
specific checkpoint commit that added the route/model/template/test/docs
changes. Public feed URLs, DB schema, telemetry ingest, validator execution,
prediction adapters, connector manifests, and consumer tracker statuses should
remain untouched.
