# Optional Evidence Gate Refresh -- Phase 131

Status date: 2026-05-16

This refresh records the current blocker-only status of optional evidence gates.
It is not evidence, does not collect evidence, does not contact external
parties, does not fetch final roots, does not move consumer statuses, and does
not authorize protected-path writes.

## Current Conclusion

`blocked_no_authorized_evidence_collection`

All optional evidence gates reviewed in Phase 131 remain blocked unless a
future retained maintainer-approved intake authorizes the exact action,
retention path, redaction plan, status rule, validation plan, and claim target.
Current repository software, release-candidate publication, install confidence,
local validators, synthetic fixtures, and prepared packets do not by themselves
prove final-root readiness, consumer submission, consumer acceptance, agency
adoption, CAL-ITP/Caltrans compliance, hosted service availability,
SLA/uptime, vendor compatibility, hardware certification, production AVL
reliability, production readiness, production-grade ETA quality, or real-world
ETA accuracy.

## Source Documents Reviewed

- `docs/future-evidence-intake-gate-pack.md`
- `docs/handoffs/phase-109.md`
- `docs/evidence/redaction-policy.md`
- `docs/evidence/consumer-submissions/README.md`
- `docs/evidence/consumer-submissions/submission-workflow.md`
- `docs/agency-owned-domain-readiness.md`
- `docs/requirements-calitp-compliance.md`
- `docs/support-boundaries.md`
- `docs/adoption/evaluator-and-contributor-kit.md`
- `docs/release-status-v0.1.0-rc.1.md`

## Gate Matrix

| Gate | Current status | Missing preconditions | Safe next action |
| --- | --- | --- | --- |
| Final-root evidence | `blocked` | No retained agency-owned or agency-approved final public feed root, representation authority, allowed fetch scope, DNS/TLS/fetch/validator evidence, retention path, redaction plan, or rollback plan is available. | Prepare a future intake naming the root, approving operator, allowed fetches, retained artifacts, redaction rules, and exact final-root claim before any collection. |
| Consumer or aggregator submission | `blocked` | All seven targets remain `prepared`; no verified official target path, representation authority, retained target-originated receipt, portal artifact, correspondence, review acknowledgement, acceptance, rejection, or blocker artifact exists for a status transition. | Keep packets prepared-only until a future target-specific intake authorizes official-path verification, contact or portal use, artifact retention, and any exact status change. |
| Real agency pilot | `blocked` | No retained kickoff authorization, participant consent, representation authority, approved public-safe feedback scope, dates, private-data exclusions, or closeout criteria exists for a new real pilot evidence track. | Use public-safe evaluator feedback templates only; open a future pilot intake before contacting agency staff or retaining pilot artifacts. |
| Real vendor/device AVL | `blocked` | No named vendor/device source, real credential handling plan, payload retention approval, vendor correspondence authority, redaction plan, or real conformance criterion exists. | Continue using synthetic/local connector fixtures; open a future vendor/device intake before using real credentials, real payloads, or vendor contact. |
| Real-world ETA-quality study | `blocked` | No authorized observed-arrival dataset, routes/trips/dates scope, sampling plan, privacy/redaction treatment, raw-record handling plan, or limited claim target exists. | Continue synthetic backtests and fixture-based QA; open a future ETA study intake before collecting observed arrival or raw telemetry records. |
| Compliance packet | `blocked` | No exact reviewer question or compliance framework intake, final-root evidence, validator-clean final-root records, source-of-truth/license proof, consumer evidence, reviewer signoff, retained packet inventory, or claim target exists. | Keep compliance language as requirements mapping; open a future compliance intake after final-root, validator, license, and consumer evidence prerequisites are available. |
| Hosted operations, paid support, SLA, or production readiness | `blocked` | No hosted-service authorization, service owner, paid support agreement, response targets, uptime measurements, incident/backup/restore evidence, production operations record, multi-tenant proof, or service commitment exists. | Continue community-only support and self-hosted evaluation docs; require separate service packaging and operations evidence before any hosted/SLA/production claim. |

## Protected Path Status

Phase 131 did not edit, generate, reformat, stage, or write files under:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/**`
- `docs/evidence/consumer-submissions/artifacts/**`
- `docs/evidence/consumer-submissions/packets/**`

## Consumer Tracker Status

The consumer tracker remains protected and unchanged. The required status is:

| Target | Required status |
| --- | --- |
| Google Maps | `prepared` |
| Apple Maps | `prepared` |
| Transit App | `prepared` |
| Bing Maps | `prepared` |
| Moovit | `prepared` |
| Mobility Database | `prepared` |
| transit.land | `prepared` |

Prepared packets are not submissions. A future status transition requires
target-originated retained evidence for the exact target, feed scope, URL root,
and date.

## Claim Boundary

Allowed current wording:

- Public `v0.1.0-rc.1` release candidate for local/self-hosted evaluation.
- Public fresh-clone rc1 install confidence passed for the Phase 117 trial.
- Optional evidence gates are documented and currently blocked.
- Prepared consumer packets exist for operator review.
- Synthetic/local connector, conformance, and realtime QA fixtures exist.

Forbidden without future retained evidence:

- stable release readiness;
- final-root readiness;
- consumer submission, review, acceptance, ingestion, listing, or display;
- agency adoption, approval, endorsement, or public launch;
- CAL-ITP/Caltrans compliance;
- hosted SaaS or hosted-service availability;
- paid support, SLA coverage, or uptime guarantees;
- production readiness or production multi-tenant readiness;
- vendor compatibility, hardware certification, or production AVL reliability;
- production-grade ETA quality or broad real-world ETA accuracy.

## Security And Redaction Notes

Future evidence work must stop before using or retaining credentials, bearer
tokens, JWTs, API keys, DNS provider tokens, database passwords, webhook URLs,
notification credentials, device tokens, private keys, private certificates,
admin URLs with embedded secrets, private correspondence, portal credentials,
private ticket links, personal data, private payloads, raw access logs, private
origin URLs, or unredacted operational details.

## Data And Migration Impact

Phase 131 is documentation-only. It adds no migration, schema change, durable
state, dependency, Go module change, route, feed contract, runtime behavior, or
release artifact.

## Decision

Keep every optional evidence gate blocked until a future retained intake makes
one narrow action safe. Continue to Phase 132 final roadmap closeout without
collecting evidence or changing protected consumer state.
