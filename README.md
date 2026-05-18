# Open Transit RT

Open Transit RT is an MIT-licensed, self-hosted transit data backend for small
agencies and civic technology teams. It moves normal GTFS and GTFS Realtime
evaluation work into the browser after a technical helper starts the local app:
import GTFS, connect vehicle data, check feeds, review readiness, and get clear
next actions.

Public explainer site:
[https://ptse8204.github.io/open-transit-rt/](https://ptse8204.github.io/open-transit-rt/)

Browser-first tutorial video:
[https://ptse8204.github.io/open-transit-rt/video.html](https://ptse8204.github.io/open-transit-rt/video.html)

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

## Start In The Browser

After a technical helper starts the app, agency staff open:

```text
http://localhost:8080/admin/local-login
```

Select **Start setup**. The local-only page creates a short private browser
session and opens the Operations Console.

Use the visible action groups in this order: **Start**, **Setup**, **GTFS**,
**Feeds**, **Realtime**, **Vehicles**, **Connectors**, **Readiness**,
**Maintenance**, and **Help**.

Normal browser review does not require manual tokens, curl, DevTools, or a
header extension. A technical helper is still needed for initial startup,
validator installation, stable HTTPS deployment, secrets, and custom connectors.

Helpful starting points:

- [UI Tour](https://ptse8204.github.io/open-transit-rt/ui-tour.html)
- [Video Guide](https://ptse8204.github.io/open-transit-rt/video.html)
- [No Command Line First Run](docs/tutorials/no-cli-agency-first-run.md)

## Technical Helper Startup

From a clean checkout:

```bash
git clone https://github.com/ptse8204/open-transit-rt.git
cd open-transit-rt
git checkout v0.1.0-rc.2
make check
make agency-app-up
```

Give agency staff the local browser setup URL printed by the startup command.
They do not need a raw admin token for normal browser review. Stop the local
app with:

```bash
make agency-app-down
```

For a validation-heavy local trial, add:

```bash
make validators-install
make validate
make test
```

More help:

- [Small Agency Quick Start](wiki/small-agency-quick-start.md)
- [Browser-First Setup](wiki/browser-first-setup.md)
- [No Command Line First Run](docs/tutorials/no-cli-agency-first-run.md)
- [Small Agency Maintenance Guide](docs/tutorials/small-agency-maintenance-guide.md)
- [Video Recording Guide](docs/tutorials/video-recording-guide.md)

## Validation

Before opening a PR, run:

```bash
go test ./...
make check
scripts/check-consumer-tracker.sh
make audit-final-claim-review
```

Validator-heavy checks, smoke tests, connector conformance, GTFS-RT
conformance, and release-package audits run as manual release gates. See
[Continuous Integration](docs/ci.md) for the exact workflow split.

## Import GTFS

Use the browser first:

```text
/admin/operations/gtfs-import
/admin/operations/gtfs-workbench
/admin/operations/gtfs-quality
/admin/operations/validation-center
```

The GTFS pages show required files, row counts, service dates, route/stop
coverage, import history, active feed version, validation issues, and next
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

Supported local connector paths include:

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
- Extension model: manifest-based sidecars, no arbitrary dynamic backend plugin
  loading, and conformance tests required.

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

The readiness page shows public feed URLs, static GTFS, Vehicle Positions, Trip
Updates, Alerts, validation, license/contact metadata, operations signals,
telemetry/device state, and consumer preparedness. Each area states the next
action and what the local review does not prove.

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
