# Phase 76 Handoff -- Design System And App Shell

## Status

Phase 76 is complete for the private Operations Console design-system and app
shell scope.

The phase kept the existing Go server-rendered Operations Console and private
route set. It did not implement a SPA, public admin route, migration, API
behavior change, evidence write, consumer tracker change, release artifact, or
stronger public claim.

## Required Closeout Report

```text
Phase: 76 -- Design System And App Shell
Sub-agents used or simulated, including intended model level:
- Master Agent, GPT-5.5 x-high, real parent orchestration.
- Context / Repo Truth Sub-Agent, GPT-5.5 x-high, simulated because the
  available sub-agent thread limit was reached.
- Planning Sub-Agent, GPT-5.5 x-high, simulated because the available
  sub-agent thread limit was reached.
- Implementation Sub-Agent, GPT-5.5 high, real read-only review.
- QA Sub-Agent, GPT-5.5 high, real read-only review.
- UI/UX Sub-Agent, GPT-5.5 high, real read-only review.
- Documentation / IA Sub-Agent, GPT-5.5 high, real read-only review.
- Claim-Boundary Sub-Agent, GPT-5.5 high, real read-only review.
- Security/Auth Sub-Agent, GPT-5.5 high, real read-only review.
Goal:
- Replace primitive Operations Console presentation with a coherent private
  app shell and reusable design tokens while preserving routes, auth, CSRF,
  claim boundaries, and server-rendered Go templates.
Changed files:
- `README.md`
- `cmd/agency-config/main_test.go`
- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_cockpit.go`
- `cmd/agency-config/operations_design_system.go`
- `cmd/agency-config/operations_first_run.go`
- `cmd/agency-config/operations_launchpad.go`
- `cmd/agency-config/operations_maintenance.go`
- `cmd/agency-config/operations_navigation.go`
- `docs/README.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-76.md`
- `docs/phase-76-design-system-and-app-shell.md`
- `docs/roadmap-status.md`
- `docs/tutorials/no-cli-agency-first-run.md`
- `docs/tutorials/small-agency-maintenance-guide.md`
- `wiki/README.md`
- `wiki/browser-first-setup.md`
- `wiki/operations-console-tour.md`
Routes added/changed:
- No route URLs added or renamed.
- Existing private Operations Console pages now share static design tokens,
  a common private control-plane header, breadcrumb, metadata row, route-stable
  navigation, and main landmark.
- Navigation labels now align to Start Here, Schedule, Realtime, Connectors,
  Health, Maintain, and Learn.
- GTFS Studio and Alerts Console links are visibly marked as separate private
  admin surfaces.
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
- `go test ./cmd/agency-config`
- `go test ./internal/auth ./internal/tenant ./cmd/feed-alerts ./cmd/gtfs-studio`
- `make validate`
- `make test`
- `RUN_LOCAL_APP=true make release-candidate-check`
Blocked checks:
- None for Phase 76 closeout.
- `RUN_LOCAL_APP=true make release-candidate-check` exited 0 and recorded
  `overall_status=needs_review` because the checkout was dirty during local
  diagnostics and release-package audit remained not checked. That is a
  release-candidate diagnostic state, not a Phase 76 blocker.
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
- Generic `Ready` UI status was changed to `Ready for local review`.
- `Consumer Submission Evidence` was changed to prepared-packet tracker
  language.
- Five-feed language was qualified as configured/local/reference where needed.
- No compliance, adoption, consumer acceptance, final-root, hosted SaaS,
  production-readiness, vendor-compatibility, hardware-certification, SLA,
  uptime, public-launch, or production-grade ETA claim was added.
Security/auth status:
- No route, auth, role, CSRF, method, agency-scope, form action, or device-token
  behavior changed.
- Focused auth/security packages and adjacent admin commands passed tests.
Accessibility status:
- Shared layout retains skip link, landmarks, `:focus-visible`, mobile table
  overflow behavior, responsive navigation, and reduced-motion CSS.
- New tests cover core route shell markers, design tokens, external-admin
  markers, and fragile decoration avoidance.
Docs/site/wiki alignment:
- README, docs home, wiki home, Operations Console tour, browser-first setup,
  no-CLI first-run tutorial, and maintenance guide now use configured feed URL
  wording and the updated route group IA.
Commit list:
- Phase 76 -- Checkpoint 000001: add design system and app shell plan
- Phase 76 -- Checkpoint 000002: implement shared layout tokens and components
- Phase 76 -- Checkpoint 000003: apply shell to core Operations Console routes
- Phase 76 -- Checkpoint 000004: add responsive and accessibility baseline checks
- Phase 76 -- Checkpoint 000005: close design system and app shell review
Master review:
- Approved after required sub-agent edits were incorporated.
Required edits:
- None remaining for Phase 76.
Decision:
- Phase 76 closed. Continue automatically to Phase 77 -- Admin Control API And
  Command Model under the authorized Phase 75-90 product track.
Next phase:
- Phase 77 -- Admin Control API And Command Model.
```
