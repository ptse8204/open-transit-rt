# Phase 75 Prompt — Consumer-Grade Control Plane Roadmap Pack

## Goal

Add this roadmap pack to the repo and link it from source-of-truth status docs
only where needed. This is a planning phase only.

## Read first

- `AGENTS.md`
- `README.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-74.md`
- `docs/open-transit-rt-master-planner-remaining-work.md`
- `docs/roadmap-status.md#review-and-recommendations`
- `docs/roadmaps/agency-first-connector-platform/00-CODEX-READ-ME-FIRST.md`
- `docs/roadmaps/agency-first-connector-platform/04-master-subagent-operating-manual.md`
- `docs/roadmaps/agency-first-connector-platform/05-validation-and-claim-boundaries.md`
- `docs/evidence/redaction-policy.md`
- `docs/evidence/consumer-submissions/status.json`

Also read the binding product-contract docs required by `AGENTS.md`:

- `docs/codex-task.md`
- `docs/architecture.md`
- `docs/conversation-summary.md`
- `docs/requirements-2a-2f.md`
- `docs/requirements-trip-updates.md`
- `docs/requirements-calitp-compliance.md`
- `docs/repo-gaps.md`
- `docs/dependencies.md`

## Implementation

Create:

```text
docs/roadmaps/consumer-grade-control-plane/README.md
docs/roadmaps/consumer-grade-control-plane/00-CODEX-READ-ME-FIRST.md
docs/roadmaps/consumer-grade-control-plane/01-roadmap-overview.md
docs/roadmaps/consumer-grade-control-plane/02-phases-and-checkpoints.md
docs/roadmaps/consumer-grade-control-plane/03-master-subagent-operating-manual.md
docs/roadmaps/consumer-grade-control-plane/04-validation-and-claim-boundaries.md
docs/roadmaps/consumer-grade-control-plane/phase-prompts/*.md
docs/roadmaps/consumer-grade-control-plane/audit-prompts/*.md
```

Update only if needed, and use planning-only wording:

- `docs/handoffs/latest.md`
- `docs/current-status.md`
- `docs/roadmap-status.md`
- `docs/open-transit-rt-master-planner-remaining-work.md`
- `README.md`

The updates should say a maintainer-authorized proposed consumer-grade
control-plane roadmap exists. They must not say Phase 76+ is active,
implemented, release-ready, or evidence-backed.

## Forbidden

- No UI implementation.
- No API implementation.
- No migrations.
- No evidence writes.
- No consumer status changes.
- No release tag/package/published image.
- No compliance/adoption/consumer/final-root/SaaS/production/vendor/SLA/ETA claim.

## Validation

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
expected = ["Google Maps", "Apple Maps", "Transit App", "Bing Maps", "Moovit", "Mobility Database", "transit.land"]
data = json.loads(Path("docs/evidence/consumer-submissions/status.json").read_text())
records = data.get("targets", [])
seen = {row["target"]: row.get("status") for row in records}
assert list(seen) == expected, seen
assert all(seen[name] == "prepared" for name in expected), seen
PY
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum
```

## Commits

```text
Phase 75 -- Checkpoint 000001: add consumer-grade control plane roadmap
Phase 75 -- Checkpoint 000002: link consumer-grade roadmap from status docs
Phase 75 -- Checkpoint 000003: close consumer-grade roadmap planning review
```

Phase 75 is closed only after Checkpoint 000003.
