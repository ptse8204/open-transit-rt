# Phase 47 Handoff — Self-Hosted Operations Notifications

## Status

Phase 47 is complete for the private local/reference operations notification
summary scope.

## What Changed

- Added `scripts/operations-notify.sh`.
- Added `make operations-notify`.
- Added a private operations notification tutorial.
- Added this Phase 47 reference and handoff.
- Updated docs navigation, roadmap wording, current status, backlog, open
  questions, deployment doctor docs, and operator smoke/support docs.
- Added tests for no-send behavior, source parsing, path privacy, redaction,
  strict mode, source size caps, output size caps, false flags, docs wording,
  and large source summaries.

## Safety Boundary

The feature writes private local drafts only. It does not send notifications,
call webhooks, email operators, call admin routes, run validators, require a
database, require Docker, contact consumers, create evidence, write
`docs/evidence`, edit GTFS, block publish, or change consumer statuses.

It makes no CAL-ITP/Caltrans compliance, consumer acceptance, agency adoption,
hosted SaaS, production readiness, vendor compatibility, or production-grade
ETA claim.

## Output

Default output is `.cache/operations-notify/<timestamp>` and contains exactly:

- `summary.json`
- `summary.md`
- `manifest.json`
- `manifest.md`
- `notification.txt`

`notification.txt` starts with `DRAFT — NOT SENT` and states that the draft was
not sent to webhook, email, consumers, agency, or public service.

## Tests

Added coverage for:

- dry-run no-network/no-send behavior with fake send commands on `PATH`;
- default and strict-mode exit behavior;
- deterministic severity precedence;
- latest timestamp source selection expectations;
- missing, malformed, oversized, hostile, and raw/private source fields;
- symlink and evidence-like path rejection;
- capped next actions, overflow counts, and bounded output sizes;
- exact output file set;
- destination booleans without destination values;
- terminal and manifest privacy;
- roadmap/docs wording boundaries.

Large-input tests record elapsed time only as local engineering diagnostics.
They are not SLA, production capacity, compliance, evidence, or readiness proof.

## Next Recommended Step

Continue self-hosted operations hardening only within private diagnostics unless
retained claim-specific evidence supports a narrower external proof path. Do
not advance consumer tracker statuses or claim compliance/readiness from Phase
47 output.
