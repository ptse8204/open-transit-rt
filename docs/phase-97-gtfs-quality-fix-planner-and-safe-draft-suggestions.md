# Phase 97 - GTFS Quality Fix Planner And Safe Draft Suggestions

## Status

Status: planned

Phase 97 turns existing GTFS importer and validator notices into a private
operator-facing fix plan. The implementation must keep all edits in source GTFS,
GTFS Studio review, import, publish, and validation flows that already exist. It
must not add automatic production GTFS edits, browser-supplied validator
execution fields, protected evidence writes, consumer-status changes, or
compliance/readiness claims.

## Master Phase Plan

1. Inspect the current GTFS Quality page, validation center, GTFS Workbench, and
   GTFS Studio draft boundaries for safe reuse.
2. Extend GTFS Quality guidance with a bounded fix planner that groups issues by
   owner, affected files, before/after validation steps, and safe draft
   suggestion handling.
3. Add a private exportable checklist surface generated from sanitized grouped
   notices, with explicit "no automatic edit" and "no production edit" language.
4. Add focused tests for HTML output, side-effect boundaries, agency isolation,
   claim boundaries, and sanitized samples.
5. Run targeted tests first, then the full phase-closeout validation set after
   code changes.

## Intended Sub-Agent Roles

| Role | Model level | Use |
| --- | --- | --- |
| Context / Repo Truth | GPT-5.5 x-high | Confirm existing GTFS Quality, Validation Center, GTFS Studio, and Workbench seams. |
| Planning | GPT-5.5 x-high | Validate checkpoint structure and safe minimum scope. |
| Implementation | GPT-5.5 high | Simulated by Master unless a worker slot is available; changes are small and route-local. |
| QA | GPT-5.5 high | Simulated by Master through targeted and full validation. |
| UI/UX | GPT-5.5 high | Simulated by Master; copy must be operator-facing and clear. |
| Documentation / IA | GPT-5.5 high | Simulated by Master; update phase and handoff docs. |
| Claim-Boundary | GPT-5.5 high | Simulated by Master; no readiness/compliance/acceptance claims. |
| Security/Auth | GPT-5.5 high | Simulated by Master; preserve auth, CSRF, no browser-supplied execution fields. |
| Data/Migration | GPT-5.5 high | Simulated by Master; no migration planned unless a safe persisted suggestion model already exists. |

## Implementation Boundaries

- The planner is private Operations Console guidance, not a public feed,
  evidence artifact, agency approval, consumer acceptance, or compliance packet.
- Draft suggestions are represented as review guidance only unless an existing
  safe draft suggestion record model is discovered. No new mutation route is
  planned for this phase.
- The GTFS Quality GET route must stay read-only.
- The GTFS Quality POST route must remain admin-only, CSRF-checked for cookie
  sessions, strict-field, and limited to the existing server-owned static
  validator rerun.
- Samples must remain sanitized and bounded. Raw reports, stdout/stderr, argv,
  private paths, tokens, URLs, credentials, and private payloads must not render.
- Protected evidence paths and `docs/evidence/consumer-submissions/status.json`
  remain untouched.

## Checkpoint 000001 Report

Checkpoint: Phase 97 -- Checkpoint 000001: add gtfs quality fix planner and safe draft suggestions plan

Sub-agents used or simulated, including intended model level: Context / Repo Truth Sub-Agent GPT-5.5 x-high and Planning Sub-Agent GPT-5.5 x-high launched; Planning Sub-Agent returned a safe advisory-only checkpoint structure; Implementation, QA, UI/UX, Documentation / IA, Claim-Boundary, Security/Auth, and Data/Migration roles simulated by Master because the current work is route-local and additional thread slots are constrained.

Changed files: `docs/phase-97-gtfs-quality-fix-planner-and-safe-draft-suggestions.md`

Validation run: `git status --short` showed only this new phase doc; `git diff --check` passed; `make check` passed; `make audit-product-acceptance` passed; `make audit-final-claim-review` passed; `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` passed; prepared-only consumer tracker assertion passed; `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum` returned no protected or migration/module changes.

Blocked checks: None planned for the plan checkpoint.

Protected path status: No protected evidence paths are modified by the plan.

Consumer tracker status: Must remain exactly seven prepared targets; no tracker writes are authorized.

Claim-boundary status: Plan explicitly avoids compliance, release readiness, production readiness, adoption, vendor compatibility, hardware certification, SLA/uptime, hosted SaaS, consumer acceptance, and production-grade ETA claims.

Security/auth status: Plan preserves existing authenticated private routes, admin-only validator rerun, CSRF guard, and strict form field rejection.

Data/migration status: No migration planned. Draft suggestion records are review guidance unless an existing safe persistence boundary is discovered.

Master review: Approved to proceed because the scope improves operator comprehension without adding unsafe data mutations or public claims.

Required edits: Implement the route-local fix planner/checklist, tests, docs, and closeout handoff.

Decision: Proceed to Checkpoint 000002.

Next checkpoint: Phase 97 -- Checkpoint 000002: implement primary scoped work

## Checkpoint 000002 Report

Checkpoint: Phase 97 -- Checkpoint 000002: implement primary scoped work

Sub-agents used or simulated, including intended model level: Context / Repo Truth Sub-Agent GPT-5.5 x-high completed read-only seam review; Planning Sub-Agent GPT-5.5 x-high completed checkpoint plan; Implementation, QA, UI/UX, Documentation / IA, Claim-Boundary, Security/Auth, and Data/Migration roles simulated by Master.

Changed files: `cmd/agency-config/operations_gtfs_quality_guidance.go`, `cmd/agency-config/operations.go`, `cmd/agency-config/operations_validation_center.go`, `cmd/agency-config/operations_gtfs_workbench.go`, `cmd/agency-config/main_test.go`, `docs/phase-97-gtfs-quality-fix-planner-and-safe-draft-suggestions.md`

Validation run: `gofmt` on changed Go files; `go test ./cmd/agency-config -run 'GTFSQuality|GTFSWorkbench|ValidationCenter|Readiness|OperationsNavigation'` passed; `go test ./internal/compliance -run 'GTFSQuality|Validation'` passed; `git diff --check` passed; `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum` returned no protected, migration, or module changes.

Blocked checks: Full `make validate`, `make test`, Docker Compose config, product acceptance audit, final claim audit, and consumer tracker assertion are reserved for the validation checkpoint and phase closeout.

Protected path status: Clean; no protected evidence path was modified.

Consumer tracker status: Prepared-only tracker status preserved; tracker file was not modified.

Claim-boundary status: Added private advisory fix planner language only. Claim flags remain false, including automatic GTFS edit, draft mutation, draft suggestion record creation, schedule publish, validator semantics change, evidence creation, consumer status movement, compliance, acceptance, and production-readiness claims.

Security/auth status: No route permission, CSRF, POST, validator command, or browser-supplied execution behavior changed. GTFS Quality GET remains read-only; existing admin-only rerun boundary is unchanged.

Data/migration status: No migrations, no new storage, no GTFS Studio draft writes, no production GTFS edits, and no feed version mutations.

Master review: Approved. The implementation is derived from existing sanitized GTFS quality groups and adds operator-facing plan/checklist clarity without expanding mutation surface.

Required edits: Run full validation, patch any failures, then record validation checkpoint.

Decision: Proceed to Checkpoint 000003.

Next checkpoint: Phase 97 -- Checkpoint 000003: run validation and patch required gaps
