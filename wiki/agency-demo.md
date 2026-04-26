# Agency Demo

The agency demo is the fastest way to see the project working end to end.

➡️ Want the shortest setup path? Start with [Local Quickstart](local-quickstart.md).

![Illustrative agency journey from GTFS import or GTFS Studio drafts through schedule publication, authenticated telemetry, validation, and public GTFS plus GTFS Realtime feeds.](assets/agency-journey-to-public-feeds.png)

*Illustrative teaching graphic, not a product screenshot.*

## Run It

```bash
make demo-agency-flow
```

## What It Shows

The demo:

- starts the local database
- imports sample GTFS
- starts the Open Transit RT services
- creates local publication metadata
- ingests authenticated vehicle telemetry
- fetches public schedule and realtime feeds
- confirms protected admin/debug routes reject anonymous access
- runs validation flow
- reads scorecard and consumer-ingestion records

## Important Boundaries

The local demo shows the project flow on one machine. Production hosting, public HTTPS availability, consumer acceptance, learned ETA quality, and full CAL-ITP/Caltrans compliance require separate deployment and evidence work.

For the detailed script reference, see [docs/tutorials/agency-demo-flow.md](../docs/tutorials/agency-demo-flow.md).

## Next Steps

- 🚀 [Plan a pilot deployment](deployment-guide.md)
- ✅ [Review readiness evidence](readiness-and-evidence.md)
- ⭐ [Support or contribute](support-and-contribute.md)
