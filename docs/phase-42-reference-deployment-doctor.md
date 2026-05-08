# Phase 42 — Reference Deployment Doctor

**Status:** implemented  
**Previous phase:** Phase 41 — Operator Smoke And Support Bundle  
**Primary goal:** add a read-only deployment doctor for the OCI/OCL-style
reference deployment path.

## What This Phase Adds

Phase 42 adds a single command:

```sh
make deployment-doctor
```

The command runs `scripts/deployment-doctor.sh`, writes private diagnostics
under `.cache/deployment-doctor/<timestamp>/`, and reports deployment blockers,
warnings, skipped checks, and unavailable checks without printing secret values
or changing deployment state.

The doctor is intended to help an operator answer:

- Which reference deployment env keys are missing or still placeholders?
- Are generated secrets present and long enough?
- Are public feed URLs anonymously fetchable?
- Are admin/private routes blocked at the public edge?
- Are services reachable on expected loopback health endpoints?
- Can the DB and migrations be checked read-only when `DATABASE_URL` is
  supplied?
- Is PostGIS available when it can be probed safely?
- Is pinned validator tooling installed?
- Are backup and restore-drill inputs present?
- Which release/git SHA is this checkout using?

## Implemented Scope

Added:

- `scripts/deployment-doctor.sh`
- `make deployment-doctor`
- `docs/deployment/reference-deployment-doctor.md`
- `docs/handoffs/phase-42.md`

Updated validation scaffolding so `make validate` checks the script, the help
path, and the Phase 42 docs.

The script:

- uses `#!/usr/bin/env sh`, `set -eu`, `umask 077`, and no `set -x`;
- does not source private env files;
- inspects already-exported environment variables only;
- derives reference env keys from `docs/deployment/oci-reference-env.example`;
- records only statuses, not values;
- defaults output to `.cache/deployment-doctor/<timestamp>/`;
- rejects unsafe output directories and symlinks;
- writes `summary.json`, `summary.md`, `manifest.json`, and `manifest.md`;
- validates generated JSON before exit;
- runs a final generated-output redaction scan;
- exits `0` by default even when blockers are reported;
- exits nonzero on blockers only when `STRICT_DOCTOR=true`.

## Boundaries

Phase 42 does not:

- create evidence packets;
- create final-root evidence;
- write to `EVIDENCE_OUTPUT_DIR`;
- contact consumers;
- change consumer statuses;
- run migrations;
- create backups;
- run restore drills;
- read or checksum full backup files;
- claim CAL-ITP/Caltrans compliance;
- claim agency adoption or approval;
- claim hosted SaaS availability;
- claim production readiness;
- claim vendor compatibility;
- claim production-grade ETA quality.

Doctor output is private operator diagnostics. It can become evidence only
after a separate evidence phase reviews, redacts, dates, inventories, and
intentionally retains specific artifacts.

## Required Checks

Run and report:

```sh
make validate
make test
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
git diff --exit-code -- docs/evidence/consumer-submissions/status.json
python3 - <<'PY'
import json
expected = {"Google Maps", "Apple Maps", "Transit App", "Bing Maps", "Moovit", "Mobility Database", "transit.land"}
data = json.load(open("docs/evidence/consumer-submissions/status.json", encoding="utf-8"))
seen = {r["target"]: r["status"] for r in data["targets"]}
assert set(seen) == expected, seen
assert all(status == "prepared" for status in seen.values()), seen
PY
docker compose -f deploy/docker-compose.yml config
make deployment-doctor
```

If a check cannot run, record the exact reason in
`docs/handoffs/phase-42.md`.
