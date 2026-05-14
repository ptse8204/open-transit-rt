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

## Checkpoint 000002 Clean-Checkout Product Validation

Clean checkout path:
`.cache/phase-92-clean-checkout`

Clean checkout source:
Detached worktree from commit `1b9b55a` (`Phase 92 -- Checkpoint 000001: add
clean checkout rc gate plan`).

| Check | Result | Notes |
| --- | --- | --- |
| `git status --short` | passed | Clean before and after product validation. |
| `git diff --check` | passed | No whitespace errors. |
| `make check` | passed | Lightweight no-network/no-Docker/no-validator-install checks passed. |
| `make validate` initial run | blocked | Pinned static GTFS validator was missing from the clean checkout's local `.cache`; command reported to run `make validators-install`. |
| `make validators-install` | passed | Installed pinned validator tooling into the clean checkout's ignored `.cache` path and pulled the pinned GTFS-RT validator image. |
| `make validate` rerun | passed | Validator tooling check and validation smoke passed. |
| `make test` | passed | `go test ./...` passed. |
| Consumer tracker JSON parse | passed | `docs/evidence/consumer-submissions/status.json` parsed as JSON. |
| Exact prepared-only consumer tracker | passed | All seven targets remain exactly `prepared`: Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, transit.land. |
| Protected path status in clean checkout | passed | No status under `docs/evidence/consumer-submissions`, `docs/evidence/captured`, `db/migrations`, `go.mod`, or `go.sum`. |
| Protected path status in main checkout | passed | Main checkout remained clean for the same protected paths. |

The initial `make validate` failure was an environment/tooling precondition in
the clean checkout. It was resolved by running the repository-supported pinned
validator installer. The installer wrote only ignored local `.cache` tooling in
the clean checkout and did not write evidence paths, consumer status records,
source code, migrations, `go.mod`, or `go.sum`.

## Checkpoint 000002 Report

Checkpoint:
Phase 92 -- Checkpoint 000002: run clean checkout product validation.

Sub-agents used or simulated, including intended model level:
Real Context / Repo Truth Sub-Agent -- GPT-5.5 x-high reported clean-checkout
commands, protected-path risks, and likely validator/tooling blockers. Real
Release / Supply-Chain Sub-Agent -- GPT-5.5 high reported Phase 92-safe
commands and deferred package commands. Real Claim-Boundary / Security
Sub-Agent -- GPT-5.5 high reported guardrails and passed claim audits. Real
Planning Sub-Agent -- GPT-5.5 x-high provided the checkpoint plan. Master
Agent -- GPT-5.5 x-high, current thread.

Changed files:
`docs/phase-92-clean-checkout-release-candidate-gate.md`.

Validation run:
From `.cache/phase-92-clean-checkout`: `git status --short` passed; `git diff
--check` passed; `make check` passed; `make validate` initially failed because
the clean checkout lacked pinned validator cache tooling; `make
validators-install` passed; `make validate` rerun passed; `make test` passed;
consumer tracker JSON parse passed; exact seven-target prepared-only tracker
check passed; protected-path status check passed. From the main checkout,
protected-path status check passed.

Blocked checks:
No CP000002 product validation check remains blocked after installing pinned
local validator tooling. Local app/five-feed diagnostics, connector/backend
diagnostics, and closeout validation remain scheduled for later Phase 92
checkpoints.

Protected path status:
No protected evidence path was edited or generated. Clean checkout and main
checkout protected-path status checks returned clean.

Consumer tracker status:
All seven targets remain exactly `prepared` in the required order. The tracker
file was not edited.

Claim-boundary status:
The validation result is local diagnostic status only. It does not claim
release readiness, compliance, adoption, consumer acceptance, production
readiness, final-root readiness, hosted-service availability, vendor
compatibility, hardware certification, SLA/uptime, or ETA quality.

Security/auth status:
No route, auth behavior, token handling, credential path, public exposure, or
admin command behavior changed.

Data/migration status:
No persistence, migration, GTFS data model, or realtime data model change is
included. `db/migrations`, `go.mod`, and `go.sum` status checks are clean.

Master review:
Approved. The clean-checkout product gate passed after resolving a local
validator tooling precondition with the repo-supported installer, and no
protected path, consumer tracker, release, or claim boundary was crossed.

Required edits:
None for CP000002.

Decision:
Proceed to CP000002 validation and commit, then CP000003 local app and
five-feed diagnostics.

Next checkpoint:
Phase 92 -- Checkpoint 000003: run local app and five-feed diagnostics.

## Checkpoint 000003 Local App And Five-Feed Diagnostics

