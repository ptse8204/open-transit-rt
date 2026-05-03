# Synthetic Multi-Agency Fixtures

These fixtures are public-safe synthetic records for Phase 27 isolation tests and documentation.

They do not represent real agency data, real private GTFS, real telemetry, real device tokens, consumer acceptance, or production multi-tenant hosting proof.

The canonical synthetic agencies are:

- `agency-a`
- `agency-b`

Use these identifiers for new multi-agency tests unless a narrower test needs a package-local fixture:

- feed versions: `feed-a`, `feed-b`
- vehicles: `bus-a-1`, `bus-b-1`
- devices: `device-a-1`, `device-b-1`
- routes: `route-a-10`, `route-b-20`
- trips: `trip-a-10`, `trip-b-20`
- public roots: `https://agency-a.example`, `https://agency-b.example`

Device tokens in tests must be generated or synthetic placeholders only. Do not commit real device tokens, admin tokens, JWT/CSRF secrets, private telemetry, private GTFS, or private operator artifacts.
