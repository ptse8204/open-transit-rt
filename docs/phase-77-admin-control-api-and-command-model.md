# Phase 77 -- Admin Control API And Command Model Plan

## Scope

Phase 77 introduces a bounded private command/result model for Operations
Console workflows. It does not expose a public API, shell commands, raw paths,
raw validator output, evidence writes, consumer tracker changes, migrations,
release artifacts, or broader production/compliance claims.

The first migrated action is the existing private validator-health `run_all`
workflow because it already has the right safety shape:

- route: `/admin/operations/validation-health`;
- role: admin-only POST;
- review: read-only GET and JSON remain available to read-only/operator/editor/admin;
- request cap: `validationHealthPostMaxBytes`;
- browser CSRF: required for cookie-auth POSTs;
- action allowlist: `run_all` only;
- field allowlist: `action` and `csrf_token` only;
- execution model: server-owned validator/tooling/artifact mappings only;
- output model: derived validation health rows and normal validation records
  only, with no raw stdout/stderr/private path exposure.

## Command Ladder

Phase 77 documents and tests this ladder for future phases:

| Level | Meaning | Browser default |
| --- | --- | --- |
| read-only refresh | Re-read existing private records or recompute a derived summary. | Allowed for read-only roles when no mutation occurs. |
| dry run | Execute a bounded check that writes no durable product state. | Admin/operator only unless existing route says otherwise; result must be explicit. |
| reversible private change | Update private configuration or records with a clear rollback path. | Admin-only unless separately justified. |
| publish/activate | Change public feed output, active schedule, or externally visible state. | Admin-only with preview, confirmation, audit, and rollback notes. |
| destructive/hard-to-reverse | Delete, rotate, restore, overwrite, or irreversibly alter data/secrets. | Disabled by default unless a later phase explicitly authorizes browser execution. |

## Master Approval

The Master Agent approves implementation when sub-agent review has no required
edits or when required edits are incorporated. Implementation is limited to:

- adding internal command/result structs and helper functions;
- documenting command safety rules;
- mapping the validator-health `run_all` POST into the command-result model;
- rendering bounded command outcome text in the existing private page;
- adding focused tests for auth, roles, CSRF, body caps, unsupported fields,
  safe result JSON shape, claim flags, and no raw command/path exposure.

## Sub-Agent Review Plan

Real or simulated reviews use the intended model levels from the authorized
Phase 75-90 track:

- Context / Repo Truth Sub-Agent, GPT-5.5 x-high: confirm current repo truth,
  protected paths, consumer tracker, validation-health route behavior, and
  relevant tests.
- Planning Sub-Agent, GPT-5.5 x-high: confirm checkpoint sequence and the
  low-risk validator-health migration.
- Implementation Sub-Agent, GPT-5.5 high: identify the smallest code slice and
  files to edit.
- QA Sub-Agent, GPT-5.5 high: confirm focused tests and baseline validation.
- UI/UX Sub-Agent, GPT-5.5 high: simulated if sub-agent capacity is exhausted;
  review command outcome wording and confirmation copy.
- Documentation / IA Sub-Agent, GPT-5.5 high: simulated if sub-agent capacity
  is exhausted; review command model docs and source-of-truth updates.
- Claim-Boundary Sub-Agent, GPT-5.5 high: confirm result statuses and wording
  do not imply compliance, production readiness, consumer acceptance, vendor
  compatibility, SLA, or ETA quality.
- Security/Auth Sub-Agent, GPT-5.5 high: confirm no auth/CSRF/role/request-cap
  regression.

## Checkpoints

```text
Phase 77 -- Checkpoint 000001: add admin control API and command model plan
Phase 77 -- Checkpoint 000002: define private command result contracts
Phase 77 -- Checkpoint 000003: add command safety tests
Phase 77 -- Checkpoint 000004: migrate one low-risk workflow to command model
Phase 77 -- Checkpoint 000005: close admin control API review
```

## Validation

Baseline Phase 77 validation:

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

Additional Phase 77 checks:

```bash
go test ./cmd/agency-config
go test ./internal/auth ./internal/tenant ./cmd/feed-alerts ./cmd/gtfs-studio
make validate
make test
docker compose -f deploy/docker-compose.yml config
```

If an environment limitation blocks a check, record the exact blocker in the
Phase 77 handoff.
