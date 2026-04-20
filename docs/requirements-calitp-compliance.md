# Requirements — Cal-ITP / Caltrans Technical Compliance and Marketplace Gaps

These requirements formalize the gap between the current starter plan and a system that is technically compliant with Caltrans transit data guidelines, plus the additional gap between technical compliance and marketplace-vendor equivalence.

---

## Scope clarification

There are two separate targets:

1. **Technical compliance target**
   - Static GTFS and GTFS-Realtime feeds meet Caltrans data-guideline expectations.
   - Feeds are public, stable, validator-clean, openly licensed, and acceptable to major trip-planning consumers.

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
Meet expectations for open publication and clear licensing.

### Required behaviors
- Associate each published feed with an explicit open license.
- Display license language on the provider-hosted landing page.
- Provide discoverability metadata:
  - feed URLs
  - provider name
  - technical contact
  - update timestamp
  - license
- Support publishing or registration workflows for:
  - Mobility Database
  - transit.land

### Acceptance criteria
- A public page exists for the feed set with license and contact info.
- Feed metadata is sufficient to register with aggregators.
- License is visible without logging in.

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

## RQ-4G — Compliance dashboard and scorecard

### Goal
The product must show compliance posture transparently.

### Required behaviors
- Per-agency compliance dashboard with sections:
  - GTFS Schedule
  - Vehicle Positions
  - Trip Updates
  - Alerts
  - validation status
  - license/discoverability
  - consumer ingestion status
- Show red/yellow/green state for:
  - not implemented
  - implemented but unhealthy
  - implemented and compliant
- Exportable compliance report for agency review

### Acceptance criteria
- A user can see exactly why an agency is not yet compliant.
- A user can download a compliance checklist or status report.

---

## Current gap summary

The current starter plan is **not yet compliant** because it does not yet fully provide:
- Trip Updates
- Alerts
- full public/stable publication workflow
- validator-clean outputs across all feeds
- open-license/discoverability workflow
- consumer ingestion workflow
- marketplace-vendor packaging

It does provide a foundation that can evolve into compliance if the requirements in this document are implemented.
