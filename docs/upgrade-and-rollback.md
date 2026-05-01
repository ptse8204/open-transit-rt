# Upgrade And Rollback

This guide covers source-tag and local Docker build upgrades. It does not claim
published production Docker images, hosted SaaS availability, paid support, SLA
coverage, consumer acceptance, agency endorsement, marketplace/vendor
equivalence, or CAL-ITP/Caltrans compliance.

## Clean Install From A Source Tag

Use tags as the normal install anchor:

```bash
git clone https://github.com/ptse8204/open-transit-rt.git
cd open-transit-rt
git checkout v0.22.0
git describe --tags --always --dirty
git rev-parse HEAD
```

Verify the local app package:

```bash
make validate
make test
make realtime-quality
make agency-app-up
curl -fsS http://localhost:8080/public/feeds.json >/tmp/open-transit-feeds.json
curl -fsS http://localhost:8080/public/gtfs/schedule.zip >/tmp/open-transit-schedule.zip
curl -fsS http://localhost:8080/public/gtfsrt/vehicle_positions.pb >/tmp/open-transit-vp.pb
make agency-app-down
```

Build a local Docker image from the same tag:

```bash
docker build -f deploy/Dockerfile.local -t open-transit-rt-local:v0.22.0 .
docker image inspect open-transit-rt-local:v0.22.0 >/tmp/open-transit-image.json
```

Use local Docker image tags for local evaluation only. Published/versioned
production Docker images are deferred.

## Version Pinning

Operators should record all of the following for each install or upgrade:

- source tag, such as `v0.22.0`;
- commit SHA from `git rev-parse HEAD`;
- dirty/clean state from `git describe --tags --always --dirty`;
- local Docker image tag, if built;
- release artifact checksum, if a tarball, binary bundle, or image archive is generated.

If an artifact file is generated, produce and retain a checksum:

```bash
shasum -a 256 path/to/artifact > path/to/artifact.sha256
```

## Backup Before Upgrade

Always take a database backup before changing source tags, binaries, images, or
migrations.

For the local Compose app, the simplest rollback path is a destructive reset:

```bash
make agency-app-down
make agency-app-reset
```

For pilot or deployment-owned environments, use the existing pilot operations
backup and restore-drill workflow:

- `docs/runbooks/small-agency-pilot-operations.md`
- `scripts/pilot-ops.sh backup --dry-run`
- `scripts/pilot-ops.sh restore-drill --dry-run`

Raw backups, database URLs with passwords, admin tokens, private keys, and
operator artifacts must not be committed as public evidence.

## Migration Run Order

Before upgrading, check the current migration state against the old version:

```bash
make migrate-status
```

Upgrade order:

1. Back up the database.
2. Record current version pinning information.
3. Stop application services or put the deployment into a maintenance window.
4. Check out the new source tag or deploy the new locally built artifacts.
5. Review release notes for migrations and rollback limits.
6. Run migrations.
7. Check migration status again.
8. Start services.
9. Run install verification checks.

For Makefile-based deployments:

```bash
make migrate-status
make migrate-up
make migrate-status
```

For compiled deployment artifacts, run the deployed `migrate` binary with the
deployment-owned environment file and `MIGRATIONS_DIR` that matches the release:

```bash
DATABASE_URL="$DATABASE_URL" MIGRATIONS_DIR=/opt/open-transit-rt/app/db/migrations /opt/open-transit-rt/bin/migrate status
DATABASE_URL="$DATABASE_URL" MIGRATIONS_DIR=/opt/open-transit-rt/app/db/migrations /opt/open-transit-rt/bin/migrate up
DATABASE_URL="$DATABASE_URL" MIGRATIONS_DIR=/opt/open-transit-rt/app/db/migrations /opt/open-transit-rt/bin/migrate status
```

## Rollback Limits

Rollback is straightforward only when no migration changed persistent state, or
when the release notes explicitly say the migration is reversible and tested.

If a release includes irreversible or untested migrations:

- do not rely on `make migrate-down` as the recovery plan;
- restore the pre-upgrade database backup;
- redeploy the previous source tag or artifacts;
- verify public feed URLs still use the same canonical paths.

`make migrate-down` rolls back one migration according to the migration files,
but it is not a substitute for backup/restore evidence.

## Restore Procedure Links

Use existing restore guidance instead of inventing ad hoc rollback steps:

- [Small-Agency Pilot Operations](runbooks/small-agency-pilot-operations.md)
- `scripts/pilot-ops.sh restore-drill --dry-run`
- `docs/evidence/redaction-policy.md`
- `SECURITY.md`

After restore, rerun the minimum verification commands:

```bash
make validate
make test
make realtime-quality
make smoke
docker compose -f deploy/docker-compose.yml config
git diff --check
```

For a local app package restore/reset, also rerun:

```bash
make agency-app-up
curl -fsS http://localhost:8080/public/feeds.json >/tmp/open-transit-feeds.json
make agency-app-down
```

## Evidence Packet Version Linkage

Future evidence packets should record:

- git tag;
- git commit SHA;
- dirty/clean state;
- release notes link;
- artifact checksums where available.

Evidence packet version metadata is supporting context only. It does not prove
consumer acceptance, compliance, hosted SaaS availability, agency endorsement,
or vendor equivalence.
