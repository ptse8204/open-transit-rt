package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"open-transit-rt/internal/auth"
	connectorpkg "open-transit-rt/internal/connectors"
)

const (
	connectorWorkbenchCSVFixture  = "examples/connectors/telemetry-csv-replay/fixtures/replay.csv"
	connectorWorkbenchHTTPFixture = "examples/connectors/telemetry-http-poller/fixtures/observations.json"
	connectorWorkbenchSuitePath   = "testdata/adapter-conformance/suite.json"
)

type connectorWorkbenchView struct {
	GeneratedAt        time.Time                             `json:"generated_at"`
	AgencyID           string                                `json:"agency_id"`
	Boundary           string                                `json:"boundary"`
	DecisionTree       []connectorWorkbenchDecisionRow       `json:"decision_tree"`
	Recipes            []connectorWorkbenchRecipe            `json:"recipes"`
	RedactionTemplates []connectorWorkbenchRedactionTemplate `json:"redaction_templates"`
	DryRunCommands     []connectorWorkbenchDryRun            `json:"dry_run_commands"`
	TelemetryPreview   connectorWorkbenchTelemetryPreview    `json:"telemetry_preview"`
	WebhookBoundary    connectorWorkbenchWebhookBoundary     `json:"webhook_boundary"`
	PredictionGuide    connectorWorkbenchGuide               `json:"prediction_guide"`
	MonitoringGuide    connectorWorkbenchGuide               `json:"monitoring_guide"`
	Conformance        connectorWorkbenchConformanceView     `json:"conformance"`
	ManifestReview     connectorWorkbenchManifestReview      `json:"manifest_review"`
	ClaimFlags         connectorWorkbenchClaimFlags          `json:"claim_flags"`
}

type connectorWorkbenchDecisionRow struct {
	ID             string   `json:"id"`
	SourceSignal   string   `json:"source_signal"`
	UseWhen        string   `json:"use_when"`
	Boundary       string   `json:"boundary"`
	FirstSafeCheck string   `json:"first_safe_check"`
	StopIf         string   `json:"stop_if"`
	NextAdminLink  string   `json:"next_admin_link"`
	DocsLinks      []string `json:"docs_links"`
	DoesNotProve   string   `json:"does_not_prove"`
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

type connectorWorkbenchRedactionTemplate struct {
	ID                 string   `json:"id"`
	Label              string   `json:"label"`
	AppliesTo          string   `json:"applies_to"`
	DataClassification string   `json:"data_classification"`
	AllowedFields      []string `json:"allowed_fields"`
	RedactFields       []string `json:"redact_fields"`
	BlockedFields      []string `json:"blocked_fields"`
	NoSendDefault      string   `json:"no_send_default"`
	FailClosedRule     string   `json:"fail_closed_rule"`
	FirstSafeCheck     string   `json:"first_safe_check"`
	DoesNotProve       string   `json:"does_not_prove"`
	DocsLinks          []string `json:"docs_links"`
}

type connectorWorkbenchManifestReview struct {
	Title            string                            `json:"title"`
	Summary          string                            `json:"summary"`
	PluginDefinition string                            `json:"plugin_definition"`
	Rows             []connectorWorkbenchManifestRow   `json:"rows"`
	LintChecks       []connectorWorkbenchManifestLint  `json:"lint_checks"`
	Diagnostics      []connectorpkg.RegistryDiagnostic `json:"diagnostics"`
}

type connectorWorkbenchManifestLint struct {
	ID             string `json:"id"`
	Label          string `json:"label"`
	Status         string `json:"status"`
	EnforcedBy     string `json:"enforced_by"`
	Blocks         string `json:"blocks"`
	OperatorAction string `json:"operator_action"`
	DoesNotProve   string `json:"does_not_prove"`
}

type connectorWorkbenchDryRun struct {
	ID                string   `json:"id"`
	Label             string   `json:"label"`
	CommandLine       string   `json:"command_line"`
	RunsWhere         string   `json:"runs_where"`
	Inputs            string   `json:"inputs"`
	ExpectedResult    string   `json:"expected_result"`
	FailureNextAction string   `json:"failure_next_action"`
	DoesNotProve      string   `json:"does_not_prove"`
	DocsLinks         []string `json:"docs_links"`
}

type connectorWorkbenchTelemetryPreview struct {
	Boundary string                            `json:"boundary"`
	Sources  []connectorWorkbenchPreviewSource `json:"sources"`
	Rows     []connectorWorkbenchPreviewRow    `json:"rows"`
	Counts   connectorWorkbenchPreviewCounts   `json:"counts"`
}

type connectorWorkbenchPreviewSource struct {
	ID             string `json:"id"`
	Label          string `json:"label"`
	FixturePath    string `json:"fixture_path"`
	Status         string `json:"status"`
	SyntheticOnly  bool   `json:"synthetic_only"`
	ObservedRows   int    `json:"observed_rows"`
	ExpectedEvents int    `json:"expected_events"`
	ExpectedDrops  int    `json:"expected_drops"`
	CommandLine    string `json:"command_line"`
	DoesNotProve   string `json:"does_not_prove"`
}

type connectorWorkbenchPreviewRow struct {
	SourceID    string `json:"source_id"`
	DeviceID    string `json:"device_id"`
	VehicleID   string `json:"vehicle_id"`
	ObservedAt  string `json:"observed_at"`
	Quality     string `json:"quality"`
	Outcome     string `json:"outcome"`
	Reason      string `json:"reason,omitempty"`
	DryRun      bool   `json:"dry_run"`
	NetworkSend bool   `json:"network_send"`
}

type connectorWorkbenchPreviewCounts struct {
	Sources            int  `json:"sources"`
	Rows               int  `json:"rows"`
	Events             int  `json:"events"`
	Drops              int  `json:"drops"`
	NetworkSendEnabled bool `json:"network_send_enabled"`
}

type connectorWorkbenchWebhookBoundary struct {
	Title     string                         `json:"title"`
	Boundary  string                         `json:"boundary"`
	Rows      []connectorWorkbenchWebhookRow `json:"rows"`
	DocsLinks []string                       `json:"docs_links"`
}

type connectorWorkbenchWebhookRow struct {
	ID             string   `json:"id"`
	Label          string   `json:"label"`
	WhatThisMeans  string   `json:"what_this_means"`
	AllowedInputs  []string `json:"allowed_inputs"`
	BlockedInputs  []string `json:"blocked_inputs"`
	FirstSafeCheck string   `json:"first_safe_check"`
	FailClosedRule string   `json:"fail_closed_rule"`
	RedactionRule  string   `json:"redaction_rule"`
	DoesNotProve   string   `json:"does_not_prove"`
	ReviewLinks    []string `json:"review_links"`
}

type connectorWorkbenchGuide struct {
	Title     string                       `json:"title"`
	Boundary  string                       `json:"boundary"`
	Rows      []connectorWorkbenchGuideRow `json:"rows"`
	DocsLinks []string                     `json:"docs_links"`
}

type connectorWorkbenchGuideRow struct {
	ID              string   `json:"id"`
	Label           string   `json:"label"`
	Status          string   `json:"status"`
	WhatThisIs      string   `json:"what_this_is"`
	Inputs          []string `json:"inputs"`
	Outputs         []string `json:"outputs"`
	FailureBehavior string   `json:"failure_behavior"`
	FirstSafeCheck  string   `json:"first_safe_check"`
	DoesNotProve    string   `json:"does_not_prove"`
	ReviewLinks     []string `json:"review_links"`
	DocsLinks       []string `json:"docs_links"`
}

type connectorWorkbenchConformanceView struct {
	Boundary       string                        `json:"boundary"`
	SuitePath      string                        `json:"suite_path"`
	Status         string                        `json:"status"`
	SyntheticOnly  bool                          `json:"synthetic_only"`
	ManifestCount  int                           `json:"manifest_count"`
	CaseCount      int                           `json:"case_count"`
	Groups         []connectorWorkbenchCaseGroup `json:"groups"`
	RunnerCommands []connectorWorkbenchDryRun    `json:"runner_commands"`
}

type connectorWorkbenchCaseGroup struct {
	ID                string                          `json:"id"`
	Label             string                          `json:"label"`
	Status            string                          `json:"status"`
	CaseCount         int                             `json:"case_count"`
	RequiredScenarios []string                        `json:"required_scenarios"`
	CommandLine       string                          `json:"command_line"`
	DoesNotProve      string                          `json:"does_not_prove"`
	Cases             []connectorWorkbenchCaseSummary `json:"cases"`
}

type connectorWorkbenchCaseSummary struct {
	ID              string   `json:"id"`
	Scenario        string   `json:"scenario"`
	FixturePath     string   `json:"fixture_path"`
	ExpectedOutcome string   `json:"expected_outcome"`
	Assertions      []string `json:"assertions"`
	Status          string   `json:"status"`
	SyntheticOnly   bool     `json:"synthetic_only"`
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
		GeneratedAt:  page.GeneratedAt,
		AgencyID:     page.AgencyID,
		Boundary:     "Private connector planning and review only. This page does not launch sidecars, run connector commands, send telemetry, contact vendors, create evidence, change consumer statuses, or prove compatibility.",
		DecisionTree: connectorWorkbenchDecisionTreeRows(),
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
				[]string{"Deployment-owned receiver outside the Operations Console", "device token issued through Devices & Tokens", "redacted local diagnostics"},
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
		RedactionTemplates: connectorWorkbenchRedactionTemplates(),
		DryRunCommands:     connectorWorkbenchDryRunCommands(),
		TelemetryPreview:   buildConnectorWorkbenchTelemetryPreview(),
		WebhookBoundary:    connectorWorkbenchWebhookBoundaryView(),
		PredictionGuide:    connectorWorkbenchPredictionGuideView(),
		MonitoringGuide:    connectorWorkbenchMonitoringGuideView(),
		Conformance:        buildConnectorWorkbenchConformanceView(),
		ManifestReview: connectorWorkbenchManifestReview{
			Title:            "Example Manifest Registry Review",
			Summary:          "Committed synthetic connector manifests only. This registry review does not accept uploads, load backend plugins, execute manifest commands, contact external systems, create retained evidence, or change consumer status.",
			PluginDefinition: safePluginDefinition,
			Rows:             connectorWorkbenchManifestRows(registry),
			LintChecks:       connectorWorkbenchManifestLintChecks(),
			Diagnostics:      registry.Diagnostics,
		},
		ClaimFlags: connectorWorkbenchClaimFlags{},
	}
}

