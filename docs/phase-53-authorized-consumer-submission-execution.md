# Phase 53 -- Authorized Consumer Submission Execution

## Status

Planning accepted. Execution must use the blocker-only branch unless real
operator authorization, official target path evidence, and target-originated
submission artifacts are supplied before implementation.

Current repository evidence does not include an operator authorization artifact,
official target path verification artifact, submitted packet artifact, or
target-originated receipt, ticket, screenshot, email, acknowledgement,
rejection, blocker, or acceptance artifact for any consumer or aggregator
target.

## Goal

Execute consumer or aggregator submission only for a selected target when an
operator is authorized and retained evidence proves the official path and
target-specific status transition.

If those artifacts are not available, close Phase 53 as blocker-documented at
the phase level without selecting a target, contacting a consumer, adding
artifacts, or changing any target status.

## Scope

- Re-check the prepared-only consumer tracker and README-only artifact
  inventory.
- Confirm whether any retained operator authorization, official path
  verification, or target-originated artifact exists locally.
- If real evidence exists, update only the named target and only to the status
  supported by retained target-originated evidence.
- If real evidence does not exist, document blocker-only closure.
- Preserve all seven consumer and aggregator targets as `prepared` unless a
  retained target-specific artifact supports a transition.

## Non-Goals

- No consumer contact.
- No portal automation or scraping.
- No live official-path browsing during blocker-only closure.
- No guessed submission path.
- No mass status changes.
- No fake receipts, tickets, screenshots, emails, acknowledgements, rejection
  notes, blocker notes, or acceptance artifacts.
- No evidence packet creation outside the already-approved consumer submission
  artifact workflow.
- No final-root, compliance, agency adoption, consumer acceptance, hosted SaaS,
  production-readiness, vendor-compatibility, marketplace, SLA/uptime, or
  production-grade ETA claim.

## Files Likely To Change

Blocker-only execution may change:

- `docs/phase-53-authorized-consumer-submission-execution.md`
- `docs/handoffs/phase-53.md`
- `docs/handoffs/latest.md`
- `docs/current-status.md`
- `docs/backlog.md`
- `docs/open-questions.md`

If, and only if, real retained target-originated evidence is supplied, a
target-specific execution may also change:

- `docs/evidence/consumer-submissions/artifacts/<target>/`
- `docs/evidence/consumer-submissions/current/<target>.md`
- `docs/evidence/consumer-submissions/README.md`
- `docs/evidence/consumer-submissions/status.json`

## Safety Boundaries

Phase 53 must not initiate contact with any consumer or aggregator. Real
submission is an operator action outside the repo.

The execution agent must not select a target, verify a private portal path,
record a submission, or update `status.json` unless retained public-safe
evidence supports the exact target, feed scope, URL root, submitted date, and
status transition.

Do not commit private URLs, personal data, raw correspondence, private portal
screenshots, credentials, cookies, Authorization headers, DB URLs, raw logs,
private telemetry, unredacted diagnostics, or placeholder artifacts.

## Evidence And Claim Boundaries

Prepared packets are not submissions. Validator success, public fetch proof,
hosted pilot evidence, final-root workflow blockers, and internal
`consumer_ingestion` records are supporting context only. They do not prove
consumer submission, review, acceptance, ingestion, listing, display,
compliance, agency endorsement, hosted SaaS availability, marketplace/vendor
equivalence, or production readiness.

Phase-level blocker-only closure must not create a target-specific `blocked`
status, because no target-specific blocker artifact exists.

## Implementation Details

Execution must follow one of two branches.

### Branch A: Evidence-Backed Target Execution

Use this branch only when all of these are present before editing:

- retained operator authorization for the selected target;
- official target-owned path verification or public target-owned path source;
- retained target-originated or operator-retained submission evidence;
- completed redaction review;
- exact feed root and submitted feed URLs;
- exact target-specific status transition supported by the artifact.

Allowed output is a one-target update only. All other target records and
statuses must remain unchanged.

### Branch B: Blocker-Only Closure

Use this branch for the current repo state.

Record that:

- no operator authorization artifact exists;
- no official path verification artifact exists;
- no target-originated artifact exists;
- no target was selected;
- no consumer or aggregator was contacted;
- no portal was automated or scraped;
- no submission path was guessed;
- no submission was recorded;
- no artifact was added;
- all seven targets remain `prepared`;
- `docs/evidence/consumer-submissions/status.json`, current target records, and
  target artifact directories remain unchanged.

Checkpoint strategy:

- `Phase 53 -- Checkpoint 000001: add authorized consumer submission plan`
- `Phase 53 -- Checkpoint 000002: close authorized consumer submission blocker`

## Tests

Required tests and checks for blocker-only closure:

- Validate `docs/evidence/consumer-submissions/status.json` JSON syntax.
- Run the exact seven-target prepared-only consumer tracker check.
- Confirm `docs/evidence/consumer-submissions/status.json` is unchanged.
- Confirm current target records and artifact directories are unchanged.
- Confirm artifact directories contain README files only.
- Run the repository validation and test suites.
- Review changed docs for unsupported submission, acceptance, ingestion,
  compliance, agency endorsement, hosted SaaS, vendor equivalence, production
  readiness, SLA/uptime, and production-grade ETA claims.

For a future evidence-backed target transition, add focused checks that exactly
one target changed, the retained artifact path exists, redaction notes are
present, and all other targets remain unchanged.

## Performance And Scale Tests

No runtime performance test is relevant for blocker-only closure.

For future retained artifacts, do not commit large raw portal captures,
correspondence dumps, or private exports. Prefer redacted summaries, bounded
public-safe artifacts, and checksums when artifact size is material.

## Docs, Status, And Handoff Updates

Blocker-only closure must update:

- this phase document;
- `docs/handoffs/phase-53.md`;
- `docs/handoffs/latest.md`;
- `docs/current-status.md`;
- `docs/backlog.md` or `docs/open-questions.md` where useful to preserve the
  next action.

The closeout must explicitly say that Phase 53 did not change
`docs/evidence/consumer-submissions/status.json`, current target records,
target artifact directories, or `docs/evidence/captured`.

## Required Verification Commands

Run and report:

```bash
make validate
make test
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
git diff --exit-code -- docs/evidence/consumer-submissions/current docs/evidence/consumer-submissions/artifacts docs/evidence/captured
find docs/evidence/consumer-submissions/artifacts -mindepth 2 -maxdepth 2 -type f ! -name README.md -print
docker compose -f deploy/docker-compose.yml config
```

The `find` command must print no files for blocker-only closure.
