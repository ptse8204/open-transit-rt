# Phase 12 — Deployment Evidence Hardening

## Status

Closed for the OCI pilot evidence scope on 2026-04-24. Step 1 repo scaffolding,
Step 2 local evidence, Step 3 validator tooling hardening, and the hosted
`oci-pilot` evidence packet are complete.

## Purpose

Phase 12 defines the deployment-side evidence needed before stronger readiness claims can be made for any real Open Transit RT installation.

This phase does **not** add backend product features. It hardens proof quality for already-implemented repo capabilities by requiring deployment artifacts that demonstrate stable public hosting, operational controls, and evidence retention.

## Scope

Phase 12 focuses on collecting and retaining reproducible deployment evidence for:

1. Stable HTTPS public feed root.
2. Reverse proxy and TLS termination proof.
3. Production validator records for schedule and realtime feeds.
4. Monitoring and alerting evidence.
5. Backup/restore and operations runbooks.
6. Compliance scorecard export evidence.

All claims must follow `docs/prompts/calitp-truthfulness.md` and clearly distinguish:
- repo proof,
- deployment/operator proof,
- third-party/consumer proof.

## Evidence Requirements

### 1) Stable HTTPS Public Feed Root

Required proof for a real deployment:
- one canonical HTTPS root host for public feeds,
- stable paths for:
  - `/public/gtfs/schedule.zip`
  - `/public/feeds.json`
  - `/public/gtfsrt/vehicle_positions.pb`
  - `/public/gtfsrt/trip_updates.pb`
  - `/public/gtfsrt/alerts.pb`
- no-login anonymous retrieval,
- URL permanence across at least one publish/rollback cycle.

Suggested artifacts:
- dated curl transcript,
- response headers (including cache/metadata headers where applicable),
- before/after publish and rollback fetch comparisons,
- deployment inventory noting canonical hostname ownership.

### 2) Reverse Proxy / TLS Proof

Required proof:
- reverse proxy routing map from public HTTPS paths to internal services,
- TLS certificate validity and renewal posture,
- redirect behavior from HTTP to HTTPS if HTTP is exposed,
- expected behavior for admin/debug routes (auth-protected/non-public).

Suggested artifacts:
- redacted proxy config snippets,
- certificate chain and expiry evidence,
- periodic renewal check logs,
- architecture diagram showing trust boundary.

### 3) Production Validator Records

Required proof:
- latest static GTFS validation record for deployed `schedule.zip`,
- latest GTFS-RT validation records for Vehicle Positions, Trip Updates, and Alerts,
- recorded timestamps, validator IDs/versions, and result levels,
- explicit retention location for historical runs.

Rules:
- validator success is required evidence for quality, but is **not** proof of consumer acceptance,
- warning/error levels must be preserved without selective omission.

### 4) Monitoring / Alerting Evidence

Required proof:
- monitored public feed availability,
- freshness and generation-latency visibility,
- alert rules and delivery path (pager/email/chat),
- at least one acknowledged alert lifecycle example.

Suggested artifacts:
- dashboard screenshots or exported panels,
- alert rule definitions,
- incident timeline showing detect → acknowledge → resolve,
- SLO/SLA note for current pilot posture.

### 5) Backup/Restore And Operations Runbooks

Required proof:
- documented backup schedule and retention,
- restore procedure with step-by-step commands,
- last successful restore drill with timestamp and operator identity,
- incident runbook for feed outage and validator failure.

Suggested artifacts:
- runbook markdown docs,
- restore drill log transcript,
- post-drill issue list and follow-up actions.

### 6) Scorecard Export Evidence

Required proof:
- exportable compliance scorecard snapshots from deployment,
- timestamped history demonstrating repeatable generation,
- storage location and retention policy for exported scorecards.

Suggested artifacts:
- exported JSON/CSV files,
- cron/job definitions,
- operator playbook for manual export fallback.

## Acceptance Criteria

Phase 12 is complete only when all are true:

