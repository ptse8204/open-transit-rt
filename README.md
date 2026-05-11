# Open Transit RT

Open Transit RT is an MIT-licensed open-source backend for small transit
agencies, civic technologists, and developer integrators who want to evaluate a
self-hosted path for GTFS and GTFS Realtime publication without buying a full
CAD/AVL or rider app platform.

The product starts with a practical center: import or author static GTFS,
ingest vehicle telemetry, match vehicles conservatively to scheduled service,
and publish GTFS-RT Vehicle Positions first. Trip Updates stay pluggable behind
a prediction adapter, and Alerts are part of the public feed set.

## Status At A Glance

| Area | Current repo support |
| --- | --- |
| Static GTFS | GTFS ZIP import plus typed GTFS Studio draft/publish workflows |
| Vehicle telemetry | Authenticated `POST /v1/telemetry` with device bearer tokens |
| Matching | Conservative deterministic matching that prefers unknown over false certainty |
| Public feeds | `/public/feeds.json`, schedule ZIP, Vehicle Positions, Trip Updates, and Alerts |
| Trip Updates | Internal deterministic predictor plus replaceable adapter boundary |
| Alerts | DB-backed Service Alerts authoring and GTFS-RT Alerts publication |
| Evaluation tooling | Local app package, public GTFS onboarding, synthetic telemetry, connector checks, release-candidate checks |
| Readiness | CAL-ITP-style readiness workflows without claiming compliance |

## Start Here

| If you are... | Use this path |
| --- | --- |
| An agency or civic technologist asking whether this is useful | [Can My Agency Use This?](wiki/can-my-agency-use-this.md) |
| Evaluating locally for 30 minutes | `make agency-app-up`, then [Agency First Run](docs/tutorials/agency-first-run.md) |
| Trying a public GTFS ZIP | `make agency-pilot-up AGENCY_ID=agency GTFS_URL=https://example.org/gtfs.zip` |
| Testing synthetic telemetry | `make telemetry-simulator` |
| Connecting GPS, AVL, or another telemetry source | [Integration Adapter Kit](docs/integration-adapter-kit.md) and [Connector Cookbook](wiki/connector-cookbook.md) |
| Reviewing release-candidate readiness from a fresh clone | `make release-candidate-check` |
| Reviewing CAL-ITP-style readiness | [Plain-English Readiness Guide](wiki/calitp-readiness-plain-english.md) |
| Contributing | [How Agencies Can Help](wiki/how-agencies-can-help.md) and [CONTRIBUTING.md](CONTRIBUTING.md) |

For the full local command map, run:

```bash
make help
```

For a lightweight no-network evaluator check, run:

```bash
make check
```

Agencies can help the project by trying the local/reference workflow, testing
with their public GTFS ZIP, contributing connector examples, reviewing
deployment docs, sharing non-sensitive feedback, or sponsoring a later pilot.
Formal agency approval, final feed-root evidence, and consumer acceptance are
not required to use or improve the software; they are future evidence
milestones only for agencies that choose public launch or compliance claims.

## 30-Minute Local Demo

Start the local evaluation stack:

```bash
make agency-app-up
```

Then send synthetic telemetry through the real authenticated ingest path:

```bash
make telemetry-simulator
```

The local app starts at `http://localhost:8080`, imports the committed demo GTFS
fixture, publishes the five public feed paths, and prints admin/device next
steps. Stop it with:

```bash
make agency-app-down
```

This is a local evaluator workflow. It is not hosted SaaS availability,
production readiness, agency approval, consumer acceptance, or CAL-ITP/Caltrans
compliance.

## One-Day Public GTFS Trial

Use the reusable agency onboarding helper when you want to import a public GTFS
ZIP without manual database edits:

```bash
make agency-pilot-up AGENCY_ID=agency GTFS_URL=https://example.org/gtfs.zip
```

The helper downloads the ZIP into ignored `.cache/` storage, imports it,
verifies the five public feed paths, checks that the fetched schedule matches
the imported GTFS, and prints validator status or blockers. Publication
metadata remains local/reference placeholder metadata unless the operator
supplies agency-reviewed contact and license values.

