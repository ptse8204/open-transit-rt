# CAL-ITP / Caltrans Truthfulness Guardrail

Use this guardrail whenever drafting code comments, docs, README content, deployment guides, or handoff language.

## Allowed Claims

You may claim:

- the repo is **technically aligned with California transit data expectations**
- the repo **implements the code and workflow foundations needed for CAL-ITP / Caltrans-style GTFS and GTFS Realtime readiness**
- the repo **supports stable URLs, validation workflow, license/contact metadata, and consumer-ingestion workflow records**

## Forbidden Claims Without Evidence

Do **not** claim:

- that the repo already fully meets CAL-ITP / Caltrans requirements
- that the feeds are accepted by Google Maps, Apple Maps, Transit App, or other major consumers
- that deployment is production-ready everywhere
- that the project is equivalent to a full marketplace/service vendor
- that consumer ingestion has occurred unless there is real evidence for the specific deployment

## Required Evidence Before Stronger Claims

Do not claim compliance or completed readiness unless the implementation plus deployment evidence supports at least:

- stable public URLs
- public publication
- all required feeds published
- explicit open data license and contact metadata
- canonical validator success
- truthful consumer-ingestion workflow evidence
- any third-party consumer acceptance that is claimed

## Required Writing Pattern

When in doubt, prefer wording like:

- "supports"
- "implements the technical foundations for"
- "provides the workflow needed for"
- "can be deployed toward"

Do not prefer wording like:

- "is compliant"
- "is accepted by"
- "meets all requirements"
- "production ready for all agencies"

unless the repo and deployment evidence truly support it.
