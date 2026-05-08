# Phase 37 Handoff

## Phase

Phase 37 — Reusable Agency Onboarding Flow

## Status

- Complete for the reusable local/reference onboarding scope.
- Active phase after this handoff: Phase 38 — Integration Adapter Kit.

Phase 37 added an opt-in onboarding command. It did not create final-root
evidence, external evidence packets, consumer artifacts, agency approval or
adoption claims, consumer acceptance claims, CAL-ITP/Caltrans compliance
claims, production-readiness claims, vendor-compatibility claims, or
production-grade ETA-quality claims.

## What Was Implemented

- Added `scripts/agency-pilot-onboard.sh`.
- Added `make agency-pilot-up AGENCY_ID=... GTFS_URL=...`.
- Updated `deploy/docker-compose.yml` so local app services can use explicit
  env interpolation for agency/public feed settings while preserving default
  `demo-agency` behavior.
- Added `docs/tutorials/reusable-agency-onboarding.md`.
- Updated README, deployment, tutorial, roadmap, backlog, current-status,
  open-question, phase-reference, and latest-handoff docs.
- Post-completion hardening patched running-mode behavior so it safely upserts
  the requested agency/admin roles through `DATABASE_URL`, requires an explicit
  running-mode `ADMIN_BASE_URL`, rejects `.`/`..`/leading-dot agency IDs, adds
  a no-network dry-run check to `make validate`, and passes additional
  advanced options through `make agency-pilot-up`.

The onboarding script:

- validates agency ID, URL, mode, timeout, and metadata inputs;
- supports `--dry-run` without Docker, network, database, validators, or
  secrets;
- downloads GTFS into ignored `.cache/agency-pilot/<agency-id>/` storage;
- records source checksum;
- seeds the requested agency/admin roles without using `scripts/seed-dev.sql`;
- imports the requested GTFS with configurable timeout;
- bootstraps explicit publication metadata from CLI/env values or obvious
  local/reference placeholders;
- starts local Compose services directly without calling `make agency-app-up`;
- fetches all five public paths;
- verifies the fetched schedule summary against the imported source summary;
- runs validators, skips validators, or records blockers.

## What Was Designed But Intentionally Not Implemented Yet

- No browser wizard for GTFS ZIP/URL onboarding was added.
- No production multi-tenant packaging was added.
- No agency-owned final-root evidence workflow was run.
- No consumer submission or consumer-status workflow was run.
- No device credential creation for arbitrary agencies was added by default.

## Schema And Interface Changes

- No database schema changed.
- No migrations changed.
- No public feed URL contracts changed.
- No API contracts changed.
- Runtime behavior changed only through the new opt-in onboarding command and
  Compose env interpolation defaults.

## Dependency Changes

- No new external dependency was added.
- The script uses existing local/operator tools already documented for this
  repo: Docker Compose, `curl`, `python3`, `unzip`, and pinned validators when
  available.

## Migrations Added

- None.

## Tests Added And Results

- `make validate` now checks that `scripts/agency-pilot-onboard.sh` exists,
  parses with `sh -n`, that `--help` exits successfully without requiring
  Docker, network, database, validators, or secrets, and that a no-network
  `--dry-run` invocation succeeds.
- Optional local fixture smoke passed using a generated ZIP from
  `testdata/gtfs/valid-small` served from ignored `.cache/` storage at
  `http://127.0.0.1:18080/valid-small.zip`. The run used `--skip-validators`,
  imported `demo-agency`, fetched all five public paths, and passed the
  schedule identity check. Local Compose services were stopped afterward
  without deleting the database volume.

## Checks Run And Blocked Checks

- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed.
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` — passed.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- Read-only consumer status check — passed; seven targets remain `prepared`.
- Optional local fixture smoke — passed with validators skipped.

No external evidence was created unless a future handoff explicitly documents
otherwise.

## Known Issues

- Publication metadata is placeholder metadata unless operators supply
  agency-approved `TECHNICAL_CONTACT_EMAIL`, `FEED_LICENSE_NAME`, and
  `FEED_LICENSE_URL` values.
- Validator tooling remains optional for onboarding unless
  `--strict-validators` is supplied.
- Running mode requires `--admin-base-url` or `ADMIN_BASE_URL`; use a loopback,
  VPN, SSH tunnel, or otherwise private/admin-protected URL.
- Agency-owned final-root proof, consumer acceptance, production readiness,
  vendor compatibility, and production-grade ETA quality remain unproven.

## Exact Next-Step Recommendation

- First files to read:
  - `docs/handoffs/latest.md`
  - `docs/current-status.md`
  - `docs/phase-38-integration-adapter-kit.md`
  - `docs/dependencies.md`
- First files likely to edit:
  - `docs/phase-38-integration-adapter-kit.md`
  - `docs/dependencies.md`
  - adapter-related tutorials under `docs/tutorials/`
- Commands to run before coding:
  - `make validate`
  - `make test`
  - `git diff --check`
- Known blockers:
  - no real vendor AVL evidence;
  - no runtime external predictor integration;
  - no consumer target-originated evidence;
  - no agency-owned final-root evidence.
- Recommended first implementation slice:
  - define the reusable integration adapter kit scope without adding named
    vendor or consumer acceptance claims.
