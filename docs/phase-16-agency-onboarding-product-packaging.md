# Phase 16 — Agency Onboarding And Product Packaging

## Status

Planned phase. Not implemented until `docs/handoffs/latest.md` marks it active.

## Purpose

Phase 16 makes Open Transit RT easier for a small, non-expert agency to try and operate. The current repo proves substantial backend capability, but it still expects the operator to understand Go services, environment variables, validators, tokens, and multiple ports.

This phase turns the project from a developer/backend system into a more approachable pilot package.

## Scope

1. Full local app packaging.
2. One-command agency demo.
3. Simpler setup documentation.
4. First-run agency checklist.
5. Device onboarding simplification.
6. Better operator feedback after commands.

## Required Work

### 1) Full App Docker Compose

Add or improve a deployment profile that can run:

- Postgres/PostGIS;
- agency-config;
- telemetry-ingest;
- feed-vehicle-positions;
- feed-trip-updates;
- feed-alerts;
- gtfs-studio;
- optional reverse proxy;
- optional validator helpers.

This may use local builds if production images are not yet published.

### 2) One-Command Demo

Create or refine one command that runs the complete demo and prints:

- public feed URLs;
- admin URL;
- generated admin token instructions;
- device telemetry token instructions;
- validation status;
- where to find logs;
- what to do next.

### 3) Agency First-Run Guide

Add a human-friendly `docs/tutorials/agency-first-run.md` that avoids assuming deep developer knowledge.

It should explain:

- what GTFS is;
- what GTFS Realtime is;
- what a device token does;
- what public feed URLs are;
- why validation matters;
- what consumers still need to accept separately.

### 4) Device Onboarding

Simplify docs and scripts for:

- creating a device token;
- binding a device to a vehicle;
- rotating/rebinding a device;
- sending a sample telemetry event;
- using a simulator when no hardware exists.

### 5) Operator-Friendly Output

Scripts should print clear next actions and avoid dumping internal-only implementation detail.

## Acceptance Criteria

Phase 16 is complete only when:

- a non-expert can run one command or a short sequence and see working feeds;
- the full app can run without manually starting six separate Go processes;
- agency onboarding docs exist and are understandable;
- device onboarding is clear;
- no unsupported compliance or acceptance claims are introduced;
- existing demo and tests still pass.

## Required Checks

```bash
make validate
make test
make smoke
make demo-agency-flow
docker compose -f deploy/docker-compose.yml config
git diff --check
```

## Explicit Non-Goals

Phase 16 does not:

- build a full hosted SaaS;
- add fares, payments, rider apps, or CAD/dispatch;
- claim consumer acceptance;
- replace Phase 12 deployment evidence or Phase 13 consumer evidence;
- require Kubernetes.
