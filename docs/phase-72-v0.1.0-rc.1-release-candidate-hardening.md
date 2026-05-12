# Phase 72 -- v0.1.0-rc.1 Release Candidate Hardening

## Status

Complete for the bounded `v0.1.0-rc.1` hardening review. This is not a release
tag, release-ready pass, validator-clean claim, production-readiness claim, or
public launch.

Checkpoint 000001 is complete for documentation-only planning and status
alignment. Checkpoint 000002 is complete for the approved temporary local
evaluator gate.

CP000002 recorded the primary checkout as dirty with approved docs changes,
created a clean local-only evaluator snapshot from that state, ran the approved
evaluator commands there, and found one local pinned-validator tooling blocker
in `make release-candidate-check`. It did not fix code, install validators, run
out-of-scope app/browser/feed-fetch/release/evidence/connector/publishing
commands, draft release notes, create retained evidence, or change consumer
tracker statuses.

Checkpoint 000003 is complete for diagnostic hardening only. It forces the
release-candidate validator tooling row to run `scripts/check-validators.sh`
with `VALIDATOR_TOOLING_MODE=pinned`, preserves dry-run `not_checked`
behavior, records validator-tooling failures as `blocker` rows with the exact
final `scripts/check-validators.sh` blocker line in `summary.json` and
`summary.md`, and keeps full command output in `check-log.txt`. It did not
install validators, did not run `make validators-install`, and did not convert
missing pinned validator tooling into release-ready, validator-clean, or
compliance evidence.

Checkpoint 000004 is complete for local agency-app and Operations Console
verification. `make agency-app-up` passed, the nine required private
Operations Console routes returned authenticated local `200` responses, the
five required public feed paths returned anonymous local `200` responses, and
`feeds.json` parsed with `python3 -m json.tool`. Browser automation with safe
bearer-header support was unavailable in this session, so terminal
authenticated GET checks using an `Authorization` header were recorded as the
safe substitute without exposing the admin token. The CP000004
`RUN_LOCAL_APP=true` release-candidate diagnostic passed its
`local_app_five_feeds` and `validators_check` rows, while the overall
diagnostic remained `needs_review` because release-package auditing was
intentionally not enabled. The local app was stopped with
`make agency-app-down`.

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

## Goal

Prepare Open Transit RT for a truthful `v0.1.0-rc.1` maintainer review by
running one bounded, repeatable hardening sequence from planning through
validation, browser walkthrough, local diagnostics, claim audit, and closeout.

The phase should make blockers explicit, preserve the existing product
boundaries, and keep all release-candidate outputs local/private until a
maintainer separately authorizes tagging, publishing, or retained evidence.

## Non-Goals

- No release tag, GitHub release, package upload, registry push, or hosted
  deployment.
- No final release notes, tag, package, or public release; the CP000006 draft
  is local and pre-tag only.
- No retained evidence writes.
- No consumer tracker status changes.
- No external party contact, portal automation, or consumer submission.
- No final public root proof, agency approval proof, or off-host production
  claim.
- No database migrations, schema changes, public feed URL changes, telemetry
  ingest changes, GTFS-RT semantic changes, prediction adapter changes, or
  connector manifest schema changes.
- No claim of CAL-ITP/Caltrans compliance, consumer submission/review/
  acceptance/listing/display/ingestion, agency adoption or approval,
  final-root readiness, hosted SaaS, paid support, SLA/uptime, production
  readiness, production multi-tenant hosting, vendor compatibility, hardware
  certification, production AVL reliability, production-grade ETA quality, or
  public launch completion.

## Checkpoint Sequence

1. `CP000001 -- add release-candidate hardening plan`
   - Add this Phase 72 plan.
   - Add `docs/handoffs/phase-72.md`.
   - Narrowly update release-candidate readiness, roadmap status, and latest
     handoff docs with the initial Phase 72 launch state.
   - Do not run full validation or draft release notes.
2. `CP000002 -- run clean-checkout local evaluator gate`
   - Run the clean-checkout source-state review and local evaluator gate as
     allowed by the local environment.
   - Record exact blockers without converting skipped or failed checks into
     claims.
   - Complete: evaluator checks passed except `make release-candidate-check`,
     which is blocked by missing pinned static GTFS validator tooling.
