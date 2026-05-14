# Phase 90 -- Final Control Plane Closeout And Future Evidence Gate Stubs

## Purpose

Phase 90 closes the authorized Phase 75-90 Consumer-Grade Control Plane product
track and leaves the repository ready for maintainer review.

This phase is a final closeout and planning-artifact phase. It may summarize
product capabilities, routes, validation, blockers, protected-path status,
consumer tracker status, and future optional evidence gates. It must not
collect evidence, contact external parties, tag or publish a release, create
packages or images, move consumer statuses, or add stronger public claims.

## Current Truth

- Phases 0-89 are complete for their bounded scopes.
- Phase 72 remains `needs_review`, not release-ready.
- Phase 74 CP000008 remains the latest GitHub Pages publication at commit
  `a8b250e`.
- Phase 89 closed the local `v0.1.0-rc.1` gate as `needs_review`: local
  product, route, connector, backend, product-acceptance, and claim-boundary
  diagnostics passed where authorized, but release packaging, package audit,
  tagging, image publication, and release actions remain blocked/not checked.
- Evidence/adoption/compliance tracks remain optional and require separate
  written authorization.
- All seven consumer targets remain exactly `prepared`: Google Maps, Apple
  Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land.

## Sub-Agent Plan

Real sub-agent spawning is unavailable in this run because the app thread
limit has been reached. Use simulated roles with the intended model levels:

- Master Agent -- GPT-5.5 x-high, simulated: owns scope, checkpoint sequence,
  protected paths, consumer tracker, claim boundaries, and final approval.
- Context / Repo Truth Sub-Agent -- GPT-5.5 x-high, simulated: confirm Phase
  89 closeout, source-of-truth status, protected paths, consumer tracker, and
  Phase 90 deliverables before edits.
- Planning Sub-Agent -- GPT-5.5 x-high, simulated: keep Phase 90 to final
  closeout docs and future gate stubs only.
- Implementation Sub-Agent -- GPT-5.5 high, simulated: update only approved
  docs and final closeout artifacts.
- QA Sub-Agent -- GPT-5.5 high, simulated: run baseline validation, product
  and claim audits, exact prepared-only tracker checks, and protected-path
  checks.
- UI/UX Sub-Agent -- GPT-5.5 high, simulated: review route and feature
  inventories for operator-facing clarity without adding UI.
- Documentation / IA Sub-Agent -- GPT-5.5 high, simulated: make the final
  status, route inventory, feature inventory, validation matrix, blocker
  matrix, and future gate stubs easy to scan.
- Claim-Boundary Sub-Agent -- GPT-5.5 high, simulated: block unsupported
  compliance, adoption, consumer, final-root, hosted-service, production,
  vendor, hardware, SLA/uptime, ETA-quality, and release-ready claims.
- Security/Auth Sub-Agent -- GPT-5.5 high, simulated: verify no new route,
  command, credential, portal, evidence, release, package, or browser
  mutation action is added.
- Data/Migration Sub-Agent -- GPT-5.5 high, simulated for persistence review;
  no migration is expected.

## Master Approval Before Implementation

Approved bounded scope:

- Add a final Phase 75-90 status artifact with route inventory, feature
  inventory, validation matrix, blocker matrix, protected-path review,
  consumer tracker review, claim-boundary review, and future optional evidence
  gate stubs.
- Add the Phase 90 handoff.
- Update current source-of-truth docs to mark Phase 90 complete only after
  final validation passes or blockers are recorded.
- Keep each evidence gate as a future stub that says separate written
  authorization is required.

Required edits before implementation: none.

## Checkpoints

### Checkpoint 000001 -- Plan

Add this Phase 90 plan and link it from current source-of-truth docs as the
active final closeout plan.

Expected files:

- `docs/phase-90-final-control-plane-closeout-and-future-evidence-gate-stubs.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/roadmap-status.md`
- `docs/open-transit-rt-master-planner-remaining-work.md`

### Checkpoint 000002 -- Final Inventories And Evidence Gate Stubs

Create the final control-plane status artifact.

Expected content:

- final control-plane status summary;
- route inventory;
- feature inventory;
- blocker matrix;
- protected-path review;
- consumer tracker review;
- claim-boundary review;
- future optional evidence gate stubs for final-root proof, consumer
  submission, real agency pilot, real vendor/device AVL, real-world ETA
  quality, and compliance packet, each requiring separate written
  authorization.

### Checkpoint 000003 -- Final Validation Matrix

Run and record final validation.

Expected checks:

- `git status --short`
- `git diff --check`
- `make check`
- `make validate`
- `make test`
- `RUN_LOCAL_APP=true make release-candidate-check`
- `make external-connection-check`
- `make adapter-conformance`
- `make test-connector-examples`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`

Blocked unless separately authorized:

- `make release-package`
- `make audit-release-package`

### Checkpoint 000004 -- Closeout

Write Phase 90 handoff and update current source-of-truth docs. Confirm:

- Phases 75-90 completed;
- no protected path writes;
- exact prepared-only consumer tracker;
- no release tag;
- no published package;
- no published image;
- no retained evidence collection;
- no forbidden claims;
- future evidence gates remain authorization-gated.

## Protected Paths

Do not modify or generate files under:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/**`
- `docs/evidence/consumer-submissions/artifacts/**`
- `docs/evidence/consumer-submissions/packets/**`

## Claim Boundaries

Phase 90 may say:

- private browser-first control-plane product work completed for the
  authorized Phase 75-90 track;
- local diagnostics passed where run;
- release-candidate status remains `needs_review`;
- future evidence gates require separate written authorization.

Phase 90 must not claim:

- CAL-ITP/Caltrans compliance;
- agency adoption or approval;
- consumer submission, review, acceptance, ingestion, listing, or display;
- final-root readiness;
- hosted SaaS;
- paid support;
- SLA or uptime guarantee;
- production readiness;
- vendor compatibility;
- hardware certification;
- production-grade ETA quality;
- real-world ETA accuracy;
- public launch completion;
- release readiness.

## Closeout Report Requirements

The closeout must include:

```text
Phase:
Sub-agents used or simulated, including intended model level:
Goal:
Changed files:
Routes added/changed:
Commands added/changed:
Migrations:
Validation run:
Blocked checks:
Known blockers:
Protected path status:
Consumer tracker status:
Claim-boundary status:
Security/auth status:
Accessibility status:
Docs/site/wiki alignment:
Commit list:
Master review:
Required edits:
Decision:
Next phase:
```
