# Open Transit RT

Open Transit RT is an MIT-licensed open-source backend for small transit
agencies, civic technologists, and developer integrators who want a
self-hosted path for GTFS and GTFS Realtime operations.

The product path is: import or author static GTFS, publish GTFS and all three
GTFS Realtime feed types, ingest vehicle telemetry through documented
boundaries, monitor feed health, review CAL-ITP-style readiness, and keep
stronger deployment or consumer claims separate from local evaluation.

Public explainer site:
[https://ptse8204.github.io/open-transit-rt/](https://ptse8204.github.io/open-transit-rt/)

Current maintainer next step: continue adoption-first product hardening toward
a `v0.1.0-rc.1` agency evaluation release candidate. The default path is
browser-first operations, docs, remote diagnostics, off-host validation, and
connector usability, not real agency pilot evidence. See
[Review And Recommendations](docs/roadmap-status.md#review-and-recommendations).

## Start In The Browser

No-developer review should start from the private local URL provided by a
technical helper, normally:

```text
http://localhost:8080/admin/operations
```

Click **Agency Operations Cockpit / Start Here** first. It shows setup
progress, primary action cards, ordered first-run tasks, the five public feed
URLs, maintenance tasks, and what each page does not prove.

## Technical Helper Startup

From a clean checkout:

```bash
make check
make agency-app-up
```

These commands are startup and health-check steps for a technical helper. They
are not the first step for no-developer review. The helper should leave the app
running and provide the private local browser URL and any local admin-token
instructions printed by `make agency-app-up`.

The local app root normally responds at:

```text
http://localhost:8080
```

Open the private Operations Console:

```text
http://localhost:8080/admin/operations
```

Stop the local app with:

```bash
make agency-app-down
```

## What You Can Do In The UI

- Review agency and publication metadata.
- Follow the setup wizard.
- Import GTFS through browser upload or safe URL import.
- Inspect `/public/feeds.json`, schedule, Vehicle Positions, Trip Updates, and
  Alerts paths.
- Review feed health and validation health.
- Review GTFS quality guidance.
- Bind or review device/telemetry state.
- Try synthetic telemetry through documented commands.
- Review Connector Hub and connector test guidance.
- Review CAL-ITP-style readiness and Help.
- Use the Maintenance Center to see weekly/monthly tasks and technical-helper
  cases.

## Private Operations Route Map

These private browser routes are the acceptance-critical navigation surface:

```text
/admin/operations
/admin/operations/setup-wizard
/admin/operations/gtfs-import
/admin/operations/feed-health
/admin/operations/readiness
/admin/operations/gtfs-quality
/admin/operations/validation-health
/admin/operations/devices
/admin/operations/telemetry
/admin/operations/telemetry-simulator
/admin/operations/connectors
/admin/operations/connectors/tests
/admin/operations/maintenance
/admin/operations/help
```

## 30-Minute Local Demo

Use the local app package when a technical helper can start the product shape
quickly:

```bash
make agency-app-up
```

Then start no-developer review from:

```text
http://localhost:8080/admin/operations
```

Follow **Agency Operations Cockpit / Start Here**. The local package imports
the committed demo GTFS fixture, publishes local feed paths, and prints the
next private UI and token instructions.

Detailed guide: [Small Agency Quick Start](wiki/small-agency-quick-start.md).

## Browser-First GTFS Review

Use the browser path first when possible:

```text
/admin/operations/gtfs-import
/admin/operations/feed-health
/admin/operations/gtfs-quality
/admin/operations/validation-health
/admin/operations/maintenance
```

The private UI can import/review GTFS, show active feed version and source
details, inspect feed health, explain validator state, review telemetry
readiness, and show the next maintenance action. See
[No Command Line First Run](docs/tutorials/no-cli-agency-first-run.md) and
[Small Agency Maintenance Guide](docs/tutorials/small-agency-maintenance-guide.md).

Use the reusable onboarding helper only as a technical-helper fallback for a
scripted public GTFS ZIP import path:

```bash
make agency-pilot-up AGENCY_ID=agency GTFS_URL=https://example.org/gtfs.zip
```

Then review the private UI:

- `/admin/operations`
- `/admin/operations/gtfs-import`
- `/admin/operations/feed-health`
- `/admin/operations/readiness`
- `/admin/operations/connectors`

Detailed guide: [Reusable Agency Onboarding](docs/tutorials/reusable-agency-onboarding.md).

## Public Feed URLs

An active local or deployment instance exposes these anonymous feed paths:

```text
/public/feeds.json
/public/gtfs/schedule.zip
/public/gtfsrt/vehicle_positions.pb
/public/gtfsrt/trip_updates.pb
/public/gtfsrt/alerts.pb
```

Admin, debug, validation, scorecard, device, evidence, and authoring routes
must stay private and authenticated.

## Connect Telemetry / GPS / AVL

Vehicle, GPS, AVL, CSV, or sidecar sources should transform observations into
the authenticated telemetry boundary:

```text
POST /v1/telemetry
Bearer device token required
JSON telemetry payload required
```

Start with:

- [Connector Cookbook](wiki/connector-cookbook.md)
- [Integration Adapter Kit](docs/integration-adapter-kit.md)
- [Device And AVL Integration](docs/tutorials/device-avl-integration.md)
- [External Adapter Conformance](docs/tutorials/external-adapter-conformance.md)

## Readiness And Validation

Open Transit RT supports CAL-ITP-style readiness workflows through private UI
and local checks. Use:

- `/admin/operations/feed-health`
- `/admin/operations/validation-health`
- `/admin/operations/readiness`
- [CAL-ITP Readiness Plain English](wiki/calitp-readiness-plain-english.md)
- [Release-Candidate Readiness](docs/release-candidate-readiness.md)
- [External Connection Readiness](docs/external-connection-readiness.md)
- [Off-Host Public Feed Validation](docs/deployment/off-host-validation.md)
- [OCI Reference Check](docs/deployment/oci-reference-check.md)
- [Product Screenshots](docs/assets/product-screenshots/README.md)
- [Product Diagrams](docs/assets/product-diagrams/README.md)

Validator output, readiness rows, and release-candidate checks are supporting
signals. They are not compliance or consumer-acceptance proof by themselves.
Screenshots and diagrams are local/demo documentation aids only, not retained
evidence or production/adoption/compliance proof.

## What This Does Not Prove

Local evaluation does not prove:

- CAL-ITP/Caltrans compliance;
- agency adoption or approval;
- consumer submission, review, acceptance, ingestion, listing, or display;
- agency-owned final-root readiness;
- hosted SaaS availability;
- production readiness;
- vendor compatibility or hardware certification;
- SLA or uptime coverage;
- production-grade ETA quality.

Formal external evidence is optional and future. It is not required for local
evaluation or open-source contribution.

## Documentation

- [Product Explainer Site](https://ptse8204.github.io/open-transit-rt/)
- [Small Agency Quick Start](wiki/small-agency-quick-start.md)
- [Browser-First Setup](wiki/browser-first-setup.md)
- [No Command Line First Run](docs/tutorials/no-cli-agency-first-run.md)
- [Small Agency Maintenance Guide](docs/tutorials/small-agency-maintenance-guide.md)
- [Operations Console Tour](wiki/operations-console-tour.md)
- [Wiki Home](wiki/README.md)
- [Documentation Home](docs/README.md)
- [Architecture](docs/architecture.md)
- [Dependencies](docs/dependencies.md)
- [Current Status](docs/current-status.md)
- [Latest Handoff](docs/handoffs/latest.md)
- [Phase 61+ Product Roadmap](docs/roadmaps/agency-first-connector-platform/README.md)
- [Adoption Productization Roadmap](docs/roadmaps/agency-first-connector-platform/adoption-productization-roadmap.md)

## Contributing

Agencies and contributors can help by trying the local workflow, testing with a
public GTFS ZIP, improving docs, adding synthetic connector examples, and
reporting exact blockers.

See [How Agencies Can Help](wiki/how-agencies-can-help.md) and
[CONTRIBUTING.md](CONTRIBUTING.md).
