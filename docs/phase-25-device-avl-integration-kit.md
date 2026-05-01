# Phase 25 — Device And AVL Integration Kit

## Status

Planned Track B phase. Not implemented until selected in `docs/handoffs/latest.md`.

## Purpose

Make real vehicle telemetry onboarding practical for agency devices, vendors, or simple GPS emitters.

Open Transit RT can ingest authenticated telemetry and publish realtime feeds, but agencies and vendors need clear integration guidance, token lifecycle docs, payload examples, and troubleshooting.

## Scope

1. Device API integration guide.
2. Vendor AVL adapter boundary.
3. Device token lifecycle playbook.
4. Payload examples and simulator guidance.
5. Clock, timezone, GPS accuracy guidance.
6. Troubleshooting and redaction rules.

## Required Work

### 1) Telemetry API Guide

Document:

- endpoint path;
- authentication header;
- required fields;
- optional fields;
- timestamp requirements;
- lat/lon/bearing/speed expectations;
- examples with curl and JSON;
- common error responses.

### 2) Device Token Lifecycle

Document:

- create/rotate/rebind behavior;
- one-time token display;
- storage expectations;
- revocation/rotation after compromise;
- device-to-vehicle binding rules;
- operator responsibilities.

### 3) Vendor AVL Boundary

Describe how a vendor integration should work:

- transform vendor payloads into Open Transit RT telemetry events;
- keep vendor credentials outside repo;
- avoid writing vendor-specific coupling into core matching;
- use adapter scripts or integration services where appropriate.

### 4) Simulator And Testing

Improve docs/scripts if needed for:

- local simulator events;
- dry-run payload display;
- sample route movement;
- no-hardware demo path;
- warnings that simulator is not production AVL proof.

### 5) Troubleshooting

Document common failures:

- bad token;
- wrong agency/device/vehicle;
- timestamp too old/future;
- GPS accuracy problems;
- no assignment;
- stale feed output;
- validator pass but no consumer acceptance.

## Acceptance Criteria

Phase 25 is complete only when:

- an agency or vendor can understand how to send telemetry;
- token lifecycle is clear and safe;
- examples are reproducible and redacted;
- simulator/demo is clearly labeled;
- no production AVL quality claim is introduced.

## Required Checks

```bash
make validate
make test
git diff --check
```

If telemetry scripts or local app helpers change:

```bash
make smoke
make demo-agency-flow
make agency-app-up
make agency-app-down
```

## Explicit Non-Goals

Phase 25 does not:

- certify hardware vendors;
- add a proprietary AVL integration as core behavior;
- claim real-world device reliability without pilot evidence;
- commit vendor credentials;
- add rider apps, fares, or CAD/dispatch.

## Likely Files

- `docs/tutorials/device-avl-integration.md`
- `docs/tutorials/device-token-lifecycle.md`
- `scripts/device-onboarding.sh` only if safe improvements are needed
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-25.md`
