# Public GTFS Local/Pilot Evidence Templates

Status: template only. Do not fill this directory with fake run evidence.

Use these templates for Phase 33 public GTFS local/pilot evidence. Actual
attempted-run evidence belongs in a dated packet:

```text
docs/evidence/captured/public-gtfs-local-pilot/YYYY-MM-DD/
```

Outcome B packets should include a blocker summary and command/log inventory
only. Outcome C packets should include public-safe retained summaries,
checksums, fetch summaries, validator summaries or blockers, telemetry
simulator/dry-run summary if applicable, admin/private boundary check, and a
claim-boundary README.

Do not commit raw GTFS ZIP files here. Store raw public GTFS downloads under
`.cache/` or another ignored local path by default.

Current templates:

- `public-gtfs-local-pilot-packet-template.md`
