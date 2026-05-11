# Validation And Claim Boundaries

## Baseline Validation

Run after most checkpoints:

```bash
git diff --check
make check
make audit-final-claim-review
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
```

Run the exact consumer tracker check:

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

Run when relevant:

```bash
go test ./cmd/agency-config
make test
make external-connection-check
make adapter-conformance
make test-connector-examples
docker compose -f deploy/docker-compose.yml config
```

## Protected Paths

Do not edit without explicit maintainer approval:

```text
docs/evidence/captured/**
docs/evidence/consumer-submissions/status.json
docs/evidence/consumer-submissions/current/**
docs/evidence/consumer-submissions/artifacts/**
docs/evidence/consumer-submissions/packets/**
```

## Forbidden Claims

Do not claim:

- CAL-ITP/Caltrans compliance;
- agency adoption or approval;
- consumer submission, review, acceptance, ingestion, listing, or display;
- hosted SaaS;
- paid support or SLA;
- universal production readiness;
- vendor compatibility or hardware certification;
- production-grade ETA quality.

## Allowed Language

Allowed:

- “supports CAL-ITP-style readiness workflows”;
- “helps agencies evaluate a self-hosted path”;
- “connector-ready through sidecars, manifests, command adapters, and conformance tests”;
- “prepared consumer packets exist”;
- “optional evidence tracks require explicit authorization.”
