# Device And AVL Integration

This guide is for an agency, device installer, AVL vendor, or simple GPS emitter that needs to send vehicle telemetry into Open Transit RT.

It documents the current repository API. It does not prove real-world device reliability, certified hardware or vendor support, consumer acceptance, CAL-ITP/Caltrans compliance, hosted SaaS availability, or production AVL quality.

All identifiers and tokens below are synthetic demo values. Do not commit real private device IDs, vehicle IDs, vendor IDs, tokens, raw private telemetry, private logs, or vendor payloads in public docs.

## Telemetry Endpoint

Send one vehicle observation at a time:

```text
POST /v1/telemetry
Authorization: Bearer <device-token>
Content-Type: application/json
```

For the local app package, the default base URL is:

```text
http://localhost:8080
```

The full local demo endpoint is:

```text
http://localhost:8080/v1/telemetry
```

The local app seed includes a synthetic demo binding for:

```text
agency_id=demo-agency
device_id=device-1
vehicle_id=bus-1
device token=dev-device-token
```

The demo token is for local evaluation only. Real deployments must generate and store their own device tokens privately.

## Payload Fields

Required fields:

| Field | Meaning | Current validation |
| --- | --- | --- |
| `agency_id` | Agency receiving the telemetry. | Must be non-empty and exist in the database. |
| `device_id` | Device credential identifier. | Must be non-empty and match the token binding. |
| `vehicle_id` | Vehicle identifier for this observation. | Must be non-empty and match the token binding. |
| `timestamp` | Observation time from the device. | Must be RFC3339 with timezone or offset. |
| `lat` | Latitude. | Must be between `-90` and `90`. |
| `lon` | Longitude. | Must be between `-180` and `180`. |

Optional fields:

| Field | Meaning | Guidance |
| --- | --- | --- |
| `driver_id` | Optional operator/driver identifier. | Omit unless the value is public-safe and approved for this deployment. |
| `bearing` | Direction of travel in degrees. | Send measured bearing when available. Do not fabricate it. Numeric `0` is valid true north when actually observed. |
| `speed_mps` | Speed in meters per second. | Send measured speed when available. Do not convert from stale or guessed data. |
| `accuracy_m` | GPS accuracy estimate in meters. | Send the device-reported estimate when available. Low accuracy can reduce matching confidence. |
| `trip_hint` | Static GTFS `trip_id` hint. | Use only when the device or vendor knows the scheduled trip. Do not guess. |

## Time, Clock, And GPS Expectations

The `timestamp` should be the observation time, not the time a batch or adapter forwards the payload. Device clocks should be synced with a reliable time source.

Timestamps that are too old can produce stale telemetry, unknown assignments, missing trip descriptors, stale Vehicle Positions, or withheld Trip Updates. Timestamps in the future can make freshness and agency-local service-day review confusing even if the payload is accepted.

The agency timezone matters for matching service days, after-midnight trips, repeated trip instances, and frequency-based service. Open Transit RT resolves service day from the agency timezone configured with the active GTFS feed.

Latitude and longitude must be WGS84 decimal degrees. GPS accuracy affects matching confidence. Bearing and speed help when measured, but guessed bearing or speed can make matching harder to troubleshoot.

Telemetry gaps can suppress realtime output after configured stale thresholds. A valid accepted payload is only one input into Vehicle Positions and Trip Updates; matching still depends on current GTFS, assignments, service state, and freshness.

## Example Request

This request uses only synthetic local demo values:

```bash
curl -i -X POST "http://localhost:8080/v1/telemetry" \
  -H "Authorization: Bearer dev-device-token" \
  -H "Content-Type: application/json" \
  --data '{
    "agency_id": "demo-agency",
    "device_id": "device-1",
    "vehicle_id": "bus-1",
    "timestamp": "2026-05-02T16:30:00Z",
    "lat": 49.2827,
    "lon": -123.1207,
    "bearing": 120.0,
    "speed_mps": 8.4,
    "accuracy_m": 7.5,
    "trip_hint": "trip-10-0800"
  }'
```

The helper script can print the same kind of payload without sending it:

```bash
scripts/device-onboarding.sh sample --dry-run
scripts/device-onboarding.sh simulate --dry-run
```

Dry-run output is a local demonstration, not production AVL proof.

## Confirmed Response Behavior

The telemetry ingest handler currently returns JSON for stored telemetry outcomes with these fields:

