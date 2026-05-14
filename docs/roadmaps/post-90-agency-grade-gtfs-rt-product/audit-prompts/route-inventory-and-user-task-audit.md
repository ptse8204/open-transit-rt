# Audit Prompt — Route Inventory And User-Task Coherence

Audit route inventory against:

- `cmd/agency-config/main.go`
- `cmd/agency-config/operations.go`
- Operations Console navigation
- README/docs/wiki route maps
- `docs/phase-90-control-plane-final-status.md`

Required columns:

```text
route
area
html_or_json
role
methods
task
next_action_visible
claim_boundary
test_coverage
docs_coverage
gap
```

Check no public admin route, no raw private data leakage, no unsupported claims,
and proper auth/method/cache behavior.
