# Evidence Redaction Policy

This policy applies to public evidence under `docs/evidence/`, especially
captured deployment packets.

## Public-Safe Evidence

Public evidence may include:

- public hostnames and public feed URLs;
- public HTTP status codes, headers, and checksums;
- GTFS and GTFS-RT validation status;
- TLS certificate metadata and redirect behavior;
- public feed discovery metadata;
- redacted operational summaries;
- loopback-only local demo logs when they do not include credentials or personal
  data.

## Redact Or Exclude

Do not commit:

- raw credentials, bearer tokens, JWTs, API keys, DNS provider tokens, database
  passwords, webhook URLs, notification credentials, or device tokens;
- private keys, private certificates, ACME account keys, SSH keys, or key paths
  that reveal private operator layouts;
- admin URLs containing embedded secrets;
- unredacted public client IP logs or raw access logs unless the packet
  explicitly justifies why the values are public-safe;
- internal hostnames, instance names, private infrastructure inventory, or
  private origin URLs unless the packet explicitly justifies why the values are
  public-safe;
- private correspondence, consumer portal credentials, private ticket links, or
  personal data.

Use placeholders such as `<redacted-admin-token>`,
`<redacted-instance-host>`, `<redacted-client-ip>`,
`<redacted-private-origin>`, and `<redacted-webhook-url>` when the surrounding
context is still useful.

## Checksums And Inventories

Every committed evidence archive must have an inventory entry that lists its
path, purpose, checksum, contents, and keep/remove decision. Opaque archives
must not be kept.

If any file under `docs/evidence/captured/**/artifacts/` changes, refresh the
relevant `SHA256SUMS.txt` or per-file checksum artifact and update markdown that
references changed hashes, filenames, or contents.

Template files under `docs/runbooks/templates/` are not evidence. Do not fill
them with fake incidents, fake alert delivery proof, fake rotation records, fake
restore events, or placeholder operational artifacts.

For the OCI pilot packet, rerun:

```bash
EVIDENCE_PACKET_DIR=docs/evidence/captured/oci-pilot/2026-04-24 make audit-hosted-evidence
```

## Secret Response

If a real secret is found:

1. Redact or remove it from the working tree.
2. Document the required rotation or revocation.
3. Check whether the file was committed in git history.
4. Record whether history cleanup may be needed.
5. Do not rewrite git history without explicit maintainer approval.

Deleting a file is not enough when a real secret was exposed. Rotate or revoke
the credential and verify the old value no longer works.
