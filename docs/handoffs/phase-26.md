# Phase

Phase 26 — Admin UX Setup Wizard

## Status

- Complete for the Phase 26 browser-guided setup checklist scope.
- Active phase after this handoff: Phase 27 — Multi-Agency Isolation Prototype.

## What Was Implemented

- Expanded `/admin/operations/setup` from a minimal checklist into a server-rendered guided setup page with status, status source, evidence signal, and next action for common agency setup tasks.
- Added setup sections for agency metadata, license/contact metadata, GTFS import/GTFS Studio path, publication bootstrap, device token setup, first telemetry event, first validation run, public feed verification, Alerts setup, consumer packet review, and evidence/readiness review.
- Added an admin-only setup publication metadata form that reuses the existing `BootstrapPublication` repository/API semantics and derives `agency_id` from the authenticated principal.
- Added an admin-only setup validation form that accepts only feed type from the browser and maps feed type to an allowlisted server-side validator ID.
- Added publication config read support so the setup page can show current public base URL, feed base URL, technical contact, license, publication environment, and read-only agency ID.
- Improved the consumer Operations Console view so Phase 20 docs/evidence tracker `prepared` packet state is shown separately from runtime DB consumer workflow records.
- Added evidence links for the OCI pilot packet, Phase 23 agency-owned-domain blocker, real-agency GTFS evidence scaffold, device/AVL evidence scaffold, consumer tracker/status JSON, California readiness summary, compliance checklist, runbook, and redaction policy.

## What Was Designed But Intentionally Not Implemented Yet

- Browser GTFS ZIP upload remains deferred because this phase did not add upload size limits, storage handling, validation gating, or a full security review for that surface.
- Manual assignment override/review UI remains deferred because a safe browser view needs carefully bounded summaries and must avoid raw diagnostics, private notes, and new mutation semantics.
- Consumer submission automation remains out of scope.
- No public feed URL, GTFS-RT protobuf, telemetry API, device API, Trip Updates adapter, external integration, consumer status, or evidence-claim changes were made.

## Schema And Interface Changes

- Added `compliance.PublicationConfig` as a safe read model for current publication metadata.
- Added `PostgresRepository.PublicationConfig(ctx, agencyID)` to read existing `feed_config` fields for authenticated admin UI display.
- Added setup-only form handling under `/admin/operations/setup`; no public route or public feed behavior changed.
- Existing `/admin/publication/bootstrap` JSON semantics were not changed.
- Existing `/admin/validation/run` JSON semantics were not changed; the browser setup form uses a narrower server-side mapping path.

## Dependency Changes

- None.

## Migrations Added

- None.

## Setup Wizard/Checklist Routes And Pages Added

- `GET /admin/operations/setup`: authenticated read-only/operator/editor/admin setup checklist and guidance page.
- `POST /admin/operations/setup` with `action=publication_bootstrap`: admin-only publication metadata bootstrap/update form using existing repository semantics.
- `POST /admin/operations/setup` with `action=run_validation`: admin-only validation form; browser supplies only `feed_type`.

## Metadata/Publication Behavior

- The setup page renders current publication metadata from `feed_config` where available.
- Agency ID is displayed as read-only authenticated-principal context and is not trusted from the browser.
- Form fields are trimmed and length-bounded before calling `BootstrapPublication`.
- Missing license, contact, feed base URL, or public base URL fields remain visibly missing.
- Demo/local metadata is not described as agency-approved.

## GTFS Import UX Decision

- Browser ZIP upload is deferred.
- The setup page links operators to CLI GTFS ZIP import guidance, GTFS Studio, validation triage, and active feed verification instead.

## Validation UX Behavior

- The setup page shows feed validation state from feed discovery validation records.
- Browser validation accepts only feed type.
- Server mapping is:
  - `schedule` -> `static-mobilitydata`
  - `vehicle_positions`, `trip_updates`, `alerts` -> `realtime-mobilitydata`
- Validation wording is bounded to supporting evidence only and does not imply consumer acceptance, ingestion, or compliance.

## Device/Telemetry Setup Behavior

- The setup page shows whether device bindings exist from the device binding store.
- It shows whether latest telemetry exists and whether stale latest rows exist from the telemetry repository summary.
- It links to `/admin/operations/devices` for the existing rotate/rebind flow and references dry-run helper commands for sample/simulated telemetry.
- It does not render long-lived secrets, token hashes, or raw telemetry payloads.

