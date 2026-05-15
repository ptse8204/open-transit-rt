# Phase 99 -- Prediction / ETA Conformance And Backtesting V2

## Scope

Phase 99 improves private prediction conformance and synthetic/local realtime
quality backtesting diagnostics. The phase is an extension of the existing
private Prediction and ETA Lab, `internal/prediction`, and
`internal/realtimequality` boundaries.

Allowed work:

- expand synthetic realtime-quality backtest fixtures for edge cases;
- add aggregate-only metrics or private view rows for after-midnight,
  frequency/headway, service-calendar, withheld, blocked, unknown, ambiguous,
  and shadow/fail-closed cases;
- improve private Prediction Lab conformance reporting and operator next
  actions;
- add or update tests for fixture shape, aggregate metrics, private route
  rendering, and claim boundaries;
- update documentation and tutorials for local/synthetic diagnostics.

Not allowed:

- real-world ETA accuracy claims;
- production-grade ETA quality claims;
- external predictor network contact;
- browser-triggered predictor execution;
- retained evidence collection;
- protected evidence path writes;
- consumer tracker status movement;
- migrations or persistent raw observed-event storage;
- release, tag, package publication, hosted-service, vendor, hardware, SLA,
  public-launch, compliance, consumer-acceptance, or release-readiness claims.

## Existing Surfaces

Phase 84 added private read-only Prediction and ETA Lab routes:

- `/admin/operations/prediction-lab`;
- `/admin/operations/prediction-lab.json`.

Existing realtime-quality backtesting lives in:

- `cmd/realtime-quality-backtest`;
- `internal/realtimequality`;
- `testdata/realtime-quality-backtest`.

The current backtest command writes aggregate private diagnostics only under
`.cache/realtime-quality-backtest/**` and rejects evidence-like output paths.
Phase 99 must preserve that model.

## Master-Approved Plan

1. Add this plan and record the Phase 99 claim/security/data boundaries.
2. Implement the smallest code/test/docs change that makes conformance and
   backtest-v2 coverage clearer to operators.
3. Run focused prediction/realtime-quality tests and the required baseline
   validation; patch any changed-code failures.
4. Close Phase 99 with a handoff, status-doc updates, protected-path status,
   prepared-only consumer tracker confirmation, and exact blockers.

The Master Agent approves implementation only if it remains private,
synthetic/local, aggregate-only, no-send, no-migration, and no-claim.

## Sub-Agent Plan

| Role | Intended model | Use in Phase 99 |
| --- | --- | --- |
| Context / Repo Truth Sub-Agent | GPT-5.5 x-high | Read-only inspection of Prediction Lab, realtime-quality backtest tooling, prediction adapter seams, tests, fixtures, and docs. |
| Planning Sub-Agent | GPT-5.5 x-high | Read-only checkpoint plan, validation plan, and guardrail review. |
| Implementation Sub-Agent | GPT-5.5 high | Simulated by Master unless a bounded disjoint edit becomes useful. |
| QA Sub-Agent | GPT-5.5 high | Simulated by Master through focused tests and full required validation. |
| UI/UX Sub-Agent | GPT-5.5 high | Simulated by Master for private Prediction Lab wording and table shape. |
| Documentation / IA Sub-Agent | GPT-5.5 high | Simulated by Master for phase docs, tutorials, roadmap status, and handoff. |
| Claim-Boundary Sub-Agent | GPT-5.5 high | Simulated by Master with final claim audit and forbidden wording review. |
| Security/Auth Sub-Agent | GPT-5.5 high | Simulated by Master; preserve private GET-only route behavior and no browser command execution. |
| Data/Migration Sub-Agent | GPT-5.5 high | Simulated by Master because no persistence or migration is planned. Stop before adding persistence. |

## Checkpoints

```text
Phase 99 -- Checkpoint 000001: add prediction / eta conformance and backtesting v2 plan
Phase 99 -- Checkpoint 000002: implement primary scoped work
Phase 99 -- Checkpoint 000003: run validation and patch required gaps
Phase 99 -- Checkpoint 000004: close prediction / eta conformance and backtesting v2 review
```

## Validation Plan

Focused checks:

```bash
go test ./cmd/agency-config -run 'PredictionLab|Realtime|OperationsNavigation|RouteTitles'
go test ./cmd/realtime-quality-backtest ./internal/realtimequality ./internal/prediction
make realtime-quality
make realtime-quality-backtest
```

