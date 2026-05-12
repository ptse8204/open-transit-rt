# Phase 72 Handoff -- v0.1.0-rc.1 Release Candidate Hardening

## Status

Complete for the bounded `v0.1.0-rc.1` hardening review. This is not a release
tag, release-ready pass, validator-clean claim, production-readiness claim, or
public launch.

Checkpoint 000001 is complete for docs-only planning/status alignment.
Checkpoint 000002 is complete for the approved temporary evaluator gate. It
recorded the primary checkout as dirty, created a clean local-only evaluator
snapshot from that state, ran the approved CP000002 commands there, and found
one local pinned-validator tooling blocker in `make release-candidate-check`.
Checkpoint 000003 is complete for diagnostic hardening only.

CP000002 did not fix code, did not run out-of-scope release, app, browser,
feed-fetch, validator-install, evidence, connector, publishing, or external
contact commands, did not draft release notes, did not create retained
evidence, and did not change consumer tracker statuses.

CP000003 did not install validators, did not run `make validators-install`,
did not edit protected evidence paths or consumer statuses, did not change
migrations, module files, runtime routes, schemas, public feed URLs, tags,
packages, retained evidence, or release notes, and did not make release-ready,
validator-clean, compliance, production-readiness, consumer, or agency claims.

Checkpoint 000004 is complete for local app startup, required private
Operations Console route checks, and five anonymous public feed fetches.
`make agency-app-up` passed, all nine required private routes returned
authenticated local `200` responses, all five public feed paths returned
anonymous local `200` responses, and `feeds.json` parsed with
`python3 -m json.tool`. Browser automation with safe bearer-header support was
unavailable, so terminal authenticated GET checks using an `Authorization`
header were used as the safe substitute without exposing the admin token. The
CP000004 `RUN_LOCAL_APP=true` release-candidate diagnostic passed its
`local_app_five_feeds` and `validators_check` rows and recorded
`overall_status=needs_review` because release-package auditing was
intentionally not enabled. `make agency-app-down` stopped the app.

Checkpoint 000005 is complete for local synthetic connector and adapter
conformance gates. `make external-connection-check`,
`make adapter-conformance`, and `make test-connector-examples` all exited `0`.
No connector/adaptor code, schemas, fixtures, examples, protected evidence
paths, or consumer tracker statuses were changed.

Checkpoint 000006 is complete for a local pre-tag rc1 release notes and known
blockers draft at `docs/release-notes-v0.1.0-rc.1-draft.md`. It did not tag,
publish, package, create retained evidence, change consumer statuses, install
validators, or run release packaging.

Checkpoint 000007 is complete for closeout. Final focused checks passed where
run, the local-app release-candidate diagnostic exited `0` with
`overall_status=needs_review`. Phase 72 is not release-ready. Remaining
`needs_review` items are the dirty primary checkout, release package audit
`not_checked`, no tag/package/published image, no retained evidence, and the
CP000004 browser automation limitation; terminal authenticated checks were the
recorded substitute. Consumer tracker statuses remain unchanged as a protected
boundary.

Current default next work is Phase 73 CP000006: close agency UI acceptance
review. Phase 73 CP000001 is complete for documentation-only agency UI
acceptance planning, CP000002 is complete for local no-developer browser
walkthrough review, CP000003 is complete for local technical-helper walkthrough
review, CP000004 is complete for narrow UI copy, route-label,
Devices/Telemetry boundary-copy, and browser-first tutorial patching, and
CP000005 is complete for small-agency docs/wiki navigation freeze.

## CP000001 Scope

- Add `docs/phase-72-v0.1.0-rc.1-release-candidate-hardening.md`.
- Add `docs/handoffs/phase-72.md`.
- Update `docs/release-candidate-readiness.md`,
  `docs/roadmap-status.md`, and `docs/handoffs/latest.md` narrowly to record
  the initial Phase 72 launch state and boundary.
- Do not edit code, scripts, migrations, module files, Makefile, README, wiki,
  site files, protected evidence paths, consumer status JSON, or current/
  packet/artifact/captured evidence.
