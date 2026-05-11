# Phase 67 -- Product Polish, Accessibility, And In-App Help

## Status

In progress. Checkpoint 000001 added this scoped plan for Operations Console
information architecture, accessibility-oriented markup and responsive layout,
contextual in-app help, and phase closeout. Checkpoint 000002 improved the
private Operations Console information architecture with operator-intent
navigation groups, active page state, route-stability tests, and a current
Phase 61+ roadmap link in the Launchpad decision gate. Phase 67 must stay
inside the private Operations Console product surface. It must not create retained
evidence, write protected evidence paths, contact external parties, change
consumer statuses, change public feed URLs, change telemetry ingest,
change GTFS-RT protobuf semantics, change validator execution semantics,
change connector manifest schema, change prediction adapter behavior, weaken
auth boundaries, or claim compliance, agency adoption, consumer acceptance,
hosted SaaS, production readiness, vendor compatibility, hardware
certification, production AVL reliability, production-grade ETA quality, SLA,
or public launch completion.

## Goal

Make Open Transit RT feel more coherent and understandable for a small agency
operator who opens the private Operations Console without reading phase
history. The console should answer where to start, where to troubleshoot, what
each section means, and which statements are only private readiness or
supporting signals.

## Checkpoints

- Completed: `Phase 67 -- Checkpoint 000001: add product polish and accessibility plan`
- Completed: `Phase 67 -- Checkpoint 000002: improve operations console information architecture`
- Planned: `Phase 67 -- Checkpoint 000003: improve accessibility and mobile layout`
- Planned: `Phase 67 -- Checkpoint 000004: implement in-app help system`
- Planned: `Phase 67 -- Checkpoint 000005: close product polish accessibility and help`

## Existing State

- The Operations Console is server-rendered from `cmd/agency-config`, mostly in
  `operations.go` plus feature-specific files.
- Existing private routes already cover dashboard, Launchpad, setup wizard,
  Connector Hub, connector test instructions, browser GTFS import, feed
  health, readiness checklist v2, feeds/validation, GTFS quality, validator
  health, reliability, telemetry, telemetry simulator, devices, consumers,
  evidence links, setup, checklist, GTFS Studio, and Alerts Console.
- Existing tests already verify private route scope, agency scoping, GET-only
  behavior for read-only pages, no-store headers, HTML escaping, and many
  claim-boundary strings.
- Navigation is currently a flat list. It exposes all major routes, but it is
  hard to scan on small screens and does not group work by operator intent.
- Shared layout has basic responsive CSS, but it can better support landmarks,
  skip links, active navigation, focus visibility, and table/form ergonomics.
- Help text exists throughout individual pages, but there is no consistent
  contextual help model for GTFS, GTFS-RT, connectors, readiness, validators,
  telemetry, and claim/evidence boundaries.

## Checkpoint Scope

### Operations Console Information Architecture

- Keep all existing route paths stable.
- Replace or supplement the flat navigation with grouped navigation that
  reflects operator intent:
  - Start;
  - GTFS and feeds;
  - realtime operations;
  - connectors;
  - readiness and diagnostics;
  - records and boundaries.
- Add an active navigation signal with `aria-current="page"` for the current
  section where practical.
- Keep navigation plain language and private/admin scoped.
- Update any stale Operations Console links that point operators to historical
  Post-60 pages when the Phase 61+ roadmap is the current forward product
  roadmap.

### Accessibility And Mobile Layout

- Improve the shared Operations Console layout with semantic landmarks,
  `lang="en"`, viewport metadata, skip-to-main link, labeled navigation, a
  stable main content region, and visible keyboard focus treatment.
- Keep tables readable on narrow screens without changing their data model.
- Improve form button types and label ergonomics where this does not change
  POST semantics or authorization behavior.
- Do not claim WCAG compliance or certified accessibility. Phase 67 may say
  accessibility-oriented markup, keyboard-friendly navigation, and responsive
  layout were improved.

