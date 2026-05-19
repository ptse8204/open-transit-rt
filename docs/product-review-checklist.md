# Product Review Checklist

Use this before merging user-facing site, docs, or Operations Console changes.
It is a product quality checklist, not a phase ledger.

## Language

- The first screen says what the user can do next.
- Page labels use operator language: import a schedule, check feed URLs,
  connect vehicle data, fix issues, maintain deployment.
- Administrator or deployment-owner work is named directly.
- Unsupported claims are covered once in a concise Limits section or collapsed
  advanced detail.
- Primary HTML does not show raw internal safety flag names.

## Layout

- The main action, current state, or form appears before secondary diagnostics.
- Core console pages show no more than three summary cards before the action
  path.
- Large diagnostics, route inventories, and detailed caveats sit below the
  workflow or behind disclosure.
- Public pages use task flows, tables, accordions, or action rows instead of
  broad card grids when comparison or sequence is clearer.
- Keyboard focus, labels, landmarks, and mobile text wrapping remain readable.

## Product Path

- Public entry points distinguish local evaluation, self-hosted/reference
  operation, and future authorized external sharing.
- `/admin/operations*` remains private, authenticated, no-store product UI.
- Screenshots or video references match current app screens and contain no
  tokens, secrets, raw logs, or private payloads.
- Protected evidence paths and consumer statuses are unchanged unless a
  maintainer explicitly authorized that work.

## Checks

Run the lightweight gate before review:

```bash
make check
make audit-product-language
make audit-ui-layout
make product-ui-smoke
```
