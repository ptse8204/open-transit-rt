# Extension Governance

Open Transit RT extensions must keep the core product stable, auditable, and
safe for small agencies. This policy covers connector manifests, sidecars,
public API compatibility, deprecation, security review, release train planning,
and the post-110 roadmap.

This policy does not introduce dynamic backend plugin loading. It does not
certify vendors, hardware, deployments, consumer acceptance, compliance,
production readiness, hosted service availability, SLA/uptime, or ETA quality.

## Extension Model

Supported extension shapes:

| Shape | Boundary | Default Posture |
| --- | --- | --- |
| Telemetry sidecar | Sends normalized observations to authenticated `POST /v1/telemetry`. | Disabled until configured by an operator. |
| Prediction sidecar | Runs behind `internal/prediction.Adapter` for Trip Updates. | Fail closed; Vehicle Positions remain independent. |
| Validator wrapper | Uses allowlisted validator IDs or offline fixtures. | No raw commands in manifests. |
| Monitoring export | Sends redacted summaries to operator-owned monitoring. | No-send by default. |
| Consumer discovery helper | Prepares metadata or URLs for review. | No portal automation and no tracker mutation. |
| Documentation recipe | Explains a synthetic/local setup or integration pattern. | Must be redaction-first and claim-bounded. |

Unsupported extension shapes:

- dynamically loaded Go plugins;
- arbitrary shell command plugins;
- browser/portal automation that submits to consumers;
- code that mutates protected evidence paths or consumer tracker state;
- default notification sending;
- direct database writes from extensions;
- public feed contract changes hidden inside connector examples;
- vendor/device compatibility claims without a separate authorized evidence
  gate.

## Connector Manifest Compatibility

The current connector schema is `open-transit-rt.connector.v1`.

Compatibility policy:

- `v1` manifests should remain additive where possible.
- New optional fields may be added only when old manifests continue to decode.
- New required fields require either a schema version bump or a documented
  migration path.
- New connector types require docs, example manifests, synthetic fixtures, and
  conformance cases.
- Existing safety fields must not be weakened: disabled-by-default,
  fail-closed behavior, redaction policy, claim boundary, docs link, and
  conformance cases remain part of the review contract.
- Unsupported positive claims remain rejected unless maintainers explicitly add
  a narrow safe phrase and update claim audits.

Manifest review checklist:

- schema version and connector ID are valid;
- connector type is supported;
- mode is disabled by default;
- no status mutation, consumer submission automation, notification send by
  default, raw validator command, unsafe URL, private path, or secret value is
  present;
- input/output contracts are narrow and documented;
- failure behavior is bounded and fail-closed;
- redaction policy lists fields that may expose secrets or private payloads;
- claim boundary avoids vendor, hardware, compliance, consumer, production,
  hosted-service, SLA/uptime, release, or ETA-quality claims;
- conformance cases use synthetic or public-safe fixtures.

## Public API Stability

Public contracts are narrower than internal packages. Treat these as stable
operator-facing boundaries once documented:

- public GTFS and GTFS-RT feed URLs;
- `feeds.json` structure;
- authenticated telemetry ingest payloads;
- connector manifest schema;
- prediction adapter input/output boundary;
- admin command result shape for documented command routes;
- documented CLI and Makefile targets used by operators.

Compatibility rules:

- Prefer additive changes for public feed metadata, telemetry ingest payloads,
  manifest fields, command result JSON, and operator CLI output.
- Do not remove or rename public feed paths without a documented deprecation
  window.
- Do not couple Trip Updates predictor behavior into telemetry ingest or
  Vehicle Positions publication.
- Do not require a heavy frontend stack or dynamic plugin runtime to preserve
  existing operator workflows.
- Internal package names and unexported implementation details are not public
  API, but changing them should preserve documented service boundaries.

## Deprecation Policy

Deprecations should be explicit and reversible until the removal release.

Deprecation steps:

1. Document the replacement path and risk.
2. Add tests or audits that confirm both old and new paths during the
   transition.
3. Mark docs, manifests, flags, or CLI output as deprecated without removing
   them immediately.
4. Keep at least one release-candidate review window before removal when a
   public or operator-facing contract is affected.
5. Record the removal in release notes and upgrade guidance.