3. `CP000003 -- harden release-candidate diagnostics and blockers`
   - Complete: hardened release-candidate diagnostic wording, blocker
     recording, and environment/tooling blocker handling found by CP000002.
   - Validator tooling diagnostics now force pinned mode even if the ambient
     shell sets `VALIDATOR_TOOLING_MODE=stub`.
   - Missing or misconfigured pinned validator tooling remains a local
     environment/tooling blocker, not validator-clean evidence.
4. `CP000004 -- verify browser-first agency operations walkthrough`
   - Complete: local app startup passed, all nine required private Operations
     Console routes returned authenticated local `200`, all five public feed
     paths returned anonymous local `200`, and in-app Browser unavailability
     was recorded with terminal authenticated checks as the safe substitute.
5. `CP000005 -- verify connector and adapter conformance gates`
   - Complete: `make external-connection-check`, `make adapter-conformance`,
     and `make test-connector-examples` all passed in the primary checkout.
   - No connector/adapter gate blocker required a code, schema, fixture, or
     example patch.
   - Preserve protected paths, consumer tracker statuses, public route
     boundaries, and claim boundaries.
6. `CP000006 -- prepare rc1 release notes and known blockers`
   - Complete: drafted local pre-tag release notes and a known blockers matrix
     at `docs/release-notes-v0.1.0-rc.1-draft.md`.
   - Did not generate or audit a release package.
   - Stated `None` for unchanged migration, security, dependency, operations,
     evidence, and claim sections.
7. `CP000007 -- close rc1 hardening review`
   - Complete: final focused checks and claim audit ran as local diagnostics.
   - Updated Phase 72 status and handoff with exact results, blockers, and
     intentional deferrals.
   - Did not tag, publish, package, or create retained evidence.

## Browser Walkthrough Routes

Checkpoint 000004 should review these required private/admin routes when local
app startup is available:

- `/admin/operations`
- `/admin/operations/gtfs-import`
- `/admin/operations/feed-health`
- `/admin/operations/gtfs-quality`
- `/admin/operations/validation-health`
- `/admin/operations/devices`
- `/admin/operations/telemetry`
- `/admin/operations/telemetry-simulator`
- `/admin/operations/connectors`
- `/admin/operations/maintenance`

Supplemental private routes such as setup wizard, readiness, Connector Tests,
and Help may be reviewed when useful, but they are not substitutes for the
required walkthrough list above.

The same checkpoint should fetch these local public feed paths when the app is
running:

- `/public/feeds.json`
- `/public/gtfs/schedule.zip`
- `/public/gtfsrt/vehicle_positions.pb`
- `/public/gtfsrt/trip_updates.pb`
- `/public/gtfsrt/alerts.pb`

## Validation Matrix

