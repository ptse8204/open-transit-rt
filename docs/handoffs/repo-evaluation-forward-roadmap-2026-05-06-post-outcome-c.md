# Repo Evaluation And Forward Roadmap — Post-Phase-33 Outcome C

Generated: 2026-05-06

## Current conclusion

Open Transit RT has moved past the earlier Phase 33 blocker. The current repo baseline is:

```text
Phases 0 through 33 are closed for their documented scopes.
Phase 33 is complete as Outcome C — public-GTFS local/pilot run completed with public-safe retained summaries.
```

This is a meaningful improvement. The repo now has retained local/pilot evidence that it can ingest a real public static GTFS ZIP, publish it locally, verify that the fetched schedule is the imported public GTFS rather than the sample feed, fetch all five public paths, retain validator results or blockers, run a synthetic telemetry dry-run, and check admin/private route boundaries.

The Outcome C evidence remains limited to local/pilot public static GTFS dataset handling. It does not prove agency adoption, agency approval, official agency feed status, agency-owned final-root readiness, consumer evidence, compliance, hosted SaaS, production readiness, real realtime data, vendor AVL compatibility, real-world ETA accuracy, or production-grade ETA quality.

## What is now fixed

The previous strongest technical blocker was large public GTFS import/publish at LA Metro scale. That blocker is now recorded as fixed inside the Phase 33 Outcome C path:

- `cmd/gtfs-import` supports configurable timeout handling through `-timeout` and `GTFS_IMPORT_TIMEOUT`.
- The default import timeout is 15 minutes; `-timeout 0` disables the import timeout.
- Large `gtfs_stop_time` and `gtfs_shape_point` publish inserts use `pgx.CopyFrom`.
- Publish-failure reporting uses a fresh short context so a canceled import context does not leave the import row stuck at `started`.
- Focused import timeout and failure-report tests were added.
- The LA Metro Bus public GTFS import completed locally and published `gtfs-import-1` for local `LACMTA` evidence setup.

## What remains missing

### 1. Agency-owned or agency-approved final-root proof

Still missing:

- candidate final public feed root;
- retained agency/operator approval;
- DNS proof;
- TLS proof;
- HTTP-to-HTTPS redirect proof if HTTP is exposed;
- anonymous public fetch proof for all five public feed URLs;
- validator records at the final root;
- redacted proxy/config summary;
- final-root evidence packet README and checksums;
- prepared consumer packet refresh using the final root.

The DuckDNS OCI pilot remains hosted/operator pilot evidence only. It is not agency-owned or agency-approved final-root proof.

### 2. Consumer or aggregator evidence

All seven consumer/aggregator targets remain `prepared` only. There is still no evidence for:

- submitted;
- under review;
- accepted;
- rejected;
- blocked by a target-originated artifact;
- ingestion;
- listing;
- display;
- adoption.

Do not edit `docs/evidence/consumer-submissions/status.json`, target records, or target-specific artifact directories unless retained, redacted, target-originated evidence exists.

### 3. Real agency pilot evidence

The repo has agency pilot materials, but no real agency pilot evidence. Missing:

- pilot authorization;
- agency/operator retained artifacts;
- real pilot kickoff/closeout evidence;
- agency feedback;
- agency-approved metadata/license/contact proof;
- real operator decision to continue, pause, or stop.

### 4. Real deployment operations evidence refresh

The OCI pilot evidence remains useful, but the repo does not yet have a newer retained operations refresh after Phase 33. Missing or stale for a current deployment:

- current public HTTPS root evidence beyond OCI pilot scope;
- current backup/restore evidence refresh;
- current feed monitor evidence;
- current scorecard export evidence;
- current alert lifecycle evidence;
- current validator evidence at a selected deployment root.

### 5. Static GTFS validator execution for the Phase 33 public dataset

Outcome C attempted static GTFS validation, but the validator did not execute because Java was unavailable in the local environment. This is a tooling/environment blocker, not a data-quality pass or fail.

Future work should either:

