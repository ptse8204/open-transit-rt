# Phase 62 -- Guided Setup And Browser GTFS Import

## Status

Planned.

Phase 62 reduces command-line dependence by improving the private Operations
Console setup flow and adding an authenticated browser GTFS import path. It
must reuse the existing GTFS import/publish pipeline, preserve public feed
contracts, avoid schema changes, keep raw uploads and downloaded ZIP files out
of committed repository paths, and avoid compliance or evidence claims.

## Checkpoints

- `Phase 62 -- Checkpoint 000001: add guided setup and browser GTFS import plan`
- `Phase 62 -- Checkpoint 000002: implement guided setup wizard`
- `Phase 62 -- Checkpoint 000003: implement browser GTFS import and validation flow`
- `Phase 62 -- Checkpoint 000004: close guided setup and browser GTFS import`

## Planned Scope

### Guided Setup Wizard

- Add a private setup wizard view under the Operations Console.
- Keep the wizard read-oriented for all admin roles, with mutations limited to
  existing admin-only setup actions.
- Show setup stages for agency profile, publication metadata, GTFS, feed
  outputs, telemetry, validators, connectors, and readiness review.
- Tie each stage to existing private routes, public-safe docs, current status,
  next action, and a claim boundary.
- Add JSON output for the wizard model so operators can export private
  diagnostics without creating evidence.

### Browser GTFS Import

- Add an authenticated admin-only browser flow for GTFS import.
- Support safe URL import and multipart ZIP upload where runtime dependencies
  are available.
- Write raw GTFS ZIP bytes only to temporary runtime storage and delete them
  after the import attempt.
- Reuse `internal/gtfs.ImportService.ImportZip` so validation, rollback, import
  records, active feed activation, and validation-report behavior stay the same
  as the CLI path.
- Show import result status, feed version, counts, warning/error totals, and
  next actions without exposing raw ZIP contents or private paths.
- Keep the CLI import path available.

## Non-Goals

- No database migrations.
- No public feed URL changes.
- No change to GTFS-RT protobuf semantics.
- No change to telemetry ingest.
- No retained evidence.
- No `docs/evidence` writes.
- No consumer status changes.
- No external contact performed by this implementation work.
- No compliance, agency approval, consumer acceptance, hosted SaaS, vendor
  compatibility, production-readiness, or production-grade ETA claim.

## Files Expected To Change

- `cmd/agency-config/main.go`
- `cmd/agency-config/operations.go`
- new or existing `cmd/agency-config/operations_gtfs_import.go`
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
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json`
- `git diff --exit-code -- docs/evidence/captured`

Run broader checks when relevant and available:

- `make external-connection-check`
- `make adapter-conformance`
- `make test-connector-examples`
