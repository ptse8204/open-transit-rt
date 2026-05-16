# Phase 120 -- GTFS-RT Feed Usefulness And Reliability V2

## Goal

Improve GTFS-RT feed usefulness and reliability review for Vehicle Positions,
Trip Updates, and Alerts without making production readiness, SLA, compliance,
consumer acceptance, vendor compatibility, hardware certification, or
production-grade ETA quality claims.

## Current Repo Context

- Phase 115 published public `v0.1.0-rc.1` for local/self-hosted evaluation.
- Phase 116 recorded public download replay and the source-archive `make check`
  limitation for the already-published rc1 archives.
- Phase 117 verified public fresh-clone install confidence from the rc1 tag.
- Phase 118 completed Web Design Skill UX validation.
- Phase 119 aligned public docs and quickstarts to the actual rc1 release and
  install-confidence state.
- Existing feed packages already emit Vehicle Positions, Trip Updates, Alerts,
  and schedule metadata; existing Operations Console routes already show
  private feed health and realtime review state.

## Scope

- Add or improve local/synthetic feed usefulness and reliability diagnostics for
  GTFS-RT outputs.
- Prefer structured helper logic, tests, and bounded Operations Console/readme
  surfaces over ad hoc wording.
- Keep Vehicle Positions first, and keep Trip Updates pluggable and
  fail-closed when unsupported.
- Preserve all protected evidence paths and prepared-only consumer tracker
  statuses.

## Protected Paths

