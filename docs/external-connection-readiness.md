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

## Claim Boundary

External connection readiness is a review of local quality controls. It does
not prove CAL-ITP/Caltrans compliance, consumer acceptance, agency adoption or
approval, final-root readiness, hosted SaaS, paid support/SLA, production
readiness, production multi-tenant hosting, vendor compatibility, hardware
certification, public launch, production AVL reliability, or production-grade
ETA quality.

Validator success is a supporting signal only. It is not compliance proof,
consumer acceptance, or a correctness guarantee.
