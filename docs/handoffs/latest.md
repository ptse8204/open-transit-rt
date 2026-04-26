# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 15 — Targeted Public Repo Hygiene And Evidence Redaction Review is
complete for the documented targeted scope.

Phase 12 remains closed for the OCI pilot hosted/operator evidence scope. Phase
13 remains closed for the initial consumer-submission evidence tracker
structure. Phase 14 remains closed for the public launch polish and repo
simplification scope. Do not reopen those phases unless a blocking truthfulness
or safety issue directly affects the next task.

## Phase Status

- Phases 0 through 11 are closed for their documented scope.
- Phase 12 is closed for the OCI pilot hosted/operator evidence scope.
- Phase 13 is closed for the consumer-submission evidence tracker structure.
- Phase 14 simplified README/docs navigation, kept public reader guides in
  `wiki/`, made `docs/README.md` the documentation hub, and added reviewed
  teaching visuals.
- Phase 15 used `839efd6` (`Phase 14 -- Checkpoint 4 -- Security Cleanup`) as
  the earlier scrub baseline and completed a targeted delta-focused public repo
  hygiene review. It did not perform or claim a complete historical security
  audit.
- All seven current consumer/aggregator records are still `not_started`.
- No current repo evidence supports submitted, under-review, accepted,
  rejected, or blocked claims for any consumer target.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `README.md`
4. `wiki/README.md`
5. `docs/README.md`
6. `docs/handoffs/phase-15.md`
7. `docs/phase-15-public-repo-security-hygiene.md`
8. `SECURITY.md`
9. `docs/evidence/redaction-policy.md`
10. `docs/evidence/archive-inventory.md`
11. `docs/compliance-evidence-checklist.md`
12. `docs/consumer-submission-evidence.md`
13. `docs/evidence/consumer-submissions/README.md`
14. `docs/evidence/captured/oci-pilot/2026-04-24/README.md`
15. `docs/prompts/calitp-truthfulness.md`
16. `docs/tutorials/README.md`
17. `docs/tutorials/local-quickstart.md`
18. `docs/tutorials/agency-demo-flow.md`
19. `docs/tutorials/deploy-with-docker-compose.md`
20. `docs/tutorials/production-checklist.md`
21. `docs/tutorials/calitp-readiness-checklist.md`
22. `docs/assets/README.md`
23. `docs/dependencies.md`
24. `docs/decisions.md`

## Current Objective

Use the simplified README as the public front door and preserve the Phase 15
evidence safety boundary. Future docs work should keep public reader guides in
`wiki/`, detailed evidence/history in `docs/`, and committed captured evidence
aligned with `docs/evidence/redaction-policy.md`.

Do not claim CAL-ITP/Caltrans compliance, production readiness,
marketplace/vendor equivalence, agency endorsement, or consumer acceptance from
repo evidence, validator success, public fetch proof, workflow records, stars,
or the Phase 12 hosted packet alone.

## Exact First Commands

```bash
make validate
make test
git diff --check
```

If Docker is available and user-facing docs, README, or scripts change, also
run:

```bash
make smoke
make demo-agency-flow
```

## Current Evidence Boundary

- The OCI pilot packet at `docs/evidence/captured/oci-pilot/2026-04-24/`
  includes hosted/operator proof for public HTTPS feed fetches, TLS/redirect
  behavior, auth boundaries, clean hosted validation, monitoring/alert
  lifecycle, backup/restore, deployment rollback, and scorecard export job
  history.
- Phase 15 redacted unnecessary raw public client IP, remote port, and OCI
  instance-host details from committed OCI operator artifacts while preserving
  public host, feed URLs, validation status, TLS evidence, and safe proof
  artifacts.
- `docs/evidence/redaction-policy.md` is the evidence safety rule for future
  packets. Public evidence may include public URLs, validation status, TLS
  metadata, checksums, public headers/status, and redacted operational
  summaries. It must not include raw credentials, bearer tokens, admin URLs with
  secrets, private SSH paths, unredacted IP logs, private keys, database
  passwords, or internal hostnames unless explicitly justified as public-safe.
- `docs/evidence/archive-inventory.md` lists every committed archive under
  `docs/evidence/captured/**`; do not keep opaque archives without listed
  contents.
- The Phase 13 tracker at `docs/evidence/consumer-submissions/README.md` links
  to the Phase 12 packet as supporting evidence only.
- Validator success and public fetch proof are not consumer acceptance.
- Consumer-ingestion workflow records are not third-party acceptance.
- Acceptance may be claimed only for the named consumer, feed scope, URL root,
  and evidence date shown in a retained evidence artifact.

## Phase 15 Security Notes

- PATH scanner binaries `gitleaks` and `trufflehog` were unavailable.
- Docker was available and `docker run --rm -v "$PWD:/repo:ro"
  zricethezav/gitleaks:latest dir /repo --redact --verbose --no-banner` was
  attempted.
- The first Docker gitleaks directory scan found real secrets in ignored local
  `.cache/` files. Those local files were removed from the working tree.
- `git ls-files` and targeted `git log --all -- <paths>` checks found no tracked
  or historical git records for those `.cache` secret paths, so destructive git
  history cleanup is not indicated from this finding.
- The operator should rotate or revoke the affected admin tokens, device token,
  admin JWT secret, CSRF secret, device token pepper, ACME account key, and TLS
  private key before further pilot use.
- After local secret removal, the Docker gitleaks directory scan reported no
  leaks. Manual high-risk searches over tracked and non-cache working-tree files
  did not find committed private keys, cloud tokens, GitHub tokens, Slack tokens,
  OpenAI-style API keys, or literal Bearer credentials.

## First Files Likely To Edit

- `docs/phase-16-agency-onboarding-product-packaging.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-16.md`
- public onboarding/package docs selected by Phase 16

## Constraints To Preserve

- Keep README understandable by a non-technical agency reader in under 3
  minutes.
- Keep README concise and easy to scan, ideally under 150 to 200 lines unless
  examples genuinely require more.
- Keep claims evidence-bounded and truthful.
- Keep captions explicit about illustrative versus exact-behavior visuals.
- Keep alt text descriptive and useful.
- Do not change backend runtime behavior, API contracts, database schema, public
  feed URLs, external integrations, or consumer-submission status unless the
  active phase explicitly requires it.
- Do not add consumer submission APIs unless explicitly required and supported
  by a public documented API.
- Do not automate fake submissions.
- Do not invent acceptance, rejection, receipt, or blocker evidence.
- Do not perform destructive git history rewriting without explicit maintainer
  approval.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the
phase explicitly documents a reason to diverge.

## Future Roadmap

After Phase 15, use `docs/roadmap-post-phase-14.md` as the roadmap source of
truth.

The next planned phase after Phase 15 is:

- Phase 16 — Agency Onboarding Product Packaging

Future roadmap docs:

- `docs/phase-16-agency-onboarding-product-packaging.md`
- `docs/phase-17-deployment-automation-pilot-operations.md`
- `docs/phase-18-admin-ux-agency-operations-console.md`
- `docs/phase-19-realtime-quality-eta-improvement.md`
- `docs/phase-20-consumer-submission-calitp-readiness.md`
- `docs/phase-21-community-governance-multi-agency.md`
