# Phase 52 -- Final Public Root Evidence Workflow

## Status

Closed blocker-only for the approved scope. Dedicated final-root workflow
docs, templates, collector, audit tooling, Make targets, and local-only script
tests were added. No real final public root or redacted approval artifact was
available in repo evidence, so no retained final-root evidence packet was
created and `docs/evidence/captured` remains unchanged for Phase 52.

## Goal

Add a guarded workflow for agency-owned or agency-approved final public feed
root evidence. The workflow must support either a real retained final-root
packet or a truthful blocker-only closure when no real root and approval
artifact are available.

Final-root evidence does not by itself prove compliance, consumer acceptance,
consumer ingestion, agency adoption, hosted SaaS availability, production
readiness, SLA/uptime, vendor compatibility, or production-grade ETA quality.

## Scope

- Add final-root approval, DNS, TLS, redirect, public-fetch, validator,
  proxy/config, checksum, and redaction workflow docs/templates.
- Add dedicated final-root collector and audit tooling with stricter
  preconditions than the hosted pilot evidence helpers.
- Default all collection output to ignored `.cache` storage.
- Allow `docs/evidence/captured` retention only with explicit opt-in, a real
  final root, and a real redacted approval artifact.
- Support blocker-only output when no real root or approval artifact exists.
- Preserve the consumer tracker exactly unchanged.

## Non-Goals

- No consumer contact.
- No consumer submission.
- No prepared packet refresh in blocker-only closure.
- No consumer tracker status change.
- No runtime public route.
- No public unauthenticated API change.
- No compliance, agency adoption, hosted SaaS, production-readiness,
  SLA/uptime, vendor-compatibility, consumer-acceptance, or
  production-grade ETA claim.
- No private DNS provider payloads, cookies, Authorization headers, DB URLs,
  TLS private keys, ACME material, raw logs, unredacted correspondence, or
  private diagnostics.

## Files Likely To Change

- `scripts/collect-final-root-evidence.sh`
- `scripts/audit-final-root-evidence.sh`
- `Makefile`
- `docs/phase-52-final-public-root-evidence-workflow.md`
- `docs/handoffs/phase-52.md`
- `docs/handoffs/latest.md`
- `docs/current-status.md`
- `docs/backlog.md`
- `docs/open-questions.md`
- `docs/agency-owned-domain-readiness.md`
- `docs/final-root-operator-request.md`
- `docs/evidence/README.md`
- `docs/evidence/templates/final-root-approval-template.md`
- `docs/evidence/templates/final-root-public-fetch-template.md`
- `docs/evidence/templates/final-root-validator-template.md`
- `docs/evidence/templates/final-root-packet-readme-template.md`

`Taskfile.yml` exists, but Phase 52 should not edit it unless the execution
agent finds an established requirement to mirror new Make targets there.

## Safety Boundaries

`scripts/collect-final-root-evidence.sh` must default to:

```text
.cache/final-root-evidence/<UTC timestamp>/
```

Allowed flags:

- `--help`
- `--blocker-only`
- `--dry-run`
- `--retain-captured`

Allowed environment variables:

- `FINAL_ROOT_BASE_URL`
- `FINAL_ROOT_ENVIRONMENT_NAME`
- `FINAL_ROOT_APPROVAL_ARTIFACT`
- `FINAL_ROOT_APPROVAL_SUMMARY`
- `CAPTURE_DATE_UTC`
- `OUTPUT_DIR`
- `FORCE`
- `ALLOW_CAPTURED_EVIDENCE_WRITE`
- `ADMIN_BASE_URL`
- `ADMIN_TOKEN`
- `CONNECT_TIMEOUT_SECONDS`
- `REQUEST_TIMEOUT_SECONDS`
- `MAX_FEED_BYTES`

Writes under `docs/evidence/captured` require all of:

- `--retain-captured`;
- `ALLOW_CAPTURED_EVIDENCE_WRITE=true`;
- valid `FINAL_ROOT_BASE_URL`;
- real readable redacted `FINAL_ROOT_APPROVAL_ARTIFACT`.

The collector must reject credential-bearing URLs, query strings, fragments,
non-HTTPS non-loopback roots, symlink outputs, traversal, raw DNS provider
exports, cookies, Authorization headers, DB URLs, private TLS keys, ACME
material, raw logs, private payloads, and unredacted diagnostics.

## Evidence And Claim Boundaries

Templates under `docs/evidence/templates` are not evidence.

`docs/evidence/captured` must remain unchanged in blocker-only closure. A real
retained packet may support only this narrow claim:

```text
Final-root approval and technical fetch/validator evidence exists for the
exact recorded root and date.
```

It does not prove compliance, consumer acceptance, consumer ingestion, agency
adoption, hosted SaaS, production readiness, SLA/uptime, vendor compatibility,
or ETA quality.

## Implementation Details

Use dedicated final-root scripts rather than reusing
`collect-hosted-evidence.sh` / `audit-hosted-evidence.sh`. Existing hosted
evidence tooling can create pending hosted-pilot scaffolds; final-root evidence
needs stricter approval/root gates so a placeholder packet cannot look like
agency-approved proof.

Blocker-only collector output must write exactly:

