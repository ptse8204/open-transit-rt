# Public Release Policy

Phase 115 is the actual public release-candidate phase.

## Authorized if gates pass

```bash
git tag -a v0.1.0-rc.1 -m "Open Transit RT v0.1.0-rc.1"
git push origin v0.1.0-rc.1
gh release create v0.1.0-rc.1 \
  --title "Open Transit RT v0.1.0-rc.1" \
  --notes-file docs/release-notes-v0.1.0-rc.1-draft.md \
  <release-assets>
```

Use actual asset paths from the repo tooling. Attach source archive, checksum
manifest, SBOM/provenance files if generated and safe.

## Required gates before publishing

- clean worktree;
- release package generated from intended release commit;
- package audit passed;
- source archive public-distribution review passed;
- release notes bounded and truthful;
- claim audit passed;
- product acceptance audit passed;
- protected path check clean;
- consumer tracker prepared-only check passed;
- no secrets/private paths/protected evidence leak into public assets;
- no unsupported claims.

## If GitHub credentials or tools are unavailable

Do not fake publication.

Close Phase 115 as `blocked_publication_credentials` and continue the roadmap
using prepared-but-unpublished release status. Record exact command attempted,
tooling state, and what would unlock publication.
