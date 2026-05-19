package main

import (
	"net/http"
	"time"

	"open-transit-rt/internal/auth"
	connectorpkg "open-transit-rt/internal/connectors"
)

const safePluginDefinition = "In Open Transit RT, a plugin is an optional sidecar, command adapter, manifest, or connector process. It is not arbitrary dynamic code loaded into the backend."

type connectorHubView struct {
	GeneratedAt      time.Time              `json:"generated_at"`
	AgencyID         string                 `json:"agency_id"`
	Boundary         string                 `json:"boundary"`
	PluginDefinition string                 `json:"plugin_definition"`
	Catalog          []connectorCatalogRow  `json:"catalog"`
	Categories       []connectorCategory    `json:"categories"`
	Registry         connectorpkg.Registry  `json:"registry"`
	ClaimFlags       connectorHubClaimFlags `json:"claim_flags"`
}

type connectorCatalogRow struct {
	ID             string   `json:"id"`
	Group          string   `json:"group"`
	Label          string   `json:"label"`
	Status         string   `json:"status"`
	StartWith      string   `json:"start_with"`
	BrowserReview  string   `json:"browser_review"`
	FirstSafeCheck string   `json:"first_safe_check"`
	DoesNotProve   string   `json:"does_not_prove"`
	DocsLinks      []string `json:"docs_links"`
}

type connectorCategory struct {
	ID                 string   `json:"id"`
	Label              string   `json:"label"`
	Status             string   `json:"status"`
	Summary            string   `json:"summary"`
	ConnectorShape     string   `json:"connector_shape"`
	Inputs             []string `json:"inputs"`
	Outputs            []string `json:"outputs"`
	FailureBehavior    string   `json:"failure_behavior"`
	CommandSuggestions []string `json:"command_suggestions"`
	AdminLinks         []string `json:"admin_links"`
	DocsLinks          []string `json:"docs_links"`
	ClaimBoundary      string   `json:"claim_boundary"`
}

type connectorHubClaimFlags struct {
	DynamicBackendPluginLoadingEnabled bool `json:"dynamic_backend_plugin_loading_enabled"`
	VendorCompatibilityClaimed         bool `json:"vendor_compatibility_claimed"`
	HardwareCertificationClaimed       bool `json:"hardware_certification_claimed"`
	ConsumerStatusesChanged            bool `json:"consumer_statuses_changed"`
	ExternalEvidenceCreated            bool `json:"external_evidence_created"`
	ComplianceClaimed                  bool `json:"compliance_claimed"`
	ProductionReadinessClaimed         bool `json:"production_readiness_claimed"`
	HostedSaaSClaimed                  bool `json:"hosted_saas_claimed"`
	ProductionGradeETAClaimed          bool `json:"production_grade_eta_claimed"`
}

func (h *handler) renderConnectorHub(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "connectors")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "connectors", page)
}

func (h *handler) renderConnectorHubJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "connectors")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.ConnectorHub)
}

