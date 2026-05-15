# Phase 102 -- Device / AVL Fleet Onboarding V2

## Goal

Make fleet and device onboarding practical for less technical operators while
preserving existing safety boundaries: no token value exposure beyond the
existing one-time rotate/rebind result, no real credentials, no real vendor
payloads, no external contact, no evidence collection, no consumer status
changes, and no vendor compatibility, hardware certification, production AVL
reliability, compliance, hosted-service, SLA, consumer acceptance, adoption, or
production readiness claim.

## Current Surface

- `/admin/operations/devices` is the primary private device credential and
  device-to-vehicle binding surface.
- `/admin/operations/telemetry`, `/admin/operations/realtime`, and
  `/admin/operations/telemetry-simulator` provide related freshness, assignment,
  and synthetic send diagnostics.
- `internal/devices` stores hashed credentials and binding metadata; token
  values are only returned by rotate/rebind and are not listable.
- `internal/telemetry`, `cmd/telemetry-ingest`, `cmd/telemetry-simulator`,
  `cmd/avl-vendor-adapter`, and `scripts/device-onboarding.sh` already provide
  the local/synthetic onboarding mechanics.
- Existing docs cover token lifecycle and device/AVL integration, but the
  private console needs clearer operator-facing inventory review, bulk import
  planning, token lifecycle guidance, unknown-device triage, binding review, and
  technical-helper handoff guidance.

## Master-Approved Plan

1. Add the Phase 102 plan and checkpoint report.
2. Add small private Operations Console additions to the existing device page:
   inventory review, bulk onboarding plan, token lifecycle checklist,
   freshness/unknown-device triage, binding review, and technical-helper
   handoff rows.
3. Keep implementation metadata-only. Do not add migrations unless a safety
   issue is discovered; do not add token or raw payload rendering.
4. Add focused tests that prove non-admin visibility, admin controls, safe
   guidance text, no token/raw payload leakage, and no unsupported claims.
5. Update operator docs so device/AVL onboarding can be followed without using
   real credentials or implying vendor/hardware certification.
6. Run focused checks, then required baseline/code-change validation.
7. Close with `docs/handoffs/phase-102.md`, `docs/handoffs/latest.md`, and
   roadmap/status updates.

## Non-Goals

- No real AVL device or vendor integration.
- No browser collection of device token values.
- No live telemetry send from the private console.
- No external contact, portal use, or evidence collection.
- No protected-path writes.
- No consumer tracker edits.
- No migration or durable schema change unless required to fix a safety issue.

## Checkpoint Plan

- `Phase 102 -- Checkpoint 000001: add device / avl fleet onboarding v2 plan`
- `Phase 102 -- Checkpoint 000002: implement primary scoped work`
- `Phase 102 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 102 -- Checkpoint 000004: close device / avl fleet onboarding v2 review`

## Focused Validation Targets

- `go test ./cmd/agency-config -run 'Device|Telemetry|Realtime|OperationsNavigation|RouteTitles|Help'`
- `go test ./internal/devices ./internal/telemetry ./cmd/telemetry-ingest ./cmd/telemetry-simulator ./cmd/avl-vendor-adapter`
- `scripts/device-onboarding.sh help`
- `scripts/telemetry-simulator.sh --help`
- `scripts/telemetry-simulator.sh --list-scenarios`

Because this phase changes code/docs/tests, closeout also requires:

- `git status --short`
- `git diff --check`
- `make check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact prepared-only consumer tracker assertion
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`

## Checkpoint Report -- 000001

Checkpoint:
Phase 102 -- Checkpoint 000001: add device / avl fleet onboarding v2 plan.

Sub-agents used or simulated, including intended model level:
Real Context / Repo Truth Sub-Agent -- GPT-5.5 x-high, active for device/AVL
surface inspection. Real Planning Sub-Agent -- GPT-5.5 x-high, active for
checkpoint planning. Implementation, QA, UI/UX, Documentation / IA,
Claim-Boundary, Security/Auth, Data/Migration, and Release/Supply-Chain roles
are simulated by the Master Agent for this plan checkpoint. Master Agent --
GPT-5.5 x-high, current thread.

Changed files:
`docs/phase-102-device-avl-fleet-onboarding-v2.md`.

Validation run:
Initial `git status --short` before planning returned no output. After adding
the plan, `git status --short` showed only
`docs/phase-102-device-avl-fleet-onboarding-v2.md`; `git diff --check` passed;
`python3 -m json.tool docs/evidence/consumer-submissions/status.json
>/dev/null` passed; the exact prepared-only consumer tracker assertion passed;
and `git status --short -- docs/evidence/consumer-submissions
docs/evidence/captured db/migrations go.mod go.sum` returned no output.

