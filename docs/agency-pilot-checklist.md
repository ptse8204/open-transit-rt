# Agency Pilot Checklist

Use this checklist before and during an Open Transit RT agency pilot. It is an
operator planning tool, not evidence of public launch, consumer acceptance,
CAL-ITP/Caltrans compliance, hosted SaaS availability, paid support, agency
endorsement, marketplace/vendor equivalence, or production readiness.

## Agency Data Prerequisites

| Check | Required result |
| --- | --- |
| Agency owner | A named agency/operator pilot owner is responsible for decisions and follow-up. |
| GTFS source | The agency can identify the GTFS ZIP or GTFS authoring path to evaluate. |
| Permission | The agency or feed owner permits use of the GTFS source for the pilot. |
| Public-safe handling | The team knows whether the GTFS may be committed, kept private, or summarized only. |
| Timezone | Agency timezone and agency-local service day expectations are known. |
| Service scope | Routes, stops, service dates, after-midnight service, frequency service, and blocks are reviewed when relevant. |

## GTFS Ownership And Permission

- Record who owns or approves the GTFS source.
- Confirm whether the source is already public or agency-provided.
- Do not commit private GTFS unless explicitly reviewed and approved as
  public-safe.
- Use `docs/tutorials/real-agency-gtfs-onboarding.md` for import review.

## License And Contact Metadata

Confirm agency-approved values for:

- agency name;
- agency URL;
- technical contact email;
- license name;
- license URL;
- feed root or planned feed root;
- publication environment label.

Demo metadata is not agency approval.

## Domain, DNS, And TLS Plan

| Check | Required result |
| --- | --- |
| Local-only pilot | Localhost URLs are labeled demo/evaluation only. |
| Pilot root | Pilot URLs are labeled hosted/operator or pilot evidence only. |
| Agency-owned root | Agency-owned or agency-approved root has an approval path before final-root proof is claimed. |
| TLS | HTTPS setup and certificate handling are owned by the deployment operator. |
| Public fetch | Anonymous public fetch proof is collected only for public-safe roots. |

Use `docs/agency-owned-domain-readiness.md` before claiming agency-owned
production-domain proof.

## Telemetry And Device Plan

Choose one current pilot path:

- no telemetry path yet, documented as blocked;
- simulator path using synthetic local values;
- simple device path using deployment-owned tokens;
- vendor/AVL adapter review using private or redacted mapping;
- Phase 29B synthetic dry-run adapter demonstration.

Do not commit device tokens, vendor credentials, private device IDs, private
vehicle IDs, raw telemetry, private AVL payloads, or private logs.

## Validator Tooling

- Run `make validators-install validators-check` where practical.
- Review static GTFS and GTFS Realtime validation results.
- Treat missing validators as a blocker or operations finding, not a pass.
- Do not treat validator success as consumer acceptance or compliance proof.

## Operations And Backup Plan

Before a pilot continues beyond local demo:

- name the operations owner;
- identify database owner and backup location outside public docs;
- define validation cadence;
- define feed monitoring cadence;
- define incident response owner;
- review restore-drill expectations;
- keep private backup paths, DB URLs, credentials, and raw logs out of public
  evidence.

Use `docs/runbooks/production-operations-hardening.md` for the operations
cadence.

## Security And Redaction Plan

Review:

- `SECURITY.md`;
- `docs/evidence/redaction-policy.md`;
- `docs/support-boundaries.md`.

Do not commit:

- tokens;
- credentials;
- DB URLs with passwords;
- private keys;
- private contacts;
- private contracts;
- raw telemetry payloads;
- private logs;
- private portal screenshots;
- private ticket links;
- private operator artifacts;
- `.cache` files.

## Consumer Submission Plan

All seven consumer and aggregator targets remain `prepared` unless retained,
redacted, target-originated evidence supports a target-specific status change.

Before any real submission:

- confirm agency/operator authorization;
- verify the official target path from a current official source;
- confirm final feed root and metadata;
- retain redacted evidence;
- update only the target supported by evidence.

