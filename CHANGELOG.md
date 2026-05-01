# Changelog

All notable Open Transit RT release changes should be recorded here.

This project uses git tags as the release source of truth. Until maintainers
document a stronger compatibility policy, use semantic-version-style tags such
as `v0.22.0`.

## Unreleased

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
