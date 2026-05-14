# Phase 82 -- Feed Health And Validation Center Plan

## Scope

Phase 82 adds a private browser-first Feed Health And Validation Center that
combines existing feed status, validator-health summaries, GTFS quality
triage, readiness signals, reliability context, consumer prepared-tracker
status, and blocker history into one operator review surface.

The phase may improve:

- a private `Feed Health And Validation Center` route and optional private
  JSON view;
- a five public feed URL panel for `feeds.json`, static GTFS, Vehicle
  Positions, Trip Updates, and Alerts;
- a clear distinction between feed health, validation health, GTFS quality,
  reliability, and readiness signals;
- validation history summaries for schedule and realtime feeds;
- GTFS issue drilldowns with likely owner, affected files, safe fix path,
  verification steps, and escalation guidance;
- readiness timeline rows using only `ok`, `needs_review`, `missing`,
  `blocked`, and `unknown`;
- blocker rows for missing metadata, missing artifacts, failed/stale/not-run
  validators, unhealthy reliability, missing contact/license, and prepared-only
  consumer tracker status;
- plain-language empty states for first-run, missing tooling, stale reports,
  and off-host validation needs.

The phase does not change public feed semantics, validator execution
semantics, schedule publish gates, telemetry ingest, GTFS-Realtime generation,
consumer tracker status, evidence retention, release packaging, release
tagging, hosted service state, or public routes. It does not create or retain
evidence, contact consumers, contact agencies, contact vendors, or make
CAL-ITP/Caltrans compliance, agency approval/adoption, consumer
submission/review/acceptance/ingestion/listing/display, final-root readiness,
hosted SaaS, SLA/uptime, production-readiness, vendor-compatibility,
hardware-certification, public-launch, or production-grade ETA claims.

## Implementation Boundary

Use the existing private Operations Console and current sanitized view models:

- `GET /admin/operations/feed-health`
- `GET /admin/operations/feed-health.json`
- `GET /admin/operations/validation-health`
- `GET /admin/operations/validation-health.json`
- `GET /admin/operations/gtfs-quality`
- `GET /admin/operations/readiness`
- `GET /admin/operations/readiness.json`
- `GET /admin/operations/reliability`
- `GET /admin/operations/realtime`
- `GET /admin/operations/gtfs-workbench`

Phase 82 may add:

- `GET /admin/operations/validation-center`
- `GET /admin/operations/validation-center.json`

No migration is planned. The first implementation pass should compose existing
derived models:

- `operationsFeedHealthView`;
- `compliance.ValidationHealthSummary`;
- `operationsGTFSQualityGuidanceView`;
- `compliance.GTFSQualityTriage`;
- `operationsReadinessV2View`;
- `compliance.ReliabilitySummary`;
- `operationsRealtimeView`;
- existing consumer prepared tracker summaries.

If a checkpoint needs durable validation history or blocker-history support
beyond existing capped records, stop and re-plan before adding a migration.

## Security And Data Boundary

The new Center must use sanitized DTOs only. It must not reuse the raw
`/admin/validation/run` response because that path can return full
`ValidationResult.Report` data including raw reports, stdout, stderr, argv, and
private paths.

The Center must not expose:

- raw validator reports;
- stdout or stderr;
- argv, command, binary, URL, timeout, artifact, output, report, or path
  fields supplied by a browser;
- raw telemetry payloads;
- raw incident details;
- private file paths;
- private hostnames;
- DB URLs;
- bearer tokens, cookies, CSRF tokens, device tokens, or token hashes;
- raw vendor or external payloads.

If GTFS quality samples are shown in the unified drilldown, the checkpoint must
first verify or harden the sample scrubbing boundary so private paths, tokens,
auth headers, DB URLs, and raw validator output cannot appear in HTML or JSON.

The Center is read-only by default. Phase 82 should not add a new POST route.
Existing validator-run POST routes remain admin-only, body-capped,
cookie-CSRF-protected, agency-scoped, and server-owned. Any future browser
action must use the Phase 77 command result model and be separately reviewed.

## Master Approval

The Master Agent approves implementation only under these constraints from
sub-agent review:

- Add the Center as a private Health navigation entry while preserving existing
  routes.
- Keep Go server-rendered HTML as the default; buildless JavaScript may only
  enhance already-rendered private content.
- Show exactly five public feed paths in the URL panel:
  `/public/feeds.json`, `/public/gtfs/schedule.zip`,
  `/public/gtfsrt/vehicle_positions.pb`,
  `/public/gtfsrt/trip_updates.pb`, and `/public/gtfsrt/alerts.pb`.
- Explain feed health versus validation:
  feed health is output/freshness/reliability metadata, validation is content
  validator feedback, `feeds.json` is discovery metadata, and internal import
  validation is not the canonical static validator.
