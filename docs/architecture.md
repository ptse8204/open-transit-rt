# Architecture

## Product shape

Open Transit RT should be a modular platform, not a monolith rewrite of every open-source transit tool.

### Main modules
1. GTFS Studio
2. GTFS importer/publisher
3. Telemetry ingest
4. State engine
5. Vehicle Positions publisher
6. Prediction adapter
7. Trip Updates publisher
8. Alerts publisher
9. Validator + monitoring
10. Admin web

## Core boundary

The critical modular contract is:
- your system emits GTFS-RT Vehicle Positions
- a prediction engine consumes that feed
- the prediction engine emits GTFS-RT Trip Updates

That makes Trip Updates pluggable.

## MVP phases

### MVP A
- import GTFS
- ingest telemetry
- rule-based trip matcher
- Vehicle Positions feed
- basic monitoring

### MVP B
- Trip Updates via TheTransitClock adapter
- unmatched vehicle review tools
- better feed QA

### MVP C
- GTFS Studio interactive editor
- multi-agency hosting
- alerts workflow
- stronger onboarding and monitoring
