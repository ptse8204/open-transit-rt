# Open Transit RT Master Plan — Self-Hosted Agency Reuse

## Purpose

The project is not trying to win external proof first. The purpose is to build a working open-source transit data backend that a small agency, operator, or civic technologist can adapt and reuse.

The near-term proof target is the existing OCI/OCL-style pilot server and repeatable self-hosted deployment, not agency-owned final-root evidence, consumer acceptance, or public launch claims.

## Corrected Strategy

Stop organizing the roadmap around external proof first.

Instead, organize the roadmap around:

```text
1. make the repo understandable;
2. make one server deployment repeatable;
3. make agency onboarding mostly one-command;
4. make CAL-ITP-style readiness visible in product workflows;
5. make integration with existing solutions adapter-based;
6. keep evidence only as internal regression/deployment proof.
```

## Current Baseline

The repo already has:

- GTFS import and publish;
- large real public GTFS handling from Phase 33 Outcome C;
- stable public feed paths for schedule, feeds.json, Vehicle Positions, Trip Updates, and Alerts;
- validator install/check workflows;
- telemetry ingest with bearer-token device auth;
- conservative matching and GTFS-RT feeds;
- Alerts;
- local app flow;
- OCI pilot evidence;
- pilot operations runbook;
- Operations Console;
- device/AVL guide and synthetic adapter;
- external predictor adapter boundary;
- prepared consumer packets.

## Important Misses

### 1. Root README is wrong for agency adoption

The live root `README.md` currently appears to be a roadmap export README, not the public Open Transit RT product front door. This is a major product gap.

Fix it first.

The root README should explain:

- what Open Transit RT is;
- who it is for;
- how to run it locally;
- how to deploy it on the OCI/OCL pilot server pattern;
- how to import agency GTFS;
- how to publish five feed URLs;
- how to connect telemetry/AVL;
- what CAL-ITP-style readiness means;
- what the project does not claim.

### 2. Current handoff still over-prioritizes external evidence forks

`docs/handoffs/latest.md` still frames the next objective around final-root proof, consumer submissions, real agency pilot evidence, and other external evidence forks.

For this project direction, the next objective should become:

```text
Make Open Transit RT easy to self-host, adapt, and integrate using the existing OCI/OCL pilot server as the reference deployment.
```

External evidence stays out of scope unless explicitly requested later.

### 3. Local app and public-GTFS run are not unified enough

`make agency-app-up` uses the demo agency and sample GTFS. The Phase 33 LA Metro run proved real GTFS handling, but it used a more manual local setup.

The agency reuse target should be:

```bash
make agency-pilot-up AGENCY_ID=... GTFS_URL=...
```

or an equivalent script that:

- downloads GTFS to ignored storage;
- creates/seeds the agency;
- imports GTFS with a suitable timeout;
- starts the local/pilot services;
- publishes the five feed paths;
- runs validators or records blockers;
- prints a support summary.

### 4. OCI/OCL pilot deployment should become the reference deployment

The existing OCI pilot has useful hosted/operator evidence and runbooks. Turn it into a reusable reference deployment profile.

This means:

- one documented server layout;
- env templates;
- systemd or compose profile;
- Caddy/reverse-proxy config;
- secrets generation;
- validator install;
- backup/restore;
- feed monitor;
- scorecard export;
- update/rollback commands;
- redaction-safe evidence capture.

### 5. CAL-ITP should be treated as workflow compatibility, not compliance claim

The repo should support CAL-ITP-style readiness by making these visible and repeatable:

- stable public feed URLs;
- static GTFS;
- all three GTFS-RT feed types;
- open license/contact metadata;
- validation records;
- feed freshness;
- consumer packet preparation;
- operations evidence.

Do not claim CAL-ITP/Caltrans compliance unless a future retained evidence path supports it.

### 6. Existing solution integration should be first-class

The project should integrate with existing systems through adapters:

- AVL/device vendors -> telemetry ingest adapter;
- external predictors -> `internal/prediction.Adapter`;
- validators -> pinned validator wrappers;
- monitoring -> Prometheus-format metrics and deployment-owned tools;
- consumers -> standard public GTFS/GTFS-RT URLs and packet templates.

Do not build a closed all-in-one CAD/AVL system.

## New Master Roadmap

### Phase 35 — README And Roadmap Realignment

Goal: fix the repo’s public face and align the roadmap with self-hosted agency reuse.

Required work:

- restore root `README.md` as product front door;
- move export/roadmap wording into docs, not root README;
- update latest handoff current objective away from external proof-first;
- patch stale `docs/phase-plan.md` static-validator wording;
- preserve all truth boundaries.

### Phase 36 — OCI/OCL Reference Deployment Productization

Goal: make the existing OCI/OCL pilot server pattern a repeatable reference deployment.

Required work:

- add or refine `docs/deployment/oci-reference-deployment.md`;
- add env examples with placeholders only;
- document install/update/rollback;
- document service supervision;
- document Caddy/reverse proxy public/private routing;
- document validators on the server;
- document backup/restore/feed-monitor/scorecard workflows;
- add a server smoke checklist;
- add an operator support bundle command or instructions.

### Phase 37 — Reusable Agency Onboarding Flow

Goal: make a small agency able to run its own GTFS locally or on the reference server without manual DB surgery.

Required work:

- create a one-command or guided script for public GTFS onboarding;
- support agency ID, GTFS URL, metadata, and import timeout;
- produce clear output: feed URLs, admin URL, validator status, next steps;
- make local and server flows consistent;
- add tests for script argument validation where practical.

### Phase 38 — Integration Adapter Kit

Goal: make existing solution integration clean and documented.

Required work:

- turn the AVL adapter guidance into a reusable adapter kit;
- add sample mappings for common payload shapes without naming unsupported vendors;
- document how to connect external AVL through `POST /v1/telemetry`;
- document external predictor adapter lifecycle;
- add conformance tests or fixtures for adapter inputs/outputs;
- preserve fail-closed behavior.

### Phase 39 — CAL-ITP-Style Readiness Workflow In Product

Goal: make CAL-ITP-style readiness a product workflow inside the Operations Console and docs.

Required work:

- expose a readiness checklist in admin UI;
- show stable URL, metadata, validation, freshness, GTFS-RT completeness, and consumer packet status;
- add plain-language remediation steps;
- keep wording as “readiness support,” not compliance.

### Future Phase — Server-Based Regression Evidence

Goal: use the OCI/OCL server as the recurring proof target for code health.

Required work:

- define a recurring server smoke run;
- validate five feeds;
- run static and GTFS-RT validators;
- run backup/restore drill;
- run feed monitor;
- run telemetry dry-run or simulator;
- retain public-safe server evidence for regression tracking only.

## What Not To Do

Do not prioritize:

- final-root evidence;
- consumer submissions;
- external agency proof;
- public launch;
- compliance claims;
- production ETA claims.

Those can come later. The current project goal is reusable working open-source code.

## Definition Of Success

A small agency or civic technologist should be able to:

```text
1. clone the repo;
2. run the local app;
3. import their GTFS;
4. see five public feed paths;
5. run validators;
6. connect a telemetry source or adapter;
7. deploy the same pattern on the OCI/OCL reference server;
8. understand CAL-ITP-style readiness gaps;
9. reuse or modify the system without proprietary lock-in.
```
