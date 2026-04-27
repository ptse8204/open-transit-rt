# Consumer Submission Workflow

This workflow is for operators turning prepared packet drafts into real external
submission evidence.

It does not submit feeds, automate portals, scrape external services, verify
private portal state, or prove consumer acceptance. All tracked targets remain
`prepared` unless retained redacted evidence supports a later status.

## Targets

Use this workflow for:

- Google Maps;
- Apple Maps;
- Transit App;
- Bing Maps;
- Moovit;
- Mobility Database;
- transit.land.

## Official Submission Path Verification

Do not guess submission paths.

Before submitting to a target, an operator must verify the current official
submission or contact path using one of:

- a current official page controlled by the named target;
- target-originated correspondence;
- a current target-owned portal page visible to the authorized operator;
- retained target documentation supplied through an official support channel.

Record the verification in the target packet or current record only after
redaction review. The record must include:

- target name;
- official URL or contact path, if public;
- verification date;
- operator;
- source type;
- retained evidence path or note explaining why the source is public and does
  not need a private artifact;
- redaction notes.

If the verified source is private, do not commit a private URL, private ticket
link, personal data, portal credentials, session identifiers, or screenshots
with account data. Commit only a redacted summary or redacted artifact that
preserves the relevant fact.

## Pre-Submission Checklist

Before any real submission, confirm every item below.

| Check | Required result |
| --- | --- |
| Agency identity | The submitting operator is authorized to represent the agency or pilot owner. |
| Agency permission | Written or otherwise retained permission exists for submitting the feed set. |
| Final public feed root | The operator has confirmed whether the submission uses the DuckDNS OCI pilot root or a final agency-owned root. |
| License/contact metadata | `feeds.json` and packet metadata contain agency-approved open-license and technical-contact values. |
| Validator records | Current canonical validator records exist for schedule, Vehicle Positions, Trip Updates, and Alerts for the submitted URL root. |
| Feed reachability | All five public URLs are reachable without login: `feeds.json`, schedule ZIP, Vehicle Positions, Trip Updates, and Alerts. |
| Phase 19 wording | Replay metrics are described only as measurement evidence, not production-grade ETA proof. |
| Redactions | No secrets, portal credentials, private ticket links, personal data, or private operator artifacts are committed. |
| Target requirements | The operator has reviewed target-specific requirements from an official source. |
| Evidence storage | The operator has identified `docs/evidence/consumer-submissions/artifacts/<target>/` or a private non-repo location for retained proof. |

If any item is not satisfied, do not submit until it is resolved or record a
`blocked` status only with evidence of the blocker.

## Evidence Intake And Status Transitions

Status changes must be evidence-backed and target-specific. Validator success,
public fetch proof, prepared packets, and internal `consumer_ingestion` records
are supporting context only; they are not third-party submission or acceptance.

| Transition | Required evidence before changing status |
| --- | --- |
| `prepared -> submitted` | Redacted target receipt, ticket, confirmation email, portal screenshot, or operator-retained submission record showing the named target, submitted date, feed root, and submitted feed URLs. |
| `submitted -> under_review` | Target-originated acknowledgement that review is in progress, such as an email, ticket update, portal state, or support response. |
| `submitted -> accepted` | Target-originated acceptance confirmation naming the feed scope, URL root, date, and any conditions. |
| `under_review -> accepted` | Target-originated acceptance confirmation naming the feed scope, URL root, date, and any conditions. |
| `submitted -> rejected` | Target-originated rejection, change request, or failed review result with the rejection reason. |
| `under_review -> rejected` | Target-originated rejection, change request, or failed review result with the rejection reason. |
| `any status -> blocked` | A blocker note with source, date, operator, affected target, required next action, and whether the blocker comes from the target, agency, deployment, validation, domain, license, or redaction review. |

Rejected or blocked status does not imply acceptance. Accepted status may be
claimed only for the exact target, feed types, URL root, environment, evidence
artifact, and date shown in the record.

## Updating Records

When evidence exists:

1. Store the redacted artifact in
   `docs/evidence/consumer-submissions/artifacts/<target>/` or keep the raw
   artifact private and commit only a redacted summary.
2. Update the named target current record under
   `docs/evidence/consumer-submissions/current/`.
3. Update `docs/evidence/consumer-submissions/README.md`.
4. Update `docs/evidence/consumer-submissions/status.json`.
5. Verify the human-readable tracker and `status.json` agree for target name,
   status, packet path, prepared timestamp, and evidence references.
6. Update `docs/current-status.md` and `docs/handoffs/latest.md` if the next
   operator action changed.

Do not update other targets by inference. One target's receipt, rejection, or
acceptance does not change another target's status.

## Claim Boundary

Allowed without later evidence:

- "Prepared packets exist for operator review."
- "The operator workflow describes how to verify official submission paths."
- "The OCI pilot has hosted/operator evidence for the recorded pilot URL root."

Not allowed without retained evidence:

- a current claim that any target has been submitted to;
- a current claim that any target is under review;
- a current claim that any target accepted, ingested, listed, or displays the
  feeds;
- a claim that the DuckDNS OCI pilot proves agency-owned production URL proof;
- a claim of CAL-ITP/Caltrans compliance;
- a claim of agency endorsement, hosted SaaS availability, vendor equivalence,
  or production-grade ETA quality.