Blocked checks:
Implementation, focused code tests, and closeout baseline checks are not yet
run because this checkpoint only approves the Phase 102 plan. Release-candidate
and package checks are out of scope for Phase 102.

Protected path status:
No protected evidence path is part of the plan. The plan forbids protected path
writes.

Consumer tracker status:
The consumer tracker is not part of the plan. The seven targets must remain in
order and `prepared`.

Claim-boundary status:
The plan explicitly forbids vendor compatibility, hardware certification,
production AVL reliability, compliance, hosted-service, SLA, consumer
acceptance, adoption, release readiness, and production readiness claims.

Security/auth status:
The plan preserves the existing role-gated private console, avoids browser
token collection, avoids raw payload display, and treats one-time token display
as limited to the existing admin rotate/rebind response.

Data/migration status:
No migration or durable schema change is planned.

Master review:
Approved. The smallest safe Phase 102 implementation is private guidance and
metadata-only review surfaces over existing device bindings and telemetry.

Required edits:
Add private device onboarding guidance, update tests, update docs, and record
validation results.

Decision:
Proceed to implementation checkpoint 000002 after plan validation and commit.

Next checkpoint:
Phase 102 -- Checkpoint 000002: implement primary scoped work.

## Checkpoint Report -- 000002

Checkpoint:
Phase 102 -- Checkpoint 000002: implement primary scoped work.

Sub-agents used or simulated, including intended model level:
Real Context / Repo Truth Sub-Agent -- GPT-5.5 x-high; real Planning
Sub-Agent -- GPT-5.5 x-high. Implementation, QA, UI/UX, Documentation / IA,
Claim-Boundary, Security/Auth, Data/Migration, and Release/Supply-Chain roles
were simulated by the Master Agent. Master Agent -- GPT-5.5 x-high, current
thread.

Changed files:
`cmd/agency-config/operations.go`;
`cmd/agency-config/operations_devices.go`;
`cmd/agency-config/main_test.go`;
`docs/tutorials/device-token-lifecycle.md`;
`docs/tutorials/device-avl-integration.md`;
`docs/phase-102-device-avl-fleet-onboarding-v2.md`.

Validation run:
`go test ./internal/devices ./internal/telemetry ./cmd/telemetry-ingest
./cmd/telemetry-simulator ./cmd/avl-vendor-adapter` passed. `go test
./cmd/agency-config -run
'Device|Telemetry|Realtime|OperationsNavigation|RouteTitles|Help'` initially
failed because new helper text included generic `authorization` wording on the
one-time-token POST page; the copy was patched to `request credential headers`
and the focused agency-config test then passed. `git diff --check` passed.
`scripts/device-onboarding.sh help` passed. `scripts/telemetry-simulator.sh
--help` passed. `scripts/telemetry-simulator.sh --list-scenarios` passed. `git
status --short -- docs/evidence/consumer-submissions docs/evidence/captured
db/migrations go.mod go.sum` returned no output. `python3 -m json.tool
docs/evidence/consumer-submissions/status.json >/dev/null` passed, and the
exact prepared-only consumer tracker assertion passed.

Blocked checks:
Full closeout baseline, `make validate`, `make test`, and docker compose
configuration validation are reserved for checkpoint 000003. Release-candidate,
package, publication, real AVL, external vendor, evidence, and consumer action
checks remain out of scope for Phase 102.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched. The
protected-path status check returned no output.

Consumer tracker status:
`docs/evidence/consumer-submissions/status.json` was not edited. The required
prepared-only assertion passed and will run again in checkpoint 000003.

Claim-boundary status:
The implementation is private, metadata-only guidance. It makes no vendor
compatibility, hardware certification, production AVL reliability, compliance,
consumer acceptance, hosted-service, SLA, release readiness, production
readiness, agency adoption, public launch, production-grade ETA, or real-world
ETA accuracy claim.

Security/auth status:
No auth behavior, credential storage, token recovery, browser token collection,
bulk token generation, external send, public route, raw payload display, or
admin role boundary changed. The existing POST-only one-time token behavior
remains covered by tests.

Data/migration status:
No migration, durable fleet inventory schema, unknown-device queue, telemetry
contract change, public feed contract change, Trip Updates coupling, or module
dependency change was added.

Master review:
Approved. The implementation uses existing bindings and latest accepted
telemetry to add operator guidance without widening persistence or secret
handling.

Required edits:
Run full phase validation, patch any validation failures caused by Phase 102,
then record checkpoint 000003.

Decision:
Proceed to validation checkpoint 000003.

Next checkpoint:
Phase 102 -- Checkpoint 000003: run validation and patch required gaps.
