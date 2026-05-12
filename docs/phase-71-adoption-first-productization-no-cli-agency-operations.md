# Phase 71 -- Adoption-First Productization And No-CLI Agency Operations

## Goal

Make Open Transit RT easier for small agencies, civic technologists, and
low-skill operators to adopt as a low-cost self-hosted GTFS and GTFS-Realtime
integration platform.

The phase turns the private Operations Console into the main routine agency
interface. A non-expert operator should be able to import or review GTFS,
inspect the five public feed paths, understand validator and GTFS quality
state, review telemetry/device readiness, review connector boundaries, and
know the next maintenance action without relying primarily on command-line
tools.

## Non-Goals

This phase does not:

- start, plan, or document a real agency pilot;
- collect retained evidence;
- write to `docs/evidence` except for unchanged navigation references already
  present in existing docs;
- contact agencies, vendors, consumers, aggregators, or portals;
- submit feeds to consumers or change consumer target statuses;
- claim CAL-ITP/Caltrans compliance;
- claim agency adoption, agency approval, consumer acceptance, final-root
  readiness, hosted SaaS availability, production readiness, vendor
  compatibility, hardware certification, SLA or uptime coverage, or
  production-grade ETA quality;
- add arbitrary dynamic backend plugin loading;
- introduce a heavy frontend stack.

## Current OCI Diagnostic Lessons

The May 12, 2026 OCI reference diagnostic is product feedback only. It showed:

- repo gates passed in that run;
- the current app built, pushed, migrated, restarted, and checked on OCI;
- all five services were active and loopback health returned `200`;
- real public Marin Transit static GTFS imported successfully;
- all five public feed paths returned HTTP `200`;
- schedule identity checks passed;
- static/realtime validator execution on the tiny OCI host remained blocked by
  runtime or resource constraints;
- backup and restore-drill values were not configured;
- no device credential existed for real telemetry sends;
- the public URL was a diagnostic root, not an agency-owned final public root.

Phase 71 treats those outcomes as product usability input. It does not convert
them into evidence, adoption, compliance, production, or consumer claims.

## Small-Agency Adoption Assumptions

- Operators may know GTFS and routes but may not be comfortable with shell
  commands.
- A technical helper may still be needed for installation, DNS/TLS, server
  secrets, backups, validator installation, and device credential placement.
- Routine decisions should be browser-first: what is missing, what is stale,
  what is blocked, what to fix first, and what is safe to share publicly.
- CLI helpers remain valuable as fallback and automation paths, but the UI
  should explain when they are required and why.

## UI Surfaces To Improve

- `/admin/operations`
- `/admin/operations/feed-health`
- `/admin/operations/feed-health.json`
- `/admin/operations/gtfs-import`
- `/admin/operations/gtfs-quality`
- `/admin/operations/validation-health`
- `/admin/operations/devices`
- `/admin/operations/telemetry`
- `/admin/operations/telemetry-simulator`
- `/admin/operations/connectors`
- `/admin/operations/readiness`
- `/admin/operations/maintenance`
- `/admin/operations/maintenance.json`

The pages stay private/admin or authenticated operations routes. No public admin
route is added.

## Docs, Wiki, And Site Surfaces To Improve

- `README.md`
- `docs/README.md`
- `wiki/README.md`
- `docs/roadmap-status.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/deployment/README.md`
- `docs/deployment/oci-reference-deployment.md`
- `docs/deployment/oci-reference-check.md`
- `docs/deployment/off-host-validation.md`
- `docs/tutorials/agency-first-run.md`
- `docs/tutorials/no-cli-agency-first-run.md`
- `docs/tutorials/small-agency-maintenance-guide.md`
- `docs/tutorials/small-agency-acceptance-script.md`
- `docs/tutorials/agency-launchpad.md`
- `docs/tutorials/self-hosted-operator-trial.md`
- `docs/assets/product-screenshots/README.md`
- `docs/assets/product-diagrams/README.md`
- `docs/roadmaps/agency-first-connector-platform/adoption-productization-roadmap.md`

