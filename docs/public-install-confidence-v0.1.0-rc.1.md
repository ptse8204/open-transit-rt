# Public Install Confidence -- `v0.1.0-rc.1`

Status date: 2026-05-16

This artifact records Phase 117 independent public fresh-clone install
confidence for Open Transit RT `v0.1.0-rc.1`.

It is a local install-confidence diagnostic only. It is not retained evidence,
release publication, production readiness, compliance proof, consumer
acceptance, agency approval, hosted service availability, vendor compatibility,
SLA/uptime, or ETA-quality proof.

## Conclusion

`public_fresh_clone_passed_after_validator_install_patch`

The first public fresh-clone trial reached the published tag and passed
`make check`, bootstrap preflight, `make test`, local app startup, and all five
local public feed fetches. It failed `make validate` because a fresh clone did
not have pinned validator tooling installed under its local `.cache`.

The install-confidence harness was patched so validate-enabled trials run
`make validators-install` before `make validate` by default. The public
fresh-clone trial then passed end to end.

## Primary Public Tag Trial

Command:

```bash
INSTALL_CONFIDENCE_MODE=clone \
INSTALL_CONFIDENCE_SOURCE=https://github.com/ptse8204/open-transit-rt.git \
INSTALL_CONFIDENCE_REF=v0.1.0-rc.1 \
INSTALL_CONFIDENCE_RUN_LOCAL_APP=true \
INSTALL_CONFIDENCE_RUN_VALIDATE=true \
INSTALL_CONFIDENCE_RUN_TEST=true \
INSTALL_CONFIDENCE_OUTPUT_DIR=.cache/install-confidence/phase117-public-tag-v0.1.0-rc.1-rerun \
INSTALL_CONFIDENCE_FORCE=true \
scripts/install-confidence.sh
```

Result:

- Output directory:
  `.cache/install-confidence/phase117-public-tag-v0.1.0-rc.1-rerun`
- Generated at: `20260516T032351Z`
- Source: `https://github.com/ptse8204/open-transit-rt.git`
- Ref: `v0.1.0-rc.1`
- Checked-out commit: `497f99a97baff630af147c83a7e1249bb08e32da`
- Describe: `v0.1.0-rc.1`
- Overall status: `passed`
- Go: `go version go1.26.2 darwin/amd64`
- Docker: `Docker version 29.4.2, build 055a478`
- Docker Compose: `Docker Compose version v5.1.3`

Steps:

| Step | Status |
| --- | --- |
| `git_clone` | passed |
| `git_checkout` | passed |
| `make-check` | passed |
| `bootstrap-check` | passed |
| `validators-install` | passed |
| `make-validate` | passed |
| `make-test` | passed |
| `agency-app-up` | passed |
| `fetch_feeds_json` | passed |
| `fetch_schedule_zip` | passed |
| `fetch_vehicle_positions_pb` | passed |
| `fetch_trip_updates_pb` | passed |
| `fetch_alerts_pb` | passed |

Fetched local feed artifact SHA-256 values:

| Artifact | SHA-256 |
| --- | --- |
| `feeds_json` | `f00e9578ddf9e2f56ec43142728c2e4f797e3e6900de63f4880e5401fcb20372` |
| `schedule_zip` | `2bd426d0e2c2c81ddbf8d1b8d856d45674a66c2fb09019efed0922bbd265f541` |
| `vehicle_positions_pb` | `fdeb603e30dd97e433905a4d69533d84c3c0c0758462c090f1ade318f85629ed` |
| `trip_updates_pb` | `fdeb603e30dd97e433905a4d69533d84c3c0c0758462c090f1ade318f85629ed` |
| `alerts_pb` | `1075870153fb2cf18d9994f183e6faaca0899c1379bd1e6f4ba29f92779a9bbd` |

## First Attempt Blocker

Output directory:
`.cache/install-confidence/phase117-public-tag-v0.1.0-rc.1`

The first attempt failed overall only because `make validate` reported missing
pinned validator tooling in the fresh clone:

```text
missing pinned tooling: static GTFS validator not installed at .../.cache/validators/gtfs-validator-7.1.0-cli.jar; run make validators-install
```

Other first-attempt steps passed:

- `git_clone`
- `git_checkout`
- `make-check`
- `bootstrap-check`
- `make-test`
- `agency-app-up`
- all five local public feed fetches

## Harness Patch

`scripts/install-confidence.sh` now has:

- `INSTALL_CONFIDENCE_RUN_VALIDATORS_INSTALL`
- default value matching `INSTALL_CONFIDENCE_RUN_VALIDATE`
- an automatic `make validators-install` step before `make validate` when
  validation is enabled
- summary/status output for the validators-install step

This keeps the default light path unchanged when validation is not requested,
while making validation-enabled public fresh-clone trials self-contained.

## Relationship To Phase 116

Phase 116 recorded that the published rc1 source archive extraction still
fails `make check` because protected consumer tracker state is intentionally
excluded from public source archives. The Phase 117 public fresh-clone path
does not have that blocker because a git clone of the tag includes tracked
protected files in the local checkout.

## Claim And Boundary Notes

- No protected evidence path was edited, generated, reformatted, or touched.
- `docs/evidence/consumer-submissions/status.json` was not edited.
- The consumer tracker remains exactly seven targets in order and all
  `prepared`.
- This fresh-clone trial does not prove production readiness, compliance,
  adoption, consumer acceptance, hosted service availability, vendor
  compatibility, hardware certification, SLA/uptime, or ETA quality.
