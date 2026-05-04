# Phase 29B Handoff

## Phase

Phase 29B — AVL / Vendor Adapter Pilot Implementation

## Status

- Complete for the synthetic, dry-run-only AVL/vendor adapter pilot scope.
- Active phase after this handoff: Phase 30 — Consumer Submission Execution.

## What Was Implemented

- Added `internal/avladapter` for strict synthetic vendor payload and mapping transforms into the existing `telemetry.Event` contract.
- Added dry-run-only `cmd/avl-vendor-adapter` with required `--dry-run`, optional deterministic `--reference-time`, telemetry JSON array on stdout, and diagnostics JSON array on stderr.
- Added synthetic fixtures under `testdata/avl-vendor/` for valid, source mismatch, duplicate mapping, empty mapped IDs, missing/invalid coordinate, stale/future timestamps, unknown vendor vehicle, low GPS accuracy, mixed batch, duplicate/out-of-order dry-run observations, optional trip hint, and malformed payload.
- Added focused tests for mapping authority, telemetry contract validation, diagnostics shape, source mismatch, duplicate mappings, empty mapped identifiers, secret-like unknown mapping fields, fixed-time diagnostics, partial dry-run output, and no-success output as `[]`.
- Updated device/AVL docs, evidence template/discoverability, dependencies, decisions, current status, and latest handoff guidance.

## What Was Designed But Intentionally Not Implemented Yet

- No network send mode.
- No real vendor runtime adapter.
- No real vendor payload, private AVL payload, credential, endpoint URL, token, or private identifier handling.
- No evidence claim that any vendor or hardware is certified, compatible, production-ready, or reliable.

## Schema And Interface Changes

- No database schema changes.
- No public API changes.
- No telemetry request shape changes.
- No device token lifecycle changes.
- No GTFS-RT protobuf contract changes.
- No Trip Updates adapter behavior changes.

## Dependency Changes

- No external runtime dependency was added.
- `docs/dependencies.md` now records the Phase 29B synthetic adapter pilot as dry-run transform tooling only, not a named vendor integration.

## Migrations Added

- None.

## Tests Added And Results

- Added focused tests under `internal/avladapter` and `cmd/avl-vendor-adapter`.
- Focused `go test ./internal/avladapter ./cmd/avl-vendor-adapter` passed.

## Checks Run And Blocked Checks

- `go test ./internal/avladapter ./cmd/avl-vendor-adapter` — passed.
- `go run ./cmd/avl-vendor-adapter help` — passed.
- `go run ./cmd/avl-vendor-adapter --dry-run --reference-time 2026-05-04T12:00:00Z --mapping testdata/avl-vendor/mapping.json testdata/avl-vendor/valid.json` — passed.
- `go run ./cmd/avl-vendor-adapter --mapping testdata/avl-vendor/mapping.json testdata/avl-vendor/valid.json` — failed as expected with send-mode-not-implemented wording.
- Targeted `testdata/avl-vendor` secret-like fixture scan — passed with no matches.
- Broader Phase 29B docs/evidence scan — reviewed; matches were redaction rules and negative claim-boundary wording only.
- `make validate` — passed.
- `make test` — passed.
- `make realtime-quality` — passed.
- `make smoke` — passed.
- `make test-integration` — passed.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- `git diff --check` — passed.
- Blocked checks: none.

## Known Issues

- The adapter is synthetic and dry-run-only.
- Partial stdout from a nonzero dry run is transform output only, not submitted telemetry, production integration evidence, successful vendor compatibility proof, or database ingest status.
- Duplicate/out-of-order adapter diagnostics are batch-level dry-run observations only, not telemetry ingest outcomes.
- Phase 29B does not prove real vendor compatibility, certified hardware support, production AVL reliability, consumer acceptance, CAL-ITP/Caltrans compliance, agency endorsement, hosted SaaS availability, or marketplace/vendor equivalence.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/handoffs/phase-29b.md`
  - `docs/evidence/consumer-submissions/submission-workflow.md`
  - `docs/evidence/consumer-submissions/status.json`
  - `docs/phase-30-consumer-submission-execution.md`
- First files likely to edit:
  - Phase 30 consumer submission records and packets only when retained, redacted, target-originated evidence exists.
- Commands to run before coding:
  - `make validate`
  - `make test`
  - `git diff --check`
- Known blockers:
  - Consumer or aggregator statuses must not advance beyond `prepared` without retained, redacted, target-originated evidence.
  - Do not contact external portals, automate submissions, guess submission paths, or claim acceptance without authorized operator action and evidence.
- Recommended first implementation slice:
  - Start Phase 30 — Consumer Submission Execution only if maintainers have authorized target-specific submission work and have or can retain redacted target-originated evidence.
