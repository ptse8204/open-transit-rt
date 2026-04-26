# Consumer Submission Evidence

This is the Phase 13 evidence layer for downstream consumer and aggregator submissions.

It tracks what has actually been prepared, submitted, reviewed, accepted, rejected, or blocked for named consumers. It does not submit feeds, automate portal workflows, or create acceptance evidence.

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
| `prepared` | Packet prepared only, no submission. | "A submission packet has been prepared, but not submitted." |
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
- next operator action for every target

## Required Evidence Fields Per Target

Each current target record and reusable template must include:

- target name
- status
- status effective timestamp
- operator
- feed root submitted
- exact feed URLs submitted
- submission packet artifact
- validation evidence reference
- Phase 12 evidence packet reference
- correspondence / receipt / ticket / portal screenshot reference, if any
- redaction notes
- next action
- allowed public wording

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
6. Update handoff/status docs when the next action changes.

## Current Phase 13 Result

As of the initial Phase 13 tracker creation, no redacted third-party submission, review, acceptance, rejection, or blocker evidence is present in the repository for the seven tracked targets. All current records therefore start as `not_started`.
