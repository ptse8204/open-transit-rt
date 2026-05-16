# Community Support And Issue Triage Kit

Use this kit to report bugs, send release-candidate feedback, triage issues,
and prepare support-bundle context without exposing private data.

Support for Open Transit RT is community/best-effort unless a separate private
agreement exists outside this repository. This kit does not claim paid support,
SLA coverage, response-time guarantees, hosted service availability,
production readiness, compliance, adoption, agency approval, consumer
acceptance, vendor compatibility, hardware certification, final-root
readiness, or production-grade ETA quality.

## Triage Lanes

| Lane | Use when | Good public-safe input | Typical next step |
| --- | --- | --- | --- |
| Bug | Existing documented behavior fails | exact command, route, fixture, expected result, actual result, redacted output | Reproduce locally or request missing public-safe detail |
| Docs | A page is unclear, stale, or incomplete | file/page path, confusing text, suggested fix | Patch docs or label as good first issue |
| Feature | A scoped GTFS/GTFS-RT/admin/operator workflow is missing | operator need, smallest useful outcome, boundaries | Clarify scope and split if broad |
| Release feedback | `v0.1.0-rc.1` install, download, quickstart, or local app trial feedback | checked-out tag, command sequence, blocker, environment summary | Decide whether docs patch, current-source patch, or rc2 gate review is needed |
| Connector feedback | Synthetic or redacted connector path needs improvement | manifest name, fixture, adapter-conformance case, no-send behavior | Add fixture/docs/test or ask for redacted reproduction |
| Support-bundle review | A maintainer needs private diagnostics summarized safely | public-safe summary, manifest categories, redaction status | Keep raw bundle private; share only bounded facts |
| Security/private data | Vulnerability, leaked secret, private artifact, or unsafe evidence | do not post publicly | Use `SECURITY.md` or GitHub private security advisory |

## Reporter Checklist

Include:

- exact command, Make target, browser route, or file path;
- expected result and actual result;
- local environment summary without secrets;
- fixture path or synthetic input when possible;
- redacted error text;
- whether the issue affects install, GTFS import, public feeds, telemetry,
  Trip Updates, Alerts, validation, connectors, maintenance, release docs, or
  support-bundle workflow.

Do not include:

- tokens, JWTs, CSRF values, API keys, database URLs, private keys, ACME
  material, or private certificates;
- admin URLs with embedded secrets;
- private portal screenshots, private ticket links, private correspondence, or
  procurement details;
- raw logs with credentials, raw telemetry, raw GTFS, backup dumps, raw
  support bundles, or unredacted operator artifacts;
- adoption, compliance, production readiness, consumer acceptance, hosted
  service, SLA, vendor, hardware, final-root, or ETA-quality claims.

## Maintainer Triage Checklist

For each issue:

1. Confirm the issue is in project scope: GTFS, GTFS Studio, telemetry ingest,
   conservative matching, GTFS-RT feeds, Alerts, validation, monitoring,
   connectors, admin/operator workflows, release-candidate installability, or
   docs.
2. Check for obvious private data. If present, stop public discussion and
   follow the security/private-data process.
3. Ask for the smallest missing public-safe reproduction detail.
4. Apply labels that match the work type, such as `bug`, `documentation`,
   `good first issue`, `connector`, `readiness`, `agency-feedback`, or
   `security-review-needed`.
5. Route release-candidate feedback to docs patch, current-source patch, or
   rc2 gate review. Do not promise an rc2.
6. Close as out-of-scope when the request asks for rider apps, fare payments,
   passenger accounts, CAD/dispatch replacement, hosted SaaS promises, paid
   support, SLA/uptime guarantees, procurement/legal commitments, consumer
   submission automation, or unsupported public claims.

## Release-Candidate Feedback

Use the release-feedback issue template for `v0.1.0-rc.1` feedback about:

- fresh clone or tag checkout;
- published release page, checksum, archive, or download replay;
- `make check`, `make agency-app-up`, validators, `make validate`, or
  `make test`;
- local app startup and private Operations Console access;
- all five local public feed paths;
- docs mismatch between README, wiki, release notes, and quickstart.

Record exact blockers. Do not convert local success into production readiness,
hosted availability, compliance, consumer acceptance, or stable release
readiness.

## Support-Bundle Guidance

Support bundles are private diagnostics. They can help a maintainer understand
a bug, but the raw bundle should not be posted publicly.

Before sharing any support context:

- run only from an operator-controlled environment;
- review the manifest and excluded categories;
- remove private paths, DB URLs, tokens, cookies, raw telemetry, raw logs, raw
  validator reports, backup dumps, and private hostnames;
- summarize the narrow facts needed for the issue;
- keep the bundle out of `docs/evidence` unless a separate retained-evidence
  phase authorizes it.

## Public Reply Patterns

Use short, bounded replies:

- "Please add the exact command and redacted error text."
- "This looks like connector feedback; can you reproduce it with a synthetic
  fixture or manifest?"
- "This touches private deployment data. Please do not post the raw file
  publicly; summarize the relevant public-safe facts."
- "This request is outside project scope because it asks for hosted operations
  or SLA coverage."
- "This may affect rc2 gate review, but it does not by itself require or prove
  a new release candidate."

Do not promise response times, paid support, release dates, consumer outcomes,
or hosted operations.
