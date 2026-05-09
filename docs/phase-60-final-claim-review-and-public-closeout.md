# Phase 60 -- Final Claim Review And Public Closeout

## Status

Planning accepted. Execution must add a final claim review record and local
audit tooling, then close public/status wording without strengthening any
unsupported claim.

## Goal

Decide exactly what Open Transit RT can truthfully claim from retained
evidence and the current official requirements context. Remove, bound, or
explicitly negate public wording that is not directly supported by retained
claim-specific evidence.

## Planning Summary

Phase 60 is a final truthfulness and public-closeout phase. It should add:

- a final claim review document;
- local audit tooling and mutation tests;
- status, roadmap, and handoff closeout updates;
- only narrow public-doc wording corrections where unsupported claims could be
  misunderstood.

The phase must not collect evidence, contact external parties, change
consumer statuses, or publish a launch/compliance statement.

## Scope

- Add a final claim-to-evidence table.
- Add an unsupported-claims table.
- Record the official-source context from Phase 54 and the Phase 60 spot
  check.
- Add an audit-only local script that scans public/status docs for unsupported
  positive claims, unsafe private strings, consumer tracker drift, and missing
  claim-review sections.
- Add mutation-style tests for the audit helper.
- Add Make targets for the audit/test helper.
- Review and update README, docs index/status, public messaging drafts,
  California readiness docs, compliance checklist docs, roadmap/status docs,
  backlog, open questions, and handoffs as needed.

## Non-Goals

- No evidence collection.
- No `docs/evidence` writes.
- No consumer contact.
- No portal automation.
- No consumer status changes.
- No public launch.
- No CAL-ITP/Caltrans compliance claim.
- No agency adoption, endorsement, or approval claim.
- No consumer submission, review, acceptance, ingestion, listing, display, or
  adoption claim.
- No hosted SaaS, hosted service, paid support, SLA, or uptime guarantee
  claim.
- No universal production-readiness or production multi-tenant hosting claim.
- No vendor compatibility, hardware certification, or marketplace approval
  claim.
- No production-grade ETA or real-world ETA accuracy claim.

## Allowed Claims After Current Review

- Open Transit RT is an open-source backend project for GTFS import/authoring,
  authenticated telemetry ingest, conservative matching, and GTFS Realtime
  feed publication.
- It implements technical foundations for stable feed paths, validation
  workflows, metadata, consumer packet preparation, and CAL-ITP-style
  readiness review.
- The OCI DuckDNS packet is hosted/operator pilot evidence only.
- Phase 33 proves local/pilot handling of one real public static GTFS dataset
  only.
- Seven consumer and aggregator packet drafts are `prepared` only.
- Phase 59 real pilot closeout is blocker-only because required retained
  pilot authorization and public-safe closeout artifacts are absent.

## Unsupported Claims

The following must either be absent or explicitly negated/bounded:

- Open Transit RT is CAL-ITP/Caltrans compliant.
- Open Transit RT has publicly launched.
- An agency adopted, endorsed, approved, or deployed Open Transit RT.
- An agency-owned or agency-approved final public root is ready or proven.
- Any consumer or aggregator target has a submitted, under-review, accepted,
  rejected, blocked, listed, displayed, ingested, or adopted status.
- Open Transit RT is hosted SaaS, a hosted service offering, paid support, or
  SLA-backed.
- Open Transit RT guarantees uptime or proves production readiness.
- Open Transit RT proves production multi-tenant hosting readiness.
- Open Transit RT is vendor compatible, hardware certified, or marketplace
  approved.
- Open Transit RT proves production-grade ETA quality or real-world ETA
  accuracy.

## Official Requirements Context

Phase 54 re-checked current public Caltrans / Cal-ITP and FTA sources on
May 9, 2026 and updated repository mappings. Phase 60 planning also spot
checked the same official pages on May 9, 2026:

- Caltrans California Transit Data Guidelines page showed Version 4.0 dated
  December 11, 2024.
- Caltrans FAQ page showed Version 4.0.
- Caltrans Website Model Language still included technical-contact guidance.
- FTA 2025 NTD Reporting Policy Manual page showed last updated
  Wednesday, April 15, 2026.

