# Requirements — Trip Updates and Production-Quality Realtime

These requirements formalize the gaps between basic trip matching and a production-quality Trip Updates system.

---

## RQ-3A — Trustworthy stop-level ETA prediction

### Goal
Generate stop-level arrival and departure predictions that are good enough to be trusted by riders and downstream consumers.

### Required behaviors
- For every actively assigned in-service trip, produce future stop predictions.
- Predictions must be ordered by increasing `stop_sequence`.
- Predictions must refresh whenever telemetry changes materially.
- Feed refresh must be frequent enough to keep realtime data fresh.
- Distinguish observed delay from downstream predictions.
- Support either `time` or `delay` appropriately based on trip type and available schedule data.
- For frequency-based trips, prefer absolute prediction times rather than delay-only modeling.

### Model requirements
- Maintain:
  - current position along shape
  - last observed stop progress
  - schedule deviation
  - dwell-time assumptions
  - segment travel-time history
- Support backtesting and offline evaluation against historical telemetry.

### Quality metrics
- MAE by stop
- MAE by route
- MAE by time-of-day
- prediction coverage rate
- percent of in-service trips with at least one future `stop_time_update`

### Acceptance criteria
- For in-progress trips, Trip Updates include at least one future stop prediction.
- Predictions remain ordered and validator-clean.
- Historical backtesting can report quality by route, stop, and time-of-day.

---

## RQ-3B — Production handling of detours, short turns, cancellations, and added trips

### Goal
Support common operational disruptions in a way that produces correct public realtime data.

### Required behaviors
- Detect or accept operator input for:
  - detour
  - short turn
  - canceled trip
  - added trip
  - replaced trip
- When a trip is canceled, emit:
  - a `TripUpdate` marked `CANCELED`
  - a corresponding `Alert`
- Support route deviation tolerance for detours.
- Support modified predictions after a short turn or detour.
- Allow Alerts to be linked to affected lines, trips, stops, and time ranges.

### Acceptance criteria
- A canceled trip appears as canceled in Trip Updates and has a corresponding Alert.
- A detoured vehicle can remain matched without being rejected purely for shape deviation when an applicable detour condition exists.
- An added trip can be surfaced without corrupting scheduled-trip coverage metrics.

---

## RQ-3C — Ambiguity resolution for weak-signal or hintless service patterns

### Goal
The matcher must handle ambiguous cases without route hints, trip hints, or perfect telemetry.

### Required behaviors
- Candidate scoring must consider:
  - service day validity
  - nearest shape distance
  - recent movement direction
  - stop-sequence progress
  - schedule fit
  - previous assignment continuity
  - block continuity
- The system must expose confidence scores and reasons.
- Confidence thresholds must be configurable.
- Below threshold, the system should output `unknown` rather than a speculative assignment.

### Acceptance criteria
- For weak-signal cases, system either produces a high-confidence correct assignment or explicitly marks the vehicle unknown.
- The UI can show why a vehicle was assigned or not assigned.

---

## RQ-3D — Full system-wide Trip Updates quality

### Goal
Trip Updates must be operationally complete, coverage-aware, and measurable at the system level.

### Required behaviors
- Coverage reporting for:
  - percent of active trips with Trip Updates
  - percent of active trips with future stop updates
  - percent of canceled trips represented properly
  - percent of vehicles matched to trips
- Detect when trips are active in static GTFS but missing realtime representation.
- Support trip-level quality scoring.
- Distinguish:
  - no realtime data
  - matched but no prediction
  - prediction available
  - canceled
  - added
  - replaced

### Acceptance criteria
- Dashboard shows coverage against active scheduled service.
- Missing realtime coverage is measurable, not hidden.
- Feed remains valid even when coverage is incomplete.

---

## RQ-3E — Human override workflows for bad matches and bad predictions

### Goal
Operations staff must be able to repair bad matching and bad ETA output quickly.

### Required behaviors
- Override trip assignment
- Override service state
- Mark trip canceled
- Mark vehicle swapped
- Mark trip added
- Recompute downstream predictions after override
- Provide incident queue for:
  - unmatched vehicles
  - low-confidence assignments
  - missing predictions
  - detour suspects
  - stale telemetry

### Acceptance criteria
- Operator can resolve a bad match and see corrected public output quickly.
- Incident queue can be filtered by agency, severity, route, and age.

---

## RQ-3F — Consumer-grade confidence and consistency

### Goal
Realtime output must be reliable enough for major consumers and end users.

### Required behaviors
- Stable entity IDs across feed iterations
- Fresh timestamps
- valid protobuf by default
- consistent `FeedHeader.timestamp`
- low invalid-response rate
- no contradictory trip descriptors across successive refreshes
- support `Last-Modified` and normal HTTP caching behavior
- deterministic ordering where practical for easier downstream debugging

### Quality targets
- invalid responses < 1%
- entity age measurable and bounded
- stable IDs through a trip lifecycle
- feed latency within configured SLA

### Acceptance criteria
- Fewer than 1% invalid responses over a rolling window.
- Entity age SLA for Trip Updates and Vehicle Positions is measurable and enforced.
- Feed can be consumed repeatedly without entity ID churn during a trip.

---

## RQ-3G — Architecture requirement for Trip Updates pluggability

### Goal
Trip Updates should remain replaceable without rewriting telemetry ingest or Vehicle Positions publishing.

### Required behaviors
- Define a narrow prediction adapter boundary:
  - input:
    - active GTFS feed version
    - current telemetry
    - current assignments
    - Vehicle Positions feed URL or data stream
  - output:
    - Trip Updates feed
    - optional prediction diagnostics
- Support multiple prediction backends:
  - deterministic in-house engine
  - external engine such as TheTransitClock
  - future ML-based predictor
- Keep matching ownership outside the predictor unless explicitly configured.

### Acceptance criteria
- A second predictor implementation can be swapped in without changing telemetry ingest contracts.
- Public Trip Updates endpoint remains stable even if the backend predictor changes.
