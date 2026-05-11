# Changelog

All notable Open Transit RT release changes should be recorded here.

This project uses git tags as the release source of truth. Until maintainers
document a stronger compatibility policy, use semantic-version-style tags such
as `v0.1.0-rc.1` for the first release-candidate gate and `v0.1.0` only after
that gate is reviewed.

## Unreleased

- Added the canonical post-60 review/recommendations section that prioritizes
  product quality, `v0.1.0-rc.1` readiness, and external-connection maturity
  while keeping real pilots, final-root proof, consumer submission, and vendor
  proof optional authorization-gated evidence tracks.
- Added product screenshot and diagram documentation manifests for local/demo
  visual aids that are not retained evidence or claim proof.
- Added docs-first production operations hardening guidance, including operations cadence, incident/response templates, alert delivery proof, capacity checks, secret rotation, operator handover, backup/restore hardening, validator failure response, and evidence refresh boundaries.
- Added release and distribution hardening documentation, including release checklist, release notes template, upgrade/rollback guide, version pinning guidance, and evidence version-linkage guidance.

## Release Note Rules

Each release entry should include:

- user-facing changes;
- migrations, or `None`;
- operations changes, or `None`;
- security notes, or `None`;
- dependency changes, or `None`;
- evidence or claim changes, or `None`;
- known limitations;
- checks run and blocked checks.

Do not use changelog entries to claim CAL-ITP/Caltrans compliance, consumer
submission/review/acceptance, hosted SaaS availability, agency endorsement,
paid support, SLA coverage, marketplace/vendor equivalence, or production-grade
ETA quality unless retained evidence supports that exact claim.
