# Phase 60 -- Final Claim Review And Public Closeout

## Status

Complete.

Phase 60 added the final claim review record, local audit tooling, mutation
tests, Make targets, and public/status closeout wording. The phase did not
create evidence, write `docs/evidence`, contact any external party, change
consumer status, refresh consumer artifacts, or strengthen any unsupported
claim.

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

## Final Claim Review

Phase 60 reviewed the retained repository evidence and public/status wording
against the Phase 54 official-source mapping and the Phase 60 claim boundary.
The review classifies each public claim as an implementation capability,
local/demo evidence, hosted/operator pilot evidence, public-GTFS local/pilot
evidence, prepared-only consumer packet state, official requirements context,
or missing/blocked evidence.

No retained evidence currently supports a claim of public launch,
CAL-ITP/Caltrans compliance, agency adoption, agency-owned final-root proof,
consumer submission/review/acceptance/ingestion/listing/display, hosted SaaS,
paid support, SLA/uptime, production readiness, production multi-tenant
hosting, vendor compatibility, marketplace approval, or production-grade ETA
quality.

## Claim-To-Evidence Table

| Claim area | Allowed wording after Phase 60 | Evidence classification | Boundary |
| --- | --- | --- | --- |
| Project purpose | Open Transit RT is an open-source backend project for GTFS import/authoring, authenticated telemetry ingest, conservative matching, and GTFS Realtime feed publication. | Implementation capability | This is not a hosted service, SaaS, paid support, or SLA claim. |
| Vehicle Positions first | The repo implements GTFS-RT Vehicle Positions as the first high-quality public output path and keeps Trip Updates pluggable. | Implementation capability | This does not prove production readiness or consumer ingestion. |
| Stable paths and metadata | The repo implements technical foundations for stable feed paths, validation workflows, metadata, consumer packet preparation, and CAL-ITP-style readiness review. | Implementation capability and official requirements context | This does not prove compliance for a real deployment. |
| Local app and demo flows | Local commands and fixtures can exercise GTFS import, telemetry ingest, matching, and public feed paths. | Local demo/prototype evidence | Local proof is not public launch, agency adoption, or production proof. |
| OCI DuckDNS packet | The OCI DuckDNS packet is hosted/operator pilot evidence only. | Hosted/operator OCI pilot evidence | It is not agency-owned final-root proof, compliance proof, consumer proof, or hosted SaaS availability. |
| Public-GTFS local/pilot packet | Phase 33 proves local/pilot handling of one real public static GTFS dataset only. | Public-GTFS local/pilot evidence | It does not prove official agency feed status, final-root readiness, real realtime data, or ETA quality. |
| Consumer packets | Seven consumer and aggregator packet drafts are `prepared` only. | Prepared-only consumer packet state | No target has submitted, under-review, accepted, rejected, blocked, ingested, listed, displayed, or adopted status. |
| Official requirements mapping | Phase 54 refreshed current Caltrans / Cal-ITP and FTA source mappings on May 9, 2026. | Official requirements context | Requirements context is not deployment evidence and does not prove compliance. |
| Real pilot closeout | Phase 59 real pilot closeout is blocker-only because required retained pilot authorization and public-safe closeout artifacts are absent. | Missing/blocked evidence | No real pilot evidence packet was created and no agency or operator claim was strengthened. |

## Unsupported Claims

The following remain unsupported and must either be absent from public/status
docs or explicitly negated/bounded:

| Unsupported claim | Current retained evidence state | Required before claim could change |
| --- | --- | --- |
| Open Transit RT is CAL-ITP/Caltrans compliant. | Missing. | Deployment-specific retained proof for stable public URLs, validator-clean schedule/realtime feeds, open license/contact/discoverability, and consumer acceptance where required by the claim. |
| Open Transit RT has publicly launched. | Missing. | Retained launch approval and public-safe publication record. |
| An agency adopted, endorsed, approved, or deployed Open Transit RT. | Missing. | Retained public-safe agency/operator authorization and approval artifacts for the exact claim. |
| An agency-owned or agency-approved final public root is ready or proven. | Missing. | Retained final-root approval, DNS, TLS, redirect, public fetch, validator, proxy/config, README, and checksum evidence. |
| A consumer or aggregator target has a status beyond `prepared`. | Missing. | Target-originated or operator-retained public-safe evidence for the exact target and status. |
| Open Transit RT is hosted SaaS, a hosted service offering, paid support, or SLA-backed. | Missing. | Approved service offering, support, and SLA evidence outside current repo-only scope. |
| Open Transit RT guarantees uptime or proves production readiness. | Missing. | Deployment-specific operations, monitoring, incident, backup/restore, availability, and validation evidence. |
| Open Transit RT proves production multi-tenant hosting readiness. | Missing. | Retained multi-tenant deployment and data-isolation proof for a production environment. |
| Open Transit RT is vendor compatible, hardware certified, or marketplace approved. | Missing. | Retained vendor, hardware, or marketplace-originated evidence for the exact named claim. |
| Open Transit RT proves production-grade ETA quality or real-world ETA accuracy. | Missing. | Real-world observed arrival/departure quality evidence and retained route/time-period metrics. |

## Retained Evidence Boundary

Phase 60 did not add, modify, or refresh retained evidence. It intentionally
left these paths unchanged:

- `docs/evidence/consumer-submissions/status.json`;
- `docs/evidence/consumer-submissions/current/`;
- `docs/evidence/consumer-submissions/artifacts/`;
- `docs/evidence/consumer-submissions/packets/`;
- `docs/evidence/captured/`.

The local audit verifies that the consumer tracker still has exactly seven
targets in the expected order and that each target remains `prepared`. It also
verifies that consumer artifact directories contain only README files.

## Maintainer Signoff

Maintainer review placeholder:

- Reviewer:
- Review date:
- Decision:
- Notes:

## Execution Closeout

Phase 60 execution added:

- `scripts/audit-final-claim-review.sh`;
- `scripts/test-final-claim-review.sh`;
- `make audit-final-claim-review`;
- `make test-final-claim-review`;
- validation scaffolding for the new scripts and handoff;
- final closeout notes in bounded public/status docs;
- `docs/handoffs/phase-60.md`.

The audit helper is local-only and read-only. It does not browse the web, fetch
live feeds, contact consumers, contact agencies, send email, trigger webhooks,
automate portals, create evidence, or change repo state outside ignored test
fixtures created by `scripts/test-final-claim-review.sh`.
