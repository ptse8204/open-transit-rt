# Phase 84 -- Prediction And ETA Lab Plan

## Scope

Phase 84 adds a private browser-first Prediction and ETA Lab for explaining
why Trip Updates and ETA-like outputs were emitted, withheld, or failed closed.
The Lab is an operator review aid only. It improves transparency around the
existing deterministic predictor, optional external predictor shadow/fail-closed
mode, and local/private realtime-quality backtesting summaries.

The phase may add:

- private `GET /admin/operations/prediction-lab`;
- private `GET /admin/operations/prediction-lab.json`;
- deterministic predictor diagnostics from existing Trip Updates diagnostics;
- withheld Trip Updates reason explanations and next actions;
- external predictor shadow/fail-closed review using bounded stored
  diagnostics only;
- local/private aggregate backtest result browsing for safe `.cache`
  summaries only;
- fixed operator-shell command guidance for existing backtest commands;
- ETA quality caveats and future evidence-gate checklist language.

The phase does not run prediction from the browser, contact external predictor
sidecars, start sidecars, upload observed-arrival files, persist raw
observed-arrival records, write retained evidence, add public routes, add
consumer submission APIs, change consumer statuses, tag, package, publish, or
make production-grade ETA, real-world ETA accuracy, compliance, consumer,
vendor, hardware, SLA, hosted-service, public-launch, or release-readiness
claims.

## Implementation Boundary

Use existing private Operations Console patterns and existing prediction
diagnostics/backtest boundaries:

- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_realtime.go`
- `cmd/agency-config/operations_feed_health.go`
- `cmd/agency-config/operations_navigation.go`
- `internal/prediction`
- `internal/feed/tripupdates`
- `internal/realtimequality`
- `cmd/realtime-quality-backtest`
- `testdata/realtime-quality-backtest`

The default route name is `/admin/operations/prediction-lab` with JSON at
`/admin/operations/prediction-lab.json`. The nav label is `Prediction & ETA
Lab`, placed under the private Realtime group after `Realtime Center`.

No migration is planned. Stop and re-plan before adding persisted Lab state,
uploaded inputs, durable backtest history, or a new mutable command route.

## Security And Data Boundary

The Lab is private and read-only by default. It must use existing
authentication, role, agency-query, GET-only, and no-store patterns.

The Lab must not expose:

- raw telemetry payloads;
- raw observed arrival/departure records;
- raw prediction sample rows;
- raw GTFS-RT payloads;
- raw external predictor request/response bodies;
- raw stdout, stderr, argv, command output, or private file paths;
- external hostnames, sidecar URLs, bearer tokens, cookies, headers, DB URLs,
  token hashes, private keys, or credential values;
- raw `details_json`, `payload_json`, or assignment score details.

The Lab may show fixed local/operator-shell command guidance only as
instructions. Viewing the page executes nothing and contacts nobody.

Backtest browsing must read only safe aggregate result files from the local
`.cache/realtime-quality-backtest/**` shape when present. It must not read raw
input files, write new outputs, follow symlinks outside `.cache`, or treat
backtest summaries as retained evidence.

## Master Approval

The Master Agent approves implementation under these constraints from
sub-agent review:

- Keep Realtime Center as the operational overview and add Prediction & ETA
  Lab as the private diagnostic drilldown.
- Keep deterministic prediction as the default and Vehicle Positions
  independent of predictor availability.
- Use route names that avoid proof language: no `eta-quality`,
  `prediction-accuracy`, or `consumer-ready` route names.
- Keep external predictors optional, disabled/shadow/fail-closed by default,
  and never contacted by a browser request.
- Treat backtests as aggregate local/private diagnostics only.
- Use tables for decision reasons, shadow review, backtest summaries, and
  operator review rows.
- Keep the page usable without JavaScript.
- Keep every claim flag false.
- Use wording such as `private prediction diagnostics`, `withheld-output
  explanation`, `deterministic fallback review`, `shadow comparison summary`,
  `fail-closed diagnostic state`, and `local aggregate diagnostics`.

## Sub-Agent Review Plan

Real or simulated reviews use the intended model levels from the authorized
Phase 75-90 track:

- Context / Repo Truth Sub-Agent, GPT-5.5 x-high: simulated locally because
  the real agent did not return before the plan gate; review confirms Phase 83
  is closed, Phase 84 is next, protected paths are gated, consumer tracker
  must stay prepared-only, and existing prediction/realtime/backtest code is
  available.
- Planning Sub-Agent, GPT-5.5 x-high: Volta approved the
  `/admin/operations/prediction-lab` route, no-migration default, checkpoint
  sequence, tests, and stop conditions.
- Implementation Sub-Agent, GPT-5.5 high: simulated in the main rollout unless
  agent capacity becomes available.
- QA Sub-Agent, GPT-5.5 high: Godel defined private route, JSON allowlist,
  forbidden-claim, no-write, prepared-tracker, backtest-safety, and UI
  regression tests.
- UI/UX Sub-Agent, GPT-5.5 high: Hubble proposed the Realtime nav placement,
  page IA, labels, empty/error states, tables, accessibility, and no-JS
  requirements.
- Documentation / IA Sub-Agent, GPT-5.5 high: Avicenna identified plan,
  tutorial, docs/status, and closeout alignment work.
- Claim-Boundary Sub-Agent, GPT-5.5 high: Fermat approved allowed wording,
  blocked wording, and review checks.
- Security/Auth Sub-Agent, GPT-5.5 high: simulated locally because the agent
  thread limit blocked the real slot; review requires private GET-only routes,
  no browser command execution, no raw/private output exposure, no external
  predictor calls, and no evidence writes.
- Data/Migration Sub-Agent: not planned because Phase 84 should add no
  migration or persisted model. Stop and re-plan if persistence becomes
  necessary.

All required edits from these reviews are incorporated into this plan.

## Checkpoints

```text
Phase 84 -- Checkpoint 000001: add prediction and ETA lab plan
Phase 84 -- Checkpoint 000002: add deterministic predictor diagnostics view
Phase 84 -- Checkpoint 000003: add external predictor shadow review UI
Phase 84 -- Checkpoint 000004: add backtesting result browser
Phase 84 -- Checkpoint 000005: add ETA quality caveats and withheld explanations
Phase 84 -- Checkpoint 000006: close prediction and ETA lab review
```

## Acceptance Criteria

- Lab routes are private, authenticated, agency-scoped, no-store, GET-only,
  and unavailable under `/public`.
- JSON is read-only, bounded, schema-stable, and all claim flags remain false.
- The page has no forms, no POST route, no command execution control, no
  upload field, no credential field, no browser predictor test, no external URL
  test, no sidecar start action, no evidence write, and no consumer status
  mutation.
- Deterministic diagnostics explain adapter status, active feed, candidate
  counts, emitted counts, withheld counts, confidence/staleness/unknown
  assignment metrics, and next actions.
- Withheld reasons are mapped to plain-language meaning, operator next action,
  and what the reason does not prove.
- External predictor review uses stored bounded diagnostics or planned mode
  guidance only. It does not contact sidecars, expose URLs/tokens/hosts, or
  claim named predictor support.
- Backtest browsing uses safe aggregate `.cache/realtime-quality-backtest`
  output only, is bounded, rejects unsafe path shapes, and omits raw input
  rows/private paths.
- Consumer tracker rows remain prepared-only and do not imply submission,
  review, acceptance, listing, display, ingestion, or approval.
- Protected evidence and consumer packet paths remain untouched.
- All seven consumer targets remain exactly `prepared`.

## Validation

Baseline Phase 84 validation:

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
```

Additional Phase 84 checks:

```bash
go test ./cmd/agency-config -run 'PredictionLab|Realtime|OperationsNavigation|RouteTitles'
go test ./cmd/realtime-quality-backtest ./internal/realtimequality ./internal/prediction ./internal/feed/tripupdates
make validate
make test
make realtime-quality
make realtime-quality-backtest
```

Run `RUN_LOCAL_APP=true make release-candidate-check` when route/UI changes are
in place and local app startup is safe. If an environment limitation blocks a
check, record the exact blocker in the Phase 84 handoff without converting it
into a release, compliance, consumer, production, vendor, hardware, SLA, or
ETA-quality claim.
