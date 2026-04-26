# Phase 15 Kickoff Prompt

Paste this into a fresh Codex instance after Phase 14 is complete.

```text
Read and obey AGENTS.md first.

Then read, in order:
1. docs/current-status.md
2. docs/handoffs/latest.md
3. docs/roadmap-post-phase-14.md
4. docs/phase-15-public-repo-security-hygiene.md
5. docs/compliance-evidence-checklist.md
6. docs/evidence/README.md
7. docs/evidence/captured/oci-pilot/2026-04-24/README.md
8. docs/prompts/calitp-truthfulness.md
9. README.md
10. docs/dependencies.md

You are starting Phase 15 — Public Repo Security Hygiene And Artifact Redaction only.

Phase 14 must be closed before this phase starts. If Phase 14 is not closed, stop and update the handoff rather than doing Phase 15 work.

Goal:
Make the public repository safe to promote by auditing committed evidence, artifacts, generated files, docs, and root-level files for secrets, accidental local files, and unnecessary sensitive operational detail.

In scope:
1. Secret scanning or manual equivalent.
2. Evidence artifact redaction review.
3. Remove tracked .DS_Store and other accidental local files.
4. Inspect root-level zip/archive files and remove or justify them.
5. Add or update .gitignore for common local/generated artifacts.
6. Add or update SECURITY.md.
7. Add docs/evidence/redaction-policy.md.
8. Update docs/current-status.md and docs/handoffs/latest.md.
9. Add docs/handoffs/phase-15.md using the repo handoff template.

Out of scope:
- backend product features;
- runtime behavior changes;
- public feed URL changes;
- hiding safe and useful evidence;
- claiming compliance or consumer acceptance.

Before editing, run and record:
- make validate
- make test
- git diff --check

If available, run and record at least one of:
- gitleaks detect --source . --redact --verbose
- trufflehog git file://. --only-verified

If these tools are unavailable, document the blocker and run a manual high-risk search for token/secret/key/password/private-key/env patterns.

After editing, run and record:
- make validate
- make test
- git diff --check

Run make smoke and make demo-agency-flow if README, scripts, or user-facing docs are changed.

Acceptance criteria:
- public evidence artifacts are reviewed;
- accidental local files are removed or justified;
- redaction policy exists;
- security disclosure path exists;
- any discovered secret has a documented rotate/revoke requirement;
- no unsupported claims are introduced;
- Phase 15 handoff is accurate.
```
