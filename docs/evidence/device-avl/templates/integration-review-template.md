# Device / AVL Integration Review Template

This template is for future public-safe evidence only. Do not fill it with fake device data, fake vendor approvals, fake hardware certifications, private AVL payloads, credentials, or raw private telemetry.

## Summary

| Field | Value |
| --- | --- |
| Review date |  |
| Reviewer |  |
| Agency or pilot environment |  |
| Integration type | Simulator / pilot device / real device / vendor adapter / agency sidecar |
| Evidence label | Synthetic / simulator / pilot / real-device / production-directed |
| Source |  |
| Source permission scope |  |
| Permission recorded | Yes / No |
| Public-safe to commit | Yes / No |
| Redaction review completed | Yes / No |
| Payload redaction review completed | Yes / No |
| Identifier mapping review completed | Yes / No |
| Synthetic or public-safe identifiers only | Yes / No |

## Identifiers

| Identifier type | Value or redaction | Public-safe or synthetic? | Notes |
| --- | --- | --- | --- |
| Agency ID |  |  |  |
| Device ID |  |  |  |
| Vehicle ID |  |  |  |
| Trip hint, if used |  |  |  |
| Vendor identifier, if relevant |  |  |  |

Do not include private device serial numbers, private vehicle IDs, vendor account IDs, or tokens unless the agency has explicitly approved them as public-safe. Prefer synthetic identifiers in committed evidence.

For Phase 29B synthetic adapter review, record the mapping source and fixture names, not private payloads or credentials. Dry-run adapter output is transform-review output only unless a separate approved telemetry submission record exists.

## Integration Path

Describe how telemetry reached Open Transit RT:

- endpoint used:
- authentication handled by:
- adapter, sidecar, or vendor-owned middleware:
- timestamp source:
- GPS source:
- field validation before forwarding:
- dry-run only or submitted to `/v1/telemetry`:
- token storage reviewed:

## Commands Or Checks Run

List redacted commands or checks. Do not include bearer tokens, admin tokens, DB URLs with passwords, private URLs, or raw private logs.

```text

```

## Results

| Check | Result | Evidence path or note |
| --- | --- | --- |
| Dry-run transform reviewed |  |  |
| Diagnostics reviewed |  |  |
| Telemetry accepted |  |  |
| Operations Console freshness reviewed |  |  |
| Vehicle Positions reviewed |  |  |
| Trip Updates reviewed, if relevant |  |  |
| Validator reviewed, if relevant |  |  |

## Redactions

List every redaction and why it was needed:

-

## What This Proves

-

## What This Does Not Prove

- Dry-run adapter output does not prove telemetry was submitted.
- Partial stdout from a failed dry run does not prove vendor compatibility.
- This does not prove certified hardware or vendor support unless separate retained evidence says so.
- This does not prove production AVL reliability without real operating evidence.
- This does not prove consumer acceptance.
- This does not prove CAL-ITP/Caltrans compliance.

## Next Action

-
