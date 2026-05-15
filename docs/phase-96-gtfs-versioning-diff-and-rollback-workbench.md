# Phase 96 -- GTFS Versioning, Diff, And Rollback Workbench

## Scope

Phase 96 improves the private GTFS Workbench so operators and technical
helpers can compare schedule versions, understand import-to-active change
signals, and review safe rollback steps without silently editing production
GTFS.

The phase stays inside private Operations Console product work. It does not
add public admin routes, execute rollback from the browser, collect retained
evidence, contact outside systems, move consumer statuses, publish release
artifacts, or claim compliance, release readiness, consumer acceptance,
production readiness, final-root readiness, hosted-service availability,
vendor compatibility, hardware certification, SLA/uptime, or ETA quality.

## Current Workbench Truth

- `/admin/operations/gtfs-workbench` and
  `/admin/operations/gtfs-workbench.json` are private, authenticated,
  agency-scoped, no-store, GET-only review routes.
- The Workbench already summarizes active schedule state, recent imports,
  source checksum and byte-size change signals, GTFS quality, validation
  health, bounded preview tables, GTFS Studio draft publish state, recent
  feed versions, rollback guidance, and all-false claim flags.
- Existing tests assert bounded JSON, no private path leaks, no mutation form,
  no POST action, no public route, and no forbidden claim wording.
- Existing persistence already stores immutable published GTFS rows by
  `feed_version_id`, feed-version lifecycle rows, active published feed
  pointers, import history, draft/publish history, and audit logs. Phase 96 can
  compute read-only diffs on demand without a migration.

## Deliverables

- Add or improve private Workbench active-vs-previous schedule comparison.
- Add import-to-active diff summary rows that show feed version identity,
  source checksum, byte count, validator/import status, and lifecycle mismatch
  signals.
- Add file-level and row-count diff coverage for core GTFS tables:
  `routes.txt`, `stops.txt`, `trips.txt`, `stop_times.txt`, `calendar.txt`,
  `calendar_dates.txt`, `shapes.txt`, and `frequencies.txt`.
- Add route, trip, stop, and service-calendar change summaries with bounded
  samples and operator next actions.
- Add rollback readiness guidance with explicit blockers, warnings, candidate
  version, required validation/feed-health/realtime review steps, and a
  draft-only command design where executable rollback is not safe to add.
- Add or update focused tests for private route behavior, bounded JSON,
  forbidden claims, and no rollback mutation path.
- Add `docs/handoffs/phase-96.md` at closeout and update source-of-truth docs.

## Non-Goals

- No automatic GTFS edits.
- No browser rollback POST route.
- No rollback mutation command unless a later phase explicitly authorizes and
  reviews it.
- No durable diff artifact tables or migration in this phase.
- No retained evidence, protected-path writes, consumer tracker changes, public
  publication, release action, or external contact.
- No claim that a diff is complete proof of source correctness, validator-clean
  status, compliance, consumer ingestion, or rollback safety.

## Data And Diff Design

Use existing `feed_version_id` scoping and compare versions only within the
same `agency_id`.

Primary version selection:

- active schedule feed version from publication discovery;
- previous visible retired/staged feed version from recent feed-version
  history;
- latest import feed version from recent import history for import-to-active
  review.

Safe first-pass diff behavior:

- read only existing GTFS preview/feed-version/import records;
- compare row counts for core GTFS files;
- surface bounded representative added, removed, or changed IDs only;
- preserve raw GTFS time text and do not coerce after-midnight times into
  clock-only values;
- treat frequency rows as service-affecting changes;
- keep all private paths, raw ZIP bytes, raw validator reports, credentials,
  and unbounded row lists out of HTML and JSON.

If repository-level table diff support is needed, add it as an explicit
read-only interface and Postgres repository method without changing the schema.

## Checkpoints

### Checkpoint 000001 -- Plan

Deliverables:

- Add this phase plan.
- Record sub-agent workflow, no-migration decision, validation plan, and Master
  approval before implementation.

Validation:

