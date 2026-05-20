# Phase 161-180 Production Login, Setup, Dashboard, And Connector Configuration

Date started: 2026-05-20

This roadmap moves the private Operations Console from local/demo sign-in
toward a self-hosted production access model and shorter operator workflows.
It is a software capability track only. It is not compliance proof, consumer
acceptance, agency adoption, hosted service availability, vendor
compatibility, SLA/uptime coverage, production AVL reliability, production ETA
quality, or final-root evidence.

## Direction

- Keep `/admin/local-login` local/demo-only and production-disabled.
- Add username/password login first, backed by database users, roles,
  password credentials, and one-time bootstrap/reset links.
- Keep SSO/OIDC deferred. A future provider should verify external identity,
  map it to an internal agency subject and roles, then issue the same signed
  `admin_session` cookie.
- Make `/admin/operations` dashboard-first, with the top issues visible before
  diagnostic tables.
- Keep setup wizard progress skippable but persistently visible until complete.
- Split settings and advanced diagnostics into focused pages.
- Separate connector examples from configured per-agency connector instances.
- Prioritize connector configuration in this order: Vehicle/GPS/AVL,
  prediction, validators, monitoring/export, and discovery/sharing prep.

## Phase 161 Baseline

Phase 161 locks the current production auth boundary as executable regression
coverage before adding browser password login:

- production disables `/admin/local-login`;
- anonymous `/admin/operations` remains `401`;
- pre-rotation JWTs are rejected when old secrets are not configured;
- new Bearer admin JWTs are accepted;
- new `admin_session` cookies are accepted;
- unsafe cookie-authenticated POSTs without CSRF remain `403`;
- the reference public edge exposes public feed paths only and does not expose
  admin routes;
- public feed endpoints remain anonymous.

The local demo sign-in behavior remains available outside production when
explicitly enabled on loopback.
