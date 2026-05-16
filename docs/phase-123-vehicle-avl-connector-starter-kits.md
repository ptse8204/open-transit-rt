# Phase 123 -- Vehicle AVL Connector Starter Kits

## Goal

Add vehicle/AVL connector starter kits for CSV, GPS polling, webhook-sidecar,
and synthetic-only paths, all redaction-first and no-vendor-claim.

Phase 123 is not a vendor compatibility, hardware certification, production
AVL reliability, production readiness, compliance, consumer acceptance,
SLA/uptime, or hosted-service proof phase.

## Current Repo Context

- Existing examples cover CSV replay, HTTP/GPS-style polling, monitoring,
  validator, and prediction sidecar patterns.
- `cmd/avl-vendor-adapter` and `internal/avladapter` provide a synthetic
  vendor payload transformation and optional send path.
- Device/AVL docs explain `/v1/telemetry`, token lifecycle, simulator, and
  troubleshooting boundaries.
- Phase 123 should close the missing webhook-sidecar starter-kit gap and
  document the starter-kit matrix.

## Scope

- Add or reconcile a webhook-sidecar starter kit with synthetic fixtures and
  tests.
- Add redaction-first connector starter-kit documentation tying CSV, GPS
  polling, webhook sidecar, and synthetic-only paths together.
- Keep all examples disabled by default, synthetic/local, and explicit that
  they do not prove vendor compatibility or production AVL reliability.

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

- Webhook-sidecar vehicle/AVL connector starter kit or equivalent audit.
- Starter-kit docs covering CSV, GPS polling, webhook-sidecar, and
  synthetic-only paths.
- Tests and connector manifest validation.
- `docs/handoffs/phase-123.md`
- Source-of-truth status updates for Phase 123 closeout.

## Implementation Plan

1. Add this Phase 123 plan and commit checkpoint 000001.
2. Implement the missing webhook-sidecar starter kit and docs alignment.
3. Run relevant connector, GTFS-RT, claim-boundary, and baseline validation;
   patch repo-caused failures.
4. Close Phase 123 with handoff/status docs and continue immediately to
   Phase 124.

## Checkpoint Plan

- `Phase 123 -- Checkpoint 000001: add vehicle avl connector starter kits plan`
- `Phase 123 -- Checkpoint 000002: implement or audit primary scoped work`
- `Phase 123 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 123 -- Checkpoint 000004: close vehicle avl connector starter kits review`

## Checkpoint Report -- 000001

Checkpoint:
Phase 123 -- Checkpoint 000001: add vehicle AVL connector starter kits plan.

Goal status:
Active. Phase 122 is closed and Phase 123 has started.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, GTFS-RT Domain, Connector, Claim-Boundary,
Security/Auth, Data/Migration, Documentation / IA, Web Design Skill, Release,
and Install Confidence roles are simulated by the Master Agent.

Changed files:
`docs/phase-123-vehicle-avl-connector-starter-kits.md`.

Validation run:
Initial inspection reviewed the Phase 123 roadmap prompt, existing connector
examples, connector manifest validator, AVL adapter command, device AVL guide,
and protected/claim boundaries.

Blocked checks:
Implementation, tests, connector validation, and closeout validation are
scheduled for later Phase 123 checkpoints.

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
The plan keeps starter kits disabled-by-default, redaction-first, and
synthetic/local. It does not change route auth, CSRF behavior, credential
handling, token handling, public exposure, private payload handling, or
operator command behavior.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change is
planned.

Release/publication status:
The public rc1 prerelease remains published. Phase 123 does not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
Phase 118 Web Design Skill artifact remains complete. Phase 123 does not plan
visual UX changes.

Master review:
Approved. The plan closes the webhook-sidecar starter-kit gap while preserving
redaction, connector, and no-vendor-claim boundaries.

Required edits:
Commit checkpoint 000001, then implement the scoped starter-kit work.

Decision:
Proceed to checkpoint 000001 validation and commit.

Next checkpoint:
Phase 123 -- Checkpoint 000002: implement or audit primary scoped work.

## Checkpoint Report -- 000002

Checkpoint:
Phase 123 -- Checkpoint 000002: implement or audit primary scoped work.

Goal status:
Active. Phase 123 implemented the scoped connector starter-kit work and is
ready for full checkpoint validation.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, GTFS-RT Domain, Connector, Claim-Boundary,
Security/Auth, Data/Migration, Documentation / IA, Web Design Skill, Release,
and Install Confidence roles are simulated by the Master Agent.

Changed files:
`examples/connectors/telemetry-webhook-sidecar/README.md`,
`examples/connectors/telemetry-webhook-sidecar/connector.json`,
`examples/connectors/telemetry-webhook-sidecar/fixtures/webhook.json`,
`examples/connectors/telemetry-webhook-sidecar/main.go`,
`examples/connectors/telemetry-webhook-sidecar/main_test.go`,
`docs/connectors/vehicle-avl-starter-kits.md`,
`docs/tutorials/device-avl-integration.md`,
`internal/connectors/examples_test.go`, and
`internal/connectors/registry_test.go`.

