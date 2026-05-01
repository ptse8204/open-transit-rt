# Phase 22 — Release And Distribution Hardening

## Status

Planned Track B phase. Not implemented until `docs/handoffs/latest.md` marks it active.

## Purpose

Make Open Transit RT easier to version, release, install, upgrade, and roll back.

The repo is now capable and documented, but agencies and contributors need a stable release process. A release should explain what changed, what checks passed, whether migrations are required, and how an operator can recover if an upgrade fails.

## Scope

1. Versioning and tagging model.
2. Changelog and release notes.
3. Release checklist.
4. Install/upgrade/rollback documentation.
5. Docker image / artifact distribution plan.
6. Migration compatibility guidance.
7. Release validation commands.
8. Release evidence and claim-change tracking.

## Required Work

### 1) Versioning And Tags

Document:

- tag naming convention;
- semantic-versioning-like expectations, if used;
- how to identify the deployed version;
- how release branches or tags relate to `main`;
- how evidence packets reference a version.

### 2) Changelog And Release Notes

Add or improve:

- `CHANGELOG.md`;
- release note template;
- release checklist.

Release notes should include:

- user-facing changes;
- migrations;
- operations changes;
- security notes;
- dependency changes;
- evidence/claim changes;
- known limitations;
- required checks.

### 3) Install, Upgrade, Rollback

Document:

- clean install path;
- upgrade path;
- migration run order;
- backup-before-upgrade rule;
- rollback limits;
- restore procedure linkage;
- local app version verification.

### 4) Distribution Artifacts

Define the current artifact strategy:

- source tag;
- local Docker build;
- optional versioned Docker images if implemented later;
- checksums where useful;
- how operators should pin versions.

Do not imply published production images exist unless they do.

### 5) Release Validation

Document required pre-release checks and evidence-safe review:

```bash
make validate
make test
make realtime-quality
make smoke
docker compose -f deploy/docker-compose.yml config
git diff --check
```

Add additional checks only when needed.

## Acceptance Criteria

Phase 22 is complete only when:

- maintainers have a repeatable release checklist;
- changelog/release note structure exists;
- version/tag guidance exists;
- install/upgrade/rollback docs exist;
- migration/evidence/security release notes are covered;
- release docs do not claim compliance, consumer acceptance, hosted SaaS, or vendor equivalence.

## Required Checks

```bash
make validate
make test
git diff --check
```

If deployment/release docs change materially, also run:

```bash
make realtime-quality
make smoke
docker compose -f deploy/docker-compose.yml config
```

## Explicit Non-Goals

Phase 22 does not:

- change backend runtime behavior;
- add migrations unless explicitly justified;
- change public feed URLs;
- publish production Docker images unless separately approved;
- claim production readiness for all agencies;
- claim consumer acceptance or compliance;
- add external integrations.

## Likely Files

- `CHANGELOG.md`
- `docs/release-process.md`
- `docs/release-checklist.md`
- `docs/upgrade-and-rollback.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-22.md`
- `README.md` only for a short release/install link if needed
