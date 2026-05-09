# BYOD Hardware Intake Template

Use this template to inventory operator-owned or vendor-owned AVL hardware
before connecting it to Open Transit RT.

This template is not a hardware certification, vendor compatibility statement,
paid support offer, hosted service offer, service-level commitment, compliance
claim, or production-readiness claim.

## Operator

- Agency/operator name: `<operator>`
- Review date: `<YYYY-MM-DD>`
- Reviewer: `<name/role>`
- Deployment context: `<local/reference/pilot/other>`

## Device Inventory

| Vehicle ID | Device ID | Hardware model | Firmware | Time source | GPS source | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| `<vehicle>` | `<device>` | `<model>` | `<version>` | `<clock>` | `<gps>` | `<notes>` |

## Telemetry Contract Review

- Target endpoint: `/v1/telemetry`
- Authentication: device bearer token managed by the operator
- Required fields confirmed: `<yes/no>`
- Timestamp format confirmed as UTC: `<yes/no>`
- Latitude/longitude precision reviewed: `<yes/no>`
- Duplicate and out-of-order behavior reviewed: `<yes/no>`
- Stale telemetry behavior reviewed: `<yes/no>`

## Data-Quality Prerequisites

- Device clock synchronization approach: `<details>`
- Vehicle-to-device binding source: `<details>`
- GPS outage handling: `<details>`
- Low-quality position filtering before send: `<details>`
- Operator review process for unknown/unmatched vehicles: `<details>`

## Redaction Notes

Do not attach raw device tokens, private vendor payloads, raw logs, driver
names, private support tickets, or database URLs. Use redacted summaries and
checksums when a later approved evidence workflow requires artifacts.
