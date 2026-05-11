# Agency Evaluation Checklist

Use this checklist for a local/reference trial. It is designed for agencies,
civic technologists, and developer integrators who want to see whether Open
Transit RT fits their operating context.

This checklist is not adoption proof, production readiness proof, consumer
acceptance, final-root evidence, or CAL-ITP/Caltrans compliance.

## 30-Minute Demo

```bash
make check
make agency-app-up
```

Confirm:

- the local app starts at `http://localhost:8080`;
- `/admin/operations` shows **Start Here**;
- `/public/feeds.json` is available;
- the schedule ZIP and three GTFS-Realtime protobuf paths respond;
- Feed Health, Readiness, Connector Hub, Telemetry Simulator, and Help are
  visible in the private UI.

If a developer is available, optionally run:

```bash
make telemetry-simulator
```

Then review telemetry and feed health in the private UI.

## One-Day Local/Reference Trial

1. Pick a public GTFS ZIP and the `agency_id` from `agency.txt`.
2. Run:

   ```bash
   make agency-pilot-up AGENCY_ID=agency GTFS_URL=https://example.org/gtfs.zip
   ```

3. Review the five public paths:

   ```text
   /public/feeds.json
   /public/gtfs/schedule.zip
   /public/gtfsrt/vehicle_positions.pb
   /public/gtfsrt/trip_updates.pb
   /public/gtfsrt/alerts.pb
   ```

4. Run or intentionally skip validators based on local tooling availability.
5. Run synthetic telemetry or the synthetic AVL dry-run.
6. Review `/admin/operations`, `/admin/operations/feed-health`, and
   `/admin/operations/readiness` through a private admin boundary.
7. Keep trial notes private unless they are redacted and intended for public
   contribution.

Detailed guide: [Small-Agency Acceptance Script](../docs/tutorials/small-agency-acceptance-script.md).

## Connectors And Telemetry

Before using a real GPS/AVL source:

- review [Connector Cookbook](connector-cookbook.md);
- run `make external-connection-check`;
- run `make adapter-conformance`;
- keep real tokens and raw private payloads out of committed files;
- map device and vehicle identifiers through deployment-owned configuration.

## Readiness Review

Use [CAL-ITP Readiness Plain English](calitp-readiness-plain-english.md) to
understand the gap between local software capability and public readiness
claims.

Formal agency approval, final feed-root evidence, and consumer acceptance are
not required to use or improve the software; they are future evidence
milestones only for agencies that choose public launch or compliance claims.
