# Security Policy

## Reporting Security Issues

Please do not open a public issue for a suspected vulnerability, leaked secret,
credential, private key, or unsafe evidence artifact.

Use GitHub private vulnerability reporting for this repository if it is enabled:

`https://github.com/ptse8204/open-transit-rt/security/advisories/new`

If that private advisory form is unavailable, contact the maintainer privately
through the GitHub repository owner profile before publishing details.

Include:

- a short description of the issue;
- affected files, URLs, or commands if known;
- whether any credential, token, private key, or private endpoint may be exposed;
- enough reproduction detail to verify the issue without sharing additional
  secrets.

If no private reporting channel is available for a fork or downstream
deployment, contact that repository owner or operator through a private channel
before publishing details.

## Scope

The current public security scope is the Open Transit RT source repository,
documentation, committed evidence packets, deployment scripts, local-development
flows, and public feed publication guidance.

Operator deployments are responsible for their own infrastructure secrets,
runtime environment variables, TLS private keys, DNS provider tokens, database
passwords, admin tokens, device credentials, and monitoring credentials.

## Evidence And Secret Handling

Public evidence may include public feed URLs, validation status, TLS certificate
metadata, public HTTP response headers, redacted operator summaries, and
checksums for committed artifacts.

Public evidence must not include raw credentials, bearer tokens, admin URLs with
embedded secrets, private SSH paths, unredacted IP logs, private keys, database
passwords, or internal hostnames unless the repository explicitly documents why
that detail is public-safe.

If a real secret is found, remove or redact it, rotate or revoke the affected
credential, and assess whether git history cleanup is required. Do not rewrite
git history in this repository without explicit maintainer approval.