This is requirements context only. It is not deployment evidence and does not
prove compliance for this repository or any deployment.

## Files Likely To Change

- `docs/phase-60-final-claim-review-and-public-closeout.md`
- `scripts/audit-final-claim-review.sh`
- `scripts/test-final-claim-review.sh`
- `Makefile`
- `README.md`
- `docs/README.md`
- `docs/roadmap-status.md`
- `docs/public-launch-checklist.md`
- `docs/public-share-copy.md`
- `docs/california-readiness-summary.md`
- `docs/compliance-evidence-checklist.md`
- `docs/current-status.md`
- `docs/backlog.md`
- `docs/open-questions.md`
- `docs/roadmap-to-calitp-compliance-and-gap-closure.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-60.md`

Execution must not change:

- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/`
- `docs/evidence/consumer-submissions/artifacts/`
- `docs/evidence/consumer-submissions/packets/`
- `docs/evidence/captured/`

## Safety Boundaries

The audit helper must be local-only. It must not browse the web, fetch live
feeds, contact consumers, contact agencies, contact vendors, send email,
trigger webhooks, automate portals, create evidence, or change repo state
outside explicit generated test fixtures under ignored temporary locations.

The audit should scan bounded text files only and skip binary, archive,
protobuf, database dump, raw log, and oversized files. It must keep evidence
paths read-only and must treat `docs/evidence/captured` and consumer artifacts
as preservation targets, not output paths.

## Evidence And Claim Boundaries

Every retained claim in the final review document must be mapped to one of:

- implementation capability;
- local demo/prototype evidence;
- hosted/operator OCI pilot evidence;
- public-GTFS local/pilot evidence;
- prepared-only consumer packet state;
- official requirements context;
- missing/blocked evidence.

No claim may imply final-root proof, compliance, agency adoption, consumer
acceptance, hosted service availability, production readiness, SLA/uptime,
vendor compatibility, marketplace approval, public launch, or production-grade
ETA quality unless retained claim-specific evidence exists. No such stronger
evidence is available in the current repository state.

## Implementation Details

1. Add the final claim review document with a claim-to-evidence table,
   unsupported-claims table, official-source context, retained-evidence
   boundary, and maintainer signoff placeholder.
2. Add `scripts/audit-final-claim-review.sh` to:
   - verify required Phase 60 document sections;
   - verify required public boundary phrases;
   - reject unsupported positive claim wording in public/status docs;
   - reject unsafe private strings in reviewed public docs;
   - verify the seven-target prepared-only consumer tracker;
   - verify consumer artifact directories contain no non-README files.
3. Add `scripts/test-final-claim-review.sh` with local mutation tests for
   unsupported claims, unsafe strings, tracker drift, missing sections, and a
   passing fixture.
4. Add `make audit-final-claim-review` and `make test-final-claim-review`.
5. Update public/status docs only where clearer wording is needed.
6. Add `docs/handoffs/phase-60.md`, update `docs/handoffs/latest.md`, and
   close the roadmap/status docs.

## Tests

- `sh -n scripts/audit-final-claim-review.sh scripts/test-final-claim-review.sh`
- `./scripts/test-final-claim-review.sh`
- `make audit-final-claim-review`
- `make validate`
- `make test`
- `make smoke`
- required consumer/evidence preservation checks
- `docker compose -f deploy/docker-compose.yml config`
- `INTEGRATION_TESTS=1 make test-integration` when the local DB is available

## Performance And Scale Tests

No runtime performance or scale test is relevant. The audit helper should use
bounded text reads and avoid large or binary artifacts so it remains suitable
for local validation.

## Required Verification Commands

Run from the repository root:

```sh
sh -n scripts/audit-final-claim-review.sh scripts/test-final-claim-review.sh
./scripts/test-final-claim-review.sh
make audit-final-claim-review
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

The artifact `find` command must print no files.

## Master Plan Review

Accepted. The plan has clear scope, non-goals, safety boundaries, tests,
validation commands, docs/handoff updates, consumer tracker preservation
checks, and no evidence-writing or claim-upgrade path.
