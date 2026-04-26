# Local Quickstart

Use this path when you want to run Open Transit RT on your machine.

## Prerequisites

- Go matching `go.mod`
- Docker with Compose support
- `curl`, `zip`, and `unzip`
- Java for the static GTFS validator
- `python3` for the GTFS Realtime validator wrapper

## Start The Local Environment

```bash
cp .env.example .env
make dev
make validators-install
make validators-check
```

`make dev` starts Postgres/PostGIS, applies migrations, seeds local demo data, and prints local tokens.

## Run The End-To-End Demo

```bash
make demo-agency-flow
```

The demo imports sample GTFS, starts the services, publishes local feed metadata, ingests authenticated telemetry, fetches public feeds, checks protected routes, runs validation, and reads readiness records.

For the internal command-by-command reference, see [docs/tutorials/local-quickstart.md](../docs/tutorials/local-quickstart.md).
