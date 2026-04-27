# Phase 21 — Community, Governance, And Multi-Agency Scale

## Status

Complete for the approved docs/process/governance/teaching-visual scope.

## Purpose

Phase 21 prepares Open Transit RT for wider open-source collaboration and multi-agency use. It focuses on contributor trust, governance, maintainability, and agency-scale operations.

## Scope

1. Contribution and governance docs.
2. Issue templates and support boundaries.
3. Security disclosure and release process.
4. Multi-agency deployment strategy.
5. Maintainership and roadmap communication.
6. Teaching visuals for contribution, governance workflow, multi-agency strategy, evidence maturity, and support boundaries.

## Required Work

### 1) Contributor Experience

Add or improve:

- `CONTRIBUTING.md`;
- `CODE_OF_CONDUCT.md`;
- issue templates;
- PR checklist;
- coding conventions;
- test expectations;
- docs contribution rules.

### 2) Governance

Document:

- maintainer role;
- decision process;
- who can merge PRs;
- who can cut releases;
- who can approve docs/evidence wording;
- how competing design decisions are resolved;
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

### 5) Teaching Visuals

Add illustrative documentation graphics under `docs/assets/` for:

- contribution paths;
- community workflow;
- single-agency versus future multi-agency options;
- evidence maturity;
- support boundaries.

These visuals are teaching graphics, not screenshots or proof artifacts. They must not imply compliance, consumer acceptance, agency endorsement, hosted SaaS availability, paid support, SLA coverage, vendor equivalence, or universal production readiness.

## Acceptance Criteria

Phase 21 is complete only when:

- contributors have a clear path to help;
- agencies know how to report issues safely;
- security disclosure path exists;
- multi-agency deployment assumptions are documented;
- support boundaries are clear;
- project scope remains focused.
- teaching visuals are added, referenced with useful alt text, and documented in `docs/assets/README.md`.

## Required Checks

```bash
make validate
make test
git diff --check
make realtime-quality
make smoke
docker compose -f deploy/docker-compose.yml config
```

## Explicit Non-Goals

Phase 21 does not:

- create a legal foundation or company;
- promise paid support;
- imply official agency endorsement;
- expand into fares, payments, rider apps, or CAD/dispatch.