- `blocker.json`
- `blocker.md`
- `manifest.json`
- `manifest.md`

Real-evidence output must write exactly:

- `README.md`
- `approval.md`
- `dns-tls-redirect.md`
- `public-fetches.md`
- `validator-record.md`
- `proxy-config-summary.md`
- `redaction-notes.md`
- `manifest.json`
- `manifest.md`
- `SHA256SUMS.txt`
- `artifacts/public/*`
- `artifacts/tls/*`
- `artifacts/dns/*`
- `artifacts/validation/*`
- `artifacts/operator-supplied/*`

`scripts/audit-final-root-evidence.sh` contract:

- Make target: `make audit-final-root-evidence`
- Env vars: `FINAL_ROOT_PACKET_DIR`, optional `AUDIT_MODE=real|blocker`
- Flags: `--help`, `--blocker-only`

Real audit must fail if approval evidence is missing, root mismatch exists,
placeholders remain, feed artifacts/checksums are missing, validator status is
missing/unavailable/failed, unsafe strings are detected, `SHA256SUMS.txt` does
not match, required redaction notes are missing, or the consumer tracker
changed.

Blocker audit may pass only when blocker files truthfully record no real root
or no approval artifact, no captured evidence directory was created, and claim
flags remain false.

Checkpoint strategy:

- `Phase 52 -- Checkpoint 000001: add final-root evidence workflow plan`
- `Phase 52 -- Checkpoint 000002: add final-root workflow docs and templates`
- `Phase 52 -- Checkpoint 000003: add final-root collector and audit tooling`
- `Phase 52 -- Checkpoint 000004: close blocker-only or retain audited final-root packet`
- `Phase 52 -- Checkpoint 000005: update status and handoff`

## Tests

Tests must use local fixtures or local HTTP servers only. Unit tests must not
fetch external final-root URLs.

Add tests for:

- help output;
- blocker-only exact files;
- dry-run exact files;
- real-mode exact files with a local HTTP server;
- no captured write by default;
- captured write only with explicit opt-in plus approval artifact;
- URL validation;
- root mismatch;
- symlink/path rejection;
- placeholder detection;
- unsafe string detection;
- checksum validation;
- validator missing/failure audit failure;
- blocker audit pass;
- consumer tracker preservation.

## Performance And Scale Tests

- Stream checksums for schedule ZIP and protobuf artifacts.
- Enforce `MAX_FEED_BYTES`.
- Bound validator outputs and manifest sizes.
- Use curl connect/request timeouts.
- Add a local large-artifact fixture test proving checksums and manifest
  generation do not repeatedly load full artifacts.

## Docs, Status, And Handoff Updates

Update status docs with one of two outcomes:

- blocker-only: no final root/approval, no captured evidence, real audit not
  run, consumer tracker unchanged;
- retained evidence: exact packet path, root, approval artifact summary, audit
  result, checksums, and limitations.

Update `docs/evidence/archive-inventory.md` only if a real committed archive is
retained.

## Required Verification Commands

```bash
sh -n scripts/collect-final-root-evidence.sh scripts/audit-final-root-evidence.sh
make collect-final-root-evidence
FINAL_ROOT_PACKET_DIR=.cache/final-root-evidence/<timestamp> AUDIT_MODE=blocker make audit-final-root-evidence
make validate
make test
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
python3 - <<'PY'
import json
from pathlib import Path

expected = [
    "Google Maps",
    "Apple Maps",
    "Transit App",
    "Bing Maps",
    "Moovit",
    "Mobility Database",
    "transit.land",
]

data = json.loads(Path("docs/evidence/consumer-submissions/status.json").read_text())
records = data.get("targets", [])
seen = {row["target"]: row.get("status") for row in records}
assert list(seen) == expected, seen
assert all(seen[name] == "prepared" for name in expected), seen
PY
git diff --exit-code -- docs/evidence/consumer-submissions/status.json
docker compose -f deploy/docker-compose.yml config
```

For blocker-only closure:

```bash
git diff --exit-code -- docs/evidence/captured
```

For real retained evidence only:

```bash
FINAL_ROOT_PACKET_DIR=docs/evidence/captured/<environment>/<UTC-date> AUDIT_MODE=real make audit-final-root-evidence
```

## Exact Consumer Tracker Preservation Checks

```bash
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
python3 - <<'PY'
import json
from pathlib import Path

expected = [
    "Google Maps",
    "Apple Maps",
    "Transit App",
    "Bing Maps",
    "Moovit",
    "Mobility Database",
    "transit.land",
]

data = json.loads(Path("docs/evidence/consumer-submissions/status.json").read_text())
records = data.get("targets", [])
seen = {row["target"]: row.get("status") for row in records}
assert list(seen) == expected, seen
assert all(seen[name] == "prepared" for name in expected), seen
PY
git diff --exit-code -- docs/evidence/consumer-submissions/status.json
```

## Plan Risks And Blockers

- No retained agency-owned or agency-approved final public root is currently
  available in repo evidence.
- Approval artifacts may contain private data and require redaction before
  retention.
- Validator tooling may be unavailable in some environments.
- Existing hosted evidence patterns can allow pending placeholders, so Phase
  52 must keep stricter dedicated final-root scripts and audits.
