# Audit Prompt — Release Candidate Gate

Use only for release-readiness phases.

Run standard checks. Only run package creation/audit when explicitly authorized:

```bash
make release-package
make audit-release-package
```

Final decision must be one of:

```text
ready_for_tag_candidate
needs_review
blocked
```
