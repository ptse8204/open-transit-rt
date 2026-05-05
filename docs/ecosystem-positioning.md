# Ecosystem Positioning

Open Transit RT is open-source agency-side tooling for GTFS schedules, vehicle telemetry, and GTFS Realtime publication. It is designed to fit into the transit data ecosystem without implying official affiliation, sponsorship, certification, acceptance, or endorsement by any agency, Caltrans/CAL-ITP, consumer, vendor, validator, or standards body.

## GTFS And GTFS Realtime

Open Transit RT uses GTFS as the schedule data model and publishes GTFS Realtime outputs. The first production-directed realtime output is Vehicle Positions. Trip Updates remain behind a replaceable prediction adapter, and Alerts are supported as a separate feed path.

The project does not redefine GTFS or GTFS Realtime. It helps operators prepare, validate, and publish feeds that can be evaluated against those standards.

## Validators

The repo includes workflows for pinned GTFS and GTFS Realtime validator tooling and records validation outcomes through server-side allowlisted validator IDs.

Validator results are evidence for a specific artifact at a specific time. They do not by themselves prove consumer acceptance, agency endorsement, compliance, or production readiness.

## Caltrans/CAL-ITP-Style Readiness

Open Transit RT can support work toward Caltrans/CAL-ITP-style transit data readiness when paired with real deployment, operations, metadata, validation, agency ownership, and retained evidence.

The repo does not claim CAL-ITP or Caltrans compliance. It is not officially affiliated with, sponsored by, certified by, or endorsed by Caltrans or CAL-ITP.

## Downstream Consumers And Aggregators

The repo has prepared packet drafts for seven consumer and aggregator targets: Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land.

Prepared packet drafts are operator-review materials only. No repo evidence currently supports submitted, under-review, accepted, rejected, blocked, ingestion, listing, display, or adoption status for any target.

## Agency-Owned Domains

Stable public feed roots are important for real evaluation and consumer submission. The repo currently has no agency-owned or agency-approved final public feed root evidence.

The OCI DuckDNS host is pilot evidence only. It should not be described as agency-owned production-domain proof.

## TheTransitClock And External Predictor Adapters

Phase 29A evaluated an external predictor adapter boundary and reviewed TheTransitClock as a candidate-style external predictor. This work confirms adapter boundaries and safety checks; it does not start, call, bundle, endorse, certify, or prove production behavior for TheTransitClock or any other external predictor.

Trip Updates must remain pluggable so an internal deterministic predictor, an external predictor, or a later replacement can be evaluated without coupling predictor internals to telemetry ingest or Vehicle Positions publication.

## AVL And Vendor Adapters

Phase 29B added a synthetic, dry-run-only AVL/vendor adapter pilot. The adapter transforms synthetic fixture payloads into existing telemetry event JSON and prints diagnostics.

That evidence does not prove real vendor compatibility, certified hardware support, production AVL reliability, vendor endorsement, or marketplace equivalence.

## Other Open-Source Transit Tooling

Open Transit RT is intended to complement standards, validators, and other open-source transit tools. It focuses on agency-side backend workflows: GTFS import or authoring, authenticated telemetry ingest, conservative matching, feed publication, validation, Operations Console setup, runbooks, and evidence review.

It is not a rider app, fare-payment system, CAD/dispatch system, passenger account platform, or generic analytics product.

## No Logo Or Affiliation Rule

Do not use agency, Caltrans/CAL-ITP, consumer, vendor, validator, or standards-body logos in public project materials unless retained permission exists. Do not use wording that implies affiliation, sponsorship, certification, acceptance, deployment approval, or endorsement unless retained evidence supports that exact claim.
