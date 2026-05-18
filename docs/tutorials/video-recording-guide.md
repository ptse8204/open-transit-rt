# Video Recording Guide

This guide defines repeatable, public-safe recording workflows for Open Transit
RT tutorials. It is for maintainers and technical helpers who want to record
short demos from the local browser-first product flow.

Do not commit large binary video files to the repository unless a maintainer
explicitly authorizes it. Store drafts outside the repo. Publish finished
videos only through an authorized release asset, documentation hosting, or
project media location.

## Safety Rules

- Use only local demo data, public GTFS data, or synthetic fixtures.
- Do not show bearer tokens, passwords, database URLs, private URLs, private
  portal records, private tickets, raw support bundles, raw logs with secrets,
  private vehicle identifiers, or unredacted operator artifacts.
- Do not show real agency logos, real staff names, or real device/vendor
  payloads without explicit permission and license review.
- Do not contact consumers, agencies, vendors, or portals during recording.
- Do not create retained evidence, change consumer tracker status, publish a
  release, push a site update, or make production/compliance/consumer claims.
- Add captions or a transcript before publication.

## Before Recording

1. Start from a clean local demo state or a controlled reference environment.
2. Start the local app:

   ```bash
   make agency-app-up
   ```

3. Open the local browser sign-in page:

   ```text
   http://localhost:8080/admin/local-login
   ```

4. Select **Start setup** and confirm the browser lands on the private
   Operations Console.
5. Confirm the visible data is public-safe.
6. Keep any admin token off screen.
7. Close browser tabs that show private paths, raw logs, credentials, or
   unrelated local work.

## Published Tutorial Asset

The public site includes a small MP4 tutorial asset and captions generated
from actual local UI screenshots:

```text
site/assets/open-transit-rt-interface-tour.mp4
site/assets/open-transit-rt-interface-tour.vtt
```

The asset was generated locally with `ffmpeg` from Playwright screenshots of
the local/demo browser UI. It does not capture a real desktop profile, secret,
private URL, raw log, agency record, external portal, consumer workflow, or
vendor payload.

Local recording tooling detected for future authorized screen captures:

- `ffmpeg`
- `ffprobe`
- macOS `screencapture`
- macOS `osascript`

Direct desktop capture was not used for this committed asset because the local
browser screenshot workflow is safer and avoids accidentally recording private
desktop context. If a maintainer later authorizes a narrated screen recording,
use the checklist in this guide, record only local/demo data, and publish
through an authorized release or documentation workflow.

### Published Asset Transcript

Open the local sign-in page after the app starts. Choose **Start setup**; no
tokens or browser extensions are needed for normal review.

Start on the Operations Console. The first page shows the current state and the
next action to take.

Import or review GTFS. Confirm the active schedule, required files, import
history, and validation guidance.

Check the five feed URLs: `feeds.json`, static GTFS, Vehicle Positions, Trip
Updates, and Alerts.

Connect vehicle data through device setup and authenticated telemetry paths.
Stored token values are not shown.

Review realtime state: fresh vehicles, stale rows, Vehicle Positions, Trip
Updates diagnostics, and Alerts.

Choose connector support from local-tested paths before any real endpoint,
credential, or vendor payload is configured.

Review what needs attention and decide whether the next step belongs to staff
or a technical helper.

Use Maintenance and Help for routine checks, backup and restore guidance,
support summaries, setup, tooling, secrets, connectors, and deployment work.

This tutorial is local demo only. It does not prove CAL-ITP/Caltrans
compliance, production readiness, consumer acceptance, vendor compatibility,
hardware certification, SLA/uptime, production AVL reliability, production ETA
quality, or real-world ETA accuracy.

## Reset After Recording

Use the least destructive reset that fits the session:

```bash
make agency-app-down
```

When a clean local demo database is required for another take:

```bash
make agency-app-reset
make agency-app-up
```

`make agency-app-reset` is destructive for local demo containers, volume state,
and logs. Do not run it against a deployment.

## Storyboard 1: Three-Minute Overview

Goal: explain what the project does and what local review does not prove.

| Time | Show | Say |
| --- | --- | --- |
| 0:00-0:25 | Homepage or README | Open Transit RT is a self-hosted backend for GTFS and GTFS Realtime evaluation. |
| 0:25-0:55 | `/admin/operations` | Agency staff start from the browser after a helper starts the app. |
| 0:55-1:25 | Feed paths | The local app exposes `feeds.json`, static GTFS, Vehicle Positions, Trip Updates, and Alerts. |
| 1:25-1:55 | Connector Hub | Vehicle/GPS/AVL, prediction, validator, monitoring, and discovery connectors stay explicit. |
| 1:55-2:25 | Readiness | CAL-ITP-style readiness helps prepare local review, but does not prove compliance. |
| 2:25-3:00 | What is not proven | Local evaluation does not prove production readiness, consumer acceptance, vendor compatibility, SLA/uptime, or ETA quality. |

