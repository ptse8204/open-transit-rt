# Phase 48 Handoff — AVL Adapter Runtime Path

## Status

Phase 48 is complete for the private AVL adapter runtime send-mode scope.

## Implemented

- Added mutually exclusive `--dry-run` and `--send` modes to
  `cmd/avl-vendor-adapter`.
- Preserved dry-run no-network stdout telemetry JSON and stderr diagnostics
  JSON behavior.
- Added strict `avl-adapter-send.v1` manifest parsing with unknown-field
  rejection, env-only token references, safe env-var validation, duplicate
  credential rejection, missing credential mapping rejection, and secret-like
  value rejection.
- Added the Phase 48 `AVL_ADAPTER_*` send config contract.
- Added preflight validation for mapping, payload transform, manifest, target
  URL, token env presence, output path, stale/future blockers, and optional
  warning gates before network I/O.
- Added per-record `POST /v1/telemetry` sending with bearer auth, JSON content
  type, `201`/`202` success handling, non-retryable `400`/`401`/`404`/`405`,
  retryable network/timeout/`408`/`429`/`5xx` handling, injectable sleeper, and
  first-terminal-failure stop behavior.
- Added redacted private send outputs:
  `summary.json`, `summary.md`, `manifest.json`, `manifest.md`, and
  `diagnostics.json`.
- Added redaction scans for generated files and terminal output, including
  token values, authorization/bearer/cookie patterns, DB URLs, private-key
  markers, raw response-body markers, and raw vendor payload identifiers.

## Boundaries

Phase 48 preserves `/v1/telemetry` as the only runtime ingest boundary and does
not change `/v1/telemetry` payload or auth semantics. It adds no public API,
admin route, queue, scheduler, daemon, webhook receiver, consumer workflow,
evidence packet, named vendor support, real vendor payload, credential value,
consumer-status change, compliance claim, hosted SaaS claim, production
readiness claim, vendor-compatibility claim, production AVL reliability claim,
or production-grade ETA claim.

## Verification

Master closeout verification passed:

- `go test ./internal/avladapter ./cmd/avl-vendor-adapter`
- `go test ./internal/avladapter ./cmd/avl-vendor-adapter ./cmd/telemetry-ingest ./internal/devices`
- `make validate`
- `make test`
- `git diff --check`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json`
- `docker compose -f deploy/docker-compose.yml config`

The consumer tracker remains unchanged with exactly seven targets, all
`prepared`. No files under `docs/evidence` were edited.

## Next Phase

Phase 49 — External Predictor Runtime Adapter should start with a fresh
read-only planning sub-agent pass before implementation.