Do not deprecate by silently changing behavior, changing public feed URLs,
weakening validation, hiding low-confidence matching, or converting optional
evidence tracks into assumed claims.

## Security Review Process

Extension and API changes need security review when they involve:

- credentials, tokens, JWTs, device tokens, database URLs, webhooks, TLS, or
  private keys;
- telemetry ingest, device binding, agency scope, role checks, CSRF, or admin
  commands;
- external HTTP calls, polling, sidecars, validator wrappers, or monitoring
  exports;
- public feed exposure, final-root guidance, consumer workflows, or evidence
  artifacts;
- raw logs, private payloads, screenshots, correspondence, or real
  agency/vendor/device data.

Security review checklist:

- secrets are env references only and never committed;
- request bodies and responses stay bounded and redacted;
- browser POST routes preserve CSRF expectations;
- sidecars fail closed and cannot mutate internal status directly;
- connector examples use synthetic fixtures;
- public docs avoid stronger claims;
- protected paths and consumer tracker files remain unchanged unless a
  separately authorized evidence/status workflow permits them;
- incident or secret-response guidance in `SECURITY.md` and
  `docs/evidence/redaction-policy.md` is followed.

## Maintainer Release Train Proposal

Recommended cadence:

- Patch review: narrowly scoped bug fixes, docs corrections, claim-boundary
  fixes, or security hardening.
- Release-candidate review: grouped operator-facing changes after local
  validation and release-candidate diagnostics.
- Minor review: additive product capabilities, connector types, public API
  additions, or workflow improvements.
- Major review: incompatible manifest/API changes or service-boundary changes
  after an explicit decision record.

Release train gates:

- source tree is clean;
- `make check`, claim audits, and prepared-only tracker checks pass;
- required code/test/connector/validator checks are run for changed areas;
- release notes list migrations, operations changes, security notes,
  dependency changes, limitations, blocked checks, and claim boundaries;
- local package or release-candidate diagnostics remain local until a
  maintainer explicitly cuts a release;
- no tag, GitHub Release, package publication, or image publication occurs
  without maintainer action.

## Governance Review Levels

| Change | Review Level |
| --- | --- |
| Typo/docs link fix with no claim change | One maintainer or normal PR review. |
| Synthetic connector example or fixture | Connector review plus conformance checks. |
| New connector type or manifest field | Architecture review, docs update, tests, and conformance cases. |
| Public feed or telemetry ingest contract change | Architecture decision, migration/compatibility plan, and full validation. |
| Security/auth/data-integrity change | Security review and focused tests. |
| Evidence or consumer status wording | Claim-boundary review and protected-path/status checks. |
| Release publication | Maintainer release process only. |

## Post-110 Roadmap

Recommended future work remains separated by gate:

| Track | Next Safe Work | Boundary |
| --- | --- | --- |
| Release review | Maintain `v0.1.0-rc.1` diagnostics and blocker matrix from a clean checkout. | No tag or publication without maintainer release action. |
| Connector ecosystem | Add synthetic/local recipes, manifest linting, and adapter conformance cases. | No vendor compatibility or hardware claims without evidence authorization. |
| Operator usability | Keep improving private Operations Console IA, docs, and training. | No hosted-service or production-readiness claims. |
| Static GTFS workflows | Continue draft/published separation, diff, rollback, and safe fix planning. | No silent production edits. |
| Realtime quality | Expand synthetic backtests, withheld-case metrics, and fail-closed predictor review. | No real-world ETA-quality claim without an authorized study. |
| Evidence gates | Use `docs/future-evidence-intake-gate-pack.md`. | No complete intake, no evidence work. |
| API compatibility | Add decision records before incompatible public contract changes. | No hidden breaking changes. |

## Required Validation For Extension PRs

Use the smallest relevant set:

```bash
git diff --check
make check
make audit-product-acceptance
make audit-final-claim-review
make external-connection-check
make adapter-conformance
make test-connector-examples
```

If Go code, runtime behavior, routes, migrations, public feed behavior, build
behavior, examples, or tests changed, also run:

```bash
make validate
make test
docker compose -f deploy/docker-compose.yml config
```

Every extension PR should state whether protected paths changed, whether the
consumer tracker changed, what checks ran, and what the change does not prove.
