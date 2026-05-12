# Phase 69 -- Maintainer Product Acceptance And UI-First Agency Usability Trial

## Status

Complete for maintainer product acceptance and UI-first agency usability
scope. This phase is a product acceptance phase. It is not an evidence intake
phase, not a public launch, not a consumer submission, and not compliance
proof.

## Why This Phase Exists

Phases 61 through 67 made the private Operations Console, Connector Hub, setup
wizard, browser GTFS import, feed health dashboard, readiness dashboard,
connector guidance, telemetry simulator guide, GTFS quality guidance,
release-candidate/installability docs, accessibility-oriented polish, and
in-app help available.

Phase 68+ closed as blocker-only and authorization-gated because the required
written authorization and evidence-intake artifacts were not supplied.

Phase 69 bridges the gap between "features exist" and "a real maintainer,
small-agency operator, civic technologist, or developer integrator can evaluate
the product from a clean checkout without reading phase history." The phase
should make the repo feel like one coherent open-source product path:

1. start the local app;
2. open the private UI;
3. follow a visible Agency Operations Cockpit / Start Here path;
4. import or publish GTFS;
5. check the five public feed paths;
6. review validation health, feed health, readiness, connectors, and telemetry;
7. understand what remains before deployment, compliance, consumer, agency, or
   production claims.

## Why This Is Not Phase 68+ Evidence Intake

Phase 68+ remains closed blocker-only and authorization-gated. Phase 69 does
not collect retained evidence, contact external parties, fetch final public
roots, submit feeds to consumers, automate portals, or move any consumer status.

Any future evidence intake still requires explicit written authorization,
an exact claim target, allowed tools, public-safe retention rules, redaction
rules, and stop conditions before work begins. Formal agency approval,
final feed-root evidence, consumer acceptance, compliance proof, vendor proof,
real AVL proof, and production ETA proof are not required for local product
evaluation or open-source contribution.

## Success Conditions

Phase 69 succeeds when a maintainer can truthfully say:

> A small agency, civic technologist, or developer integrator can evaluate Open
> Transit RT from a clean checkout, start the local UI, follow a browser-first
> setup path, import GTFS, inspect public feed URLs, review feed health and
> readiness, understand connector options, and know exactly what remains before
> production deployment or compliance/consumer claims.

The phase is complete only after the UI-first path, README/wiki/docs navigation,
small-agency walkthrough, capability-versus-evidence wording, audit helpers,
and closeout validation pass review.

## Acceptance Questions

1. Can someone run this from a clean clone without the maintainer explaining it?
2. Can a small-agency operator use the private UI without living in the CLI?
3. Can a civic technologist understand and test a simple telemetry connector
   path?
4. Can the repo produce a release-candidate-style evaluation that is clear,
   safe, and honest?

## Existing UI Routes To Use

Prefer improving the existing server-rendered Operations Console routes instead
of adding new routes:

- `/admin/operations`
- `/admin/operations/launchpad`
- `/admin/operations/setup-wizard`
- `/admin/operations/gtfs-import`
- `/admin/operations/feed-health`
- `/admin/operations/readiness`
- `/admin/operations/connectors`
- `/admin/operations/telemetry-simulator`
- `/admin/operations/gtfs-quality`
- `/admin/operations/validation-health`
- `/admin/operations/reliability`
- `/admin/operations/help`

Add `/admin/operations/first-run` and `/admin/operations/first-run.json` only
if the existing dashboard and launchpad cannot support the acceptance workflow
cleanly. Any private UI work must keep GET pages authenticated/private,
Cache-Control `no-store` on private diagnostic pages touched by the workflow,
admin-only POSTs admin-only, CSRF behavior intact, role checks intact, and all
claim flags false.

## User-Facing Documentation Path

After README, the first public reader path should be:

1. `wiki/small-agency-quick-start.md`
2. `wiki/browser-first-setup.md`
3. `wiki/operations-console-tour.md`
4. `docs/tutorials/small-agency-acceptance-script.md`
5. `wiki/connector-cookbook.md`
6. `wiki/calitp-readiness-plain-english.md`
7. `docs/release-candidate-readiness.md`

`README.md` should be the product front door. `wiki/README.md` should route
nontechnical readers by task. `docs/README.md` should remain the deeper
documentation hub and separate public user docs, operator docs, integrator
docs, maintainer docs, release-candidate docs, and evidence/claim-boundary
docs.

## Stale Docs To Inspect And Clean

Phase 69 must inspect and patch stale capability-versus-evidence wording where
needed in:

- `docs/requirements-calitp-compliance.md`
- `docs/california-readiness-summary.md`
- `docs/compliance-evidence-checklist.md`
- `docs/roadmap-status.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/roadmaps/agency-first-connector-platform/README.md`
- `docs/roadmaps/agency-first-connector-platform/05-validation-and-claim-boundaries.md`
- `docs/repo-gaps.md`
- `docs/backlog.md`
- `docs/open-questions.md`

The intended distinction is:

- software capability exists for GTFS import/publication, Vehicle Positions,
  Trip Updates, Alerts, validation workflows, feed health, readiness workflows,
  telemetry ingest, and connector boundaries;
- deployment/compliance evidence is still required for stronger public claims;
- external proof tracks are optional and authorization-gated.

## Files Expected To Change

- `docs/phase-69-maintainer-product-acceptance-and-ui-first-agency-usability-trial.md`
- `docs/handoffs/phase-69.md`
- `cmd/agency-config/operations.go`
- optional new narrow `cmd/agency-config/operations_*` helper for the first-run
  model if it keeps existing routes cleaner
