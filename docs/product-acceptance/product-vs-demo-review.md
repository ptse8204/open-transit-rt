# Product Versus Demo Review

Date: May 18, 2026

This review checks whether the Product Reality Pass improved real Open Transit
RT workflows, not only the local demo. It is based on the current code paths for
local sign-in, admin auth, the Operations Console, public feed routes, telemetry
ingest, validator wrappers, connector examples, and the public website.

This record is product acceptance context only. It is not retained external
evidence and does not prove production readiness, compliance, consumer
acceptance, agency approval, hosted service availability, vendor compatibility,
hardware certification, production AVL reliability, or ETA quality.

## Short Answer

The browser-first work now changes real product surfaces:

- `/admin/operations` uses the same private Operations Console in local and
  self-hosted deployments.
- GTFS import, feed checks, realtime review, device credentials, connector
  review, validator health, maintenance, and help are real authenticated admin
  routes.
- Public feed URLs are implemented routes, not static website mockups.
- Telemetry still enters through the real authenticated `POST /v1/telemetry`
  boundary.
- Connector support remains bounded to local-supported examples, wrappers,
  sidecars, and conformance checks unless a deployment owner wires a real
  external source.

The local sign-in handoff is deliberately local-demo only. A self-hosted
deployment still needs deployment-owned admin access, secrets, network controls,
TLS, public feed root configuration, telemetry source setup, and operations
monitoring.

## Flow Review

| Flow | Local demo | Self-hosted deployment | What this means |
| --- | --- | --- | --- |
| Local admin sign-in | Works at `/admin/local-login` when `LOCAL_ADMIN_LOGIN_ENABLED=true` and the host is localhost. | Not exposed in production. A deployment must provide its own admin access path and secrets. | Normal local users can enter the browser without tokens after startup. Production auth is not weakened. |
| Operations Console workflow | Works after local sign-in with the seeded demo agency. | Works after authenticated admin access with the configured agency. | The straight-line workflow is product UI, not a website-only demo. |
| Start setup | Shows setup, feed URL, validation, and next-action state from configured records. | Works when publication metadata, agency records, and backing services are configured. | It helps users see what to do next; deployment setup is still owned outside the browser. |
| Import GTFS | Browser upload and safe URL import use the same validation and publish path as the backend importer. | Works when storage, database, size limits, and network access are configured. | This is real product behavior. It is not a static guide. |
| Check feeds | Shows `/public/feeds.json`, static GTFS, Vehicle Positions, Trip Updates, and Alerts state. | Works with a configured public base URL and active feed artifacts. | It reviews current feed availability but does not prove consumer ingestion. |
| Connect vehicles | Uses device bindings, token rotation, simulator guidance, and telemetry freshness pages. | Works with deployment-owned device tokens and telemetry sources. | Real devices still require deployment-owned AVL or connector setup. |
| Telemetry ingest | Local demo seeds a sample device binding and can use synthetic/local payloads. | `POST /v1/telemetry` requires a valid device token and binding. | The ingest boundary is real and protected. Vendor payload adapters remain deployment-owned unless implemented as tested sidecars. |
| Realtime review | Reviews local Vehicle Positions, Trip Updates, Alerts, assignment, freshness, and prediction signals. | Reviews the same records for the configured deployment. | It improves operations visibility; it does not prove field reliability or ETA accuracy. |
| Connector support | Public and admin pages list supported local examples and conformance checks. | Deployment owners can run sidecars, transforms, validators, and monitoring helpers with their own credentials and endpoints. | The product documents the boundary plainly without claiming vendor compatibility. |
| Validator health | Can review or run allowlisted validator checks when tooling is installed. | Same path works when server-owned validator tooling is installed and configured. | Browser requests cannot provide validator commands, paths, URLs, or arbitrary arguments. |
| Share public URLs | Shows and copies configured public feed URLs. | Works when DNS, TLS, base URL, and feed publication are configured. | It helps prepare sharing but does not submit to or change any consumer portal. |
| Maintenance and help | Shows local state, recovery prompts, support bundle guidance, and glossary help. | Works against configured deployment records. | It reduces operator reading but does not replace deployment operations ownership. |

## Local-Demo Only

These flows are intentionally limited to the local evaluation package:

- `/admin/local-login`, guarded by `LOCAL_ADMIN_LOGIN_ENABLED=true`,
  `!appconfig.IsProduction()`, a localhost host check, signed one-time state,
  and a short `admin_session` cookie.
- `scripts/agency-local-app.sh` startup output, seeded demo agency, seeded
  sample GTFS, seeded local device binding, and local browser setup URL.
- Tutorial screenshots and the public interface walkthrough, which use safe
  local/demo data.
- Local proxy assumptions around `http://localhost:8080`.
- Any demo simulator path or local sample payload used to make Vehicle
  Positions visible during evaluation.

## Deployment-Owned Setup

A real self-hosted deployment still needs a technical owner for:

- DNS, TLS, reverse proxy, and public feed root configuration.
- Admin access, secrets, role bindings, cookie/session policy, and network
  controls.
- Database provisioning, migrations, backups, restore drills, and retention.
- Validator installation and pinned server-side validator mappings.
- Real telemetry source setup, device tokens, sidecar hosting, and credential
  rotation.
- Monitoring destinations, alert routing, uptime checks, log retention, and
  operational escalation.
- Public data governance such as license/contact metadata and a decision about
  whether to share URLs with consumers.

## Code Work Still Needed

These are product gaps, not documentation gaps:

- A production admin onboarding flow, such as deployment-owned SSO, invitations,
  or a managed admin user setup path.
- A browser-native deployment setup review that explains missing DNS, TLS,
  secrets, validators, feed base URL, and monitoring without requiring users to
  inspect environment variables.
- Stronger GTFS repair assistance after import and validation, especially for
  calendar, shape, stop-time, block, and realtime matching issues.
- A clearer validator report explanation layer for nontechnical operators.
- Connector runtime management that can show sidecar status and setup progress
  without exposing arbitrary browser command execution.
- Better self-hosted maintenance workflows for backup status, restore drills,
  retention, and alerting.
- Consumer discovery preparation that remains local/prepared-only until a
  deployment owner chooses to contact or submit to a consumer.

## Unsupported Claims

The following remain unsupported and must stay out of product copy unless they
become implemented and verified:

- CAL-ITP or Caltrans compliance.
- Production readiness.
- Agency adoption, approval, or endorsement.
- Consumer submission, review, acceptance, ingestion, listing, or display.
- Final-root ownership or final public launch readiness.
- Hosted service availability, paid support, SLA coverage, or uptime guarantee.
- Vendor compatibility or hardware certification.
- Production AVL reliability.
- Production-grade ETA quality or real-world ETA accuracy.

## Verdict

The correction pass now improves the actual product experience: local users can
sign in from the browser after startup, the Operations Console follows a
straight-line workflow, the website shows real rendered UI, the public pages
carry the main user path, and connector support is split into current versus
planned categories.

The remaining boundary is clear: Open Transit RT supports local and self-hosted
evaluation with browser-first operations, but production deployment, real AVL
source integration, consumer outreach, compliance, and public launch decisions
remain deployment-owned work.
