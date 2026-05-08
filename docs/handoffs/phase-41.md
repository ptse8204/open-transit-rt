# Phase 41 Handoff

## Phase

Phase 41 — Operator Smoke And Support Bundle

## Status

Complete for the local/reference operator smoke and redaction-safe support
bundle scope.

Active phase after this handoff: Phase 42 — Reference Deployment Doctor is the
next proposed roadmap phase if the maintainer continues the
CAL-ITP-style/self-hosted gap-closure roadmap.

## What Was Implemented

- Added `scripts/operator-smoke.sh` and `make operator-smoke`.
- Added `scripts/support-bundle.sh` and `make support-bundle`.
- Added `docs/tutorials/operator-smoke-and-support-bundle.md`.
- Updated validation scaffolding so `make validate` checks the new scripts and
  docs.
- Updated the Phase 41 phase doc and navigation/status docs.

The operator smoke helper:

- writes under `.cache/operator-smoke/<timestamp>/` by default;
- validates `PUBLIC_BASE_URL` and safe `ADMIN_BASE_URL` rules;
- fetches the five public feed paths and records status, size, checksum, and
  redirect metadata;
- fails when public feeds are unavailable, empty, or exceed `MAX_FEED_BYTES`;
- checks unauthenticated `/admin/operations/readiness` on the public base URL
  and fails on any `2xx`;
- uses only `Authorization: Bearer "$ADMIN_TOKEN"` for authenticated admin
  requests, never cookie auth;
- checks authenticated `/admin/operations/readiness` only through a safe admin
  URL when `ADMIN_TOKEN` is supplied;
- always records validator tooling state through `scripts/check-validators.sh`;
- runs allowlisted validator API calls only when safe and authenticated;
- runs the deterministic synthetic AVL dry-run fixture;
- prints `external_evidence_created=false` and
  `consumer_statuses_changed=false`.

The support bundle helper:

- writes under `.cache/support-bundles/<timestamp>/` by default;
- records unavailable app, admin, database, and validator checks without
  failing the whole bundle;
- records validator tooling status without requiring pinned tooling to exist;
- includes authenticated readiness only when `ADMIN_TOKEN` is supplied and
  `INCLUDE_ADMIN_READINESS=true`;
- stores summaries, not full validation reports or raw feed bodies;
- sanitizes migration-status output when `DATABASE_URL` is supplied;
- includes manifest files for included and excluded diagnostics;
- runs a final redaction scan for secret-shaped values.

## What Was Designed But Intentionally Not Implemented Yet

- No deployment doctor/preflight workflow; that is Phase 42.
- No production evidence packet generation.
- No final-root evidence workflow.
- No consumer submission automation or status update.
- No real vendor AVL adapter send mode.
- No runtime external predictor integration.

## Schema And Interface Changes

- No database schema changes.
- No HTTP API changes.
- No public feed URL or GTFS-RT protobuf contract changes.
- New local operator command interfaces:
  - `make operator-smoke`
  - `make support-bundle`

## Dependency Changes

- No new external dependency was added.
- The new scripts use existing repo tools and common local tools: `sh`, `curl`,
  `python3`, `go`, git, and optional validator/Docker tooling where already
  documented.

## Migrations Added

- None.

## Tests Added And Results

- Added `make validate` checks for:
  - `scripts/operator-smoke.sh`
  - `scripts/support-bundle.sh`
  - script POSIX syntax via `sh -n`
  - script help paths
  - Phase 41 tutorial and handoff files

## Checks Run And Blocked Checks

- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed after fixing one trailing-space issue.
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` — passed.
- Read-only exact seven-target consumer status check — passed; Google Maps,
  Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and
  transit.land all remain `prepared`.
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json` —
  passed; the tracker file is unchanged.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- `make smoke` — passed.
- `make support-bundle` — passed. The bundle recorded public feed checks as
  unavailable because no local/reference app was running, recorded validator
  tooling as passed, skipped authenticated readiness by default, ran the
  synthetic AVL dry-run, and passed the generated-file redaction scan.

Blocked:

- `make operator-smoke SKIP_VALIDATORS=true` — not run because no
  local/reference app was running at `http://localhost:8080`; direct probe of
  `/public/feeds.json` returned curl connection failure / HTTP `000`. Operator
  smoke is intentionally strict and should be run after `make agency-app-up` or
  another local/reference deployment is actually serving the public feed paths.

## Known Issues

- No agency-owned or agency-approved final-root evidence exists.
- Consumer and aggregator targets remain `prepared` only.
- Support bundles and smoke output are private diagnostics by default, not
  evidence packets.

## Exact Next-Step Recommendation

- First files to read:
  - `docs/handoffs/latest.md`
  - `docs/roadmap-to-calitp-compliance-and-gap-closure.md`
  - `docs/phase-41-operator-smoke-support-bundle.md`
  - `docs/tutorials/operator-smoke-and-support-bundle.md`
- First files likely to edit:
  - `scripts/deployment-doctor.sh`
  - `Makefile`
  - `docs/deployment/reference-deployment-doctor.md`
- Commands to run before coding:
  - `make validate`
  - `make test`
  - `git diff --check`
- Known blockers:
  - None for local/reference diagnostic tooling.
  - Final-root, real agency adoption, consumer acceptance, and compliance
    claims still require retained claim-specific evidence.
- Recommended first implementation slice:
  - Start Phase 42 with a read-only deployment doctor that checks environment
    presence, service health, DB migration status, validator tooling, public
    feed edge, and admin boundary safety without creating evidence.
