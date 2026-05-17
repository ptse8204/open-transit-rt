# External Connector Runtime Integration Phase Plan

This plan is the durable implementation guide for the focused external
connector runtime roadmap. Each phase must have its own commit and must keep
the evidence and claim boundaries unchanged.

## Grounding

Read before planning or editing:

- `AGENTS.md`
- `README.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/connectors/catalog.md`
- `docs/external-connection-readiness.md`
- `docs/integration-adapter-kit.md`
- `docs/connectors/plugin-contract.md`
- `docs/tutorials/device-avl-integration.md`
- `docs/tutorials/external-adapter-conformance.md`
- `docs/evidence/redaction-policy.md`
- `.github/workflows/test.yml`
- `.github/workflows/release-gates.yml`

Inspect current implementation under:

- `cmd/adapter-conformance/`
- `cmd/avl-vendor-adapter/`
- `cmd/agency-config/`
- `cmd/telemetry-ingest/`
- `examples/connectors/`
- `internal/admincontrol/`
- `internal/avladapter/`
- `internal/connectors/`
- `internal/prediction/`
- `internal/telemetry/`
- `testdata/adapter-conformance/`
- `testdata/avl-vendor/`

Before UI, website, README visual structure, tutorial-flow, docs-UX, or
product-flow changes, load:

```text
/Users/edwintse/.agents/skills/web-design-engineer/SKILL.md
```

## Hard Boundaries

Do not:

- write to protected evidence paths;
- change `docs/evidence/consumer-submissions/status.json`;
- contact external agencies, vendors, consumers, validators, portals, or live
  connector services;
- use real private agency data;
- commit credentials, tokens, private payloads, webhook URLs, private endpoint
  URLs, database URLs, or private paths;
- expose arbitrary command execution in the browser;
- claim vendor compatibility, hardware certification, CAL-ITP/Caltrans
  compliance, production readiness, consumer acceptance, final-root readiness,
  hosted service availability, SLA/uptime, production AVL reliability,
  production-grade ETA quality, or real-world ETA accuracy.

All connector examples, tests, and conformance fixtures should be local,
synthetic, redacted, and no-contact by default.

## Baseline Validation

Run the narrow validation for each phase plus the phase-specific commands
below. When code changes touch shared behavior, broaden to the full set.

Common baseline:

```bash
git diff --check
go test ./...
make check
make test
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```

Connector baseline:

```bash
make external-connection-check
make adapter-conformance
make test-connector-examples
make gtfsrt-conformance
```

Do not run live external validators, portals, vendor endpoints, or consumer
systems as part of this roadmap.

## Commit Format

Use one commit per phase:

```text
External Connector Runtime -- Phase 01: baseline connector runtime boundaries
External Connector Runtime -- Phase 02: improve telemetry adapter runtimes
External Connector Runtime -- Phase 03: harden telemetry and Vehicle Positions
External Connector Runtime -- Phase 04: implement predictor shadow diagnostics
External Connector Runtime -- Phase 05: improve monitoring export runtime
External Connector Runtime -- Phase 06: expand connector conformance fixtures
External Connector Runtime -- Phase 07: improve browser connector setup and health
External Connector Runtime -- Phase 08: close connector runtime integration roadmap
```

## Phase 01 - Runtime Boundary Baseline

Goal: Reconfirm the current connector runtime boundaries and freeze the first
implementation contract before adding new runtime behavior.

Required work:

- Inventory all current connector examples, conformance fixture types,
  Make targets, private Connector Workbench routes, telemetry ingest contracts,
  prediction adapter contracts, and monitoring/export helpers.
- Identify which connector work is runtime behavior, which is local dry-run
  tooling, and which remains documentation-only.
- Document deployment-owned config ownership for endpoints, tokens,
  destination URLs, retry settings, and retention.
- Confirm no runtime path requires dynamic backend plugin loading.
- Add or update tests that catch unsafe manifest fields, private endpoint
  leakage, and accidental evidence writes if gaps exist.
- Keep named vendor and named predictor compatibility claims out of public
  wording.

Acceptance:

- A maintainer can see the runtime boundary, default disabled/no-send state,
  and first safe check for each connector family.
- No browser page can run arbitrary commands.
- Protected evidence paths and consumer tracker status are unchanged.

Validation:

