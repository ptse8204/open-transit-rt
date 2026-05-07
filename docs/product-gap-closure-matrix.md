# Product Gap Closure Matrix — Self-Hosted Agency Reuse

| Gap | Why it matters | Current state | Next action | Done when |
| --- | --- | --- | --- | --- |
| Root README is wrong | Agencies see the wrong first impression | Root README appears to be a roadmap export | Restore product README | README explains product, quickstart, OCI/OCL deployment, CAL-ITP-style readiness |
| Public GTFS onboarding is too manual | Agencies need easy reuse | Phase 33 run worked but was not fully productized | Add guided public GTFS onboarding command/runbook | User can import a public GTFS with agency ID and URL without manual DB work |
| OCI/OCL server path not productized enough | This is the target proof/deployment path | OCI pilot exists with evidence/runbook | Create reference deployment profile | Fresh operator can deploy/update/rollback on similar server |
| CAL-ITP workflow not productized enough | Project goal is CAL-ITP-style integration | Docs exist, workflow scattered | Add readiness workflow in docs/UI | Operators see gaps and next steps in one place |
| Existing solution integration needs clearer adapter kit | Agencies may already have AVL/predictor tools | Adapter boundaries exist | Package adapter kit | New adapter can be built from examples/tests |
| Static validator warnings remain | Warnings may matter for readiness quality | 0 system errors, 3 warnings | Decide whether to remediate or document | Warning policy is explicit |
| Real telemetry flow is not agency-easy | Vehicle Positions need telemetry | Device scripts exist | Improve onboarding UX/support bundle | Agency can create token and test telemetry quickly |
| Operations support bundle missing | Troubleshooting should be easy | Logs/runbooks exist | Add support summary/export guidance | Operator can provide redacted bundle for help |
