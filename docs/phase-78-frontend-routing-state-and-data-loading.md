# Phase 78 -- Frontend Routing, State, And Data Loading Plan

## Scope

Phase 78 adds a small progressive interaction layer to the private Operations
Console. The server-rendered Go templates remain the source of truth, and the
browser layer only improves review ergonomics for already-rendered private
diagnostics.

This phase does not introduce a SPA, React/Vue/Svelte, npm build chain, CDN,
external map/API key, public admin route, public diagnostic JSON route, secret
storage, evidence write, consumer tracker change, migration, release artifact,
or stronger public claim.

## Architecture Decision

Use buildless JavaScript served as an embedded, allowlisted private asset under
`/admin/operations/assets/operations.js`.

Rules:

- keep all core content, links, forms, and native details usable without
  JavaScript;
- serve the asset only through authenticated private admin routing;
- set `Cache-Control: no-store`, `Content-Type: application/javascript`, and
  `X-Content-Type-Options: nosniff`;
- do not use `http.FileServer`, filesystem path routing, CDNs, import maps, or
  generated bundles;
- use same-origin relative URLs only under `/admin/operations/`;
- do not fetch `/public/*`, `/v1/events`, `/public/gtfsrt/*.json`, or external
  URLs from the browser layer;
- do not store CSRF tokens, bearer tokens, cookies, device tokens, raw JSON
  responses, logs, private paths, hostnames, credentials, URLs, commands,
  agency IDs, vehicle IDs, or row contents in browser storage;
- allow local-only preferences only for UI state such as filter, sort, compact
  mode, and opt-in refresh interval.

## Implementation Slice

Phase 78 applies the pattern narrowly:

- copy affordances for configured/local/reference values already visible in the
  page;
- bounded table/card filtering and sorting using rendered text only;
- route-level loading, error, empty, and live-result regions;
- an explicit read-only `validation_health.refresh` control that posts the
  Phase 77 command result with form-encoded `action=refresh` plus CSRF token;
- manual refresh for private JSON summaries, with any polling disabled by
  default, opt-in, same-origin only, bounded, and paused while the page is
  hidden.

Telemetry auto-polling and reusable browser JSON from `/v1/events` are out of
scope until a sanitized private DTO exists.

## Sub-Agent Review

Real or simulated reviews use the intended model levels from the authorized
Phase 75-90 track:

- Context / Repo Truth Sub-Agent, GPT-5.5 x-high, real read-only review:
  confirmed Phase 77 is closed, the working tree/protected paths are clean,
  the Operations Console is still Go server-rendered, and relevant private JSON
  routes are under `/admin/operations`.
- Planning Sub-Agent, GPT-5.5 x-high, real read-only review: approved only if
  the runtime is private, embedded, allowlisted, buildless, no-SPA, and free of
  secret storage or public JSON reuse.
- UI/UX Sub-Agent, GPT-5.5 high, real read-only review: recommended a shared
  Review tools layer for safe copy, filters, sorting, visible reset, and live
  feedback.
- Security/Auth Sub-Agent, GPT-5.5 high, real read-only review: required
  private same-origin fetches, authenticated asset serving, no filesystem asset
  server, `nosniff`, no `/public/*` or `/v1/events` reuse, and tests around
  body caps, CSRF, and sanitized outputs before mutating refresh affordances.
- QA Sub-Agent, GPT-5.5 high, real read-only review: recommended focused Go
  tests, `node --check`, and Node built-in tests for the buildless runtime.
- Claim-Boundary Sub-Agent, GPT-5.5 high, real read-only review: required all
  loading/error/empty/refresh/copy wording to stay private, local, diagnostic,
  configured, or supporting-signal only.
- Documentation / IA Sub-Agent, GPT-5.5 high, simulated by the Master Agent:
  review source-of-truth status updates and handoff after implementation.
- Implementation Sub-Agent, GPT-5.5 high, simulated by the Master Agent:
  keep edits scoped to Operations Console assets/templates/tests/docs.

## Checkpoints

```text
Phase 78 -- Checkpoint 000001: add frontend interaction architecture plan
Phase 78 -- Checkpoint 000002: add progressive UI runtime and asset policy
Phase 78 -- Checkpoint 000003: add private task progress pattern
Phase 78 -- Checkpoint 000004: apply interaction pattern to selected routes
Phase 78 -- Checkpoint 000005: close frontend interaction review
```

## Validation

Baseline Phase 78 validation:

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

Additional Phase 78 checks:

```bash
go test ./cmd/agency-config ./internal/admincontrol ./internal/auth
go test ./cmd/agency-config -run 'TestOperations|TestValidationHealth'
node --check cmd/agency-config/operations_admin.js
node --test cmd/agency-config/operations_admin_test.mjs
make validate
make test
RUN_LOCAL_APP=true make release-candidate-check
```

If an environment limitation blocks a check, record the exact blocker in the
Phase 78 handoff.
