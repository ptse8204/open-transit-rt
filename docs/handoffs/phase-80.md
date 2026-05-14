# Phase 80 Handoff -- GTFS Workbench

## Phase

Phase 80 -- GTFS Workbench.

## Sub-Agents Used Or Simulated

Real sub-agents were used where available, with the intended model levels from
the Phase 75-90 track:

- Master Agent: GPT-5.5 x-high, main rollout.
- Context / Repo Truth Sub-Agent: GPT-5.5 x-high, simulated locally because
  the agent thread limit was reached for the context role.
- Planning Sub-Agent: GPT-5.5 x-high, Mencius.
- Security/Auth Sub-Agent: GPT-5.5 high, Mendel.
- UI/UX Sub-Agent: GPT-5.5 high, Nietzsche.
- QA Sub-Agent: GPT-5.5 high, Carver.
- Documentation / IA Sub-Agent: GPT-5.5 high, Newton.
- Claim-Boundary Sub-Agent: GPT-5.5 high, Sartre.
- Implementation Sub-Agent: GPT-5.5 high, simulated in the main rollout.
- Data/Migration Sub-Agent: not used because Phase 80 added no migration or
  destructive persisted model.

All reviews were resolved before implementation and closeout. No required edits
remain.

## Goal

Make static schedule work understandable and browser-first by adding a private
GTFS Workbench that coordinates active schedule state, import history, bounded
schedule previews, GTFS quality, validation, draft publish review, feed output,
and rollback guidance without adding public admin routes or new browser
mutation paths.

## Changed Files

- `cmd/agency-config/main.go`
- `cmd/agency-config/main_test.go`
- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_gtfs_workbench.go`
- `cmd/agency-config/operations_navigation.go`
- `cmd/gtfs-studio/main.go`
- `cmd/gtfs-studio/main_test.go`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-80.md`
- `docs/phase-80-gtfs-workbench.md`
- `docs/roadmap-status.md`
- `docs/tutorials/real-agency-gtfs-onboarding.md`
- `internal/compliance/model.go`
- `internal/compliance/postgres.go`

## Routes Added/Changed

Added private authenticated read-only Operations Console routes:

- `GET /admin/operations/gtfs-workbench`
- `GET /admin/operations/gtfs-workbench.json`

Changed private authenticated surfaces:

- `GET /admin/operations`
- `GET /admin/operations/gtfs-import`
- `GET /admin/gtfs-studio`
- `GET /admin/gtfs-studio/drafts/{draft_id}`
- GTFS Studio browser POST handlers now enforce CSRF for cookie-auth
  mutations when `CSRF_SECRET` is configured.

No public admin route was added. The Workbench has no POST route.

## Commands Added/Changed

No Makefile, Taskfile, CLI, release, evidence, consumer, or package commands
were added or changed.

## Migrations

None.

Phase 80 added read-only repository methods over existing `feed_version`,
`gtfs_import`, `gtfs_*`, `gtfs_draft`, and `gtfs_draft_publish` tables.

## Validation Run

- `git status --short`
- `git diff --check`
- `make check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`
- focused Workbench/GTFS/setup/validation tests:
  `go test ./cmd/agency-config -run 'GTFSWorkbench|GTFSImport|GTFSQuality|ValidationHealth|Setup'`
- `go test ./cmd/gtfs-studio`
- `go test ./internal/compliance ./internal/gtfs ./cmd/gtfs-import ./cmd/gtfs-studio`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`
- `RUN_LOCAL_APP=true make release-candidate-check`

## Blocked Checks

No Phase 80 closeout check is blocked.

`RUN_LOCAL_APP=true make release-candidate-check` wrote
`.cache/release-candidate-check/20260514T015120Z` with overall
`needs_review` because the helper detected the intentionally dirty closeout
worktree and kept its repository validation, Go unit test, HTTP smoke, and
release-package rows as bounded `not_checked` diagnostics. `make validate`,
`make test`, Docker Compose config, and the helper's local app startup/five
public feed diagnostics were run separately or passed within the helper. This
is not a release-ready pass and did not create a tag, package, published image,
or release artifact.

## Known Blockers

- Phase 72 remains `needs_review`, not release-ready.
- No release tag, package, published image, final-root proof, consumer
  submission, compliance packet, real agency pilot, real vendor/device proof,
  SLA proof, or production ETA-quality proof exists.
- Optional evidence tracks remain authorization-gated.

## Protected Path Status

Protected evidence paths were not modified:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/**`
- `docs/evidence/consumer-submissions/artifacts/**`
- `docs/evidence/consumer-submissions/packets/**`

`db/migrations`, `go.mod`, and `go.sum` were unchanged.

## Consumer Tracker Status

All seven consumer targets remain exactly `prepared`:

- Google Maps
- Apple Maps
- Transit App
- Bing Maps
- Moovit
- Mobility Database
- transit.land

## Claim-Boundary Status

Phase 80 added no CAL-ITP/Caltrans compliance, agency adoption/approval,
consumer submission/review/acceptance/ingestion/listing/display,
final-root readiness, hosted SaaS, paid support, SLA/uptime, production
readiness, vendor compatibility, hardware certification, public launch, or
production-grade ETA claim.

The Workbench labels local operator review, importer outcomes, validator
health, preview rows, draft publish history, and schedule history as private
diagnostics/supporting signals only.

## Security/Auth Status

- Workbench routes are authenticated, private, GET-only, no-store, and
  agency-scoped through the existing principal checks.
- Workbench JSON is read-only and bounded.
- Workbench has no POST route, no browser-supplied backend command model, no
  raw validator report exposure, no raw ZIP bytes, no arbitrary path, and no
  external submission target.
- GTFS Studio publish/discard/create/update/remove mutations now require CSRF
  for cookie-auth requests when the CSRF secret is configured.
- Existing GTFS import and validator run actions remain admin-only and
  server-owned.

## Accessibility Status

The Workbench uses the existing Operations Console shell, headings, status
chips, tables, landmarks, skip link, focus styling, mobile table overflow, and
buildless review filters. No SPA or frontend dependency was added.

## Docs/Site/Wiki Alignment

Status docs now describe Phase 80 as complete and point to Phase 81 as the
next authorized product-track phase. Public docs continue to describe the
Consumer-Grade Control Plane roadmap as proposed/authorized product work, not
implemented public release or evidence proof.

## Commit List

- `d5222f5` Phase 80 -- Checkpoint 000001: add GTFS Workbench plan
- `47f7c85` Phase 80 -- Checkpoint 000002: add import diff and schedule summaries
- `e35a9c1` Phase 80 -- Checkpoint 000003: add GTFS preview tables and filters
- `7bc4d72` Phase 80 -- Checkpoint 000004: improve safe draft publish review
- `e06fc30` Phase 80 -- Checkpoint 000005: add rollback and schedule history UX
- Phase 80 -- Checkpoint 000006: close GTFS Workbench review

## Master Review

The Master Agent reviewed the real/simulated sub-agent reports, code changes,
tests, docs, protected paths, consumer tracker state, and claim boundaries.
Phase 80 is bounded to private product work and is safe to close after
validation passes.

## Required Edits

None.

## Decision

Close Phase 80 after the final closeout commit.

## Next Phase

Phase 81 -- Realtime Operations Center.
