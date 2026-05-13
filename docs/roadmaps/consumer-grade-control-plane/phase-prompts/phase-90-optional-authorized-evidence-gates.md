# Phase 90+ Prompt — Optional Authorized Evidence Gates

## Goal

Collect claim-specific retained evidence only when explicitly authorized.

## Required authorization before any Phase 90+ work

The maintainer must provide:

```text
exact claim target
allowed tools
public-safe retention rules
redaction rules
stop conditions
operator/agency authorization
specific target or root if applicable
```

## Possible tracks

- `Phase 90A` — final-root evidence gate.
- `Phase 90B` — authorized consumer submission gate.
- `Phase 90C` — real agency pilot gate.
- `Phase 90D` — real vendor/device AVL gate.
- `Phase 90E` — real-world ETA quality gate.
- `Phase 90F` — deployment/compliance packet gate.

## Default outcome without authorization

If any required authorization element is missing, close blocker-only. Do not collect evidence, contact anyone, fetch final roots, submit to portals, change consumer statuses, or write protected evidence.

## Validation

Baseline validation plus evidence-track-specific audits only when authorized.

## Commit pattern

```text
Phase 90X -- Checkpoint 000001: add authorized <track> plan
Phase 90X -- Checkpoint 000002: close <track> blocker
```

or, when fully authorized and evidence exists:

```text
Phase 90X -- Checkpoint 000001: add authorized <track> plan
Phase 90X -- Checkpoint 000002: collect redacted <track> evidence
Phase 90X -- Checkpoint 000003: audit <track> evidence and claims
Phase 90X -- Checkpoint 000004: close <track> review
```

Consumer status movement has an additional hard gate: it requires retained,
redacted, target-originated evidence for the named target and exact feed scope.
Operator authorization alone is not enough to move a target beyond
`prepared`.
