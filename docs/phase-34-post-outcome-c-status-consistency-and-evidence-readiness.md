# Phase 34 — Post-Outcome-C Status Consistency And Evidence Readiness

## Purpose

Phase 33 is now complete as Outcome C. The repo has local/pilot public static GTFS dataset handling evidence for the LA Metro Bus public GTFS run.

Phase 34 should make the repo easier for future Codex agents, maintainers, and agency evaluators to use after Outcome C. It should fix stale status docs, preserve the new truth state, and prepare the next real retained-evidence path without inventing external evidence.

## Scope

Phase 34 is a docs-only documentation, evidence-readiness, and repeatability
phase.

It should cover:

1. post-Outcome-C status consistency;
2. public-GTFS local/pilot repeatability guidance;
3. final-root operator request packaging;
4. static validator tooling blocker clarification;
5. forward roadmap selection.

For this Phase 34 implementation, do not add scripts, Makefile targets, runtime
code, schema changes, migrations, APIs, consumer tracker changes, final-root
evidence packets, or target artifacts.

## Required first reading

A fresh Codex agent should read:

```text
AGENTS.md
README.md
docs/current-status.md
docs/handoffs/latest.md
docs/handoffs/phase-33.md
docs/phase-33-public-gtfs-local-pilot-evidence.md
docs/evidence/captured/public-gtfs-local-pilot/2026-05-06/README.md
docs/evidence/captured/public-gtfs-local-pilot/2026-05-06/command-log-inventory-2026-05-06.md
docs/evidence/captured/public-gtfs-local-pilot/2026-05-06/retained-summaries-2026-05-06.md
docs/roadmap-status.md
docs/track-b-productization-roadmap.md
docs/repo-gaps.md
docs/phase-plan.md
docs/evidence/consumer-submissions/status.json
docs/agency-owned-domain-readiness.md
docs/compliance-evidence-checklist.md
docs/evidence/redaction-policy.md
```

## Current baseline

Phase 33 proves only:

```text
Open Transit RT can handle a current public static GTFS dataset in a recorded local/pilot environment.
```

The Outcome C packet proves:

- raw LA Metro Bus GTFS was downloaded to ignored `.cache/` storage;
- source URL, access date, license/terms URL, and checksum were recorded;
- the public GTFS was imported through `cmd/gtfs-import -agency-id LACMTA -timeout 15m`;
- the local schedule feed was fetched;
- the fetched schedule was verified as the imported LA Metro public GTFS rather than the repo sample feed;
- all five local public paths were fetched;
- realtime protobuf endpoints were valid empty publications;
- static GTFS validation was attempted but blocked by missing Java;
- GTFS-RT validation passed for empty valid protobuf feeds;
- telemetry simulator dry-run printed synthetic payloads and sent no telemetry;
- public/admin/private boundary checks were recorded.

## Non-goals

Phase 34 must not claim or imply:

- agency adoption;
- agency endorsement;
- agency approval;
- official agency feed status;
- agency-owned final-root proof;
- consumer submission;
- consumer review;
- consumer acceptance;
- consumer ingestion, listing, or display;
- Caltrans/CAL-ITP compliance;
- hosted SaaS availability;
- paid support or SLA coverage;
- production readiness;
- production multi-tenant hosting;
- real vendor AVL compatibility;
- real LA Metro realtime data;
- real-world ETA accuracy;
- production-grade ETA quality;
- public launch completion.

## Required implementation

### 1. Patch stale status and roadmap docs

Update these files so they consistently reflect Phase 33 Outcome C:

```text
docs/current-status.md
docs/roadmap-status.md
docs/track-b-productization-roadmap.md
docs/repo-gaps.md
docs/phase-plan.md
docs/README.md
docs/phase-34-post-outcome-c-status-consistency-and-evidence-readiness.md
docs/future-roadmap-post-outcome-c.md
```

Minimum expected corrections:

