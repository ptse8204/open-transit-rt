# Phase 108 -- Post-RC Bug Bash And Stabilization

## Goal

Stabilize the local post-RC product surface by rerunning the route audit,
validation, test, connector, and claim-boundary checks, fixing only scoped
bugs or documentation drift discovered during the review, and refreshing the
known blocker matrix without adding feature breadth or stronger claims.

## Current Repo Context

- Phase 89 and Phase 92 release-candidate gates remain `needs_review`.
- Phase 95 produced and audited a local `.cache` source package only; no tag,
  GitHub Release, package publication, image publication, or release action is
  authorized.
- Phase 107 aligned public docs, contributor onboarding, and architecture
  guidance, so Phase 108 should focus on rough edges, blockers, route drift,
  validation status, and stale wording.
- Protected evidence paths and the prepared-only consumer tracker remain hard
  boundaries.

## Scope

- Rerun the private Operations Console route inventory audit and record the
  result.
- Rerun repository validation, tests, connector checks, release-candidate
  diagnostics, product acceptance audit, and final claim audit.
- Triage test flakes only if they appear during repeated or focused local
  checks; do not hide failures with broad skips or retries.
- Patch stale copy, IA, docs references, blocker wording, or route drift when
  the fix is obvious, local, and bounded.
- Update the local `v0.1.0-rc.1` known blockers matrix to include post-RC bug
  bash status.
- Close with a Phase 108 handoff and source-of-truth status updates.

## Boundaries

- No protected evidence path writes.
- No consumer tracker status changes.
- No external contact, real credentials, real private payloads, or real
  agency/vendor/device data.
- No `git tag`, GitHub Release, package publication, image publication, public
  announcement, or remote release action.
- No claims of release readiness, compliance, adoption, consumer acceptance,
  final-root readiness, hosted SaaS, SLA/uptime, production readiness, vendor
  compatibility, hardware certification, production AVL reliability,
  production-grade ETA quality, or real-world ETA accuracy.

## Bug Bash Checklist

| Area | Review Action | Safe Output |
| --- | --- | --- |
| Route inventory | Run route inventory audit and strict docs mode. | Pass/fail result and exact blocker if any. |
| Private task flows | Check Operations Console route/task labels for stale references. | Copy/IA patch or no-change note. |
| Release blockers | Reconcile Phase 89, 92, 95, and 107 blocker language. | Updated local blockers matrix. |
| Validation | Run baseline, full validation, tests, connector checks, and compose config. | Command results only. |
| Test flakes | If a failure appears, isolate with focused repeated runs. | Deterministic fix or exact blocker. |
| Claim audit | Run product acceptance and final claim audits. | Pass/fail result with no stronger claims. |
| Protected paths | Check protected evidence paths, migrations, and module files. | Clean status or hard-stop/blocker note. |
| Consumer tracker | Parse and assert exact seven prepared targets. | Prepared-only confirmation. |

## Checkpoint Plan

- `Phase 108 -- Checkpoint 000001: add post-rc bug bash and stabilization plan`
- `Phase 108 -- Checkpoint 000002: implement primary scoped work`
- `Phase 108 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 108 -- Checkpoint 000004: close post-rc bug bash and stabilization review`

## Validation Plan

Checkpoint 000003 and closeout should run:

```bash
git status --short
git diff --check
make audit-operations-route-inventory
OPERATIONS_ROUTE_AUDIT_STRICT_DOCS=true scripts/audit-operations-route-inventory.sh
make check
make audit-product-acceptance
make audit-final-claim-review
make validate
make test
make external-connection-check
make adapter-conformance
make test-connector-examples
RUN_LOCAL_APP=true make release-candidate-check
docker compose -f deploy/docker-compose.yml config
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum
```

The exact seven-target prepared-only consumer tracker assertion must also run
at closeout.

## Checkpoint Report -- 000001

Checkpoint:
Phase 108 -- Checkpoint 000001: add post-rc bug bash and stabilization plan.

Sub-agents used or simulated, including intended model level:
Planning Sub-Agent -- GPT-5.5 x-high returned a bounded stabilization plan.
Context / Repo Truth Sub-Agent -- GPT-5.5 x-high was attempted, timed out, and
was shut down without edits; Context / Repo Truth was therefore simulated by
the Master Agent through direct inspection of Phase 89, Phase 92, Phase 95,
Phase 107, release-note, route-audit, and handoff docs. Implementation, QA,
UI/UX, Documentation / IA, Claim-Boundary, Security/Auth, Data/Migration, and
Release/Supply-Chain roles are simulated by the Master Agent for this plan
checkpoint. Master Agent -- GPT-5.5 x-high, current thread.

Changed files:
`docs/phase-108-post-rc-bug-bash-and-stabilization.md`.

Validation run:
Initial repository inspection reviewed the Phase 108 roadmap prompt, Phase 89
gate result, Phase 92 clean-checkout gate, Phase 95 draft release notes,
Phase 107 handoff, release-candidate readiness doc, route inventory audit
helper references, current status, latest handoff, roadmap status, and master
planner status. After adding the plan, `git status --short` showed only this
new Phase 108 plan doc; `git diff --check` passed; `python3 -m json.tool
docs/evidence/consumer-submissions/status.json >/dev/null` passed; the exact
prepared-only consumer tracker assertion passed; and `git status --short --
docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod
go.sum` returned no output.

Blocked checks:
Implementation, route audit rerun, validation rerun, connector checks,
release-candidate diagnostic, and closeout baseline are scheduled for later
Phase 108 checkpoints. Release actions, public publication, protected evidence
path writes, external contact, real credentials, consumer actions, and
stronger claims are out of scope.

Protected path status:
No protected evidence path is part of the plan. The plan forbids protected
path writes.

Consumer tracker status:
The consumer tracker is not part of the plan. The seven targets must remain in
order and `prepared`.

Claim-boundary status:
The plan explicitly forbids release readiness, compliance, adoption, consumer
acceptance, final-root readiness, hosted SaaS, SLA/uptime, production
readiness, vendor compatibility, hardware certification, production AVL
reliability, production-grade ETA quality, and real-world ETA accuracy claims.

Security/auth status:
The plan does not change auth, token handling, public exposure, private
payload handling, external contact, or route permissions.

Data/migration status:
No migration, schema change, durable state change, dependency change, public
feed contract change, or Go module change is planned.

Master review:
Approved. The smallest safe Phase 108 implementation is blocker/documentation
reconciliation plus route, validation, connector, and claim-boundary reruns,
with targeted fixes only when a reproducible bug or drift is found.

Required edits:
Patch scoped docs/copy/blocker drift if found, update the release-note blocker
matrix, run the planned checks, and record exact results.

Decision:
Proceed to checkpoint 000001 validation and commit, then checkpoint 000002
implementation.

Next checkpoint:
Phase 108 -- Checkpoint 000002: implement primary scoped work.
