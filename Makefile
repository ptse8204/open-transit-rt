SHELL := /bin/sh
.DEFAULT_GOAL := help

DATABASE_URL ?= postgres://postgres:postgres@localhost:55432/open_transit_rt?sslmode=disable
TEST_DATABASE_URL ?= postgres://postgres:postgres@localhost:55432/open_transit_rt_test?sslmode=disable
MIGRATIONS_DIR ?= db/migrations
DOCKER_COMPOSE ?= docker compose -f deploy/docker-compose.yml

migrate-up migrate-down migrate-status migrate-redo run-telemetry-ingest test-integration: export DATABASE_URL := $(DATABASE_URL)
migrate-up migrate-down migrate-status migrate-redo test-integration: export MIGRATIONS_DIR := $(MIGRATIONS_DIR)
test-integration: export TEST_DATABASE_URL := $(TEST_DATABASE_URL)

.PHONY: help check check-links check-stable-filter build build-linux-amd64 deps db-up db-down migrate-up migrate-down migrate-status migrate-redo seed dev bootstrap demo-agency-flow agency-app-up agency-app-down agency-app-logs agency-app-reset agency-pilot-up capture-ui-tour product-ui-smoke telemetry-simulator operator-smoke support-bundle deployment-doctor validator-health operations-notify operations-reliability oci-reference-check validate-public-feeds multi-agency-hosting test-multi-agency-hosting install-confidence test-install-confidence release-candidate-check test-release-candidate-check external-connection-check adapter-conformance gtfsrt-conformance test-connector-examples caltrans-readiness-check release-package audit-release-package test-release-package audit-vendor-equivalent-pack test-vendor-equivalent-pack collect-hosted-evidence audit-hosted-evidence collect-final-root-evidence audit-final-root-evidence test-final-root-evidence generate-compliance-evidence-packet audit-compliance-evidence-packet test-compliance-evidence-packet audit-final-claim-review test-final-claim-review audit-product-acceptance test-product-acceptance audit-product-language test-product-language audit-ui-layout audit-product-roadmap-baseline audit-operations-route-inventory test-operations-route-inventory pilot-ops-help run-agency-config run-telemetry-ingest run-feed-vehicle-positions run-feed-trip-updates run-feed-alerts run-gtfs-studio fmt lint test test-integration smoke validate realtime-quality realtime-quality-backtest validators-install validators-check oci-build oci-setup oci-push oci-units oci-deploy oci-status oci-start oci-stop oci-restart oci-logs oci-update-dns oci-collect

help:
	@printf '%s\n' 'Open Transit RT command map'
	@printf '%s\n' ''
	@printf '%s\n' 'Local evaluation:'
	@printf '%s\n' '  make agency-app-up              Start the local evaluator app package'
	@printf '%s\n' '  make capture-ui-tour            Capture current tutorial screenshots from local app'
	@printf '%s\n' '  make telemetry-simulator        Send synthetic telemetry through authenticated ingest'
	@printf '%s\n' '  make agency-pilot-up            Import a supplied public GTFS URL for local/reference review'
	@printf '%s\n' '  make agency-app-down            Stop the local evaluator app package'
	@printf '%s\n' ''
	@printf '%s\n' 'Lightweight checks:'
	@printf '%s\n' '  make check                      No-network/no-Docker/no-validator-install evaluator check'
	@printf '%s\n' '  make check-stable-filter        Verify stable branch filtering rules locally'
	@printf '%s\n' '  make audit-product-language     Check primary product wording guardrails'
	@printf '%s\n' '  make test-product-language      Test product wording guardrails'
	@printf '%s\n' '  make audit-ui-layout            Check static public/console layout guardrails'
	@printf '%s\n' '  make audit-product-roadmap-baseline Check product-quality roadmap baseline'
	@printf '%s\n' '  make test                       Go unit tests'
	@printf '%s\n' '  make validate                   Full repo validation; requires pinned validators'
	@printf '%s\n' ''
	@printf '%s\n' 'Connector checks:'
	@printf '%s\n' '  make external-connection-check  Validate connector manifests and examples'
	@printf '%s\n' '  make adapter-conformance        Run offline adapter conformance fixtures'
	@printf '%s\n' '  make gtfsrt-conformance        Test offline GTFS-RT protobuf conformance harness'
	@printf '%s\n' '  make test-connector-examples    Test synthetic connector examples'
	@printf '%s\n' ''
	@printf '%s\n' 'Release/readiness:'
	@printf '%s\n' '  make release-candidate-check    Local release-candidate diagnostic summary'
	@printf '%s\n' '  make product-ui-smoke           Render private product routes with reference settings'
	@printf '%s\n' '  make install-confidence         Run local fresh-clone/archive install diagnostics'
	@printf '%s\n' '  make oci-reference-check        Private OCI/reference deployment diagnostic summary'
	@printf '%s\n' '  make validate-public-feeds      Off-host five-feed fetch and validator diagnostic'
	@printf '%s\n' '  make test-release-candidate-check Test release-candidate diagnostic boundaries'
	@printf '%s\n' '  make release-package            Generate a local .cache source package'
	@printf '%s\n' '  make audit-release-package      Audit an existing local release package'
	@printf '%s\n' '  make test-release-package       Test local release package helpers'
	@printf '%s\n' '  make caltrans-readiness-check   Local CAL-ITP-style readiness gap summary'
	@printf '%s\n' '  make audit-final-claim-review   Read-only claim and consumer tracker audit'
	@printf '%s\n' '  make audit-product-acceptance   Read-only product acceptance path audit'
	@printf '%s\n' '  make audit-operations-route-inventory Read-only private route inventory audit'
	@printf '%s\n' ''
	@printf '%s\n' 'Environment setup:'
	@printf '%s\n' '  scripts/bootstrap-dev.sh --check Local bootstrap preflight without starting services'
	@printf '%s\n' '  make deps                       Download Go modules'
	@printf '%s\n' '  make validators-install         Install pinned validators into .cache'
	@printf '%s\n' '  make db-up                      Start local Postgres/PostGIS with Docker Compose'
	@printf '%s\n' ''
	@printf '%s\n' 'Evidence helpers are authorization-gated. Do not run collection targets unless a maintainer explicitly approves that evidence scope.'

