# Browser-First Setup

This guide shows the private UI path first. It is for evaluators who want to
understand Open Transit RT from the browser before using deeper CLI tools.

## Technical Helper Startup

A technical helper starts the local runtime from the checkout:

```bash
make check
make agency-app-up
```

These are startup and health-check commands, not no-developer evaluator steps.
The helper should leave the app running and provide the private local browser
URL and any local admin-token instructions printed by `make agency-app-up`.

No-developer evaluators start from the provided private local URL, normally:

```text
http://localhost:8080/admin/operations
```

## Use Agency Operations Cockpit / Start Here First

On the Operations Console home, start with **Agency Operations Cockpit / Start Here**.

Treat it as the single first-run cockpit. It groups the path into:

1. Start in the browser.
2. Open **Agency Operations Cockpit / Start Here**.
3. Review setup.
4. Import or review GTFS.
5. Check the five public feed URLs.
6. Review feed health, readiness, validation, telemetry, connectors, and
   maintenance.
7. Understand what remains before deployment or stronger claims.

Each row shows status, current signal, what it means, next action, UI link,
docs link, and what the row does not prove.

## No Developer Today

Use this path when you only want to review the product from a browser. Start
from the private local URL provided by the technical helper:

```text
http://localhost:8080/admin/operations
```

- open the setup wizard;
- review publication metadata;
- use browser GTFS import if you have an admin role;
- open feed health;
- review GTFS quality triage after import;
- open readiness;
- open Device Credentials;
- open Telemetry Freshness;
- open Telemetry Simulator guidance;
- open Connector Hub;
- open Connector Tests;
- open Maintenance Center;
- open Operations Console Help when labels are unclear.

Admin-only actions remain admin-only. The UI does not bypass CSRF, role checks,
or private route boundaries.

## Technical Helper Available

Use this path when someone can run local commands:

```bash
make check
make agency-app-up
make telemetry-simulator
```

Then return to the browser and review:

- `/admin/operations/devices`
- `/admin/operations/telemetry`
- `/admin/operations/feed-health`
- `/admin/operations/readiness`
- `/admin/operations/connectors/tests`

The browser points to commands; it does not execute them.

## Import Or Review GTFS

Open:

```text
/admin/operations/gtfs-import
```

Admins can upload a GTFS ZIP or import from a safe URL. After an import, review:

```text
/admin/operations/gtfs-quality
/admin/operations/validation-health
/admin/operations/feed-health
```

If browser upload is not appropriate, use the CLI import tutorials as the
fallback.

## Check The Five Feed Paths

Use the Agency Operations Cockpit / Start Here copy section and the Feed
Health command center:

```text
/admin/operations/feed-health
```

Confirm these paths are listed:

- `/public/feeds.json`
- `/public/gtfs/schedule.zip`
- `/public/gtfsrt/vehicle_positions.pb`
- `/public/gtfsrt/trip_updates.pb`
- `/public/gtfsrt/alerts.pb`

Local listing is not a consumer or compliance claim.

Feed Health should show exactly those five paths with a current signal,
freshness, validator context, health context, and next action for each row.

## Review Readiness And Help

Open:

```text
/admin/operations/readiness
/admin/operations/help
```

Readiness rows show the private signal, why it matters, next action, and claim
boundary. Operations Console Help topics explain GTFS, GTFS Realtime,
connectors, validators, telemetry, readiness, and evidence boundaries.

## Review Device Credentials And Telemetry Freshness

Open:

```text
/admin/operations/devices
/admin/operations/telemetry
/admin/operations/telemetry-simulator
```

Device Credentials shows device bindings and token status without exposing
token values. Telemetry Freshness shows latest accepted telemetry, stale state,
assignment state, match confidence, or unknown reasons when available.
Telemetry Simulator remains a technical-helper command guide.

## Review Connector Tests And Maintenance Center

Open:

```text
/admin/operations/connectors
/admin/operations/connectors/tests
/admin/operations/maintenance
```

Connector Tests show synthetic manifest and conformance guidance. Maintenance
Center shows active feed, import, five-feed check, validator, backup/restore,
telemetry freshness, service-health, and support-summary rows where configured.

## What Needs A Technical Helper

Ask for technical help when you need Docker troubleshooting, real deployment
hosting, validator installation, custom connector code, device-token handling,
or future authorized evidence work.

## What This Does Not Prove

Browser-first setup is a private local evaluation. It does not prove agency
approval, consumer acceptance, compliance, final-root readiness, hosted SaaS,
production readiness, vendor compatibility, SLA coverage, or production-grade
ETA quality.
