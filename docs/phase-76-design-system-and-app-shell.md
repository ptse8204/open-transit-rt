# Phase 76 -- Design System And App Shell Plan

## Scope

Phase 76 improves the private Operations Console shell and design system while
preserving the existing Go server-rendered architecture. It does not add a SPA,
public admin route, migration, evidence write, consumer tracker change, or
backend behavior change.

## Master Approval

The Master Agent approves the Phase 76 plan for implementation with these
constraints:

- keep existing private route URLs stable;
- keep current auth, role, CSRF, and device-token rules unchanged;
- use shared CSS tokens and reusable template vocabulary inside the existing
  Operations Console layout;
- preserve browser-first operator language without making compliance,
  production, adoption, consumer, hosted SaaS, vendor, SLA, hardware, or ETA
  quality claims;
- add focused tests for shell consistency, responsive behavior, accessibility,
  and claim-boundary-safe styling.

## Sub-Agent Review Plan

Real or simulated reviews use the intended model levels from the authorized
Phase 75-90 track:

- Context / Repo Truth Sub-Agent, GPT-5.5 x-high: confirm the current shell,
  docs, protected paths, consumer tracker, and route contracts.
- Planning Sub-Agent, GPT-5.5 x-high: confirm checkpoint sequence and the
  no-SPA implementation path.
- Implementation Sub-Agent, GPT-5.5 high: apply the shared layout tokens,
  components, and route-stable shell refinements.
- QA Sub-Agent, GPT-5.5 high: run focused tests and the baseline validation
  commands.
- UI/UX Sub-Agent, GPT-5.5 high: check operator hierarchy, responsive behavior,
  status semantics, and nontechnical readability.
- Documentation / IA Sub-Agent, GPT-5.5 high: confirm docs/status handoff
  alignment.
- Claim-Boundary Sub-Agent, GPT-5.5 high: confirm no forbidden claim or consumer
  status movement.
- Security/Auth Sub-Agent, GPT-5.5 high: confirm no auth, CSRF, public route, or
  browser command execution regression.

Implementation starts only when reviews have no required edits or when required
edits are incorporated before the implementation checkpoint.

## Checkpoints

```text
Phase 76 -- Checkpoint 000001: add design system and app shell plan
Phase 76 -- Checkpoint 000002: implement shared layout tokens and components
Phase 76 -- Checkpoint 000003: apply shell to core Operations Console routes
Phase 76 -- Checkpoint 000004: add responsive and accessibility baseline checks
Phase 76 -- Checkpoint 000005: close design system and app shell review
```

## Validation

Baseline Phase 76 validation:

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

Additional Phase 76 checks:

```bash
go test ./cmd/agency-config
make validate
make test
RUN_LOCAL_APP=true make release-candidate-check
```

If an environment limitation blocks a check, record the exact blocker in the
Phase 76 handoff.

## Closeout Notes

Phase 76 closes as private Operations Console design-system and app-shell work
only. It:

- kept Go server-rendered templates and existing private route URLs;
- added static shared design tokens for color, spacing, typography, focus,
  status chips, cards, tables, forms, responsive layout, and reduced motion;
- added shared shell markers for the private control plane header, breadcrumb,
  metadata row, navigation, and main landmark;
- aligned route groups to Start Here, Schedule, Realtime, Connectors, Health,
  Maintain, and Learn;
- marked GTFS Studio and Alerts Console links as separate private admin
  surfaces;
- changed generic `Ready` UI status to `Ready for local review`;
- changed consumer and feed wording to prepared-packet and configured
  local/reference language where needed;
- added focused tests for shared shell coverage, design tokens, responsive
  behavior, accessibility landmarks, external admin surface markers, and
  overclaim-prone wording.

Phase 76 did not add public admin routes, migrations, API behavior, evidence
writes, consumer tracker changes, release artifacts, hosted SaaS claims,
production-readiness claims, vendor claims, SLA claims, or ETA-quality claims.
