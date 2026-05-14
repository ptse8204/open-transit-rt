# Autonomous Run Policy

This policy applies only when the maintainer has explicitly authorized the
post-90 autonomous run.

## Sequence

Run phases in order:

```text
Phase 91 -> Phase 92 -> Phase 93 -> Phase 94 -> Phase 95 -> Phase 96 ->
Phase 97 -> Phase 98 -> Phase 99 -> Phase 100 -> Phase 101 -> Phase 102 ->
Phase 103 -> Phase 104 -> Phase 105 -> Phase 106 -> Phase 107 -> Phase 108 ->
Phase 109 -> Phase 110
```

After a phase closes, start the next phase immediately unless a hard-stop
condition prevents safe continuation.

## Hard Stops

Stop or convert the phase to a blocker-only closeout when a required action
would:

- modify protected evidence paths;
- move a consumer tracker status beyond `prepared`;
- contact an agency, vendor, consumer, portal, map provider, aggregator, or
  other external service for proof;
- require real credentials, real private payloads, or real vendor/agency data;
- publish a tag, GitHub Release, container image, package, or public
  announcement;
- make a forbidden public claim;
- continue past a security, auth, or data-integrity issue that makes further
  work unsafe.

If the hard stop applies only to one phase, close that phase truthfully as
blocked or `needs_review` and continue to the next safe product phase.

## Checkpoint Discipline

Every checkpoint must end in a git commit using:

```text
Phase XX -- Checkpoint 000001: add <phase> plan
Phase XX -- Checkpoint 000002: <implementation checkpoint>
Phase XX -- Checkpoint 000003: <validation or audit checkpoint>
Phase XX -- Checkpoint 000004: close <phase> review
```

Checkpoint numbers reset per phase. Do not squash phases together.

Every phase closeout must update `docs/handoffs/latest.md`, add
`docs/handoffs/phase-XX.md`, and record validation, blockers, protected-path
status, consumer tracker status, claim-boundary status, security/auth status,
data/migration status, master review, decision, and next checkpoint.

## If Runtime Limits Interrupt Work

Commit safe completed work first. Then update `docs/handoffs/latest.md` with an
exact resume point and include a resume prompt that starts at the next
unfinished checkpoint.
