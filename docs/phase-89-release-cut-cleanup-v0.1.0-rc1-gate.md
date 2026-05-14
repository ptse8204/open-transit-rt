# Phase 89 -- Release-Cut Cleanup / v0.1.0-rc.1 Gate

## Purpose

Phase 89 runs a serious local release-candidate gate after the consumer-grade
control-plane work. It is a review phase only unless a maintainer separately
authorizes exact release actions.

This phase must not tag, publish, distribute packages, publish images, create
release artifacts, collect retained evidence, contact external services for
proof, move consumer statuses, or claim release readiness unless the gate
passes and a maintainer separately authorizes the exact release action.

## Current Truth

- Phases 0-88 are complete for their bounded scopes.
- Phase 72 remains `needs_review`, not release-ready.
- Phase 74 CP000008 remains the latest GitHub Pages publication at commit
  `a8b250e`.
- Phase 88 closed with private in-app training, role tours, glossary, recovery
  guidance, quick tasks, handoff checklist, and operator training guide.
- Evidence/adoption/compliance tracks remain optional and require separate
  written authorization.
- Release tag, package publication, image publication, and release-ready claims
  remain unauthorized.
- All seven consumer targets remain exactly `prepared`: Google Maps, Apple
  Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land.

## Sub-Agent Plan

- Context / Repo Truth Sub-Agent -- GPT-5.5 x-high, simulated: confirm current
  phase status, prior Phase 72 blockers, protected paths, package/release
  scripts, and validation targets before edits.
- Planning Sub-Agent -- GPT-5.5 x-high, simulated: keep checkpoints as local
  diagnostics, draft notes, and blockers; do not create a release.
- Implementation Sub-Agent -- GPT-5.5 high, simulated: update only docs,
  bounded diagnostic summaries, release notes draft, and route inventory files
  needed for review.
- QA Sub-Agent -- GPT-5.5 high, simulated: run baseline validation,
  release-candidate diagnostics, connector/backend checks, exact prepared-only
  tracker checks, and protected-path checks.
- UI/UX Sub-Agent -- GPT-5.5 high, simulated: review private UI route coverage
  and route inventory; do not add new UI unless a gate defect requires a
  bounded fix.
- Documentation / IA Sub-Agent -- GPT-5.5 high, simulated: prepare clear
  release notes draft, known blockers matrix, and status handoff.
- Claim-Boundary Sub-Agent -- GPT-5.5 high, simulated: block release-ready,
  production, hosted-service, compliance, consumer, final-root, vendor,
  hardware, SLA/uptime, and ETA-quality claims.
- Security/Auth Sub-Agent -- GPT-5.5 high, simulated: preserve private route,
  auth, CSRF, token, and no-secret boundaries during route checks.
- Data/Migration Sub-Agent -- GPT-5.5 high, simulated for migration review; no
  migration is expected.

## Master Approval Before Implementation

Approved bounded scope:

- Run local product, route, backend, connector, and claim-boundary diagnostics.
- Record exact pass/fail/not-checked status and blockers.
- Prepare draft release notes and a known blockers matrix.
- Record release package and package audit as blocked/not checked unless a
  maintainer separately authorizes package creation.
- Do not tag, publish, distribute packages, publish images, push branches,
  create release artifacts, or claim release readiness.

Required edits before implementation: none.

## Checkpoints

### Checkpoint 000001 -- Plan

Deliver this plan and link it from current source-of-truth docs as the active
Phase 89 plan.

Expected files:

- `docs/phase-89-release-cut-cleanup-v0.1.0-rc1-gate.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/roadmap-status.md`
- `docs/open-transit-rt-master-planner-remaining-work.md`

### Checkpoint 000002 -- Clean-Checkout Local Product Gate

Run local product validation from the current branch state:

- `git status --short`
- `git diff --check`
- `make check`
- `make validate`
- `make test`
- exact seven-target prepared-only tracker check
- protected-path status check

Expected output:

- update a Phase 89 diagnostic/status doc with pass/fail/not-checked rows;
- document dirty working tree state honestly;
- do not create retained evidence.

### Checkpoint 000003 -- Frontend And Accessibility Gate

Run local route/UI diagnostics:

- `RUN_LOCAL_APP=true make release-candidate-check`
- private Operations Console major route inventory review;
- five local public feed path diagnostics from the existing checker;
- accessibility/status-boundary review from tests and route output.

Expected output:

- route inventory;
- local route diagnostic result;
- blockers or residual risks.

### Checkpoint 000004 -- Connector And Backend Diagnostics Gate

Run backend and connector diagnostics:

- `make external-connection-check`
- `make adapter-conformance`
- `make test-connector-examples`
- relevant backend/package tests already covered by `make test`.

Expected output:

- connector/backend diagnostics matrix;
- note that synthetic/local connector checks do not prove real vendor, device,
  hardware, or production operations outcomes.

### Checkpoint 000005 -- RC1 Notes, Package Status, And Blockers Matrix

Prepare review artifacts:

- draft `v0.1.0-rc.1` release notes;
- known blockers matrix;
- release package status rows.

Package boundary:

- `make release-package` and `make audit-release-package` are blocked/not
  checked unless a maintainer separately authorizes package creation.
- Do not create or publish a package as part of Phase 89 without that separate
  authorization.

### Checkpoint 000006 -- Closeout

Write Phase 89 handoff and update current source-of-truth docs. Confirm:

- no protected path writes;
- exact prepared-only consumer tracker;
- no release tag;
- no published package;
- no published image;
- no retained evidence collection;
- no forbidden claims.

## Protected Paths

Do not modify or generate files under:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/**`
- `docs/evidence/consumer-submissions/artifacts/**`
- `docs/evidence/consumer-submissions/packets/**`

## Claim Boundaries

Phase 89 may say:

- local release-candidate diagnostics;
- route inventory;
- draft release notes;
- known blockers matrix;
- package diagnostics blocked/not checked without separate authorization;
- `needs_review` unless every gate passes and release action is separately
  authorized.

Phase 89 must not claim:

- release readiness unless the full gate passes and release action is
  separately authorized;
- CAL-ITP/Caltrans compliance;
- agency adoption or approval;
- consumer submission, review, acceptance, ingestion, listing, or display;
- final-root readiness;
- hosted SaaS;
- paid support;
- SLA or uptime;
- production readiness;
- vendor compatibility;
- hardware certification;
- production-grade ETA quality;
- real-world ETA accuracy;
- public launch completion.

## Validation

Run at least:

```bash
git status --short
git diff --check
make check
make validate
make test
RUN_LOCAL_APP=true make release-candidate-check
make external-connection-check
make adapter-conformance
make test-connector-examples
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
```

Blocked unless separately authorized:

```bash
make release-package
make audit-release-package
```

Optional local config review when safe:

```bash
docker compose -f deploy/docker-compose.yml config
```

## Closeout Report Requirements

The closeout must include:

```text
Phase:
Sub-agents used or simulated, including intended model level:
Goal:
Changed files:
Routes added/changed:
Commands added/changed:
Migrations:
Validation run:
Blocked checks:
Known blockers:
Protected path status:
Consumer tracker status:
Claim-boundary status:
Security/auth status:
Accessibility status:
Docs/site/wiki alignment:
Commit list:
Master review:
Required edits:
Decision:
Next phase:
```
