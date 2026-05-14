# Phase 83 -- Connector Workbench Plan

## Scope

Phase 83 adds a private browser-first Connector Workbench for choosing,
reviewing, and testing local/synthetic connector paths without claiming real
vendor compatibility or production connector readiness.

The phase may improve:

- a private `Connector Workbench` route and optional private JSON view;
- recipe chooser rows for CSV telemetry, GPS/API polling, webhook transforms,
  synthetic telemetry, external prediction sidecars, monitoring/export, and
  off-host validation;
- manifest registry review using committed example manifests only;
- dry-run command cards for operator-shell execution outside the browser;
- telemetry normalization preview from committed synthetic fixtures;
- webhook and vendor-transform boundary guidance without real vendor payloads;
- prediction sidecar shadow/fail-closed review while keeping Vehicle Positions
  independent of predictor availability;
- validator and monitoring/export recipes;
- synthetic adapter-conformance coverage guidance from committed testdata.

The phase does not contact external systems, execute connector commands from a
browser request, start sidecars, upload manifests, load arbitrary backend
plugins, send telemetry from the Workbench, add public admin routes, add
consumer submission APIs, write retained evidence, change consumer tracker
statuses, tag, package, publish, or make CAL-ITP/Caltrans compliance, agency
adoption/approval, consumer submission/review/acceptance/ingestion/listing/
display, final-root readiness, hosted SaaS, SLA/uptime, production-readiness,
vendor-compatibility, hardware-certification, public-launch, real-world
ETA-accuracy, or production-grade ETA claims.

## Implementation Boundary

Use the existing private Operations Console, existing connector examples, and
current adapter-conformance data:

- `GET /admin/operations/connectors`
- `GET /admin/operations/connectors.json`
- `GET /admin/operations/connectors/tests`
- `GET /admin/operations/connectors/tests.json`
- `examples/connectors/**`
- `testdata/connectors/**`
- `testdata/adapter-conformance/**`
- `internal/connectors`

Phase 83 may add:

- `GET /admin/operations/connectors/workbench`
- `GET /admin/operations/connectors/workbench.json`

No migration is planned. The first implementation pass should compose static
or derived in-memory DTOs from committed examples, manifests, and synthetic
fixtures. If a checkpoint needs durable connector run history, uploaded
manifests, generated artifact retention, or command-result persistence, stop
and re-plan before adding a migration or storage model.

## Security And Data Boundary

The Workbench is private and read-only by default. It must not add a POST
route. It must use the existing authentication, role, agency-query, and
no-store patterns from Connector Hub and Connector Tests.

The Workbench must not expose:

- raw vendor payloads;
- credentials, API tokens, bearer tokens, cookies, CSRF tokens, device tokens,
  token hashes, auth headers, or DB URLs;
- raw validator reports, stdout, stderr, argv, command output, private
  filesystem paths, private hostnames, or browser-supplied paths;
- external webhook destinations, live GPS/API URLs, or target portal URLs as
  executable browser actions;
- browser-side "test connection", "send telemetry", "start sidecar", "run
  connector", "submit packet", or "contact vendor" controls.

The Workbench may show fixed local/operator-shell commands only as copyable
instructions. Viewing the page executes nothing and contacts nobody.

Security review found one required hardening item before expanding manifest
display surfaces: manifest validation must reject unsafe URL/private-endpoint
strings across all displayable manifest fields, not only selected docs links or
contract fields. Phase 83 must resolve that item before or alongside the new
Workbench route.

## Master Approval

The Master Agent approves implementation only under these constraints from
sub-agent review:

- Keep Connector Hub as the registry/category overview, Connector Workbench as
  guided recipe review, and Connector Tests as fixed offline command guidance.
- Add the Workbench under the existing private Connectors navigation group.
- Keep Go server-rendered HTML as source of truth. Buildless JavaScript may
  only enhance copy buttons, filters, sorting, or local UI preferences.
- Label the registry as an example manifest registry, not installed
  connectors.
- Show all recipe choices with `what this is`, `what you need`, `what to run`,
  `good result`, `if it fails`, and `does not prove`.
- Keep conformance status local/synthetic with statuses such as `covered`,
  `needs_review`, `missing`, `blocked`, and `unknown`.
- Keep Vehicle Positions independent from predictor availability and describe
  prediction review without ETA-quality claims.
- Keep validator recipes behind server-owned allowlisted IDs and monitoring
  examples no-send/redacted by default.
- Keep every claim flag false.
- Use wording such as `private connector planning`, `local/synthetic
  conformance`, `operator-shell dry run`, `no-send preview`, and
  `fail-closed review`.

## Sub-Agent Review Plan

Real or simulated reviews use the intended model levels from the authorized
Phase 75-90 track:

