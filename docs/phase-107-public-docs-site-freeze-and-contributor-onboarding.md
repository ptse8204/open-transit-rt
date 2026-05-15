# Phase 107 -- Public Docs/Site Freeze And Contributor Onboarding

## Goal

Align README, docs, wiki, public site guidance, contributor onboarding, connector
contribution guidance, architecture overview, and visual asset policy around
the current agency-grade GTFS-RT product story without making launch,
adoption, compliance, release-readiness, production-readiness, hosted-service,
or consumer-acceptance claims.

## Current Repo Context

- README, docs home, and wiki home already center the browser-first path and
  private Operations Console route map.
- `CONTRIBUTING.md`, `.github/ISSUE_TEMPLATE/*`, and
  `.github/pull_request_template.md` already include secret-safety and
  truthfulness rules.
- `docs/architecture.md` is still terse and MVP-era; it should reflect the
  current modular product shape, private Operations Console surfaces,
  connector/prediction boundaries, multi-agency/role posture, and claim
  boundaries without overclaiming.
- `docs/assets/product-screenshots/README.md` and
  `docs/assets/product-diagrams/README.md` already define local/demo visual
  policy and should remain bounded.
- The public GitHub Pages site is on the `gh-pages` branch; this phase should
  not push or publish a site update unless explicitly authorized. Main-branch
  docs can describe the freeze policy.

## Scope

- Refresh README/docs/wiki route-map language where needed to include the
  current private route inventory and Phase 106 staff-training path.
- Add a docs/site freeze checklist that explains what public pages should say,
  what they must not claim, and how screenshots/diagrams should be handled.
- Improve contributor onboarding for first issues and first PRs.
- Add or improve connector contribution guidance for synthetic/local fixtures,
  manifests, adapter conformance, redaction, and no real-vendor claims.
- Update architecture overview so new contributors understand the current
  service boundaries and product modules.
- Avoid changing protected evidence paths, consumer tracker status, public feed
  contracts, runtime behavior, database schema, or release artifacts.

## Boundaries

- No public launch claim.
- No agency adoption or approval claim.
- No consumer submission, review, ingestion, listing, display, or acceptance
  claim.
- No CAL-ITP/Caltrans compliance claim.
- No final-root readiness claim.
- No hosted SaaS, SLA/uptime, production-readiness, vendor-compatibility,
  hardware-certification, release-readiness, or production-grade ETA claim.
- No generated screenshots or diagrams unless locally generated, reviewed,
  demo-only, and clearly captioned. Default is documentation policy refresh, not
  new visual capture.
- No site publish, `gh-pages` push, tag, GitHub Release, package publication,
  image push, external contact, real credentials, or retained evidence.

## Implementation Plan

1. Add this Phase 107 plan and commit checkpoint 000001.
2. Patch public docs and contributor guides:
   - README route map or docs list if stale.
   - `docs/README.md` and `wiki/README.md` for contributor/training/docs freeze
     navigation.
   - `CONTRIBUTING.md` for first issue/first PR paths and connector
     contribution rules.
   - `docs/architecture.md` for current architecture and boundaries.
   - Add a docs/site freeze checklist if no existing page cleanly covers it.
3. Run docs/claim audits and focused checks.
4. Close with `docs/handoffs/phase-107.md`, source-of-truth status updates,
   protected-path checks, and consumer tracker assertion.

## Non-Goals

- No code behavior change.
- No frontend redesign or heavy web stack.
- No new public route, private route, migration, dependency, evidence write,
  or consumer tracker change.
- No generated evidence or site publication.

## Checkpoint Plan

- `Phase 107 -- Checkpoint 000001: add public docs/site freeze and contributor onboarding plan`
- `Phase 107 -- Checkpoint 000002: implement primary scoped work`
- `Phase 107 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 107 -- Checkpoint 000004: close public docs/site freeze and contributor onboarding review`

## Focused Validation Targets

- `git diff --check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact prepared-only consumer tracker assertion
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`

Because this phase is expected to be docs-only unless a required test helper is
discovered, closeout also runs:

