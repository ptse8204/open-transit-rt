# Public Launch Checklist

Use this checklist before reviewing or publishing any public-facing Open Transit RT wording. Phase 32 produced draft launch materials only; this checklist is not evidence that an announcement was posted, an agency was contacted, a consumer was contacted, or a public launch occurred.

## Scope Check

- [ ] Material is marked as draft until reviewed.
- [ ] No social post, email, reporter contact, agency contact, consumer contact, or announcement is represented as completed unless retained evidence exists.
- [ ] No private agency data, private contacts, raw telemetry, raw logs, private ticket links, portal screenshots with private data, or private operator artifacts are included.
- [ ] No secrets, tokens, DB URLs with passwords, private keys, webhook URLs, or generated credentials are included.
- [ ] No fake, placeholder, or inferred evidence is presented as real.
- [ ] Pilot wording is separated from production wording.
- [ ] Consumer packet wording remains prepared-only.
- [ ] Agency-owned final-root proof is listed as missing.
- [ ] Redaction review follows [Redaction Policy](evidence/redaction-policy.md) and [Security Policy](../SECURITY.md).
- [ ] Internal links and referenced docs paths were checked.

## Unsupported Claim Check

- [ ] No agency endorsement or agency adoption claim.
- [ ] No consumer submission, review, acceptance, ingestion, listing, display, or adoption claim.
- [ ] No CAL-ITP/Caltrans compliance claim.
- [ ] No hosted SaaS, paid support, or SLA claim.
- [ ] No universal production readiness or production multi-tenant hosting claim.
- [ ] No marketplace/vendor equivalence claim.
- [ ] No production-grade ETA quality or real-world ETA accuracy claim.
- [ ] No certified hardware, vendor support, or production AVL reliability claim.

## No Logo Or Affiliation Rule

Do not use agency, Caltrans/CAL-ITP, consumer, vendor, validator, or standards-body logos unless retained permission exists. Do not use wording that implies affiliation, sponsorship, certification, acceptance, deployment approval, or endorsement unless retained evidence supports that exact claim.

## Claim-To-Evidence Table

| Claim | Allowed wording | Evidence source | Forbidden overclaim |
| --- | --- | --- | --- |
| Local demo works | The repo has local demo tooling and `make agency-app-up` / `make demo-agency-flow` can exercise the committed local app flow when checks pass. | [README](../README.md), [Agency First Run](tutorials/agency-first-run.md), [Agency Demo Flow](tutorials/agency-demo-flow.md), recorded command results in the current handoff. | The demo proves production deployment, agency adoption, consumer acceptance, compliance, or hosted SaaS availability. |
| OCI pilot evidence exists | The repo has OCI DuckDNS pilot evidence for the captured pilot scope. | [OCI Pilot Evidence Packet](evidence/captured/oci-pilot/2026-04-24/README.md), [Compliance Evidence Checklist](compliance-evidence-checklist.md). | The OCI host is an agency-owned final feed root, production deployment proof, or consumer-ready final public endpoint. |
| Consumer packets are prepared | Seven consumer and aggregator packet drafts are prepared for operator review. | [Consumer Status JSON](evidence/consumer-submissions/status.json), [Consumer Packet Index](evidence/consumer-submissions/packets/README.md). | Packets were submitted, are under review, were accepted, were rejected, are blocked, are ingested, are listed, are displayed, or prove adoption. |
| AVL/vendor adapter is synthetic dry-run only | The AVL/vendor adapter pilot transforms synthetic fixtures in dry-run mode only. | [Device And AVL Integration](tutorials/device-avl-integration.md), [AVL/vendor fixtures](../testdata/avl-vendor/), Phase 29B handoff. | The project has certified hardware support, production vendor compatibility, vendor endorsement, or production AVL reliability evidence. |
| External predictor adapter was evaluated | Phase 29A evaluated the external predictor adapter boundary and candidate-style external predictor fit. | [Phase 29A External Predictor Adapter Evaluation](phase-29a-external-predictor-adapter-evaluation.md), Phase 29A handoff. | The project runs TheTransitClock in production, proves production-grade ETA quality, or is endorsed by an external predictor project. |
| Agency pilot package exists | The repo has agency pilot docs for evaluation planning. | [Agency Pilot Program](agency-pilot-program.md), [Agency Pilot Checklist](agency-pilot-checklist.md), [Agency Feedback Template](agency-feedback-template.md). | An agency adopted the project, endorsed it, completed a pilot, or approved production use. |
| Agency-owned final root is missing | No agency-owned or agency-approved final public feed root exists in repo evidence. | [Agency-Owned Domain Readiness](agency-owned-domain-readiness.md), [California Readiness Summary](california-readiness-summary.md), [Current Status](current-status.md). | The OCI DuckDNS host is agency-owned final-root proof or enough for final consumer submission claims. |
| Consumer acceptance is missing | No consumer or aggregator has submitted, under-review, accepted, rejected, blocked, ingestion, listing, display, or adoption evidence. | [Consumer Status JSON](evidence/consumer-submissions/status.json), [Consumer Submission Workflow](evidence/consumer-submissions/submission-workflow.md), [Consumer Submission Evidence](consumer-submission-evidence.md). | Any consumer has accepted, ingested, listed, displayed, rejected, blocked, or adopted the feeds. |

## Final Review Questions

- Would a non-technical reader understand that a GitHub star is only a bookmark or support signal, not an endorsement?
- Would an agency understand what evidence is missing before stronger claims?
- Would a contributor understand where to help without seeing private data?
- Would a consumer or validator name be read as affiliation, certification, or acceptance?
- Are all dates, statuses, and evidence references current and source-backed?
