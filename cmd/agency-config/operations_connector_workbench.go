package main

import (
	"net/http"
	"time"

	"open-transit-rt/internal/auth"
	connectorpkg "open-transit-rt/internal/connectors"
)

type connectorWorkbenchView struct {
	GeneratedAt    time.Time                        `json:"generated_at"`
	AgencyID       string                           `json:"agency_id"`
	Boundary       string                           `json:"boundary"`
	Recipes        []connectorWorkbenchRecipe       `json:"recipes"`
	ManifestReview connectorWorkbenchManifestReview `json:"manifest_review"`
	ClaimFlags     connectorWorkbenchClaimFlags     `json:"claim_flags"`
}

type connectorWorkbenchRecipe struct {
	ID             string   `json:"id"`
	Label          string   `json:"label"`
	OperatorStory  string   `json:"operator_story"`
	Status         string   `json:"status"`
	WhatThisIs     string   `json:"what_this_is"`
	WhatYouNeed    []string `json:"what_you_need"`
	RunsWhere      string   `json:"runs_where"`
	FirstSafeCheck string   `json:"first_safe_check"`
	GoodResult     string   `json:"good_result"`
	IfItFails      string   `json:"if_it_fails"`
	DoesNotProve   string   `json:"does_not_prove"`
	AdminLinks     []string `json:"admin_links"`
	DocsLinks      []string `json:"docs_links"`
	ManifestIDs    []string `json:"manifest_ids"`
}

type connectorWorkbenchManifestReview struct {
	Title            string                            `json:"title"`
	Summary          string                            `json:"summary"`
	PluginDefinition string                            `json:"plugin_definition"`
	Rows             []connectorWorkbenchManifestRow   `json:"rows"`
	Diagnostics      []connectorpkg.RegistryDiagnostic `json:"diagnostics"`
}

type connectorWorkbenchManifestRow struct {
	SourcePath           string   `json:"source_path"`
	ConnectorID          string   `json:"connector_id"`
	DisplayName          string   `json:"display_name"`
	ConnectorType        string   `json:"type"`
	Mode                 string   `json:"mode"`
	DisabledByDefault    bool     `json:"disabled_by_default"`
	FailClosed           bool     `json:"fail_closed"`
	SecretStorage        string   `json:"secret_storage"`
	InputContracts       []string `json:"input_contracts"`
	OutputContracts      []string `json:"output_contracts"`
	ConformanceCaseCount int      `json:"conformance_case_count"`
	DocsLink             string   `json:"docs_link"`
	Boundary             string   `json:"boundary"`
	FirstCheck           string   `json:"first_check"`
	DoesNotProve         string   `json:"does_not_prove"`
}

type connectorWorkbenchClaimFlags struct {
	BackendCommandExecutionEnabled     bool `json:"backend_command_execution_enabled"`
	BrowserNetworkSendEnabled          bool `json:"browser_network_send_enabled"`
	ManifestCommandExecutionEnabled    bool `json:"manifest_command_execution_enabled"`
	DynamicBackendPluginLoadingEnabled bool `json:"dynamic_backend_plugin_loading_enabled"`
	ExternalNetworkContacted           bool `json:"external_network_contacted"`
	ExternalEvidenceCreated            bool `json:"external_evidence_created"`
	ConsumerStatusesChanged            bool `json:"consumer_statuses_changed"`
	ComplianceClaimed                  bool `json:"compliance_claimed"`
	VendorCompatibilityClaimed         bool `json:"vendor_compatibility_claimed"`
	HardwareCertificationClaimed       bool `json:"hardware_certification_claimed"`
	ProductionReadinessClaimed         bool `json:"production_readiness_claimed"`
	HostedSaaSClaimed                  bool `json:"hosted_saas_claimed"`
	SLAClaimed                         bool `json:"sla_claimed"`
	ProductionGradeETAClaimed          bool `json:"production_grade_eta_claimed"`
}

func (h *handler) renderConnectorWorkbench(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "connector-workbench")
	w.Header().Set("Cache-Control", "no-store")
	renderOperationsTemplate(w, "connector-workbench", page)
}