Detailed evaluator guides:

- [Reusable Agency Onboarding](docs/tutorials/reusable-agency-onboarding.md)
- [Self-Hosted Operator Trial](docs/tutorials/self-hosted-operator-trial.md)
- [Agency Launchpad](docs/tutorials/agency-launchpad.md)
- [Public GTFS Local/Pilot Runbook](docs/tutorials/public-gtfs-local-pilot.md)

## Public Feed Contract

Open Transit RT publishes these anonymous feed paths for an active local or
deployment instance:

```text
/public/feeds.json
/public/gtfs/schedule.zip
/public/gtfsrt/vehicle_positions.pb
/public/gtfsrt/trip_updates.pb
/public/gtfsrt/alerts.pb
```

Admin, JSON debug, validation, scorecard, device, and alert-authoring routes
must remain protected behind admin authentication and deployment network
controls.

## Integration Boundaries

- AVL/device data should enter through the existing telemetry contract or a
  sidecar that transforms external observations before calling
  authenticated `POST /v1/telemetry`.
- External predictors must stay behind `internal/prediction.Adapter`. Vehicle
  Positions generation remains independent of external predictor availability.
- Validators are pinned tooling invoked through allowlisted validator IDs.
  Validator success alone is not consumer acceptance or compliance.
- Monitoring is deployment-owned. The repo exposes lightweight metrics and
  local diagnostics, but does not provision a monitoring service or SLA.
- Consumers and aggregators fetch standard public GTFS/GTFS-RT URLs. Prepared
  packets exist, but target statuses stay `prepared` without retained
  target-originated evidence.

Start with:

- [Integration Adapter Kit](docs/integration-adapter-kit.md)
- [Connector Cookbook](wiki/connector-cookbook.md)
- [Device And AVL Integration](docs/tutorials/device-avl-integration.md)
- [External Adapter Conformance](docs/tutorials/external-adapter-conformance.md)

## CAL-ITP-Style Readiness

Open Transit RT supports technical foundations for California transit data
readiness: stable public feed paths, static GTFS, all three GTFS Realtime feed
types, license/contact metadata, validation workflows, feed discovery, and
consumer packet preparation.

Use these docs for the current boundary:

- [CAL-ITP Readiness Plain English](wiki/calitp-readiness-plain-english.md)
- [CAL-ITP Readiness Checklist](docs/tutorials/calitp-readiness-checklist.md)
- [California Readiness Summary](docs/california-readiness-summary.md)
- [Phase 60 Final Claim Review](docs/phase-60-final-claim-review-and-public-closeout.md)

The repo does not claim CAL-ITP/Caltrans compliance.

## What This Is Not

Open Transit RT is not:

- a hosted SaaS service;
- a paid support or SLA offering;
- proof of agency endorsement, adoption, or approval;
- proof of agency-owned final-root readiness;
- proof of consumer submission, review, acceptance, ingestion, listing, or
  display;
- proof of CAL-ITP/Caltrans compliance;
- proof of production readiness or production multi-tenant hosting;
- proof of vendor AVL compatibility or certified hardware support;
- proof of production-grade ETA quality;
- a rider-facing app, fare-payment system, passenger account system, or CAD /
  dispatch replacement.

All seven tracked consumer and aggregator targets remain `prepared` only unless
future retained, redacted, target-originated evidence supports a specific status
change.

## Documentation

- [Wiki Home](wiki/README.md)
- [Documentation Home](docs/README.md)
- [Current Status](docs/current-status.md)
- [Latest Handoff](docs/handoffs/latest.md)
- [Phase 61+ Product Roadmap](docs/roadmaps/agency-first-connector-platform/README.md)
- [Post-60 Product Roadmap](docs/post-60-product-roadmap.md)
- [Architecture](docs/architecture.md)
- [Dependencies](docs/dependencies.md)
- [Roadmap Status](docs/roadmap-status.md)
- [License](LICENSE)
