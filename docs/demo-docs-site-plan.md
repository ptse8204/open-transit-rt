# Demo And Documentation Site Plan

## Status

Phase 66 Checkpoint 000005 adds this plan only. It does not build, host,
publish, or announce a website. It does not create retained evidence, contact
agencies or consumers, change consumer tracker status, or claim public launch,
hosted SaaS, production readiness, agency approval, consumer acceptance, vendor
compatibility, or CAL-ITP/Caltrans compliance.

## Goal

Make the public documentation easier to turn into a future static demo/docs
site without changing the product boundary. A small agency evaluator should be
able to understand what Open Transit RT is, try the local demo, review
connectors, understand readiness boundaries, and find operator setup docs
without reading phase history.

## Recommended Information Architecture

| Section | Source Today | Purpose | Boundary |
| --- | --- | --- | --- |
| Start | `README.md`, `wiki/README.md` | Explain the product and fastest local evaluator path | No hosted-service or public-launch claim |
| Try Locally | `wiki/agency-demo.md`, `docs/tutorials/agency-first-run.md`, `docs/tutorials/local-quickstart.md` | Guide a 30-minute local demo | Local evaluator workflow only |
| Agency Fit | `wiki/can-my-agency-use-this.md`, `docs/agency-one-pager.md` | Help agencies decide whether a self-hosted path fits | No agency approval or adoption claim |
| Connectors | `docs/integration-adapter-kit.md`, `wiki/connector-cookbook.md`, `docs/connectors/plugin-contract.md` | Explain sidecars, manifests, command adapters, and examples | No dynamic plugin loading or vendor compatibility claim |
| Readiness | `wiki/calitp-readiness-plain-english.md`, `docs/release-candidate-readiness.md`, `docs/roadmap-status.md` | Explain readiness signals, validators, and release-candidate checks | No compliance or production-readiness claim |
| Deploy | `wiki/deployment-guide.md`, `docs/tutorials/deploy-with-docker-compose.md`, `docs/upgrade-and-rollback.md` | Point to self-hosted evaluation and local/source release paths | No hosted SaaS or registry-published app image claim |
| Help Improve | `wiki/how-agencies-can-help.md`, `wiki/support-and-contribute.md`, `CONTRIBUTING.md` | Invite safe feedback, issues, docs, connector examples, and authorized pilots | No external contact automation |

## Page Requirements

- Use plain language and a small-agency operator lens.
- Prefer existing Markdown pages and wiki pages over a new frontend stack.
- Link to command-level docs rather than duplicating long command sequences.
- Keep screenshots or visuals illustrative unless they are generated from a
  current running app and explicitly marked as local/private diagnostics.
- Keep evidence and consumer-status wording behind the existing evidence
  boundary docs.
- Keep the safe plugin definition visible on connector-facing pages:
  "In Open Transit RT, a plugin is an optional sidecar, command adapter,
  manifest, or connector process. It is not arbitrary dynamic code loaded into
  the backend."

## Future Static Site Readiness

A future maintainer may choose a lightweight static documentation site only
after deciding:

- hosting target, if any;
- source directory and build command;
- asset ownership and alt-text review;
- URL and redirect policy;
- no-secrets and no-evidence publishing audit;
- release-note wording;
- public-launch review using `docs/public-launch-checklist.md`;
- claim-boundary audit using `make audit-final-claim-review`.

Until that decision exists, repository Markdown and wiki pages are the
canonical public-friendly documentation surfaces.

## Non-Goals

- No marketing site implementation.
- No hosted website or public launch.
- No outreach, announcement, consumer contact, or agency contact.
- No screenshots or evidence capture.
- No consumer tracker or evidence status change.
- No claim that any agency, consumer, vendor, or regulator has approved,
  accepted, listed, certified, adopted, or endorsed Open Transit RT.
- No new frontend framework, analytics product, tracking pixel, mailing list,
  portal, account system, or paid-support funnel.

## Validation

For documentation-only updates, run:

```sh
git diff --check
make check
make audit-final-claim-review
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
```

Also preserve:

```sh
git diff --exit-code -- docs/evidence/consumer-submissions/status.json
git diff --exit-code -- docs/evidence/captured
git diff --exit-code -- db/migrations go.mod go.sum
```

## Rollback

Rollback is a normal docs revert of the checkpoint commit that added this plan
and navigation links. No database, feed URL, runtime, evidence, or consumer
status rollback should be required.
