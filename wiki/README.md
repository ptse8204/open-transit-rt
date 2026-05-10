# Open Transit RT Wiki

Welcome. This wiki is the public guide for Open Transit RT.

Open Transit RT helps small transit agencies and civic technologists evaluate a
self-hosted path for GTFS and GTFS Realtime publication using open-source
backend tools. It is MIT-licensed software, not a hosted SaaS product, not a
CAD/AVL replacement, and not proof of consumer acceptance by itself.

[Star the repo](https://github.com/ptse8204/open-transit-rt) if this work is
useful to you.

![Illustrative documentation guide showing paths for trying locally, running the agency demo, planning deployment, reviewing evidence, and contributing.](assets/docs-choose-your-path.png)

*Illustrative docs navigation graphic, not an app screenshot.*

## Start With This Command

```bash
make agency-app-up
```

That starts the local evaluator app at `http://localhost:8080`. It is the
fastest way to see the product shape before reading deeper docs.

## Choose Your Path

| If you want to... | Read this |
| --- | --- |
| Decide whether this fits your agency or project | [Can My Agency Use This?](can-my-agency-use-this.md) |
| Understand the system | [How It Works](how-it-works.md) |
| Try the 30-minute local demo | [Agency Demo](agency-demo.md) |
| Run a one-day local/reference trial | [Agency Adoption Checklist](agency-adoption-checklist.md) |
| Connect GPS, AVL, CSV, or sidecar telemetry | [Connector Cookbook](connector-cookbook.md) |
| Review CAL-ITP-style readiness plainly | [CAL-ITP Readiness Plain English](calitp-readiness-plain-english.md) |
| Plan a self-hosted deployment | [Deployment Guide](deployment-guide.md) |
| Understand readiness and evidence boundaries | [Readiness And Evidence](readiness-and-evidence.md) |
| Help improve the project | [How Agencies Can Help](how-agencies-can-help.md) |
| Contribute code or docs | [Support And Contribute](support-and-contribute.md) |

## Important Boundaries

Open Transit RT can publish stable public feed paths for static GTFS,
`feeds.json`, Vehicle Positions, Trip Updates, and Alerts. Formal agency
approval, final feed-root evidence, and consumer acceptance are not required to
use or improve the software; they are future evidence milestones only for
agencies that choose public launch or compliance claims.

For deeper implementation notes and maintainer references, see
[docs](../docs/README.md).
