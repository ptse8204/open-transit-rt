# Phase Handoff

## Phase

Phase 25 — Device And AVL Integration Kit

## Status

- Complete for the docs/process and template-only evidence scope.
- Active phase after this handoff: Phase 26 — Admin UX Setup Wizard is recommended next.

## What Was Implemented

- Added a telemetry API and AVL integration guide for authenticated device, vendor, adapter, and simulator telemetry onboarding.
- Added a device token lifecycle guide for bearer credential handling, rotate/rebind behavior, one-time token display, secure storage, compromise response, and binding rules.
- Added vendor AVL adapter boundary guidance without naming or integrating any vendor.
- Added simulator/no-hardware testing guidance, clock/timezone/GPS quality guidance, and explicit troubleshooting guidance.
- Added template-only device/AVL evidence scaffolding.
- Updated tutorial, evidence, README, phase status, current status, decisions, and latest handoff docs to point to the new Phase 25 materials.

## What Was Designed But Intentionally Not Implemented Yet

- No backend API behavior, protobuf contract, prediction logic, public feed URL, consumer status, dependency, or runtime integration changed.
- No named external vendor, proprietary AVL integration, certified hardware path, or production AVL reliability claim was added.
- No real agency device data, vendor payload, private telemetry, credential, fake evidence, hardware certification, consumer acceptance, agency endorsement, hosted SaaS, marketplace-equivalence, or CAL-ITP/Caltrans compliance claim was added.
- No changes were made to `scripts/device-onboarding.sh`; the existing helper already provides help, sample, simulate, dry-run, and rebind flows.

## Telemetry API Docs Added

- Added `docs/tutorials/device-avl-integration.md`.
- The guide documents:
  - `POST /v1/telemetry`
  - `Authorization: Bearer <device-token>`
  - required fields: `agency_id`, `device_id`, `vehicle_id`, `timestamp`, `lat`, `lon`
  - optional fields: `driver_id`, `bearing`, `speed_mps`, `accuracy_m`, `trip_hint`
  - RFC3339 timestamp requirements
  - WGS84 latitude/longitude expectations
  - bearing, speed, accuracy, and trip-hint guidance
  - example curl and JSON payloads using only synthetic demo identifiers
  - verification through response JSON, `/admin/operations/telemetry`, protected `/v1/events`, and Vehicle Positions output
  - what the integration proves and does not prove
- Response examples are limited to behavior confirmed from `cmd/telemetry-ingest` and tests. Error bodies are documented as plain troubleshooting text, not a versioned JSON API contract.

## Token Lifecycle Docs Added

- Added `docs/tutorials/device-token-lifecycle.md`.
- The guide documents:
  - device tokens as bearer credentials
  - agency/device/vehicle binding checks
  - local seeded demo token behavior
  - `/admin/devices/rebind`, `/admin/operations/devices`, and `scripts/device-onboarding.sh rebind`
  - one-time token display
  - immediate old-token invalidation after rotation/rebind
  - secure storage expectations
  - rotation after suspected compromise
  - operator responsibilities and audit logging
  - what never to log or commit

## Vendor AVL Boundary Guidance

- Vendor payloads should be transformed into Open Transit RT telemetry events before forwarding to `/v1/telemetry`.
- Vendor credentials, private payloads, private device identifiers, private vehicle identifiers, and private logs stay outside the public repo unless reviewed and explicitly approved as public-safe.
- Vendor-specific assumptions must not be embedded into core matching, Vehicle Positions generation, or Trip Updates prediction.
- Acceptable integration shapes include agency-owned adapter scripts, deployment-owned sidecars, vendor-owned middleware, and private operator integration processes.
- `docs/decisions.md` records this as ADR-0026.
- `docs/dependencies.md` was intentionally not changed because no named external vendor, adapter implementation, or dependency status was introduced.

## Simulator And Testing Guidance

- `docs/tutorials/device-avl-integration.md` documents:
  - `scripts/device-onboarding.sh sample --dry-run`
  - `scripts/device-onboarding.sh simulate --dry-run`
  - no-hardware demo limits
  - how to view telemetry freshness in `/admin/operations/telemetry`
  - how to view public Vehicle Positions effects
  - why simulator success is not production AVL proof

## Troubleshooting And Redaction Guidance

- The integration guide includes a troubleshooting table with symptom, likely cause, how to check, next action, and what not to claim yet for:
  - bad token or missing `Authorization`
  - wrong agency/device/vehicle
  - timestamp too old
  - timestamp in the future
  - invalid lat/lon
  - low GPS accuracy
  - telemetry accepted but no assignment
  - Vehicle Positions stale or missing
  - Trip Updates withheld
  - simulator works but real hardware is still unproven
  - validator passes but consumer acceptance is not proven
- The docs explicitly state that `/v1/events` is an authenticated admin/debug review path, not a public feed or consumer-facing endpoint.
- Redaction guidance prohibits committing device tokens, admin tokens, JWT/CSRF secrets, DB URLs with passwords, vendor credentials, private AVL payloads, raw private telemetry, private IDs, private operator notes, private logs with credentials, and `.cache` files.

