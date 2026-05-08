# Phase 38 Handoff

## Phase

Phase 38 — Integration Adapter Kit

## Status

Complete for the navigation and conformance scope.

Phase 38 added reusable adapter-kit documentation and synthetic conformance
fixtures. It did not create final-root evidence, external evidence packets,
consumer artifacts, agency approval or adoption claims, consumer acceptance
claims, CAL-ITP/Caltrans compliance claims, production-readiness claims,
vendor-compatibility claims, production AVL reliability claims, or
production-grade ETA-quality claims.

## What Changed

- Added `docs/integration-adapter-kit.md` as the central map for adapter
  decisioning, Phase 37 onboarding-to-telemetry flow, `/v1/telemetry` pointers,
  synthetic AVL examples, mapping/diagnostics, external predictor lifecycle,
  validator boundaries, monitoring boundaries, consumer/feed workflow
  boundaries, redaction/evidence boundaries, and claim limits.
- Added neutral synthetic fixtures under `testdata/avl-vendor/`:
  - `minimal-gps.json`
  - `full-gps.json`
  - `multi-vehicle-mapping.json`
  - `multi-vehicle-gps.json`
  - `duplicate-batch.json`
  - `out-of-order-batch.json`
- Added `testdata/avl-vendor/README.md` with fixture purpose, expected
  diagnostic codes, valid/warning/failing classification, and evidence
  boundaries.
- Refreshed `cmd/avl-vendor-adapter` help text from Phase 29B-specific wording
  to dry-run adapter kit wording.
- Added focused Go tests for new fixtures, CLI help/output boundary wording,
  no-send-mode behavior, and stale/future/low-accuracy/duplicate/out-of-order
  diagnostics.
- Updated root README, docs navigation, tutorial navigation, reusable agency
  onboarding tutorial, device/AVL tutorial, roadmap/status/backlog/open-question
  docs, `docs/dependencies.md`, `make validate`, and this handoff.

## Go Behavior

Runtime behavior did not change. The only Go changes are help text for
`cmd/avl-vendor-adapter` and focused tests.

The CLI remains dry-run only. Network send mode still fails and no telemetry is
submitted by `cmd/avl-vendor-adapter`.

## Boundaries Reviewed

- `docs/dependencies.md` was updated to record Phase 38 adapter kit support as
  synthetic dry-run/developer integration support only.
- `docs/decisions.md` was reviewed; no architecture-significant decision
  changed, so no edit was needed.

## Not Added

- No network send mode.
- No named vendor support.
- No real vendor payloads.
- No credentials.
- No runtime external predictor integration.
- No Prometheus/Grafana assets.
- No OpenTelemetry SDK/exporter wiring.
- No consumer APIs.
- No consumer status changes.
- No final-root evidence.
- No external evidence.

## Phase 37 Status

Phase 37 onboarding remained closed. The Phase 37 no-network dry-run validation
check remains part of `make validate` and stayed green in the Phase 38 check
run.

## Consumer Status

All seven consumer and aggregator targets remain `prepared`. No target was
submitted, reviewed, accepted, rejected, listed, ingested, or moved to
`blocked`.

## Checks Run

- `gofmt` for changed Go files.
- `go test ./internal/avladapter ./cmd/avl-vendor-adapter ./internal/prediction ./internal/feed/tripupdates`
- `make realtime-quality`
- `make validate`
- `make test`
- `git diff --check`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- Read-only consumer status check confirming all seven targets remain
  `prepared`.

## Next Recommended Phase

Phase 39 — CAL-ITP-Style Readiness Workflow.

