# Phase 89 Prompt — Release-Cut Cleanup / v0.1.0-rc.1 Gate

## Goal

After the consumer-grade control-plane work, rerun a release-candidate gate for `v0.1.0-rc.1`.

This phase may prepare a candidate, but it must not tag, publish, distribute a
package, or claim release-ready unless the gate passes and the maintainer
explicitly authorizes the exact release action. Local package diagnostics may
run as review artifacts only when explicitly authorized; otherwise record them
as `blocked` or `not_checked`.

## Scope

- Clean checkout review.
- Frontend/control-plane smoke review.
- Backend validation.
- Connector conformance.
- Release package audit.
- Draft release notes.
- Known blockers matrix.
- Protected-path and claim-boundary review.

## Validation

```bash
git status --short
git diff --check
make check
make validate
make test
RUN_LOCAL_APP=true make release-candidate-check
make external-connection-check
make adapter-conformance
make test-connector-examples
make audit-product-acceptance
make audit-final-claim-review
# only with explicit maintainer authorization; otherwise record blocked/not_checked
make release-package
make audit-release-package
docker compose -f deploy/docker-compose.yml config
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
python3 - <<'PY'
import json
from pathlib import Path
expected = ["Google Maps", "Apple Maps", "Transit App", "Bing Maps", "Moovit", "Mobility Database", "transit.land"]
data = json.loads(Path("docs/evidence/consumer-submissions/status.json").read_text())
records = data.get("targets", [])
seen = {row["target"]: row.get("status") for row in records}
assert list(seen) == expected, seen
assert all(seen[name] == "prepared" for name in expected), seen
PY
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum
```

## Commits

```text
Phase 89 -- Checkpoint 000001: add post-control-plane rc1 gate plan
Phase 89 -- Checkpoint 000002: run clean-checkout local product gate
Phase 89 -- Checkpoint 000003: run frontend and accessibility gate
Phase 89 -- Checkpoint 000004: run connector and backend diagnostics gate
Phase 89 -- Checkpoint 000005: prepare rc1 notes package and blockers matrix
Phase 89 -- Checkpoint 000006: close rc1 gate review
```
