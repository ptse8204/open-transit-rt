# Phase 21 — Community, Governance, And Multi-Agency Scale

## Status

Planned phase. Not implemented until `docs/handoffs/latest.md` marks it active.

## Purpose

Phase 21 prepares Open Transit RT for wider open-source collaboration and multi-agency use. It focuses on contributor trust, governance, maintainability, and agency-scale operations.

## Scope

1. Contribution and governance docs.
2. Issue templates and support boundaries.
3. Security disclosure and release process.
4. Multi-agency deployment strategy.
5. Maintainership and roadmap communication.

## Required Work

### 1) Contributor Experience

Add or improve:

- `CONTRIBUTING.md`;
- issue templates;
- PR checklist;
- coding conventions;
- test expectations;
- docs contribution rules.

### 2) Governance

Document:

- maintainer role;
- decision process;
- release process;
- how agencies can request features;
- what is out of scope.

### 3) Multi-Agency Strategy

Review and document:

- single-agency deployment model;
- multi-agency deployment options;
- agency-scoped auth boundaries;
- data isolation expectations;
- hosted service considerations.

### 4) Public Communication

Add a clear roadmap status page that avoids overclaiming and helps contributors understand what matters next.

## Acceptance Criteria

Phase 21 is complete only when:

- contributors have a clear path to help;
- agencies know how to report issues safely;
- security disclosure path exists;
- multi-agency deployment assumptions are documented;
- support boundaries are clear;
- project scope remains focused.

## Required Checks

```bash
make validate
make test
git diff --check
```

## Explicit Non-Goals

Phase 21 does not:

- create a legal foundation or company;
- promise paid support;
- imply official agency endorsement;
- expand into fares, payments, rider apps, or CAD/dispatch.