| Check | Planned checkpoint | CP000001 status | Current result | Boundary |
| --- | --- | --- | --- | --- |
| Primary `git status --short` | CP000002 | not_run_in_cp000001 | passed; primary checkout dirty with approved docs changes | Source-state signal only |
| Primary `git describe --tags --always --dirty` | CP000002 | not_run_in_cp000001 | passed; `b0e354d-dirty` | Source metadata only |
| Primary `git rev-parse HEAD` | CP000002 | not_run_in_cp000001 | passed; `b0e354d29b14ba5c5ab0afab3269ec5eb3405fa7` | Source metadata only |
| Primary protected source-path review | CP000002 | not_run_in_cp000001 | passed; empty output for evidence, migrations, and module paths | Scope guard only |
| Evaluator `git status --short` | CP000002 | not_run_in_cp000001 | passed; clean local-only evaluator snapshot | Source-state signal only |
| Evaluator `git describe --tags --always --dirty` | CP000002 | not_run_in_cp000001 | passed; `83e0183` | Source metadata only |
| Evaluator `git rev-parse HEAD` | CP000002 | not_run_in_cp000001 | passed; `83e01831c035404e53e921b19c6b7e89e9541816` | Source metadata only |
| `make check` | CP000002 and CP000007 | not_run_in_cp000001 | passed in evaluator and CP000007 primary | Local repo validation only |
| `scripts/release-candidate-check.sh --dry-run` | CP000002 | not_run_in_cp000001 | passed in evaluator; wrote `.cache/release-candidate-check/20260512T062130Z` | Private diagnostic rehearsal only |
| `make release-candidate-check` | CP000002 or CP000003 | not_run_in_cp000001 | blocker in evaluator; missing pinned static GTFS validator tooling | Private `.cache` diagnostic only |
| `scripts/bootstrap-dev.sh --check` | CP000002 | not_run_in_cp000001 | passed in evaluator | Local setup preflight only |
| `make validate` | CP000003 or CP000007 | not_run_in_cp000001 | passed in CP000007 primary | Validator tooling signal only |
| `make test` | CP000002, CP000003, or CP000007 | not_run_in_cp000001 | passed in evaluator and CP000007 primary | Local test signal only |
| `docker compose -f deploy/docker-compose.yml config` | CP000002 and CP000007 | not_run_in_cp000001 | passed in evaluator and CP000007 primary | Compose syntax/config signal only |
| `RUN_LOCAL_APP=true make release-candidate-check` | CP000004 and CP000007 | not_run_in_cp000001 | CP000007 command exited 0 with `overall_status=needs_review`; output `.cache/release-candidate-check/20260512T105655Z` | Local app/feed diagnostic only |
| Five local public feed fetches | CP000004 | not_run_in_cp000001 | not_run_in_cp000002 by scope | Local fetch signal only |
| Browser walkthrough routes | CP000004 | not_run_in_cp000001 | not_run_in_cp000002 by scope | Private UI review only |
| `make telemetry-simulator` | CP000004 or CP000005 | not_run_in_cp000001 | not_run_in_cp000002 by scope | Synthetic telemetry only |
| `make external-connection-check` | CP000005 and CP000007 | not_run_in_cp000001 | passed in primary; connector manifests remain sidecar/manifest/conformance bounded | Synthetic/local connector boundary only |
| `make adapter-conformance` | CP000005 and CP000007 | not_run_in_cp000001 | passed in primary; `testdata/adapter-conformance` suite passed | Synthetic adapter conformance only |
| `make test-connector-examples` | CP000005 and CP000007 | not_run_in_cp000001 | passed in primary; `go test ./examples/connectors/...` passed | Synthetic examples only |
| `make audit-final-claim-review` | CP000007 | not_run_in_cp000001 | passed in CP000007 primary | Unsupported-claim guard only |
| Consumer tracker prepared-only review | CP000007 | not_run_in_cp000001 | passed in CP000007 primary | Status drift guard only |
| Protected-path diff review | CP000007 | not_run_in_cp000001 | passed in CP000007 primary for evidence, migration, and module paths | Scope guard only |

## CP000004 Agency-App Walkthrough Results

CP000004 changed only Phase 72 status/handoff docs. It did not change code,
scripts, migrations, module files, route contracts, telemetry/protobuf/
prediction/connector schemas, release notes, tags, packages, publishing files,
protected evidence paths, or consumer submission tracker files.

Baseline source and protected-path review:

| Check | Result |
| --- | --- |
| `git status --short` | Dirty primary checkout with pre-existing modified code/docs/scripts/wiki files and untracked Phase 72 docs; not treated as clean. |
| Protected paths status | Empty output for `docs/evidence/captured`, `docs/evidence/consumer-submissions/status.json`, `docs/evidence/consumer-submissions/current`, `docs/evidence/consumer-submissions/artifacts`, `docs/evidence/consumer-submissions/packets`, `db/migrations`, `go.mod`, and `go.sum`. |

Local app lifecycle:

| Command | Result |
| --- | --- |
| `make agency-app-up` | passed; started local Postgres/PostGIS, built the local app image, applied migrations, seeded demo agency/device records, started services and proxy, imported sample GTFS, bootstrapped publication metadata, and reported `http://localhost:8080` running. |
| `make agency-app-down` | passed; stopped and removed local app containers and the local Compose network. |

Authentication and browser boundary:

- Admin token was generated only into a local shell variable, parsed from the
  `token=` line, used only in an `Authorization` header for terminal GET
  checks, and unset after each shell block.
- The token was not printed, stored, pasted into URLs, or pasted into
  browser-visible fields.
- Public feed fetches used no authentication.
- Browser automation with safe bearer-header support was unavailable in this
  session. Terminal authenticated GET checks were used as the safe substitute.