## Alerts/Override/Consumer/Evidence Behavior

- Alerts setup links to `/admin/alerts/console` and explains that Alerts feed availability does not prove consumer acceptance.
- Manual override/review UI is deferred.
- Consumer packet/status display uses the Phase 20 docs/evidence tracker for all seven `prepared` packet targets.
- Runtime DB consumer records are shown separately as deployment workflow records and do not override docs tracker truth.
- Evidence/readiness links make OCI pilot evidence, Phase 23 blocker, real-agency GTFS scaffold, device/AVL scaffold, consumer packets, California readiness summary, compliance checklist, and redaction policy easier to find.

## Auth/CSRF/Security Behavior

- All setup routes remain under existing admin auth middleware.
- Read-only setup view allows read-only/operator/editor/admin roles.
- Unsafe setup actions require admin role.
- Cookie-authenticated unsafe setup form posts require existing CSRF validation.
- Browser validation form does not accept validator command, validator path, output path, artifact path, or arbitrary validator ID.
- Publication setup derives `agency_id` from the authenticated principal and rejects conflicting form agency IDs if supplied.
- Errors render as bounded text rather than raw debug output.

## Tests Added And Results

- Added handler tests for setup missing-state rendering, setup publication role boundaries and server-derived agency ID, conflicting browser agency ID rejection, setup validation feed-type mapping, setup CSRF enforcement, and docs/evidence tracker consumer wording.
- Existing Operations Console tests were kept for unauthenticated access, empty state, safe telemetry diagnostics, one-time token display, device role boundaries, cookie CSRF, consumer non-acceptance wording, and readiness behavior.
- Focused check while implementing: `go test ./cmd/agency-config` — passed.

## Checks Run And Blocked Checks

- Pre-edit/planning `make validate` — passed.
- Pre-edit/planning `make test` — passed.
- Pre-edit/planning `make smoke` — passed.
- Pre-edit/planning `make demo-agency-flow` — passed.
- Pre-edit/planning `make realtime-quality` — passed.
- Pre-edit/planning `docker compose -f deploy/docker-compose.yml config` — passed.
- Pre-edit/planning `git diff --check` — passed.
- Focused implementation `go test ./cmd/agency-config` — passed.
- Post-edit `make validate` — passed.
- Post-edit `make test` — passed.
- Post-edit `make smoke` — passed.
- Post-edit `make demo-agency-flow` — passed.
- Post-edit `make realtime-quality` — passed.
- Post-edit `docker compose -f deploy/docker-compose.yml config` — passed.
- Post-edit `git diff --check` — passed.
- Local app profile `make agency-app-up` — passed.
- Local app profile `make agency-app-down` — passed.
- Local app profile `docker compose -f deploy/docker-compose.yml --profile app config` — passed.
- Blocked commands: none.

## Known Issues

- The setup checklist is still a compact server-rendered checklist, not a full multi-page wizard with saved step progress.
- Browser GTFS ZIP upload remains deferred.
- Manual override/review UI remains deferred.
- Validation report details remain summarized; deep triage still uses validator docs and stored reports.
- The consumer tracker is represented from the repo’s Phase 20 docs/evidence tracker constants in the UI rather than runtime parsing of `status.json`.
- The setup page helps operators find evidence but does not create agency-owned-domain proof, real-agency GTFS evidence, device/AVL evidence, consumer submission evidence, or compliance evidence by itself.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/track-b-productization-roadmap.md`
  - `docs/handoffs/phase-26.md`
  - `docs/prompts/calitp-truthfulness.md`
  - `SECURITY.md`
- First files likely to edit:
  - Multi-agency isolation docs or tests selected for Phase 27.
  - `internal/auth/`, `internal/compliance/`, `internal/telemetry/`, `internal/state/`, and service handlers only if Phase 27 needs explicit agency-isolation assertions.
- Commands to run before coding:
  - `make validate`
  - `make test`
  - `make realtime-quality`
  - `make smoke`
  - `docker compose -f deploy/docker-compose.yml config`
  - `git diff --check`
- Known blockers:
  - No Phase 27 detailed phase brief exists yet beyond the Track B roadmap.
  - Multi-agency production readiness must not be claimed without explicit isolation tests and deployment evidence.
- Recommended first implementation slice:
  - Define Phase 27’s exact multi-agency isolation acceptance criteria, then add tests that prove existing admin/public/feed/telemetry/device/compliance paths cannot cross agency boundaries under authenticated and unauthenticated access patterns.