Use `docs/evidence/consumer-submissions/submission-workflow.md`.

## Staff And Operator Roles

| Role | Needed for | Required before |
| --- | --- | --- |
| Pilot owner | Scope, decisions, follow-up, closeout. | Kickoff. |
| GTFS/data owner | GTFS source, permission, metadata, validation triage. | Real GTFS import. |
| Operations/IT owner | Hosting, DNS, TLS, database, backups, monitoring, secrets. | Any hosted or production-directed pilot. |
| Device/AVL owner | Device tokens, hardware, vendor mapping, simulator path. | Telemetry review. |
| Public communications owner | Public-safe wording and approval. | Any public launch discussion. |
| Security/privacy reviewer | Redaction and vulnerability handling. | Evidence publication or private data review. |

## Responsibility Matrix

| Area | Agency/operator owns | Maintainer/community can help with | Out of scope |
| --- | --- | --- | --- |
| GTFS ownership | Permission, source approval, license, metadata, public-safe handling. | Review import docs, validation findings, and redaction-safe templates. | Creating fake permission, committing private GTFS, or claiming agency endorsement. |
| Domain/DNS/TLS | Domain approval, DNS, certificates, reverse proxy, public root permanence. | Point to readiness checklist and public feed path expectations. | Hosting guarantee, agency-owned proof without evidence, or managed DNS/TLS service. |
| Device/AVL credentials | Token storage, rotation, vendor credentials, device mapping, private adapter config. | Explain telemetry contract and synthetic adapter boundary. | Certified hardware support, vendor compatibility claims, or credential custody. |
| Operations/backups | Database operations, backups, restores, monitoring, alerting, incident response. | Review runbooks and reproducible operations issues. | Paid operations, SLA response, hosted SaaS availability, or tenant-safe production hosting proof. |
| Consumer submissions | Authorization, official-path verification, actual submission, retained target evidence. | Explain prepared packets and status-transition rules. | Portal contact, submission automation, guessed paths, or consumer acceptance claims. |
| Incident response | Detect, triage, mitigate, preserve private raw evidence, publish redacted summaries only. | Review public-safe bug reports and docs/runbook gaps. | Taking over deployment operations or handling private secrets in public channels. |
| Support expectations | Decide whether community support is sufficient for the pilot. | Best-effort bug review, docs clarification, fixture review, and architecture guidance. | Paid support, response targets, legal commitments, procurement commitments, or SLA coverage. |

## Launch And Readiness Review

Before discussing a pilot publicly, complete the public launch readiness checklist
below. This checklist does not approve public launch, agency endorsement,
consumer acceptance, CAL-ITP/Caltrans compliance, hosted SaaS availability, or
production readiness. It only helps decide whether public messaging is safe and
truthful.

| Check | Required result |
| --- | --- |
| Agency permission | Agency or operator approves any public mention. |
| Approved wording | Public wording is reviewed by the agency/operator and maintainer where relevant. |
| No private data | No private GTFS, private contacts, private logs, raw telemetry, credentials, or private artifacts are included. |
| Pilot wording | Messaging clearly says pilot, evaluation, demo, or prepared packet as applicable. |
| No compliance overclaim | No CAL-ITP/Caltrans compliance claim is made without supporting evidence. |
| No consumer overclaim | No submission, review, acceptance, ingestion, listing, display, or adoption claim is made without target-originated evidence. |
| Current blockers | Domain, validation, telemetry, operations, consumer, and support blockers are current. |
| Redaction review | `SECURITY.md` and redaction policy have been applied. |
| Support boundary | No paid support, SLA, hosted SaaS, or production-operation promise is implied. |

## Exit Criteria

At closeout, choose one:

- Continue: the operator owns the next concrete action and understands remaining
  evidence gaps.
- Pause: blocker is documented and no stronger claim is made.
- Close: pilot evaluation ends with a public-safe summary and claim boundaries.

Use the closeout summary in `docs/agency-pilot-program.md`.

