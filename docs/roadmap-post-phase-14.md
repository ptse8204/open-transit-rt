# Post-Phase-14 Roadmap

This roadmap extends Open Transit RT beyond the current Phase 14 public-launch polish track. It is designed to give future Codex instances durable repo-native context so work can continue phase-by-phase without relying on chat history.

## Current State Assumption

- Phases 0 through 15 are closed for their documented scope.
- Phase 16 is the next planned phase unless `docs/handoffs/latest.md` says otherwise.
- Phase 12 produced hosted/operator evidence for the OCI pilot.
- Phase 13 created the consumer-submission evidence layer, but all consumer/aggregator records remain `not_started` until real third-party artifacts exist.

## Roadmap Principles

1. **Do not overclaim.** Do not claim CAL-ITP/Caltrans compliance, consumer acceptance, marketplace/vendor equivalence, or universal production readiness without evidence.
2. **Keep README welcoming.** The README should be the public front door. Phase history, evidence packets, implementation detail, and operations runbooks belong in `docs/`.
3. **Make agencies successful.** Future phases should reduce operator burden, not just add more backend features.
4. **Preserve architecture boundaries.** Open Transit RT remains the source of truth for GTFS management, telemetry, assignments, publication, validation records, workflow records, and audit state. Optional predictors stay behind `internal/prediction.Adapter`.
5. **Work one phase per Codex instance.** Each phase should update `docs/current-status.md`, `docs/handoffs/latest.md`, and a phase-specific handoff.
6. **Evidence beats assertions.** If a claim cannot be proven by code, deployment evidence, or third-party artifacts, record it as a future requirement.

## Phase Sequence

### Phase 14 — Public Launch Polish And Repo Simplification

Existing plan: `docs/phase-14-public-launch-polish.md`.

Goal: make the repository easier to understand publicly. Simplify README, improve docs navigation, add truthful teaching visuals, and improve support/star wording without changing runtime behavior.

### Phase 15 — Targeted Public Repo Hygiene And Evidence Redaction Review

Status: complete for the targeted delta-focused scope.

Goal: make the public repo safe to promote. Review files added or changed since the earlier scrub baseline, audit committed evidence, generated artifacts, zips, local files, and docs for accidental secrets or reconnaissance-heavy details.

Plan doc: `docs/phase-15-public-repo-security-hygiene.md`.

### Phase 16 — Agency Onboarding And Product Packaging

Goal: make the project usable by a non-expert small agency. Package a one-command product demo, full app Docker Compose, setup walkthrough, and first-use agency workflow.

Plan doc: `docs/phase-16-agency-onboarding-product-packaging.md`.

### Phase 17 — Deployment Automation And Pilot Operations

Goal: turn the OCI pilot evidence into a repeatable deployment model. Add deployment runbooks, production profiles, backup/restore automation, validator schedules, and operational playbooks.

Plan doc: `docs/phase-17-deployment-automation-pilot-operations.md`.

### Phase 18 — Admin UX And Agency Operations Console

Goal: reduce command-line dependence. Build a minimal web console for setup, feed health, validation, telemetry freshness, device credentials, alerts, and consumer evidence status.

Plan doc: `docs/phase-18-admin-ux-agency-operations-console.md`.

### Phase 19 — Realtime Quality And ETA Improvement

Goal: return to the hardest original product problem: trip matching and Trip Updates quality. Add replay evaluation, quality metrics, stronger diagnostics, and optional predictor contracts without losing conservative behavior.

Plan doc: `docs/phase-19-realtime-quality-eta-improvement.md`.

### Phase 20 — Consumer Submission Execution And CAL-ITP Readiness Program

Goal: move from evidence trackers to actual external submission work. Prepare packets, submit to consumers/aggregators, record outcomes, and produce a truthful readiness report.

Plan doc: `docs/phase-20-consumer-submission-calitp-readiness.md`.

### Phase 21 — Community, Governance, And Multi-Agency Scale

Goal: prepare the project for outside contributors and multiple agencies. Add contribution governance, issue templates, security policy, maintainership notes, and multi-agency deployment strategy.

Plan doc: `docs/phase-21-community-governance-multi-agency.md`.

## How Future Codex Instances Should Use This Roadmap

For each phase:

1. Read `AGENTS.md` first.
2. Read `docs/current-status.md` and `docs/handoffs/latest.md`.
3. Read this roadmap.
4. Read the phase-specific plan doc.
5. Execute only that phase.
6. Update status and handoff docs.
7. Run the required checks.
8. Do not mark the phase closed unless acceptance criteria are met.

## Important Non-Goals Across Future Phases

- Do not add rider apps, payments, fare collection, passenger accounts, or CAD/dispatch replacement unless a new roadmap explicitly changes scope.
- Do not claim consumer acceptance without artifacts from the named consumer.
- Do not claim CAL-ITP/Caltrans compliance from repo-only evidence.
- Do not let optional external predictors become the source of truth.
- Do not expose secrets or private operator evidence in public docs.