check:
	@echo "Running lightweight no-network/no-Docker/no-validator-install checks..."
	@if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then git diff --check; else echo "Skipping git diff --check outside a git worktree."; fi
	@scripts/check-consumer-tracker.sh >/dev/null
	@for f in testdata/connectors/valid/*.json testdata/connectors/invalid/*.json examples/connectors/*/connector.json examples/connectors/*/fixtures/*.json testdata/adapter-conformance/suite.json testdata/adapter-conformance/fixtures/*.json testdata/gtfsrt-conformance/*.json testdata/telemetry-simulator/*.json; do python3 -m json.tool "$$f" >/dev/null; done
	@for s in scripts/bootstrap-dev.sh scripts/agency-local-app.sh scripts/agency-pilot-onboard.sh scripts/capture-ui-tour.sh scripts/product-ui-smoke.sh scripts/install-confidence.sh scripts/test-install-confidence.sh scripts/release-candidate-check.sh scripts/oci-reference-check.sh scripts/validate-public-feeds.sh scripts/external-connection-check.sh scripts/caltrans-readiness-check.sh scripts/audit-final-claim-review.sh scripts/audit-product-acceptance.sh scripts/audit-product-language.sh scripts/audit-ui-layout.sh scripts/audit-product-roadmap-baseline.sh scripts/check-consumer-tracker.sh scripts/check-internal-links.sh scripts/check-stable-filter.sh scripts/test-product-acceptance.sh scripts/test-product-language.sh scripts/audit-operations-route-inventory.sh scripts/test-operations-route-inventory.sh; do sh -n "$$s"; done
	@scripts/bootstrap-dev.sh --help >/dev/null
	@scripts/agency-local-app.sh --help >/dev/null
	@scripts/install-confidence.sh --help >/dev/null
	@scripts/audit-final-claim-review.sh >/dev/null
	@scripts/audit-product-acceptance.sh >/dev/null
	@scripts/audit-product-language.sh >/dev/null
	@scripts/audit-ui-layout.sh >/dev/null
	@scripts/audit-product-roadmap-baseline.sh >/dev/null
	@scripts/audit-operations-route-inventory.sh >/dev/null
	@scripts/product-ui-smoke.sh >/dev/null
	@scripts/check-internal-links.sh >/dev/null
	@scripts/check-stable-filter.sh --skip-ref-check >/dev/null
	@scripts/release-candidate-check.sh --help >/dev/null
	@scripts/caltrans-readiness-check.sh --help >/dev/null
	@echo "Lightweight check passed. Heavier follow-ups when your environment supports them: make test, make validate, make release-candidate-check, make external-connection-check, make adapter-conformance."

check-links:
	@scripts/check-internal-links.sh

check-stable-filter:
	@scripts/check-stable-filter.sh

build:
	go build ./...

build-linux-amd64:
	./scripts/oci-pilot.sh build

deps:
	go mod download

db-up:
	$(DOCKER_COMPOSE) up -d postgres

db-down:
	$(DOCKER_COMPOSE) down

migrate-up:
	@go run ./cmd/migrate up

migrate-down:
	@go run ./cmd/migrate down

migrate-status:
	@go run ./cmd/migrate status

migrate-redo:
	@go run ./cmd/migrate redo

seed:
	$(DOCKER_COMPOSE) exec -T postgres psql -U postgres -d open_transit_rt < scripts/seed-dev.sql

dev bootstrap:
	./scripts/bootstrap-dev.sh

demo-agency-flow:
	./scripts/demo-agency-flow.sh

agency-app-up:
	./scripts/agency-local-app.sh up

agency-app-down:
	./scripts/agency-local-app.sh down

agency-app-logs:
	./scripts/agency-local-app.sh logs

agency-app-reset:
	./scripts/agency-local-app.sh reset

agency-pilot-up:
	@./scripts/agency-pilot-onboard.sh

capture-ui-tour:
	@./scripts/capture-ui-tour.sh

product-ui-smoke:
	@./scripts/product-ui-smoke.sh

telemetry-simulator:
	@./scripts/telemetry-simulator.sh

operator-smoke:
	@./scripts/operator-smoke.sh

support-bundle:
	@./scripts/support-bundle.sh

deployment-doctor:
	@./scripts/deployment-doctor.sh

validator-health:
	@./scripts/validator-health.sh

operations-notify:
	@./scripts/operations-notify.sh

operations-reliability:
	@./scripts/operations-reliability.sh

oci-reference-check:
	@./scripts/oci-reference-check.sh

validate-public-feeds:
	@./scripts/validate-public-feeds.sh

multi-agency-hosting:
	@./scripts/multi-agency-hosting.sh

test-multi-agency-hosting:
	@./scripts/test-multi-agency-hosting.sh

install-confidence:
	@./scripts/install-confidence.sh

test-install-confidence:
	@./scripts/test-install-confidence.sh

release-candidate-check:
	@./scripts/release-candidate-check.sh

test-release-candidate-check:
	@go test ./cmd/agency-config -run TestReleaseCandidateCheck

external-connection-check:
	@./scripts/external-connection-check.sh

adapter-conformance:
	@go run ./cmd/adapter-conformance run --suite testdata/adapter-conformance

test-connector-examples:
	@go test ./examples/connectors/...

caltrans-readiness-check:
	@./scripts/caltrans-readiness-check.sh

release-package:
	@RELEASE_PACKAGE_ALLOW_DIRTY=$${RELEASE_PACKAGE_ALLOW_DIRTY:-true} ./scripts/release-package.sh

audit-release-package:
	@./scripts/audit-release-package.sh

test-release-package:
	@./scripts/test-release-package.sh

audit-vendor-equivalent-pack:
	@./scripts/audit-vendor-equivalent-pack.sh

test-vendor-equivalent-pack:
	@./scripts/test-vendor-equivalent-pack.sh

collect-hosted-evidence:
	./scripts/collect-hosted-evidence.sh

audit-hosted-evidence:
	./scripts/audit-hosted-evidence.sh

collect-final-root-evidence:
	./scripts/collect-final-root-evidence.sh

audit-final-root-evidence:
	./scripts/audit-final-root-evidence.sh

test-final-root-evidence:
	./scripts/test-final-root-evidence.sh

generate-compliance-evidence-packet:
	./scripts/generate-compliance-evidence-packet.sh

audit-compliance-evidence-packet:
	./scripts/audit-compliance-evidence-packet.sh

test-compliance-evidence-packet:
	./scripts/test-compliance-evidence-packet.sh

audit-final-claim-review:
	./scripts/audit-final-claim-review.sh

test-final-claim-review:
	./scripts/test-final-claim-review.sh

audit-product-acceptance:
	./scripts/audit-product-acceptance.sh

test-product-acceptance:
	./scripts/test-product-acceptance.sh

audit-product-language:
	./scripts/audit-product-language.sh

test-product-language:
	./scripts/test-product-language.sh

audit-ui-layout:
	./scripts/audit-ui-layout.sh

audit-product-roadmap-baseline:
	./scripts/audit-product-roadmap-baseline.sh

audit-operations-route-inventory:
	./scripts/audit-operations-route-inventory.sh

test-operations-route-inventory:
	./scripts/test-operations-route-inventory.sh

pilot-ops-help:
	./scripts/pilot-ops.sh help

run-agency-config:
	PORT=8081 go run ./cmd/agency-config

run-telemetry-ingest:
	@PORT=8082 go run ./cmd/telemetry-ingest

run-feed-vehicle-positions:
	PORT=8083 go run ./cmd/feed-vehicle-positions

run-feed-trip-updates:
	PORT=8084 go run ./cmd/feed-trip-updates

run-feed-alerts:
	PORT=8085 go run ./cmd/feed-alerts

run-gtfs-studio:
	PORT=8086 go run ./cmd/gtfs-studio

fmt:
	gofmt -w ./cmd ./internal ./examples

test:
	go test ./...

realtime-quality:
	go test ./internal/realtimequality

realtime-quality-backtest:
	go run ./cmd/realtime-quality-backtest --observed testdata/realtime-quality-backtest/observed-events.json --predictions testdata/realtime-quality-backtest/prediction-samples.json

gtfsrt-conformance:
	go test ./internal/gtfsrtconformance ./cmd/gtfsrt-conformance

test-integration: migrate-status
	@echo "Phase 9 production-closure integration: database is reachable; DB-backed telemetry, matcher, Vehicle Positions, GTFS import, GTFS Studio, Trip Updates diagnostics, prediction operations, Alerts, publication, compliance, device auth, assignment race, and hardening tests use isolated temporary databases when supported."
	@INTEGRATION_TESTS=1 go test ./...

validators-install:
	./scripts/install-validators.sh

validators-check:
	./scripts/check-validators.sh

smoke:
	@echo "Running hardening HTTP smoke coverage..."
	@./scripts/check-validators.sh
	go test ./cmd/agency-config ./cmd/telemetry-ingest ./cmd/feed-vehicle-positions ./cmd/feed-trip-updates ./cmd/feed-alerts ./cmd/gtfs-studio ./internal/auth ./internal/devices ./internal/compliance ./internal/state ./internal/tenant

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run ./...; else echo "optional lint skipped: golangci-lint is not installed; future CI should make this required once configured"; fi

validate:
	@echo "Phase 9 production-closure validation smoke: checking scaffold, auth, device credentials, pinned validators, telemetry, matcher, Vehicle Positions, GTFS import, GTFS Studio, Trip Updates prediction operations, Alerts, schedule publication, and compliance workflow files."
	@./scripts/check-validators.sh
	@test -f db/migrations/000001_initial_schema.sql
	@test -f db/migrations/000002_telemetry_ingest_foundation.sql
	@test -f db/migrations/000003_deterministic_matching.sql
	@test -f db/migrations/000004_gtfs_import_pipeline.sql
	@test -f db/migrations/000005_gtfs_studio_drafts.sql
	@test -f db/migrations/000006_prediction_operations.sql
	@test -f db/migrations/000007_phase_8_alerts_compliance.sql
	@test -f db/migrations/000008_production_hardening.sql
	@test -f internal/auth/jwt.go
	@test -f internal/devices/devices.go
	@test -f internal/feed/vehicle_positions.go
	@test -f internal/feed/tripupdates/trip_updates.go
	@test -f internal/feed/alerts/alerts.go
	@test -f internal/feed/schedule/schedule.go
	@test -f internal/alerts/model.go
	@test -f internal/compliance/model.go
	@test -f cmd/adapter-conformance/main.go
	@test -f cmd/agency-config/operations_launchpad.go
	@test -f tools/validators/validators.lock.json
	@test -f scripts/install-validators.sh
	@test -f scripts/check-validators.sh
	@test -f scripts/agency-local-app.sh
	@test -f scripts/agency-pilot-onboard.sh
	@test -f scripts/operator-smoke.sh
	@test -f scripts/support-bundle.sh
	@test -f scripts/deployment-doctor.sh
	@test -f scripts/validator-health.sh
	@test -f scripts/operations-notify.sh
	@test -f scripts/operations-reliability.sh
	@test -f scripts/multi-agency-hosting.sh
	@test -f scripts/test-multi-agency-hosting.sh
	@test -f scripts/release-candidate-check.sh
	@test -f scripts/oci-reference-check.sh
	@test -f scripts/validate-public-feeds.sh
	@test -f scripts/external-connection-check.sh
	@test -f scripts/caltrans-readiness-check.sh
	@test -f scripts/release-package.sh
	@test -f scripts/audit-release-package.sh
	@test -f scripts/test-release-package.sh
	@test -f scripts/audit-vendor-equivalent-pack.sh
	@test -f scripts/test-vendor-equivalent-pack.sh
	@test -f scripts/collect-final-root-evidence.sh
	@test -f scripts/audit-final-root-evidence.sh
	@test -f scripts/test-final-root-evidence.sh
	@test -f scripts/generate-compliance-evidence-packet.sh
	@test -f scripts/audit-compliance-evidence-packet.sh
	@test -f scripts/test-compliance-evidence-packet.sh
	@test -f scripts/audit-product-acceptance.sh
	@test -f scripts/test-product-acceptance.sh
	@test -f scripts/telemetry-simulator.sh
	@test -f cmd/realtime-quality-backtest/main.go
	@sh -n scripts/agency-pilot-onboard.sh
	@sh -n scripts/operator-smoke.sh
	@sh -n scripts/support-bundle.sh
	@sh -n scripts/deployment-doctor.sh
	@sh -n scripts/validator-health.sh
	@sh -n scripts/operations-notify.sh
	@sh -n scripts/operations-reliability.sh
	@sh -n scripts/multi-agency-hosting.sh
	@sh -n scripts/test-multi-agency-hosting.sh
	@sh -n scripts/release-candidate-check.sh
	@sh -n scripts/oci-reference-check.sh
	@sh -n scripts/validate-public-feeds.sh
	@sh -n scripts/external-connection-check.sh
	@sh -n scripts/caltrans-readiness-check.sh
	@sh -n scripts/release-package.sh
	@sh -n scripts/audit-release-package.sh
	@sh -n scripts/test-release-package.sh
	@sh -n scripts/audit-vendor-equivalent-pack.sh
	@sh -n scripts/test-vendor-equivalent-pack.sh
	@sh -n scripts/collect-final-root-evidence.sh
	@sh -n scripts/audit-final-root-evidence.sh
	@sh -n scripts/test-final-root-evidence.sh
	@sh -n scripts/generate-compliance-evidence-packet.sh
	@sh -n scripts/audit-compliance-evidence-packet.sh
	@sh -n scripts/test-compliance-evidence-packet.sh
	@sh -n scripts/audit-product-acceptance.sh
	@sh -n scripts/test-product-acceptance.sh
	@sh -n scripts/telemetry-simulator.sh
	@python3 -c 'from pathlib import Path; s=Path("scripts/deployment-doctor.sh").read_text(); assert "\"/admin/gtfs-studio\"" in s and "\"/admin/gtfs\"" not in s'
	@scripts/agency-pilot-onboard.sh --help >/dev/null
	@scripts/agency-pilot-onboard.sh --agency-id dryrun-agency --gtfs-url http://127.0.0.1/example.zip --dry-run >/dev/null
	@scripts/operator-smoke.sh --help >/dev/null
	@scripts/support-bundle.sh --help >/dev/null
	@scripts/deployment-doctor.sh --help >/dev/null
	@scripts/validator-health.sh --help >/dev/null
	@OUTPUT_DIR=.cache/validate/validator-health FORCE=true scripts/validator-health.sh --dry-run >/dev/null
	@scripts/operations-notify.sh --help >/dev/null
	@OUTPUT_DIR=.cache/validate/operations-notify FORCE=true VALIDATOR_HEALTH_SUMMARY=.cache/validate/missing-validator/summary.json DEPLOYMENT_DOCTOR_SUMMARY=.cache/validate/missing-doctor/summary.json scripts/operations-notify.sh --dry-run >/dev/null
	@scripts/operations-reliability.sh --help >/dev/null
	@OUTPUT_DIR=.cache/validate/operations-reliability FORCE=true VALIDATOR_HEALTH_SUMMARY=.cache/validate/missing-validator/summary.json DEPLOYMENT_DOCTOR_SUMMARY=.cache/validate/missing-doctor/summary.json OPERATIONS_NOTIFY_SUMMARY=.cache/validate/missing-notify/summary.json scripts/operations-reliability.sh --dry-run >/dev/null
	@scripts/multi-agency-hosting.sh --help >/dev/null
	@OUTPUT_DIR=.cache/validate/multi-agency-hosting FORCE=true scripts/multi-agency-hosting.sh >/dev/null
	@scripts/release-candidate-check.sh --help >/dev/null
	@scripts/oci-reference-check.sh --help >/dev/null
	@scripts/validate-public-feeds.sh --help >/dev/null
	@rm -rf .cache/validate/oci-reference-check .cache/validate/validate-public-feeds
	@OUTPUT_DIR=.cache/validate/oci-reference-check FORCE=true scripts/oci-reference-check.sh --public-base-url https://feeds.example.org --dry-run >/dev/null
	@OUTPUT_DIR=.cache/validate/validate-public-feeds FORCE=true scripts/validate-public-feeds.sh --public-base-url https://feeds.example.org --dry-run >/dev/null
	@python3 -m json.tool .cache/validate/oci-reference-check/summary.json >/dev/null
	@python3 -m json.tool .cache/validate/validate-public-feeds/summary.json >/dev/null
	@python3 -c 'import json; from pathlib import Path; s=json.loads(Path(".cache/validate/oci-reference-check/summary.json").read_text()); assert all(v is False for v in s["claim_flags"].values()), s["claim_flags"]; p=json.loads(Path(".cache/validate/validate-public-feeds/summary.json").read_text()); assert len(p["rows"]) == 5 and all(v is False for v in p["claim_flags"].values()), p'
	@rm -rf .cache/validate/release-candidate-check
	@OUTPUT_DIR=.cache/validate/release-candidate-check FORCE=true scripts/release-candidate-check.sh --dry-run >/dev/null
	@python3 -m json.tool .cache/validate/release-candidate-check/summary.json >/dev/null
	@python3 -m json.tool .cache/validate/release-candidate-check/manifest.json >/dev/null
	@python3 -c 'import json; from pathlib import Path; d=Path(".cache/validate/release-candidate-check"); assert sorted(p.name for p in d.iterdir() if p.is_file()) == ["check-log.txt","manifest.json","manifest.md","summary.json","summary.md"]; s=json.loads((d/"summary.json").read_text()); assert all(v is False for v in s["claim_flags"].values()), s["claim_flags"]'
	@scripts/external-connection-check.sh --help >/dev/null
	@scripts/external-connection-check.sh >/dev/null
	@scripts/caltrans-readiness-check.sh --help >/dev/null
	@rm -rf .cache/caltrans-readiness-check/validate
	@OUTPUT_DIR=.cache/caltrans-readiness-check/validate FORCE=true scripts/caltrans-readiness-check.sh --dry-run >/dev/null
	@python3 -m json.tool .cache/caltrans-readiness-check/validate/summary.json >/dev/null
	@python3 -m json.tool .cache/caltrans-readiness-check/validate/manifest.json >/dev/null
	@python3 -c 'import json; from pathlib import Path; d=Path(".cache/caltrans-readiness-check/validate"); assert sorted(p.name for p in d.iterdir() if p.is_file()) == ["gap-review.txt","manifest.json","manifest.md","summary.json","summary.md"]; s=json.loads((d/"summary.json").read_text()); assert all(v is False for v in s["claim_flags"].values()), s["claim_flags"]; forbidden={"ok","passed","compliant","certified","accepted","ingested","listed","displayed"}; assert not any(row["status"] in forbidden for row in s["rows"]), s["rows"]; assert s["consumer_tracker"]["prepared_only"] is True, s["consumer_tracker"]'
	@scripts/release-package.sh --help >/dev/null
	@scripts/audit-release-package.sh --help >/dev/null
	@rm -rf .cache/validate/release-package
	@RELEASE_PACKAGE_VERSION=v0.0.0-validate RELEASE_PACKAGE_OUTPUT_DIR=.cache/validate/release-package RELEASE_PACKAGE_FORCE=true RELEASE_PACKAGE_ALLOW_DIRTY=true scripts/release-package.sh >/dev/null
	@RELEASE_PACKAGE_DIR=.cache/validate/release-package scripts/audit-release-package.sh >/dev/null
	@scripts/audit-vendor-equivalent-pack.sh --help >/dev/null
	@scripts/audit-vendor-equivalent-pack.sh >/dev/null
	@scripts/collect-final-root-evidence.sh --help >/dev/null
	@scripts/audit-final-root-evidence.sh --help >/dev/null
	@rm -rf .cache/validate/final-root-evidence
	@OUTPUT_DIR=.cache/validate/final-root-evidence FORCE=true scripts/collect-final-root-evidence.sh --blocker-only >/dev/null
	@FINAL_ROOT_PACKET_DIR=.cache/validate/final-root-evidence AUDIT_MODE=blocker scripts/audit-final-root-evidence.sh >/dev/null
	@scripts/generate-compliance-evidence-packet.sh --help >/dev/null
	@scripts/audit-compliance-evidence-packet.sh --help >/dev/null
	@rm -rf .cache/validate/compliance-evidence-packet
	@COMPLIANCE_PACKET_OUTPUT_DIR=.cache/validate/compliance-evidence-packet COMPLIANCE_PACKET_FORCE=true scripts/generate-compliance-evidence-packet.sh >/dev/null
	@COMPLIANCE_PACKET_DIR=.cache/validate/compliance-evidence-packet scripts/audit-compliance-evidence-packet.sh >/dev/null
	@scripts/audit-final-claim-review.sh --help >/dev/null
	@scripts/audit-final-claim-review.sh >/dev/null
	@scripts/audit-product-acceptance.sh --help >/dev/null
	@scripts/audit-product-acceptance.sh >/dev/null
	@test -f scripts/test-final-claim-review.sh
	@test -f scripts/test-product-acceptance.sh
	@sh -n scripts/audit-final-claim-review.sh scripts/test-final-claim-review.sh scripts/audit-product-acceptance.sh scripts/test-product-acceptance.sh
	@test -f docs/phase-60-final-claim-review-and-public-closeout.md
	@test -f docs/handoffs/phase-60.md
	@test -f docs/phase-55-compliance-evidence-packet-generator.md
	@test -f docs/handoffs/phase-55.md
	@test -f docs/phase-56-multi-agency-hosting-hardening.md
	@test -f docs/handoffs/phase-56.md
	@test -f docs/phase-57-release-packaging-and-supply-chain.md
	@test -f docs/release-candidate-readiness.md
	@test -f docs/deployment/oci-reference-check.md
	@test -f docs/deployment/off-host-validation.md
	@test -f docs/tutorials/no-cli-agency-first-run.md
	@test -f docs/tutorials/small-agency-maintenance-guide.md
	@test -f docs/roadmaps/agency-first-connector-platform/adoption-productization-roadmap.md
	@test -f docs/caltrans-readiness-gap-report.md
	@test -f docs/connectors/plugin-contract.md
	@test -f docs/external-connection-readiness.md
	@test -f docs/phase-58-optional-marketplace-vendor-equivalent-pack.md
	@test -f docs/vendor-equivalent-pack/README.md
	@test -f docs/vendor-equivalent-pack/byod-hardware-intake-template.md
	@test -f docs/vendor-equivalent-pack/implementation-plan-template.md
	@test -f docs/vendor-equivalent-pack/support-boundaries-template.md
	@test -f docs/vendor-equivalent-pack/sla-kpi-template.md
	@test -f docs/vendor-equivalent-pack/procurement-response-template.md
	@test -f docs/evidence/templates/final-root-approval-template.md
	@test -f docs/evidence/templates/final-root-public-fetch-template.md
	@test -f docs/evidence/templates/final-root-validator-template.md
	@test -f docs/evidence/templates/final-root-packet-readme-template.md
	@scripts/telemetry-simulator.sh --help >/dev/null
	@scripts/telemetry-simulator.sh --list-scenarios >/dev/null
	@OUTPUT_DIR=.cache/validate/telemetry-simulator scripts/telemetry-simulator.sh --scenario on-route --dry-run --force >/dev/null
	@go run ./cmd/adapter-conformance help >/dev/null
	@go run ./cmd/adapter-conformance manifest --suite testdata/adapter-conformance >/dev/null
	@go run ./cmd/adapter-conformance run --suite testdata/adapter-conformance >/dev/null
	@go run ./cmd/gtfsrt-conformance --help >/dev/null
	@go run ./cmd/realtime-quality-backtest --help >/dev/null
	@rm -rf .cache/validate/realtime-quality-backtest
	@go run ./cmd/realtime-quality-backtest --observed testdata/realtime-quality-backtest/observed-events.json --predictions testdata/realtime-quality-backtest/prediction-samples.json --output-dir .cache/validate/realtime-quality-backtest --generated-at 2026-05-09T20:00:00Z >/dev/null
	@test -f docs/integration-adapter-kit.md
	@test -f docs/phase-39-calitp-readiness-workflow.md
	@test -f docs/handoffs/phase-39.md
	@test -f docs/tutorials/calitp-readiness-checklist.md
	@test -f docs/phase-40-guided-self-hosted-operator-trial.md
	@test -f docs/tutorials/self-hosted-operator-trial.md
	@test -f docs/handoffs/phase-40.md
	@test -f docs/phase-41-operator-smoke-support-bundle.md
	@test -f docs/tutorials/operator-smoke-and-support-bundle.md
	@test -f docs/tutorials/agency-launchpad.md
	@test -f docs/handoffs/phase-41.md
	@test -f docs/deployment/reference-deployment-doctor.md
	@test -f docs/phase-42-reference-deployment-doctor.md
	@test -f docs/handoffs/phase-42.md
	@test -f docs/phase-43-operator-ux-setup-v2.md
	@test -f docs/handoffs/phase-43.md
	@test -f docs/phase-44-telemetry-simulator-and-device-trial.md
	@test -f docs/tutorials/telemetry-simulator-and-device-trial.md
	@test -f docs/tutorials/external-adapter-conformance.md
	@test -f docs/handoffs/phase-44.md
	@test -f docs/phase-45-gtfs-quality-triage-loop.md
	@test -f docs/tutorials/gtfs-validation-triage.md
	@test -f docs/handoffs/phase-45.md
	@test -f docs/phase-46-validator-automation-and-health-gates.md
	@test -f docs/handoffs/phase-46.md
	@test -f docs/phase-47-self-hosted-operations-notifications.md
	@test -f docs/tutorials/self-hosted-operations-notifications.md
	@test -f docs/handoffs/phase-47.md
	@test -f docs/phase-48-avl-adapter-runtime-path.md
	@test -f docs/handoffs/phase-48.md
	@test -f docs/phase-49-external-predictor-runtime-adapter.md
	@test -f docs/phase-50-realtime-quality-backtesting.md
	@test -f docs/roadmap-to-calitp-compliance-and-gap-closure.md
	@scripts/check-consumer-tracker.sh >/dev/null
	@test -f testdata/avl-vendor/README.md
	@test -f testdata/avl-vendor/minimal-gps.json
	@test -f testdata/avl-vendor/full-gps.json
	@test -f testdata/avl-vendor/multi-vehicle-mapping.json
	@test -f testdata/avl-vendor/multi-vehicle-gps.json
	@test -f testdata/avl-vendor/duplicate-batch.json
	@test -f testdata/avl-vendor/out-of-order-batch.json
	@test -f testdata/avl-vendor/send-manifest.json
	@test -f testdata/connectors/valid/valid-telemetry-source.json
	@test -f testdata/connectors/valid/valid-prediction.json
	@test -f testdata/connectors/valid/valid-validator.json
	@test -f testdata/connectors/valid/valid-monitoring-export.json
	@test -f testdata/connectors/valid/valid-consumer-discovery.json
	@test -f testdata/connectors/invalid/invalid-secret.json
	@for f in testdata/connectors/valid/*.json testdata/connectors/invalid/*.json; do python3 -m json.tool "$$f" >/dev/null; done
	@test -f examples/connectors/telemetry-http-poller/README.md
	@test -f examples/connectors/telemetry-http-poller/connector.json
	@test -f examples/connectors/telemetry-http-poller/main.go
	@test -f examples/connectors/telemetry-http-poller/fixtures/observations.json
	@test -f examples/connectors/telemetry-csv-replay/README.md
	@test -f examples/connectors/telemetry-csv-replay/connector.json
	@test -f examples/connectors/telemetry-csv-replay/main.go
	@test -f examples/connectors/telemetry-csv-replay/fixtures/replay.csv
	@test -f examples/connectors/predictor-sidecar-stub/README.md
	@test -f examples/connectors/predictor-sidecar-stub/connector.json
	@test -f examples/connectors/predictor-sidecar-stub/main.go
	@test -f examples/connectors/predictor-sidecar-stub/fixtures/prediction-input.json
	@test -f examples/connectors/consumer-discovery-metadata/README.md
	@test -f examples/connectors/consumer-discovery-metadata/connector.json
	@test -f examples/connectors/consumer-discovery-metadata/main.go
	@test -f examples/connectors/consumer-discovery-metadata/fixtures/feeds.json
	@test -f examples/connectors/monitoring-export/README.md
	@test -f examples/connectors/monitoring-export/connector.json
	@test -f examples/connectors/monitoring-export/main.go
	@test -f examples/connectors/monitoring-export/fixtures/metrics.json
	@for f in examples/connectors/*/connector.json examples/connectors/*/fixtures/*.json; do python3 -m json.tool "$$f" >/dev/null; done
	@go test ./examples/connectors/...
	@test -f testdata/adapter-conformance/suite.json
	@test -f testdata/adapter-conformance/fixtures/telemetry-malformed.json
	@test -f testdata/adapter-conformance/fixtures/telemetry-stale.json
	@test -f testdata/adapter-conformance/fixtures/telemetry-future.json
	@test -f testdata/adapter-conformance/fixtures/telemetry-wrong-agency.json
	@test -f testdata/adapter-conformance/fixtures/telemetry-unknown-device.json
	@test -f testdata/adapter-conformance/fixtures/telemetry-low-quality.json
	@test -f testdata/adapter-conformance/fixtures/telemetry-duplicate.json
	@test -f testdata/adapter-conformance/fixtures/telemetry-out-of-order.json
	@test -f testdata/adapter-conformance/fixtures/telemetry-missing-required-field.json
	@test -f testdata/adapter-conformance/fixtures/telemetry-invalid-coordinate.json
	@test -f testdata/adapter-conformance/fixtures/prediction-timeout.json
	@test -f testdata/adapter-conformance/fixtures/prediction-malformed.json
	@test -f testdata/adapter-conformance/fixtures/prediction-stale.json
	@test -f testdata/adapter-conformance/fixtures/prediction-wrong-agency.json
	@test -f testdata/adapter-conformance/fixtures/prediction-low-confidence.json
	@test -f testdata/adapter-conformance/fixtures/prediction-missing-vehicle-positions-ref.json
	@test -f testdata/adapter-conformance/fixtures/prediction-public-mutation-attempt.json
	@test -f testdata/adapter-conformance/fixtures/validator-allowlist.json
	@test -f testdata/adapter-conformance/fixtures/validator-raw-command.json
	@test -f testdata/adapter-conformance/fixtures/monitoring-redaction.json
	@test -f testdata/adapter-conformance/fixtures/monitoring-no-send.json
	@test -f testdata/adapter-conformance/fixtures/monitoring-unredacted-destination.json
	@test -f testdata/adapter-conformance/fixtures/consumer-discovery-feed-url-metadata.json
	@test -f testdata/adapter-conformance/fixtures/consumer-discovery-status-mutation.json
	@test -f testdata/adapter-conformance/fixtures/consumer-discovery-submission-automation.json
	@for f in testdata/adapter-conformance/suite.json testdata/adapter-conformance/fixtures/*.json; do python3 -m json.tool "$$f" >/dev/null; done
	@test -f testdata/telemetry-simulator/README.md
	@test -f testdata/telemetry-simulator/on-route.json
	@test -f testdata/telemetry-simulator/stale.json
	@test -f testdata/telemetry-simulator/out-of-order.json
	@test -f testdata/telemetry-simulator/unknown-device.json
	@test -f testdata/telemetry-simulator/low-quality-gps.json
	@test -f testdata/telemetry-simulator/after-midnight.json
	@test -f testdata/telemetry-simulator/block-transition.json
	@for f in testdata/telemetry-simulator/*.json; do python3 -m json.tool "$$f" >/dev/null; done
	@test -f scripts/device-onboarding.sh
	@test -f scripts/pilot-ops.sh
	@test -f deploy/Dockerfile.local
	@test -f deploy/Caddyfile.local
	@python3 -c 'from pathlib import Path; import re; s=Path("deploy/Caddyfile.local").read_text(); lines=[line.strip() for line in s.splitlines() if line.strip() and not line.strip().startswith("#")]; assert "@local_root {" in lines and "path /" in lines, "Caddyfile.local must define an explicit exact local-root matcher"; assert "respond @local_root \"Open Transit RT local app is running. Public feeds are under /public/ and admin routes require auth.\" 200" in lines, "Caddyfile.local must return the local app message at exact /"; assert "respond \"not found\" 404" in lines, "Caddyfile.local must include an explicit unmatched 404 fallback"; assert not any(re.fullmatch(r"respond\s+\"[^\"]*\"\s+200", line) for line in lines), "Caddyfile.local must not contain an unconditional 200 catch-all"; assert [line for line in lines if line.startswith("respond ")][-1] == "respond \"not found\" 404", "Caddyfile.local final respond must be the unmatched 404 fallback"'
	@test -f deploy/systemd/open-transit-validator-cycle.service
	@test -f deploy/systemd/open-transit-backup.service
	@test -f deploy/systemd/open-transit-feed-monitor.service
	@test -f deploy/systemd/open-transit-scorecard-export.service
	@test -f docs/runbooks/small-agency-pilot-operations.md
	@test -f docs/tutorials/agency-first-run.md
	@test -f internal/prediction/model.go
	@test -f internal/prediction/deterministic.go
	@test -f internal/prediction/postgres_operations.go
	@test -f internal/realtimequality/replay.go
	@test -f internal/realtimequality/backtest.go
	@test -f internal/avladapter/adapter.go
	@test -f cmd/avl-vendor-adapter/main.go
	@test -f internal/gtfs/importer.go
	@test -f internal/gtfs/draft.go
	@test -f cmd/feed-vehicle-positions/main.go
	@test -f cmd/feed-trip-updates/main.go
	@test -f cmd/feed-alerts/main.go
	@test -f cmd/gtfs-import/main.go
	@test -f cmd/gtfs-studio/main.go
	@test -d testdata/gtfs/valid-small
	@test -d testdata/gtfs/after-midnight
	@test -d testdata/gtfs/frequency-based
	@test -d testdata/gtfs/malformed
	@test -d testdata/telemetry
	@test -d testdata/replay
	@test -d testdata/realtime-quality-backtest
	@test -d testdata/avl-vendor
	@echo "Validation smoke passed. Canonical validators run through server-side allowlisted IDs when configured."

# ---------------------------------------------------------------------------
# OCI Pilot targets — delegate to scripts/oci-pilot.sh
# Set OCI_HOST, OCI_USER, OCI_KEY, DUCKDNS_TOKEN as needed.
# ---------------------------------------------------------------------------

oci-build:
	./scripts/oci-pilot.sh build

oci-setup:
	./scripts/oci-pilot.sh setup

oci-push:
	./scripts/oci-pilot.sh push

oci-units:
	./scripts/oci-pilot.sh units

oci-deploy:
	./scripts/oci-pilot.sh deploy

oci-status:
	./scripts/oci-pilot.sh status

oci-start:
	./scripts/oci-pilot.sh start

oci-stop:
	./scripts/oci-pilot.sh stop

oci-restart:
	./scripts/oci-pilot.sh restart

oci-logs:
	./scripts/oci-pilot.sh logs

oci-update-dns:
	./scripts/oci-pilot.sh update-dns

oci-collect:
	./scripts/oci-pilot.sh collect
