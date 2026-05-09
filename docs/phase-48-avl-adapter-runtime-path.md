# Phase 48 -- AVL Adapter Runtime Path

## Status

Planning approved. Implementation has not started.

## Goal

Phase 48 adds a private, authorized send-mode pattern to the existing AVL
adapter example while preserving `/v1/telemetry` as the only runtime ingest
boundary.

The phase turns the synthetic dry-run adapter into a safer reference client for
deployment-owned AVL processes. It must not add a public API, named vendor
support, evidence packets, consumer workflow changes, or production AVL claims.

## Scope

- Add mutually exclusive `--dry-run` and `--send` modes to
  `cmd/avl-vendor-adapter`.
- Preserve current dry-run no-network stdout/stderr JSON behavior.
- Add strict send manifest validation with env-only token references.
- Send one transformed `telemetry.Event` per `POST /v1/telemetry` request.
- Write bounded private diagnostics under `.cache/avl-vendor-adapter/` by
  default.
- Add retry/backoff only for retryable per-record send failures.
- Update Phase 48 docs, status, handoff, dependency, and decision language.
- Keep `make validate` checks no-network.

## Non-Goals

- No new ingest API, public route, admin route, queue, scheduler, daemon, or
  webhook receiver.
- No `/v1/telemetry` auth or payload-contract change.
- No weakening of device token binding, admin auth, CSRF, role checks, or agency
  isolation.
- No real vendor payloads, real IDs, credentials, private telemetry, or private
  logs.
- No writes under `docs/evidence`.
- No consumer tracker changes.
- No CAL-ITP/Caltrans compliance, consumer acceptance, agency adoption, hosted
  SaaS, production-readiness, vendor-compatibility, production AVL reliability,
  or production-grade ETA claim.

## Runtime Contract

CLI shape:

```bash
go run ./cmd/avl-vendor-adapter --dry-run --mapping <mapping.json> <payload.json>

go run ./cmd/avl-vendor-adapter --send \
  --mapping <mapping.json> \
  --manifest <send-manifest.json> \
  <payload.json>
```

Configuration contract:

| Name | Meaning | Default |
| --- | --- | --- |
| `AVL_ADAPTER_TELEMETRY_URL` | Target URL env named by manifest `telemetry_url_env`; must end in `/v1/telemetry`. | Required for send mode |
| `AVL_ADAPTER_OUTPUT_DIR` | Optional output directory. | `.cache/avl-vendor-adapter/<timestamp>` |
| `AVL_ADAPTER_TIMEOUT` | Request timeout. | `10s` |
| `AVL_ADAPTER_MAX_RETRIES` | Retry count per event. | `2` |
| `AVL_ADAPTER_RETRY_INITIAL_DELAY` | Initial retry delay. | `250ms` |
| `AVL_ADAPTER_RETRY_MAX_DELAY` | Maximum retry delay. | `2s` |
| `AVL_ADAPTER_FAIL_ON_WARNINGS` | Blocks non-stale/non-future warnings when `true`. | `false` |
| `AVL_ADAPTER_REFERENCE_TIME` | RFC3339 reference time for stale/future review. | Current UTC |
| `AVL_ADAPTER_STALE_THRESHOLD` | Stale timestamp threshold. | `90s` |
| `AVL_ADAPTER_FUTURE_THRESHOLD` | Future timestamp threshold. | `30s` |

Token values are read only from env vars named by each manifest credential row's
`token_env`.

## Send Manifest

The send manifest uses this strict JSON schema:

```json
{
  "schema_version": "avl-adapter-send.v1",
  "telemetry_url_env": "AVL_ADAPTER_TELEMETRY_URL",
  "credentials": [
    {
      "agency_id": "demo-agency",
      "device_id": "device-1",
      "vehicle_id": "bus-1",
      "token_env": "AVL_ADAPTER_DEVICE_TOKEN",
      "notes": "Synthetic local reference binding."
    }
  ]
}
```

Validation rules:

- Reject unknown fields at every level.
- Require exact `schema_version = "avl-adapter-send.v1"`.
- Require env names to match `^[A-Z_][A-Z0-9_]*$`.
- Reject duplicate `(agency_id, device_id, vehicle_id)` credential rows.
- Reject missing credential mapping for any transformed event.
- Reject token-looking values, bearer strings, authorization-header values, DB
  URLs, private-key markers, or credential values in manifest fields.
- Treat `notes` as optional public-safe text only.
- Never write token env values to output.

## Send Behavior

- Full preflight validates mapping, payload transform, manifest, target URL,
  token env presence, output path, and stale/future blockers before any send.
- Any transform, stale/future, manifest, path, config, or credential blocker
  sends zero records and exits nonzero.
- Dry-run may keep stale/future timestamps as warnings. Send mode always treats
  stale/future transformed records as batch blockers.
