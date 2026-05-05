# Agency One-Pager

Open Transit RT is an open-source transit data platform for small agencies that need a practical way to manage GTFS schedules, ingest vehicle telemetry, and publish GTFS Realtime feeds.

This page is a short evaluation summary. It is not an endorsement, compliance finding, procurement offer, hosted SaaS offer, paid support offer, SLA, or production-readiness claim.

## Problem

Small transit agencies often need public GTFS and GTFS Realtime feeds, but the path from schedules and vehicle data to stable public feeds can be difficult to evaluate. Common blockers include limited staff time, unclear data ownership, DNS/TLS setup, GTFS validation errors, device or AVL integration gaps, and uncertainty about consumer submission requirements.

## Open Transit RT Solution

Open Transit RT provides a mostly Go backend with Postgres/PostGIS for:

- importing or authoring static GTFS;
- ingesting authenticated vehicle telemetry;
- making conservative deterministic trip assignments;
- publishing GTFS Realtime Vehicle Positions first;
- keeping Trip Updates behind a replaceable prediction adapter;
- publishing basic Service Alerts;
- running validation, scorecard, and evidence workflows;
- supporting local demos and pilot evaluation.

## Who It Helps

- Agencies evaluating whether they can publish and operate their own feeds.
- Operators who need a lightweight local demo before committing real data.
- Civic technologists helping agencies understand GTFS and GTFS Realtime paths.
- Contributors improving validation, docs, fixtures, adapters, and operations runbooks.

## What Works Today

The repo has local demo tooling, GTFS import, GTFS Studio draft/publish flows, authenticated telemetry ingest, Vehicle Positions publication, Trip Updates diagnostics behind an adapter boundary, Alerts publication, validation records, scorecard workflows, Operations Console setup views, replay-quality checks, pilot operations runbooks, and prepared consumer packet drafts.

Current evidence includes local validation/demo workflows and the OCI DuckDNS pilot evidence packet. The OCI host is pilot evidence only, not agency-owned final-root proof.

## Pilot Path

1. Review the [Agency Pilot Program](agency-pilot-program.md).
2. Walk through the [Agency Pilot Checklist](agency-pilot-checklist.md).
3. Start the local app with `make agency-app-up`.
4. Import a public-safe GTFS ZIP or use the committed demo feed.
5. Review public feed URLs, validation records, and Operations Console setup status.
6. Test device telemetry or the documented AVL/vendor dry-run path.
7. Record public-safe feedback with the [Agency Feedback Template](agency-feedback-template.md).

## Requirements

An agency evaluation needs:

- permission to use the GTFS data being tested;
- public-safe contact, license, and feed metadata;
- an operator owner for hosting, DNS/TLS, secrets, backups, and incident response;
- a device, simulator, or AVL data path for telemetry review;
- a redaction process before anything becomes public evidence.

## Readiness Boundaries

Open Transit RT can support work toward Caltrans/CAL-ITP-style readiness, but this repo does not prove compliance by itself. Stronger claims require deployment-specific records, agency-approved feed roots, validator records for that feed scope, operations evidence, and any required third-party confirmation.

No agency-owned or agency-approved final feed root exists in repo evidence today.

## Evidence Boundaries

- Consumer and aggregator packet drafts are prepared only.
- No consumer or aggregator submission, review, acceptance, ingestion, listing, display, rejection, or blocker evidence exists in the repo.
- Phase 29A external predictor work is adapter evaluation only.
- Phase 29B AVL/vendor work is synthetic dry-run transform evidence only.
- The agency pilot package does not prove agency adoption.

## Next Steps For An Agency

- Start with the local demo and pilot checklist.
- Confirm data ownership, publication metadata, and redaction rules.
- Decide whether the next retained evidence target is agency-owned final-root proof, real agency pilot evidence, real deployment operations evidence, or an authorized target-specific consumer submission.
- Pause stronger public claims until that real retained evidence exists.
