# Audit Prompt — Public Release Blocking Audit

Use in Phases 112, 115, and 116.

Check:

- local source package generated from intended commit;
- package audit passed;
- source archive has no public-distribution blockers;
- release notes are bounded and truthful;
- no secrets or private paths in release metadata/assets;
- protected paths untouched;
- consumer tracker prepared-only;
- no unsupported claims;
- GitHub Release body says release candidate and local/self-hosted evaluation;
- release assets match checksums.

If a gate fails, do not publish.
