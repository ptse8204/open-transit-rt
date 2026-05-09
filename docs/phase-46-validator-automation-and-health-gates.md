# Phase 46 — Validator Automation And Health Gates

Phase 46 adds private validator-health diagnostics for local/reference
operators. It helps an authenticated operator see whether the allowlisted
validators are installed, runnable, stale, blocked, missing artifacts, or
producing results that need review.

This phase does not claim compliance, create evidence, auto-edit GTFS, block
publishing, submit to consumers, or change consumer statuses.

## Operator Surface

- HTML: `/admin/operations/validation-health`
- JSON: `/admin/operations/validation-health.json`
- Roles for GET: `read_only`, `operator`, `editor`, and `admin`
- POST: admin-only `action=run_all`

Both GET routes are private authenticated Operations Console routes and set
`Cache-Control: no-store`. The JSON route emits only the bounded
`ValidationHealthSummary` contract.

The POST route accepts only `action=run_all` and `csrf_token`, caps form bodies
at 64 KiB, requires CSRF when browser cookie auth is in use, and rejects
browser-supplied validator IDs, commands, paths, URLs, argv/args, artifacts,
reports, and timeouts.

## Health Model

The JSON summary includes:

- `generated_at`
- `agency_id`
- `overall_status`
- `tooling_status`
- `feeds`
- `external_evidence_created=false`
- `consumer_statuses_changed=false`
- `compliance_claimed=false`
- `production_readiness_claimed=false`

Feed rows are fixed and ordered:

1. `schedule`
2. `vehicle_positions`
3. `trip_updates`
4. `alerts`

Static schedule health uses only `static-mobilitydata`. Realtime health uses
only `realtime-mobilitydata` for Vehicle Positions, Trip Updates, and Alerts.
Open Transit RT internal GTFS import validation remains context in the GTFS
quality triage page and is not canonical validator health.

The model does not retain or render raw validator reports, stdout, stderr,
argv, raw historical report JSON, token-like fields, cookies, database URLs,
admin URLs, or private paths.

## Script

`scripts/validator-health.sh` and `make validator-health` write private
diagnostics under `.cache/validator-health/<timestamp>` by default:

- `summary.json`
- `summary.md`
- `manifest.json`
- `manifest.md`

Dry-run mode does not require network, a database, Docker, `ADMIN_TOKEN`, or a
running app:

```bash
scripts/validator-health.sh --dry-run
```

Authenticated mode requires `ADMIN_TOKEN` and a safe `ADMIN_BASE_URL`, except
that `PUBLIC_BASE_URL` may be used as the admin base only when it is loopback.
Non-loopback admin URLs must be HTTPS. Admin requests do not follow redirects.

`RUN_VALIDATORS=true` triggers the admin-only `run_all` POST. Without
`ADMIN_TOKEN`, the script does not call private admin routes and instead records
local validator tooling status from `scripts/check-validators.sh`.

`STRICT_VALIDATOR_HEALTH=true` exits non-zero on blocked, failed,
missing-tooling, misconfigured-tooling, artifact-unavailable, or stale health
states.

The script refuses output under `docs/evidence` or evidence-like paths even
when `ALLOW_UNIGNORED_OUTPUT_DIR=true`.

## Deployment Doctor

The reference deployment doctor may GET
`/admin/operations/validation-health.json` only when `ADMIN_TOKEN` and a safe
`ADMIN_BASE_URL` are present. It stores summary fields only and never POSTs the
validator-health route.

## Claim Boundaries

Phase 46 output is private diagnostics only. It is not:

- external evidence;
- final-root evidence;
- a consumer submission artifact;
- consumer acceptance or ingestion proof;
- CAL-ITP/Caltrans compliance proof;
- agency adoption or approval proof;
- hosted SaaS proof;
- production readiness proof;
- vendor compatibility proof;
- production-grade ETA proof.

## Local Diagnostics

Benchmark commands for this phase:

```bash
go test ./internal/compliance -run TestBuildValidationHealthManyReports -bench BenchmarkBuildValidationHealth -benchmem
go test ./cmd/agency-config -run TestValidationHealthPage -bench BenchmarkRenderValidationHealth -benchmem
go test ./cmd/agency-config -run TestValidationHealthJSON -bench BenchmarkRenderValidationHealthJSON -benchmem
go test ./internal/compliance -run TestBuildValidationHealthHostileHistory -bench BenchmarkBuildValidationHealthHostileHistory -benchmem
```

Benchmark results are local engineering diagnostics only. They are not
production capacity, SLA, evidence, compliance, consumer-readiness, or
production-readiness proof.