func connectorWorkbenchDecisionTreeRows() []connectorWorkbenchDecisionRow {
	return []connectorWorkbenchDecisionRow{
		connectorWorkbenchDecisionRowView(
			"csv_vehicle_locations",
			"Vehicle location rows arrive as a file or spreadsheet export.",
			"Use the CSV replay example first, then review column mapping and timestamp/device identity before any deployment-owned import path.",
			"CSV replay adapter around committed or operator-redacted rows; no browser upload and no automatic ingest.",
			"make test-connector-examples",
			"Stop if rows include private identifiers, unredacted locations from a real fleet, credentials, or unclear device ownership.",
			"/admin/operations/telemetry-simulator",
			[]string{"examples/connectors/telemetry-csv-replay/README.md", "docs/tutorials/device-avl-integration.md"},
			"Production data quality, hardware certification, real device proof, compliance, consumer acceptance, or production readiness.",
		),
		connectorWorkbenchDecisionRowView(
			"gps_polling_api",
			"A source exposes a GPS polling API.",
			"Use the synthetic HTTP poller shape to review response normalization, retries, credential ownership, and no-send defaults.",
			"Deployment-owned sidecar or command adapter; credentials stay outside manifests and outside browser pages.",
			"make external-connection-check",
			"Stop if the manifest contains live URLs, bearer values, private endpoint text, or source payload bodies.",
			"/admin/operations/connectors/tests",
			[]string{"examples/connectors/telemetry-http-poller/README.md", "docs/integration-adapter-kit.md"},
			"Named API support, vendor compatibility, external connectivity, production AVL reliability, compliance, or public launch.",
		),
		connectorWorkbenchDecisionRowView(
			"avl_can_post",
			"An AVL system can POST observations.",
			"Treat the receiver as deployment-owned infrastructure and transform observations before authenticated telemetry ingest.",
			"Webhook receiver outside the Operations Console; this page only reviews allowed inputs, blocked inputs, redaction, and fail-closed behavior.",
			"go run ./cmd/adapter-conformance telemetry --suite testdata/adapter-conformance",
			"Stop if authentication, agency mapping, device binding, timestamp handling, or redaction is not reviewed.",
			"/admin/operations/devices",
			[]string{"docs/tutorials/device-avl-integration.md", "docs/evidence/redaction-policy.md"},
			"Vendor support, hardware certification, production AVL reliability, SLA, compliance, consumer acceptance, or agency approval.",
		),
		connectorWorkbenchDecisionRowView(
			"synthetic_only",
			"No external source is available or authorized.",
			"Use simulator and committed conformance fixtures to exercise local behavior without contacting a third party.",
			"Synthetic-only local diagnostics through existing simulator and adapter examples.",
			"make telemetry-simulator",
			"Stop before using real credentials, real payloads, retained evidence, or any public claim.",
			"/admin/operations/telemetry-simulator",
			[]string{"docs/tutorials/telemetry-simulator-and-device-trial.md", "docs/tutorials/external-adapter-conformance.md"},
			"Agency adoption, real device proof, consumer acceptance, compliance, production readiness, or ETA quality.",
		),
		connectorWorkbenchDecisionRowView(
			"prediction_sidecar",
			"Trip Updates prediction needs an external sidecar review.",
			"Keep Vehicle Positions independent and evaluate predictor output in shadow or fail-closed mode using sanitized synthetic inputs.",
			"Optional prediction adapter boundary behind `internal/prediction.Adapter`; no public feed mutation from this page.",
			"go run ./cmd/adapter-conformance prediction --suite testdata/adapter-conformance",
			"Stop if sidecar output is stale, malformed, wrong-agency, low-confidence, lacks Vehicle Positions reference, or attempts public mutation.",
			"/admin/operations/prediction-lab",
			[]string{"examples/connectors/predictor-sidecar-stub/README.md", "docs/requirements-trip-updates.md"},
			"Production-grade ETA quality, real-world accuracy, named predictor compatibility, consumer acceptance, or release readiness.",
		),
		connectorWorkbenchDecisionRowView(
			"monitoring_export",
			"Operators want summaries outside the console.",
			"Start with redacted no-send monitoring/export output and keep delivery deployment-owned and disabled by default.",
			"Monitoring/export summary adapter with destination redaction and no notification send from examples.",
			"go run ./cmd/adapter-conformance monitoring --suite testdata/adapter-conformance",
			"Stop if destination, contact address, token, private endpoint, or unredacted incident detail appears in output.",
			"/admin/operations/maintenance",
			[]string{"examples/connectors/monitoring-export/README.md", "docs/tutorials/self-hosted-operations-notifications.md"},
			"SLA coverage, uptime guarantee, hosted service availability, retained evidence, production readiness, or compliance.",
		),
		connectorWorkbenchDecisionRowView(
			"off_host_validation",
			"The local host is too small for validator tooling.",
			"Use server-owned validator IDs and operator-run off-host checks; record blockers privately without turning them into evidence.",
			"Operator-run validator workflow using allowlisted IDs and redacted artifact references.",
			"go run ./cmd/adapter-conformance validator --suite testdata/adapter-conformance",
			"Stop if a raw validator command, private artifact path, or unsupported validator ID is required.",
			"/admin/operations/validation-center",
			[]string{"examples/connectors/validator-allowlist/README.md", "docs/tutorials/gtfs-validation-triage.md"},
			"Validator-clean proof, CAL-ITP/Caltrans compliance, consumer acceptance, public launch, or production readiness.",
		),
	}
}

