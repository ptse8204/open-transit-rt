# Commands Run

- Environment: `local-demo`
- Capture date (UTC): 2026-04-22
- Operator: Codex local run

## Required Repo Checks

| Command | Result |
| --- | --- |
| `make validators-check` | Passed. Pinned static GTFS JAR and Docker-backed GTFS-RT wrapper installed. |
| `make validate` | Passed. Scaffold and pinned validator tooling check passed. |
| `make test` | Passed. `go test ./...` completed successfully. |
| `make smoke` | Passed. Hardening HTTP smoke packages completed successfully. |
| `make demo-agency-flow` | Passed. Local demo imported GTFS, started services, fetched public feeds, checked protected routes, ran validation flow, and exported scorecard. |
| `make test-integration` | Passed. Migration status succeeded and `INTEGRATION_TESTS=1 go test ./...` completed successfully. |
| `docker compose -f deploy/docker-compose.yml config` | Passed. Compose renders PostGIS service on host port `55432`. |
| `git diff --check` | Passed after documentation and evidence edits. |

## Evidence Collection Commands

Important commands used to collect this packet:

```sh
make demo-agency-flow
```

```sh
curl -sS -D - -o "$out" -w 'curl_status=%{http_code}\ncontent_type=%{content_type}\nsize_download=%{size_download}\ntime_total=%{time_total}\n' "$url"
```

```sh
curl -sS -X POST http://localhost:8081/admin/validation/run \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  --data '{"validator_id":"realtime-mobilitydata","feed_type":"trip_updates"}'
```

```sh
curl -sS -X POST http://localhost:8081/admin/validation/run \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  --data '{"validator_id":"realtime-mobilitydata","feed_type":"alerts"}'
```

```sh
docker compose -f deploy/docker-compose.yml exec -T postgres pg_dump -U postgres -d open_transit_rt --no-owner --no-privileges > "$DUMP"
docker compose -f deploy/docker-compose.yml exec -T postgres dropdb -U postgres --if-exists open_transit_rt_restore_drill_20260422
docker compose -f deploy/docker-compose.yml exec -T postgres createdb -U postgres open_transit_rt_restore_drill_20260422
docker compose -f deploy/docker-compose.yml exec -T postgres psql -U postgres -d open_transit_rt_restore_drill_20260422 < "$DUMP"
```

```sh
curl -sS -X POST http://localhost:8081/admin/compliance/scorecard \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  --data '{}'
```

## Environment Notes

- Go: `go version go1.26.2 darwin/amd64`.
- Docker: Docker version `29.4.0`, Compose version `v5.1.2`.
- Java: command exists at `/usr/bin/java`, but no Java Runtime is installed. Static validator execution failed because of this.
- Pinned GTFS-RT validator image is installed and checked, but runtime validation failed because the current wrapper invocation passed unsupported CLI arguments to the image.
