# Small Agency Quick Start

This is the shortest path for a small transit agency, civic technologist, or
developer integrator who wants to evaluate Open Transit RT from a clean checkout
without reading maintainer phase history.

This is a local evaluation path. It does not collect retained evidence, contact
external systems, submit feeds to consumers, move consumer statuses, or prove
CAL-ITP/Caltrans compliance, agency adoption, consumer acceptance, final-root
readiness, hosted SaaS availability, vendor compatibility, SLA coverage, or
production-grade ETA quality.

## Who This Is For

Use this guide if you want to answer:

- Can I start the software locally?
- Can I use the private browser UI to understand the next step?
- Can I import or review GTFS without living in the CLI?
- Can I inspect the five public feed paths?
- Can I understand telemetry and connector options before involving a vendor or
  developer?

## What You Need

- A clean checkout of the repository.
- Docker with Compose support for the local app package.
- A shell that can run `make`.
- Optional: a public GTFS ZIP URL or local GTFS ZIP for your own trial.

You do not need agency approval, consumer confirmation, final public root
evidence, vendor credentials, real AVL hardware, or production hosting to run a
local evaluation.

## 30-Minute Local Demo

From a clean checkout:

```bash
git clone https://github.com/ptse8204/open-transit-rt.git
cd open-transit-rt
make check
make agency-app-up
```

`make agency-app-up` starts the local evaluator package, imports the committed
small GTFS fixture, publishes local feed paths, and prints the local app URL.
The URL is normally:

```text
http://localhost:8080
```

If startup fails, run:

```bash
make agency-app-logs
scripts/bootstrap-dev.sh --check
```

Common blockers are Docker not running, ports `8080` or `55432` already in use,
or a stale local demo database volume.

## Open The UI

Open the local URL printed by the command, normally:

```text
http://localhost:8080/admin/operations
```

Look for the **Start Here** section at the top of the Operations Console. It
shows a no-developer browser path, a developer path, ordered first-run tasks,
and the five public feed URLs.

## Follow The Setup Wizard

Open:

```text
/admin/operations/setup-wizard
```

Review agency profile, publication metadata, GTFS, feeds, telemetry, validators,
connectors, and readiness. The wizard is private and read-only. Admin-only
actions still live on their existing admin pages.

## Import GTFS In The Browser

Open:

```text
/admin/operations/gtfs-import
```

An admin can upload a GTFS ZIP or import from a safe URL. Read-only users can
review the page and then use GTFS quality and feed health pages after an import.

For larger or scripted imports, a technical helper can use the documented CLI
fallbacks in the deeper tutorials.

## Check Public Feed URLs

The local demo normally exposes:

```text
/public/feeds.json
/public/gtfs/schedule.zip
/public/gtfsrt/vehicle_positions.pb
/public/gtfsrt/trip_updates.pb
/public/gtfsrt/alerts.pb
```

In the Operations Console, the Start Here section and Feeds page show the
configured URLs. The local URLs are useful for local evaluation; a public
deployment still needs stable HTTPS hosting, source-of-truth listing, validation
records, and any authorized evidence required for stronger claims.

## Review Feed Health

Open:

```text
/admin/operations/feed-health
```

Review feeds.json, GTFS Schedule, Vehicle Positions, Trip Updates, and Alerts
in plain language. Each row explains the current signal, what it means, what to
do next, and what it does not prove.

## Review Readiness

Open:

```text
/admin/operations/readiness
```

This is a private CAL-ITP-style readiness review. It helps you see missing
metadata, validation, feed health, telemetry, reliability, connector, and
consumer-workflow signals. It is not compliance proof.

## Try Synthetic Telemetry

Open:

```text
/admin/operations/telemetry-simulator
```

The page describes simulator scenarios and copyable shell commands. The browser
does not collect device tokens or execute commands. A developer can run:

```bash
make telemetry-simulator
```

Then review:

```text
/admin/operations/telemetry
/admin/operations/feed-health
```

## Review Connectors

Open:

```text
/admin/operations/connectors
```

The Connector Hub explains sidecar, manifest, command-adapter, telemetry,
prediction, validator, monitoring, and discovery paths. Start with synthetic
examples and local conformance checks before connecting any external GPS/AVL or
prediction system.

## What Needs A Technical Helper

A technical helper is useful for:

- Docker or port conflicts during local startup;
- importing large GTFS ZIP files;
- diagnosing GTFS validator errors;
- configuring stable public HTTPS hosting;
- securing device tokens;
- writing a GPS/AVL transform into `POST /v1/telemetry`;
- running validator, connector, release-candidate, or deployment checks;
- preparing any future authorized evidence intake.

## What This Does Not Prove

The local quick start does not prove:

- agency approval or adoption;
- agency-owned final public root readiness;
- CAL-ITP/Caltrans compliance;
- consumer submission, review, acceptance, ingestion, listing, or display;
- hosted SaaS availability;
- production readiness;
- vendor compatibility or hardware certification;
- SLA or uptime coverage;
- production-grade ETA quality.

## Where To Go Next

- [Browser-first setup](browser-first-setup.md)
- [Operations Console tour](operations-console-tour.md)
- [Connector Cookbook](connector-cookbook.md)
- [CAL-ITP readiness in plain English](calitp-readiness-plain-english.md)
- [Small-agency acceptance script](../docs/tutorials/small-agency-acceptance-script.md)
