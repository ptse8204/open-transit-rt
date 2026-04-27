# Docs Assets

These are repo-owned documentation visuals. They are used to teach the current Open Transit RT architecture, local workflows, and public/admin boundaries without making unsupported product or compliance claims.

Public-facing copies of selected PNG assets also live under `wiki/assets/` so `/wiki` pages do not depend on internal relative paths.

## Visual Review Rule

Every generated or generated-assisted image must be manually reviewed before it is referenced in docs.

- Labels must match project behavior.
- Captions must say whether the image is illustrative or exact-behavior based.
- Alt text must describe the useful content of the image, not just repeat the filename.
- Generated drafts may be refined into checked-in SVG/PNG assets when text or labels need cleanup.
- Do not check in visuals that imply hosted SaaS, SSO, CAD/AVL replacement, consumer acceptance, marketplace/vendor equivalence, full CAL-ITP/Caltrans compliance, or universal production readiness.

## Current Assets

### `agency-journey-to-public-feeds.png`

Source:
- `agency-journey-to-public-feeds.svg`
- Regenerated image draft created with the image-generation tool, then manually reviewed and refined into a simpler SVG-derived PNG because generated text and dense labels can be unreliable.

Type:
- Illustrative teaching graphic, not a screenshot.

Used in:
- `README.md`
- `docs/tutorials/agency-demo-flow.md`
- `wiki/agency-demo.md`

Alt text:
- Illustrative agency journey from GTFS import or GTFS Studio drafts through schedule publication, authenticated telemetry, validation, and public GTFS plus GTFS Realtime feeds.

Prompt/spec:
- Create a clean, modern, easy-to-understand agency path graphic. Show the simple path from preparing GTFS, publishing a schedule, adding vehicle data, validating, and publishing public feeds. Keep labels large and readable. Do not show fake UI or unsupported claims.

Truthfulness notes:
- The visual says the path is illustrative and does not claim hosted SaaS, CAD/AVL replacement, consumer acceptance, or full compliance.

### `docs-choose-your-path.png`

Source:
- `docs-choose-your-path.svg`
- Regenerated image draft created with the image-generation tool, then manually reviewed and refined into a simpler SVG-derived PNG.

Type:
- Illustrative docs navigation graphic.

Used in:
- `docs/README.md`
- `wiki/README.md`

Alt text:
- Illustrative documentation guide showing paths for trying locally, running the agency demo, planning deployment, reviewing evidence, and contributing.

Prompt/spec:
- Create a clean docs navigation graphic with five large destination cards: Try locally, Agency demo, Deploy, Evidence, and Contribute. Keep labels large and readable. Do not imply SaaS hosting, consumer acceptance, full compliance, or agency endorsement.

Truthfulness notes:
- The visual points to documentation paths only. It is not a product UI and does not imply agency endorsement.

### `data-flow-through-system.png`

Source:
- `data-flow-through-system.svg`
- Regenerated image draft created with the image-generation tool, then manually reviewed and refined into a simpler SVG-derived PNG.

Type:
- Illustrative system explainer.

Used in:
- `docs/README.md`
- `wiki/how-it-works.md`

Alt text:
- Illustrative data-flow diagram showing GTFS import, GTFS Studio drafts, vehicle telemetry, Open Transit RT state, validation, and public feed outputs.

Prompt/spec:
- Create a clean four-column system explainer showing Inputs, Open Transit RT, Validation, and Public feeds. Keep labels large and readable. Do not imply SSO, SaaS hosting, CAD/AVL replacement, full compliance, consumer acceptance, or external predictor integration.

Truthfulness notes:
- The visual summarizes current boundaries. Trip Updates remain behind the prediction adapter; validation checks generated artifacts rather than proving consumer acceptance.

### `architecture-overview.png`

Source:
- `architecture-overview.svg`

Type:
- Exact-behavior architecture summary rendered from a reviewed SVG spec.

