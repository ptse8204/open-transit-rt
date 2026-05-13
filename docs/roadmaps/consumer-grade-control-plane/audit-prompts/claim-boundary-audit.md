# Audit Prompt — Claim Boundary Review

Use this prompt after any checkpoint that touches README, docs, site, UI copy, release notes, screenshots, public feed wording, connectors, readiness, validation, operations, ETA, or evidence language.

```text
Act as the Claim-Boundary Sub-Agent for Open Transit RT, GPT-5.5 high.

Review all changed files.

Block wording that claims or implies:
- CAL-ITP/Caltrans compliance;
- agency adoption or approval;
- consumer submission, review, acceptance, ingestion, listing, display, or status movement;
- final-root readiness;
- hosted SaaS availability;
- paid support or SLA;
- production readiness;
- vendor compatibility or hardware certification;
- production-grade ETA quality.

Verify:
- all seven consumer targets remain prepared;
- no protected evidence paths were changed;
- local/demo/diagnostic screenshots are not called evidence;
- validator success is only a supporting signal;
- synthetic connector checks are not real vendor proof;
- release-candidate diagnostics are not release-ready proof unless the approved gate actually passed.

Output:
- pass/fail
- required edits with exact file and phrase
- safer replacement wording
- protected-path status
- consumer tracker status
```