- If preflight passes, send one `telemetry.Event` POST per transformed record in
  input order.
- Treat `201 Created` and `202 Accepted` as success.
- Do not retry `400`, `401`, `404`, or `405`.
- Retry network errors, timeouts, `408`, `429`, and `5xx` per event.
- Stop on the first terminal failure or retry exhaustion, mark later records as
  `skipped_after_failure`, and exit nonzero.
- Reruns may produce duplicate or out-of-order telemetry; existing
  `/v1/telemetry` storage semantics classify those cases.

## Output Contract

Default output is `.cache/avl-vendor-adapter/<timestamp>`. Send mode writes
exactly:

- `summary.json`
- `summary.md`
- `manifest.json`
- `manifest.md`
- `diagnostics.json`

`summary.json` includes:

- `generated_at`
- `mode`
- `dry_run`
- `telemetry_url_path`
- `telemetry_target_loopback`
- optional `telemetry_host_ref`
- `output_label`
- optional `output_ref`
- transformed/sent/succeeded/failed/skipped counts
- retry totals
- duration fields
- false claim flags

The false flags are:

- `external_evidence_created=false`
- `consumer_statuses_changed=false`
- `compliance_claimed=false`
- `production_readiness_claimed=false`
- `hosted_saas_claimed=false`
- `agency_adoption_claimed=false`
- `consumer_acceptance_claimed=false`
- `vendor_compatibility_claimed=false`
- `production_avl_reliability_claimed=false`
- `production_grade_eta_claimed=false`

`diagnostics.json` uses only redaction-safe fields such as `record_index`,
`credential_ref`, outcome, attempts, duration, safe parsed success fields, and
response SHA-256. It must not include raw response bodies, raw telemetry
payloads, raw vendor payloads, raw agency/device/vehicle IDs, token values,
authorization headers, cookies, DB URLs, or private-key material.

Stdout prints only redaction-safe counts and a repo-relative `.cache` output
path when possible. Stderr prints only fatal redacted errors.

## Safety And Redaction

Send mode must target only `/v1/telemetry`. Non-loopback credentialed sends
require HTTPS. Target URLs reject userinfo, query strings, fragments, and wrong
paths.

Output paths reject symlinks and `docs/evidence`. Generated files and terminal
strings are scanned before successful completion for token values,
Authorization headers, raw response bodies, DB URLs, private-key markers,
bearer/cookie values, and raw payload fields.

`credential_ref` is a deterministic redacted reference for local correlation. It
is not public evidence and must not be described as proof of integration.

## Implementation Notes

- Keep CLI parsing in `cmd/avl-vendor-adapter`.
- Put manifest/config validation, redacted references, target URL validation,
  output writing, redaction scanning, and send client helpers in
  `internal/avladapter`.
- Use injectable HTTP and sleeper/backoff hooks so retry tests do not sleep.
- Keep helper APIs internal and scoped to the adapter; do not move vendor logic
  into telemetry ingest, matching, Vehicle Positions, or Trip Updates.

## Tests

Required focused tests:

- Dry-run performs no network I/O and preserves current stdout/stderr JSON.
- Send posts to `/v1/telemetry` with bearer auth and JSON content type.
- Stale and future timestamps each block send mode with zero HTTP requests.
- Other warnings send by default; `AVL_ADAPTER_FAIL_ON_WARNINGS=true` blocks
  them.
- The exact send output file set and required fields are present.
- Diagnostics use `record_index` and `credential_ref`, not raw IDs.
- Summary omits raw host and absolute private output path.
- Manifest rejects wrong schema version, unknown fields, unsafe env names,
  duplicate credentials, missing credential mappings, token values, auth-looking
  values, and private-key markers.
- `201` and `202` succeed.
- `400`, `401`, `404`, and `405` do not retry.
- Retryable failures retry through injected sleeper hooks.
- First terminal failure stops later sends and exits nonzero.
- No raw response-body leakage.
- Redaction scan rejects token values, Authorization headers, raw response
  bodies, DB URLs, private-key markers, bearer/cookie values, and raw payload
  fields in generated files and terminal strings.
- Output rejects symlinks and `docs/evidence`.
- Consumer tracker remains exactly seven targets, all `prepared`.

No benchmark is required. Any timing fields are local engineering diagnostics
only, not SLA, readiness, compliance, production capacity, or AVL reliability
proof.

## Required Verification

```bash
go test ./internal/avladapter ./cmd/avl-vendor-adapter ./cmd/telemetry-ingest ./internal/devices
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

## Phase 48 Closeout Requirements

Phase 48 is not closed until implementation is reviewed by the master agent,
all required checks pass or blockers are documented truthfully, the Phase 48
handoff exists, `docs/handoffs/latest.md` is updated, roadmap/status docs are
consistent, no forbidden claims are introduced, and the consumer tracker remains
unchanged.
