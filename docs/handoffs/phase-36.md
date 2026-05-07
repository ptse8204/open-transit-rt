# Phase 36 Handoff

## Phase

Phase 36 — OCI/OCL Reference Deployment Productization

## Status

- Phase 36 is complete for OCI/OCL reference deployment productization.
- Active phase after this handoff: Phase 37 — Reusable Agency Onboarding Flow.

Phase 36 is docs-only. No runtime behavior changed. No external evidence was
created. No consumer statuses changed.

## What Was Implemented

- Added `docs/deployment/README.md` as the deployment index for the
  self-hosted OCI/OCL-style reference path.
- Added `docs/deployment/oci-reference-deployment.md` with the operator
  copy/paste path, server prerequisites, user/group layout, directory layout,
  placeholder-only env guidance, secret generation, database setup/migrations,
  systemd supervision, Caddy or equivalent reverse proxy guidance,
  public/private route boundaries, validator install/check, GTFS import, five
  public feed URL verification, backup/restore, feed monitor, scorecard export,
  update/rollback, and redacted support bundle guidance.
- Added `docs/deployment/oci-reference-env.example` with grouped placeholder
  sections and generated-secret placeholders only.
- Added `docs/deployment/oci-reference-smoke-checklist.md` for repeatable
  operator verification without converting smoke output into evidence.
- Updated navigation, status, roadmap, backlog, open-question, and latest
  handoff docs to close Phase 36 and make Phase 37 the recommended next phase.
- Updated `docs/phase-36-oci-reference-deployment-productization.md` from a
  pre-work phase stub into a closed phase reference.

## What Was Designed But Intentionally Not Implemented Yet

- No managed-hosting, SaaS, Kubernetes, or container-platform packaging path
  was added.
- No agency onboarding wizard or reusable intake flow was added; that remains
  the Phase 37 recommendation.
- No new evidence packet, final-root proof, consumer submission, or external
  acceptance workflow was created.

## Schema And Interface Changes

- No schemas changed.
- No migrations changed.
- No APIs changed.
- No public feed contracts changed.
- No runtime behavior changed unless explicitly documented; this phase
  documented only the reference deployment path.

## Dependency Changes

- No dependency files changed.
- Existing PostgreSQL/PostGIS, Caddy or equivalent proxy, systemd, Java, Docker
  or equivalent validator runtime, and MobilityData validator boundaries remain
  as previously documented.

## Migrations Added

- None.

## Tests Added And Results

- No tests were added because Phase 36 was documentation-only.
- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed.
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` — passed.
- Read-only consumer status check — passed; 7 targets found and all remain
  `prepared`.

## Checks Run And Blocked Checks

- No blocked checks.

## Known Issues

- External-proof gaps remain unresolved: no agency-owned final-root evidence,
  no target-originated consumer evidence, no real agency adoption evidence, no
  real vendor AVL evidence, and no production-grade ETA evidence.
- The reference deployment docs are operator guidance only; they do not prove
  that a deployment was run, accepted, compliant, agency-approved, or
  production-ready.

## Exact Next-Step Recommendation

Proceed to Phase 37 — Reusable Agency Onboarding Flow.

Use the Phase 36 reference deployment docs as the self-hosted server target for
the onboarding flow. Keep external-proof tracks as future optional paths unless
retained, redacted, claim-specific evidence exists.
