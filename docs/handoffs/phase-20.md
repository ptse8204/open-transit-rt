# Phase Handoff Template

## Phase

Phase 20 — Consumer Submission Execution And CAL-ITP Readiness Program

## Status

- Complete for the approved docs/evidence packet-preparation and readiness-summary scope.
- Active phase after this handoff: Phase 21 — Community, Governance, And Multi-Agency Scale, if the roadmap still applies.

## What Was Implemented

- Added complete prepared packet drafts for Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land under `docs/evidence/consumer-submissions/packets/`.
- Added packet freshness fields: prepared at, prepared by, evidence snapshot, OCI packet reference, `feeds.json` snapshot reference, validator records reference, and Phase 19 replay/quality summary reference.
- Added submission-method fields to every packet and marked submission method/contact as `not verified` instead of guessing.
- Added operator warnings to every packet: do not submit from repo docs alone; review feed URLs, license/contact metadata, validation status, agency identity, consumer-specific requirements, and redactions before actual submission.
- Added packet completeness checklist table at `docs/evidence/consumer-submissions/packets/README.md`.
- Added machine-readable tracker snapshot at `docs/evidence/consumer-submissions/status.json`.
- Updated all seven current target records and the human-readable tracker to `prepared` only because complete packets exist.
- Added `docs/california-readiness-summary.md`, including the missing agency-owned stable URL/domain proof gap.
- Added `docs/marketplace-vendor-gap-review.md`.
- Updated status, evidence, compliance, phase, and docs index pages to reflect Phase 20 truthfully.

## What Was Designed But Intentionally Not Implemented Yet

- No external portal was contacted.
- No submission was automated.
- No submission path was guessed.
- No backend API behavior was added.
- No consumer submission, under-review, acceptance, rejection, blocker, ingestion, CAL-ITP/Caltrans compliance, marketplace/vendor equivalence, hosted SaaS availability, agency endorsement, or production-grade ETA claim was added.
- Official consumer submission paths remain unverified and must be verified by the operator outside the repo before any real submission.

## Schema And Interface Changes

- No database schema, runtime API, public feed URL, GTFS-RT protobuf contract, Trip Updates adapter, unauthenticated surface, or backend behavior changed.
- Added docs/evidence-only `status.json`; it is not a backend API.

## Dependency Changes

- None.

## Migrations Added

- None.

## Tests Added And Results

- Added no code tests.
- Added docs/evidence consistency artifacts and ran local consistency checks:
  - `status.json` agrees with the human-readable tracker for target name, status, packet path, prepared timestamp, and evidence reference values.
  - every packet includes all completeness-required fields.
  - `status.json` parses as valid JSON.

## Checks Run And Blocked Checks

Pre-edit baseline:

- `make validate` — passed.
- `make realtime-quality` — passed.
- `make test` — passed.
- `make smoke` — passed.
- `make test-integration` — initially blocked because Postgres was not listening on local port `55432`.
- `make db-up` — passed and started local Postgres.
- `make test-integration` — passed after `make db-up`.
- `git diff --check` — passed.
- `docker compose -f deploy/docker-compose.yml config` — passed.

Docs consistency checks:

- `python3` tracker/status consistency check — passed.
- `python3` packet completeness check — passed.

Final post-edit checks:

- `make validate` — passed.
- `make realtime-quality` — passed.
- `make test` — passed.
- `make smoke` — passed.
- `make test-integration` — passed.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json` — passed.
- targeted redaction scan over new/edited evidence docs — passed.
- `git diff --check` — passed.

Blocked checks:

- None after starting local Postgres with `make db-up`.

## Known Issues

- All consumer and aggregator targets are `prepared` only; no target has retained third-party submission, review, acceptance, rejection, or blocker evidence.
- Official submission URL/contact values are intentionally `not verified`.
- The OCI pilot DuckDNS hostname is useful pilot evidence but is not agency-owned stable URL/domain proof.
- Phase 19 replay metrics are measurement evidence only and do not prove production-grade ETA quality.
- Marketplace/vendor-equivalent support, SLA/SLO, procurement, onboarding, hardware, incident response, training, and managed-service commitments remain gaps.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/handoffs/phase-20.md`
  - `docs/evidence/consumer-submissions/README.md`
  - `docs/evidence/consumer-submissions/status.json`
  - `docs/california-readiness-summary.md`
  - `docs/marketplace-vendor-gap-review.md`
- First files likely to edit:
  - `docs/evidence/consumer-submissions/current/<target>.md` only after real redacted third-party evidence exists
  - `docs/evidence/consumer-submissions/status.json`
  - `docs/evidence/consumer-submissions/README.md`
  - `docs/handoffs/latest.md`
  - `docs/current-status.md`
- Commands to run before coding:
  - `make validate`
  - `make test`
  - `git diff --check`
  - Run `make realtime-quality`, `make smoke`, `make test-integration`, and `docker compose -f deploy/docker-compose.yml config` if evidence/readiness docs or app/demo docs change materially.
- Known blockers:
  - Consumer status cannot move beyond `prepared` without retained redacted evidence from the named consumer or aggregator.
  - Agency-owned stable URL/domain proof remains missing for stronger California readiness language.
- Recommended first implementation slice:
  - Phase 21 should focus on community/governance/multi-agency readiness without changing consumer evidence claims. If an operator obtains real consumer artifacts first, update the specific target record and `status.json` only to the evidence-backed status.