func (h *handler) renderConnectorWorkbenchJSON(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.RequireRole(w, r, auth.RoleReadOnly, auth.RoleOperator, auth.RoleEditor, auth.RoleAdmin)
	if !ok || !auth.RequireAgencyQueryMatch(w, r, principal) {
		return
	}
	page := h.buildOperationsPage(r, principal, "connector-workbench")
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, page.ConnectorWorkbench)
}

func buildConnectorWorkbench(page operationsPage) connectorWorkbenchView {
	registry := connectorpkg.LoadExampleRegistry()
	return connectorWorkbenchView{
		GeneratedAt: page.GeneratedAt,
		AgencyID:    page.AgencyID,
		Boundary:    "Private connector planning and review only. This page does not launch sidecars, run connector commands, send telemetry, contact vendors, create evidence, change consumer statuses, or prove compatibility.",
		Recipes: []connectorWorkbenchRecipe{
			connectorWorkbenchRecipeView(
				"csv_telemetry_sandbox",
				"I have a CSV of vehicle locations",
				"Replay a bounded CSV fixture through the synthetic telemetry adapter path.",
				"covered",
				"Use the CSV replay example to understand column mapping, timestamps, device IDs, and fail-closed behavior before any real feed is configured.",
				[]string{"Committed CSV fixture", "agency-owned device mapping", "operator shell for local dry run"},
				"Operator shell using committed examples only.",
				"make test-connector-examples",
				"Synthetic rows normalize or are rejected with bounded diagnostics.",
				"Review timestamp, device, agency, and required-field mapping before changing fixtures or adapter code.",
				"Real fleet reliability, hardware certification, real vendor proof, production readiness, compliance, consumer acceptance, or ETA quality.",
				[]string{"/admin/operations/telemetry-simulator", "/admin/operations/devices"},
				[]string{"examples/connectors/telemetry-csv-replay/README.md", "docs/tutorials/device-avl-integration.md"},
				[]string{"example.telemetry-csv-replay"},
			),
			connectorWorkbenchRecipeView(
				"api_polling_recipe",
				"I have a GPS API",
				"Review the HTTP polling sidecar shape without storing credentials or contacting a live endpoint from the browser.",
				"covered",
				"Use the synthetic HTTP poller example to review request ownership, response normalization, retries, and no-send defaults.",
				[]string{"Synthetic observation fixture", "deployment-owned credential reference outside the repo", "operator shell for dry run"},
				"Sidecar or command adapter outside the web request.",
				"make external-connection-check",
				"Committed manifests and examples remain disabled by default and fail closed.",
				"Keep live URLs, credentials, and private payloads out of manifests and browser requests.",
				"Named API support, vendor compatibility, real network behavior, production readiness, compliance, or consumer acceptance.",
				[]string{"/admin/operations/connectors", "/admin/operations/connectors/tests"},
				[]string{"examples/connectors/telemetry-http-poller/README.md", "docs/integration-adapter-kit.md"},
				[]string{"example.telemetry-http-poller"},
			),
			connectorWorkbenchRecipeView(
				"webhook_transform_boundary",
				"I have an AVL source that can POST",
				"Review where an adapter can transform inbound observations before authenticated telemetry ingest.",
				"needs_review",
				"Use this as a boundary checklist for deployment-owned webhook receivers and transform adapters.",
				[]string{"Deployment-owned receiver outside the Operations Console", "device token issued through Device Credentials", "redacted local diagnostics"},
				"Outside this browser page; the Workbench does not receive or forward posted payloads.",
				"go run ./cmd/adapter-conformance telemetry --suite testdata/adapter-conformance",
				"Synthetic malformed, stale, future, duplicate, low-quality, and wrong-agency cases fail closed.",
				"Review auth, timestamp, agency, and payload redaction before connecting any external source.",
				"Real vendor support, hardware certification, production AVL quality, SLA, compliance, consumer acceptance, or public launch.",
				[]string{"/admin/operations/devices", "/admin/operations/telemetry"},
				[]string{"docs/tutorials/device-avl-integration.md", "docs/evidence/redaction-policy.md"},
				[]string{"example.telemetry-http-poller", "example.telemetry-csv-replay"},
			),
			connectorWorkbenchRecipeView(
				"synthetic_only",
				"I want synthetic telemetry only",
				"Use the simulator and committed conformance fixtures to exercise local behavior without external sources.",
				"covered",
				"Keep evaluation local by using synthetic telemetry scenarios and no-send connector examples.",
				[]string{"Local app", "synthetic scenario fixtures", "device credentials for intentional local telemetry exercises"},
				"Private console review plus operator shell for local checks.",
				"make telemetry-simulator",
				"Local diagnostics show matched, unmatched, stale, and fail-closed behavior without external proof.",
				"Use Realtime Center to review resulting private diagnostics and keep all external evidence gates closed.",
				"Agency approval, real device proof, production readiness, compliance, consumer acceptance, or ETA quality.",
				[]string{"/admin/operations/telemetry-simulator", "/admin/operations/realtime"},
				[]string{"docs/tutorials/telemetry-simulator-and-device-trial.md", "docs/tutorials/external-adapter-conformance.md"},
				[]string{"example.telemetry-csv-replay", "example.telemetry-http-poller"},
			),
			connectorWorkbenchRecipeView(
				"predictor_sidecar",
				"I want an external predictor",
				"Review prediction adapter shape, shadow mode, and fail-closed behavior while Vehicle Positions stay independent.",
				"covered",
				"Use the predictor sidecar stub and prediction conformance cases to review bounded request/response handling.",
				[]string{"Active feed version", "current telemetry", "current assignments", "Vehicle Positions data or URL", "operator shell for local checks"},
				"Optional prediction sidecar behind internal/prediction.Adapter.",
				"go run ./cmd/adapter-conformance prediction --suite testdata/adapter-conformance",
				"Timeout, malformed, stale, wrong-agency, and low-confidence cases are withheld or fail closed.",
				"Keep deterministic fallback available and never block Vehicle Positions on predictor availability.",
				"Production-grade ETA quality, real-world accuracy, named predictor support, consumer acceptance, or production readiness.",
				[]string{"/admin/operations/realtime", "/admin/operations/feed-health"},
				[]string{"examples/connectors/predictor-sidecar-stub/README.md", "docs/requirements-trip-updates.md"},
				[]string{"example.predictor-sidecar-stub"},
			),
			connectorWorkbenchRecipeView(
				"monitoring_export",
				"I want monitoring summaries",
				"Review redacted no-send monitoring/export patterns for private operational summaries.",
				"covered",
				"Use the monitoring export example to understand redaction, no-send defaults, and deployment-owned delivery boundaries.",
				[]string{"Local reliability or validator summary", "redaction policy", "deployment-owned delivery config outside this page"},
				"Operator shell or deployment process; browser review only.",
				"go run ./cmd/adapter-conformance monitoring --suite testdata/adapter-conformance",
				"Monitoring fixtures preserve redaction and no-send defaults.",
				"Configure any delivery outside the Workbench with explicit deployment ownership and redaction review.",
				"SLA coverage, uptime guarantee, hosted service availability, production readiness, compliance, or retained evidence.",
				[]string{"/admin/operations/maintenance", "/admin/operations/reliability"},
				[]string{"examples/connectors/monitoring-export/README.md", "docs/tutorials/self-hosted-operations-notifications.md"},
				[]string{"example.monitoring-export"},
			),
			connectorWorkbenchRecipeView(
				"public_feed_url_verification",
				"I want off-host validation",
				"Review validator and public feed URL checks as operator-run diagnostics, not proof.",
				"covered",
				"Use server-owned validator IDs and operator-run off-host checks when local server resources are too small.",
				[]string{"Configured public feed base URL", "validator tooling or off-host validator machine", "operator shell"},
				"Operator shell outside the browser.",
				"make validate-public-feeds",
				"Feed URLs can be fetched and validator blockers are recorded as private diagnostics.",
				"Treat missing Java, Docker, pinned tools, or stale reports as blocker rows, not success.",
				"Validator-clean proof, CAL-ITP/Caltrans compliance, consumer acceptance, production readiness, or public launch.",
				[]string{"/admin/operations/validation-center", "/admin/operations/feed-health"},
				[]string{"examples/connectors/validator-allowlist/README.md", "docs/tutorials/gtfs-validation-triage.md"},
				[]string{"example.validator-allowlist"},
			),
		},
		ManifestReview: connectorWorkbenchManifestReview{
			Title:            "Example Manifest Registry Review",
			Summary:          "Committed synthetic connector manifests only. This registry review does not accept uploads, load backend plugins, execute manifest commands, contact external systems, create retained evidence, or change consumer status.",
			PluginDefinition: safePluginDefinition,
			Rows:             connectorWorkbenchManifestRows(registry),
			Diagnostics:      registry.Diagnostics,
		},
		ClaimFlags: connectorWorkbenchClaimFlags{},
	}
}

