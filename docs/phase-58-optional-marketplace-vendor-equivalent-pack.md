# Phase 58 -- Optional Marketplace / Vendor-Equivalent Pack

## Status

Planning accepted. Execution may add docs/templates, local audit tooling, tests,
Make targets, and status/handoff updates for an optional BYOD/procurement
packet. It must not claim marketplace approval, paid support, vendor
compatibility, hosted service availability, SLA coverage, or production
readiness.

## Goal

Provide a public-safe, reusable operator documentation pack for BYOD/hardware
planning, implementation templates, support boundaries, SLA/KPI template
language, and procurement-oriented responses without converting templates into
claims.

## Scope

- Add an optional vendor-equivalent documentation pack under `docs/`.
- Include BYOD/hardware intake guidance.
- Include implementation planning template.
- Include support boundary template.
- Include SLA/KPI template language that is explicitly uncommitted unless a
  later operator contract exists.
- Include procurement response template language.
- Add a local audit script and tests that fail on unsupported approval, support,
  SLA, vendor compatibility, hosted service, compliance, production-readiness,
  agency adoption, consumer acceptance, or marketplace claims.
- Add Make targets and validation scaffolding.
- Update roadmap/status/backlog/open-question/handoff docs.

## Non-Goals

- No marketplace submission.
- No marketplace approval claim.
- No vendor certification claim.
- No hardware certification claim.
- No paid support claim.
- No SLA/uptime commitment.
- No hosted SaaS or hosted service claim.
- No production-readiness claim.
- No compliance claim.
- No agency adoption or consumer acceptance claim.
- No consumer contact, portal automation, or consumer status change.
- No retained evidence creation.
- No `docs/evidence` writes.

## Files Likely To Change

- `docs/vendor-equivalent-pack/README.md`
- `docs/vendor-equivalent-pack/byod-hardware-intake-template.md`
- `docs/vendor-equivalent-pack/implementation-plan-template.md`
- `docs/vendor-equivalent-pack/support-boundaries-template.md`
- `docs/vendor-equivalent-pack/sla-kpi-template.md`
- `docs/vendor-equivalent-pack/procurement-response-template.md`
- `scripts/audit-vendor-equivalent-pack.sh`
- `scripts/test-vendor-equivalent-pack.sh`
- `Makefile`
- `docs/phase-58-optional-marketplace-vendor-equivalent-pack.md`
- `docs/current-status.md`
- `docs/backlog.md`
- `docs/open-questions.md`
- `docs/roadmap-to-calitp-compliance-and-gap-closure.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-58.md`
- `docs/dependencies.md`
- `docs/decisions.md`

Execution must not change:

- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/`
- `docs/evidence/consumer-submissions/artifacts/`
- `docs/evidence/consumer-submissions/packets/`
- `docs/evidence/captured/`

## Safety Boundaries

Templates must use placeholders only. They must not include real agency names,
real vendor names, real customer names, real consumer names as contacted
parties, private URLs, credentials, prices, support commitments, uptime
commitments, marketplace IDs, ticket IDs, or private correspondence.

Any SLA/KPI wording must be template-only and must say it is not a commitment
unless an operator signs a separate agreement and retains claim-specific
evidence. Procurement language must distinguish existing repository
capabilities from future/operator-owned commitments.

The audit script must scan only repository template files and must not contact
external marketplaces, vendors, consumers, agencies, or procurement systems.

## Evidence And Claim Boundaries

All claim flags remain false:

- marketplace approval;
- paid support;
- SLA/uptime;
- vendor compatibility;
- hardware certification;
- hosted SaaS/hosted service;
- production readiness;
- compliance;
- agency adoption;
- consumer acceptance/ingestion;
- production-grade ETA.

Phase 58 creates templates only. It creates no retained evidence and writes
nothing under `docs/evidence`.

Consumer tracker state remains unchanged; all seven targets remain
`prepared`.

## Implementation Details

Add `docs/vendor-equivalent-pack/` with:

- `README.md`: pack index, use boundaries, required review before reuse.
- `byod-hardware-intake-template.md`: device inventory, telemetry contract,
  device-token handling, data-quality prerequisites, and redaction guidance.
- `implementation-plan-template.md`: deployment, GTFS, AVL, validation,
  operator workflow, rollback, and acceptance-review placeholders.
- `support-boundaries-template.md`: maintainer/community/operator boundaries,
  escalation placeholders, no paid support/SLA unless separately contracted.
- `sla-kpi-template.md`: template metrics and review cadence with explicit
  non-commitment language.
- `procurement-response-template.md`: reusable response sections with
  claim-gated placeholders.

Add `scripts/audit-vendor-equivalent-pack.sh`:

- verify required templates exist;
- fail on unsafe private strings;
- fail on positive claims such as "marketplace approved", "certified
  hardware", "vendor compatible", "SLA-backed", "guaranteed uptime",
  "paid support included", "production ready", "compliant", "consumer
  accepted", or "hosted SaaS available";
- require explicit boundary phrases in the README and SLA/KPI template;
- verify consumer tracker remains seven prepared targets.

Add `scripts/test-vendor-equivalent-pack.sh`:

- exercise audit success;
- mutate a temporary copy with unsupported claim wording and verify audit
  failure;
- mutate a temporary copy with missing template and verify audit failure;
- verify the script help path.

Add Make targets:

- `audit-vendor-equivalent-pack`
- `test-vendor-equivalent-pack`

## Tests

- Shell syntax tests for the new scripts.
- Local audit script tests.
- Existing `make validate` must include script and docs existence checks.
- Consumer tracker preservation checks must pass.

## Performance And Scale Tests

The audit scans a small bounded docs directory and the consumer tracker JSON
only. No benchmark is needed.

## Docs, Status, And Handoff Updates

Close Phase 58 by updating:

- this phase document;
- `docs/handoffs/phase-58.md`;
- `docs/handoffs/latest.md`;
- `docs/current-status.md`;
- `docs/backlog.md`;
- `docs/open-questions.md`;
- `docs/roadmap-to-calitp-compliance-and-gap-closure.md`;
- `docs/dependencies.md`;
- `docs/decisions.md`.

The handoff must explicitly state that Phase 58 is templates/audit only and no
marketplace approval, vendor compatibility, certified hardware, paid support,
SLA/uptime, hosted service, production-readiness, compliance, agency adoption,
consumer acceptance, retained evidence, or production-grade ETA claim was
created.

## Required Verification Commands

Run and report:

```bash
sh -n scripts/audit-vendor-equivalent-pack.sh scripts/test-vendor-equivalent-pack.sh
./scripts/test-vendor-equivalent-pack.sh
make audit-vendor-equivalent-pack
make validate
make test
make smoke
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
git diff --exit-code -- docs/evidence/consumer-submissions/current docs/evidence/consumer-submissions/artifacts docs/evidence/consumer-submissions/packets docs/evidence/captured
find docs/evidence/consumer-submissions/artifacts -mindepth 2 -maxdepth 2 -type f ! -name README.md -print
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured
docker compose -f deploy/docker-compose.yml config
```

Run `INTEGRATION_TESTS=1 make test-integration` if the local database is
available and record any environment blocker truthfully.

The `find` command must print no files for the current Phase 58 state.