Phase closeout baseline:

```bash
git status --short
git diff --check
make check
make audit-product-acceptance
make audit-final-claim-review
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
python3 - <<'PY'
import json
from pathlib import Path

expected = [
    "Google Maps",
    "Apple Maps",
    "Transit App",
    "Bing Maps",
    "Moovit",
    "Mobility Database",
    "transit.land",
]

data = json.loads(Path("docs/evidence/consumer-submissions/status.json").read_text())
records = data.get("targets", [])
seen = {row["target"]: row.get("status") for row in records}
assert list(seen) == expected, seen
assert all(seen[name] == "prepared" for name in expected), seen
PY
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum
make validate
make test
docker compose -f deploy/docker-compose.yml config
```

Connector and release-candidate checks are not required for Phase 99 unless the
implementation unexpectedly changes connector or release surfaces.

## Implementation Summary

Checkpoint 000002 added an aggregate synthetic conformance report to realtime
quality backtest summaries. The report derives only from local synthetic
fixture rows and records:

- after-midnight service coverage;
- frequency/headway service coverage;
- service-calendar start-date/start-time instance coverage;
- unknown, ambiguous, and fail-closed withheld-output coverage;
- shadow/fail-closed external predictor evaluation coverage.

The private Prediction Lab backtest table now surfaces only a bounded
conformance signal such as `synthetic_covered (5/5 synthetic cases)`. The
browser still reads exact aggregate `.cache/realtime-quality-backtest/**`
outputs only and does not expose raw fixture rows, execute commands, contact
predictors, mutate feeds, write evidence, or make ETA-quality claims.

The committed synthetic fixture set now includes unknown-assignment,
ambiguous-assignment, and external-predictor fail-closed withheld rows in
addition to the existing after-midnight, frequency, stale, missing, block, and
shadow samples.

## Checkpoint 000001 Report

Checkpoint: Phase 99 -- Checkpoint 000001: add prediction / eta conformance and backtesting v2 plan

Sub-agents used or simulated, including intended model level: Context / Repo Truth Sub-Agent GPT-5.5 x-high and Planning Sub-Agent GPT-5.5 x-high were launched read-only; Implementation, QA, UI/UX, Documentation / IA, Claim-Boundary, Security/Auth, and Data/Migration are simulated by the Master Agent for this planning checkpoint.

Changed files: `docs/phase-99-prediction-eta-conformance-and-backtesting-v2.md`

Validation run: `git status --short` before edits; source/doc read of Phase 99 prompt, Prediction Lab, realtime-quality backtest tooling, current status, latest handoff, roadmap status, and post-90 validation/operating manuals.

Blocked checks: Full validation deferred until implementation and closeout checkpoints.

Protected path status: No protected evidence path edits planned or made.

Consumer tracker status: Must remain exactly seven prepared-only targets; no status edits planned or made.

Claim-boundary status: Plan explicitly forbids real-world ETA accuracy, production-grade ETA quality, compliance, release-readiness, consumer acceptance, vendor, hardware, hosted-service, SLA, public-launch, and adoption claims.

Security/auth status: Plan preserves private read-only authenticated Operations Console behavior; no browser predictor execution, external contact, credentials, or raw private rows.

Data/migration status: No migration, raw observed-event persistence, or data-model change planned.

Master review: Approved to proceed under private, synthetic/local, aggregate-only, no-send, no-migration, no-claim constraints.

Required edits: Implement a bounded conformance/backtest-v2 improvement with tests and docs.

Decision: Continue to Checkpoint 000002.

Next checkpoint: Phase 99 -- Checkpoint 000002: implement primary scoped work

## Checkpoint 000002 Report

Checkpoint: Phase 99 -- Checkpoint 000002: implement primary scoped work

Sub-agents used or simulated, including intended model level: Context / Repo Truth Sub-Agent GPT-5.5 x-high and Planning Sub-Agent GPT-5.5 x-high completed read-only review. Implementation, QA, UI/UX, Documentation / IA, Claim-Boundary, Security/Auth, and Data/Migration were simulated by the Master Agent for this bounded implementation checkpoint.