| Field | Meaning |
| --- | --- |
| `accepted` | `true` only when the event is the latest accepted observation for the vehicle stream. |
| `ingest_status` | Confirmed values include `accepted`, `duplicate`, and `out_of_order`. |
| `agency_id` | Agency from the stored event. |
| `vehicle_id` | Vehicle from the stored event. |
| `observed_at` | Parsed observation timestamp. |
| `received_at` | Database receive timestamp. |

Confirmed success status codes:

| HTTP status | Confirmed meaning |
| --- | --- |
| `201 Created` | The event was accepted as the latest observation. |
| `202 Accepted` | The event was stored but was not accepted as the latest observation, such as duplicate or out-of-order telemetry. |

Example `201 Created` body shape:

```json
{
  "accepted": true,
  "ingest_status": "accepted",
  "agency_id": "demo-agency",
  "vehicle_id": "bus-1",
  "observed_at": "2026-05-02T16:30:00Z",
  "received_at": "2026-05-02T16:30:01Z"
}
```

For duplicate or out-of-order cases, treat `202 Accepted` plus `accepted=false` and `ingest_status` as troubleshooting signals. Do not build external integrations that depend on unversioned debug wording beyond the confirmed fields above.

Confirmed error status codes:

| HTTP status | Confirmed condition |
| --- | --- |
| `400 Bad Request` | Invalid JSON or invalid telemetry payload, including timezone-less timestamp or invalid latitude/longitude. |
| `401 Unauthorized` | Missing device token, invalid token, or token binding mismatch. |
| `404 Not Found` | Unknown agency. |
| `405 Method Not Allowed` | Any method other than `POST` on `/v1/telemetry`. |

Error bodies are plain troubleshooting text from the HTTP handler, not a versioned JSON API contract.

## Verify Telemetry Was Accepted

Use the response first. A fresh observation should return `201 Created` and `accepted=true`.

Then use the authenticated Operations Console:

```text
http://localhost:8080/admin/operations/telemetry
```

The telemetry view shows vehicle ID, device ID, observed time, received time, age, freshness, assignment state, route, trip, confidence, reason codes, and assignment time. It omits raw payloads and token fields.

The protected debug event list is also available to admins:

```text
/v1/events?agency_id=demo-agency&limit=25
/admin/debug/telemetry/events?agency_id=demo-agency&limit=25
```

`/v1/events` is an authenticated admin/debug review path. It is not a public feed, not a consumer-facing endpoint, and not a GTFS Realtime endpoint.

For public feed effects, fetch Vehicle Positions:

```text
/public/gtfsrt/vehicle_positions.pb
```

Public protobuf feeds are anonymous. JSON debug feed endpoints are protected and should not be exposed as production public surfaces.

## Vendor AVL Adapter Boundary

Vendor AVL integrations should transform vendor-specific payloads into the Open Transit RT telemetry event shape before forwarding to `/v1/telemetry`.

Use one of these patterns:

- an agency-owned adapter script;
- a sidecar service operated by the deployment owner;
- vendor-owned middleware that calls Open Transit RT;
- a private integration process that validates and redacts before forwarding.

Adapter responsibilities:

- map vendor vehicle/device identifiers to Open Transit RT `device_id` and `vehicle_id`;
- keep vendor credentials outside this repo;
- validate required fields before forwarding;
- preserve observation timestamp and measured GPS fields;
- avoid sending stale batches as if they were current observations;
- record only redacted integration evidence after review.

Do not add vendor-specific coupling to core matching, Vehicle Positions generation, or Trip Updates prediction. Do not claim certified vendor support, marketplace equivalence, or production hardware compatibility without retained public-safe evidence.

## Troubleshooting

