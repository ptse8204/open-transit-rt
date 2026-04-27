# Phase Handoff Template

## Phase

Phase 21 — Community, Governance, And Multi-Agency Scale

## Status

- Complete for the approved docs/process/governance/teaching-visual scope.
- Active phase after this handoff: next roadmap phase should be selected from `docs/roadmap-post-phase-14.md` or a new maintainer-approved plan.

## What Was Implemented

- Added contributor, code of conduct, issue template, PR template, governance, release process, support-boundary, multi-agency strategy, and roadmap/status docs.
- Added safe-reporting warnings for tokens, DB URLs, private keys, admin URLs with secrets, private portal screenshots, private ticket links, raw logs with credentials, and unredacted operator artifacts.
- Added contributor start points for docs improvements, tutorial fixes, replay fixtures, bug reproduction, tests, operator runbooks, and issue triage.
- Added explicit governance authority for who can merge PRs, cut releases, approve docs/evidence wording, and resolve competing design decisions.
- Added release-note guidance for user-facing changes, migrations, operations changes, security notes, evidence/claim changes, known limitations, and required checks.
- Added teaching visuals:
  - `docs/assets/how-to-contribute-paths.png` used in `CONTRIBUTING.md` and `docs/README.md`.
  - `docs/assets/community-workflow.png` used in `CONTRIBUTING.md` and `docs/governance.md`.
  - `docs/assets/single-vs-multi-agency.png` used in `docs/multi-agency-strategy.md`.
  - `docs/assets/evidence-maturity-ladder.png` used in `docs/roadmap-status.md` and `docs/compliance-evidence-checklist.md`.
  - `docs/assets/support-boundaries.png` used in `docs/support-boundaries.md` and `wiki/support-and-contribute.md`.
- Updated `docs/assets/README.md` with filename, purpose, usage, alt text, generation method, prompt/spec, and truthfulness notes for the new images.

## What Was Designed But Intentionally Not Implemented Yet

- No backend feature changes.
- No database schema changes.
- No runtime API or public feed URL changes.
- No consumer submission automation.
- No external integrations.
- No GitHub label creation automation.
- No legal foundation, paid support model, SLA process, hosted SaaS service, procurement package, or formal foundation/company governance was added.
- No consumer status moved beyond `prepared`.

## Schema And Interface Changes

- No database schema, runtime API, public feed URL, GTFS-RT protobuf contract, Trip Updates adapter, unauthenticated surface, or backend behavior changed.
- Added repository process files and documentation assets only.

## Dependency Changes

- None.

## Migrations Added

- None.

## Tests Added And Results

- No Go tests were added because the change set is Markdown, GitHub templates, process docs, and PNG documentation assets.
- Manual documentation verification included:
  - generated image assets exist under `docs/assets/`;
  - consuming Markdown pages include descriptive alt text;
  - `docs/assets/README.md` documents each new image with filename, purpose, usage, alt text, generation method, prompt/spec, and truthfulness note.
- Issue template YAML parse check passed.

## Checks Run And Blocked Checks

- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed.
- `make realtime-quality` — passed.
- `make smoke` — passed.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- visual documentation verification script — passed.
- issue template YAML parse script — passed.

Blocked commands:

- None currently known.

## Known Issues

- GitHub labels referenced in issue templates may need to be created or colorized manually in GitHub.
- The project still has no paid support, SLA coverage, formal legal foundation, or hosted SaaS offering.
- True multi-tenant hosting remains future work; current docs describe single-agency/pilot support plus agency-scoped foundations.
- Consumer/aggregator targets remain `prepared` only; no target has submitted, under-review, accepted, rejected, or blocked evidence.
- OCI pilot evidence remains pilot evidence, not agency-owned production proof.
- Production-grade ETA quality remains unproven beyond current replay metrics.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `CONTRIBUTING.md`
  - `docs/governance.md`
  - `docs/roadmap-status.md`
  - `docs/multi-agency-strategy.md`
- First files likely to edit:
  - GitHub labels in repository settings, if maintainers want labels to match templates.
  - `docs/roadmap-post-phase-14.md` if maintainers define a post-Phase-21 sequence.
  - `docs/evidence/consumer-submissions/current/<target>.md` only after real target-originated evidence exists.
- Commands to run before coding:
  - `make validate`
  - `make test`
  - `git diff --check`
  - `make realtime-quality` when realtime/readiness docs change materially.
  - `make smoke` when validation, hardening, support, or operations docs change materially.
  - `docker compose -f deploy/docker-compose.yml config` when deployment assumptions change.
- Known blockers:
  - Consumer status cannot move beyond `prepared` without retained redacted evidence from the named consumer or aggregator.
  - Agency-owned stable URL/domain proof remains missing.
  - True multi-tenant hosting requires additional implementation and isolation proof.
- Recommended first implementation slice:
  - Choose the next roadmap phase explicitly. A practical next step is a maintainer-approved post-Phase-21 roadmap update that decides whether to prioritize multi-agency isolation tests, production operations hardening, or real consumer-submission evidence intake after target-originated artifacts exist.