## Storyboard 2: Ten-Minute Local Setup

Goal: show a technical helper starting the app and handing off browser review.

| Time | Show | Say |
| --- | --- | --- |
| 0:00-0:45 | README startup block | Use the current release candidate and local commands. |
| 0:45-2:30 | Terminal with `make check` | The check is a local repo sanity check, not external approval. |
| 2:30-5:30 | Terminal with `make agency-app-up` | The helper starts local services, imports demo GTFS, and prints next steps. |
| 5:30-6:15 | Local sign-in page | Open `/admin/local-login` and select **Start setup**. |
| 6:15-8:30 | `/admin/operations` | Agency staff start from the private browser session. |
| 8:30-10:00 | Stop/reset guidance | Use `make agency-app-down`; use reset only for local demo cleanup. |

## Storyboard 3: Browser-First GTFS Import

Goal: show GTFS import and quality review without using command-line-first
operations.

| Time | Show | Say |
| --- | --- | --- |
| 0:00-0:45 | `/admin/operations/gtfs-import` | Admins can upload GTFS or use a safe URL import path. |
| 0:45-2:30 | Import result | Import status and validation feedback are shown in the browser. |
| 2:30-4:30 | `/admin/operations/gtfs-workbench` | Staff review required files, row counts, service dates, and active feed version. |
| 4:30-5:30 | GTFS quality triage | Issues are grouped by likely owner, meaning, suggested fix path, and next action. |
| 5:30-6:00 | Boundary | Import success does not prove agency approval or validator-clean public deployment. |

## Storyboard 4: Feed Health And Readiness Review

Goal: explain public feed review and CAL-ITP-style readiness without claiming
compliance.

| Time | Show | Say |
| --- | --- | --- |
| 0:00-1:30 | `/admin/operations/feed-health` | Review all five feed paths and local health context. |
| 1:30-2:45 | `/admin/operations/validation-health` | Review validator records, stale state, and blockers. |
| 2:45-4:45 | `/admin/operations/readiness` | Read the ten readiness areas and next actions. |
| 4:45-5:30 | `/admin/operations/consumers` | Draft packet records do not prove submission, review, acceptance, ingestion, listing, or display. |
| 5:30-6:00 | Boundary | Stronger public statements require deployment-specific and external evidence. |

## Storyboard 5: Connector And AVL Overview

Goal: show connector categories and the authenticated telemetry boundary.

| Time | Show | Say |
| --- | --- | --- |
| 0:00-1:00 | `/admin/operations/connectors` | Pick the closest connector category before integration work. |
| 1:00-2:30 | `/admin/operations/connectors/workbench` | Review recipes, redaction templates, examples, and conformance coverage. |
| 2:30-3:45 | `/admin/operations/telemetry-simulator` | Preview synthetic telemetry scenarios before intentional sends. |
| 3:45-4:30 | `/admin/operations/devices` | Device token creation is admin-only and token values stay off screen. |
| 4:30-5:00 | Boundary | Synthetic examples do not prove vendor compatibility, hardware certification, or production AVL reliability. |

## Storyboard 6: Maintenance And Support Workflow

Goal: show routine operations and troubleshooting paths.

| Time | Show | Say |
| --- | --- | --- |
| 0:00-1:30 | `/admin/operations/maintenance` | Review weekly/monthly tasks and technical-helper cases. |
| 1:30-2:30 | `/admin/operations/help` | Use task-based help when labels or next actions are unclear. |
| 2:30-3:30 | Support bundle guide | Support bundles are local diagnostics, not retained evidence. |
| 3:30-4:30 | Update guidance | Use documented release and deployment guidance before updates. |
| 4:30-5:00 | Boundary | Maintenance checks do not prove SLA, uptime, hosted service availability, or production readiness. |

## Publication Checklist

- Verify the video uses only public-safe local/demo data.
- Add captions or transcript.
- Add a short description with unsupported-claim boundaries.
- Store the source video outside the repo.
- If publishing through release assets or docs hosting, record the authorized
  location and checksum in maintainer notes outside protected evidence paths
  unless a separate evidence workflow authorizes retention.

## Static Tutorial Page

The public-site source includes a companion page:

```text
site/video.html
```

It embeds the screenshot-based MP4, captions, and public transcript. Maintainer
recording checklists and storyboards stay in this guide, not on the public
page.

## What This Does Not Prove

Recording or publishing a tutorial does not prove CAL-ITP/Caltrans compliance,
production readiness, agency adoption or approval, consumer submission, review,
acceptance, ingestion, listing, or display, final-root readiness, hosted
service availability, vendor compatibility, hardware certification, SLA/uptime,
production AVL reliability, production-grade ETA quality, or real-world ETA
accuracy.