func connectorWorkbenchRecipeView(id string, label string, story string, status string, what string, need []string, runsWhere string, firstCheck string, good string, fail string, doesNotProve string, adminLinks []string, docsLinks []string, manifestIDs []string) connectorWorkbenchRecipe {
	return connectorWorkbenchRecipe{
		ID:             firstNonEmpty(id, "connector_recipe"),
		Label:          firstNonEmpty(label, "Connector recipe"),
		OperatorStory:  firstNonEmpty(story, "Review a local/synthetic connector recipe."),
		Status:         firstNonEmpty(status, checklistStatusNeedsReview),
		WhatThisIs:     firstNonEmpty(what, "Private connector planning guidance."),
		WhatYouNeed:    cleanLaunchpadList(need),
		RunsWhere:      firstNonEmpty(runsWhere, "Operator shell outside the browser."),
		FirstSafeCheck: firstNonEmpty(firstCheck, "make external-connection-check"),
		GoodResult:     firstNonEmpty(good, "Synthetic/local checks pass or produce bounded diagnostics."),
		IfItFails:      firstNonEmpty(fail, "Review the synthetic fixture, mapping, or adapter boundary."),
		DoesNotProve:   firstNonEmpty(doesNotProve, "Compatibility, compliance, production readiness, or consumer acceptance."),
		AdminLinks:     safeAdminLinks(adminLinks),
		DocsLinks:      safeDocsLinks(docsLinks),
		ManifestIDs:    cleanLaunchpadList(manifestIDs),
	}
}

