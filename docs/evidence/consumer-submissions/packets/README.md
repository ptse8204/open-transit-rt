# Consumer Submission Packet Index

This directory contains Phase 20 prepared packet drafts for consumer and aggregator submission review.

These packets do not submit feeds, automate portal workflows, verify consumer-specific submission paths, or prove consumer acceptance. They are reviewable operator packets only.

## Operator Warning

Do not submit from repo docs alone. Before any actual submission, the operator must review feed URLs, license/contact metadata, validation status, agency identity, consumer-specific requirements, and redactions.

Use `../submission-workflow.md` before any real submission. The workflow
explains official-path verification, pre-submission checks, evidence intake,
status transitions, and artifact storage.

## Completeness Rule

A target may move to `prepared` only when its packet includes:

- all five public feed URLs;
- license and contact metadata;
- Phase 12 hosted evidence link;
- validator evidence link;
- redaction note;
- next action;
- allowed public wording.

If any of these are missing, the target must remain `not_started`.

## Packet Completeness Checklist

| Target | All five feed URLs present | License/contact metadata present | Phase 12 evidence link present | Validator evidence link present | Redaction note present | Next action present | Allowed wording present | Resulting status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Google Maps | yes | yes | yes | yes | yes | yes | yes | `prepared` |
| Apple Maps | yes | yes | yes | yes | yes | yes | yes | `prepared` |
| Transit App | yes | yes | yes | yes | yes | yes | yes | `prepared` |
| Bing Maps | yes | yes | yes | yes | yes | yes | yes | `prepared` |
| Moovit | yes | yes | yes | yes | yes | yes | yes | `prepared` |
| Mobility Database | yes | yes | yes | yes | yes | yes | yes | `prepared` |
| transit.land | yes | yes | yes | yes | yes | yes | yes | `prepared` |

## Packet Paths

- `google-maps/README.md`
- `apple-maps/README.md`
- `transit-app/README.md`
- `bing-maps/README.md`
- `moovit/README.md`
- `mobility-database/README.md`
- `transit-land/README.md`

## Artifact Paths

Target-specific artifact intake directories live under `../artifacts/`.
They must contain only README files unless real redacted target-originated
evidence exists.
