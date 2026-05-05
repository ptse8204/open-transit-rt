# Agency Feedback Template

Use this template during or after an Open Transit RT agency pilot. Keep feedback
public-safe. Do not include secrets, private logs, raw telemetry, tokens, portal
credentials, private ticket links, private GTFS, screenshots with private data,
or unredacted operator artifacts.

## Feedback Metadata

| Field | Response |
| --- | --- |
| Agency or evaluator name, if public-safe |  |
| Pilot date range |  |
| Feedback author role |  |
| Public-safe contact, if approved |  |
| Environment reviewed | Local demo / pilot root / agency-approved root / other. |
| Redaction review completed | Yes / no / not applicable. |

## Onboarding Friction

- What step was hardest to start?
- Which prerequisite was unclear?
- What would have made kickoff easier?
- Did the pilot timeline feel realistic as a planning aid?

## Docs Clarity

- Which document helped most?
- Which document was confusing or incomplete?
- Were support, evidence, and consumer-submission boundaries clear?
- Were any terms hard to understand?

## Setup Difficulty

- Which command or local setup step failed, if any?
- What operating system and Docker/Compose setup was used, if public-safe?
- What exact public-safe error message can be shared?
- What private details were intentionally omitted?

## GTFS Import Issues

- Was the GTFS source approved for pilot use?
- Did import succeed, fail validation, or remain blocked?
- Which validation findings need agency review?
- Were timezone, service day, after-midnight, frequency, shape, or block issues
  found?

## Validation Issues

- Were validator tools installed and checked?
- Which feed types were validated?
- Were findings errors, warnings, informational notices, or tooling blockers?
- What evidence can be shared publicly, if any?

## Device Or AVL Issues

- Was the telemetry path simulator, simple device, vendor/AVL adapter, or blocked?
- Were tokens handled privately?
- Did telemetry reach the Operations Console?
- Did Vehicle Positions reflect the expected vehicle state?
- What device, clock, GPS, mapping, or freshness issue blocked progress?

## Operations Console Feedback

- Which console view was most useful?
- Which status was unclear?
- Could the operator find feed URLs, validation state, telemetry freshness,
  device bindings, and evidence links?
- What workflow should be easier?

## Runbook And Operations Feedback

- Was there a named operations owner?
- Were backup, restore, monitoring, validation cadence, and incident response
  responsibilities clear?
- Which operations task felt too heavy for the agency?
- What private operations evidence was retained outside the repo?

## Missing Features

- What agency workflow was not supported?
- Is the request inside the project scope: GTFS, telemetry, matching,
  GTFS-RT feeds, validation, operations, admin workflows, or evidence?
- Would the feature affect public feed contracts, security, evidence, or
  multi-agency boundaries?

## Support Requests

- What help is requested?
- Is the request reproducible with public-safe commands or fixtures?
- Is the request best-effort community support, or does it require paid support,
  SLA coverage, hosted operations, legal review, or procurement work outside the
  repo scope?

## Bug Reports

Include:

- exact command or endpoint path;
- expected result;
- actual result;
- public-safe logs or error text;
- fixture or synthetic input, if available;
- whether the issue affects public feeds, admin views, telemetry, validation, or
  evidence.

Do not include credentials, raw private logs, raw telemetry, private GTFS,
private portal screenshots, private ticket links, or unredacted operator
artifacts.

## Training Gaps

- Which topic needs more explanation?
- Which role needs a simpler guide?
- Which training session should be split, shortened, or expanded?
- What should be added to the kickoff agenda or closeout summary?

## Claim Boundary Review

Confirm that feedback does not claim, unless separately backed by retained
evidence:

- agency endorsement;
- paid support or SLA coverage;
- hosted SaaS availability;
- consumer submission, review, acceptance, ingestion, listing, display, or
  adoption;
- CAL-ITP/Caltrans compliance;
- production multi-tenant hosting;
- production-grade ETA quality;
- vendor or hardware compatibility.