```bash
git diff --check
go test ./...
make check
make external-connection-check
make adapter-conformance
make test-connector-examples
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```

## Phase 02 - Vehicle / GPS / AVL Runtime Adapters

Goal: Make vehicle telemetry connector runtimes more useful while preserving
authenticated ingest and fail-closed behavior.

Required work:

- Improve the CSV replay adapter so a technical helper can run a repeatable
  local replay into dry-run output and, only with explicit deployment-owned
  config, authenticated `POST /v1/telemetry`.
- Improve the HTTP polling adapter around bounded polling config, redacted
  diagnostics, timeout handling, stale/future timestamp rejection, and
  no-send dry-run defaults.
- Improve the webhook sidecar adapter shape around fixed request parsing,
  bounded body size, auth expectations, duplicate/out-of-order review, and
  deployment-owned send config.
- Improve the generic JSON transform adapter with clearer mapping validation,
  required-field errors, coordinate/time validation, and redacted diagnostics.
- Keep all examples synthetic and avoid named vendor compatibility wording.

Acceptance:

- Each AVL adapter family has a safe dry-run path and an explicit
  deployment-owned path toward `/v1/telemetry`.
- Malformed, stale, future-dated, wrong-agency, unknown-device, low-quality,
  duplicate, out-of-order, missing-field, and invalid-coordinate cases remain
  fail-closed or explicitly diagnostic.
- Tests cover adapter behavior without contacting external systems.

Validation:

```bash
git diff --check
go test ./cmd/avl-vendor-adapter ./internal/avladapter ./examples/connectors/...
make external-connection-check
make adapter-conformance
make test-connector-examples
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```

## Phase 03 - Authenticated Telemetry And Vehicle Positions Hardening

Goal: Strengthen the authenticated telemetry boundary that runtime adapters
call and keep Vehicle Positions safe when connector inputs are incomplete.

Required work:

- Recheck `/v1/telemetry` request validation, auth failures, device-token
  failure visibility, agency-scope rejection, stale/future handling, and
  duplicate/out-of-order diagnostics.
- Improve operator diagnostics without exposing tokens, raw private payloads,
  auth headers, private paths, or database details.
- Make accepted and rejected observations easier to review in private browser
  surfaces where the data is already safe to display.
- Preserve conservative matching: unknown is better than false certainty.
- Ensure Vehicle Positions continue independently of external predictor
  availability.
- Recheck Vehicle Positions trip descriptor suppression for stale telemetry,
  low-confidence assignment, missing active GTFS, and ambiguous matching.

Acceptance:

- Runtime adapters have a clear, authenticated server boundary to call.
- Operators can diagnose common ingest failures without command-line access to
  raw secrets or private payloads.
- Public feed behavior remains conservative when telemetry is unavailable or
  low confidence.
- Vehicle Positions do not depend on an external predictor or external
  prediction sidecar health.

Validation:

```bash
git diff --check
go test ./cmd/telemetry-ingest ./internal/telemetry ./internal/devices ./internal/feed ./internal/state
make check
make test
make smoke
make external-connection-check
make adapter-conformance
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```

## Phase 04 - External HTTP Predictor Shadow Mode

Goal: Make external predictor evaluation repeatable without mutating public
Trip Updates by default.

Required work:

- Recheck and improve generic external HTTP predictor configuration, health
  checks, timeout behavior, byte caps, redirect rejection, malformed output
  handling, stale output handling, wrong-agency/wrong-feed rejection, and
  fail-closed fallback.
- Improve shadow-mode comparison diagnostics so maintainers can compare
  deterministic and external outputs without raw telemetry or private
  assignment details.
- Keep Vehicle Positions independent when a predictor is missing, slow,
  malformed, or unavailable.
- Do not make production-grade ETA quality, named predictor compatibility, or
  real-world ETA accuracy claims.

Acceptance:

- Shadow mode can be tested locally with synthetic fixtures and a local stub.
- Public Trip Updates remain unchanged unless a later deployment-owned config
  explicitly enables a supported path.
- Failure cases are bounded, redacted, and fail closed.

Validation:

```bash
git diff --check
go test ./internal/prediction ./internal/feed/tripupdates ./examples/connectors/predictor-sidecar-stub
make check
make adapter-conformance
make gtfsrt-conformance
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```

