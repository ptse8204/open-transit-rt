# Phase 94 -- Operations Console Architecture Refactor

## Scope

Phase 94 reduces fragility in private Operations Console route, navigation,
template, and audit growth. The phase is intentionally refactor-first: no
public admin routes, no new product behavior, and no route behavior changes
unless a focused test covers the expected stable behavior.

## Current Architecture Risk

Route truth is repeated across:

- `cmd/agency-config/main.go` explicit `mux.Handle` registrations and the
  `/admin/operations/` catch-all.
- `cmd/agency-config/operations.go` `operationsRoot` switch cases and
  renderer selection.
- `cmd/agency-config/operations_navigation.go` navigation groups, labels,
  hrefs, sections, external-surface markers, and page-title switch.
- `scripts/audit-operations-route-inventory.sh` canonical route inventory.
- Tests and route-map docs.

This duplication was acceptable during rapid Control Plane growth, but it now
creates avoidable drift risk before additional Phase 96-110 feature work.

## Target Refactor

The primary implementation target is a central private Operations Console route
registry with shared route metadata:

- section ID;
- canonical HTML path;
- optional canonical JSON path;
- nav label;
- page title;
- nav group;
- external-admin-surface marker;
- cache/no-store expectation;
- method posture for read-only, command, and existing bounded mutation routes.

The registry should feed navigation and titles first, then support the local
route inventory audit helper. Handler registration and renderer dispatch may
remain explicit where it protects method, body-size, CSRF, and role behavior.

## Checkpoints

### Checkpoint 000001 -- Plan

Deliverables:

- Add this plan.
- Record Phase 94 boundaries, intended sub-agent roles, validation, and
  Master approval.

Validation:

- `git status --short`
- `git diff --check`
- `make check`
- `make audit-operations-route-inventory`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- consumer tracker JSON parse and exact prepared-only assertion
- protected-path status check

### Checkpoint 000002 -- Primary Scoped Refactor

Deliverables:

- Add a central Operations Console route registry in Go.
- Refactor `operationsNavGroups` and `operationsPageTitle` to use registry
  metadata.
- Preserve external admin surface labels for GTFS Studio and Alerts Console.
- Add route-registry helpers that expose canonical HTML routes, JSON routes,
  command routes, and external admin surfaces for tests/audits.
- Keep renderer dispatch and route registration behavior stable unless the
  change is mechanical and covered by focused tests.

Expected files:

- `cmd/agency-config/operations_route_registry.go`
- `cmd/agency-config/operations_navigation.go`
- `cmd/agency-config/main_test.go`
- optionally `scripts/audit-operations-route-inventory.sh`
- this phase doc

Validation:

- `gofmt`
- focused `go test ./cmd/agency-config -run
  'TestOperationsConsoleNavigationIsGroupedAndRouteStable|TestOperationsRouteTitlesAndFirstClickLabelOrder|TestOperationsRouteRegistry'`
- `make audit-operations-route-inventory`
- `make test-operations-route-inventory`
- baseline claim/protected-path checks

### Checkpoint 000003 -- Validation And Gap Patching

Deliverables:

- Run heavier validations for code/script changes.
- Patch required drift or audit gaps found by the registry integration.
- Record exact blocked checks, if any.

Validation:

- `git status --short`
- `git diff --check`
- `make check`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`
- `make audit-operations-route-inventory`
- `make test-operations-route-inventory`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- consumer tracker JSON parse and exact prepared-only assertion
- protected-path status check

### Checkpoint 000004 -- Closeout

Deliverables:

- Add `docs/handoffs/phase-94.md`.
- Update `docs/handoffs/latest.md`, `docs/current-status.md`,
  `docs/roadmap-status.md`, and
  `docs/open-transit-rt-master-planner-remaining-work.md`.
- Record final Phase 94 result, validation, blockers, protected-path status,
  consumer tracker status, claim-boundary status, security/auth status, and
  data/migration status.

## Hard Boundaries

- Do not modify protected evidence paths.
- Do not edit or reformat
  `docs/evidence/consumer-submissions/status.json`.
- Do not move consumer statuses beyond `prepared`.
- Do not contact external services, consumers, agencies, vendors, or portals.
- Do not use real credentials, private payloads, real AVL data, or real agency
  data.
- Do not tag, publish, create a GitHub Release, push an image, or publish a
  package.
- Do not claim release readiness, compliance, adoption, consumer acceptance,
  final-root readiness, hosted-service availability, production readiness,
  vendor compatibility, hardware certification, SLA/uptime, or ETA quality.
- Do not weaken admin auth, role checks, CSRF behavior, request-body caps, or
  no-store cache handling.

## Master Approval

The Master Agent approves this Phase 94 plan. Implementation may proceed after
Checkpoint 000001 is committed, using the central route registry as the narrow
refactor seam and preserving behavior unless focused tests prove equivalence.

## Checkpoint 000001 Report

Checkpoint:
Phase 94 -- Checkpoint 000001: add operations console architecture refactor
plan.

Sub-agents used or simulated, including intended model level:
Real Context / Repo Truth Sub-Agent -- GPT-5.5 x-high, Planning Sub-Agent --
GPT-5.5 x-high, UI/UX / Documentation IA Sub-Agent -- GPT-5.5 high, and
Claim-Boundary / Security Sub-Agent -- GPT-5.5 high were launched for Phase
94. Implementation and QA roles are simulated by the Master Agent until the
plan checkpoint is committed. Master Agent -- GPT-5.5 x-high, current thread.

Changed files:
`docs/phase-94-operations-console-architecture-refactor.md`.

Validation run:
`git status --short`; `git diff --check`; `make check`; `make
audit-operations-route-inventory`; `make audit-product-acceptance`; `make
audit-final-claim-review`; consumer tracker JSON parse; exact prepared-only
consumer tracker assertion; protected-path status check.

Blocked checks:
None for this docs-only planning checkpoint.

Protected path status:
No protected evidence path is edited or generated.

Consumer tracker status:
`docs/evidence/consumer-submissions/status.json` is not edited. All seven
consumer targets must remain exactly `prepared`.

Claim-boundary status:
This plan describes private architecture refactoring only. It makes no release
readiness, compliance, adoption, consumer acceptance, production readiness,
final-root readiness, hosted-service availability, vendor compatibility,
hardware certification, SLA/uptime, or ETA-quality claim.

Security/auth status:
The plan explicitly preserves private admin auth, role checks, CSRF behavior,
request-size limits, and no-store cache handling.

Data/migration status:
No persistence, migration, GTFS data model, tenant model, or realtime data
model change is planned.

Master review:
Approved. The route registry is a narrow, testable seam that addresses the
Phase 91 drift finding without adding product scope.

Required edits:
None.

Decision:
Commit CP000001, then implement the primary scoped refactor.

Next checkpoint:
Phase 94 -- Checkpoint 000002: implement primary scoped work.
