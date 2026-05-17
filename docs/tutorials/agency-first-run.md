# Agency First Run

Start here instead if you are agency staff reviewing the product from the
browser: [No Command Line First Run](no-cli-agency-first-run.md).

This guide is technical-helper detail for starting Open Transit RT locally and
understanding the startup output without learning every backend service first.

The local package is a demo/evaluation flow. It is not hosted SaaS, not a production deployment, and not proof that Google Maps, Apple Maps, Transit App, or any other consumer has accepted the feeds.

## Start The Local App

Prerequisites:

- Docker with Compose support
- `curl`
- Go is useful for normal development, but the local app package builds the Go services into a local container image

For development bootstrap preflight before running `make dev`, use:

```bash
scripts/bootstrap-dev.sh --check
```

The preflight checks local tools, required repository files, Docker daemon
availability, Docker Compose config, and likely port conflicts without starting
services or changing database state. If it reports a missing tool or blocked
Docker daemon, resolve that blocker before starting the local app or dev
bootstrap.

Run:

```bash
make agency-app-up
```

The command starts Postgres/PostGIS, builds local service binaries, applies migrations, seeds demo records, imports `testdata/gtfs/valid-small`, publishes it as the active local feed, bootstraps publication metadata, waits for service readiness, and verifies the public feed URLs.

At the end it prints the public feed root, feed URLs, Operations Console URL, GTFS Studio/admin URL, token instructions, log command, validation next step, and the exact next action.

If startup fails, run:

```bash
make agency-app-logs
scripts/bootstrap-dev.sh --check
docker compose -f deploy/docker-compose.yml ps
```

Common local blockers are Docker not running, host port `55432` already in
use, host port `8080` already in use for the local app proxy, missing Go for
development commands, or pinned validator tooling not installed. These are
setup blockers, not evidence or production-readiness signals.

If a previous local database volume is reused and the demo state looks stale,
run `make agency-app-reset` only after confirming you want to remove local demo
database state and container logs. If feed readiness times out, inspect
`make agency-app-logs` before rerunning startup.

Stop the local app:

```bash
make agency-app-down
```

Reset local demo state:

```bash
make agency-app-reset
```

Reset is destructive. It removes the local Compose containers, the Postgres volume, local demo database state, and container logs.

## What GTFS Means

GTFS is the schedule data format used by trip planners. It describes routes, stops, trips, stop times, calendars, and related schedule details.

In this repo, a GTFS ZIP can be imported, validated by the internal importer, and published as the active schedule feed. The local app imports the small committed fixture so you can see the feed URLs work.

For a real agency ZIP, use [Real Agency GTFS Onboarding](real-agency-gtfs-onboarding.md) before importing. That guide covers source permission, metadata approval, validation triage, publish review, privacy/redaction checks, and final-root claim boundaries. Do not treat the local demo fixture or demo metadata as agency-approved.

If you already have a reviewed agency ID and GTFS URL, use
[Reusable Agency Onboarding](reusable-agency-onboarding.md). That flow imports
the requested GTFS without running the demo sample import path.

## What GTFS Realtime Means

GTFS Realtime is the live data format used beside the static GTFS schedule. Open Transit RT publishes:

- Vehicle Positions: where vehicles are now
- Trip Updates: conservative trip prediction output
- Alerts: service notices

Vehicle Positions are the first production-directed output. Trip Updates remain behind a replaceable prediction adapter.

## What A Device Token Does

A device token is a Bearer credential used by a vehicle device or simulator to send telemetry.

The token is bound to:

- agency
- device ID
- vehicle ID

If a device is rebound to a different vehicle, the token rotates and the previous token stops working.

For local testing:

```bash
scripts/device-onboarding.sh sample
scripts/device-onboarding.sh simulate --dry-run
scripts/device-onboarding.sh simulate
```

For a repeatable synthetic simulator that sends events through the real
authenticated ingest path:

```bash
make telemetry-simulator
RUN_MATCHER=true make telemetry-simulator
```