func connectorWorkbenchManifestRows(registry connectorpkg.Registry) []connectorWorkbenchManifestRow {
	rows := make([]connectorWorkbenchManifestRow, 0, len(registry.Entries))
	for _, entry := range registry.Entries {
		rows = append(rows, connectorWorkbenchManifestRow{
			SourcePath:           entry.SourcePath,
			ConnectorID:          entry.ConnectorID,
			DisplayName:          entry.DisplayName,
			ConnectorType:        entry.ConnectorType,
			Mode:                 entry.ModeName,
			DisabledByDefault:    entry.DisabledByDefault,
			FailClosed:           entry.FailureBehavior.FailClosed,
			SecretStorage:        entry.RedactionPolicy.SecretStorage,
			InputContracts:       connectorRegistryContractNames(entry.InputContracts),
			OutputContracts:      connectorRegistryContractNames(entry.OutputContracts),
			ConformanceCaseCount: len(entry.ConformanceCases),
			DocsLink:             entry.DocsLink,
			Boundary:             "Example manifest review only; manifests are static sidecar/adapter descriptions, not dynamic backend plugins.",
			FirstCheck:           "make external-connection-check",
			DoesNotProve:         "Real integration proof, vendor compatibility, hardware certification, compliance, consumer acceptance, production readiness, hosted service, SLA, or ETA quality.",
		})
	}
	return rows
}

func connectorRegistryContractNames(contracts []connectorpkg.RegistryContract) []string {
	names := make([]string, 0, len(contracts))
	for _, contract := range contracts {
		names = append(names, firstNonEmpty(contract.Name, "contract"))
	}
	return cleanLaunchpadList(names)
}
