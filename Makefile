SHELL := /bin/sh

DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/open_transit_rt?sslmode=disable
TEST_DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/open_transit_rt_test?sslmode=disable
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
	@echo "Phase 0 has deterministic fixtures under testdata/. Runtime seeding is implemented in a later phase."

dev bootstrap:
	./scripts/bootstrap-dev.sh

run-agency-config:
	PORT=8081 go run ./cmd/agency-config

run-telemetry-ingest:
	PORT=8082 go run ./cmd/telemetry-ingest

run-feed-vehicle-positions:
	PORT=8083 go run ./cmd/feed-vehicle-positions

fmt:
	gofmt -w ./cmd ./internal

test:
	go test ./...

test-integration:
	INTEGRATION_TESTS=1 TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test ./...

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run ./...; else echo "golangci-lint not installed; skipping lint"; fi

validate:
	@echo "Static GTFS and GTFS-RT validators are documented but not wired in Phase 0."
