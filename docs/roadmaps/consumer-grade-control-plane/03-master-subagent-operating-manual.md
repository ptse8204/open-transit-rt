# Master/Sub-Agent Operating Manual — Consumer-Grade Control Plane

## Model assignment

Use this exact operating model for future Codex work:

| Role | Model level |
| --- | --- |
| Master Agent | GPT-5.5 x-high |
| Context / Repo Truth Sub-Agent | GPT-5.5 x-high |
| Planning Sub-Agent | GPT-5.5 x-high |
| Implementation Sub-Agent | GPT-5.5 high |
| QA Sub-Agent | GPT-5.5 high |
| UI/UX Sub-Agent | GPT-5.5 high |
| Documentation / IA Sub-Agent | GPT-5.5 high |
| Claim-Boundary Sub-Agent | GPT-5.5 high |

If Codex can spawn real sub-agents, assign those model levels. If not, simulate them in clearly labeled sections.

## Master Agent responsibilities

The Master Agent owns:

- source-of-truth verification;
- phase selection;
- checkpoint sequencing;
- plan approval;
- implementation approval;
- protected-path enforcement;
- consumer tracker enforcement;
- claim-boundary enforcement;
- final closeout commit acceptance.

The Master Agent must not allow implementation until all sub-agent reviews have no required edits.

## Context / Repo Truth Sub-Agent

Read and report:

- current branch and working tree status;
- current phase/handoff truth;
- protected path status;
- consumer tracker exact status;
- relevant routes and APIs;
- relevant docs and tests;
- current validation commands;
- current UI limitations.

Minimum reading:

```text
AGENTS.md
README.md
docs/current-status.md
docs/handoffs/latest.md
docs/handoffs/phase-74.md
docs/open-transit-rt-master-planner-remaining-work.md
docs/roadmap-status.md
docs/roadmaps/agency-first-connector-platform/00-CODEX-READ-ME-FIRST.md
docs/roadmaps/agency-first-connector-platform/04-master-subagent-operating-manual.md
docs/roadmaps/agency-first-connector-platform/05-validation-and-claim-boundaries.md
docs/evidence/redaction-policy.md
docs/evidence/consumer-submissions/status.json
```

Also read the binding product-contract docs required by `AGENTS.md`:

```text
docs/codex-task.md
docs/architecture.md
docs/conversation-summary.md
docs/requirements-2a-2f.md
docs/requirements-trip-updates.md
docs/requirements-calitp-compliance.md
docs/repo-gaps.md
docs/dependencies.md
```

## Planning Sub-Agent

For each checkpoint, produce:

```text
Checkpoint name:
Goal:
Files to read:
Files likely to edit:
Files forbidden to edit:
Route/API changes:
Data/schema changes:
UI/UX acceptance criteria:
QA commands:
Protected path checks:
Consumer tracker check:
Claim-boundary checks:
Stop conditions:
Expected commit message:
```

## Implementation Sub-Agent

Rules:

- Execute only approved checkpoint scope.
- No extra phase scope.
- No evidence writes.
- No consumer status changes.
- No raw secrets, raw tokens, private hostnames, private file paths, or raw external payloads in docs or UI.
- No dynamic backend plugin loading unless separately planned and security-reviewed.
- Preserve auth, role checks, CSRF behavior, request-size caps, and server-owned command mappings.
- Prefer progressive UI on existing private routes over speculative SPA rewrites.
- Classify every browser-controlled workflow using the safe command ladder in
  `01-roadmap-overview.md`.
- For risky maintenance actions, default to browser guidance and
  technical-helper handoff unless a separate phase authorizes browser
  execution.

## UI/UX Sub-Agent

Review questions:

- Can a small-agency operator understand the page without phase history?
- Is there one clear primary action?
- Are blockers explained in plain language?
- Are empty states useful?
- Are labels consistent with `Agency Operations Cockpit / Start Here`?
- Does the UI say what it does not prove?
- Is keyboard navigation clear?
- Is mobile/tablet layout usable enough for review?
- Are dangerous actions confirmed and reversible where possible?
- Does confirmation copy name the exact object, agency scope, public/private
  impact, public-feed impact, rollback path, and claim boundary?

## Documentation / IA Sub-Agent

Review:

- README alignment;
- docs navigation;
- wiki navigation;
- route names and labels;
- tutorial alignment;
- in-app help links;
- screenshot captions;
- public site wording if touched.
- source-of-truth links that call future roadmaps proposed/planning-only unless
  separately implemented.

## QA Sub-Agent

Baseline validation for future scoped work:

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
```

Add phase-specific commands when code changes:

```bash
make validate
make test
docker compose -f deploy/docker-compose.yml config
```

Add UI-specific checks when UI changes:

```bash
RUN_LOCAL_APP=true make release-candidate-check
make audit-product-acceptance
```

If a command is unavailable or blocked, record exact blocker and do not claim pass.

## Claim-Boundary Sub-Agent

Block any wording that implies:

- CAL-ITP/Caltrans compliance;
- agency adoption or approval;
- consumer submission/review/acceptance/ingestion/listing/display;
- final-root readiness;
- hosted SaaS availability;
- paid support or SLA;
- production readiness;
- vendor compatibility or hardware certification;
- production-grade ETA quality.

Allowed wording examples:

- `supports browser-first GTFS/GTFS-Realtime operations workflows`;
- `provides local/reference diagnostics`;
- `shows readiness signals and blockers`;
- `runs synthetic connector checks`;
- `keeps evidence/adoption tracks authorization-gated`.

## Required checkpoint report

After every checkpoint:

```text
Checkpoint:
Sub-agents used or simulated, including intended model level:
Changed files:
Validation run:
Blocked checks:
Protected path status:
Consumer tracker status:
Claim-boundary status:
Master review:
Required edits:
Decision:
Commit created:
Next checkpoint:
```

A phase is not closed unless `Decision: closed` and `Commit created:` records the final closeout commit.
