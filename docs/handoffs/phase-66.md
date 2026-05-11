# Phase 66 Handoff -- Release Candidate And Installability

## Status

Complete.

Phase 66 made Open Transit RT easier to install, evaluate, package locally,
and review as a release candidate from a clean checkout. It kept all work
bounded to local/self-hosted evaluator workflows and did not create retained
evidence, publish artifacts, push images, change consumer statuses, or make
stronger production, hosted-service, compliance, adoption, consumer, vendor, or
ETA-quality claims.

## Checkpoints

- `Phase 66 -- Checkpoint 000001: add release candidate and installability plan`
- `Phase 66 -- Checkpoint 000002: prepare first release candidate workflow`
- `Phase 66 -- Checkpoint 000003: improve installer and bootstrap UX`
- `Phase 66 -- Checkpoint 000004: document Docker image publishing decision`
- `Phase 66 -- Checkpoint 000005: add demo site or documentation website plan`
- `Phase 66 -- Checkpoint 000006: close release candidate and installability`

## Changed Files

- `Makefile`
- `README.md`
- `cmd/agency-config/bootstrap_dev_script_test.go`
- `cmd/agency-config/release_candidate_script_test.go`
- `docs/current-status.md`
- `docs/decisions.md`
- `docs/demo-docs-site-plan.md`
- `docs/demo-video-outline.md`
- `docs/dependencies.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-66.md`
- `docs/phase-66-release-candidate-and-installability.md`
- `docs/public-launch-checklist.md`
- `docs/README.md`
- `docs/release-candidate-readiness.md`
- `docs/release-checklist.md`
- `docs/release-notes-template.md`
- `docs/release-process.md`
- `docs/roadmap-status.md`
- `docs/tutorials/agency-first-run.md`
- `docs/tutorials/deploy-with-docker-compose.md`
- `docs/tutorials/local-quickstart.md`
- `docs/upgrade-and-rollback.md`
- `scripts/agency-local-app.sh`
- `scripts/bootstrap-dev.sh`
- `scripts/release-candidate-check.sh`
- `wiki/README.md`

## Product Outcome

- `make release-candidate-check` now exposes source metadata, ordered
  release-candidate review steps, release-note inputs, and a local package
  audit matrix while still writing exactly five private `.cache` diagnostic
  files.
- Release docs now separate pre-tag release-candidate review from final
  tagging and record local package/image blockers without stronger claims.
- `scripts/bootstrap-dev.sh --check` provides a non-mutating local bootstrap
  preflight for tools, required repo paths, Docker daemon availability, Docker
  Compose config, and likely port conflicts.
- `scripts/agency-local-app.sh` now reports Docker Compose, port, build,
  existing-volume, migration, seed, and readiness blockers more clearly.
- Docker image guidance now explicitly keeps Phase 66 source/local-image only:
  no registry publication, no official published app image, no image push, and
  no hosted-service or production-image claim.
- `docs/demo-docs-site-plan.md` defines a repository-native future docs/demo
  site information architecture without building, hosting, announcing, or
  launching a website.

## Claim Boundary

Phase 66 created no retained evidence, wrote nothing under protected evidence
paths, contacted no external party, changed no consumer status, added no
public route, added no migration, changed no public feed URL, changed no
telemetry ingest contract, changed no GTFS-RT protobuf semantics, changed no
validator execution semantics, changed no connector manifest schema, changed
no prediction adapter behavior, and weakened no auth or public/private route
boundary.

It added no CAL-ITP/Caltrans compliance, final-root, agency approval, agency
adoption, consumer submission/review/acceptance/listing/display/ingestion,
hosted SaaS, paid support, service-level or uptime proof, public launch,
production readiness, vendor compatibility, hardware certification,
production AVL reliability, real realtime proof, published production image,
or production-grade ETA claim.

All seven consumer and aggregator targets remain `prepared`.

## Verification

All listed checks passed from `/Users/edwintse/Downloads/open-transit-rt`.

- `git diff --check`
- `go test ./cmd/agency-config -run 'Test(BootstrapDev|ReleaseCandidateCheck)'`
- `make check`
- `make test`
- `OUTPUT_DIR=.cache/validate/release-candidate-check FORCE=true scripts/release-candidate-check.sh --dry-run`
- `scripts/bootstrap-dev.sh --check`
- `make test-release-package`
- `make external-connection-check`
- `make adapter-conformance`
- `make test-connector-examples`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json`
- `git diff --exit-code -- docs/evidence/captured`
- `git diff --exit-code -- db/migrations go.mod go.sum`
- `docker compose -f deploy/docker-compose.yml config`

`scripts/bootstrap-dev.sh --check` reported a non-fatal local setup hint that
host port `55432` already had a listener. The preflight completed successfully
and did not start services or mutate database state.

## Next Work

Start Phase 67 -- Product Polish, Accessibility, and In-App Help.

Recommended first checkpoint:

`Phase 67 -- Checkpoint 000001: add product polish and accessibility plan`
