# Validation And Claim Boundaries

## Baseline validation

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

## Heavier validation when code changes

```bash
make validate
make test
docker compose -f deploy/docker-compose.yml config
```

## Connector validation

```bash
make external-connection-check
make adapter-conformance
make test-connector-examples
```

## Release-candidate validation

```bash
RUN_LOCAL_APP=true make release-candidate-check
```

Phase 95 specifically authorizes local `.cache` package generation and audit
when existing repo tooling supports it:

```bash
make release-package
make audit-release-package
```

No phase authorizes `git tag`, `git push --tags`, GitHub Release creation,
image publication, package registry publication, or a release-ready claim.

## Forbidden claims

No compliance, adoption, consumer action, final-root readiness, hosted SaaS,
paid support, SLA/uptime, production readiness, vendor compatibility, hardware
certification, production-grade ETA, real-world ETA accuracy, public launch, or
release-ready claim without separate proof and authorization.
