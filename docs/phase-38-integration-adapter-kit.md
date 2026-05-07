# Phase 38 — Integration Adapter Kit

## Goal

Make Open Transit RT easy to integrate with existing AVL, device, and prediction systems.

## Required work

- document adapter architecture in one place;
- provide telemetry adapter examples;
- add mapping fixtures for common payload shapes;
- document `POST /v1/telemetry` contract for adapters;
- document external predictor adapter lifecycle;
- add adapter conformance tests where practical;
- add troubleshooting matrix.

## Boundaries

Adapters are integration points. They do not prove certified vendor compatibility, production AVL reliability, or production-grade ETA quality.
