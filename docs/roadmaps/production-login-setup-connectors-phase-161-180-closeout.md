# Phase 161-180 Production Login, Setup, Dashboard, And Connector Configuration Closeout

Date: 2026-05-21

This closeout summarizes the roadmap that added the missing product layer for
self-hosted agency administration: production username/password login,
first-admin setup links, dashboard-first navigation, a persistent setup wizard,
focused admin pages, and per-agency connector configuration state.

It records software capability only. It is not retained evidence, not a
release publication, not CAL-ITP/Caltrans compliance proof, not consumer
acceptance, not agency adoption, not hosted service availability, not vendor
compatibility, not hardware certification, not SLA or uptime coverage, not
production AVL reliability, not production-grade ETA quality, and not
final-root readiness.

## Completed Commits

- `7085d71` Phase 161: production auth boundary and revised roadmap baseline.
- `33713de` Phase 162: password credentials and first-admin bootstrap tokens.
- `dbe2c61` Phase 163: production username/password login and logout.
- `18dc466` Phase 164: admin user and role management.
- `19cda3d` Phase 165: auth status UX and future SSO placeholder.
- `8161188` Phase 166: dashboard top-issues redesign.
- `914120e` Phase 167: skippable first-run setup wizard.
- `28becbb` Phase 168: focused admin config pages.
- `1cc8cdd` Phase 169: connector instance model and truthful states.
- `82068bd` Phase 170: Vehicle / GPS / AVL connector setup wizard.
- `df0a8bc` Phase 171: connector dry-run jobs and redacted results.
- `e229240` Phase 172: connector activation gate tied to device bindings.
- `a2fd76e` Phase 173: prediction connector configuration review.
- `d05ee2a` Phase 174: validator connector configuration review.
- `9c5bfbe` Phase 175: monitoring and discovery connector configuration.
- `b7c5271` Phase 176: shorter admin pages and refined menu categories.
- `19b631d` Phase 177: setup reminders and dashboard health summaries.
- `2d41893` Phase 178: docs refresh for production login, setup, and connectors.
- `5878a6c` Phase 179: release gate rows for auth, setup, and connectors.

Phase 180 is this closeout, status, and handoff update.

## Product Capability Added

- Auth: production keeps `/admin/local-login` disabled, anonymous private
  Operations Console requests return `401`, old JWTs fail after secret
  rotation, new Bearer and `admin_session` JWTs work, and cookie-authenticated
  unsafe POSTs still require CSRF. Username/password login now uses DB-backed
  users, roles, Argon2id password hashes, first-admin setup links, reset links,
  failed-login handling, scoped `admin_session` cookies, and logout.
- Admin access: Users & Roles pages let admins create, disable, reset, and map
  roles for staff without direct SQL. Session identity, role, auth mode,
  password-login status, and deferred SSO status are visible in the private UI.
- Setup and dashboard: `/admin/operations` is dashboard-first, showing the top
  three issues or compact healthy summaries. The setup wizard has ten skippable
  steps, required/recommended/optional completion buckets, and a
  session-scoped reminder that reappears when the next blocker changes.
- Information architecture: private navigation is grouped as Dashboard, Setup,
  Data, Realtime, Connectors, Operations, and Admin. Settings moved into
  focused config pages for agency profile, public feed URLs, login, deployment,
  and advanced review.
- Connectors: examples are no longer treated as configured connectors.
  Per-agency connector instances use explicit states:
  `example_available`, `not_configured`, `configured_not_tested`,
  `dry_run_passed`, `ready_for_activation`, `active`, and `blocked`.
- Vehicle data connectors: Vehicle / GPS / AVL setup supports generic JSON
  transform, HTTP polling, webhook/sidecar, CSV replay, and vendor-shaped
  payload mapping through the adapter boundary. Metadata and secret reference
  labels are stored; secret values and live private endpoints are not.
- Connector safety gates: dry-run jobs record bounded redacted summaries and
  counts without raw payload retention. Vehicle connector activation requires
  mapping, a passed dry-run, device bindings, safe `/v1/telemetry` target
  shape, stale/future/quality rules, secret reference labels, and redaction
  scan status before marking ready for deployment-owned activation.
- Remaining connector family: Prediction Setup stores deterministic, external
  HTTP shadow, or external HTTP fail-closed metadata without enabling external
  prediction by default. Validator Setup stores allowlisted validator metadata
  without raw commands. Monitoring/export is no-send by default. Discovery
  reviews `/public/feeds.json`, license/contact, and public URL readiness
  without portal automation or consumer status mutation.
