# Self-Hosted Operator Trial

This tutorial guides one local/reference evaluation path across deployment
prep, GTFS onboarding, public feed checks, readiness review, validator
handling, and the synthetic AVL adapter dry-run.

It is an operator trial checklist, not evidence. Operator output stays
local/private by default. Logs, screenshots, validator output, copied
summaries, and support notes must not be committed unless a later evidence
phase reviews, redacts, and retains them.

This trial does not prove agency approval, agency adoption, consumer
acceptance, final-root proof, CAL-ITP/Caltrans compliance, hosted SaaS
availability, production readiness, vendor compatibility, or production-grade
ETA quality.

## Related Context

- [Phase 39 Handoff](../handoffs/phase-39.md)
- [Phase 40 Handoff](../handoffs/phase-40.md)
- [OCI/OCL Reference Deployment](../deployment/oci-reference-deployment.md)
- [Reusable Agency Onboarding](reusable-agency-onboarding.md)
- [Integration Adapter Kit](../integration-adapter-kit.md)
- [Operator Smoke And Support Bundle](operator-smoke-and-support-bundle.md)
- [CAL-ITP Readiness Checklist](calitp-readiness-checklist.md)
- [Evidence Redaction Policy](../evidence/redaction-policy.md)

## 1. Prepare The Deployment Path

For a local Compose evaluation, run the onboarding command in its default
local-compose mode. It starts the required local services, applies migrations,
imports the requested GTFS, verifies public feeds, and prints admin next steps.

For an already-running reference deployment, follow the
[OCI/OCL Reference Deployment](../deployment/oci-reference-deployment.md)
first, then run onboarding in `--mode running` with the deployment's private
admin URL, admin token, and database URL.

Do not expose admin, debug, GTFS Studio, or Operations Console routes on the
anonymous public feed edge.

## 2. Choose A GTFS Source

For a real local/reference evaluation, provide the agency ID from
`agency.txt` and a public GTFS ZIP URL:

```bash
make agency-pilot-up \
  AGENCY_ID=example-agency \
  GTFS_URL=https://example.org/path/to/gtfs.zip
```

Publication metadata is local/reference placeholder metadata unless the
operator supplies agency-approved contact and license values. Placeholder
metadata is not agency-approved metadata, final-root evidence, consumer
evidence, compliance proof, or production-readiness proof.

### No-External-Network Fixture Option

When the goal is only to prove the command path works without depending on an
external GTFS download, serve the committed demo fixture from a temporary local
HTTP server:

```bash
tmpdir="$(mktemp -d)"
(cd testdata/gtfs/valid-small && zip -qr "$tmpdir/valid-small.zip" .)
python3 -m http.server 18080 --directory "$tmpdir"
```

In another terminal, run:

```bash
make agency-pilot-up \
  AGENCY_ID=demo-agency \
  GTFS_URL=http://127.0.0.1:18080/valid-small.zip \
  SKIP_VALIDATORS=true
```

This `demo-agency` fixture path is local demo evaluation only. It is not real
agency proof, agency approval, official agency feed status, final-root proof,
or consumer evidence.

## 3. Verify The Five Public Paths

For a repeatable command that fetches these paths, records sizes/checksums,
checks the admin boundary, records validator tooling state, and runs the
synthetic AVL dry-run fixture, use
[Operator Smoke And Support Bundle](operator-smoke-and-support-bundle.md).

The onboarding output prints the public base URL. Verify that these paths
return non-empty responses:

```text
/public/feeds.json
/public/gtfs/schedule.zip
/public/gtfsrt/vehicle_positions.pb
/public/gtfsrt/trip_updates.pb
/public/gtfsrt/alerts.pb
```

For manual checks:

```bash
BASE=http://localhost:8080
curl -fsS "$BASE/public/feeds.json" -o /tmp/open-transit-feeds.json
curl -fsS "$BASE/public/gtfs/schedule.zip" -o /tmp/open-transit-schedule.zip
curl -fsS "$BASE/public/gtfsrt/vehicle_positions.pb" -o /tmp/open-transit-vehicle-positions.pb
curl -fsS "$BASE/public/gtfsrt/trip_updates.pb" -o /tmp/open-transit-trip-updates.pb
curl -fsS "$BASE/public/gtfsrt/alerts.pb" -o /tmp/open-transit-alerts.pb
```

