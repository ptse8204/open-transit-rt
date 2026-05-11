# Phase 66 -- Release Candidate And Installability

## Status

In progress. Checkpoint 000001 added this plan for release-candidate workflow,
installer/bootstrap UX, Docker image publishing decision, demo/docs website
planning, and closeout. Checkpoint 000002 prepared the first release-candidate
workflow by adding the ordered review sequence, validation matrix, release-note
inputs, package audit matrix, and private diagnostic summary fields. Phase 66
must keep release candidate and installability work bounded to evaluator
workflows and local/self-hosted packaging. It must not create retained
evidence, publish artifacts, push images, change consumer statuses, or claim
hosted SaaS, production readiness, SLA/uptime, agency adoption, consumer
acceptance, public launch, vendor compatibility, hardware certification,
production AVL reliability, production-grade ETA quality, or CAL-ITP/Caltrans
compliance.

## Goal

Make Open Transit RT easier to install, evaluate, package, and review from a
clean checkout. A maintainer or small-agency evaluator should be able to
understand what to install, which checks to run, what a blocker means, and what
release-candidate artifacts are local diagnostics rather than proof of
production readiness.

## Checkpoints

- Completed: `Phase 66 -- Checkpoint 000001: add release candidate and installability plan`
- Completed: `Phase 66 -- Checkpoint 000002: prepare first release candidate workflow`
- Planned: `Phase 66 -- Checkpoint 000003: improve installer and bootstrap UX`
- Planned: `Phase 66 -- Checkpoint 000004: document Docker image publishing decision`
- Planned: `Phase 66 -- Checkpoint 000005: add demo site or documentation website plan`
- Planned: `Phase 66 -- Checkpoint 000006: close release candidate and installability`

## Existing State

- `make check`, `make test`, `make validate`, `make release-candidate-check`,
  `make release-package`, `make audit-release-package`,
  `make test-release-package`, `make external-connection-check`,
  `make adapter-conformance`, and `make test-connector-examples` already
  exist.
- `scripts/release-candidate-check.sh` writes private local diagnostics under
  `.cache/release-candidate-check/` and explicitly avoids tagging, publishing,
  pushing images, evidence writes, consumer status changes, and stronger
  claims.
- `scripts/release-package.sh`, `scripts/audit-release-package.sh`, and
  `scripts/test-release-package.sh` already support ignored local source
  package diagnostics.
- `scripts/bootstrap-dev.sh`, `make bootstrap`, `make agency-app-up`, and
  `docs/tutorials/local-quickstart.md` provide local setup paths, but missing
  tool messages and first-run guidance can be clearer.
- `docs/release-candidate-readiness.md`, `docs/release-checklist.md`,
  `docs/release-process.md`, `docs/release-notes-template.md`, and
  `docs/tutorials/deploy-with-docker-compose.md` already establish release and
  local Docker guidance.
- Current Docker guidance supports source tags, local Docker builds, and
  self-hosted evaluation. Published/versioned production Docker images remain
  deferred unless a later explicit decision changes that.

## Checkpoint Scope

### Release Candidate Workflow

- Improve the release-candidate review path so a clean checkout can produce a
  bounded local readiness summary.
- Clarify the order of commands, expected outputs, missing-tool blockers, and
  release notes inputs.
- Keep generated `.cache` diagnostics private and local.
- Do not tag, publish, upload, push images, create GitHub releases, or claim
  production readiness.

### Installer And Bootstrap UX

- Improve bootstrap/preflight messages for missing tools such as Docker, Go,
  Make, Python, Git, curl, Java, and pinned validator tooling.
- Prefer actionable next steps over raw command failure text.
- Keep `bootstrap-dev.sh` and Makefile behavior compatible with existing local
  evaluator workflows.
- Do not change migrations, database schema, telemetry ingest, public feed
  routes, or auth boundaries.

### Docker Image Publishing Decision

- Document whether Phase 66 keeps Docker image work source/local-only or opens
  a future image publication track.
- If image publication remains deferred, state why and what future
  authorization/release process would be required.
