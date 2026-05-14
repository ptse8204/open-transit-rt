# Phase 95 -- v0.1.0-rc.1 Candidate Cut

## Scope

Phase 95 prepares and audits a local `v0.1.0-rc.1` candidate package using the
existing release tooling. The autonomous kickoff authorizes local `.cache`
package generation and audit for this phase only. It does not authorize a git
tag, GitHub Release creation, public image push, package-registry publication,
public announcement, retained evidence creation, consumer status movement, or
release-ready claim.

## Current Release Truth

- Phase 72 remains `needs_review`, not release-ready.
- Phase 89 remains the current local `v0.1.0-rc.1` gate result and closed as
  `needs_review`.
- Phase 92 ran a clean-checkout local release-candidate gate and also closed as
  `needs_review`.
- Phase 95 is the first post-90 phase authorized to run `make
  release-package` and `make audit-release-package`.
- Generated package outputs must remain under ignored `.cache` paths.

## Local Package Tooling

Existing release package tooling:

- `make release-package`
- `make audit-release-package`
- `make test-release-package`
- `scripts/release-package.sh`
- `scripts/audit-release-package.sh`

Expected local output:

```text
.cache/release-package/v0.1.0-rc.1/
```

Expected package contents:

- `artifacts/open-transit-rt-v0.1.0-rc.1.source.tar.gz`
- `checksums/SHA256SUMS.txt`
- `summary.json`
- `summary.md`
- `manifest.json`
- `manifest.md`
- `provenance.json`
- `provenance.md`
- `sbom.json`
- `image.json`

## Checkpoints

### Checkpoint 000001 -- Plan

Deliverables:

- Add this phase plan.
- Record authorization boundaries, local package path expectations, validation
  plan, and Master approval.

Validation:

