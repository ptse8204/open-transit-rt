# No Command Line First Run

This guide is for a small agency operator who wants to review Open Transit RT
from the browser first.

Open Transit RT is self-hosted open-source software for GTFS and GTFS-Realtime
publication workflows. It helps an operator import or review a GTFS schedule,
check the five configured feed paths, understand validator and data-quality
state, review telemetry/device readiness, and decide the next maintenance
action.

This guide does not prove agency adoption, CAL-ITP/Caltrans compliance,
consumer acceptance, final-root ownership, hosted service availability,
production readiness, vendor compatibility, SLA coverage, or production-grade
ETA quality.

## Before You Start

A technical helper may still need to start the local app or reference
deployment. After that, the routine review path should start in the private
browser UI:

```text
/admin/operations
```

Use Agency Operations Cockpit / Start Here as the first screen. It shows setup
progress, action cards, next actions, and claim boundaries in one place.

## Browser-First Product Path

Use this order before interpreting readiness or asking for stronger claims:

1. Start in the browser.
2. Open **Agency Operations Cockpit / Start Here**.
3. Review setup.
4. Import or review GTFS.
5. Check the five configured feed URLs.
6. Review feed health, readiness, validation, telemetry, connectors, and
   maintenance.
7. Understand what remains before deployment or stronger claims.

## 1. Open Agency Operations Cockpit / Start Here

Go to:

```text
/admin/operations
```

Review:

- agency metadata;
- GTFS imported state;
- active feed version;
- five configured feed paths;
- validators;
- telemetry;
- readiness;
- maintenance.

Each action card has a current signal, the next action, a private admin link,
and what the card does not prove.

## 2. Import Or Review GTFS

Open:

```text
/admin/operations/gtfs-import
```

Admins can import a GTFS ZIP by browser upload or by a safe HTTP(S) URL. The
page shows the source type, source URL when applicable, checksum, byte count,
import timestamp, active feed version, and schedule identity summary when the
current import path can provide it.

After import, review:

- the GTFS Workbench **Agency Review Summary**, which explains required files,
  row counts, service dates, routes/stops/trips, import history, what changed,
  and validation triage in one table;
- the GTFS Workbench **Validation Issue Triage**, which shows likely owner,
  plain-English meaning, suggested fix path, safe next action, and what the
  row does not prove;
- routes, stops, trips, stop times, and shapes counts;
- validation/import warnings grouped by file;
- GTFS quality next actions;
- feed health next actions;
- validator health next actions.

If staged comparison or browser rollback is not available in the current
runtime, the UI says so instead of pretending rollback exists. A technical
helper may use the documented CLI rollback path when needed.

## 3. Check The Five Public Feed Paths

Open:

```text
/admin/operations/feed-health
```

The page tracks exactly these public paths:

```text
/public/feeds.json
/public/gtfs/schedule.zip
/public/gtfsrt/vehicle_positions.pb
/public/gtfsrt/trip_updates.pb
/public/gtfsrt/alerts.pb
```

For each path, review the configured URL, recorded HTTP status, byte count,
content type, checksum, generated/checked time, validator state, health state,
next action, and what it does not prove.

Use this page before asking a technical helper to run curl commands.

## 4. Understand Realtime Usefulness

Stay on Feed Health and review the realtime usefulness section:

- Vehicle Positions: empty or non-empty, vehicle count when known, stale or
  suppressed rows, latest telemetry age when available.
- Trip Updates: empty or non-empty, deterministic or external adapter state,
  withheld counts/reasons when available, fallback/error state when available.
- Alerts: empty or non-empty, active alert count when available, Alerts Console
  link.

If Vehicle Positions are empty, go to Devices and Telemetry Simulator. If Trip
Updates are empty, review telemetry, active GTFS, assignment confidence, and
Trip Updates diagnostics before assuming there is a problem.

For the broader private realtime review, open:

```text
/admin/operations/realtime
/admin/operations/prediction-lab
```

## 5. Review GTFS Quality And Validators

Open:

