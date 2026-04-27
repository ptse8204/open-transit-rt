## Summary

Describe the change and why it is needed.

## Scope Check

- [ ] This stays inside Open Transit RT scope.
- [ ] This does not add rider apps, payments, passenger accounts, or CAD/dispatch replacement.
- [ ] This does not change public feed URLs unless explicitly documented and approved.
- [ ] This does not change GTFS-RT contracts unless explicitly documented and approved.
- [ ] This does not add consumer submission automation or contact external portals.

## Tests And Checks

- [ ] `make validate`
- [ ] `make test`
- [ ] `git diff --check`
- [ ] `make realtime-quality` when realtime/readiness behavior or docs changed materially.
- [ ] `make smoke` when hardening, HTTP, validation, readiness, or operations docs changed materially.
- [ ] `docker compose -f deploy/docker-compose.yml config` when deployment docs or Compose assumptions changed.
- [ ] I documented any blocked checks and why.

## Docs, Evidence, And Truthfulness

- [ ] I updated docs for behavior, setup, operations, support, or evidence changes.
- [ ] I did not claim CAL-ITP/Caltrans compliance, consumer acceptance, agency endorsement, hosted SaaS availability, paid support/SLA coverage, vendor equivalence, or universal production readiness without retained evidence.
- [ ] I kept prepared packets as prepared-only unless target-originated evidence supports a status change.
- [ ] I followed `docs/evidence/redaction-policy.md` for evidence changes.

## Security And Private Data

- [ ] I did not commit tokens, DB URLs with passwords, private keys, admin URLs with secrets, private portal screenshots, private ticket links, raw logs with credentials, or unredacted operator artifacts.
- [ ] I checked whether this change needs `security-review-needed`.

## Architecture

- [ ] I considered whether this needs an ADR in `docs/decisions.md`.
- [ ] I preserved Trip Updates pluggability, Vehicle Positions-first posture, draft/published GTFS separation, and conservative matching behavior.

