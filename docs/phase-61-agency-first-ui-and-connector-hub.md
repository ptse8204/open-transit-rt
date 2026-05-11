# Phase 61 -- Agency-First UI And Connector Hub

## Status

Complete.

Phase 61 adopted the maintainer-approved Phase 61+ roadmap naming and added a
private Connector Hub to the Operations Console. The phase did not reopen
Phases 0 through 60, write retained evidence, contact external parties, change
consumer statuses, or add compliance, adoption, hosted-service, vendor,
hardware, public-launch, SLA, production-readiness, or production-grade ETA
claims.

## Checkpoints

- `Phase 61 -- Checkpoint 000001: add agency-first connector platform roadmap`
- `Phase 61 -- Checkpoint 000002: implement agency-first UI and connector hub`
- `Phase 61 -- Checkpoint 000003: close agency-first UI and connector hub`

## Implementation Summary

- Added the canonical roadmap pack under
  `docs/roadmaps/agency-first-connector-platform/`.
- Updated source-of-truth docs so future product work uses Phase 61+ naming
  instead of stale "not Phase 61" wording.
- Improved the private Operations Console dashboard with agency-first primary
  action cards.
- Added read-only private routes:
  - `/admin/operations/connectors`
  - `/admin/operations/connectors.json`
- Added Connector Hub categories for:
  - telemetry / GPS / AVL source connectors;
  - prediction engine connectors;
  - validator connectors;
  - monitoring / export connectors;
  - consumer / discovery workflows.
- Recorded the safe plugin definition:
  "In Open Transit RT, a plugin is an optional sidecar, command adapter,
  manifest, or connector process. It is not arbitrary dynamic code loaded into
  the backend."

## Claim Boundary

Connector Hub is private authenticated operator guidance only. It creates no
evidence, contacts no external party, changes no consumer status, runs no
external connector, and records no approval, compatibility, compliance,
hosted-service, SLA, public-launch, or production-readiness outcome.

All connector paths remain optional sidecar, manifest, command-adapter, or
deployment-owned process paths. Phase 61 added no dynamic backend plugin
loading and no vendor compatibility claim.

## Verification

Verification was run from `/Users/edwintse/Downloads/open-transit-rt`.

- `git diff --check`
- `go test ./cmd/agency-config`
- `make check`
- `make audit-final-claim-review`
- `make external-connection-check`
- `make adapter-conformance`
- `make test-connector-examples`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json`
- `git diff --exit-code -- docs/evidence/captured`

## Next Phase

Phase 62 -- Guided Setup and Browser GTFS Import.
