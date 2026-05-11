# Open Transit RT

Open Transit RT is an MIT-licensed open-source backend for small transit
agencies, civic technologists, and developer integrators who want a
self-hosted path for GTFS and GTFS Realtime operations.

The product path is: import or author static GTFS, publish GTFS and all three
GTFS Realtime feed types, ingest vehicle telemetry through documented
boundaries, monitor feed health, review CAL-ITP-style readiness, and keep
stronger deployment or consumer claims separate from local evaluation.

## Start The Software UI

From a clean checkout:

```bash
make check
make agency-app-up
```

The local app normally opens at:

```text
http://localhost:8080
```

Open the private Operations Console:

```text
http://localhost:8080/admin/operations
```

Click **Start Here** first. It shows the no-developer path, developer path,
ordered first-run tasks, and the five public feed URLs.

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

## 30-Minute Local Demo

Use the local app package when you want to see the product shape quickly:

```bash
make agency-app-up
```

Then open:

```text
http://localhost:8080/admin/operations
```

Follow **Start Here**. The local package imports the committed demo GTFS
fixture, publishes local feed paths, and prints the next private UI and token
instructions.

Detailed guide: [Small Agency Quick Start](wiki/small-agency-quick-start.md).

## One-Day Public GTFS Trial

Use the reusable onboarding helper when you want to test a public GTFS ZIP:

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

Validator output, readiness rows, and release-candidate checks are supporting
signals. They are not compliance or consumer-acceptance proof by themselves.

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

- [Small Agency Quick Start](wiki/small-agency-quick-start.md)
- [Browser-First Setup](wiki/browser-first-setup.md)
- [Operations Console Tour](wiki/operations-console-tour.md)
- [Wiki Home](wiki/README.md)
- [Documentation Home](docs/README.md)
- [Architecture](docs/architecture.md)
- [Dependencies](docs/dependencies.md)
- [Current Status](docs/current-status.md)
- [Latest Handoff](docs/handoffs/latest.md)
- [Phase 61+ Product Roadmap](docs/roadmaps/agency-first-connector-platform/README.md)

## Contributing

Agencies and contributors can help by trying the local workflow, testing with a
public GTFS ZIP, improving docs, adding synthetic connector examples, and
reporting exact blockers.

See [How Agencies Can Help](wiki/how-agencies-can-help.md) and
[CONTRIBUTING.md](CONTRIBUTING.md).
