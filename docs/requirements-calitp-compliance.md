# Requirements — Cal-ITP / Caltrans Technical Compliance and Marketplace Gaps

These requirements formalize the gap between the current Open Transit RT
software capabilities, the deployment evidence required before stronger
Caltrans/CAL-ITP-style public claims, and the separate non-code work required
for marketplace-vendor equivalence.

---

## Phase 54 official-source refresh

Phase 54 refreshed this mapping against official public sources available on
May 9, 2026. This refresh records requirement mappings only. It is not
deployment evidence, does not prove compliance, and does not change consumer
submission status.

| Source | URL | Visible version/date | Guidance category | Mapping impact |
| --- | --- | --- | --- | --- |
| Caltrans, California Transit Data Guidelines | `https://dot.ca.gov/cal-itp/california-transit-data-guidelines` | Version 4.0, December 11, 2024 | Current Caltrans guidance | Confirms the compliance characteristics for GTFS Schedule and GTFS Realtime: stable public URL, regular canonical validation with no errors, major trip-planner acceptance, and open license. Confirms realtime completeness covers Trip Updates, Vehicle Positions, and Service Alerts. |
| Caltrans, California Transit Data Guidelines FAQ | `https://dot.ca.gov/cal-itp/california-transit-data-guidelines-faqs-v4_0` | Version 4.0 | Current Caltrans FAQ | Confirms aggregator publication paths, trip-planner publication guidance, open-license guidance, validator references, and public-data rationale. |
| Caltrans, California Integrated Travel Project GTFS overview | `https://dot.ca.gov/cal-itp/cal-itp-gtfs` | No page-level version visible during Phase 54 review | Cal-ITP overview | Confirms GTFS/GTFS-RT as the California transit-data standard context and Cal-ITP assistance/contact framing. |
| Caltrans, Critical GTFS Validation Errors | `https://dot.ca.gov/cal-itp/critical-gtfs-validation-errors` | No page-level version visible during Phase 54 review | Validator context | Confirms Caltrans-recognized validator tooling context and critical-error framing. |
| Caltrans, Website Model Language | `https://dot.ca.gov/cal-itp/website-model-language` | No page-level version visible during Phase 54 review | Website/source-of-truth and licensing model language | Confirms stable feed URL, provider website listing, open usage terms, and designated technical-contact expectations. |
| FTA, 2025 NTD Reporting Policy Manual | `https://www.transit.dot.gov/ntd/2025-ntd-reporting-policy-manual` | Page last updated April 15, 2026; PDF applies to Report Year 2025 | Federal GTFS reporting context | Confirms applicable NTD reporters with fixed-route modes must maintain a public-domain GTFS dataset and a publicly accessible persistent, machine-readable, non-password-protected link for GTFS ZIP collection. |

Phase 54 also keeps the following separations explicit:

- Caltrans compliance mapping is separate from marketplace/vendor-equivalence
  positioning.
- Major trip-planner acceptance is a third-party outcome, not something repo
  code or validator results can prove.
- Transitland and Mobility Database availability are data-availability and
  discoverability recommendations, not evidence that a consumer accepted a
  feed.
- If realtime API keys are used, registration must be straightforward,
  automated, transparent, and suitable for HTTPS request use; authenticated
  realtime access still must not become a login wall or arbitrary approval
  gate.

---

## Scope clarification

There are two separate targets:

1. **Technical compliance target**
   - Static GTFS and GTFS-Realtime feeds meet Caltrans data-guideline expectations.
   - Feeds are public, stable, regularly validator-clean with no errors,
     openly licensed, discoverable from provider or regional source-of-truth
     pages, and accepted by major trip-planning consumers.

2. **Marketplace-vendor equivalence target**
   - The product also behaves like a California Mobility Marketplace vendor offering, including support, integrations, optional hardware path, documentation, and service packaging.

This document separates those two targets.

---

## RQ-4A — Complete realtime feed set

### Goal
Publish the full GTFS-Realtime feed set expected for completeness:
- Trip Updates
- Vehicle Positions
- Service Alerts

### Required behaviors
- Publish all three feed types for the active GTFS schedule feed.
- Keep them publicly available.
- Keep them synchronized against the same active published GTFS feed version.
- Track health, freshness, and coverage for all three.
- Preserve trip ID consistency across GTFS Schedule, Vehicle Positions, Trip
  Updates, and Alerts where GTFS-RT references scheduled service.

### Acceptance criteria
- All three feed URLs are live and documented.
- Each feed validates cleanly.
- Coverage metrics exist for each feed type.

---

## RQ-4B — Stable public production URLs

### Goal
Meet the requirement that feeds are publicly available at stable URLs and fetchable automatically by trip-planning applications.

### Required behaviors
- Stable canonical URLs for:
  - GTFS Schedule
  - Trip Updates
  - Vehicle Positions
  - Alerts
- URL permanence across dataset refreshes
- HTTPS only
- no login wall
- if optional realtime API keys are used, registration is automated,
  nondiscriminatory, transparent about terms/rate limits, and discoverable from
  the provider or regional partner website
- support provider-domain hosting or an operationally equivalent stable hostname
- publish metadata page that lists:
  - feed URLs
  - agency name
  - technical contact
  - last updated time
  - license

### Acceptance criteria
- Consumers do not need a new URL after a publish.
- Feed URLs can be listed on the agency site and remain valid over time.
- Rollback does not change the feed URL.

