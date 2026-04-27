# Support Boundaries

Open Transit RT is independent open-source software. Maintainer help is community support, not paid support, SLA coverage, hosted SaaS availability, agency endorsement, or universal production-readiness assurance.

![Illustrative support boundary diagram showing maintainer help, operator-owned responsibilities, and community-only support.](assets/support-boundaries.png)

## Maintainers Can Help With

Maintainers can reasonably help with:

- reproducible bugs in the repository;
- focused PR review;
- docs corrections;
- validator or fixture findings;
- questions about supported Make targets;
- redaction-safe evidence structure;
- architecture and scope clarification.

Useful reports include exact commands, local fixture names, endpoint paths, expected behavior, actual behavior, and redacted logs.

## Operators Must Own

Operators and deployments own:

- runtime secrets and token rotation;
- DNS, TLS, reverse proxies, and hosting accounts;
- database operations, backups, and restores;
- monitoring, alerting, and incident response;
- production validation cadence;
- agency identity, license, and contact metadata;
- consumer or aggregator submissions;
- private logs, private portals, private tickets, and correspondence.

Do not post private deployment artifacts in public issues.

## Community-Only Support

Community discussion can help with:

- examples and patterns;
- peer notes from local evaluation;
- public-safe deployment lessons;
- doc improvements;
- issue triage.

Community discussion does not create support commitments, paid support, response targets, or production operating guarantees.

## Reporting Deployment Issues Safely

When reporting a deployment issue, include:

- command or endpoint path;
- public-safe environment summary;
- redacted logs;
- public feed URL only if it is intended to be public;
- validator status without private admin tokens.

Do not include tokens, DB URLs, private keys, admin URLs with secrets, private portal screenshots, private ticket links, raw logs with credentials, or unredacted operator artifacts.

If a report involves a vulnerability or leaked secret, follow `SECURITY.md` instead of opening a public issue.

## Feature Requests

Feature requests should explain the agency or operator problem and the smallest useful outcome. Requests remain out of scope if they ask for rider apps, fare payments, passenger accounts, CAD/dispatch replacement, consumer submission automation, hosted SaaS promises, paid support, or marketplace/vendor-equivalent commitments.