- `cmd/agency-config/main_test.go`
- `README.md`
- `wiki/README.md`
- `docs/README.md`
- `wiki/small-agency-quick-start.md`
- `wiki/browser-first-setup.md`
- `wiki/operations-console-tour.md`
- `docs/tutorials/small-agency-acceptance-script.md`
- selected public/wiki/docs pages listed in the checkpoint prompts
- `scripts/audit-product-acceptance.sh`
- `scripts/test-product-acceptance.sh`
- `Makefile`
- closeout/status docs listed in Checkpoint 000007

## Files Not To Edit

Do not modify:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/**`
- `docs/evidence/consumer-submissions/artifacts/**`
- `docs/evidence/consumer-submissions/packets/**`

Also avoid `db/migrations/**`, `go.mod`, and `go.sum` for this phase unless a
future maintainer explicitly changes the scope. This phase must not add DB
migrations, protobuf contract changes, telemetry ingest semantic changes,
connector schema changes, or heavy frontend dependencies.

## Checkpoints

- `Phase 69 -- Checkpoint 000001: add product acceptance and UI-first agency trial plan`
- `Phase 69 -- Checkpoint 000002: improve operations console first-run experience`
- `Phase 69 -- Checkpoint 000003: add browser-first small agency acceptance walkthrough`
- `Phase 69 -- Checkpoint 000004: clean readme wiki and docs navigation`
- `Phase 69 -- Checkpoint 000005: fix capability versus evidence docs`
- `Phase 69 -- Checkpoint 000006: add product acceptance audit helpers`
- `Phase 69 -- Checkpoint 000007: close product acceptance review`

## Execution Closeout

Phase 69 is complete for the approved product acceptance and UI-first agency
usability scope.

It improved the private Operations Console home and launchpad with a visible
Agency Operations Cockpit / Start Here path, ordered first-run tasks, five
public feed URLs, no-developer and developer paths, and local demo /
deployment / evidence-claim boundaries.
It added browser-first small-agency walkthroughs, reshaped README/wiki/docs
navigation around task-based product evaluation, clarified
capability-versus-evidence wording, and added local product acceptance audit
helpers.

Phase 69 created no retained evidence, contacted no external party, changed no
consumer status, and made no compliance, agency adoption, consumer acceptance,
final-root, hosted SaaS, production-readiness, vendor-compatibility, SLA, or
ETA-quality claim.

The product acceptance audits are local-only, read-only, no-network, no-Docker,
and do not require a running app.

## Closeout Validation

Checkpoint 000007 validation completed with:

- `git diff --check`: passed
- `make audit-product-acceptance`: passed
- `make test-product-acceptance`: passed
- `make check`: passed
- `make validate`: passed
- `make test`: passed
- `make audit-final-claim-review`: passed
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`: passed
- exact seven-target prepared-only consumer tracker check: passed
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json`: passed
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured`: passed with no output
- `docker compose -f deploy/docker-compose.yml config`: passed

The optional live local app walkthrough was not run during closeout; the local
app startup can pull/build container assets and is not required for the
read-only product acceptance audit.

## Master And Sub-Agent Review Workflow

The master agent owns scope, repo truth, checkpoint sequencing, protected-path
boundaries, and final approval.

For each checkpoint:

1. Master reads source-of-truth docs.
2. Context / Repo Truth sub-agent reports current repo truth.
3. Planning sub-agent drafts or refines the checkpoint plan.
4. Master reviews and approves the plan.
5. Implementation sub-agent executes only that checkpoint.
6. QA sub-agent validates.
7. UI/UX sub-agent reviews usability impact.
8. Documentation / Information Architecture sub-agent reviews docs consistency.
9. Claim-Boundary sub-agent reviews claims and protected paths.
10. Master decides pass, patch inside Phase 69, or blocker.

If real sub-agents are unavailable for any role, the role may be simulated in a
clearly labeled section, but the master must still perform the same review gate.

## Validation Commands

Run as many as feasible and report exact results:

```bash
git diff --check
make audit-product-acceptance
make test-product-acceptance
make check
make validate
make test
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
git diff --exit-code -- docs/evidence/consumer-submissions/status.json
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured
```

Optional if the environment supports it:

```bash
docker compose -f deploy/docker-compose.yml config
make agency-app-up
make agency-app-down
```

The product acceptance audit helper added in Checkpoint 000006 must remain
local-only, read-only, no-network, no-Docker, and must not require a running app.

## Stop Conditions

Stop or document a blocker if the work would require:

- retained evidence collection;
- final-root evidence fetching or verification;
- consumer submission or consumer status movement;
- agency, vendor, consumer, portal, or external-system contact;
- portal automation;
- public unauthenticated admin/debug/operations routes;
- protected evidence path writes;
- DB migrations;
- protobuf contract changes;
- telemetry ingest semantic changes;
- connector schema changes;
- heavy frontend stack introduction;
- fake screenshots or generated images pretending to be product screenshots;
- any stronger public claim without retained, authorized proof for that exact
  claim.

## Forbidden Claims

Phase 69 must not claim:

- CAL-ITP/Caltrans compliance;
- agency adoption or approval;
- consumer submission, review, acceptance, ingestion, listing, or display;
- agency-owned final-root readiness;
- hosted SaaS availability;
- production readiness;
- vendor compatibility;
- hardware certification;
- SLA or uptime coverage;
- production-grade ETA quality.

Phase 69 may say that the private UI-first path, README/wiki/docs navigation,
and small-agency acceptance workflow were improved for local product
evaluation, with claim boundaries preserved.
