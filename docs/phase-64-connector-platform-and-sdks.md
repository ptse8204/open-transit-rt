# Phase 64 -- Connector Platform And SDKs

## Status

In progress. Checkpoint 000001 added this plan and kept implementation scoped
to private connector-platform UX, local synthetic examples, and offline
conformance checks. Checkpoint 000002 added the private Connector Hub manifest
registry UI and bounded JSON registry model from committed synthetic example
manifests. Checkpoint 000003 added private connector test instructions for
fixed offline checks, without backend command execution.

Phase 64 turns the existing connector/plugin architecture into a visible,
testable developer platform for telemetry, prediction, validator,
monitoring/export, and consumer/discovery workflows. It must preserve the
existing connector manifest schema unless a later checkpoint identifies a
narrow, reviewed need; keep connectors as optional sidecars, command adapters,
manifests, or connector processes; and avoid dynamic backend plugin loading,
real credentials, real vendor payloads, retained evidence, consumer status
changes, or vendor compatibility claims.

In Open Transit RT, a plugin is an optional sidecar, command adapter, manifest,
or connector process. It is not arbitrary dynamic code loaded into the
backend.

## Checkpoints

- Completed: `Phase 64 -- Checkpoint 000001: add connector platform and SDK plan`
- Completed: `Phase 64 -- Checkpoint 000002: implement connector manifest registry UI`
- Completed: `Phase 64 -- Checkpoint 000003: implement connector test runner UI`
- Planned: `Phase 64 -- Checkpoint 000004: improve telemetry connector SDK examples`
- Planned: `Phase 64 -- Checkpoint 000005: improve prediction connector SDK examples`
- Planned: `Phase 64 -- Checkpoint 000006: improve monitoring export connector examples`
- Planned: `Phase 64 -- Checkpoint 000007: close connector platform and SDKs`

## Checkpoint Scope

### Connector Manifest Registry UI

- Add a private Operations Console registry view from safe committed connector
  manifests under `examples/connectors/*/connector.json`.
- Add a private JSON export for the same bounded registry model.
- Show connector ID, type, display name, mode, disabled-by-default state,
  contract summaries, failure behavior, redaction posture, docs link,
  conformance case count, and explicit claim boundaries.
- Link the registry from the existing Connector Hub and Operations Console
  navigation where useful.
- Treat unreadable or invalid committed examples as registry diagnostics, not
  runtime plugin loading.
- Do not accept user-uploaded manifests, execute manifest commands, load Go
  plugins, mutate connector state, add public routes, or change the manifest
  schema.

Checkpoint 000002 added `internal/connectors.LoadExampleRegistry` and surfaced
its bounded registry model inside the existing private Connector Hub JSON and
HTML views. The registry reads only committed
`examples/connectors/*/connector.json` manifests, records diagnostics for
missing or invalid committed examples, and exposes bounded summaries of
contracts, failure behavior, redaction posture, claim boundaries, docs links,
and synthetic conformance cases. It does not accept runtime paths, load
plugins, execute commands, write files, contact networks, change schemas, or
mutate status.

### Connector Test Runner UI

- Add a private Operations Console page that explains the allowlisted connector
  checks already available in the repository:
  - `make external-connection-check`
  - `make adapter-conformance`
  - `make test-connector-examples`
- Implement generated copyable instructions only. Do not add backend command
  execution in Phase 64 without a later explicit re-approval and a new scoped
  checkpoint.
- Show what each check validates, what it does not prove, and the safe next
  action when a check fails.
- Keep the page private and diagnostics-only. It must not contact external
  parties, run arbitrary commands, write retained evidence, change consumer
  statuses, or claim compatibility, compliance, production readiness, or
  ETA quality.

Checkpoint 000003 added `/admin/operations/connectors/tests` and
`/admin/operations/connectors/tests.json` as private GET-only generated
instructions for fixed offline checks. The page lists manifest/example checks,
full adapter conformance, targeted conformance sections, and connector example
tests. It records false claim flags for backend command execution, manifest
command execution, external network contact, evidence creation, consumer
status changes, compatibility, compliance, production readiness, and ETA
quality. It does not execute commands from the web request, read
manifest-provided commands, capture output, write files, contact external
systems, or mutate status.

### Telemetry Connector SDK Examples

- Improve the synthetic telemetry HTTP poller and CSV replay examples so a
  developer can understand input mapping, fail-closed decisions, redaction, and
  dry-run output without reading phase history.
