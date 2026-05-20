# API, Feed, And Extension Contracts

This page records the current Open Transit RT integration contracts for local
and self-hosted evaluation. These are release-candidate contracts, not an
indefinite compatibility promise and not evidence of external adoption,
consumer ingestion, compliance, or production readiness.

## Contract Stability Rules

- Keep public feed paths stable unless a future versioned path is introduced.
- Keep `POST /v1/telemetry` backward-compatible for accepted fields whenever
  possible; reject unsafe or ambiguous inputs instead of accepting partial data.
- Keep private admin JSON companion routes authenticated, no-store, and
  same-agency scoped.
- Keep connector manifests and adapter conformance fixtures versioned.
- Keep Trip Updates prediction behind the Go `internal/prediction.Adapter`
  boundary and the external HTTP predictor DTOs versioned by schema name.
- Document any breaking change before release-candidate packaging and add a
  transition note or new versioned path/schema where practical.

## Public Feed Paths

| Path | Contract | Notes |
| --- | --- | --- |
| `/public/feeds.json` | Public JSON feed discovery metadata. | Lists schedule, Vehicle Positions, Trip Updates, Alerts, license/contact metadata, and readiness booleans. Metadata does not prove source-of-truth listing or consumer acceptance. |
| `/public/gtfs/schedule.zip` | Public static GTFS ZIP. | Generated from the active published feed version, not draft data. |
| `/public/gtfsrt/vehicle_positions.pb` | Public GTFS-RT Vehicle Positions protobuf. | Vehicle Positions remain independent of Trip Updates predictor health. |
| `/public/gtfsrt/trip_updates.pb` | Public GTFS-RT Trip Updates protobuf. | Output remains behind the prediction adapter boundary and may be empty or withheld. |
| `/public/gtfsrt/alerts.pb` | Public GTFS-RT Alerts protobuf. | Alerts are separate from telemetry ingest and prediction. |
| `/public/agencies/{agency_id}/...` | Agency-scoped public feed variants where implemented. | Agency IDs must remain path-safe and must not bypass tenant checks. |

Public JSON debug routes for GTFS-RT remain authenticated diagnostics, not
public integration contracts.

## Telemetry Ingest Contract

`POST /v1/telemetry` is the authenticated device/adapter ingest boundary.

Request body:

- JSON object;
- `agency_id`, `device_id`, `vehicle_id`, `timestamp`, `lat`, and `lon`;
- optional `bearing`, `speed_mps`, and `accuracy_m`.

Runtime rules:

- bearer device token required;
- body size capped by server code;
- agency, device, and vehicle binding must match active credentials;
- timestamp must include timezone and must not be future-dated beyond the
  configured threshold;
- invalid coordinates, negative speed/accuracy, invalid bearing, duplicate,
  out-of-order, stale, unknown-device, wrong-agency, and low-quality GPS cases
  remain rejected or classified without exposing raw tokens or payloads.

Response body is JSON and includes accepted/classified status and `received_at`.
It must not include raw credentials, raw headers, token hashes, or private
payloads.

## Private Admin JSON Companion Routes

These routes are private, authenticated, same-agency scoped, and no-store. They
exist for browser workflows, smoke tests, and local automation, not public data
sharing:

- `/admin/operations.json`
- `/admin/operations/launchpad.json`
- `/admin/operations/checklist.json`
- `/admin/operations/feed-health.json`
- `/admin/operations/validation-center.json`
- `/admin/operations/readiness.json`
- `/admin/operations/telemetry-simulator.json`
- `/admin/operations/prediction-lab.json`
- `/admin/operations/connectors/tests.json`
- `/admin/operations/connectors/workbench.json`
- `/admin/operations/help.json`
- `/admin/operations/gtfs-workbench.json`
- `/admin/operations/validation-health.json`
- `/admin/operations/validation-health/refresh.json`
- `/admin/operations/reliability.json`
- `/admin/operations/maintenance.json`

Any new private HTML workflow that has a JSON companion should preserve the
same auth, no-store, agency-scope, redaction, and claim-boundary behavior.

## Connector And Adapter Contracts

Connector manifests use:

```text
open-transit-rt.connector.v1
```

Adapter conformance suites use:

```text
open-transit-rt.adapter_conformance.v1
```

Current conformance types are telemetry, prediction, validator, monitoring,
and consumer discovery. Fixtures stay synthetic-only and local/offline by
default. They must not include real credentials, private endpoint URLs, private
vendor payloads, portal automation, or status mutation.

## Prediction Adapter Contract

Trip Updates prediction remains behind:

```go
type Adapter interface {
    Name() string
    PredictTripUpdates(ctx context.Context, request Request) (Result, error)
}
```

The request includes the active GTFS feed version, current telemetry, current
assignments, and the Vehicle Positions feed URL or URL reference. The result
returns Trip Updates plus diagnostics. External HTTP predictor DTOs are
sanitized and versioned through the connector/adapter conformance fixtures.

Prediction diagnostics do not prove production-grade ETA quality, real-world
ETA accuracy, consumer display, consumer ingestion, or vendor compatibility.

## Local Contract Gate

Run:

```bash
make api-contract-check
```

The gate verifies this page against route registrations, connector manifests,
adapter conformance suite metadata, and prediction adapter symbols. It is a
local consistency check only; it does not contact agencies, vendors, consumers,
portals, or live validators.