func buildConnectorHub(page operationsPage) connectorHubView {
	categories := []connectorCategory{
		connectorCategoryView(
			"telemetry_source",
			"Vehicle / GPS / AVL connectors",
			"available",
			"Use CSV replay, HTTP polling, webhook sidecars, generic JSON transforms, vendor-shaped synthetic examples, and authenticated telemetry POST paths.",
			"sidecar or command adapter that normalizes observations before authenticated POST /v1/telemetry",
			[]string{"external observations", "agency-owned mapping", "device credential reference"},
			[]string{"Open Transit RT telemetry.Event records", "private redacted adapter diagnostics"},
			"Fail closed before ingest when mapping, credentials, timestamps, or payload quality are unsafe.",
			[]string{"make telemetry-simulator", "make adapter-conformance", "make test-connector-examples"},
			[]string{"/admin/operations/connectors/workbench", "/admin/operations/devices", "/admin/operations/telemetry"},
			[]string{"docs/connectors/catalog.md", "docs/connectors/vehicle-avl-starter-kits.md", "docs/integration-adapter-kit.md", "docs/tutorials/device-avl-integration.md"},
			"Telemetry connectors do not prove real vendor compatibility, hardware certification, production AVL reliability, or agency adoption.",
		),
		connectorCategoryView(
			"prediction",
			"Prediction connectors",
			"available",
			"Review the deterministic built-in predictor, external HTTP predictor adapter, shadow mode, fail-closed behavior, and TheTransitClock candidate boundary.",
			"optional internal adapter or external HTTP sidecar behind internal/prediction.Adapter",
			[]string{"active GTFS feed version", "latest telemetry", "current assignments", "Vehicle Positions feed data or URL"},
			[]string{"Trip Updates feed", "bounded diagnostics"},
			"Return deterministic output or valid empty Trip Updates with adapter diagnostics when an optional sidecar fails.",
			[]string{"make adapter-conformance", "go test ./internal/prediction"},
			[]string{"/admin/operations/feeds", "/admin/operations/reliability"},
			[]string{"docs/connectors/catalog.md", "docs/requirements-trip-updates.md", "docs/integration-adapter-kit.md", "docs/tutorials/external-adapter-conformance.md"},
			"Prediction connectors do not prove production-grade ETA quality, consumer acceptance, or named predictor compatibility.",
		),
		connectorCategoryView(
			"validator",
			"Validator connectors",
			"available",
			"Review MobilityData static GTFS validation, MobilityData GTFS Realtime validation, allowlisted validator IDs, and private validation health.",
			"allowlisted command adapter configured by deployment-owned environment",
			[]string{"server-derived schedule ZIP", "server-derived realtime protobuf", "validator ID"},
			[]string{"normalized validation_report rows", "private validator health diagnostics"},
			"Missing or failing tooling records not-run, missing, warning, failed, or blocked states instead of pretending success.",
			[]string{"make validators-check", "make validator-health"},
			[]string{"/admin/operations/validation-health", "/admin/operations/gtfs-quality"},
			[]string{"docs/connectors/catalog.md", "docs/dependencies.md", "docs/tutorials/gtfs-validation-triage.md", "docs/release-candidate-readiness.md"},
			"Validator records are supporting signals only; they are not CAL-ITP/Caltrans compliance or consumer acceptance.",
		),
		connectorCategoryView(
			"monitoring_export",
			"Monitoring / export connectors",
			"available",
			"Review local health summaries, operations notify drafts, monitoring/export helpers, and the deployment-owned monitoring boundary.",
			"deployment-owned no-send exporter or redacted summary writer",
			[]string{"validator health summary", "deployment doctor summary", "reliability summary"},
			[]string{"private redacted monitoring/export summaries"},
			"Default examples write local private summaries only and require deployment-owned configuration before any external send.",
			[]string{"make operations-reliability", "make operations-notify"},
			[]string{"/admin/operations/reliability", "/admin/operations/validation-health"},
			[]string{"docs/connectors/catalog.md", "docs/tutorials/self-hosted-operations-notifications.md", "docs/deployment/reference-deployment-doctor.md"},
			"Monitoring connectors do not prove SLA coverage, uptime guarantees, hosted service availability, or production readiness.",
		),
		connectorCategoryView(
			"consumer_discovery",
			"Consumer / discovery connectors",
			"available",
			"Review /public/feeds.json, static GTFS, Vehicle Positions, Trip Updates, Alerts, and consumer packet preparedness without submission or acceptance claims.",
			"docs workflow and packet generator, not portal automation",
			[]string{"stable feed URL metadata", "license/contact metadata", "validation status"},
			[]string{"prepared consumer packet records", "private readiness summaries"},
			"Do not contact targets or change statuses without retained target-originated or operator-retained authorization artifacts.",
			[]string{"make external-connection-check", "make audit-final-claim-review"},
			[]string{"/admin/operations/consumers", "/admin/operations/readiness"},
			[]string{"docs/connectors/catalog.md", "docs/consumer-submission-evidence.md", "docs/evidence/consumer-submissions/README.md", "docs/compliance-evidence-checklist.md"},
			"Prepared packets are not submission, review, acceptance, ingestion, listing, display, compliance, or public launch proof.",
		),
		connectorCategoryView(
			"future_extension_model",
			"Future connector extension model",
			"policy",
			"Use manifest-based sidecars, keep arbitrary dynamic backend plugin loading disabled, and require conformance tests for new connector shapes.",
			"manifest-based sidecar or command adapter with explicit contracts",
			[]string{"connector manifest", "synthetic fixture", "redaction plan", "conformance cases"},
			[]string{"reviewed connector example", "private diagnostics", "bounded docs"},
			"Reject secrets, raw commands, status mutation, external submission automation, unreviewed runtime loading, and unsupported claims.",
			[]string{"make external-connection-check", "make adapter-conformance"},
			[]string{"/admin/operations/connectors/workbench", "/admin/operations/connectors/tests"},
			[]string{"docs/connectors/catalog.md", "docs/connectors/plugin-contract.md", "docs/connectors/contributing-connectors.md", "docs/extension-governance.md"},
			"Extension policy does not show runtime distribution support, third-party compatibility, compliance, consumer acceptance, or production readiness.",
		),
	}
	return connectorHubView{
		GeneratedAt:      page.GeneratedAt,
		AgencyID:         page.AgencyID,
		Boundary:         "Private authenticated Connectors only; viewing it creates no evidence, contacts no external party, changes no consumer status, and records no approval, compatibility, compliance, hosted-service, SLA, production-readiness, or ETA-quality outcome. Treat it as the starting point for manifest, redaction, fail-closed, and adapter-conformance review.",
		PluginDefinition: safePluginDefinition,
		Catalog:          connectorCatalogRows(),
		Categories:       categories,
		Registry:         connectorRegistryForSection(page.Section),
		ClaimFlags:       connectorHubClaimFlags{},
	}
}

