# Phase 41 — Operator Smoke And Support Bundle

**Status:** implemented
**Previous phase:** Phase 40 — Guided Self-Hosted Operator Trial  
**Primary goal:** make the Phase 40 trial repeatable, diagnosable, and safe to share with maintainers without creating external evidence or leaking private data.

## Why This Phase Matters

Phase 40 gives operators a guided checklist. The next product gap is that the
same trial is still too manual. A small agency, operator, or civic technologist
should be able to run one safe command to answer:

- Are my five public feed paths responding?
- Is the admin readiness route reachable only through the private/admin path?
- Are validators installed, skipped, blocked, or passing?
- Does the synthetic AVL adapter dry-run still work?
- What should I send to a maintainer without leaking secrets?

This phase turns the guided trial into an operator smoke workflow and a
redaction-safe support bundle. It does not create final-root evidence,
consumer evidence, compliance evidence, agency evidence, vendor evidence, or
production evidence.

## Required Work

### 1. Add operator smoke helper

Added:

```text
scripts/operator-smoke.sh
make operator-smoke
```

The helper supports local/reference deployments with environment
variables such as:

```text
PUBLIC_BASE_URL=http://localhost:8080
ADMIN_BASE_URL=http://localhost:8080
ADMIN_TOKEN=<optional admin token>
SKIP_VALIDATORS=true|false
STRICT_VALIDATORS=true|false
OUTPUT_DIR=.cache/operator-smoke/<timestamp>
```

Implemented smoke checks:

- fetch `/public/feeds.json`;
- fetch `/public/gtfs/schedule.zip`;
- fetch `/public/gtfsrt/vehicle_positions.pb`;
- fetch `/public/gtfsrt/trip_updates.pb`;
- fetch `/public/gtfsrt/alerts.pb`;
- write checksums and sizes for fetched public files;
- check the unauthenticated admin boundary at
  `/admin/operations/readiness`, where `401`, `403`, or `404` is expected and
  any `2xx` response fails smoke;
- check pinned validator tooling state through `scripts/check-validators.sh`;
- when `ADMIN_TOKEN` is available and `ADMIN_BASE_URL` is safe, call only
  allowlisted validation APIs or record a clear skipped/blocker status;
- when `ADMIN_TOKEN` is available, check `/admin/operations/readiness` through
  `ADMIN_BASE_URL` and record only status/summary output;
- run the synthetic AVL dry-run fixture from `cmd/avl-vendor-adapter`;
- print a copy/paste support summary;
- clearly state `external_evidence_created=false` and
  `consumer_statuses_changed=false`.

### 2. Add redaction-safe support bundle helper

Added:

```text
scripts/support-bundle.sh
make support-bundle
```

The helper collects only safe diagnostics by default:

- command versions;
- git commit SHA or release version;
- OS/runtime summary;
- repo validation/validator check status;
- public feed fetch sizes/checksums/status codes;
- readiness page status only when authenticated and explicitly requested;
- selected service health endpoint results;
- migration status summary when `DATABASE_URL` is explicitly supplied;
- support manifest with included/excluded files.
- final redaction scan over generated files for secret-shaped values.

It must exclude:

- raw telemetry;
- private vendor payloads;
- device tokens;
- admin tokens;
- JWTs;
- CSRF secrets;
- DB passwords;
- cookies;
- private keys;
- ACME material;
- webhook URLs;
- notification credentials;
- unredacted logs;
- raw `.env` files;
- raw database dumps.

Output must go under ignored local storage, for example:

```text
.cache/support-bundles/<timestamp>/
```

### 3. Add operator docs

Added:

```text
docs/tutorials/operator-smoke-and-support-bundle.md
```

The tutorial should explain:

- when to use smoke checks;
- when to use support bundles;
- which values are safe to share;
- which values must never be shared;
- why support bundles are not evidence packets;
- how a future evidence phase would review/redact/retain artifacts if needed.

### 4. Update navigation and handoffs

Updated after implementation:

- `README.md`
- `docs/README.md`
- `docs/tutorials/README.md`
- `docs/current-status.md`
- `docs/backlog.md`
- `docs/open-questions.md`
- `docs/track-b-productization-roadmap.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-41.md`

## Boundaries

Phase 41 must not:

- create external evidence;
- create final-root evidence;
- contact consumers;
- change consumer statuses;
- add named vendor support;
- add real vendor payloads;
- add runtime external predictor integration;
- claim CAL-ITP/Caltrans compliance;
- claim agency adoption or approval;
- claim consumer acceptance;
- claim hosted SaaS availability;
- claim production readiness;
- claim production-grade ETA quality.

## Definition Of Done

Phase 41 is complete when:

- `make operator-smoke` exists and runs safely in local/reference mode;
- `make support-bundle` exists and writes redaction-safe local output only;
- output paths are ignored by git;
- docs explain the difference between smoke output, support bundles, and
  retained evidence;
- docs/navigation/status/handoff files are updated;
- consumer tracker statuses remain unchanged;
- all claim boundaries are preserved.

## Required Checks

Run:

```bash
make validate
make test
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
```

Run when relevant:

```bash
docker compose -f deploy/docker-compose.yml config
make smoke
make operator-smoke SKIP_VALIDATORS=true
```

If any check is blocked, document the exact reason in `docs/handoffs/phase-41.md`.

## Suggested Codex Kickoff Prompt

```text
Implement Phase 41 — Operator Smoke And Support Bundle.

Read first:
- AGENTS.md
- docs/current-status.md
- docs/handoffs/latest.md
- docs/phase-40-guided-self-hosted-operator-trial.md
- docs/tutorials/self-hosted-operator-trial.md
- docs/integration-adapter-kit.md
- docs/evidence/redaction-policy.md
- docs/roadmap-to-calitp-compliance-and-gap-closure.md
- docs/phase-41-operator-smoke-support-bundle.md

Add scripts/operator-smoke.sh and make operator-smoke. The smoke helper should
check the five public feed paths, validator tooling status, optional admin
readiness route access through ADMIN_BASE_URL with ADMIN_TOKEN, optional
allowlisted validation API calls, and the synthetic AVL dry-run fixture. It
must write local output under ignored .cache/operator-smoke/<timestamp>/ and
print a redaction-safe support summary.

Add scripts/support-bundle.sh and make support-bundle. The support bundle must
collect only safe diagnostics by default and must exclude secrets, tokens, raw
telemetry, private vendor payloads, ACME material, private keys, DB passwords,
cookies, JWTs, CSRF secrets, webhook URLs, notification credentials,
unredacted logs, raw .env files, and raw database dumps. Output must go under
ignored .cache/support-bundles/<timestamp>/.

Add docs/tutorials/operator-smoke-and-support-bundle.md and update README,
docs navigation, current-status, backlog, open-questions,
track-b-productization-roadmap, and latest handoff. Add docs/handoffs/phase-41.md
after implementation.

Do not create external evidence, final-root evidence, consumer submission
artifacts, real vendor artifacts, or stronger public claims. Do not change
consumer statuses. Do not claim CAL-ITP/Caltrans compliance, agency adoption,
consumer acceptance, hosted SaaS, production readiness, vendor compatibility,
or production-grade ETA quality.

Run and report:
make validate
make test
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null

Run docker compose -f deploy/docker-compose.yml config and make smoke if touched
surfaces make them relevant. If any command is blocked, record the exact blocker
in the Phase 41 handoff.
```