Command:

```bash
RUN_LOCAL_APP=true make release-candidate-check
```

Clean checkout path:
`.cache/phase-92-clean-checkout`

Diagnostic output path:
`.cache/phase-92-clean-checkout/.cache/release-candidate-check/20260514T231229Z`

| Check | Result | Notes |
| --- | --- | --- |
| `RUN_LOCAL_APP=true make release-candidate-check` | passed | Command exited `0` and wrote private local `.cache` diagnostics only. |
| Local app startup and five public feeds | passed | The helper started the local app and fetched `/public/feeds.json`, `/public/gtfs/schedule.zip`, `/public/gtfsrt/vehicle_positions.pb`, `/public/gtfsrt/trip_updates.pb`, and `/public/gtfsrt/alerts.pb`. |
| Diagnostic summary counts | needs_review | Summary recorded `35` passed, `0` blocker, `0` needs_review, and `4` not_checked rows; overall status is `not_checked` because release package audit and bounded follow-up commands are outside this helper's run. |
| Release package audit | not_checked | Phase 92 did not run package generation or package audit. These are deferred to Phase 95. |
| Repository validation inside helper | not_checked | The helper intentionally does not run `make validate` or `make test`; CP000002 ran both directly from the clean checkout. |
| Java runtime row | needs_review | The helper reported the Java tool row as `passed` while the detail contained the macOS system-stub message `Unable to locate a Java Runtime`; CP000002 independently passed the repo-supported pinned validator check after local validator installation. |
| `make agency-app-down` | passed | Local app containers and network were stopped and removed after diagnostics. |
| Protected path status | passed | No protected evidence path, consumer tracker file, migration, `go.mod`, or `go.sum` status was created by this checkpoint. |

The diagnostic summary's `overall_status=not_checked` is expected for Phase 92
because release package auditing is intentionally out of scope and the helper
does not run the full validation/test suite internally. CP000002 already ran
`make validate` and `make test` from the same clean checkout.

## Checkpoint 000003 Report

Checkpoint:
Phase 92 -- Checkpoint 000003: run local app and five-feed diagnostics.

Sub-agents used or simulated, including intended model level:
Real Context / Repo Truth Sub-Agent -- GPT-5.5 x-high, Real Release /
Supply-Chain Sub-Agent -- GPT-5.5 high, Real Claim-Boundary / Security
Sub-Agent -- GPT-5.5 high, and Real Planning Sub-Agent -- GPT-5.5 x-high
informed this checkpoint. Master Agent -- GPT-5.5 x-high, current thread.

Changed files:
`docs/phase-92-clean-checkout-release-candidate-gate.md`.

Validation run:
From `.cache/phase-92-clean-checkout`: `RUN_LOCAL_APP=true make
release-candidate-check` passed; diagnostic summary JSON parsed through
`python3 -m json.tool`; diagnostic summary showed the local app and five public
feeds passed; `make agency-app-down` passed; protected-path status check
passed.

Blocked checks:
No CP000003 local app/five-feed command remains blocked. Release package audit
is intentionally `not_checked` and deferred to Phase 95. Repository validation
and unit tests were intentionally not run by the helper and were already run
directly in CP000002.

Protected path status:
No protected evidence path was edited or generated. Diagnostics stayed under
the ignored clean-checkout `.cache` directory.

Consumer tracker status:
The tracker file was not edited. Exact prepared-only tracker validation passed
in CP000002 and remains scheduled again for closeout.

Claim-boundary status:
The local app and feed fetch result is local diagnostic status only. It does
not claim release readiness, public launch, public feed adoption, compliance,
consumer acceptance, production readiness, final-root readiness,
hosted-service availability, vendor compatibility, hardware certification,
SLA/uptime, or ETA quality.

Security/auth status:
No route, auth behavior, token handling, credential path, public exposure, or
admin command behavior changed. The local app stack was stopped after the
diagnostic.

Data/migration status:
The local app ran migrations in the temporary local stack and reported no
migrations to run at current version `9`; no repository migration file changed.

Master review:
Approved. The local app and five-feed diagnostic completed without crossing
release, evidence, consumer-status, or claim boundaries. The diagnostic
summary's `not_checked` overall state is correctly bounded to package and
helper-internal validation omissions, not product failure.

Required edits:
None for CP000003.

Decision:
Proceed to CP000003 validation and commit, then CP000004 connector backend and
claim-boundary diagnostics.

Next checkpoint:
Phase 92 -- Checkpoint 000004: run connector backend and claim-boundary diagnostics.