- Release gate: `make release-candidate-check` now includes rows for the
  production auth boundary, password login, first-admin bootstrap handling,
  logout, cookie CSRF rejection, dashboard issue priority, setup reminder,
  examples versus configured connector instances, and connector dry-run
  redaction.

## Validation Snapshot

Phase 179 ran the broad release gate after adding the new rows. The local
release-candidate diagnostic under `.cache` reported 64 passed rows, zero
blockers, one dirty-worktree `needs_review` row, and five intentionally
bounded follow-up rows: `validate`, `test`, `smoke`, local app five-feed fetch,
and release package audit. The first three follow-up commands were run
directly and passed. Local app startup and package audit remain opt-in release
review steps.

The following commands passed during the Phase 179 gate review:

- `git diff --check`
- `go test ./...`
- `make check`
- `make test`
- `make smoke`
- `make validate`
- `make release-candidate-check`
- `make test-release-candidate-check`
- `make check-links`
- `make product-ui-smoke`
- `make audit-product-language`
- `make audit-ui-layout`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `make audit-operations-route-inventory`
- `make external-connection-check`
- `make adapter-conformance`
- `make test-connector-examples`
- `make gtfsrt-conformance`
- `make api-contract-check`
- `scripts/check-consumer-tracker.sh`

Phase 180 reran a closeout validation set after status and closeout docs were
updated:

- `git diff --check`
- `make check`
- `make check-links`
- `make product-ui-smoke`
- `make test-release-candidate-check`
- `make audit-product-language`
- `make audit-final-claim-review`
- `scripts/check-consumer-tracker.sh`
- `OUTPUT_DIR=.cache/phase180-release-candidate-full FORCE=true scripts/release-candidate-check.sh`

The Phase 180 release-candidate diagnostic again reported 64 passed rows, zero
blockers, one dirty-worktree `needs_review` row, and five intentionally
bounded follow-up rows: `validate`, `test`, `smoke`, local app five-feed fetch,
and release package audit. The first three follow-up commands passed during
Phase 179 immediately before closeout.

## Protected Boundaries

No phase in this roadmap intentionally edited `docs/evidence/**`.
`docs/evidence/consumer-submissions/status.json` remained unchanged, and the
consumer tracker remained exactly seven `prepared` targets. The roadmap did
not contact agencies, vendors, consumers, portals, live IdPs, live validators,
or live connector services.

Secrets, private endpoints, DB URLs, raw telemetry payloads, private logs, and
raw connector dry-run payloads are not part of the browser configuration model
added here. Browser pages save metadata and deployment-owned reference labels;
actual secret provisioning and connector activation remain deployment-owned.

## Remaining SSO/OIDC Work

SSO/OIDC is not implemented. It remains moderate to high complexity because a
future phase must add provider configuration, issuer discovery,
redirect/callback handling, state and nonce, Authorization Code + PKCE, token
exchange, ID token validation, JWKS fetching and key rotation, claim mapping,
email verification handling, agency/role mapping, account linking, logout,
provider-specific group normalization, stub-provider tests, and replay,
mix-up, redirect URI, and token leakage defenses.

The intended architecture is unchanged: an OIDC/SSO provider should verify the
external identity, map it to an internal agency subject and roles, then issue
the same internal signed `admin_session` used by username/password login and
Bearer-token operator checks.

## Known Risks And Limits

- Password login is now productized, but deployment owners still need secure
  TLS, secrets management, backups, DB lifecycle, and operator procedures.
- Connector instances are configured metadata and readiness state, not proof
  of live vendor compatibility or production AVL reliability.
- The browser does not start, supervise, or execute external connector
  processes.
- Dry-run summaries are redacted diagnostics, not retained evidence packets.
- Prediction configuration does not prove production-grade ETA quality or
  real-world ETA accuracy.
- Validator configuration and passing local validators do not prove
  CAL-ITP/Caltrans compliance, consumer acceptance, or final-root readiness.

## Recommended Next Track

The recommended next software track is deeper realtime correctness: observed
arrival/departure evaluation, delay propagation, cancellation behavior,
frequency service, after-midnight service, block continuity, repeated trip
instances, and conservative Trip Updates quality measures.

If the maintainer wants to prioritize access integrations instead, SSO/OIDC is
the next identity track. If a deployment owner has authorized real connector
data, real connector runtime hardening can proceed under explicit
secrets/evidence boundaries. Release publication remains a separate maintainer
decision after release-package and optional local app checks are run from a
clean checkout.
