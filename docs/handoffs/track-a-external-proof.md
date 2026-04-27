# Track A Handoff

## Phase

Track A — External Proof And Adoption

## Status

- Complete for the approved docs-only operator workflow, evidence intake, artifact-directory, and agency-owned domain readiness scope.
- Active phase after this handoff: none selected.

## What Was Implemented

- Added `docs/evidence/consumer-submissions/submission-workflow.md` for official submission-path verification, pre-submission checks, evidence intake, status transitions, record updates, and claim boundaries.
- Added README-only target artifact intake directories under `docs/evidence/consumer-submissions/artifacts/` for Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land.
- Added `docs/agency-owned-domain-readiness.md` for moving from the DuckDNS OCI pilot to agency-owned production-domain proof.
- Updated evidence, readiness, roadmap, docs index, current status, and latest handoff docs to point operators to the new workflow.
- Kept all seven consumer and aggregator targets at `prepared` only.

## What Was Designed But Intentionally Not Implemented Yet

- No external portal was contacted.
- No submission was automated.
- No target submission path was marked verified.
- No backend behavior, helper script, external integration, public feed URL, GTFS-RT behavior, prediction logic, or consumer status was changed.
- No placeholder screenshots, fake receipts, fake tickets, fake emails, example correspondence, or private target artifacts were created.

## Schema And Interface Changes

- No database schema, runtime API, public feed URL, GTFS-RT protobuf contract, Trip Updates adapter, unauthenticated surface, or backend behavior changed.
- Added docs/evidence-only workflow and artifact README files.

## Dependency Changes

- None.

## Migrations Added

- None.

## Tests Added And Results

- Added no code tests.
- Verified `docs/evidence/consumer-submissions/status.json` remains valid JSON.
- Verified the human-readable tracker and `status.json` agree for target name, status, packet path, prepared timestamp, and evidence references.
- Verified per-target artifact directories contain README files only.
- Ran a context-aware forbidden-claim scan that allows definitions, warnings, negated statements, and transition rules, while checking that no current target status or current claim advanced beyond `prepared`.
- Ran a targeted redaction scan over new/edited evidence docs.

## Checks Run And Blocked Checks

- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed.
- `make realtime-quality` — passed.
- `make smoke` — passed.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json` — passed.
- Tracker/status consistency check — passed.
- Artifact README-only check — passed.
- Context-aware forbidden-claim scan — passed.
- Targeted redaction scan — passed.

Blocked checks:

- None.

## Known Issues

- All consumer and aggregator targets are still `prepared` only; no target has retained third-party submission, review, acceptance, rejection, or blocker evidence.
- Official submission paths remain unverified because no current official target source or operator-retained target evidence was added.
- The OCI pilot DuckDNS hostname remains pilot evidence, not agency-owned stable URL/domain proof.
- Agency-owned final-root validator records, refreshed packets, and final-root consumer submissions remain missing.
- Phase 19 replay metrics remain measurement evidence only and do not prove production-grade ETA quality.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/handoffs/track-a-external-proof.md`
  - `docs/evidence/consumer-submissions/submission-workflow.md`
  - `docs/evidence/consumer-submissions/README.md`
  - `docs/evidence/consumer-submissions/status.json`
  - `docs/agency-owned-domain-readiness.md`
- First files likely to edit:
  - `docs/evidence/consumer-submissions/current/<target>.md` only after real redacted target-originated evidence exists
  - `docs/evidence/consumer-submissions/artifacts/<target>/` only after real redacted target-originated evidence exists
  - `docs/evidence/consumer-submissions/README.md`
  - `docs/evidence/consumer-submissions/status.json`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
- Commands to run before coding:
  - `make validate`
  - `make test`
  - `git diff --check`
  - Run `make realtime-quality`, `make smoke`, and `docker compose -f deploy/docker-compose.yml config` if evidence/readiness docs change materially.
- Known blockers:
  - Consumer status cannot move beyond `prepared` without retained redacted target-originated evidence from the named consumer or aggregator.
  - Agency-owned stable URL/domain proof remains missing.
  - Official target submission paths cannot be marked verified without current official source evidence or operator-retained target evidence.
- Recommended first implementation slice:
  - If a real operator is ready, verify one target's official submission path using `docs/evidence/consumer-submissions/submission-workflow.md`, retain redacted evidence, and update only that target.