```text
/admin/operations/gtfs-workbench
/admin/operations/gtfs-quality
/admin/operations/validation-center
/admin/operations/validation-health
```

GTFS Quality explains likely owners, affected files, what to fix first,
verification steps, and escalation triggers.

GTFS Workbench is the staff review page for the active schedule. Start with
Agency Review Summary and Validation Issue Triage before using the bounded
preview tables. The Workbench remains read-only: it does not import, edit,
publish, run validators, execute rollback, create evidence, contact external
systems, or change consumer status.

Validator Health distinguishes:

- internal import validation;
- canonical MobilityData static GTFS validation;
- GTFS-Realtime validation.

Admins can request allowlisted server-side validator runs from the browser.
The browser cannot supply validator commands, paths, URLs, argument lists,
artifacts, binaries, or timeouts.

## 6. Review Device Credentials, Telemetry Freshness, And Simulator

Open:

```text
/admin/operations/devices
/admin/operations/telemetry
/admin/operations/telemetry-simulator
```

Devices & Tokens shows device bindings, token status without token values,
vehicle binding, and latest token use. Telemetry Freshness shows latest
accepted telemetry time, assignment state, match confidence or unknown reason
when available, and stale telemetry state. Telemetry Simulator shows a
browser-only synthetic dry-run preview for committed fixture summaries, plus
fixed technical-helper commands for private shell dry-runs or intentional local
sends.

Token creation can happen from the private Devices & Tokens page for admins,
but the one-time token must still be stored outside the browser after
creation. Simulator sending may still need a technical helper because device
tokens stay in the operator shell and are not collected by the browser preview.

## 7. Review Connectors

Open:

```text
/admin/operations/connectors
/admin/operations/connectors/tests
```

Open Transit RT connectors are bounded adapters, manifests, sidecars, or
connector processes. They are not arbitrary dynamic backend plugins.

Use connector pages to understand telemetry, predictor, validator,
monitoring/export, and discovery boundaries before connecting private systems.

## 8. Review Maintenance Center

Open:

```text
/admin/operations/maintenance
```

The Maintenance Center summarizes deployed version presence, active feed
version, last GTFS import, last five-feed check, validator state, backup and
restore-drill configuration presence, telemetry freshness, service health
availability, support-summary instructions, and weekly/monthly tasks.

If a value is not available or not configured, the UI says so. It does not
turn missing data into OK.

## 9. Review Operations Console Help

Open:

```text
/admin/operations/help
```

Help & Tutorials explains Start Here, Devices & Tokens, Telemetry
Freshness, Telemetry Simulator, Connector Tests, Maintenance Center, GTFS,
GTFS Realtime, validators, readiness, and evidence boundaries.

The app shell groups routes as Start Here, Setup, GTFS Workbench, Feed Health,
Validation, Realtime, Devices / AVL, Prediction / ETA Lab, Connectors, Alerts,
Readiness, Maintenance, Help / Tutorials, and Support / Troubleshooting. GTFS
Studio and Alerts Console remain separate private tools when linked from the
Operations Console.

## What Still Needs A Technical Helper

Use a technical helper for:

- starting Docker or a server deployment;
- changing DNS, TLS, reverse proxy, firewall, or systemd setup;
- installing pinned validators, Java, Docker, or off-host validation tooling;
- large GTFS imports or rollback execution when browser rollback is not
  implemented;
- secure storage of one-time device tokens after browser rotation/rebind;
- simulator sends that require private shell credentials;
- GPS/AVL adapter development;
- `make oci-reference-check`, `make validate-public-feeds`, or support-bundle
  runs when a local shell is needed;
- any future evidence intake, consumer submission, or public claim workflow.

## Safe To Share vs Private

Safe public feed paths are the five anonymous `/public/...` paths. Private
admin pages, support summaries, logs, tokens, database URLs, raw telemetry,
validator raw reports, and `.cache` diagnostics are private unless a separate
retained-evidence approval exists.

Routine agency setup and review should start in the browser. Shell commands
remain a technical-helper path for deployment, diagnostics, validators, and
secure token handling.
