# Public Feed Root and TLS Proof

- Environment: `local-demo`
- Capture date (UTC): 2026-04-22
- Operator: Codex local run
- Canonical HTTPS host: Missing. No hosted HTTPS deployment hostname was configured or provided.
- Local feed root used for evidence: `http://localhost:8090`

## Evidence Category

- Repo-proven capability: public feed paths exist and return artifacts through a local loopback proxy.
- Deployment/operator proof: not complete because this packet does not include public HTTPS hosting, DNS, TLS, or production proxy evidence.
- Third-party proof: none.

## Stable Public Paths

Anonymous local HTTP fetches were captured at 2026-04-22T21:36:10Z. Raw artifacts are stored under `artifacts/public/`.

| Path | Status | Content type | Bytes | SHA-256 |
| --- | --- | --- | ---: | --- |
| `/public/gtfs/schedule.zip` | 200 | `application/zip` | 1960 | `0956ed037a40ca9d2cca94a501bea1547d27dbd25a195c7ebefe1a34ffc78194` |
| `/public/feeds.json` | 200 | `application/json` | 2408 | `4c274d137dddfa1d65b5e5c8c1ebbf15ee6f385fa8298fa32baa7495e13b5ac3` |
| `/public/gtfsrt/vehicle_positions.pb` | 200 | `application/x-protobuf` | 63 | `6b7033acef5a560e5bd6480009b747b905097d5561a64427c45f990eb85496f3` |
| `/public/gtfsrt/trip_updates.pb` | 200 | `application/x-protobuf` | 15 | `74629d6bafb7f50c9ed5d68e5af20ba83047eb7d24bd2df9d24c43bf9a322279` |
| `/public/gtfsrt/alerts.pb` | 200 | `application/x-protobuf` | 135 | `564f35b5f5b8791b7dc59c0e40487298486ce339285a8af7e9f0781244b26b02` |

### Key Headers

`/public/gtfs/schedule.zip` returned:

```text
HTTP/1.1 200 OK
Content-Length: 1960
Content-Type: application/zip
Date: Wed, 22 Apr 2026 21:36:10 GMT
Etag: "gtfs-import-7-0956ed037a40ca9d2cca94a501bea1547d27dbd25a195c7ebefe1a34ffc78194"
Last-Modified: Wed, 22 Apr 2026 21:33:47 GMT
X-Checksum-Sha256: 0956ed037a40ca9d2cca94a501bea1547d27dbd25a195c7ebefe1a34ffc78194
```

Realtime protobuf paths returned `Content-Type: application/x-protobuf`, `HTTP/1.1 200 OK`, and `Last-Modified: Wed, 22 Apr 2026 21:36:10 GMT`.

## Publish / Rollback URL Stability

- Before publish proof: Missing. No pre-publish fetch was captured in this pass.
- After publish proof: local demo imported `testdata/gtfs/valid-small`, published feed version `gtfs-import-7`, and the local fetches above prove the local URLs served artifacts after import and publication metadata bootstrap.
- After rollback proof: Missing. No rollback operation was performed in this pass.
- URL changed? Not proven for publish/rollback. Local URLs stayed constant during this capture only.

## Reverse Proxy / TLS

- Routing map reference: `reverse-proxy-tls-2026-04-22.md`.
- TLS termination owner: Missing for this environment.
- Certificate issuer + validity: Missing for this environment.
- Renewal check process: Missing for this environment.
- HTTP-to-HTTPS behavior: Missing for this environment. The local proxy is HTTP only.

## Admin/Debug Protection Boundary

Anonymous checks at 2026-04-22T21:36:10Z returned:

| URL | Status |
| --- | ---: |
| `http://localhost:8081/admin/compliance/scorecard` | 401 |
| `http://localhost:8082/v1/events?limit=10` | 401 |
| `http://localhost:8083/public/gtfsrt/vehicle_positions.json` | 401 |
| `http://localhost:8084/public/gtfsrt/trip_updates.json` | 401 |
| `http://localhost:8085/public/gtfsrt/alerts.json` | 401 |
| `http://localhost:8086/admin/gtfs-studio` | 401 |

## Notes / Gaps

- Pending: real public HTTPS hostname, DNS ownership, TLS certificate, redirect proof, production routing config, and publish/rollback URL permanence.
- Known blocker: no deployment hostname or TLS-terminating proxy configuration was available in the workspace.
