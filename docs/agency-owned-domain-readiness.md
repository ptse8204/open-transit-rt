# Agency-Owned Domain Readiness

This checklist explains what must be proven before moving from the DuckDNS OCI
pilot root to stronger agency-owned production-domain evidence.

The current OCI pilot root, `https://open-transit-pilot.duckdns.org`, is useful
pilot evidence. It is not agency-owned production-domain proof.

## External Reference Points

Caltrans' current California Transit Data Guidelines v4.0 describe GTFS Schedule
and GTFS Realtime compliance in terms that include public stable URLs, regular
validator success with no errors, open license publication, and ingestion by
major trip planners. For GTFS Realtime completeness, the three standard feed
types are Trip Updates, Vehicle Positions, and Alerts.

The Caltrans FAQ says stable URLs should not change and should use the transit
provider's domain when possible.

Official pages:

- [California Transit Data Guidelines v4.0](https://dot.ca.gov/cal-itp/california-minimum-general-transit-feed-specification-gtfs-guidelines)
- [California Transit Data Guidelines FAQ](https://dot.ca.gov/cal-itp/california-transit-data-guidelines-faqs)

## Required Proof

Before claiming agency-owned production-domain readiness, collect evidence for:

| Area | Required evidence |
| --- | --- |
| Agency-owned URL root | Agency-controlled domain or agency-approved stable hostname, plus retained approval for use as the public feed root. |
| Final public feed URLs | Final URLs for `feeds.json`, schedule ZIP, Vehicle Positions, Trip Updates, and Alerts. |
| TLS proof | HTTPS certificate metadata, redirect behavior if HTTP is exposed, and anonymous public fetch proof. |
| Validator records | Current no-error canonical validator records for schedule and all three realtime feeds at the final root. |
| Metadata | Agency-approved license, technical contact, provider identity, and discoverability metadata at the final root. |
| Updated packets | Prepared packet drafts refreshed to use the final root, final metadata, and final validator records. |
| Consumer submission | Real submissions use the final root unless the operator documents why the pilot root is intentionally being submitted. |
| Migration plan | Redirect, communication, or resubmission plan if any previously shared URL changes. |

## Operator Checklist

1. Confirm the agency-owned root and final feed paths.
2. Confirm TLS and anonymous reachability for all five public URLs.
3. Run and retain current validator records for the final root.
4. Refresh `feeds.json` evidence and packet metadata for the final root.
5. Review redactions before committing any domain, TLS, or operator artifacts.
6. Update prepared packets only after all final-root fields are known.
7. Submit to consumers or aggregators only through verified official paths and
   only when authorized.
8. Record target-originated evidence before changing any target status.

## Claim Boundary

Allowed now:

- "The DuckDNS OCI pilot provides pilot evidence for the recorded URL root."
- "Agency-owned domain readiness requirements are documented."

Not supported yet:

- "The OCI pilot proves agency-owned production-domain compliance."
- "The final agency-owned feed root has validator-clean production evidence."
- "Consumers accepted or ingested the final agency-owned feed root."
- "Open Transit RT is CAL-ITP/Caltrans compliant."