- A deployment evidence packet exists for one real environment with dated artifacts.
- Stable HTTPS feed root and URL permanence are demonstrated across update and rollback.
- Reverse proxy and TLS evidence is captured and reviewable.
- Production validator records exist for schedule + all three realtime feeds.
- Monitoring/alerting evidence includes at least one real alert lifecycle.
- Backup/restore runbooks exist and a restore drill result is recorded.
- Scorecard export evidence is retained with timestamped outputs.
- `docs/current-status.md` and `docs/handoffs/latest.md` are updated truthfully.

## Commands To Run

These commands validate repo-side workflows while deployment evidence is collected:

```bash
make validators-check
make validate
make test
make smoke
make demo-agency-flow
make test-integration
docker compose -f deploy/docker-compose.yml config
git diff --check
```

Deployment operators may also run environment-specific fetch/validator/monitoring commands, but those are deployment-owned and must be recorded in deployment runbooks.

## Explicit Non-Goals

Phase 12 does **not**:
- claim full CAL-ITP/Caltrans compliance,
- claim consumer acceptance by Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, transit.land, or any other downstream consumer,
- add new realtime prediction algorithms,
- introduce new backend product features,
- reopen implementation work from Phases 9–11,
- claim universal production readiness for all agencies.

## Deliverables

- `docs/phase-12-deployment-evidence-hardening.md` (this plan).
- Updated handoff/current-status pointers for next execution.
- Identified artifact folders or links where deployment evidence will be stored in a future execution pass.
## Step 1 (Repo Scaffolding) Completion Note

Phase 12 Step 1 is complete as a repository documentation pass. It added concrete runbooks under `docs/runbooks/` and evidence templates/placeholders under `docs/evidence/`.

This step does not include real hosted deployment artifacts. Deployment/operator evidence collection remains pending in later Phase 12 execution slices.

## Step 2 (Local Evidence Packet) Completion Note

Phase 12 Step 2 created `docs/evidence/captured/local-demo/2026-04-22/`.

This packet contains real local artifacts for:

- anonymous local HTTP fetches of all five public feed paths;
- local reverse proxy routing through the demo proxy;
- protected admin/debug anonymous 401 checks;
- validator records for schedule plus Vehicle Positions, Trip Updates, and Alerts;
- a local Postgres dump/restore drill with restored feed fetch checks;
- a manual scorecard export.

The packet is intentionally marked partial. It does not include hosted HTTPS evidence, TLS/certificate renewal proof, production monitoring alert lifecycle, clean validator results, production backup retention, rollback URL permanence, or third-party consumer confirmation.

## Step 3 (Hosted Closure Tooling) Completion Note

Phase 12 Step 3 hardened the repo-side tooling needed for a future hosted closure pass:

- the pinned GTFS-RT validator wrapper generated by `make validators-install` now drives the MobilityData validator webapp API instead of passing unsupported CLI flags to the pinned image;
- `make validators-check` now verifies Java for the static validator and Docker, `curl`, and `python3` for the GTFS-RT wrapper;
- no hosted production/pilot evidence was collected in this step.

## Hosted OCI Pilot Closure Note

Phase 12 hosted evidence was collected at `docs/evidence/captured/oci-pilot/2026-04-24/`.

The packet includes:

- public HTTPS fetch proof for `schedule.zip`, `feeds.json`, Vehicle Positions, Trip Updates, and Alerts;
- public-edge auth boundary checks and SSH-tunneled admin auth checks;
- TLS certificate and HTTP-to-HTTPS redirect evidence;
- clean hosted validator records for schedule and all three realtime feeds;
- pre-update, post-update, transient-update, and deployment data-restore rollback `feeds.json` snapshots;
- final current-live recheck showing active `gtfs-import-3`, all required hosted validators passed, and `canonical_validation_complete=true`;
- operator-supplied reverse proxy, monitoring, alert lifecycle, backup, restore, and scorecard job artifacts.

Closure gate:

```bash
EVIDENCE_PACKET_DIR=docs/evidence/captured/oci-pilot/2026-04-24 make audit-hosted-evidence
```

Result: passed.

This closes Phase 12 for hosted/operator evidence. It does not claim Cal-ITP
compliance or third-party consumer acceptance.