func connectorCatalogRows() []connectorCatalogRow {
	return []connectorCatalogRow{
		connectorCatalogRowView(
			"csv_replay_adapter",
			"Vehicle / GPS / AVL",
			"CSV replay adapter",
			"covered",
			"`examples/connectors/telemetry-csv-replay` with synthetic or operator-redacted rows.",
			"/admin/operations/connectors/workbench and /admin/operations/telemetry-simulator",
			"make test-connector-examples",
			"Production data quality, real device proof, vendor compatibility, hardware certification, compliance, consumer acceptance, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/connectors/vehicle-avl-starter-kits.md", "docs/tutorials/telemetry-simulator-and-device-trial.md"},
		),
		connectorCatalogRowView(
			"http_polling_adapter",
			"Vehicle / GPS / AVL",
			"HTTP polling adapter",
			"covered",
			"`examples/connectors/telemetry-http-poller` with committed synthetic observations.",
			"/admin/operations/connectors/workbench",
			"make external-connection-check",
			"Live endpoint behavior, named API support, vendor compatibility, production AVL reliability, compliance, or public launch.",
			[]string{"docs/connectors/catalog.md", "docs/integration-adapter-kit.md", "docs/tutorials/device-avl-integration.md"},
		),
		connectorCatalogRowView(
			"webhook_sidecar_adapter",
			"Vehicle / GPS / AVL",
			"Webhook sidecar adapter",
			"covered",
			"`examples/connectors/telemetry-webhook-sidecar` as a dry-run transform shape.",
			"/admin/operations/connectors/workbench and /admin/operations/devices",
			"go run ./cmd/adapter-conformance telemetry --suite testdata/adapter-conformance",
			"Vendor support, hardware certification, production AVL reliability, SLA, compliance, consumer acceptance, or agency approval.",
			[]string{"docs/connectors/catalog.md", "docs/connectors/vehicle-avl-starter-kits.md", "docs/evidence/redaction-policy.md"},
		),
		connectorCatalogRowView(
			"generic_json_transform_adapter",
			"Vehicle / GPS / AVL",
			"Generic JSON transform adapter",
			"covered",
			"`cmd/avl-vendor-adapter --dry-run` with `testdata/avl-vendor` fixtures.",
			"/admin/operations/connectors/workbench and /admin/operations/realtime",
			"go test ./cmd/avl-vendor-adapter",
			"Vendor compatibility, real AVL reliability, production readiness, consumer acceptance, compliance, or ETA quality.",
			[]string{"docs/connectors/catalog.md", "docs/connectors/vehicle-avl-starter-kits.md", "docs/integration-adapter-kit.md"},
		),
		connectorCatalogRowView(
			"vendor_shaped_synthetic_examples",
			"Vehicle / GPS / AVL",
			"Vendor-shaped synthetic examples",
			"covered",
			"`testdata/avl-vendor` fixture families for mapping, stale, future, duplicate, and low-quality cases.",
			"/admin/operations/connectors/workbench",
			"go test ./cmd/avl-vendor-adapter",
			"Named vendor support, hardware certification, real source behavior, production AVL reliability, or agency approval.",
			[]string{"docs/connectors/catalog.md", "docs/connectors/vehicle-avl-starter-kits.md", "docs/evidence/redaction-policy.md"},
		),
		connectorCatalogRowView(
			"authenticated_telemetry_post",
			"Vehicle / GPS / AVL",
			"Authenticated telemetry POST",
			"available",
			"`POST /v1/telemetry` with a deployment-owned device bearer token.",
			"/admin/operations/devices and /admin/operations/telemetry",
			"make telemetry-simulator",
			"Real device proof, vendor compatibility, production AVL reliability, compliance, consumer acceptance, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/tutorials/device-avl-integration.md", "docs/tutorials/device-token-lifecycle.md"},
		),
		connectorCatalogRowView(
			"deterministic_builtin_predictor",
			"Prediction",
			"Deterministic built-in predictor",
			"available",
			"Existing deterministic prediction adapter boundary.",
			"/admin/operations/prediction-lab and /admin/operations/realtime",
			"go test ./internal/prediction",
			"Production-grade ETA quality, real-world ETA accuracy, consumer acceptance, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/requirements-trip-updates.md", "docs/tutorials/prediction-eta-lab.md"},
		),
		connectorCatalogRowView(
			"external_http_predictor_adapter",
			"Prediction",
			"External HTTP predictor adapter",
			"review_only",
			"`examples/connectors/predictor-sidecar-stub` with sanitized synthetic request/response fixtures.",
			"/admin/operations/connectors/workbench and /admin/operations/prediction-lab",
			"go run ./cmd/adapter-conformance prediction --suite testdata/adapter-conformance",
			"Named predictor compatibility, live service behavior, production readiness, consumer acceptance, or ETA quality.",
			[]string{"docs/connectors/catalog.md", "docs/requirements-trip-updates.md", "docs/tutorials/external-adapter-conformance.md"},
		),
		connectorCatalogRowView(
			"shadow_mode_predictor",
			"Prediction",
			"Shadow-mode predictor",
			"review_only",
			"Compare external predictor output while public Trip Updates stay on the reviewed path.",
			"/admin/operations/prediction-lab",
			"go run ./examples/connectors/predictor-sidecar-stub",
			"External predictor acceptance, production-grade ETA quality, real-world accuracy, or public-feed improvement.",
			[]string{"docs/connectors/catalog.md", "docs/tutorials/prediction-eta-lab.md", "docs/requirements-trip-updates.md"},
		),
		connectorCatalogRowView(
			"fail_closed_predictor_behavior",
			"Prediction",
			"Fail-closed predictor behavior",
			"covered",
			"Timeout, malformed, stale, wrong-agency, low-confidence, missing Vehicle Positions, and public-mutation cases.",
			"/admin/operations/connectors/tests",
			"make adapter-conformance",
			"ETA quality, real-world accuracy, named predictor compatibility, production readiness, or consumer acceptance.",
			[]string{"docs/connectors/catalog.md", "docs/tutorials/external-adapter-conformance.md", "docs/requirements-trip-updates.md"},
		),
		connectorCatalogRowView(
			"thetransitclock_candidate_notes",
			"Prediction",
			"TheTransitClock candidate notes",
			"candidate_only",
			"Treat TheTransitClock as a future candidate behind the external predictor adapter boundary only.",
			"/admin/operations/connectors/workbench",
			"make adapter-conformance",
			"TheTransitClock integration, compatibility, certification, production ETA quality, or real-world accuracy.",
			[]string{"docs/connectors/catalog.md", "docs/dependencies.md", "docs/requirements-trip-updates.md"},
		),
		connectorCatalogRowView(
			"mobilitydata_static_gtfs_validator",
			"Validator",
			"MobilityData static GTFS validator",
			"available",
			"Server-owned static GTFS validator tooling when installed by an administrator.",
			"/admin/operations/validation-health and /admin/operations/validation-center",
			"make validators-check",
			"Validator-clean proof, CAL-ITP/Caltrans compliance, consumer acceptance, public launch, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/dependencies.md", "docs/tutorials/gtfs-validation-triage.md"},
		),
		connectorCatalogRowView(
			"mobilitydata_gtfs_realtime_validator",
			"Validator",
			"MobilityData GTFS Realtime validator",
			"available",
			"Server-owned GTFS Realtime validator wrapper when installed by an administrator.",
			"/admin/operations/validation-health and /admin/operations/realtime",
			"make gtfsrt-conformance",
			"Validator-clean proof, consumer acceptance, production readiness, public launch, or compliance.",
			[]string{"docs/connectors/catalog.md", "docs/dependencies.md", "docs/tutorials/external-adapter-conformance.md"},
		),
		connectorCatalogRowView(
			"allowlisted_validator_ids",
			"Validator",
			"Allowlisted validator IDs",
			"covered",
			"`examples/connectors/validator-allowlist` with safe validator ID/feed-type mappings only.",
			"/admin/operations/connectors/tests",
			"go run ./cmd/adapter-conformance validator --suite testdata/adapter-conformance",
			"Raw validator command safety beyond the allowlist, validator-clean output, compliance, or consumer acceptance.",
			[]string{"docs/connectors/catalog.md", "docs/connectors/plugin-contract.md", "docs/tutorials/gtfs-validation-triage.md"},
		),
		connectorCatalogRowView(
			"private_validation_health",
			"Validator",
			"Private validation health",
			"available",
			"Private health rows for validator installation, last run, blockers, and next action.",
			"/admin/operations/validation-health",
			"make validate",
			"External approval, consumer acceptance, compliance, production readiness, or public launch.",
			[]string{"docs/connectors/catalog.md", "docs/tutorials/gtfs-validation-triage.md", "docs/release-candidate-readiness.md"},
		),
		connectorCatalogRowView(
			"local_health_summaries",
			"Monitoring / Export",
			"Local health summaries",
			"available",
			"Private feed, validator, telemetry, maintenance, and reliability summaries.",
			"/admin/operations/feed-health and /admin/operations/maintenance",
			"make operations-reliability",
			"SLA coverage, uptime guarantee, hosted service availability, production readiness, or retained evidence.",
			[]string{"docs/connectors/catalog.md", "docs/tutorials/small-agency-maintenance-guide.md", "docs/runbooks/monitoring-and-alerting.md"},
		),
		connectorCatalogRowView(
			"operations_notify_draft",
			"Monitoring / Export",
			"Operations notify draft",
			"review_only",
			"Deployment-owned notification draft workflow with no send by default.",
			"/admin/operations/maintenance",
			"make operations-notify",
			"Notification delivery, incident response maturity, SLA, uptime, hosted service availability, or evidence creation.",
			[]string{"docs/connectors/catalog.md", "docs/tutorials/self-hosted-operations-notifications.md", "docs/evidence/redaction-policy.md"},
		),
		connectorCatalogRowView(
			"monitoring_export_helper",
			"Monitoring / Export",
			"Monitoring/export helper",
			"covered",
			"`examples/connectors/monitoring-export` with redacted no-send batches.",
			"/admin/operations/connectors/workbench",
			"go run ./cmd/adapter-conformance monitoring --suite testdata/adapter-conformance",
			"Notification delivery, SLA coverage, uptime guarantee, hosted service availability, production readiness, or retained evidence.",
			[]string{"docs/connectors/catalog.md", "docs/tutorials/self-hosted-operations-notifications.md", "docs/connectors/redaction-first-recipes.md"},
		),
		connectorCatalogRowView(
			"deployment_owned_monitoring_boundary",
			"Monitoring / Export",
			"Deployment-owned monitoring boundary",
			"needs_review",
			"Any real monitoring destination, credential, retention, and alert routing stays outside this browser page.",
			"/admin/operations/maintenance",
			"make external-connection-check",
			"Hosted service availability, paid support, SLA coverage, uptime guarantee, production readiness, or agency approval.",
			[]string{"docs/connectors/catalog.md", "docs/runbooks/monitoring-and-alerting.md", "docs/support-boundaries.md"},
		),
		connectorCatalogRowView(
			"public_feeds_json",
			"Consumer / Discovery",
			"`/public/feeds.json`",
			"available",
			"Public feed metadata endpoint for a running local or self-hosted instance.",
			"/admin/operations/feeds and /admin/operations/readiness",
			"make smoke",
			"Consumer submission, review, acceptance, ingestion, listing, display, compliance, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/release-candidate-readiness.md", "docs/consumer-submission-evidence.md"},
		),
		connectorCatalogRowView(
			"static_gtfs_url",
			"Consumer / Discovery",
			"Static GTFS URL",
			"available",
			"`/public/gtfs/schedule.zip` when an active feed version exists.",
			"/admin/operations/feeds and /admin/operations/feed-health",
			"make validate-public-feeds",
			"Validator-clean proof, consumer acceptance, public launch, compliance, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/tutorials/gtfs-validation-triage.md", "docs/release-candidate-readiness.md"},
		),
		connectorCatalogRowView(
			"vehicle_positions_url",
			"Consumer / Discovery",
			"Vehicle Positions URL",
			"available",
			"`/public/gtfsrt/vehicle_positions.pb` when realtime publishing is configured.",
			"/admin/operations/realtime and /admin/operations/feed-health",
			"make gtfsrt-conformance",
			"Consumer ingestion, public display, production AVL reliability, compliance, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/tutorials/telemetry-simulator-and-device-trial.md", "docs/requirements-calitp-compliance.md"},
		),
		connectorCatalogRowView(
			"trip_updates_url",
			"Consumer / Discovery",
			"Trip Updates URL",
			"available",
			"`/public/gtfsrt/trip_updates.pb` behind the pluggable prediction boundary.",
			"/admin/operations/realtime and /admin/operations/prediction-lab",
			"make gtfsrt-conformance",
			"Production-grade ETA quality, real-world accuracy, consumer acceptance, compliance, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/requirements-trip-updates.md", "docs/tutorials/prediction-eta-lab.md"},
		),
		connectorCatalogRowView(
			"alerts_url",
			"Consumer / Discovery",
			"Alerts URL",
			"available",
			"`/public/gtfsrt/alerts.pb` plus the separate private Alerts Console.",
			"/admin/operations/realtime and /admin/alerts/console",
			"make gtfsrt-conformance",
			"Consumer acceptance, complete disruption operations, public display, compliance, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/requirements-calitp-compliance.md", "docs/release-candidate-readiness.md"},
		),
		connectorCatalogRowView(
			"consumer_packet_preparedness",
			"Consumer / Discovery",
			"Consumer packet preparedness",
			"prepared_only",
			"Prepared URL and metadata review without external submission automation or status movement.",
			"/admin/operations/consumers and /admin/operations/readiness",
			"scripts/check-consumer-tracker.sh",
			"Consumer submission, review, acceptance, ingestion, listing, display, compliance, or public launch.",
			[]string{"docs/connectors/catalog.md", "docs/consumer-submission-evidence.md", "docs/evidence/consumer-submissions/README.md"},
		),
		connectorCatalogRowView(
			"manifest_based_sidecars",
			"Future Extension Model",
			"Manifest-based sidecars",
			"policy",
			"Connector manifests describe optional sidecars and command adapters; backend contracts stay explicit.",
			"/admin/operations/connectors/workbench",
			"make external-connection-check",
			"Runtime installation, automatic backend extension, vendor compatibility, compliance, or production readiness.",
			[]string{"docs/connectors/catalog.md", "docs/connectors/plugin-contract.md", "docs/extension-governance.md"},
		),
		connectorCatalogRowView(
			"no_dynamic_backend_plugin_loading",
			"Future Extension Model",
			"No arbitrary dynamic backend plugin loading",
			"policy",
			safePluginDefinition,
			"/admin/operations/connectors",
			"make external-connection-check",
			"Third-party runtime distribution support, remote plugin installation, or unreviewed extension execution.",
			[]string{"docs/connectors/catalog.md", "docs/connectors/plugin-contract.md", "docs/extension-governance.md"},
		),
		connectorCatalogRowView(
			"conformance_tests_required",
			"Future Extension Model",
			"Conformance tests required",
			"policy",
			"New connector shapes need manifest validation, synthetic fixtures, fail-closed behavior, and docs.",
			"/admin/operations/connectors/tests",
			"make adapter-conformance",
			"Real integration proof, production readiness, vendor compatibility, compliance, or consumer acceptance.",
			[]string{"docs/connectors/catalog.md", "docs/connectors/contributing-connectors.md", "docs/tutorials/external-adapter-conformance.md"},
		),
	}
}

