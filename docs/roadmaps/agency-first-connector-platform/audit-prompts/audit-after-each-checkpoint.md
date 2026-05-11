# Audit Prompt After Each Checkpoint

Use after any checkpoint commit.

```text
Audit the just-completed checkpoint.

Confirm:
1. Commit message matches `Phase XX -- Checkpoint 00000N`.
2. Scope stayed inside the approved phase.
3. Protected evidence and consumer files were not changed.
4. Consumer tracker remains exactly seven targets and all `prepared`.
5. No unsupported compliance/adoption/consumer/vendor/production/ETA claims were added.
6. Required tests passed or environment blockers are documented.
7. UI changes are understandable to a small agency operator.
8. If gaps exist, patch inside the same phase using the next checkpoint number.

Return verdict:
- Pass
- Pass with small patch required
- Do not proceed; scope mismatch
```