- provide a documented Java/static-validator preflight for public-GTFS local pilot runs; or
- re-run the static validator in an environment where Java is available; or
- keep the existing blocker language if no such environment exists.

### 6. Public-GTFS local/pilot repeatability

Outcome C succeeded, but it used repo service binaries with `AGENCY_ID=LACMTA` and a local public proxy rather than the default `make agency-app-up` path, because the default app imports the repo sample feed for `demo-agency`.

Missing repeatability work:

- a documented public-GTFS local pilot flow;
- a public agency setup pattern that does not imply agency approval;
- optional script or make target for repeating a public-GTFS local run with an arbitrary public GTFS ZIP;
- clear separation between sample demo flow and public-GTFS evaluation flow.

### 7. Status and roadmap consistency

Several repo docs appear stale or internally inconsistent after Outcome C:

- `docs/current-status.md` still begins by calling the project an early-stage starter, even though later lines document Phases 0-33 and Outcome C.
- `docs/roadmap-status.md` still says current evidence includes an attempted public-GTFS blocker and that completed public-GTFS evidence is missing, while later lines say Phase 33 Outcome C is complete.
- `docs/track-b-productization-roadmap.md` still says the recommended next phase is Phase 32, even though Phase 32 and Phase 33 are closed.
- `docs/repo-gaps.md` still lists starter-repo gaps that the repo has already filled.
- `docs/phase-plan.md` should point to the post-Outcome-C roadmap or list the next roadmap phase so future agents do not infer that Phase 33 is the end of the plan.
- `docs/README.md` labels the retained public-GTFS evidence as an “attempt,” which should be updated to “Outcome C evidence” or similar.

## Recommended immediate next phase

```text
Phase 34 — Post-Outcome-C Status Consistency And Evidence Readiness
```

Phase 34 should not fabricate new external evidence. It should:

1. patch stale status/roadmap docs;
2. add a final-root operator request package;
3. document a repeatable public-GTFS local/pilot run path;
4. clarify the static validator Java/tooling blocker;
5. select the next external-evidence fork without weakening truthfulness boundaries.

## Recommended forward roadmap

| Phase | Name | Purpose | Evidence boundary |
| --- | --- | --- | --- |
| 34 | Post-Outcome-C Status Consistency And Evidence Readiness | Patch stale docs, add final-root request package, document repeatable public-GTFS local pilot flow. | No new external evidence unless actually collected. |
| 35 | Agency-Owned Or Agency-Approved Final-Root Proof | Obtain approval, DNS/TLS, five public feed fetches, validator records, proxy/config summary, checksums. | Only if real retained operator/agency evidence exists. |
| 36 | Real Agency Pilot Execution | Run the agency pilot package with a permissioned agency/operator. | No adoption claim without agency-retained evidence. |
| 37 | Real Deployment Operations Evidence Refresh | Refresh backup, restore, monitoring, validation, scorecard, alert lifecycle, and admin-boundary evidence. | Deployment/operator proof only, not consumer acceptance. |
| 38 | Real Device Or Vendor AVL Pilot | Connect authorized real device/vendor telemetry through the adapter boundary. | No vendor compatibility claim without real authorized evidence. |
| 39 | Realtime Quality And External Predictor Runtime Evaluation | Add real observed-arrival metrics and/or optional runtime predictor behind the adapter. | No production ETA-quality claim without real measurements. |
| 40 | Authorized Consumer Submission | Submit one target only with operator authorization and official-path verification. | Target status changes only from target-originated evidence. |
| 41 | Production Packaging And Multi-Agency Readiness | Harden multi-agency operations, tenant-safe evidence/export/backup, release artifacts. | No hosted SaaS or production multi-tenant claim without proof. |
| 42 | Evidence-Bounded Public Launch | Launch only when public wording matches retained evidence. | No launch claim until launch actually occurs. |

## Preferred next Codex entry point

Use:

```text
docs/handoffs/phase-34-kickoff.md
```

That file should be treated as a kickoff prompt for a fresh Codex agent.