- Do not run the full release-candidate gate.
- Do not create final release notes, tags, packages, or public release
  artifacts; CP000006 release-note output is a local pre-tag draft only.

## CP000002 Scope

- Run primary checkout source review only in the dirty primary checkout.
- Create a temporary local-only evaluator clone/snapshot commit from the
  current approved primary state.
- Run the approved evaluator commands in the temporary clean snapshot.
- Record both primary dirty state and evaluator clean state.
- Edit only this handoff and the Phase 72 plan/results matrix.
- Do not edit code, scripts, `Makefile`, protected evidence paths, consumer
  status JSON, migrations, `go.mod`, `go.sum`, release notes, or other files.

## CP000003 Scope

- Harden release-candidate diagnostic handling for validator blockers.
- Force the release-candidate `validators_check` row to invoke
  `scripts/check-validators.sh` in pinned mode so ambient
  `VALIDATOR_TOOLING_MODE=stub` cannot make the row pass.
- Preserve dry-run behavior where `validators_check` is `not_checked`.
- Keep validator-tooling failures as `blocker` rows and make the full
  diagnostic exit nonzero when any blocker remains.
- Include the exact final `scripts/check-validators.sh` blocker line in
  `summary.json` and `summary.md` while keeping full command output in
  `check-log.txt`.
- Keep output limited to exactly `summary.json`, `summary.md`,
  `manifest.json`, `manifest.md`, and `check-log.txt`.
- Update only approved Phase 72/readiness/handoff docs for this behavior and
  result boundary.

## Current Repo Truth

- Phase 71 is complete for adoption-first productization and no-CLI agency
  operations.
- Phase 72 is complete for bounded `v0.1.0-rc.1` release-candidate hardening
  review with `needs_review` diagnostics.
- CP000002 ran the bounded evaluator gate in a temporary clean snapshot, not in
  the primary dirty checkout, except for primary source/protected-path review
  commands.
- CP000003 hardened the primary release-candidate diagnostic script and test
  coverage in the primary dirty checkout.
- CP000004 verified the local app, required private Operations Console routes,
  and five anonymous public feed paths in the primary dirty checkout.
- The primary checkout remains dirty with approved docs changes; do not call it
  clean.
- The evaluator snapshot was clean at local-only commit
  `83e01831c035404e53e921b19c6b7e89e9541816`.
- Release-candidate work remains local/private until a maintainer separately
  authorizes tagging, publishing, retained evidence, or external submission.
- No current Phase 72 result proves production readiness, CAL-ITP/Caltrans
  compliance, agency approval, consumer acceptance, final-root readiness,
  hosted SaaS, SLA/uptime, vendor compatibility, production AVL reliability,
  or production-grade ETA quality.

## Protected Paths

