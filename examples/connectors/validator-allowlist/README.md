# Validator Allowlist Example

Synthetic-only validator connector example. It demonstrates how a connector can
accept a server-owned validator ID/feed-type mapping and safe `fixture://`
artifact reference without accepting raw validator commands, private paths, or
operator-supplied executable arguments.

## Run

```sh
go run ./examples/connectors/validator-allowlist
```

The command reads `fixtures/request.json` and prints a dry-run allowlist
decision. The output records that raw validator commands, network sends, status
mutation, and evidence writes are disabled.

## Adapter Shape

Open Transit RT validator execution remains server-owned. Operators should use
allowlisted validator IDs such as `static-mobilitydata` for `schedule` and
`realtime-mobilitydata` for `vehicle_positions`, `trip_updates`, and `alerts`.
Connector manifests and browser workflows must not carry executable command
strings or private artifact paths.

An allowlisted decision is only an adapter safety signal. It is not a
validator-clean feed result, CAL-ITP/Caltrans compliance proof, consumer
acceptance, agency approval, production readiness, or public launch evidence.

## Boundaries

- synthetic fixture data only
- no real credentials
- no private artifact paths or non-`fixture://` artifact references
- no raw validator command strings
- no network sends
- no status mutation
- no evidence writes
- no compliance or consumer acceptance claim
