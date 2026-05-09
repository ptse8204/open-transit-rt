# Phase 56 Handoff -- Multi-Agency Hosting Hardening

## Status

Phase 56 is complete for the approved repository-boundary hardening scope.

## What Changed

- Added `internal/tenant` for conservative agency ID validation and
  `/public/agencies/{agency_id}/...` path parsing.
- Added validated path-routed public feed endpoints while preserving existing
  single-agency public routes:
  - `/public/agencies/{agency_id}/feeds.json`
  - `/public/agencies/{agency_id}/gtfs/schedule.zip`
  - `/public/agencies/{agency_id}/gtfsrt/vehicle_positions.pb`
  - `/public/agencies/{agency_id}/gtfsrt/trip_updates.pb`
  - `/public/agencies/{agency_id}/gtfsrt/alerts.pb`
- Kept per-agency JSON/debug routes out of the public route surface. Existing
  JSON debug routes remain authenticated.
- Updated local and OCI Caddy routing so the OCI anonymous public edge exposes
  only public feed paths.
- Added `scripts/multi-agency-hosting.sh`,
  `scripts/test-multi-agency-hosting.sh`, `make multi-agency-hosting`, and
  `make test-multi-agency-hosting`.
- Updated multi-agency, deployment, dependency, decision, roadmap, backlog,
  open-question, status, and latest handoff docs.

## Boundaries Preserved

Phase 56 proves repository-level routing/proxy/tooling boundaries through tests
and docs only. It does not certify production multi-tenant hosting.

No hosted SaaS, production multi-tenant hosting, SLA/uptime,
production-readiness, compliance, agency adoption, consumer acceptance, vendor
compatibility, marketplace approval, or production-grade ETA claim was created.

No retained evidence was created. `docs/evidence/captured`,
`docs/evidence/consumer-submissions/status.json`, consumer current records, and
consumer artifact/packet directories remain unchanged. All seven consumer and
aggregator targets remain `prepared`.

Tenant restore into a shared live database remains blocked. Phase 56 private
diagnostics are not backups, not restore artifacts, and not evidence.

## Verification

Final verification was run from `/Users/edwintse/Downloads/open-transit-rt`.

- `sh -n scripts/multi-agency-hosting.sh scripts/test-multi-agency-hosting.sh`
- `./scripts/test-multi-agency-hosting.sh`
- `go test ./cmd/agency-config ./cmd/feed-vehicle-positions ./cmd/feed-trip-updates ./cmd/feed-alerts ./cmd/telemetry-ingest ./internal/auth ./internal/compliance ./internal/server ./internal/state ./internal/tenant`
- `make validate`
- `make test`
- `make smoke`
- `git diff --check`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json`
- `git diff --exit-code -- docs/evidence/consumer-submissions/current docs/evidence/consumer-submissions/artifacts docs/evidence/consumer-submissions/packets docs/evidence/captured`
- consumer artifact directory scan printed no files
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured`
- `docker compose -f deploy/docker-compose.yml config`
- `INTEGRATION_TESTS=1 make test-integration`

## Next Phase

Proceed to Phase 57 -- Release Packaging And Supply Chain. Keep release,
image, checksum, SBOM, and provenance work evidence-bounded and avoid hosted
service claims.