- `git status --short`
- `git diff --check`
- `make check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- consumer tracker JSON parse and exact prepared-only assertion
- protected-path status check

### Checkpoint 000002 -- Implement Primary Scoped Work

Deliverables:

- Add read-only active-vs-previous and import-to-active comparison model fields
  to the private GTFS Workbench JSON.
- Add HTML sections for file-level row-count diffs and route/trip/stop/service
  summaries.
- Keep output bounded, private, no-store, GET-only, agency-scoped, and
  claim-bounded.
- Update focused tests for JSON keys, HTML labels, bounded samples, no forms,
  no POST, no path leaks, and no forbidden claims.

Validation:

- `gofmt` on changed Go files.
- `go test ./cmd/agency-config -run 'GTFSWorkbench|GTFSImport|GTFSQuality|ValidationHealth|Setup'`
- `git diff --check`
- protected-path and consumer tracker checks.

### Checkpoint 000003 -- Rollback Guidance And Validation Patch

Deliverables:

- Add rollback readiness guidance with candidate version, blockers, warnings,
  operator approval step, validator/feed-health review, realtime assignment
  implications, and audit expectation.
- Add draft-only rollback command design text where executable rollback is not
  safe to add.
- Patch any validation, copy, claim-boundary, or test gaps found in CP000002.

Validation:

- focused Workbench tests;
- `make validate`;
- `make test`;
- `docker compose -f deploy/docker-compose.yml config`;
- baseline claim/protected-path checks.

### Checkpoint 000004 -- Closeout

Deliverables:

- Add `docs/handoffs/phase-96.md`.
- Update `docs/current-status.md`, `docs/handoffs/latest.md`,
  `docs/roadmap-status.md`, and
  `docs/open-transit-rt-master-planner-remaining-work.md`.
- Record final validation, blockers, protected-path status, consumer tracker
  status, claim-boundary status, security/auth status, and data/migration
  status.

Validation:

- `git status --short`
- `git diff --check`
- `make check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- consumer tracker JSON parse and exact prepared-only assertion
- protected-path status check
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`

## Hard Boundaries

- Do not modify, generate, rewrite, reformat, or touch protected evidence paths.
- Do not edit `docs/evidence/consumer-submissions/status.json`.
- Do not move any consumer or aggregator status beyond `prepared`.
- Do not contact agencies, vendors, consumers, portals, map providers,
  aggregators, or external services for evidence or proof.
- Do not use real credentials, private payloads, private AVL/vendor data, or
  real agency data.
- Do not add a public admin route.
- Do not execute rollback from the Workbench.
- Do not collapse draft GTFS and published feed versions into one model.
- Do not claim compliance, release readiness, public launch, consumer
  acceptance, production readiness, final-root readiness, hosted SaaS, vendor
  compatibility, hardware certification, SLA/uptime, or ETA quality.

## Master Approval

The Master Agent approves this Phase 96 plan. Implementation may proceed after
Checkpoint 000001 is committed. The approved implementation posture is
read-only diff/review with no migration, no browser rollback execution, bounded
operator guidance, and explicit claim boundaries.

## Checkpoint 000001 Report

Checkpoint:
Phase 96 -- Checkpoint 000001: add GTFS versioning, diff, and rollback
workbench plan.

Sub-agents used or simulated, including intended model level:
Real Context / Repo Truth Sub-Agent -- GPT-5.5 x-high returned published
feed-version, Workbench, and rollback seam findings. Real Planning Sub-Agent
-- GPT-5.5 x-high returned the four-checkpoint plan. Real Data/Migration
Sub-Agent -- GPT-5.5 high returned a no-migration recommendation for read-only
on-demand diffs. Real Claim-Boundary / Security/Auth QA Sub-Agent -- GPT-5.5
high returned protected-path, auth, no-silent-edit, and wording constraints.
UI/UX, Documentation, Implementation, and closeout QA roles are simulated by
the Master Agent for this checkpoint. Master Agent -- GPT-5.5 x-high, current
thread.

Changed files:
`docs/phase-96-gtfs-versioning-diff-and-rollback-workbench.md`.

Validation run:
`git status --short`; `git diff --check`; `make check`; `make
audit-product-acceptance`; `make audit-final-claim-review`; `python3 -m
json.tool docs/evidence/consumer-submissions/status.json >/dev/null`; exact
prepared-only consumer tracker assertion; `git status --short --
docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod
go.sum`.

Blocked checks:
Code tests, `make validate`, `make test`, and compose config are deferred
until implementation checkpoints because this checkpoint is planning-only.

Protected path status:
No protected evidence path is edited or generated by this plan.

Consumer tracker status:
`docs/evidence/consumer-submissions/status.json` is not edited. All seven
consumer targets must remain exactly `prepared`.

Claim-boundary status:
The plan authorizes private read-only GTFS review improvements only. It makes
no compliance, release-ready, consumer, final-root, hosted-service,
production-readiness, vendor, hardware, SLA/uptime, or ETA-quality claim.

Security/auth status:
The plan preserves existing private authenticated Workbench routes, no-store
behavior, agency scoping, GET-only behavior, and no browser rollback
execution.

Data/migration status:
No migration is planned. Existing versioned GTFS, feed-version, import,
draft/publish, and audit tables are sufficient for read-only on-demand diff
and rollback guidance.

Master review:
Approved. The plan is inside Phase 96 scope, keeps draft and published GTFS
separate, preserves protected paths and prepared-only consumer statuses, and
uses a no-migration read-only implementation posture.

Required edits:
None before committing this checkpoint.

Decision:
Commit Checkpoint 000001, then proceed to CP000002 implementation.

Next checkpoint:
Phase 96 -- Checkpoint 000002: implement primary scoped work.