Unauthenticated admin boundary check:

| Route | HTTP status | Result |
| --- | ---: | --- |
| `/admin/operations` without an admin token | 401 | Admin route stayed non-public. |

Private Operations Console route checks:

| Route | HTTP status | Content type | Bytes | Page-specific text |
| --- | ---: | --- | ---: | --- |
| `/admin/operations` | 200 | `text/html; charset=utf-8` | 40873 | `Agency Operations Cockpit / Start Here` present |
| `/admin/operations/gtfs-import` | 200 | `text/html; charset=utf-8` | 12045 | `GTFS Import` present |
| `/admin/operations/feed-health` | 200 | `text/html; charset=utf-8` | 23661 | `Feed Health` present |
| `/admin/operations/gtfs-quality` | 200 | `text/html; charset=utf-8` | 16216 | `GTFS Quality` present |
| `/admin/operations/validation-health` | 200 | `text/html; charset=utf-8` | 12310 | `Validator Health` present |
| `/admin/operations/devices` | 200 | `text/html; charset=utf-8` | 11323 | `Devices` present |
| `/admin/operations/telemetry` | not_run_in_cp000004 | not recorded | not recorded | Superseded by Phase 73 CP000003 route-list alignment and local `200` check. |
| `/admin/operations/telemetry-simulator` | 200 | `text/html; charset=utf-8` | 13698 | `Telemetry Simulator` present |
| `/admin/operations/connectors` | 200 | `text/html; charset=utf-8` | 16611 | `Connector` present |
| `/admin/operations/maintenance` | 200 | `text/html; charset=utf-8` | 15078 | `Maintenance` present |

The first route-check script attempt used zsh-incompatible Bash array
expansion and exited before producing route results. The corrected Bash run
above is the recorded CP000004 route result and was not a product blocker.

Anonymous public feed fetches:

| Path | HTTP status | Content type | Bytes | Additional check |
| --- | ---: | --- | ---: | --- |
| `/public/feeds.json` | 200 | `application/json` | 2477 | `python3 -m json.tool` parse passed |
| `/public/gtfs/schedule.zip` | 200 | `application/zip` | 1913 | fetched without auth |
| `/public/gtfsrt/vehicle_positions.pb` | 200 | `application/x-protobuf` | 15 | fetched without auth |
| `/public/gtfsrt/trip_updates.pb` | 200 | `application/x-protobuf` | 15 | fetched without auth |
| `/public/gtfsrt/alerts.pb` | 200 | `application/x-protobuf` | 135 | fetched without auth |

Optional local release-candidate diagnostic:

| Command | Result | Boundary |
| --- | --- | --- |
| `OUTPUT_DIR=.cache/phase-72-cp000004/release-candidate-check FORCE=true RUN_LOCAL_APP=true scripts/release-candidate-check.sh` | passed command; wrote private diagnostics to `.cache/phase-72-cp000004/release-candidate-check`; `overall_status=needs_review`; `validators_check=passed`; `local_app_five_feeds=passed`; `release_package_audit=not_checked` | Private local diagnostic only; not release readiness, validator-clean evidence, package evidence, compliance evidence, or production evidence. |

CP000004 intentionally did not run or change:

- `make telemetry-simulator`
- connector/adapter gates scheduled for CP000005
- release-package generation or audit
- release notes
- tags, publishing, image pushes, retained evidence, external contacts, or
  consumer submissions

These results are local app and private-route verification only. They do not
prove release readiness, validator-clean feeds, production readiness,
CAL-ITP/Caltrans compliance, consumer submission/review/acceptance/ingestion/
listing/display, agency adoption/approval, final-root readiness, hosted SaaS,
SLA/uptime, vendor compatibility/hardware certification, production AVL
reliability, or production-grade ETA quality.

## CP000005 Connector And Adapter Conformance Results

CP000005 verified local synthetic connector and adapter conformance gates only.
It changed no connector/adaptor code, schemas, fixtures, examples, protected
evidence paths, consumer tracker statuses, migrations, module files, public
routes, release artifacts, tags, packages, or retained evidence.

Baseline source and protected-path review:

