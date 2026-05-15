# Phase 106 -- Staff Training, Demo Datasets, And Adoption Kit

## Goal

Make Open Transit RT easier to teach and evaluate for small-agency staff by
adding bounded, demo-only training paths, scenario guidance, recovery drills,
trainer scripts, and technical-helper checklists. The work must remain private,
local, synthetic, and non-evidentiary.

## Current Repo Context

- Phase 88 already added private Operations Console Help with role tours,
  first-week checklist, glossary, recovery rows, quick tasks, staff handoff
  checklist, and `docs/operator-training-guide.md`.
- Phase 93 already trialed browser task flows for no-developer evaluators,
  operations staff, technical helpers, release reviewers, and connector
  evaluators.
- Phase 101 through Phase 105 added connector, fleet/device, monitoring,
  deployment, and multi-agency hardening that should now be reflected in
  training language.
- Existing synthetic fixtures live under `testdata/gtfs`,
  `testdata/telemetry-simulator`, `testdata/replay`,
  `testdata/adapter-conformance`, and related fixture directories.
- Existing demo documentation includes `docs/tutorials/agency-demo-flow.md`,
  `docs/tutorials/agency-first-run.md`,
  `docs/tutorials/telemetry-simulator-and-device-trial.md`,
  `docs/tutorials/external-adapter-conformance.md`, and the tutorials index.

## Scope

- Add or update demo-only staff training guidance.
- Add a bounded synthetic scenario catalog that references committed fixtures
  without inventing real agency or vendor data.
- Add role-based training paths for no-developer evaluators, directors,
  daily operators, technical helpers, integrators, trainers, and maintainers
  when useful.
- Add a first-week operator guide, common mistakes and recovery drills, trainer
  script, and technical-helper checklist.
- Optionally extend the existing private `/admin/operations/help` and
  `/admin/operations/help.json` read-only model with adoption-kit/training
  sections if that materially improves in-app discoverability.
- Add or update focused tests for any changed private Help model or UI.

## Boundaries

- Do not create retained evidence.
- Do not write protected evidence paths.
- Do not contact agencies, vendors, consumers, aggregators, portals, or other
  external services.
- Do not use real agency, vendor, device, credential, or private payload data.
- Do not move consumer statuses beyond `prepared`.
- Do not tag, publish, create releases, push images, or distribute packages.
- Do not claim adoption, agency approval, consumer acceptance, compliance,
  final-root readiness, hosted SaaS, SLA/uptime, production readiness,
  vendor compatibility, hardware certification, public launch, or
  production-grade ETA quality.
- Keep browser help read-only, authenticated, `no-store`, and free of forms,
  mutation routes, backend command execution, `.cache` reads, or public admin
  routes.

## Implementation Plan

1. Extend the operator training guide and tutorial navigation with a demo-only
   scenario catalog and role-based training paths.
2. Add an adoption-kit/trainer-script section that uses local/synthetic
   fixtures only and states what each exercise does not prove.
3. Add a technical-helper checklist for safe handoff inputs: page, blocker,
   owner, intended next action, fixture or docs reference, and whether separate
   authorization is needed.
4. If in-app discoverability is worth the small code change, add a read-only
   private Help section for demo scenarios and trainer checklist data with
   all-false claim flags.
5. Update focused tests to enforce safe links, required training rows, no
   public routes, no forms, no command execution, and no unsafe claims.
6. Run focused and baseline validation, then close with handoff/status docs.

## Non-Goals

- No evidence collection or adoption proof.
- No real pilot, real agency trial, vendor/device proof, public launch, or
  external submission workflow.
- No new CLI command unless a tiny documentation-only lint helper becomes
  necessary; default is docs plus private in-app read-only model.
- No new migration, persistent training table, user progress tracking, hosted
  analytics, notification send, or credential capture.
- No changes to public feed contracts or Trip Updates prediction boundaries.

## Checkpoint Plan

- `Phase 106 -- Checkpoint 000001: add staff training, demo datasets, and adoption kit plan`
- `Phase 106 -- Checkpoint 000002: implement primary scoped work`
- `Phase 106 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 106 -- Checkpoint 000004: close staff training, demo datasets, and adoption kit review`

## Focused Validation Targets

- `go test ./cmd/agency-config -run 'OperationsHelp|OperationsSharedLayoutRendersContextualHelp|Training|Demo'`
- `git diff --check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`

