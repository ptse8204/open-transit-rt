# Phase 81 -- Realtime Operations Center Plan

## Scope

Phase 81 adds a private browser-first Realtime Operations Center that brings
telemetry freshness, device status, conservative assignment state, Vehicle
Positions, Trip Updates diagnostics, Alerts lifecycle context, realtime issue
queues, and simulator guidance into one operator surface.

The phase may improve:

- a private `Realtime Operations Center` route and optional private JSON view;
- fleet freshness summaries for fresh, stale, not-seen, and no-telemetry
  states;
- device binding status and safe next actions without exposing tokens;
- assignment explanation for matched, unknown, low-confidence, ambiguous,
  degraded, stale, and manual-override cases;
- Vehicle Positions status that stays independent of Trip Updates and external
  predictors;
- Trip Updates diagnostics that expose emitted counts, withheld reasons,
  stale/unknown/ambiguous rates, and conservative fallback behavior;
- Alerts lifecycle context that links to the existing Alerts Console without
  claiming consumer display or agency approval;
- a read-only issue queue for stale telemetry, unknown or low-confidence
  assignment, withheld Trip Updates, missing feed/validator signals, and Alerts
  follow-up;
- simulator guidance that remains synthetic/local and operator-shell based.

The phase does not change `/v1/telemetry`, public GTFS-RT feed semantics,
device-token verification, Alerts Console mutation semantics, consumer tracker
status, evidence retention, release packaging, release tagging, hosted service
state, or public routes. It does not add a public fleet map, browser command
execution, arbitrary backend plugin loading, external portal contact, vendor
claim, hardware certification claim, SLA claim, production-readiness claim,
consumer acceptance claim, compliance claim, public launch claim, or
production-grade ETA claim.

## Implementation Boundary

Use the existing private Operations Console and realtime sources:

- `GET /admin/operations/telemetry`
- `GET /admin/operations/devices`
- `GET /admin/operations/telemetry-simulator`
- `GET /admin/operations/feed-health`
- `GET /admin/operations/feeds`
- `/admin/alerts/console` as the existing Alerts lifecycle surface.

Phase 81 may add:

- `GET /admin/operations/realtime`
- `GET /admin/operations/realtime.json`

No migration is planned. The first implementation pass should compose existing
sanitized view models and repositories:

- latest accepted telemetry and current assignments;
- device bindings;
- feed discovery and feed health snapshots;
- validation health summaries;
- Trip Updates diagnostics;
- reliability summaries;
- committed simulator scenario metadata;
- existing Alerts feed and Alerts Console links.

The Realtime Center is read-only by default. If a later checkpoint proposes a
new browser action, it must use the Phase 77 command result model, be private,
role-checked, CSRF-safe for cookie auth, body-capped, agency-scoped from the
authenticated principal, bounded to a server-owned action, and separately
reviewed before implementation. Phase 81 should not add new POST routes.

## Master Approval

The Master Agent approves implementation only under these constraints from
sub-agent review:

- Add **Realtime Center** as a private Realtime navigation entry while
  preserving all existing route paths.
- Keep Go server-rendered HTML as the default; buildless JavaScript may only
  enhance already-rendered private content.
- Keep Vehicle Positions independent from Trip Updates and prediction adapter
  availability.
- Treat `unknown`, `stale`, `degraded`, `ambiguous`, and `withheld` as safe
  conservative states rather than failures that must be hidden.
- Use existing sanitized DTOs. Do not render raw telemetry payloads,
  `payload_json`, score-details blobs, token hashes, one-time tokens except the
  existing immediate device rebind result, raw validator output, private paths,
  DB URLs, auth headers, cookies, bearer tokens, or private vendor payloads.
- Keep public protobuf routes unchanged and do not add public admin/debug/fleet
  routes.
- Keep simulator guidance synthetic/local and operator-shell based; the browser
  must not execute shell commands or send telemetry.
- Use claim-safe wording such as `freshness status`, `assignment confidence`,
  `unknown/degraded reason`, `withheld Trip Updates`, `local/reference
  diagnostic`, and `adapter boundary`.

## Sub-Agent Review Plan

Real or simulated reviews use the intended model levels from the authorized
Phase 75-90 track:

