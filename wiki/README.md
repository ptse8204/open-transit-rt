# Open Transit RT Wiki

Welcome. This wiki is the task-based public guide for Open Transit RT.

Open Transit RT helps small transit agencies, civic technologists, and
developer integrators evaluate a self-hosted path for GTFS and GTFS Realtime
operations. It is open-source software, not hosted SaaS, not a CAD/AVL
replacement, and not proof of agency approval or consumer acceptance.

[Star the repo](https://github.com/ptse8204/open-transit-rt) if this work is
useful to you.

![Illustrative documentation guide showing paths for trying locally, running the agency demo, planning deployment, reviewing evidence, and contributing.](assets/docs-choose-your-path.png)

*Illustrative docs navigation graphic, not an app screenshot.*

## Start Here

1. [Small Agency Quick Start](small-agency-quick-start.md)
2. [Browser-First Setup](browser-first-setup.md)
3. [Operations Console Tour](operations-console-tour.md)
4. [Review And Recommendations](../docs/roadmap-status.md#review-and-recommendations)

Fast local command:

```bash
make agency-app-up
```

Then open:

```text
http://localhost:8080/admin/operations
```

Click **Start Here** in the private Operations Console.

## Choose Your Task

| If you want to... | Read this |
| --- | --- |
| Decide whether this fits your agency or project | [Can My Agency Use This?](can-my-agency-use-this.md) |
| Try the browser-first local path | [Small Agency Quick Start](small-agency-quick-start.md) |
| Understand what to click in the private UI | [Operations Console Tour](operations-console-tour.md) |
| Use your own public GTFS ZIP | [Agency Evaluation Checklist](agency-adoption-checklist.md) |
| Connect GPS, AVL, CSV, or sidecar telemetry | [Connector Cookbook](connector-cookbook.md) |
| Review release-candidate and connector maturity | [Review And Recommendations](../docs/roadmap-status.md#review-and-recommendations) |
| Review CAL-ITP-style readiness plainly | [CAL-ITP Readiness Plain English](calitp-readiness-plain-english.md) |
| Understand readiness and evidence boundaries | [Readiness And Evidence](readiness-and-evidence.md) |
| Plan a self-hosted evaluator or pilot deployment | [Deployment Guide](deployment-guide.md) |
| Help improve the project | [How Agencies Can Help](how-agencies-can-help.md) |
| Contribute code or docs | [Support And Contribute](support-and-contribute.md) |

## What You Can Do In The UI

- Follow a Start Here first-run path.
- Review setup and publication metadata.
- Import GTFS through the browser if you have an admin role.
- Check the five public feed paths.
- Review feed health, GTFS quality, validation health, readiness, connectors,
  telemetry simulator guidance, and Help.

## What Still Needs A Technical Helper

You may need technical help for Docker setup, large GTFS imports, validator
installation, stable HTTPS deployment, device-token handling, custom GPS/AVL
connectors, and any future authorized evidence intake.

## Important Boundaries

Local UI evaluation does not prove CAL-ITP/Caltrans compliance, agency adoption
or approval, consumer acceptance, final-root readiness, hosted SaaS,
production readiness, vendor compatibility, SLA coverage, or production-grade
ETA quality.

Formal agency approval, final feed-root evidence, consumer acceptance, real
agency pilots, and vendor proof are optional future evidence tracks when
authorized. They are not required to use or improve the software.

## Maintainers And Project History

For deeper implementation notes, handoffs, release-candidate checks, phase
history, and evidence/claim-boundary references, see [docs](../docs/README.md).
