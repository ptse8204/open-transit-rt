# Phase 49 -- External Predictor Runtime Adapter

## Status

Planning approved. Implementation has not started.

## Goal

Phase 49 adds an optional, disabled-by-default runtime path behind
`internal/prediction.Adapter` for a generic operator-owned HTTP predictor
sidecar. The default Trip Updates adapter remains deterministic, and Vehicle
Positions must continue independently when the external adapter is disabled,
misconfigured, slow, unavailable, or returns unsafe output.

This phase must not add TheTransitClock-specific runtime code, start external
processes, weaken Trip Updates output validation, or create production-grade
ETA claims.

## Scope

- Add `TRIP_UPDATES_ADAPTER=external_http` as a private controlled runtime mode.
- Add `TRIP_UPDATES_ADAPTER=external_http_shadow` so deterministic output
  remains public while external output contributes only redacted diagnostics.
- Keep `TRIP_UPDATES_ADAPTER=deterministic` as the default.
- Add shared adapter factory and config validation used by both
  `cmd/feed-trip-updates` and `cmd/agency-config`.
- Add strict generic HTTP sidecar validation:
  - fixed path `/v1/predict/trip-updates`;
  - explicit allowlisted host;
  - no userinfo, query, fragment, or redirects;
  - HTTPS except loopback test stubs;
  - bounded timeout, request size, and response size.
- Add a versioned JSON request/response contract using dedicated sanitized DTOs.
- Preserve existing Trip Updates builder output validation and require external
  outputs to carry confidence plus agency/feed scope.

## Non-Goals

- No TheTransitClock process start, Java/Maven/Tomcat invocation, vendored
  code, named external service call, or named vendor compatibility claim.
- No public unauthenticated API, public route, admin route, queue, scheduler,
  daemon, webhook receiver, or consumer workflow.
- No GTFS auto-editing, publish blocking, evidence packet creation, or writes
  under `docs/evidence`.
- No consumer status changes.
- No CAL-ITP/Caltrans compliance, consumer acceptance, agency adoption, hosted
  SaaS, production-readiness, vendor-compatibility, production AVL reliability,
  real-world ETA accuracy, or production-grade ETA claim.

## Sanitized Request Contract

The external request must use a dedicated DTO. Implementation must never marshal
`telemetry.StoredEvent`, `telemetry.Event`, or `state.Assignment` directly to
the external predictor because those types can contain `PayloadJSON`,
device/driver fields, `ScoreDetails`, and override internals.

Allowed top-level request fields only:

- `schema_version`
- `agency_id`
- `feed_version_id`
- `generated_at`
- `vehicle_positions_url`
- `telemetry`
- `assignments`

Allowed telemetry DTO fields only:

- `vehicle_id`
- `timestamp`
- `lat`
- `lon`
- optional `bearing`
- optional `speed_mps`
- optional `accuracy_m`
- optional `trip_hint` if already present in accepted telemetry

Forbidden telemetry fields include `device_id`, `driver_id`, `payload_json`,
raw vendor payloads, tokens, auth fields, headers, cookies, and credential
material.

Allowed assignment DTO fields only:

- `vehicle_id`
- `feed_version_id`
- `state`
- `service_date`
- `route_id`
- `trip_id`
- `block_id`
- `start_date`
- `start_time`
- `current_stop_sequence`
- `shape_dist_traveled`
- `confidence`
- `assignment_source`
- `reason_codes`
- `degraded_state`
- `manual_override_active`

Forbidden assignment fields include `score_details`, `manual_override_id`,
audit details, raw override reason text, and any internal override payload.

## Runtime Contract

Configuration keys:

| Name | Meaning |
| --- | --- |
| `TRIP_UPDATES_ADAPTER` | `deterministic`, `noop`, `external_http`, or `external_http_shadow`. |
| `TRIP_UPDATES_EXTERNAL_HTTP_URL` | Generic predictor URL; path must be `/v1/predict/trip-updates`. |
| `TRIP_UPDATES_EXTERNAL_HTTP_ALLOWED_HOSTS` | Comma-separated exact host allowlist. |
| `TRIP_UPDATES_EXTERNAL_HTTP_TIMEOUT_SECONDS` | Per-call timeout. |
| `TRIP_UPDATES_EXTERNAL_HTTP_MAX_REQUEST_BYTES` | Encoded request cap. |
| `TRIP_UPDATES_EXTERNAL_HTTP_MAX_RESPONSE_BYTES` | Encoded response cap. |
| `TRIP_UPDATES_EXTERNAL_HTTP_TOKEN_ENV` | Optional env-var name containing bearer token. |

