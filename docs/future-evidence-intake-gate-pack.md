# Future Evidence Intake Gate Pack

This pack defines authorization gates for future optional evidence work. It is
not evidence, does not collect evidence, does not contact external parties,
does not fetch final roots, does not move consumer statuses, and does not
authorize protected-path writes.

Core rule: no complete written intake, no evidence work.

## Universal Intake Fields

Every future evidence gate requires a retained maintainer-approved intake that
includes:

| Field | Requirement |
| --- | --- |
| Authorization | Approver, date, exact scope, allowed actions, forbidden actions, and stop conditions. |
| Claim target | The narrow claim under review and the stronger claims that remain unsupported. |
| Named scope | Agency, public root, consumer target, vendor/device source, ETA study scope, or compliance packet scope. |
| Representation authority | Who may represent the agency/operator and who may contact external parties, if any. |
| Allowed tools | Exact scripts, validators, browser/manual workflows, portals, accounts, and network actions allowed. |
| Retention plan | Which artifacts may be committed, which must stay local/private, approved paths, inventories, and checksum expectations. |
| Redaction plan | Required removal or masking for credentials, private URLs, raw logs, device identifiers, personal data, correspondence, screenshots, and operational details. |
| Status rules | Any tracker or status transition criteria, rollback instructions, and reviewer signoff needed before a change. |
| Validation plan | Required local checks and external checks, if explicitly authorized. |
| Stop rules | Conditions requiring immediate pause, redaction, rollback, or maintainer review. |

Missing, ambiguous, stale, or unsafe intake fields block the gate.

## Universal Stop Rules

Stop before work if any of these are true:

- The action would modify protected paths without explicit authorization.
- The action would move a consumer status beyond `prepared` without retained
  target-originated evidence and maintainer approval.
- The action requires real credentials, real private payloads, private
  correspondence, personal data, or real vendor/agency/device data not covered
  by the intake.
- The action would contact an agency, vendor, consumer, portal, map provider,
  aggregator, or external service without named authorization.
- The action would publish a tag, GitHub Release, package, image, public
  announcement, or stronger public claim.
- Redaction cannot be verified.
- The claim target is broader than the retained artifacts can support.

## Gate A -- Final-Root Intake

Purpose: decide whether a future final public feed root evidence collection
may begin.

Required intake:

- named agency/operator and public root URL;
- representation authority for the final-root claim;
- allowed fetch scope for DNS, TLS, redirects, `feeds.json`, static GTFS,
  Vehicle Positions, Trip Updates, Alerts, and validators;
- retention path and checksum/inventory requirements;
- redaction rules for hostnames, headers, logs, private origins, and operator
  notes;
- rollback plan if the root, files, or source-of-truth metadata change during
  review.

Minimum preconditions before collection:

- active published GTFS and GTFS-RT endpoints are configured in the product;
- operator confirms the root is intended for public review;
- validator tooling and public-fetch commands are listed in the intake;
- final-root claim wording is limited to the exact reviewed root.

Forbidden without separate authorization:

- fetching or validating a public root;
- writing retained root evidence;
- claiming final-root readiness, public launch, compliance, consumer
  acceptance, hosted service availability, SLA/uptime, or production readiness.

## Gate B -- Consumer Submission Intake

Purpose: decide whether a future consumer or aggregator submission workflow may
begin for one named target.

Required intake:

- named target from the prepared tracker;
- target-originated official path or instructions;
- representation authority for submission/contact;
- allowed portal, account, manual workflow, or correspondence path;
- exact target status transition under review;
- required target-originated retained artifact for any status change;
- rollback plan for tracker or packet mistakes.

Minimum preconditions before collection:

- target packet is already `prepared`;
- public feed root gate is either separately satisfied or explicitly scoped
  out by the target;
- maintainer approves whether external contact or portal use is allowed;
- artifact retention and redaction rules are explicit.

Forbidden without separate authorization:

- contacting the target;
- automating a portal;
- editing `docs/evidence/consumer-submissions/status.json`;
- writing under protected consumer-submission packet/current/artifact paths;
- claiming submission, review, ingestion, listing, display, acceptance, or
  rejection.

## Gate C -- Real Agency Pilot Intake

Purpose: decide whether a future real agency/operator pilot closeout can retain
public-safe artifacts.

Required intake:

