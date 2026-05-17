# Consumer Discovery Metadata Example

Synthetic-only consumer/discovery connector example. It reviews feed URL and
metadata shape for local preparedness without submitting to any consumer,
contacting a portal, changing consumer tracker status, writing evidence, or
claiming acceptance.

## Run

```sh
go run ./examples/connectors/consumer-discovery-metadata
```

The command reads `fixtures/feeds.json` and prints a dry-run metadata decision.
The output keeps submission, status mutation, network send, and evidence write
flags disabled.

## Adapter Shape

Consumer/discovery connectors may prepare local metadata around:

- `/public/feeds.json`;
- static GTFS URL;
- Vehicle Positions URL;
- Trip Updates URL;
- Alerts URL;
- license/contact metadata.

They must not automate consumer submissions, mutate prepared-only tracker
status, contact external portals, or write retained evidence. Prepared metadata
is a local review aid only.

## Boundaries

- synthetic fixture data only
- no portal credentials
- no consumer contact
- no status mutation
- no evidence writes
- no submission, review, acceptance, ingestion, listing, display, compliance,
  public launch, production-readiness, or hosted-service claim
