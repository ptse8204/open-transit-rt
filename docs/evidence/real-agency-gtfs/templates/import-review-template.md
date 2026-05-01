# Real Agency GTFS Import Review Template

Status: template only. Do not fill this with fake evidence.

## Agency And Source

| Field | Value |
| --- | --- |
| Agency or feed owner |  |
| GTFS source URL or private source note |  |
| Source type | Public GTFS / agency-provided / other |
| Permission or license note |  |
| Operator or reviewer |  |
| Review date |  |

## Metadata Approval

| Field | Approved value | Approved by | Approval date | Notes |
| --- | --- | --- | --- | --- |
| Agency name |  |  |  |  |
| Agency URL |  |  |  |  |
| Timezone |  |  |  |  |
| Technical contact email |  |  |  |  |
| License name |  |  |  |  |
| License URL |  |  |  |  |
| Public feed root |  |  |  |  |

## Redaction And Privacy Review

| Check | Result | Notes |
| --- | --- | --- |
| No private contracts |  |  |
| No private contact info |  |  |
| No private operator notes |  |  |
| No private ticket links or portal screenshots |  |  |
| No non-public vehicle or device identifiers |  |  |
| No raw private telemetry |  |  |
| No credentials, tokens, private keys, or DB URLs with passwords |  |  |
| Public GTFS license or permission reviewed |  |  |

## Validation Output Handling

Raw validation outputs must be reviewed before commit.

If raw outputs include private paths, private contacts, private operator notes, non-public data, credentials, or private agency material, keep the raw output private and commit only a redacted summary.

Do not commit private agency data or fake validation artifacts.

| Output | Status | Evidence path or private-retention note |
| --- | --- | --- |
| Open Transit RT import validation |  |  |
| Static GTFS canonical validation |  |  |
| Post-publish `schedule.zip` check |  |  |

## Import Result

| Field | Value |
| --- | --- |
| Import command or workflow used |  |
| Import timestamp |  |
| Import result | Passed / failed / blocked |
| Active feed version after import |  |
| Validation report reference |  |
| Follow-up fixes required |  |

## Publish Approval

| Field | Value |
| --- | --- |
| Publish approved by |  |
| Approval date |  |
| Publication environment | local / pilot / production-directed |
| Notes |  |

## Public Feed Verification

| Feed | URL reviewed | Result | Notes |
| --- | --- | --- | --- |
| `feeds.json` |  |  |  |
| `schedule.zip` |  |  |  |
| Vehicle Positions |  |  |  |
| Trip Updates |  |  |  |
| Alerts |  |  |  |

## Phase 23 Final-Root Status

| Question | Answer |
| --- | --- |
| Is the public feed root agency-owned or agency-approved? |  |
| Is retained approval evidence available? |  |
| Were DNS, TLS, public fetch, and validator records collected for the final root? |  |
| If not, what is the current limit of the evidence? |  |

If no agency-owned or agency-approved root exists, this packet may support local/demo or pilot review only. It does not prove agency-domain production readiness.

## Final Notes

Summarize what was proven, what was not proven, and the next operator action. Do not claim consumer submission, consumer acceptance, CAL-ITP/Caltrans compliance, agency endorsement, hosted SaaS availability, or production-grade ETA quality unless separate retained evidence supports that exact claim.