func connectorWorkbenchDecisionRowView(id string, signal string, useWhen string, boundary string, check string, stopIf string, adminLink string, docsLinks []string, doesNotProve string) connectorWorkbenchDecisionRow {
	adminLinks := safeAdminLinks([]string{adminLink})
	nextAdminLink := "/admin/operations/connectors/tests"
	if len(adminLinks) > 0 {
		nextAdminLink = adminLinks[0]
	}
	return connectorWorkbenchDecisionRow{
		ID:             firstNonEmpty(id, "connector_decision"),
		SourceSignal:   firstNonEmpty(signal, "Review the available source shape."),
		UseWhen:        firstNonEmpty(useWhen, "Choose the safest local/synthetic boundary."),
		Boundary:       firstNonEmpty(boundary, "Private connector review boundary."),
		FirstSafeCheck: firstNonEmpty(check, "make external-connection-check"),
		StopIf:         firstNonEmpty(stopIf, "Stop if credentials, private payloads, external contact, evidence, or stronger claims are required."),
		NextAdminLink:  nextAdminLink,
		DocsLinks:      safeDocsLinks(docsLinks),
		DoesNotProve:   firstNonEmpty(doesNotProve, "Compatibility, compliance, production readiness, or consumer acceptance."),
	}
}

func connectorWorkbenchRedactionTemplates() []connectorWorkbenchRedactionTemplate {
	return []connectorWorkbenchRedactionTemplate{
		connectorWorkbenchRedactionTemplateView(
			"telemetry_source",
			"Telemetry Source Template",
			"CSV replay, GPS polling API, or AVL POST transform.",
			"Device and vehicle observations; treat as private until operator-redacted.",
			[]string{"agency_id", "device_id", "vehicle_id", "timestamp", "lat", "lon", "quality or accuracy"},
			[]string{"source record id", "operator contact", "private endpoint", "authorization header", "source payload body"},
			[]string{"bearer values", "API keys", "database URLs", "private paths", "unbounded vendor fields"},
			"send_enabled=false and network_send=false until deployment owner enables an adapter outside this page.",
			"Reject before ingest when identity, timestamp, agency, coordinates, quality, duplicate, or ordering checks fail.",
			"go run ./cmd/adapter-conformance telemetry --suite testdata/adapter-conformance",
			"Accepted telemetry, Vehicle Positions output, real device proof, vendor compatibility, hardware certification, compliance, or production readiness.",
			[]string{"docs/tutorials/device-avl-integration.md", "docs/evidence/redaction-policy.md"},
		),
		connectorWorkbenchRedactionTemplateView(
			"prediction_sidecar",
			"Prediction Sidecar Template",
			"External predictor shadow or fail-closed review.",
			"Sanitized prediction input and output; treat as private diagnostics until reviewed.",
			[]string{"active feed version id", "assignment ids", "Vehicle Positions reference", "synthetic timestamps", "confidence"},
			[]string{"source private payload", "private URL", "operator notes", "raw predictor diagnostics"},
			[]string{"public feed mutation flags", "consumer submission fields", "credentials", "unvalidated Trip Updates"},
			"public_mutation=false; Vehicle Positions publishing remains independent.",
			"Withhold output when predictor response is stale, malformed, wrong-agency, low-confidence, or missing Vehicle Positions reference.",
			"go run ./cmd/adapter-conformance prediction --suite testdata/adapter-conformance",
			"Production-grade ETA quality, real-world accuracy, named predictor support, consumer acceptance, release readiness, or production readiness.",
			[]string{"docs/requirements-trip-updates.md", "docs/tutorials/prediction-eta-lab.md"},
		),
		connectorWorkbenchRedactionTemplateView(
			"validator_off_host",
			"Validator / Off-Host Template",
			"Server-owned validator ID review and off-host validation planning.",
			"Validator selection metadata; treat artifacts and host paths as private deployment details.",
			[]string{"validator_id", "feed_type", "artifact_ref", "run intent", "blocked-check reason"},
			[]string{"private artifact path", "hostnames", "operator workstation path", "validator command arguments"},
			[]string{"raw commands", "arbitrary argv", "private files", "unsupported validator ids", "evidence paths"},
			"validator execution remains operator-run or server-owned; no browser execution.",
			"Block checks that require raw commands, unsupported IDs, private paths, or evidence writes.",
			"go run ./cmd/adapter-conformance validator --suite testdata/adapter-conformance",
			"Validator-clean result, compliance, consumer acceptance, production readiness, or public launch.",
			[]string{"docs/tutorials/gtfs-validation-triage.md", "docs/dependencies.md"},
		),
		connectorWorkbenchRedactionTemplateView(
			"monitoring_export",
			"Monitoring / Export Template",
			"Health digests, monitoring summaries, and deployment-owned notification drafts.",
			"Operational health metadata; treat destination details and incident notes as private.",
			[]string{"feed freshness bucket", "validator count", "stale telemetry count", "redacted incident category"},
			[]string{"notification destination", "person name", "email address", "private endpoint", "unredacted incident text"},
			[]string{"webhook tokens", "send-enabled defaults", "SLA evidence", "retained evidence paths"},
			"send_by_default=false and destination_redacted=true in examples.",
			"Block export when destination, token, private endpoint, or unredacted incident detail appears.",
			"go run ./cmd/adapter-conformance monitoring --suite testdata/adapter-conformance",
			"SLA coverage, uptime guarantee, hosted service availability, retained evidence, compliance, or production readiness.",
			[]string{"docs/tutorials/self-hosted-operations-notifications.md", "examples/connectors/monitoring-export/README.md"},
		),
	}
}

func connectorWorkbenchRedactionTemplateView(id string, label string, appliesTo string, classification string, allowed []string, redact []string, blocked []string, noSend string, failClosed string, check string, doesNotProve string, docsLinks []string) connectorWorkbenchRedactionTemplate {
	return connectorWorkbenchRedactionTemplate{
		ID:                 firstNonEmpty(id, "redaction_template"),
		Label:              firstNonEmpty(label, "Redaction template"),
		AppliesTo:          firstNonEmpty(appliesTo, "Connector review."),
		DataClassification: firstNonEmpty(classification, "Private until operator-redacted."),
		AllowedFields:      cleanLaunchpadList(allowed),
		RedactFields:       cleanLaunchpadList(redact),
		BlockedFields:      cleanLaunchpadList(blocked),
		NoSendDefault:      firstNonEmpty(noSend, "send_enabled=false"),
		FailClosedRule:     firstNonEmpty(failClosed, "Fail closed when input safety is unclear."),
		FirstSafeCheck:     firstNonEmpty(check, "make adapter-conformance"),
		DoesNotProve:       firstNonEmpty(doesNotProve, "Compatibility, compliance, production readiness, or consumer acceptance."),
		DocsLinks:          safeDocsLinks(docsLinks),
	}
}

