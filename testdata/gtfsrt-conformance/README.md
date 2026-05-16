# GTFS-RT Conformance Fixtures

This directory contains synthetic/local fixture metadata for the offline
GTFS-RT conformance harness.

The fixtures are not evidence, are not consumer submissions, and do not prove
compliance, production readiness, SLA/uptime, vendor compatibility, hardware
certification, consumer acceptance, production-grade ETA quality, or
real-world ETA accuracy.

Use:

```bash
make gtfsrt-conformance
```

Future binary protobuf fixtures should stay synthetic and outside
`docs/evidence/**`.