If the GitHub Pages site must be changed from the `gh-pages` branch, this phase
will document the required updates from main rather than silently changing a
separate publication branch.

## Code And Scripts Touched

Expected code and script files:

- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_feed_health.go`
- `cmd/agency-config/operations_gtfs_import.go`
- `cmd/agency-config/operations_devices.go`
- `cmd/agency-config/operations_telemetry_simulator.go`
- `cmd/agency-config/operations_navigation.go`
- `cmd/agency-config/operations_maintenance.go`
- `cmd/agency-config/main.go`
- `cmd/agency-config/main_test.go`
- `scripts/oci-reference-check.sh`
- `scripts/validate-public-feeds.sh`
- `Makefile`

## Required Tests And Checks

Focused checkpoint checks:

- route protection and method tests for new private routes;
- HTML content tests for stable card or row IDs and plain-language next
  actions;
- JSON shape tests for feed health and maintenance;
- all-false claim flag tests;
- unsafe-field rejection tests for validator browser requests;
- script syntax and dry-run tests for new scripts.

Final requested checks:

- `git status --short`
- `git diff --check`
- `make check`
- `make test`
- `make validate`
- `make external-connection-check`
- `make adapter-conformance`
- `make test-connector-examples`
- `make audit-final-claim-review`
- `make audit-product-acceptance`
- `docker compose -f deploy/docker-compose.yml config`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact consumer tracker prepared-only guard
- protected evidence path guards

If validator tooling, Docker, network, or local services are unavailable, the
final report must list the exact blocked command and distinguish missing
tooling from validation failure.

## Protected Paths

The following paths must remain unchanged unless a separate explicit evidence
authorization exists:

- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current`
- `docs/evidence/consumer-submissions/artifacts`
- `docs/evidence/consumer-submissions/packets`
- `docs/evidence/captured`

The consumer tracker must keep these target statuses exactly `prepared`:

- Google Maps
- Apple Maps
- Transit App
- Bing Maps
- Moovit
- Mobility Database
- transit.land

## Claim Boundaries

Allowed wording:

- Open Transit RT supports self-hosted GTFS and GTFS-Realtime publication
  workflows.
- Open Transit RT supports CAL-ITP-style readiness workflows.
- Open Transit RT provides local/reference deployment diagnostics.
- Open Transit RT provides private web UI surfaces for setup, feed health,
  GTFS quality, telemetry, connectors, and readiness.
- Open Transit RT is being productized for low-cost small-agency adoption.

Forbidden wording:

- compliance achieved;
- consumer submitted, accepted, ingested, listed, displayed, or reviewed;
- agency adopted, approved, endorsed, or final-root ready;
- hosted SaaS, SLA, production ready, vendor compatible, hardware certified,
  or production-grade ETA quality.

## Success Criteria

An evaluator can use the web UI to import or review GTFS, inspect five public
feed paths, understand GTFS quality, understand validator state, review
telemetry/device readiness, understand realtime feed state, review connectors,
and know the next maintenance action without relying primarily on command-line
tools.

The phase is successful only if:

- the private Operations Console clearly shows setup progress, primary actions,
  status labels, next actions, and claim boundaries;
- feed health shows exactly the five public paths and does not invent missing
  HTTP, checksum, validator, or freshness data;
- GTFS import/update pages show source, counts, active feed version, quality,
  validator, feed-health, and rollback limitations truthfully;
- validator and GTFS quality pages explain internal import validation versus
  canonical validators and preserve server-side allowlisted execution only;
- device, vehicle, telemetry, and simulator pages explain how to make Vehicle
  Positions non-empty and why Trip Updates may be empty without exposing
  tokens;
- maintenance has private HTML and JSON summaries with all-false claim flags;
- OCI/reference diagnostics and off-host validation write ignored `.cache`
  output only and print no secrets;
- docs/wiki navigation lets non-experts start from a browser-first workflow;
- protected evidence paths and consumer statuses remain unchanged.