- Context / Repo Truth Sub-Agent, GPT-5.5 x-high: Laplace reviewed existing
  realtime routes, stores, public feeds, data models, extension points, and
  risks.
- Planning Sub-Agent, GPT-5.5 x-high: Peirce approved the checkpoint sequence,
  route additions, no-migration boundary, tests, and validation commands.
- Implementation Sub-Agent, GPT-5.5 high: simulated in the main rollout unless
  agent capacity becomes available.
- QA Sub-Agent, GPT-5.5 high: Huygens defined focused route, JSON, privacy,
  bounded-output, state-matrix, and no-overclaim tests.
- UI/UX Sub-Agent, GPT-5.5 high: Jason proposed a dense private operations
  center organized around today status, fleet freshness, assignments, Vehicle
  Positions, Trip Updates, Alerts, and review queue.
- Documentation / IA Sub-Agent, GPT-5.5 high: simulated locally because the
  agent thread limit was reached; update the phase plan, handoff, latest
  handoff, roadmap status, current status, and maintained planner docs.
- Claim-Boundary Sub-Agent, GPT-5.5 high: Euler approved Phase 81 wording only
  if it avoids production AVL reliability, real-world ETA accuracy, consumer
  acceptance, compliance, vendor compatibility, SLA, release readiness, and
  public launch claims.
- Security/Auth Sub-Agent, GPT-5.5 high: Boole approved planning if routes stay
  private, agency-scoped, no-store, read-only, sanitized, and CSRF-safe for any
  future unsafe cookie-auth action.
- Data/Migration Sub-Agent: not planned unless implementation requires a new
  persisted model. Stop and re-plan if a migration becomes necessary.

All required edits from these reviews are incorporated into this plan.

## Checkpoints

```text
Phase 81 -- Checkpoint 000001: add realtime operations center plan
Phase 81 -- Checkpoint 000002: add fleet and telemetry freshness overview
Phase 81 -- Checkpoint 000003: add assignment and matching explanation views
Phase 81 -- Checkpoint 000004: add Trip Updates and Alerts status views
Phase 81 -- Checkpoint 000005: add realtime issue queue and simulator guidance
Phase 81 -- Checkpoint 000006: close realtime operations center review
```

## Acceptance Criteria

- Realtime Center routes are private, authenticated, agency-scoped, no-store,
  GET-only, and unavailable under `/public/operations`.
- JSON is read-only, bounded, schema-stable, and all claim flags remain false.
- The page does not expose mutation controls, command execution controls, token
  fields, raw telemetry payloads, raw score details, raw incident details,
  private paths, raw validator output, DB URLs, cookies, authorization headers,
  or external submission targets.
- Fleet status and telemetry rows are row-capped, text-capped, escaped, and
  deterministic.
- Assignment explanations include conservative state, degraded state, source,
  confidence, and reason codes where available without treating low-confidence
  matches as certain.
- Vehicle Positions status stays visible even when trip assignment is unknown.
- Trip Updates diagnostics explain withheld, missing, fallback, and
  no-eligible-prediction states without promising ETA quality.
- Alerts lifecycle context distinguishes draft/published/archive follow-up from
  consumer display or agency approval.
- Simulator guidance links to the existing synthetic/local guide and does not
  send telemetry from the browser.
- Protected evidence and consumer packet paths remain untouched.
- All seven consumer targets remain exactly `prepared`.

## Validation

Baseline Phase 81 validation:

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

Additional Phase 81 checks:

```bash
go test ./cmd/agency-config -run 'Realtime|Telemetry|Device|FeedHealth|ValidationHealth|TelemetrySimulator'
go test ./cmd/agency-config ./cmd/feed-alerts ./internal/alerts ./internal/prediction ./internal/state ./internal/telemetry ./internal/compliance
make validate
make test
make telemetry-simulator
docker compose -f deploy/docker-compose.yml config
```

Run `RUN_LOCAL_APP=true make release-candidate-check` when route/UI changes are
in place and local app startup is safe. If an environment limitation blocks a
check, record the exact blocker in the Phase 81 handoff without converting it
into a release, compliance, consumer, production, vendor, SLA, or ETA-quality
claim.