| Check | Result |
| --- | --- |
| `git status --short` | Dirty primary checkout with existing modified code/docs/scripts/wiki files and untracked Phase 72 docs; not treated as clean. |
| Protected paths status | Empty output for `docs/evidence/captured`, `docs/evidence/consumer-submissions/status.json`, `docs/evidence/consumer-submissions/current`, `docs/evidence/consumer-submissions/artifacts`, `docs/evidence/consumer-submissions/packets`, `db/migrations`, `go.mod`, and `go.sum`. |
| Consumer tracker | `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` passed; exact seven-target prepared-only check passed for Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land. |

Connector and adapter gate results:

| Command | Exit | Status | Result |
| --- | ---: | --- | --- |
| `make external-connection-check` | 0 | passed | JSON-validates connector manifests/fixtures, ran `go test ./internal/connectors`, ran `go test ./examples/connectors/...`, and reported connector manifests remain sidecar/manifest/conformance bounded. |
| `make adapter-conformance` | 0 | passed | `go run ./cmd/adapter-conformance run --suite testdata/adapter-conformance` passed the synthetic suite. |
| `make test-connector-examples` | 0 | passed | `go test ./examples/connectors/...` passed for committed synthetic connector examples. |

Boundary:

- Results are local synthetic connector/adaptor checks only.
- No external systems, vendor endpoints, consumer portals, agency systems, or
  retained evidence paths were contacted or written.
- Passing these gates does not prove vendor compatibility, hardware
  certification, production AVL reliability, production readiness,
  CAL-ITP/Caltrans compliance, consumer acceptance, release readiness, or
  production-grade ETA quality.

## Initial Known Blockers Matrix

CP000001 rows were intentionally `not_run_in_cp000001` because Checkpoint
000001 was planning/status only. CP000002 updates the current result where an
approved evaluator or guard command ran.

| Area | CP000001 status | CP000002 status | Current note |
| --- | --- | --- | --- |
| Clean checkout source review | not_run_in_cp000001 | passed | Primary was dirty; evaluator snapshot was clean at local commit `83e01831c035404e53e921b19c6b7e89e9541816`. |
| `make check` | not_run_in_cp000001 | passed | Evaluator command exited 0. |
| Full release-candidate diagnostic | not_run_in_cp000001 | blocker | CP000002 evaluator `make release-candidate-check` exited 2 because pinned static GTFS validator tooling was missing. CP000004 primary `RUN_LOCAL_APP=true` diagnostic exited 0 with `overall_status=needs_review` because release-package auditing was intentionally not enabled. |
| Local app startup | not_run_in_cp000001 | not_run_in_cp000002 | CP000004 `make agency-app-up` passed and `make agency-app-down` stopped the local app. |
| Browser Operations Console walkthrough | not_run_in_cp000001 | not_run_in_cp000002 | CP000004 browser automation with safe bearer-header support was unavailable; terminal authenticated GET checks substituted safely and all nine required routes returned local `200`. |
| Five public feed fetches | not_run_in_cp000001 | not_run_in_cp000002 | CP000004 anonymous local fetches returned `200` for all five required public paths; `feeds.json` parsed with `python3 -m json.tool`. |
| Static and realtime validator execution | not_run_in_cp000001 | not_run_in_cp000002 | `make validate` and validator install were explicitly out of CP000002 scope. |
| Java/pinned validator tooling availability | not_run_in_cp000001 | blocker | Exact CP000002 evaluator blocker: `missing pinned tooling: static GTFS validator not installed at /private/var/folders/_g/bvzl9cms7cx1d0wdpc981n9w0000gn/T/tmp.BA6ULYuxdo/open-transit-rt-evaluator/.cache/validators/gtfs-validator-7.1.0-cli.jar; run make validators-install`. CP000004 primary optional diagnostic reported `validators_check=passed`; this is local tooling status only. |
| Docker daemon and Compose availability | not_run_in_cp000001 | passed | `scripts/bootstrap-dev.sh --check` and `docker compose -f deploy/docker-compose.yml config` passed in evaluator. |
| Telemetry simulator | not_run_in_cp000001 | not_run_directly_in_cp000002 | Only the release-candidate diagnostic's dry-run helper coverage ran; standalone telemetry simulator command was out of scope. |
| External connector checks | not_run_in_cp000001 | not_run_in_cp000002 | CP000005 `make external-connection-check` exited 0 in the primary checkout. |
| Adapter conformance | not_run_in_cp000001 | not_run_in_cp000002 | CP000005 `make adapter-conformance` and `make test-connector-examples` exited 0 in the primary checkout. |
| Release notes draft | not_run_in_cp000001 | not_run_in_cp000002 | CP000006 added local pre-tag draft `docs/release-notes-v0.1.0-rc.1-draft.md`. |
| Local release package diagnostics | not_run_in_cp000001 | not_run_in_cp000002 | Not run in CP000006 or CP000007; release package, checksum, SBOM, and provenance remain `None`/`not_checked`. |
| Final claim audit | not_run_in_cp000001 | not_run_directly_in_cp000002 | CP000007 `make audit-final-claim-review` passed. |
| Consumer tracker prepared-only audit | not_run_in_cp000001 | passed | CP000007 exact seven-target prepared-only check passed. |
| Protected-path diff audit | not_run_in_cp000001 | passed | CP000007 protected path status output was empty for evidence, migration, and module paths. |

