# Phase 86 -- Multi-Agency, Roles, Audit, And Accessibility

## Goal

Harden the private Operations Console around agency scope visibility, role
clarity, auditability, and accessibility without claiming production
multi-tenancy or adding public admin surfaces.

Phase 86 is product work only. It may improve private UI, route behavior,
tests, and documentation. It must not collect evidence, move consumer tracker
statuses, publish release artifacts, contact external services, or claim
production multi-tenant readiness.

## Scope

- Agency scope visibility and safe agency-switcher/status guidance, using the
  current authenticated principal as the source of truth.
- Role and permission explanations for `admin`, `editor`, `operator`, and
  `read_only`.
- Access-denied UX that explains missing permission or agency mismatch without
  leaking cross-agency data.
- Audit log browser only if existing records and scoped queries support it
  safely; otherwise add a docs/model view that states audit browsing is not yet
  available.
- Accessibility pass across core private Operations Console pages: keyboard
  navigation, focus order, labels, table captions, skip links, status
  semantics, reduced motion, responsive/mobile behavior, and no-JavaScript
  fallback where practical.

## Non-Goals

- No production multi-tenant hosting claim.
- No global admin powers.
- No cross-agency data visibility from client-supplied `agency_id`.
- No migration unless a later checkpoint proves it is required and the
  Data/Migration review approves a backward-compatible path.
- No public admin routes.
- No evidence writes, consumer status changes, external contacts, release tags,
  packages, published images, or hosted-service claims.

## Starting Repo Truth

- Auth principals carry `Subject`, `AgencyID`, `Roles`, and `Method`.
- Roles are `admin`, `editor`, `operator`, and `read_only` in
  `internal/auth`.
- Role lookup is scoped by authenticated claim agency in
  `PostgresRoleStore.RolesForSubject`.
- `auth.RequireRole` rejects conflicting `agency_id` query values and provides
  private `403` responses for role or agency mismatch.
- `docs/multi-agency-strategy.md` documents multi-agency boundaries and says
  future audit readers require review.
- The private Operations Console is server-rendered Go with progressive
  buildless JavaScript only where already introduced.

## Simulated Sub-Agent Reviews

Real additional sub-agents were unavailable because the agent thread limit had
already been reached, so the Master Agent simulated the required roles with the
intended model levels.

- Context / Repo Truth Sub-Agent, GPT-5.5 x-high: confirmed Phase 85 is closed,
  Phase 86 is next, protected evidence paths are off limits, all consumer
  targets remain prepared, roles are already agency-scoped, and
  `docs/multi-agency-strategy.md` is the current strategy reference.
- Planning Sub-Agent, GPT-5.5 x-high: approved a conservative sequence:
  planning doc, scope/role UI, access-denied UX, audit model/browser decision,
  accessibility pass, closeout.
- UI/UX Sub-Agent, GPT-5.5 high: required plain-language role labels, visible
  agency scope, no confusing switcher if switching is not actually supported,
  and accessible tables/forms.
- Documentation / IA Sub-Agent, GPT-5.5 high: required source-of-truth status
  alignment and a closeout handoff.
- Claim-Boundary Sub-Agent, GPT-5.5 high: required explicit "not production
  multi-tenant readiness" wording and blocked hosted SaaS, SLA, adoption,
  consumer, compliance, release, vendor, hardware, and ETA claims.
- Security/Auth Sub-Agent, GPT-5.5 high: required agency query conflict tests,
  role checks, CSRF preservation for unsafe methods, no cross-agency leaks, and
  no global admin expansion.
- QA Sub-Agent, GPT-5.5 high: required focused route tests plus baseline
  audits, exact prepared-only consumer tracker checks, and protected-path
  checks after every checkpoint.
- Implementation Sub-Agent, GPT-5.5 high: will execute only approved checkpoint
  scopes.
- Data/Migration Sub-Agent, GPT-5.5 high: not used for CP000001; will be
  required if a later checkpoint proposes persistence or migration work.

All required edits from these reviews are incorporated into this plan.

## Checkpoints

```text
Phase 86 -- Checkpoint 000001: add multi-agency roles audit accessibility plan
Phase 86 -- Checkpoint 000002: add agency scope and switcher improvements
Phase 86 -- Checkpoint 000003: add role permission and access-denied UX
Phase 86 -- Checkpoint 000004: add audit log browser or scoped audit model docs
Phase 86 -- Checkpoint 000005: run accessibility and keyboard navigation review
Phase 86 -- Checkpoint 000006: close multi-agency roles audit accessibility review
```

## Acceptance Criteria

- Core private Operations Console pages show the authenticated agency scope and
  current role context without allowing client-side agency switching to override
  server-owned authorization.
- Any agency switcher/status UI is truthful: it either uses supported server
  scope safely or clearly states that cross-agency switching is not supported
  in this private reference mode.
- Role and permission copy explains what each role can do without exposing
  secrets or implying hosted admin service availability.
- Access-denied responses are private, bounded, no-store, and explain likely
  next steps without leaking whether another agency has data.
- Audit log browsing is added only if existing scoped records can be queried
  safely. Otherwise the UI/doc states that audit-log browsing needs a future
  persistence/review phase.
- Accessibility review covers keyboard focus, skip links, labels, table
  captions, status semantics, reduced-motion preference, responsive layout, and
  no-JavaScript fallback for core workflows.
- Consumer tracker rows remain prepared-only.
- Protected evidence and consumer packet paths remain untouched.

## Validation

Baseline Phase 86 validation:

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

Additional Phase 86 checks:

```bash
go test ./cmd/agency-config -run 'Operations|Role|Agency|Access|Audit|Accessibility'
go test ./internal/auth ./internal/tenant
make validate
make test
RUN_LOCAL_APP=true make release-candidate-check
```
