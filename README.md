# Open Transit RT

[![Star Open Transit RT on GitHub](https://img.shields.io/github/stars/ptse8204/open-transit-rt?style=social)](https://github.com/ptse8204/open-transit-rt)

Open Transit RT is an open-source backend for small transit agencies that want a practical way to publish GTFS schedules and GTFS Realtime feeds.

It is for agencies, civic technologists, operators, contributors, and evaluators who want to understand or run a lightweight transit data stack without starting from a large commercial platform.

![Illustrative diagram showing an agency moving from GTFS setup and vehicle telemetry to validation and public GTFS plus GTFS Realtime feeds.](wiki/assets/agency-journey-to-public-feeds.png)

*Illustrative teaching graphic, not a product screenshot. It shows the intended project flow at a high level.*

## What This Is

- Open-source tooling to help agencies publish GTFS and GTFS Realtime feeds.
- A mostly Go backend with Postgres/PostGIS, GTFS import, GTFS Studio drafts, vehicle telemetry ingest, conservative trip assignment, feed publishing, Alerts, validation records, and readiness evidence tracking.
- A project that can support work toward Caltrans/CAL-ITP-style transit data readiness when paired with real hosting, operations, and evidence.

## What This Is Not

- Not a hosted SaaS product.
- Not a CAD/AVL replacement.
- Not proof of consumer acceptance by itself.
- Not a claim of full CAL-ITP/Caltrans compliance or universal production readiness.

## What It Can Do Today

Open Transit RT can help you:

- import an existing static GTFS ZIP or publish typed GTFS Studio drafts
- receive authenticated vehicle telemetry
- keep vehicle matching conservative when trip assignment is uncertain
- publish public GTFS and GTFS Realtime feed outputs
- keep Trip Updates behind a replaceable prediction adapter
- author and publish basic service alerts
- run validation, scorecard, and evidence workflows
- try the full local agency demo from committed project files

## Try It Locally

Prerequisites for the simplest local trial: Docker with Compose support and `curl`.

```bash
make agency-app-up
```

The local app package starts the full stack behind `http://localhost:8080`, imports the committed sample GTFS, bootstraps feed metadata, verifies public feed URLs, and prints next steps.

Validators are optional for startup. To install and check them:

```bash
make validators-install validators-check
```

For the deeper executable demo:

```bash
make demo-agency-flow
```

For a guided version, start with the [Agency First Run](docs/tutorials/agency-first-run.md), [Local Quickstart](wiki/local-quickstart.md), or [Agency Demo](wiki/agency-demo.md).

## Where To Go Next

| Need | Start here |
| --- | --- |
| 🧭 Understand the system | [How It Works](wiki/how-it-works.md) |
| 🚌 Try the full local package | [Agency First Run](docs/tutorials/agency-first-run.md) |
| 💻 Run it on your machine | [Local Quickstart](wiki/local-quickstart.md) |
| 🚌 Walk through the agency demo | [Agency Demo](wiki/agency-demo.md) |
| 🚀 Plan a small deployment | [Deployment Guide](wiki/deployment-guide.md) |
| ✅ Review readiness and evidence | [Readiness And Evidence](wiki/readiness-and-evidence.md) |
| 📚 Browse all documentation | [Documentation Home](docs/README.md) |
| 🤝 Support or contribute | [Support And Contribute](wiki/support-and-contribute.md) |

## Evidence And Boundaries

This repo includes local validation workflows, an executable agency demo, documentation for deployment readiness, and captured OCI pilot evidence. These are useful signals for evaluating the project.

They are not the same as third-party consumer acceptance, an agency endorsement, full CAL-ITP/Caltrans compliance, or a guarantee that every deployment is production-ready. Those claims require deployment-specific records and external confirmation.

For details, see:

- [Readiness And Evidence](wiki/readiness-and-evidence.md)
- [Documentation Home](docs/README.md)

## Support The Project

⭐ **[Star Open Transit RT](https://github.com/ptse8204/open-transit-rt)** if this project is useful to you. A GitHub star is like a like or bookmark: it helps more people discover the project and supports continued independent open-source work. A star is not an official agency endorsement.

Useful contributions include focused issues, reproducible demo failures, validator findings, clearer docs, deployment notes, and small-agency workflow feedback.
