# Browser-First Setup

This guide shows the private UI path first. It is for evaluators who want to
understand Open Transit RT from the browser before using deeper CLI tools.

## Start The Local UI

Run:

```bash
make check
make agency-app-up
```

Open the local URL printed by the command, normally:

```text
http://localhost:8080/admin/operations
```

## Use Start Here First

On the Operations Console home, start with **Start Here**.

It groups the first-run path into:

1. Set agency/publication metadata.
2. Import or publish GTFS.
3. Check the five public feed paths.
4. Run or review validation health.
5. Add or test telemetry.
6. Review Vehicle Positions, Trip Updates, and Alerts.
7. Review readiness.
8. Review connectors.
9. Review support and release-candidate checks.

Each row shows status, current signal, what it means, next action, UI link,
docs link, and what the row does not prove.

## No Developer Today

Use this path when you only want to review the product from a browser:

- open the setup wizard;
- review publication metadata;
- use browser GTFS import if you have an admin role;
- open feed health;
- open readiness;
- open Connector Hub;
- open telemetry simulator guidance;
- open Help when labels are unclear.

Admin-only actions remain admin-only. The UI does not bypass CSRF, role checks,
or private route boundaries.

## Developer Available

Use this path when someone can run local commands:

```bash
make check
make agency-app-up
make telemetry-simulator
```

Then return to the browser and review:

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
/admin/operations/feed-health
```

If browser upload is not appropriate, use the CLI import tutorials as the
fallback.

## Check The Five Feed Paths

Use the Start Here copy section or open:

```text
/admin/operations/feeds
```

Confirm these paths are listed:

- `/public/feeds.json`
- `/public/gtfs/schedule.zip`
- `/public/gtfsrt/vehicle_positions.pb`
- `/public/gtfsrt/trip_updates.pb`
- `/public/gtfsrt/alerts.pb`

Local listing is not a consumer or compliance claim.

## Review Readiness And Help

Open:

```text
/admin/operations/readiness
/admin/operations/help
```

Readiness rows show the private signal, why it matters, next action, and claim
boundary. Help topics explain GTFS, GTFS Realtime, connectors, validators,
telemetry, readiness, and evidence boundaries.

## What Needs A Technical Helper

Ask for technical help when you need Docker troubleshooting, real deployment
hosting, validator installation, custom connector code, device-token handling,
or future authorized evidence work.

## What This Does Not Prove

Browser-first setup is a private local evaluation. It does not prove agency
approval, consumer acceptance, compliance, final-root readiness, hosted SaaS,
production readiness, vendor compatibility, SLA coverage, or production-grade
ETA quality.
