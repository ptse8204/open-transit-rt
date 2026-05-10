# Evidence Track Router

Stage 8 adds a router for optional evidence work. It is documentation only.

Core rule: no intake, no evidence phase.

This router does not start an evidence phase, collect artifacts, run evidence
tools, contact external parties, automate portals, edit consumer status,
mutate retained packets, or create a claim. It only defines which optional
track may be used after a maintainer approves a complete intake.

The router does not support claims of CAL-ITP/Caltrans compliance, consumer
acceptance, agency adoption or approval, final-root proof, hosted SaaS, paid
support or SLA coverage, production readiness, production multi-tenant
hosting, vendor compatibility, hardware certification, public launch,
production AVL reliability, or production-grade ETA quality.

## Intake Requirement

Every optional evidence track requires a written intake before any collection,
tool execution, external fetch, portal action, target contact, status update,
or retained artifact write.

The intake must define:

- Authorization: who approved the work, the exact scope, the date, and whether
  any external fetch, target contact, portal use, or retained write is allowed.
- Public-safe retention: where artifacts may be stored, what must remain
  private, whether `.cache`-only diagnostics are required, and which files may
  be committed.
- Redaction: the required removal or masking for secrets, tokens, private
  hostnames, private paths, raw logs, account identifiers, device identifiers,
  personal data, correspondence, screenshots, and operational details.
- Exact claim target: the narrow claim being evaluated and the stronger claims
  that remain explicitly unsupported.
- Required artifacts: the minimum public-safe retained artifacts needed before
  the claim target can be advanced.
- Allowed tools: the exact scripts, commands, validators, browsers, portals,
  accounts, manual workflows, and network actions that are permitted.
- Stop conditions: missing authorization, unclear representation authority,
  sensitive data exposure, ambiguous claim scope, unsafe retention, missing
  redaction review, unexpected external-party workflow, failed validation, or
  any request to broaden the claim without a new intake.

If any intake field is missing, unclear, stale, or unsafe, stop. Repository
access, prior discussion, local diagnostics, prepared packets, or validator
success do not imply authorization.

## Track C: Final Public Root And Authorization

Use Track C only for a specific final public feed root and a specific
agency/operator authorization target.

The intake must name the root, the authorized operator or agency role, allowed
fetch scope, retention path, redaction requirements, and the exact final-root
claim under review.

Required artifacts may include public-safe authorization records, DNS/TLS and
redirect summaries, five-feed public fetch summaries, final-root validator
summaries, source-of-truth page references, packet README files, inventories,
and checksums. The intake decides which artifacts are required.

Track C does not prove final-root readiness, agency approval, public launch,
consumer acceptance, CAL-ITP/Caltrans compliance, production readiness, or
hosted service availability unless separate retained artifacts support that
exact claim after review.

## Track D: Consumer Or Aggregator Status

Use Track D only for a named consumer or aggregator target and a specific
target-originated status question.

The intake must name the target, the allowed official path, representation
authority, whether contact or portal use is allowed, what status transition is
being evaluated, and what retained target-originated artifact is required.

Required artifacts may include target instructions, public-safe submission
receipts, target-originated emails or dashboard states, redacted ticket records,
or blocker messages. Operator notes alone are not target-originated status.

Track D does not authorize consumer submissions, portal automation, external
contact, status mutation, or any change to
`docs/evidence/consumer-submissions/status.json`. Target statuses remain
`prepared` unless a separately approved transition is backed by retained
target-originated evidence for that exact target and status.

## Track E: Real Agency Pilot Closeout And Feedback

Use Track E only for a real agency or operator pilot closeout with explicit
authorization to retain public-safe feedback artifacts.

The intake must define the pilot scope, dates, authorized participants, what
feedback may be retained, what must remain private, and the exact pilot
closeout claim under review.

Required artifacts may include kickoff authorization, scoped pilot notes,
operator feedback, issue summaries, continue/pause/close decisions, and
redaction-reviewed closeout summaries. Private correspondence, personal data,
and operationally sensitive details must be excluded or redacted.

Track E does not prove agency adoption, endorsement, approval, public launch,
production readiness, hosted service availability, consumer acceptance, or
CAL-ITP/Caltrans compliance.

## Track F: Device, AVL, Or Vendor Integration

Use Track F only for a specific telemetry source, device class, AVL system, or
vendor integration boundary with authorization to review public-safe artifacts.

The intake must define the integration boundary, sample retention rules, real
credential handling, private payload handling, allowed adapters or sidecars,
test environment, and the exact compatibility or reliability question under
review.

Required artifacts may include redacted manifest records, mapping summaries,
synthetic-or-approved sample descriptions, send-mode diagnostics, failure-mode
reviews, and operator signoff on what may be public. Real credentials, raw
private payloads, device secrets, private vehicle identifiers, and vendor
account details must not be committed.

Track F does not prove vendor compatibility, hardware certification, field
readiness, production AVL reliability, production readiness, consumer
acceptance, or production-grade ETA quality.

## Track G: Production Operations Reliability

Use Track G only for a specific deployment and a specific operations reliability
claim target.

The intake must define the deployment scope, retention window, public-safe
summaries, backup and restore artifact rules, alerting and incident records
allowed for review, and the exact reliability claim being evaluated.

Required artifacts may include redacted uptime/freshness summaries, backup and
restore drill records, incident closeout summaries, feed monitoring summaries,
validator-health summaries, support-bundle inventories, and operational
runbook records. Raw access logs, private hostnames, credentials, internal IP
inventories, and alert destinations must be excluded or redacted.

Track G does not prove production readiness, SLA coverage, uptime guarantees,
hosted SaaS availability, paid support, production multi-tenant hosting,
consumer acceptance, or CAL-ITP/Caltrans compliance.

## Track H: Real-World Realtime And ETA Quality

Use Track H only for real-world realtime quality, observed-arrival comparison,
or ETA-quality review with approved inputs and retention rules.

The intake must define the data source, observation window, routes or trips,
privacy treatment, comparison method, prediction adapter, metrics, fixtures,
and exact quality claim under review.

Required artifacts may include public-safe aggregate metrics, fixture
inventories, methodology notes, validator context, withheld-case summaries,
coverage summaries, and limitation notes. Raw rider data, personal data, raw
private telemetry, private vehicle identifiers, and unredacted operational logs
must not be committed.

Track H does not prove production-grade ETA quality, broad ETA accuracy,
production AVL reliability, consumer acceptance, production readiness, or
CAL-ITP/Caltrans compliance.

## Routing Rules

Route to the narrowest track that matches the exact claim target. If a request
spans multiple tracks, each track needs its own intake and stop conditions.

Evidence from one track must not be reused to imply a broader claim in another
track. Validator success, public fetchability, prepared consumer packets,
local diagnostics, and `.cache` summaries are supporting signals only. They are
not compliance, correctness, acceptance, production, reliability, or adoption
proof by themselves.

If the claim target is too broad, rewrite it as a narrower evidence question or
stop before collecting artifacts.
