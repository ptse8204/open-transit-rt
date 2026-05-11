# Phase 65 -- Operator Workflow And Data Quality UX

## Status

In progress. Checkpoint 000001 added this plan and kept implementation scoped
to private operator workflow UX, existing device credential behavior, generated
telemetry simulator guidance, and GTFS quality fix guidance. Checkpoint 000002
improved the private device and vehicle onboarding UI. Checkpoint 000003 added
the private telemetry simulator guide UI. Phase 65 must
preserve telemetry ingest contracts, device credential semantics, GTFS import
and publish boundaries, validator execution semantics, public feed URLs,
consumer tracker state, protected evidence paths, and unsupported-claim
boundaries.

## Goal

Make day-to-day agency operations easier for small teams by improving device
and vehicle onboarding, synthetic telemetry simulator guidance, GTFS quality
triage, and operator troubleshooting. The work should answer "what should I do
next?" without requiring operators to read phase history or raw diagnostics.

## Checkpoints

- Completed: `Phase 65 -- Checkpoint 000001: add operator workflow and data quality UX plan`
- Completed: `Phase 65 -- Checkpoint 000002: implement device and vehicle onboarding UI`
- Completed: `Phase 65 -- Checkpoint 000003: implement telemetry simulator UI`
- Planned: `Phase 65 -- Checkpoint 000004: implement GTFS quality fix guidance UI`
- Planned: `Phase 65 -- Checkpoint 000005: close operator workflow and data quality UX`

## Existing State

- `/admin/operations/devices` lists device bindings and supports the existing
  admin-only rotate/rebind flow. The one-time token is shown only on POST and
  not repeated on GET.
- `/admin/operations/telemetry` lists latest accepted telemetry by vehicle and
  joins current assignment summaries when available. It omits raw payloads,
  token fields, full score details, and private debug blobs.
- `cmd/telemetry-simulator` and `scripts/telemetry-simulator.sh` run
  synthetic scenarios through authenticated `/v1/telemetry`, with dry-run mode
  and private `.cache/telemetry-simulator/` diagnostics.
- `/admin/operations/gtfs-quality` summarizes canonical static validator and
  internal importer issues, caps raw notice rendering, and allows an admin-only
  rerun of the server-owned static validator mapping.
- `docs/tutorials/device-token-lifecycle.md`,
  `docs/tutorials/telemetry-simulator-and-device-trial.md`, and
  `docs/tutorials/gtfs-validation-triage.md` already explain the command-line
  paths and claim boundaries.

## Checkpoint Scope

### Device And Vehicle Onboarding UI

- Improve `/admin/operations/devices` so a small agency can see binding state,
  latest telemetry for each binding when available, freshness, and a plain
  next action.
- Keep rotate/rebind as the existing mutation path. Do not add token lookup,
  token recovery, token replay, or token display after the initial POST.
- Hide or clearly disable the rotate/rebind form for non-admin users because
  the handler already restricts POST to admins.
- Make the page clearer about when to create a new `device_id`, when to
  rotate/rebind an existing device, when to move a device to a different
  `vehicle_id`, and when to inspect telemetry freshness.
- Add private JSON only if it is useful for the same bounded device-onboarding
  model. Any JSON must exclude token values, token hashes, raw telemetry
  payloads, private operator notes, and vendor identifiers.
- Preserve admin-only POST, role checks, agency scoping, CSRF behavior, audit
  logging, and existing `internal/devices.Store` semantics.

Checkpoint 000002 added guided onboarding use-case cards to
`/admin/operations/devices`, hid the rotate/rebind form from non-admin users,
and replaced the raw binding table with per-device rows showing credential
dates, latest accepted telemetry freshness, assignment summary, and next
action. It derives those rows from existing device binding, latest telemetry,
and assignment summaries only. It does not expose token values after the
one-time POST, token hashes, raw telemetry payloads, private debug fields, or
hardware-specific identifiers.

### Telemetry Simulator UI

- Add or improve a private Operations Console page that guides operators
  through synthetic simulator scenarios and generated commands.
- Prefer generated instructions and scenario summaries over backend command
  execution. Do not run shell commands, start servers, send telemetry, create
  `.cache` output, or capture simulator output from a web request.
- Surface available synthetic scenarios, prerequisites, dry-run mode, optional
  matcher diagnostics, target URL rules, device-token handling, and where
  private diagnostics are written when the CLI is run.
- Link the page from Dashboard, Launchpad, setup, devices, telemetry, and
  readiness where it helps the operator workflow.
- Keep simulator results as local/private diagnostics only. Synthetic simulator
  runs must not be described as real fleet reliability, real vendor AVL
  compatibility, real realtime data, production-grade ETA quality, compliance,
  or consumer acceptance evidence.

