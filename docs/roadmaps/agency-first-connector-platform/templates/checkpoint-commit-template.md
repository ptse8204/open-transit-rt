# Checkpoint Commit Template

```text
Phase <XX> -- Checkpoint <00000N>: <short outcome>
```

Commit body template:

```text
Scope:
- <what changed>

Validation:
- <commands run>

Boundaries:
- no evidence writes
- no consumer status changes
- no unsupported claims
- protected files unchanged
```
