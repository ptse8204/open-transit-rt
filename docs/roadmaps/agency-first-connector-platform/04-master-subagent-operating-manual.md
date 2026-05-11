# Master / Sub-Agent Operating Manual

## Roles

### Master Agent

Owns phase scope, claim boundaries, and final approval. It does not let implementation begin until the phase plan is specific.

### Planning Sub-Agent

Creates a detailed phase plan:

- files to read;
- files to edit;
- files not to edit;
- tests to run;
- commit sequence;
- stop conditions.

### Implementation Sub-Agent

Executes only the approved plan. It must not expand scope.

### QA Sub-Agent

Runs validation commands, records environment blockers, and checks protected files.

### UI/UX Sub-Agent

Checks whether a small agency operator could understand the workflow without heavy command-line or phase-history knowledge.

### Claim-Boundary Sub-Agent

Searches for unsupported claims and blocks changes to evidence/consumer status files.

## Workflow Per Phase

1. Read source-of-truth docs.
2. Planning sub-agent drafts plan.
3. Master reviews plan and tightens it.
4. Implementation sub-agent executes first checkpoint.
5. QA validates.
6. UI/UX reviews usability.
7. Claim-boundary reviews claims and protected paths.
8. Master decides pass/fail.
9. If pass, close phase with closeout checkpoint.
10. If fail, patch inside the same phase using next checkpoint number.

## Commit Sequence

Every phase starts at checkpoint `000001`.

Example:

```text
Phase 62 -- Checkpoint 000001: add guided setup and browser GTFS import plan
Phase 62 -- Checkpoint 000002: implement guided setup wizard v1
Phase 62 -- Checkpoint 000003: implement browser GTFS import and validation flow
Phase 62 -- Checkpoint 000004: close guided setup and browser GTFS import
```

Patch after close:

```text
Phase 62 -- Checkpoint 000005: fix browser GTFS import audit gaps
```

## Progress Updates

After each checkpoint, report:

- checkpoint name;
- changed files;
- validations run;
- blockers;
- protected path status;
- claim-boundary status;
- next checkpoint.
