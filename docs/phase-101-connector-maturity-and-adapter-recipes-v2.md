# Phase 101 -- Connector Maturity And Adapter Recipes V2

## Scope

Phase 101 improves synthetic/local connector recipes, manifest safety review,
and offline adapter conformance coverage. The phase extends the existing
private Connector Workbench, Connector Tests page, manifest contract, example
sidecars, and conformance suite without changing production connector runtime
behavior.

Allowed work:

- expand operator-facing recipes for CSV vehicle locations, GPS polling APIs,
  AVL sources that can POST, synthetic-only telemetry, external prediction
  sidecars, monitoring/export summaries, and off-host validation;
- add redaction-first recipe templates and a clearer connection decision tree;
- harden safe connector manifest linting for local/example manifests;
- expand adapter conformance cases using committed synthetic fixtures only;
- update focused tests, examples, tutorials, and connector docs.

Not allowed:

- real vendor compatibility, hardware certification, or production AVL
  reliability claims;
- external network contact, vendor portal automation, consumer portal
  automation, or notification send behavior;
- real agency/vendor/device payloads, credentials, private URLs, or retained
  evidence collection;
- protected evidence path writes or consumer tracker status changes;
- public feed contract changes, telemetry ingest contract changes, Trip
  Updates hard-coupling, or migration work unless a safety issue requires
  stopping and re-planning;
- release, tag, package publication, GitHub Release creation, compliance,
  consumer acceptance, production readiness, hosted-service, SLA, or
  production-grade ETA claims.

## Existing Surfaces

Connector maturity work already has bounded surfaces:

- `cmd/agency-config/operations_connector_workbench.go` builds the private
  Connector Workbench and JSON view.
- `cmd/agency-config/operations_connector_tests.go` builds private fixed
  connector test instructions.
- `internal/connectors/manifest.go` validates `open-transit-rt.connector.v1`
  manifests, rejects secrets/private endpoints/raw commands/unsafe claims, and
  loads committed example registry data.
- `cmd/adapter-conformance` validates offline synthetic conformance suites.
- `examples/connectors/**` contains synthetic local examples for CSV replay,
  HTTP polling, prediction sidecar, validator allowlist, and monitoring
  export.
- `testdata/adapter-conformance/**` and `testdata/connectors/**` contain
  committed synthetic fixtures and manifest validation cases.
- `docs/connectors/plugin-contract.md`,
  `docs/tutorials/external-adapter-conformance.md`,
  `docs/tutorials/device-avl-integration.md`, and `wiki/connector-cookbook.md`
  describe the current connector boundaries.

Phase 101 should improve these existing seams. It should not introduce dynamic
plugin loading, browser command execution, live connector execution, or new
credential storage.

## Master-Approved Plan

1. Add this plan and record the Phase 101 claim/security/data boundaries.
2. Implement a bounded V2 improvement across connector recipe guidance,
   manifest linting, synthetic fixtures, and tests. Prefer adding
   machine-checkable local metadata over broad runtime rewrites.
3. Run connector-specific validation plus required baseline/heavier checks;
   patch changed-code failures only.
4. Close Phase 101 with a handoff, status-doc updates, protected-path status,
   prepared-only consumer tracker confirmation, and exact blockers.

The Master Agent approves implementation only if it remains private/offline,
synthetic/local, no-send, no-evidence, no-claim, and compatible with the
existing sidecar/manifest/conformance model.

## Sub-Agent Plan

