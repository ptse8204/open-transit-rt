# Phase 56 -- Multi-Agency Hosting Hardening

## Status

Complete for the approved repository-boundary hardening scope. Phase 56 added
tenant-safe public feed route validation, path-routed public feed endpoints,
proxy route exposure checks, private `.cache` diagnostics, tests, and docs. It
does not claim hosted SaaS availability or production multi-tenant hosting.

## Goal

Harden repository-level multi-agency boundaries for public feed routing,
backup/restore/export/evidence review, and operations documentation while
keeping all public claims evidence-bounded.

## Scope

- Add a shared agency route validation helper.
- Add explicit path-routed public feed contracts:
  - `/public/agencies/{agency_id}/feeds.json`
  - `/public/agencies/{agency_id}/gtfs/schedule.zip`
  - `/public/agencies/{agency_id}/gtfsrt/vehicle_positions.pb`
  - `/public/agencies/{agency_id}/gtfsrt/trip_updates.pb`
  - `/public/agencies/{agency_id}/gtfsrt/alerts.pb`
- Preserve existing single-agency public URLs unchanged.
- Preserve authenticated admin/debug agency boundaries.
- Add production/local Caddy routing for only the new public feed paths.
- Add private multi-agency hosting diagnostics/tooling under ignored `.cache`.
- Document backup, restore, export, and evidence boundaries for separate
  deployment, separate database, and shared database models.
- Add focused tests for route parsing, public route isolation, Caddy public-edge
  exposure, tooling path safety, and consumer tracker preservation.

## Non-Goals

- No hosted SaaS claim.
- No production multi-tenant hosting claim.
- No SLA/uptime guarantee.
- No production-readiness claim.
- No compliance claim.
- No agency adoption or consumer acceptance claim.
- No global-admin runtime model.
- No tenant restore into a shared live database.
- No migration unless implementation discovers a narrow unavoidable schema gap.
- No public admin/debug route.
- No query-routed protobuf endpoints.
- No consumer contact, portal automation, submission, or status transition.
- No `docs/evidence/captured` writes.

## Files Likely To Change

- `internal/tenant/...`
- `cmd/agency-config/main.go`
- `cmd/agency-config/main_test.go`
- `cmd/feed-vehicle-positions/main.go`
- `cmd/feed-vehicle-positions/main_test.go`
- `cmd/feed-trip-updates/main.go`
- `cmd/feed-trip-updates/main_test.go`
- `cmd/feed-alerts/main.go`
- `cmd/feed-alerts/main_test.go`
- `deploy/Caddyfile.local`
- `deploy/oci/Caddyfile`
- `scripts/multi-agency-hosting.sh`
- `scripts/test-multi-agency-hosting.sh`
- `Makefile`
- `docs/phase-56-multi-agency-hosting-hardening.md`
- `docs/multi-agency-strategy.md`
- `docs/deployment/README.md`
- `docs/deployment/oci-reference-deployment.md`
- `docs/deployment/reference-deployment-doctor.md`
- `docs/dependencies.md`
- `docs/decisions.md`
- `docs/current-status.md`
- `docs/backlog.md`
- `docs/open-questions.md`
- `docs/roadmap-to-calitp-compliance-and-gap-closure.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-56.md`

Execution must not change:

- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/`
- `docs/evidence/consumer-submissions/artifacts/`
- `docs/evidence/consumer-submissions/packets/`
- `docs/evidence/captured/`

## Safety Boundaries

Agency ID path segments must reject:

- empty values;
- `.` and `..`;
- leading dot;
- slashes;
- encoded slashes;
- backslashes;
- path traversal;
- query/fragment-derived agency IDs;
- unsafe characters outside a conservative ASCII set.

Public per-agency feed routes must route only by validated path agency ID.
Existing single-agency public routes remain service-instance scoped.

Admin reads and writes must continue deriving agency from authenticated
principals, not client-supplied agency IDs. Public JSON/debug endpoints remain
authenticated. No new public debug routes may be added.

Public Caddy routes may expose only public feed paths. `/admin/*`,
`/admin/debug/*`, `/v1/events`, `/v1/telemetry`, `/metrics`, GTFS Studio, and
JSON debug endpoints must remain off the anonymous public edge.

Backup/export diagnostics must default to ignored `.cache`, mode `0700`, and
reject symlink, traversal, and evidence-like paths. Any tenant export added in
Phase 56 is a private redacted support diagnostic, not a restorable backup and
not evidence.

## Evidence And Claim Boundaries

Claim flags remain false for:

- compliance;
- consumer acceptance;
- consumer ingestion;
- agency adoption;
- hosted SaaS;
- production readiness;
- SLA/uptime;
- vendor compatibility;
- production-grade ETA.

Phase 56 creates no retained evidence and writes nothing under
`docs/evidence/captured`.

OCI pilot evidence remains hosted/operator pilot evidence only. Consumer
tracker state remains unchanged; all seven targets remain `prepared`.

## Implementation Details

Add a small shared tenant routing helper, for example `internal/tenant`, that:

- validates agency ID path segments;
- parses `/public/agencies/{agency_id}/...` paths;
- exposes tests for valid and invalid values.

Public routing:

- Keep existing `/public/feeds.json`, `/public/gtfs/schedule.zip`, and
  `/public/gtfsrt/*.pb` routes unchanged.
- Add path-routed agency feed discovery and schedule ZIP in `cmd/agency-config`.
- Add path-routed protobuf routes in Vehicle Positions, Trip Updates, and
  Alerts services.
- Do not add path-routed public JSON debug routes.
- For path-routed Trip Updates, allow the configured Vehicle Positions URL to
  use either `/public/gtfsrt/vehicle_positions.pb` or
  `/public/agencies/{agency_id}/gtfsrt/vehicle_positions.pb`.
- Builders/factories must build or query only the requested agency, not list all
  agencies.

Reverse proxy:

- Add local and OCI Caddy public route matchers for
  `/public/agencies/*/feeds.json`, `/public/agencies/*/gtfs/schedule.zip`, and
  `/public/agencies/*/gtfsrt/*.pb`.
- Keep final fallback `404`.
- Do not expose admin, telemetry ingest, metrics, Studio, or JSON debug on the
  anonymous public OCI edge.

Operations model:

- Preferred isolation order is separate deployment per agency, then separate DB
  per agency, then shared DB only with explicit tested route/export boundaries.
- Restore remains full database restore into a restore database. Tenant restore
  into a shared live DB is blocked until a later approved phase.
- Any Phase 56 script should write only `.cache/multi-agency-hosting/<timestamp>`
  diagnostics, with claim flags false and no retained evidence.

Checkpoint strategy:

- `Phase 56 -- Checkpoint 000001: add multi-agency hosting plan`
- `Phase 56 -- Checkpoint 000002: add tenant route validation and public routes`
- `Phase 56 -- Checkpoint 000003: add multi-agency diagnostics and proxy checks`
- `Phase 56 -- Checkpoint 000004: close multi-agency hosting handoff and status`

## Tests

Required tests:

- Unit tests for valid/invalid agency route parsing.
- Public route tests proving agency A path returns agency A data and agency B
  path returns agency B data.
- Existing single-agency routes still work.
- Admin/auth tests proving conflicting `agency_id` remains `403`.
- Tests proving public JSON debug routes remain authenticated and no per-agency
  JSON debug route is added.
- Caddy route tests proving anonymous public edge exposes only public feed
  paths and not admin/debug/metrics/telemetry/Studio routes.
- Script tests for `.cache` defaults, symlink/traversal/evidence-like path
  rejection, no claim flags, and no evidence writes.
- DB-backed integration tests where existing coverage supports it.

## Performance And Scale Tests

- Route parsing must be O(1).
- Path-routed feed generation must build only the requested agency.
- Add a bounded synthetic `agency-a` / `agency-b` smoke fixture.
- Private export diagnostics must cap file size and row counts, and must not
  copy raw backups, GTFS ZIPs, protobufs, raw telemetry, raw logs, or evidence
  artifacts.

## Docs, Status, And Handoff Updates

Close Phase 56 by updating:

- this phase document;
- `docs/handoffs/phase-56.md`;
- `docs/handoffs/latest.md`;
- `docs/current-status.md`;
- `docs/backlog.md`;
- `docs/open-questions.md`;
- `docs/roadmap-to-calitp-compliance-and-gap-closure.md`;
- `docs/multi-agency-strategy.md`;
- relevant deployment docs;
- `docs/dependencies.md`;
- `docs/decisions.md`.

The handoff must say Phase 56 proves repository-level boundaries through tests
and docs only. It must explicitly state no hosted SaaS, production
multi-tenant hosting, SLA/uptime, production-readiness, compliance, agency
adoption, consumer acceptance, or retained evidence claim was created.

## Implementation Summary

Phase 56 added `internal/tenant` for conservative agency path segment
validation and `/public/agencies/{agency_id}/...` route parsing. The public
feed services now preserve existing single-agency routes and add these
validated path-routed protobuf/static/discovery routes:

- `/public/agencies/{agency_id}/feeds.json`
- `/public/agencies/{agency_id}/gtfs/schedule.zip`
- `/public/agencies/{agency_id}/gtfsrt/vehicle_positions.pb`
- `/public/agencies/{agency_id}/gtfsrt/trip_updates.pb`
- `/public/agencies/{agency_id}/gtfsrt/alerts.pb`

Per-agency JSON/debug routes were not added. Existing debug JSON routes remain
authenticated. The OCI Caddyfile exposes only anonymous public feed paths and
keeps admin, telemetry, metrics, Studio, and JSON debug routes off the public
edge.

Phase 56 also added `scripts/multi-agency-hosting.sh`,
`scripts/test-multi-agency-hosting.sh`, and Make targets for private local
diagnostics. The diagnostic writes exactly `summary.json`, `summary.md`,
`manifest.json`, and `manifest.md` under ignored `.cache` output by default,
rejects symlink/traversal/evidence-like output paths, records false claim
flags, and treats tenant restore into a shared live database as blocked.

This phase created no retained evidence, wrote nothing under
`docs/evidence/captured`, changed no consumer tracker state, contacted no
consumer, and made no hosted SaaS, production multi-tenant hosting,
production-readiness, SLA/uptime, compliance, agency adoption, consumer
acceptance, vendor compatibility, marketplace approval, or production-grade ETA
claim.

## Required Verification Commands

Run and report:

```bash
sh -n scripts/multi-agency-hosting.sh scripts/test-multi-agency-hosting.sh
./scripts/test-multi-agency-hosting.sh
go test ./cmd/agency-config ./cmd/feed-vehicle-positions ./cmd/feed-trip-updates ./cmd/feed-alerts ./cmd/telemetry-ingest ./internal/auth ./internal/compliance ./internal/server ./internal/state
make validate
make test
make smoke
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
git diff --exit-code -- docs/evidence/consumer-submissions/current docs/evidence/consumer-submissions/artifacts docs/evidence/consumer-submissions/packets docs/evidence/captured
find docs/evidence/consumer-submissions/artifacts -mindepth 2 -maxdepth 2 -type f ! -name README.md -print
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured
docker compose -f deploy/docker-compose.yml config
```

Run `INTEGRATION_TESTS=1 make test-integration` if the local database is
available and record any environment blocker truthfully.

The `find` command must print no files for the current Phase 56 state.
