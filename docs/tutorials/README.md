# Tutorials

These tutorials describe the current Phase 9/Phase 10 repository surface. They are intentionally evidence-bounded: they document what the repo can run today and do not claim hosted production readiness, consumer acceptance, or completed CAL-ITP/Caltrans compliance.

Start here:

- [Local Quickstart](local-quickstart.md)
- [Agency Demo Flow](agency-demo-flow.md)
- [Deploy With Docker Compose](deploy-with-docker-compose.md)
- [Production Checklist](production-checklist.md)
- [CAL-ITP Readiness Checklist](calitp-readiness-checklist.md)

For the detailed Phase 11 evidence separation, see [Compliance Evidence Checklist](../compliance-evidence-checklist.md).

Rules for future edits:

- Every command must be runnable from the committed repo or clearly marked as deployment-specific.
- Every endpoint and environment variable must match the actual codebase.
- Public protobuf endpoints are anonymous; JSON debug, admin, and GTFS Studio routes are protected.
- Use "supports" and "technical foundations for" when describing compliance readiness unless deployment and external-consumer evidence supports stronger wording.
