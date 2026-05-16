# Phase 120 Handoff -- GTFS-RT Feed Usefulness And Reliability V2

## Status

Phase 120 is complete for bounded GTFS-RT feed usefulness and reliability
review V2.

The implementation added safe aggregate Vehicle Positions review summaries for
local/private diagnostics:

- protobuf entity count;
- published trip descriptor count;
- trip descriptor omission counts by safe reason;
- stale and suppressed vehicle counts;
- unmatched vehicle count;
- assignment telemetry mismatch count;
- truncation and telemetry-read counts.

These aggregates are available from `VehiclePositionsSnapshot.ReviewSummary()`,
Vehicle Positions debug JSON, and persisted private feed-health `details_json`
for new Vehicle Positions health snapshots.

Phase 120 did not change public protobuf semantics, route auth, CSRF behavior,
credentials, token handling, public exposure, migrations, durable schema,
consumer statuses, protected evidence paths, release tags, release assets, or
published release state.

## Completed Checkpoints

- Phase 120 -- Checkpoint 000001: add gtfs rt feed usefulness and reliability
  v2 plan.
- Phase 120 -- Checkpoint 000002: implement or audit primary scoped work.
- Phase 120 -- Checkpoint 000003: run validation and patch required gaps.
- Phase 120 -- Checkpoint 000004: close gtfs rt feed usefulness and reliability
  v2 review.

## Product Result

Vehicle Positions review now gives operators and developers a compact,
redaction-safe reason breakdown for why trip descriptors were not published.
That makes local feed usefulness review more actionable while preserving the
core conservative matching behavior: prefer unknown or omitted trip context
over false certainty.

## Changed Files

- `internal/feed/vehicle_positions.go`
- `internal/feed/vehicle_positions_health.go`
- `internal/feed/vehicle_positions_test.go`
- `docs/phase-120-gtfs-rt-feed-usefulness-and-reliability-v2.md`
- `docs/handoffs/phase-120.md`
- `docs/handoffs/latest.md`
- `docs/current-status.md`
- `docs/roadmap-status.md`
- `docs/open-transit-rt-master-planner-remaining-work.md`

## Validation

Passed:

- `gofmt -w internal/feed/vehicle_positions.go internal/feed/vehicle_positions_health.go internal/feed/vehicle_positions_test.go`
- `go test ./internal/feed/...`
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

Blocked:

- None for Phase 120.

## Protected Path Status

No protected evidence path was edited, generated, reformatted, or touched.

## Consumer Tracker Status

`docs/evidence/consumer-submissions/status.json` was not edited. The exact
seven targets remain present in order and all remain `prepared`:

- Google Maps
- Apple Maps
- Transit App
- Bing Maps
- Moovit
- Mobility Database
- transit.land

## Claim Boundary Status

Phase 120 makes no stable release readiness, production readiness, compliance,
adoption, agency approval, consumer acceptance, consumer
ingestion/listing/display, final-root readiness, hosted service availability,
paid support, SLA/uptime, vendor compatibility, hardware certification,
production AVL reliability, production-grade ETA quality, or real-world ETA
accuracy claim.

## Security/Auth Status

No application security behavior changed. New aggregates omit raw telemetry
payloads, private paths, tokens, credentials, URLs, stdout/stderr, argv,
private debug blobs, and host-specific data.

## Data/Migration Status

No migration, schema, durable state, dependency, or Go module change was
added. Existing private `feed_health_snapshot.details_json` records remain
compatible; new Vehicle Positions snapshots may include additional bounded
aggregate keys.

## Release/Publication Status

The Phase 115 public `v0.1.0-rc.1` prerelease remains published. Phase 120 did
not publish, republish, retag, upload assets, or patch the public rc1 release.

## Install Confidence Status

Phase 117 public fresh-clone install confidence remains passed. Phase 120 is
current-source hardening after rc1 and is not part of the published rc1 tag.

## Web Design Skill Status

Phase 118 Web Design Skill artifact remains complete. Phase 120 did not change
visual UI templates.

## Commit List

- `dcc5307` -- Phase 120 -- Checkpoint 000001: add gtfs rt feed usefulness and
  reliability v2 plan
- `fc11708` -- Phase 120 -- Checkpoint 000002: implement or audit primary
  scoped work
- `8e0dddc` -- Phase 120 -- Checkpoint 000003: run validation and patch
  required gaps
- Phase 120 -- Checkpoint 000004: close gtfs rt feed usefulness and reliability
  v2 review

## Checkpoint Report

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