## Script Changes

- None.
- Existing helper behavior was documented and verified with syntax/help/dry-run checks.

## Evidence Scaffold Added

- Added `docs/evidence/device-avl/README.md`.
- Added `docs/evidence/device-avl/templates/integration-review-template.md`.
- The scaffold is template-only and explicitly forbids fake telemetry evidence, fake vendor approvals, fake hardware certifications, private AVL payloads, credentials, and raw private telemetry.

## Schema And Interface Changes

- None.

## Dependency Changes

- None.
- `docs/dependencies.md` was left unchanged.

## Migrations Added

- None.

## Tests Added And Results

- No automated tests were added because Phase 25 is documentation and template-only evidence work.
- Existing validation and test commands are recorded below.

## Checks Run And Blocked Checks

- Pre-edit/planning `make validate` — passed.
- Pre-edit/planning `make test` — passed.
- Pre-edit/planning `make realtime-quality` — passed.
- Pre-edit/planning `make smoke` — passed.
- Pre-edit/planning `docker compose -f deploy/docker-compose.yml config` — passed.
- Pre-edit/planning `git diff --check` — passed.
- Pre-edit/planning `sh -n scripts/device-onboarding.sh` — passed.
- Pre-edit/planning `scripts/device-onboarding.sh help` — passed.
- Pre-edit/planning `scripts/device-onboarding.sh sample --dry-run` — passed.
- Pre-edit/planning `scripts/device-onboarding.sh simulate --dry-run` — passed.
- Post-edit `make validate` — passed.
- Post-edit `make test` — passed.
- Post-edit `make realtime-quality` — passed.
- Post-edit `make smoke` — passed.
- Post-edit `docker compose -f deploy/docker-compose.yml config` — passed.
- Post-edit `git diff --check` — passed.
- Post-edit `sh -n scripts/device-onboarding.sh` — passed.
- Post-edit `scripts/device-onboarding.sh help` — passed.
- Post-edit `scripts/device-onboarding.sh sample --dry-run` — passed.
- Post-edit `scripts/device-onboarding.sh simulate --dry-run` — passed.
- Post-edit targeted docs secret/example scan — passed.

Blocked or intentionally not run:

- `make demo-agency-flow` — not run because no telemetry scripts, local app helpers, demo flow, backend behavior, fixtures, or local app behavior changed.
- `make agency-app-up` / `make agency-app-down` — not run because no local app behavior changed.
- `docker compose -f deploy/docker-compose.yml --profile app config` — not run because no local app behavior changed.
- `make test-integration` — not run because no telemetry ingest code, schema, fixtures, or integration behavior changed.
- `EVIDENCE_PACKET_DIR=<packet> make audit-hosted-evidence` — not run because no hosted or final-root evidence packet was created.

## Known Issues

- No real device or vendor AVL integration evidence exists in the repo.
- No certified hardware or certified vendor support evidence exists.
- Simulator/no-hardware examples remain local demo evidence only.
- No agency-owned or agency-approved final public feed root is available in repo evidence; Phase 23 remains blocker-documented only.
- Consumer and aggregator targets remain `prepared` only with no submitted, under-review, accepted, rejected, blocked, ingestion, listing, or display evidence.
- Production-grade ETA quality remains unproven beyond existing replay measurement.

## Known Remaining Device/AVL Integration Gaps

- A future operator still needs real public-safe device or vendor evidence, source/permission notes, redaction review, token rotation records, telemetry freshness evidence, Vehicle Positions review, and clear simulator/pilot/real-device labeling.
- A future vendor-specific integration, if added, must define an adapter outside core matching and update dependency docs only when a named dependency or integration status actually changes.
- The current docs explain existing API and helper paths; they do not add a browser setup wizard, hosted device management service, hardware certification path, or richer non-expert onboarding UI.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/handoffs/phase-25.md`
  - `docs/phase-26-admin-ux-setup-wizard.md`
  - `docs/tutorials/device-avl-integration.md`
  - `docs/tutorials/device-token-lifecycle.md`
  - `docs/evidence/redaction-policy.md`
  - `SECURITY.md`
- First files likely to edit:
  - `docs/phase-26-admin-ux-setup-wizard.md`
  - `cmd/agency-config/`
  - `cmd/gtfs-studio/`
  - `cmd/feed-alerts/`
  - `docs/tutorials/agency-first-run.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/handoffs/phase-26.md`
- Commands to run before coding:
  - `make validate`
  - `make test`
  - `make smoke`
  - `make demo-agency-flow`
  - `git diff --check`
- Known blockers:
  - No agency-owned final root exists in repo evidence.
  - No real public-safe device or vendor AVL evidence exists yet.
  - Consumer statuses must remain `prepared` unless retained target-originated evidence supports a named target transition.
- Recommended first implementation slice:
  - Start Phase 26 with a browser-guided setup checklist in the Operations Console that links existing GTFS, publication metadata, validation, device token, telemetry freshness, and evidence workflows without changing public feed URLs, protobuf contracts, consumer statuses, or unsupported readiness claims.

