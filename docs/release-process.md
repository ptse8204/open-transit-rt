# Release Process

This document defines the expected release process for Open Transit RT. It is a maintainer process, not a promise of hosted service, paid support, SLA coverage, consumer acceptance, agency endorsement, or production readiness for every deployment.

## Versioning And Tags

Use git tags for release points. Until the project declares a stronger compatibility policy, prefer simple semantic version tags such as `v0.21.0`.

Patch tags should be used for narrowly scoped fixes. Minor tags may include docs, operations, packaging, or runtime improvements. Major tags should be reserved for future incompatible contracts after maintainers document those contracts.

## Who Can Cut A Release

Maintainers with repository write access may cut releases. A release should not be tagged until required checks pass or blocked checks are documented.

## Pre-Release Checklist

Run the relevant checks:

```bash
make validate
make test
git diff --check
```

Also run when relevant:

```bash
make realtime-quality
make smoke
docker compose -f deploy/docker-compose.yml config
```

For release candidates with deployment or evidence changes, review:

- `SECURITY.md`;
- `docs/evidence/redaction-policy.md`;
- `docs/prompts/calitp-truthfulness.md`;
- `docs/current-status.md`;
- `docs/handoffs/latest.md`;
- `docs/dependencies.md`;
- `docs/decisions.md`.

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

