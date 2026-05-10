# Can My Agency Use This?

Open Transit RT is a good fit when an agency, civic technologist, or developer
integrator wants to evaluate a self-hosted open-source backend for schedule and
realtime feed publication.

## Good Fit

Open Transit RT is worth evaluating if you have or want:

- a public GTFS ZIP, or a small schedule you want to author and publish;
- GTFS-Realtime Vehicle Positions as the first high-quality output;
- Trip Updates and Alerts available through documented backend boundaries;
- an operator-owned deployment, local trial, or civic technology partner;
- a telemetry source that can be transformed into `POST /v1/telemetry`;
- a preference for conservative matching that uses unknown instead of false
  certainty;
- validation, readiness, and connector checks that can run before any public
  launch claim.

## Not A Fit Yet

This project is not currently a fit if you need:

- hosted SaaS availability;
- paid support or SLA coverage;
- a certified hardware or named vendor compatibility guarantee;
- a full CAD/dispatch system;
- fare payments, rider accounts, or a rider-facing app;
- proof of consumer acceptance, listing, ingestion, or display;
- a blanket production-readiness or CAL-ITP/Caltrans compliance guarantee.

## Try It Locally

Run the 30-minute local evaluator path:

```bash
make agency-app-up
make telemetry-simulator
```

This starts the local app, imports the committed demo GTFS fixture, publishes
the five public feed paths, and sends synthetic telemetry through real device
auth. See [Agency Demo](agency-demo.md) and
[Agency First Run](../docs/tutorials/agency-first-run.md).

## Try Your Public GTFS ZIP

Run a local/reference trial with your public GTFS ZIP:

```bash
make agency-pilot-up AGENCY_ID=agency GTFS_URL=https://example.org/gtfs.zip
```

Use [Agency Adoption Checklist](agency-adoption-checklist.md) and
[Reusable Agency Onboarding](../docs/tutorials/reusable-agency-onboarding.md)
for the full one-day evaluator path.

## Connect A GPS Or AVL Source

Start with [Connector Cookbook](connector-cookbook.md). The normal path is to
write a sidecar, script, or private operator process that transforms source
observations into the authenticated Open Transit RT telemetry contract.

Do not commit real credentials, raw private payloads, or private vehicle
identifiers to this repository.

## Readiness Boundary

Open Transit RT supports CAL-ITP-style readiness workflows. It does not claim
CAL-ITP/Caltrans compliance. Formal agency approval, final feed-root evidence,
and consumer acceptance are not required to use or improve the software; they
are future evidence milestones only for agencies that choose public launch or
compliance claims.
