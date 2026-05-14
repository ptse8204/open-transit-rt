# Phase 85 -- Operations And Maintenance Center V2 Plan

## Scope

Phase 85 improves the private Operations Console maintenance surface for
small-host operations. The goal is to make routine maintenance understandable
from the browser while keeping destructive or environment-specific work in an
operator shell.

The phase may add or improve:

- backup readiness and restore-drill guidance;
- upgrade and rollback review checklists;
- deployment doctor summary browsing from safe local `.cache` summaries;
- operations reliability summary browsing from safe local `.cache` summaries;
- notification draft status from safe local `.cache` summaries;
- support-bundle generation guidance and redaction warnings;
- maintenance cadence calendar/checklist rows;
- safe disk, DB, validator, and tooling status rows when sourced from existing
  non-secret signals;
- operator-facing docs for private maintenance review.

The phase does not run destructive backup, restore, rollback, migration, or
upgrade actions from the browser. It does not execute shell commands from the
browser, upload private outputs, collect retained evidence, contact external
services, send notifications, move consumer statuses, tag, package, publish, or
make production-readiness, hosted-service, SLA/uptime, compliance, consumer,
agency, vendor, hardware, public-launch, release-readiness, or ETA-quality
claims.

## Implementation Boundary

Use existing private Operations Console and maintenance tooling:

- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_maintenance.go`
- `cmd/agency-config/operations_navigation.go`
- `scripts/support-bundle.sh`
- `scripts/deployment-doctor.sh`
- `scripts/operations-notify.sh`
- `scripts/operations-reliability.sh`
- `docs/tutorials/small-agency-maintenance-guide.md`
- `docs/deployment/reference-deployment-doctor.md`
- `docs/upgrade-and-rollback.md`
- `docs/evidence/redaction-policy.md`

Keep the existing route names:

- private `GET /admin/operations/maintenance`;
- private `GET /admin/operations/maintenance.json`.

No migration is planned. Stop and re-plan before adding persisted maintenance
state, a mutable maintenance command API, browser-triggered shell execution, or
destructive operations.

## Security And Data Boundary

The Maintenance Center remains private, authenticated, agency-scoped,
GET-only, and no-store. It may display bounded summaries and fixed
operator-shell guidance, but it must not expose:

- raw logs, stdout, stderr, argv, or command output;
- database URLs, restore URLs, bearer tokens, cookies, JWTs, CSRF values,
  webhook URLs, private keys, ACME material, or secret env values;
- raw telemetry payloads, vendor payloads, GTFS-RT payloads, validator raw
  reports, raw support-bundle files, raw backup dumps, or private file paths;
- evidence-like retained outputs under `docs/evidence/**`;
- arbitrary script names, command arguments, file paths, plugin payloads, or
  browser-supplied shell input.

Safe `.cache` readers, if added, must validate exact file names, schema or
expected field shapes, output root shape, symlink status, size caps, redaction
claim flags, and supported summary fields before rendering. They must not
create directories, run helpers, follow symlinks, or read raw artifacts.

## Master Approval

The Master Agent approves implementation under these constraints from
simulated sub-agent review because real sub-agent spawning was blocked by the
thread limit at Phase 85 start:

- Build on the existing Maintenance Center route rather than creating a new
  public or parallel app surface.
- Keep all maintenance actions as read-only review or fixed operator-shell
  guidance.
- Prefer status cards and tables for backup readiness, restore drills,
  upgrade/rollback, deployment doctor, reliability, notification drafts,
  support bundle, redaction, and cadence rows.
- Use existing `.cache` summary contracts only when they can be read safely and
  redacted.
- Keep no-JavaScript fallback because all core content should be
  server-rendered.
- Keep every claim flag false.

## Sub-Agent Review Plan

Real or simulated reviews use the intended model levels from the authorized
Phase 75-90 track:

- Context / Repo Truth Sub-Agent, GPT-5.5 x-high: simulated locally because
  real-agent spawning was blocked by the thread limit. Review confirms Phase 84
  is closed, Phase 85 is active, the existing Maintenance route is private and
  GET-only, maintenance scripts already write ignored `.cache` summaries, and
  protected evidence/consumer paths remain gated.
- Planning Sub-Agent, GPT-5.5 x-high: simulated locally. Review approves the
  checkpoint sequence below and the no-migration/no-browser-execution default.
- Implementation Sub-Agent, GPT-5.5 high: simulated in the main rollout unless
  agent capacity becomes available.
- QA Sub-Agent, GPT-5.5 high: simulated locally. Review requires private route
  tests, JSON shape tests, no-form/no-POST checks, safe `.cache` reader tests,
  forbidden private string checks, claim-flag tests, and baseline validation.
- UI/UX Sub-Agent, GPT-5.5 high: simulated locally. Review requires
  nontechnical operator labels, bounded cards/tables, clear missing/blocked
  states, no-JS completeness, and stable route/navigation behavior.
- Documentation / IA Sub-Agent, GPT-5.5 high: simulated locally. Review
  requires tutorial/status alignment and clear redaction/authorization
  boundaries.
- Claim-Boundary Sub-Agent, GPT-5.5 high: simulated locally. Review blocks
  production readiness, hosted SaaS, SLA/uptime, compliance, agency adoption,
  consumer acceptance, vendor compatibility, hardware certification,
  release-readiness, public-launch, and evidence claims.
- Security/Auth Sub-Agent, GPT-5.5 high: simulated locally. Review requires no
  browser command execution, no destructive action, no raw/private output
  exposure, no external sends, and no evidence writes.
- Data/Migration Sub-Agent: not planned because Phase 85 should add no
  migration or persisted model. Stop and re-plan if persistence becomes
  necessary.

All required edits from these reviews are incorporated into this plan.

## Checkpoints

```text
Phase 85 -- Checkpoint 000001: add operations maintenance center v2 plan
Phase 85 -- Checkpoint 000002: add maintenance summary readers
Phase 85 -- Checkpoint 000003: add backup restore and upgrade review panels
Phase 85 -- Checkpoint 000004: add support bundle redaction and cadence guidance
Phase 85 -- Checkpoint 000005: add maintenance infrastructure check summaries
Phase 85 -- Checkpoint 000006: close operations maintenance center v2 review
```

## Acceptance Criteria

- Existing Maintenance routes remain private, authenticated, agency-scoped,
  no-store, GET-only, and unavailable under `/public`.
- The browser executes no commands and performs no backup, restore, rollback,
  migration, package, notification, or external send action.
- JSON is bounded, read-only, schema-stable, and all claim flags remain false.
- Backup and restore panels distinguish configuration presence from actual
  successful backup/restore proof.
- Upgrade/rollback guidance stays checklist-only and does not tag, package, or
  publish.
- Deployment doctor, operations reliability, notification, and support-bundle
  summaries are shown only from safe local `.cache` summary files when present.
- Redaction warnings and private-output boundaries are visible before support
  bundle or diagnostic sharing guidance.
- Missing, blocked, stale, and not-configured states remain visible.
- Consumer tracker rows remain prepared-only.
- Protected evidence and consumer packet paths remain untouched.

## Validation

Baseline Phase 85 validation:

```bash
git status --short
git diff --check
make check
make audit-product-acceptance
make audit-final-claim-review
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
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum
```

Additional Phase 85 checks:

```bash
go test ./cmd/agency-config -run 'Maintenance|OperationsNavigation|RouteTitles'
sh -n scripts/support-bundle.sh scripts/deployment-doctor.sh scripts/operations-notify.sh scripts/operations-reliability.sh
make validate
make test
```

When route/UI work lands and local startup is safe, also run:

```bash
RUN_LOCAL_APP=true make release-candidate-check
```
