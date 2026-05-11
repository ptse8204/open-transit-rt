# Phase 68+ -- Optional Authorized Evidence Tracks

## Status

Checkpoint 000001 documents Phase 68+ as authorization-gated and
blocker-only in the current repository state. No explicit written
authorization, intake packet, or public-safe retained artifacts are available
for agency trial evidence, final public root evidence, consumer submission
evidence, real AVL/vendor evidence, or real-world ETA quality evidence.

This document is scaffolding and blocker documentation only. It does not
collect evidence, contact external parties, write protected evidence paths,
move consumer statuses, or close any evidence gap.

## Scope

Phase 68+ is reserved for optional evidence tracks that may proceed only after
maintainer-approved written intake exists for the exact claim target, allowed
tools, artifact retention plan, redaction plan, and stop conditions.

Until that intake exists, the correct action is to close the track as blocked
or keep it gated without taking external action.

## No-Authorization Finding

The current repository state does not contain explicit written authorization
for Phase 68+ evidence work. Specifically, it does not contain a retained
maintainer intake authorizing an implementation agent to:

- contact an agency, vendor, consumer, aggregator, or public-data operator;
- browse private portals or automate submissions;
- collect retained evidence from external systems;
- verify or change final public root status;
- move consumer or aggregator submission statuses;
- represent real-world AVL/vendor compatibility;
- represent real-world ETA quality.

Because authorization and intake are absent, CP000001 may document blockers
only.

## Track Gates

| Track | Current state | Gate required before work | CP000001 boundary |
| --- | --- | --- | --- |
| Phase 68 -- Agency trial | Blocked/gated | Written agency/operator authorization, named agency scope, public-safe artifact plan, redaction plan, allowed contact/tools, and stop conditions | Do not contact an agency, record agency feedback, or claim trial progress. |
| Phase 69 -- Final public root | Blocked/gated | Written authorization for the exact root, approved DNS/TLS/redirect/fetch/validator evidence plan, retention location, redaction rules, and rollback instructions | Do not fetch or verify a final root, write evidence, or claim final-root proof. |
| Phase 70 -- Consumer submission | Blocked/gated | Written selected-target authorization, target path evidence, exact feed scope, approved operator action plan, public-safe retained artifacts, and status-transition criteria | Do not contact consumers, automate portals, add artifacts, or move statuses. |
| Phase 71 -- AVL/vendor | Blocked/gated | Written vendor/device authorization, payload/data-sharing approval, credential handling plan, redaction plan, conformance criteria, and stop conditions | Do not ingest real vendor data, add compatibility evidence, or claim vendor support. |
| Phase 72 -- ETA quality | Blocked/gated | Written study authorization, real observed-arrival/departure data approval, metric definitions, privacy/redaction plan, and artifact retention rules | Do not collect real-world ETA data, compute retained results, or claim ETA quality. |

## Required Future Intake Fields

Any future authorized evidence track must start with a retained intake record
that includes:

- maintainer authorization date and approver;
- track number and exact claim target;
- named agency, root, target consumer, vendor, device class, or ETA study
  scope, as applicable;
- allowed actions and explicitly forbidden actions;
- allowed tools and accounts;
- external contacts, if any, and who is permitted to contact them;
- public-safe artifact types and approved retention paths;
- redaction requirements for private URLs, credentials, personal data,
  correspondence, screenshots, logs, telemetry, and diagnostics;
- success criteria and blocker criteria;
- status-transition rules, if the track can affect a tracker;
- rollback instructions for any doc/status change;
- stop conditions.

Without those fields, the track remains blocked/gated.

## Protected Paths Not Touched

CP000001 must not change:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/**`
- `db/migrations/**`
- `go.mod`
- `go.sum`
- `docs/current-status.md`
- `docs/roadmap-status.md`
- `docs/handoffs/latest.md`

CP000001 also must not edit code, scripts, migrations, generated evidence,
consumer tracker records, or retained evidence directories.

## Claim Boundaries

Phase 68+ CP000001 may claim only that optional evidence tracks are documented
as authorization-gated and currently blocked because explicit written
authorization and required intake artifacts are absent.

Phase 68+ CP000001 must not claim or imply:

- evidence completion beyond scaffolding/blocker documentation;
- agency adoption, agency approval, agency endorsement, or agency trial
  completion;
- final public root proof;
- consumer submission, review, acceptance, ingestion, listing, display, or
  rejection;
- Caltrans, Cal-ITP, or other compliance status;
- hosted SaaS availability, public launch, paid support, SLA, or production
  readiness;
- real vendor compatibility, hardware certification, or production AVL
  reliability;
- production-grade ETA quality or real-world ETA accuracy.

## Validations

For CP000001, use docs-only validation:

- `git diff --check`
- confirm the protected paths above are unchanged;
- review the changed docs for unsupported evidence, compliance, adoption,
  final-root, consumer-status, hosted-service, vendor, hardware, AVL,
  production ETA, SLA, or public-launch claims.

No live feed validation, external browsing, portal automation, consumer
submission, vendor test, or ETA study is authorized by this checkpoint.

## Rollback

If CP000001 needs to be rolled back, revert only:

- `docs/phase-68-plus-optional-authorized-evidence-tracks.md`
- `docs/roadmaps/agency-first-connector-platform/02-phases-and-checkpoints.md`
- `docs/roadmaps/agency-first-connector-platform/phase-prompts/phase-68-plus-optional-evidence-tracks.md`

Do not change protected evidence paths, consumer tracker records, migrations,
Go module files, code, scripts, status closeout files, or handoff files as
part of this checkpoint rollback.
