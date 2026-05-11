# Telemetry Simulator And Device Trial

This tutorial covers the synthetic telemetry simulator added in Phase 44.
It sends synthetic events through the real authenticated `POST /v1/telemetry`
ingest path. It is for local/reference diagnostics only.

It does not create evidence packets, change consumer statuses, prove vendor
compatibility, prove production AVL reliability, prove real realtime data,
prove production-grade ETA quality, or prove CAL-ITP/Caltrans compliance.

## What It Tests

The simulator can exercise:

- device bearer-token authentication;
- `/v1/telemetry` request validation;
- accepted, duplicate, and out-of-order ingest statuses;
- device binding rejection for an unknown device payload;
- optional DB-backed matching after accepted HTTP ingest;
- optional private Vehicle Positions debug snapshots from the same builder used
  by the feed service.

The simulator does not bypass ingest. The optional matcher step reads accepted
events back from Postgres after the HTTP request succeeds.

## Local Quick Run

Start the local app:

```bash
make agency-app-up
```

Run the default synthetic on-route scenario:

```bash
make telemetry-simulator
```

For a full local matcher and Vehicle Positions debug snapshot:

```bash
RUN_MATCHER=true make telemetry-simulator
```

Stop the local app when finished:

```bash
make agency-app-down
```

If a local verification run starts the app, always run `make agency-app-down`
afterward, even if a simulator step fails.

## Device Token Rules

The simulator uses:

```text
Authorization: Bearer <device-token>
```

For local demo runs, the seeded token is supplied through `DEVICE_TOKEN`:

```bash
DEVICE_TOKEN=dev-device-token
```

The local default is accepted only for loopback targets. Plain HTTP is allowed
only for loopback non-dry-run sends. For a non-loopback reference deployment,
use HTTPS and set a real device token issued for that deployment:

```bash
TARGET=https://reference.example.org \
DEVICE_TOKEN=replace-with-private-device-token \
SCENARIO=on-route \
make telemetry-simulator
```

Do not commit or share device tokens. The simulator does not write
Authorization headers or token values into diagnostics. Prefer environment
variables such as `DEVICE_TOKEN` for credentials; CLI flags are intended for
non-secret convenience values.

## Scenarios

The private Operations Console also has a browser guide:

```text
/admin/operations/telemetry-simulator
/admin/operations/telemetry-simulator.json
```

That page reads committed synthetic scenario metadata, shows target rules,
credential-handling guidance, private diagnostics policy, and copyable
operator-shell commands. It is guidance only: it does not execute commands,
send telemetry, collect device tokens, read `.cache` diagnostics, create
evidence, or change consumer statuses.

List scenarios:

```bash
scripts/telemetry-simulator.sh --list-scenarios
```

Current synthetic scenarios live under `testdata/telemetry-simulator/`:

| Scenario | Purpose |
| --- | --- |
| `on-route` | Local demo event near `trip-10-0800` on `valid-small`. |
| `stale` | Event older than the default matcher stale threshold. |
| `out-of-order` | Newer event first, then an older event for the same vehicle. |
| `unknown-device` | Payload device ID that does not match the bearer token binding. |
| `low-quality-gps` | Off-shape, low-quality GPS point for conservative diagnostics. |
| `after-midnight` | Requires a synthetic overnight feed and matching device binding. |
| `block-transition` | Requires a synthetic reference feed with block-linked trips. |

Some scenarios require a reference deployment with a matching synthetic feed
and device binding. Each fixture has a `requires` list.

## Private Diagnostics

Diagnostics default to:

```text
.cache/telemetry-simulator/<timestamp>/
```

The directory contains private local summaries such as:

- `summary.json`
- `events.json`
- `assignments.json` when `RUN_MATCHER=true`
- `vehicle_positions_debug.json` when `RUN_MATCHER=true`

The output is intentionally private diagnostics. It must not be committed as an
evidence packet without a future evidence-specific review, redaction, and claim
mapping process.

Custom output directories must resolve under repo `.cache/` unless
`ALLOW_UNIGNORED_OUTPUT_DIR=true` is set. The simulator always rejects
`docs/evidence`, rejects symlink output directories, creates new output
directories with mode `0700`, and runs a final redaction scan over generated
files.

## Dry Run

Preview payloads without sending telemetry:

```bash
DRY_RUN=true make telemetry-simulator
```

Dry run output is useful for checking scenario identity, timestamps, target URL
shape, and diagnostic shape. Dry runs may validate a non-loopback `http://`
target because no credentials are sent. It is not ingest proof.

## Reference Deployment Use

For a reference deployment:

```bash
TARGET=https://reference.example.org \
DEVICE_TOKEN=replace-with-private-device-token \
AGENCY_ID=agency-id \
DEVICE_ID=device-id \
VEHICLE_ID=vehicle-id \
SCENARIO=on-route \
make telemetry-simulator
```

To run matching diagnostics, also provide DB access from the operator machine:

```bash
TARGET=https://reference.example.org \
DEVICE_TOKEN=replace-with-private-device-token \
DATABASE_URL=postgres://... \
RUN_MATCHER=true \
SCENARIO=on-route \
make telemetry-simulator
```

The DB-backed matcher step is optional. It is useful for private diagnostics,
but the simulator's ingest test remains the authenticated HTTP request.
Prefer `DATABASE_URL` as an environment variable and avoid passing database
URLs through shell history or command-line flags.

## Boundary Checklist

- Use synthetic fixtures only.
- Use real device bearer-token auth.
- Post only to `/v1/telemetry`.
- Do not use real vendor payloads.
- Do not include private telemetry, credentials, private device IDs, or private
  vehicle IDs in committed fixtures.
- Do not create evidence packets from simulator output.
- Do not write simulator diagnostics outside `.cache/` unless
  `ALLOW_UNIGNORED_OUTPUT_DIR=true` is explicitly set for a private operator
  path.
- Do not change `docs/evidence/consumer-submissions/status.json`.
- Do not claim vendor compatibility, production AVL reliability, real realtime
  data, production-grade ETA quality, or CAL-ITP/Caltrans compliance.