- `git status --short`
- `make check`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`

## Checkpoint Report -- 000001

Checkpoint:
Phase 107 -- Checkpoint 000001: add public docs/site freeze and contributor
onboarding plan.

Sub-agents used or simulated, including intended model level:
Context / Repo Truth Sub-Agent -- GPT-5.5 x-high and Planning Sub-Agent --
GPT-5.5 x-high were attempted, timed out, and were shut down without edits.
Context / Repo Truth and Planning were therefore simulated by the Master Agent
using direct repository inspection. Implementation, QA, UI/UX, Documentation /
IA, Claim-Boundary, Security/Auth, Data/Migration, and Release/Supply-Chain
roles are simulated by the Master Agent for this plan checkpoint. Master Agent
-- GPT-5.5 x-high, current thread.

Changed files:
`docs/phase-107-public-docs-site-freeze-and-contributor-onboarding.md`.

Validation run:
Initial repository inspection reviewed README, `CONTRIBUTING.md`, docs home,
wiki home, architecture overview, issue/PR templates, visual asset policy,
roadmap prompt, current status, and latest handoff. After adding the plan,
`git status --short` showed only this Phase 107 plan doc; `git diff --check`
passed; `python3 -m json.tool
docs/evidence/consumer-submissions/status.json >/dev/null` passed; the exact
prepared-only consumer tracker assertion passed; and `git status --short --
docs/evidence/consumer-submissions docs/evidence/captured db/migrations
go.mod go.sum` returned no output.

Blocked checks:
Implementation, focused audits, and closeout baseline checks are not yet run
because this checkpoint only approves the Phase 107 plan. Site publication,
release actions, evidence collection, external contact, real credentials, and
consumer actions are out of scope.

Protected path status:
No protected evidence path is part of the plan. The plan forbids protected path
writes.

Consumer tracker status:
The consumer tracker is not part of the plan. The seven targets must remain in
order and `prepared`.

Claim-boundary status:
The plan explicitly forbids public launch, adoption, agency approval, consumer
acceptance, compliance, final-root readiness, hosted SaaS, SLA/uptime,
production readiness, vendor compatibility, hardware certification, release
readiness, and production-grade ETA claims.

Security/auth status:
The plan is docs-first and forbids public route changes, credential capture,
raw private data, external contact, protected evidence writes, and site
publication.

Data/migration status:
No migration, schema change, durable state, dependency, public feed contract,
or runtime behavior change is planned.

Master review:
Approved. The smallest safe Phase 107 implementation is docs/site/contributor
alignment plus architecture overview cleanup, keeping public wording bounded.

Required edits:
Patch docs and contributor guidance, run claim audits and baseline checks, and
record validation results.

Decision:
Proceed to implementation checkpoint 000002 after plan validation and commit.

Next checkpoint:
Phase 107 -- Checkpoint 000002: implement primary scoped work.

## Checkpoint Report -- 000002

Checkpoint:
Phase 107 -- Checkpoint 000002: implement primary scoped work.

Sub-agents used or simulated, including intended model level:
Implementation, QA, UI/UX, Documentation / IA, Claim-Boundary, Security/Auth,
and Data/Migration roles are simulated by the Master Agent using direct
repository inspection because the attempted Context / Repo Truth and Planning
sub-agents timed out without edits. Master Agent -- GPT-5.5 x-high, current
thread.

Changed files:
`docs/architecture.md`; `docs/public-docs-site-freeze-checklist.md`;
`docs/contributor-first-issues.md`;
`docs/connectors/contributing-connectors.md`; `README.md`; `docs/README.md`;
`wiki/README.md`; `CONTRIBUTING.md`; `docs/integration-adapter-kit.md`;
`examples/README.md`; `testdata/README.md`;
`docs/phase-107-public-docs-site-freeze-and-contributor-onboarding.md`.

Validation run:
`git diff --check` passed. `make audit-product-acceptance` passed. `make
audit-final-claim-review` passed.

Blocked checks:
Full closeout baseline, `make check`, `make validate`, `make test`, Docker
Compose config, and final protected-path checks are deferred to checkpoint
000003/000004. Site publication, release actions, evidence collection,
external contact, real credentials, and consumer actions are out of scope.

Protected path status:
No protected evidence path was modified or required.

Consumer tracker status:
The consumer tracker was not modified. The seven targets remain required to
stay in order and `prepared`.

Claim-boundary status:
The implementation added docs/contributor guidance only and retained explicit
no-public-launch, no-adoption, no-consumer, no-compliance, no-final-root,
no-hosted-service, no-SLA, no-production-readiness, no-vendor, no-hardware, and
no-ETA-quality boundaries.

Security/auth status:
No runtime route, auth behavior, credential handling, raw private data,
external contact, site publication, or protected evidence path changed.

Data/migration status:
No migration, schema, durable state, dependency, public feed contract, or
runtime behavior changed.

Master review:
Approved. Phase 107 implementation remains docs-only and claim-bounded.

Required edits:
Run full focused validation and baseline checks. Patch only failures caused by
this phase.

Decision:
Proceed to validation/audit checkpoint 000003.

Next checkpoint:
Phase 107 -- Checkpoint 000003: run validation and patch required gaps.