## CP000007 Closeout Results

CP000007 closed the Phase 72 hardening review with exact local diagnostic
results. It did not create a release tag, package, image, checksum, SBOM,
provenance file, retained evidence, consumer submission, or consumer status
change.

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
| `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` plus exact seven-target prepared-only check | 0 | passed | Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land remain exactly `prepared`. |
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

These closeout results are local diagnostics only. They do not prove release
readiness, validator-clean feeds, production readiness, CAL-ITP/Caltrans
compliance, consumer submission/review/acceptance/ingestion/listing/display,
agency adoption/approval, final-root readiness, hosted SaaS, SLA/uptime, vendor
compatibility, hardware certification, production AVL reliability, or
production-grade ETA quality.

## CP000003 Diagnostic-Hardening Results

CP000003 changed only the approved release-candidate diagnostic script, focused
tests, a narrow validation-health test-fake concurrency fix, and Phase
72/readiness/handoff docs. It did not install validators or run
`make validators-install`.

Implemented behavior:

- `validators_check` now runs
  `VALIDATOR_TOOLING_MODE=pinned scripts/check-validators.sh` inside
  `scripts/release-candidate-check.sh`.
- Ambient `VALIDATOR_TOOLING_MODE=stub` no longer converts the
  release-candidate validator tooling check into a pass.
- Dry-run output still records `validators_check` as `not_checked`.
- Validator-tooling failures remain `blocker` rows and keep the full
  diagnostic exit nonzero.
- `summary.json` and `summary.md` include the exact final non-empty blocker
  line from `scripts/check-validators.sh`; `check-log.txt` retains the full
  command output.
- The release-candidate output directory remains limited to exactly
  `summary.json`, `summary.md`, `manifest.json`, `manifest.md`, and
  `check-log.txt`.
- Master-side `make test` initially exposed a test-only concurrent map write
  in `cmd/agency-config/main_test.go`; CP000003 fixed the shared
  validation-health test fakes with mutex guards and copied returned slices.
  The required edit did not change production validation semantics, runtime
  routes, migrations, module files, schemas, protected evidence paths, or
  consumer statuses.

Focused CP000003 validation in the primary dirty checkout:

| Command | Result | Boundary |
| --- | --- | --- |
| `sh -n scripts/release-candidate-check.sh scripts/check-validators.sh` | passed | Shell syntax only |
| `go test ./cmd/agency-config -run ReleaseCandidate` | passed | Focused release-candidate script tests |
| `make test-release-candidate-check` | passed | Make wrapper for focused tests |
| `OUTPUT_DIR=.cache/release-candidate-check-cp000003/dry-run FORCE=true scripts/release-candidate-check.sh --dry-run` | passed; `validators_check` was `not_checked`; five expected files only | Dry-run diagnostic only |
| Forced negative diagnostic with ambient `VALIDATOR_TOOLING_MODE=stub` and invalid `GTFS_VALIDATOR_PATH` | expected nonzero exit; `validators_check` was `blocker`; summary included actionable missing/misconfigured pinned tooling wording | Local environment/tooling blocker only |
| `make check` | passed | Lightweight local repo validation only |
| Initial master-side `make test` | blocker, then fixed | Exposed test-only validation-health fake shared-state writes; no production behavior was changed |
| `go test ./cmd/agency-config -run TestValidationHealthConcurrentRunAllDoesNotPanicOrLeak -count=100` | passed | Required-edit concurrency regression check |
| `go test ./cmd/agency-config -race -run TestValidationHealthConcurrentRunAllDoesNotPanicOrLeak -count=5` | passed | Required-edit race check |
| `go test ./cmd/agency-config` | passed | Package-level regression check |
| `make test` | passed after required-edit patch | Full local Go test suite |
| `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` | passed | Consumer tracker JSON parse only |
| Exact seven-target prepared-only tracker check | passed | Status drift guard only |
| `git status --short -- docs/evidence/captured docs/evidence/consumer-submissions db/migrations go.mod go.sum` | passed with empty output | Protected-path guard only |

