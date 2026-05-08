# Phase 38 — Integration Adapter Kit

## Status

Complete for the integration adapter kit scope.

Phase 38 added a central integration map, expanded synthetic AVL fixture
coverage, documented fixture diagnostics, refreshed dry-run adapter CLI
boundary wording, and added focused conformance tests.

It did not add network send mode, named vendor support, real vendor payloads,
credentials, runtime external predictor integration, Prometheus/Grafana
deployment assets, OpenTelemetry SDK/exporter wiring, consumer submission APIs,
consumer status changes, final-root evidence, or external evidence packets.

## Goal

Make Open Transit RT easier to integrate with existing AVL/device systems,
external predictors, validators, monitoring, and consumer-facing feed workflows
through clear adapter boundaries and reusable examples.

## What Changed

- Added `docs/integration-adapter-kit.md` as the central map for telemetry,
  synthetic AVL, prediction, validator, monitoring, consumer workflow, and
  evidence/redaction boundaries.
- Linked the adapter kit from the root README, docs navigation, tutorial index,
  reusable agency onboarding guide, and device/AVL integration tutorial.
- Added `testdata/avl-vendor/README.md` with a fixture manifest, expected
  diagnostics, valid/warning/failing classification, synthetic-ID boundary, and
  no-vendor-evidence warning.
- Added neutral synthetic AVL fixtures for minimal GPS, full GPS, multi-vehicle
  mapping, multi-vehicle GPS, duplicate batch, and out-of-order batch cases.
- Refreshed `cmd/avl-vendor-adapter` help text so it describes a dry-run
  adapter kit example rather than a Phase 29B pilot.
- Added focused Go tests for new fixtures, dry-run diagnostics, no-send-mode
  behavior, and CLI boundary wording.
- Updated `docs/dependencies.md` to record Phase 38 as synthetic
  dry-run/developer integration support only.

## Boundaries Preserved

- `/v1/telemetry` remains the telemetry ingest contract.
- `internal/prediction.Adapter` remains the Trip Updates prediction boundary.
- TheTransitClock and other predictors remain deferred optional integrations.
- Prometheus/Grafana and OpenTelemetry remain deferred optional integrations.
- Consumer APIs and real vendor AVL integrations remain deferred.
- Consumer and aggregator statuses remain `prepared`.

## Checks

Required Phase 38 checks:

```bash
go test ./internal/avladapter ./cmd/avl-vendor-adapter ./internal/prediction ./internal/feed/tripupdates
make realtime-quality
make validate
make test
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
```

Also run a read-only consumer status check confirming all seven targets remain
`prepared`.

## Next

If Phase 38 closes successfully, the next recommended phase is Phase 39 —
CAL-ITP-Style Readiness Workflow.
