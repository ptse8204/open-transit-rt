# Phase 27 — Multi-Agency Isolation Prototype

## Status

Planned Track B phase. Not implemented until selected in `docs/handoffs/latest.md`.

## Purpose

Move from documented multi-agency strategy to testable isolation foundations.

The project currently supports agency-scoped records and auth concepts, but true multi-agency hosting requires stronger proof that data, admin actions, feeds, and evidence cannot leak across agencies.

## Scope

1. Multi-agency fixtures.
2. Agency-scoped auth tests.
3. Data isolation tests.
4. Per-agency feed discovery tests.
5. Public feed root strategy for multiple agencies.
6. Backup/restore implications.
7. Consumer packet implications.
8. Documentation of remaining true multi-tenant gaps.

## Required Work

### 1) Test Fixtures

Create two or more agency fixtures with distinct:

- agency IDs;
- GTFS feed versions;
- vehicle IDs;
- device IDs;
- telemetry events;
- public metadata.

### 2) Auth Boundary Tests

Test that users for one agency cannot access or mutate another agency’s:

- publication metadata;
- device bindings;
- validation records;
- scorecards;
- consumer records;
- operations console data.

### 3) Feed Isolation

Test or document:

- per-agency `feeds.json` behavior;
- per-agency feed roots or query strategy;
- public URL implications;
- consumer packet implications.

### 4) Operations Isolation

Document:

- backup/restore per agency versus whole database;
- evidence packets per agency;
- monitoring per agency;
- incident response implications.

## Acceptance Criteria

Phase 27 is complete only when:

- multi-agency assumptions are backed by tests or explicit documented gaps;
- agency-scoped auth boundary tests exist for key admin workflows;
- feed discovery isolation is tested or blocked with clear reason;
- docs truthfully state whether multi-agency support is prototype or production-ready.

## Required Checks

```bash
make validate
make test
make test-integration
git diff --check
```

If console/local app behavior changes:

```bash
make smoke
make agency-app-up
make agency-app-down
```

## Explicit Non-Goals

Phase 27 does not:

- claim production multi-tenant hosting;
- implement hosted SaaS;
- create paid support or SLA commitments;
- change consumer statuses;
- weaken agency auth boundaries.

## Likely Files

- `internal/*_test.go`
- `cmd/agency-config/*_test.go`
- `testdata/`
- `docs/multi-agency-strategy.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-27.md`
