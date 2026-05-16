# Vehicle AVL Connector Starter Kits

This page maps the local starter-kit paths for vehicle telemetry connectors.
All paths are disabled by default, synthetic/local unless explicitly configured
by an operator, and redaction-first. They do not prove vendor compatibility,
hardware certification, production AVL reliability, consumer display,
compliance, production readiness, hosted operation, SLA/uptime, or ETA quality.

## Starter Matrix

| Pattern | Starter kit | Use when | Default behavior |
| --- | --- | --- | --- |
| CSV replay | `examples/connectors/telemetry-csv-replay` | A technical helper has a small synthetic export and wants deterministic normalization. | Dry-run summary only; no network send. |
| GPS/API polling | `examples/connectors/telemetry-http-poller` | A deployment-owned sidecar can poll a source and normalize a batch. | Dry-run summary only; no private endpoint or token in the manifest. |
| Webhook sidecar | `examples/connectors/telemetry-webhook-sidecar` | A deployment-owned sidecar receives webhook-style batches and transforms them before calling `/v1/telemetry`. | Dry-run summary only; no listener or send path. |
| Synthetic-only | `scripts/telemetry-simulator.sh` and `testdata/telemetry-simulator` | Operators need local training, stale/unknown/low-quality review, or reproducible examples. | Synthetic dry-run or local simulator flow only. |
| Vendor payload transform | `cmd/avl-vendor-adapter` and `testdata/avl-vendor` | A technical helper needs a redaction-first example for mapping vendor-shaped payloads into the telemetry event model. | Dry-run by default; send mode requires explicit manifest and environment configuration. |

## Redaction Rules

- Keep connector manifests free of tokens, private endpoints, raw payloads,
  vendor secrets, webhook signatures, database URLs, and private paths.
- Store real credentials outside the repo and reference them by deployment
  environment only.
- Prefer dry-run output until the operator has reviewed device bindings,
  agency IDs, vehicle IDs, timestamp handling, coordinates, stale thresholds,
  and low-quality GPS behavior.
- Treat every example as a starter pattern, not a vendor certification or
  production reliability claim.

## Review Sequence

1. Validate the connector manifest with `make external-connection-check`.
2. Run the example tests with `make test-connector-examples`.
3. Dry-run a synthetic fixture and confirm `dry_run=true` and
   `network_send=false` in generated telemetry events.
4. Review `/admin/operations/connectors/workbench`,
   `/admin/operations/devices`, and `/admin/operations/realtime` before
   enabling any deployment-owned sender.
5. If a sender is later configured, keep diagnostics under `.cache/` or another
   private runtime path and do not treat them as public evidence.
