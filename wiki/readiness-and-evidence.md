# Readiness And Evidence

Open Transit RT supports deployment toward Caltrans/CAL-ITP-style transit data readiness. A specific deployment still needs real evidence before making stronger claims.

## What The Repo Can Support

- stable public GTFS and GTFS Realtime feed paths
- `feeds.json` discovery metadata
- validation workflow records
- license and technical contact metadata
- compliance scorecard snapshots
- consumer-ingestion workflow records

## Evidence Still Needed For Stronger Claims

Before claiming a deployment is compliant or consumer-ready, collect:

- public HTTPS fetch proof for each feed
- successful static GTFS validation
- successful GTFS Realtime validation for Vehicle Positions, Trip Updates, and Alerts
- complete open-license and technical-contact metadata
- deployment monitoring and operational evidence
- actual consumer submission, review, or acceptance evidence if that status is claimed

Consumer-ingestion records inside the app are workflow records. They are not third-party acceptance.

## Source Records

- [Compliance Evidence Checklist](../docs/compliance-evidence-checklist.md)
- [Consumer Submission Evidence](../docs/consumer-submission-evidence.md)
- [Consumer Submission Tracker](../docs/evidence/consumer-submissions/README.md)
- [OCI Pilot Evidence Packet](../docs/evidence/captured/oci-pilot/2026-04-24/README.md)