These local files are trial output. Do not commit them as evidence.

## 4. Review Readiness

Use the admin URL and token printed by onboarding, or use the private admin
boundary from the reference deployment:

```text
/admin/operations/readiness
```

The readiness page must stay behind admin authentication, SSH tunnel, VPN, or
another private/admin-protected boundary. Do not expose
`/admin/operations/readiness` on the public edge.

Review the readiness rows for public URLs, static GTFS, Vehicle Positions,
Trip Updates, Alerts, license/contact metadata, validation, telemetry
freshness, operations status, and consumer packet preparedness. The page is
read-only and does not run validators, contact consumers, create evidence, or
claim CAL-ITP/Caltrans compliance.

## 5. Run, Skip, Or Blocker-Document Validators

Install and check pinned validator tooling when validation should run:

```bash
make validators-install
make validators-check
```

By default, `make agency-pilot-up` treats validators as best-effort. Missing
or misconfigured pinned tooling records a blocker in the command output but
does not fail the whole onboarding flow.

Use strict validators when validator blockers should fail the trial command:

```bash
make agency-pilot-up \
  AGENCY_ID=example-agency \
  GTFS_URL=https://example.org/path/to/gtfs.zip \
  STRICT_VALIDATORS=true
```

Use skipped validators when the goal is only import and public-path
verification:

```bash
make agency-pilot-up \
  AGENCY_ID=example-agency \
  GTFS_URL=https://example.org/path/to/gtfs.zip \
  SKIP_VALIDATORS=true
```

Validator output is a local quality signal for this trial. It is not consumer
acceptance, CAL-ITP/Caltrans compliance, agency approval, final-root proof, or
production-readiness proof.

## 6. Run The Synthetic AVL Dry-Run Adapter

Run the synthetic adapter fixture:

```bash
go run ./cmd/avl-vendor-adapter --dry-run \
  --reference-time 2026-05-04T12:00:00Z \
  --mapping testdata/avl-vendor/mapping.json \
  testdata/avl-vendor/minimal-gps.json
```

The command prints transformed Open Transit RT telemetry JSON to stdout and
diagnostics JSON to stderr. It does not send telemetry, prove ingest status,
prove real vendor compatibility, or prove production AVL reliability.

Use [Integration Adapter Kit](../integration-adapter-kit.md) before mapping
any private device, AVL, or vendor payload outside this public repo.

## 7. Review Next Actions Without Creating Evidence

Use trial output only to decide what to fix next:

- replace placeholder license/contact metadata with agency-approved values;
- fix GTFS import or validator blockers;
- configure private device credentials before sending real telemetry;
- review readiness rows and row-level next actions;
- keep support notes private until a future evidence phase approves redaction
  and retention.

Do not create external evidence packets, final-root evidence, consumer
submission artifacts, target-originated artifact directories, or stronger
public claims from this trial.

## 8. Teardown And Cleanup

For the local Compose flow, stop the local app package:

```bash
make agency-app-down
```

If you used the temporary `python3 -m http.server`, stop it with `Ctrl+C` in
the terminal where it is running. Then remove the temporary directory:

```bash
rm -rf "$tmpdir"
```

For `--mode running` or a reference deployment, use the deployment supervisor
or the [OCI/OCL Reference Deployment](../deployment/oci-reference-deployment.md)
guide instead of `make agency-app-down`.

Use `--reset-local-state` only when you intentionally want to reset local
Compose state for the onboarding helper:

```bash
scripts/agency-pilot-onboard.sh \
  --agency-id demo-agency \
  --gtfs-url http://127.0.0.1:18080/valid-small.zip \
  --reset-local-state
```

The reset asks for typed confirmation. Use `--force` only for intentional
scripted reset flows.