- `docs/current-status.md` should no longer open by calling the project an early-stage starter without qualification.
- `docs/roadmap-status.md` should not say public-GTFS evidence is still only an attempted blocker or missing; it should say Outcome C exists and is limited to local/pilot public static GTFS handling.
- `docs/track-b-productization-roadmap.md` should not recommend Phase 32 as next. It should show Phases 22-33 as closed and Phase 34 or the next evidence fork as the current recommendation.
- `docs/repo-gaps.md` should be refreshed, retired, or clearly marked historical. It should not list already-completed starter scaffolding as current missing work.
- `docs/phase-plan.md` should point future agents to the post-Outcome-C roadmap and latest handoff.
- `docs/README.md` should label the public-GTFS evidence packet as Outcome C evidence rather than merely an attempt.
- `docs/phase-34-post-outcome-c-status-consistency-and-evidence-readiness.md`
  should state Phase 34 is docs-only and should not allow scripts or Makefile
  targets.
- `docs/future-roadmap-post-outcome-c.md` should mark Phase 34 complete once
  closed and make the next path a retained-evidence fork.

### 2. Add a final-root operator request package

Add:

```text
docs/final-root-operator-request.md
```

This should be plain-language and operator-facing. It should explain:

- why a final public feed root matters;
- acceptable examples of agency-owned or agency-approved roots;
- approval evidence required;
- DNS/TLS/redirect proof required;
- all five final public feed URLs required;
- validator evidence required;
- redaction and security rules;
- what not to send;
- how DuckDNS OCI pilot evidence should remain labeled;
- what happens after final-root proof exists.

This file is a request package, not evidence.

### 3. Add a repeatable public-GTFS local/pilot guide

Add or update:

```text
docs/tutorials/public-gtfs-local-pilot.md
```

This guide should help future maintainers repeat the Phase 33 style run without implying agency approval.

It should cover:

- choosing a public GTFS dataset;
- checking source/catalog/license facts;
- downloading raw GTFS only to ignored storage;
- recording checksum;
- creating local-only agency/admin setup for `agency_id` matching;
- using `cmd/gtfs-import` with `-timeout`;
- avoiding the default demo sample feed when testing a public agency dataset;
- fetching the five public paths;
- verifying fetched `schedule.zip` as the imported public dataset;
- running static and realtime validators or documenting blockers;
- running telemetry dry-run only;
- checking admin/private route boundaries;
- retaining public-safe summaries;
- preserving claim boundaries.

### 4. Clarify static validator environment blocker

The Phase 33 Outcome C packet records that static GTFS validation did not execute because Java was unavailable. Phase 34 should make the next action clearer by documenting one of these paths:

- install/check Java before static validation; or
- run the static validator in a known Java-capable environment; or
- keep the blocker explicitly documented when Java is unavailable.

Do not retroactively claim that the static GTFS validator passed unless a real retained no-error static validator record exists.

### 5. Preserve external-evidence boundaries

Do not change these unless real retained external evidence exists:

```text
docs/evidence/consumer-submissions/status.json
consumer target records
target-specific consumer artifact directories
final-root evidence packets
OCI pilot final-root wording
```

## Acceptance criteria

Phase 34 is complete when:

- stale status/roadmap docs no longer contradict Phase 33 Outcome C;
- a final-root operator request package exists;
- a public-GTFS local/pilot repeatability guide exists;
- the Phase 33 static validator blocker is clearly explained and not overclaimed;
- `docs/handoffs/latest.md` identifies the next retained-evidence fork;
- no consumer statuses are advanced;
- no final-root evidence is fabricated;
- no agency/consumer/compliance/production claims are added;
- required checks pass or blockers are documented.

## Required checks

Run:

```bash
make validate
make test
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
```

Also confirm all seven consumer targets remain `prepared` and run a targeted
wording scan for forbidden positive claims.

Run when touched surfaces require them:

```bash
make realtime-quality
make smoke
make test-integration
docker compose -f deploy/docker-compose.yml config
```

Run local app or public-GTFS pilot commands only if the implementation actually changes those flows.

## Handoff requirement

Add:

```text
docs/handoffs/phase-34.md
```

The handoff should state:

- what files changed;
- which stale docs were corrected;
- whether any script/Makefile behavior changed;
- what checks ran;
- what remained blocked;
- what evidence was or was not created;
- which evidence fork should be pursued next.
