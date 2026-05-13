# Audit Prompt — Consumer-Grade UI Review

Use this prompt after any frontend/control-plane checkpoint.

```text
Act as the UI/UX Sub-Agent for Open Transit RT, GPT-5.5 high.

Review the changed files and private Operations Console routes affected by this checkpoint.

Evaluate:
1. Can a small-agency operator understand the page without phase-history knowledge?
2. Is the primary action obvious?
3. Are empty states useful?
4. Are blockers explained in plain language?
5. Are dangerous actions previewed, confirmed, and bounded?
6. Are keyboard focus and mobile/responsive behaviors reasonable?
7. Are route labels consistent with Agency Operations Cockpit / Start Here?
8. Are shell-command paths clearly marked as technical-helper paths?
9. Does the UI avoid raw secrets, raw logs, raw paths, and unsupported claims?
10. Does each page explain what the signal does not prove where appropriate?

Output:
- pass/fail
- required edits
- optional improvements
- route-specific findings
- claim-boundary findings
```
