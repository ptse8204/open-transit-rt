# Phase 35 Handoff

## Phase

Phase 35 — README And Roadmap Realignment

## Status

- Complete for docs-only README and roadmap realignment.
- Active phase after this handoff: Phase 36 — OCI/OCL Reference Deployment
  Productization.

Phase 35 restored the root README as the Open Transit RT product front door and
made self-hosted agency reuse the default roadmap. It did not create new
evidence or strengthen any public claim.

## What Was Implemented

- Replaced root `README.md` roadmap-export wording with a product front door
  covering:
  - what Open Transit RT is;
  - who it is for;
  - current GTFS, telemetry, matching, GTFS-RT, validation, local app, and
    pilot-ops capabilities;
  - three paths: local evaluation, real public GTFS local/pilot run, and the
    OCI/OCL-style reference deployment path;
  - CAL-ITP-style readiness pointers;
  - integration boundaries for AVL/device adapters, external predictors,
    validators, monitoring, and consumers;
  - clear claim boundaries.
- Updated `docs/current-status.md`, `docs/handoffs/latest.md`,
  `docs/future-roadmap-post-outcome-c.md`, `docs/roadmap-status.md`, and
  `docs/track-b-productization-roadmap.md` so the default next work is
  self-hosted agency reuse and OCI/OCL reference deployment productization.
- Updated `docs/phase-plan.md` so Phase 34 records both the original static
  GTFS validator missing-Java blocker and the later Homebrew Java 17 retry:
  exit code `0`, system error count `0`, and 3 warning notices. The retry is
  not described as validator-clean, no-warning, compliance, final-root,
  consumer, or production evidence.
- Updated `docs/README.md` and `docs/tutorials/README.md` with findable links
  for the self-hosted master plan, Phase 35, Phase 36, public GTFS onboarding,
  pilot operations, and readiness docs.
- Updated `docs/backlog.md` to record Phase 35 completion and Phase 36 as the
  next backlog direction.
- Updated `docs/open-questions.md` to record that the next deployment target is
  the OCI/OCL-style single-server reference path; managed container and
  Kubernetes paths remain later options.

Changed files:

- `README.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-35.md`
- `docs/future-roadmap-post-outcome-c.md`
- `docs/roadmap-status.md`
- `docs/track-b-productization-roadmap.md`
- `docs/phase-plan.md`
- `docs/README.md`
- `docs/tutorials/README.md`
- `docs/backlog.md`
- `docs/open-questions.md`

## What Was Intentionally Not Changed

- No runtime code changed.
- No schemas, migrations, APIs, or public feed contracts changed.
- No consumer statuses changed.
- No final-root evidence changed.
- No external evidence packets changed.
- No validator artifacts changed.
- No external-proof docs were deleted; they are reframed as future optional
  proof tracks.

## Evidence And Claim Boundary

- All seven consumer and aggregator targets remain `prepared` only.
- The OCI DuckDNS pilot remains hosted/operator pilot evidence only, not
  agency-owned final-root proof.
- No CAL-ITP/Caltrans compliance claim was added.
- No consumer acceptance, agency endorsement/adoption, production readiness,
  hosted SaaS, vendor compatibility, or production-grade ETA claim was added.
- Phase 35 created no external evidence.

## Checks Run

- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed.
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` — passed.
- Read-only consumer tracker status check — passed; 7 targets found and all
  remain `prepared`.

## Known Issues

- Phase 36 reference deployment docs still need to be productized from the
  existing small-agency pilot operations path.
- External-proof gaps remain unresolved: no agency-owned final-root evidence,
  no target-originated consumer evidence, no real agency adoption evidence, no
  real vendor AVL evidence, and no production-grade ETA evidence.

## Exact Next-Step Recommendation

Proceed to Phase 36 — OCI/OCL Reference Deployment Productization.

Use the existing OCI/OCL-style pilot server pattern as the reference deployment
path for self-hosted agency reuse. Keep external-proof tracks as future optional
paths unless retained, redacted, claim-specific evidence exists.
