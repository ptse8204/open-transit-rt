# Audit Prompt — Claim Boundary Review

Scan changed files for forbidden claims.

Run:

```bash
make audit-final-claim-review
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
```

Verify all seven targets remain exactly `prepared`.