### In-App Help System

- Add private GET-only help routes:
  - `/admin/operations/help`
  - `/admin/operations/help.json`
- Add a bounded static or derived help model for:
  - GTFS;
  - GTFS-RT;
  - connectors;
  - readiness;
  - validators;
  - telemetry;
  - claim and evidence boundaries.
- Add contextual help in the shared layout so each major section can point to
  the most relevant help topics.
- Keep help private, read-only, agency-scoped, no-store, and all-false for
  claim flags.
- Do not execute commands, read `.cache`, write evidence, collect secrets,
  contact external systems, or move any consumer tracker status from help
  routes.

### Closeout

- Mark this phase complete only after implementation checkpoints pass focused
  tests and closeout validation.
- Add `docs/handoffs/phase-67.md`.
- Update status and handoff docs to point next work to Phase 68+ authorization
  gated evidence scaffolding, not to evidence collection.

## Files Expected To Change

- `cmd/agency-config/operations.go`
- optional `cmd/agency-config/operations_navigation.go`
- optional `cmd/agency-config/operations_help.go`
- `cmd/agency-config/main_test.go`
- `docs/phase-67-product-polish-accessibility-in-app-help.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-67.md`
- `docs/roadmap-status.md`
- optional docs index or historical roadmap pointer updates where wording
  still obscures the current Phase 61+ roadmap.

Protected paths remain untouched:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/**`
- `db/migrations/**`
- `go.mod`
- `go.sum`

## Non-Goals

- No new frontend framework.
- No database migration.
- No public route or public feed URL change.
- No telemetry ingest contract change.
- No GTFS-RT protobuf semantic change.
- No validator execution semantic change.
- No connector manifest schema change.
- No prediction adapter behavior change.
- No auth boundary or public/private route boundary weakening.
- No arbitrary dynamic plugin loading.
- No external contact.
- No retained evidence writes.
- No consumer status changes.
- No evidence collection from Phase 68+ tracks.
- No claim of CAL-ITP/Caltrans compliance, consumer submission/review/
  acceptance/listing/display/ingestion, agency adoption or approval,
  final-root proof, hosted SaaS, paid support, SLA, universal production
  readiness, production multi-tenant hosting, vendor compatibility, hardware
  certification, production AVL reliability, production-grade ETA quality, or
  public launch completion.

## Validation Plan

Run after focused UI/help checkpoints:

- `git diff --check`
- `go test ./cmd/agency-config`
- `make check`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json`
- `git diff --exit-code -- docs/evidence/captured`
- `git diff --exit-code -- db/migrations go.mod go.sum`

Run at closeout:

- `git diff --check`
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

If Docker, Java, network, validator tooling, or pinned images are unavailable,
record the exact blocker and continue with checks that do not depend on that
environment.

## Claim Boundary

Phase 67 may say the private Operations Console is easier to navigate, more
responsive, more keyboard-friendly, and more understandable through contextual
help. It may say the app supports CAL-ITP-style readiness workflows as private
readiness review.

Phase 67 must not say or imply compliance, consumer acceptance, agency
approval, hosted service availability, production readiness, vendor
compatibility, hardware certification, production AVL reliability,
production-grade ETA quality, public launch completion, or accessibility
certification.

Help content, feed health, readiness rows, validator state, telemetry state,
and connector conformance are private or supporting signals only. They do not
create retained evidence and must not move consumer tracker records beyond
`prepared`.

## Rollback Path

Phase 67 should remain private UI, tests, and documentation work. If rollback
is needed, revert the specific checkpoint commit that added navigation,
accessibility markup, help routes, tests, or status docs. Because Phase 67
does not change DB schema, public feed URLs, telemetry ingest, GTFS-RT
semantics, validator execution semantics, connector manifests, prediction
adapter behavior, evidence files, or consumer tracker statuses, no data
rollback should be required.
