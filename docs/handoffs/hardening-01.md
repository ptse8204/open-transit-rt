# Phase Handoff Template

All future phase handoff files must use this structure unless the phase explicitly documents a reason to diverge.

## Phase

Post-Phase-8 Hardening 01 — Pilot-readiness security and operability slice

## Status

- Complete
- Active phase after this handoff: no Phase 9 is defined; next work should continue hardening/operability.

## What Was Implemented

- Added HS256 admin JWT validation with required `sub`, `agency_id`, `iat`, `exp`, `iss`, and `aud`; DB-backed role loading from `agency_user` / `role_binding`; optional `admin_session` cookie support; and CSRF validation for all cookie-authenticated unsafe admin methods.
- Protected admin routes in `agency-config`, `feed-alerts`, `gtfs-studio`, and protected JSON debug routes for Vehicle Positions, Trip Updates, Alerts, and telemetry events.
- Kept public `.pb` feed endpoints anonymous and stable.
- Moved admin actor and agency scope to auth context. Request `agency_id` fields/query params that conflict with auth scope return `403`.
- Replaced validator command execution with server-side `validator_id` allowlists. Requests may provide only `validator_id`, `feed_type`, and optional `feed_version_id`.
- Validator execution now uses generated/fetched server-owned local artifacts/temp files, argv-based `exec.CommandContext`, timeouts, stdout/stderr caps, report-size caps, and temp/output confinement.
- Realtime validation derives protobuf bytes from server-owned feed URLs and writes local temp files for Vehicle Positions, Trip Updates, and Alerts validation.
- Added opaque telemetry device Bearer token verification with peppered HMAC hashes, active credential status checks, and agency/device/vehicle binding checks.
- Added admin-managed `POST /admin/devices/rebind`, which rotates the token/binding and audit-logs the change.
- Added partial unique current-assignment index plus per-agency/per-vehicle advisory transaction lock in `SaveAssignment`.
- Split scorecard behavior: GET reads latest stored scorecard; POST recomputes/stores.
- Made prediction review incidents idempotent with dedupe key, last-seen timestamp, and occurrence counter.
- Added schedule ZIP `ETag`, `X-Checksum-SHA256`, and max payload guard.
- Added `make smoke` / `task smoke`.

## What Was Designed But Intentionally Not Implemented Yet

- Full login UI, SSO/OIDC, and session management UI.
- Server-side `jti` replay tracking; JWTs emit `jti`, but replay storage is deferred.
- CI/prod installation of exact canonical validator binaries/images. Static validator version is documented as MobilityData GTFS Validator `v7.1.0`; GTFS Realtime validator must still be pinned by immutable digest in automation.
- Full Prometheus metrics and request-log middleware.
- External consumer submission API integrations.

## Schema And Interface Changes

- New admin auth header: `Authorization: Bearer <admin-jwt>`.
- Optional browser admin auth uses signed `admin_session` JWT cookie plus CSRF token for unsafe methods.
- Telemetry ingest requires `Authorization: Bearer <device-token>`.
- `/admin/validation/run` request shape is only `validator_id`, `feed_type`, and optional `feed_version_id`.
- `/admin/devices/rebind` returns a one-time plaintext replacement device token.
- JSON debug routes require admin read auth; `.pb` routes remain anonymous.

## Dependency Changes

- No new Go module dependency was added.
- Validator integration now depends on server-owned `GTFS_VALIDATOR_PATH` and `GTFS_RT_VALIDATOR_PATH`/`GTFS_RT_VALIDATOR_ARGS`, not request-provided shell text.
- Production services require `ADMIN_JWT_SECRET`, `ADMIN_JWT_ISSUER`, `ADMIN_JWT_AUDIENCE`, `CSRF_SECRET`, and `DEVICE_TOKEN_PEPPER`.

## Migrations Added

- `db/migrations/000008_production_hardening.sql`
  - adds `device_credential.vehicle_id` and `last_used_at`
  - adds active device credential lookup index
  - deduplicates current assignments and adds `vehicle_trip_assignment_current_uidx`
  - adds `incident.dedupe_key`, `last_seen_at`, and `occurrence_count`
  - adds active prediction-review dedupe unique index

## Tests Added And Results

- Added JWT, secret-rotation, CSRF, device-token, validator allowlist, unauthenticated admin/debug rejection, telemetry device-token rejection, and concurrent assignment tests.
- Latest observed result during implementation: `go test ./...` passed.

## Checks Run And Blocked Checks

- Run during implementation:
  - `gofmt -w ./cmd ./internal`
  - `go test ./...`
- Still expected before final handoff/PR:
  - `make validate`
  - `make smoke`
  - `make test-integration`
  - `docker compose -f deploy/docker-compose.yml config`
  - `git diff --check`

## Known Issues

- `SCHEDULE_ZIP_MAX_BYTES` bounds payload size and the endpoint emits ETag/checksum, but ZIP generation is still on-demand and not yet a shared materialized cache service.
- Realtime validator local artifact generation depends on each feed service exposing/building bytes in-process; `agency-config` stores `not_run` when realtime artifacts are unavailable.
- Admin JWT issuance is via local CLI/bootstrap; production identity integration is not included.

## Exact Next-Step Recommendation

- First files to read:
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `internal/auth/`
  - `internal/devices/`
  - `internal/compliance/validation.go`
- First files likely to edit:
  - CI/deploy validator setup
  - `internal/server/` for structured request logging and metrics
  - deployment docs for reverse proxy/TLS
- Commands to run before coding:
  - `make test`
  - `make smoke`
  - `make test-integration`
  - `make validate`
- Known blockers:
  - Docker must be running for DB-backed integration checks.
  - Canonical validator binaries/images are not installed by repo automation yet.
- Recommended first implementation slice:
  - Add structured request logging, `/metrics`, and stronger readiness checks, then pin validator distributions in CI/prod automation.
