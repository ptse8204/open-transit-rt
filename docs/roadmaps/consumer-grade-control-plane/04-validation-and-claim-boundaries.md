# Validation And Claim Boundaries — Consumer-Grade Control Plane

## Protected paths

Do not modify or generate files under:

```text
docs/evidence/captured/**
docs/evidence/consumer-submissions/status.json
docs/evidence/consumer-submissions/current/**
docs/evidence/consumer-submissions/artifacts/**
docs/evidence/consumer-submissions/packets/**
```

## Consumer tracker check

All seven targets remain exactly `prepared`:

```bash
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
```

## Baseline validation

For docs-only planning phases:

```bash
git status --short
git diff --check
make check
make audit-product-acceptance
make audit-final-claim-review
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
# run exact seven-target prepared-only consumer tracker check
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum
```

For code/UI phases, add:

```bash
make validate
make test
docker compose -f deploy/docker-compose.yml config
```

For local app/UI route review when authorized:

```bash
make agency-app-up
# verify authenticated private routes and anonymous public feed paths
make agency-app-down
```

For release-candidate phases:

```bash
RUN_LOCAL_APP=true make release-candidate-check
make external-connection-check
make adapter-conformance
make test-connector-examples
make audit-release-package
```

## Claim-boundary phrasing

Never say:

- compliant;
- accepted;
- adopted;
- approved;
- production ready;
- hosted SaaS;
- SLA-backed;
- vendor-compatible;
- hardware-certified;
- consumer-listed;
- production-grade ETA.

Use safer alternatives:

- readiness signal;
- local/reference diagnostic;
- browser-first workflow;
- synthetic connector check;
- prepared packet;
- optional evidence track;
- requires retained evidence before stronger claims.

## UI-specific safety rules

The consumer-grade UI must not leak:

- bearer tokens;
- device secrets;
- raw validator stdout/stderr;
- raw private file paths;
- raw external payloads;
- private hostnames;
- unredacted support bundles;
- raw response bodies from external connectors.

UI command forms must:

- use authenticated private routes;
- preserve role checks;
- preserve CSRF for browser cookie flows;
- use bearer-token rules only where already safe and intended;
- cap request body size;
- reject client-supplied validator commands, paths, argv, shell fragments, output paths, and artifact paths;
- use server-owned mappings and configured timeouts;
- render derived summaries rather than raw internals.

Every UI/content checkpoint must also audit labels, badges, empty states,
screenshots, readiness panels, release notes, and help text for unsupported
claims before closeout.

## Evidence boundary

Do not convert local UI walkthroughs, screenshots, validator results, local app checks, public GTFS trial runs, connector tests, or release-candidate diagnostics into retained evidence unless a separate evidence phase has explicit written authorization and a public-safe retention plan.
