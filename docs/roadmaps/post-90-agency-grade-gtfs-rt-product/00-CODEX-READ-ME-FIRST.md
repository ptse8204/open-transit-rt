# Codex Read Me First — Post-90 Roadmap

## Read first

```text
AGENTS.md
README.md
docs/current-status.md
docs/handoffs/latest.md
docs/handoffs/phase-90.md
docs/phase-90-control-plane-final-status.md
docs/phase-89-rc1-gate-results.md
docs/roadmap-status.md
docs/evidence/redaction-policy.md
docs/evidence/consumer-submissions/status.json
```

## Current truth

- Phases 75-90 are complete.
- Release readiness remains `needs_review`.
- Evidence/adoption/compliance tracks remain authorization-gated.
- This roadmap starts only when a maintainer explicitly authorizes it.
- The autonomous Phase 91-110 authorization uses
  `CODEX-KICKOFF-AUTONOMOUS-PHASE-91-TO-110.md`.

## Authorized autonomous sequence

When the autonomous kickoff is the active prompt, continue phase by phase:

```text
Phase 91 -> Phase 92 -> ... -> Phase 110
```

Do not stop after Phase 91, after roadmap-pack reconciliation, after a
documentation-only phase, after a validation run, or after a `needs_review`
diagnostic unless a hard-stop condition cannot be converted into a safe
blocker-only closeout.

## Protected paths

Do not modify or generate files under:

```text
docs/evidence/captured/**
docs/evidence/consumer-submissions/status.json
docs/evidence/consumer-submissions/current/**
docs/evidence/consumer-submissions/artifacts/**
docs/evidence/consumer-submissions/packets/**
```

## Consumer tracker

All seven targets must remain exactly `prepared`:

```text
Google Maps
Apple Maps
Transit App
Bing Maps
Moovit
Mobility Database
transit.land
```

## Forbidden claims

No CAL-ITP/Caltrans compliance, agency adoption/approval, consumer submission,
consumer acceptance, consumer ingestion/listing/display, final-root readiness,
hosted SaaS, production readiness, vendor compatibility, hardware certification,
SLA/uptime, production-grade ETA, real-world ETA accuracy, public launch, or
release readiness.

## Required master/sub-agent mode

Use the operating model in `03-master-subagent-operating-manual.md`. If real
sub-agents are unavailable, simulate the roles in labeled sections.
