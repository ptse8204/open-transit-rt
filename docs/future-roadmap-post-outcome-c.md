# Future Roadmap — Post-Phase-33 Outcome C

## Baseline

Phase 33 is complete as Outcome C for local/pilot public static GTFS dataset handling. The repo has now proven more than toy-fixture import, but it still has not proven agency adoption, final-root readiness, consumer acceptance, production readiness, real realtime data, real device/vendor AVL, or production-grade ETA quality.

This roadmap starts after Outcome C.

## Evidence ladder after Outcome C

```text
1. Code and local fixtures                         done
2. Hosted/operator OCI pilot evidence              done for pilot root only
3. Prepared consumer packets                       done, all targets prepared only
4. Public GTFS local/pilot evidence                done as Phase 33 Outcome C
5. Status/repeatability/final-root request patch    done as Phase 34
6. Agency-owned/final-root proof                   future external evidence
7. Real agency pilot                               future external evidence
8. Real deployment operations refresh              future deployment evidence
9. Real device/vendor AVL evidence                 future external evidence
10. Real-world realtime quality evidence           future measurement evidence
11. Authorized consumer submission                 future target-originated evidence
12. Consumer acceptance/ingestion/display          future target-originated evidence
13. Evidence-bounded public launch                 future public action
```

## Immediate roadmap

| Phase | Name | When to do it | What it produces | What it must not claim |
| --- | --- | --- | --- | --- |
| 34 | Post-Outcome-C Status Consistency And Evidence Readiness | Complete | Consistent docs, final-root request package, repeatable public-GTFS local pilot guide, static-validator blocker clarification. | No new external evidence was created. |
| 35 | Agency-Owned Or Agency-Approved Final-Root Proof | When an agency/operator can approve a root | DNS/TLS/redirect/fetch/validator/proxy/checksum packet for an approved final root. | No consumer acceptance or compliance claim from final-root proof alone. |
| 36 | Real Agency Pilot Execution | When an agency/operator is ready to pilot | Pilot kickoff, agency feedback, local or hosted evidence, closeout, decision record. | No adoption claim unless the agency explicitly supports it with retained evidence. |
| 37 | Real Deployment Operations Evidence Refresh | When a pilot or hosted environment is running | Backup, restore, monitoring, validator, scorecard, alert lifecycle, role/admin boundary evidence. | No final-root or consumer claim unless those artifacts exist. |
| 38 | Real Device Or Vendor AVL Pilot | When authorized real device/vendor data is available | Redacted adapter/device evidence, ingest proof, freshness/quality notes. | No certified vendor compatibility or production reliability claim. |
| 39 | Real-World Realtime Quality And Predictor Evaluation | When observed-arrival data or external predictor runtime is available | Route/time-period metrics, ETA error summaries, fallback behavior, adapter health evidence. | No production-grade ETA claim without sustained real-world evidence. |
| 40 | Authorized Consumer Submission | When operator selects a target and official path is verified | Target-originated submission receipt/review/response artifacts. | No target status change without target-originated evidence. |
| 41 | Production Packaging And Multi-Agency Readiness | After core evidence paths are stronger | Tenant-safe ops, release packaging, upgrade/rollback, capacity notes, deployment patterns. | No hosted SaaS, SLA, or production multi-tenant claim without proof. |
| 42 | Evidence-Bounded Public Launch | After public wording can be supported | Published copy, outreach log, launch evidence packet. | No “launched” claim before actual launch. |

## Decision tree

### If an agency-owned or agency-approved final root is available

Go to Phase 35.

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

### If no final root exists but an agency/operator is interested

Go to Phase 36.

Use the agency pilot package. Keep the wording as pilot/evaluation unless the agency provides retained approval for stronger language.

### If no agency is available but a pilot environment is running

Go to Phase 37.

Refresh operations evidence. This can improve credibility without pretending agency adoption or consumer acceptance.

### If real device/vendor payloads are available with authorization

Go to Phase 38.

Keep vendor/device data private and redacted. Use the adapter boundary. Do not commit credentials, raw telemetry, or private correspondence.

### If real observed-arrival data or an external predictor runtime is available

Go to Phase 39.

Compare behavior against real observations, document fallback behavior, and preserve the default conservative predictor path unless the adapter proves better under retained evidence.

### If a consumer target is selected and official path is verified

Go to Phase 40.

Update only the selected target status and only from retained target-originated artifacts.

## Current high-priority evidence forks

Phase 34 closed the status/repeatability/final-root request maintenance step.
The next path should be selected from these retained-evidence forks:

1. agency-owned/final-root proof;
2. authorized target-specific consumer submission evidence;
3. real agency pilot evidence;
4. real deployment operations refresh;
5. real device/vendor AVL evidence;
6. real-world realtime quality evidence.

## Claims policy

The safest public wording after Outcome C is:

```text
Open Transit RT has retained local/pilot evidence that it can import and publish a real public static GTFS dataset using the LA Metro Bus public GTFS feed. This evidence is local/pilot only and does not prove agency adoption, official feed status, final-root readiness, consumer acceptance, compliance, production readiness, real realtime data, or ETA quality.
```

Do not shorten this into a stronger claim.