func connectorWorkbenchDryRunCommands() []connectorWorkbenchDryRun {
	return []connectorWorkbenchDryRun{
		connectorWorkbenchDryRunView(
			"csv-replay-dry-run",
			"CSV telemetry replay dry run",
			"go run ./examples/connectors/telemetry-csv-replay",
			"Operator shell outside the browser.",
			connectorWorkbenchCSVFixture,
			"Synthetic CSV rows normalize to dry-run telemetry events or bounded drops.",
			"Fix CSV columns, timestamps, device identity, coordinates, or quality in the committed synthetic fixture.",
			"Accepted telemetry, Vehicle Positions output, real device proof, vendor compatibility, production readiness, compliance, or consumer acceptance.",
			[]string{"examples/connectors/telemetry-csv-replay/README.md", "docs/tutorials/device-avl-integration.md"},
		),
		connectorWorkbenchDryRunView(
			"http-poller-dry-run",
			"GPS API polling dry run",
			"go run ./examples/connectors/telemetry-http-poller",
			"Operator shell outside the browser.",
			connectorWorkbenchHTTPFixture,
			"Synthetic observation rows normalize to dry-run telemetry events or bounded drops.",
			"Fix fixture shape, timestamp parsing, device identity, coordinates, or quality before any live endpoint is configured.",
			"Live API behavior, named vendor support, vendor compatibility, production readiness, compliance, or consumer acceptance.",
			[]string{"examples/connectors/telemetry-http-poller/README.md", "docs/integration-adapter-kit.md"},
		),
		connectorWorkbenchDryRunView(
			"telemetry-conformance",
			"Telemetry conformance cases",
			"go run ./cmd/adapter-conformance telemetry --suite testdata/adapter-conformance",
			"Operator shell outside the browser.",
			"testdata/adapter-conformance/fixtures",
			"Malformed, stale, future, wrong-agency, unknown-device, low-quality, duplicate, and out-of-order cases fail closed.",
			"Review the named synthetic case and adapter boundary expectation.",
			"Real fleet reliability, hardware certification, real vendor support, production AVL quality, compliance, or ETA quality.",
			[]string{"docs/tutorials/external-adapter-conformance.md", "docs/tutorials/device-avl-integration.md"},
		),
		connectorWorkbenchDryRunView(
			"example-tests",
			"Synthetic connector example tests",
			"make test-connector-examples",
			"Operator shell outside the browser.",
			"examples/connectors and committed fixtures",
			"All committed synthetic connector examples compile and pass local dry-run tests.",
			"Fix the example or fixture without adding network sends, private payloads, credentials, or stronger claims.",
			"Real integration proof, production readiness, vendor compatibility, compliance, consumer acceptance, or retained evidence.",
			[]string{"examples/README.md", "docs/integration-adapter-kit.md"},
		),
	}
}

func connectorWorkbenchDryRunView(id string, label string, commandLine string, runsWhere string, inputs string, expected string, failure string, doesNotProve string, docsLinks []string) connectorWorkbenchDryRun {
	return connectorWorkbenchDryRun{
		ID:                firstNonEmpty(id, "connector-dry-run"),
		Label:             firstNonEmpty(label, "Connector dry run"),
		CommandLine:       firstNonEmpty(commandLine, "make test-connector-examples"),
		RunsWhere:         firstNonEmpty(runsWhere, "Operator shell outside the browser."),
		Inputs:            firstNonEmpty(inputs, "Committed synthetic fixtures only."),
		ExpectedResult:    firstNonEmpty(expected, "Synthetic/local rows pass or produce bounded diagnostics."),
		FailureNextAction: firstNonEmpty(failure, "Review the synthetic fixture or adapter boundary."),
		DoesNotProve:      firstNonEmpty(doesNotProve, "Compatibility, compliance, production readiness, or consumer acceptance."),
		DocsLinks:         safeDocsLinks(docsLinks),
	}
}

func connectorWorkbenchWebhookBoundaryView() connectorWorkbenchWebhookBoundary {
	return connectorWorkbenchWebhookBoundary{
		Title:    "Webhook And AVL Transform Boundaries",
		Boundary: "Boundary guidance only. The Workbench is not a webhook receiver, does not accept vendor payloads, does not hold credentials, does not forward observations, and does not call telemetry ingest.",
		Rows: []connectorWorkbenchWebhookRow{
			connectorWorkbenchWebhookRowView(
				"receiver_owner",
				"Receiver is deployment-owned",
				"Any external POST receiver belongs in deployment-owned infrastructure outside the Operations Console.",
				[]string{"committed synthetic fixtures", "server-issued device credential reference", "deployment-owned endpoint documented outside this page"},
				[]string{"credentials in manifests", "payload samples from live systems", "browser-supplied receiver URLs", "portal or vendor contact from this page"},
				"go run ./cmd/adapter-conformance telemetry --suite testdata/adapter-conformance",
				"Reject before ingest when authentication, agency mapping, device identity, timestamp, coordinate, or quality checks are unsafe.",
				"Store only redacted local diagnostics; do not print private payloads, tokens, headers, private URLs, or filesystem paths.",
				"Named vendor support, vendor compatibility, hardware certification, production AVL reliability, compliance, consumer acceptance, or public launch.",
				[]string{"/admin/operations/devices", "/admin/operations/telemetry"},
			),
			connectorWorkbenchWebhookRowView(
				"transform_mapping",
				"Transform before telemetry ingest",
				"Adapters should convert source-specific observations into the narrow telemetry event shape before any intentional send.",
				[]string{"agency ID", "device ID", "vehicle ID", "observed timestamp", "latitude/longitude", "quality signal"},
				[]string{"unbounded vendor fields", "operator-provided argv", "client-supplied file paths", "direct browser telemetry sends"},
				"make test-connector-examples",
				"Drop or withhold observations with missing identity, invalid timestamps, stale/future timestamps, low quality, or invalid coordinates.",
				"Log reason codes and counts, not source payload bodies or credential material.",
				"Accepted telemetry, Vehicle Positions output, Trip Updates quality, production readiness, or evidence.",
				[]string{"/admin/operations/connectors/workbench", "/admin/operations/realtime"},
			),
			connectorWorkbenchWebhookRowView(
				"credential_boundary",
				"Credentials stay server-owned",
				"Device and adapter credentials must be provisioned through existing private admin paths or deployment configuration, not through manifests or this page.",
				[]string{"environment references", "device token issuance through private Devices & Tokens", "redacted credential presence indicators"},
				[]string{"API keys in JSON", "bearer values in HTML", "credential upload fields", "token hashes in operator-visible tables"},
				"make external-connection-check",
				"Block connector manifests or examples that contain secret-like values, private endpoints, raw commands, or status mutation.",
				"Show only credential ownership and next action; never render secret values.",
				"Credential validity, live external connectivity, vendor approval, production readiness, or SLA.",
				[]string{"/admin/operations/devices", "/admin/operations/maintenance"},
			),
			connectorWorkbenchWebhookRowView(
				"review_before_send",
				"Review before any intentional send",
				"After local synthetic checks pass, a technical helper should review the adapter boundary before any deployment-owned send path is enabled.",
				[]string{"passing synthetic conformance results", "redaction plan", "rollback plan", "operator-owned go/no-go notes outside protected evidence paths"},
				[]string{"retained evidence without authorization", "consumer status movement", "automatic submissions", "public claims based on local checks"},
				"make adapter-conformance",
				"Keep the adapter disabled or shadow-only until authentication, redaction, fail-closed behavior, and rollback are reviewed.",
				"Keep support bundles redacted and avoid retaining real external payloads without separate written authorization.",
				"Agency approval, external proof, compliance, consumer acceptance, public launch, or production readiness.",
				[]string{"/admin/operations/maintenance", "/admin/operations/validation-center"},
			),
		},
		DocsLinks: safeDocsLinks([]string{"docs/tutorials/device-avl-integration.md", "docs/evidence/redaction-policy.md", "docs/integration-adapter-kit.md"}),
	}
}

