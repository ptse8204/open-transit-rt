# Codex Task Brief

Use this repo as the starting point.

## Goal
Build a production-directed MVP for **Open Transit RT**, a low-cost transit realtime stack for small agencies.

## What to prioritize
1. Keep the codebase mostly in Go.
2. Treat Trip Updates as a pluggable module.
3. Make Vehicle Positions the first high-quality output.
4. Keep GTFS Studio and existing-GTFS import both in scope.
5. Avoid premature rider-app or CAD/dispatch scope.

## Immediate implementation tasks

### Task 1: Persist telemetry
- replace in-memory event storage with Postgres
- add repository interfaces
- wire schema from `db/schema.sql`
- add migration strategy

### Task 2: Build deterministic trip matcher
Implement rule-based matching using:
- active service day
- candidate trips by route and calendar
- nearest shape matching
- stop sequence progress
- continuity from previous assignment
- block transition awareness

### Task 3: Publish true GTFS-RT Vehicle Positions
- generate protobuf output
- expose `/public/gtfsrt/vehicle_positions.pb`
- keep JSON debug endpoint for inspection

### Task 4: Add GTFS import pipeline
- ingest `gtfs.zip`
- validate required files
- normalize into database tables
- activate a published feed version

### Task 5: Add GTFS Studio draft model
- draft tables
- publish flow into canonical feed version
- first UI can be minimal server-rendered HTML or simple SPA
