# Governance

Open Transit RT is an independent open-source project. Governance is intentionally lightweight, but decisions that affect architecture, evidence wording, security posture, release process, or multi-agency assumptions must be explicit.

![Illustrative community workflow: Issue, Discuss, Implement, Test, PR, Review, Merge, Release Docs Update.](assets/community-workflow.png)

## Maintainer Role

Maintainers are responsible for:

- preserving the product scope in `AGENTS.md`;
- reviewing pull requests;
- protecting security and evidence boundaries;
- keeping public wording truthful;
- deciding when a change needs an ADR;
- cutting releases and writing release notes.

## Authority

- Maintainers with repository write access can merge PRs after review.
- Maintainers can cut releases and create version tags.
- Maintainers must approve docs/evidence wording that changes readiness, compliance, consumer, security, support, or deployment claims.
- Maintainers may close or redirect issues that expose secrets, private operator artifacts, unsupported claims, or out-of-scope work.

## Decision Process

Most changes can be resolved in the issue or PR discussion. Architecture-significant decisions should be recorded in `docs/decisions.md`.

Competing design decisions are resolved by maintainers after considering:

- direct user or maintainer instructions;
- `AGENTS.md`;
- binding requirements docs;
- existing architecture decisions;
- tests and reproducible evidence;
- long-term maintainability for small agencies.

When a decision changes service boundaries, persistence model, feed contracts, security posture, evidence interpretation, or multi-agency assumptions, update the relevant docs and add or amend an ADR.

## Evidence And Claim Approval

Evidence and readiness wording needs maintainer review when it changes any claim about:

- consumer packet status;
- submission, review, acceptance, rejection, or blocker state;
- CAL-ITP/Caltrans readiness or compliance;
- marketplace/vendor equivalence;
- agency endorsement;
- hosted/operator evidence;
- production ETA quality;
- support commitments.

Prepared packets are not submissions. Validator success and public fetch proof are not consumer acceptance. OCI pilot evidence is pilot evidence, not agency-owned production proof.

## Agency Feature Requests

Agencies and operators should open feature requests with public-safe context:

- the operational problem;
- agency size or service pattern, if safe to share;
- affected workflow;
- expected public feed or operator outcome;
- redacted examples.

Do not post private credentials, private portal screenshots, private ticket links, or raw logs with credentials.

## Security And Conduct

Security issues and leaked secrets follow `SECURITY.md`. Conduct expectations follow `CODE_OF_CONDUCT.md`.

Do not open public issues for suspected vulnerabilities, leaked credentials, private keys, private operator artifacts, or unsafe evidence material.

## Out Of Scope

The project does not add rider apps, fare payments, passenger accounts, CAD/dispatch replacement, consumer submission automation, hosted SaaS promises, paid support commitments, or marketplace/vendor-equivalent claims as part of governance work.

