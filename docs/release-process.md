# Release Process

This document defines the expected release process for Open Transit RT. It is a maintainer process, not a promise of hosted service, paid support, SLA coverage, consumer acceptance, agency endorsement, or production readiness for every deployment.

## Versioning And Tags

Use git tags for release points. Until the project declares a stronger compatibility policy, prefer simple semantic version tags such as `v0.22.0`.

Patch tags should be used for narrowly scoped fixes. Minor tags may include docs, operations, packaging, or runtime improvements. Major tags should be reserved for future incompatible contracts after maintainers document those contracts.

Releases are cut from `main` using tags. Release branches are not used unless a later release process explicitly documents them.

Use these commands to identify the deployed source version:

```bash
git describe --tags --always --dirty
git rev-parse HEAD
```

Evidence packets and release notes should record the git tag, commit SHA, dirty/clean state, release notes link, and artifact checksums where available.

## Distribution Artifacts

The current supported distribution anchors are:

- source tags;
- exact commit SHAs;
- local Docker image builds from a checked-out tag;
- checksums for any generated release artifact.

Published/versioned production Docker images are deferred. Do not claim a production image exists unless a future release adds and documents one.

Operators can pin a local Docker image to the release tag:

```bash
docker build -f deploy/Dockerfile.local -t open-transit-rt-local:v0.22.0 .
```

## Who Can Cut A Release

Maintainers with repository write access may cut releases. A release should not be tagged until required checks pass or blocked checks are documented.

## Pre-Release Checklist

Run and record all required final checks:

```bash
make validate
make test
make realtime-quality
make smoke
docker compose -f deploy/docker-compose.yml config
git diff --check
```

For release candidates with deployment or evidence changes, review:

- `SECURITY.md`;
- `docs/evidence/redaction-policy.md`;
- `docs/prompts/calitp-truthfulness.md`;
- `docs/current-status.md`;
- `docs/handoffs/latest.md`;
- `docs/dependencies.md`;
- `docs/decisions.md`.

Use `docs/release-checklist.md` for the full release procedure.

## Release Notes

Release notes should list:

- user-facing changes;
- migrations;
- operations changes;
- security notes;
- evidence or claim changes;
- known limitations;
- required checks and their result.

If no migrations, security notes, or evidence/claim changes exist, say so explicitly.

Use `docs/release-notes-template.md` for the release note shape.

## Install, Upgrade, And Rollback

Use `docs/upgrade-and-rollback.md` for operator-facing install, upgrade, migration, and restore guidance.

Release notes must say whether a release requires migrations. Operators must back up before upgrading, check migration status before and after running migrations, and treat database restore as the recovery path when migrations are irreversible or untested.

## Evidence And Claims

Release notes must not convert implementation, validation, replay, or pilot evidence into unsupported claims.

Do not claim:

- CAL-ITP/Caltrans compliance;
- consumer submission, review, ingestion, or acceptance;
- agency endorsement;
- hosted SaaS availability;
- marketplace/vendor equivalence;
- paid support or SLA coverage;
- universal production readiness;
- production-grade ETA quality.

Any release note that changes consumer or evidence status must point to retained, redacted, target-originated evidence.

## After Tagging

After tagging, update status or handoff docs when a phase closes. Keep `docs/current-status.md` and `docs/handoffs/latest.md` aligned with the actual release scope.