func connectorWorkbenchWebhookRowView(id string, label string, means string, allowed []string, blocked []string, check string, failClosed string, redaction string, doesNotProve string, reviewLinks []string) connectorWorkbenchWebhookRow {
	return connectorWorkbenchWebhookRow{
		ID:             firstNonEmpty(id, "webhook-boundary"),
		Label:          firstNonEmpty(label, "Webhook boundary"),
		WhatThisMeans:  firstNonEmpty(means, "Review this boundary before any deployment-owned integration work."),
		AllowedInputs:  cleanLaunchpadList(allowed),
		BlockedInputs:  cleanLaunchpadList(blocked),
		FirstSafeCheck: firstNonEmpty(check, "make adapter-conformance"),
		FailClosedRule: firstNonEmpty(failClosed, "Fail closed before ingest when the input is unsafe."),
		RedactionRule:  firstNonEmpty(redaction, "Keep diagnostics redacted."),
		DoesNotProve:   firstNonEmpty(doesNotProve, "Compatibility, compliance, production readiness, or consumer acceptance."),
		ReviewLinks:    safeAdminLinks(reviewLinks),
	}
}

func connectorWorkbenchPredictionGuideView() connectorWorkbenchGuide {
	return connectorWorkbenchGuide{
		Title:    "Prediction Sidecar Guide",
		Boundary: "Prediction connector review only. Vehicle Positions stay independent, Trip Updates can be withheld, and shadow/fail-closed modes do not prove ETA quality or real-world accuracy.",
		Rows: []connectorWorkbenchGuideRow{
			connectorWorkbenchGuideRowView(
				"deterministic_fallback",
				"Deterministic fallback",
				"available",
				"Keep the internal deterministic path as the safe baseline when an external predictor is absent or disabled.",
				[]string{"active GTFS feed version", "current telemetry", "current vehicle assignments"},
				[]string{"bounded Trip Updates diagnostics", "valid empty or withheld Trip Updates when confidence is insufficient"},
				"Prefer unknown/withheld output over false certainty and keep Vehicle Positions publishing independent.",
				"go test ./internal/prediction",
				"Production-grade ETA quality, real-world accuracy, consumer acceptance, production readiness, or external predictor support.",
				[]string{"/admin/operations/realtime", "/admin/operations/feed-health"},
				[]string{"docs/requirements-trip-updates.md", "docs/integration-adapter-kit.md"},
			),
			connectorWorkbenchGuideRowView(
				"external_http_shadow",
				"External HTTP shadow review",
				"review_only",
				"Run a sidecar-shaped adapter against synthetic input while public Trip Updates remain unchanged.",
				[]string{"sanitized prediction request", "Vehicle Positions reference", "synthetic assignments"},
				[]string{"sidecar diagnostics", "withheld public mutation flag", "comparison notes"},
				"Ignore sidecar output for public feeds until sanitized request/response and validator checks are reviewed.",
				"go run ./examples/connectors/predictor-sidecar-stub",
				"Named predictor compatibility, live service behavior, production readiness, consumer acceptance, or ETA quality.",
				[]string{"/admin/operations/connectors/workbench", "/admin/operations/realtime"},
				[]string{"examples/connectors/predictor-sidecar-stub/README.md", "docs/tutorials/external-adapter-conformance.md"},
			),
			connectorWorkbenchGuideRowView(
				"external_http_fail_closed",
				"External HTTP fail-closed adapter",
				"disabled_by_default",
				"Use fail-closed conformance cases before any deployment-owned external predictor is enabled.",
				[]string{"timeout fixture", "malformed fixture", "stale fixture", "wrong-agency fixture", "low-confidence fixture"},
				[]string{"withheld output", "bounded diagnostics", "no Vehicle Positions dependency"},
				"Timeout, malformed, stale, wrong-agency, and low-confidence responses must withhold output or fall back safely.",
				"go run ./cmd/adapter-conformance prediction --suite testdata/adapter-conformance",
				"Production-grade ETA quality, real-world ETA accuracy, named predictor support, or release readiness.",
				[]string{"/admin/operations/realtime", "/admin/operations/connectors/tests"},
				[]string{"docs/requirements-trip-updates.md", "docs/tutorials/external-adapter-conformance.md"},
			),
		},
		DocsLinks: safeDocsLinks([]string{"docs/requirements-trip-updates.md", "docs/tutorials/external-adapter-conformance.md", "docs/integration-adapter-kit.md"}),
	}
}

func connectorWorkbenchMonitoringGuideView() connectorWorkbenchGuide {
	return connectorWorkbenchGuide{
		Title:    "Monitoring Export Guide",
		Boundary: "Monitoring connector review only. Examples are redacted and no-send by default; delivery requires deployment-owned configuration outside this page.",
		Rows: []connectorWorkbenchGuideRow{
			connectorWorkbenchGuideRowView(
				"no_send_export_batch",
				"No-send export batch",
				"covered",
				"Review a redacted dry-run monitoring batch without notification delivery or network send.",
				[]string{"synthetic metric counts", "synthetic incident summary", "redaction policy"},
				[]string{"redacted dry-run export", "send_enabled=false", "network_send=false"},
				"Do not send notifications or write evidence from example monitoring code.",
				"go run ./examples/connectors/monitoring-export",
				"SLA coverage, uptime guarantee, hosted service availability, production readiness, retained evidence, or consumer acceptance.",
				[]string{"/admin/operations/maintenance", "/admin/operations/reliability"},
				[]string{"examples/connectors/monitoring-export/README.md", "docs/tutorials/self-hosted-operations-notifications.md"},
			),
			connectorWorkbenchGuideRowView(
				"redaction_review",
				"Redaction review",
				"covered",
				"Check that monitoring/export fixtures keep private contacts, source bodies, and credential material out of operator-visible output.",
				[]string{"synthetic monitoring fixture", "redaction rules", "no-send defaults"},
				[]string{"redacted summary rows", "bounded diagnostic fields", "no status mutation"},
				"Block outputs that expose secrets, private paths, source bodies, notification destinations, or status mutation.",
				"go run ./cmd/adapter-conformance monitoring --suite testdata/adapter-conformance",
				"Retained evidence, production monitoring maturity, SLA/uptime proof, hosted service availability, or compliance.",
				[]string{"/admin/operations/maintenance", "/admin/operations/connectors/tests"},
				[]string{"docs/evidence/redaction-policy.md", "docs/tutorials/external-adapter-conformance.md"},
			),
			connectorWorkbenchGuideRowView(
				"deployment_delivery",
				"Deployment-owned delivery",
				"needs_review",
				"Delivery to a monitoring system is a deployment decision outside this browser page.",
				[]string{"deployment-owned destination config", "operator authorization", "redaction plan", "rollback plan"},
				[]string{"private monitoring setup notes", "local support bundle pointers", "delivery disabled by default"},
				"Keep delivery disabled until the destination, auth, redaction, and rollback are reviewed by the operator.",
				"make external-connection-check",
				"Hosted service availability, paid support, SLA coverage, uptime guarantee, production readiness, or agency approval.",
				[]string{"/admin/operations/maintenance", "/admin/operations/validation-center"},
				[]string{"docs/tutorials/self-hosted-operations-notifications.md", "docs/runbooks/small-agency-pilot-operations.md"},
			),
		},
		DocsLinks: safeDocsLinks([]string{"docs/tutorials/self-hosted-operations-notifications.md", "docs/evidence/redaction-policy.md", "docs/integration-adapter-kit.md"}),
	}
}

