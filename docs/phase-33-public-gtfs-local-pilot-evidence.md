# Phase 33 — Public GTFS Local/Pilot Evidence

## Status

Complete as Outcome B — attempted public-GTFS run blocked.

Phase 33 is an intermediate development evidence phase. It is meant to prove,
or prepare a repeatable path to prove, that Open Transit RT can ingest, publish,
validate, and review a real current public agency GTFS dataset in a local or
pilot environment.

The 2026-05-06 LA Metro Bus GTFS local attempt is documented at
`docs/evidence/captured/public-gtfs-local-pilot/2026-05-06/`. It blocked during
repo-supported import/publish when the large `stop_times.txt` dataset exceeded
the current importer context. Outcome C evidence was not completed.

Do not describe Phase 33 as agency adoption, agency endorsement, agency
approval, official agency feed status, agency-owned final-root proof, consumer
submission, consumer review, consumer acceptance, consumer ingestion/listing/
display, Caltrans/CAL-ITP compliance, hosted SaaS availability, production
readiness, production multi-tenant hosting, real vendor AVL compatibility,
real-world ETA accuracy, or production-grade ETA quality.

## Closure Outcomes

Record exactly one outcome.

### Outcome A — Template-Only Closure

Use Outcome A when no public-GTFS run is attempted.

Required:

- Phase 33 docs added.
- Evidence packet template added.
- Status/navigation docs updated.
- No fake run evidence added.
- Final-root blocker unchanged.
- Consumer statuses unchanged.
- `make validate`, `make test`, and `git diff --check` passed or blockers are
  documented.

Outcome A is template-only. Do not call it evidence completed.

### Outcome B — Attempted-Run Blocked Closure

Use Outcome B when a core public-GTFS run step blocks.

Required:

- Everything from Outcome A.
- Attempted public-GTFS run documented.
- Blocker recorded truthfully.
- No fake logs, validator outputs, or feed evidence added.

Core steps whose failure means Outcome B:

- public GTFS download;
- source checksum;
- import through the repo-supported flow;
- local/pilot app startup;
- published `/public/gtfs/schedule.zip` fetch;
- proof that the fetched schedule is the imported public GTFS, not the repo
  sample feed;
- five-path public fetch.

### Outcome C — Public-GTFS Local/Pilot Evidence Closure

Use Outcome C only when all core import/publish/fetch evidence exists.

Required:

- LA Metro GTFS or approved fallback downloaded only into an ignored path.
- Source URL, access date, license/terms URL, and source checksum recorded.
- Import completed through the repo-supported flow.
- `/public/gtfs/schedule.zip` fetched from the local/pilot root.
- The fetched schedule ZIP is unzipped in an ignored temp path and verified as
  the imported public GTFS, not the sample feed, using public-safe summaries of
  `agency.txt`, `routes.txt`, and service-date coverage from `calendar.txt` or
  `calendar_dates.txt`.
- All five public paths are fetched:
  - `/public/feeds.json`
  - `/public/gtfs/schedule.zip`
  - `/public/gtfsrt/vehicle_positions.pb`
  - `/public/gtfsrt/trip_updates.pb`
  - `/public/gtfsrt/alerts.pb`
- A dated packet is created under
  `docs/evidence/captured/public-gtfs-local-pilot/YYYY-MM-DD/`.
- Claim boundaries are explicit.

Validators, telemetry simulator/dry-run, and admin/private/debug route checks
should be run if available. If unavailable, document the blocker. Do not claim
that blocked checks passed.

Outcome C proves local/pilot handling of a real public static GTFS dataset. It
does not prove final-root readiness, consumer acceptance, compliance, official
agency status, production realtime quality, or real LA Metro realtime data.

## Evidence Boundaries

Keep these evidence categories separate:

- Public GTFS local/pilot evidence proves real public dataset handling in the
  recorded local or pilot environment.
- Agency-owned final-root evidence proves official/stable agency-domain
  readiness only when retained owner/approval, DNS/TLS, public fetch, validator,
  and checksum evidence exists for that final root.
- Consumer evidence proves only what a retained, redacted, target-originated
  artifact supports for that specific target and feed scope.

The post-Phase-32 final-root blocker remains unchanged unless a future run has
an agency-owned or agency-approved final public feed root plus retained approval
evidence. The DuckDNS OCI pilot remains hosted/operator pilot evidence only.

## Attempted Run Procedure

Preferred source: LA Metro Bus GTFS.

Record catalog facts with a checked-at timestamp. Use Mobility Database
primarily for official-feed, route-count, service-range, producer URL, and
dataset-size facts. Use Transitland secondarily for current URL, last fetch,
license URL, and attribution fields. Treat catalog facts as time-sensitive, not
permanent.

Store the raw GTFS ZIP under `.cache/` or another ignored path. Do not commit
the raw GTFS ZIP by default. If a future maintainer wants to commit the ZIP,
first review license/terms and intentionally accept the file as public-safe.

Core run steps:

1. Download the public GTFS ZIP to an ignored path.
2. Record source URL, access date, license/terms URL, and SHA-256 checksum.
3. Start the repo-supported local or pilot app flow.
4. Import the public GTFS through the repo-supported GTFS import flow.
5. Fetch `/public/gtfs/schedule.zip` from the local or pilot root.
6. Unzip the fetched schedule in an ignored temp path.
7. Summarize `agency.txt`, `routes.txt`, and service-date coverage to prove
   the active published schedule is the imported public GTFS, not the sample
   feed.
8. Fetch all five public paths.
9. Run validators if available; otherwise document the blocker.
10. Run telemetry simulator or dry-run if supported; otherwise document the
    blocker.
11. Check admin/private/debug route boundaries if available; otherwise document
    the blocker.

GTFS-RT fetches prove endpoint availability and protobuf publication only.
They may be empty, simulated, or withheld. Do not imply real LA Metro realtime
data, real LA Metro vehicle telemetry, real LA Metro alerts operation,
real-world ETA accuracy, or production-grade Trip Updates.

## Required Checks

Run after edits:

```bash
make validate
make test
git diff --check
```

If a public-GTFS run is attempted, record local app, import, fetch, validator,
evidence-packet, and teardown commands run or blocked.

## Files Not To Change

Do not change unless real external evidence exists:

- `docs/evidence/consumer-submissions/status.json`
- consumer target records
- target-specific artifact directories
- final-root evidence
- OCI pilot final-root wording

## Explicit Non-Goals

Phase 33 does not:

- claim agency adoption;
- claim agency endorsement;
- claim agency approval;
- claim official agency feed status;
- claim agency-owned final-root proof;
- claim consumer submission, review, acceptance, ingestion, listing, or display;
- claim Caltrans/CAL-ITP compliance;
- claim hosted SaaS availability;
- claim production readiness or production multi-tenant hosting;
- claim real vendor AVL compatibility;
- claim real-world ETA accuracy;
- claim production-grade ETA quality.