Do not modify, reformat, delete, stage, or generate files under:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/**`
- `docs/evidence/consumer-submissions/artifacts/**`
- `docs/evidence/consumer-submissions/packets/**`

The consumer tracker must remain exactly seven targets in order and all
`prepared`.

## Deliverables

- A bounded implementation or audit artifact that improves GTFS-RT feed
  usefulness/reliability review for local operators.
- Focused tests or validation covering the changed behavior.
- `docs/handoffs/phase-120.md`
- Source-of-truth status updates for Phase 120 closeout.

## Implementation Plan

1. Add this Phase 120 plan and commit checkpoint 000001.
2. Inspect feed/realtimequality/Operations Console code and implement one
   conservative, test-backed usefulness/reliability improvement.
3. Run relevant GTFS-RT, connector, claim-boundary, and baseline validation;
   patch repo-caused failures.
4. Close Phase 120 with handoff/status docs and continue immediately to
   Phase 121.

## Checkpoint Plan

- `Phase 120 -- Checkpoint 000001: add gtfs rt feed usefulness and reliability v2 plan`
- `Phase 120 -- Checkpoint 000002: implement or audit primary scoped work`
- `Phase 120 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 120 -- Checkpoint 000004: close gtfs rt feed usefulness and reliability v2 review`

## Checkpoint Report -- 000001

Checkpoint:
Phase 120 -- Checkpoint 000001: add GTFS-RT feed usefulness and reliability
v2 plan.

Goal status:
Active. Phase 119 is closed and Phase 120 has started.

Sub-agents used or simulated:
The environment refused an additional GTFS-RT explorer because the agent
thread limit was reached. Context / Repo Truth, Planning, Implementation, QA,
GTFS-RT Domain, Claim-Boundary, Security/Auth, Data/Migration,
Documentation / IA, Web Design Skill, Release, and Install Confidence roles
are simulated by the Master Agent for this checkpoint.

Changed files:
`docs/phase-120-gtfs-rt-feed-usefulness-and-reliability-v2.md`.

Validation run:
Initial inspection reviewed the Phase 120 roadmap prompt, validation
boundaries, master/sub-agent workflow, feed/realtimequality packages,
Operations Console route files, scripts, and Makefile GTFS-RT targets.

Blocked checks:
Implementation, tests, connector validation, and closeout validation are
scheduled for later Phase 120 checkpoints.

Protected path status:
No protected evidence path is part of the plan. The plan forbids protected
path writes.

Consumer tracker status:
The consumer tracker is not part of the plan. The seven targets must remain in
order and exactly `prepared`.

Claim-boundary status:
The plan explicitly forbids stable release readiness, production readiness,
compliance, adoption, agency approval, consumer acceptance, consumer
ingestion/listing/display, final-root readiness, hosted service availability,
paid support, SLA/uptime, vendor compatibility, hardware certification,
production AVL reliability, production-grade ETA quality, and real-world ETA
accuracy claims.

Security/auth status:
The plan does not change route auth, CSRF behavior, credential handling, token
handling, public exposure, private payload handling, or operator command
behavior without focused review.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change is
planned.

Release/publication status:
The public rc1 prerelease remains published. Phase 120 does not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
Phase 118 Web Design Skill artifact remains complete. Phase 120 does not plan
visual UX changes unless implementation inspection shows a tightly scoped
Operations Console review improvement.

Master review:
Approved. The plan keeps Phase 120 inside local/synthetic GTFS-RT usefulness
and reliability review, preserving all evidence and claim boundaries.

Required edits:
Commit checkpoint 000001, then implement or audit the scoped GTFS-RT
improvement.

Decision:
Proceed to checkpoint 000001 validation and commit.

Next checkpoint:
Phase 120 -- Checkpoint 000002: implement or audit primary scoped work.

## Checkpoint Report -- 000002

Checkpoint:
Phase 120 -- Checkpoint 000002: implement or audit primary scoped work.

Goal status:
Active. Phase 120 implemented a bounded GTFS-RT Vehicle Positions usefulness
and reliability review improvement.

Sub-agents used or simulated:
The agent thread limit still prevents new real sub-agents. Context / Repo
Truth, Implementation, QA, GTFS-RT Domain, Claim-Boundary, Security/Auth,
Data/Migration, Documentation / IA, Web Design Skill, Release, and Install
Confidence roles were simulated by the Master Agent.

Changed files:
`internal/feed/vehicle_positions.go`;
`internal/feed/vehicle_positions_health.go`;
`internal/feed/vehicle_positions_test.go`;
`docs/phase-120-gtfs-rt-feed-usefulness-and-reliability-v2.md`.

Validation run:
Passed:

- `gofmt -w internal/feed/vehicle_positions.go internal/feed/vehicle_positions_health.go internal/feed/vehicle_positions_test.go`
- `go test ./internal/feed/...`
- `git diff --check`
- protected-path git status check

Blocked checks:
Full Phase 120 baseline validation, connector/GTFS-RT validation, and claim
audits are scheduled for checkpoint 000003.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.

Consumer tracker status:
The tracker was not edited. The exact seven consumer targets remain in order
and all remain `prepared`.

Claim-boundary status:
The implementation adds only private/local diagnostic aggregates. It makes no
production readiness, SLA, uptime, compliance, adoption, consumer acceptance,
vendor compatibility, hardware certification, production AVL reliability,
production-grade ETA quality, or real-world ETA accuracy claim.

Security/auth status:
No route auth, CSRF, credential, token, public exposure, private payload, or
operator command behavior changed. New aggregates intentionally omit raw
payloads, URLs, credentials, private paths, tokens, and vehicle payload JSON.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change was
added. Existing `details_json` health payloads gain bounded aggregate keys for
new snapshots only.

Release/publication status:
The public rc1 prerelease remains published. Phase 120 does not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed. The changed
code is post-rc1 current-source hardening and is not part of the published rc1
tag.

Web design skill status:
Phase 118 Web Design Skill artifact remains complete. CP000002 did not change
visual UI templates.

Master review:
Approved. The change improves Vehicle Positions review by adding safe
aggregate counts for protobuf inclusion, published trip descriptors,
stale/suppressed vehicles, assignment telemetry mismatches, and trip
descriptor omission reasons without changing public feed semantics.

Required edits:
Commit checkpoint 000002, then run Phase 120 validation and patch any
repo-caused failures.

Decision:
Proceed to checkpoint 000002 commit and checkpoint 000003 validation.

Next checkpoint:
Phase 120 -- Checkpoint 000003: run validation and patch required gaps.

## Checkpoint Report -- 000003

Checkpoint:
Phase 120 -- Checkpoint 000003: run validation and patch required gaps.

Goal status:
Active. Phase 120 implementation validation passed.

Sub-agents used or simulated:
Context / Repo Truth, QA, GTFS-RT Domain, Connector, Claim-Boundary,
Security/Auth, Data/Migration, Documentation / IA, Web Design Skill, Release,
and Install Confidence roles were simulated by the Master Agent because the
agent thread limit prevents new real sub-agents.

Changed files:
`docs/phase-120-gtfs-rt-feed-usefulness-and-reliability-v2.md`.

Validation run:
Passed:

- `git status --short`
- `git diff --check`
- `make check`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`
- `make external-connection-check`
- `make adapter-conformance`
- `make test-connector-examples`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json`
- `scripts/check-consumer-tracker.sh`
- protected-path git status check

Blocked checks:
No Phase 120 validation blocker remains.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.

Consumer tracker status:
The tracker was not edited. The exact seven consumer targets remain in order
and all remain `prepared`.

Claim-boundary status:
Product acceptance and final claim audits passed. Phase 120 remains bounded to
local/private GTFS-RT usefulness and reliability diagnostics and makes no
stable release readiness, production readiness, compliance, adoption, agency
approval, consumer acceptance, consumer ingestion/listing/display, final-root
readiness, hosted service, paid support, SLA/uptime, vendor compatibility,
hardware certification, production AVL reliability, production-grade ETA
quality, or real-world ETA accuracy claim.

Security/auth status:
No route auth, CSRF, credential handling, token handling, public exposure,
private payload handling, or operator command behavior changed.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change was
added.

Release/publication status:
The public rc1 prerelease remains published. Phase 120 did not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed. Phase 120
current-source hardening is not part of the published rc1 tag.

Web design skill status:
Phase 118 Web Design Skill artifact remains complete. Phase 120 did not change
visual UI templates.

Master review:
Approved. Validation found no required patch after the Vehicle Positions
aggregate review-summary implementation.

Required edits:
Commit checkpoint 000003, then close Phase 120 with handoff and status docs.

Decision:
Proceed to checkpoint 000003 commit and Phase 120 closeout.

Next checkpoint:
Phase 120 -- Checkpoint 000004: close GTFS-RT feed usefulness and reliability
v2 review.

## Checkpoint Report -- 000004

Checkpoint:
Phase 120 -- Checkpoint 000004: close GTFS-RT feed usefulness and reliability
v2 review.

Goal status:
Active. Phase 120 is closed and the goal continues to Phase 121.

Sub-agents used or simulated:
Context / Repo Truth, Planning, Implementation, QA, GTFS-RT Domain, Connector,
Claim-Boundary, Security/Auth, Data/Migration, Documentation / IA, Web Design
Skill, Release, and Install Confidence roles were simulated by the Master
Agent because the agent thread limit prevented new real sub-agents.

Changed files:
`docs/handoffs/phase-120.md`; `docs/handoffs/latest.md`;
`docs/current-status.md`; `docs/roadmap-status.md`;
`docs/open-transit-rt-master-planner-remaining-work.md`;
`docs/phase-120-gtfs-rt-feed-usefulness-and-reliability-v2.md`.

Validation run:
Full Phase 120 validation passed before closeout docs. Focused closeout
validation passed after closeout docs: `git diff --check`, `make check`,
`make audit-product-acceptance`, `make audit-final-claim-review`,
`scripts/check-consumer-tracker.sh`, and protected-path git status.

Blocked checks:
No Phase 120 check remains blocked.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.

Consumer tracker status:
The tracker was not edited. The exact seven consumer targets remain in order
and all remain `prepared`.

Claim-boundary status:
Phase 120 remains bounded to local/private GTFS-RT usefulness and reliability
diagnostics and makes no stronger public claim.

Security/auth status:
No application security behavior changed.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change was added.

Release/publication status:
The public rc1 prerelease remains published. No release action was taken.

Install confidence status:
Public fresh-clone rc1 install confidence remains passed.

Web design skill status:
Phase 118 Web Design Skill artifact remains complete.

Master review:
Approved. Phase 120 closes with redaction-safe Vehicle Positions aggregate
review diagnostics and full validation passed.

Required edits:
Commit checkpoint 000004, then continue directly to Phase 121.

Decision:
Proceed to checkpoint 000004 commit and continue to Phase 121.

Next checkpoint:
Phase 121 -- Checkpoint 000001: add GTFS-RT interoperability conformance
harness plan.