## Phase 05 - Monitoring / Export Helper Runtime

Goal: Improve monitoring/export as a deployment-owned runtime helper while
keeping no-send defaults and redacted summaries.

Required work:

- Define the summary export format for health digests, connector status,
  validation posture, telemetry freshness, and private operations summaries.
- Keep destination config deployment-owned and off by default.
- Add or improve dry-run output for file, stdout, or local fixture
  destinations before any live delivery path is considered.
- Redact endpoint URLs, tokens, destination IDs, raw private payloads, private
  paths, and private infrastructure details.
- Preserve no SLA/uptime, hosted-service, paid-support, or production
  readiness claims.

Acceptance:

- Operators can generate a redacted local summary export.
- No notification or monitoring destination is contacted by default.
- The helper does not write retained evidence or mutate consumer status.

Validation:

```bash
git diff --check
go test ./examples/connectors/monitoring-export ./examples/connectors/sdk/monitoring
make check
make external-connection-check
make adapter-conformance
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```

## Phase 06 - Connector Conformance Expansion

Goal: Make adapter conformance catch realistic local failure modes before a
connector reaches runtime evaluation.

Required work:

- Expand telemetry fixture coverage where runtime work adds new branches.
- Expand prediction fixture coverage for health checks, timeout/malformed
  output, stale output, wrong agency/feed, missing Vehicle Positions reference,
  and public mutation attempts.
- Expand monitoring/export fixtures for no-send defaults, redaction, summary
  format, and blocked unredacted destinations.
- Keep validator and consumer/discovery fixtures no-contact, no-portal, and
  status-preserving.
- Make first safe checks explicit for every fixture family.

Acceptance:

- `make adapter-conformance` remains offline and deterministic.
- Added fixtures increase runtime confidence without contacting external
  systems or claiming external proof.
- The suite fails clearly when a connector attempts unsafe sends, status
  mutation, arbitrary commands, or unredacted private output.

Validation:

```bash
git diff --check
go test ./cmd/adapter-conformance ./internal/connectors
make check
make external-connection-check
make adapter-conformance
make test-connector-examples
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```

## Phase 07 - Browser Connector Setup And Health

Goal: Make connector setup and health review clear from the private browser UI
after a technical helper starts the app.

Required work:

- Improve `/admin/operations/connectors`,
  `/admin/operations/connectors/workbench`, and
  `/admin/operations/connectors/tests` around connector setup guidance,
  connector health review, first safe checks, and runtime boundaries.
- Add copyable deployment-owned config checklists for AVL, prediction, and
  monitoring/export connectors without rendering secrets or private
  destinations.
- Add warnings for secrets, private payloads, private endpoint URLs, and
  private paths.
- Keep fixed commands as operator-shell guidance only; do not execute them
  from the browser.
- Keep wording focused on local/self-hosted evaluation and synthetic/local
  conformance checks.

Acceptance:

- Agency staff can review connector status and next actions from the browser.
- Technical helpers can copy safe config checklists without exposing secrets.
- The UI does not present vendor compatibility, production readiness, SLA,
  compliance, or ETA-quality claims.

Validation:

```bash
git diff --check
go test ./cmd/agency-config ./internal/admincontrol ./internal/connectors
make check
make test
make smoke
make external-connection-check
make adapter-conformance
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```

## Phase 08 - Runtime Roadmap Closeout And Release Gate

Goal: Close the runtime integration roadmap with validation, status updates,
and a concrete next product-quality recommendation.

Required work:

- Re-run connector-specific validations and document what passed, what was
  skipped, and why.
- Update human docs, connector catalog, integration adapter kit, current
  status, and latest handoff as needed.
- Confirm protected evidence paths are untouched.
- Confirm the consumer tracker remains exactly seven prepared-only targets.
- Confirm unsupported claims remain unsupported.
- Decide whether the next safe step is another runtime hardening slice, a
  release-candidate gate, or a separately authorized evidence track.

Acceptance:

- The roadmap closeout is truthful, actionable, and narrow.
- All runtime work remains local/synthetic or deployment-owned by default.
- No external proof or claim has been implied.

Validation:

```bash
git diff --check
go test ./...
make check
make test
make smoke
make external-connection-check
make adapter-conformance
make test-connector-examples
make gtfsrt-conformance
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```
