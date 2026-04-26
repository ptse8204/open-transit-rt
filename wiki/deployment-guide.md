# Deployment Guide

This guide describes the current small-agency pilot shape. It is not a complete hosted production package.

![Agency deployment diagram showing public internet through TLS reverse proxy, anonymous public feed paths, protected admin/debug paths, Go services, Postgres/PostGIS, pinned validators, and optional predictor adapter boundary.](assets/agency-deployment.png)

*Exact-behavior deployment-boundary diagram rendered from a reviewed spec.*

## Deployment Shape

A pilot deployment should provide:

- PostgreSQL with PostGIS
- the Open Transit RT Go services
- a TLS reverse proxy
- stable public feed paths
- protected admin and debug routes
- real production secrets
- pinned validator tooling or equivalent pinned validator artifacts

## Public Paths

Expose these paths anonymously over stable HTTPS:

```text
/public/gtfs/schedule.zip
/public/feeds.json
/public/gtfsrt/vehicle_positions.pb
/public/gtfsrt/trip_updates.pb
/public/gtfsrt/alerts.pb
```

Keep admin/debug/JSON routes behind admin auth and deployment network controls.

## Useful Internal References

- [Deploy With Docker Compose](../docs/tutorials/deploy-with-docker-compose.md)
- [Production Checklist](../docs/tutorials/production-checklist.md)
- [Reverse Proxy And TLS Runbook](../docs/runbooks/reverse-proxy-and-tls.md)
- [Monitoring And Alerting Runbook](../docs/runbooks/monitoring-and-alerting.md)
