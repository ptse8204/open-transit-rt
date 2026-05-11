# Audit Prompt After Each Phase

Use after a phase closeout checkpoint.

```text
Audit the completed phase.

Confirm:
1. The phase has a plan/implementation/closeout sequence or equivalent approved checkpoint history.
2. Source-of-truth docs mention the closed phase accurately.
3. The next phase is identified but not started without maintainer approval.
4. Protected paths are unchanged.
5. Claim boundaries are safe.
6. Validation commands passed or environment blockers are precise.
7. UI/UX is materially easier for agencies or integrators.

If accepted, recommend the next phase prompt from `phase-prompts/`.
If not accepted, propose one patch checkpoint in the same phase.
```
