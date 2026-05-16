# Validation And Claim Boundaries

## Baseline validation for every phase closeout

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
expected = ["Google Maps","Apple Maps","Transit App","Bing Maps","Moovit","Mobility Database","transit.land"]
data = json.loads(Path("docs/evidence/consumer-submissions/status.json").read_text())
records = data.get("targets", [])
seen = {row["target"]: row.get("status") for row in records}
assert list(seen) == expected, seen
assert all(seen[name] == "prepared" for name in expected), seen
PY
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum
```

## If code/runtime/UI/scripts/tests changed

```bash
make validate
make test
docker compose -f deploy/docker-compose.yml config
```

## Release phases

```bash
make test-release-package
make release-package
make audit-release-package
RUN_LOCAL_APP=true make release-candidate-check
```

## Connector / GTFS-RT phases

```bash
make external-connection-check
make adapter-conformance
make test-connector-examples
```

## Forbidden claims

Do not claim:

- CAL-ITP/Caltrans compliance;
- agency adoption or approval;
- consumer submission/review/acceptance/ingestion/listing/display;
- final-root readiness;
- hosted SaaS or hosted-service availability;
- paid support;
- SLA or uptime guarantee;
- production readiness;
- vendor compatibility;
- hardware certification;
- production-grade ETA quality;
- real-world ETA accuracy;
- stable release readiness.

Allowed release wording after Phase 115, if published:

```text
Public v0.1.0-rc.1 release candidate for local/self-hosted evaluation.
```
