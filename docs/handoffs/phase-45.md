# Phase 45 Handoff — GTFS Quality Triage Loop

## Summary

Phase 45 added a private `/admin/operations/gtfs-quality` Operations Console
page that turns stored static GTFS validator and internal import validation
notices into bounded operator triage actions.

The page separates Canonical MobilityData static validator results from Open
Transit RT internal import validation. It uses neutral statuses only:
`blocking`, `needs_review`, `informational`, and `unknown`.

## Boundary

- Validator output is diagnostics/supporting signal only.
- The route does not claim consumer acceptance or CAL-ITP/Caltrans compliance.
- The route does not create, export, download, retain, or write evidence
  packets and does not write to `docs/evidence`.
- The route does not auto-edit agency GTFS.
- Benchmarks below are local engineering diagnostics only. They are not
  production capacity, SLA, evidence, compliance, or consumer-readiness proof.

## Implementation Notes

- Added `compliance.CanonicalStaticValidatorName` and
  `compliance.InternalGTFSImportValidatorName`.
- Added bounded derived GTFS quality triage structs and builder logic under
  `internal/compliance`.
- Added latest validation report lookup for the Operations Console.
- Added GET/POST route handling, strict POST form checking, POST body cap,
  no-store headers, source-separated UI sections, stale active-feed detection,
  and admin-only rerun action.
- Successful rerun stores only the normal validation result row through the
  existing validation flow, against the active schedule feed version.
- Stored raw report fields can remain in `validation_report.report_json`; the
  derived triage model and HTML omit raw report, stdout, stderr, argv, commands,
  temp paths, and private artifacts.

## Local Diagnostics

Benchmark commands were run as local diagnostics only:

```sh
go test ./internal/compliance -run TestBuildGTFSQualityTriageLargeCanonicalReport -bench BenchmarkBuildGTFSQualityTriage -benchmem
go test ./internal/compliance -run TestGTFSQualityTriageHostileReport -bench BenchmarkBuildGTFSQualityTriageHostileReport -benchmem
go test ./cmd/agency-config -run TestGTFSQualityPageLargeReportRender -bench BenchmarkRenderGTFSQualityPage -benchmem
```

Results:

- `BenchmarkBuildGTFSQualityTriage/10000-14`: 62,764,616 ns/op;
  5,368,463 B/op; 61,124 allocs/op.
- `BenchmarkBuildGTFSQualityTriage/50000-14`: 298,176,010 ns/op;
  26,567,902 B/op; 301,125 allocs/op.
- `BenchmarkBuildGTFSQualityTriageHostileReport-14`: 428,385,375 ns/op;
  411,641,653 B/op; 448,161 allocs/op.
- `BenchmarkRenderGTFSQualityPage-14`: 334,188,694 ns/op;
  387,201,437 B/op; 353,482 allocs/op.

The first benchmark command also matched
`BenchmarkBuildGTFSQualityTriageHostileReport` because its name shares the
same prefix; the exact hostile-report command was still run separately and the
separate value above is the one recorded for that benchmark.

## Required Checks

Run before closing:

```sh
make validate
make test
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
python3 - <<'PY'
import json
from pathlib import Path
expected = ["Google Maps", "Apple Maps", "Transit App", "Bing Maps", "Moovit", "Mobility Database", "transit.land"]
data = json.loads(Path("docs/evidence/consumer-submissions/status.json").read_text())
records = data.get("targets", data if isinstance(data, list) else [])
seen = {record["target"]: record.get("status") for record in records}
assert list(seen) == expected, seen
assert all(seen[name] == "prepared" for name in expected), seen
PY
git diff --exit-code -- docs/evidence/consumer-submissions/status.json
docker compose -f deploy/docker-compose.yml config
```

The consumer tracker must remain byte-for-byte unchanged.