- `git status --short`
- `git diff --check`
- `make check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- consumer tracker JSON parse and exact prepared-only assertion
- protected-path status check

### Checkpoint 000002 -- Refresh Candidate Docs

Deliverables:

- Refresh `docs/release-notes-v0.1.0-rc.1-draft.md` for Phase 95.
- Add draft-only tag command text.
- Add draft-only GitHub Release text.
- Keep all release action wording explicitly blocked.

Validation:

- `git diff --check`
- `make audit-final-claim-review`
- consumer tracker JSON parse and exact prepared-only assertion
- protected-path status check

### Checkpoint 000003 -- Generate And Audit Local Candidate Package

Deliverables:

- Run:

  ```bash
  RELEASE_PACKAGE_VERSION=v0.1.0-rc.1 RELEASE_PACKAGE_OUTPUT_DIR=.cache/release-package/v0.1.0-rc.1 RELEASE_PACKAGE_ALLOW_DIRTY=false RELEASE_PACKAGE_STRICT=true RELEASE_PACKAGE_FORCE=true make release-package
  RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1 make audit-release-package
  RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1 RUN_RELEASE_PACKAGE=true RUN_LOCAL_APP=true make release-candidate-check
  ```

- Record package path, source archive path, checksum manifest status,
  SBOM/provenance status, package audit result, and release-candidate diagnostic
  result.
- Keep generated artifacts in `.cache` only.

Validation:

- `make test-release-package`
- `make release-package`
- `make audit-release-package`
- `RUN_LOCAL_APP=true make release-candidate-check` with package audit enabled
- `git status --short -- .cache docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`
- baseline claim/protected-path checks

### Checkpoint 000004 -- Closeout

Deliverables:

- Add `docs/handoffs/phase-95.md`.
- Update `docs/current-status.md`, `docs/handoffs/latest.md`,
  `docs/roadmap-status.md`, and
  `docs/open-transit-rt-master-planner-remaining-work.md`.
- Record final candidate package/audit status, blockers, protected path status,
  consumer tracker status, claim-boundary status, security/auth status, and
  data/migration status.

Validation:

- `git status --short`
- `git diff --check`
- `make check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- consumer tracker JSON parse and exact prepared-only assertion
- protected-path status check
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`

## Draft-Only Release Action Text

Draft tag command, not to run in Phase 95:

```bash
git tag -a v0.1.0-rc.1 <commit> -m "Open Transit RT v0.1.0-rc.1"
```

Draft GitHub Release title, not to create in Phase 95:

```text
Open Transit RT v0.1.0-rc.1
```

Draft GitHub Release boundary sentence:

```text
This release candidate is for local evaluator review of the self-hosted Open
Transit RT codebase. It is not a hosted service, compliance certification,
consumer acceptance claim, agency endorsement, public launch claim, final-root
readiness claim, production readiness claim, release-ready claim, vendor
compatibility claim, hardware certification, SLA/uptime claim, or
production-grade ETA claim.
```

## Hard Boundaries

- Do not run `git tag`, `git push --tags`, `gh release create`, image push, or
  package-registry publication.
- Do not create or modify protected evidence paths.
- Do not edit or reformat
  `docs/evidence/consumer-submissions/status.json`.
- Do not move consumer statuses beyond `prepared`.
- Do not contact external services, consumers, agencies, vendors, or portals.
- Do not use real credentials, private payloads, real AVL data, or real agency
  data.
- Do not claim release readiness, compliance, adoption, consumer acceptance,
  final-root readiness, hosted-service availability, production readiness,
  vendor compatibility, hardware certification, SLA/uptime, or ETA quality.

## Master Approval

The Master Agent approves this Phase 95 plan. Implementation may proceed after
Checkpoint 000001 is committed. The package generation checkpoint must use
local ignored `.cache` output only and must record audit results without
publishing or tagging.

## Checkpoint 000001 Report

Checkpoint:
Phase 95 -- Checkpoint 000001: add v0.1.0-rc.1 candidate cut plan.

Sub-agents used or simulated, including intended model level:
Real Context / Repo Truth Sub-Agent -- GPT-5.5 x-high, Planning Sub-Agent --
GPT-5.5 x-high, Release / Supply-Chain Sub-Agent -- GPT-5.5 high, and
Claim-Boundary / Security QA Sub-Agent -- GPT-5.5 high were launched for Phase
95. Implementation, QA, and Documentation roles are simulated by the Master
Agent until the plan checkpoint is committed. Master Agent -- GPT-5.5 x-high,
current thread.

Changed files:
`docs/phase-95-v0-1-0-rc-1-candidate-cut.md`.

Validation run:
`git status --short`; `git diff --check`; `make check`; `make
audit-product-acceptance`; `make audit-final-claim-review`; consumer tracker
JSON parse; exact prepared-only consumer tracker assertion; protected-path
status check.

Blocked checks:
Package generation, package audit, local app release-candidate diagnostics,
and closeout validation are scheduled for later checkpoints.

Protected path status:
No protected evidence path is edited or generated.

Consumer tracker status:
`docs/evidence/consumer-submissions/status.json` is not edited. All seven
consumer targets must remain exactly `prepared`.

Claim-boundary status:
This plan authorizes local `.cache` package diagnostics only. It makes no
release readiness, compliance, adoption, consumer acceptance, production
readiness, final-root readiness, hosted-service availability, vendor
compatibility, hardware certification, SLA/uptime, or ETA-quality claim.

Security/auth status:
No route, auth, credential, CSRF, token, or protected data behavior changes are
planned. Release package tooling must not include secrets or private payloads.

Data/migration status:
No persistence, migration, GTFS data model, tenant model, or realtime data
model change is planned.

Master review:
Approved. The checkpoint sequence refreshes candidate docs first, cuts the
local source archive from a clean checkpoint, audits the package locally, and
keeps all release actions blocked.

Required edits:
None.

Decision:
Commit CP000001 and continue to candidate documentation refresh.

Next checkpoint:
Phase 95 -- Checkpoint 000002: refresh candidate docs.

## Checkpoint 000002 Candidate Documentation Refresh

Documentation refresh result:

- Refreshed `docs/release-notes-v0.1.0-rc.1-draft.md` for the Phase 95
  package-authorized checkpoint sequence.
- Added draft-only tag command text and draft-only GitHub Release body text.
- Recorded that Phase 95 package output remains local `.cache` diagnostics and
  not publication.
- Recorded a supply-chain publication blocker: the local source archive is
  generated from `git archive HEAD` and therefore includes tracked repository
  files; any public distribution requires a separate review for tracked
  evidence-path material.
- Tightened the package command to use `RELEASE_PACKAGE_ALLOW_DIRTY=false` and
  `RELEASE_PACKAGE_STRICT=true` for the local candidate cut.

## Checkpoint 000002 Report

Checkpoint:
Phase 95 -- Checkpoint 000002: refresh candidate docs.

Sub-agents used or simulated, including intended model level:
Real Release / Supply-Chain Sub-Agent -- GPT-5.5 high identified the stale
local package and source archive publication hazard. Real Planning Sub-Agent
-- GPT-5.5 x-high recommended clean strict package commands. Real Context /
Repo Truth Sub-Agent -- GPT-5.5 x-high and Claim-Boundary / Security QA
Sub-Agent -- GPT-5.5 high informed boundaries. Implementation and QA roles
were simulated by the Master Agent. Master Agent -- GPT-5.5 x-high, current
thread.

Changed files:
`docs/release-notes-v0.1.0-rc.1-draft.md`;
`docs/phase-95-v0-1-0-rc-1-candidate-cut.md`.

Validation run:
`git diff --check`; `make audit-final-claim-review`; `make
audit-product-acceptance`; consumer tracker JSON parse; exact prepared-only
consumer tracker assertion; protected-path status check.

Blocked checks:
Package generation, package audit, package-enabled release-candidate
diagnostics, and closeout validation are scheduled for later checkpoints.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched by
tracked changes.

Consumer tracker status:
`docs/evidence/consumer-submissions/status.json` was not edited. The exact
seven consumer targets remain required to stay `prepared`.

Claim-boundary status:
Draft tag and GitHub Release text are explicitly draft-only and blocked from
execution in Phase 95. The docs make no release readiness, compliance,
adoption, consumer acceptance, production readiness, final-root readiness,
hosted-service availability, vendor compatibility, hardware certification,
SLA/uptime, or ETA-quality claim.

Security/auth status:
No route, auth, credential, CSRF, token, or protected data behavior changed.
Release package outputs remain local-only `.cache` diagnostics.

Data/migration status:
No persistence, migration, GTFS data model, tenant model, or realtime data
model change is included.

Master review:
Approved. The release notes are current for the local package-authorized
checkpoint, and they explicitly block tag, release, public package, image,
evidence, and consumer-status actions.

Required edits:
None.

Decision:
Commit CP000002, then generate and audit the local candidate package from the
clean committed checkpoint.

Next checkpoint:
Phase 95 -- Checkpoint 000003: generate and audit local candidate package.
