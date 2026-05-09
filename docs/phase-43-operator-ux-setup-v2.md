# Phase 43 — Operator UX Setup V2

## Status

Complete for the private Operations Console setup/readiness UX scope.

## What Changed

Phase 43 first patched a Phase 42/local routing issue:

- the deployment doctor now checks `/admin/gtfs-studio` instead of the retired
  exact `/admin/gtfs` route;
- `deploy/Caddyfile.local` keeps the exact local root page at `/` returning
  `200`;
- matched public, admin, and service proxy routes keep their existing
  behavior;
- unmatched local paths now return `404` instead of the local root message.

The authenticated Operations Console now includes a private operator checklist:

- HTML: `/admin/operations/checklist`
- JSON: `/admin/operations/checklist.json`

Both routes are admin-authenticated only, use `Cache-Control: no-store`, reject
conflicting `agency_id` query values, and are derived from one shared internal
checklist model.

The checklist has deterministic group order:

1. `setup`
2. `feeds`
3. `validation`
4. `telemetry`
5. `operations`
6. `consumer_workflow`

Each row has a stable ID, neutral status, source, current signal, next action,
claim boundary, repo-relative docs links, and heuristic labels. JSON statuses
are limited to `ok`, `needs_review`, `missing`, `blocked`, and `unknown`.

## Boundaries

This checklist is private operator diagnostics. It is not evidence, not an
evidence packet, not compliance proof, not agency approval, not consumer
acceptance, and not production readiness.

Phase 43 did not:

- add public checklist routes;
- add schema migrations or approval flags;
- create evidence packets;
- read `.cache` diagnostics as evidence;
- contact consumers;
- change consumer statuses;
- claim CAL-ITP/Caltrans compliance;
- claim agency adoption, consumer acceptance, hosted SaaS, production
  readiness, vendor compatibility, or production-grade ETA quality.

## JSON Flags

The JSON export includes these explicit false flags:

- `external_evidence_created`
- `final_root_evidence_created`
- `consumer_statuses_changed`
- `compliance_claimed`
- `production_readiness_claimed`
- `agency_approval_claimed`
- `consumer_acceptance_claimed`

All are always `false` in Phase 43.

## Heuristic Labels

The JSON export uses repo-safe underscore labels only:

- `missing`
- `placeholder_like`
- `operator_entered_unverified`
- `approval_unknown`
- `approval_artifact_not_retained`
- `local_only`
- `pilot_or_reference_root`
- `final_root_candidate_unverified`
- `no_final_root_evidence`

The HTML view renders human labels, but it does not emit approval,
compliance, consumer readiness, or production readiness labels.

## Verification

The Phase 43 implementation adds focused tests for:

- checklist route registration and method rejection;
- role access for `read_only`, `operator`, `editor`, and `admin`;
- unauthenticated rejection;
- agency scoping;
- JSON schema shape, deterministic group/row ordering, row IDs, and false
  flags;
- HTML and JSON safety against raw secrets, tokens, cache paths, raw
  telemetry payloads, and private filesystem paths;
- HTML escaping of script-like metadata;
- metadata and URL heuristic classifiers;
- setup/readiness/dashboard navigation to the checklist routes;
- deployment-doctor route regression from `/admin/gtfs` to
  `/admin/gtfs-studio`;
- local Caddy fallback shape.

The consumer tracker remains prepared-only for Google Maps, Apple Maps,
Transit App, Bing Maps, Moovit, Mobility Database, and transit.land.