func connectorWorkbenchGuideRowView(id string, label string, status string, what string, inputs []string, outputs []string, failure string, check string, doesNotProve string, reviewLinks []string, docsLinks []string) connectorWorkbenchGuideRow {
	return connectorWorkbenchGuideRow{
		ID:              firstNonEmpty(id, "connector-guide-row"),
		Label:           firstNonEmpty(label, "Connector guide row"),
		Status:          firstNonEmpty(status, checklistStatusNeedsReview),
		WhatThisIs:      firstNonEmpty(what, "Private connector guidance."),
		Inputs:          cleanLaunchpadList(inputs),
		Outputs:         cleanLaunchpadList(outputs),
		FailureBehavior: firstNonEmpty(failure, "Fail closed and keep diagnostics bounded."),
		FirstSafeCheck:  firstNonEmpty(check, "make adapter-conformance"),
		DoesNotProve:    firstNonEmpty(doesNotProve, "Compatibility, compliance, production readiness, or consumer acceptance."),
		ReviewLinks:     safeAdminLinks(reviewLinks),
		DocsLinks:       safeDocsLinks(docsLinks),
	}
}

type connectorWorkbenchSuite struct {
	SchemaVersion      string                        `json:"schema_version"`
	SyntheticOnly      bool                          `json:"synthetic_only"`
	ConnectorManifests []string                      `json:"connector_manifests"`
	Cases              []connectorWorkbenchSuiteCase `json:"cases"`
}

type connectorWorkbenchSuiteCase struct {
	ID              string   `json:"id"`
	Type            string   `json:"type"`
	Scenario        string   `json:"scenario"`
	Fixture         string   `json:"fixture"`
	ExpectedOutcome string   `json:"expected_outcome"`
	Assertions      []string `json:"assertions"`
	SyntheticOnly   bool     `json:"synthetic_only"`
}

func buildConnectorWorkbenchConformanceView() connectorWorkbenchConformanceView {
	view := connectorWorkbenchConformanceView{
		Boundary:       "Synthetic conformance viewer only. These rows describe committed offline cases and fixed operator-shell commands; the browser does not execute the suite, contact external systems, write evidence, or change statuses.",
		SuitePath:      connectorWorkbenchSuitePath,
		Status:         "covered",
		RunnerCommands: connectorWorkbenchConformanceCommands(),
	}
	suite, err := readConnectorWorkbenchSuite()
	if err != nil {
		view.Status = checklistStatusBlocked
		return view
	}
	view.SyntheticOnly = suite.SyntheticOnly
	view.ManifestCount = len(suite.ConnectorManifests)
	view.CaseCount = len(suite.Cases)
	view.Groups = connectorWorkbenchCaseGroups(suite.Cases)
	return view
}

func connectorWorkbenchConformanceCommands() []connectorWorkbenchDryRun {
	return []connectorWorkbenchDryRun{
		connectorWorkbenchDryRunView(
			"adapter-conformance-full",
			"Full synthetic conformance suite",
			"make adapter-conformance",
			"Operator shell outside the browser.",
			"testdata/adapter-conformance/suite.json and committed fixtures",
			"Telemetry, prediction, validator, and monitoring cases are present, synthetic, offline, and fail closed where required.",
			"Open the named case in the committed suite and fix the fixture or adapter boundary expectation.",
			"External network behavior, real validator execution, vendor compatibility, production readiness, compliance, consumer acceptance, or ETA quality.",
			[]string{"docs/tutorials/external-adapter-conformance.md", "docs/connectors/plugin-contract.md"},
		),
		connectorWorkbenchDryRunView(
			"adapter-conformance-telemetry",
			"Telemetry cases",
			"go run ./cmd/adapter-conformance telemetry --suite testdata/adapter-conformance",
			"Operator shell outside the browser.",
			"telemetry cases under testdata/adapter-conformance/fixtures",
			"Malformed, stale, future, wrong-agency, unknown-device, low-quality, duplicate, and out-of-order inputs reject offline.",
			"Review the failing synthetic telemetry case and keep any fix fail-closed.",
			"Production AVL reliability, real device proof, vendor compatibility, hardware certification, compliance, or consumer acceptance.",
			[]string{"docs/tutorials/external-adapter-conformance.md", "docs/tutorials/device-avl-integration.md"},
		),
		connectorWorkbenchDryRunView(
			"adapter-conformance-prediction",
			"Prediction cases",
			"go run ./cmd/adapter-conformance prediction --suite testdata/adapter-conformance",
			"Operator shell outside the browser.",
			"prediction cases under testdata/adapter-conformance/fixtures",
			"Timeout, malformed, stale, wrong-agency, and low-confidence outputs withhold or fail closed.",
			"Review sanitized prediction input/output handling while keeping Vehicle Positions independent.",
			"Production-grade ETA quality, real-world accuracy, named predictor support, production readiness, or release readiness.",
			[]string{"docs/requirements-trip-updates.md", "docs/tutorials/external-adapter-conformance.md"},
		),
		connectorWorkbenchDryRunView(
			"adapter-conformance-validator-monitoring",
			"Validator and monitoring cases",
			"go run ./cmd/adapter-conformance validator --suite testdata/adapter-conformance && go run ./cmd/adapter-conformance monitoring --suite testdata/adapter-conformance",
			"Operator shell outside the browser.",
			"validator and monitoring cases under testdata/adapter-conformance/fixtures",
			"Validator IDs stay allowlisted and monitoring/export stays redacted and no-send.",
			"Fix server-owned validator IDs, redaction expectations, or no-send defaults in committed synthetic fixtures.",
			"Validator-clean feeds, compliance, SLA/uptime proof, hosted service availability, production readiness, or retained evidence.",
			[]string{"docs/tutorials/gtfs-validation-triage.md", "docs/tutorials/self-hosted-operations-notifications.md"},
		),
	}
}

func readConnectorWorkbenchSuite() (connectorWorkbenchSuite, error) {
	raw, err := os.ReadFile(connectorWorkbenchFixtureAbs(connectorWorkbenchSuitePath))
	if err != nil {
		return connectorWorkbenchSuite{}, err
	}
	var suite connectorWorkbenchSuite
	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&suite); err != nil {
		return connectorWorkbenchSuite{}, err
	}
	if suite.SchemaVersion != "open-transit-rt.adapter_conformance.v1" || !suite.SyntheticOnly || len(suite.Cases) == 0 {
		return connectorWorkbenchSuite{}, errors.New("invalid synthetic conformance suite")
	}
	return suite, nil
}

