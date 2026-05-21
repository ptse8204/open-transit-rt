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

## Closeout

Phases 161-180 are complete. The closeout is
`docs/roadmaps/production-login-setup-connectors-phase-161-180-closeout.md`.

The completed roadmap added username/password login, first-admin setup links,
DB-backed user and role management, session/status UX, a dashboard-first
Operations Console, a persistent skippable setup wizard, focused config pages,
truthful connector instance state, Vehicle / GPS / AVL connector setup,
redacted dry-run review, activation readiness gates, prediction/validator/
monitoring/discovery connector configuration, shorter navigation, and release
gate rows for the new auth/setup/connector risks.

SSO/OIDC remains deferred and should be implemented later as an identity source
that verifies external identity, maps it to internal agency subjects and roles,
and issues the same internal signed `admin_session`. The roadmap does not
prove compliance, production readiness, consumer acceptance, agency adoption,
final-root readiness, hosted service availability, vendor compatibility, SLA
coverage, production AVL reliability, production-grade ETA quality, or
real-world ETA accuracy.