| Issue | Symptom | Likely cause | How to check | Next action | What not to claim yet |
| --- | --- | --- | --- | --- | --- |
| Bad token or missing `Authorization` header | `401 Unauthorized`; telemetry does not appear in Operations Console. | Missing `Authorization: Bearer <device-token>`, wrong token, expired old token after rebind, or copied token with whitespace. | Check the request headers locally without printing the token into public logs. Review `/admin/operations/devices` for the device binding and last-used time. | Rotate/rebind the device token and update the device or adapter secret store. | Do not claim the device is integrated or sending live telemetry. |
| Wrong agency/device/vehicle | `401 Unauthorized` for binding mismatch, or `404 Not Found` for unknown agency. | Payload `agency_id`, `device_id`, or `vehicle_id` does not match the token binding or known agency. | Compare the payload identifiers to the binding shown in `/admin/operations/devices`. | Correct the adapter mapping or rebind the device to the intended vehicle. | Do not claim a vendor mapping is complete. |
| Timestamp too old | Request may be stored, but Operations Console shows stale telemetry; Vehicle Positions may omit trip descriptor or suppress the vehicle after stale thresholds. | Device clock is behind, adapter forwards delayed batches, or network buffering is too long. | Check `observed_at`, `received_at`, and age in `/admin/operations/telemetry`. | Sync device clock and forward observations promptly. | Do not claim realtime freshness or production AVL reliability. |
| Timestamp in the future | Telemetry may look fresh even when the observation time is wrong; matching service-day review may be confusing. | Device clock is ahead or adapter used an incorrect timezone conversion. | Compare payload `timestamp` with trusted UTC time and agency local time. | Fix clock sync and timezone handling before collecting evidence. | Do not claim matching quality or service-day correctness from that sample. |
| Invalid lat/lon | `400 Bad Request` with invalid telemetry payload. | Latitude outside `-90..90`, longitude outside `-180..180`, missing coordinate, or non-numeric coordinate. | Validate the JSON payload before sending. | Fix device parser or coordinate mapping. | Do not claim telemetry ingest works for that device payload. |
| Low GPS accuracy | Telemetry is accepted but assignment confidence is low, unknown, or degraded. | Device has poor sky view, reports low accuracy, or vendor payload has imprecise coordinates. | Review `accuracy_m` in the payload source and assignment reason/degraded state in `/admin/operations/telemetry`. | Improve antenna/device placement, forward real accuracy, and avoid fabricated precision. | Do not claim production matching quality. |
| Telemetry accepted but no assignment | `201 Created` but Operations Console shows no assignment, unknown, low confidence, ambiguous, or missing schedule data. | No active GTFS, stale data, off-shape location, weak trip evidence, no matching candidate, or conservative threshold behavior. | Review active feed status, telemetry freshness, assignment state, confidence, reason codes, and degraded state in `/admin/operations/telemetry`. | Fix GTFS, device mapping, trip hints, shape quality, or manual override workflow as appropriate. | Do not claim the vehicle is matched to a trip. |
| Vehicle Positions stale or missing | Public Vehicle Positions protobuf omits the vehicle or includes it without trip descriptor. | Telemetry is stale, suppressed after stale threshold, no assignment, assignment/telemetry mismatch, low confidence, or non-service state. | Review `/admin/operations/telemetry` and protected Vehicle Positions JSON debug if authorized. | Send fresh telemetry and resolve assignment or service-state issues. | Do not claim public realtime coverage for that vehicle. |
| Trip Updates withheld | Vehicle Positions may exist, but Trip Updates output is empty or diagnostic counts show withheld reasons. | Trip Updates adapter conservatively withholds weak, stale, degraded, unsupported, or unmatched cases. | Review `/admin/operations/feeds` Trip Updates quality diagnostics. | Improve telemetry freshness, assignment confidence, GTFS data, and operational overrides. | Do not claim production-grade ETA quality. |
| Simulator works but real hardware still unproven | `scripts/device-onboarding.sh simulate` produces local accepted events, but no hardware has been tested. | Simulator uses synthetic local identifiers and controlled payloads. | Label evidence as simulator or no-hardware only. | Test the real device or vendor adapter with approved redacted evidence collection. | Do not claim hardware reliability, vendor support, or production AVL proof. |
| Validator passes but consumer acceptance is not proven | Static or realtime validation records pass, but consumer records remain `prepared` or not accepted. | Validators check feed format; consumers have separate ingestion workflows. | Review consumer evidence records and `docs/evidence/consumer-submissions/submission-workflow.md`. | Collect target-originated submission or acceptance evidence before status changes. | Do not claim Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, transit.land, or any other consumer accepted the feed. |

## Redaction Rules

Do not commit:

- device tokens;
- admin tokens;
- JWT or CSRF secrets;
- DB URLs with passwords;
- vendor credentials;
- private AVL payloads;
- raw private telemetry;
- private device, vehicle, or vendor identifiers;
- private operator notes;
- private logs with credentials;
- `.cache` files.

If device or AVL evidence is collected later, it must record source, permission, redaction review, whether identifiers are public-safe or synthetic, whether data is simulator/pilot/real-device, and what the evidence proves and does not prove.