- named agency/operator, authorized participants, and dates;
- pilot scope and out-of-scope activities;
- what feedback may be retained and what must remain private;
- issue, support, and training artifacts allowed for public-safe summaries;
- redaction rules for names, contact details, internal processes, URLs,
  screenshots, and operational notes;
- decision criteria for continue, pause, or close.

Minimum preconditions before collection:

- kickoff authorization exists;
- participant consent and representation authority are clear;
- feedback retention is scoped to public-safe summaries;
- private support or operational data is excluded.

Forbidden without separate authorization:

- contacting agency staff;
- recording real feedback;
- writing retained pilot evidence;
- claiming adoption, agency approval, endorsement, public launch, compliance,
  production readiness, hosted service availability, or SLA/uptime.

## Gate D -- Real Vendor / Device AVL Intake

Purpose: decide whether a future real telemetry source, AVL integration, device
class, or vendor review may retain public-safe artifacts.

Required intake:

- named vendor, device class, AVL source, or telemetry boundary;
- credential handling plan and explicit prohibition on committing secrets;
- private payload handling and sample retention rules;
- allowed adapters, sidecars, polling, POST, or transform paths;
- conformance criteria and failure-mode review;
- redaction rules for vehicle identifiers, device identifiers, account
  details, payloads, URLs, logs, screenshots, and vendor correspondence.

Minimum preconditions before collection:

- synthetic/local connector path is already documented;
- real credential storage and rotation responsibilities are clear;
- sample data retention is either forbidden or explicitly public-safe;
- stop conditions cover unexpected payload fields and vendor restrictions.

Forbidden without separate authorization:

- using real credentials;
- ingesting or retaining real vendor/device payloads;
- contacting a vendor;
- writing retained integration evidence;
- claiming vendor compatibility, hardware certification, production AVL
  reliability, production readiness, consumer acceptance, or ETA quality.

## Gate E -- Real-World ETA-Quality Study Intake

Purpose: decide whether a future observed-arrival, realtime quality, or
ETA-quality study may begin.

Required intake:

- exact routes, trips, dates, observation window, and data sources;
- prediction adapter or deterministic baseline under review;
- metric definitions, sampling method, exclusions, and confidence thresholds;
- privacy and redaction treatment for telemetry, vehicle IDs, riders, drivers,
  raw observations, and operational logs;
- public-safe aggregate retention plan;
- limitation and non-generalization language.

Minimum preconditions before collection:

- synthetic backtests and conformance fixtures are current;
- real observed data is authorized for the stated scope;
- raw private records remain excluded unless explicitly approved and redacted;
- the claim target is limited to the study scope.

Forbidden without separate authorization:

- collecting real-world observed-arrival data;
- retaining raw telemetry or private operational records;
- contacting operators or riders;
- claiming production-grade ETA quality, broad real-world ETA accuracy,
  production AVL reliability, consumer acceptance, compliance, or production
  readiness.

## Gate F -- Compliance Packet Intake

Purpose: decide whether a future compliance-oriented packet may collect and
retain public-safe artifacts.

Required intake:

- exact compliance framework, checklist, or reviewer question;
- claim target and exclusions;
- required artifacts and source-of-truth records;
- validator versions, public feed roots, metadata requirements, and license
  expectations;
- retention, checksum, inventory, and redaction rules;
- reviewer signoff and rollback criteria.

Minimum preconditions before collection:

- local validators and product acceptance checks are current;
- final-root and consumer gates are not implied unless separately authorized;
- artifacts are mapped to requirements without overclaiming;
- unsupported requirements are listed as blockers.

Forbidden without separate authorization:

- collecting compliance evidence;
- contacting reviewers or agencies;
- writing retained packet evidence;
- claiming CAL-ITP/Caltrans compliance, validator-clean public feeds,
  consumer acceptance, final-root readiness, production readiness, hosted
  service availability, SLA/uptime, or release readiness.

## Gate Review Output

A future gate review may output only:

- `authorized`: complete intake exists and the exact next action is safe;
- `blocked`: required intake, authority, artifact, redaction, or validation
  precondition is missing;
- `needs_review`: the intake exists but the scope, redaction, retention, or
  claim target needs maintainer review;
- `not_applicable`: the requested claim does not fit this repository boundary.

The review output must include protected-path status, consumer tracker status,
claim-boundary status, security/auth status, data/migration status, and the
next safe action.