The simulator writes private diagnostics under ignored `.cache/` storage and
does not create evidence packets. See
[Telemetry Simulator And Device Trial](telemetry-simulator-and-device-trial.md).

For payload fields, response behavior, troubleshooting, and vendor adapter boundaries, see [Device And AVL Integration](device-avl-integration.md). For rotation, rebinding, secure storage, and compromise response, see [Device Token Lifecycle](device-token-lifecycle.md).

To rotate or bind a local device token:

```bash
scripts/device-onboarding.sh rebind --device-id device-1 --vehicle-id bus-1
```

The rebind command prints the returned one-time token because the existing API intentionally returns it. Store it securely if you need to use it. Do not commit it.

## What Public Feed URLs Are

Public feed URLs are the stable URLs that trip planners and validators fetch.

The local app uses:

```text
http://localhost:8080/public/gtfs/schedule.zip
http://localhost:8080/public/feeds.json
http://localhost:8080/public/gtfsrt/vehicle_positions.pb
http://localhost:8080/public/gtfsrt/trip_updates.pb
http://localhost:8080/public/gtfsrt/alerts.pb
```

These local URLs are for demo packaging only. Production deployments need stable HTTPS URLs and deployment-owned reverse proxy controls.

## What The Operations Console Shows

The local app exposes the authenticated Operations Console at:

```text
http://localhost:8080/admin/operations
```

The guided setup checklist is available at:

```text
http://localhost:8080/admin/operations/setup
```

Start with **Agency Operations Cockpit / Start Here** on `/admin/operations`,
then use **Agency Setup** at `/admin/operations/setup-wizard`. Agency Setup
shows progress, review blocks, safe diagnostics, role visibility,
technical-helper escalation cards, feed URLs, validation status, telemetry
freshness, safe device binding information, Alerts links, prepared-only
consumer tracker status, setup next actions, and the status source for each
setup step. Admin routes still require an admin token; the local proxy does
not make them public.

The Agency Operations Cockpit / Start Here path on `/admin/operations` is the
first routine review screen. It summarizes setup progress, GTFS import state, active feed
version, five feed paths, validator state, telemetry, readiness, connectors,
maintenance, and the next action for every card.

## Why Validation Matters

Validation checks whether schedule and realtime feeds follow the expected GTFS and GTFS Realtime rules. Passing validation is necessary for a credible feed, but it is not the same as consumer acceptance.

The local app does not fail startup just because validator tooling is missing. To install and check the pinned validator tools:

```bash
make validators-install
make validators-check
```

The executable agency demo runs the validation workflow:

```bash
make demo-agency-flow
```

For common real GTFS import and validator failures, see [GTFS Validation Triage](gtfs-validation-triage.md).

When a small server cannot run validators comfortably, ask a technical helper
to run off-host validation from an operator machine:

```bash
PUBLIC_BASE_URL=http://localhost:8080 make validate-public-feeds
```

See [Off-Host Public Feed Validation](../deployment/off-host-validation.md).

## Routine Maintenance

Open:

```text
http://localhost:8080/admin/operations/maintenance
```

Use this page for weekly/monthly feed health, validator, GTFS, telemetry,
alerts, backup/restore configuration, and support-summary tasks. See
[Small Agency Maintenance Guide](small-agency-maintenance-guide.md).

## What Consumers Still Need To Accept Separately

Trip planners and aggregators have their own ingestion processes. This repo can track consumer workflow records, but those records are not acceptance evidence by themselves.

Do not claim that a feed has been accepted by a consumer unless there is explicit evidence from that consumer for the specific feed and URL root.

## What The Demo Proves

The local app proves that this repository can run the current stack locally, import sample GTFS, publish public feed paths, protect admin/debug surfaces with auth, accept token-authenticated telemetry, and expose feed discovery.

It does not prove:

- full CAL-ITP/Caltrans compliance
- universal production readiness
- hosted SaaS availability
- consumer acceptance
- agency endorsement
- learned ETA quality
