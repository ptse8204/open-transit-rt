# Phase 79 -- Agency Setup V3 Plan

## Scope

Phase 79 makes the private first-run setup path clearer for nontechnical
agency staff while preserving the existing Go server-rendered Operations
Console, authenticated private routes, role checks, CSRF protections, and
claim boundaries.

The phase may improve:

- the private setup wizard progress model;
- the `/admin/operations/setup` agency profile and publication metadata page;
- safe metadata guidance for agency profile, public base URL, feed base URL,
  license, technical contact, and publication environment;
- environment/config diagnostics that do not expose secrets;
- admin account and role visibility without showing tokens or private
  credential material;
- technical-helper escalation cards for local app startup, validator/tooling,
  feed-root configuration, backup/restore, and support bundle work;
- private JSON shape tests for setup progress and claim flags.

The phase does not add public admin routes, contact external parties, collect
or retain evidence, move consumer statuses, claim final-root readiness, claim
agency approval, publish a release artifact, or create a hosted SaaS,
production, vendor, SLA, hardware, compliance, consumer, or ETA-quality claim.

## Implementation Boundary

Use existing private surfaces:

- `GET /admin/operations/setup-wizard`
- `GET /admin/operations/setup-wizard.json`
- `GET /admin/operations/setup`
- admin-only `POST /admin/operations/setup`

Do not add a new persistence model unless current publication metadata cannot
represent the safe fields. No migration is planned for Phase 79.

Any mutation remains admin-only, agency-scoped, CSRF-protected for cookie auth,
body-capped by the existing setup route, and limited to the existing
publication bootstrap/update workflow. Browser setup must not accept raw
commands, file paths, validator binaries, evidence destinations, secrets, or
external submission targets.

## Master Approval

The Master Agent approves implementation only after required sub-agent edits
are incorporated. Implementation is limited to:

- expanding setup wizard V3 JSON with grouped operator paths, blockers, and
  technical-helper escalation;
- improving setup page copy and layout for safe metadata review;
- adding safe environment/config/role diagnostics derived from existing private
  config and principal data;
- adding focused route, shape, rendering, auth, CSRF, and claim-boundary tests;
- updating closeout/status docs.

## Sub-Agent Review Plan

Real or simulated reviews use the intended model levels from the authorized
Phase 75-90 track:

- Context / Repo Truth Sub-Agent, GPT-5.5 x-high: confirm current repo truth,
  protected paths, consumer tracker, setup routes, setup wizard JSON, and
  existing tests.
- Planning Sub-Agent, GPT-5.5 x-high: confirm checkpoint sequence and safe
  field scope.
- Implementation Sub-Agent, GPT-5.5 high: simulated if thread capacity is
  exhausted; identify the smallest setup V3 code slice.
- QA Sub-Agent, GPT-5.5 high: confirm focused tests and validation commands.
- UI/UX Sub-Agent, GPT-5.5 high: review nontechnical setup language, field
  labels, and escalation cards.
- Documentation / IA Sub-Agent, GPT-5.5 high: simulated if thread capacity is
  exhausted; review source-of-truth updates.
- Claim-Boundary Sub-Agent, GPT-5.5 high: confirm no final-root, agency
  approval, compliance, consumer, hosted SaaS, production, vendor, SLA, or ETA
  claim.
- Security/Auth Sub-Agent, GPT-5.5 high: confirm role, CSRF, agency-scope,
  body cap, and no-secret display boundaries.

## Checkpoints

```text
Phase 79 -- Checkpoint 000001: add agency setup v3 plan
Phase 79 -- Checkpoint 000002: implement agency profile and metadata review
Phase 79 -- Checkpoint 000003: improve browser GTFS source import review
Phase 79 -- Checkpoint 000004: add setup progress and blocker next actions
Phase 79 -- Checkpoint 000005: close agency setup v3 review
```

## Validation

Baseline Phase 79 validation:

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

Additional Phase 79 checks:

```bash
go test ./cmd/agency-config -run 'Test.*(Setup|SetupWizard|OperationsConsole|Progressive)'
go test ./cmd/agency-config ./internal/auth ./internal/tenant
make validate
make test
```

If an environment limitation blocks a check, record the exact blocker in the
Phase 79 handoff.