Used in:
- Available as a deeper architecture reference; not currently used on the README front door.

Alt text:
- Architecture overview showing schedule inputs, vehicle state, feed builders, public feed outputs, and validation/readiness evidence as four clear areas.

Prompt/spec:
- Clean architecture summary with four bounded areas: schedule inputs, vehicle state, feed builders, and public feeds. Use large readable labels, no crossing lines, and no phase-history wording.

### `agency-deployment.png`

Source:
- `agency-deployment.svg`

Type:
- Exact-behavior deployment-shape diagram rendered from a reviewed SVG spec.

Used in:
- Deployment/readiness documentation as needed.

Alt text:
- Small agency deployment diagram showing public internet, TLS reverse proxy, anonymous public feeds, protected admin tools, Open Transit RT services, pinned validators, prediction adapter, and Postgres/PostGIS.

Prompt/spec:
- Clean small-agency deployment diagram. Show TLS reverse proxy, anonymous public feeds, protected admin tools, Open Transit RT services, Postgres/PostGIS, pinned validators, and optional predictor adapter. Keep text bounded and avoid overlapping connectors. Do not show unsupported SSO, consumer acceptance, or production-ready claims.

### `quickstart-flow.png`

Source:
- `quickstart-flow.svg`

Type:
- Exact-behavior local workflow diagram rendered from a reviewed SVG spec.

Used in:
- `docs/tutorials/README.md`
- `docs/tutorials/local-quickstart.md`

Alt text:
- Exact-behavior local quickstart flow grouped into setup, publish-and-ingest, and review-output steps.

Prompt/spec:
- Clean numbered quickstart flow grouped into setup, publish-and-ingest, and review-output steps. Use large bounded labels and avoid long connector lines.

### `public-vs-admin-endpoints.png`

Source:
- `public-vs-admin-endpoints.svg`

Type:
- Exact-behavior endpoint-boundary diagram rendered from a reviewed SVG spec.

Used in:
- `docs/tutorials/agency-demo-flow.md`

Alt text:
- Exact-behavior endpoint boundary diagram showing anonymous public feed routes separated from protected admin, debug, telemetry, validation, and scorecard examples.

Prompt/spec:
- Clean two-column endpoint boundary. Public side includes `schedule.zip`, `feeds.json`, `vehicle_positions.pb`, `trip_updates.pb`, and `alerts.pb`. Protected side includes examples such as GTFS Studio, JSON debug feeds, telemetry events, validation, and scorecard. Keep labels bounded and readable.

### `how-to-contribute-paths.png`

Source/generation method:
- Generated with the image generation tool, copied into `docs/assets/`, and manually reviewed for label truthfulness.

Type:
- Illustrative teaching graphic, not a screenshot.

Purpose:
- Show contribution paths such as reporting a bug, improving docs, suggesting a feature, submitting code, and helping with evidence/runbooks.

Used in:
- `CONTRIBUTING.md`
- `docs/README.md`

Alt text:
- Illustrative contribution paths: report a bug, improve docs, suggest a feature, submit code, and help with evidence runbooks.

Prompt/spec:
- Create an illustrative infographic showing contribution paths for an open-source transit backend project. Use a clean horizontal map of five paths leading into a central project hub. Include labels: Report Bug, Improve Docs, Suggest Feature, Submit Code, Evidence Runbooks. Keep it polished, flat, readable, and not a fake UI screenshot.

Truthfulness notes:
- The visual shows ways to contribute only. It does not imply compliance, consumer acceptance, agency endorsement, hosted SaaS availability, paid support, SLA coverage, vendor equivalence, or production readiness.

### `community-workflow.png`

Source/generation method:
- Generated with the image generation tool, copied into `docs/assets/`, and manually reviewed for label truthfulness.

Type:
- Illustrative teaching graphic, not a screenshot.

Purpose:
- Show the community workflow from issue discussion through implementation, testing, PR review, merge, and release/docs update.

