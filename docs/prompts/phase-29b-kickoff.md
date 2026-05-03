# Prompt — Start Phase 29B AVL / Vendor Adapter Pilot Implementation

```text
Read and obey AGENTS.md first.

Then read:
1. docs/current-status.md
2. docs/handoffs/latest.md
3. docs/handoffs/phase-29a.md
4. docs/phase-29b-avl-vendor-adapter-pilot.md
5. docs/tutorials/device-avl-integration.md
6. docs/tutorials/device-token-lifecycle.md
7. docs/evidence/device-avl/README.md
8. docs/evidence/redaction-policy.md
9. SECURITY.md
10. docs/dependencies.md
11. docs/decisions.md

You are starting Phase 29B — AVL / Vendor Adapter Pilot Implementation only.

Goal:
Create or document a minimal vendor/AVL adapter pilot pattern using synthetic vendor payloads, dry-run behavior, and the existing POST /v1/telemetry boundary.

Do not certify vendors or commit real vendor payloads/credentials.

Focus on:
- generic vendor payload to Open Transit RT telemetry contract
- synthetic vendor payload fixtures
- dry-run transform behavior if safe
- tests for synthetic payload transformations
- credential/redaction boundaries
- docs/handoffs/phase-29b.md

Run:
- make validate
- make test
- make realtime-quality
- make smoke
- make test-integration
- docker compose -f deploy/docker-compose.yml config
- git diff --check

If scripts are added:
- sh -n <script>
- <script> help
- <script> --dry-run <synthetic fixture>
```
