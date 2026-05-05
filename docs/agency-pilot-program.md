# Agency Pilot Program

This guide packages Open Transit RT for small-agency pilot evaluation. It helps
an agency, operator, or evaluator decide what is needed, what the pilot can
prove, what remains outside the pilot, and how to close the pilot without
overstating the result.

The pilot package does not create paid support, SLA coverage, hosted SaaS
availability, agency endorsement, consumer acceptance, CAL-ITP/Caltrans
compliance, marketplace/vendor equivalence, production multi-tenant hosting, or
production-grade ETA proof.

## What Open Transit RT Does

Open Transit RT provides an open-source backend for:

- importing or authoring static GTFS;
- receiving token-authenticated vehicle telemetry;
- conservatively matching vehicles to service when evidence is strong enough;
- publishing GTFS and GTFS Realtime feed paths for schedule, Vehicle Positions,
  Trip Updates, and Alerts;
- keeping Trip Updates behind a replaceable prediction adapter;
- running validation, scorecard, operations, and evidence workflows;
- helping operators understand readiness gaps before public or consumer-facing
  claims are made.

Vehicle Positions remain the first production-directed realtime output. Trip
Updates exist and are measurable, but current deterministic output is
conservative and must not be described as production-grade ETA quality without
real-world evaluation evidence.

## What Open Transit RT Does Not Do

Open Transit RT is not:

- a rider-facing mobile app;
- a fare payment system;
- a CAD/dispatch replacement;
- a hosted SaaS service;
- a paid support or SLA-backed product;
- proof that an agency endorses the project;
- proof that any consumer or aggregator has accepted, ingested, listed, or
  displays a feed;
- proof of CAL-ITP/Caltrans compliance;
- a certified hardware, vendor AVL, or marketplace-equivalent offering.

## Pilot Goals

A successful pilot helps the agency and maintainers understand:

- whether the agency can run the local demo and understand the feed surfaces;
- whether agency-approved GTFS can be imported, validated, reviewed, and
  published in a controlled environment;
- whether the agency has a public-safe license, contact, and feed-root plan;
- whether a device, AVL adapter, or simulator path can send usable telemetry;
- whether Vehicle Positions behavior is understandable and useful;
- whether Trip Updates diagnostics and withheld reasons are visible;
- whether Alerts and operations workflows are understandable;
- whether evidence can be collected or blockers can be documented;
- whether support load and operator responsibilities are realistic.

Pilot success criteria are evaluation criteria only. They are not compliance,
consumer acceptance, production readiness, or agency endorsement.

## Suggested Timeline

This timeline is a planning aid, not an SLA, delivery guarantee, hosted service
promise, or production-readiness claim.

| Stage | Focus | Expected output |
| --- | --- | --- |
| Preflight | Permissions, GTFS source, operator owner, public-safe data handling, and kickoff preparation. | Pilot owner named, GTFS permission path understood, support boundary accepted. |
| Week 1 | Local demo and GTFS review. | `make agency-app-up` or equivalent local walkthrough completed, GTFS source reviewed for permission and metadata. |
| Week 2 | Real GTFS import and validation review. | Import attempted with approved GTFS or blocker documented; validation issues triaged. |
| Week 3 | Telemetry/device or simulator review. | Device, AVL adapter, or simulator path tested or documented as blocked. |
| Week 4 | Operations, evidence, and readiness review. | Operations Console, validation, evidence, support load, and public messaging boundaries reviewed. |
| Exit | Continue, pause, or close pilot. | Closeout summary completed with next operator action and claim boundaries. |

Pilots may move faster, take longer, or stop early when permission, staffing,
data, device, validation, security, or domain blockers appear.

## Responsibilities

Agency and operator responsibilities include:

- obtaining permission to use and review GTFS and metadata;
- approving public-safe license, contact, and feed-root information;
- owning DNS, TLS, hosting, database, secrets, backups, monitoring, and incident
  response for any deployment;
- keeping device tokens, admin tokens, DB URLs, portal credentials, private logs,
  raw telemetry, and private operator artifacts outside public evidence;
- deciding whether and when any consumer submission is authorized;
- operating the system after the pilot if the agency continues.

Maintainer and community help is best-effort and may include:

- reproducible bug review;
- docs corrections;
- fixture and validator finding review;
- architecture and scope clarification;
- redaction-safe evidence structure feedback;
- guidance on supported Make targets and local workflows.

No community or maintainer participation creates paid support, response targets,
SLA coverage, hosted operations, legal commitments, procurement commitments, or
agency-specific production guarantees.

Use `docs/support-boundaries.md` for the detailed support boundary.

## Evidence Boundaries

The pilot may collect public-safe evidence such as:

- local demo command results;
- public feed URL fetch summaries;
- validator status summaries;
- redacted operations notes;
- evidence blocker summaries;
- public-safe screenshots only after redaction review.

Do not collect or commit:

- real private GTFS unless explicitly reviewed and approved as public-safe;
- private contracts;
- private contacts;
- tokens;
- DB URLs with passwords;
- device tokens;
- admin tokens;
- JWT or CSRF secrets;
- private keys;
- raw telemetry payloads;
- private logs;
- portal screenshots with private data;
- private ticket links;
- private operator artifacts;
- `.cache` files.

Follow `SECURITY.md` and `docs/evidence/redaction-policy.md`.

## Consumer Submission Boundary

