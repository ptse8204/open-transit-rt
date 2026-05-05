# Demo Video Outline

This is a draft outline for a short review video. It is not proof of production deployment, consumer submission, consumer acceptance, agency endorsement, CAL-ITP/Caltrans compliance, hosted SaaS availability, or public launch.

## Goal

Show how a reviewer can try Open Transit RT locally and understand the evidence boundary before considering an agency pilot.

## Suggested Length

6 to 9 minutes.

## Script

### 1. Opening

- State that Open Transit RT is open-source tooling for GTFS, telemetry ingest, and GTFS Realtime publication.
- State that the demo runs locally from repo files.
- State the boundary: local demo evidence is not production proof or consumer acceptance.

### 2. Local App Startup

Show:

```bash
make agency-app-up
```

Narration:

- The command starts the local app package behind `http://localhost:8080`.
- It applies migrations, seeds local data, imports the committed sample GTFS, bootstraps publication metadata, and prints next steps.
- This is local evaluation tooling, not hosted SaaS.

### 3. GTFS Import Or Demo Feed

Show either:

- the local package importing the committed sample GTFS; or
- the documented real GTFS onboarding path in [Real Agency GTFS Onboarding](tutorials/real-agency-gtfs-onboarding.md).

Narration:

- The demo feed helps reviewers understand the workflow.
- Real agency data needs permission, validation review, metadata approval, and redaction before evidence is retained.

### 4. Public Feed URLs

Show the printed local feed URLs:

- `/public/gtfs/schedule.zip`
- `/public/feeds.json`
- `/public/gtfsrt/vehicle_positions.pb`
- `/public/gtfsrt/trip_updates.pb`
- `/public/gtfsrt/alerts.pb`

Narration:

- Public protobuf feed endpoints are part of the local app surface.
- JSON debug and admin views remain protected.
- Local URLs are not agency-owned final feed roots.

### 5. Operations Console Setup Checklist

Show `/admin/operations/setup`.

Narration:

- The checklist summarizes publication metadata, feed discovery, validation records, device bindings, telemetry state, and evidence links.
- It helps operators find gaps; it does not certify readiness.

### 6. Device Telemetry Or Dry-Run Adapter Path

Show one of:

```bash
./scripts/device-onboarding.sh sample-telemetry
```

or:

```bash
go run ./cmd/avl-vendor-adapter -input testdata/avl-vendor/valid.json -mapping testdata/avl-vendor/mapping.json
```

Narration:

- The device path shows authenticated telemetry ingest for the local stack.
- The AVL/vendor adapter command is synthetic dry-run transform evidence only.
- It does not prove real vendor support, certified hardware support, or production AVL reliability.

### 7. Validation And Evidence View

Show:

- validation records or scorecard views in the Operations Console;
- [Compliance Evidence Checklist](compliance-evidence-checklist.md);
- [California Readiness Summary](california-readiness-summary.md).

Narration:

- Validation and scorecard records help reviewers understand current behavior.
- Compliance and readiness claims require deployment-specific and external evidence not present in the repo today.

### 8. Consumer Packet Boundary

Show:

- [Consumer Submission Tracker](evidence/consumer-submissions/README.md);
- [Consumer Status JSON](evidence/consumer-submissions/status.json);
- prepared packet directory links.

Narration:

- Packets exist for operator review.
- All seven consumer and aggregator targets remain prepared only.
- No target has submitted, under-review, accepted, rejected, blocked, ingestion, listing, display, or adoption evidence.

### 9. Pilot Package Next Step

Show:

- [Agency One-Pager](agency-one-pager.md);
- [Agency Pilot Program](agency-pilot-program.md);
- [Agency Pilot Checklist](agency-pilot-checklist.md);
- [Agency Feedback Template](agency-feedback-template.md).

Narration:

- The next useful step is an agency evaluation with public-safe retained evidence.
- Stronger claims should wait for real retained evidence, such as agency-owned final-root proof, authorized target-specific consumer submission evidence, real agency pilot evidence, or real deployment operations evidence.

### 10. Closing

- Invite reviewers to star or bookmark the repo if useful.
- Invite contributors to review docs, fixtures, runbooks, adapter examples, and public-safe evidence wording.
- Repeat that this outline is draft launch material only and does not mean an announcement was posted or a public launch occurred.
