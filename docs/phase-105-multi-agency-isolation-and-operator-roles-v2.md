# Phase 105 -- Multi-Agency Isolation And Operator Roles V2

## Goal

Strengthen repository-level agency scoping, role behavior, access-denied
guidance, metadata-only audit review, and tenant-safe public feed route checks
without claiming production multi-tenant hosting.

## Current Surface

- `internal/auth` owns principals, roles, admin middleware, access-denied
  responses, CSRF checks, and agency query conflict rejection.
- `internal/tenant` owns public agency path parsing and agency ID validation.
- `cmd/agency-config` owns the private Operations Console, route registry,
  Access & Roles review, metadata-only audit review, GTFS import/setup,
  validation, devices, alerts, and public feed discovery routes.
- Feed services already expose single-agency public routes and path-routed
  `/public/agencies/{agency_id}/...` routes for tenant-safe public feed
  delivery.
- Existing tests cover many role and agency boundaries, but coverage is
  scattered. Phase 105 should add table-driven checks and close verified gaps
  without adding a new tenancy system.

## Master-Approved Plan

1. Add the Phase 105 plan and checkpoint report.
2. Implement tests-first hardening over existing seams:
   - private Operations route agency-scope tests from the route registry;
   - role matrix tests for selected review and mutation routes;
   - access-denied UX checks for HTML/API responses, escaping, no-store, and
     no cross-agency data disclosure;
   - tenant-safe public feed route tests for path agency precedence, invalid
     encoded agency IDs, and no per-agency debug JSON exposure;
   - bounded metadata-only audit review improvements if tests reveal safe gaps.
3. Keep audit output metadata-only. Do not render raw actor secrets, reasons,
   old/new JSON diffs, tokens, paths, payloads, database URLs, cookies, or
   credentials.
4. Do not add migrations, row-level security, hosted identity provider behavior,
   new public admin routes, live external calls, real credentials, consumer
   status changes, or protected evidence writes.
5. Run focused checks, then required baseline/code-change validation.
6. Close with `docs/handoffs/phase-105.md`, `docs/handoffs/latest.md`, and
   roadmap/status updates.

## Non-Goals

- No production multi-tenant hosting claim.
- No hosted SaaS, SLA, uptime, compliance, consumer-acceptance,
  agency-adoption, deployment-success, vendor, hardware, final-root,
  release-readiness, public-launch, or production-grade ETA claim.
- No consumer tracker edits or status movement.
- No retained evidence or protected-path writes.
- No new database migration, row-level security model, hosted login/SSO, JWT
  replay store, public admin route, or external identity integration.
- No exposure of raw audit details, raw support artifacts, credentials, private
  paths, or cross-agency payloads.

## Checkpoint Plan

- `Phase 105 -- Checkpoint 000001: add multi-agency isolation and operator roles v2 plan`
- `Phase 105 -- Checkpoint 000002: implement primary scoped work`
- `Phase 105 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 105 -- Checkpoint 000004: close multi-agency isolation and operator roles v2 review`

## Focused Validation Targets

- `go test ./internal/auth ./internal/tenant`
- `go test ./cmd/agency-config -run 'Agency|Access|Audit|Role|Route|Public'`
- `go test ./cmd/feed-vehicle-positions ./cmd/feed-trip-updates ./cmd/feed-alerts ./cmd/gtfs-studio ./internal/compliance`
- `scripts/audit-operations-route-inventory.sh`
- `scripts/test-multi-agency-hosting.sh`

Because this phase is expected to change code/docs/tests, closeout also
requires:

- `git status --short`
- `git diff --check`
- `make check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact prepared-only consumer tracker assertion
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`

## Checkpoint Report -- 000001

Checkpoint:
Phase 105 -- Checkpoint 000001: add multi-agency isolation and operator roles
v2 plan.

Sub-agents used or simulated, including intended model level:
Real Context / Repo Truth Sub-Agent -- GPT-5.5 x-high; real Planning
Sub-Agent -- GPT-5.5 x-high. Implementation, QA, UI/UX, Documentation / IA,
Claim-Boundary, Security/Auth, Data/Migration, and Release/Supply-Chain roles
are simulated by the Master Agent for this plan checkpoint. Master Agent --
GPT-5.5 x-high, current thread.

Changed files:
`docs/phase-105-multi-agency-isolation-and-operator-roles-v2.md`.

Validation run:
Initial Phase 105 repository inspection found existing auth, tenant route,
operations route registry, access UX, audit review, and public feed route
seams. Sub-agents completed read-only inspection and planning. After adding the
plan, `git status --short` showed only
`docs/phase-105-multi-agency-isolation-and-operator-roles-v2.md`; `git diff
--check` passed; `python3 -m json.tool
docs/evidence/consumer-submissions/status.json >/dev/null` passed; the exact
prepared-only consumer tracker assertion passed; and `git status --short --
docs/evidence/consumer-submissions docs/evidence/captured db/migrations
go.mod go.sum` returned no output.

Blocked checks:
Implementation, focused tests, and closeout baseline checks are not yet run
because this checkpoint only approves the Phase 105 plan. Release-candidate
checks, package generation/audit, evidence collection, publication, hosted
identity, production tenancy changes, and consumer actions are out of scope.

Protected path status:
No protected evidence path is part of the plan. The plan forbids protected path
writes.

Consumer tracker status:
The consumer tracker is not part of the plan. The seven targets must remain in
order and `prepared`.

Claim-boundary status:
The plan explicitly forbids production multi-tenant hosting, hosted SaaS, SLA,
uptime, compliance, consumer-acceptance, agency-adoption, deployment-success,
vendor, hardware, final-root, release-readiness, public-launch, and
production-grade ETA claims.

Security/auth status:
The plan preserves existing admin auth, role checks, CSRF behavior, agency
query conflict rejection, private access UX, and no-store/private-route
boundaries.

Data/migration status:
No migration, row-level security model, hosted identity store, audit schema
change, durable tenancy state, or module dependency change is planned.

Master review:
Approved. The smallest safe Phase 105 implementation is tests-first hardening
and bounded metadata-only UX over existing auth, tenant, route, and audit
seams.

Required edits:
Add scoped tests, patch verified gaps only, update docs if behavior or
operator guidance changes, and record validation results.

Decision:
Proceed to implementation checkpoint 000002 after plan validation and commit.

Next checkpoint:
Phase 105 -- Checkpoint 000002: implement primary scoped work.
