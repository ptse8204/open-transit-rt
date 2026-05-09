# Phase 59 -- Real Pilot Closeout

## Status

Planning accepted. Execution must use the blocker-only branch unless a
maintainer supplies real retained pilot authorization and public-safe pilot
closeout artifacts before implementation.

## Goal

Close a real authorized pilot with public-safe outcome evidence when such a
pilot exists. If no real pilot is available, close Phase 59
blocker-documented only without creating fake pilot evidence or stronger
claims.

## Current Inspection Result

No retained Phase 59 pilot authorization record, kickoff note,
agency/operator feedback record, operations closeout, or continue/pause/close
decision artifact exists in the current repository.

Existing retained artifacts are limited to earlier hosted/operator OCI pilot
evidence and local public-GTFS pilot evidence. Those packets remain useful for
their documented scopes only. They do not prove real pilot authorization,
agency feedback, agency adoption, final-root readiness, consumer acceptance,
CAL-ITP/Caltrans compliance, production readiness, vendor compatibility,
SLA/uptime, or production-grade ETA quality.

## Scope

- Inspect for retained real pilot authorization and public-safe closeout
  artifacts before any evidence claim.
- If real artifacts exist, close exactly one named pilot with:
  - authorization record;
  - kickoff notes;
  - GTFS import and validation results;
  - telemetry or adapter notes;
  - operations records;
  - agency/operator feedback;
  - blocker, continue, or pause decision;
  - closeout summary.
- If real artifacts do not exist, close blocker-documented only.
- Update phase, status, roadmap, backlog, open-question, and handoff docs.

## Non-Goals

- No fake evidence.
- No consumer contact.
- No portal automation.
- No consumer submission.
- No consumer status change.
- No compliance claim.
- No agency adoption, endorsement, or approval claim.
- No production-readiness claim.
- No hosted SaaS claim.
- No SLA or uptime claim.
- No vendor compatibility claim.
- No marketplace approval or paid support claim.
- No production-grade ETA claim.
- No raw private GTFS, telemetry, logs, credentials, private contacts,
  private ticket links, private screenshots, or unredacted correspondence.

## Files Likely To Change

- `docs/phase-59-real-pilot-closeout.md`
- `docs/handoffs/phase-59.md`
- `docs/handoffs/latest.md`
- `docs/current-status.md`
- `docs/backlog.md`
- `docs/open-questions.md`
- `docs/roadmap-to-calitp-compliance-and-gap-closure.md`

Execution must not change:

- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/`
- `docs/evidence/consumer-submissions/artifacts/`
- `docs/evidence/consumer-submissions/packets/`
- `docs/evidence/captured/`

## Safety Boundaries

Default Phase 59 execution is blocker-only documentation. No retained evidence
packet may be created unless all of these are true:

- a real pilot authorization artifact exists;
- the artifact is redacted and public-safe;
- pilot kickoff, operations, feedback, and closeout-decision artifacts exist;
- artifacts contain no secrets, private credentials, raw private telemetry,
  raw logs, private contacts, private URLs, cookies, Authorization headers,
  database URLs, private TLS material, or unredacted correspondence;
- a maintainer explicitly approves retention under `docs/evidence/captured`;
- retained artifacts pass redaction and claim-boundary review.

Prior OCI/local pilot artifacts must remain labeled as earlier-scope
operator/local evidence only.

## Evidence And Claim Boundaries

For the current blocker-only branch, the only allowed Phase 59 claim is that
no real authorized pilot closeout artifacts were available in the repository.

If a future real pilot closeout is retained, the packet may claim only what
that named pilot and retained public-safe artifacts directly support. Pilot
participation alone must not be represented as agency adoption, agency
approval, consumer acceptance, production readiness, compliance, SLA/uptime,
hosted service availability, vendor compatibility, marketplace approval, or
production-grade ETA proof.

All seven consumer and aggregator targets remain `prepared` unless retained
target-originated evidence supports a target-specific transition.

## Implementation Details

Current execution should:

1. Keep Phase 59 docs-only because required real pilot artifacts are absent.
2. Add a Phase 59 handoff documenting blocker-only closure.
3. Update current status, latest handoff, backlog, open questions, and roadmap
   wording so Phase 60 starts from the correct evidence boundary.
4. Leave `docs/evidence`, consumer submission trackers, consumer artifacts,
   and consumer packets unchanged.

Do not add a pilot evidence collector during the blocker-only path. Phase 52
and Phase 55 already provide guarded evidence workflow examples; adding a new
Phase 59 collector without real pilot inputs would increase surface area
without closing the real-pilot blocker.

## Tests

No code or script changes are planned for the blocker-only branch.

Required verification remains:

- `make validate`
- `make test`
- `make smoke`
- `git diff --check`
- consumer tracker JSON parse
- exact seven-target prepared-only consumer tracker check
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json`
- `git diff --exit-code -- docs/evidence/consumer-submissions/current docs/evidence/consumer-submissions/artifacts docs/evidence/consumer-submissions/packets docs/evidence/captured`
- consumer artifact directory scan
- `docker compose -f deploy/docker-compose.yml config`
- `INTEGRATION_TESTS=1 make test-integration` when the local DB is available

## Performance And Scale Tests

No performance or scale test is relevant to the blocker-only branch.

For a future real pilot closeout, public-safe runtime observations may include
import duration, GTFS entity counts, validator runtime/status, feed fetch byte
sizes, and operations diagnostic timestamps. Those observations are local
pilot facts only and must not be described as SLA, production capacity,
production readiness, compliance, or ETA-quality proof.

## Required Verification Commands

Run from the repository root:

```sh
make validate
make test
make smoke
git diff --check
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
git diff --exit-code -- docs/evidence/consumer-submissions/status.json
git diff --exit-code -- docs/evidence/consumer-submissions/current docs/evidence/consumer-submissions/artifacts docs/evidence/consumer-submissions/packets docs/evidence/captured
find docs/evidence/consumer-submissions/artifacts -mindepth 2 -maxdepth 2 -type f ! -name README.md -print
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured
docker compose -f deploy/docker-compose.yml config
INTEGRATION_TESTS=1 make test-integration
```

The artifact `find` command must print no files for the current blocker-only
branch.

## Master Plan Review

Accepted for blocker-only execution. The plan has clear scope, non-goals,
safety boundaries, claim/evidence boundaries, docs updates, validation
commands, consumer tracker preservation checks, and no evidence-writing risk
for the current repository state.
