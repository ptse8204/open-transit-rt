# Phase 29B — AVL / Vendor Adapter Pilot Implementation

## Status

Planned Track B integration phase. Not implemented until selected in `docs/handoffs/latest.md`.

## Purpose

Turn the Phase 25 Device and AVL Integration Kit from documentation into a small, safe, testable adapter pilot pattern.

Phase 25 documented how devices and AVL vendors can send telemetry into Open Transit RT. Phase 29B should provide a minimal adapter pilot that transforms synthetic vendor-style payloads into Open Transit RT telemetry events without introducing vendor lock-in, real credentials, or certified hardware claims.

This phase is about **adapter pattern implementation and synthetic tests**, not real vendor certification or production AVL reliability.

## Why This Phase Exists

Open Transit RT should be able to connect to devices, simple GPS emitters, and vendor AVL feeds. But those integrations must remain outside core matching and feed generation. Vendor-specific assumptions should be isolated in adapters or sidecars that transform external payloads into the existing `POST /v1/telemetry` event shape.

## Scope

1. Generic AVL/vendor adapter contract.
2. Synthetic vendor payload fixtures.
3. Adapter dry-run mode.
4. Adapter-to-telemetry transformation tests.
5. Credential and redaction boundaries.
6. Optional local no-hardware integration flow.
7. Documentation for deployment-owned adapter patterns.
8. Evidence scaffold for adapter pilot reviews.

## Required Work

### 1) Adapter Contract

Define a generic adapter contract that maps external AVL/device payloads into Open Transit RT telemetry events.

The contract should cover:

- source agency or deployment;
- vendor/device identifier;
- Open Transit RT `agency_id`;
- Open Transit RT `device_id`;
- Open Transit RT `vehicle_id`;
- observation timestamp;
- latitude and longitude;
- optional bearing;
- optional speed;
- optional accuracy;
- optional trip hint;
- source payload retention policy;
- validation errors;
- redaction expectations.

The adapter contract should not require any named vendor.

### 2) Synthetic Vendor Payload Fixtures

Add synthetic fixtures such as:

- valid vendor payload;
- missing coordinate;
- stale timestamp;
- future timestamp;
- unknown vendor vehicle;
- low GPS accuracy;
- batch payload;
- duplicate/out-of-order event;
- optional trip hint;
- malformed payload.

Use only synthetic identifiers:

- `vendor-demo`
- `vendor-device-1`
- `vendor-vehicle-1`
- `demo-agency`
- `device-1`
- `bus-1`

Do not commit real vendor payloads or private device identifiers.

### 3) Adapter Implementation Shape

Implement the smallest safe adapter pilot if appropriate.

Acceptable shapes:

- a small CLI helper under `scripts/` or `cmd/` that reads synthetic JSON and prints/transforms Open Transit RT telemetry JSON;
- a dry-run-only adapter utility;
- a test-only adapter package;
- a documented sidecar pattern without runtime code if implementation would be too broad.

Preferred behavior:

- dry-run prints the transformed telemetry payload without sending it;
- send mode, if implemented, must require explicit base URL and token from environment;
- no secret is printed;
- no default production endpoint is assumed;
- adapter errors are clear;
- batch handling is deterministic.

If send mode is too much, defer it and keep the adapter as transform/test-only.

### 4) Telemetry Submission Boundary

If the adapter sends telemetry, it must call the existing `POST /v1/telemetry` endpoint and use the existing bearer device-token mechanism.

Do not change:

- telemetry API shape;
- device API shape;
- token lifecycle;
- GTFS-RT contracts;
- prediction logic.

### 5) Tests

Add tests that prove:

- valid synthetic vendor payload transforms to valid Open Transit RT telemetry;
- missing required fields fail;
- bad coordinates fail;
- stale/future timestamp behavior is labeled correctly;
- vendor vehicle/device mapping works;
- unknown vendor mapping fails;
- dry-run does not send telemetry;
- no secrets are printed;
- batch order is stable;
- adapter output can be accepted by the existing telemetry payload validation shape, where feasible.

### 6) Documentation

Add or update:

- `docs/phase-29b-avl-vendor-adapter-pilot.md`
- `docs/tutorials/device-avl-integration.md`
- `docs/tutorials/device-token-lifecycle.md` only if token handling guidance changes
- `docs/evidence/device-avl/README.md`
- `docs/evidence/device-avl/templates/integration-review-template.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-29b.md`
- `docs/dependencies.md` only if a named external dependency or integration status is introduced
- `docs/decisions.md` only if adapter architecture policy changes

### 7) Evidence And Redaction

If adding adapter evidence templates, include fields for:

- adapter type;
- source system;
- source permission;
- payload redaction review;
- identifier mapping review;
- token storage review;
- dry-run output;
- telemetry accepted evidence;
- Operations Console freshness review;
- Vehicle Positions review;
- what the evidence proves;
- what it does not prove.

Do not add fake evidence, fake vendor approvals, fake hardware certifications, or private AVL payloads.

## Acceptance Criteria

Phase 29B is complete only when:

- a generic adapter contract exists;
- synthetic vendor payload fixtures exist;
- transform/dry-run behavior exists or is explicitly deferred with reason;
- tests prove the adapter boundary for synthetic payloads if code is added;
- vendor credentials remain outside the repo;
- no real vendor payloads are committed;
- no certified vendor or hardware support is claimed;
- no public feed URL, GTFS-RT contract, consumer status, or evidence claim changes are introduced;
- `docs/handoffs/phase-29b.md` exists and uses the repo handoff template.

## Required Checks

```bash
make validate
make test
make realtime-quality
make smoke
make test-integration
docker compose -f deploy/docker-compose.yml config
git diff --check
```

If scripts are added or changed:

```bash
sh -n <changed scripts>
<new script> help
<new script> --dry-run <synthetic fixture>
```

If local app/demo behavior changes:

```bash
make demo-agency-flow
make agency-app-up
make agency-app-down
docker compose -f deploy/docker-compose.yml --profile app config
```

## Explicit Non-Goals

Phase 29B does not:

- certify any hardware vendor;
- add a proprietary vendor integration to core behavior;
- commit vendor credentials;
- commit private AVL payloads;
- claim production AVL reliability;
- change telemetry API contract;
- change device token lifecycle;
- change public feed URLs;
- change GTFS-RT protobuf contracts;
- change Trip Updates prediction logic;
- claim consumer acceptance or CAL-ITP/Caltrans compliance;
- add hosted SaaS or paid support claims.

## Security And Privacy Boundaries

Do not commit:

- vendor credentials;
- device tokens;
- admin tokens;
- JWT/CSRF secrets;
- DB URLs with passwords;
- private AVL payloads;
- raw private telemetry;
- private vehicle/device/vendor identifiers;
- private logs;
- `.cache` files.

Synthetic fixtures must be clearly synthetic.

## Likely Files

- `docs/phase-29b-avl-vendor-adapter-pilot.md`
- `docs/handoffs/phase-29b.md`
- `docs/tutorials/device-avl-integration.md`
- `docs/evidence/device-avl/README.md`
- `docs/evidence/device-avl/templates/integration-review-template.md`
- `testdata/avl-vendor/`
- `scripts/` only if adding a safe adapter helper
- `internal/` or `cmd/` only if adding a small testable adapter package or CLI
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/dependencies.md` only if dependency status changes
- `docs/decisions.md` only if adapter policy changes
