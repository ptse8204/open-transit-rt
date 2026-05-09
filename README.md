# Open Transit RT

Open Transit RT is an open-source backend for small transit agencies that need
to publish GTFS and GTFS Realtime feeds without buying a full CAD/AVL or rider
app platform.

It is built around a practical first goal: import or author static GTFS, ingest
vehicle telemetry, match vehicles conservatively to scheduled service, and
publish stable GTFS-RT Vehicle Positions first. Trip Updates stay pluggable so
an agency can use the internal deterministic predictor, an external predictor,
or a later replacement without rewriting telemetry ingest or Vehicle Positions.

## Who It Is For

- Small agencies evaluating a low-cost self-hosted realtime stack.
- Civic technologists helping an agency publish schedule and realtime feeds.
- Operators who need stable GTFS/GTFS-RT URLs, validation workflows, and basic
  admin tools.
- Developers who want adapter boundaries for AVL/device data, validators,
  monitoring, and future prediction engines.

## What Works Today

The repository has code and docs for:

- GTFS ZIP import and typed GTFS Studio draft/publish workflows.
- Authenticated vehicle telemetry ingest with device bearer tokens.
- Conservative deterministic trip matching that prefers unknown over false
  certainty.
- Public schedule, `feeds.json`, Vehicle Positions, Trip Updates, and Alerts
  endpoints.
- A pluggable Trip Updates prediction adapter boundary.
- DB-backed Service Alerts authoring and GTFS-RT Alerts publication.
- Pinned static and realtime validator workflows.
- Local app packaging through `make agency-app-up`.
- Small-agency pilot operations helpers for validation, backup, restore drills,
  feed monitoring, and scorecard export.
- Read-only reference deployment diagnostics through `make deployment-doctor`.
- A private authenticated operator checklist at `/admin/operations/checklist`
  and `/admin/operations/checklist.json` for setup/readiness next actions.
- Documentation for CAL-ITP-style readiness workflows without claiming
  compliance.

## Three Paths

### 1. Try Locally

Start the local evaluation stack:

```bash
make agency-app-up
```

This starts the local app at `http://localhost:8080`, imports the committed demo
GTFS fixture, publishes the five public feed paths, and prints admin/device next
steps. See [Agency First Run](docs/tutorials/agency-first-run.md).

Stop it with:

```bash
make agency-app-down
```

### 2. Try A Real Public GTFS Feed

Use the reusable agency onboarding helper when you want to provide an agency ID
and public GTFS ZIP URL without manual database edits:

```bash
make agency-pilot-up AGENCY_ID=agency GTFS_URL=https://example.org/gtfs.zip
```

The helper downloads the GTFS ZIP into ignored `.cache/` storage, imports it,
verifies the five public feed paths, checks that the fetched schedule matches
the imported GTFS, and prints validator status or blockers. Publication
metadata is local/reference placeholder metadata unless the operator supplies
agency-approved values.

Use the public-GTFS local/pilot runbook when you need a fuller repeatability
guide without implying agency approval or consumer acceptance:

- [Reusable Agency Onboarding](docs/tutorials/reusable-agency-onboarding.md)
- [Self-Hosted Operator Trial](docs/tutorials/self-hosted-operator-trial.md)
- [Operator Smoke And Support Bundle](docs/tutorials/operator-smoke-and-support-bundle.md)
- [Phase 43 Operator UX Setup V2](docs/phase-43-operator-ux-setup-v2.md)
- [Public GTFS Local/Pilot Runbook](docs/tutorials/public-gtfs-local-pilot.md)

That guide records source URL, checksum, import summary, fetched schedule proof,
five-path fetches, validator results or blockers, and claim boundaries.

### 3. Deploy Using The OCI/OCL-Style Reference Path

The current self-hosted reference path is the existing OCI/OCL-style pilot
server pattern: compiled Go services, Postgres/PostGIS, Caddy or equivalent
reverse proxy, systemd services/timers, pinned validators, backups, feed
monitoring, and scorecard export.

Start with:

- [OCI/OCL Reference Deployment](docs/deployment/oci-reference-deployment.md)
- [Reference Deployment Env Example](docs/deployment/oci-reference-env.example)
- [Reference Deployment Smoke Checklist](docs/deployment/oci-reference-smoke-checklist.md)
- [Reference Deployment Doctor](docs/deployment/reference-deployment-doctor.md)
- [Self-Hosted Operator Trial](docs/tutorials/self-hosted-operator-trial.md)
- [Operator Smoke And Support Bundle](docs/tutorials/operator-smoke-and-support-bundle.md)
- [Small-Agency Pilot Operations](docs/runbooks/small-agency-pilot-operations.md)
- [Self-Hosted Agency Reuse Master Plan](docs/master-plan-self-hosted-agency-reuse.md)
- [Phase 36 Reference Deployment Productization](docs/phase-36-oci-reference-deployment-productization.md)
- [Phase 37 Reusable Agency Onboarding Flow](docs/phase-37-agency-reusable-onboarding-flow.md)

The existing OCI DuckDNS pilot remains hosted/operator pilot evidence only. It
is not agency-owned final-root proof.

## Integration Boundaries

- AVL/device data should enter through the existing telemetry contract or an
  adapter that transforms external payloads before calling telemetry ingest.
  Start with the [Integration Adapter Kit](docs/integration-adapter-kit.md),
  then use [Device And AVL Integration](docs/tutorials/device-avl-integration.md)
  for the detailed telemetry tutorial.
- External predictors must stay behind `internal/prediction.Adapter`. Vehicle
  Positions generation remains independent of external predictor availability.
- Validators are pinned tooling invoked through allowlisted validator IDs and
  repo-supported install/check workflows. Validator success alone is not
  consumer acceptance or compliance.
- Monitoring is deployment-owned. The repo exposes lightweight metrics and
  pilot operations helpers, but does not provision Prometheus, Grafana, or SLO
  operations by itself.
- Consumers and aggregators fetch standard public GTFS/GTFS-RT URLs. Prepared
  packets exist, but target statuses must not move beyond `prepared` without
  retained target-originated evidence.

## CAL-ITP-Style Readiness

Open Transit RT supports technical foundations for California transit data
readiness: stable public feed paths, static GTFS, all three GTFS Realtime feed
types, license/contact metadata, validation records, feed discovery, and
consumer packet preparation.

Use these docs for the current boundary:

- [California Readiness Summary](docs/california-readiness-summary.md)
- [Compliance Evidence Checklist](docs/compliance-evidence-checklist.md)
- [CAL-ITP Readiness Checklist](docs/tutorials/calitp-readiness-checklist.md)

The repo does not claim CAL-ITP/Caltrans compliance.

## What This Is Not

Open Transit RT is not:

- a hosted SaaS claim;
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

- [Documentation Home](docs/README.md)
- [Current Status](docs/current-status.md)
- [Latest Handoff](docs/handoffs/latest.md)
- [Architecture](docs/architecture.md)
- [Dependencies](docs/dependencies.md)
- [Roadmap Status](docs/roadmap-status.md)
