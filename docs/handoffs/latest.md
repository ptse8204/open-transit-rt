# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 29B — AVL / Vendor Adapter Pilot Implementation is complete for the synthetic, dry-run-only adapter pilot scope.

Phases 0 through 29B are closed for their documented scopes. Track A is also closed for its docs-only external-proof workflow scope. Do not reopen earlier phases unless a blocking truthfulness, safety, security, realtime-quality, evidence, agency-boundary, auth, data-isolation, agency-domain, device/AVL onboarding, admin-UX, operations-hardening, or submission-readiness issue directly requires it.

The recommended next implementation phase is Phase 30 — Consumer Submission Execution. Phase 30 must not advance consumer or aggregator statuses without retained, redacted, target-originated evidence.

## Phase 29B Summary

- Added `internal/avladapter` for strict synthetic vendor payload and mapping transforms into the existing `telemetry.Event` contract.
- Added dry-run-only `cmd/avl-vendor-adapter` with required `--dry-run`, optional `--reference-time`, telemetry JSON array on stdout, and diagnostics JSON array on stderr.
- Added synthetic fixtures under `testdata/avl-vendor/` for valid, source mismatch, duplicate mapping, empty mapped IDs, missing/invalid coordinate, stale/future timestamps, unknown vendor vehicle, low GPS accuracy, mixed batch, duplicate/out-of-order dry-run observations, optional trip hint, and malformed payload.
- Added focused adapter and CLI tests.
- Updated Phase 29B docs, device/AVL tutorial guidance, evidence template/discoverability, decisions, dependencies, and current status.

## Truthfulness And Evidence Boundary

- Phase 29B is synthetic adapter pattern and dry-run transform evidence only.
- The mapping file is the authority for Open Transit RT `agency_id`, `device_id`, and `vehicle_id`; vendor payload IDs are lookup keys only.
- Partial stdout from a nonzero dry run is transform output only, not submitted telemetry, production integration evidence, successful vendor compatibility proof, or database ingest status.
- Duplicate/out-of-order adapter diagnostics are batch-level dry-run observations only, not telemetry ingest acceptance statuses.
- `trip_hint` is a hint only, not assignment proof, ETA proof, consumer-facing correctness proof, or evidence that the vehicle is matched.
- No network send mode exists.
- No real vendor AVL data, credentials, endpoint URLs, tokens, private identifiers, named vendor dependency, or runtime vendor integration was added.
- No public feed URLs changed.
- No GTFS-RT protobuf contracts changed.
- No consumer statuses changed.
- No auth boundaries changed.
- No database schema changed.
- No Trip Updates adapter behavior changed.

Do not claim hosted SaaS availability, paid support/SLA coverage, universal production readiness, production multi-tenant hosting, consumer acceptance, CAL-ITP/Caltrans compliance, agency endorsement, marketplace/vendor equivalence, real-world ETA accuracy, production-grade ETA quality, certified hardware support, vendor compatibility, or production AVL reliability.

All seven consumer and aggregator targets remain `prepared` only. No target has submitted, under-review, accepted, rejected, or blocked evidence.

