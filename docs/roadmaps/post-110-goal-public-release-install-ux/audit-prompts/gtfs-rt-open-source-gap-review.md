# Audit Prompt — GTFS-RT Open Source Gap Review

Use in Phases 120-125.

Review whether the implementation helps a self-hosted small agency produce
useful GTFS-RT, not merely valid endpoints.

Check:

- Vehicle Positions freshness/usefulness;
- Trip Updates emission, withholding, and fallback explanations;
- Alerts lifecycle and cancellation linkage;
- malformed/stale/unknown/ambiguous edge cases;
- validator and conformance coverage;
- connector fixture clarity;
- conservative omission rules;
- no ETA-quality or production-readiness overclaim.
