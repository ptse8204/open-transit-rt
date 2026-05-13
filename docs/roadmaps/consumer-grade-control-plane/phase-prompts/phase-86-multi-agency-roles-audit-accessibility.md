# Phase 86 Prompt — Multi-Agency, Roles, Audit, And Accessibility

## Goal

Make tenant/role visibility safer and improve accessibility for the control plane.

## Scope

- Agency switcher/scope visibility improvements.
- Role and permission explanation UI.
- Access-denied UX.
- Audit log browser for admin operations if existing audit records support it; otherwise add docs/model only.
- Accessibility review: keyboard navigation, focus order, labels, table captions, skip links, reduced motion, responsive behavior.

## Do not

- Claim production multi-tenant hosting.
- Add global admin powers without explicit design review.
- Leak cross-agency data.
- Add migrations unless required and separately reviewed.
- Write evidence.

## Validation

Baseline validation plus:

```bash
make validate
make test
RUN_LOCAL_APP=true make release-candidate-check
```

## Commits

```text
Phase 86 -- Checkpoint 000001: add multi-agency roles audit accessibility plan
Phase 86 -- Checkpoint 000002: add agency scope and switcher improvements
Phase 86 -- Checkpoint 000003: add role permission and access-denied UX
Phase 86 -- Checkpoint 000004: add audit log browser or scoped audit model docs
Phase 86 -- Checkpoint 000005: run accessibility and keyboard navigation review
Phase 86 -- Checkpoint 000006: close multi-agency roles audit accessibility review
```
