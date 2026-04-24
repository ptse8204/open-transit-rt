# Hosted Monitoring and Alerting Evidence

- Environment: `oci-pilot`
- Capture date (UTC): 2026-04-24
- Operator: Codex operator session using OCI pilot admin credentials

## Metrics And Alerting Artifacts

- Monitor job definition, timer, service status, and feed monitor history: `artifacts/operator-supplied/monitoring-feed-monitor-history.txt`
- Controlled alert lifecycle drill: `artifacts/operator-supplied/alert-lifecycle-controlled-drill.txt`

## Summary

The pilot has a systemd timer for public feed availability monitoring. The controlled alert lifecycle drill used the same availability rule, intentionally probed an invalid path to create a detected alert state, acknowledged it, mitigated by confirming all configured public feed paths returned 200, and resolved the drill.
