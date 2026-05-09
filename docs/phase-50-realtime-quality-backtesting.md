# Phase 50 -- Realtime Quality Backtesting

## Status

Planning approved. Implementation has not started.

## Goal

Phase 50 adds a private CLI/library workflow for comparing synthetic or
operator-supplied observed stop arrivals/departures against prediction samples.
It produces route/time-period ETA quality diagnostics and maturity gates as
local engineering output only.

Phase 50 must not persist backtest history in the database, add Operations
Console storage, create evidence, block publishing, or claim production-grade
ETA quality.

## Scope

- Add a backtesting library under `internal/realtimequality`.
- Add `cmd/realtime-quality-backtest`.
- Add synthetic public-safe fixtures under `testdata/realtime-quality-backtest`.
- Add `make realtime-quality-backtest`.
- Produce bounded route/time-period metrics:
  - MAE;
  - median absolute error;
  - p90 absolute error;
  - prediction coverage;
  - future stop coverage;
  - stale, missing, and withheld counts;
  - diagnostic maturity gates.
- Preserve existing `make realtime-quality` replay regression coverage.
- Preserve existing Operations Console Trip Updates diagnostics without adding
  a new stored backtest view.

## Non-Goals

- No DB persistence, migration, aggregate table, or stored backtest history.
- No Operations Console latest aggregate view.
- No public unauthenticated API.
- No publish blocking.
- No GTFS auto-editing.
- No evidence packet creation or writes under `docs/evidence`.
- No consumer contact or consumer tracker status changes.
- No named predictor runtime expansion beyond Phase 49.
- No final-root, compliance, agency adoption, consumer acceptance, hosted SaaS,
  vendor compatibility, production-readiness, real-world ETA accuracy, or
  production-grade ETA claim.

## Expected Outputs

The default output directory is:

```text
.cache/realtime-quality-backtest/<timestamp>
```

The CLI must write exactly these files:

- `summary.json`
- `summary.md`
- `metrics.json`
- `metrics.md`
- `manifest.json`

Output path validation must reject:

- `docs/evidence` and any child path;
- evidence-like paths;
- symlink ancestors;
- unsafe path traversal.

Outputs must be bounded in row count, group count, and size. They are private
local diagnostics, not evidence packets.

## Input Boundaries

The CLI may read local operator files explicitly supplied by path, but it must
not copy raw operator input files into output directories and must not persist
raw rows.

Committed fixtures must remain synthetic and public-safe.

Outputs must not include:

- raw observed rows;
- raw prediction rows;
- raw telemetry;
- raw GTFS-RT payloads;
- private device IDs;
- driver IDs;
- vendor payloads;
- credentials;
- authorization or cookie headers;
- DB URLs;
- private file paths;
- raw logs.

Outputs may include bounded aggregate diagnostics and redacted manifest
metadata such as input kind, schema version, record counts, and safe checksums.

## Local Schemas

Define versioned local schemas for:

- observed stop events;
- prediction samples.

Observed stop events should include only bounded matching keys and times:

- agency/feed scope;
- route ID;
- trip ID;
- start date;
- start time;
- stop ID;
- stop sequence;
- event type;
- observed time;
- optional public-safe service pattern label.

Prediction samples should include only bounded prediction keys and values:

- generated time;
- adapter name;
- agency/feed scope;
- trip instance;
- stop sequence;
- predicted arrival/departure;
- confidence;
- schedule relationship;
- optional withheld reason.

## Metrics

Join predictions to observations by agency/feed scope, trip instance, stop
sequence, and event type.

Group metrics by:

- overall;
- route;
- agency-local time period;
- route plus time period.

Compute:

- lead time;
- absolute error seconds;
- MAE;
- median absolute error;
- p90 absolute error;
- coverage denominator and numerator;
- future stop coverage;
- missing prediction count;
- missing observation count;
- stale prediction count;
- withheld-by-reason counts.

Maturity gates must use diagnostic labels only:

- `insufficient_data`
- `diagnostic_pass`
- `diagnostic_watch`

Forbidden gate labels:

- `ready`
- `production_ready`
- `compliant`
- `accepted`

External/shadow predictor comparison may be summarized only as diagnostic
deltas when local samples are provided. It must not be treated as proof of
better ETA quality.

## Safety Boundaries

- Vehicle Positions remain the first production-directed realtime output.
- Trip Updates remain behind `internal/prediction.Adapter`.
- Phase 50 must not alter public feed contracts.
- Phase 50 must not introduce public unauthenticated routes.
- Phase 50 must not contact consumers or external portals.
- Phase 50 must not write under `docs/evidence`.
- Phase 50 benchmark/backtest output is local engineering diagnostics only, not
  SLA, capacity, readiness, compliance, consumer acceptance, or ETA-quality
  proof.

## Tests

- Schema validation and malformed input tests.
- Synthetic fixture tests for matched, missing, stale, withheld,
  after-midnight, frequency, block-transition, manual override, and
  zero-denominator cases.
- Metric tests for MAE, median, p90, coverage, future stop coverage,
  route/time-period grouping, and diagnostic gate labels.
- CLI tests for exact output files, `.cache` default output, path rejection,
  symlink rejection, no raw input copying, redaction, bounded output,
  deterministic JSON ordering, and forbidden claim label absence.
- Preserve existing `make realtime-quality` replay coverage.

## Performance And Scale

- Use indexed maps for joins, not pairwise prediction-by-observation scans.
- Add a synthetic scale test or benchmark for large local inputs.
- Document benchmark output as local engineering diagnostics only.

## Docs And Handoff Updates

Implementation should update:

- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-50.md`
- `docs/backlog.md`
- `docs/open-questions.md`
- `docs/decisions.md` only if a meaningful architecture decision is added

Do not update `docs/evidence`.

## Required Verification

Phase 50 must run and report:

```bash
go test ./internal/realtimequality ./cmd/realtime-quality-backtest
make realtime-quality
make realtime-quality-backtest
make validate
make test
make smoke
make test-integration
git diff --check
docker compose -f deploy/docker-compose.yml config
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
git diff --exit-code -- docs/evidence/consumer-submissions/status.json
```

Additional broader `docs/evidence` no-diff safety scans may be used, but they
must not replace the exact `status.json` consumer tracker check above.

## Closeout Requirements

Phase 50 is not closed until implementation is reviewed by the master agent,
required checks pass or blockers are documented truthfully, the Phase 50
handoff exists, `docs/handoffs/latest.md` is updated, roadmap/status docs are
consistent, no forbidden claims are introduced, no `docs/evidence` files are
edited, and the consumer tracker remains exactly seven `prepared` targets.