These results are diagnostic-hardening evidence only. They are not validator
installation evidence, validator-clean feed evidence, release-ready evidence,
consumer submission evidence, production-readiness evidence, or
CAL-ITP/Caltrans compliance evidence.

## Protected Paths

Phase 72 did not edit or generate files under these protected paths. Future
work must not edit or generate files under these paths unless later explicit
authorization changes the scope:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/**`
- `docs/evidence/consumer-submissions/artifacts/**`
- `docs/evidence/consumer-submissions/packets/**`
- `db/migrations/**`
- `go.mod`
- `go.sum`
- code, scripts, `Makefile`, README, wiki, and site files during CP000001

Generated `.cache` diagnostics from Phase 72 checkpoints remain private local
diagnostics only. They are not retained evidence.

## Consumer Tracker Boundary

All seven consumer and aggregator targets remained `prepared` through CP000007.
Phase 72 reviewed the prepared-only status for drift, but it did not move any
target to submitted, under review, accepted, rejected, listed, ingested, or any
other stronger state.

Prepared packet workflows remained operator-review material only. Phase 72 did
not submit packets, automate portals, contact consumers, refresh target records
from external systems, or create target-originated evidence. Future packet
submission remains authorization-gated.

## Claim Boundaries

Phase 72 closeout wording may say:

- release-candidate hardening completed a bounded local review;
- CP000001 created a bounded plan and handoff;
- CP000002 through CP000007 ran local/private release-candidate diagnostics,
  walkthrough checks, connector/adaptor checks, release-note drafting, and
  closeout review as documented;
- local checks can support maintainer review when they are actually run.

Phase 72 closeout wording must not say that CP000001 or the phase itself
proves:

- the release-candidate gate passed;
- the checkout is clean or release-ready;
- validator-clean feeds;
- production readiness;
- CAL-ITP/Caltrans compliance;
- consumer submission, review, acceptance, ingestion, listing, or display;
- agency adoption, agency approval, final-root readiness, hosted SaaS, paid
  support, SLA/uptime, vendor compatibility, hardware certification,
  production AVL reliability, or production-grade ETA quality.

## Stop Conditions

Stop and hand back to the maintainer before continuing if any checkpoint would
require:

- editing protected evidence paths or consumer status JSON;
- contacting an external party or submitting to a consumer/aggregator;
- tagging a release, publishing a package, pushing an image, or creating a
  GitHub release;
- changing migrations, module dependencies, public feed URLs, telemetry ingest
  contracts, prediction adapter contracts, connector manifest schemas, or
  GTFS-RT semantics outside an explicitly approved blocker fix;
- retaining public evidence or making a stronger public claim;
- using a real agency, real device/vendor, final-root, credential, or consumer
  artifact without written authorization and redaction rules.

## Success Criteria

Phase 72 closeout criteria/results:

- the phase plan and handoff accurately track checkpoint status;
- CP000002 through CP000006 produce exact pass/blocker/draft results without
  overclaiming;
- CP000005 records exact local synthetic connector/adapter gate results without
  required code fixes;
- CP000006 drafts release notes with exact source/check/blocker inputs and
  unchanged sections explicitly marked;
- CP000007 records final validation, protected-path review, consumer tracker
  prepared-only review, claim audit status, and intentional deferrals;
- no protected evidence path or consumer tracker status is changed;
- no release tag, public artifact, external submission, retained evidence, or
  stronger compliance/production/consumer/adoption claim is created without
  separate authorization.
