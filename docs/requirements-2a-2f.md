# Requirements 2A–2F

These requirements formalize the missing non-v1 needs identified in the architecture review. They are written to be portable into Codex as binding requirements.

---

## RQ-2A — Service-day and time-disambiguation correctness

### Goal
The system must correctly resolve the active service day and the exact scheduled trip instance for every vehicle, including trips that cross midnight and trips defined with `frequencies.txt`.

### Why this exists
Trip matching is not only route matching. Correct GTFS Realtime requires identifying the correct trip instance, which may require `trip_id`, `start_time`, and `start_date` for frequency-based or repeated trips.

### Required behaviors
- Compute agency-local service day using the agency timezone, not server timezone.
- Support trips with times beyond `24:00:00` in static GTFS.
- Distinguish multiple instances of the same `trip_id` by `start_date` and `start_time` when needed.
- Support `frequencies.txt`, including both `exact_times=0` and `exact_times=1`.
- Preserve continuity when a vehicle remains on the same trip instance across feed refreshes.
- Carry trip-instance identity into both Vehicle Positions and Trip Updates.

### Data requirements
- Store agency timezone and service-date logic explicitly.
- Store trip instance identity as:
  - `trip_id`
  - `start_date`
  - `start_time`
- Maintain a service-day resolution layer separate from raw UTC timestamps.

### Acceptance criteria
- For a sample GTFS with after-midnight service, the matcher assigns the correct service date and trip instance for telemetry at `00:30`, `01:15`, and `02:10` local time.
- For a frequency-based feed, the system emits trip descriptors that include `trip_id`, `start_time`, and `start_date` when required.
- Unit and integration tests cover overnight, repeated-trip, and frequency-based scenarios.

---

## RQ-2B — Manual override and operations control

### Goal
The system must support human correction of bad or ambiguous assignments, because real operations include deadheading, vehicle swaps, short turns, missed trips, and ad hoc dispatch decisions that telemetry alone cannot always infer.

### Required behaviors
- Operators can manually assign a vehicle to:
  - a route
  - a trip instance
  - an out-of-service state
  - a deadhead state
  - a layover state
- Operators can clear or replace an incorrect assignment.
- Manual overrides can be time-bounded and expire automatically.
- Every override is audit-logged with:
  - operator identity
  - timestamp
  - old state
  - new state
  - reason
- Matching engine must respect an active override and stop auto-reassigning until the override expires or is cleared.
- Operators can mark:
  - canceled trip
  - added trip
  - vehicle swap
  - detour state
  - short turn

### Data requirements
- `manual_override` table with:
  - `agency_id`
  - `vehicle_id`
  - `override_type`
  - `trip_id`
  - `start_date`
  - `start_time`
  - `state`
  - `expires_at`
  - `reason`
  - `created_by`
  - `created_at`
- `audit_log` table for all privileged operational actions.

### Acceptance criteria
- An operator can manually force a vehicle onto a trip instance and that assignment appears in public feeds within one refresh cycle.
- Audit history is queryable by agency and vehicle.
- Clearing the override returns the vehicle to automatic matching.

---

## RQ-2C — Feed quality observability and SLOs

### Goal
The system must expose measurable quality indicators for feed freshness, validity, matching quality, and coverage.

### Required metrics
- percent of active vehicles with recent telemetry
- percent of active vehicles with trip assignment
- percent of active scheduled trips represented in Trip Updates
- stale-vehicle count by threshold
- unmatched-vehicle count
- low-confidence assignment count
- protobuf validation pass/fail rate
- HTTP endpoint availability
- p50 and p95 feed generation latency
- percent invalid feed responses
- percent entities older than SLA threshold

### Required behaviors
- Dashboard per agency and per feed type
- Alerting on threshold breaches
- Historical storage for trend analysis
- Metrics export in a standard format such as Prometheus
- Daily validation reports for static GTFS and realtime feeds
- Feed-health API for operational automation

### Suggested SLOs
- Feed endpoint availability: >= 99.0%
- Invalid realtime responses: < 1%
- Vehicle Positions freshness target: <= 90 seconds
- Trip Updates freshness target: <= 90 seconds
- Public feed generation latency: p95 <= 5 seconds

### Acceptance criteria
- System can show whether Vehicle Positions and Trip Updates entities are younger than configured freshness thresholds.
- System can show whether feed endpoints exceed uptime and validity targets over rolling windows.
- Historical metrics can be queried by route, agency, and day.

---

## RQ-2D — Public publishing, permanence, and data-governance workflow

### Goal
The application must support production publication of GTFS and GTFS-Realtime feeds through stable public URLs, open licensing, and discoverable documentation.

### Required behaviors
- Static GTFS and each GTFS-Realtime feed have stable public URLs.
- URLs remain constant across dataset updates.
- Feeds are served over HTTPS.
- Feed endpoints do not require login.
- Provider website can expose:
  - feed URLs
  - technical contact
  - open license
  - last updated time
  - version or revision metadata
- Optional API keys, if used, must remain automated and non-discriminatory.
- Feed publishing supports staged rollout and rollback without URL changes.
- Publishing workflow separates draft, staged, and active feed versions.

### Data requirements
- `published_feed` metadata model with:
  - canonical public URL
  - feed type
  - license
  - contact email
  - website metadata
  - revision timestamp
  - activation status

### Acceptance criteria
- A feed update does not change the public URL.
- The feed can be fetched anonymously over HTTPS.
- The website can display the license and contact info associated with the feed.
- Deployment supports rollback without URL change.

---

## RQ-2E — Failure modes and degraded-state handling

### Goal
The system must behave predictably when data is missing, delayed, noisy, or contradictory.

### Required failure modes
- missing shapes
- missing stop_times
- duplicate telemetry
- out-of-order telemetry
- long telemetry gaps
- clock skew on devices
- inconsistent vehicle IDs
- route detours
- short turns
- deadheading
- vehicle swaps
- block transition ambiguity
- GTFS import with partial validity

### Required behaviors
- Prefer `unknown` over false certainty when evidence is weak.
- Preserve last known assignment only within bounded confidence and time windows.
- Mark stale vehicles as stale or suppress them after configurable TTL.
- Allow detour tolerance when a detour flag or alert is active.
- All failure modes generate internal incidents, warnings, or degraded-state markers.
- Import publish is atomic: invalid imports do not partially activate.
- Public feed must remain valid protobuf even if some vehicles are unmatched or stale.

### Acceptance criteria
- Integration tests cover each failure mode above.
- A malformed or ambiguous case does not emit a confidently wrong trip descriptor.
- Feed remains valid protobuf even when many vehicles are unmatched.

---

## RQ-2F — Multi-tenancy, auth, and auditability

### Goal
The application must be safely operable for multiple agencies and devices, even if v1 keeps auth simple.

### Required behaviors
- Agency-scoped data model
- Device credentials per device
- Key rotation support
- Operator roles:
  - admin
  - editor
  - operator
  - read-only
- Audit logging for:
  - GTFS publishes
  - manual overrides
  - alert edits
  - device reassignment
  - feed configuration changes
- Per-agency endpoint configuration and branding metadata
- Hard separation between draft and published feed versions
- Per-agency access control on every read/write path

### Data requirements
- `agency_user`
- `device_credential`
- `role_binding`
- `audit_log`
- `feed_config`

### Acceptance criteria
- One agency cannot read or mutate another agency’s draft or realtime state.
- Device tokens can be revoked without affecting other devices.
- All privileged actions are audit logged.