- Keep fixtures synthetic and generic. Do not add real AVL/vendor payloads,
  real device identifiers, credentials, private endpoints, or vendor-specific
  compatibility claims.
- Preserve the authenticated `/v1/telemetry` contract and keep example sends
  disabled by default.

### Prediction Connector SDK Examples

- Improve the synthetic predictor sidecar stub so it documents the prediction
  adapter boundary, Vehicle Positions independence, failure fallback, shadow
  or dry-run diagnostics, and withheld/empty-output behavior.
- Do not add a named predictor integration, change default Trip Updates
  behavior, change `internal/prediction.Adapter`, or claim production-grade ETA
  quality.

### Validator And Monitoring/Export Connector Examples

- Improve validator connector guidance around server-owned allowlisted
  validator IDs, offline fixtures, normalized reports, and failure states.
- Improve the monitoring/export example so redaction and no-send defaults are
  easier to adapt for deployment-owned monitoring.
- Keep examples local, synthetic, generic, and disabled by default. Do not add
  raw validator commands, arbitrary argv, webhook sends, real monitoring
  destinations, SLA/uptime claims, or retained evidence writes.

## Non-Goals

- No database migrations.
- No public feed URL changes.
- No change to GTFS-RT protobuf semantics.
- No telemetry ingest contract changes.
- No validator execution semantic changes.
- No connector manifest schema changes unless a later checkpoint documents and
  tests a narrow required adjustment.
- No prediction adapter behavior changes.
- No auth boundary changes.
- No public/private route boundary changes.
- No arbitrary dynamic backend plugin loading.
- No real credentials, real vendor payloads, private endpoints, or webhook
  destinations.
- No retained evidence.
- No `docs/evidence` writes.
- No consumer status changes.
- No external contact.
- No CAL-ITP/Caltrans compliance, agency approval/adoption, final-root,
  consumer submission/review/acceptance/listing/display/ingestion, hosted
  SaaS, paid support, SLA/uptime, production-readiness, public-launch, vendor
  compatibility, hardware-certification, production AVL reliability, or
  production-grade ETA claim.

## Files Expected To Change

- `cmd/agency-config/main.go`
- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_connectors.go`
- new or existing `cmd/agency-config/operations_connector_registry.go`
- new or existing `cmd/agency-config/operations_connector_tests.go`
- `cmd/agency-config/main_test.go`
- `examples/connectors/telemetry-http-poller/README.md`
- `examples/connectors/telemetry-http-poller/*`
- `examples/connectors/telemetry-csv-replay/README.md`
- `examples/connectors/telemetry-csv-replay/*`
- `examples/connectors/predictor-sidecar-stub/README.md`
- `examples/connectors/predictor-sidecar-stub/*`
- `examples/connectors/monitoring-export/README.md`
- `examples/connectors/monitoring-export/*`
- optional synthetic validator example files under `examples/connectors/`
- `internal/connectors/examples_test.go` when example inventory changes
- `docs/connectors/plugin-contract.md`
- `docs/integration-adapter-kit.md`
- `docs/tutorials/external-adapter-conformance.md`
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
- `go test ./internal/connectors`
- `go test ./cmd/adapter-conformance`
- `go test ./examples/connectors/...`
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

Run DB-backed integration checks only if a later checkpoint touches behavior
that depends on the database. If Java, Docker, validator tooling, or pinned
images are unavailable, record the environment blocker and continue with
non-environment-dependent checks.

## Claim Boundary

Phase 64 may say the repository is connector-ready or plugin-friendly through
sidecars, manifests, command adapters, and conformance tests. It may say the
registry and examples are private diagnostics, synthetic conformance aids, or
supporting signals.

Phase 64 must not say or imply that Open Transit RT is CAL-ITP/Caltrans
compliant, accepted by a consumer, adopted or approved by an agency, final-root
approved, a hosted SaaS, SLA-backed, universally production ready,
vendor-compatible, hardware certified, public-launch complete, production AVL
reliable, or production-grade ETA proven.

Passing connector checks remains a local quality signal only. It is not
external evidence and must not move consumer tracker records beyond
`prepared`.

## Rollback Path

Phase 64 should remain code, docs, and synthetic-example work. If rollback is
needed, revert the specific checkpoint commit that added the route, UI,
example, fixture, or documentation change. Public feed URLs, DB schema,
telemetry ingest, GTFS-RT semantics, validator execution semantics, connector
manifest schema, prediction adapter behavior, auth boundaries, protected
evidence paths, and consumer tracker statuses should remain untouched.
