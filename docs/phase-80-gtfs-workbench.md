# Phase 80 -- GTFS Workbench Plan

## Scope

Phase 80 turns GTFS import, review, validation triage, draft authoring links,
and feed-output review into a private browser-first GTFS Workbench for
operators. The workbench coordinates existing Operations Console surfaces
rather than replacing the importer, GTFS Studio, validator runners, or
published feed model.

The phase may improve:

- a private `GTFS Workbench` route and optional private JSON view;
- active schedule, latest import, and draft schedule summaries;
- schedule ZIP identity, source checksum, import history, and row counts;
- required-file checklist for `agency.txt`, `routes.txt`, `stops.txt`,
  `trips.txt`, `stop_times.txt`, and service calendars;
- bounded previews for agency, routes, stops, trips, calendar/service, and
  frequency/service signals where existing private data supports them;
- validation triage grouped by source, severity, likely owner, affected files,
  safe fix path, and verification step;
- draft/publish review checklist that keeps draft GTFS separate from active
  published feed versions;
- schedule history and rollback guidance without browser rollback execution;
- docs and tests for private route, JSON shape, escaping, bounded preview,
  claim-boundary, and no-raw-output behavior.

The phase does not add public admin routes, auto-fix GTFS, publish schedules
without the existing confirmation path, contact external validators or map
services from the browser, collect or retain evidence, move consumer statuses,
claim validator-clean status, claim CAL-ITP/Caltrans compliance, claim agency
approval, claim consumer acceptance, claim final-root readiness, publish a
release artifact, or create a hosted SaaS, production, vendor, SLA, hardware,
or ETA-quality claim.

## Implementation Boundary

Use the existing private Operations Console and GTFS surfaces:

- `GET /admin/operations/gtfs-import`
- admin-only `POST /admin/operations/gtfs-import`
- `GET /admin/operations/gtfs-quality`
- admin-only `POST /admin/operations/gtfs-quality`
- `GET /admin/operations/validation-health`
- `GET /admin/operations/feed-health`
- `GET /admin/operations/feeds`
- `/admin/gtfs-studio` as the existing draft authoring surface.

Phase 80 may add:

- `GET /admin/operations/gtfs-workbench`
- `GET /admin/operations/gtfs-workbench.json`

No migration is planned. The first implementation pass should use existing
`feed_version`, `gtfs_import`, `gtfs_*`, `gtfs_draft`, `gtfs_draft_publish`,
and `validation_report` data where the configured runtime exposes it. If the
current repository interfaces cannot support a durable field, record the gap
instead of adding a risky schema change.

Any browser workbench command must be private, role-checked, CSRF-safe for
cookie auth, body-capped, agency-scoped from the authenticated principal, and
bounded to a server-owned action. Phase 80 should avoid new POST routes unless
they reuse an already established safe action. The browser must not accept raw
commands, validator binaries, argument lists, artifact paths, local paths,
output paths, raw URLs beyond the existing GTFS import URL field, evidence
destinations, secrets, or external submission targets.

## Master Approval

The Master Agent approves implementation with these required constraints from
sub-agent review:

- Add **GTFS Workbench** as a private Schedule entry while preserving existing
  route paths.
- Keep Go server-rendered HTML as the default; any buildless JavaScript may
  only enhance already-rendered private content.
- Make the operator sequence explicit:
  set metadata, import or edit GTFS, review quality, run validation, review
  feed health, then decide the next operator action.
- Keep internal importer checks distinct from MobilityData static validator
  diagnostics.
- Keep draft data and active published feed versions separate.
- Replace risky wording such as `Accepted source` with `Allowed source`, and
  avoid broad `ready`, `valid`, `accepted`, `approved`, `compliant`, or
  production wording.
- Use status language such as `ok`, `needs_review`, `missing`, `blocked`, and
  `unknown`; `ready for local review` may appear only as a bounded operator
  review status, never as production, compliance, final-root, consumer, or
  release readiness.

## Sub-Agent Review Plan

Real or simulated reviews use the intended model levels from the authorized
Phase 75-90 track:

- Context / Repo Truth Sub-Agent, GPT-5.5 x-high: simulated locally because
  the agent thread limit was reached; confirm GTFS import, GTFS quality, GTFS
  Studio, validation, feed discovery, tests, protected paths, and consumer
  tracker truth.
- Planning Sub-Agent, GPT-5.5 x-high: confirm checkpoint sequence and no
  migration/evidence boundary.
- Implementation Sub-Agent, GPT-5.5 high: simulated if thread capacity is
  exhausted; identify the smallest private route/model/template slices.
- QA Sub-Agent, GPT-5.5 high: confirm focused Workbench route/JSON/rendering,
  preview, no-leak, and no-overclaim tests.
- UI/UX Sub-Agent, GPT-5.5 high: confirm nontechnical workbench IA, labels,
  no-JS fallback, and stable route linking.
- Documentation / IA Sub-Agent, GPT-5.5 high: confirm README/docs/wiki/status
  updates.
- Claim-Boundary Sub-Agent, GPT-5.5 high: confirm no validator-clean,
  compliance, adoption, consumer, final-root, hosted SaaS, production, vendor,
  SLA, hardware, or ETA-quality claim.
- Security/Auth Sub-Agent, GPT-5.5 high: confirm private route, role, CSRF,
  agency-scope, upload/URL, command, path, output, and JSON boundaries.
- Data/Migration Sub-Agent: not planned unless implementation requires a new
  persisted model. Stop and re-plan if a migration becomes necessary.

## Checkpoints

```text
Phase 80 -- Checkpoint 000001: add GTFS Workbench plan
Phase 80 -- Checkpoint 000002: add import diff and schedule summaries
Phase 80 -- Checkpoint 000003: add GTFS preview tables and filters
Phase 80 -- Checkpoint 000004: improve safe draft publish review
Phase 80 -- Checkpoint 000005: add rollback and schedule history UX
Phase 80 -- Checkpoint 000006: close GTFS Workbench review
```

## Acceptance Criteria

- Workbench routes are private, authenticated, agency-scoped, no-store, and
  unavailable under `/public/operations`.
- Workbench JSON is read-only, bounded, schema-stable, and claim flags remain
  false.
- Non-admin users can review but do not see new mutation forms.
- Required-file and preview sections are row-capped, text-capped, escaped, and
  deterministic.
- The UI explains whether each section reads or writes; read-only previews do
  not silently import, publish, run validators, create evidence, or contact
  external systems.
- Internal importer checks and MobilityData validator diagnostics are visibly
  separate.
- Publish review remains guidance unless the user uses the existing admin
  confirmation path in GTFS Studio.
- Rollback remains guidance unless a separately reviewed safe browser rollback
  action already exists.
- No raw validator report, stdout, stderr, argv, temp path, private path, raw
  ZIP bytes, token, DB URL, cookie, or authorization header appears in HTML or
  JSON.
- Protected evidence and consumer packet paths remain untouched.
- All seven consumer targets remain exactly `prepared`.

## Validation

Baseline Phase 80 validation:

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

Additional Phase 80 checks:

```bash
go test ./cmd/agency-config -run 'GTFSWorkbench|GTFSImport|GTFSQuality|ValidationHealth|Setup'
go test ./internal/compliance ./internal/gtfs ./cmd/gtfs-import ./cmd/gtfs-studio
make validate
make test
```

Run `RUN_LOCAL_APP=true make release-candidate-check` when the route/UI changes
are in place and local app startup is safe. If an environment limitation blocks
a check, record the exact blocker in the Phase 80 handoff.
