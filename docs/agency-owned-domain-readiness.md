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

## Phase 23 Blocker Record

Phase 23 closed as blocker-documented only. No agency-owned or
agency-approved final public feed root was available or approved, so final-root
evidence was not collected.

| Area | Phase 23 status |
| --- | --- |
| Candidate final root | Blocked — not available. |
| Domain owner / approving operator | Blocked — no final root or approval artifact available. |
| DNS proof | Not collected — final root unavailable. |
| TLS proof | Not collected — final root unavailable. |
| HTTP to HTTPS redirect proof | Not collected — final root unavailable. |
| All five final feed URLs | Not collected — final root unavailable. |
| Final-root validator records | Not collected — final root unavailable. |
| Prepared packet refresh | Not performed — final-root evidence unavailable. |
| Migration / redirect proof from DuckDNS | Not collected — no final root exists to migrate to. |

## Post-Phase-32 Follow-Up Blocker Record

The post-Phase-32 final-root evidence follow-up also closed as
blocker-documented only. The selected status was `No root`: no agency-owned or
agency-approved final public feed root was available, and no operator approval
artifact was available for using a candidate root as the final public feed
root.

No evidence packet was created and no prepared consumer packet references were
refreshed.

| Area | Follow-up status |
| --- | --- |
| Final root exists | No. |
| Root used | None. |
| Domain owner / approval evidence | None available. |
| DNS proof | Not collected - final root unavailable. |
| TLS certificate metadata | Not collected - final root unavailable. |
| HTTP to HTTPS redirect proof | Not collected - final root unavailable. |
| Public feed fetch proof | Not collected - final root unavailable. |
| All five final feed URLs | Not collected - final root unavailable. |
| Validator records | Not collected - final root unavailable. |
| Redacted proxy/config summary | Not collected - final root unavailable. |
| Evidence packet README/checksums | Not created - no real final-root artifacts exist. |
| Prepared packet refresh | Not performed - final-root evidence unavailable. |

Exact next actions:

1. Identify a candidate agency-owned or agency-approved root.
2. Obtain retained operator approval for using that root as the public feed
   root.
3. Configure DNS for the approved root.
4. Configure TLS and HTTP-to-HTTPS redirect behavior if HTTP is exposed.
5. Deploy all five public feed URLs at the approved root.
6. Collect DNS, TLS, redirect, public fetch, validator, proxy/config, README,
   and checksum evidence under a dated final-root evidence packet.
7. Run the hosted evidence packet audit only after the packet contains real
   final-root artifacts.

## Future Agency-Owned Domain Next Actions

1. Identify the candidate agency-owned or agency-approved root.
2. Confirm agency/operator approval for that root and for use in submissions.
3. Configure DNS.
4. Configure TLS.
5. Deploy all five public feed URLs at the final root.
6. Run validators against the final root.
7. Collect the final-root evidence packet.
8. Refresh prepared packets with final-root evidence.
9. Update the Track A submission workflow/status only with retained evidence.

## Claim Boundary

Allowed now:

- "The DuckDNS OCI pilot provides pilot evidence for the recorded URL root."
- "Agency-owned domain readiness requirements are documented."
- "Phase 23 documented that final-root proof is blocked because no
  agency-owned or agency-approved root is available."

Not supported yet:

- "The OCI pilot proves agency-owned production-domain compliance."
- "The final agency-owned feed root has validator-clean production evidence."
- "Consumers accepted or ingested the final agency-owned feed root."
- "Open Transit RT is CAL-ITP/Caltrans compliant."
