# Operations Console Tour

The Operations Console is the private browser UI for local evaluation and
operator review. It is not public, not a portal automation surface, and not an
evidence collector.

Open it locally after startup:

```text
http://localhost:8080/admin/operations
```

## Browser-First Product Path

Use the same seven-step path from the README, wiki home, docs home, and public
site:

1. Start in the browser.
2. Open **Agency Operations Cockpit / Start Here**.
3. Review setup.
4. Import or review GTFS.
5. Check the five public feed URLs.
6. Review feed health, readiness, validation, telemetry, connectors, and
   maintenance.
7. Understand what remains before deployment or stronger claims.

## Home

The home page starts with **Agency Operations Cockpit / Start Here**. Use it
as the first-run acceptance cockpit. It shows:

- no-developer and technical-helper paths;
- ordered first-run tasks;
- five copyable public feed URLs;
- links into browser GTFS import, feed health, telemetry, readiness, and
  Connector Hub;
- local demo / deployment / evidence boundary;
- all-false claim flags.

## Launchpad

Open:

```text
/admin/operations/launchpad
```

Launchpad keeps a broader private operator workflow with setup, GTFS, metadata,
five feeds, telemetry, validators, readiness, connector conformance, support
bundle, and decision-gate sections.

## Setup Wizard

Open:

```text
/admin/operations/setup-wizard
```

The wizard explains the staged setup path. It is read-only and points to the
existing pages that perform admin-only work.

## GTFS Import

Open:

```text
/admin/operations/gtfs-import
```

Admins can upload a GTFS ZIP or import from a safe URL. The import path uses the
existing importer and validation records, then points operators to GTFS quality,
validator health, and the five public feed paths. It does not create retained
evidence or agency approval.

## Feed Health

Open:

```text
/admin/operations/feed-health
```

Feed Health is the five-path command center for:

- `/public/feeds.json`
- `/public/gtfs/schedule.zip`
- `/public/gtfsrt/vehicle_positions.pb`
- `/public/gtfsrt/trip_updates.pb`
- `/public/gtfsrt/alerts.pb`

Use this page to see whether exact paths, URLs, validation context, freshness,
health context, and next actions are visible.

## Readiness

Open:

```text
/admin/operations/readiness
```

Readiness is a private CAL-ITP-style checklist. It distinguishes UI signals from
missing deployment evidence and keeps missing data visible.

## Connector Hub

Open:

```text
/admin/operations/connectors
```

Connector Hub explains optional sidecars, manifests, command adapters,
telemetry sources, prediction boundaries, validator adapters, monitoring/export,
and consumer/discovery boundaries. Treat it as the starting point for manifest,
redaction, fail-closed, and synthetic conformance review. It does not load
arbitrary backend plugins or prove vendor compatibility.

## Device Credentials

Open:

```text
/admin/operations/devices
```

Device Credentials shows device bindings, token status without token values,
vehicle binding, latest token use, and safe next actions for token rotation or
rebinding. Treat adjacent telemetry pages as follow-up review after credentials
are understood.

## Telemetry Freshness

Open:

```text
/admin/operations/telemetry
```

Telemetry Freshness shows latest accepted telemetry, stale state, assignment
state, match confidence, degraded or unknown reasons, and related feed-health
next actions when available. Use it before assuming a Vehicle Positions or Trip
Updates issue is caused by feed generation.

## Telemetry Simulator

Open:

```text
/admin/operations/telemetry-simulator
```

This page explains synthetic scenarios and copyable shell commands that use the
authenticated `/v1/telemetry` boundary from an operator shell. It does not
collect tokens or execute commands in the browser.

## Connector Tests

Open:

```text
/admin/operations/connectors/tests
```

Connector Tests explain synthetic manifest, sidecar, adapter, and conformance
checks. Use the Connector Hub first for the boundary overview, then this page
for local test guidance before connecting private systems.

## GTFS Quality

Open:

```text
/admin/operations/gtfs-quality
```

Use this page to review static validator and internal importer findings by
likely owner, affected files, safe fix path, verification step, and escalation
trigger. Admins can rerun the allowlisted static validator from the existing
admin action.

## Validation Health

Open:

```text
/admin/operations/validation-health
```

Validation Health shows validator tooling state and feed artifact availability.
It runs only server-side allowlisted validator actions when an admin posts the
existing form.

## Reliability

Open:

```text
/admin/operations/reliability
```

Reliability shows private diagnostics from existing records. It does not prove
SLA coverage or uptime.

## Maintenance Center

Open:

```text
/admin/operations/maintenance
```

Maintenance Center shows active feed, last import, latest five-feed check,
validator state, backup/restore configuration presence, telemetry freshness,
service-health availability, support-summary instructions, and weekly/monthly
tasks. It keeps missing values visible instead of turning them into OK.

## Help

Open:

```text
/admin/operations/help
```

Operations Console Help explains the concepts behind each major section and
links back to the right UI pages and docs.

## What Needs A Technical Helper

The UI is enough for first evaluation, but a technical helper is still needed
for local Docker issues, real deployment hosting, validator installation,
custom telemetry connectors, token security, and any future authorized evidence
intake.

## What This Does Not Prove

Using the Operations Console does not prove CAL-ITP/Caltrans compliance, agency
approval or adoption, consumer acceptance, final-root readiness, hosted SaaS,
production readiness, vendor compatibility, hardware certification, SLA/uptime
coverage, or production-grade ETA quality.
