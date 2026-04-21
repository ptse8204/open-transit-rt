SHELL := /bin/sh

DATABASE_URL ?= postgres://postgres:postgres@localhost:55432/open_transit_rt?sslmode=disable
TEST_DATABASE_URL ?= postgres://postgres:postgres@localhost:55432/open_transit_rt_test?sslmode=disable
MIGRATIONS_DIR ?= db/migrations
DOCKER_COMPOSE ?= docker compose -f deploy/docker-compose.yml

.PHONY: build deps db-up db-down migrate-up migrate-down migrate-status migrate-redo seed dev bootstrap run-agency-config run-telemetry-ingest run-feed-vehicle-positions fmt lint test test-integration validate

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

run-agency-config:
	PORT=8081 go run ./cmd/agency-config

run-telemetry-ingest:
	DATABASE_URL="$(DATABASE_URL)" PORT=8082 go run ./cmd/telemetry-ingest

run-feed-vehicle-positions:
	PORT=8083 go run ./cmd/feed-vehicle-positions

fmt:
	gofmt -w ./cmd ./internal

test:
	go test ./...

test-integration: migrate-status
	@echo "Phase 2 integration: database is reachable; DB-backed telemetry and matcher tests use isolated temporary databases when supported."
	INTEGRATION_TESTS=1 TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test ./...

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run ./...; else echo "optional lint skipped: golangci-lint is not installed; future CI should make this required once configured"; fi

validate:
	@echo "Phase 2 validation smoke: checking scaffold, durable telemetry, and deterministic matcher files only; canonical GTFS and GTFS-RT validators are documented but not wired yet."
	@test -f db/migrations/000001_initial_schema.sql
	@test -f db/migrations/000002_telemetry_ingest_foundation.sql
	@test -f db/migrations/000003_deterministic_matching.sql
	@test -d testdata/gtfs/valid-small
	@test -d testdata/gtfs/after-midnight
	@test -d testdata/gtfs/frequency-based
	@test -d testdata/gtfs/malformed
	@test -d testdata/telemetry
	@echo "Validation smoke passed. Future phases must wire canonical validators before any compliance claim."