The OCI pilot DuckDNS hostname remains pilot evidence, not agency-owned stable URL/domain proof.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-29b.md`
4. `docs/phase-30-consumer-submission-execution.md`
5. `docs/evidence/consumer-submissions/submission-workflow.md`
6. `docs/evidence/consumer-submissions/status.json`
7. `docs/evidence/redaction-policy.md`
8. `SECURITY.md`
9. `docs/california-readiness-summary.md`
10. `docs/compliance-evidence-checklist.md`
11. `README.md`
12. `docs/dependencies.md`
13. `docs/decisions.md`

## Current Objective

Start Phase 30 only when maintainers are ready and authorized. Phase 30 should execute consumer or aggregator submission workflows only when official target paths are verified and target-originated evidence can be retained and redacted. Product improvements, validator success, or prepared packets alone must not advance target statuses.

## Exact First Commands

```bash
make validate
make test
git diff --check
```

Run these when Phase 30 work touches relevant surfaces:

```bash
make realtime-quality
make smoke
make test-integration
docker compose -f deploy/docker-compose.yml config
```

## Checks Run For Phase 29B

- Focused `go test ./internal/avladapter ./cmd/avl-vendor-adapter` — passed.
- `go run ./cmd/avl-vendor-adapter help` — passed.
- `go run ./cmd/avl-vendor-adapter --dry-run --reference-time 2026-05-04T12:00:00Z --mapping testdata/avl-vendor/mapping.json testdata/avl-vendor/valid.json` — passed.
- `go run ./cmd/avl-vendor-adapter --mapping testdata/avl-vendor/mapping.json testdata/avl-vendor/valid.json` — failed as expected with send-mode-not-implemented wording.
- Targeted `testdata/avl-vendor` secret-like fixture scan — passed with no matches.
- Broader Phase 29B docs/evidence scan — reviewed; matches were redaction rules and negative claim-boundary wording only.
- `make validate` — passed.
- `make test` — passed.
- `make realtime-quality` — passed.
- `make smoke` — passed.
- `make test-integration` — passed.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- `git diff --check` — passed.
- Blocked checks: none.

## Current Evidence And Security Boundary

- The OCI pilot packet at `docs/evidence/captured/oci-pilot/2026-04-24/` remains the current hosted/operator evidence packet.
- Phase 23 did not create final-root evidence. No agency-owned or agency-approved final public feed root is available in repo evidence.
- Phase 24 real-agency GTFS evidence scaffolding is template-only until real agency-approved, public-safe evidence exists.
- Phase 25 device/AVL evidence scaffolding is template-only until real public-safe device or AVL integration evidence exists.
- Phase 29B synthetic fixtures are not real vendor AVL evidence.
- Phase 20 prepared packets are operator review artifacts only; they are not submissions.
- Consumer-ingestion workflow records and docs tracker records are not third-party acceptance unless retained evidence from the named target exists.
- Do not rely on old local `.cache` credentials.
- Do not commit secrets, generated tokens, private keys, ACME material, admin tokens, device tokens, JWT secrets, CSRF secrets, DB passwords, webhook URLs, notification credentials, raw telemetry payloads, unredacted correspondence, private portal credentials, private ticket links, raw logs with credentials, private backup paths, or raw private operator artifacts.

## First Files Likely To Edit For Phase 30

- `docs/phase-30-consumer-submission-execution.md`
- `docs/evidence/consumer-submissions/submission-workflow.md`
- target-specific files under `docs/evidence/consumer-submissions/current/`
- target-specific packet or artifact directories under `docs/evidence/consumer-submissions/`
- `docs/handoffs/phase-30.md`
- `docs/handoffs/latest.md`
- `docs/current-status.md`

## Constraints To Preserve

- Keep Trip Updates pluggable and Vehicle Positions first.
- Preserve admin auth, role checks, CSRF behavior, and token/secret handling.
- Do not expose admin/debug/JSON surfaces on the production public edge.
- Do not add consumer submission APIs unless explicitly approved and backed by current target documentation.
- Do not automate submissions, contact external portals, guess submission paths, or invent acceptance/rejection/compliance evidence.
- Keep `prepared` conditional on packet completeness.
- Do not describe Open Transit RT as hosted SaaS, paid support, SLA-backed, agency-endorsed, marketplace/vendor equivalent, universally production ready, production multi-tenant hosted, production-grade ETA proven, real-world ETA-accuracy proven, certified hardware supported, or vendor-compatible.

## Exact Next-Step Recommendation

Start Phase 30 — Consumer Submission Execution only when maintainers have authorized target-specific submission work. Do not advance any target beyond `prepared` without retained, redacted, target-originated evidence.