- Do not push images, publish registry tags, or imply hosted SaaS or production
  support.

### Demo Site Or Documentation Website Plan

- Add a public-friendly docs/demo-site plan that explains how evaluators can
  learn the product without treating the plan as a public launch.
- Prefer repository documentation structure over a heavy frontend or marketing
  site unless a later phase explicitly authorizes one.
- Keep the plan claim-bounded: no public launch completion, agency adoption,
  consumer acceptance, hosted SaaS, or production readiness.

### Closeout

- Mark Phase 66 complete only after implementation checkpoints pass validation
  and claim-boundary review.
- Record protected-path, consumer tracker, and unsupported-claim results.
- Point the next active phase to Phase 67.

## Files Expected To Change

- `Makefile`
- `scripts/release-candidate-check.sh`
- `scripts/bootstrap-dev.sh`
- optional new preflight helper under `scripts/`
- `docs/release-candidate-readiness.md`
- `docs/release-checklist.md`
- `docs/release-process.md`
- `docs/release-notes-template.md`
- `docs/tutorials/local-quickstart.md`
- `docs/tutorials/deploy-with-docker-compose.md`
- `docs/decisions.md`
- optional new docs/demo or website plan under `docs/`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/roadmap-status.md`
- this phase file
- closeout handoff `docs/handoffs/phase-66.md`

Protected paths remain untouched:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/**`
- `db/migrations/**`
- `go.mod`
- `go.sum`

## Non-Goals

- No release tagging, GitHub release creation, registry publication, package
  repository upload, or image push.
- No hosted SaaS, paid support, SLA/uptime, public launch, production
  readiness, agency adoption, consumer acceptance, vendor compatibility,
  hardware certification, production AVL reliability, production-grade ETA
  quality, or CAL-ITP/Caltrans compliance claim.
- No retained evidence writes.
- No consumer status changes.
- No external contact.
- No database migrations.
- No public feed URL changes.
- No telemetry ingest contract changes.
- No GTFS-RT protobuf semantic changes.
- No validator execution semantic changes unless a checkpoint narrowly
  approves missing-tool messaging around existing checks.
- No connector manifest schema changes.
- No prediction adapter behavior changes.
- No auth boundary or public/private route boundary weakening.

## Validation Plan

- `git diff --check`
- `go test ./cmd/agency-config`
- `go test ./cmd/agency-config -run TestReleaseCandidateCheck`
- `make check`
- `make test`
- `make release-candidate-check` when environment blockers are not expected,
  or `scripts/release-candidate-check.sh --dry-run` for bounded local review
- `make test-release-package`
- `make external-connection-check`
- `make adapter-conformance`
- `make test-connector-examples`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json`
- `git diff --exit-code -- docs/evidence/captured`
- `git diff --exit-code -- db/migrations go.mod go.sum`
- `docker compose -f deploy/docker-compose.yml config`

If Java, Docker, network, validator tooling, or pinned images are unavailable,
record the exact blocker and continue with non-environment-dependent checks.

## Claim Boundary

Phase 66 may say the project is easier to install, evaluate, package locally,
and review as a release candidate. It may say the release-candidate workflow
helps maintainers identify blockers.

Phase 66 must not say or imply that a release candidate is production-ready,
hosted, SLA-backed, agency-approved, consumer-accepted, public-launch complete,
vendor-compatible, hardware-certified, production AVL reliable,
production-grade ETA proven, or CAL-ITP/Caltrans compliant.

Release-candidate checks, release packages, Docker Compose configuration,
local Docker builds, and docs/demo plans are private or local evaluation
signals only. They are not retained external evidence and must not move
consumer tracker records beyond `prepared`.

## Rollback Path

Phase 66 should remain scripts/docs/private diagnostic work. If rollback is
needed, revert the specific checkpoint commit that added the workflow, helper,
decision, docs-site plan, test, or documentation change. Public feed URLs, DB
schema, telemetry ingest, GTFS-RT semantics, validator execution semantics,
connector manifest schema, prediction adapter behavior, auth boundaries,
protected evidence paths, and consumer tracker statuses should remain
untouched.
