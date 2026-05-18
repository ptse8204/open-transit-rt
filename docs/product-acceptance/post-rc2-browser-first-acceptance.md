# Post-rc2 Browser-First Product Acceptance

Date: 2026-05-18 UTC / 2026-05-18 Pacific

This acceptance review checks whether a small-agency user can use the browser
for normal Operations Console work after a technical helper starts the local
app.

Open Transit RT `v0.1.0-rc.2` remains a public release candidate for
local/self-hosted evaluation. This review does not prove production readiness,
CAL-ITP/Caltrans compliance, consumer acceptance, agency adoption, hosted
service availability, vendor compatibility, hardware certification, SLA/uptime,
production AVL reliability, production-grade ETA quality, or real-world ETA
accuracy.

## Web Design Engineer Skill

The Web Design Engineer skill was loaded before writing this acceptance page.
It shaped the record toward a user-journey walkthrough instead of a raw route
dump: the report starts with whether the browser task flow works, summarizes
next actions and remaining command-line boundaries, and keeps internal
diagnostic detail in the local `.cache/product-acceptance/` artifacts.

## Local Startup

Command:

```bash
make agency-app-up
```

Result: passed.

The local evaluator stack started at `http://localhost:8080`, reused the
existing local demo database volume, applied migrations, seeded demo agency
records, imported the sample GTFS feed, bootstrapped publication metadata, and
verified all five public feed files.

## Browser-First Access Correction

The local/demo app now exposes a localhost-only admin sign-in handoff at
`/admin/local-login` when explicitly enabled by the local app configuration.
Normal users open that page, choose **Start setup**, and receive a short browser
session for the private console. They do not need to generate an admin token,
paste a cookie, use DevTools, install a header extension, or call curl.

The correction preserves the private boundary: unauthenticated
`/admin/operations` still returns `401`, cookie-auth unsafe POSTs still require
CSRF, and the local sign-in route is disabled for production mode.

## Action-First Console Correction

The Operations Console Start page now leads with a next-action queue and short
current-status cards. The route groups are:

- Start
- Setup
- GTFS
- Feeds
- Realtime
- Vehicles
- Connectors
- Readiness
- Maintenance
- Help

Longer first-run tables, feed URL details, console inventory, and safety flags
remain available behind details panels for review and audit without forcing a
normal user to read them before choosing an action.

## Browser Route Walkthrough

The walkthrough generated an admin token in memory and did not print it. These
routes returned authenticated `200` and `Cache-Control: no-store`:

| Workflow | Route | Page |
| --- | --- | --- |
| Start Here / Dashboard | `/admin/operations` | Start Here |
| Setup | `/admin/operations/setup-wizard` | Agency Setup |
| Setup Details | `/admin/operations/setup` | Setup Details |
| GTFS Workbench | `/admin/operations/gtfs-workbench` | GTFS Workbench |
| GTFS Import | `/admin/operations/gtfs-import` | Import GTFS |
| Feed URLs | `/admin/operations/feeds` | Feed URLs |
| Feed Health | `/admin/operations/feed-health` | Feed Health |
| Validation | `/admin/operations/validation-center` | Validation Center |
| Validator Health | `/admin/operations/validation-health` | Validator Health |
| Realtime | `/admin/operations/realtime` | Realtime |
| Devices / AVL | `/admin/operations/devices` | Devices & Tokens |
| Telemetry | `/admin/operations/telemetry` | Telemetry Freshness |
| Telemetry Simulator | `/admin/operations/telemetry-simulator` | Telemetry Simulator |
| Connectors | `/admin/operations/connectors` | Connector Hub |
| Connector Workbench | `/admin/operations/connectors/workbench` | Connector Workbench |
| Connector Checks | `/admin/operations/connectors/tests` | Connector Checks |
| Prediction Lab | `/admin/operations/prediction-lab` | Prediction Lab |
| Readiness | `/admin/operations/readiness` | Readiness |
| Maintenance | `/admin/operations/maintenance` | Maintenance |
| Reliability Review | `/admin/operations/reliability` | Reliability Review |
| Help | `/admin/operations/help` | Help |
| Support Audit | `/admin/operations/audit` | Audit History |
| Evidence Guidance | `/admin/operations/evidence` | Evidence Guidance |

Additional browser surfaces returned `200`:

| Surface | Route |
| --- | --- |
| GTFS Studio | `/admin/gtfs-studio` |
| Alerts | `/admin/alerts/console` |

The unauthenticated check returned `401` for `/admin/operations`, preserving the
private admin boundary.

## Post-rc2 Polish Recheck

On 2026-05-17, the local app was rebuilt and the browser-first route walkthrough
was repeated after copy and responsive-layout polish. The recheck confirmed:

- all 23 Operations Console routes listed above returned authenticated `200`;
- GTFS Studio and Alerts returned authenticated `200`;
- unauthenticated `/admin/operations` returned `401`;
- a local Playwright run loaded all 25 private browser routes at desktop
  `1280x900` and
  mobile `390x844` viewports, for 50 checked route/viewport pairs;
- the local browser sweep reported no missing page title/H1, page-level
  horizontal overflow on the checked routes, console errors, stale internal route
  wording, or favicon 404s;
- the normal operator path now uses agency workflow wording instead of
  release/phase planning phrasing.

This recheck remains local/self-hosted evaluation only. It does not create
external evidence, contact external systems, move consumer status, or prove any
unsupported compliance, production, consumer, vendor, SLA, hardware, AVL, or ETA
claim.

## Public Feed Checks

Startup and follow-up fetches returned `200` for:

- `/public/feeds.json`
- `/public/gtfs/schedule.zip`
- `/public/gtfsrt/vehicle_positions.pb`
- `/public/gtfsrt/trip_updates.pb`
- `/public/gtfsrt/alerts.pb`

## Browser-First Result

After a technical helper starts the local app, normal users can use the browser
to:

- sign in through the local/demo browser handoff;
- review the dashboard and next actions;
- complete setup review;
- import and review GTFS;
- inspect feed URLs and feed health;
- review validation and validator health;
- inspect realtime Vehicle Positions, Trip Updates, and Alerts status;
- review devices, tokens, telemetry freshness, and simulator guidance;
- review connector categories, workbench recipes, and connector checks;
- review CAL-ITP-style readiness without overclaiming;
- inspect maintenance, reliability, help, support, audit, and evidence guidance.

## Remaining Command-Line Boundaries

The browser is not allowed to run arbitrary commands. These remain technical
helper or maintainer work:

- starting and stopping the local app package;
- installing pinned validators;
- release packaging and release-candidate gates;
- Docker, TLS, DNS, reverse proxy, and deployment monitoring setup;
- secure storage and rotation of real device secrets;
- evidence retention, external contact, consumer submission, or consumer status
  movement.

## Local Artifacts

Local diagnostic output was written under `.cache/product-acceptance/`:

- `phase14-route-walkthrough.csv`
- `phase14-auth-boundary.txt`
- `phase14-external-and-public-routes.txt`

These local files are not retained evidence and are not committed.
