# Phase 51 -- Operations Reliability And SLO Readiness

## Status

Closed for the approved private operations reliability diagnostics scope.

## Goal

Add private operations reliability diagnostics for feed freshness, incidents,
backup/restore readiness, alerting readiness, availability sampling, and
long-running operations workflows.

Phase 51 must not create evidence or proof, and must not claim SLA coverage,
uptime guarantees, production readiness, compliance, consumer acceptance,
agency adoption, hosted SaaS availability, vendor compatibility, or
production-grade ETA quality.

## Scope

Phase 51 is a mixed runtime/code, script, docs, and tests phase.

- Add GET-only `/admin/operations/reliability`.
- Add GET-only `/admin/operations/reliability.json`.
- Use existing admin Operations Console auth roles: read-only, operator,
  editor, and admin.
- Enforce agency query scoping consistently with existing admin operations
  routes.
- Add no public route.
- Add no migrations.
- Use only existing `feed_health_snapshot` and `incident` database state.
- Add Vehicle Positions health snapshot persistence using the existing
  `feed_health_snapshot` schema only.
- Keep Vehicle Positions persistence best-effort, non-blocking,
  bounded/redacted, and unable to change public feed response status.
- Keep schedule health persistence out of scope; missing schedule health must
  report `unknown` or `missing`.
- Add `scripts/operations-reliability.sh`.
- Add `make operations-reliability`.

## Non-Goals

- No new monitoring stack.
- No public unauthenticated API.
- No evidence packet creation.
- No writes under `docs/evidence`.
- No consumer contact.
- No consumer tracker status change.
- No publish blocking.
- No GTFS auto-editing.
- No migrations.
- No raw backup/log inspection.
- No raw payload capture.
- No webhook sending.
- No SLA, uptime guarantee, production-readiness, compliance, hosted SaaS,
  agency adoption, vendor compatibility, consumer acceptance, or
  production-grade ETA claim.

## Safety Boundaries

All outputs are private diagnostics. `scripts/operations-reliability.sh`
defaults to:

```text
.cache/operations-reliability/<UTC timestamp>/
```

The script must write exactly:

- `summary.json`
- `summary.md`
- `manifest.json`
- `manifest.md`
- `reliability-review.txt`

The script must reject `docs/evidence`, evidence-like output paths, symlinked
sources, oversized inputs, raw logs, raw backup dumps, DB URLs, webhook values,
secrets, and private payloads. It must not send notifications or call mutating
admin endpoints.

## Evidence And Claim Boundaries

Every runtime and script summary must include these exact flags, all `false`:

- `external_evidence_created`
- `final_root_evidence_created`
- `consumer_statuses_changed`
- `compliance_claimed`
- `production_readiness_claimed`
- `sla_claimed`
- `uptime_guarantee_claimed`
- `hosted_saas_claimed`
- `agency_adoption_claimed`
- `consumer_acceptance_claimed`
- `vendor_compatibility_claimed`
- `production_grade_eta_claimed`

Use language such as `availability sample`, `local diagnostic sampling`,
`diagnostic threshold`, and `needs review`. Do not use language that implies an
uptime evidence workflow, SLA satisfaction, guaranteed uptime, production
readiness, compliance, or acceptance.

## Reliability Summary Model

The reliability summary must have fixed high-level sections:

- `feeds`
- `incidents`
- `backup_restore`
- `alerting`
- `availability_sampling`
- `long_running_operations`
- `claim_flags`

Allowed statuses:

- `ok`: observed source exists and is within configured diagnostic thresholds.
- `needs_review`: observed source exists but is stale, degraded, failing a
  threshold, or requires operator action.
- `missing`: expected safe private diagnostic source is absent.
- `unknown`: source is not instrumented or no database row exists.
- `unhealthy`: observed source reports explicit failure.

Missing data must never become `ok`.

Feed rows must use this fixed order:

- `schedule`
- `vehicle_positions`
- `trip_updates`
- `alerts`

Incident rollups must come from the existing `incident` table only. Include
counts by status, severity, and type, oldest open incident age, and capped
recent items. Recent items may include sanitized ID, type, severity, status,
opened/updated timestamps, and sanitized title/category only. Do not include
raw `details_json`, raw payloads, private logs, tokens, hostnames, or
free-form private text.

Backup/restore, alerting, and availability inputs must come only from safe
private summaries such as:

- `.cache/deployment-doctor/.../summary.json`
- `.cache/operations-notify/.../summary.json`
- `.cache/validator-health/.../summary.json`
- `.cache/operations-reliability/.../summary.json`

Explicit operator-declared booleans or paths may be represented only as safe,
redacted presence/status fields. Do not write raw backup dumps, DB URLs,
webhook values, raw logs, private paths, or secrets.

If no safe source exists, report `missing` or `needs_review`.

## Files Likely To Change

- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_checklist.go`
- `cmd/feed-vehicle-positions/main.go`
- `internal/compliance/`
- `internal/feed/`
- `scripts/operations-reliability.sh`
- `scripts/operations-notify.sh` only if it safely consumes the new summary
- `scripts/deployment-doctor.sh` only if it safely consumes the new summary
- `Makefile`
- `docs/current-status.md`
- `docs/backlog.md`
- `docs/open-questions.md`
- `docs/dependencies.md`
- `docs/decisions.md` if an architecture-significant decision is added
- `docs/handoffs/phase-51.md`
- `docs/handoffs/latest.md`

Do not edit `docs/evidence`.

## Tests

Go tests:

- Reliability summary fixed section order, fixed feed order, allowed statuses,
  and missing-data behavior.
- Incident rollup sanitization, count accuracy, oldest open age, capped recent
  items, and no raw details output.
- Admin route auth, GET-only behavior, `Cache-Control: no-store`, agency query
  scoping, JSON shape, and no public route.
- Vehicle Positions health snapshot persistence success and failure, proving
  persistence failure does not alter public feed status.

Script tests:

- Exact output directory contract.
- Exact five output files.
- Evidence-path rejection.
- Symlink rejection.
- Oversized source rejection.
- Redaction scan.
- No notification sending.
- Claim flags all false.
- Missing safe sources become `missing` or `needs_review`.

Docs and repository checks:

- No forbidden claim language.
- No `docs/evidence` changes.
- Consumer tracker unchanged.

## Performance And Scale

- Add bounded aggregation tests for large `feed_health_snapshot` and
  `incident` histories.
- Add script scale tests with large diagnostic summaries to prove capped
  output.
- Vehicle Positions persistence must use bounded details and must be tested or
  benchmarked so it does not materially alter feed behavior.

## Required Verification

Master-required minimum:

```bash
make validate
make test
git diff --check
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
docker compose -f deploy/docker-compose.yml config
```

Phase-specific checks:

```bash
sh -n scripts/operations-reliability.sh
make operations-reliability
go test ./internal/compliance ./cmd/agency-config ./cmd/feed-vehicle-positions
make test-integration
git diff --exit-code -- docs/evidence
```

`make test-integration` is required because Phase 51 touches DB-backed reads
and Vehicle Positions health persistence using existing tables.

## Exact Consumer Tracker Preservation Checks

```bash
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
git diff --exit-code -- docs/evidence
```

## Open Plan Risks

- Vehicle Positions health persistence must avoid turning public feed requests
  into DB-dependent behavior.
- Schedule health remains out of scope, so it will likely report `unknown`
  until a later approved phase.
- Backup/restore and alerting readiness depend on safe private diagnostics;
  absent diagnostics must be shown as `missing` or `needs_review`.
- Local availability sampling is not evidence of uptime and must stay framed
  as private diagnostic sampling.
