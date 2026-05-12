# Product Diagrams

These diagrams make the technical docs easier to read. They are documentation
aids only, not retained evidence and not proof of production readiness,
CAL-ITP/Caltrans compliance, agency adoption or approval, consumer submission
or acceptance, final-root readiness, hosted SaaS availability, SLA coverage,
vendor compatibility, hardware certification, or production-grade ETA quality.

Use Markdown or Mermaid so the source stays reviewable.

## Self-Hosted Operator Flow

```mermaid
flowchart TD
    A["Clean checkout"] --> B["make check"]
    B --> C["make agency-app-up"]
    C --> D["Operations Console: Agency Operations Cockpit / Start Here"]
    D --> E["Import demo or public-safe GTFS"]
    E --> F["Review five public feed paths"]
    F --> G["Review feed health, GTFS quality, and validator health"]
    G --> H["Run synthetic telemetry or adapter dry run"]
    H --> I["Review readiness and claim boundaries"]
    I --> J["Record blockers or local next actions"]
```

**What this proves:** a local/demo operator can follow the product workflow and
find the right private UI surfaces.

**What this does not prove:** production readiness, agency approval, final-root
readiness, consumer acceptance, compliance, hosted SaaS, SLA coverage, vendor
compatibility, or production-grade ETA quality.

## No-CLI Agency Operations Flow

```mermaid
flowchart TD
    A["Private /admin/operations"] --> B["Agency Operations Cockpit"]
    B --> C["Import or review GTFS"]
    B --> D["Review five public feed paths"]
    B --> E["Review GTFS quality and validators"]
    B --> F["Review devices, telemetry, and realtime usefulness"]
    B --> G["Review connector readiness"]
    B --> H["Review Maintenance Center"]
    H --> I["Weekly/monthly next action"]
    I --> J["Technical helper only when deployment, validator, token, or diagnostic shell work is needed"]
```

**What this proves:** a local/demo operator has a browser-first path for
routine setup, feed review, quality review, realtime usefulness, connector
review, and maintenance decisions.

**What this does not prove:** agency adoption, compliance, consumer
acceptance, final-root readiness, hosted service availability, production
readiness, vendor compatibility, SLA coverage, or production-grade ETA
quality.

## External-Connection Flow

```mermaid
flowchart LR
    A["AVL, GPS, CSV, or device source"] --> B["Operator-owned adapter or sidecar"]
    B --> C["Quality, freshness, agency, and redaction checks"]
    C --> D["Authenticated POST /v1/telemetry"]
    D --> E["Open Transit RT telemetry state"]
    E --> F["Conservative matcher"]
    F --> G["Vehicle Positions feed"]
    E --> H["Prediction adapter boundary"]
    H --> I["Deterministic predictor fallback"]
    H --> J["Optional external predictor shadow/fail-closed mode"]
    I --> K["Trip Updates feed"]
    J --> K
    E --> L["Validator, monitoring, and readiness surfaces"]
```

**What this proves:** connectors have a documented local boundary and a safer
evaluation path.

**What this does not prove:** named vendor compatibility, hardware
certification, real AVL reliability, consumer acceptance, production readiness,
or production-grade ETA quality.

## Release-Candidate Readiness Flow

```mermaid
flowchart TD
    A["Clean checkout"] --> B["Source status: git status and commit SHA"]
    B --> C["make check"]
    C --> D["make validate"]
    D --> E["make test"]
    E --> F["Local app startup"]
    F --> G["Five public feed fetches"]
    G --> H["Public GTFS trial when allowed"]
    H --> I["Validator health"]
    I --> J["Telemetry simulator"]
    J --> K["make external-connection-check"]
    K --> L["make adapter-conformance"]
    L --> M["make audit-final-claim-review"]
    M --> N["v0.1.0-rc.1 review decision"]
    N --> O["Full v0.1.0 only after RC blockers are resolved"]
```

**What this proves:** a maintainer has a repeatable local release-candidate
diagnostic path and blocker list.

**What this does not prove:** a release has been published, production
readiness, compliance, consumer status, agency adoption or approval,
agency-owned final-root readiness, hosted SaaS availability, SLA coverage,
vendor compatibility, or production-grade ETA quality.

## Review Rules

- Keep diagrams product- and workflow-focused.
- Do not call diagrams evidence.
- Do not add agency logos, vendor logos, real transit photos, private URLs,
  private identifiers, emails, IPs, tokens, or secrets.
- Link to the canonical
  [Review And Recommendations](../../roadmap-status.md#review-and-recommendations)
  instead of duplicating the full scorecard in every doc.
