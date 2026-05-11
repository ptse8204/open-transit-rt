# External Connection Readiness

This readiness guide helps maintainers evaluate optional connectors without
making unsupported external-integration claims.

## Default Path

1. Review the connector manifest against
   [Connector Plugin Contract](connectors/plugin-contract.md).
2. Run local manifest checks:

   ```sh
   make external-connection-check
   ```

3. Run the adapter conformance suite when available for the connector type.
4. Keep diagnostics under ignored local storage or the connector's own private
   operator environment.
5. Update public wording only when retained, claim-specific evidence supports
   the exact claim.

## External-Connection Readiness Summary

| Surface | Ready shape | Gate |
| --- | --- | --- |
| AVL/device input | Adapter, sidecar, script, or private process transforms observations before calling authenticated `POST /v1/telemetry`. | Reject malformed, stale, future-dated, wrong-agency, unknown-device, low-quality, duplicate, or out-of-order observations. |
| Vehicle Positions | Vehicle Positions continue without any external predictor. | Emit optional fields only behind reliability, freshness, and confidence gates. |
| Trip Updates | Deterministic prediction remains the safe fallback. | External predictors run only behind `internal/prediction.Adapter` and are tested in shadow or fail-closed modes. |
| Validator tooling | Server-owned allowlisted validator IDs, pinned tooling checks, and private validator-health summaries. | No arbitrary command execution from manifests or browsers. |
| Monitoring/export | Redacted summaries and optional export surfaces. | No-send by default unless deployment-owned configuration explicitly enables delivery. |
| Feed-consumer URL/metadata expectations | Public feed URLs and `/public/feeds.json` metadata can be reviewed locally. | No portal automation, no submissions, and no target status changes. |
| Redaction | Fixtures, manifests, logs, screenshots, and summaries avoid secrets and private identifiers. | Review before committing any integration material. |

## Readiness Questions

- Is the connector optional and disabled unless the operator configures it?
- Does it fail closed when upstream data is malformed, stale, future-dated,
  wrong-agency, unknown-device, low-quality, duplicate, or out of order?
- Does it avoid raw secrets, private paths, raw private payloads, and private
  deployment URLs in manifests, fixtures, logs, and summaries?
- Does Vehicle Positions continue without an external predictor?
- Does deterministic prediction remain the default?
- Are validator calls allowlisted rather than arbitrary commands?
- Are monitoring exports no-send by default and redacted?
- Are consumer/discovery workflows free of portal automation and status
  mutation?
- Are product screenshots and diagrams clearly labeled as local/demo
  documentation aids rather than evidence?

## Vehicle Positions And Trip Updates Gates

Vehicle Positions are the first high-quality public realtime output. Optional
Vehicle Positions fields should stay omitted until freshness, reliability,
assignment confidence, and downstream-safety checks justify emitting them.

Trip Updates remain pluggable. Deterministic prediction is the safe fallback.
External predictor adapters must be tested in shadow or fail-closed modes with
sanitized DTOs, bounded timeouts, redacted diagnostics, active-feed validation,
and rejection of stale, wrong-agency, impossible-stop, low-confidence, or
unsupported output.

> **What this proves:** the integration boundary is safer to evaluate with
> synthetic/local data.
>
> **What this does not prove:** named vendor compatibility, hardware
> certification, production AVL reliability, production-grade ETA quality,
> consumer acceptance, production readiness, or compliance.

## Claim Boundary

External connection readiness is a review of local quality controls. It does
not prove CAL-ITP/Caltrans compliance, consumer acceptance, agency adoption or
approval, final-root readiness, hosted SaaS, paid support/SLA, production
readiness, production multi-tenant hosting, vendor compatibility, hardware
certification, public launch, production AVL reliability, or production-grade
ETA quality.

Validator success is a supporting signal only. It is not compliance proof,
consumer acceptance, or a correctness guarantee.
