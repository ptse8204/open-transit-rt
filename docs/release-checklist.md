# Release Checklist

This checklist is for maintainers cutting an Open Transit RT release. It is not
a hosted-service, paid-support, SLA, consumer-acceptance, agency-endorsement, or
compliance promise.

## Release Source

Releases are cut from `main` using git tags. Release branches are not used
unless a later release process explicitly documents them.

Before tagging:

```bash
git checkout main
git pull --ff-only
git status --short
git describe --tags --always --dirty
git rev-parse HEAD
```

The working tree should be clean before tagging. If it is not clean, do not
produce release artifacts from that checkout.

## Version And Pinning

Choose a tag such as `v0.22.0` and record:

- source tag, for example `v0.22.0`;
- commit SHA from `git rev-parse HEAD`;
- dirty/clean state from `git describe --tags --always --dirty`;
- local Docker image tag, if built;
- artifact checksum, if generated.

Operators pin releases by checking out the source tag or exact commit SHA. If
they build a local image, they should tag it with the release tag and optionally
the short commit SHA:

```bash
docker build -f deploy/Dockerfile.local \
  -t open-transit-rt-local:v0.22.0 \
  -t open-transit-rt-local:67e6c95 \
  .
```

Published/versioned production Docker images are deferred. Current distribution
guidance supports source tags and local Docker builds only.

## Required Final Checks

Run and record all required checks:

```bash
make validate
make test
make realtime-quality
make smoke
docker compose -f deploy/docker-compose.yml config
git diff --check
```

If a check cannot run, record the exact command, failure reason, and whether the
release is blocked.

For releases that affect operations docs, evidence packets, deployment helpers,
or runbooks, also perform a context-aware scan of changed docs for:

- tokens or bearer values;
- DB URLs with passwords;
- private keys or TLS/ACME private material;
- webhook URLs or notification credentials;
- sensitive private backup paths;
- raw logs or private operator artifacts;
- unsupported claims such as paid support/SLA coverage, hosted SaaS,
  CAL-ITP/Caltrans compliance, consumer acceptance, agency endorsement,
  production multi-tenant hosting, marketplace/vendor equivalence, or
  production-grade ETA quality.

Negated boundary language such as "no SLA coverage" is allowed and should not
be treated as a claim.

## Clean Install Verification

Use a clean checkout of the tag:

```bash
git clone https://github.com/ptse8204/open-transit-rt.git
cd open-transit-rt
git checkout v0.22.0
git describe --tags --always --dirty
```

Minimum install proof for the local app package:

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

Minimum local image proof:

```bash
docker build -f deploy/Dockerfile.local -t open-transit-rt-local:v0.22.0 .
docker image inspect open-transit-rt-local:v0.22.0 >/tmp/open-transit-image.json
docker run --rm open-transit-rt-local:v0.22.0 sh -c '/app/bin/migrate 2>&1 | grep -q "usage: migrate"'
```

The image proof confirms a built binary exists in the image. Full runtime proof
comes from `make agency-app-up`.

## Release Notes

Prepare notes from `docs/release-notes-template.md`. The notes must explicitly
state `None` for migrations, operations changes, security notes, dependency
changes, or evidence/claim changes when none apply.

Release notes that mention evidence must follow:

- `SECURITY.md`;
- `docs/evidence/redaction-policy.md`;
- `docs/prompts/calitp-truthfulness.md`.

For operations-impacting releases, release notes should record:

- whether migrations are required;
- whether a pre-upgrade backup is required and where private backup evidence is retained;
- migration status before and after upgrade;
- validator and public feed verification after upgrade;
- rollback limits, especially irreversible or untested migrations;
- evidence packet version linkage;
- any secret rotation, incident, restore, or handover docs changed.

## Tagging

After checks pass and notes are ready:

```bash
git tag -a v0.22.0 -m "Open Transit RT v0.22.0"
git show v0.22.0 --stat
```

Push only after confirming the tag points at the intended clean `main` commit:

```bash
git push origin v0.22.0
```

After tagging, update phase handoff/status docs only when the release closes a
documented phase.
