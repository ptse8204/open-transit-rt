# Master / Sub-Agent Operating Manual

## Required role flow

Use real sub-agents when available. If the environment cannot spawn them,
simulate them in clearly labeled sections.

| Role | Model | Used for |
| --- | --- | --- |
| Master Agent | GPT-5.5 x-high | All plans, checkpoint approval, phase closeout |
| Context / Repo Truth Sub-Agent | GPT-5.5 x-high | Current repo facts and source-of-truth alignment |
| Planning Sub-Agent | GPT-5.5 x-high | Phase plan and checkpoint sequence |
| Implementation Sub-Agent | GPT-5.5 high | Code/docs/test/script implementation |
| QA Sub-Agent | GPT-5.5 high | Tests, validation, route checks, install replay |
| UI/UX Sub-Agent | GPT-5.5 high | Task flow and usability review |
| Web Design Skill Sub-Agent | GPT-5.5 high + `web-design-engineer` | UX phases |
| Documentation / IA Sub-Agent | GPT-5.5 high | README/docs/wiki/site alignment |
| Claim-Boundary Sub-Agent | GPT-5.5 high | Unsupported-claim prevention |
| Security/Auth Sub-Agent | GPT-5.5 high | Auth, CSRF, credentials, release safety |
| Data/Migration Sub-Agent | GPT-5.5 high | Persistence, schema, migration impact |
| Release/Supply-Chain Sub-Agent | GPT-5.5 high | Package, tag, release, SBOM, checksums |
| Install Confidence Sub-Agent | GPT-5.5 high | Fresh clone and release replay |
| GTFS-RT Domain Sub-Agent | GPT-5.5 high | Realtime feed usefulness, fixtures, interoperability |

## Master approval

The Master Agent must approve every phase plan before implementation.
The Master Agent must approve every checkpoint after sub-agent review.

Do not stop to ask the maintainer between phases. The active `/goal` is the
continuation contract.

## Checkpoint report format

Every checkpoint report must include:

```text
Checkpoint:
Goal status:
Sub-agents used or simulated:
Changed files:
Validation run:
Blocked checks:
Protected path status:
Consumer tracker status:
Claim-boundary status:
Security/auth status:
Data/migration status:
Release/publication status:
Install confidence status:
Web design skill status:
Master review:
Required edits:
Decision:
Next checkpoint:
```

## Commit rule

Every checkpoint ends in a commit. Every phase ends in a closeout commit.

Use:

```text
Phase XXX -- Checkpoint 000001: add <phase> plan
Phase XXX -- Checkpoint 000002: implement or audit primary scoped work
Phase XXX -- Checkpoint 000003: run validation and patch required gaps
Phase XXX -- Checkpoint 000004: close <phase> review
```

Longer phases may add checkpoints. Do not mix phase scopes in one commit.
