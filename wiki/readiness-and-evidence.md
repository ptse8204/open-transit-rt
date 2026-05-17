# Readiness And Evidence

Open Transit RT helps with the technical foundations for Caltrans/CAL-ITP-style transit data readiness. Stronger claims need real deployment evidence and, where relevant, third-party consumer evidence.

For agency-facing wording without compliance jargon, start with
[CAL-ITP Readiness Plain English](calitp-readiness-plain-english.md).

For the browser path, open `/admin/operations/readiness` after following
[Small Agency Quick Start](small-agency-quick-start.md). The UI shows readiness
signals and missing evidence separately. The workflow map covers public feed
URLs, static GTFS, Vehicle Positions, Trip Updates, Alerts, validation,
license/contact metadata, uptime and operations signals, telemetry/device
state, and consumer preparedness.

## What Open Transit RT Can Support

- stable public GTFS and GTFS Realtime feed paths
- `feeds.json` discovery metadata
- private Operations Console feed health and readiness rows
- validation workflow records
- license and technical contact metadata
- deployment-specific scorecard snapshots
- consumer-ingestion workflow records

## Evidence Needed For Stronger Claims

Before claiming a deployment is compliant or consumer-ready, collect:

- public HTTPS fetch proof for each feed
- successful static GTFS validation
- successful GTFS Realtime validation for Vehicle Positions, Trip Updates, and Alerts
- complete open-license and technical-contact metadata
- deployment monitoring and operational evidence
- actual consumer submission, review, or acceptance evidence if that status is claimed

Consumer-ingestion records inside the app are workflow records. They are not third-party acceptance.

Formal agency approval, final feed-root evidence, and consumer acceptance are
not required to use or improve the software; they are future evidence
milestones only for agencies that choose public launch or compliance claims.

## Detailed Evidence Links

- [Compliance Evidence Checklist](../docs/compliance-evidence-checklist.md)
- [Consumer Submission Evidence](../docs/consumer-submission-evidence.md)
- [Consumer Submission Tracker](../docs/evidence/consumer-submissions/README.md)
- [OCI Pilot Evidence Packet](../docs/evidence/captured/oci-pilot/2026-04-24/README.md)

## Next Steps

- [Can My Agency Use This?](can-my-agency-use-this.md)
- [Review deployment guidance](deployment-guide.md)
- [Support or contribute](support-and-contribute.md)
