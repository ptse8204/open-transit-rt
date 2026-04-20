.PHONY: build run-agency-config run-telemetry-ingest run-feed-vehicle-positions fmt test

build:
	go build ./...

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
