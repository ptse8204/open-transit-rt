# Synthetic AVL Adapter Fixtures

These fixtures exercise the dry-run adapter kit example in
`internal/avladapter` and `cmd/avl-vendor-adapter`.

All identifiers are synthetic. These fixtures are not real vendor payloads, not
vendor compatibility evidence, not production AVL reliability evidence, and not
telemetry ingest status proof.

Use only diagnostic codes defined by `internal/avladapter`.

## Reference Command

```bash
go run ./cmd/avl-vendor-adapter --dry-run \
  --reference-time 2026-05-04T12:00:00Z \
  --mapping testdata/avl-vendor/mapping.json \
  testdata/avl-vendor/minimal-gps.json
```

## Fixture Manifest

| Fixture | Expected result | Expected diagnostic codes | Class | Purpose |
| --- | --- | --- | --- | --- |
| `mapping.json` | Valid mapping for one synthetic vehicle. | none | valid | Baseline mapping authority for most payload fixtures. |
| `multi-vehicle-mapping.json` | Valid mapping for two synthetic vehicles. | none | valid | Demonstrates multiple mapped vehicles without real identifiers. |
| `minimal-gps.json` | One transformed telemetry event. | none | valid | Minimal required GPS fields. |
| `full-gps.json` | One transformed telemetry event. | none | valid | Optional `bearing`, `speed_mps`, `accuracy_m`, and `trip_hint`. |
| `multi-vehicle-gps.json` | Two transformed telemetry events. | none | valid | Multiple vehicles using `multi-vehicle-mapping.json`. |
| `stale-timestamp.json` | One transformed telemetry event plus warning. | `stale_timestamp` | warning-bearing | Stale observation review relative to the reference time. |
| `future-timestamp.json` | One transformed telemetry event plus warning. | `future_timestamp` | warning-bearing | Future observation review relative to the reference time. |
| `low-gps-accuracy.json` | One transformed telemetry event plus warning. | `low_gps_accuracy` | warning-bearing | Low GPS quality review. |
| `duplicate-batch.json` | Two transformed telemetry events plus warning. | `duplicate_observation` | warning-bearing | Duplicate observed time within one dry-run batch. |
| `out-of-order-batch.json` | Two transformed telemetry events plus warning. | `out_of_order_observation` | warning-bearing | Older observation after a newer observation in one dry-run batch. |
| `duplicate-out-of-order.json` | Three transformed telemetry events plus warnings. | `duplicate_observation`, `out_of_order_observation` | warning-bearing | Combined legacy duplicate/out-of-order dry-run case. |
| `optional-trip-hint.json` | One transformed telemetry event with `trip_hint`. | none | valid | Shows that `trip_hint` is preserved only as a telemetry hint. |
| `batch-mixed.json` | One transformed event, diagnostics, and nonzero command exit. | `unknown_vendor_mapping` | intentionally failing | Mixed valid and invalid batch behavior. |
| `unknown-vendor-vehicle.json` | No transformed events and nonzero command exit. | `unknown_vendor_mapping` | intentionally failing | Unknown mapping rejection. |
| `source-mismatch.json` | No transformed events and nonzero command exit. | `vendor_source_mismatch` | intentionally failing | Vendor source mismatch for known synthetic IDs. |
| `missing-coordinate.json` | No transformed events and nonzero command exit. | `missing_required_field` | intentionally failing | Missing required coordinate field. |
| `invalid-coordinate.json` | No transformed events and nonzero command exit. | `invalid_coordinate` | intentionally failing | Coordinate range validation. |
| `malformed.json` | No transformed events and nonzero command exit. | `invalid_payload_json` | intentionally failing | Invalid JSON payload. |
| `duplicate-mapping.json` | Invalid mapping. | `duplicate_mapping` | intentionally failing | Duplicate mapping row rejection. |
| `empty-mapped-ids.json` | Invalid mapping. | `empty_mapped_identifier` | intentionally failing | Empty emitted Open Transit RT IDs. |

## Boundary Notes

- Mapping files are the authority for emitted `agency_id`, `device_id`, and
  `vehicle_id`.
- Vendor-looking identifiers in payloads are lookup keys only.
- Warnings are dry-run review output only.
- Duplicate and out-of-order fixture warnings are not database ingest statuses.
- Do not add real vendor names, credentials, endpoint URLs, tokens, private
  device identifiers, private vehicle identifiers, or raw private telemetry to
  this directory.

