# Phase 92 -- Clean Checkout Release-Candidate Gate

This phase runs a serious local clean-checkout release-candidate diagnostic
gate and records exact blockers. It is a review artifact only: it is not a
release tag, package, publication, retained evidence packet, compliance proof,
consumer proof, final-root proof, production-readiness proof, hosted-service
proof, SLA/uptime proof, vendor proof, hardware proof, or ETA-quality proof.

## Scope

- Prepare and run the clean-checkout product validation gate.
- Run local app and five public feed diagnostics from the clean checkout.
- Run synthetic/local connector backend checks and claim-boundary audits.
- Record exact pass, needs-review, not-checked, and blocked statuses.
- Close with protected-path, consumer tracker, claim-boundary, security/auth,
  and data/migration review.

## Out Of Scope

- Protected evidence path writes.
- Consumer tracker status changes.
- External contact with agencies, vendors, consumers, portals, map providers,
  aggregators, or other external services for proof.
- Real credentials, real private payloads, or real agency/vendor/device data.
- `git tag`, `git push --tags`, GitHub Release creation, public image
  publication, package publication, or public announcements.
- Phase 95 package generation/audit commands such as `make release-package` or
  `make audit-release-package`.
- Release-ready, compliance, adoption, consumer acceptance, production
  readiness, final-root readiness, hosted SaaS, vendor compatibility, hardware
  certification, SLA/uptime, or production-grade ETA claims.

## Checkpoints

```text
Phase 92 -- Checkpoint 000001: add clean checkout rc gate plan
Phase 92 -- Checkpoint 000002: run clean checkout product validation
Phase 92 -- Checkpoint 000003: run local app and five-feed diagnostics
Phase 92 -- Checkpoint 000004: run connector backend and claim-boundary diagnostics
Phase 92 -- Checkpoint 000005: record rc gate result and blockers
Phase 92 -- Checkpoint 000006: close clean checkout rc gate
```

## Clean Checkout Strategy

Use a local ignored checkout under `.cache/phase-92-clean-checkout` so generated
diagnostics stay local and out of source control. Prefer `git worktree add
--detach .cache/phase-92-clean-checkout HEAD`; if worktree creation is
unavailable, use a local clone under `.cache` and record the fallback.

The clean checkout must run only local repository diagnostics. It may write
temporary `.cache` diagnostics inside the clean checkout, but it must not write
under `docs/evidence/**`, move consumer statuses, publish artifacts, or contact
external parties for proof.

## Validation Plan

Checkpoint 000002 clean-checkout product validation:

```bash
git status --short
git diff --check
make check
make validate
make test
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

Checkpoint 000003 local app and five-feed diagnostics:

```bash
RUN_LOCAL_APP=true make release-candidate-check
```

The five local feed paths are:

- `/public/feeds.json`
- `/public/gtfs/schedule.zip`
- `/public/gtfsrt/vehicle_positions.pb`
- `/public/gtfsrt/trip_updates.pb`
- `/public/gtfsrt/alerts.pb`

Checkpoint 000004 connector backend and claim-boundary diagnostics:

```bash
make external-connection-check
make adapter-conformance
make test-connector-examples
docker compose -f deploy/docker-compose.yml config
make audit-product-acceptance
make audit-final-claim-review
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
```

## Result Taxonomy

| Status | Meaning |
| --- | --- |
| `passed` | The local command or review completed with no required follow-up. |
| `needs_review` | The local diagnostic completed but does not prove release readiness, public acceptance, production readiness, compliance, real-world quality, or external compatibility. |
| `not_checked` | The command or artifact was outside Phase 92 scope or could not run safely. |
| `blocked` | The command failed or a required safe precondition was missing. |

## Checkpoint 000001 Report

Checkpoint:
Phase 92 -- Checkpoint 000001: add clean checkout rc gate plan.

Sub-agents used or simulated, including intended model level:
Real Planning Sub-Agent -- GPT-5.5 x-high returned the checkpoint plan. Real
Context / Repo Truth Sub-Agent -- GPT-5.5 x-high, Release / Supply-Chain
Sub-Agent -- GPT-5.5 high, and Claim-Boundary / Security Sub-Agent --
GPT-5.5 high are running for Phase 92. Master Agent -- GPT-5.5 x-high,
current thread.

Changed files:
`docs/phase-92-clean-checkout-release-candidate-gate.md`;
`docs/roadmaps/post-90-agency-grade-gtfs-rt-product/02-phases-and-checkpoints.md`.

Validation run:
`git status --short` was clean before edits. `git diff --check` passed.

Blocked checks:
Clean-checkout product validation, local app/five-feed diagnostics,
connector/backend diagnostics, and phase closeout validation are scheduled for
later Phase 92 checkpoints.

Protected path status:
No protected evidence path was edited or required.

Consumer tracker status:
The tracker is not edited in Phase 92; the exact seven prepared-only check is
scheduled for product validation and closeout.

Claim-boundary status:
The plan is bounded to local diagnostics and explicit blockers. It does not
claim release readiness, compliance, adoption, consumer acceptance, production
readiness, final-root readiness, hosted-service availability, vendor
compatibility, hardware certification, SLA/uptime, or ETA quality.

Security/auth status:
No route, auth behavior, token handling, credential path, public exposure, or
admin command behavior changed.

Data/migration status:
No persistence, migration, GTFS data model, or realtime data model change is
included.

Master review:
Approved. The phase plan uses the required autonomous checkpoint sequence,
keeps Phase 95 packaging out of Phase 92, and preserves protected-path,
consumer-status, no-publish, and no-claim boundaries.

Required edits:
None for CP000001.

Decision:
Proceed to CP000001 validation and commit, then CP000002 clean-checkout product
validation.

Next checkpoint:
Phase 92 -- Checkpoint 000002: run clean checkout product validation.