Because this phase is expected to change docs and may change private Help UI
code/tests, closeout also requires:

- `git status --short`
- `make check`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact prepared-only consumer tracker assertion
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`

## Checkpoint Report -- 000001

Checkpoint:
Phase 106 -- Checkpoint 000001: add staff training, demo datasets, and adoption
kit plan.

Sub-agents used or simulated, including intended model level:
Real Planning Sub-Agent -- GPT-5.5 x-high. Context / Repo Truth Sub-Agent --
GPT-5.5 x-high was attempted but timed out and was shut down without edits;
Context / Repo Truth was therefore simulated by the Master Agent using direct
repository inspection. Implementation, QA, UI/UX, Documentation / IA,
Claim-Boundary, Security/Auth, Data/Migration, and Release/Supply-Chain roles
are simulated by the Master Agent for this plan checkpoint. Master Agent --
GPT-5.5 x-high, current thread.

Changed files:
`docs/phase-106-staff-training-demo-datasets-and-adoption-kit.md`.

Validation run:
Initial repository inspection reviewed existing Phase 106 roadmap prompts,
Phase 88 Help/training implementation, operator training guide, tutorials,
synthetic fixture directories, and current source-of-truth status docs. After
adding the plan, `git status --short` showed only the Phase 106 plan doc;
`git diff --check` passed; `python3 -m json.tool
docs/evidence/consumer-submissions/status.json >/dev/null` passed; the exact
prepared-only consumer tracker assertion passed; and `git status --short --
docs/evidence/consumer-submissions docs/evidence/captured db/migrations
go.mod go.sum` returned no output.

Blocked checks:
Implementation, focused tests, and closeout baseline checks are not yet run
because this checkpoint only approves the Phase 106 plan. Release-candidate
checks, package generation/audit, evidence collection, publication, external
contact, real credentials, and consumer actions are out of scope.

Protected path status:
No protected evidence path is part of the plan. The plan forbids protected path
writes.

Consumer tracker status:
The consumer tracker is not part of the plan. The seven targets must remain in
order and `prepared`.

Claim-boundary status:
The plan explicitly forbids adoption, agency approval, consumer acceptance,
compliance, final-root readiness, hosted SaaS, SLA/uptime, production
readiness, vendor compatibility, hardware certification, public launch, and
production-grade ETA claims.

Security/auth status:
The plan preserves existing private authenticated Help behavior and forbids
public admin routes, forms, mutation routes, command execution, `.cache` reads,
credential capture, and raw private data.

Data/migration status:
No migration, persistent training state, user progress tracking, hosted
analytics, notification send, public feed contract change, or module dependency
change is planned.

Master review:
Approved. The smallest safe Phase 106 implementation is docs-first training
and demo scenario guidance, with a narrowly scoped private Help model extension
only if it improves discoverability.

Required edits:
Implement the bounded training/adoption-kit scope, update tests if private Help
code changes, and record validation results.

Decision:
Proceed to implementation checkpoint 000002 after plan validation and commit.

Next checkpoint:
Phase 106 -- Checkpoint 000002: implement primary scoped work.

## Checkpoint Report -- 000002

Checkpoint:
Phase 106 -- Checkpoint 000002: implement primary scoped work.

Sub-agents used or simulated, including intended model level:
Implementation, QA, UI/UX, Documentation / IA, Claim-Boundary, Security/Auth,
and Data/Migration roles are simulated by the Master Agent using the approved
plan and the real Planning Sub-Agent output. Master Agent -- GPT-5.5 x-high,
current thread.

Changed files:
`cmd/agency-config/operations_help.go`; `cmd/agency-config/operations.go`;
`cmd/agency-config/main_test.go`; `docs/operator-training-guide.md`;
`docs/tutorials/README.md`; `docs/tutorials/staff-training-demo-kit.md`;
`testdata/training-demo/README.md`; `testdata/training-demo/scenarios.json`;
`docs/phase-106-staff-training-demo-datasets-and-adoption-kit.md`.

Validation run:
`gofmt` ran on changed Go files. `go test ./cmd/agency-config -run
'OperationsHelp|OperationsSharedLayoutRendersContextualHelp|Training|Demo'`
passed. `python3 -m json.tool testdata/training-demo/scenarios.json
>/dev/null` passed. `git diff --check` passed. `make
audit-product-acceptance` passed. `make audit-final-claim-review` passed.

Blocked checks:
Full closeout baseline, `make validate`, `make test`, Docker Compose config,
and final protected-path checks are deferred to checkpoint 000003/000004.
Release-candidate checks, package generation/audit, evidence collection,
publication, external contact, real credentials, and consumer actions are out
of scope.

Protected path status:
No protected evidence path was modified or required.

Consumer tracker status:
The consumer tracker was not modified. The seven targets remain required to
stay in order and `prepared`.

Claim-boundary status:
The implementation adds training/demo guidance only and repeats no-evidence,
no-adoption, no-compliance, no-consumer, no-public-launch, no-hosted-service,
no-SLA, no-vendor, no-hardware, and no-ETA-quality boundaries.

Security/auth status:
The private Help additions remain authenticated, read-only, `no-store`, and
free of forms, mutation routes, backend command execution, `.cache` reads,
credential collection, public admin routes, and raw private data.

Data/migration status:
No migration, persistent training state, user progress tracking, hosted
analytics, notification send, public feed contract change, or module dependency
change occurred. The new scenario catalog references only committed synthetic
fixtures under `testdata/`.

Master review:
Approved. The implementation keeps Phase 106 inside training/adoption-kit
boundaries while improving in-app discoverability and docs.

Required edits:
Run the full focused validation and baseline checks. Patch only failures
caused by this phase.

Decision:
Proceed to validation/audit checkpoint 000003.

Next checkpoint:
Phase 106 -- Checkpoint 000003: run validation and patch required gaps.

## Checkpoint Report -- 000003

Checkpoint:
Phase 106 -- Checkpoint 000003: run validation and patch required gaps.

Sub-agents used or simulated, including intended model level:
QA, Documentation / IA, Claim-Boundary, Security/Auth, Data/Migration, and
Release/Supply-Chain roles are simulated by the Master Agent for validation.
Master Agent -- GPT-5.5 x-high, current thread.

Changed files:
`docs/phase-106-staff-training-demo-datasets-and-adoption-kit.md`.

Validation run:
Focused checks passed: `go test ./cmd/agency-config -run
'OperationsHelp|OperationsSharedLayoutRendersContextualHelp|Training|Demo'`;
`python3 -m json.tool testdata/training-demo/scenarios.json >/dev/null`; and
`git diff --check`.

Baseline/code-change checks passed: `git status --short`; `python3 -m
json.tool docs/evidence/consumer-submissions/status.json >/dev/null`; the exact
prepared-only consumer tracker assertion; `git status --short --
docs/evidence/consumer-submissions docs/evidence/captured db/migrations
go.mod go.sum`; `make check`; `make audit-product-acceptance`; `make
audit-final-claim-review`; `make validate`; `make test`; `docker compose -f
deploy/docker-compose.yml config`; final `git status --short`; final protected
path status check; and final `git diff --check`.

Blocked checks:
None for this phase. Release-candidate checks, package generation/audit,
retained evidence, real agency/vendor/device data, real credentials, external
contact, consumer submission, public publication, and tag/release/package/image
publication remain out of scope.

Protected path status:
`git status --short -- docs/evidence/consumer-submissions
docs/evidence/captured db/migrations go.mod go.sum` returned no output. No
protected evidence path was modified.

Consumer tracker status:
The JSON syntax check and exact prepared-only assertion passed. All seven
targets remain in order and `prepared`.

Claim-boundary status:
`make audit-product-acceptance` and `make audit-final-claim-review` passed.
No adoption, agency approval, consumer acceptance, compliance, final-root
readiness, hosted SaaS, SLA/uptime, production readiness, vendor
compatibility, hardware certification, public launch, or production-grade ETA
claim was added.

Security/auth status:
Validation covered the private authenticated Help model and rendering. The new
training sections are read-only, no-store, local/synthetic, and do not add
forms, mutation routes, command execution, `.cache` reads, credential capture,
public admin routes, or raw private data exposure.

Data/migration status:
No migration, persistent training state, user progress tracking, hosted
analytics, notification send, public feed contract change, `go.mod`, or
`go.sum` change occurred.

Master review:
Approved. No required validation gaps remain for the Phase 106 scope.

Required edits:
Prepare the Phase 106 handoff and closeout docs.

Decision:
Proceed to closeout checkpoint 000004.

Next checkpoint:
Phase 106 -- Checkpoint 000004: close staff training, demo datasets, and
adoption kit review.