- Keep validator success as a supporting private signal only.
- Keep blocker rows and readiness timeline rows derived from existing private
  records unless a later migration is explicitly approved.
- Keep every claim flag false.
- Use wording such as `private validation diagnostics`, `supporting validation
  signal`, `operator readiness signals`, and `blocker history`.

## Sub-Agent Review Plan

Real or simulated reviews use the intended model levels from the authorized
Phase 75-90 track:

- Context / Repo Truth Sub-Agent, GPT-5.5 x-high: Plato reviewed existing
  feed-health, validation-health, GTFS-quality, readiness, reliability, and
  route surfaces.
- Planning Sub-Agent, GPT-5.5 x-high: Planck approved the proposed
  `validation-center` route, no-migration default, checkpoint sequence, tests,
  and stop conditions.
- Implementation Sub-Agent, GPT-5.5 high: simulated in the main rollout unless
  agent capacity becomes available.
- QA Sub-Agent, GPT-5.5 high: Aristotle defined focused private-route, JSON,
  five-feed URL, redaction, state-matrix, GTFS issue, timeline, blocker,
  consumer tracker, and protected-path tests.
- UI/UX Sub-Agent, GPT-5.5 high: Chandrasekhar proposed a dense private
  operations page organized around summary status, five feed URLs,
  feed-health-versus-validation distinction, validation history, issue triage,
  blocker timeline, and claim boundary footer.
- Documentation / IA Sub-Agent, GPT-5.5 high: Popper identified plan-time and
  closeout docs plus source-of-truth wording risks.
- Security/Auth Sub-Agent, GPT-5.5 high: Dewey required sanitized DTOs, no raw
  `/admin/validation/run` reuse, no raw validator output, no client-supplied
  execution fields, and no evidence writes.
- Claim-Boundary Sub-Agent, GPT-5.5 high: simulated locally because the agent
  thread limit was reached; the review blocks compliance, consumer acceptance,
  release readiness, validator-clean proof, hosted SaaS, production readiness,
  vendor/hardware certification, SLA/uptime, and ETA-quality claims.
- Data/Migration Sub-Agent: not planned unless implementation requires a new
  persisted model. Stop and re-plan if a migration becomes necessary.

All required edits from these reviews are incorporated into this plan.

## Checkpoints

```text
Phase 82 -- Checkpoint 000001: add feed health and validation center plan
Phase 82 -- Checkpoint 000002: unify feed status and validator history views
Phase 82 -- Checkpoint 000003: add validation issue drilldowns and fix-owner guidance
Phase 82 -- Checkpoint 000004: add readiness timeline and blocker history
Phase 82 -- Checkpoint 000005: close feed health and validation center review
```

## Acceptance Criteria

- Validation Center routes are private, authenticated, agency-scoped,
  no-store, GET-only, and unavailable under `/public/operations`.
- JSON is read-only, bounded, schema-stable, and all claim flags remain false.
- The page does not expose mutation controls, command execution controls,
  browser-supplied validator command fields, raw validator reports,
  stdout/stderr, argv, raw private paths, DB URLs, auth headers, cookies,
  bearer tokens, CSRF tokens, device tokens, raw telemetry payloads, raw
  incident details, or external submission targets.
- Five public feed URL rows are exactly ordered and distinguish configured URL
  from route path.
- Validation history distinguishes schedule static validation from realtime
  validation and internal importer validation.
- GTFS quality drilldowns include owner/fix/verify/escalation guidance while
  staying sample-capped and private-text scrubbed.
- Readiness timeline and blocker rows use only allowed statuses and include
  safe next actions plus `does not prove` boundaries.
- Consumer tracker rows remain prepared-only and do not imply submission,
  review, acceptance, listing, display, ingestion, or approval.
- Protected evidence and consumer packet paths remain untouched.
- All seven consumer targets remain exactly `prepared`.

## Validation

Baseline Phase 82 validation:

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

Additional Phase 82 checks:

```bash
go test ./cmd/agency-config -run 'ValidationCenter|FeedHealth|ValidationHealth|GTFSQuality|Readiness|Reliability|OperationsNavigation'
go test ./internal/compliance -run 'ValidationHealth|GTFSQuality|Reliability|Validation'
go test ./cmd/agency-config ./internal/compliance
make validate
make test
docker compose -f deploy/docker-compose.yml config
```

Run `RUN_LOCAL_APP=true make release-candidate-check` when route/UI changes are
in place and local app startup is safe. If an environment limitation blocks a
check, record the exact blocker in the Phase 82 handoff without converting it
into a release, compliance, consumer, production, vendor, SLA, or ETA-quality
claim.
