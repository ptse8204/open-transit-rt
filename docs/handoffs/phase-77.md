# Phase 77 Handoff -- Admin Control API And Command Model

## Status

Phase 77 is complete for the private Admin Control API and command-model
scope.

The phase added a bounded internal command result contract and migrated one
low-risk, read-only validation-health refresh workflow into that model. It did
not expose a public API, add arbitrary command execution, add a migration,
write retained evidence, move consumer tracker status, publish a release
artifact, or make stronger public claims.

## Required Closeout Report

```text
Phase: 77 -- Admin Control API And Command Model
Sub-agents used or simulated, including intended model level:
- Master Agent, GPT-5.5 x-high, real parent orchestration.
- Context / Repo Truth Sub-Agent, GPT-5.5 x-high, real read-only review.
- Planning Sub-Agent, GPT-5.5 x-high, real read-only review.
- Implementation Sub-Agent, GPT-5.5 high, real read-only review.
- QA Sub-Agent, GPT-5.5 high, real read-only review.
- UI/UX Sub-Agent, GPT-5.5 high, simulated by the Master Agent for the
  command outcome wording and existing private page copy.
- Documentation / IA Sub-Agent, GPT-5.5 high, simulated by the Master Agent
  for command-model docs and source-of-truth alignment.
- Claim-Boundary Sub-Agent, GPT-5.5 high, real read-only review.
- Security/Auth Sub-Agent, GPT-5.5 high, real read-only review.
Goal:
- Create a safe private command/query model for browser-controlled backend
  workflows and migrate one low-risk workflow without exposing public APIs,
  raw validator output, arbitrary backend execution, evidence writes, or
  stronger claims.
Changed files:
- `cmd/agency-config/main.go`
- `cmd/agency-config/main_test.go`
- `cmd/agency-config/operations.go`
- `docs/README.md`
- `docs/admin-command-model.md`
- `docs/current-status.md`
- `docs/decisions.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-77.md`
- `docs/phase-77-admin-control-api-and-command-model.md`
- `docs/roadmap-status.md`
- `internal/admincontrol/model.go`
- `internal/admincontrol/model_test.go`
Routes added/changed:
- Added private `POST /admin/operations/validation-health/refresh.json`.
- The new route is authenticated, agency-scoped, role-checked, request-capped,
  strict about unsupported execution fields, and `Cache-Control: no-store`.
- Cookie-auth POSTs require CSRF.
- The existing private `/admin/operations/validation-health` page remains
  route-stable and now explains the command-model boundary.
- No public route was added or changed.
Commands added/changed:
- None.
Migrations:
- None.
Validation run:
- `git status --short`
- `git diff --check`
- `make check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`
- `go test ./internal/admincontrol ./cmd/agency-config`
- `go test ./internal/auth ./internal/tenant ./cmd/feed-alerts ./cmd/gtfs-studio`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`
Blocked checks:
- None.
Known blockers:
- Phase 72 remains `needs_review`, not release-ready.
- No release tag, package, published image, retained evidence, final-root
  proof, consumer submission, vendor proof, SLA proof, or production ETA proof
  exists.
Protected path status:
- No status under `docs/evidence/consumer-submissions`,
  `docs/evidence/captured`, `db/migrations`, `go.mod`, or `go.sum`.
Consumer tracker status:
- Exact seven-target check passed. Google Maps, Apple Maps, Transit App,
  Bing Maps, Moovit, Mobility Database, and transit.land remain `prepared`.
Claim-boundary status:
- Product acceptance and final claim audits passed.
- `validation_health.refresh` returns all-false claim flags and writes
  nothing.
- `validation_health.run_all` is documented as a private diagnostic write
  because it may store normal `validation_report` rows; it is not evidence or
  compliance proof.
- No CAL-ITP/Caltrans compliance, adoption, consumer acceptance, final-root,
  hosted SaaS, production-readiness, vendor-compatibility, hardware
  certification, SLA, uptime, public-launch, or production-grade ETA claim was
  added.
Security/auth status:
- The new refresh JSON route is private, role-checked for read-only/operator/
  editor/admin access, agency-scoped, body-capped, CSRF-protected for cookie
  auth, and strict about unsupported command/path/tooling fields.
- The existing validator run workflow remains admin-only.
- Command results are bounded and redact raw reports, stdout/stderr,
  argument arrays, private paths, tokens, cookies, and database URLs.
Accessibility status:
- No route layout or control pattern changed outside the existing private
  Operations Console shell.
- The validation-health page added text cards inside the existing accessible
  shell; no new custom widget or keyboard-only dependency was introduced.
Docs/site/wiki alignment:
- `docs/admin-command-model.md` is now the command-model reference.
- `docs/README.md`, `docs/current-status.md`, `docs/handoffs/latest.md`,
  `docs/roadmap-status.md`, and `docs/decisions.md` align on the Phase 77
  private-command boundary.
Commit list:
- Phase 77 -- Checkpoint 000001: add admin control API and command model plan
- Phase 77 -- Checkpoint 000002: define private command result contracts
- Phase 77 -- Checkpoint 000003: add command safety tests
- Phase 77 -- Checkpoint 000004: migrate validation health refresh command
- Phase 77 -- Checkpoint 000005: close admin control API review
Master review:
- Approved after the Planning, QA, Security/Auth, and Claim-Boundary required
  edits were incorporated.
Required edits:
- None remaining for Phase 77.
Decision:
- Phase 77 closed. Continue automatically to Phase 78 -- Frontend Routing,
  State, And Data Loading under the authorized Phase 75-90 product track.
Next phase:
- Phase 78 -- Frontend Routing, State, And Data Loading.
```
