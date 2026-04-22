# Phase Handoff Template

## Phase

Phase 12 — Deployment Evidence Hardening (Step 2: Local evidence packet)

## Status

- Partially complete for Step 2.
- Active phase after this handoff: Phase 12.
- Phase 12 is not fully closed.

## What Evidence Was Collected

- Evidence packet: `docs/evidence/captured/local-demo/2026-04-22/`
- Environment: `local-demo`, a developer workstation environment using local Docker/Postgres, Go services, and a temporary loopback HTTP proxy at `http://localhost:8090`.
- Public feed artifacts:
  - `artifacts/public/schedule.zip`
  - `artifacts/public/feeds.json`
  - `artifacts/public/vehicle_positions.pb`
  - `artifacts/public/trip_updates.pb`
  - `artifacts/public/alerts.pb`
- Local anonymous fetch proof for all five public feed paths, including headers, timestamps, bytes, and checksums in `public-feed-proof-2026-04-22.md`.
- Local reverse proxy route map and protected admin/debug 401 checks in `reverse-proxy-tls-2026-04-22.md`.
- Validator records for schedule, Vehicle Positions, Trip Updates, and Alerts in `validator-record-2026-04-22.md` and `artifacts/validation/`.
- Monitoring evidence limited to local service request logs and scorecard health fields in `monitoring-alert-2026-04-22.md`.
- Backup/restore evidence in `backup-restore-drill-2026-04-22.md`, including one manual `pg_dump`, restore into `open_transit_rt_restore_drill_20260422`, restored row counts, and public feed fetch checks against the restored database.
- Manual scorecard export evidence in `scorecard-export-2026-04-22.md` and `artifacts/scorecard/`.
- Artifact checksums in `SHA256SUMS.txt`.

## What Evidence Is Still Missing

- Public HTTPS feed root with real DNS/hostname ownership.
- TLS certificate, validity, renewal posture, and HTTP-to-HTTPS redirect evidence.
- Production reverse proxy or load balancer configuration.
- Publish/update and rollback URL permanence evidence from a hosted deployment.
- Clean production validator records. The captured validator records are real but failed.
- Monitoring dashboard exports, alert rules, notification destination, and a real alert lifecycle.
- Production backup schedule, retention policy, encrypted/storage access boundary, and production restore drill.
- Scheduled scorecard export/job evidence.
- Third-party consumer or aggregator confirmation.

## Blockers Remaining

- No hosted deployment hostname or TLS-terminating proxy was available in the workspace.
- Java was unavailable, so the local MobilityData static GTFS validator run failed.
- The pinned GTFS-RT Docker wrapper reached the validator image, but the image rejected the wrapper's `--schedule` argument.
- No Prometheus/Grafana, alert manager, pager/email/chat destination, or incident/ticket system was configured.
- No production backup service or retention policy was configured.
- No consumer submission or acceptance artifact exists.

## Environment This Evidence Applies To

`local-demo` only. It is local operator/repo evidence and must not be described as hosted production evidence, CAL-ITP compliance, or consumer acceptance.

## Commands Run

Required repo-side checks:

- `make validators-check` — passed.
- `make validate` — passed.
- `make test` — passed.
- `make smoke` — passed.
- `make demo-agency-flow` — passed.
- `make test-integration` — passed.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- `git diff --check` — passed after evidence/document edits.

Evidence collection commands included:

- local public feed `curl -D -` fetches through `http://localhost:8090`;
- `/admin/validation/run` calls for Trip Updates and Alerts after the demo's schedule and Vehicle Positions validation calls;
- local `pg_dump`, isolated restore database creation, and `psql` restore;
- public feed fetches against the restored database;
- manual `/admin/compliance/scorecard` POST and GET.

See `docs/evidence/captured/local-demo/2026-04-22/commands-run-2026-04-22.md`.

## Artifacts Created

- `docs/evidence/captured/local-demo/2026-04-22/README.md`
- `docs/evidence/captured/local-demo/2026-04-22/public-feed-proof-2026-04-22.md`
- `docs/evidence/captured/local-demo/2026-04-22/reverse-proxy-tls-2026-04-22.md`
- `docs/evidence/captured/local-demo/2026-04-22/validator-record-2026-04-22.md`
- `docs/evidence/captured/local-demo/2026-04-22/monitoring-alert-2026-04-22.md`
- `docs/evidence/captured/local-demo/2026-04-22/backup-restore-drill-2026-04-22.md`
- `docs/evidence/captured/local-demo/2026-04-22/scorecard-export-2026-04-22.md`
- `docs/evidence/captured/local-demo/2026-04-22/commands-run-2026-04-22.md`
- `docs/evidence/captured/local-demo/2026-04-22/SHA256SUMS.txt`
- raw public, validation, scorecard, and log artifacts under `docs/evidence/captured/local-demo/2026-04-22/artifacts/`

## Whether Phase 12 Is Fully Closed

Phase 12 is only partially closed. Step 2 collected real local evidence and blockers, but hosted deployment evidence remains absent.

## Exact Recommendation For Phase 13

Do not start Phase 13 implementation work yet if Phase 13 depends on deployment-readiness claims. First run another Phase 12 hosted evidence pass against an actual HTTPS environment and fix validator execution so schedule, Vehicle Positions, Trip Updates, and Alerts can produce clean production validator records.

If Phase 13 is documentation-only planning, keep it explicitly separate from compliance/readiness claims and reference this packet as local evidence only.