Do not edit or generate files under:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/**`
- `docs/evidence/consumer-submissions/artifacts/**`
- `docs/evidence/consumer-submissions/packets/**`
- `db/migrations/**`
- `go.mod`
- `go.sum`

For CP000001 specifically, also avoid code, scripts, `Makefile`, README, wiki,
and site files.

## Consumer Tracker State

Prepared-only. All seven consumer and aggregator targets remain `prepared`.
CP000002 parsed `docs/evidence/consumer-submissions/status.json` and ran the
exact seven-target prepared-only tracker check successfully. Phase 72 did not
submit, contact, automate portals, refresh target-originated records, or move
any target beyond `prepared`. Future consumer or aggregator work remains
separately authorization-gated.

## Checkpoint Sequence

1. `CP000001 -- add release-candidate hardening plan` -- complete for
   docs-only setup.
2. `CP000002 -- run clean-checkout local evaluator gate` -- complete with one
   local pinned-validator tooling blocker.
3. `CP000003 -- harden release-candidate diagnostics and blockers` --
   complete for diagnostic hardening only.
4. `CP000004 -- verify browser-first agency operations walkthrough` --
   complete with terminal authenticated checks substituting for unavailable
   in-app Browser review.
5. `CP000005 -- verify connector and adapter conformance gates` -- complete for
   local synthetic connector/adaptor checks.
6. `CP000006 -- prepare rc1 release notes and known blockers` -- complete for
   local pre-tag draft notes and blocker matrix.
7. `CP000007 -- close rc1 hardening review` -- complete for bounded closeout
   with `needs_review` blockers/deferrals recorded.

## CP000001 Review

Status: complete.

Review placeholder:

- Confirm only approved docs changed.
- Confirm CP000001 recorded the initial Phase 72 launch state without claiming
  checks passed.
- Confirm all initial blocker matrix rows are `not_run_in_cp000001`.
- Confirm release notes were deferred at CP000001.
- Confirm protected evidence paths and consumer tracker statuses were not
  edited.

## CP000002 Review

Status: complete for the approved temporary evaluator gate with one local
pinned-validator tooling blocker.

Temporary evaluator path:

- `/var/folders/_g/bvzl9cms7cx1d0wdpc981n9w0000gn/T/tmp.BA6ULYuxdo/open-transit-rt-evaluator`

Primary checkout source review:

| Command | Exit | Result |
| --- | ---: | --- |
| `git status --short` | 0 | Dirty with approved docs changes: modified `docs/handoffs/latest.md`, `docs/handoffs/phase-69.md`, `docs/handoffs/phase-71.md`, `docs/open-transit-rt-master-planner-remaining-work.md`, `docs/release-candidate-readiness.md`, `docs/roadmap-status.md`, `docs/tutorials/small-agency-acceptance-script.md`, `wiki/browser-first-setup.md`; untracked `docs/handoffs/phase-72.md`, `docs/phase-72-v0.1.0-rc.1-release-candidate-hardening.md`. |
| `git describe --tags --always --dirty` | 0 | `b0e354d-dirty` |
| `git rev-parse HEAD` | 0 | `b0e354d29b14ba5c5ab0afab3269ec5eb3405fa7` |
| `git status --short -- docs/evidence/captured docs/evidence/consumer-submissions db/migrations go.mod go.sum` | 0 | Empty output; no protected source-path changes reported. |

Evaluator source review:

| Command | Exit | Result |
| --- | ---: | --- |
| `git status --short` | 0 | Empty output; evaluator snapshot clean. |
| `git describe --tags --always --dirty` | 0 | `83e0183` |
| `git rev-parse HEAD` | 0 | `83e01831c035404e53e921b19c6b7e89e9541816` |

Evaluator command results:

| Command | Exit | Status | Result / blocker |
| --- | ---: | --- | --- |
| `make check` | 0 | passed | Lightweight no-network/no-Docker/no-validator-install checks passed. |
| `scripts/release-candidate-check.sh --dry-run` | 0 | passed | Private dry-run diagnostics written to `.cache/release-candidate-check/20260512T062130Z`. |
| `make release-candidate-check` | 2 | blocker | Private diagnostics written to `.cache/release-candidate-check/20260512T062134Z`; summary overall status `blocker`. Exact blocker from `check-log.txt`: `missing pinned tooling: static GTFS validator not installed at /private/var/folders/_g/bvzl9cms7cx1d0wdpc981n9w0000gn/T/tmp.BA6ULYuxdo/open-transit-rt-evaluator/.cache/validators/gtfs-validator-7.1.0-cli.jar; run make validators-install`. Classified as local environment/tooling blocker. |
| `scripts/bootstrap-dev.sh --check` | 0 | passed | Go, Docker, Compose render, and local bootstrap `DATABASE_URL` preflight passed. |
| `make test` | 0 | passed | `go test ./...` passed. |
| `docker compose -f deploy/docker-compose.yml config` | 0 | passed | Compose configuration rendered successfully. |

Post-check guards in primary:

| Command | Exit | Result |
| --- | ---: | --- |
| `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` | 0 | Passed. |
| Exact seven-target prepared-only tracker check | 0 | Passed for Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land. |
| `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured` | 0 | Empty output; no protected evidence path changes reported. |

Commands intentionally not run in CP000002:

- `RUN_LOCAL_APP=true make release-candidate-check`
- `make agency-app-up` / `make agency-app-down`
- browser walkthroughs
- five feed fetches
- `make validate`
- `make validators-install`
- release-package / audit-release-package
- `make external-connection-check`, `make adapter-conformance`,
  `make test-connector-examples`
- evidence helpers, tags, publishing, image pushes, external contacts

## CP000003 Review

Status: complete for diagnostic hardening after required-edit validation.

Implemented behavior:

- `scripts/release-candidate-check.sh` now runs the validator tooling row as
  `VALIDATOR_TOOLING_MODE=pinned scripts/check-validators.sh`.
- Dry-run still records `validators_check` as `not_checked`.
- A failed validator tooling row remains `blocker`; the overall diagnostic
  exits nonzero when the summary overall status is `blocker`.
- `validators_check` failure detail in `summary.json` and `summary.md` is the
  exact final non-empty blocker line emitted by `scripts/check-validators.sh`.
- `check-log.txt` keeps the full command output.
- The output file set remains exactly `summary.json`, `summary.md`,
  `manifest.json`, `manifest.md`, and `check-log.txt`.

Validation and guard results are recorded in the CP000003 validation matrix
section after command execution. Missing pinned validator tooling remains a
local environment/tooling blocker. CP000003 did not install validators and did
not create validator-clean or release-ready evidence.

Master-side review initially found a `make test` blocker in an existing
validation-health concurrency test: `fakeRealtimeArtifacts` and
`fakePublicationStore` in `cmd/agency-config/main_test.go` used shared mutable
test-fake state during concurrent requests. CP000003 fixed that test-only
state with mutex guards and copied returned slices. The fix did not change
production validation semantics, runtime routes, migrations, module files,
schemas, protected evidence paths, or consumer statuses.

CP000003 validation results:

| Command | Exit | Status | Result / blocker |
| --- | ---: | --- | --- |
| `git status --short` | 0 | passed | Primary checkout remains dirty with pre-existing docs changes plus CP000003 script/test/doc edits. |
| `sh -n scripts/release-candidate-check.sh scripts/check-validators.sh` | 0 | passed | Shell syntax passed. |
| `go test ./cmd/agency-config -run ReleaseCandidate` | 0 | passed | Focused release-candidate tests passed. |
| `make test-release-candidate-check` | 0 | passed | Make wrapper passed. |
| `OUTPUT_DIR=.cache/release-candidate-check-cp000003/dry-run FORCE=true scripts/release-candidate-check.sh --dry-run` | 0 | passed | Dry-run diagnostics wrote exactly five files; `validators_check` stayed `not_checked`. |
| Forced negative diagnostic with ambient `VALIDATOR_TOOLING_MODE=stub` and invalid `GTFS_VALIDATOR_PATH` | 1 | passed_expected_failure | Diagnostics wrote exactly five files; `validators_check` was `blocker`; `summary.json`, `summary.md`, and `check-log.txt` included the actionable `misconfigured pinned tooling` blocker detail. |
| `make check` | 0 | passed | Lightweight no-network/no-Docker/no-validator-install checks passed. |
| Initial master-side `make test` | 2 | blocker_then_fixed | Exposed a test-only concurrent map write in `cmd/agency-config/main_test.go`; CP000003 patched the shared validation-health test fakes. |
| `go test ./cmd/agency-config -run TestValidationHealthConcurrentRunAllDoesNotPanicOrLeak -count=100` | 0 | passed | Required-edit concurrency regression check passed. |
| `go test ./cmd/agency-config -race -run TestValidationHealthConcurrentRunAllDoesNotPanicOrLeak -count=5` | 0 | passed | Race check for the required-edit test fake patch passed. |
| `go test ./cmd/agency-config` | 0 | passed | Full agency-config package tests passed after the required edit. |
| `make test` | 0 | passed | Full repository Go test suite passed after the required edit. |
| `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` | 0 | passed | Consumer tracker JSON parsed. |
| Exact seven-target prepared-only tracker check | 0 | passed | Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land remained `prepared`. |
| `git status --short -- docs/evidence/captured docs/evidence/consumer-submissions db/migrations go.mod go.sum` | 0 | passed | Empty output; no protected evidence, migration, or module-path changes reported. |

Commands intentionally not run in CP000003:

- `make validators-install`
- `make validate`
- `RUN_LOCAL_APP=true make release-candidate-check`
- app startup, browser walkthroughs, five feed fetches
- release-package / audit-release-package
- connector/adapter gates scheduled for CP000005
- tags, publishing, image pushes, retained evidence, external contacts

## CP000004 Review

Status: complete for local app, required route, and public feed verification.

Scope and source guards:

| Command | Exit | Status | Result |
| --- | ---: | --- | --- |
| `git status --short` | 0 | passed | Primary checkout was dirty before CP000004 with pre-existing modified code/docs/scripts/wiki files and untracked Phase 72 docs; not treated as clean. |
| `git status --short -- docs/evidence/captured docs/evidence/consumer-submissions/status.json docs/evidence/consumer-submissions/current docs/evidence/consumer-submissions/artifacts docs/evidence/consumer-submissions/packets db/migrations go.mod go.sum` | 0 | passed | Empty output; no protected evidence, migration, or module-path changes reported at checkpoint start. |

Local app lifecycle:

| Command | Exit | Status | Result |
| --- | ---: | --- | --- |
| `make agency-app-up` | 0 | passed | Started local Postgres/PostGIS, built the local app image, applied migrations, seeded demo records, started local services and proxy, imported sample GTFS, bootstrapped publication metadata, and reported `http://localhost:8080` running. |
| `make agency-app-down` | 0 | passed | Stopped and removed local app containers and the local Compose network. |

Auth and browser boundary:

- Admin token was generated only into a local shell variable, parsed from the
  `token=` line, used only in terminal `Authorization` headers, and unset
  after the shell block.
- The token was not printed, stored, documented, placed in a URL, or pasted
  into browser-visible fields.
- Public feed fetches used no auth.
- Browser automation with safe bearer-header support was unavailable in this
  session. Terminal authenticated GET checks were used as the safe substitute.

Unauthenticated admin boundary check:

| Route | HTTP status | Result |
| --- | ---: | --- |
| `/admin/operations` without an admin token | 401 | Admin route stayed non-public. |

Required private route checks:

| Route | HTTP status | Content type | Bytes | Page-specific text |
| --- | ---: | --- | ---: | --- |
| `/admin/operations` | 200 | `text/html; charset=utf-8` | 40873 | `Agency Operations Cockpit / Start Here` present |
| `/admin/operations/gtfs-import` | 200 | `text/html; charset=utf-8` | 12045 | `GTFS Import` present |
| `/admin/operations/feed-health` | 200 | `text/html; charset=utf-8` | 23661 | `Feed Health` present |
| `/admin/operations/gtfs-quality` | 200 | `text/html; charset=utf-8` | 16216 | `GTFS Quality` present |
| `/admin/operations/validation-health` | 200 | `text/html; charset=utf-8` | 12310 | `Validator Health` present |
| `/admin/operations/devices` | 200 | `text/html; charset=utf-8` | 11323 | `Devices` present |
| `/admin/operations/telemetry-simulator` | 200 | `text/html; charset=utf-8` | 13698 | `Telemetry Simulator` present |
| `/admin/operations/connectors` | 200 | `text/html; charset=utf-8` | 16611 | `Connector` present |
| `/admin/operations/maintenance` | 200 | `text/html; charset=utf-8` | 15078 | `Maintenance` present |

An initial route-check script attempt used zsh-incompatible Bash array
expansion and exited before route results were produced. The corrected Bash
run above is the recorded route result and was not a product blocker.

Required anonymous public feed checks:

| Path | HTTP status | Content type | Bytes | Additional check |
| --- | ---: | --- | ---: | --- |
| `/public/feeds.json` | 200 | `application/json` | 2477 | `python3 -m json.tool` parse passed |
| `/public/gtfs/schedule.zip` | 200 | `application/zip` | 1913 | fetched without auth |
| `/public/gtfsrt/vehicle_positions.pb` | 200 | `application/x-protobuf` | 15 | fetched without auth |
| `/public/gtfsrt/trip_updates.pb` | 200 | `application/x-protobuf` | 15 | fetched without auth |
| `/public/gtfsrt/alerts.pb` | 200 | `application/x-protobuf` | 135 | fetched without auth |

Optional diagnostic:

| Command | Exit | Status | Result |
| --- | ---: | --- | --- |
| `OUTPUT_DIR=.cache/phase-72-cp000004/release-candidate-check FORCE=true RUN_LOCAL_APP=true scripts/release-candidate-check.sh` | 0 | passed_needs_review | Private diagnostics written to `.cache/phase-72-cp000004/release-candidate-check`; `overall_status=needs_review`; `validators_check=passed`; `local_app_five_feeds=passed`; `release_package_audit=not_checked`. |

Commands intentionally not run in CP000004:

- `make telemetry-simulator`
- `make external-connection-check`, `make adapter-conformance`, and
  `make test-connector-examples`
- release-package generation or audit
- release notes
- tags, publishing, image pushes, retained evidence, external contact, or
  consumer submissions

CP000004 created no release-readiness, validator-clean, production-readiness,
CAL-ITP/Caltrans compliance, consumer submission/review/acceptance/ingestion/
listing/display, agency adoption/approval, final-root readiness, hosted SaaS,
SLA/uptime, vendor compatibility/hardware certification, production AVL
reliability, or production-grade ETA-quality claim.

## CP000005 Review

Status: complete for local synthetic connector and adapter conformance gates.

Scope and source guards:

| Command | Exit | Status | Result |
| --- | ---: | --- | --- |
| `git status --short` | 0 | passed | Primary checkout was dirty with existing modified code/docs/scripts/wiki files and untracked Phase 72 docs; not treated as clean. |
| `git status --short -- docs/evidence/captured docs/evidence/consumer-submissions/status.json docs/evidence/consumer-submissions/current docs/evidence/consumer-submissions/artifacts docs/evidence/consumer-submissions/packets db/migrations go.mod go.sum` | 0 | passed | Empty output; no protected evidence, migration, or module-path changes reported. |
| `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` plus exact seven-target prepared-only check | 0 | passed | Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land remained exactly `prepared`. |

Connector and adapter gates:

| Command | Exit | Status | Result |
| --- | ---: | --- | --- |
| `make external-connection-check` | 0 | passed | Validated connector manifests/fixtures, ran `go test ./internal/connectors`, ran `go test ./examples/connectors/...`, and reported connector manifests remain sidecar/manifest/conformance bounded. |
| `make adapter-conformance` | 0 | passed | `testdata/adapter-conformance` synthetic suite passed. |
| `make test-connector-examples` | 0 | passed | `go test ./examples/connectors/...` passed for committed synthetic connector examples. |

Boundary:

CP000005 is a local synthetic connector/adaptor gate only. It did not contact
external systems, automate portals, submit to consumers, publish packages or
images, create retained evidence, change public feed URLs, change connector
schemas, or change consumer statuses. It does not prove vendor compatibility,
hardware certification, production AVL reliability, production readiness,
release readiness, CAL-ITP/Caltrans compliance, consumer acceptance, or
production-grade ETA quality.

## CP000006 Review

Status: complete for local pre-tag release notes and known blockers draft.

Deliverable:

- `docs/release-notes-v0.1.0-rc.1-draft.md`

Source metadata recorded in the draft:

| Check | Result |
| --- | --- |
| `git rev-parse HEAD` | `b0e354d29b14ba5c5ab0afab3269ec5eb3405fa7` |
| `git describe --tags --always --dirty` | `b0e354d-dirty` |
| Source state | Dirty primary checkout with approved Phase 71/72 work; not a clean tagged source state. |

CP000006 recorded the known blockers from CP000002 through CP000005:

- unresolved pinned static GTFS validator tooling blocker in the CP000002
  evaluator full `make release-candidate-check`;
- dirty primary checkout;
- no release package, checksum, SBOM, provenance, tag, GitHub release, or
  published image;
- CP000007 final validation and claim-boundary closeout completed with
  remaining `needs_review` blockers/deferrals recorded;
- CP000004 browser automation with safe bearer-header support unavailable;
  terminal authenticated GET checks were used as the safe substitute.

Boundary:

CP000006 is docs/local draft work only. It did not generate a package, run
release packaging, install validators, tag, publish, push images, create
retained evidence, contact external parties, change consumer statuses, or make
release-readiness, validator-clean, compliance, adoption, consumer, final-root,
hosted-service, production, vendor, hardware, SLA, uptime, or ETA-quality
claims.

## CP000007 Review

Status: complete for bounded rc1 hardening closeout.

Final validation:

| Command | Exit | Status | Result |
| --- | ---: | --- | --- |
| `git status --short` | 0 | needs_review | Primary checkout remained dirty with approved Phase 71/72 code, docs, script, and wiki changes plus untracked Phase 72 docs. |
| `git diff --check` | 0 | passed | No whitespace errors. |
| `make check` | 0 | passed | Lightweight no-network/no-Docker/no-validator-install checks passed. |
| `make validate` | 0 | passed | Validation smoke passed; pinned validator tooling check passed in the primary environment. |
| `make test` | 0 | passed | `go test ./...` passed. |
| `RUN_LOCAL_APP=true make release-candidate-check` | 0 | needs_review | Wrote `.cache/release-candidate-check/20260512T105655Z`; summary `overall_status=needs_review`, `git_clean=needs_review`, `local_app_five_feeds=passed`, `validators_check=passed`, `claim_audit=passed`, `release_package_audit=not_checked`. |
| `make agency-app-down` | 0 | passed | Local app containers and network were removed after the local-app diagnostic. |
| `make external-connection-check` | 0 | passed | Connector manifests remain sidecar/manifest/conformance bounded. |
| `make adapter-conformance` | 0 | passed | `testdata/adapter-conformance` suite passed. |
| `make test-connector-examples` | 0 | passed | `go test ./examples/connectors/...` passed. |
| `make audit-product-acceptance` | 0 | passed | Product acceptance audit passed. |
| `make audit-final-claim-review` | 0 | passed | Final claim review audit passed. |
| `docker compose -f deploy/docker-compose.yml config` | 0 | passed | Compose config rendered successfully. |
| `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` plus exact seven-target prepared-only check | 0 | passed | Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land remained exactly `prepared`. |
| `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum` | 0 | passed | Empty output; protected evidence, migration, and module paths stayed clean. |
| `docker compose -f deploy/docker-compose.yml --profile app ps` | 0 | passed | No Open Transit RT app containers remained after cleanup. |

Remaining blockers and deferrals:

- Primary checkout remains dirty, so this is not a clean tagged source state.
- The CP000002 temporary evaluator full `make release-candidate-check` blocker
  is still recorded: missing pinned static GTFS validator tooling in that
  evaluator snapshot.
- CP000007 primary `RUN_LOCAL_APP=true make release-candidate-check` exited `0`
  but remained `needs_review` because the checkout was dirty and release-package
  auditing was not enabled.
- Release package, artifact checksums, SBOM/provenance, GitHub release,
  published image, and tag remain `None`/`not_checked`.
- Browser automation with safe bearer-header support was unavailable in
  CP000004; terminal authenticated checks remain the recorded safe substitute.

Boundary:

Phase 72 is complete for bounded rc1 hardening review only. It did not tag,
publish, package, push images, create retained evidence, contact external
parties, change consumer statuses, or prove release readiness, validator-clean
feeds, production readiness, CAL-ITP/Caltrans compliance, consumer
submission/review/acceptance/ingestion/listing/display, agency adoption/
approval, final-root readiness, hosted SaaS, SLA/uptime, vendor compatibility,
hardware certification, production AVL reliability, or production-grade ETA
quality.

## Next Checkpoint

Phase 73 -- Checkpoint 000006: close agency UI acceptance review after CP000005
completion.

Expected next work:

- Continue Phase 73 from CP000006 after CP000005 completion and Master review.
- Preserve all claim boundaries and protected paths.
- Do not tag, publish, package, or create retained evidence unless separately
  authorized.
