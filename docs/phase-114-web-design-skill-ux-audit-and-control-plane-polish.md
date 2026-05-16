# Phase 114 -- Web Design Skill UX Audit And Control Plane Polish

## Goal

Use the Web Design Skill to audit and polish the private Operations Console
before the public release-candidate gate.

Phase 114 must improve first-run/operator clarity without changing public feed
contracts, auth boundaries, GTFS-RT semantics, release status, consumer
statuses, or evidence status.

## Current Repo Context

- Phase 113 is closed with bounded local fresh-clone and local source-archive
  install-confidence passes.
- Publication remains `blocked_public_distribution_review` because the source
  archive contains tracked protected evidence and consumer-submission paths.
- The Operations Console is intentionally Go server-rendered with buildless
  progressive JavaScript under authenticated private admin routes.
- The Web Design Skill file at
  `/Users/edwintse/.agents/skills/web-design-engineer/SKILL.md` was read for
  this phase.

## Web Design Skill Decisions

Design Decisions:

- Color palette: preserve the existing restrained operations palette and use
  existing semantic status colors rather than adding new decorative colors.
- Typography: preserve the repo's system-font server-rendered UI style; keep
  headings compact inside dashboards and reserve large display type for public
  docs only.
- Spacing system: preserve compact operator-console density with small
  multiples and improve grouping where it reduces scan burden.
- Border-radius strategy: keep cards and controls small-radius and utilitarian.
- Shadow hierarchy: keep elevation minimal; prefer borders, grouping, and
  table structure over decorative depth.
- Motion style: no new motion unless it improves state feedback; the Phase 114
  target is clarity and reliable no-JS behavior.

## Web Design Skill Findings To Triage

The Phase 111 Web Design Skill UX sub-agent identified these Phase 114/118
candidates:

- Start Here is dense for first-time operators.
- Primary actions are mixed with JSON/debug links.
- Missing feed URLs can produce a copy affordance for the literal value
  `missing`.
- Wide tables are hard to scan.
- Repeated explanatory panels can push core tables below the fold.
- Realtime labels such as `VP/TU/Alerts` should use plain-language feed names.
- Accessibility basics are solid and should be preserved.

## Scope

- Add the required Phase 114 UX artifact under `docs/ux/`.
- Inspect and, where safe, polish the private Operations Console release /
  first-run / feed-health path.
- Fix repo-caused UX defects that could confuse a nontechnical operator before
  release evaluation, especially missing-copy affordances and unclear realtime
  labels.
- Keep changes small and compatible with no-JS, private-auth, and buildless
  JavaScript behavior.

## Boundaries

- Do not touch protected evidence paths:
  - `docs/evidence/captured/**`
  - `docs/evidence/consumer-submissions/status.json`
  - `docs/evidence/consumer-submissions/current/**`
  - `docs/evidence/consumer-submissions/artifacts/**`
  - `docs/evidence/consumer-submissions/packets/**`
- Do not move consumer statuses.
- Do not add a frontend build system or heavy frontend framework.
- Do not add public admin routes, evidence collection, external contact,
  release publication, tag creation, asset upload, or stronger public claims.
- Do not claim production readiness, compliance, adoption, consumer acceptance,
  hosted service availability, SLA/uptime, vendor compatibility, hardware
  certification, production AVL reliability, production-grade ETA quality, or
  real-world ETA accuracy.

## Deliverables

- `docs/ux/web-design-skill-review-phase-114.md`
- Focused Operations Console polish, if code inspection confirms safe changes.
- Tests for any meaningful UI/runtime behavior change.
- `docs/handoffs/phase-114.md`
- Source-of-truth status updates.

## Implementation Plan

1. Add this Phase 114 plan and commit checkpoint 000001.
2. Audit relevant Operations Console code and write the Web Design Skill UX
   review artifact.
3. Implement narrowly scoped UX polish with tests.
4. Run focused route/UI tests and baseline validation; patch required gaps.
5. Close Phase 114 with handoff/status docs and continue immediately to Phase
   115.

## Checkpoint Plan

- `Phase 114 -- Checkpoint 000001: add web design skill ux audit and control plane polish plan`
- `Phase 114 -- Checkpoint 000002: implement or audit primary scoped work`
- `Phase 114 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 114 -- Checkpoint 000004: close web design skill ux audit and control plane polish review`

## Focused Validation Targets

- `go test ./cmd/agency-config -run 'FirstRun|FeedHealth|RealtimeOperations|Operations|ValidationCenter'`
- JavaScript runtime tests if touched.
- `git diff --check`
- `make check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact prepared-only consumer tracker assertion
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`

Because Phase 114 may touch UI runtime code, also run:

- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`

## Checkpoint Report -- 000001

Checkpoint:
Phase 114 -- Checkpoint 000001: add web design skill ux audit and control
plane polish plan.

Goal status:
Active. Phase 113 is closed and Phase 114 has started.

Sub-agents used or simulated:
Web Design Skill UX sub-agent findings from Phase 111 were incorporated. The
Web Design Skill file was read for this phase. Planning, Implementation, QA,
Documentation / IA, Claim-Boundary, Security/Auth, Data/Migration,
Release/Supply-Chain, Install Confidence, and GTFS-RT Domain roles are
simulated by the Master Agent for this plan checkpoint.

Changed files:
`docs/phase-114-web-design-skill-ux-audit-and-control-plane-polish.md`.

Validation run:
Initial Phase 114 inspection reviewed the Web Design Skill, Phase 114 prompt,
UX validation policy, master/sub-agent operating manual, and claim-boundary
validation policy. Focused checkpoint validation is scheduled before commit.

Blocked checks:
Implementation, route tests, browser/visual review, and full validation are
scheduled for later Phase 114 checkpoints. Release publication, retained
evidence, external contact, consumer status movement, protected path writes,
and stronger claims remain out of Phase 114 scope.

Protected path status:
No protected evidence path is part of the plan. The plan forbids protected
path writes.

Consumer tracker status:
The consumer tracker is not part of the plan. The seven targets must remain in
order and exactly `prepared`.

Claim-boundary status:
The plan explicitly forbids stable release readiness, production readiness,
compliance, adoption, consumer acceptance, final-root, hosted-service,
paid-support, SLA/uptime, vendor, hardware, production AVL, and ETA-quality
claims.

Security/auth status:
The plan keeps private Operations Console auth boundaries unchanged and does
not change CSRF, credentials, tokens, public exposure, private payload
handling, external contact, or command execution behavior.

Data/migration status:
No migration, schema, durable state, dependency, public feed contract, or Go
module change is planned.

Release/publication status:
No release action is planned for Phase 114. Publication remains blocked by
Phase 112 source archive public-distribution review.

Install confidence status:
Phase 113 local install-confidence diagnostics passed for bounded local clone
and archive replay. No new install-confidence claim is made by this plan.

Web design skill status:
The skill was used for Phase 114 planning. The required UX review artifact is
scheduled for checkpoint 000002.

Master review:
Approved. The phase is scoped to small, release-relevant operator UX polish
and preserves the repo's private server-rendered control-plane architecture.

Required edits:
Run checkpoint 000001 validation and commit, then audit the relevant
Operations Console code and create the UX review artifact.

Decision:
Proceed to checkpoint 000001 validation and commit.

Next checkpoint:
Phase 114 -- Checkpoint 000002: implement or audit primary scoped work.

## Checkpoint Report -- 000002

Checkpoint:
Phase 114 -- Checkpoint 000002: implement or audit primary scoped work.

Goal status:
Active. Phase 114 UX audit and primary polish are implemented.

Sub-agents used or simulated:
Web Design Skill UX findings from Phase 111 were incorporated. The
Web Design Skill file was read and applied. Implementation, QA,
Documentation / IA, Claim-Boundary, Security/Auth, Data/Migration,
Release/Supply-Chain, Install Confidence, and GTFS-RT Domain roles were
simulated by the Master Agent for this checkpoint.

Changed files:
`cmd/agency-config/operations_first_run.go`;
`cmd/agency-config/operations.go`;
`cmd/agency-config/operations_admin.js`;
`cmd/agency-config/operations_admin_test.mjs`;
`cmd/agency-config/main_test.go`;
`docs/ux/web-design-skill-review-phase-114.md`;
`docs/phase-114-web-design-skill-ux-audit-and-control-plane-polish.md`.

Validation run:
`node --test cmd/agency-config/operations_admin_test.mjs` passed. Focused
Operations Console Go tests passed. `make agency-app-up` started the local app.
Authenticated local HTML review confirmed the plain-language realtime label and
no `data-copy-value="missing"` marker. Playwright CLI captured desktop and
mobile Operations Console screenshots into ignored `.cache/phase-114-ux/`.

Blocked checks:
Full baseline validation is scheduled for checkpoint 000003. Public release
publication remains blocked by the Phase 112 source archive
public-distribution review. Release publication, retained evidence, external
contact, consumer status movement, protected path writes, and stronger claims
remain out of Phase 114 scope.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.

Consumer tracker status:
The tracker was not edited. The exact seven consumer targets must remain in
order and all remain `prepared`.

Claim-boundary status:
The UX artifact and code changes preserve forbidden claim boundaries and make
no stronger public claim.

Security/auth status:
No route, auth behavior, CSRF behavior, credential handling, token handling,
public exposure, private payload handling, external contact, or command
execution behavior changed. The Playwright auth storage state was local,
ignored, and not committed.

Data/migration status:
No migration, schema, durable state, dependency, public feed contract, or Go
module change was added.

Release/publication status:
No release action was taken. Publication remains blocked by the Phase 112
public-distribution review.

Install confidence status:
No new install-confidence claim was made. Phase 113 remains the current local
install-confidence result.

Web design skill status:
Phase 114 used the Web Design Skill and added
`docs/ux/web-design-skill-review-phase-114.md`.

Master review:
Approved. The implementation fixes release-relevant copy affordance and
plain-language label issues while preserving private no-JS-safe behavior.

Required edits:
Commit checkpoint 000002, then run full Phase 114 validation.

Decision:
Proceed to checkpoint 000002 commit and checkpoint 000003 validation.

Next checkpoint:
Phase 114 -- Checkpoint 000003: run validation and patch required gaps.