Implementation summary:
Added a disabled-by-default synthetic webhook sidecar example that reads a
local fixture, normalizes webhook-style AVL observations through the shared
telemetry connector SDK, and emits dry-run summaries with `dry_run=true` and
`network_send=false`. Added a starter-kit matrix covering CSV replay, GPS/API
polling, webhook sidecars, synthetic-only telemetry, and vendor-shaped payload
transform examples. Updated connector registry tests for the sixth committed
example manifest.

Validation run:
`gofmt` on touched Go files passed. `python3 -m json.tool` passed for the new
connector manifest and fixture. `go test ./internal/connectors
./examples/connectors/...` passed. `make test-connector-examples` passed.
`make external-connection-check` passed. `git diff --check` passed.
`scripts/check-consumer-tracker.sh` passed.

Blocked checks:
None for this checkpoint. Full repo validation is scheduled for checkpoint
000003.

Protected path status:
`git status --short -- docs/evidence/consumer-submissions
docs/evidence/captured db/migrations go.mod go.sum` returned no output. No
protected evidence path, migration, or module file was modified.

Consumer tracker status:
`scripts/check-consumer-tracker.sh` reported exactly seven prepared-only
targets.

Claim-boundary status:
The starter kit and docs explicitly avoid vendor compatibility, hardware
certification, production AVL reliability, production readiness, compliance,
consumer acceptance, hosted service, SLA/uptime, and ETA-quality claims.

Security/auth status:
The webhook sidecar is local fixture input only, has no listener, no send path,
no token handling, no private endpoint, and no raw private payload in the
manifest. Redaction fields cover authorization, signatures, raw payloads,
private endpoints, and webhook secrets.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change was made.

Release/publication status:
The public rc1 prerelease remains published. Phase 123 did not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
Phase 118 Web Design Skill artifact remains complete. Phase 123 did not make
visual UX changes.

Master review:
Approved for full validation. The implementation closes the webhook-sidecar
starter-kit gap without expanding beyond synthetic, disabled-by-default, and
redaction-first connector boundaries.

Required edits:
Run checkpoint 000003 full validation and patch any repo-caused failures.

Decision:
Proceed to checkpoint 000002 commit.

Next checkpoint:
Phase 123 -- Checkpoint 000003: run validation and patch required gaps.

## Checkpoint Report -- 000003

Checkpoint:
Phase 123 -- Checkpoint 000003: run validation and patch required gaps.

Goal status:
Active. Phase 123 implementation passed full validation after patching stale
admin JSON shape test expectations.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, GTFS-RT Domain, Connector, Claim-Boundary,
Security/Auth, Data/Migration, Documentation / IA, Web Design Skill, Release,
and Install Confidence roles are simulated by the Master Agent.

Changed files:
`cmd/agency-config/main_test.go` and this phase report.

Validation run:
`git status --short` was clean at checkpoint start. `git diff --check`
passed. `python3 -m json.tool docs/evidence/consumer-submissions/status.json`
passed. `scripts/check-consumer-tracker.sh` passed. `make check` passed.
`make validate` passed. `docker compose -f deploy/docker-compose.yml config`
passed. `make audit-product-acceptance` passed. `make
audit-final-claim-review` passed. `make external-connection-check` passed.
`make adapter-conformance` passed. `make test-connector-examples` passed.
`make gtfsrt-conformance` passed. `make test` initially failed in
`cmd/agency-config` because connector hub/workbench tests still expected five
committed example manifests. After updating those expectations to include
`example.telemetry-webhook-sidecar`, `go test ./cmd/agency-config` passed and
`make test` passed.

Blocked checks:
None. The full Phase 123 validation set passed after the repo-caused test
expectation patch.

Protected path status:
`git status --short -- docs/evidence/consumer-submissions
docs/evidence/captured db/migrations go.mod go.sum` returned no output. No
protected evidence path, migration, or module file was modified.

Consumer tracker status:
`scripts/check-consumer-tracker.sh` reported exactly seven prepared-only
targets.

Claim-boundary status:
Claim audits passed. The new connector example remains a synthetic adapter
contract and does not claim vendor compatibility, hardware certification,
production AVL reliability, production readiness, compliance, consumer
acceptance, hosted service, SLA/uptime, production-grade ETA quality, or
real-world ETA accuracy.

Security/auth status:
No route auth, CSRF behavior, credential handling, token handling, public
exposure, private payload handling, or operator command behavior was changed.
The admin tests now assert that the new manifest is present in private
connector hub/workbench JSON and HTML while preserving no-form, no-command,
no-public-route, disabled-by-default, and fail-closed boundaries.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change was made.

Release/publication status:
The public rc1 prerelease remains published. Phase 123 did not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
Phase 118 Web Design Skill artifact remains complete. Phase 123 did not make
visual UX changes.

Master review:
Approved. The stale tests were corrected to align with the new committed
connector manifest, and all relevant gates passed afterward.

Required edits:
Close Phase 123 with handoff and status updates.

Decision:
Proceed to checkpoint 000003 commit.

Next checkpoint:
Phase 123 -- Checkpoint 000004: close vehicle avl connector starter kits review.
