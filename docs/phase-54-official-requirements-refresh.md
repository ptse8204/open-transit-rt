# Phase 54 -- Official Requirements Refresh

## Status

Planning accepted. Execution must be docs-only and must refresh requirement
mappings from current official public sources without creating compliance
evidence or a compliance claim.

## Goal

Re-check current Caltrans / Cal-ITP transit data guidance and align repository
requirement mappings with the current official public guidance where available.

This phase records what the official sources say and how Open Transit RT maps
to them. It does not prove any deployment satisfies those sources.

## Scope

- Review current official Caltrans / Cal-ITP public guidance.
- Record source title, URL, visible version/date, access date, and any access
  blocker.
- Update stale source references and requirement mappings.
- Separate compliance checks from recommended, beyond-compliance, and
  experimental guidance.
- Preserve all existing evidence, consumer-status, public-route, auth, and
  claim boundaries.

## Non-Goals

- No implementation code.
- No migrations.
- No public route or admin route changes.
- No validator behavior change.
- No evidence packet creation.
- No `docs/evidence` or `docs/evidence/captured` writes.
- No consumer contact, portal automation, submission-path guessing, target
  selection, artifact creation, or consumer status change.
- No CAL-ITP/Caltrans compliance claim.
- No final-root, agency adoption, consumer acceptance, hosted SaaS,
  production-readiness, vendor-compatibility, SLA/uptime, marketplace, or
  production-grade ETA claim.

## Files Likely To Change

- `docs/phase-54-official-requirements-refresh.md`
- `docs/requirements-calitp-compliance.md`
- `docs/california-readiness-summary.md`
- `docs/compliance-evidence-checklist.md`
- `docs/roadmap-to-calitp-compliance-and-gap-closure.md`
- `docs/repo-gaps.md`
- `docs/backlog.md`
- `docs/open-questions.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-54.md`

## Safety Boundaries

Official public webpages may be cited as references in docs. They must not be
treated as retained deployment evidence.

Execution must not write under `docs/evidence`, `docs/evidence/captured`,
`docs/evidence/consumer-submissions/current`, or
`docs/evidence/consumer-submissions/artifacts`.

Execution must not edit `docs/evidence/consumer-submissions/status.json`.

Do not add new public or admin routes, relax auth, alter public feed URLs, or
change runtime behavior.

## Evidence And Claim Boundaries

Refreshed source mapping is not compliance evidence.

Validator success alone does not prove compliance or consumer acceptance.
Prepared packets are not submissions. OCI pilot evidence is not agency-owned
final-root proof. Final-root evidence, if later retained, still does not prove
consumer acceptance or compliance by itself.

All seven consumer and aggregator targets must remain `prepared`.

## Official Source Strategy

Prefer official public Caltrans pages under `dot.ca.gov/cal-itp`.

Sources verified during planning on 2026-05-09:

- Caltrans, California Transit Data Guidelines:
  `https://dot.ca.gov/cal-itp/california-transit-data-guidelines`
- Caltrans, California Transit Data Guidelines FAQ v4.0:
  `https://dot.ca.gov/cal-itp/california-transit-data-guidelines-faqs-v4_0`
- Caltrans, California Integrated Travel Project (Cal-ITP) GTFS overview:
  `https://dot.ca.gov/cal-itp/cal-itp-gtfs`
- Caltrans, Critical GTFS Validation Errors:
  `https://dot.ca.gov/cal-itp/critical-gtfs-validation-errors`
- Caltrans, Website Model Language:
  `https://dot.ca.gov/cal-itp/website-model-language`
- FTA, 2025 NTD Reporting Policy Manual:
  `https://www.transit.dot.gov/ntd/2025-ntd-reporting-policy-manual`

The execution pass should cite source facts conservatively. Treat GTFS.org and
MobilityData pages as supporting references only where Caltrans points to them.

If an official source cannot be reached during execution, document the URL,
time, and error, and leave any unsupported mapping unchanged or marked blocked.

## Implementation Details

Add or update a Phase 54 source review section that records:

- source title;
- URL;
- visible version/date when available;
- access date;
- relevant guidance category;
- repository mapping impact.

Update mappings for at least:

- stable public GTFS Schedule URL;
- stable public GTFS Realtime URLs;
- all three standard GTFS Realtime feed types: Trip Updates, Vehicle Positions,
  and Service Alerts;
- canonical validator/no-error expectations;
- major trip-planner acceptance as a separate third-party requirement;
- open license visibility;
- provider website source-of-truth links;
- technical contact / feed contact;
- aggregator availability through Transitland and Mobility Database;
- realtime service completeness and trip ID consistency;
- API-key registration constraints if a deployment chooses authenticated
  realtime access.

Keep marketplace/vendor-equivalence requirements separate from Caltrans data
guideline mapping.

Checkpoint strategy:

- `Phase 54 -- Checkpoint 000001: add official requirements refresh plan`
- `Phase 54 -- Checkpoint 000002: refresh official requirements mapping`
- `Phase 54 -- Checkpoint 000003: close requirements refresh handoff`

## Tests

Required docs-only checks:

- `make validate`
- `make test`
- `git diff --check`
- JSON syntax check for the consumer tracker
- exact seven-target prepared-only consumer tracker check
- diff guards proving consumer tracker, consumer current/artifact dirs, and
  captured evidence are unchanged
- review changed docs for unsupported compliance, consumer, final-root,
  readiness, hosted SaaS, vendor, SLA, or ETA-quality claims

No unit test changes are expected because Phase 54 is docs-only.

## Performance And Scale Tests

Not relevant. Phase 54 changes no runtime path, feed generation, DB behavior,
or validator execution behavior.

## Docs, Status, And Handoff Updates

Close the phase by updating:

- `docs/handoffs/phase-54.md`
- `docs/handoffs/latest.md`
- `docs/current-status.md`
- `docs/backlog.md`
- `docs/open-questions.md`

The handoff must say that Phase 54 refreshed official-source mappings only and
did not create evidence, change consumer statuses, or claim compliance.

## Required Verification Commands

Run and report:

```bash
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
git diff --exit-code -- docs/evidence/consumer-submissions/current docs/evidence/consumer-submissions/artifacts docs/evidence/captured
find docs/evidence/consumer-submissions/artifacts -mindepth 2 -maxdepth 2 -type f ! -name README.md -print
docker compose -f deploy/docker-compose.yml config
```

The `find` command must print no files.
