# Captured Deployment Evidence (Operator-Owned)

Add real deployment artifacts here when available.

Current captured packet:

- `local-demo/2026-04-22/`: real local demo evidence packet. This is useful repo/operator evidence for the local loopback environment, but it is not hosted HTTPS deployment proof.
- `hosted-pending/2026-04-22/`: hosted evidence intake packet. It contains collection commands and pending artifact slots, not completed hosted proof.
- `oci-pilot/2026-04-24/`: hosted OCI pilot packet for the recorded pilot scope. It is deployment/operator proof, not CAL-ITP compliance or consumer acceptance.

Use one directory per environment, for example:

- `pilot-agency-prod/`
- `staging/`

If full raw artifacts cannot be committed, add a redacted summary plus a reference to secure storage.

For Phase 17 evidence refresh, do not call a packet complete until:

```sh
EVIDENCE_PACKET_DIR=docs/evidence/captured/<environment>/<UTC-date> make audit-hosted-evidence
```

passes.
