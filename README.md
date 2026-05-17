# Open Transit RT

Open Transit RT is an MIT-licensed, self-hosted transit data backend for small
agencies and civic technology teams. It helps you import or author GTFS,
ingest vehicle telemetry, publish GTFS Realtime feeds, review feed health, and
prepare for public-data workflows from a browser after a technical helper
starts the local app.

Public explainer site:
[https://ptse8204.github.io/open-transit-rt/](https://ptse8204.github.io/open-transit-rt/)

Current public release candidate:
[`v0.1.0-rc.2`](https://github.com/ptse8204/open-transit-rt/releases/tag/v0.1.0-rc.2).
Use it for local or self-hosted evaluation. It is not a stable release and it
does not prove production readiness, compliance, agency adoption, consumer
acceptance, hosted service availability, vendor compatibility, SLA/uptime, or
production-grade ETA quality. See the
[release status](docs/release-status-v0.1.0-rc.2.md) and
[download replay](docs/release-download-replay-v0.1.0-rc.2.md).

## Who It Is For

Open Transit RT is for:

- small agencies that want a practical path from GTFS to GTFS Realtime;
- technical helpers running a local or self-hosted evaluator;
- connector developers adapting GPS, AVL, CSV, prediction, validator, or
  monitoring systems;
- maintainers improving an open-source transit backend.

It is not a rider app, fare system, CAD/dispatch system, hosted SaaS product,
or proof that any agency, vendor, consumer, or regulator has accepted a feed.

## Normal User Flow

After a technical helper starts the app, agency staff should start here:

```text
http://localhost:8080/admin/operations
```

Then follow this order in the browser:

1. Open **Agency Operations Cockpit / Start Here**.
2. Review agency setup and publication metadata.
3. Import or review GTFS.
4. Check the five feed URLs.
5. Review feed health and validation.
6. Connect or review vehicle telemetry.
7. Review connectors, readiness, maintenance, and help.
8. Read what the local review does not prove before making external claims.

The key private pages are:

```text
/admin/operations
/admin/operations/setup-wizard
/admin/operations/gtfs-workbench
/admin/operations/gtfs-import
/admin/operations/feeds
/admin/operations/feed-health
/admin/operations/validation-center
/admin/operations/realtime
/admin/operations/devices
/admin/operations/telemetry
/admin/operations/connectors
/admin/operations/connectors/workbench
/admin/operations/readiness
/admin/operations/maintenance
/admin/operations/help
```

GTFS Studio and the Alerts Console are also private admin tools:

```text
/admin/gtfs-studio
/admin/alerts/console
```

## Technical Helper Startup

From a clean checkout:

```bash
git clone https://github.com/ptse8204/open-transit-rt.git
cd open-transit-rt
git checkout v0.1.0-rc.2
make check
make agency-app-up
```

Give agency staff the private browser URL printed by the startup command. Stop
the local app with:

```bash
make agency-app-down
```

For a validation-heavy local trial:

```bash
make validators-install
make validate
make test
```

Detailed guides:

- [Small Agency Quick Start](wiki/small-agency-quick-start.md)
- [Browser-First Setup](wiki/browser-first-setup.md)
- [No Command Line First Run](docs/tutorials/no-cli-agency-first-run.md)
- [Small Agency Maintenance Guide](docs/tutorials/small-agency-maintenance-guide.md)

## CI And Validation

Fast GitHub Actions run the same lightweight baseline contributors should run
before opening a PR:

```bash
go test ./...
make check
scripts/check-consumer-tracker.sh
make audit-final-claim-review
```

The Go version follows `go.mod`; GitHub Actions reads it with
`actions/setup-go` and `go-version-file: go.mod`. Validator-heavy checks,
`make smoke`, connector conformance, GTFS-RT conformance, the standalone
product-acceptance gate, and release-package audits stay in the manual
release-gates workflow because they need pinned validators or broader
release-candidate context. The product-acceptance audit also runs inside
`make check`.

See [Continuous Integration](docs/ci.md) for the exact workflow split.

## Import GTFS

Use the browser first:

```text
/admin/operations/gtfs-import
/admin/operations/gtfs-workbench
/admin/operations/gtfs-quality
/admin/operations/validation-center
```

The GTFS Workbench shows required files, row counts, service dates, route/stop
coverage, import history, active feed version, validation issues, and safe next
actions.

For a scripted fallback with a public GTFS ZIP:

```bash
make agency-pilot-up AGENCY_ID=agency GTFS_URL=https://example.org/gtfs.zip
```

## Public Feed URLs

An active local or deployment instance exposes:

```text
/public/feeds.json
/public/gtfs/schedule.zip
/public/gtfsrt/vehicle_positions.pb
/public/gtfsrt/trip_updates.pb
/public/gtfsrt/alerts.pb
```

Admin, debug, validation, scorecard, device, evidence, and authoring routes
must stay private and authenticated.

## Connect Vehicle Data

Connector review starts in the browser:

```text
/admin/operations/connectors
/admin/operations/connectors/workbench
/admin/operations/connectors/tests
```

The connector catalog covers:

- Vehicle / GPS / AVL connectors: CSV replay adapter, HTTP polling adapter,
  webhook sidecar adapter, generic JSON transform adapter, vendor-shaped
  synthetic examples, and authenticated `POST /v1/telemetry`.
- Prediction connectors: deterministic built-in predictor, external HTTP
  predictor adapter, shadow-mode predictor, fail-closed behavior, and
  TheTransitClock candidate notes only.
- Validator connectors: MobilityData static GTFS validator, MobilityData GTFS
  Realtime validator, allowlisted validator IDs, and private validation health.
- Monitoring/export connectors: local health summaries, operations notify
  draft, monitoring/export helper, and deployment-owned monitoring boundary.
- Consumer/discovery connectors: `/public/feeds.json`, static GTFS URL,
  Vehicle Positions URL, Trip Updates URL, Alerts URL, and prepared packet
  review without submission or acceptance claims.
- Future connector extension model: manifest-based sidecars, no arbitrary
  dynamic backend plugin loading, and conformance tests required.

Start with:

- [Connector Catalog](docs/connectors/catalog.md)
- [Connector Cookbook](wiki/connector-cookbook.md)
- [Integration Adapter Kit](docs/integration-adapter-kit.md)
- [Device And AVL Integration](docs/tutorials/device-avl-integration.md)
- [External Adapter Conformance](docs/tutorials/external-adapter-conformance.md)

## Review Readiness

Open Transit RT supports CAL-ITP-style readiness workflows. In the browser,
review:

```text
/admin/operations/feed-health
/admin/operations/validation-health
/admin/operations/readiness
```

The readiness page covers public feed URLs, static GTFS, Vehicle Positions,
Trip Updates, Alerts, validation, license/contact metadata, operations signals,
telemetry/device state, and consumer preparedness. Each area explains what the
local review helps prepare and what it does not prove.

More detail:

- [CAL-ITP Readiness Plain English](wiki/calitp-readiness-plain-english.md)
- [CAL-ITP-Style Readiness Checklist](docs/tutorials/calitp-readiness-checklist.md)
- [External Connection Readiness](docs/external-connection-readiness.md)

## What Is Not Proven Yet

Local evaluation does not prove:

- CAL-ITP/Caltrans compliance;
- agency adoption or approval;
- consumer submission, review, acceptance, ingestion, listing, or display;
- agency-owned final-root readiness;
- hosted SaaS availability;
- production readiness;
- vendor compatibility or hardware certification;
- SLA or uptime coverage;
- production AVL reliability;
- production-grade ETA quality or real-world ETA accuracy.

Formal external evidence is optional and future. It is not required for local
evaluation or open-source contribution.

## Documentation

Use the role-based docs index first:

- [Docs Index](docs/index.md)
- [Wiki Home](wiki/README.md)
- [Product Explainer Site](https://ptse8204.github.io/open-transit-rt/)
- [Operations Console Tour](wiki/operations-console-tour.md)
- [Evaluator And Contributor Kit](docs/adoption/evaluator-and-contributor-kit.md)
- [Product Acceptance](docs/product-acceptance/post-rc2-browser-first-acceptance.md)
- [Branching And Release Policy](docs/branching-and-release-policy.md)
- [Continuous Integration](docs/ci.md)

## Contributing

Agencies and contributors can help by trying the local workflow, testing with a
public GTFS ZIP, improving docs, adding synthetic connector examples, and
reporting exact blockers.

Start with:

- [How Agencies Can Help](wiki/how-agencies-can-help.md)
- [Contributing](CONTRIBUTING.md)
- [Contributor First Issues](docs/contributor-first-issues.md)
- [Contributing Connectors](docs/connectors/contributing-connectors.md)
