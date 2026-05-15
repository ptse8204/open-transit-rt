# Contributor First Issues

Use this guide to choose small, safe first contributions to Open Transit RT.
It complements `CONTRIBUTING.md`.

## Good First Issue Shapes

Prefer one of these:

- Fix one broken command, link, or prerequisite note in a tutorial.
- Improve one troubleshooting entry for Docker, Java, validators, ports, or
  local app startup.
- Add one synthetic fixture with no credentials, real vendor payload, private
  agency data, or private endpoint.
- Add one focused test for existing documented behavior.
- Improve connector examples or manifests using synthetic inputs only.
- Clarify one claim boundary without making the underlying claim stronger.
- Improve wiki or docs navigation for one task path.

## Avoid For A First PR

Avoid first PRs that touch:

- database migrations;
- public feed URL contracts;
- GTFS-RT protobuf semantics;
- auth, token, or secret handling;
- protected evidence paths;
- `docs/evidence/consumer-submissions/status.json`;
- real agency, vendor, device, or portal records;
- release tags, packages, images, or GitHub Releases;
- broad public wording about compliance, launch, adoption, consumers,
  production readiness, hosted service, vendor compatibility, hardware, SLA, or
  ETA quality.

## First PR Flow

1. Read `AGENTS.md`, `README.md`, and `CONTRIBUTING.md`.
2. Pick one small scope.
3. Search for nearby tests or docs examples.
4. Make the smallest change that fixes the issue.
5. Run:

   ```bash
   git diff --check
   make check
   ```

6. If code changes, also run:

   ```bash
   make test
   ```

7. If connector examples change, also run:

   ```bash
   make external-connection-check
   make adapter-conformance
   make test-connector-examples
   ```

8. If readiness, evidence, release, or public-claim wording changes, also run:

   ```bash
   make audit-product-acceptance
   make audit-final-claim-review
   ```

9. In the PR, state what changed, what checks ran, and whether any evidence or
   consumer status was touched. For first PRs, the answer should normally be no.

## Safe Fixture Checklist

Before committing a fixture:

- Use synthetic IDs such as `demo-agency`, `vehicle-1`, and `device-1`.
- Keep timestamps deterministic.
- Keep files small and documented.
- Do not include device tokens, admin tokens, JWTs, CSRF secrets, database URLs,
  private hostnames, portal links, private operator names, raw private payloads,
  or screenshots.
- Add a README if the fixture directory is not self-explanatory.

## When To Ask For Maintainer Review First

Ask before work that changes:

- feed contracts;
- migrations;
- auth or role behavior;
- prediction adapter semantics;
- public claim boundaries;
- evidence handling;
- release process;
- connector manifest schema;
- data retention or redaction behavior.
