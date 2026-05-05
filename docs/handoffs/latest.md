# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 30 — Consumer Submission Execution is closed as Outcome B — blocker-documented closure only.

No authorized submission, official-path verification evidence, or target-originated artifact was available.

Phases 0 through 30 are closed for their documented scopes. Track A is also closed for its docs-only external-proof workflow scope. Do not reopen earlier phases unless a blocking truthfulness, safety, security, realtime-quality, evidence, agency-boundary, auth, data-isolation, agency-domain, device/AVL onboarding, admin-UX, operations-hardening, or submission-readiness issue directly requires it.

The recommended next implementation phase is Phase 31 — Agency Pilot Program Package. Phase 31 must proceed from the prepared-only consumer state and must not assume submission, review, acceptance, rejection, blocker, ingestion, listing, display, or adoption evidence exists.

## Phase 29B Summary

- Added `internal/avladapter` for strict synthetic vendor payload and mapping transforms into the existing `telemetry.Event` contract.
- Added dry-run-only `cmd/avl-vendor-adapter` with required `--dry-run`, optional `--reference-time`, telemetry JSON array on stdout, and diagnostics JSON array on stderr.
- Added synthetic fixtures under `testdata/avl-vendor/` for valid, source mismatch, duplicate mapping, empty mapped IDs, missing/invalid coordinate, stale/future timestamps, unknown vendor vehicle, low GPS accuracy, mixed batch, duplicate/out-of-order dry-run observations, optional trip hint, and malformed payload.
- Added focused adapter and CLI tests.
- Updated Phase 29B docs, device/AVL tutorial guidance, evidence template/discoverability, decisions, dependencies, and current status.

## Phase 30 Outcome B Summary

- Phase 30 closed as Outcome B — blocker-documented closure only.
- No authorized submission, official-path verification evidence, or target-originated artifact was available.
- No Phase 30 target was selected.
- Target selection is deferred until an operator is authorized and either official-path verification or target-originated evidence can be retained.
- No individual target status changed to `blocked` because no target-specific blocker artifact exists.
- `docs/evidence/consumer-submissions/status.json` and all current target records under `docs/evidence/consumer-submissions/current/` were left unchanged.
- Artifact directories remain README-only; no receipts, screenshots, tickets, correspondence, blocker notes, or placeholder artifacts were added.
- Mobility Database and transit.land may be considered as future candidate suggestions once authorized, but they were not selected in Phase 30.

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
3. `docs/handoffs/phase-30.md`
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

Start Phase 31 from the prepared-only consumer state. Phase 31 must not assume consumer or aggregator submission, review, acceptance, rejection, blocker, ingestion, listing, display, or adoption evidence exists.

Consumer or aggregator submission work remains available only when a future operator is authorized, a target is selected, official target paths are verified, and target-originated evidence can be retained and redacted. Product improvements, validator success, or prepared packets alone must not advance target statuses.

## Exact First Commands

```bash
make validate
make test
git diff --check
```

Run these when Phase 31 work touches relevant surfaces:

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

## Checks Run For Phase 30

- Pre-edit `make validate` — passed.
- Pre-edit `make test` — passed.
- Pre-edit `git diff --check` — passed.
- Post-edit `make validate` — passed.
- Post-edit `python3 -m json.tool docs/evidence/consumer-submissions/status.json` — passed.
- Post-edit tracker/status consistency check — passed; all seven targets remain `prepared`.
- Post-edit `make test` — passed.
- Post-edit `make realtime-quality` — passed.
- Post-edit `make smoke` — passed.
- Post-edit `make test-integration` — passed.
- Post-edit `docker compose -f deploy/docker-compose.yml config` — passed.
- Post-edit `git diff --check` — passed.
- Post-edit targeted artifact/tracker scans — passed.
- Post-edit context-aware forbidden-claim and redaction-sensitive term scans — reviewed; matches are allowed negative/boundary/security contexts.
- Blocked commands: none.

## Current Evidence And Security Boundary

- The OCI pilot packet at `docs/evidence/captured/oci-pilot/2026-04-24/` remains the current hosted/operator evidence packet.
- Phase 23 did not create final-root evidence. No agency-owned or agency-approved final public feed root is available in repo evidence.
- Phase 24 real-agency GTFS evidence scaffolding is template-only until real agency-approved, public-safe evidence exists.
- Phase 25 device/AVL evidence scaffolding is template-only until real public-safe device or AVL integration evidence exists.
- Phase 29B synthetic fixtures are not real vendor AVL evidence.
- Phase 20 prepared packets are operator review artifacts only; they are not submissions.
- Phase 30 did not select a target, verify an official path, submit a packet, add artifacts, or change consumer statuses.
- Consumer-ingestion workflow records and docs tracker records are not third-party acceptance unless retained evidence from the named target exists.
- Do not rely on old local `.cache` credentials.
- Do not commit secrets, generated tokens, private keys, ACME material, admin tokens, device tokens, JWT secrets, CSRF secrets, DB passwords, webhook URLs, notification credentials, raw telemetry payloads, unredacted correspondence, private portal credentials, private ticket links, raw logs with credentials, private backup paths, or raw private operator artifacts.

## First Files Likely To Edit For Phase 31

- `docs/phase-31-agency-pilot-program-package.md`
- `docs/handoffs/latest.md`
- `docs/current-status.md`

Do not edit target-specific consumer records, `docs/evidence/consumer-submissions/status.json`, or artifact directories unless retained, redacted, target-originated evidence supports a target-specific status transition.

## Constraints To Preserve

- Keep Trip Updates pluggable and Vehicle Positions first.
- Preserve admin auth, role checks, CSRF behavior, and token/secret handling.
- Do not expose admin/debug/JSON surfaces on the production public edge.
- Do not add consumer submission APIs unless explicitly approved and backed by current target documentation.
- Do not automate submissions, contact external portals, guess submission paths, or invent acceptance/rejection/compliance evidence.
- Keep `prepared` conditional on packet completeness.
- Do not describe Open Transit RT as hosted SaaS, paid support, SLA-backed, agency-endorsed, marketplace/vendor equivalent, universally production ready, production multi-tenant hosted, production-grade ETA proven, real-world ETA-accuracy proven, certified hardware supported, or vendor-compatible.

## Exact Next-Step Recommendation

Start Phase 31 — Agency Pilot Program Package from the prepared-only consumer state. Do not assume submission, review, acceptance, rejection, blocker, ingestion, listing, display, or adoption evidence exists.
