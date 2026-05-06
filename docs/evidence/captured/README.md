# Captured Deployment Evidence (Operator-Owned)

Add real deployment artifacts here when available.

Current captured packet:

- `local-demo/2026-04-22/`: real local demo evidence packet. This is useful repo/operator evidence for the local loopback environment, but it is not hosted HTTPS deployment proof.
- `hosted-pending/2026-04-22/`: hosted evidence intake packet. It contains collection commands and pending artifact slots, not completed hosted proof.
- `oci-pilot/2026-04-24/`: hosted OCI pilot packet for the recorded pilot scope. It is deployment/operator proof, not CAL-ITP compliance or consumer acceptance.
- `public-gtfs-local-pilot/2026-05-06/`: Phase 33 LA Metro Bus GTFS local/pilot packet. It is Outcome C local/pilot evidence for public GTFS import, publication, five-path fetch, schedule proof, and retained public-safe summaries. It is not agency adoption, final-root proof, consumer evidence, compliance proof, production readiness, or real realtime/ETA-quality proof.
- `public-gtfs-local-pilot/templates/`: reusable templates for future public GTFS local/pilot packets. Do not store actual run evidence in the templates directory.

Use one directory per environment, for example:

- `pilot-agency-prod/`
- `staging/`

If full raw artifacts cannot be committed, add a redacted summary plus a reference to secure storage.

For Phase 17 evidence refresh, do not call a packet complete until:

```sh
EVIDENCE_PACKET_DIR=docs/evidence/captured/<environment>/<UTC-date> make audit-hosted-evidence
```

passes.