Checkpoint 000003 added `/admin/operations/telemetry-simulator` and
`/admin/operations/telemetry-simulator.json` as private GET-only guide
surfaces. The UI reads committed synthetic scenario metadata from
`testdata/telemetry-simulator`, shows target rules, credential-handling
guidance, private diagnostics policy, scenario requirements, and copyable
operator-shell commands. It does not execute commands, send telemetry, collect
device credentials, read `.cache` diagnostics, expose raw scenario payloads,
write evidence, or change consumer statuses.

### GTFS Quality Fix Guidance UI

- Improve `/admin/operations/gtfs-quality` so common validator/importer issue
  groups include practical next actions for operators.
- Make it clearer which issues usually need source GTFS re-export, GTFS Studio
  review, service-calendar review, shape/stop-time cleanup, block/frequency
  review, or technical escalation.
- Keep the page diagnostic and advisory. Do not auto-edit GTFS, mutate drafts,
  publish schedules, infer agency approval, or convert validator output into a
  compliance or consumer-readiness claim.
- Preserve the existing admin-only validator rerun path, server-owned
  validator ID mapping, body cap, private-output redaction, and stale-result
  handling.
- Update operator docs only where they reduce ambiguity between importer
  errors, canonical validator notices, GTFS Studio edits, and source-system
  fixes.

### Closeout

- Update status and handoff docs to mark Phase 65 complete only after
  implementation checkpoints pass validation and claim-boundary review.
- Record protected-path, consumer tracker, and unsupported-claim results.
- Point the next active phase to Phase 66.

## Files Expected To Change

- `cmd/agency-config/main.go`
- `cmd/agency-config/main_test.go`
- `cmd/agency-config/operations.go`
- optional new `cmd/agency-config/operations_devices.go`
- optional new `cmd/agency-config/operations_telemetry_simulator.go`
- optional new `cmd/agency-config/operations_gtfs_quality_guidance.go`
- `cmd/telemetry-simulator/main.go` only if read-only scenario metadata
  helpers are needed and no CLI behavior changes are required
- `docs/tutorials/device-token-lifecycle.md`
- `docs/tutorials/telemetry-simulator-and-device-trial.md`
- `docs/tutorials/gtfs-validation-triage.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/roadmap-status.md`
- this phase file
- closeout handoff `docs/handoffs/phase-65.md`

Protected paths remain untouched:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/**`
- `db/migrations/**`
- `go.mod`
- `go.sum`

## Non-Goals

- No database migrations.
- No telemetry ingest contract changes.
- No device credential schema or token hashing changes.
- No public feed URL changes.
- No GTFS-RT protobuf semantic changes.
- No validator execution semantic changes.
- No connector manifest schema changes.
- No prediction adapter behavior changes.
- No auth boundary or public/private route boundary weakening.
- No arbitrary backend command execution from the Operations Console.
- No simulator web action that sends telemetry, writes diagnostics, starts
  services, or captures command output.
- No automatic GTFS edits, draft mutations, or schedule publishes from the
  quality guidance UI.
- No retained evidence writes.
- No consumer status changes.
- No external contact.
- No real vendor payloads, real private telemetry, private device IDs,
  credentials, or webhook destinations.
- No CAL-ITP/Caltrans compliance, agency approval/adoption, final-root,
  consumer submission/review/acceptance/listing/display/ingestion, hosted
  SaaS, paid support, SLA/uptime, production-readiness, public-launch, vendor
  compatibility, hardware-certification, production AVL reliability, real
  realtime data, or production-grade ETA claim.

## Validation Plan

- `git diff --check`
- `go test ./cmd/agency-config`
- `go test ./cmd/telemetry-simulator` if simulator helpers change
- `go test ./internal/devices ./internal/telemetry ./internal/gtfs` when
  related boundaries are touched
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

Run DB-backed integration checks only if a checkpoint touches behavior that
requires the database. If Java, Docker, validator tooling, or pinned images are
unavailable, record the environment blocker and continue with
non-environment-dependent checks.

## Claim Boundary

Phase 65 may say the private Operations Console is easier for operators to
use, that simulator guidance supports synthetic local/reference diagnostics,
and that GTFS quality guidance helps operators decide what to fix next.

Phase 65 must not say or imply that Open Transit RT is CAL-ITP/Caltrans
compliant, accepted by a consumer, adopted or approved by an agency, final-root
approved, a hosted SaaS, SLA-backed, universally production ready,
vendor-compatible, hardware certified, public-launch complete, production AVL
reliable, real-realtime proven, or production-grade ETA proven.

Simulator outputs, device binding rows, telemetry freshness rows, and
validator notices are private diagnostics or supporting signals only. They are
not retained external evidence and must not move consumer tracker records
beyond `prepared`.

## Rollback Path

Phase 65 should remain private UI, docs, and synthetic-guidance work. If
rollback is needed, revert the specific checkpoint commit that added the UI,
generated guidance, test, or documentation change. Public feed URLs, DB schema,
telemetry ingest, GTFS-RT semantics, validator execution semantics, connector
manifest schema, prediction adapter behavior, auth boundaries, protected
evidence paths, and consumer tracker statuses should remain untouched.
