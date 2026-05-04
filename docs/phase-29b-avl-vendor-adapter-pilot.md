# Phase 29B — AVL / Vendor Adapter Pilot Implementation

## Status

Complete for a synthetic, dry-run-only vendor/AVL adapter pilot pattern.

## Purpose

Phase 29B turns the Phase 25 Device and AVL Integration Kit from documentation into a small, safe, testable adapter pattern.

The implementation transforms synthetic vendor-style payload fixtures into the existing Open Transit RT telemetry event shape. It does not send telemetry, certify vendors, add real vendor payloads, add credentials, change runtime APIs, or prove production AVL reliability.

## Implemented Adapter Shape

The pilot has two code surfaces:

- `internal/avladapter`: strict JSON mapping and payload transform logic.
- `cmd/avl-vendor-adapter`: a dry-run-only CLI that prints transformed telemetry JSON.

The CLI requires `--dry-run` in Phase 29B:

```bash
go run ./cmd/avl-vendor-adapter --dry-run \
  --reference-time 2026-05-04T12:00:00Z \
  --mapping testdata/avl-vendor/mapping.json \
  testdata/avl-vendor/valid.json
```

Running without `--dry-run` fails with a clear message that send mode is not implemented. Phase 29B intentionally does not add network send mode.

## Mapping Contract

`testdata/avl-vendor/mapping.json` is the authority for Open Transit RT identifiers. Vendor payload IDs are lookup keys only. A vendor payload cannot override mapped `agency_id`, `device_id`, or `vehicle_id`; extra identifier fields in payload records are rejected by strict JSON decoding.

Mapping file shape:

```json
{
  "mappings": [
    {
      "vendor_source": "vendor-demo",
      "vendor_device_id": "vendor-device-1",
      "vendor_vehicle_id": "vendor-vehicle-1",
      "agency_id": "demo-agency",
      "device_id": "device-1",
      "vehicle_id": "bus-1",
      "notes": "Synthetic Phase 29B mapping for dry-run adapter tests."
    }
  ]
}
```

Allowed mapping row fields:

| Field | Required | Meaning |
| --- | --- | --- |
| `vendor_source` | Yes | Synthetic vendor/source namespace. |
| `vendor_device_id` | Yes | Vendor-side device lookup key. |
| `vendor_vehicle_id` | Yes | Vendor-side vehicle lookup key. |
| `agency_id` | Yes | Open Transit RT agency ID to emit. |
| `device_id` | Yes | Open Transit RT device ID to emit. |
| `vehicle_id` | Yes | Open Transit RT vehicle ID to emit. |
| `notes` | No | Public-safe notes. |

Mapping files must not contain tokens, endpoint URLs, auth headers, passwords, secrets, credentials, database URLs, private keys, real vendor account IDs, private device identifiers, or private vehicle identifiers. Unknown mapping fields are rejected.

Hard-error diagnostics are emitted for:

- duplicate rows with the same `vendor_source`, `vendor_device_id`, and `vendor_vehicle_id`;
- empty mapped `agency_id`, `device_id`, or `vehicle_id`;
- unknown mapping fields or malformed mapping JSON.

## Payload Contract

Synthetic vendor payloads use this envelope:

```json
{
  "vendor_source": "vendor-demo",
  "records": [
    {
      "vendor_device_id": "vendor-device-1",
      "vendor_vehicle_id": "vendor-vehicle-1",
      "observed_at": "2026-05-04T12:00:00Z",
      "lat": 49.2827,
      "lon": -123.1207,
      "bearing": 120.0,
      "speed_mps": 8.4,
      "accuracy_m": 7.5,
      "trip_hint": "trip-10-0800"
    }
  ]
}
```

Every transformed output record is marshaled and unmarshaled as the existing `telemetry.Event` type and must pass `telemetry.Event.Valid()`. No new telemetry request shape was introduced.

`trip_hint` is only a hint. It is not assignment proof, ETA proof, consumer-facing correctness proof, or evidence that a vehicle is matched.

## Output And Diagnostics

Output streams are stable and tested:

- stdout: transformed Open Transit RT telemetry JSON array.
- stderr: diagnostics as one JSON array.

If no records transform successfully, stdout is `[]`.

Diagnostic fields:

| Field | Meaning |
| --- | --- |
| `code` | Stable diagnostic code. |
| `severity` | `error` for hard errors or `warning` for review warnings. |
| `message` | Human-readable explanation. |
| `index` | Zero-based source record index when the diagnostic applies to a payload or mapping row. |

Warnings do not prevent output. Hard errors make the CLI exit nonzero. In a mixed batch, valid records are still printed to stdout in input order, invalid records are omitted, and hard-error diagnostics are printed to stderr.

Partial stdout from a nonzero exit is dry-run transform output only. It is not submitted telemetry, production integration evidence, successful vendor compatibility proof, or database ingest status.

Duplicate and out-of-order diagnostics are dry-run batch-level observations only. They are not telemetry ingest `accepted`, `duplicate`, or `out_of_order` database outcomes.

Stale and future timestamp diagnostics use `--reference-time <RFC3339>` for deterministic dry-run review. If omitted, the CLI uses current UTC time.

## Synthetic Fixtures

Synthetic fixtures live in `testdata/avl-vendor/`:

- valid vendor payload;
- source mismatch;
- duplicate mapping;
- empty mapped identifiers;
- missing coordinate;
- invalid coordinate;
- stale timestamp;
- future timestamp;
- unknown vendor vehicle;
- low GPS accuracy;
- mixed-validity batch;
- duplicate/out-of-order dry-run observations;
- optional trip hint;
- malformed payload.

All fixtures use synthetic identifiers such as `vendor-demo`, `vendor-device-1`, `vendor-vehicle-1`, `demo-agency`, `device-1`, and `bus-1`.

## Non-Goals

Phase 29B does not:

- send telemetry over the network;
- change `/v1/telemetry`;
- change device token lifecycle;
- change public feed URLs;
- change GTFS-RT protobuf contracts;
- change matching, Vehicle Positions, or Trip Updates behavior;
- add a named vendor dependency;
- commit vendor credentials or real AVL payloads;
- certify vendor or hardware support;
- prove production AVL reliability;
- change consumer statuses;
- prove consumer acceptance;
- prove CAL-ITP/Caltrans compliance.

## Validation

Focused checks:

```bash
go test ./internal/avladapter ./cmd/avl-vendor-adapter
go run ./cmd/avl-vendor-adapter help
go run ./cmd/avl-vendor-adapter --dry-run --reference-time 2026-05-04T12:00:00Z --mapping testdata/avl-vendor/mapping.json testdata/avl-vendor/valid.json
```

Full Phase 29B closure checks:

```bash
make validate
make test
make realtime-quality
make smoke
make test-integration
docker compose -f deploy/docker-compose.yml config
git diff --check
```
