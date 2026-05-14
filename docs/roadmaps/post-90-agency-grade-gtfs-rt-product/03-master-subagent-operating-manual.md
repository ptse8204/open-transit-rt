# Master / Sub-Agent Operating Manual

## Model assignments

| Role | Model level |
| --- | --- |
| Master Agent | GPT-5.5 x-high |
| Context / Repo Truth Sub-Agent | GPT-5.5 x-high |
| Planning Sub-Agent | GPT-5.5 x-high |
| Implementation Sub-Agent | GPT-5.5 high |
| QA Sub-Agent | GPT-5.5 high |
| UI/UX Sub-Agent | GPT-5.5 high |
| Documentation / IA Sub-Agent | GPT-5.5 high |
| Claim-Boundary Sub-Agent | GPT-5.5 high |
| Security/Auth Sub-Agent | GPT-5.5 high |
| Data/Migration Sub-Agent | GPT-5.5 high when persistence is touched |
| Release Sub-Agent | GPT-5.5 high when release tooling is touched |

## Workflow

1. Context / Repo Truth reads source-of-truth docs and reports current state.
2. Planning drafts a checkpoint-level plan.
3. Security/Auth, Documentation/IA, UI/UX, and Claim-Boundary review the plan.
4. Master approves before implementation starts.
5. Implementation executes only approved scope.
6. QA runs validation and records blockers.
7. Master closes only when all required edits are resolved.

## Checkpoint report

```text
Checkpoint:
Sub-agents used or simulated:
Changed files:
Validation run:
Blocked checks:
Protected path status:
Consumer tracker status:
Claim-boundary status:
Security/auth status:
Data/migration status:
UI/UX status:
Documentation/IA status:
Master review:
Required edits:
Decision:
Next checkpoint:
```
