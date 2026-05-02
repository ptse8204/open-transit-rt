# Device Token Lifecycle

This guide explains how Open Transit RT device tokens work for telemetry onboarding.

Device tokens are bearer credentials. Anyone with a valid token for a bound agency, device, and vehicle can send telemetry for that binding, so tokens must be stored like passwords.

All examples use synthetic local demo values. Do not commit real private device IDs, vehicle IDs, vendor IDs, tokens, raw private telemetry, or private logs in public docs.

## What A Device Token Is

A device token is an opaque credential sent with telemetry requests:

```text
Authorization: Bearer <device-token>
```

The token is checked against the active binding for:

- `agency_id`
- `device_id`
- `vehicle_id`

If the payload identifiers do not match the token binding, telemetry ingest returns `401 Unauthorized` and does not store the event.

## Local Demo Seed

The local demo seed includes this synthetic binding:

```text
agency_id=demo-agency
device_id=device-1
vehicle_id=bus-1
device token=dev-device-token
```

This token is for local evaluation only. It is documented so the local demo can run from committed files. Production deployments must generate their own tokens and keep them private.

## Create, Rotate, And Rebind

The current supported operation is rotate/rebind. It creates a credential when none exists for the device and rotates the token when a credential already exists.

Browser flow:

```text
/admin/operations/devices
```

JSON admin API:

```text
POST /admin/devices/rebind
Authorization: Bearer <admin-token>
Content-Type: application/json
```

Synthetic example body:

```json
{
  "device_id": "device-1",
  "vehicle_id": "bus-1",
  "reason": "local device onboarding"
}
```

Local helper:

```bash
scripts/device-onboarding.sh rebind --device-id device-1 --vehicle-id bus-1
```

The rebind API intentionally returns the new token once. Store it immediately. The Operations Console and helper script may display it only at the time of rotation.

After rotation or rebind, the previous token for that device binding stops working immediately.

## Secure Storage Expectations

Store tokens in deployment-owned secret storage or private environment files with restricted permissions. For devices or vendor middleware, use the vendor or device secret mechanism rather than hardcoding tokens in source code.

Do not store tokens in:

- public docs;
- issue comments;
- screenshots;
- shell history;
- example config committed to git;
- raw logs;
- telemetry payload archives;
- `.cache` files.

If a token may have been exposed, rotate it and update the device or adapter configuration. Treat old logs, screenshots, and copied commands as suspect until reviewed.

## Binding Rules

One active token is bound to one device ID and one vehicle ID for an agency.

Use rebind when:

- a physical device is moved to a different vehicle;
- a vehicle identifier changes;
- a vendor mapping is corrected;
- a token is suspected compromised;
- a device is replaced but the logical `device_id` stays the same.

Use a new `device_id` when the agency wants to track a different logical device identity. Keep identifiers stable enough for operations review, but do not expose private serial numbers or vendor account IDs in public docs unless they are approved public-safe identifiers.

## Operator Responsibilities

Operators should:

- record why a token was rotated or rebound;
- confirm the device sends `agency_id`, `device_id`, and `vehicle_id` matching the binding;
- confirm old tokens fail after rotation;
- confirm fresh telemetry appears in `/admin/operations/telemetry`;
- keep token material out of public evidence;
- review device and vendor identifiers before committing screenshots or artifacts.

Device rebinding is audit logged. Audit records should identify the operational action without storing the raw token.

## When A Vehicle Or Device Changes

If a device moves from one vehicle to another, rotate/rebind the device to the new `vehicle_id` before sending telemetry from the new vehicle.

If a vehicle is renamed in GTFS or agency operations, update the binding to the approved public-safe vehicle identifier expected by Open Transit RT.

If a vendor changes its internal unit ID, update the adapter mapping privately. Do not spread vendor-specific IDs through core matching or public evidence unless they are approved public-safe identifiers.

## What Never To Commit

Do not commit:

- device tokens;
- admin tokens;
- JWT or CSRF secrets;
- `DEVICE_TOKEN_PEPPER`;
- DB URLs with passwords;
- vendor credentials;
- private AVL payloads;
- raw private telemetry;
- private device IDs, vehicle IDs, or vendor IDs;
- private operator notes;
- private logs with credentials;
- `.cache` files.

Use synthetic values such as `demo-agency`, `device-1`, `bus-1`, `trip-10-0800`, and `dev-device-token` in public examples.