func connectorCatalogRowView(id string, group string, label string, status string, startWith string, browserReview string, firstCheck string, doesNotProve string, docsLinks []string) connectorCatalogRow {
	return connectorCatalogRow{
		ID:             firstNonEmpty(id, "connector_catalog_row"),
		Group:          firstNonEmpty(group, "Connector"),
		Label:          firstNonEmpty(label, "Connector row"),
		Status:         firstNonEmpty(status, checklistStatusNeedsReview),
		StartWith:      firstNonEmpty(startWith, "Review the Connector Workbench before adapting this connector shape."),
		BrowserReview:  firstNonEmpty(browserReview, "/admin/operations/connectors/workbench"),
		FirstSafeCheck: firstNonEmpty(firstCheck, "make external-connection-check"),
		DoesNotProve:   firstNonEmpty(doesNotProve, "Compatibility, compliance, production readiness, or consumer acceptance."),
		DocsLinks:      safeDocsLinks(docsLinks),
	}
}

func connectorRegistryForSection(section string) connectorpkg.Registry {
	if section != "connectors" {
		return connectorpkg.Registry{}
	}
	return connectorpkg.LoadExampleRegistry()
}

func connectorCategoryView(id string, label string, status string, summary string, shape string, inputs []string, outputs []string, failure string, commands []string, adminLinks []string, docsLinks []string, boundary string) connectorCategory {
	return connectorCategory{
		ID:                 id,
		Label:              label,
		Status:             firstNonEmpty(status, checklistStatusNeedsReview),
		Summary:            firstNonEmpty(summary, "Review connector documentation before use."),
		ConnectorShape:     firstNonEmpty(shape, "sidecar, command adapter, manifest, or connector process"),
		Inputs:             cleanLaunchpadList(inputs),
		Outputs:            cleanLaunchpadList(outputs),
		FailureBehavior:    firstNonEmpty(failure, "Fail closed and keep diagnostics private."),
		CommandSuggestions: safeCommandSuggestions(commands),
		AdminLinks:         safeAdminLinks(adminLinks),
		DocsLinks:          safeDocsLinks(docsLinks),
		ClaimBoundary:      firstNonEmpty(boundary, privateBoundary()),
	}
}