Changed files: `internal/realtimequality/backtest.go`, `internal/realtimequality/browser.go`, `internal/realtimequality/backtest_test.go`, `internal/realtimequality/browser_test.go`, `cmd/agency-config/operations.go`, `cmd/agency-config/main_test.go`, `testdata/realtime-quality-backtest/observed-events.json`, `testdata/realtime-quality-backtest/prediction-samples.json`, `testdata/realtime-quality-backtest/README.md`, `docs/tutorials/prediction-eta-lab.md`, `docs/phase-99-prediction-eta-conformance-and-backtesting-v2.md`

Validation run: `python3 -m json.tool testdata/realtime-quality-backtest/observed-events.json >/dev/null`; `python3 -m json.tool testdata/realtime-quality-backtest/prediction-samples.json >/dev/null`; `gofmt -w internal/realtimequality/backtest.go internal/realtimequality/browser.go internal/realtimequality/backtest_test.go internal/realtimequality/browser_test.go cmd/agency-config/main_test.go`; `git diff --check`; `go test ./internal/realtimequality ./cmd/realtime-quality-backtest ./internal/prediction`; `go test ./cmd/agency-config -run 'PredictionLab|Realtime|OperationsNavigation|RouteTitles'`

Blocked checks: Full phase closeout validation deferred to Checkpoint 000003.

Protected path status: No protected evidence path edits made.

Consumer tracker status: No consumer tracker edits made; prepared-only status must be rechecked in Checkpoint 000003.

Claim-boundary status: Added wording and tests keep synthetic conformance local/private and explicitly non-evidentiary; no ETA quality, compliance, consumer, vendor, hardware, hosted-service, SLA, release, or public-launch claim added.

Security/auth status: Private browser surface remains read-only; no browser command execution, external predictor contact, credentials, raw private rows, raw telemetry, raw predictions, or private paths added.

Data/migration status: No migration, persistence model, or raw observed-event storage added.

Master review: Approved. The implementation follows the smallest safe seam: aggregate synthetic conformance metadata and private display only.

Required edits: Run full required validation and patch any changed-code failures.

Decision: Continue to Checkpoint 000003.

Next checkpoint: Phase 99 -- Checkpoint 000003: run validation and patch required gaps

## Checkpoint 000003 Report

Checkpoint: Phase 99 -- Checkpoint 000003: run validation and patch required gaps

Sub-agents used or simulated, including intended model level: Context / Repo Truth Sub-Agent GPT-5.5 x-high and Planning Sub-Agent GPT-5.5 x-high completed read-only review. QA, Claim-Boundary, Security/Auth, Data/Migration, UI/UX, Documentation / IA, and Implementation were simulated by the Master Agent for validation and audit.

Changed files: `docs/phase-99-prediction-eta-conformance-and-backtesting-v2.md`

Validation run: `git status --short`; `git diff --check`; `go test ./internal/realtimequality ./cmd/realtime-quality-backtest ./internal/prediction`; `go test ./cmd/agency-config -run 'PredictionLab|Realtime|OperationsNavigation|RouteTitles'`; `make realtime-quality`; `make realtime-quality-backtest`; `make check`; `make audit-product-acceptance`; `make audit-final-claim-review`; `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`; exact prepared-only consumer tracker assertion; `make validate`; `make test`; `docker compose -f deploy/docker-compose.yml config`; `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`; final `git status --short`; final `git diff --check`

Blocked checks: None for Phase 99.

Protected path status: Protected evidence and consumer-submission paths remained untouched; protected-path status check returned no output.

Consumer tracker status: `docs/evidence/consumer-submissions/status.json` parsed successfully and remained exactly seven targets, all `prepared`.

Claim-boundary status: Product acceptance and final claim audits passed; no forbidden ETA quality, compliance, consumer, vendor, hardware, hosted-service, SLA, public-launch, release, or adoption claim detected.

Security/auth status: No browser command execution, external predictor contact, credentials, raw private rows, raw telemetry, raw predictions, public API, or feed mutation added.

Data/migration status: No migration or persistence changes; `db/migrations`, `go.mod`, and `go.sum` status checks returned no output.

Master review: Approved. Validation passed with no required code patches.

Required edits: Add Phase 99 closeout handoff/status updates.

Decision: Continue to Checkpoint 000004.

Next checkpoint: Phase 99 -- Checkpoint 000004: close prediction / eta conformance and backtesting v2 review
