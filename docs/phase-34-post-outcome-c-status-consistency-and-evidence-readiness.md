# Phase 34 — Post-Outcome-C Status Consistency And Evidence Readiness

## Status

Complete for docs-only status consistency and evidence-readiness scope.

See `docs/handoffs/phase-34.md` and `docs/handoffs/latest.md`.

This document is retained as the Phase 34 scope/closure reference.

## Purpose

Phase 33 is now complete as Outcome C. The repo has local/pilot public static GTFS dataset handling evidence for the LA Metro Bus public GTFS run.

Phase 34 made the repo easier for future Codex agents, maintainers, and agency
evaluators to use after Outcome C. It fixed stale status docs, preserved the
new truth state, and prepared the next real retained-evidence path without
inventing external evidence.

## Scope

Phase 34 was a docs-only documentation, evidence-readiness, and repeatability
phase.

The closed scope covered:

1. post-Outcome-C status consistency;
2. public-GTFS local/pilot repeatability guidance;
3. final-root operator request packaging;
4. static validator tooling blocker clarification;
5. forward roadmap selection.

This Phase 34 implementation added no scripts, Makefile targets, runtime code,
schema changes, migrations, APIs, consumer tracker changes, final-root evidence
packets, or target artifacts.

## Continuation Reading

Future continuation read list:

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
- original static GTFS validation was attempted but blocked because
  `/usr/bin/java` could not locate a Java runtime;
- a later post-Phase-34 static validator retry used Homebrew Java 17 against the
  already-fetched schedule ZIP and reported process exit code `0`, system error
  count `0`, and 3 warning notices;
- GTFS-RT validation passed for empty valid protobuf feeds;
- telemetry simulator dry-run printed synthetic payloads and sent no telemetry;
- public/admin/private boundary checks were recorded.

## Non-goals

Phase 34 added no claim or implication of:

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

## Closed Implementation Reference

### 1. Patch stale status and roadmap docs

These files were updated so they consistently reflect Phase 33 Outcome C:

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

Closure corrections:

- `docs/current-status.md` no longer opens by calling the project an early-stage starter without qualification.
- `docs/roadmap-status.md` no longer says public-GTFS evidence is only an attempted blocker or missing; it says Outcome C exists and is limited to local/pilot public static GTFS handling.
- `docs/track-b-productization-roadmap.md` no longer recommends Phase 32 as next. It shows Phases 22-33 as closed and Phase 34 or the next evidence fork as the current recommendation.
- `docs/repo-gaps.md` was refreshed so it no longer lists already-completed starter scaffolding as current missing work.
- `docs/phase-plan.md` points future agents to the post-Outcome-C roadmap and latest handoff.
- `docs/README.md` labels the public-GTFS evidence packet as Outcome C evidence rather than merely an attempt.
- `docs/phase-34-post-outcome-c-status-consistency-and-evidence-readiness.md`
  states Phase 34 is docs-only and does not allow scripts or Makefile targets.
- `docs/future-roadmap-post-outcome-c.md` marks Phase 34 complete and makes the
  next path a retained-evidence fork.

### 2. Add a final-root operator request package

Added:

```text
docs/final-root-operator-request.md
```

This plain-language operator-facing request package explains:

- why a final public feed root matters;
- acceptable examples of agency-owned or agency-approved roots;
- approval evidence required;
- DNS/TLS/redirect proof required;
- all five final public feed URLs required;
- validator evidence required;
- redaction and security rules;
- what not to send;
- how DuckDNS OCI pilot evidence remains labeled;
- what happens after final-root proof exists.

This file is a request package, not evidence.

### 3. Add a repeatable public-GTFS local/pilot guide

Added or updated:

```text
docs/tutorials/public-gtfs-local-pilot.md
```

This guide helps future maintainers repeat the Phase 33 style run without
implying agency approval.

It covers:

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

The Phase 33 Outcome C packet records that the original static GTFS validation
attempt did not execute because `/usr/bin/java` could not locate a Java runtime.
After Phase 34, a retry using Homebrew Java 17 executed the pinned static
validator against the already-fetched schedule ZIP and reported process exit
code `0`, system error count `0`, and 3 warning notices.

Future run guidance requires one of these paths to be explicit:

- install/check Java before static validation; or
- run the static validator in a known Java-capable environment; or
- keep the blocker explicitly documented when Java is unavailable.

Do not convert an executed validator retry with warnings into a
validator-clean, no-warning, compliance, or production-readiness claim.

### 5. Preserve external-evidence boundaries

These remained unchanged because no real retained external evidence existed for
them:

```text
docs/evidence/consumer-submissions/status.json
consumer target records
target-specific consumer artifact directories
final-root evidence packets
OCI pilot final-root wording
```

## Closure Criteria

Phase 34 closed after:

- stale status/roadmap docs no longer contradicted Phase 33 Outcome C;
- a final-root operator request package exists;
- a public-GTFS local/pilot repeatability guide exists;
- the original Phase 33 static validator blocker and later Homebrew Java 17
  retry are clearly explained and not overclaimed;
- `docs/handoffs/latest.md` identified the next retained-evidence fork;
- no consumer statuses were advanced;
- no final-root evidence was fabricated;
- no agency/consumer/compliance/production claims were added;
- required checks passed or blockers were documented.

## Closure Checks

Closure checks:

```bash
make validate
make test
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
```

All seven consumer targets remained `prepared`, and a targeted wording scan
reviewed forbidden positive claims.

Additional checks listed for future phases when touched surfaces require them:

```bash
make realtime-quality
make smoke
make test-integration
docker compose -f deploy/docker-compose.yml config
```

Run local app or public-GTFS pilot commands only in future phases that actually
change those flows.

## Handoff Record

Added:

```text
docs/handoffs/phase-34.md
```

The handoff states:

- what files changed;
- which stale docs were corrected;
- whether any script/Makefile behavior changed;
- what checks ran;
- what remained blocked;
- what evidence was or was not created;
- which evidence fork is recommended next.
