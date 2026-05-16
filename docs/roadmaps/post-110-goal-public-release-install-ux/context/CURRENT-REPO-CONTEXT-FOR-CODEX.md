# Current Repo Context For Codex

## State to preserve

- Phase 110 is complete for long-term extensibility and plugin governance.
- The authorized Phase 91-110 post-90 agency-grade GTFS-RT roadmap is closed.
- Phase 95 generated and audited a local `.cache` `v0.1.0-rc.1` candidate package.
- Phase 108 reran post-RC stabilization and kept the release-candidate state bounded.
- There is still no public release tag, GitHub Release, package publication, or published image unless a later phase creates one.
- Optional evidence tracks remain authorization-gated.
- All seven consumer targets must remain exactly `prepared`.

## Protected paths

Do not modify or generate:

```text
docs/evidence/captured/**
docs/evidence/consumer-submissions/status.json
docs/evidence/consumer-submissions/current/**
docs/evidence/consumer-submissions/artifacts/**
docs/evidence/consumer-submissions/packets/**
```

## Public release principle

A public `v0.1.0-rc.1` release candidate may be published only after release
gates pass. The release must be described as a public release candidate for
local/self-hosted evaluation only. It is not production readiness, compliance,
consumer acceptance, agency adoption, hosted service availability, SLA/uptime,
vendor compatibility, hardware certification, final-root readiness, or
production-grade ETA quality.