func connectorWorkbenchCaseGroups(cases []connectorWorkbenchSuiteCase) []connectorWorkbenchCaseGroup {
	order := []struct {
		id       string
		label    string
		required []string
		command  string
		boundary string
	}{
		{
			id:       "telemetry",
			label:    "Telemetry connector cases",
			required: []string{"malformed", "stale", "future", "wrong-agency", "unknown-device", "low-quality", "duplicate", "out-of-order", "missing-required-field", "invalid-coordinate"},
			command:  "go run ./cmd/adapter-conformance telemetry --suite testdata/adapter-conformance",
			boundary: "Real device proof, production AVL reliability, vendor compatibility, hardware certification, compliance, or consumer acceptance.",
		},
		{
			id:       "prediction",
			label:    "Prediction connector cases",
			required: []string{"timeout", "malformed", "stale", "wrong-agency", "low-confidence", "missing-vehicle-positions-ref", "public-mutation-attempt"},
			command:  "go run ./cmd/adapter-conformance prediction --suite testdata/adapter-conformance",
			boundary: "Production-grade ETA quality, real-world accuracy, named predictor support, production readiness, or consumer acceptance.",
		},
		{
			id:       "validator",
			label:    "Validator connector cases",
			required: []string{"allowlist", "raw-command"},
			command:  "go run ./cmd/adapter-conformance validator --suite testdata/adapter-conformance",
			boundary: "Validator-clean feeds, CAL-ITP/Caltrans compliance, consumer acceptance, public launch, or production readiness.",
		},
		{
			id:       "monitoring",
			label:    "Monitoring/export connector cases",
			required: []string{"redaction", "no-send", "unredacted-destination"},
			command:  "go run ./cmd/adapter-conformance monitoring --suite testdata/adapter-conformance",
			boundary: "SLA coverage, uptime guarantee, hosted service availability, production readiness, or retained evidence.",
		},
	}
	byType := map[string][]connectorWorkbenchSuiteCase{}
	for _, tc := range cases {
		byType[tc.Type] = append(byType[tc.Type], tc)
	}
	groups := make([]connectorWorkbenchCaseGroup, 0, len(order))
	for _, item := range order {
		groupCases := byType[item.id]
		group := connectorWorkbenchCaseGroup{
			ID:                item.id,
			Label:             item.label,
			Status:            conformanceGroupStatus(groupCases, item.required),
			CaseCount:         len(groupCases),
			RequiredScenarios: cleanLaunchpadList(item.required),
			CommandLine:       item.command,
			DoesNotProve:      item.boundary,
		}
		for _, tc := range groupCases {
			group.Cases = append(group.Cases, connectorWorkbenchCaseSummary{
				ID:              firstNonEmpty(tc.ID, item.id+"-case"),
				Scenario:        firstNonEmpty(tc.Scenario, "scenario"),
				FixturePath:     safeConformanceFixturePath(tc.Fixture),
				ExpectedOutcome: firstNonEmpty(tc.ExpectedOutcome, "review"),
				Assertions:      cleanLaunchpadList(tc.Assertions),
				Status:          conformanceCaseStatus(tc),
				SyntheticOnly:   tc.SyntheticOnly,
			})
		}
		groups = append(groups, group)
	}
	return groups
}

func conformanceGroupStatus(cases []connectorWorkbenchSuiteCase, required []string) string {
	if len(cases) == 0 {
		return checklistStatusMissing
	}
	seen := map[string]bool{}
	for _, tc := range cases {
		if tc.SyntheticOnly {
			seen[tc.Scenario] = true
		}
	}
	for _, scenario := range required {
		if !seen[scenario] {
			return checklistStatusNeedsReview
		}
	}
	return "covered"
}

func conformanceCaseStatus(tc connectorWorkbenchSuiteCase) string {
	if tc.SyntheticOnly && tc.ID != "" && tc.Fixture != "" && tc.ExpectedOutcome != "" && len(tc.Assertions) > 0 {
		return "covered"
	}
	return checklistStatusNeedsReview
}

