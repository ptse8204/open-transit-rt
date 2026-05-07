# Future Roadmap — Post-Phase-33 Outcome C

## Baseline

Phase 33 is complete as Outcome C for local/pilot public static GTFS dataset handling. The repo has now proven more than toy-fixture import, but it still has not proven agency adoption, final-root readiness, consumer acceptance, production readiness, real realtime data, real device/vendor AVL, or production-grade ETA quality.

This roadmap starts after Outcome C. Phase 35 realigned the default path away
from external-proof chasing and toward reusable self-hosted agency deployment.

## Evidence ladder after Outcome C

```text
1. Code and local fixtures                         done
2. Hosted/operator OCI pilot evidence              done for pilot root only
3. Prepared consumer packets                       done, all targets prepared only
4. Public GTFS local/pilot evidence                done as Phase 33 Outcome C
5. Status/repeatability/final-root request patch    done as Phase 34
6. README and roadmap realignment                  done as Phase 35
7. OCI/OCL reference deployment productization     current recommended path
8. Reusable agency onboarding flow                 future productization
9. Integration adapter kit                         future productization
10. CAL-ITP-style readiness workflow in product    future productization
11. Agency-owned/final-root proof                  future optional external evidence
12. Authorized consumer submission                 future target-originated evidence
13. Consumer acceptance/ingestion/display          future target-originated evidence
14. Evidence-bounded public launch                 future public action
```

## Immediate roadmap

| Phase | Name | When to do it | What it produces | What it must not claim |
| --- | --- | --- | --- | --- |
| 34 | Post-Outcome-C Status Consistency And Evidence Readiness | Complete | Consistent docs, final-root request package, repeatable public-GTFS local pilot guide, static-validator blocker clarification. | No new external evidence was created. |
| 35 | README And Roadmap Realignment | Complete | Product-front-door README and self-hosted agency reuse continuation path. | No runtime change or new evidence. |
| 36 | OCI/OCL Reference Deployment Productization | Next recommended phase | Repeatable reference deployment docs, env placeholders, update/rollback, validators, backup/restore, feed monitor, scorecard, and smoke checklist. | No hosted SaaS, production readiness, final-root, or compliance claim. |
| 37 | Agency Reusable Onboarding Flow | After Phase 36 reference path is clear | Guided local/server onboarding flow for agency ID, GTFS URL, metadata, import timeout, feed URLs, validators, and support summary. | No agency adoption or official-feed claim from tooling alone. |
| 38 | Integration Adapter Kit | After onboarding flow | Reusable AVL/device adapter guidance, fixture patterns, external predictor adapter lifecycle, validator/monitoring/consumer boundaries. | No certified vendor compatibility or production AVL reliability claim. |
| 39 | CAL-ITP-Style Readiness Workflow | After deployment/onboarding basics | Product-facing readiness checklist for stable URLs, metadata, validation, freshness, GTFS-RT completeness, and consumer packet state. | No CAL-ITP/Caltrans compliance claim. |
| 40 | Agency-Owned Or Agency-Approved Final-Root Proof | Optional later, when an agency/operator can approve a root | DNS/TLS/redirect/fetch/validator/proxy/checksum packet for an approved final root. | No consumer acceptance or compliance claim from final-root proof alone. |
| 41 | Authorized Consumer Submission | Optional later, when operator selects a target and official path is verified | Target-originated submission receipt/review/response artifacts. | No target status change without target-originated evidence. |

## Decision tree

### If continuing the default productization roadmap

Go to Phase 36.

Productize the OCI/OCL-style pilot server as the reference deployment path for
self-hosted agency reuse. Keep the DuckDNS OCI pilot labeled as hosted/operator
pilot evidence only.

### If an agency-owned or agency-approved final root is available later

Use the optional final-root proof track after the self-hosted reference path is
clear or when a real operator can provide approved artifacts.

Required evidence:

- approval artifact;
- DNS proof;
- TLS proof;
- redirect proof if HTTP is exposed;
- five final feed URL fetches;
- static and realtime validator records;
- redacted proxy/config summary;
- evidence packet README and checksums;
- prepared packet refresh only after final-root values are known.

### If no final root exists but an agency/operator is interested later

Use the agency pilot package and the future reusable onboarding flow. Keep the
wording as pilot/evaluation unless the agency provides retained approval for
stronger language.

### If no agency is available but a pilot environment is running later

Refresh operations evidence only as a regression/deployment proof path. This
can improve credibility without pretending agency adoption or consumer
acceptance.

### If real device/vendor payloads are available with authorization later

Keep vendor/device data private and redacted. Use the adapter boundary. Do not commit credentials, raw telemetry, or private correspondence.

### If real observed-arrival data or an external predictor runtime is available later

Compare behavior against real observations, document fallback behavior, and preserve the default conservative predictor path unless the adapter proves better under retained evidence.

### If a consumer target is selected and official path is verified

Go to Phase 40.

Update only the selected target status and only from retained target-originated artifacts.

## Current high-priority productization path

Phase 35 closed the README and roadmap realignment step. The next recommended
phase is Phase 36 — OCI/OCL Reference Deployment Productization.

The current self-hosted continuation path is:

1. OCI/OCL reference deployment productization;
2. agency reusable onboarding flow;
3. integration adapter kit;
4. CAL-ITP-style readiness workflow in product.

External-proof tracks remain available later when retained, redacted,
claim-specific artifacts exist. They are not the default roadmap.

## Claims policy

The safest public wording after Outcome C is:

```text
Open Transit RT has retained local/pilot evidence that it can import and publish a real public static GTFS dataset using the LA Metro Bus public GTFS feed. This evidence is local/pilot only and does not prove agency adoption, official feed status, final-root readiness, consumer acceptance, compliance, production readiness, real realtime data, or ETA quality.
```

Do not shorten this into a stronger claim.
