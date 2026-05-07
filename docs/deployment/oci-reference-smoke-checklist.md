# OCI/OCL Reference Smoke Checklist

This guide is documentation for reproducing a self-hosted deployment pattern. It is not evidence that a deployment was run, accepted, compliant, agency-approved, or production-ready.

Use this checklist after installing or updating the self-hosted OCI/OCL-style
reference deployment. Results are private operator verification until reviewed
and redacted. This checklist does not create final-root proof, consumer
submission evidence, compliance evidence, agency approval, or production
readiness evidence.

## Setup

```sh
BASE=https://feeds.example.org
OUT=/opt/open-transit-rt/evidence/smoke-$(date -u +%Y-%m-%d)
mkdir -p "$OUT"
```

## Server And Service Checks

- Confirm deployment-owned DNS resolves to the server.
- Confirm ports `80` and `443` are open at the firewall and cloud edge.
- Confirm Open Transit RT services bind to loopback, not public interfaces.
- Confirm systemd services are active:

```sh
sudo systemctl is-active open-transit-agency-config.service
sudo systemctl is-active open-transit-telemetry-ingest.service
sudo systemctl is-active open-transit-feed-vehicle-positions.service
sudo systemctl is-active open-transit-feed-trip-updates.service
sudo systemctl is-active open-transit-feed-alerts.service
```

## Local Health Checks

Run from the server:

```sh
curl -fsS http://127.0.0.1:8081/healthz
curl -fsS http://127.0.0.1:8082/healthz
curl -fsS http://127.0.0.1:8083/healthz
curl -fsS http://127.0.0.1:8084/healthz
curl -fsS http://127.0.0.1:8085/healthz
curl -fsS http://127.0.0.1:8081/readyz
```

If any check fails, inspect the matching systemd unit and environment file
before continuing.

## Five Public Feed URL Checks

Run from outside the server:

```sh
curl -fsS -D "$OUT/feeds.headers" "$BASE/public/feeds.json" -o "$OUT/feeds.json"
curl -fsS -D "$OUT/schedule.headers" "$BASE/public/gtfs/schedule.zip" -o "$OUT/schedule.zip"
curl -fsS -D "$OUT/vehicle_positions.headers" "$BASE/public/gtfsrt/vehicle_positions.pb" -o "$OUT/vehicle_positions.pb"
curl -fsS -D "$OUT/trip_updates.headers" "$BASE/public/gtfsrt/trip_updates.pb" -o "$OUT/trip_updates.pb"
curl -fsS -D "$OUT/alerts.headers" "$BASE/public/gtfsrt/alerts.pb" -o "$OUT/alerts.pb"
shasum -a 256 "$OUT"/*
```

Expected result: all five public URLs return successfully with non-error HTTP
status, retained headers, retained bytes, and checksums.

## Public/Private Boundary Checks

Run from outside the server:

```sh
curl -o /dev/null -sS -w '%{http_code}\n' "$BASE/admin/"
curl -o /dev/null -sS -w '%{http_code}\n' "$BASE/admin/debug/"
curl -o /dev/null -sS -w '%{http_code}\n' "$BASE/v1/events"
curl -o /dev/null -sS -w '%{http_code}\n' "$BASE/metrics"
```

Expected result: public edge denies these routes, commonly with `404` from the
reference Caddyfile. Admin/debug/studio access belongs behind SSH tunnel, VPN,
or another operator-controlled auth-aware boundary.

## Validator Checks

From the repo checkout used for the deployment:

```sh
make validators-check
```

After feeds are reachable, run the private validator-cycle dry-run:

```sh
ENVIRONMENT_NAME=reference \
EVIDENCE_OUTPUT_DIR="$OUT" \
ADMIN_BASE_URL=http://127.0.0.1:8081 \
ADMIN_TOKEN=replace-with-generated-secret \
/opt/open-transit-rt/app/scripts/pilot-ops.sh validator-cycle --dry-run
```

Validator output is an operator check unless it is intentionally retained,
reviewed, and redacted as evidence.

## Operations Helper Dry-Runs

Run:

```sh
ENVIRONMENT_NAME=reference EVIDENCE_OUTPUT_DIR="$OUT" \
  /opt/open-transit-rt/app/scripts/pilot-ops.sh backup --dry-run

ENVIRONMENT_NAME=reference EVIDENCE_OUTPUT_DIR="$OUT" \
  /opt/open-transit-rt/app/scripts/pilot-ops.sh restore-drill --dry-run

ENVIRONMENT_NAME=reference EVIDENCE_OUTPUT_DIR="$OUT" PUBLIC_BASE_URL="$BASE" \
  /opt/open-transit-rt/app/scripts/pilot-ops.sh feed-monitor --dry-run

ENVIRONMENT_NAME=reference EVIDENCE_OUTPUT_DIR="$OUT" \
  ADMIN_BASE_URL=http://127.0.0.1:8081 ADMIN_TOKEN=replace-with-generated-secret \
  /opt/open-transit-rt/app/scripts/pilot-ops.sh scorecard-export --dry-run
```

Expected result: every helper prints the target environment and planned action
without exposing secrets. Enable timers only after dry-runs are understood and
private env files are populated correctly.

## Redaction Review

Before sharing results:

- Remove or redact tokens, database URLs, webhook URLs, private hosts, and raw
  access logs.
- Keep raw database backups and raw GTFS ZIPs private unless reviewed.
- Keep operator evidence under `/opt/open-transit-rt/evidence/` private until
  `docs/evidence/redaction-policy.md` review is complete.
- Do not convert this smoke output into claims about consumer ingestion,
  compliance, agency approval, final-root proof, or production readiness.
