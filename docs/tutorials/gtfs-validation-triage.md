# GTFS Validation Triage

This guide explains common GTFS import and validation issues in plain language. It is for agency operators and implementers reviewing real schedule data.

Importer errors are produced by Open Transit RT before activation. Canonical validator errors come from external validator tooling. Both matter, but neither proves consumer acceptance.

## Common Issues

| Issue | What it means | Why it matters | How to fix it |
| --- | --- | --- | --- |
| Missing required files | The ZIP does not include a required GTFS file such as `agency.txt`, `routes.txt`, `stops.txt`, `trips.txt`, or `stop_times.txt`. | A schedule cannot be published without the core GTFS tables. | Export a complete GTFS ZIP from the source scheduling tool. |
| Missing usable service source | Neither `calendar.txt` nor `calendar_dates.txt` provides usable service. | Trips cannot be tied to actual service dates. | Add valid weekday service in `calendar.txt` or service additions in `calendar_dates.txt`. |
| Invalid `route_type` | A route has a value outside the supported GTFS route type range. | Trip planners and validators may reject or misclassify the route. | Correct the route type to the agency's actual service mode. |
| Bad `stop_times` references | `stop_times.txt` references a missing trip or stop. | The system cannot tell where a trip goes or when it stops. | Fix IDs so every stop time references an existing `trip_id` and `stop_id`. |
| Missing trips, stops, or routes references | GTFS rows point to IDs that are not present in their parent table. | Broken references can block import or produce unusable public schedules. | Re-export from the scheduling system or repair the missing IDs consistently. |
| Calendar problems | Service IDs are missing, date ranges are reversed, weekdays are all off, or exceptions remove all service. | The feed may appear to have no operating service. | Review service IDs, start/end dates, weekday flags, and exception dates. |
| Unusable service calendars | The feed validates structurally but has no service in the intended review window. | A public feed with no active service is not useful for launch. | Add the intended service period or confirm the feed is intentionally future-dated. |
| Times beyond `24:00:00` | Late-night service uses GTFS after-midnight times such as `25:10:00`. | These are valid GTFS, but they must match the agency service day and timezone. | Keep valid after-midnight times; fix only impossible or accidental values. |
| Shape ordering | Shape points are missing sequence values or are ordered incorrectly. | Matching and map display can become unreliable. | Export shapes with stable `shape_pt_sequence` order for every `shape_id`. |
| `shape_dist_traveled` issues | Distances are missing, decreasing, or inconsistent across shapes and stop times. | Consumers and matching logic may have weaker progress information. | Correct distances or omit them consistently if the agency cannot provide reliable values. |
| Missing or unexpected `block_id` | Trips that should stay with the same vehicle do not share a block, or unrelated trips do. | Block continuity helps trip matching and vehicle assignment review. | Add or correct `block_id` values according to the agency's vehicle work. |
| Malformed `frequencies.txt` | Frequency rows have invalid times, bad headways, missing trip IDs, or invalid `exact_times`. | Headway service needs clear trip-instance handling. | Fix referenced trips, time windows, headways, and `exact_times` values. |
| Duplicate IDs | A table repeats an ID that should be unique. | The importer or consumers may not know which row is authoritative. | Make IDs unique within each GTFS table. |
| Timezone/date mistakes | `agency_timezone` or service dates do not match the agency's real operating calendar. | After-midnight trips and service-day matching can be wrong. | Confirm timezone and date ranges with the agency schedule owner. |
| Empty or all-suppressed service | Import succeeds or partially validates, but no trips are active for review. | Public feeds may look available while showing no useful service. | Check calendars, date ranges, exceptions, and whether the review date is inside the service period. |
| Validator errors versus importer errors | The importer blocks activation for Open Transit RT's minimum contract; canonical validators check broader GTFS quality. | Passing one check does not guarantee passing the other. | Resolve importer errors first, then run and review canonical validator output. |

## How To Triage

1. Start with importer errors. If import activation failed, no public schedule update should be treated as complete.
2. Fix required files and broken references before tuning warnings.
3. Review calendars and timezone early; many downstream issues come from bad service dates.
4. Check shapes, frequencies, and blocks after the core schedule imports.
5. Run canonical validators before stronger readiness claims.
6. Keep raw validation output private until redaction review is complete.

## When To Ask For Technical Help

Ask for technical help when:

- the same import error remains after the agency re-exports GTFS
- the feed uses complicated after-midnight, frequency-based, interlined, or seasonal service
- validation output mentions many broken references across several files
- the agency is unsure whether a field is public-safe to commit
- raw validation output includes private paths, contacts, operator notes, or non-public data
- `feeds.json` or `schedule.zip` does not match the active feed after a successful import
- canonical validator output disagrees with the Open Transit RT importer in a way the operator cannot explain
- the operator wants to make a public readiness, compliance, or consumer-acceptance claim

When asking for help, share a redacted summary first. Do not paste private GTFS, private logs, credentials, private contacts, or raw validation output that has not been reviewed.
