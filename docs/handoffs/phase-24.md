# Phase Handoff

## Phase

Phase 24 — Real Agency Data Onboarding

## Status

- Complete for the docs/process and evidence-template scope.
- Active phase after this handoff: Phase 25 — Device And AVL Integration Kit is recommended next.

## What Was Implemented

- Added a real-agency GTFS onboarding guide for intake, approval, validation, publish review, redaction, and final public-feed review.
- Added a plain-language GTFS validation triage guide with common import/validation issues and when to ask for technical help.
- Added a template-only real-agency GTFS evidence scaffold.
- Updated tutorial, evidence, README, phase status, current status, and latest handoff docs to point to the new onboarding path.

## What Was Designed But Intentionally Not Implemented Yet

- No real agency GTFS was imported or committed.
- No real validation outputs, approvals, import evidence, screenshots, or placeholder artifacts were added.
- No backend import behavior, GTFS Studio behavior, Operations Console behavior, public feed URL, consumer status, or Phase 23 final-root status changed.
- No consumer submission, consumer acceptance, agency endorsement, CAL-ITP/Caltrans compliance, hosted SaaS, vendor-equivalence, or production-grade ETA claim was added.

## Onboarding Docs Added

- `docs/tutorials/real-agency-gtfs-onboarding.md` covers:
  - GTFS ZIP source and agency permission review
  - license/contact metadata review
  - agency identity and timezone review
  - routes, stops, trips, stop times, calendars, shapes, frequencies, blocks, and service date checks
  - import/publish path table for `cmd/gtfs-import`, GTFS Studio typed draft publish, validation review, public feed verification, and `/admin/operations`
  - publish review checklist
  - Phase 23-aware final public-feed review

## Validation Triage Docs Added

- `docs/tutorials/gtfs-validation-triage.md` explains:
  - missing required files
  - invalid route types
  - broken references
  - calendar and service-date problems
  - after-midnight times
  - shapes and `shape_dist_traveled`
  - blocks, frequencies, duplicates, timezone/date mistakes, and empty service
  - importer errors versus canonical validator errors
  - when to ask for technical help

## Publish Review And Metadata Approval Docs Added

- The onboarding guide and evidence template now include explicit approval fields for:
  - agency name
  - agency URL
  - timezone
  - technical contact email
  - license name and URL
  - public feed root
  - approver
  - approval date
  - notes
- `docs/tutorials/production-checklist.md` now points production-directed pilots to the real-agency onboarding approval and publish-review process.

## Real Data Redaction And Privacy Guidance

- The onboarding guide and evidence template prohibit private contracts, private contact information, private operator notes, private ticket links, non-public vehicle/device identifiers, raw private telemetry, credentials, tokens, private keys, and DB URLs with passwords.
- Raw validation outputs must be reviewed before commit. If they contain private paths, private contacts, operator notes, non-public data, credentials, or private agency material, keep raw output private and commit only a redacted summary.
- The evidence scaffold states that it contains templates only until real agency-approved, public-safe evidence exists.

## Evidence Templates Added

- Added `docs/evidence/real-agency-gtfs/README.md`.
- Added `docs/evidence/real-agency-gtfs/templates/import-review-template.md`.
- Linked the scaffold from `docs/evidence/README.md`.
- The scaffold explicitly forbids placeholder GTFS ZIPs, fake validation outputs, fake approvals, fake import results, screenshots, private agency data, and private operator artifacts.

## Schema And Interface Changes

- None.

## Dependency Changes

- None.

## Migrations Added

- None.

## Tests Added And Results

- No automated tests were added because Phase 24 is documentation and evidence-template work only.
- Existing validation and test commands are recorded below.

## Checks Run And Blocked Checks

- Pre-edit `make validate` — passed.
- Pre-edit `make test` — passed.
- Pre-edit `make realtime-quality` — passed.
- Pre-edit `make smoke` — passed.
- Pre-edit `docker compose -f deploy/docker-compose.yml config` — passed.
- Pre-edit `git diff --check` — passed.
- Post-edit `make validate` — passed.
- Post-edit `make test` — passed.
- Post-edit `make realtime-quality` — passed.
- Post-edit `make smoke` — passed.
- Post-edit `docker compose -f deploy/docker-compose.yml config` — passed.
- Post-edit `git diff --check` — passed.

Blocked or intentionally not run:

- `make test-integration` — not run because no GTFS import code, fixtures, demo flow, or local app behavior changed.
- `make demo-agency-flow` — not run because no GTFS import code, fixtures, demo flow, or local app behavior changed.
- `make agency-app-up` / `make agency-app-down` — not run because no local app behavior changed.
- `docker compose -f deploy/docker-compose.yml --profile app config` — not run because no local app behavior changed.
- `EVIDENCE_PACKET_DIR=<packet> make audit-hosted-evidence` — not run because no hosted or final-root evidence packet was created.

## Known Issues

- No real agency-approved GTFS source is present in repo evidence.
- No agency-owned or agency-approved final public feed root is available in repo evidence; Phase 23 remains blocker-documented only.
- Real-agency import evidence remains template-only.
- Consumer and aggregator targets remain `prepared` only with no submitted, under-review, accepted, rejected, blocked, ingestion, listing, or display evidence.
- The OCI pilot DuckDNS hostname remains hosted/operator pilot evidence only.

## Known Remaining Real-Agency Onboarding Gaps

- A future operator still needs a real public-safe GTFS source, permission/license note, metadata approval, validation outputs, import result, publish approval, public feed verification, and redaction notes.
- A future final-root review still needs agency-owned or agency-approved root evidence, DNS/TLS proof, public fetch proof, validator records, and packet refreshes before stronger agency-domain production claims.
- The current docs explain existing CLI/UI paths; they do not add a browser upload wizard or richer non-expert setup flow.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/handoffs/phase-24.md`
  - `docs/phase-25-device-avl-integration-kit.md`
  - `docs/track-b-productization-roadmap.md`
  - `docs/tutorials/real-agency-gtfs-onboarding.md`
  - `docs/evidence/redaction-policy.md`
  - `SECURITY.md`
- First files likely to edit:
  - `docs/phase-25-device-avl-integration-kit.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - future device/vendor onboarding docs selected by Phase 25
- Commands to run before coding:
  - `make validate`
  - `make test`
  - `make realtime-quality`
  - `make smoke`
  - `docker compose -f deploy/docker-compose.yml config`
  - `git diff --check`
- Known blockers:
  - No agency-owned final root exists in repo evidence.
  - No real agency-approved GTFS import packet exists yet.
  - Consumer statuses must remain `prepared` unless retained target-originated evidence supports a named target transition.
- Recommended first implementation slice:
  - Start Phase 25 with a docs-first real device/vendor telemetry intake kit that preserves token security, private telemetry redaction, and conservative realtime-quality claims.
