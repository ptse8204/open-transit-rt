# Independent Install Confidence Policy

Independent install confidence must prove that a non-maintainer can get through
the control plane without project-history context.

## Fresh clone harness

Run from a temporary directory outside the active checkout:

```bash
tmpdir="$(mktemp -d)"
cd "$tmpdir"
git clone https://github.com/ptse8204/open-transit-rt.git
cd open-transit-rt
make check
scripts/bootstrap-dev.sh --check
make agency-app-up
```

Then verify:

```text
/admin/operations
/admin/operations/setup-wizard
/admin/operations/gtfs-import
/admin/operations/gtfs-workbench
/admin/operations/feed-health
/admin/operations/validation-center
/admin/operations/realtime
/admin/operations/connectors/workbench
/admin/operations/maintenance
/admin/operations/help
```

## Release archive replay

After Phase 115, fetch the public tag or release archive and repeat a local
install/startup path.

Record:

- commands;
- environment;
- time to first UI;
- token/admin-login friction;
- Docker/Compose blockers;
- docs confusion;
- UI confusion;
- whether a no-maintainer evaluator can proceed.
