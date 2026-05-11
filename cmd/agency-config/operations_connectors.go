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
	Categories       []connectorCategory    `json:"categories"`
	Registry         connectorpkg.Registry  `json:"registry"`
	ClaimFlags       connectorHubClaimFlags `json:"claim_flags"`
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
			"Telemetry / GPS / AVL source",
			"available",
			"Transform GPS, AVL, CSV, or device observations before sending them to authenticated telemetry ingest.",
			"sidecar or command adapter that calls authenticated POST /v1/telemetry",
			[]string{"external observations", "agency-owned mapping", "device credential reference"},
			[]string{"Open Transit RT telemetry.Event records", "private redacted adapter diagnostics"},
			"Fail closed before ingest when mapping, credentials, timestamps, or payload quality are unsafe.",
			[]string{"make telemetry-simulator", "make adapter-conformance"},
			[]string{"/admin/operations/devices", "/admin/operations/telemetry"},
			[]string{"docs/integration-adapter-kit.md", "docs/tutorials/device-avl-integration.md", "docs/tutorials/telemetry-simulator-and-device-trial.md"},
			"Telemetry connectors do not prove real vendor compatibility, hardware certification, production AVL reliability, or agency adoption.",
		),
		connectorCategoryView(
			"prediction_engine",
			"Prediction engine",
			"available",
			"Keep Trip Updates replaceable behind the prediction adapter while Vehicle Positions remain independent.",
			"optional internal adapter or external HTTP sidecar behind internal/prediction.Adapter",
			[]string{"active GTFS feed version", "latest telemetry", "current assignments", "Vehicle Positions feed data or URL"},
			[]string{"Trip Updates feed", "bounded diagnostics"},
			"Return deterministic output or valid empty Trip Updates with adapter diagnostics when an optional sidecar fails.",
			[]string{"make adapter-conformance", "go test ./internal/prediction"},
			[]string{"/admin/operations/feeds", "/admin/operations/reliability"},
			[]string{"docs/requirements-trip-updates.md", "docs/integration-adapter-kit.md", "docs/tutorials/external-adapter-conformance.md"},
			"Prediction connectors do not prove production-grade ETA quality, consumer acceptance, or named predictor compatibility.",
		),
		connectorCategoryView(
			"validator",
			"Validator connector",
			"available",
			"Run static and realtime validators through server-owned allowlisted validator IDs.",
			"allowlisted command adapter configured by deployment-owned environment",
			[]string{"server-derived schedule ZIP", "server-derived realtime protobuf", "validator ID"},
			[]string{"normalized validation_report rows", "private validator health diagnostics"},
			"Missing or failing tooling records not-run, missing, warning, failed, or blocked states instead of pretending success.",
			[]string{"make validators-check", "make validator-health"},
			[]string{"/admin/operations/validation-health", "/admin/operations/gtfs-quality"},
			[]string{"docs/dependencies.md", "docs/tutorials/gtfs-validation-triage.md", "docs/release-candidate-readiness.md"},
			"Validator records are supporting signals only; they are not CAL-ITP/Caltrans compliance or consumer acceptance.",
		),
		connectorCategoryView(
			"monitoring_export",
			"Monitoring / export connector",
			"available",
			"Export redacted local diagnostics to deployment-owned monitoring without sending anything by default.",
			"deployment-owned no-send exporter or redacted summary writer",
			[]string{"validator health summary", "deployment doctor summary", "reliability summary"},
			[]string{"private redacted monitoring/export summaries"},
			"Default examples write local private summaries only and require deployment-owned configuration before any external send.",
			[]string{"make operations-reliability", "make operations-notify"},
			[]string{"/admin/operations/reliability", "/admin/operations/validation-health"},
			[]string{"docs/tutorials/self-hosted-operations-notifications.md", "docs/deployment/reference-deployment-doctor.md"},
			"Monitoring connectors do not prove SLA coverage, uptime guarantees, hosted SaaS availability, or production readiness.",
		),
		connectorCategoryView(
			"consumer_discovery",
			"Consumer / discovery workflow",
			"available",
			"Prepare feed URL and metadata packets while leaving external submissions and status changes authorization-gated.",
			"docs workflow and packet generator, not portal automation",
			[]string{"stable feed URL metadata", "license/contact metadata", "validation status"},
			[]string{"prepared consumer packet records", "private readiness summaries"},
			"Do not contact targets or change statuses without retained target-originated or operator-retained authorization artifacts.",
			[]string{"make external-connection-check", "make audit-final-claim-review"},
			[]string{"/admin/operations/consumers", "/admin/operations/readiness"},
			[]string{"docs/consumer-submission-evidence.md", "docs/evidence/consumer-submissions/README.md", "docs/compliance-evidence-checklist.md"},
			"Prepared packets are not submission, review, acceptance, ingestion, listing, display, compliance, or public launch proof.",
		),
	}
	return connectorHubView{
		GeneratedAt:      page.GeneratedAt,
		AgencyID:         page.AgencyID,
		Boundary:         "Private authenticated Connector Hub only; viewing it creates no evidence, contacts no external party, changes no consumer status, and records no approval, compatibility, compliance, hosted-service, SLA, production-readiness, or ETA-quality outcome. Treat it as the starting point for manifest, redaction, fail-closed, and adapter-conformance review.",
		PluginDefinition: safePluginDefinition,
		Categories:       categories,
		Registry:         connectorRegistryForSection(page.Section),
		ClaimFlags:       connectorHubClaimFlags{},
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
