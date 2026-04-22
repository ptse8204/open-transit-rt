SHELL := /bin/sh

DATABASE_URL ?= postgres://postgres:postgres@localhost:55432/open_transit_rt?sslmode=disable
TEST_DATABASE_URL ?= postgres://postgres:postgres@localhost:55432/open_transit_rt_test?sslmode=disable
MIGRATIONS_DIR ?= db/migrations
DOCKER_COMPOSE ?= docker compose -f deploy/docker-compose.yml

.PHONY: build deps db-up db-down migrate-up migrate-down migrate-status migrate-redo seed dev bootstrap demo-agency-flow run-agency-config run-telemetry-ingest run-feed-vehicle-positions run-feed-trip-updates run-feed-alerts run-gtfs-studio fmt lint test test-integration smoke validate validators-install validators-check

build:
	go build ./...

deps:
	go mod download

db-up:
	$(DOCKER_COMPOSE) up -d postgres

db-down:
	$(DOCKER_COMPOSE) down

migrate-up:
	DATABASE_URL="$(DATABASE_URL)" MIGRATIONS_DIR="$(MIGRATIONS_DIR)" go run ./cmd/migrate up

migrate-down:
	DATABASE_URL="$(DATABASE_URL)" MIGRATIONS_DIR="$(MIGRATIONS_DIR)" go run ./cmd/migrate down

migrate-status:
	DATABASE_URL="$(DATABASE_URL)" MIGRATIONS_DIR="$(MIGRATIONS_DIR)" go run ./cmd/migrate status

migrate-redo:
	DATABASE_URL="$(DATABASE_URL)" MIGRATIONS_DIR="$(MIGRATIONS_DIR)" go run ./cmd/migrate redo

seed:
	$(DOCKER_COMPOSE) exec -T postgres psql -U postgres -d open_transit_rt < scripts/seed-dev.sql

dev bootstrap:
	./scripts/bootstrap-dev.sh

demo-agency-flow:
	./scripts/demo-agency-flow.sh

run-agency-config:
	PORT=8081 go run ./cmd/agency-config

run-telemetry-ingest:
	DATABASE_URL="$(DATABASE_URL)" PORT=8082 go run ./cmd/telemetry-ingest

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

test-integration: migrate-status
	@echo "Phase 9 production-closure integration: database is reachable; DB-backed telemetry, matcher, Vehicle Positions, GTFS import, GTFS Studio, Trip Updates diagnostics, prediction operations, Alerts, publication, compliance, device auth, assignment race, and hardening tests use isolated temporary databases when supported."
	INTEGRATION_TESTS=1 TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test ./...

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
	@test -f internal/prediction/model.go
	@test -f internal/prediction/deterministic.go
	@test -f internal/prediction/postgres_operations.go
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
	@echo "Validation smoke passed. Canonical validators run through server-side allowlisted IDs when configured."
