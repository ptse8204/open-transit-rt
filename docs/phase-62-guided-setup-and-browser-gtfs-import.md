# Phase 62 -- Guided Setup And Browser GTFS Import

## Status

Complete.

Phase 62 reduces command-line dependence by improving the private Operations
Console setup flow and adding an authenticated browser GTFS import path. It
must reuse the existing GTFS import/publish pipeline, preserve public feed
contracts, avoid schema changes, keep raw uploads and downloaded ZIP files out
of committed repository paths, and avoid compliance or evidence claims.

Checkpoint 000002 added the private setup wizard. Checkpoint 000003 adds the
admin-only browser GTFS import route while keeping the CLI import path
available.

## Implementation Notes

Checkpoint 000003 adds `/admin/operations/gtfs-import` for admin-only GTFS ZIP
upload or safe URL import. The route derives agency and actor from the
authenticated principal, writes ZIP bytes only to temporary runtime storage,
deletes the temporary file after the import attempt, and calls
`internal/gtfs.ImportService.ImportZip` so validation, import records, active
feed activation, failure behavior, and validation-report storage match the CLI
path.

The page shows bounded status, import ID, feed version, validation counts,
stored-report status, row counts, and next actions. It does not render raw ZIP
bytes, temporary paths, raw reports, validator argv, credentials, or private
URLs. URL import blocks embedded credentials, non-HTTP(S) schemes, private or
internal hosts unless explicitly allowed for local testing, secret-looking
queries, oversized downloads, and unsafe redirects.

## Checkpoints

- `Phase 62 -- Checkpoint 000001: add guided setup and browser GTFS import plan`
- `Phase 62 -- Checkpoint 000002: implement guided setup wizard`
- `Phase 62 -- Checkpoint 000003: implement browser GTFS import and validation flow`
- `Phase 62 -- Checkpoint 000004: close guided setup and browser GTFS import`

## Closeout Summary

Phase 62 is complete for the guided setup and browser GTFS import scope.

- Added `/admin/operations/setup-wizard` and
  `/admin/operations/setup-wizard.json` as private setup guidance for agency
  profile, publication metadata, GTFS, feeds, telemetry, validators,
  connectors, and readiness.
- Added `/admin/operations/gtfs-import` as an admin-only browser import route
  for GTFS ZIP upload or safe URL import.
- Reused `internal/gtfs.ImportService.ImportZip` for browser imports so
  validation, import records, active feed activation, failure behavior, and
  validation-report storage match the CLI path.
- Kept the CLI import path and GTFS Studio draft/publish path available.
- Updated Operations Console dashboard, launchpad, setup, status, and handoff
  wording so browser import is visible without implying readiness or evidence.

Phase 62 added no migrations, public feed URL changes, GTFS-RT semantic
changes, telemetry ingest contract changes, validator execution semantic
changes, connector manifest schema changes, prediction adapter behavior
changes, consumer tracker status changes, public routes, retained evidence, or
protected evidence writes.

## Claim Boundary

The setup wizard and browser import page are private authenticated operator
workflow surfaces. Viewing them creates no evidence, contacts no external
party, opens no public route, changes no consumer status, and records no agency
approval, compliance, consumer acceptance, hosted-service, vendor compatibility,
SLA, public-launch, production-readiness, production AVL reliability, or
production-grade ETA outcome.

Browser import stores raw ZIP bytes only in temporary runtime storage for the
import attempt and removes them afterward. The rendered page shows bounded
import and validation status only; it does not render raw ZIP contents,
temporary paths, raw reports, validator argv, credentials, or private URLs.

All seven consumer and aggregator targets remain `prepared`.

## Verification

Verification was run from `/Users/edwintse/Downloads/open-transit-rt`; all
listed commands passed.

- `git diff --check`
- `go test ./cmd/agency-config`
- `go test ./cmd/gtfs-import ./internal/gtfs`
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
- `docker compose -f deploy/docker-compose.yml config`
- `make db-up`
- `make test-integration`

The browser import HTTP handler is covered with an injected import runner for
route, auth, CSRF, temp-file, URL-safety, and rendering behavior. The
DB-backed import, validation-report, active-feed activation, rollback, and
audit behavior remains covered by the existing GTFS importer integration tests
that ran through `make test-integration`.

## Next Phase

Phase 63 -- Feed Health and Readiness UX.

## Original Planned Scope

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
- `cmd/agency-config/operations_gtfs_import.go`
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