Used in:
- `CONTRIBUTING.md`
- `docs/governance.md`

Alt text:
- Illustrative community workflow: Issue, Discuss, Implement, Test, PR, Review, Merge, Release Docs Update.

Prompt/spec:
- Create an illustrative workflow graphic for open-source governance. Show connected stations on a transit line labeled Issue, Discuss, Implement, Test, PR, Review, Merge, Release Docs Update. Use a polished flat style and avoid browser chrome or fake product UI.

Truthfulness notes:
- The visual illustrates process only. It does not imply merge entitlement, release timing, compliance, consumer acceptance, agency endorsement, hosted SaaS availability, paid support, or production readiness.

### `single-vs-multi-agency.png`

Source/generation method:
- Generated with the image generation tool, copied into `docs/assets/`, and manually reviewed for label truthfulness.

Type:
- Illustrative teaching graphic, not a screenshot.

Purpose:
- Compare the current single-agency/local-demo/pilot deployment model with possible future multi-agency hosted options.

Used in:
- `docs/multi-agency-strategy.md`

Alt text:
- Illustrative comparison of the current single-agency/local-demo/pilot model and future multi-agency hosted options that require code changes.

Prompt/spec:
- Create a two-column comparison graphic. Left column: Current Single Agency, Local Demo, Pilot Host, One Feed Root. Right column: Future Multi Agency Options, Agency Boundary, Shared Infrastructure, Separate Feed Roots, Code Changes Needed. Make the future side conditional/planned and avoid fake UI.

Truthfulness notes:
- The visual explicitly marks future multi-agency hosting as requiring code changes. It does not imply hosted SaaS availability, agency endorsement, consumer acceptance, compliance, paid support, SLA coverage, vendor equivalence, or production readiness.

### `evidence-maturity-ladder.png`

Source/generation method:
- Generated with the image generation tool, copied into `docs/assets/`, and manually reviewed for label truthfulness.

Type:
- Illustrative teaching graphic, not a screenshot or evidence artifact.

Purpose:
- Show evidence maturity stages from code existence through hosted evidence, prepared packet, submitted, under review, and accepted.

Used in:
- `docs/roadmap-status.md`
- `docs/compliance-evidence-checklist.md`

Alt text:
- Illustrative evidence maturity ladder showing code exists, hosted evidence, prepared packet, submitted, under review, and accepted as separate evidence stages.

Prompt/spec:
- Create an evidence maturity ladder with steps labeled Code Exists, Hosted Evidence, Prepared Packet, Submitted, Under Review, Accepted. Make the latter stages visually neutral and evidence-required, not completed. Include the idea that claims need evidence.

Truthfulness notes:
- The visual separates prepared packets from submitted, under-review, and accepted states. It does not imply current acceptance, compliance, consumer ingestion, agency endorsement, hosted SaaS availability, paid support, SLA coverage, vendor equivalence, or production readiness.

### `support-boundaries.png`

Source/generation method:
- Generated with the image generation tool, copied into `docs/assets/`, and manually reviewed for label truthfulness.

Type:
- Illustrative teaching graphic, not a screenshot.

Purpose:
- Explain maintainer help, operator-owned responsibilities, and community-only support boundaries.

Used in:
- `docs/support-boundaries.md`
- `wiki/support-and-contribute.md`

Alt text:
- Illustrative support boundary diagram showing maintainer help, operator-owned responsibilities, and community-only support.

Prompt/spec:
- Create a three-column support boundaries diagram with headings Maintainer Help, Operator Owned, Community Only. Include maintainer items such as code review, docs, reproducible bugs; operator-owned items such as secrets, DNS/TLS, deployments, private logs; community-only items such as ideas, peer notes, examples. Include "No paid support or SLA."

Truthfulness notes:
- The visual states that support is bounded and does not promise paid support, SLA coverage, hosted SaaS availability, agency endorsement, consumer acceptance, compliance, vendor equivalence, or universal production readiness.
