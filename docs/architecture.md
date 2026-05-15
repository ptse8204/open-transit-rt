# Architecture

Open Transit RT is a mostly-Go, self-hosted transit data platform for small
agencies. The product boundary is GTFS import/authoring, telemetry ingest,
conservative matching, GTFS Realtime feed publication, validation, monitoring,
connectors, and private operator workflows.

It is not a rider app, fare system, CAD/dispatch replacement, hosted SaaS
offering, or evidence of agency adoption, consumer acceptance, compliance, or
production readiness.

## Product Shape

The core product path is:

1. import or author static GTFS;
2. keep draft GTFS and published feed versions separate;
3. ingest authenticated vehicle telemetry;
4. persist telemetry and operational state;
5. assign vehicles conservatively, preferring `unknown` over false certainty;
6. publish GTFS-RT Vehicle Positions first;
7. publish Trip Updates through a pluggable prediction adapter boundary;
8. publish GTFS-RT Alerts through an operational alert lifecycle;
9. validate, monitor, and review readiness in private operator surfaces.

## Main Modules

| Area | Responsibility | Boundary |
| --- | --- | --- |
| GTFS importer and publisher | Import ZIPs, record feed versions, publish active schedule outputs, and expose schedule review state. | Draft editing and published active feed state remain separate. |
| GTFS Studio | Browser authoring and review for draft schedule data. | Drafts do not become active feed versions until an explicit publish workflow succeeds. |
| Telemetry ingest | Accept authenticated `POST /v1/telemetry` observations from devices, simulators, sidecars, or adapters. | Device credentials stay deployment-owned; raw private payloads do not belong in public docs or fixtures. |
| State and matching | Store events, evaluate staleness, assignment confidence, overrides, service day, after-midnight trips, frequencies, and block continuity. | Manual overrides beat automatic matching; low confidence stays unknown. |
| Vehicle Positions publisher | Emit GTFS-RT Vehicle Positions from current telemetry and conservative assignments. | Vehicle Positions is the first production-grade realtime target. |
| Prediction adapter | Provide a narrow boundary for deterministic or external Trip Updates predictors. | Trip Updates logic must not be hard-wired into telemetry ingest or Vehicle Positions publication. |
| Trip Updates publisher | Emit or withhold GTFS-RT Trip Updates based on adapter output and conservative diagnostics. | Withheld/fallback states are valid when confidence or data quality is insufficient. |
| Alerts publisher | Author, review, publish, archive, and validate GTFS-RT Alerts. | Alerts remain independent from telemetry ingest and prediction internals. |
| Connectors and examples | Show synthetic/local adapter patterns for telemetry, prediction, validators, monitoring/export, and feed metadata. | Examples are not named vendor compatibility, hardware certification, or real device proof. |
| Validation and monitoring | Run local/off-host checks, summarize feed health, and expose private diagnostics. | Validator and health signals are supporting diagnostics, not compliance or consumer-acceptance proof. |
| Private Operations Console | Browser-first operator workflow for setup, GTFS, realtime, validation, connectors, maintenance, help, access, and audit. | Admin surfaces stay authenticated and no-store; public feed paths remain anonymous. |

## Feed Boundaries

Anonymous public feed paths are limited to publication outputs:

```text
/public/feeds.json
/public/gtfs/schedule.zip
/public/gtfsrt/vehicle_positions.pb
/public/gtfsrt/trip_updates.pb
/public/gtfsrt/alerts.pb
```

Admin, debug, validation, scorecard, device, evidence, authoring, and
Operations Console routes must remain private and authenticated.

Path-routed multi-agency public feed URLs may use:

```text
/public/agencies/{agency_id}/gtfsrt/vehicle_positions.pb
/public/agencies/{agency_id}/gtfsrt/trip_updates.pb
/public/agencies/{agency_id}/gtfsrt/alerts.pb
```

The path agency is the scope. Query parameters must not override that scope.

## Trip Updates Boundary

Trip Updates remain pluggable behind the prediction adapter boundary.

Input to a predictor can include:

- active GTFS feed version;
- current telemetry;
- current vehicle assignments and manual overrides;
- Vehicle Positions feed URL or feed data;
- operator configuration and thresholds.

Output from a predictor is:

- Trip Updates feed data;
- bounded diagnostics.

Supported implementation shapes include an internal deterministic predictor, an
external sidecar such as TheTransitClock, or a future replacement predictor.
The rest of the system should depend on the adapter contract, not on a specific
predictor implementation.

## Conservative Matching Rules

Matching and realtime logic must preserve support for:

- agency-local service day;
- after-midnight trips;
- repeated trip instances;
- `start_date` and `start_time`;
- `frequencies.txt`;
- block continuity;
- stale telemetry handling;
- low-confidence unknown state;
- manual override precedence;
- auditability of assignment changes.

Do not emit a trip descriptor unless confidence is above the configured
threshold. Unknown is safer than false certainty.

## Data And Tenancy Posture

The current repo has agency-scoped handlers, route parsing, role checks,
metadata-only audit browsing, and tenant-safe public feed paths. That is a
software isolation posture, not a production multi-tenant hosting claim.

Future durable tenancy changes must document the data model, migration, role
behavior, and failure modes before implementation.

## External Dependencies

External systems and tools are isolated behind explicit boundaries:

- Postgres/PostGIS for persistence;
- pinned GTFS and GTFS-RT validator tooling;
- optional prediction sidecars;
- optional telemetry/AVL/device adapters;
- optional monitoring/export processes.

When adding integrations, update `docs/dependencies.md`, define the adapter in
code/docs, make failure modes explicit, and add tests or stubs for the
boundary.

## Documentation And Claim Boundaries

README, docs, wiki, screenshots, diagrams, examples, tutorials, and public site
content are documentation aids. They do not prove:

- CAL-ITP/Caltrans compliance;
- agency adoption, endorsement, or approval;
- consumer submission, review, ingestion, listing, display, or acceptance;
- agency-owned final-root readiness;
- hosted SaaS availability;
- production readiness;
- vendor compatibility or hardware certification;
- SLA or uptime coverage;
- production-grade ETA quality.

Formal external evidence remains optional and authorization-gated.
