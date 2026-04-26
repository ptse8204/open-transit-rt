# Phase 13 Handoff

All future phase handoff files must use this structure unless the phase explicitly documents a reason to diverge.

## Phase

Phase 13 — Consumer Submission and Acceptance Evidence

## Status

- Complete for the initial docs/evidence tracker structure
- Active phase after this handoff: operator collection of real external evidence, or the next documented phase only after evidence-bounded review

## What Was Implemented

- Added `docs/consumer-submission-evidence.md` with status definitions, allowed claims by status, required evidence fields, acceptance-scope fields, and operator update process.
- Added `docs/evidence/consumer-submissions/README.md` with tracker freshness fields, reviewed-by field, linked Phase 12 evidence packet, current target table, status definitions, and OCI pilot feed URL references for future submission packets.
- Added current records for Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land under `docs/evidence/consumer-submissions/current/`.
- Added reusable templates for all seven targets under `docs/evidence/consumer-submissions/templates/`.
- Updated `docs/compliance-evidence-checklist.md`, `docs/current-status.md`, `docs/handoffs/latest.md`, `docs/dependencies.md`, and `docs/phase-13-consumer-submission-evidence.md`.

## What Was Designed But Intentionally Not Implemented Yet

- No consumer submission API integrations.
- No automated submissions.
- No external consumer or aggregator acceptance claims.
- No CAL-ITP compliance claim.
- No backend/runtime/product changes.
- No public feed URL changes.

## Schema And Interface Changes

- None.

## Dependency Changes

- No runtime dependency changes.
- `docs/dependencies.md` now notes that Phase 13 consumer targets are documentation-only evidence records/templates and remain workflow metadata, not external API integrations.

## Migrations Added

- None.

## Tests Added And Results

- No tests added because this was a documentation/evidence-only change.

## Checks Run And Blocked Checks

- Pre-edit checks passed:
  - `make validators-check`
  - `make validate`
  - `make test`
  - `make smoke`
  - `make demo-agency-flow`
  - `make test-integration`
  - `docker compose -f deploy/docker-compose.yml config`
  - `git diff --check`
- Post-edit checks passed:
  - `make validate`
  - `make test`
  - `git diff --check`

## Known Issues

- No real redacted third-party submission, review, acceptance, rejection, or blocker evidence is present in the repository.
- All current target records are therefore `not_started`.
- Validator success and public fetch proof remain supporting evidence only and cannot be used as consumer acceptance.

## Exact Next-Step Recommendation

- First files to read: `docs/consumer-submission-evidence.md`, `docs/evidence/consumer-submissions/README.md`, `docs/prompts/calitp-truthfulness.md`, `docs/evidence/captured/oci-pilot/2026-04-24/README.md`
- First files likely to edit: `docs/evidence/consumer-submissions/current/<target>.md`, `docs/evidence/consumer-submissions/README.md`, `docs/current-status.md`, `docs/handoffs/latest.md`
- Commands to run before coding: `make validators-check`, `make validate`, `make test`, `make smoke`, `make demo-agency-flow`, `make test-integration`, `docker compose -f deploy/docker-compose.yml config`, `git diff --check`
- Known blockers: consumer status cannot move beyond `not_started` without real redacted external artifacts from the named consumer or aggregator workflow
- Recommended first implementation slice: when an operator has a real receipt, ticket, portal screenshot, correspondence, rejection, blocker, or acceptance artifact, redact it, add it to an evidence packet, update exactly one target current record, update tracker freshness fields, and keep allowed public wording scoped to that target and evidence date