All seven consumer and aggregator targets remain `prepared` only unless a future
operator provides retained, redacted, target-originated evidence for a
target-specific status transition.

The current prepared targets are:

- Google Maps;
- Apple Maps;
- Transit App;
- Bing Maps;
- Moovit;
- Mobility Database;
- transit.land.

Prepared packets are review materials. They are not submissions, under-review
evidence, acceptance, ingestion, listing, display, or adoption evidence.

## Success Criteria

Pilot success can be recorded when the agency and maintainer/community reviewers
can truthfully say which of these evaluation criteria were met:

- GTFS import completed with an approved source, or a clear GTFS blocker was
  documented.
- Public schedule feed was reviewed at the pilot root.
- License, contact, and metadata were reviewed.
- Validation results were run or validator tooling blockers were documented.
- Device, AVL adapter, or simulator path was tested or documented as blocked.
- Vehicle Positions output was reviewed.
- Trip Updates quality diagnostics and withheld reasons were reviewed.
- Alerts path was reviewed.
- Operations Console setup checklist was understood by the operator.
- Operations runbook, backup, monitoring, and incident boundaries were reviewed.
- Evidence package was assembled or blocker-documented.
- Operator can explain what they own after the pilot.
- Support load and support boundary are understood.

## Failure Or Blocker Criteria

Treat the pilot as blocked, paused, or failed for the current scope when one or
more of these conditions prevents truthful evaluation:

- no agency or feed-owner permission;
- no public-safe GTFS source;
- no operator owner;
- no license or contact metadata plan;
- no agency-owned or agency-approved domain plan when final-root proof is a goal;
- no device, AVL, or simulator path;
- unresolved validation blockers;
- security or redaction concerns;
- no staff available to operate the system;
- support expectations require paid support, SLA coverage, hosted SaaS, or legal
  commitments that this repository does not provide;
- consumer submission is requested without authorization, official-path
  verification, and retained redacted evidence.

## Risk Register

| Risk | Description | Likelihood | Impact | Mitigation | Owner | Evidence needed |
| --- | --- | --- | --- | --- | --- | --- |
| Data ownership | GTFS source may not be agency-owned or approved for pilot use. | Medium | High | Confirm permission before import and record source/license. | Agency/operator | Permission note or public-source/license reference. |
| Private data leakage | Private GTFS, contacts, logs, telemetry, or screenshots could be committed. | Medium | High | Use redaction policy and commit only public-safe summaries. | Agency/operator with maintainer review | Redaction review note. |
| Secret leakage | Tokens, DB URLs, private keys, or portal credentials could appear in docs or evidence. | Medium | High | Keep secrets in private stores; scan docs before commit. | Agency/operator | Secret scan result or reviewer note. |
| Unstable public URL | Pilot root may not be agency-owned or stable enough for final proof. | High | High | Use agency-owned-domain readiness checklist; label pilot roots accurately. | Agency/operator | Domain approval, DNS/TLS/fetch proof when available. |
| GTFS validation failure | Real GTFS may fail import or canonical validation. | Medium | High | Triage validation errors before publication claims. | Agency/operator | Validation report or blocker summary. |
| Device/AVL reliability | Real devices, vendors, or adapters may not deliver fresh valid telemetry. | Medium | High | Start with simulator or dry-run adapter; test real path only with approved evidence. | Agency/operator | Telemetry acceptance summary or blocked-path note. |
| Trip Updates quality | Conservative deterministic Trip Updates may be withheld or insufficient for ETA expectations. | High | Medium | Review diagnostics and keep production-grade ETA claims out of scope. | Maintainer/community can explain; agency evaluates | Quality diagnostics summary. |
| Operations capacity | Agency may lack staff for backups, monitoring, validation, and incident response. | Medium | High | Review runbooks and assign operator owner before continuing. | Agency/operator | Operator role and cadence note. |
| Consumer submission delay | Prepared packets may not move forward without authorization or official-path evidence. | High | Medium | Keep targets `prepared`; follow submission workflow only when authorized. | Agency/operator | Target-originated evidence only if status changes later. |
| Support expectation | Pilot participants may expect paid support, SLA response, or hosted operations. | Medium | High | Review support boundaries at kickoff and closeout. | Maintainer/community and agency/operator | Support boundary acknowledgement. |
| Multi-agency boundary | Current evidence is not production multi-tenant hosting proof. | Medium | Medium | Avoid shared production tenant claims; use deployment/DB-scoped operations language. | Maintainer/community can clarify; agency/operator owns deployment | Isolation test references and deployment notes. |

## Pilot Closeout Summary

Use this mini-template at the end of the pilot.

| Field | Notes |
| --- | --- |
| Decision | Continue / pause / close. |
| What worked | Local demo, GTFS import, validation, telemetry, console, operations, evidence, or training items that were useful. |
| What blocked progress | Permission, GTFS, metadata, domain, device/AVL, validation, security, staffing, support, or consumer-submission blockers. |
| Evidence collected | Public-safe evidence paths or private-retained evidence summaries. |
| Evidence still missing | Final-root proof, validator records, telemetry proof, operations proof, target-originated consumer evidence, or other missing items. |
| Next operator action | One concrete action owned by the agency/operator. |
| Claim boundaries | State what cannot be claimed: agency endorsement, consumer acceptance, CAL-ITP/Caltrans compliance, hosted SaaS availability, paid support/SLA, production readiness, vendor equivalence, or production-grade ETA quality unless later evidence supports it. |

