# Phase 42 Handoff

## Phase

Phase 42 — Reference Deployment Doctor

## Status

Complete for the read-only reference deployment doctor scope.

Active phase after this handoff: Phase 43 — Operator UX Setup V2 is the next
proposed roadmap phase if the maintainer continues the
CAL-ITP-style/self-hosted gap-closure roadmap.

## What Was Implemented

- Added `scripts/deployment-doctor.sh` and `make deployment-doctor`.
- Added `docs/deployment/reference-deployment-doctor.md`.
- Added `docs/phase-42-reference-deployment-doctor.md`.
- Updated validation scaffolding so `make validate` checks the new script,
  script help path, and Phase 42 docs/handoff files.
- Updated navigation/status docs for the Phase 42 closeout.

The deployment doctor:

- is read-only and does not source private env files;
- inspects already-exported env vars only;
- derives reference env keys from `docs/deployment/oci-reference-env.example`;
- records env key statuses without values;
- checks generated-secret presence/placeholder/minimum length for
  `ADMIN_JWT_SECRET`, `CSRF_SECRET`, and `DEVICE_TOKEN_PEPPER`;
- treats `ADMIN_TOKEN` as optional unless authenticated admin checks are
  requested;
- validates URL syntax and public/admin boundary rules;
- fetches public feed bodies only to temporary files and deletes them after
  recording status, size, checksum, redirect count, effective URL, and content
  type;
- checks public-edge private/admin routes with `HEAD` first and `GET` fallback
  only on `405`;
- checks loopback `/healthz` and `/readyz` for the six reference services;
- runs read-only migration status through `go run ./cmd/migrate status` when
  `DATABASE_URL` is supplied;
- probes PostGIS only when it can do so without passing `DATABASE_URL` as a
  visible `psql` argument;
- checks pinned validator tooling through `scripts/check-validators.sh`;
- checks backup and restore-drill input readiness without creating backups,
  reading backup contents, or running restores;
- records git/release identity;
- verifies the seven consumer tracker targets remain `prepared`;
- writes `summary.json`, `summary.md`, `manifest.json`, and `manifest.md`;
- validates generated JSON;
- runs a final redaction scan over generated output.

Default `make deployment-doctor` is expected to exit `0` even in a local
checkout without deployment env, while reporting blockers, skipped checks, and
unavailable checks. `STRICT_DOCTOR=true` is the mode that fails on blockers.

## What Was Intentionally Not Implemented

- No evidence packet generation.
- No final-root evidence workflow.
- No consumer submission automation.
- No consumer status changes.
- No migration, backup, restore, or write-side deployment action.
- No runtime external predictor integration.
- No production readiness, hosted SaaS, compliance, agency adoption, vendor
  compatibility, or production-grade ETA claim.

## Schema And Interface Changes

- No database schema changes.
- No HTTP API changes.
- No public feed URL or GTFS-RT protobuf contract changes.
- New local operator command interface:
  - `make deployment-doctor`

## Dependency Changes

- No new external dependency was added.
- The script uses existing repo tools and common local tools: `sh`, `curl`,
  `python3`, `go`, git, optional `openssl`, optional `psql`, and existing
  validator tooling.

## Migrations Added

None.

## Checks Run And Blocked Checks

- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed.
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` — passed.
- Exact seven-target consumer tracker check — passed; Google Maps,
  Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and
  transit.land all remain `prepared`.
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json` —
  passed; the tracker file is unchanged.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- `make deployment-doctor` — passed. With no local/reference deployment env or
  running app exported in this checkout, the doctor exited `0` as designed and
  reported 44 blockers, 0 warnings, and 12 unavailable checks. Public feed,
  admin boundary, database, migration, and PostGIS checks were skipped because
  no deployment values were supplied; validator tooling passed; backup and
  restore readiness reported blockers because deployment-owned paths/targets
  were not supplied.

Blocked:

- None for the Phase 42 implementation checks.

## Known Issues

- No agency-owned or agency-approved final-root evidence exists.
- Consumer and aggregator targets remain `prepared` only.
- Deployment doctor output is private diagnostics by default, not evidence.

## Exact Next-Step Recommendation

- Continue with Phase 43 — Operator UX Setup V2.
- Keep Phase 43 scoped to operator remediation UX and private/exportable
  operator checklists without exposing admin routes publicly or claiming
  compliance from UI state.
