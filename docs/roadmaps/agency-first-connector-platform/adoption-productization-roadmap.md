# Adoption Productization Roadmap

This roadmap keeps the default path focused on product quality and
small-agency adoption readiness. It explicitly removes a real agency pilot as
the default next step.

Evidence, real agency pilots, final-root proof, consumer submission, vendor
proof, and production ETA proof are optional future tracks only after explicit
maintainer authorization, public-safe retention rules, redaction rules, and
stop conditions exist.

## Sequence

| Order | Workstream | Outcome |
| --- | --- | --- |
| 1 | No-CLI agency operations UI | The private Operations Console starts with Agency Operations Cockpit / Start Here, showing setup progress, GTFS, five feed paths, validators, telemetry, readiness, connectors, maintenance, and next actions. |
| 2 | GTFS import/update/rollback UX | Browser import shows source review, checksum, byte count, import counts, active feed version, schedule identity, and truthful rollback/staging limitations. |
| 3 | Feed health and validator UX | Feed Health tracks exactly five public paths; Validator Health explains internal import validation, canonical static validation, GTFS-RT validation, tooling states, stale reports, and safe browser-triggered allowlisted runs. |
| 4 | Synthetic realtime demo and telemetry onboarding | Devices, telemetry, simulator, Vehicle Positions, Trip Updates, and Alerts pages explain what is empty, stale, withheld, matched, unknown, or blocked without exposing token values. |
| 5 | Remote diagnostics and off-host validation | `make oci-reference-check` and `make validate-public-feeds` produce private `.cache` diagnostics for tiny reference servers without writing evidence or requiring validators on the server. |
| 6 | Maintenance center | `/admin/operations/maintenance` summarizes version presence, active feed, last import, five-feed check, validators, backup/restore configuration presence, telemetry freshness, service health availability, support summary, and weekly/monthly tasks. |
| 7 | Docs/wiki/site adoption upgrade | README, docs hub, wiki, tutorials, deployment docs, screenshots, diagrams, and site planning explain browser-first adoption, technical-helper cases, connectors, off-host validation, and claim boundaries. |
| 8 | v0.1.0-rc.1 agency evaluation release | A release-candidate gate checks clean checkout, UI-first operations, feed paths, validators, telemetry simulator, connectors, diagnostics, docs, protected paths, and final claim audit. |
| 9 | Optional evidence tracks only after authorization | Real agency pilot, final-root proof, consumer submission, vendor proof, and public claim tracks stay closed unless explicitly authorized. |

## Website / gh-pages Notes

The repository has a `gh-pages` branch for the static product explainer site.
This Phase 71 work did not switch branches or edit the published site payload
because the current worktree is on `main` and the safe workflow for updating
Pages assets is a separate branch-specific operation.

When the maintainer authorizes a `gh-pages` update, update the site to:

- make the first-viewport path "Agency Operations Cockpit / Start Here" and
  browser-first GTFS/GTFS-RT operations;
- add links to `docs/tutorials/no-cli-agency-first-run.md` and
  `docs/tutorials/small-agency-maintenance-guide.md`;
- mention `make validate-public-feeds` and `make oci-reference-check` as
  private diagnostics, not evidence;
- show Feed Health, Maintenance Center, GTFS import, validator health,
  devices/telemetry, connectors, and readiness as private UI surfaces;
- keep screenshots labeled local/demo and claim-bounded;
- avoid hosted-service, compliance, consumer, adoption, vendor, SLA,
  production, and ETA-quality claims.

## Success Criteria

A non-expert evaluator can use the web UI to import or review GTFS, inspect
the five public feed paths, understand GTFS quality, understand validator
state, review telemetry/device readiness, understand realtime feed state,
review connectors, and know the next maintenance action without relying
primarily on command-line tools.

The remaining command-line paths are clearly marked as technical-helper paths
for deployment, validators, secure device-token handling, remote diagnostics,
support summaries, and future authorization-gated evidence tracks.