- Context / Repo Truth Sub-Agent, GPT-5.5 x-high: simulated locally because
  the queued agent did not return in time; review confirms Phase 82 is closed,
  Phase 83 is next, protected paths are gated, consumer tracker must stay
  prepared-only, and existing Connector Hub/Tests routes are private/read-only.
- Planning Sub-Agent, GPT-5.5 x-high: Planck approved a private
  `/admin/operations/connectors/workbench` route, no-migration default,
  checkpoint sequence, tests, and stop conditions.
- Implementation Sub-Agent, GPT-5.5 high: simulated in the main rollout unless
  agent capacity becomes available.
- QA Sub-Agent, GPT-5.5 high: Aristotle defined focused route, JSON, no-form,
  no-browser-send, recipe, redaction, conformance-only, prepared-tracker, and
  protected-path tests.
- UI/UX Sub-Agent, GPT-5.5 high: Chandrasekhar proposed the Workbench IA,
  section structure, nontechnical copy pattern, no-SPA expectation, and
  accessibility requirements.
- Documentation / IA Sub-Agent, GPT-5.5 high: Popper identified plan-time,
  connector-doc, tutorial, and closeout alignment work plus safe wording.
- Security/Auth Sub-Agent, GPT-5.5 high: Dewey confirmed existing connector
  route posture and required display-field manifest URL/private-endpoint
  hardening before expanding Workbench surfaces.
- Claim-Boundary Sub-Agent, GPT-5.5 high: simulated locally; review blocks
  vendor compatibility, hardware certification, real vendor/device proof,
  evidence, compliance, consumer acceptance, hosted service, release readiness,
  SLA/uptime, and production-grade ETA claims.
- Data/Migration Sub-Agent: not planned because Phase 83 should add no
  migration or destructive persisted model. Stop and re-plan if persistence
  becomes necessary.

All required edits from these reviews are incorporated into this plan.

## Checkpoints

```text
Phase 83 -- Checkpoint 000001: add connector workbench plan
Phase 83 -- Checkpoint 000002: add connector recipe chooser and manifest review
Phase 83 -- Checkpoint 000003: add CSV and API telemetry connector sandbox
Phase 83 -- Checkpoint 000004: add webhook and vendor transform boundary guidance
Phase 83 -- Checkpoint 000005: add predictor and monitoring connector recipe UI
Phase 83 -- Checkpoint 000006: add synthetic conformance runner guidance
Phase 83 -- Checkpoint 000007: close connector workbench review
```

## Acceptance Criteria

- Workbench routes are private, authenticated, agency-scoped, no-store,
  GET-only, and unavailable under `/public/operations`.
- JSON is read-only, bounded, schema-stable, and all claim flags remain false.
- The page has no forms, no POST route, no command execution control, no
  browser network send, no external URL test, no sidecar start action, no
  evidence write, and no consumer status mutation.
- Recipe chooser includes the seven required operator stories: CSV vehicle
  locations, GPS API, AVL POST transform, synthetic telemetry only, external
  predictor, monitoring summaries, and off-host validation.
- Manifest review uses committed examples only and does not describe manifests
  as dynamically loaded backend plugins.
- Displayed manifest fields reject or scrub unsafe URL/private-endpoint,
  private-path, credential, token, DB URL, and raw payload content.
- Dry-run commands are fixed operator-shell instructions and not browser-run
  actions.
- Telemetry normalization preview uses committed synthetic fixtures only and is
  labeled as preview output, not accepted telemetry, Vehicle Positions output,
  real vendor proof, or production AVL proof.
- Conformance guidance is local/synthetic and does not imply real connector,
  vendor, hardware, SLA, production, compliance, consumer, release, or ETA
  proof.
- Consumer tracker rows remain prepared-only and do not imply submission,
  review, acceptance, listing, display, ingestion, or approval.
- Protected evidence and consumer packet paths remain untouched.
- All seven consumer targets remain exactly `prepared`.

## Validation

Baseline Phase 83 validation:

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

Additional Phase 83 checks:

```bash
go test ./cmd/agency-config -run 'Connector|Workbench|OperationsNavigation'
go test ./internal/connectors
go test ./cmd/adapter-conformance
go test ./examples/connectors/...
make validate
make test
make external-connection-check
make adapter-conformance
make test-connector-examples
docker compose -f deploy/docker-compose.yml config
```

Run `RUN_LOCAL_APP=true make release-candidate-check` when route/UI changes are
in place and local app startup is safe. If an environment limitation blocks a
check, record the exact blocker in the Phase 83 handoff without converting it
into a release, compliance, consumer, production, vendor, hardware, SLA, or
ETA-quality claim.
