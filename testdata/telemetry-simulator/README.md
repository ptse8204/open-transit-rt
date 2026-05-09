# Telemetry Simulator Fixtures

These scenarios are synthetic-only inputs for `cmd/telemetry-simulator`.
They are not vendor payloads, private telemetry, or evidence packets.

The default local scenario is `on-route`, which targets the local demo seed:

- `agency_id=demo-agency`
- `device_id=device-1`
- `vehicle_id=bus-1`
- device token `dev-device-token`
- active GTFS imported from `testdata/gtfs/valid-small`

Some scenarios require alternate synthetic reference feeds or device bindings.
Those requirements are encoded in each fixture's `requires` list so operators
can distinguish default local checks from reference-deployment trials.