`TRIP_UPDATES_EXTERNAL_HTTP_TOKEN_ENV` is an env-var name only. It must match
the strict uppercase env-name pattern, and the referenced value must exist
before any external call. The token value must never be logged, written,
persisted, returned, or included in diagnostics.

## Failure Semantics

- `external_http`: config, call, timeout, malformed-response, unsafe-output, or
  oversized-response failures return valid empty Trip Updates with
  `adapter_error` or existing adapter-output rejection diagnostics.
- `external_http_shadow`: deterministic output remains the public Trip Updates
  output; external failures record redacted shadow diagnostics only.
- Vehicle Positions, telemetry ingest, GTFS import, assignments, audit state,
  static GTFS publication, Alerts, and admin workflows must continue without
  consulting the external predictor.

## Diagnostics Boundaries

External and shadow diagnostics may include bounded counts, status, reason,
latency duration or bucket, accepted/rejected counts, and
deterministic-vs-external count deltas.

Diagnostics must not store raw request bodies, raw response bodies, raw external
Trip Updates, private host/path details, token values, headers, cookies, DB
URLs, score details, payload JSON, raw private telemetry, or raw vendor data.

## Tests

- Default adapter remains deterministic; `noop` still works.
- `external_http` and `external_http_shadow` are never called unless explicitly
  configured.
- Shared factory/config validation is used by both Trip Updates service and
  agency-config internal builder paths.
- URL validation rejects non-allowlisted hosts, userinfo, query, fragment,
  wrong path, unsafe HTTP, redirects, missing allowlist, oversized request, and
  oversized response.
- Token env-name validation rejects invalid names, missing referenced values,
  and never surfaces token values in diagnostics.
- Test-only `httptest` predictor returns valid output that maps through
  `prediction.Adapter` and existing Trip Updates validation.
- Timeout, `4xx`, `5xx`, malformed JSON, stale response timestamp, low
  confidence, wrong agency/feed, missing scope, and impossible stop sequences
  degrade safely.
- Shadow mode publishes deterministic output only and records bounded redacted
  shadow diagnostics.
- Adapter errors still produce valid GTFS-RT Trip Updates protobuf.
- Architecture tests confirm Vehicle Positions, telemetry ingest, GTFS Studio,
  and Alerts remain uncoupled from external predictor runtime packages.

## Performance And Scale Checks

- Existing `TRIP_UPDATES_MAX_VEHICLES` continues to cap input rows before
  request building.
- Request and response byte caps are enforced before sending/persisting
  diagnostics.
- Timeout tests use local stubs only.
- Any benchmark output is local engineering diagnostics only, not SLA,
  production capacity, readiness, compliance, or ETA-quality proof.

## Docs And Handoff Updates

Implementation should update:

- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-49.md`
- `docs/dependencies.md`
- `docs/decisions.md`
- `docs/backlog.md`
- `docs/open-questions.md`
- `docs/integration-adapter-kit.md`

Keep TheTransitClock documented as candidate-only unless a later phase
explicitly approves named runtime integration with dependency/license review and
retained claim-specific evidence.

## Required Verification

Phase 49 must run and report:

```bash
go test ./internal/prediction ./internal/feed/tripupdates ./cmd/feed-trip-updates ./cmd/agency-config ./internal/architecture
go test ./cmd/feed-vehicle-positions ./internal/feed
make validate
make realtime-quality
make test
make smoke
make test-integration
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

## Closeout Requirements

Phase 49 is not closed until implementation is reviewed by the master agent,
required checks pass or blockers are documented truthfully, the Phase 49
handoff exists, `docs/handoffs/latest.md` is updated, roadmap/status docs are
consistent, no forbidden claims are introduced, no `docs/evidence` files are
edited, and the consumer tracker remains exactly seven `prepared` targets.