| Role | Intended model | Use in Phase 101 |
| --- | --- | --- |
| Context / Repo Truth Sub-Agent | GPT-5.5 x-high | Read-only inspection of connector workbench, manifest validation, conformance suite, examples, fixtures, docs, and Makefile gates. |
| Planning Sub-Agent | GPT-5.5 x-high | Read-only checkpoint plan, implementation options, validation plan, and guardrail review. |
| Implementation Sub-Agent | GPT-5.5 high | Simulated by Master unless a bounded disjoint edit becomes useful. |
| QA Sub-Agent | GPT-5.5 high | Simulated by Master through focused connector tests and full required validation. |
| UI/UX Sub-Agent | GPT-5.5 high | Simulated by Master for private Connector Workbench decision-tree and recipe wording. |
| Documentation / IA Sub-Agent | GPT-5.5 high | Simulated by Master for phase docs, tutorials, wiki/status, and handoff. |
| Claim-Boundary Sub-Agent | GPT-5.5 high | Simulated by Master with connector-specific forbidden-claim review. |
| Security/Auth Sub-Agent | GPT-5.5 high | Simulated by Master; preserve private route gates, no browser sends, no command execution, and no secret rendering. |
| Data/Migration Sub-Agent | GPT-5.5 high | Simulated by Master because no persistence or migration is planned. Stop before adding persistence. |
| Release/Supply-Chain Sub-Agent | GPT-5.5 high | Not active unless release/package behavior unexpectedly changes; Phase 101 is not a release phase. |

## Checkpoints

```text
Phase 101 -- Checkpoint 000001: add connector maturity and adapter recipes v2 plan
Phase 101 -- Checkpoint 000002: implement primary scoped work
Phase 101 -- Checkpoint 000003: run validation and patch required gaps
Phase 101 -- Checkpoint 000004: close connector maturity and adapter recipes v2 review
```

## Validation Plan

Focused checks:

```bash
go test ./internal/connectors ./cmd/adapter-conformance ./cmd/agency-config
make external-connection-check
make adapter-conformance
make test-connector-examples
```

Phase closeout baseline:

```bash
git status --short
git diff --check
make check
make audit-product-acceptance
make audit-final-claim-review
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
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum
make validate
make test
docker compose -f deploy/docker-compose.yml config
```

Release-candidate and package checks are not required for Phase 101 unless the
implementation unexpectedly changes release tooling.

## Checkpoint 000001 Report

Checkpoint: Phase 101 -- Checkpoint 000001: add connector maturity and adapter recipes v2 plan

Sub-agents used or simulated, including intended model level: Context / Repo Truth Sub-Agent GPT-5.5 x-high and Planning Sub-Agent GPT-5.5 x-high were launched read-only. Implementation, QA, UI/UX, Documentation / IA, Claim-Boundary, Security/Auth, Data/Migration, and Release/Supply-Chain are simulated by the Master Agent for this planning checkpoint.

Changed files: `docs/phase-101-connector-maturity-and-adapter-recipes-v2.md`

Validation run: `git status --short` before edits; source/doc read of Phase 101 prompt, Connector Workbench, Connector Tests page, connector manifest validator/tests, adapter conformance CLI, external connection check script, connector docs/tutorials, examples, fixtures, current status, latest handoff, roadmap status, and post-90 validation/operating manuals.

Blocked checks: Full validation deferred until implementation and closeout checkpoints.

Protected path status: No protected evidence path edits planned or made.

Consumer tracker status: Must remain exactly seven prepared-only targets; no status edits planned or made.

Claim-boundary status: Plan explicitly forbids real vendor compatibility, hardware certification, production AVL reliability, compliance, consumer acceptance, production readiness, hosted-service, SLA, production-grade ETA, release-readiness, public-launch, and evidence claims.

Security/auth status: Plan preserves private authenticated Connector Workbench/Test routes, no browser command execution, no browser network sends, no secret rendering, no live endpoint storage, and no external contact.

Data/migration status: No migration, telemetry ingest contract change, public feed contract change, Trip Updates hard-coupling, or durable connector runtime state planned.

Master review: Approved to proceed under private/offline, synthetic/local, no-send, no-evidence, no-claim constraints.

Required edits: Implement bounded V2 connector recipe, manifest linting, conformance fixture, and documentation improvements with tests.

Decision: Continue to Checkpoint 000002.

Next checkpoint: Phase 101 -- Checkpoint 000002: implement primary scoped work
