# Phase 46 Handoff — Validator Automation And Health Gates

## Status

Phase 46 is complete for the private local/reference validator-health
diagnostic scope.

## What Changed

- Added bounded `ValidationHealthSummary` and fixed four-row feed health model.
- Added `/admin/operations/validation-health` and
  `/admin/operations/validation-health.json`.
- Added admin-only `run_all` validator-health POST with strict form fields,
  CSRF handling, a 64 KiB cap, server-owned validator mappings, and bounded
  partial-success summaries.
- Added `scripts/validator-health.sh` and `make validator-health`.
- Extended the reference deployment doctor to GET validation-health JSON only
  when `ADMIN_TOKEN` and a safe `ADMIN_BASE_URL` are present.
- Linked the page from dashboard, readiness, setup, feeds, checklist, and
  Operations Console navigation.

## Safety Boundary

The feature is private diagnostics only. It creates no evidence packet, writes
nothing under `docs/evidence`, auto-edits no GTFS, does not block publishing,
does not submit to consumers, does not change consumer statuses, and does not
claim CAL-ITP/Caltrans compliance, consumer acceptance, agency adoption,
hosted SaaS availability, production readiness, vendor compatibility, or
production-grade ETA quality.

## Test Coverage

Added coverage for:

- route registration, auth matrix, no-store, content type, unsupported methods;
- strict POST field rejection, body cap, and CSRF rejection;
- fixed feed order across normal, missing-artifact, and stale cases;
- JSON field allowlist and false claim flags;
- HTML/JSON consistency;
- no raw report/stdout/stderr/argv/private path/token leakage;
- partial-success validator run semantics;
- concurrent POST safety;
- validator tooling exit-code mapping;
- 1k/10k report history and hostile history output bounds;
- script help, dry-run, JSON validity, output structure, output-dir refusal,
  and dry-run no-network behavior;
- deployment-doctor GET-only validation-health guard.

## Benchmarks

Record local benchmark output here when closing the phase:

```text
go test ./internal/compliance -run TestBuildValidationHealthManyReports -bench BenchmarkBuildValidationHealth -benchmem
BenchmarkBuildValidationHealth_1k-14       12204   98226 ns/op    218736 B/op   2021 allocs/op
BenchmarkBuildValidationHealth_10k-14       1258  913884 ns/op   2168557 B/op  20052 allocs/op
BenchmarkBuildValidationHealthHostileHistory-14  1245  938710 ns/op  2172158 B/op  20117 allocs/op

go test ./cmd/agency-config -run TestValidationHealthPage -bench BenchmarkRenderValidationHealth -benchmem
BenchmarkRenderValidationHealthPage-14     2484  447941 ns/op     58978 B/op    686 allocs/op
BenchmarkRenderValidationHealthJSON-14     2978  405333 ns/op     36591 B/op    282 allocs/op

go test ./cmd/agency-config -run TestValidationHealthJSON -bench BenchmarkRenderValidationHealthJSON -benchmem
BenchmarkRenderValidationHealthJSON-14     2752  404505 ns/op     36569 B/op    282 allocs/op

go test ./internal/compliance -run TestBuildValidationHealthHostileHistory -bench BenchmarkBuildValidationHealthHostileHistory -benchmem
BenchmarkBuildValidationHealthHostileHistory-14  1170  945607 ns/op  2172771 B/op  20123 allocs/op
```

These results are local engineering diagnostics only. They are not production
capacity, SLA, evidence, compliance, consumer-readiness, or production-readiness
proof.

## Next Recommended Step

Continue the self-hosted agency reuse roadmap with the next operations or
deployment-hardening phase. Do not advance consumer tracker statuses or claim
compliance/readiness without retained claim-specific evidence.
