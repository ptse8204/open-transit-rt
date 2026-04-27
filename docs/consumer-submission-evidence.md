# Consumer Submission Evidence

This is the Phase 13 and Phase 20 evidence layer for downstream consumer and aggregator submissions.

It tracks what has actually been prepared, submitted, reviewed, accepted, rejected, or blocked for named consumers. It does not submit feeds, automate portal workflows, guess submission paths, or create acceptance evidence.

## Scope

Targets tracked in Phase 13:

- Google Maps
- Apple Maps
- Transit App
- Bing Maps
- Moovit
- Mobility Database
- transit.land

Phase 13 links consumer submission evidence back to the OCI pilot hosted/operator evidence packet at `docs/evidence/captured/oci-pilot/2026-04-24/`.

Phase 20 adds prepared packet drafts under `docs/evidence/consumer-submissions/packets/` and a machine-readable tracker snapshot at `docs/evidence/consumer-submissions/status.json`.

## Truthfulness Boundary

Validator success and public fetch proof are supporting evidence only. They show that a deployed feed was reachable or validator-clean at a recorded time. They are not consumer acceptance, consumer ingestion, CAL-ITP compliance, marketplace listing, or vendor equivalence.

Consumer-ingestion workflow records in the Open Transit RT database are internal/operator records. They are not third-party acceptance unless the evidence packet includes retained proof from the named consumer or aggregator.

## Statuses

Each target must have exactly one current status:

- `not_started`
- `prepared`
- `submitted`
- `under_review`
- `accepted`
- `rejected`
- `blocked`

## Allowed Claims By Status

| Status | Evidence meaning | Allowed public wording |
| --- | --- | --- |
| `not_started` | No submission has been made. | "No submission has been made for this consumer." |
| `prepared` | Complete packet prepared only, no submission. | "A submission packet has been prepared, but not submitted." |
| `submitted` | Submission sent, no acceptance implied. | "A submission was sent on the recorded date; acceptance is not confirmed." |
| `under_review` | Consumer review in progress, no acceptance implied. | "The named consumer has acknowledged review is in progress; acceptance is not confirmed." |
| `accepted` | Acceptance may be claimed only for the named consumer, feed scope, URL root, and evidence date. | "Accepted by `<consumer>` for `<feed types>` at `<URL root>` as of `<acceptance date>`." |
| `rejected` | Rejection documented, no acceptance claim. | "The named consumer rejected or requested changes to the submission." |
| `blocked` | Submission blocked by named missing evidence/action. | "Submission is blocked pending `<named evidence or action>`." |

## Required Current Tracker Fields

The tracker at `docs/evidence/consumer-submissions/README.md` must include:

- tracker last reviewed timestamp
- reviewed by
- linked Phase 12 evidence packet
- current record for every target
- exact status for every target
- packet path for every prepared target
- prepared timestamp for every prepared target
- evidence references for every prepared target
- next operator action for every target

The machine-readable snapshot at `docs/evidence/consumer-submissions/status.json` must agree with the tracker for target name, status, packet path, prepared timestamp, and evidence reference values.

## Required Evidence Fields Per Target

Each current target record and reusable template must include:

- target name
- status
- status effective timestamp
- operator
- prepared timestamp, when status is `prepared`
- feed root submitted
- exact feed URLs submitted
- submission packet artifact
- validation evidence reference
- Phase 12 evidence packet reference
- correspondence / receipt / ticket / portal screenshot reference, if any
- redaction notes
- next action
- allowed public wording

## Prepared Packet Completeness Rule

A target may move to `prepared` only when its packet includes:

- all five public feed URLs: `feeds.json`, schedule, Vehicle Positions, Trip Updates, and Alerts;
- license/contact metadata;
- Phase 12 hosted evidence link;
- validator evidence link;
- redaction note;
- next action;
- allowed public wording.

If a packet is incomplete, the target must remain `not_started`.

Every packet must also include:

- Prepared at;
- Prepared by;
- Evidence snapshot;
- OCI packet reference;
- `feeds.json` snapshot reference;
- validator records reference;
- Phase 19 replay/quality summary reference;
- Submission method;
- Official submission URL/contact;
- Verified as current;
- Notes.

If the official submission path is not verified, record `not verified` instead of guessing.

Every packet must warn operators not to submit from repo docs alone and to review feed URLs, license/contact metadata, validation status, agency identity, consumer-specific requirements, and redactions before actual submission.

## Required Acceptance-Scope Fields

Acceptance fields must remain empty or `N/A` unless the named consumer or aggregator has provided retained evidence:

- accepted feed types
- accepted environment
- accepted URL root
- acceptance date
- evidence artifact
- limitations / conditions

Acceptance may be claimed only for the exact target, feed types, environment, URL root, evidence artifact, and date shown in the accepted record.

## Operator Update Process

1. Copy the target template from `docs/evidence/consumer-submissions/templates/` or update the current target record under `docs/evidence/consumer-submissions/current/`.
2. Store redacted artifacts under the appropriate evidence packet path before changing a status to `submitted`, `under_review`, `accepted`, `rejected`, or `blocked`.
3. Record the status effective timestamp and operator.
4. Keep private correspondence, portal credentials, ticket links with secrets, and personal data out of the repo unless redacted.
5. Update the tracker freshness fields in `docs/evidence/consumer-submissions/README.md`.
6. Update `docs/evidence/consumer-submissions/status.json` when tracker values change.
7. Update handoff/status docs when the next action changes.

## Current Phase 20 Result

As of the Phase 20 packet preparation pass, complete prepared packets exist for all seven tracked targets. The targets are `prepared` only; no redacted third-party submission, review, acceptance, rejection, or blocker evidence is present in the repository.
