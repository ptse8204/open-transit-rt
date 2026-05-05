# Agency Pilot Kickoff Agenda

Use this script for the first Open Transit RT pilot meeting. It is designed to
keep the pilot practical, evidence-bounded, and safe for public documentation.

The kickoff does not authorize public launch, consumer submission, agency
endorsement claims, CAL-ITP/Caltrans compliance claims, hosted SaaS claims, paid
support/SLA promises, or production-readiness claims.

## Who Should Attend

Required:

- agency or operator pilot owner;
- GTFS/data owner or the person who can approve the GTFS source;
- operations or IT owner for local machine, hosting, DNS, TLS, secrets, and
  backups if deployment is in scope;
- maintainer/community participant if available.

Recommended when relevant:

- vehicle/device owner;
- AVL vendor or adapter owner;
- public communications owner;
- compliance or open-data reviewer;
- security or privacy reviewer.

## Prepare Before Kickoff

The agency/operator should prepare:

- a named pilot owner and backup contact;
- a GTFS source, plus permission or license notes;
- known license, license URL, agency URL, and technical contact values;
- whether the pilot will use only the local demo, a pilot root, or a future
  agency-owned or agency-approved root;
- whether telemetry will come from a simulator, simple device, vendor AVL path,
  or no device path yet;
- whether Docker with Compose support and `curl` are available for local demo
  work;
- known security or redaction constraints;
- expected pilot decision: continue, pause, or close after evaluation.

Maintainer/community participants should prepare:

- links to `docs/agency-pilot-program.md`,
  `docs/agency-pilot-checklist.md`, and
  `docs/tutorials/agency-first-run.md`;
- the support boundary from `docs/support-boundaries.md`;
- the consumer submission boundary from
  `docs/evidence/consumer-submissions/submission-workflow.md`;
- a reminder that all consumer targets remain `prepared` unless retained,
  redacted, target-originated evidence supports a later status change.

## What Not To Collect

Do not collect or paste into notes:

- secrets or credentials;
- device tokens;
- admin tokens;
- JWT or CSRF secrets;
- DB URLs with passwords;
- private keys;
- portal credentials;
- private portal screenshots;
- private ticket links;
- private GTFS unless explicitly approved as public-safe;
- raw telemetry payloads;
- private logs;
- private operator artifacts;
- `.cache` files.

If a sensitive value appears during the meeting, stop and move it to the private
operator-owned channel or secret store. Do not copy it into public docs, issues,
or evidence.

## 30-Minute Agenda

| Time | Topic | Outcome |
| --- | --- | --- |
| 0-5 min | Introductions and support boundary | Pilot owner, roles, and no-SLA/community-support boundary confirmed. |
| 5-10 min | Pilot goal and scope | Local demo, real GTFS, telemetry, operations, or evidence focus selected. |
| 10-15 min | Data and permission review | GTFS source, permission, license/contact, and redaction constraints identified. |
| 15-20 min | Demo and deployment path | Local demo prerequisites, domain plan, and operations owner reviewed. |
| 20-25 min | Device/AVL and consumer boundary | Telemetry path selected or blocked; consumer targets remain prepared-only. |
| 25-30 min | Decisions and next actions | Next operator action, maintainer/community help item, and closeout criteria recorded. |

## 60-Minute Agenda

| Time | Topic | Outcome |
| --- | --- | --- |
| 0-5 min | Introductions and meeting boundary | Roles confirmed; no secrets or private artifacts in notes. |
| 5-10 min | What Open Transit RT does and does not do | Product and claim boundaries reviewed. |
| 10-20 min | Local demo walkthrough | `make agency-app-up`, public feed paths, and Operations Console path reviewed. |
| 20-30 min | GTFS source and metadata review | Permission, license/contact, timezone, and validation expectations reviewed. |
| 30-40 min | Telemetry/device or simulator plan | Device token lifecycle, simulator option, and AVL adapter boundary reviewed. |
| 40-48 min | Operations and evidence | Backup, monitoring, validation, redaction, and evidence boundaries reviewed. |
| 48-54 min | Consumer submission boundary | Prepared packets and target-originated evidence rules reviewed. |
| 54-60 min | Decisions and follow-up | Timeline, owner, next action, blocker list, and closeout path recorded. |

## Walkthrough Topics

- `docs/tutorials/agency-first-run.md` for the local app package.
- `docs/tutorials/real-agency-gtfs-onboarding.md` for real GTFS intake.
- `docs/tutorials/device-token-lifecycle.md` for token handling.
- `docs/tutorials/device-avl-integration.md` for device and AVL boundaries.
- `docs/runbooks/production-operations-hardening.md` for operations cadence.
- `docs/compliance-evidence-checklist.md` for evidence categories.
- `docs/agency-owned-domain-readiness.md` for final-root proof.
- `docs/evidence/consumer-submissions/submission-workflow.md` for consumer
  status transitions.

## Decisions To Make

Record only public-safe decisions:

- pilot owner and backup;
- GTFS source and permission path;
- local demo, pilot root, or future agency-owned-root focus;
- telemetry source: simulator, simple device, AVL adapter, or blocked;
- validator tooling owner;
- operations owner;
- public messaging owner, if any;
- closeout decision date or checkpoint.

## Follow-Up Actions

After kickoff:

- run or schedule the local demo walkthrough;
- complete the pilot checklist;
- review GTFS source permission and metadata;
- choose the telemetry or simulator path;
- identify validation tooling and evidence storage location;
- record blockers without weakening claim boundaries;
- schedule closeout using the summary in `docs/agency-pilot-program.md`.

