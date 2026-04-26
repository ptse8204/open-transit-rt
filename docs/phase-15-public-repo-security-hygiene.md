# Phase 15 — Public Repo Security Hygiene And Artifact Redaction

## Status

Planned phase. Not implemented until `docs/handoffs/latest.md` marks it active.

## Purpose

Phase 15 prepares the repository for broader public attention by reducing the risk that committed docs, evidence artifacts, generated files, or local files reveal secrets, private infrastructure detail, or unnecessary attack surface.

This phase is especially important because Phase 12 committed real hosted evidence artifacts. Those artifacts are valuable proof, but they must be reviewed carefully before public promotion.

## Scope

1. Public artifact and evidence audit.
2. Secret scanning and token hygiene.
3. Removal of accidental local files and stale bundles.
4. Security policy and disclosure guidance.
5. Evidence redaction rules for future operator packets.
6. GitHub-visible repo hygiene.

## Required Work

### 1) Public Artifact Audit

Inspect:

- `docs/evidence/captured/**`
- `docs/assets/**`
- root-level archives such as `*.zip`
- `.env*` examples
- scripts and generated logs
- hidden local files such as `.DS_Store`

Look for:

- tokens, bearer credentials, JWTs, API keys, DuckDNS tokens, DB passwords, SSH keys, private certs, private IPs paired with credentials, unredacted emails that should remain private, and internal hostnames that are not necessary for public evidence.

### 2) Secret Scanning

Run at least one local secret scanner if available, such as:

```bash
gitleaks detect --source . --redact --verbose
trufflehog git file://. --only-verified
```

If a scanner is unavailable, document that blocker and perform a manual high-risk pattern search.

### 3) Remove Accidental Files

- Remove tracked `.DS_Store` and add it to `.gitignore`.
- Inspect any root-level zip bundles. Remove them unless they are intentional, current, and safe.
- Avoid committing generated local runtime logs unless they are deliberate redacted evidence.

### 4) Evidence Redaction Policy

Add or update a doc such as `docs/evidence/redaction-policy.md` explaining:

- what evidence may be public;
- what must stay private;
- what must be redacted;
- how to store secrets out of the repo;
- how to summarize private operator artifacts safely.

### 5) Security Policy

Add or update:

- `SECURITY.md`
- optional `.github/ISSUE_TEMPLATE/` docs if appropriate
- disclosure instructions for security issues

### 6) Rotate If Needed

If a real secret is found, do not only delete it. Document that the operator must rotate/revoke the secret because Git history and forks may retain it.

## Acceptance Criteria

Phase 15 is complete only when:

- secret scan or manual equivalent is recorded;
- risky accidental files are removed or justified;
- evidence artifacts are reviewed for redaction;
- `.gitignore` covers common local artifacts;
- `SECURITY.md` exists or is updated;
- evidence redaction guidance exists;
- any discovered secret has a documented rotation action;
- README remains friendly and does not expose sensitive details;
- status and handoff docs are updated truthfully.

## Required Checks

```bash
make validate
make test
git diff --check
```

Run `make smoke` and `make demo-agency-flow` if docs or scripts that affect user-facing flows are changed.

## Explicit Non-Goals

Phase 15 does not:

- add backend product features;
- change public feed URLs;
- delete useful evidence merely because it is operational;
- hide evidence that is safe and necessary for trust;
- claim compliance or consumer acceptance.
