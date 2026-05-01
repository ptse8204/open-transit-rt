# Phase Handoff

## Phase

Phase 22 — Release And Distribution Hardening

## Status

- Complete
- Active phase after this handoff: Phase 23 — Agency-Owned Deployment Proof is recommended next.

## What Was Implemented

- Added a changelog, release checklist, release notes template, and upgrade/rollback guide.
- Expanded the release process with tag, version verification, source artifact, local Docker build, release note, install, upgrade, rollback, and evidence version-linkage guidance.
- Added release and upgrade links to `docs/README.md` and a short release/install link group to `README.md`.
- Updated `docs/current-status.md`, `docs/handoffs/latest.md`, and `docs/phase-22-release-distribution-hardening.md`.

## What Was Designed But Intentionally Not Implemented Yet

- No backend behavior, `/version` endpoint, binary `--version`, OCI image labels, database migrations, public feed URL changes, consumer status changes, external integrations, compliance claims, consumer acceptance claims, hosted-SaaS claims, or vendor-equivalence claims were added.
- Published/versioned production Docker images are deferred. Current distribution guidance supports source tags and local Docker builds only.
- Release branches are deferred. Releases are cut from `main` using tags unless a later process explicitly changes that rule.

## Release/Distribution Docs Added

- `CHANGELOG.md`
- `docs/release-checklist.md`
- `docs/upgrade-and-rollback.md`
- `docs/release-notes-template.md`
- Expanded `docs/release-process.md`

## Versioning/Tagging Guidance

- Use semantic-version-style git tags such as `v0.22.0`.
- Cut releases from `main` using tags.
- Identify a deployed checkout with `git describe --tags --always --dirty` and `git rev-parse HEAD`.
- Pin installs by source tag, exact commit SHA, local Docker image tag, and release artifact checksum when generated.
- Future evidence packets should record git tag, commit SHA, dirty/clean state, release notes link, and artifact checksums where available.

## Install/Upgrade/Rollback Guidance

- Clean install guidance starts from a source tag and verifies the local app package with `make agency-app-up` plus public feed fetches.
- Local Docker image guidance builds from the checked-out tag with `deploy/Dockerfile.local`.
- Upgrade guidance requires backup before upgrade, migration status checks before and after `make migrate-up`, service restart, and install verification.
- Rollback guidance warns that irreversible or untested migrations require database restore plus redeploying the previous tag/artifacts; `make migrate-down` is not treated as the recovery plan.
- Restore procedure links point to `docs/runbooks/small-agency-pilot-operations.md`, `scripts/pilot-ops.sh restore-drill --dry-run`, `docs/evidence/redaction-policy.md`, and `SECURITY.md`.

## Schema And Interface Changes

- None.

## Dependency Changes

- None.

## Migrations Added

- None.

## Tests Added And Results

- No automated tests were added because Phase 22 is docs/process-only.

## Checks Run And Blocked Checks

- Pre-edit `make validate` — passed.
- Pre-edit `make test` — passed.
- Pre-edit `git diff --check` — passed.
- Post-edit `make validate` — passed.
- Post-edit `make test` — passed.
- Post-edit `make realtime-quality` — passed.
- Post-edit `make smoke` — passed.
- Post-edit `docker compose -f deploy/docker-compose.yml config` — passed.
- Post-edit `git diff --check` — passed.
- Blocked commands: none.

## Known Issues

- Published/versioned production Docker images are not available.
- Runtime version reporting is still manual through git tag/SHA records; no `/version` endpoint, binary `--version`, or image labels exist.
- Release artifact checksum guidance exists, but no automated release artifact generator was added.
- Evidence packet version-linkage guidance exists, but existing captured evidence packets were not retroactively rewritten.

## Exact Next-Step Recommendation

- First files to read: `AGENTS.md`, `docs/current-status.md`, `docs/handoffs/latest.md`, `docs/track-b-productization-roadmap.md`, `docs/phase-23-agency-owned-deployment-proof.md`, `docs/agency-owned-domain-readiness.md`, `docs/evidence/redaction-policy.md`, and `SECURITY.md`.
- First files likely to edit: `docs/phase-23-agency-owned-deployment-proof.md`, `docs/handoffs/latest.md`, and `docs/current-status.md`.
- Commands to run before coding: `make validate`, `make test`, `make realtime-quality`, `make smoke`, `docker compose -f deploy/docker-compose.yml config`, and `git diff --check`.
- Known blockers: agency-owned stable URL/domain proof is still missing; consumer targets remain `prepared` only.
- Recommended first implementation slice: document and gather agency-owned or agency-approved domain deployment proof without changing public feed paths or consumer statuses.
