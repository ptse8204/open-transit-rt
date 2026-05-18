# Open Transit RT Wiki

This wiki is the task-based guide for evaluating Open Transit RT from the
browser.

Open Transit RT helps small agencies and civic technology teams evaluate a
self-hosted path for GTFS and GTFS Realtime operations. It is open-source
software, not hosted SaaS, not a CAD/AVL replacement, and not proof of agency
approval or consumer acceptance.

Public explainer site:
[https://ptse8204.github.io/open-transit-rt/](https://ptse8204.github.io/open-transit-rt/)

Current release candidate:
[`v0.1.0-rc.2`](https://github.com/ptse8204/open-transit-rt/releases/tag/v0.1.0-rc.2).
Use it for local/self-hosted evaluation only. It is not a stable release.

## Start Here

1. [Small Agency Quick Start](small-agency-quick-start.md)
2. [Browser-First Setup](browser-first-setup.md)
3. [No Command Line First Run](../docs/tutorials/no-cli-agency-first-run.md)
4. [Operations Console Tour](operations-console-tour.md)
5. [Small Agency Maintenance Guide](../docs/tutorials/small-agency-maintenance-guide.md)
6. [Docs Index](../docs/index.md)

Agency staff review starts from the private local setup URL provided by an
administrator:

```text
http://localhost:8080/admin/local-login
```

Choose **Start setup**. The browser opens the private Operations Console and
shows **Start Here** first.

## Administrator Startup

```bash
git clone https://github.com/ptse8204/open-transit-rt.git
cd open-transit-rt
git checkout v0.1.0-rc.2
make check
make agency-app-up
```

This starts the local app. It is helper work, not the first step for agency
staff review.

## Choose Your Task

| If you want to... | Read this |
| --- | --- |
| Decide whether this fits your agency or project | [Can My Agency Use This?](can-my-agency-use-this.md) |
| Try the browser-first local path | [Small Agency Quick Start](small-agency-quick-start.md) |
| Start after the app is running | [No Command Line First Run](../docs/tutorials/no-cli-agency-first-run.md) |
| Understand what to click in the private UI | [Operations Console Tour](operations-console-tour.md) |
| Use your own public GTFS ZIP | [Agency Evaluation Checklist](agency-adoption-checklist.md) |
| Connect GPS, AVL, CSV, prediction, validator, monitoring, or discovery systems | [Connector Catalog](../docs/connectors/catalog.md) and [Connector Cookbook](connector-cookbook.md) |
| Review CAL-ITP-style readiness plainly | [CAL-ITP Readiness Plain English](calitp-readiness-plain-english.md) |
| Understand readiness and evidence boundaries | [Readiness And Evidence](readiness-and-evidence.md) |
| Plan a self-hosted evaluator or reference deployment | [Deployment Guide](deployment-guide.md) and [Self-hosted site guide](https://ptse8204.github.io/open-transit-rt/deploy.html) |
| Help improve the project | [How Agencies Can Help](how-agencies-can-help.md) |
| Contribute code or docs | [Support And Contribute](support-and-contribute.md) |

## What You Can Do In The UI

- Follow the Start Here action queue.
- Review setup and publication metadata.
- Import GTFS through the browser if you have an admin role.
- Check the five configured feed paths.
- Review feed health, GTFS quality, validation health, readiness, connectors,
  telemetry simulator guidance, Maintenance, and Help.

## When An Administrator Is Needed

You may need an administrator or deployment owner for Docker setup, large GTFS
imports, validator installation, stable HTTPS deployment, device-token
handling, custom GPS/AVL connectors, and any future authorized evidence intake.

## Limits

Local UI evaluation does not prove CAL-ITP/Caltrans compliance, agency adoption
or approval, consumer acceptance, final-root readiness, hosted SaaS,
production readiness, vendor compatibility, SLA coverage, uptime, production
AVL reliability, or production-grade ETA quality.

For maintainer implementation notes, release checks, and project-history
records, use the [Docs Index](../docs/index.md).
