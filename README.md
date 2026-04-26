# Open Transit RT

Open Transit RT is open-source tooling for small transit agencies that need a practical path to publish static GTFS and GTFS Realtime feeds.

It is meant for agencies, civic technologists, operators, and evaluators who want to understand or run a lightweight backend for schedule publication, vehicle telemetry, validation, and public feed URLs.

## What This Is

- Open-source tooling to help agencies publish GTFS and GTFS Realtime feeds.
- A mostly Go backend with Postgres/PostGIS, GTFS import, GTFS Studio drafts, authenticated vehicle telemetry, conservative trip assignment, public feed builders, Alerts, validation workflow records, and readiness evidence tracking.
- A repo that supports deployment toward Caltrans/CAL-ITP-style readiness when paired with real hosting, operations, and evidence.

## What This Is Not

- Not a hosted SaaS product.
- Not a CAD/AVL replacement.
- Not proof of consumer acceptance by itself.
- Not a claim of full CAL-ITP/Caltrans compliance or universal production readiness.

![Illustrative agency journey from GTFS import or GTFS Studio drafts through schedule publication, authenticated telemetry, validation, and public GTFS plus GTFS Realtime feeds.](wiki/assets/agency-journey-to-public-feeds.png)

*Illustrative teaching graphic, not a product screenshot. It shows the repo-supported path without claiming hosted SaaS, CAD/AVL replacement, consumer acceptance, or full compliance.*

## What It Can Do Today

The current repo can:

- import static GTFS ZIP files or publish typed GTFS Studio drafts
- persist authenticated vehicle telemetry
- preserve conservative vehicle assignment state
- publish stable public paths for `schedule.zip`, `feeds.json`, Vehicle Positions, Trip Updates, and Alerts
- keep Trip Updates behind a replaceable prediction adapter
- run validation workflows and store readiness/consumer-ingestion records
- run an executable local agency demo from the committed repo

Public feed paths are:

```text
/public/gtfs/schedule.zip
/public/feeds.json
/public/gtfsrt/vehicle_positions.pb
/public/gtfsrt/trip_updates.pb
/public/gtfsrt/alerts.pb
```

Admin, JSON debug, GTFS Studio, validation, scorecard, device, and alert-authoring routes require admin auth.

## Try It Locally

Prerequisites: Go, Docker with Compose, `curl`, `zip`, `unzip`, Java for the static GTFS validator, and `python3` for the GTFS-RT validator wrapper.

```bash
cp .env.example .env
make dev
make validators-install
make validators-check
make demo-agency-flow
```

The demo imports sample GTFS, starts the current services, publishes local feed metadata, ingests token-authenticated telemetry, fetches public feeds, verifies protected admin/debug access, runs validation, and reads scorecard plus consumer-ingestion records.

Start with the public wiki:

- [Wiki Home](wiki/README.md)
- [How It Works](wiki/how-it-works.md)
- [Local Quickstart](wiki/local-quickstart.md)
- [Agency Demo](wiki/agency-demo.md)
- [Deployment Guide](wiki/deployment-guide.md)
- [Readiness And Evidence](wiki/readiness-and-evidence.md)
- [Support And Contribute](wiki/support-and-contribute.md)

## Evidence And Status

Evidence is kept in docs instead of duplicated here:

- [Public Readiness And Evidence Guide](wiki/readiness-and-evidence.md)
- [Current Status](docs/current-status.md)
- [Latest Handoff](docs/handoffs/latest.md)
- [Compliance Evidence Checklist](docs/compliance-evidence-checklist.md)
- [Consumer Submission Evidence](docs/consumer-submission-evidence.md)
- [Consumer Submission Tracker](docs/evidence/consumer-submissions/README.md)
- [OCI Pilot Evidence Packet](docs/evidence/captured/oci-pilot/2026-04-24/README.md)

Maintainer and agent reference docs:

- [Internal Docs Home](docs/README.md)
- [Architecture](docs/architecture.md)
- [Dependencies](docs/dependencies.md)
- [Decisions](docs/decisions.md)
- [Repo Gaps](docs/repo-gaps.md)
- [Handoff History](docs/handoffs/)

Evidence boundary: validator success, public fetch proof, hosted operator evidence, and consumer-ingestion workflow records are useful, but they are not the same as third-party consumer acceptance.

## What Is Not Claimed Yet

Do not overstate the current repo:

- No hosted login/SSO product.
- No packaged full production platform with universal operations guarantees.
- No external predictor integration such as TheTransitClock.
- No consumer submission API integrations.
- No evidence that Google Maps, Apple Maps, Transit App, or another consumer has accepted a feed from this repo.
- No full CAL-ITP/Caltrans compliance claim.

## Support Or Contribute

If this project is useful to you, consider starring the repo. A GitHub star is a simple way to show support, similar to a like or bookmark. For an independent open-source project, stars help more people discover the work, help show that the project is useful, and support continued work on it. A star is not an official agency endorsement.

Useful contributions include focused issues, reproducible demo failures, validator findings, clearer docs, deployment notes, and small-agency workflow feedback.

Please keep requests inside the product boundary: GTFS import/Studio, telemetry ingest, deterministic matching, GTFS-RT feeds, Alerts, validation, monitoring, and admin/operator workflows.
