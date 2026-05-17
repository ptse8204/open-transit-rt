# Branching And Release Policy

Open Transit RT currently publishes `v0.1.0-rc.2` as a public release candidate
for local/self-hosted evaluation. It is not a stable release and does not prove
production readiness, CAL-ITP/Caltrans compliance, consumer acceptance, hosted
service availability, vendor compatibility, SLA/uptime, production AVL
reliability, or ETA quality.

## Branches

`main` is the active development branch. It may contain AI-agent handoffs,
phase ledgers, roadmap packs, prompt files, and other maintenance history.

`stable` is the filtered product branch. It should contain product/user-facing
source, examples, tests, release docs, website files, and practical human
documentation. It should not contain AI-agent-only docs.

`gh-pages` is the static public site branch.

## Stable Branch Baseline

The stable branch is initialized from the `v0.1.0-rc.2` release-candidate
commit because that is the latest published release-candidate baseline:

```bash
git branch stable v0.1.0-rc.2
```

The initial branch is a filtered branch, not a new release tag. It must not be
described as production-ready or compliant.

## Stable Sync Automation

`.github/workflows/update-stable.yml` filters `main` into `stable`.

The workflow:

- runs automatically after pushes to `main`;
- can be run manually with `workflow_dispatch`;
- has a `dry_run` input for previewing the filtered update;
- uses `.github/stable-sync-excludes.txt` to remove AI-agent-only docs;
- commits only when the filtered tree changes;
- pushes with a normal `git push origin HEAD:stable`;
- does not force-push.

If the remote `stable` branch has diverged, the push should fail instead of
rewriting history. A maintainer must then inspect the branch and decide how to
reconcile it.

## Stable Branch Exclusions

The stable branch excludes:

- `docs/agent/**`;
- `docs/handoffs/**`;
- `docs/prompts/**`;
- `docs/roadmaps/**`;
- root `docs/phase-*.md` phase ledgers;
- Codex task, conversation summary, master planner, and roadmap-status files;
- ignored local cache output.

The stable branch preserves:

- `README.md`;
- human docs and tutorials;
- release docs;
- website files under `site/`;
- source code;
- examples;
- tests;
- fixtures;
- GitHub workflows.

## Claim And Evidence Boundaries

Stable branch automation is not an evidence workflow. It does not contact
agencies, consumers, vendors, portals, validators, or external endpoints. It
does not create evidence packets and does not change
`docs/evidence/consumer-submissions/status.json`.

Consumer tracker status must remain seven prepared-only targets unless future
retained, redacted, target-originated evidence supports a specific change.

Do not use the stable branch to claim compliance, production readiness,
consumer acceptance, agency adoption, final-root readiness, hosted service
availability, vendor compatibility, hardware certification, SLA/uptime,
production AVL reliability, production-grade ETA quality, or real-world ETA
accuracy.

## Release Relationship

Releases are still cut from reviewed commits and tags. Creating or updating
`stable` does not publish a release, tag a release, or change the status of
`v0.1.0-rc.2`.
