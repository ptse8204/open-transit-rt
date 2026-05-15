# Phase 96 Handoff -- GTFS Versioning, Diff, And Rollback Workbench

## Status

Phase 96 is complete for the private GTFS versioning, diff, and rollback
Workbench scope. The private GTFS Workbench now shows active-vs-previous
schedule comparison, file-level row-count diffs, bounded route/stop/trip/
service/frequency sample summaries, and rollback-review guidance. It remains
read-only and does not execute rollback.

## Completed Checkpoints

- Phase 96 -- Checkpoint 000001: add GTFS versioning, diff, and rollback
  workbench plan.
- Phase 96 -- Checkpoint 000002: implement primary scoped work.
- Phase 96 -- Checkpoint 000003: run validation and patch required gaps.
- Phase 96 -- Checkpoint 000004: close GTFS versioning, diff, and rollback
  workbench review.

## Product Result

- `/admin/operations/gtfs-workbench` and
  `/admin/operations/gtfs-workbench.json` include a new
  `version_comparison` section.
- The comparison uses the active published schedule feed version and the first
  visible non-active previous feed version from bounded feed-version history.
- File-level row-count diffs cover `routes.txt`, `stops.txt`, `trips.txt`,
  `stop_times.txt`, `calendar.txt`, `calendar_dates.txt`, `shapes.txt`, and
  `frequencies.txt`.
- Bounded entity summaries cover routes, stops, trips, service calendars, and
  frequencies with capped added/removed/changed samples.
- Rollback guidance names candidate-review, realtime assignment review, and
  draft-only rollback command design limits without adding a browser rollback
  action.

## Changed Files

- `cmd/agency-config/operations_gtfs_workbench.go`
- `cmd/agency-config/operations.go`
- `cmd/agency-config/main_test.go`
- `docs/phase-96-gtfs-versioning-diff-and-rollback-workbench.md`
- `docs/handoffs/phase-96.md`
- `docs/handoffs/latest.md`
- `docs/current-status.md`
- `docs/roadmap-status.md`
- `docs/open-transit-rt-master-planner-remaining-work.md`

## Validation

Passed:

- `git status --short`
- `git diff --check`
- `go test ./cmd/agency-config -run 'GTFSWorkbench|GTFSImport|GTFSQuality|ValidationHealth|Setup'`
- `make check`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact prepared-only consumer tracker assertion
- protected-path status check

Blocked:

- `RUN_LOCAL_APP=true make release-candidate-check` was not run because Phase
  96 is not a release-candidate phase and does not require local package/app
  diagnostics.
- Release actions, retained evidence, consumer action, external contact, and
  rollback execution remain blocked by scope.

## Protected Path Status

No protected evidence path was edited, generated, reformatted, or touched. The
protected-path status check for `docs/evidence/consumer-submissions`,
`docs/evidence/captured`, `db/migrations`, `go.mod`, and `go.sum` returned no
output.

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

## Claim-Boundary Status

Phase 96 makes no compliance, release readiness, consumer ingestion, consumer
acceptance, agency approval, production readiness, final-root readiness,
hosted-service, vendor compatibility, hardware certification, SLA/uptime,
public launch, or ETA-quality claim. The new Workbench copy describes private
local review only.

## Security/Auth Status

No public admin route, auth behavior, CSRF behavior, credential path, token
handling, browser mutation, or command route changed. Existing tests continue
to verify the private Workbench route is no-store, read-only, agency-scoped,
POST-blocked, unauthenticated-blocked, and unavailable under `/public`.

## Data/Migration Status

No migration, durable diff artifact table, rollback command, feed activation
mutation, or schema change was added. The implementation reuses existing
published feed-version history and the bounded GTFS schedule preview reader.

## Commit List

- `00ccbc0` -- Phase 96 -- Checkpoint 000001: add gtfs versioning, diff, and rollback workbench plan
- `505bec5` -- Phase 96 -- Checkpoint 000002: implement primary scoped work
- `957f7fc` -- Phase 96 -- Checkpoint 000003: run validation and patch required gaps
- Phase 96 -- Checkpoint 000004: close GTFS versioning, diff, and rollback workbench review

## Checkpoint Report

Checkpoint:
Phase 96 -- Checkpoint 000004: close GTFS versioning, diff, and rollback
workbench review.

Sub-agents used or simulated, including intended model level:
Real Context / Repo Truth Sub-Agent -- GPT-5.5 x-high; real Planning
Sub-Agent -- GPT-5.5 x-high; real Implementation Review Sub-Agent -- GPT-5.5
high; real QA Sub-Agent -- GPT-5.5 high; real Claim-Boundary / Security/Auth
QA Sub-Agent -- GPT-5.5 high; real Data/Migration Sub-Agent -- GPT-5.5 high.
UI/UX and Documentation / IA closeout roles were simulated by the Master
Agent. Master Agent -- GPT-5.5 x-high, current thread.

Changed files:
`cmd/agency-config/operations_gtfs_workbench.go`;
`cmd/agency-config/operations.go`; `cmd/agency-config/main_test.go`;
`docs/phase-96-gtfs-versioning-diff-and-rollback-workbench.md`;
`docs/handoffs/phase-96.md`; `docs/handoffs/latest.md`;
`docs/current-status.md`; `docs/roadmap-status.md`;
`docs/open-transit-rt-master-planner-remaining-work.md`.

Validation run:
`git status --short`; `git diff --check`; `go test ./cmd/agency-config -run
'GTFSWorkbench|GTFSImport|GTFSQuality|ValidationHealth|Setup'`; `make check`;
`make validate`; `make test`; `docker compose -f deploy/docker-compose.yml
config`; `make audit-product-acceptance`; `make audit-final-claim-review`;
`python3 -m json.tool docs/evidence/consumer-submissions/status.json
>/dev/null`; exact prepared-only consumer tracker assertion; protected-path
status check.

Blocked checks:
`RUN_LOCAL_APP=true make release-candidate-check` was not run because Phase 96
is not a release-candidate phase. Release actions, retained evidence,
consumer action, external contact, and rollback execution remain blocked by
scope.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched. The
protected-path status check returned no output.

Consumer tracker status:
`docs/evidence/consumer-submissions/status.json` was not edited. The exact
seven consumer targets remain present in order and all remain `prepared`.

Claim-boundary status:
Phase 96 stays bounded to private local Workbench review and makes no stronger
public claim.

Security/auth status:
Existing private Workbench auth, role, agency-scope, no-store, no-POST, and
public-route-blocking behavior is preserved.

Data/migration status:
No migration or rollback mutation was added.

Master review:
Approved. Phase 96 met the authorized scope, added useful private schedule
version comparison, preserved protected paths and consumer statuses, and kept
rollback execution out of the browser.

Required edits:
None.

Decision:
Phase 96 is complete. Continue immediately to Phase 97.

Next checkpoint:
Phase 97 -- Checkpoint 000001: add GTFS quality fix planner and safe draft
suggestions plan.
