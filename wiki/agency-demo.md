# Agency Demo

The agency demo is the fastest way to see the project working end to end.

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
- starts the current services
- creates local publication metadata
- ingests authenticated vehicle telemetry
- fetches public schedule and realtime feeds
- confirms protected admin/debug routes reject anonymous access
- runs validation flow
- reads scorecard and consumer-ingestion records

## What It Does Not Prove

The local demo does not prove production hosting, public HTTPS availability, consumer acceptance, learned ETA quality, or full CAL-ITP/Caltrans compliance.

For the internal script-level reference, see [docs/tutorials/agency-demo-flow.md](../docs/tutorials/agency-demo-flow.md).