func safeConformanceFixturePath(fixture string) string {
	clean := filepath.ToSlash(filepath.Clean(fixture))
	if clean != fixture || strings.HasPrefix(clean, "../") || strings.HasPrefix(clean, "/") {
		return "invalid-fixture-path"
	}
	if !strings.HasPrefix(clean, "fixtures/") {
		return "invalid-fixture-path"
	}
	return "testdata/adapter-conformance/" + clean
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

func connectorWorkbenchManifestLintChecks() []connectorWorkbenchManifestLint {
	return []connectorWorkbenchManifestLint{
		connectorWorkbenchManifestLintView(
			"secret_and_endpoint_scan",
			"Secrets and private endpoints",
			"covered",
			"internal/connectors.DecodeManifest and make external-connection-check",
			"Secret-like values, bearer text, private endpoint strings, private local paths, and unsafe URLs in displayable manifest fields.",
			"Remove real credentials, live endpoints, private paths, and payload fragments before committing a manifest.",
			"Credential validity, live connectivity, vendor approval, production readiness, or SLA.",
		),
		connectorWorkbenchManifestLintView(
			"command_and_plugin_boundary",
			"Command and plugin boundary",
			"covered",
			"manifest unknown-field rejection, raw command rejection, and safe plugin definition review",
			"Manifest-provided commands, arbitrary argv, dynamic backend plugin loading, validator command strings, and browser-run sidecars.",
			"Use fixed operator-shell commands from Connector Tests and server-owned validator IDs instead of manifest commands.",
			"Runtime connector execution, validator-clean proof, external contact, compliance, or compatibility.",
		),
		connectorWorkbenchManifestLintView(
			"status_submission_no_send",
			"No status mutation, submission automation, or send-by-default",
			"covered",
			"internal/connectors.Manifest.Validate",
			"Consumer submission automation, consumer tracker mutation, internal status mutation, and notification sending by default.",
			"Keep connectors disabled by default, no-send by default, and operator-reviewed before any deployment-owned delivery path.",
			"Consumer acceptance, notification delivery, hosted service availability, public launch, or evidence creation.",
		),
		connectorWorkbenchManifestLintView(
			"claim_boundary",
			"Positive claim allowlist",
			"covered",
			"internal/connectors.Manifest.Validate",
			"Unsupported positive claims about compliance, agency approval, public launch, vendor compatibility, hardware certification, production readiness, hosted service, SLA, or ETA quality.",
			"Use only bounded claims such as adapter contract only, synthetic conformance tested, disabled by default, redacted diagnostics only, and no status mutation.",
			"Any real-world proof, external acceptance, or release readiness.",
		),
		connectorWorkbenchManifestLintView(
			"fixture_scope",
			"Synthetic fixture scope",
			"covered",
			"internal/connectors.Manifest.Validate and cmd/adapter-conformance",
			"Absolute paths, parent traversal, evidence paths, non-synthetic fixtures, and unbounded fixture references.",
			"Reference committed synthetic fixtures under allowed testdata or examples paths and keep cases offline.",
			"Production data quality, real device proof, retained evidence, compliance, or consumer acceptance.",
		),
	}
}

func connectorWorkbenchManifestLintView(id string, label string, status string, enforcedBy string, blocks string, action string, doesNotProve string) connectorWorkbenchManifestLint {
	return connectorWorkbenchManifestLint{
		ID:             firstNonEmpty(id, "manifest_lint"),
		Label:          firstNonEmpty(label, "Manifest lint"),
		Status:         firstNonEmpty(status, "covered"),
		EnforcedBy:     firstNonEmpty(enforcedBy, "make external-connection-check"),
		Blocks:         firstNonEmpty(blocks, "Unsafe connector manifest content."),
		OperatorAction: firstNonEmpty(action, "Fix the manifest before review."),
		DoesNotProve:   firstNonEmpty(doesNotProve, "Compatibility, compliance, production readiness, or consumer acceptance."),
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

type connectorWorkbenchRawObservation struct {
	AgencyID   string
	DeviceID   string
	VehicleID  string
	ObservedAt time.Time
	Latitude   float64
	Longitude  float64
	Quality    float64
}

type connectorWorkbenchHTTPFixtureData struct {
	SyntheticOnly bool                                `json:"synthetic_only"`
	Observations  []connectorWorkbenchHTTPObservation `json:"observations"`
}

type connectorWorkbenchHTTPObservation struct {
	AgencyID  string  `json:"agency_id"`
	DeviceID  string  `json:"device_id"`
	VehicleID string  `json:"vehicle_id"`
	Timestamp string  `json:"timestamp"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Quality   float64 `json:"quality"`
}

func buildConnectorWorkbenchTelemetryPreview() connectorWorkbenchTelemetryPreview {
	preview := connectorWorkbenchTelemetryPreview{
		Boundary: "Synthetic fixture preview only. Rows are normalized in-memory for review, not accepted by telemetry ingest, not emitted as Vehicle Positions, not sent over the network, not evidence, and not connector proof.",
	}
	preview.addSource("csv_replay", "CSV replay fixture", connectorWorkbenchCSVFixture, "go run ./examples/connectors/telemetry-csv-replay", time.Date(2026, 5, 10, 16, 1, 0, 0, time.UTC), readConnectorWorkbenchCSVObservations)
	preview.addSource("http_poller", "GPS API fixture", connectorWorkbenchHTTPFixture, "go run ./examples/connectors/telemetry-http-poller", time.Date(2026, 5, 10, 15, 1, 0, 0, time.UTC), readConnectorWorkbenchHTTPObservations)
	preview.Counts.Sources = len(preview.Sources)
	preview.Counts.Rows = len(preview.Rows)
	return preview
}

func (p *connectorWorkbenchTelemetryPreview) addSource(id string, label string, fixturePath string, commandLine string, now time.Time, read func(string) ([]connectorWorkbenchRawObservation, bool, error)) {
	observations, syntheticOnly, err := read(fixturePath)
	source := connectorWorkbenchPreviewSource{
		ID:            id,
		Label:         label,
		FixturePath:   fixturePath,
		Status:        "covered",
		SyntheticOnly: syntheticOnly,
		CommandLine:   commandLine,
		DoesNotProve:  "Accepted telemetry, Vehicle Positions output, real device proof, vendor compatibility, production readiness, compliance, or consumer acceptance.",
	}
	if err != nil {
		source.Status = checklistStatusBlocked
		source.DoesNotProve = "A missing or unreadable committed fixture does not prove connector failure or success."
		p.Sources = append(p.Sources, source)
		return
	}
	source.ObservedRows = len(observations)
	for _, observation := range observations {
		row := normalizeConnectorWorkbenchPreviewRow(id, observation, now)
		if row.Outcome == "event" {
			source.ExpectedEvents++
			p.Counts.Events++
		} else {
			source.ExpectedDrops++
			p.Counts.Drops++
		}
		p.Rows = append(p.Rows, row)
	}
	p.Sources = append(p.Sources, source)
}

func readConnectorWorkbenchCSVObservations(fixturePath string) ([]connectorWorkbenchRawObservation, bool, error) {
	raw, err := os.ReadFile(connectorWorkbenchFixtureAbs(fixturePath))
	if err != nil {
		return nil, false, err
	}
	reader := csv.NewReader(strings.NewReader(string(raw)))
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, false, err
	}
	if len(rows) < 2 || len(rows[0]) != 8 {
		return nil, false, errors.New("invalid synthetic csv fixture shape")
	}
	var observations []connectorWorkbenchRawObservation
	for _, row := range rows[1:] {
		if len(row) != 8 || row[0] != "true" {
			return nil, false, errors.New("csv rows must be synthetic and eight columns")
		}
		observedAt, _ := time.Parse(time.RFC3339, row[4])
		latitude, _ := strconv.ParseFloat(row[5], 64)
		longitude, _ := strconv.ParseFloat(row[6], 64)
		quality, _ := strconv.ParseFloat(row[7], 64)
		observations = append(observations, connectorWorkbenchRawObservation{
			AgencyID:   row[1],
			DeviceID:   row[2],
			VehicleID:  row[3],
			ObservedAt: observedAt,
			Latitude:   latitude,
			Longitude:  longitude,
			Quality:    quality,
		})
	}
	return observations, true, nil
}

func readConnectorWorkbenchHTTPObservations(fixturePath string) ([]connectorWorkbenchRawObservation, bool, error) {
	raw, err := os.ReadFile(connectorWorkbenchFixtureAbs(fixturePath))
	if err != nil {
		return nil, false, err
	}
	var fixture connectorWorkbenchHTTPFixtureData
	if err := json.Unmarshal(raw, &fixture); err != nil {
		return nil, false, err
	}
	if !fixture.SyntheticOnly || len(fixture.Observations) == 0 {
		return nil, fixture.SyntheticOnly, errors.New("http fixture must be synthetic and non-empty")
	}
	observations := make([]connectorWorkbenchRawObservation, 0, len(fixture.Observations))
	for _, observation := range fixture.Observations {
		observedAt, _ := time.Parse(time.RFC3339, observation.Timestamp)
		observations = append(observations, connectorWorkbenchRawObservation{
			AgencyID:   observation.AgencyID,
			DeviceID:   observation.DeviceID,
			VehicleID:  observation.VehicleID,
			ObservedAt: observedAt,
			Latitude:   observation.Latitude,
			Longitude:  observation.Longitude,
			Quality:    observation.Quality,
		})
	}
	return observations, true, nil
}

func normalizeConnectorWorkbenchPreviewRow(sourceID string, observation connectorWorkbenchRawObservation, now time.Time) connectorWorkbenchPreviewRow {
	row := connectorWorkbenchPreviewRow{
		SourceID:    sourceID,
		DeviceID:    firstNonEmpty(observation.DeviceID, "missing-device"),
		VehicleID:   firstNonEmpty(observation.VehicleID, "missing-vehicle"),
		ObservedAt:  observation.ObservedAt.UTC().Format(time.RFC3339),
		Quality:     strconv.FormatFloat(observation.Quality, 'f', 2, 64),
		Outcome:     "event",
		DryRun:      true,
		NetworkSend: false,
	}
	switch {
	case observation.AgencyID == "" || observation.DeviceID == "" || observation.VehicleID == "":
		row.Outcome = "drop"
		row.Reason = "missing identity"
	case observation.ObservedAt.IsZero():
		row.Outcome = "drop"
		row.Reason = "invalid timestamp"
	case observation.ObservedAt.After(now.Add(30 * time.Second)):
		row.Outcome = "drop"
		row.Reason = "future timestamp"
	case now.Sub(observation.ObservedAt) > 2*time.Minute:
		row.Outcome = "drop"
		row.Reason = "stale observation"
	case observation.Quality < 0.5:
		row.Outcome = "drop"
		row.Reason = "low quality"
	case observation.Latitude < -90 || observation.Latitude > 90 || observation.Longitude < -180 || observation.Longitude > 180:
		row.Outcome = "drop"
		row.Reason = "invalid coordinates"
	}
	return row
}

func connectorWorkbenchFixtureAbs(rel string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return rel
	}
	return filepath.Join(filepath.Dir(filename), "..", "..", filepath.FromSlash(rel))
}
