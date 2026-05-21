# Production Login, Setup, And Connectors

This guide describes the current self-hosted browser path after the app is
deployed behind a private admin boundary. It is not hosted service
documentation, SSO documentation, consumer-submission evidence, compliance
evidence, production-readiness proof, vendor compatibility proof, or
production AVL/ETA reliability proof.

## 1. Create The First Admin

In production-style deployments, `/admin/local-login` is disabled. A deployment
owner creates the first admin from an operator shell:

```bash
go run ./cmd/agency-config bootstrap-admin-link \
  --agency-id agency \
  --email admin@example.org \
  --base-url https://admin.example.org \
  --ttl 30m
```

The command creates the user and role binding, stores only a hashed bootstrap
token, and prints the setup URL once. Do not paste that URL into tracked docs,
evidence folders, issue comments, screenshots, or public logs.

Open the setup link through the private admin access path, set the password,
then sign in at:

```text
/admin/login
```

Successful password login issues the internal signed `admin_session` cookie.
Logout is:

```text
/admin/logout
```

SSO/OIDC is not implemented yet. A future SSO path must verify the external
identity, map it to an internal agency/user/role, and issue the same internal
`admin_session` cookie.

## 2. Start At The Dashboard

After login, open:

```text
/admin/operations
```

The dashboard shows the top three issues first. If fewer than three issues are
visible, it fills the first view with compact healthy/current category rows.
Setup reminders persist while setup is incomplete. Operators can dismiss the
reminder only for the current browser session/current next setup blocker; when
the blocker changes, the reminder returns.

The primary navigation groups are:

- Dashboard
- Setup
- Data
- Realtime
- Connectors
- Operations
- Admin

## 3. Use The Setup Wizard

Open:

```text
/admin/operations/setup-wizard
```

The wizard is skippable and private. It does not mutate setup state by being
viewed. It tracks required, recommended, and optional completion buckets:

- required: admin sign-in, agency profile, public feed URLs, license/contact,
  schedule import, vehicle/device source, connector configuration, and
  validation;
- recommended: maintenance and backup owner;
- optional: sharing readiness.

Use focused config pages for settings:

```text
/admin/operations/config
/admin/operations/config/agency
/admin/operations/config/feeds
/admin/operations/config/auth
/admin/operations/config/deployment
/admin/operations/config/advanced
```

## 4. Configure Connectors Truthfully

Open:

```text
/admin/operations/connectors
```

Committed example manifests are examples only. They do not count as configured
connectors. Deployment-owned connector instances have explicit states:

```text
example_available
not_configured
configured_not_tested
dry_run_passed
ready_for_activation
active
blocked
```

Configure connector families in this order:

```text
/admin/operations/connectors/vehicle-avl
/admin/operations/connectors/prediction
/admin/operations/connectors/validators
/admin/operations/connectors/monitoring
/admin/operations/connectors/discovery
```

Vehicle/GPS/AVL setup comes first because the main integration boundary is:

```text
vendor payload -> adapter/sidecar/transform -> authenticated POST /v1/telemetry
```

The browser stores metadata and reference labels only. It does not store secret
values, raw payloads, private endpoint URLs, raw validator commands, or
notification destinations. Dry-run and activation remain separate gates.

## 5. Keep The Boundary Honest

These pages help a deployment owner operate a private self-hosted instance.
They do not contact agencies, vendors, consumers, portals, live validators,
live IdPs, or live connector services, and they do not create retained
evidence.