---

## RQ-4C — Validator-clean feeds

### Goal
Public datasets should regularly pass canonical validators with no errors.

### Required behaviors
- Run MobilityData GTFS validator on static GTFS before publish.
- Run GTFS Realtime validation continuously or on schedule for each realtime feed.
- Block publish or mark unhealthy on validation failure.
- Store validation history and expose it in the admin UI.
- Differentiate:
  - blocking errors
  - warnings
  - informational notices

### Acceptance criteria
- Latest production feeds show no validator errors for compliant status.
- Validation reports are viewable by agency and timestamp.
- Failed validation prevents unsafe production activation where configured.

---

## RQ-4D — Open license and discoverability

### Goal
Meet expectations for open publication, source-of-truth visibility, aggregator
availability, and clear licensing.

### Required behaviors
- Associate each published feed with an explicit open license.
- Display license language on the provider-hosted landing page.
- Provide discoverability metadata:
  - feed URLs
  - provider name
  - technical contact
  - feed contact in `feed_info.txt` when supported by the active GTFS feed
  - update timestamp
  - license
- Publish or link the GTFS Schedule dataset from the provider website or a
  regional partner website that is in agreement with the provider's canonical
  feed.
- Publish or link all three GTFS Realtime datasets from the provider website or
  agreed regional partner website:
  - Trip Updates
  - Vehicle Positions
  - Service Alerts
- Support publishing or registration workflows for:
  - Mobility Database
  - transit.land

### Acceptance criteria
- A public page exists for the feed set with license and contact info.
- Feed metadata is sufficient to register with aggregators.
- License is visible without logging in.
- Provider or regional source-of-truth pages list the canonical schedule and
  realtime feeds.

---

## RQ-4E — Consumer ingestion workflow

### Goal
Make it operationally possible for major trip planners to actually ingest the feeds.

### Required behaviors
- Export a partner-ingestion packet containing:
  - stable URLs
  - contact info
  - license
  - validation status
  - sample fetch proof
- Track submission status per consumer:
  - Google Maps
  - Apple Maps
  - Transit App
  - Bing Maps
  - Moovit
- Store consumer onboarding history and notes.
- Support resubmission after feed fixes without changing feed URLs.

### Acceptance criteria
- Admin UI can show which consumers have been submitted, accepted, rejected, or pending.
- Feeds can be re-submitted after fixes without URL changes.

---

## RQ-4F — Marketplace-vendor equivalence versus technical compliance

### Goal
Separate technical compliance from the additional non-code work needed to resemble a California marketplace vendor.

### Technical-compliance minimum
To claim technical compliance, the system must satisfy:
- RQ-4A
- RQ-4B
- RQ-4C
- RQ-4D
- RQ-4E

### Additional requirements for vendor-equivalent positioning
If the goal is to resemble a California Mobility Marketplace vendor offering, the system must also support:
- optional hardware strategy or documented BYOD hardware path
- integration support for up to 3 journey-planning apps
- SLA/KPI definitions and reporting
- implementation plan templates and onboarding workflow
- support documentation and operational runbooks
- agency-facing documentation suitable for procurement and contracting
- service operations beyond code alone

### Acceptance criteria
- If targeting technical compliance only: all data compliance requirements pass.
- If targeting vendor-equivalent product: support, documentation, and service packaging exist beyond the software itself.

---

## RQ-4G — Readiness dashboard and scorecard

### Goal
The product must show readiness posture transparently without treating local
software signals as deployment evidence.

### Required behaviors
- Per-agency readiness dashboard with sections:
  - GTFS Schedule
  - Vehicle Positions
  - Trip Updates
  - Alerts
  - validation status
  - license/discoverability
  - consumer ingestion status
- Show status for:
  - missing software capability
  - capability available but unhealthy or unconfigured
  - capability available but missing deployment/operator evidence
  - deployment evidence reviewed for the selected scope
- Exportable readiness report for agency review

### Acceptance criteria
- A user can see exactly which software capability, deployment evidence, or
  third-party confirmation is still missing before stronger public claims.
- A user can download a readiness checklist or status report.

---

## Current capability and evidence summary

Software capability exists for GTFS import/publication, stable feed paths,
Vehicle Positions, Trip Updates, Alerts, validation workflows, feed health,
readiness workflows, telemetry ingest, and connector boundaries.

That capability is useful for local evaluation and open-source contribution,
but it is not deployment/compliance evidence for a specific agency-owned public
root. Stronger public claims still require retained, public-safe evidence for
the exact claim and deployment scope.

Remaining evidence and service gaps include:

- agency-owned or agency-approved public root proof;
- current DNS, TLS, anonymous fetch, redirect, and URL-permanence evidence;
- current validator-clean records for schedule, Vehicle Positions, Trip
  Updates, and Alerts on the deployed feeds;
- agency-approved license, contact, and source-of-truth website proof;
- consumer or aggregator submission, review, acceptance, ingestion, listing,
  or display records from the named target;
- real device/vendor AVL proof, real-world realtime quality proof, and
  production operations evidence;
- marketplace-vendor packaging, support commitments, and service operations.

External proof tracks are optional and authorization-gated. Evidence packet
generation, readiness exports, and local audits summarize signals and gaps; by
themselves they do not prove compliance, agency approval, consumer acceptance,
vendor compatibility, final-root readiness, production readiness, SLA/uptime,
or production-grade ETA quality.
