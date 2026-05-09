SHELL := /bin/sh

DATABASE_URL ?= postgres://postgres:postgres@localhost:55432/open_transit_rt?sslmode=disable
TEST_DATABASE_URL ?= postgres://postgres:postgres@localhost:55432/open_transit_rt_test?sslmode=disable
MIGRATIONS_DIR ?= db/migrations
DOCKER_COMPOSE ?= docker compose -f deploy/docker-compose.yml

migrate-up migrate-down migrate-status migrate-redo run-telemetry-ingest test-integration: export DATABASE_URL := $(DATABASE_URL)
migrate-up migrate-down migrate-status migrate-redo test-integration: export MIGRATIONS_DIR := $(MIGRATIONS_DIR)
test-integration: export TEST_DATABASE_URL := $(TEST_DATABASE_URL)

.PHONY: build build-linux-amd64 deps db-up db-down migrate-up migrate-down migrate-status migrate-redo seed dev bootstrap demo-agency-flow agency-app-up agency-app-down agency-app-logs agency-app-reset agency-pilot-up telemetry-simulator operator-smoke support-bundle deployment-doctor validator-health collect-hosted-evidence audit-hosted-evidence pilot-ops-help run-agency-config run-telemetry-ingest run-feed-vehicle-positions run-feed-trip-updates run-feed-alerts run-gtfs-studio fmt lint test test-integration smoke validate realtime-quality validators-install validators-check oci-build oci-setup oci-push oci-units oci-deploy oci-status oci-start oci-stop oci-restart oci-logs oci-update-dns oci-collect

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

collect-hosted-evidence:
	./scripts/collect-hosted-evidence.sh

audit-hosted-evidence:
	./scripts/audit-hosted-evidence.sh

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
	gofmt -w ./cmd ./internal

test:
	go test ./...

realtime-quality:
	go test ./internal/realtimequality

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
	go test ./cmd/agency-config ./cmd/telemetry-ingest ./cmd/feed-vehicle-positions ./cmd/feed-trip-updates ./cmd/feed-alerts ./cmd/gtfs-studio ./internal/auth ./internal/devices ./internal/compliance ./internal/state

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
	@test -f tools/validators/validators.lock.json
	@test -f scripts/install-validators.sh
	@test -f scripts/check-validators.sh
	@test -f scripts/agency-local-app.sh
	@test -f scripts/agency-pilot-onboard.sh
	@test -f scripts/operator-smoke.sh
	@test -f scripts/support-bundle.sh
	@test -f scripts/deployment-doctor.sh
	@test -f scripts/validator-health.sh
	@test -f scripts/telemetry-simulator.sh
	@sh -n scripts/agency-pilot-onboard.sh
	@sh -n scripts/operator-smoke.sh
	@sh -n scripts/support-bundle.sh
	@sh -n scripts/deployment-doctor.sh
	@sh -n scripts/validator-health.sh
	@sh -n scripts/telemetry-simulator.sh
	@python3 -c 'from pathlib import Path; s=Path("scripts/deployment-doctor.sh").read_text(); assert "\"/admin/gtfs-studio\"" in s and "\"/admin/gtfs\"" not in s'
	@scripts/agency-pilot-onboard.sh --help >/dev/null
	@scripts/agency-pilot-onboard.sh --agency-id dryrun-agency --gtfs-url http://127.0.0.1/example.zip --dry-run >/dev/null
	@scripts/operator-smoke.sh --help >/dev/null
	@scripts/support-bundle.sh --help >/dev/null
	@scripts/deployment-doctor.sh --help >/dev/null
	@scripts/validator-health.sh --help >/dev/null
	@OUTPUT_DIR=.cache/validate/validator-health FORCE=true scripts/validator-health.sh --dry-run >/dev/null
	@scripts/telemetry-simulator.sh --help >/dev/null
	@scripts/telemetry-simulator.sh --list-scenarios >/dev/null
	@OUTPUT_DIR=.cache/validate/telemetry-simulator scripts/telemetry-simulator.sh --scenario on-route --dry-run --force >/dev/null
	@test -f docs/integration-adapter-kit.md
	@test -f docs/phase-39-calitp-readiness-workflow.md
	@test -f docs/handoffs/phase-39.md
	@test -f docs/tutorials/calitp-readiness-checklist.md
	@test -f docs/phase-40-guided-self-hosted-operator-trial.md
	@test -f docs/tutorials/self-hosted-operator-trial.md
	@test -f docs/handoffs/phase-40.md
	@test -f docs/phase-41-operator-smoke-support-bundle.md
	@test -f docs/tutorials/operator-smoke-and-support-bundle.md
	@test -f docs/handoffs/phase-41.md
	@test -f docs/deployment/reference-deployment-doctor.md
	@test -f docs/phase-42-reference-deployment-doctor.md
	@test -f docs/handoffs/phase-42.md
	@test -f docs/phase-43-operator-ux-setup-v2.md
	@test -f docs/handoffs/phase-43.md
	@test -f docs/phase-44-telemetry-simulator-and-device-trial.md
	@test -f docs/tutorials/telemetry-simulator-and-device-trial.md
	@test -f docs/handoffs/phase-44.md
	@test -f docs/phase-45-gtfs-quality-triage-loop.md
	@test -f docs/tutorials/gtfs-validation-triage.md
	@test -f docs/handoffs/phase-45.md
	@test -f docs/phase-46-validator-automation-and-health-gates.md
	@test -f docs/handoffs/phase-46.md
	@python3 -c 'import json; from pathlib import Path; expected=["Google Maps","Apple Maps","Transit App","Bing Maps","Moovit","Mobility Database","transit.land"]; data=json.loads(Path("docs/evidence/consumer-submissions/status.json").read_text()); records=data.get("targets", []); seen={r["target"]: r.get("status") for r in records}; assert list(seen)==expected, seen; assert all(seen[name]=="prepared" for name in expected), seen'
	@test -f testdata/avl-vendor/README.md
	@test -f testdata/avl-vendor/minimal-gps.json
	@test -f testdata/avl-vendor/full-gps.json
	@test -f testdata/avl-vendor/multi-vehicle-mapping.json
	@test -f testdata/avl-vendor/multi-vehicle-gps.json
	@test -f testdata/avl-vendor/duplicate-batch.json
	@test -f testdata/avl-vendor/out-of-order-batch.json
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
