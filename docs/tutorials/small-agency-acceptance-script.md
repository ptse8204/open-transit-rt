# Small-Agency Acceptance Script

This script is for maintainers who want to walk the browser-first product path
exactly as a small agency, civic technologist, or developer integrator would see
it from a clean checkout.

This script creates no retained evidence, contacts no external party, changes no
consumer status, and makes no compliance, agency adoption, consumer acceptance,
final-root, hosted SaaS, production-readiness, vendor-compatibility, SLA, or
ETA-quality claim.

## Technical Helper Startup

The shell commands are for a technical helper who can clone the repo, run local
checks, and start Docker-backed services. They are not no-developer evaluator
steps.

```bash
git clone https://github.com/ptse8204/open-transit-rt.git
cd open-transit-rt
make check
make agency-app-up
```

The helper provides the private local browser URL and any local admin-token
instructions printed by `make agency-app-up`. No-developer evaluators start in a
private browser window at:

```text
http://localhost:8080/admin/operations
```

Stop the app at the end with:

```bash
make agency-app-down
```

## Browser Steps

| Step | Expected UI result | Common blocker | What to do next | CLI fallback | What this does not prove |
| --- | --- | --- | --- | --- | --- |
| Open the provided private local URL | The private console loads at `/admin/operations` and shows **Agency Operations Cockpit / Start Here** near the top. | Docker is not running, port `8080` is busy, the local proxy did not start, or local auth instructions were missed. | Ask the technical helper to run `make agency-app-logs` and `scripts/bootstrap-dev.sh --check`, then retry from the provided private URL. | `docker compose -f deploy/docker-compose.yml ps` | Local startup and console access do not prove hosted service availability, public launch, or agency approval. |
| Confirm the Operations Console path is `/admin/operations` | The private browser path remains `/admin/operations`; do not start from shell commands during no-developer review. | Browser is at the site root or another page. | Navigate to the provided `/admin/operations` URL. | `make agency-app-logs` | Reaching a private local route does not prove public deployment readiness. |
| Review Agency Operations Cockpit / Start Here | The page shows setup progress, primary action cards, current signals, next actions, and what each card does not prove. | Older checkout or browser cache. | Refresh and ask the technical helper to confirm the provided private URL is running the current checkpoint or evaluation build. | `/admin/operations.json` | Cockpit visibility does not prove adoption, compliance, hosted service, or production readiness. |
| Confirm Agency Operations Cockpit / Start Here tasks are visible | The Agency Operations Cockpit / Start Here path shows no-developer and developer paths, ordered first-run tasks, and five feed URLs. | Browser cache or an older checkout. | Refresh and ask the technical helper to confirm the current checkout is serving the provided private URL. | `go test ./cmd/agency-config -run FirstRun` | Visibility does not prove evidence intake or compliance. |
| Open setup wizard at `/admin/operations/setup-wizard` | The wizard lists agency profile, publication metadata, GTFS, feeds, telemetry, validators, connectors, and readiness. | Auth role is missing or agency query conflicts with the authenticated agency. | Remove conflicting `agency_id` query or use the local token for the seeded agency. | Review `/admin/operations/setup-wizard.json` if HTML rendering is unclear. | The wizard is a private guide, not agency approval. |
| Enter or review publication metadata | Setup pages show public base URL, feed base URL, license, contact, and environment state. | Read-only role cannot submit changes. | Use an admin role for changes, or keep the row marked missing/review-needed. | `make agency-app-up` seeds local demo metadata. | Metadata completeness does not prove final-root ownership. |
| Import GTFS through browser upload or safe URL import | `/admin/operations/gtfs-import` shows upload and URL import options for admins, plus validation feedback after import. | Read-only role, large file, private URL, invalid ZIP, or import validation errors. | Use an admin role, fix the ZIP/source URL, then review GTFS Quality. | Use the documented GTFS import CLI path for large or scripted imports. | Import success does not prove validator-clean status or agency approval. |
| Check `/public/feeds.json` | The local feed discovery document is reachable from the local app. | Local app not started or publication metadata missing. | Return to Agency Operations Cockpit / Start Here and Feed Health. | `curl http://localhost:8080/public/feeds.json` | Local discovery does not prove source-of-truth website listing. |
| Check `/public/gtfs/schedule.zip` | The local schedule ZIP endpoint returns a GTFS ZIP. | No active schedule feed or app startup failed. | Import/publish GTFS, then retry. | `curl -I http://localhost:8080/public/gtfs/schedule.zip` | A schedule URL does not prove open-license publication or consumer acceptance. |
| Check `/public/gtfsrt/vehicle_positions.pb` | The Vehicle Positions protobuf endpoint is reachable. | No telemetry, stale telemetry, or local realtime service not ready. | Review telemetry and feed health. | `curl -I http://localhost:8080/public/gtfsrt/vehicle_positions.pb` | Reachability does not prove real AVL reliability. |
| Check `/public/gtfsrt/trip_updates.pb` | The Trip Updates protobuf endpoint is reachable. | Prediction diagnostics missing or no active schedule/telemetry context. | Review feed health and Trip Updates quality rows. | `curl -I http://localhost:8080/public/gtfsrt/trip_updates.pb` | Reachability does not prove production-grade ETA quality. |
| Check `/public/gtfsrt/alerts.pb` | The Alerts protobuf endpoint is reachable. | No active alerts or local service not ready. | Review Alerts Console and feed health. | `curl -I http://localhost:8080/public/gtfsrt/alerts.pb` | Reachability does not prove consumer display or disruption-workflow completeness. |
| Open Feed Health | `/admin/operations/feed-health` shows feeds.json, schedule, Vehicle Positions, Trip Updates, and Alerts rows. | Feed metadata or validation records are missing. | Follow each row's next action. | `/admin/operations/feed-health.json` | Feed health rows are supporting signals only. |
| Review realtime usefulness | Feed Health explains whether Vehicle Positions, Trip Updates, and Alerts are empty, non-empty, stale, withheld, or unavailable where data exists. | Telemetry or diagnostics are not available yet. | Review Devices, Telemetry, Simulator, and Trip Updates diagnostics. | `make telemetry-simulator` when a technical helper is available. | Empty or non-empty feeds do not prove production-grade ETA quality or consumer display. |
| Open Devices & Tokens | `/admin/operations/devices` shows device bindings, token status without token values, vehicle binding, latest token use, and next actions. | No devices are configured, role is read-only, or token creation requires an admin. | Keep one-time token values out of the browser after creation and follow rotation or rebinding guidance. | No JSON fallback; use the private browser page and ask a technical helper for secure token storage or device installation. | Device rows do not prove real AVL integration, vendor compatibility, or hardware certification. |
| Open Telemetry Freshness | `/admin/operations/telemetry` shows latest accepted telemetry, stale state, assignment state, match confidence, or unknown reasons when available. | No telemetry has been sent, telemetry is stale, or the device token is missing. | Review Devices & Tokens, then use Telemetry Simulator's browser dry-run preview before any technical-helper shell send. | Review `/admin/operations/feed-health.json` for related feed impact. | Fresh local telemetry does not prove realtime coverage, production AVL reliability, or production matching quality. |
| Preview synthetic telemetry | `/admin/operations/telemetry-simulator` shows committed synthetic scenarios and a browser-only dry-run preview with redacted event summaries. | Fixtures are missing, the selected scenario is unknown, or an intentional local send is needed. | Choose a listed scenario and preview it in the browser; only a technical helper should run shell dry-runs or sends with private tokens. | `/admin/operations/telemetry-simulator.json?scenario=on-route` | Browser preview does not send telemetry, collect tokens, prove vendor compatibility, or prove production AVL reliability. |
| Open Maintenance Center | `/admin/operations/maintenance` shows active feed, import, five-feed check, validators, backup/restore configuration presence, telemetry freshness, service health availability, support summary, and weekly/monthly tasks. | Source records are not configured or not available. | Keep missing rows missing and follow the row-level next action. | `/admin/operations/maintenance.json` | Maintenance diagnostics do not prove SLA, uptime, or production readiness. |
| Open Readiness | `/admin/operations/readiness` shows private readiness rows and claim boundaries. | Missing source records remain missing. | Keep missing rows missing until the underlying source exists. | `/admin/operations/readiness.json` | Readiness review is not CAL-ITP/Caltrans compliance proof. |
| Open Connector Hub | `/admin/operations/connectors` shows connector categories and manifest registry guidance. | Connector examples or manifests are missing from the checkout. | Run local connector checks after the app walk. | `make external-connection-check` and `make adapter-conformance` | Connector conformance does not prove vendor compatibility. |
| Open Connector Tests | `/admin/operations/connectors/tests` shows synthetic manifest, sidecar, adapter, and conformance test guidance. | Connector examples or manifests are missing from the checkout. | Keep tests synthetic until a technical helper is ready to run local checks. | `/admin/operations/connectors/tests.json` | Connector test guidance does not prove vendor compatibility or real private-system integration. |
| Open Telemetry Simulator guide | `/admin/operations/telemetry-simulator` shows scenarios and commands. | No device token or simulator prerequisites. | Follow the page, then review Telemetry and Feed Health. | `make telemetry-simulator` | Synthetic telemetry does not prove real AVL proof. |
| Open Help & Tutorials | `/admin/operations/help` explains GTFS, GTFS-RT, connectors, readiness, validators, telemetry, Start Here, Devices & Tokens, Telemetry Freshness, Telemetry Simulator, Connector Checks, Maintenance Center, and claims/evidence. | None if the Operations Console is available. | Use Help links to return to the relevant UI page. | `/admin/operations/help.json` | Help is read-only guidance, not retained evidence. |
| Stop the app | Local containers stop. | Docker command fails or another process owns the services. | Run `make agency-app-logs`, then retry. | `make agency-app-down` | Stopping local services does not say anything about deployment readiness. |

## Technical Helper Checklist

Use a technical helper when a step requires:

- Docker or port troubleshooting;
- admin token handling;
- large GTFS import or validation triage;
- stable HTTPS deployment configuration;
- real device-token storage;
- a GPS/AVL transform into `POST /v1/telemetry`;
- local validator, connector, or release-candidate command execution;
- off-host validation with `make validate-public-feeds`;
- reference deployment diagnostics with `make oci-reference-check`;
- future authorized evidence intake.

## Close The Walkthrough

After stopping the app, record only command pass/fail results in the maintainer
handoff. Do not save screenshots, logs, feed payloads, or private diagnostics as
retained evidence unless a separate authorized evidence intake exists.
